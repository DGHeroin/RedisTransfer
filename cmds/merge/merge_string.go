package merge

import (
    "context"
    "sync/atomic"
    "time"
)

func HandleString(key string) error {
    t0 := time.Now()
    var sz int64
    defer func() {
        elapsedTime := time.Since(t0)
        logd("[string] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    str, err := sourceClient.Get(context.Background(), key).Result()
    if err != nil {
        return err
    }

    err = targetClient.Set(context.Background(), key, str, 0).Err()
    if err != nil {
        return err
    }
    sz = int64(len(str))

    atomic.AddUint32(&countString, 1)
    return nil
}
