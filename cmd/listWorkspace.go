package cmd

import (
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewListWorkspaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "List workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := bootstrap.CreateManager()

			var workspaces []*core.Workspace
			var err error

			if resolveRemoteAddr() != "" {
				workspaces, err = remoteListWorkspaces()
			} else {
				workspaces, err = manager.WorkspaceRepository.List()
			}
			if err != nil {
				return err
			}

			return printOutput(cmd, workspaces, func() string {
				if len(workspaces) == 0 {
					return "No workspaces found"
				}
				lines := make([]string, 0, len(workspaces))
				for _, workspace := range workspaces {
					if workspace == nil {
						continue
					}
					lines = append(lines, "- ID: "+workspace.Id+", Name: "+workspace.Name)
				}
				if len(lines) == 0 {
					return "No workspaces found"
				}
				return strings.Join(lines, "\n")
			}, nil)
		},
	}
	return cmd
}

func init() {
	listCmd.AddCommand(NewListWorkspaceCmd())
}
