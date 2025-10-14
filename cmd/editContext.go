package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func newEditContextCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use: "context",
	}
}

var editContextCmd = newEditContextCmd(bootstrap.CreateManager())
func init() {
	editCmd.AddCommand(editContextCmd)
}
