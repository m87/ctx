package cmd

import (
	"strconv"

	"github.com/m87/ctx/core"
	localstorage "github.com/m87/ctx/storage/local"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var deleteIntervalCmd = &cobra.Command{
	Use:     "interval",
	Aliases: []string{"int", "i"},
	Run: func(cmd *cobra.Command, args []string) {
		description := args[0]
		id, err := util.Id(description, false)
		util.Checkm(err, "Unable to process id "+description)
		index, err := strconv.Atoi(args[1])
		util.Checkm(err, "Unable to process index "+args[1])
		if index < 0 {
			util.Checkm(err, "Index must be greater than or equal to 0")
		} else {
			util.Check(localstorage.CreateManager().WithSession(func(session core.Session) error {
				return session.DeleteIntervalByIndex(id, index)
			}))
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteIntervalCmd)
}
