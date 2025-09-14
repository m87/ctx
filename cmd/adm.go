package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func newAdminCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "admin",
		Aliases: []string{"adm"},
		Short:   "Admin command",
	}
}

var admCmd = newAdminCmd(bootstrap.CreateManager())
func init() {
	rootCmd.AddCommand(admCmd)
}
