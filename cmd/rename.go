package cmd

import (
	"strings"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"r"},
	Short:   "Rename context",
	Run: func(cmd *cobra.Command, args []string) {
		src := strings.TrimSpace(args[0])
		srcId, err := util.Id(src, false)
		util.Checkm(err, "Unable to process id "+src)

		target := strings.TrimSpace(args[1])
		targetId, err := util.Id(target, false)
		util.Checkm(err, "Unable to process id "+target)

		mgr := ctx.CreateManager()

		util.Check(mgr.RenameContext(srcId, targetId, target))

	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
