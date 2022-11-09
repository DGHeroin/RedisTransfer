package merge

import (
    "context"
    "log"
    "sync/atomic"
    "time"
)

func HandleHash(key string) error {
    t0 := time.Now()
    defer func() {
        elapsedTime := time.Since(t0)
        logd("[hash] 成功 %v %s\n", elapsedTime, key)
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
        log.Println(err)
        return err
    }

    atomic.AddUint32(&countHash, 1)
    return nil
}
