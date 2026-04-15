package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"wzjk-cli/pkg/api"
	"wzjk-cli/pkg/config"
	"wzjk-cli/pkg/ssl"
	"wzjk-cli/pkg/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "域名管理",
	Long:  `管理您监控的域名，包括添加、删除、查看列表等操作。`,
}

var domainsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "列出所有监控的域名",
	Example: `  # 列出所有域名
  wzjk-cli domains list

  # 只显示有告警的域名
  wzjk-cli domains list --alerts-only

  # JSON 格式输出
  wzjk-cli domains list --format json`,
	RunE: runDomainsList,
}

var domainsAddCmd = &cobra.Command{
	Use:   "add <domain>",
	Short: "添加新的监控域名",
	Args:  cobra.ExactArgs(1),
	Example: `  # 添加域名（使用默认设置）
  wzjk-cli domains add example.com

  # 指定端口和告警天数
  wzjk-cli domains add example.com --port 443 --alert-days 30

  # 本地 SSL 检查（不经过服务端）
  wzjk-cli domains add example.com --local

  # 跳过 SSL 检查
  wzjk-cli domains add example.com --skip-check`,
	RunE: runDomainsAdd,
}

var domainsDeleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "删除域名监控",
	Args:    cobra.ExactArgs(1),
	Example: `  # 删除指定域名
  wzjk-cli domains delete <domain-id>

  # 强制删除，不确认
  wzjk-cli domains delete <domain-id> --force`,
	RunE: runDomainsDelete,
}

var domainsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "更新域名设置",
	Args:  cobra.ExactArgs(1),
	Example: `  # 更新告警天数
  wzjk-cli domains update <domain-id> --alert-days 14

  # 禁用监控
  wzjk-cli domains update <domain-id> --active false`,
	RunE: runDomainsUpdate,
}

var domainsCheckCmd = &cobra.Command{
	Use:   "check <domain>",
	Short: "检查域名的 SSL 证书",
	Args:  cobra.ExactArgs(1),
	Example: `  # 检查域名 SSL 证书（默认使用服务端检查）
  wzjk-cli domains check example.com

  # 本地检查（不经过服务端）
  wzjk-cli domains check example.com --local

  # 指定端口
  wzjk-cli domains check example.com --port 8443`,
	RunE: runDomainsCheck,
}

func init() {
	rootCmd.AddCommand(domainsCmd)
	domainsCmd.AddCommand(domainsListCmd)
	domainsCmd.AddCommand(domainsAddCmd)
	domainsCmd.AddCommand(domainsDeleteCmd)
	domainsCmd.AddCommand(domainsUpdateCmd)
	domainsCmd.AddCommand(domainsCheckCmd)

	// List flags
	domainsListCmd.Flags().String("format", "table", "输出格式: table, json")
	domainsListCmd.Flags().Bool("alerts-only", false, "只显示有告警的域名")

	// Add flags
	domainsAddCmd.Flags().Int("port", 443, "端口")
	domainsAddCmd.Flags().Int("alert-days", 20, "告警提前天数")
	domainsAddCmd.Flags().String("type", "http_tls", "检查类型: http_tls, https, tcp, tls")
	domainsAddCmd.Flags().Bool("skip-check", false, "跳过 SSL 检查")
	domainsAddCmd.Flags().Bool("local", false, "本地 SSL 检查（不经过服务端）")

	// Delete flags
	domainsDeleteCmd.Flags().Bool("force", false, "强制删除，不确认")

	// Update flags
	domainsUpdateCmd.Flags().Int("alert-days", 0, "告警提前天数")
	domainsUpdateCmd.Flags().String("active", "", "启用/禁用监控: true, false")

	// Check flags
	domainsCheckCmd.Flags().Int("port", 443, "端口")
	domainsCheckCmd.Flags().Bool("local", false, "本地检查（不经过服务端）")
}

func runDomainsList(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	alertsOnly, _ := cmd.Flags().GetBool("alerts-only")

	cfg, client, err := getClient()
	if err != nil {
		return err
	}

	domains, err := client.ListDomains()
	if err != nil {
		return fmt.Errorf("获取域名列表失败: %w", err)
	}

	// Sort by SSL expiry time (earliest first, nil values at the end)
	sort.Slice(domains, func(i, j int) bool {
		if domains[i].SSLValidTo == nil {
			return false
		}
		if domains[j].SSLValidTo == nil {
			return true
		}
		return domains[i].SSLValidTo.Before(*domains[j].SSLValidTo)
	})

	// Filter by alerts if requested
	if alertsOnly {
		var filtered []api.Domain
		for _, d := range domains {
			daysLeft := utils.DaysUntil(d.SSLValidTo)
			if daysLeft <= d.AlertDays {
				filtered = append(filtered, d)
			}
		}
		domains = filtered
	}

	if format == "json" {
		// JSON output
		fmt.Printf("%+v\n", domains)
		return nil
	}

	// Table output
	if len(domains) == 0 {
		fmt.Println("暂无监控的域名")
		return nil
	}

	fmt.Printf("API URL: %s\n\n", cfg.APIURL)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "域名", "端口", "过期时间", "剩余", "状态"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, d := range domains {
		daysLeft := utils.DaysUntil(d.SSLValidTo)
		daysStr := utils.FormatDaysLeft(daysLeft, d.AlertDays)

		// Status
		statusStr := color.GreenString("正常")
		if daysLeft <= d.AlertDays {
			statusStr = color.RedString("即将过期")
		} else if daysLeft <= d.AlertDays+7 {
			statusStr = color.YellowString("需注意")
		}
		if !d.IsActive {
			statusStr = color.WhiteString("已暂停")
		}

		table.Append([]string{
			utils.TruncateID(d.ID),
			d.Domain,
			strconv.Itoa(d.Port),
			utils.FormatTime(d.SSLValidTo),
			daysStr,
			statusStr,
		})
	}

	table.Render()

	fmt.Printf("\n总计: %d 个域名\n", len(domains))
	return nil
}

func runDomainsAdd(cmd *cobra.Command, args []string) error {
	domain := args[0]
	port, _ := cmd.Flags().GetInt("port")
	alertDays, _ := cmd.Flags().GetInt("alert-days")
	checkType, _ := cmd.Flags().GetString("type")
	skipCheck, _ := cmd.Flags().GetBool("skip-check")
	local, _ := cmd.Flags().GetBool("local")

	_, client, err := getClient()
	if err != nil {
		return err
	}

	// Check SSL if not skipped
	if !skipCheck {
		fmt.Println("正在检查 SSL 证书...")
		var sslInfo *api.SSLInfo
		var checkErr error

		if local {
			// Local SSL check
			checker := ssl.NewChecker(30 * time.Second)
			sslInfo, checkErr = checker.CheckSSL(domain, port)
		} else {
			// Server-side SSL check
			sslInfo, checkErr = client.CheckSSL(domain, port)
		}

		if checkErr != nil {
			fmt.Printf("SSL 检查警告: %v\n", checkErr)
		} else {
			displaySSLInfo(sslInfo)
		}
	}

	// Confirm
	confirmed := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("确认添加域名 %s?", domain),
		Default: true,
	}
	if err := survey.AskOne(prompt, &confirmed); err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("已取消")
		return nil
	}

	// Add domain
	req := api.AddDomainRequest{
		Domain:    domain,
		Port:      port,
		CheckType: checkType,
		AlertDays: alertDays,
	}

	newDomain, err := client.AddDomain(req)
	if err != nil {
		return fmt.Errorf("添加域名失败: %w", err)
	}

	fmt.Println(color.GreenString("✓ 域名添加成功"))
	fmt.Printf("ID: %s\n", newDomain.ID)
	return nil
}

func runDomainsDelete(cmd *cobra.Command, args []string) error {
	id := args[0]
	force, _ := cmd.Flags().GetBool("force")

	_, client, err := getClient()
	if err != nil {
		return err
	}

	// Confirm unless --force
	if !force {
		confirmed := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("确定要删除域名 %s 吗?", id),
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

	if err := client.DeleteDomain(id); err != nil {
		return fmt.Errorf("删除域名失败: %w", err)
	}

	fmt.Println(color.GreenString("✓ 域名已删除"))
	return nil
}

func runDomainsUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	alertDays, _ := cmd.Flags().GetInt("alert-days")
	activeStr, _ := cmd.Flags().GetString("active")

	_, client, err := getClient()
	if err != nil {
		return err
	}

	req := api.UpdateDomainRequest{}
	if alertDays > 0 {
		req.AlertDays = alertDays
	}
	if activeStr != "" {
		active := activeStr == "true"
		req.IsActive = &active
	}

	updated, err := client.UpdateDomain(id, req)
	if err != nil {
		return fmt.Errorf("更新域名失败: %w", err)
	}

	fmt.Println(color.GreenString("✓ 域名已更新"))
	fmt.Printf("ID: %s\n", updated.ID)
	return nil
}

func runDomainsCheck(cmd *cobra.Command, args []string) error {
	domain := args[0]
	port, _ := cmd.Flags().GetInt("port")
	local, _ := cmd.Flags().GetBool("local")

	fmt.Printf("正在检查 %s:%d ...\n", domain, port)

	var info *api.SSLInfo
	var err error

	if local {
		// Local SSL check
		checker := ssl.NewChecker(30 * time.Second)
		info, err = checker.CheckSSL(domain, port)
		if err != nil {
			return fmt.Errorf("本地 SSL 检查失败: %w", err)
		}
	} else {
		// Server-side SSL check
		_, client, err := getClient()
		if err != nil {
			return err
		}
		info, err = client.CheckSSL(domain, port)
		if err != nil {
			return fmt.Errorf("SSL 检查失败: %w", err)
		}
	}

	displaySSLInfo(info)
	return nil
}

func displaySSLInfo(info *api.SSLInfo) {
	fmt.Println()
	fmt.Printf("域名:    %s\n", info.Domain)
	fmt.Printf("颁发者:  %s\n", info.Issuer)
	fmt.Printf("主题:    %s\n", info.Subject)
	fmt.Printf("生效时间: %s\n", info.ValidFrom)
	fmt.Printf("过期时间: %s\n", info.ValidTo)

	daysStr := fmt.Sprintf("%d 天", info.DaysUntilExpiry)
	if info.DaysUntilExpiry <= 7 {
		daysStr = color.RedString(daysStr)
	} else if info.DaysUntilExpiry <= 30 {
		daysStr = color.YellowString(daysStr)
	} else {
		daysStr = color.GreenString(daysStr)
	}
	fmt.Printf("剩余天数: %s\n", daysStr)

	if info.DomainMismatch {
		fmt.Println()
		fmt.Println(color.RedString("⚠ 警告: 证书域名不匹配！"))
	}
	if !info.IsValid {
		fmt.Println()
		fmt.Println(color.RedString("✗ 证书无效"))
	}
	fmt.Println()
}

// getClient returns config and API client
func getClient() (*config.Config, *api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("加载配置失败: %w", err)
	}

	if cfg.Token == "" {
		return nil, nil, fmt.Errorf("未登录，请先运行: wzjk-cli login")
	}

	client := api.NewClient(cfg.APIURL, cfg.Token)
	return cfg, client, nil
}
