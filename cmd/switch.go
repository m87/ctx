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
			if len(args) == 0 {
				panic("Please provide a description or id")
			}
			id, description, isRawId, err := flags.ResolveContextId(args[0], ctxId, ctxDescription)

			if isRawId {
				util.Checkm(err, "Unable to process context id "+id)
			} else {
				util.Checkm(err, "Unable to process context "+description)
			}

			util.Check(manager.WithSession(func(session core.Session) error {
				if isRawId {
					return session.Switch(id)
				} else {
					return session.CreateIfNotExistsAndSwitch(id, description)
				}

			}))

		},
	}

	flags.AddContextIdFlags(cmd)
	return cmd
}

func init() {
	switchCmd := newSwitchCmd(bootstrap.CreateManager())
	rootCmd.AddCommand(switchCmd)
}
