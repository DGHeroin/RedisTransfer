package merge

import (
    "RedisTransfer/redis2"
    "sync"
    "time"
)

var (
    sourceAddr    string
    targetAddr    string
    sourceAuth    string
    targetAuth    string
    sourceCluster bool
    targetCluster bool
    sourceDB      int
    targetDB      int
    verbose       bool
    workerNum     int
)
var (
    sourceClient redis2.ClientX
    targetClient redis2.ClientX
)
var (
    count       uint32
    countString uint32
    countHash   uint32
    countSet    uint32
    countZSet   uint32
    countList   uint32
    countTTL    uint32
    mu          sync.Mutex
    mm          = map[string]bool{}
    uptime      = time.Now()
)
