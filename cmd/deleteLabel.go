package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newDeleteLabelCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId string
		label string
	)

	cmd := &cobra.Command{
		Use:   "label",
		Short: "Delete context label",
		Run: func(cmd *cobra.Command, args []string) {
			cid, params, err := flags.ResolveCidWithParams(args, ctxId, flags.ParamSpec{Default: label, Name: "label"})
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error {
				return session.DeleteLabelContext(cid.Id, params["label"])
			}))
		},
	}

	flags.AddContextIdFlags(cmd, &ctxId)
	cmd.Flags().StringVar(&label, "label", "", "label content")
	return cmd
}

var deleteLabelCmd = newDeleteLabelCmd(bootstrap.CreateManager())

func init() {
	deleteCmd.AddCommand(deleteLabelCmd)
}
