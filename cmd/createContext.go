package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewCreateContextCmd(manager *core.ContextManager) *cobra.Command {
	var (
		name string
	)
	createContextCmd := &cobra.Command{
		Use:   "context",
		Short: "Create a new context",
		RunE: func(cmd *cobra.Command, args []string) error {
			context := &core.Context{
				Name: strings.TrimSpace(name),
			}

			if context.Name == "" {
				return fmt.Errorf("name is required")
			}

			if resolveRemoteAddr() != "" {
				if err := remoteCreateContext(context); err != nil {
					return err
				}
			} else {
				id, err := manager.ContextRepository.Save(context)
				if err != nil {
					return err
				}
				context.Id = id
			}

			return printOutput(cmd, context, func() string {
				return "Context created successfully"
			}, nil)
		},
	}
	createContextCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the context")
	_ = createContextCmd.MarkFlagRequired("name")
	return createContextCmd
}

func init() {
	createCmd.AddCommand(NewCreateContextCmd(bootstrap.CreateManager()))
}
