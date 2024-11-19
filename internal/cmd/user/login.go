/*
 * @Author: Liu Sainan
 * @Date: 2023-12-08 22:33:52
 */

package user

import (
	"fmt"

	"github.com/spf13/cobra"
)

func LoginCmd() *cobra.Command {
	var username string
	var password string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "用户登录操作",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("用户登录成功", username)
		},
	}
	cmd.Flags().StringVarP(&username, "username", "u", "user001", "用户名称")
	cmd.Flags().StringVarP(&password, "password", "p", "", "用户密码")
	return cmd
}
