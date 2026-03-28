package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewDeleteContextCmd(manager *core.ContextManager) *cobra.Command {
	var contextId string
	deleteContextCmd := &cobra.Command{
		Use:   "context",
		Short: "Delete a context by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := strings.TrimSpace(contextId)
			if id == "" {
				return fmt.Errorf("id is required")
			}

			if resolveRemoteAddr() != "" {
				if err := remoteDeleteContext(id); err != nil {
					return err
				}
			} else {
				if err := manager.ContextRepository.Delete(id); err != nil {
					return err
				}
			}

			return printOutput(cmd, map[string]string{"id": id, "status": "deleted"}, func() string {
				return "Context deleted successfully"
			}, nil)
		},
	}
	deleteContextCmd.Flags().StringVarP(&contextId, "id", "i", "", "ID of the context to delete")
	_ = deleteContextCmd.MarkFlagRequired("id")
	return deleteContextCmd
}

func init() {
	deleteCmd.AddCommand(NewDeleteContextCmd(bootstrap.CreateManager()))
}
