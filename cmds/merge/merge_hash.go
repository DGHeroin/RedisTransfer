package merge

import (
    "RedisTransfer/log"
    "context"
    "sync/atomic"
    "time"
)

func HandleHash(key string) error {
    t0 := time.Now()
    var sz int64
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[hash] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    kv, err := sourceClient.HGetAll(context.Background(), key).Result()
    if err != nil {
        return err
    }

    var args []interface{}
    for k, v := range kv {
        args = append(args, k, v)
    }
    err = targetClient.HSet(context.Background(), key, args...).Err()
    if err != nil {
        return err
    }

    if val, err := targetClient.HLen(context.Background(), key).Result(); err == nil {
        sz = val
    }

    atomic.AddUint32(&countHash, 1)
    return nil
}
