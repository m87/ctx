package cmd

import (
	"github.com/google/uuid"
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newAddCommentCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId   string
		comment string
	)
	cmd := &cobra.Command{
		Use:     "comment",
		Aliases: []string{"c"},
		Short:   "Add comment to context",
		Long: `Add comment to context. For example:
	ctx add comment "my-context-id" "This is my comment"
	`,
		Run: func(cmd *cobra.Command, args []string) {
			cid, err := flags.ResolveContextId(args, ctxId)
			util.Check(err)
			comment, err := flags.ResolveArgument(args, flags.ConditionalIndexProvider(ctxId != "")(1), comment, "comment")
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error {
				return session.SaveContextComment(cid.Id, core.Comment{Content: comment, Id: uuid.NewString()})
			}))
		},
	}

	flags.AddContextIdFlags(cmd, &ctxId)
	cmd.Flags().StringVar(&comment, "comment", "", "comment content")
	return cmd
}

var addCommentCmd = newAddCommentCmd(bootstrap.CreateManager())

func init() {
	addCmd.AddCommand(addCommentCmd)
}
