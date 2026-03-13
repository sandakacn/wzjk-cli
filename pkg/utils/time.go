package utils

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// FormatTime formats a time pointer to a string
func FormatTime(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return t.Format("2006-01-02")
}

// FormatDateTime formats a time pointer to a datetime string
func FormatDateTime(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return t.Format("2006-01-02 15:04")
}

// DaysUntil returns the number of days until a date
func DaysUntil(t *time.Time) int {
	if t == nil {
		return 0
	}
	return int(t.Sub(time.Now()).Hours() / 24)
}

// FormatDaysLeft formats days left with color coding
func FormatDaysLeft(days int, alertDays int) string {
	if days < 0 {
		return color.RedString("已过期")
	}
	if days <= alertDays {
		return color.RedString("%d 天", days)
	}
	if days <= alertDays+7 {
		return color.YellowString("%d 天", days)
	}
	return color.GreenString("%d 天", days)
}

// FormatAvailabilityStatus formats availability status
func FormatAvailabilityStatus(available bool) string {
	if available {
		return color.GreenString("正常")
	}
	return color.RedString("异常")
}

// TruncateID truncates an ID to show only first 8 characters
func TruncateID(id string) string {
	if len(id) <= 8 {
		return id
	}
	return id[:8]
}

// FormatBool returns a colored yes/no for a boolean
func FormatBool(v bool) string {
	if v {
		return color.GreenString("是")
	}
	return color.RedString("否")
}

// FormatPlan returns a colored plan name
func FormatPlan(isPro bool) string {
	if isPro {
		return color.CyanString("Pro")
	}
	return color.WhiteString("免费")
}

// ParseDate parses a date string in YYYY-MM-DD format
func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f秒", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0f分钟", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.0f小时", d.Hours())
	}
	return fmt.Sprintf("%.0f天", d.Hours()/24)
}
