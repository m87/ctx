package cmd

import (
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var labelCmd = &cobra.Command{
	Use:     "label",
	Aliases: []string{"l", "lbl"},
	Run: func(cmd *cobra.Command, args []string) {
		context := strings.TrimSpace(args[0])
		contextId, err := util.Id(context, false)
		util.Checkm(err, "Unable to process id "+context)

		delete, _ := cmd.Flags().GetBool("delete")
		label := strings.TrimSpace(args[1])

		mgr := bootstrap.CreateManager()

		if delete {
			util.Check(
				mgr.WithSession(func(session core.Session) error {
					return session.DeleteLabelContext(contextId, label)
				}),
			)
		} else {
			util.Check(
				mgr.WithSession(func(session core.Session) error {
					return session.LabelContext(contextId, label)
				}),
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(labelCmd)
	labelCmd.Flags().BoolP("delete", "d", false, "Delete label from context")
}
