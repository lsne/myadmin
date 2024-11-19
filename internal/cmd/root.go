/*
 * @Author: Liu Sainan
 * @Date: 2023-12-07 12:58:05
 */

package cmd

import (
	"fmt"
	"log"
	"myadmin/internal/config"
	"myadmin/internal/server"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "myadmin",
	Short: "myadmin 后端",
	Long:  `一个 web 页面形式的管理平台`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.Server(); err != nil {
			log.Fatalln(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	fmt.Println("欢迎使用 myadmin !!!")
}

func init() {
	// OnInitialize 会在执行任何命令(子命令, 二级命令, 三级命令等)之前执行
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&config.FileName, "config", "c", "./config.toml", "Path to config file")

	// 装载子命令
	rootCmd.AddCommand(
		VersionCmd(),
		UserCmd(),
	)
}
