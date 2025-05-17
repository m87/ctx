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

			startDT, err := time.ParseInLocation(time.DateTime, strings.TrimSpace(args[2]), time.Local)
			util.Checkm(err, "Unable to parse start datetime")
			endDT, err := time.ParseInLocation(time.DateTime, strings.TrimSpace(args[3]), time.Local)
			util.Checkm(err, "Unable to parse end datetime")

			mgr.EditContextInterval(id, intervalIndex, ctx_model.LocalTime{Time: startDT}, ctx_model.LocalTime{Time: endDT})
		} else {
			for index, interval := range ctx.Intervals {
				fmt.Printf("[%d] %s - %s\n", index, interval.Start.Format(time.RFC3339Nano), interval.End.Format(time.RFC3339Nano))
			}
		}

	},
}

func init() {
	editContextCmd.AddCommand(editContextIntervalCmd)
}
