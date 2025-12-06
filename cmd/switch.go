package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newSwitchCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId          string
		ctxDescription string
	)

	cmd := &cobra.Command{
		Use:     "switch",
		Aliases: []string{"sw", "s"},
		Short:   "Switch context",
		Long: `Switch context:
	- switch description, created if not exists
	- switch -i id"`,
		Run: func(cmd *cobra.Command, args []string) {
			cid, err := flags.ResolveContextId(args, ctxId)
			util.Check(err)

			util.Check(manager.WithSession(func(session core.Session) error {
				if ctxDescription == "" {
					return session.Switch(cid.Id)
				} else {
					return session.CreateIfNotExistsAndSwitch(cid.Id, ctxDescription)
				}

			}))

		},
	}

	flags.AddContextIdFlags(cmd, &ctxId)
	return cmd
}

func init() {
	switchCmd := newSwitchCmd(bootstrap.CreateManager())
	rootCmd.AddCommand(switchCmd)
}
