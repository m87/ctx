package cmd

import (
	"github.com/google/uuid"
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
	"strings"
)

func newCommentContextCmd(manager *core.ContextManager) *cobra.Command {
	var (
		ctxId          string
		ctxDescription string
	)

	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Comment context",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				panic("Please provide a description or id")
			}
			id, _, _, err := flags.ResolveContextId(args[0], ctxId, ctxDescription)
			comment := strings.TrimSpace(args[1])
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error {
				return session.SaveContextComment(id, core.Comment{Id: uuid.NewString(), Content: comment})
			}))
		},
	}

	flags.AddContextIdFlags(cmd, &ctxId, &ctxDescription)
	return cmd
}

var commentCmd = newCommentContextCmd(bootstrap.CreateManager())

func init() {
	rootCmd.AddCommand(commentCmd)
}
