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
		ctxId string
	)

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "d", "rm"},
		Short:   "Delete context",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cid, err := flags.ResolveContextId(args, ctxId)
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error { return session.Delete(cid.Id) }))
		},
	}

	flags.AddContextIdFlags(cmd, &ctxId)
	return cmd
}

var deleteCmd = newDeleteContextCmd(bootstrap.CreateManager())

func init() {
	rootCmd.AddCommand(deleteCmd)
}
