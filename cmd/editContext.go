package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewEditContextCmd(manager *core.ContextManager) *cobra.Command {
	var (
		id          string
		name        string
		description string
		status      string
	)

	cmd := &cobra.Command{
		Use:   "context",
		Short: "Edit an existing context",
		RunE: func(cmd *cobra.Command, args []string) error {
			contextID := strings.TrimSpace(id)
			if contextID == "" {
				return fmt.Errorf("id is required")
			}

			context := &core.Context{Id: contextID}
			if resolveRemoteAddr() == "" {
				existing, err := manager.ContextRepository.GetById(contextID)
				if err != nil {
					return err
				}
				if existing == nil {
					return fmt.Errorf("context not found")
				}
				context = existing
			}

			if strings.TrimSpace(name) != "" {
				context.Name = strings.TrimSpace(name)
			}
			if strings.TrimSpace(description) != "" {
				context.Description = strings.TrimSpace(description)
			}
			if strings.TrimSpace(status) != "" {
				context.Status = strings.TrimSpace(status)
			}

			if resolveRemoteAddr() != "" {
				if err := remoteUpdateContext(context); err != nil {
					return err
				}
			} else {
				if _, err := manager.ContextRepository.Save(context); err != nil {
					return err
				}
			}

			return printOutput(cmd, context, func() string {
				return "Context updated successfully"
			}, nil)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID of the context to edit")
	cmd.Flags().StringVarP(&name, "name", "n", "", "New context name")
	cmd.Flags().StringVar(&description, "description", "", "New context description")
	cmd.Flags().StringVar(&status, "status", "", "New context status")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func init() {
	editCmd.AddCommand(NewEditContextCmd(bootstrap.CreateManager()))
}
