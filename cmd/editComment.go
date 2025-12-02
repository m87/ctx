package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newEditCommentCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId     string
		commentId string
		comment   string
	)
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Edit context comment",
		Args:  cobra.MaximumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			cid, err := flags.ResolveContextId(args, ctxId)
			util.Check(err)
			commentId, err = flags.ResolveArgument(args, 1, commentId, "comment id")
			util.Check(err)
			comment, err = flags.ResolveArgument(args, 2, comment, "comment")
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error {
				err := session.DeleteContextComment(cid.Id, commentId)
				if err != nil {
					return err
				}
				return session.SaveContextComment(cid.Id, core.Comment{Id: commentId, Content: comment})
			}))
		},
	}
	flags.AddContextIdFlags(cmd, &ctxId)
	cmd.Flags().StringVar(&comment, "comment", "", "comment content")
	cmd.Flags().StringVar(&commentId, "comment-id", "", "comment id")
	return cmd
}

var editCommentCmd = newEditCommentCmd(bootstrap.CreateManager())

func init() {
	editCmd.AddCommand(editCommentCmd)
}
