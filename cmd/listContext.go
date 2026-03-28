package cmd

import (
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewListContextCmd(manager *core.ContextManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "List all contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			var contexts []*core.Context
			var err error
			if resolveRemoteAddr() != "" {
				contexts, err = remoteListContexts()
			} else {
				contexts, err = manager.ContextRepository.List()
			}
			if err != nil {
				return err
			}

			return printOutput(cmd, contexts, func() string {
				if len(contexts) == 0 {
					return "No contexts found"
				}
				lines := make([]string, 0, len(contexts))
				for _, context := range contexts {
					lines = append(lines, "- ID: "+context.Id+", Name: "+context.Name)
				}
				return strings.Join(lines, "\n")
			}, nil)
		},
	}
	return cmd
}

func init() {
	listCmd.AddCommand(NewListContextCmd(bootstrap.CreateManager()))
}
