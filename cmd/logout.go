package cmd

import (
	"fmt"

	"wzjk-cli/pkg/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "退出登录",
	Long:  `清除本地保存的登录凭证。`,
	Example: `  # 退出登录
  wzjk-cli logout

  # 强制退出，不确认
  wzjk-cli logout --force`,
	RunE: runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
	logoutCmd.Flags().Bool("force", false, "强制退出，不确认")
}

func runLogout(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	// Check if logged in
	if !config.IsLoggedIn() {
		fmt.Println("您尚未登录")
		return nil
	}

	// Confirm logout unless --force is used
	if !force {
		confirmed := false
		prompt := &survey.Confirm{
			Message: "确定要退出登录吗？",
			Default: false,
		}
		if err := survey.AskOne(prompt, &confirmed); err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("已取消")
			return nil
		}
	}

	// Clear config
	if err := config.Clear(); err != nil {
		return fmt.Errorf("清除配置失败: %w", err)
	}

	fmt.Println(color.GreenString("✓ 已退出登录"))
	return nil
}
