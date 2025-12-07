package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newAddCmd(manager *core.ContextManager) *cobra.Command {
	var (
		description string
	)

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c", "add"},
		Short:   "Create new context",
		Long: `Create new context from given description. Passed description is used to generate contextId with sha256. For example:
	ctx add new-context
	ctx create "new context with spaces"
	`,
		Run: func(cmd *cobra.Command, args []string) {
			params, err := flags.ResolveParams(args, flags.ParamSpec{Default: description, Name: "description"})
			util.Check(err)
			util.Check(manager.WithSession(func(session core.Session) error {
				return session.CreateContext(util.GenerateId(params["description"]), params["description"])
			}))
		},
	}

	cmd.Flags().StringVar(&description, "description", "", "context description")
	return cmd
}

var addCmd = newAddCmd(bootstrap.CreateManager())

func init() {
	rootCmd.AddCommand(addCmd)
}
