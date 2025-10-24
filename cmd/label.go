package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newLabelContextCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "label",
		Aliases: []string{"l", "lbl"},
		Run: func(cmd *cobra.Command, args []string) {
			contextId, err := flags.ResolveContextIdLegacy(cmd)
			util.Check(err)

			delFlag, err := flags.ResolveDeleteFlag(cmd)
			util.Check(err)

			label, err := flags.ResolveLabelFlag(cmd)
			util.Check(err)

			util.Check(
				manager.WithSession(func(session core.Session) error {
					if delFlag {
						return session.DeleteLabelContext(contextId, label)
					} else {
						return session.LabelContext(contextId, label)
					}
				}),
			)
		},
	}
}

func init() {
	cmd := newLabelContextCmd(bootstrap.CreateManager())
	flags.AddContxtFlag(cmd)
	flags.AddDeleteFlag(cmd, "Delete label from context")
	flags.AddLabelFlag(cmd)

	rootCmd.AddCommand(cmd)
}
