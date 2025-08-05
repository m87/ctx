//go:build preview
package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/server"
	"github.com/spf13/cobra"
)

func NewServeCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, ars []string) {
			server.Serve(manager)
		},
	}
}

func init() {
	rootCmd.AddCommand(NewServeCmd(bootstrap.CreateManager()))
}
