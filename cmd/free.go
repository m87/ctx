package cmd

import (
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var freeCmd = &cobra.Command{
	Use:     "free",
	Aliases: []string{"f"},
	Short:   "Stop current context",
	Run: func(cmd *cobra.Command, args []string) {
		util.Check(core.CreateManager().Free())
	},
}

func init() {
	rootCmd.AddCommand(freeCmd)
}
