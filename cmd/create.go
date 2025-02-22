package cmd

import (
	"strings"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new context",
	Long: `Create new context from given description. Passed description is used to generate contextId with sha256. For example:
	ctx create new-context
	ctx create "new context with spaces"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		description := strings.TrimSpace(args[0])
		id, err := util.Id(description, false)
		util.Check(err, "Unable to process id "+description)

		util.ApplyPatch(func(state *ctx_model.State) {
			ctx.Create(state, id, description)
		})
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
