package cmd

import (
	"strings"

	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:     "archive",
	Aliases: []string{"arc", "a"},
	Short:   "Archive all contexts",
	Run: func(cmd *cobra.Command, args []string) {
		cm := core.CreateManager()
		if all, _ := cmd.Flags().GetBool("all"); all {
			util.Check(cm.ArchiveAll())
		} else {
			description := strings.TrimSpace(args[0])
			byId, _ := cmd.Flags().GetBool("id")
			id, err := util.Id(description, byId)
			util.Checkm(err, "Unable to process id "+description)
			util.Check(cm.Archive(id))
		}
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().BoolP("all", "a", false, "Archive all active contexts")
}
