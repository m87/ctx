package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewDeleteWorkspaceCmd(manager *core.ContextManager) *cobra.Command {
	var workspaceId string

	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Delete a workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			if resolveRemoteAddr() != "" {
				return remoteDeleteWorkspace(workspaceId)
			}
			return manager.WorkspaceRepository.Delete(workspaceId)
		},
	}

	cmd.Flags().StringVarP(&workspaceId, "id", "i", "", "ID of the workspace to delete")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func init() {
	deleteCmd.AddCommand(NewDeleteWorkspaceCmd(bootstrap.CreateManager()))
}
