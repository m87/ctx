package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewCreateWorkspaceCmd() *cobra.Command {
	var (
		name string
	)

	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Create a new workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := bootstrap.CreateManager()

			workspace := &core.Workspace{
				Name: name,
			}

			if resolveRemoteAddr() != "" {
				if err := remoteCreateWorkspace(workspace); err != nil {
					return err
				}
			} else {
				id, err := manager.WorkspaceRepository.Save(workspace)
				if err != nil {
					return err
				}
				workspace.Id = id
			}

			return printOutput(cmd, workspace, func() string {
				return "Workspace created successfully"
			}, nil)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the workspace")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func init() {
	createCmd.AddCommand(NewCreateWorkspaceCmd())
}
