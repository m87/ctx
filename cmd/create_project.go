package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewCreateProjectCmd() *cobra.Command {

	var (
		name        string
		parentId    string
		workspaceId string
	)

	cmd := &cobra.Command{
		Use:   "project",
		Short: "Create a new project",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := bootstrap.CreateManager()
			if err != nil {
				return err
			}
			project := &core.Project{
				Name:        name,
				ParentId:    parentId,
				WorkspaceId: workspaceId,
			}

			if resolveRemoteAddr() != "" {
				if err := remoteClient().CreateProject(project); err != nil {
					return err
				}
			} else {
				id, err := manager.ProjectRepository.Save(project)
				if err != nil {
					return err
				}
				project.Id = id
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the project")
	cmd.Flags().StringVarP(&parentId, "parent", "p", "", "Parent ID of the project")
	cmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "Workspace ID of the project")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("workspace")

	return cmd
}

func init() {
	createCmd.AddCommand(NewCreateProjectCmd())
}
