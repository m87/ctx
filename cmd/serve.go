package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/server"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	var (
		addr string
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the context server",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager, err := bootstrap.CreateSettingsManager()
			if err != nil {
				return err
			}
			manager := bootstrap.CreateManager()
			server := server.NewServer(manager, settingsManager)
			return server.Listen(addr)
		},
	}

	cmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "Address to listen on")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewServeCmd())
}
