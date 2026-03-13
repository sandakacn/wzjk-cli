package utils

import (
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// NewTable creates a new table writer with default settings
func NewTable(headers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	return table
}

// StatusOK returns a green "✓" status
func StatusOK() string {
	return color.GreenString("✓")
}

// StatusError returns a red "✗" status
func StatusError() string {
	return color.RedString("✗")
}

// StatusWarning returns a yellow "⚠" status
func StatusWarning() string {
	return color.YellowString("⚠")
}
