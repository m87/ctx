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
		Use:   "interval",
		Short: "Delete interval from context",
		Long: `Delete interval from context. For example:
	ctx delete interval "my-context" "interval-id"
	ctx delete interval "my-context" --interval-id "interval-id"
	ctx delete interval --ctx-id "my-context-id" --interval-id "interval-id"
	`,
		Run: func(cmd *cobra.Command, args []string) {
			cid, params, err := flags.ResolveCidWithParams(args, ctxId, flags.ParamSpec{Default: intervalId, Name: "interval-id"})
			util.Check(err)

			util.Check(manager.WithSession(func(session core.Session) error {
				return session.DeleteInterval(cid.Id, params["interval-id"])
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
