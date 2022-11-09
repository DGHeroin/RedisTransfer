package merge

import (
    "context"
    "sync/atomic"
    "time"
)

func HandleSet(key string) error {
    t0 := time.Now()
    defer func() {
        elapsedTime := time.Since(t0)
        logd("[set] 成功 %v %s\n", elapsedTime, key)
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
    atomic.AddUint32(&countSet, 1)
    return nil
}
