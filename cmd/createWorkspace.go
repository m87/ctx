package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
	"github.com/spf13/cobra"
)

func NewCreateWorkspaceCmd(manager *core.ContextManager) *cobra.Command {
	var (
		name string
	)

	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Create a new workspace",
		Long:  `Create a new workspace in the context management system.`,
		Run: func(cmd *cobra.Command, args []string) {

			manager.Execute(func(repository *nod.Repository) error {
				workspace := core.NewWorkspace(name)
				return repository.Save(workspace)
			})
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the workspace")
	return cmd
}

func init() {
	createCmd.AddCommand(NewCreateWorkspaceCmd(bootstrap.CreateManager()))
}
