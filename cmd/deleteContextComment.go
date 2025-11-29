package cmd

import (
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newDeleteCommentContextCmd(manager *core.ContextManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete context comment",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				panic("Please provide a description or id")
			}
			cid, err := flags.ResolveContextIdentifier(cmd, args)
			util.Check(err)
			commentId := strings.TrimSpace(args[1])
			util.Check(manager.WithSession(func(session core.Session) error {
				return session.DeleteContextComment(cid.Id, commentId)
			}))
		},
	}

	flags.AddContextIdentifierFlags(cmd)
	return cmd
}

var deleteCommentCmd = newDeleteCommentContextCmd(bootstrap.CreateManager())

func init() {
	commentCmd.AddCommand(deleteCommentCmd)
}
