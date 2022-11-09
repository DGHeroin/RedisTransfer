package merge

import (
    "context"
    "github.com/go-redis/redis/v8"
    "sync/atomic"
    "time"
)

func HandleZSet(key string) error {
    t0 := time.Now()
    defer func() {
        elapsedTime := time.Since(t0)
        logd("[zset] 成功 %v %s\n", elapsedTime, key)
    }()

    result, err := sourceClient.ZRangeWithScores(context.Background(), key, 0, -1).Result()
    if err != nil {
        return err
    }
    var args []*redis.Z
    for _, s := range result {
        args = append(args, &s)
    }
    err = targetClient.ZAdd(context.Background(), key, args...).Err()
    if err != nil {
        return err
    }
    atomic.AddUint32(&countZSet, 1)
    return nil
}
