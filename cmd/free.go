package cmd

import (
	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var freeCmd = &cobra.Command{
	Use:   "free",
	Short: "Stop current context",
	Run: func(cmd *cobra.Command, args []string) {
		util.Check(ctx.CreateManager().Free())
	},
}

func init() {
	rootCmd.AddCommand(freeCmd)
}
