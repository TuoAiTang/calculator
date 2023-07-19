package main

import (
	"github.com/spf13/cobra"
	"github.com/tuoaitang/calculator/cmd"
)

func main() {
	cobra.CheckErr(cmd.RootCmd.Execute())
}
