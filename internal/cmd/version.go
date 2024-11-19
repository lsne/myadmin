/*
 * @Author: Liu Sainan
 * @Date: 2023-12-07 13:10:54
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var _version string
var _gitsha string

func VersionCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "版本信息",
		// PersistentPreRun 会在运行当前命令之前执行, 会传递给他的子命令。 除非子命令也定义了自己的 PersistentPreRun
		PersistentPreRun: func(cmd *cobra.Command, args []string) { fmt.Println("version info:") },
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("    myadmin version: myadmin=%s GitSHA=%s\n", _version, _gitsha)
		},
	}

	return cmd
}
