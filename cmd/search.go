package cmd

import (
	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"S", "search"},
	Short:   "Search for a context by description with regex",
	Run: func(cmd *cobra.Command, args []string) {
		regex := args[0]
		mgr := ctx.CreateManager()
		ctxs, err := mgr.Search(regex)

		if err != nil {
			util.Checkm(err, "Unable to search for context "+regex)
		}

		if f, _ := cmd.Flags().GetBool("verbose"); f {
			for _, c := range ctxs {
				println(c.Id + ": " + c.Description)
				for _, i := range c.Intervals {
					println("  " + i.Start.Time.String() + " - " + i.End.Time.String())
				}
			}
		} else {
			for _, c := range ctxs {
				println(c.Id + ": " + c.Description)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
