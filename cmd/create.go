package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newCreateContextCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Aliases: []string{"new", "c"},
		Short:   "Create new context",
		Long: `Create new context from given description. Passed description is used to generate contextId with sha256. For example:
	ctx create new-context
	ctx create "new context with spaces"
	`,
		Run: func(cmd *cobra.Command, args []string) {
			cid, err := flags.ResolveContextIdentifier(cmd, args)
			util.Check(err)

			util.Check(manager.WithSession(func(session core.Session) error { return session.CreateContext(cid.Id, cid.Description) }))
		},
	}

}

func init() {
	cmd := newCreateContextCmd(bootstrap.CreateManager())
	rootCmd.AddCommand(cmd)
}
