package dump

import (
    "RedisTransfer/log"
    "context"
    "github.com/go-redis/redis/v8"
    "github.com/spf13/cobra"
    "sync"
    "sync/atomic"
    "time"
)

var (
    Cmd = &cobra.Command{
        Use:   "dump command <args>",
        Short: "dump数据",
        RunE: func(cmd *cobra.Command, args []string) error {
            var (
                key string
            )
            if len(args) >= 1 {
                key = args[0]
            }
            return handleDump(key, "")
        },
        PreRunE: func(cmd *cobra.Command, args []string) error {
            if err := connectToRedis(); err != nil {
                return err
            }
            return nil
        },
    }
)

type (
    _Scanable interface {
        Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
        ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *redis.ScanCmd
    }
)

func handleDump(match, keyType string) error {
    defer outputFile.Close()

    var wg sync.WaitGroup
    ch := make(chan string, 500)
    log.I("[合并] 工作者:[%v] 类型:[%s] match:[%s]\n", workerNum, keyType, match)
    if status {
        go func() { // monitor
            for {
                time.Sleep(time.Second)
                log.I("[已处理] key数量:%v string:%v hash:%v set:%v zset:%v list:%v ttl keys:%v keys:%v\n",
                    atomic.LoadUint32(&count), atomic.LoadUint32(&countString), atomic.LoadUint32(&countHash),
                    atomic.LoadUint32(&countSet), atomic.LoadUint32(&countZSet), atomic.LoadUint32(&countList),
                    atomic.LoadUint32(&countTTL), len(mm),
                )
            }
        }()
    }
    for i := 0; i < workerNum; i++ {
        go keyHandler(ch, &wg)
    }

    err := getSourceKeys(match, keyType, func(keys []string) {
        for _, key := range keys {
            wg.Add(1)
            ch <- key
        }
    })
    wg.Wait()
    log.I("[处理完成] key数量:%v string:%v hash:%v set:%v zset:%v list:%v ttl keys:%v keys:%v\n",
        atomic.LoadUint32(&count), atomic.LoadUint32(&countString), atomic.LoadUint32(&countHash),
        atomic.LoadUint32(&countSet), atomic.LoadUint32(&countZSet), atomic.LoadUint32(&countList),
        atomic.LoadUint32(&countTTL), len(mm),
    )
    return err
}
func keyHandler(ch chan string, wg *sync.WaitGroup) {
    for key := range ch {
        handleKey(key, wg)
    }
}
func handleKey(key string, wg *sync.WaitGroup) {
    defer func() {
        atomic.AddUint32(&count, 1)
        wg.Done()
    }()
    mu.Lock()
    mm[key] = true
    mu.Unlock()

    t, err := sourceClient.Type(context.Background(), key).Result()
    if err != nil {
        return
    }
    switch t {
    case "string":
        err = HandleString(key)
        if err == nil {
            atomic.AddUint32(&countString, 1)
        }
    case "hash":
        err = HandleHash(key)
        if err == nil {
            atomic.AddUint32(&countHash, 1)
        }
    case "set":
        err = HandleSet(key)
        if err == nil {
            atomic.AddUint32(&countSet, 1)
        }
    case "zset":
        err = HandleZSet(key)
        if err == nil {
            atomic.AddUint32(&countZSet, 1)
        }
    case "list":
        err = HandleList(key)
        if err == nil {
            atomic.AddUint32(&countList, 1)
        }
    }
    if err != nil {
        log.E("fail:%s %v\n", key, err)
        return
    }
}
func getSourceKeys(match, keyType string, fn func(keys []string)) error {
    if sourceCluster {
        return sourceClient.ForEachShard(context.Background(), func(ctx context.Context, client *redis.Client) error {
            return doScanKeys(client, match, keyType, fn)
        })
    } else {
        return doScanKeys(sourceClient, match, keyType, fn)
    }
}
func doScanKeys(client _Scanable, match, keyType string, fn func([]string)) error {
    var cursor uint64
    if match == "" {
        match = "*"
    }
    for {
        var keys []string
        var err error
        if keyType == "" {
            keys, cursor, err = client.Scan(context.Background(), cursor, match, 1000).Result()
        } else {
            keys, cursor, err = client.ScanType(context.Background(), cursor, match, 1000, keyType).Result()
        }
        if err != nil {
            panic(err)
        }
        fn(keys)
        if cursor == 0 {
            break
        }
    }
    return nil
}
