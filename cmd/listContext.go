package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewListContextCmd(manager *core.ContextManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "List all contexts",
		Run: func(cmd *cobra.Command, args []string) {

			var contexts []*core.Context
			var err error
			if RemoteAddr != "" {
				contexts, err = remoteListContexts(cmd)
			} else {
				contexts, err = manager.ContextRepository.List()
			}

			if err != nil {
				cmd.PrintErrln("Error listing contexts:", err)
				return
			}
			if len(contexts) == 0 {
				cmd.Println("No contexts found")
				return
			}
			for _, ctx := range contexts {
				cmd.Printf("- ID: %s, Name: %s\n", ctx.Id, ctx.Name)
			}
		},
	}
	return cmd
}

func init() {
	listCmd.AddCommand(NewListContextCmd(bootstrap.CreateManager()))
}
