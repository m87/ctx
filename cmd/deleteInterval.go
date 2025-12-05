package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newDeleteIntervalCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId      string
		intervalId string
	)

	cmd := &cobra.Command{
		Use:     "interval",
		Aliases: []string{"int", "i"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cid, intervalId, err := flags.ResolveCidWithResourceId(args, ctxId, intervalId, "interval id")
			util.Check(err)

			util.Check(manager.WithSession(func(session core.Session) error {
				return session.DeleteInterval(cid.Id, intervalId)
			}))
		},
	}

	flags.AddContextIdFlags(cmd, &ctxId)
	flags.AddIntervalFlag(cmd, &intervalId)
	return cmd
}

func init() {
	cmd := newDeleteIntervalCmd(bootstrap.CreateManager())
	deleteCmd.AddCommand(cmd)
}
