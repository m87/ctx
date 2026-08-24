package cmd

import (
	"fmt"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewEditProjectCmd() *cobra.Command {
	var (
		projectId string
		name      string
	)

	cmd := &cobra.Command{
		Use:   "project",
		Short: "Edit a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := bootstrap.CreateManager()
			if err != nil {
				return err
			}
			var project *core.Project
			if resolveRemoteAddr() != "" {
				project, err = remoteClient().GetProject(projectId)
				if err != nil {
					return err
				}
			} else {
				project, err = manager.ProjectRepository.GetById(projectId)
				if err != nil {
					return err
				}
			}
			if project == nil {
				return fmt.Errorf("project not found")
			}

			if name != "" {
				project.Name = name
			}

			if resolveRemoteAddr() != "" {
				return remoteClient().UpdateProject(project)
			}
			_, err = manager.ProjectRepository.Save(project)
			return err
		},
	}

	cmd.Flags().StringVar(&projectId, "id", "", "Project ID")
	cmd.Flags().StringVar(&name, "name", "", "New name of the project")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func init() {
	editCmd.AddCommand(NewEditProjectCmd())
}
