package merge

import (
    "context"
    "sync/atomic"
    "time"
)

func HandleList(key string) error {
    t0 := time.Now()
    var sz int64
    defer func() {
        elapsedTime := time.Since(t0)
        logd("[list] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    result, err := sourceClient.LRange(context.Background(), key, 0, -1).Result()
    if err != nil {
        return err
    }
    var args []interface{}
    for _, s := range result {
        args = append(args, s)
    }

    err = targetClient.LPush(context.Background(), key, args...).Err()
    if err != nil {
        return err
    }

    if val, err := targetClient.LRange(context.Background(), key, 0, -1).Result(); err == nil {
        sz = int64(len(val))
    }

    atomic.AddUint32(&countList, 1)

    return nil
}
