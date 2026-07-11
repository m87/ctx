package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/spf13/cobra"
)

func NewArchiveContextCmd() *cobra.Command {
	var (
		contextId string
	)

	cmd := &cobra.Command{
		Use:   "context",
		Short: "Archive a context",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := bootstrap.CreateManager()
			if err != nil {
				return err
			}

			if resolveRemoteAddr() != "" {
				if err := remoteArchiveContext(contextId); err != nil {
					return err
				}
			} else {
				err = manager.ArchiveContext(contextId)
				if err != nil {
					return err
				}
			}

			cmd.Println("Context archived successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&contextId, "id", "", "ID of the context to archive")
	cmd.MarkFlagRequired("id")

	return cmd
}
