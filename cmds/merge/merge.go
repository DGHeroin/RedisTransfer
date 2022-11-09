package merge

import "github.com/spf13/cobra"

var (
    Cmd = &cobra.Command{
        Use:   "merge command <args>",
        Short: "迁移数据",
    }
)

func init() {
    Cmd.PersistentFlags().StringVar(&sourceAddr, "source-addr", "127.0.0.1:6379", "源redis地址, 多个ip用','隔开")
    Cmd.PersistentFlags().StringVar(&targetAddr, "target-addr", "127.0.0.1:6379", "目标redis地址, 多个ip用','隔开")
    Cmd.PersistentFlags().IntVar(&sourceDB, "sourceClient-database", 0, "源db")
    Cmd.PersistentFlags().IntVar(&targetDB, "targetClient-database", 0, "目标db")
    Cmd.PersistentFlags().StringVar(&sourceAuth, "sourceClient-auth", "", "源密码")
    Cmd.PersistentFlags().StringVar(&targetAuth, "targetClient-auth", "", "目标密码")
    Cmd.PersistentFlags().BoolVar(&sourceCluster, "sourceClient-cluster", false, "源redis是否是集群")
    Cmd.PersistentFlags().BoolVar(&targetCluster, "targetClient-cluster", false, "目标redis是否是集群")
    Cmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "详细日志")
    Cmd.PersistentFlags().IntVar(&workerNum, "worker", 100, "工作者数量")

    Cmd.AddCommand(mergeAllCmd, mergeStringCmd, mergeHashCmd, mergeSetCmd, mergeZSetgCmd, mergeListCmd)
}
