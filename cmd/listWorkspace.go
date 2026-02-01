package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
	"github.com/spf13/cobra"
)

func NewListWorkspaceCmd(manager *core.ContextManager) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "List all workspaces",
		Run: func(cmd *cobra.Command, args []string) {
			manager.Execute(func(repository *nod.Repository) error {
				nodes, err := repository.Query().TypeEquals(core.WorkspaceType).List()
				if err != nil {
					return err
				}

				for _, node := range nodes {
					ws := node.(*core.Workspace)
					cmd.Println("Workspace ID:", ws.Id, "Name:", ws.Name)
				}
				return nil
			})
		},
	}

	return cmd
}

func init() {
	listCmd.AddCommand(NewListWorkspaceCmd(bootstrap.CreateManager()))
}
