/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	localstorage "github.com/m87/ctx/storage/local"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var editContextIntervalSplitCmd = &cobra.Command{
	Use:     "split",
	Aliases: []string{"s"},
	Short:   "Split interval",
	Run: func(cmd *cobra.Command, args []string) {
		description := strings.TrimSpace(args[0])
		id, err := util.Id(description, false)
		util.Checkm(err, "Unable to process id "+description)

		mgr := localstorage.CreateManager()
		intervalIndex := -1

		ctx, err := mgr.Ctx(id)

		if err != nil {
			panic("Context not found: " + id)
		}

		if len(args) > 1 {
			var err error
			intervalIndex, err = strconv.Atoi(args[1])
			util.Checkm(err, "Unable to parse interval index")

			if intervalIndex < 0 || intervalIndex > len(ctx.Intervals)-1 {
				panic("interval index out of range")
			}

			interval := ctx.Intervals[intervalIndex]
			loc, _ := time.LoadLocation(interval.Start.Timezone)
			split, err := time.ParseInLocation(time.DateTime, strings.TrimSpace(args[2]), loc)
			util.Checkm(err, "Unable to parse split datetime")

			if split.Before(interval.Start.Time) {
				panic("split time is before interval start time")
			}

			if split.After(interval.End.Time) {
				panic("split time is after interval end time")
			}

			mgr.SplitContextIntervalByIndex(id, intervalIndex, split)
		} else {
			for index, interval := range ctx.Intervals {
				fmt.Printf("[%d] %s - %s\n", index, interval.Start.Time.Format(time.RFC3339), interval.End.Time.Format(time.RFC3339))
			}
		}

	},
}

func init() {
	editContextIntervalCmd.AddCommand(editContextIntervalSplitCmd)
}
