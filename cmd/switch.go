package cmd

import (
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func NewSwitchCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
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

			description := strings.TrimSpace(args[0])
			byId, _ := cmd.Flags().GetBool("id")
			id, err := util.Id(description, byId)
			util.Checkm(err, "Unable to process id "+description)

			util.Check(manager.WithSession(func(session core.Session) error {
				if byId {
					return session.Switch(id)
				} else {
					return session.CreateIfNotExistsAndSwitch(id, description)
				}

			}))

		},
	}

}

func init() {
	switchCmd := NewSwitchCmd(bootstrap.CreateManager())
	switchCmd.Flags().BoolP("id", "i", false, "stop by description")
	rootCmd.AddCommand(switchCmd)
}
