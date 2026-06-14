package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/spf13/cobra"
)

func NewEditWorkspaceCmd() *cobra.Command {
	var (
		workspaceId string
		name        string
		description string
	)

	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Edit a workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := bootstrap.CreateManager()

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
			if cmd.Flags().Changed("description") {
				workspace.Description = description
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
	cmd.Flags().StringVar(&description, "description", "", "New description of the workspace")
	_ = cmd.MarkFlagRequired("id")
	return cmd

}

func init() {
	editCmd.AddCommand(NewEditWorkspaceCmd())
}
