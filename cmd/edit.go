package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func newEditCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use: "edit",
	}
}

var editCmd = newEditCmd(bootstrap.CreateManager())

func init() {
	rootCmd.AddCommand(editCmd)
}
