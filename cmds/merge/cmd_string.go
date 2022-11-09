package merge

import (
    "github.com/spf13/cobra"
)

var (
    mergeStringCmd = &cobra.Command{
        Use:   "string <args>",
        Short: "迁移所有string数据",
        RunE: func(cmd *cobra.Command, args []string) error {
            var key string
            if len(args) >= 1 {
                key = args[0]
            }
            return handleMerge(key, "string")
        },
        PreRunE: func(cmd *cobra.Command, args []string) error {
            if err := connectToRedis(); err != nil {
                return err
            }
            return nil
        },
    }
)
