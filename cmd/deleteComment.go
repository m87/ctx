package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newDeleteCommentContextCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId     string
		commentId string
	)

	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Delete context comment",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cid, params, err := flags.ResolveCidWithParams(args, ctxId, flags.ParamSpec{Default: commentId, Name: "comment-id"})
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error {
				return session.DeleteContextComment(cid.Id, params["comment-id"])
			}))
		},
	}

	flags.AddContextIdFlags(cmd, &ctxId)
	cmd.Flags().StringVar(&commentId, "comment-id", "", "comment id to delete")
	return cmd
}

var deleteCommentCmd = newDeleteCommentContextCmd(bootstrap.CreateManager())

func init() {
	deleteCmd.AddCommand(deleteCommentCmd)
}
