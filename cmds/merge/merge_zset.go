package merge

import (
    "context"
    "github.com/go-redis/redis/v8"
    "sync/atomic"
    "time"
)

func HandleZSet(key string) error {
    t0 := time.Now()
    var sz int64
    defer func() {
        elapsedTime := time.Since(t0)
        logd("[zset] %d 成功 %v 大小:%v [%s] \n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    result, err := sourceClient.ZRangeWithScores(context.Background(), key, 0, -1).Result()
    if err != nil {
        return err
    }
    var args []*redis.Z
    for _, s := range result {
        args = append(args, &redis.Z{
            Score:  s.Score,
            Member: s.Member,
        })
    }
    err = targetClient.ZAdd(context.Background(), key, args...).Err()
    if err != nil {
        return err
    }
    if val, err := targetClient.ZCard(context.Background(), key).Result(); err == nil {
        sz = val
    }

    atomic.AddUint32(&countZSet, 1)
    return nil
}
