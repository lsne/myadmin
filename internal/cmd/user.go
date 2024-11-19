/*
 * @Author: Liu Sainan
 * @Date: 2023-12-07 13:23:03
 */

package cmd

import (
	"myadmin/internal/cmd/user"

	"github.com/spf13/cobra"
)

func UserCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "user",
		Short: "用户相关操作",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	// 装载命令
	cmd.AddCommand(
		user.LoginCmd(),
	)

	return cmd
}
