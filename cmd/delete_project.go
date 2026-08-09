package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/spf13/cobra"
)

func NewDeleteProjectCmd() *cobra.Command {
	var (
		projectId string
	)

	cmd := &cobra.Command{
		Use:   "project",
		Short: "Delete a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := bootstrap.CreateManager()
			if err != nil {
				return err
			}

			if resolveRemoteAddr() != "" {
				err := remoteClient().DeleteProject(projectId)
				if err != nil {
					return err
				}
			} else {
				err := manager.DeleteProject(projectId)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&projectId, "id", "p", "", "Project ID to delete")

	return cmd
}

func init() {
	deleteCmd.AddCommand(NewDeleteProjectCmd())
}
