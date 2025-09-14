package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func newLabelIntervalCmd(manager *core.ContextManager) *cobra.Command {

	return &cobra.Command{
		Use:     "labelInterval",
		Aliases: []string{"li"},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}

func init() {
	intervalCmd.AddCommand(newLabelIntervalCmd(bootstrap.CreateManager()))
}
