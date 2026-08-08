package cmd

import (
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)


func NewListProjectCmd() *cobra.Command {

	var (
		workspaceId string
	)

	cmd := &cobra.Command{
		Use:   "project",
		Short: "List all projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := bootstrap.CreateManager()
			if err != nil {
				return err
			}

			var projects []*core.Project
			if resolveRemoteAddr() != "" {
				projects, err = remoteClient().ListProjects(workspaceId)
			} else {
				projects, err = manager.ProjectRepository.List(workspaceId)
			}
			if err != nil {
				return err
			}


			textRenderer := func() string {
				if len(projects) == 0 {
					return "No projects found"
				}
				lines := make([]string, 0, len(projects))
				for _, project := range projects {
					lines = append(lines, "- ID: "+project.Id+", Name: "+project.Name+", Parent ID: "+project.ParentId+", Workspace ID: "+project.WorkspaceId)
				}
				return "Projects:\n" + strings.Join(lines, "\n")
			}

			return printOutput(cmd, projects, textRenderer, nil)
		},
	}

	cmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "Workspace ID of the projects to list")
	cmd.MarkFlagRequired("workspace")

	return cmd
}

func init() {
	listCmd.AddCommand(NewListProjectCmd())
}
