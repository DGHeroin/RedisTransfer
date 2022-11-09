package merge

import (
    "context"
    "sync/atomic"
    "time"
)

func HandleList(key string) error {
    t0 := time.Now()
    defer func() {
        elapsedTime := time.Since(t0)
        logd("[list] 成功 %v %s\n", elapsedTime, key)
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
    atomic.AddUint32(&countList, 1)

    return nil
}
