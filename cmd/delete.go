package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newDeleteContextCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId          string
		ctxDescription string
	)

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "d", "rm"},
		Short:   "Delete context",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				panic("Please provide a description or id")
			}
			id, _, _, err := flags.ResolveContextId(args[0], ctxId, ctxDescription)
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error { return session.Delete(id) }))
		},
	}

	flags.AddContextIdFlags(cmd, &ctxId, &ctxDescription)
	return cmd
}

var deleteCmd = newDeleteContextCmd(bootstrap.CreateManager())

func init() {
	rootCmd.AddCommand(deleteCmd)
}
