package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/spf13/cobra"
)

var summarizeDayCmd = &cobra.Command{
	Use:     "day",
	Aliases: []string{"d", "day"},
	Short:   "Summarize day",
	Run: func(cmd *cobra.Command, args []string) {
		date := strings.TrimSpace(args[0])

		filter := ctx_model.EventsFilter{
			Date: date,
		}

		mgr := ctx.CreateManager()
		events := mgr.FilterEvents(filter)

		durations := map[string]time.Duration{}
		overallDuration := time.Duration(0)

		for _, e := range events {
			if e.Type == ctx_model.END_INTERVAL {
				duration, _ := time.ParseDuration(e.Data["duration"])
				if _, ok := durations[e.CtxId]; ok {
					durations[e.CtxId] = durations[e.CtxId] + duration
				} else {
					durations[e.CtxId] = duration
				}
				overallDuration = overallDuration + duration
			}
		}

		for c, d := range durations {
			ctx, _ := mgr.Ctx(c)
			fmt.Printf("- %s: %s\n", ctx.Description, d)
		}

		fmt.Printf("Overall: %s\n", overallDuration)
	},
}

func init() {
	summarizeCmd.AddCommand(summarizeDayCmd)
}
