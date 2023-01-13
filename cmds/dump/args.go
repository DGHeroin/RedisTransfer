package dump

import (
    "RedisTransfer/log"
    "RedisTransfer/redis2"
    "encoding/gob"
    "encoding/json"
    "os"
    "sync"
)

var (
    sourceAddr    string
    sourceAuth    string
    sourceCluster bool
    sourceDB      int
    codecType     string
    targetFile    string
    workerNum     int
    status        bool
)

func init() {
    Cmd.PersistentFlags().StringVar(&sourceAddr, "s-addr", "127.0.0.1:6379", "源redis地址, 多个ip用','隔开")
    Cmd.PersistentFlags().IntVar(&sourceDB, "s-database", 0, "源db")
    Cmd.PersistentFlags().StringVar(&sourceAuth, "s-auth", "", "源密码")
    Cmd.PersistentFlags().BoolVar(&sourceCluster, "s-cluster", false, "源redis是否是集群")
    Cmd.PersistentFlags().StringVar(&targetFile, "t", "dump.rd", "保存路径")
    Cmd.PersistentFlags().StringVar(&codecType, "codec", "json", "编码类型")
    Cmd.PersistentFlags().IntVar(&workerNum, "worker", 100, "工作者数量")
    Cmd.PersistentFlags().BoolVar(&log.Verbose, "verbose", false, "详细日志")
    Cmd.PersistentFlags().BoolVar(&status, "status", false, "实时数据")
}

var (
    sourceClient redis2.ClientX
    outputFile   *os.File
    mu           sync.Mutex
    outputStream interface {
        Encode(v interface{}) error
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
        if sourceCluster {
            redisType = "cluster"
        }
        sourceClient, err = redis2.NewRedisClient(redisType, sourceAddr, sourceAuth)
        if err != nil {
            return err
        }
    }
    // 保存路径
    {
        outputFile, err = os.Create(targetFile)
        if err != nil {
            return err
        }
        switch codecType {
        case "json":
            outputStream = json.NewEncoder(outputFile)
        case "gob":
            outputStream = gob.NewEncoder(outputFile)
        default:
            outputStream = json.NewEncoder(outputFile)
        }
    }
    return nil
}
