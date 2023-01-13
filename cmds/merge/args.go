package merge

import (
    "RedisTransfer/log"
    "RedisTransfer/redis2"
    "github.com/spf13/cobra"
    "sync"
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
    status        bool
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
)
var (
    Cmd = &cobra.Command{
        Use:   "merge command <args>",
        Short: "迁移数据",
    }
)

func init() {
    Cmd.PersistentFlags().StringVar(&sourceAddr, "s-addr", "127.0.0.1:6379", "源redis地址, 多个ip用','隔开")
    Cmd.PersistentFlags().StringVar(&targetAddr, "t-addr", "127.0.0.1:6379", "目标redis地址, 多个ip用','隔开")
    Cmd.PersistentFlags().IntVar(&sourceDB, "s-database", 0, "源db")
    Cmd.PersistentFlags().IntVar(&targetDB, "t-database", 0, "目标db")
    Cmd.PersistentFlags().StringVar(&sourceAuth, "s-auth", "", "源密码")
    Cmd.PersistentFlags().StringVar(&targetAuth, "t-auth", "", "目标密码")
    Cmd.PersistentFlags().BoolVar(&sourceCluster, "s-cluster", false, "源redis是否是集群")
    Cmd.PersistentFlags().BoolVar(&targetCluster, "t-cluster", false, "目标redis是否是集群")
    Cmd.PersistentFlags().BoolVar(&log.Verbose, "verbose", false, "详细日志")
    Cmd.PersistentFlags().BoolVar(&status, "status", false, "实时数据")
    Cmd.PersistentFlags().IntVar(&workerNum, "worker", 100, "工作者数量")

    Cmd.AddCommand(mergeAllCmd, mergeStringCmd, mergeHashCmd, mergeSetCmd, mergeZSetgCmd, mergeListCmd)
}
