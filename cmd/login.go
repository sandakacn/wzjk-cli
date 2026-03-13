package cmd

import (
	"fmt"

	"wzjk-cli/pkg/api"
	"wzjk-cli/pkg/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "使用 API Key 登录",
	Long:  `使用 API Key 进行身份验证。\n可以在网页版个人资料页面生成 API Key。`,
	Example: `  # 交互式登录（会提示输入 API Key）
  wzjk-cli login --api-url https://wangzhanjiankong.cn

  # 直接指定 API Key 登录
  wzjk-cli login --api-url https://wangzhanjiankong.cn --api-key <your-key>

  # 使用 --token 别名
  wzjk-cli login --api-url https://wangzhanjiankong.cn --token <your-key>`,
	RunE: runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().String("api-url", "", "API 地址（默认 http://localhost:3000）")
	loginCmd.Flags().String("api-key", "", "API Key")
	loginCmd.Flags().String("token", "", "API Key（与 --api-key 相同）")
}

func runLogin(cmd *cobra.Command, args []string) error {
	// Get flags
	apiURL, _ := cmd.Flags().GetString("api-url")
	apiKey, _ := cmd.Flags().GetString("api-key")
	if apiKey == "" {
		apiKey, _ = cmd.Flags().GetString("token")
	}

	// Load existing config to get default API URL
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{APIURL: "http://localhost:3000"}
	}

	// Use flag value or config value
	if apiURL == "" {
		apiURL = cfg.APIURL
	}

	// Warn if not using HTTPS
	if len(apiURL) > 8 && apiURL[:8] != "https://" {
		fmt.Println(color.YellowString("⚠ 警告: 建议使用 HTTPS 以保护 API Key 安全"))
		fmt.Println()
	}

	// If no API key provided, enter interactive mode
	if apiKey == "" {
		return loginInteractive(apiURL)
	}

	// Direct login with provided API key
	return loginWithAPIKey(apiURL, apiKey)
}

func loginInteractive(apiURL string) error {
	fmt.Printf("正在连接到 %s ...\n\n", apiURL)
	fmt.Println("请先在网页版生成 API Key：")
	fmt.Printf("  %s/profile\n\n", apiURL)
	fmt.Println("然后输入您的 API Key（输入内容将被隐藏）：")
	fmt.Println()

	var apiKey string
	prompt := &survey.Password{
		Message: "API Key",
	}
	if err := survey.AskOne(prompt, &apiKey); err != nil {
		return fmt.Errorf("输入失败: %w", err)
	}

	if apiKey == "" {
		return fmt.Errorf("API Key 不能为空")
	}

	fmt.Println()
	return loginWithAPIKey(apiURL, apiKey)
}

func loginWithAPIKey(apiURL, apiKey string) error {
	// Validate API key by calling login endpoint
	client := api.NewClient(apiURL, "")

	result, err := client.LoginWithAPIKey(apiKey)
	if err != nil {
		return fmt.Errorf("登录失败: %w", err)
	}

	// Save config
	newCfg := &config.Config{
		APIURL: apiURL,
		Token:  result.Token,
		User: &config.User{
			ID:    result.User.ID,
			Name:  result.User.Name,
			Email: result.User.Email,
		},
	}

	if err := config.Save(newCfg); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	fmt.Println(color.GreenString("✓ 登录成功！"))
	fmt.Printf("欢迎, %s\n", color.CyanString(result.User.Name))
	fmt.Println()
	fmt.Println("您现在可以使用以下命令：")
	fmt.Println("  wzjk-cli domains list    # 查看域名列表")
	fmt.Println("  wzjk-cli status          # 查看监控状态")

	return nil
}
