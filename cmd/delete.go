package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newDeleteContextCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "d", "rm"},
		Short:   "Delete context",
		Run: func(cmd *cobra.Command, args []string) {
			ctxId, err := flags.ResolveContextId(cmd)
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error { return session.Delete(ctxId) }))
		},
	}

}

var deleteCmd = newDeleteContextCmd(bootstrap.CreateManager())

func init() {
	flags.AddContxtFlag(deleteCmd)
	rootCmd.AddCommand(deleteCmd)
}
