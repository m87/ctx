package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/server"
	"github.com/spf13/cobra"
)

func NewServeCmd(manager *core.ContextManager) *cobra.Command {
	var (
		addr string
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the context server",
		RunE: func(cmd *cobra.Command, args []string) error {
			server := server.NewServer(manager)
			return server.Listen(addr)
		},
	}

	cmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "Address to listen on")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewServeCmd(bootstrap.CreateManager()))
}
