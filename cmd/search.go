package cmd

import (
	"fmt"
	"time"

	localstorage "github.com/m87/ctx/storage/local"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"S", "search"},
	Short:   "Search for a context by description with regex",
	Run: func(cmd *cobra.Command, args []string) {
		regex := args[0]
		mgr := localstorage.CreateManager()
		ctxs, err := mgr.Search(regex)

		if err != nil {
			util.Checkm(err, "Unable to search for context "+regex)
		}

		if f, _ := cmd.Flags().GetBool("verbose"); f {
			for _, c := range ctxs {
				println(c.Id + ": " + c.Description)
				for i, interval := range c.Intervals {
					fmt.Printf("\t[%d] %s - %s\n", i, interval.Start.Time.Format(time.DateTime), interval.End.Time.Format(time.DateTime))
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
