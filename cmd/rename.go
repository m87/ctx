package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newRenameContextCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "rename",
		Aliases: []string{"r"},
		Short:   "Rename context",
		Run: func(cmd *cobra.Command, args []string) {
			srcId, err := flags.ResolveCustomContextId(cmd, "src-ctx")
			util.Check(err)

			targetId, err := flags.ResolveCustomContextId(cmd, "target-ctx")
			util.Check(err)
			
			target, err := cmd.Flags().GetString("target-ctx")
			util.Check(err)

			manager.WithSession(func(session core.Session) error {
				return session.RenameContext(srcId, targetId, target)
			})	

		},
	}
}

func init() {
	cmd := newRenameContextCmd(bootstrap.CreateManager())
	flags.AddCustomContextFlag(cmd, "src-ctx", "s", "Source context")
	flags.AddCustomContextFlag(cmd, "target-ctx", "t", "Target context")
	rootCmd.AddCommand(cmd)
}
