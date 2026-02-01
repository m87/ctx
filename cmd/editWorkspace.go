package cmd

import (
	"fmt"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
	"github.com/spf13/cobra"
)

func NewEditWorkspaceCmd(manager *core.ContextManager) *cobra.Command {
	var (
		id   string
		name string
	)

	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Edit an existing workspace",
		Long:  `Edit an existing workspace by specifying its ID and new name.`,
		Run: func(cmd *cobra.Command, args []string) {
			manager.Execute(func(repository *nod.Repository) error {
				node, err := repository.Query().TypeEquals(core.WorkspaceType).NodeId(id).First()
				if err != nil {
					return err
				}

				if node == nil {
					return fmt.Errorf("workspace with ID %s not found", id)
				}

				ws := node.(*core.Workspace)
				ws.Name = name

				if err := repository.Save(ws); err != nil {
					return err
				}

				return nil
			})
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID of the workspace to edit")
	cmd.Flags().StringVarP(&name, "name", "n", "", "New name for the workspace")
	return cmd
}

func init() {
	editCmd.AddCommand(NewEditWorkspaceCmd(bootstrap.CreateManager()))
}
