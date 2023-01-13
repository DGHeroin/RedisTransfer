package load

import (
    "RedisTransfer/log"
    "RedisTransfer/redis2"
    "encoding/gob"
    "encoding/json"
    "os"
    "sync"
)

var (
    targetAddr    string
    targetAuth    string
    targetCluster bool
    targetDB      int
    status        bool
    workerNum     int
    sourceFile    string
    codecType     string
    preDelete     bool
)

func init() {
    Cmd.PersistentFlags().StringVar(&targetAddr, "t-addr", "127.0.0.1:6379", "目标redis地址, 多个ip用','隔开")
    Cmd.PersistentFlags().IntVar(&targetDB, "t-database", 0, "目标db")
    Cmd.PersistentFlags().StringVar(&targetAuth, "t-auth", "", "目标密码")
    Cmd.PersistentFlags().BoolVar(&targetCluster, "t-cluster", false, "目标redis是否是集群")
    Cmd.PersistentFlags().BoolVar(&log.Verbose, "verbose", false, "详细日志")
    Cmd.PersistentFlags().BoolVar(&status, "status", false, "实时数据")
    Cmd.PersistentFlags().IntVar(&workerNum, "worker", 100, "工作者数量")
    Cmd.PersistentFlags().StringVar(&sourceFile, "t", "dump.rd", "输入路径")
    Cmd.PersistentFlags().StringVar(&codecType, "codec", "json", "编码类型")
    Cmd.PersistentFlags().BoolVar(&preDelete, "delete", false, "添加前先删除")
}

var (
    targetClient redis2.ClientX
    inputFile    *os.File
    mu           sync.Mutex
    inputStream  interface {
        Decode(v interface{}) error
    }
    count       uint32
    countString uint32
    countHash   uint32
    countSet    uint32
    countZSet   uint32
    countList   uint32
    countTTL    uint32
    mm          = map[string]bool{}
)

func connectToRedis() (err error) {
    // 源
    {
        redisType := "single"
        if targetCluster {
            redisType = "cluster"
        }
        targetClient, err = redis2.NewRedisClient(redisType, targetAddr, targetAuth)
        if err != nil {
            return err
        }
    }
    // 保存路径
    {
        inputFile, err = os.Open(sourceFile)
        if err != nil {
            return err
        }
        switch codecType {
        case "json":
            inputStream = json.NewDecoder(inputFile)
        case "gob":
            inputStream = gob.NewDecoder(inputFile)
        default:
            inputStream = json.NewDecoder(inputFile)
        }
    }
    return nil
}
