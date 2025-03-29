package cmd

import (
	"github.com/m87/ctx/ctx"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List contexts",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := ctx.CreateManager()
		if j, _ := cmd.Flags().GetBool("json"); j {
			mgr.ListJson()
		} else if f, _ := cmd.Flags().GetBool("verbose"); f {
			mgr.ListFull()
		} else {
			mgr.List()
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("verbose", "v", false, "show full list")
	listCmd.Flags().BoolP("json", "j", false, "show list as json")
}
