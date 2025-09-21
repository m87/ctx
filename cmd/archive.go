package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func NewArchiveCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "archive",
		Aliases: []string{"archive", "a"},
		Short:   "Archive contexts",
	
		Run: func(cmd *cobra.Command, args []string) {
			util.Check(manager.WithSession(func(session core.Session) error { return session.Archive() }))
		},
	}

}

func init() {
	cmd := NewArchiveCmd(bootstrap.CreateManager())

	admCmd.AddCommand(cmd)
}
