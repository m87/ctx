package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/spf13/cobra"
)

func NewRestoreContextCmd() *cobra.Command {
	var (
		contextId string
	)

	cmd := &cobra.Command{
		Use:   "context",
		Short: "Restore a context",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := bootstrap.CreateManager()
			if err != nil {
				return err
			}

			if resolveRemoteAddr() != "" {
				if err := remoteRestoreContext(contextId); err != nil {
					return err
				}
			} else {
				err = manager.RestoreContext(contextId)
				if err != nil {
					return err
				}
			}

			cmd.Println("Context restored successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&contextId, "id", "", "ID of the context to restore")
	cmd.MarkFlagRequired("id")

	return cmd
}
