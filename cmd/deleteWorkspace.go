package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/spf13/cobra"
)

func NewDeleteWorkspaceCmd() *cobra.Command {
	var workspaceId string

	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Delete a workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			if resolveRemoteAddr() != "" {
				return remoteDeleteWorkspace(workspaceId)
			}
			manager, err := bootstrap.CreateManager()
			if err != nil {
				return err
			}
			return manager.DeleteWorkspace(workspaceId)
		},
	}

	cmd.Flags().StringVarP(&workspaceId, "id", "i", "", "ID of the workspace to delete")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func init() {
	deleteCmd.AddCommand(NewDeleteWorkspaceCmd())
}
