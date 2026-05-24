package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewEditWorkspaceCmd(manager *core.ContextManager) *cobra.Command {
	var (
		workspaceId string
		name        string
	)

	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Edit a workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, err := manager.WorkspaceRepository.GetById(workspaceId)
			if err != nil {
				return err
			}
			if workspace == nil {
				return nil
			}

			if name != "" {
				workspace.Name = name
			}

			if resolveRemoteAddr() != "" {
				return remoteUpdateWorkspace(workspace)
			}
			_, err = manager.WorkspaceRepository.Save(workspace)
			if err != nil {
				return err
			}
			return printOutput(cmd, workspace, func() string {
				return "Workspace updated successfully"
			}, nil)
		},
	}

	cmd.Flags().StringVarP(&workspaceId, "id", "i", "", "ID of the workspace to edit")
	cmd.Flags().StringVarP(&name, "name", "n", "", "New name of the workspace")
	_ = cmd.MarkFlagRequired("id")
	return cmd

}

func init() {
	editCmd.AddCommand(NewEditWorkspaceCmd(bootstrap.CreateManager()))
}
