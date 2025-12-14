package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newRenameContextCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId       string
		description string
	)

	cmd := &cobra.Command{
		Use:     "rename",
		Aliases: []string{"r"},
		Short:   "Rename context",
		Long: `Rename context. For example:
	ctx rename "my-context" --description "New Description"
	ctx rename --ctx-id "my-context-id" --description "New Description"
	`,
		Run: func(cmd *cobra.Command, args []string) {
			ctxId, params, err := flags.ResolveCidWithParams(args, ctxId, flags.ParamSpec{Default: description, Name: "description"})
			util.Check(err)
			manager.WithSession(func(session core.Session) error {
				return session.RenameContext(ctxId.Id, util.GenerateId(params["description"]), params["description"])
			})

		},
	}

	flags.AddContextIdFlags(cmd, &ctxId)
	cmd.Flags().StringVar(&description, "description", "", "New context description")
	return cmd
}

func init() {
	rootCmd.AddCommand(newRenameContextCmd(bootstrap.CreateManager()))
}
