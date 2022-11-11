package redis2

import (
    "context"
    "github.com/go-redis/redis/v8"
)

type single struct {
    *redis.Client
}

func (s *single) ForEachShard(ctx context.Context, fn func(ctx context.Context, client *redis.Client) error) error {
    return fn(ctx, s.Client)
}

func (s *single) ForEachMaster(ctx context.Context, fn func(ctx context.Context, client *redis.Client) error) error {
    return fn(ctx, s.Client)
}

func (s *single) ForEachSlave(ctx context.Context, fn func(ctx context.Context, client *redis.Client) error) error {
    return fn(ctx, s.Client)
}

func NewRedisSingle(addr string, password string) (ClientX, error) {
    var cli *single
    client := redis.NewClient(&redis.Options{
        Addr:         addr,
        Password:     password,
        ReadTimeout:  -1,
        WriteTimeout: -1,
        PoolSize:     5000,
        PoolTimeout:  -1,
    })

    if err := client.Ping(context.Background()).Err(); err != nil {
        return nil, err
    }

    cli = &single{client}
    return cli, nil
}
