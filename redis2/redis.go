package redis2

import (
    "context"
    . "github.com/go-redis/redis/v8"
    "strings"
    "time"
)

type ClientX interface {
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
    Get(ctx context.Context, k string) *StringCmd
    Del(ctx context.Context, keys ...string) *IntCmd
    Exists(ctx context.Context, keys ...string) *IntCmd
    Type(ctx context.Context, key string) *StatusCmd

    Close() error

    Scan(ctx context.Context, cursor uint64, match string, count int64) *ScanCmd
    ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *ScanCmd

    Incr(ctx context.Context, key string) *IntCmd
    IncrBy(ctx context.Context, key string, value int64) *IntCmd
    TTL(ctx context.Context, key string) *DurationCmd

    SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd
    Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd
    PFAdd(ctx context.Context, key string, els ...interface{}) *IntCmd
    PFCount(ctx context.Context, keys ...string) *IntCmd
    // Scripter
    Eval(ctx context.Context, script string, keys []string, args ...interface{}) *Cmd
    EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *Cmd
    ScriptExists(ctx context.Context, hashes ...string) *BoolSliceCmd
    ScriptLoad(ctx context.Context, script string) *StringCmd
    // hash
    HSet(ctx context.Context, key string, values ...interface{}) *IntCmd
    HGet(ctx context.Context, key, field string) *StringCmd
    HGetAll(ctx context.Context, key string) *StringStringMapCmd
    HLen(ctx context.Context, key string) *IntCmd
    HDel(ctx context.Context, key string, fields ...string) *IntCmd
    HIncrBy(ctx context.Context, key, field string, incr int64) *IntCmd
    HIncrByFloat(ctx context.Context, key, field string, incr float64) *FloatCmd
    HSetNX(ctx context.Context, key, field string, value interface{}) *BoolCmd
    HMGet(ctx context.Context, key string, fields ...string) *SliceCmd
    // set
    SAdd(ctx context.Context, key string, members ...interface{}) *IntCmd
    SRem(ctx context.Context, key string, members ...interface{}) *IntCmd
    SCard(ctx context.Context, key string) *IntCmd
    SIsMember(ctx context.Context, key string, member interface{}) *BoolCmd
    SMembers(ctx context.Context, key string) *StringSliceCmd
    SRandMemberN(ctx context.Context, key string, count int64) *StringSliceCmd
    ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd
    SMIsMember(ctx context.Context, key string, members ...interface{}) *BoolSliceCmd
    ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmd
    ZRandMember(ctx context.Context, key string, count int, withScores bool) *StringSliceCmd
    // zset
    ZCard(ctx context.Context, key string) *IntCmd
    ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd
    ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd
    ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
    ZScore(ctx context.Context, key, member string) *FloatCmd
    ZAdd(ctx context.Context, key string, members ...*Z) *IntCmd
    ZAddNX(ctx context.Context, key string, members ...*Z) *IntCmd
    ZRem(ctx context.Context, key string, members ...interface{}) *IntCmd
    ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmd
    // list
    LPush(ctx context.Context, key string, values ...interface{}) *IntCmd
    LPop(ctx context.Context, key string) *StringCmd
    LRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd
    // cluster
    ForEachShard(ctx context.Context, fn func(ctx context.Context, client *Client) error) error
    ForEachMaster(ctx context.Context, fn func(ctx context.Context, client *Client) error) error
    ForEachSlave(ctx context.Context, fn func(ctx context.Context, client *Client) error) error
}

func NewRedisClient(redisType, hosts, password string) (ClientX, error) {
    switch redisType {
    case "cluster":
        return NewRedisCluster(
            strings.Split(hosts, ","),
            password,
        )
    case "single":
        return NewRedisSingle(
            hosts,
            password,
        )
    default:
        return NewRedisCluster(
            strings.Split(hosts, ","),
            password,
        )
    }
}
