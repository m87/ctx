package cmd

import (
	"fmt"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func NewSearchCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "search",
		Aliases: []string{"S", "search"},
		Short:   "Search for a context by description with regex",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				util.Checkm(fmt.Errorf("no regex provided"), "Usage: ctx search <regex>")
			}
			regex := args[0]

			verbose, err := flags.ResolveVerboseFlag(cmd)
			util.Check(err)

			ctxs, err := manager.Search(regex)

			if err != nil {
				util.Checkm(err, "Unable to search for context "+regex)
			}

			if verbose {
				for _, c := range ctxs {
					println(c.Id + ": " + c.Description)
					for _, interval := range c.Intervals {
						fmt.Printf("\t[%s] %s - %s\n", interval.Id, interval.Start.Time.Format(time.DateTime), interval.End.Time.Format(time.DateTime))
					}
				}
			} else {
				for _, c := range ctxs {
					println(c.Id + ": " + c.Description)
				}
			}

		},
	}
}

func init() {
	cmd := NewSearchCmd(bootstrap.CreateManager())
	flags.AddVerboseFlag(cmd)
	rootCmd.AddCommand(cmd)
}
