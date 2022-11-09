package merge

import (
    "context"
    "sync"
    "sync/atomic"
    "time"
)

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
        loge("fail:%s\n", key)
        return
    }
    switch t {
    case "string":
        err = HandleString(key)
    case "hash":
        err = HandleHash(key)
    case "set":
        err = HandleSet(key)
    case "zset":
        err = HandleZSet(key)
    case "list":
        err = HandleList(key)
    }
    if err != nil {
        loge("fail:%s\n", key)
        return
    }
    checkTTL(key)
}

func checkTTL(key string) {
    ttl, err := sourceClient.TTL(context.Background(), key).Result()
    if err != nil {
        return
    }
    if ttl != -1 {
        err = targetClient.Expire(context.Background(), key, ttl).Err()
        if err != nil {
            return
        }
    }
    atomic.AddUint32(&countTTL, 1)
}
func keyHandler(ch chan string, wg *sync.WaitGroup) {
    for key := range ch {
        handleKey(key, wg)
    }
}

func handleMerge(match, keyType string) error {
    defer func() {
        if err := sourceClient.Close(); err != nil {
            loge("[关闭源 redis] 错误:%v\n", err)
        }
        if err := targetClient.Close(); err != nil {
            loge("[关闭目的 redis] 错误:%v\n", err)
        }
    }()

    var wg sync.WaitGroup
    ch := make(chan string, 500)
    logi("[合并] 工作者:[%v] 类型:[%s] match:[%s]\n", workerNum, keyType, match)
    go func() { // monitor
        for {
            time.Sleep(time.Second)
            logi("[已处理] key数量:%v string:%v hash:%v set:%v zset:%v list:%v ttl keys:%v keys:%v\n",
                atomic.LoadUint32(&count), atomic.LoadUint32(&countString), atomic.LoadUint32(&countHash),
                atomic.LoadUint32(&countSet), atomic.LoadUint32(&countZSet), atomic.LoadUint32(&countList),
                atomic.LoadUint32(&countTTL), len(mm),
            )
        }
    }()

    for i := 0; i < workerNum; i++ {
        go keyHandler(ch, &wg)
    }

    err := getSourceKeys(match, keyType, func(keys []string) {
        for _, key := range keys {
            wg.Add(1)
            ch <- key
        }
    })

    close(ch)

    wg.Wait()
    logi("[处理完成] key数量:%v string:%v hash:%v set:%v zset:%v list:%v ttl keys:%v keys:%v\n",
        atomic.LoadUint32(&count), atomic.LoadUint32(&countString), atomic.LoadUint32(&countHash),
        atomic.LoadUint32(&countSet), atomic.LoadUint32(&countZSet), atomic.LoadUint32(&countList),
        atomic.LoadUint32(&countTTL), len(mm),
    )
    return err
}

// [00:00:06][处理完成] key数量:53347 string:0 hash:0 set:9171 zset:2379 list:3
