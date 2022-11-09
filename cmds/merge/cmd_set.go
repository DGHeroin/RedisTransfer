package merge

import (
    "github.com/spf13/cobra"
)

var (
    mergeSetCmd = &cobra.Command{
        Use:   "set <args>",
        Short: "迁移所有set数据",
        RunE: func(cmd *cobra.Command, args []string) error {
            var key string
            if len(args) >= 1 {
                key = args[0]
            }
            return handleMerge(key, "set")
        },
        PreRunE: func(cmd *cobra.Command, args []string) error {
            if err := connectToRedis(); err != nil {
                return err
            }
            return nil
        },
    }
)
