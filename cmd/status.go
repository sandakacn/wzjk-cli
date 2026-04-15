package cmd

import (
	"fmt"

	"wzjk-cli/pkg/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看监控状态概览",
	Long:  `显示用户信息和域名监控的整体状态统计。`,
	Example: `  # 查看状态概览
  wzjk-cli status`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, client, err := getClient()
	if err != nil {
		return err
	}

	// Get user profile
	profile, err := client.GetProfile()
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// Get domains
	domains, err := client.ListDomains()
	if err != nil {
		return fmt.Errorf("获取域名列表失败: %w", err)
	}

	// Calculate stats
	total := len(domains)
	var expiringSoon, inactive int
	for _, d := range domains {
		daysLeft := utils.DaysUntil(d.SSLValidTo)
		if daysLeft <= d.AlertDays {
			expiringSoon++
		}
		if !d.IsActive {
			inactive++
		}
	}

	// Print status
	fmt.Println()
	fmt.Println(color.CyanString("╔══════════════════════════════════════╗"))
	fmt.Println(color.CyanString("║          网站监控状态概览            ║"))
	fmt.Println(color.CyanString("╚══════════════════════════════════════╝"))
	fmt.Println()

	// User info
	fmt.Println(color.WhiteString("【用户信息】"))
	fmt.Printf("  用户名:  %s\n", color.CyanString(profile.Name))
	if profile.Email != "" {
		fmt.Printf("  邮箱:    %s\n", profile.Email)
	}
	fmt.Printf("  订阅:    %s\n", utils.FormatPlan(profile.IsPro))
	fmt.Printf("  API:     %s\n", cfg.APIURL)
	fmt.Println()

	// Domain stats
	fmt.Println(color.WhiteString("【域名统计】"))
	fmt.Printf("  总计:        %d\n", total)
	if expiringSoon > 0 {
		fmt.Printf("  即将过期:    %s\n", color.RedString("%d", expiringSoon))
	} else {
		fmt.Printf("  即将过期:    %d\n", expiringSoon)
	}
	if inactive > 0 {
		fmt.Printf("  已暂停:      %s\n", color.YellowString("%d", inactive))
	} else {
		fmt.Printf("  已暂停:      %d\n", inactive)
	}
	fmt.Println()

	// Notification settings
	fmt.Println(color.WhiteString("【通知设置】"))
	fmt.Printf("  邮件通知:    %s\n", utils.FormatBool(profile.NotificationSettings.Email))
	fmt.Printf("  微信通知:    %s\n", utils.FormatBool(profile.NotificationSettings.Wechat))
	fmt.Println()

	// Alert preferences
	fmt.Println(color.WhiteString("【告警设置】"))
	fmt.Printf("  证书过期:    %s\n", utils.FormatBool(profile.AlertPreferences.ExpiryAlert))
	fmt.Println()

	// Summary
	if expiringSoon > 0 {
		fmt.Println(color.YellowString("⚠ 有 %d 个证书即将过期", expiringSoon))
	} else {
		fmt.Println(color.GreenString("✓ 所有域名证书状态正常"))
	}

	return nil
}
