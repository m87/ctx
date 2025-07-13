package cmd

import (
	"strings"

	localstorage "github.com/m87/ctx/storage/local"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "c"},
	Short:   "Create new context",
	Long: `Create new context from given description. Passed description is used to generate contextId with sha256. For example:
	ctx create new-context
	ctx create "new context with spaces"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		description := strings.TrimSpace(args[0])
		id, err := util.Id(description, false)
		util.Checkm(err, "Unable to process id "+description)

		util.Check(localstorage.CreateManager().CreateContext(id, description))
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
