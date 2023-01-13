package dump

import (
    "RedisTransfer/log"
    "RedisTransfer/types"
    "context"
    "fmt"
    "sync/atomic"
    "time"
)

func HandleString(key string) error {
    t0 := time.Now()
    var sz int64
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[string] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    str, err := sourceClient.Get(context.Background(), key).Bytes()
    if err != nil {
        return err
    }
    sz = int64(len(str))
    var info = types.DataBase{}
    info.Key = key
    info.Type = "string"
    info.String = str
    info.TTL = checkTTL(key)
    atomic.AddUint32(&countString, 1)
    return writeData(&info)
}
func HandleHash(key string) error {
    t0 := time.Now()
    var sz int64
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[hash] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    str, err := sourceClient.HGetAll(context.Background(), key).Result()
    if err != nil {
        return err
    }
    sz = int64(len(str))
    var info = types.DataBase{}
    info.Key = key
    info.Type = "hash"
    info.Hash = str
    info.TTL = checkTTL(key)
    atomic.AddUint32(&countHash, 1)
    return writeData(&info)
}
func HandleSet(key string) error {
    t0 := time.Now()
    var sz int64
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[set] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    str, err := sourceClient.SMembers(context.Background(), key).Result()
    if err != nil {
        return err
    }
    sz = int64(len(str))
    var info = types.DataBase{}
    info.Key = key
    info.Type = "set"
    info.Set = str
    info.TTL = checkTTL(key)
    atomic.AddUint32(&countSet, 1)
    return writeData(&info)
}
func HandleZSet(key string) error {
    t0 := time.Now()
    var sz int64
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[zset] %d 成功 %v 大小:%v [%s] \n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    z, err := sourceClient.ZRangeWithScores(context.Background(), key, 0, -1).Result()
    if err != nil {
        return err
    }

    var info = types.DataBase{}
    sz = int64(len(z))
    info.Key = key
    info.Type = "zset"
    info.TTL = checkTTL(key)
    for _, v := range z {
        info.ZSet = append(info.ZSet, struct {
            Score  float64 `json:"score"`
            Member string  `json:"member"`
        }{Score: v.Score, Member: fmt.Sprintf("%s", v.Member)})
    }
    atomic.AddUint32(&countZSet, 1)
    return writeData(&info)
}
func HandleList(key string) error {
    t0 := time.Now()
    var sz int64
    defer func() {
        elapsedTime := time.Since(t0)
        log.D("[list] %d 成功 %v 大小:%v [%s]\n", atomic.LoadUint32(&count), elapsedTime, sz, key)
    }()

    str, err := sourceClient.LRange(context.Background(), key, 0, -1).Result()
    if err != nil {
        return err
    }
    sz = int64(len(str))
    var info = types.DataBase{}
    info.Key = key
    info.Type = "list"
    info.List = str
    info.TTL = checkTTL(key)
    atomic.AddUint32(&countList, 1)
    return writeData(&info)
}
func checkTTL(key string) int64 {
    ttl, _ := sourceClient.TTL(context.Background(), key).Result()
    if ttl == -1 || ttl == 0 {
        return 0
    }

    atomic.AddUint32(&countTTL, 1)

    return time.Now().Add(ttl).Unix()
}
func writeData(v interface{}) error {
    mu.Lock()
    defer mu.Unlock()
    return outputStream.Encode(v)
}
