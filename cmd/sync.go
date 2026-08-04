package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/spf13/cobra"
)

func NewSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync command with remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := bootstrap.CreateManager()
			if err != nil {
				return err
			}
			err = manager.Sync(resolveRemoteAddr())
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(NewSyncCmd())
}
