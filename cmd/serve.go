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
			port, _ := cmd.Flags().GetString("port")
			server.Serve(manager, port)
		},
	}
}

func init() {
	serveCmd := NewServeCmd(bootstrap.CreateManager())
	serveCmd.Flags().StringP("port", "p", "8080", "server port")
	rootCmd.AddCommand(serveCmd)
}
