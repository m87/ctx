package cmd

import (
	"strings"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:     "merge",
	Aliases: []string{"m", "combine"},
	Short:   "Merge two contexts",
	Run: func(cmd *cobra.Command, args []string) {
		fromDescription := strings.TrimSpace(args[0])
		fromId, err := util.Id(fromDescription, false)
		util.Checkm(err, "Unable to process id "+fromDescription)

		toDescription := strings.TrimSpace(args[1])
		toId, err := util.Id(toDescription, false)
		util.Checkm(err, "Unable to process id "+toDescription)

		mgr := ctx.CreateManager()

		err = mgr.MergeContext(fromId, toId)
		util.Check(err)
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}
