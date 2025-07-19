package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func NewDeleteIntervalCmd(manager *core.ContextManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "interval",
		Aliases: []string{"int", "i"},
		Run: func(cmd *cobra.Command, args []string) {
			contextId, err := flags.ResolveContextId(cmd)
			util.Check(err)
			id, err := flags.ResolveIntervalId(cmd)
			util.Check(err)

			util.Check(manager.WithSession(func(session core.Session) error {
				return session.DeleteInterval(contextId, id)
			}))
		},
	}
	return cmd
}

func init() {
	cmd := NewDeleteIntervalCmd(bootstrap.CreateManager())
	flags.AddContxtFlag(deleteCmd)
	flags.AddIntervalFlag(deleteCmd)

	deleteCmd.AddCommand(cmd)
}
