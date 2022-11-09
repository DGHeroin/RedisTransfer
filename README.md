#Redis Transfer

一个redis迁移助手, 减少心智负担

- [x] 支持string/hash/set/zset/list数据格式
- [x] 支持单点/集群
- [x] 支持多工作者同时迁移
- [x] 跨平台

## 教程
> 单点到单点

```bash
RedisTransfer-darwin-amd64 merge all \
  --source-addr="127.0.0.1:6378" \
  --target-addr="127.0.0.1:6379"
```

> 单点到集群

```bash
RedisTransfer-darwin-amd64 merge all \
  --source-addr="127.0.0.1:6378" \
  --target-addr="127.0.0.1:6379" \
  --targetClient-cluster
```