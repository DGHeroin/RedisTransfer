package merge

import (
    "RedisTransfer/log"
    "context"
    "sync/atomic"
    "time"
)

func HandleSet(key string) error {
    t0 := time.Now()
    var sz int64
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[set] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    result, err := sourceClient.SMembers(context.Background(), key).Result()
    if err != nil {
        return err
    }
    var args []interface{}
    for _, s := range result {
        args = append(args, s)
    }
    err = targetClient.SAdd(context.Background(), key, args...).Err()
    if err != nil {
        return err
    }

    if val, err := targetClient.SCard(context.Background(), key).Result(); err == nil {
        sz = val
    }

    atomic.AddUint32(&countSet, 1)
    return nil
}
