package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newSwitchCmd(manager *core.ContextManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "switch",
		Aliases: []string{"sw", "s"},
		Short:   "Switch context",
		Long: `Switch context:
	- switch description, created if not exists
	- switch -i id"`,
		Run: func(cmd *cobra.Command, args []string) {
			cid, err := flags.ResolveContextIdentifier(cmd, args)

		util.Check(manager.WithSession(func(session core.Session) error {
				if session.ValidateContextExists(cid) {
					return session.Switch(cid)
				} else {
					return session.CreateIfNotExistsAndSwitch(cid)
				}

			}))

		},
	}

	flags.AddContextIdentifierFlags(cmd)
	return cmd
}

func init() {
	switchCmd := newSwitchCmd(bootstrap.CreateManager())
	rootCmd.AddCommand(switchCmd)
}
