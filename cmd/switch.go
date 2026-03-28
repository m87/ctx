package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewSwitchCmd(manager *core.ContextManager) *cobra.Command {
	var (
		id   string
		name string
	)

	cmd := &cobra.Command{
		Use:   "switch",
		Short: "Switch active context",
		RunE: func(cmd *cobra.Command, args []string) error {
			contextID := strings.TrimSpace(id)
			contextName := strings.TrimSpace(name)
			if contextID == "" && contextName == "" {
				return fmt.Errorf("provide --id or --name")
			}

			if resolveRemoteAddr() != "" {
				if err := remoteSwitchContext(contextID, contextName); err != nil {
					return err
				}
				return printOutput(cmd, map[string]string{"id": contextID, "name": contextName, "status": "switched"}, func() string {
					return "Context switched successfully"
				}, nil)
			}

			context := &core.Context{Id: contextID, Name: contextName}
			if contextID == "" && contextName != "" {
				contexts, err := manager.ContextRepository.List()
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
	return cmd
}

func init() {
	rootCmd.AddCommand(NewSwitchCmd(bootstrap.CreateManager()))
}
