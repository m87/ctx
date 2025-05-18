package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var editContextIntervalCmd = &cobra.Command{
	Use:   "interval",
	Short: "Edit selected interval",
	Run: func(cmd *cobra.Command, args []string) {
		description := strings.TrimSpace(args[0])
		id, err := util.Id(description, false)
		util.Checkm(err, "Unable to process id "+description)

		mgr := ctx.CreateManager()
		intervalIndex := -1

		ctx, err := mgr.Ctx(id)

		if err != nil {
			panic("Context not found: " + id)
		}

		if len(args) > 1 {
			var err error
			intervalIndex, err = strconv.Atoi(args[1])
			util.Checkm(err, "Unable to parse id")

			if intervalIndex < 0 || intervalIndex > len(ctx.Intervals)-1 {
				panic("interval index out of range")
			}
			loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
			if err != nil {
				loc = time.UTC
			}
			startDT, err := time.ParseInLocation(time.DateTime, strings.TrimSpace(args[2]), loc)
			util.Checkm(err, "Unable to parse start datetime")
			endDT, err := time.ParseInLocation(time.DateTime, strings.TrimSpace(args[3]), loc)
			util.Checkm(err, "Unable to parse end datetime")

			mgr.EditContextInterval(id, intervalIndex, ctx_model.ZonedTime{Time: startDT, Timezone: loc.String()}, ctx_model.ZonedTime{Time: endDT, Timezone: loc.String()})
		} else {
			for index, interval := range ctx.Intervals {
				fmt.Printf("[%d] %s - %s\n", index, interval.Start.Time.Format(time.RFC3339Nano), interval.End.Time.Format(time.RFC3339Nano))
			}
		}

	},
}

func init() {
	editContextCmd.AddCommand(editContextIntervalCmd)
}
