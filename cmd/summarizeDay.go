package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func roundDuration(d time.Duration, unit string) time.Duration {
	switch unit {
	case "nanosecond":
		return d.Round(time.Nanosecond)
	case "microsecond":
		return d.Round(time.Microsecond)
	case "millisecond":
		return d.Round(time.Millisecond)
	case "second":
		return d.Round(time.Second)
	case "minute":
		return d.Round(time.Minute)
	case "hour":
		return d.Round(time.Hour)
	default:
		return d.Round(time.Nanosecond)
	}
}

var summarizeDayCmd = &cobra.Command{
	Use:     "day",
	Aliases: []string{"d", "day"},
	Short:   "Summarize day",
	Run: func(cmd *cobra.Command, args []string) {
		roundUnit, _ := cmd.Flags().GetString("round")

		date := time.Now().Local()
		if len(args) > 0 {
			rawDate := strings.TrimSpace(args[0])

			if rawDate != "" {
				date, _ = time.ParseInLocation(time.DateOnly, rawDate, time.Local)
			}
		}
		mgr := ctx.CreateManager()

		durations := map[string]time.Duration{}
		overallDuration := time.Duration(0)

		mgr.ContextStore.Read(func(s *ctx_model.State) error {
			for ctxId, _ := range s.Contexts {
				d, err := mgr.GetIntervalDurationsByDate(s, ctxId, ctx_model.LocalTime{Time: date})
				util.Checkm(err, "Unable to get interval durations for context "+ctxId)
				durations[ctxId] = roundDuration(d, roundUnit)
			}
			return nil
		})

		sortedIds := make([]string, 0, len(durations))
		for k := range durations {
			sortedIds = append(sortedIds, k)
		}
		sort.Strings(sortedIds)

		for _, c := range sortedIds {
			d := durations[c]
			ctx, _ := mgr.Ctx(c)
			if d > 0 {
				fmt.Printf("- %s: %s\n", ctx.Description, d)
				overallDuration += d
				if f, _ := cmd.Flags().GetBool("verbose"); f {
					mgr.ContextStore.Read(func(s *ctx_model.State) error {
						for _, interval := range mgr.GetIntervalsByDate(s, c, ctx_model.LocalTime{Time: date}) {
							fmt.Printf("\t- %s - %s\n", interval.Start.Format(time.RFC3339Nano), interval.End.Format(time.RFC3339Nano))
						}
						return nil
					})
				}
			}
		}

		fmt.Printf("Overall: %s\n", overallDuration)
	},
}

func init() {
	summarizeCmd.AddCommand(summarizeDayCmd)
	summarizeDayCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	summarizeDayCmd.Flags().StringP("round", "r", "nanosecond", "Round to the nearest nanosecond, microsecond, millisecond, second, minute, hour")
}
