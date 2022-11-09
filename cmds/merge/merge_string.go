package merge

import (
    "context"
    "sync/atomic"
    "time"
)

func HandleString(key string) error {
    t0 := time.Now()
    defer func() {
        elapsedTime := time.Since(t0)
        logd("[string] 成功 %v %s\n", elapsedTime, key)
    }()

    str, err := sourceClient.Get(context.Background(), key).Result()
    if err != nil {
        return err
    }

    err = targetClient.Set(context.Background(), key, str, 0).Err()
    if err != nil {
        return err
    }

    atomic.AddUint32(&countString, 1)
    return nil
}
