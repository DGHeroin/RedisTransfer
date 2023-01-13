package merge

import (
    "RedisTransfer/log"
    "RedisTransfer/redis2"
    "context"
    "github.com/go-redis/redis/v8"
)

type (
    _Scanable interface {
        Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
        ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *redis.ScanCmd
    }
)

func doScanKeys(client _Scanable, match, keyType string, fn func([]string)) error {
    var cursor uint64
    if match == "" {
        match = "*"
    }
    for {
        var keys []string
        var err error
        if keyType == "" {
            keys, cursor, err = client.Scan(context.Background(), cursor, match, 1000).Result()
        } else {
            keys, cursor, err = client.ScanType(context.Background(), cursor, match, 1000, keyType).Result()
        }
        if err != nil {
            panic(err)
        }
        fn(keys)
        if cursor == 0 {
            break
        }
    }
    return nil
}

func getSourceKeys(match, keyType string, fn func(keys []string)) error {
    if sourceCluster {
        return sourceClient.ForEachShard(context.Background(), func(ctx context.Context, client *redis.Client) error {
            return doScanKeys(client, match, keyType, fn)
        })
    } else {
        return doScanKeys(sourceClient, match, keyType, fn)
    }
}
func connectToRedis() (err error) {
    // 源
    {
        redisType := "single"
        if sourceCluster {
            redisType = "cluster"
        }
        sourceClient, err = redis2.NewRedisClient(redisType, sourceAddr, sourceAuth)
        if err != nil {
            return err
        }
        log.I("[初始化源] %v %v\n", sourceAddr, redisType)
    }
    // 目标
    {
        redisType := "single"
        if targetCluster {
            redisType = "cluster"
        }
        targetClient, err = redis2.NewRedisClient(redisType, targetAddr, targetAuth)
        if err != nil {
            return err
        }
        log.I("[初始化目标] %v %v\n", targetAddr, redisType)
    }
    return nil
}
