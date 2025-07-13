package cmd

import (
	"strings"

	localstorage "github.com/m87/ctx/storage/local"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "d", "rm"},
	Short:   "Delete context",
	Run: func(cmd *cobra.Command, args []string) {
		description := strings.TrimSpace(args[0])
		id, err := util.Id(description, false)
		util.Checkm(err, "Unable to process id "+description)

		util.Check(localstorage.CreateManager().Delete(id))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
