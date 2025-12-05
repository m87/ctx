package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newAddLabelCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId string
		label string
	)
	cmd := &cobra.Command{
		Use:     "label",
		Aliases: []string{"l"},
		Short:   "Add label to context",
		Long: `Add label to context. For example:
	ctx add label "my-context-id" "production"
	`,
		Run: func(cmd *cobra.Command, args []string) {
			cid, label, err := flags.ResolveCidWithResourceId(args, ctxId, label, "label")
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error {
				return session.LabelContext(cid.Id, label)
			}))
		},
	}

	flags.AddContextIdFlags(cmd, &ctxId)
	cmd.Flags().StringVar(&label, "label", "", "label content")
	return cmd
}

var addLabelCmd = newAddLabelCmd(bootstrap.CreateManager())

func init() {
	addCmd.AddCommand(addLabelCmd)
}
