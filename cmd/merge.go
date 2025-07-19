package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func NewMergeCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "merge <from> <to>",
		Aliases: []string{"m", "combine"},
		Short:   "Merge two contexts",
		Long:    "Merge the context with id <from> into the context with id <to>.",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fromId, err := flags.ResolveCustomContextId(cmd, "from")
			util.Check(err)
			toId, err := flags.ResolveCustomContextId(cmd, "to")
			util.Check(err)

			util.Check(manager.MergeContext(fromId, toId))
		},
	}
}

func init() {
	cmd := NewMergeCmd(bootstrap.CreateManager())
	flags.AddCustomContextFlag(cmd, "from", "f", "Source context to merge from")
	flags.AddCustomContextFlag(cmd, "to", "t", "Target context to merge into")
	rootCmd.AddCommand(cmd)
}
