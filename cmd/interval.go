package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func newIntervalCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use: "interval",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}

var intervalCmd = newIntervalCmd(bootstrap.CreateManager())
func init() {
	rootCmd.AddCommand(intervalCmd)
}
