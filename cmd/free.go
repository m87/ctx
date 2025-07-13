package cmd

import (
	localstorage "github.com/m87/ctx/storage/local"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var freeCmd = &cobra.Command{
	Use:     "free",
	Aliases: []string{"f"},
	Short:   "Stop current context",
	Run: func(cmd *cobra.Command, args []string) {
		util.Check(localstorage.CreateManager().Free())
	},
}

func init() {
	rootCmd.AddCommand(freeCmd)
}
