package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tuoaitang/calculator/finance/pension"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "calc",
	Short: "计算器",
	Long: `当前支持的任务命令：
1. calc pension 15 6500 (计算缴费年限为15年,当前月薪为6500元的情况下养老金待遇) 
	`,
}

func init() {
	pensionCmd := &cobra.Command{
		Use:   "pension",
		Short: "退休金待遇",
		Long:  "退休金待遇",
		Run:   pension.CalcPension,
	}
	RootCmd.AddCommand(pensionCmd)
}
