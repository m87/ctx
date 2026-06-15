package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewSwitchCmd() *cobra.Command {
	var (
		id          string
		name        string
		workspaceID string
	)

	cmd := &cobra.Command{
		Use:   "switch",
		Short: "Switch active context",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := bootstrap.CreateManager()

			contextID := strings.TrimSpace(id)
			contextName := strings.TrimSpace(name)
			selectedWorkspaceID := strings.TrimSpace(workspaceID)
			if contextID == "" && contextName == "" {
				return fmt.Errorf("provide --id or --name")
			}

			if resolveRemoteAddr() != "" {
				if contextID == "" && selectedWorkspaceID == "" {
					return fmt.Errorf("workspace is required when switching by name")
				}
				if err := remoteSwitchContext(contextID, contextName, selectedWorkspaceID); err != nil {
					return err
				}
				return printOutput(cmd, map[string]string{"id": contextID, "name": contextName, "status": "switched"}, func() string {
					return "Context switched successfully"
				}, nil)
			}

			context := &core.Context{Id: contextID, Name: contextName, WorkspaceId: selectedWorkspaceID}
			if contextID == "" && contextName != "" {
				if selectedWorkspaceID == "" {
					return fmt.Errorf("workspace is required when switching by name")
				}
				contexts, err := manager.ContextRepository.ListByWorkspace(selectedWorkspaceID)
				if err != nil {
					return err
				}
				for _, candidate := range contexts {
					if strings.EqualFold(strings.TrimSpace(candidate.Name), contextName) {
						context = candidate
						break
					}
				}
			}

			if err := manager.SwitchContext(context); err != nil {
				return err
			}

			return printOutput(cmd, context, func() string {
				return "Context switched successfully"
			}, nil)
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Context ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Context name")
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func init() {
	rootCmd.AddCommand(NewSwitchCmd())
}
