package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newFreeCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "free",
		Aliases: []string{"f"},
		Short:   "Stop current context",
		Long: `Stop current context. For example:
		ctx free
		`,
		Run: func(cmd *cobra.Command, args []string) {
			util.Check(manager.WithSession(func(session core.Session) error { return session.Free() }))
		},
	}
}

func init() {
	cmd := newFreeCmd(bootstrap.CreateManager())
	rootCmd.AddCommand(cmd)
}
