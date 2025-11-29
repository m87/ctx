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

func newCommentContextCmd(manager *core.ContextManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Comment context",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cid, err := flags.ResolveContextIdentifier(cmd, args)
			util.Check(err)
			comment := strings.TrimSpace(args[1])
			util.Check(manager.WithSession(func(session core.Session) error {
				return session.SaveContextComment(cid.Id, core.Comment{Id: uuid.NewString(), Content: comment})
			}))
		},
	}

	flags.AddContextIdFlags(cmd)
	return cmd
}

var commentCmd = newCommentContextCmd(bootstrap.CreateManager())

func init() {
	rootCmd.AddCommand(commentCmd)
}
