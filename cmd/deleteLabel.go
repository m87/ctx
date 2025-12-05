package cmd

import (
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newDeleteLabelCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId string
	)

	cmd := &cobra.Command{
		Use:   "label",
		Short: "Delete context label",
		Args:  cobra.RangeArgs(2, 2),
		Run: func(cmd *cobra.Command, args []string) {
			cid, err := flags.ResolveContextId(args, ctxId)
			label := strings.TrimSpace(args[1])
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error {
				return session.DeleteLabelContext(cid.Id, label)
			}))
		},
	}

	flags.AddContextIdFlags(cmd, &ctxId)
	return cmd
}

var deleteLabelCmd = newDeleteLabelCmd(bootstrap.CreateManager())

func init() {
	deleteCmd.AddCommand(deleteLabelCmd)
}
