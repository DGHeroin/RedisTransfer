package load

import (
    "RedisTransfer/log"
    "RedisTransfer/types"
    "context"
    "github.com/go-redis/redis/v8"
    "sync/atomic"
    "time"
)

func HandleString(info *types.DataBase) error {
    key := info.Key
    var sz = len(info.String)

    t0 := time.Now()
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[string] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    _, err := targetClient.Set(context.Background(), key, info.String, 0).Result()
    if err == nil {
        checkTTL(info)
    }
    return err
}
func HandleHash(info *types.DataBase) error {
    key := info.Key
    var sz = len(info.Hash)

    t0 := time.Now()
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[hash] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    _, err := targetClient.HSet(context.Background(), key, info.Hash).Result()
    if err == nil {
        checkTTL(info)
    }
    return err
}
func HandleSet(info *types.DataBase) error {
    key := info.Key
    var sz = len(info.Set)

    t0 := time.Now()
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[set] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    _, err := targetClient.SAdd(context.Background(), key, info.Set).Result()
    if err == nil {
        checkTTL(info)
    }
    return err
}
func HandleZSet(info *types.DataBase) error {
    key := info.Key
    var sz = len(info.Hash)

    t0 := time.Now()
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[zset] %d 成功 %v 大小:%v [%s] \n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()
    var ss []*redis.Z
    for _, s := range info.ZSet {
        ss = append(ss, &redis.Z{
            Score:  s.Score,
            Member: s.Member,
        })
    }

    _, err := targetClient.ZAdd(context.Background(), key, ss...).Result()
    if err == nil {
        checkTTL(info)
    }
    return err
}
func HandleList(info *types.DataBase) error {
    key := info.Key
    var sz = len(info.Hash)

    t0 := time.Now()
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[list] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    _, err := targetClient.LPush(context.Background(), key, info.List).Result()
    if err == nil {
        checkTTL(info)
    }
    return err
}
func checkTTL(info *types.DataBase) {
    if info.TTL <= 0 {
        return
    }
    ttl := time.Unix(info.TTL, 0).Sub(time.Now())
    targetClient.Expire(context.Background(), info.Key, ttl)
    atomic.AddUint32(&countTTL, 1)
}
