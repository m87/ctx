/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func NewEditContextIntervalSplitCmd(manager *core.ContextManager) *cobra.Command {

	return &cobra.Command{
		Use:     "split",
		Aliases: []string{"s"},
		Short:   "Split interval",
		Run: func(cmd *cobra.Command, args []string) {
			description := strings.TrimSpace(args[0])
			id, err := util.Id(description, false)
			util.Checkm(err, "Unable to process id "+description)

			ctx, err := manager.Ctx(id)

			if err != nil {
				panic("Context not found: " + id)
			}

			if len(args) > 1 {
				var err error
				intervalId := args[1]
				util.Checkm(err, "Unable to parse interval index")

				interval := ctx.Intervals[intervalId]
				loc, _ := time.LoadLocation(interval.Start.Timezone)
				split, err := time.ParseInLocation(time.DateTime, strings.TrimSpace(args[2]), loc)
				util.Checkm(err, "Unable to parse split datetime")

				if split.Before(interval.Start.Time) {
					panic("split time is before interval start time")
				}

				if split.After(interval.End.Time) {
					panic("split time is after interval end time")
				}

				manager.SplitContextIntervalById(id, intervalId, split)
			} else {
				for _, interval := range ctx.Intervals {
					fmt.Printf("[%s] %s - %s\n", interval.Id, interval.Start.Time.Format(time.RFC3339), interval.End.Time.Format(time.RFC3339))
				}
			}

		},
	}

}
func init() {
	cmd := NewEditContextIntervalSplitCmd(bootstrap.CreateManager())
	editContextIntervalCmd.AddCommand(cmd)
}
