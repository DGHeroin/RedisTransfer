package cmds

import (
    "RedisTransfer/cmds/dump"
    "RedisTransfer/cmds/load"
    "RedisTransfer/cmds/merge"
    "github.com/spf13/cobra"
    "os"
)

var (
    rootCmd = &cobra.Command{
        Use: "RedisTransfer",
        Long: `redis迁移助手
`,
    }
)

func init() {
    rootCmd.AddCommand(merge.Cmd)
    rootCmd.AddCommand(dump.Cmd)
    rootCmd.AddCommand(load.Cmd)
}
func Run() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(-1)
    }
}
