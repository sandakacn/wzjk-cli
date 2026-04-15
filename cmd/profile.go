package cmd

import (
	"fmt"

	"wzjk-cli/pkg/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "查看用户信息",
	Long:  `显示当前登录用户的详细信息。`,
	Example: `  # 查看用户信息
  wzjk-cli profile`,
	RunE: runProfile,
}

func init() {
	rootCmd.AddCommand(profileCmd)
}

func runProfile(cmd *cobra.Command, args []string) error {
	cfg, client, err := getClient()
	if err != nil {
		return err
	}

	profile, err := client.GetProfile()
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	fmt.Println()
	fmt.Println(color.CyanString("═══ 用户信息 ═══"))
	fmt.Println()

	if cfg.User != nil {
		fmt.Printf("用户ID:   %s\n", cfg.User.ID)
	}
	fmt.Printf("用户名:   %s\n", color.CyanString(profile.Name))
	if profile.Email != "" {
		fmt.Printf("邮箱:     %s\n", profile.Email)
	}
	fmt.Printf("订阅计划: %s\n", utils.FormatPlan(profile.IsPro))

	fmt.Println()
	fmt.Println(color.CyanString("═══ 账号绑定 ═══"))
	fmt.Println()
	fmt.Printf("微信登录:     %s\n", utils.FormatBool(profile.HasWechat))
	fmt.Printf("微信服务号:   %s\n", utils.FormatBool(profile.HasWechatService))
	fmt.Printf("登录方式:     %s\n", profile.LoginProvider)

	fmt.Println()
	fmt.Println(color.CyanString("═══ 通知设置 ═══"))
	fmt.Println()
	fmt.Printf("邮件通知:     %s\n", utils.FormatBool(profile.NotificationSettings.Email))
	fmt.Printf("微信通知:     %s\n", utils.FormatBool(profile.NotificationSettings.Wechat))

	fmt.Println()
	fmt.Println(color.CyanString("═══ 告警偏好 ═══"))
	fmt.Println()
	fmt.Printf("证书过期告警:     %s\n", utils.FormatBool(profile.AlertPreferences.ExpiryAlert))
	fmt.Println()

	return nil
}
