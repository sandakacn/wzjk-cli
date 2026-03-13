package cmd

import (
	"fmt"
	"os"

	"wzjk-cli/internal/version"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wzjk-cli",
	Short: "网站监控系统 CLI 工具",
	Long: `wzjk-cli 是网站监控系统的命令行工具，用于管理域名监控、查看证书状态等。

常用命令:
  wzjk-cli login                    # 微信扫码登录
  wzjk-cli domains list             # 列出所有域名
  wzjk-cli domains add example.com  # 添加新域名
  wzjk-cli status                   # 查看监控状态`,
	Version: version.GetVersion(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().String("api-url", "", "API 地址 (默认使用配置文件中的设置)")
}
