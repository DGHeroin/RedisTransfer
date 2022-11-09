package redis2

import (
    "context"
    "github.com/go-redis/redis/v8"
)

type cluster struct {
    *redis.ClusterClient
}

func NewRedisCluster(addr []string, password string) (ClientX, error) {
    var cli *cluster
    client := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs:        addr,
        Password:     password,
        ReadTimeout:  -1,
        WriteTimeout: -1,
        PoolSize:     5000,
        PoolTimeout:  -1,
    })
    err := client.ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
        return shard.Ping(ctx).Err()
    })
    if err != nil {
        return nil, err
    }
    cli = &cluster{}
    return cli, nil
}
