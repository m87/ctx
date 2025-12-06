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
			cid, params, err := flags.ResolveCidWithParams(args, ctxId, flags.ParamSpec{Default: commentId, Name: "comment-id"}, flags.ParamSpec{Default: comment, Name: "comment"})
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error {
				err := session.DeleteContextComment(cid.Id, params["comment-id"])
				if err != nil {
					return err
				}
				return session.SaveContextComment(cid.Id, core.Comment{Id: params["comment-id"], Content: params["comment"]})
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
