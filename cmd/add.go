package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newAddCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "add",
		Aliases: []string{"new", "c", "create"},
		Short:   "Create new context",
		Long: `Create new context from given description. Passed description is used to generate contextId with sha256. For example:
	ctx create new-context
	ctx create "new context with spaces"
	`,
		Run: func(cmd *cobra.Command, args []string) {
			description, err := flags.GetStringArg(args, 0, "description")
			util.Check(err)
			id, err := flags.ResolveArgumentAsContextId(args, 0, "description")
			util.Check(err)

			util.Check(manager.WithSession(func(session core.Session) error { return session.CreateContext(id, description) }))
		},
	}

}

var addCmd = newAddCmd(bootstrap.CreateManager())

func init() {
	rootCmd.AddCommand(addCmd)
}
