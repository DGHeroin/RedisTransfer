package load

import (
    "RedisTransfer/log"
    "RedisTransfer/types"
    "context"
    "github.com/spf13/cobra"
    "io"
    "sync"
    "sync/atomic"
    "time"
)

var (
    Cmd = &cobra.Command{
        Use:   "load command <args>",
        Short: "load",
        RunE: func(cmd *cobra.Command, args []string) error {
            var key string
            if len(args) >= 1 {
                key = args[0]
            }
            return handleLoad(key, "")
        },
        PreRunE: func(cmd *cobra.Command, args []string) error {
            if err := connectToRedis(); err != nil {
                return err
            }
            return nil
        },
    }
)

func handleLoad(match, keyType string) (err error) {
    defer inputFile.Close()

    var wg sync.WaitGroup
    ch := make(chan *types.DataBase, 500)
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

    for {
        var info = &types.DataBase{}
        err = inputStream.Decode(info)
        if err != nil {
            if err == io.EOF {
                err = nil
            }
            break
        }
        wg.Add(1)
        ch <- info
    }
    wg.Wait()
    log.I("[处理完成] key数量:%v string:%v hash:%v set:%v zset:%v list:%v ttl keys:%v keys:%v\n",
        atomic.LoadUint32(&count), atomic.LoadUint32(&countString), atomic.LoadUint32(&countHash),
        atomic.LoadUint32(&countSet), atomic.LoadUint32(&countZSet), atomic.LoadUint32(&countList),
        atomic.LoadUint32(&countTTL), len(mm),
    )
    return
}
func keyHandler(ch chan *types.DataBase, wg *sync.WaitGroup) {
    for info := range ch {
        handleKey(info, wg)
    }
}

func handleKey(info *types.DataBase, wg *sync.WaitGroup) {
    defer func() {
        atomic.AddUint32(&count, 1)
        wg.Done()
    }()
    mu.Lock()
    mm[info.Key] = true
    mu.Unlock()
    var err error

    if preDelete {
        targetClient.Del(context.Background(), info.Key)
    }
    switch info.Type {
    case "string":
        err = HandleString(info)
    case "hash":
        err = HandleHash(info)
    case "set":
        err = HandleSet(info)
    case "zset":
        err = HandleZSet(info)
    case "list":
        err = HandleList(info)
    }
    if err != nil {
        log.E("fail:%s %v\n", info, err)
        return
    }
}
