package cmd

import (
	"strings"

	"github.com/google/uuid"
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newCommentAddCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a comment to a context",
		Long: `Add a comment to a context identified by its id or description.
		For example:
		- comment add description "This is a comment"
		- comment add --ctx-id contextId "This is a comment"`,
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cid, err := flags.ResolveContextId(args, ctxId)
			comment := strings.TrimSpace(args[1])
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error {
				return session.SaveContextComment(cid.Id, core.Comment{Id: uuid.NewString(), Content: comment})
			}))
		},
	}

	flags.AddContextIdFlags(cmd, &ctxId)
	return cmd
}

var commentAddCmd = newCommentAddCmd(bootstrap.CreateManager())

func init() {
	commentCmd.AddCommand(commentAddCmd)
}
