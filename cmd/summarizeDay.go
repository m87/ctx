package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	ctxtime "github.com/m87/ctx/time"
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

type DaySummary struct {
	Contexts []core.Context `json:"contexts"`
	Duration time.Duration  `json:"duration"`
}

func newSummarizeDayCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "day",
		Aliases: []string{"d", "day"},
		Short:   "Summarize day",
		Run: func(cmd *cobra.Command, args []string) {
			roundUnit, _ := cmd.Flags().GetString("round")
			loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
			if err != nil {
				loc = time.UTC
			}

			manager.WithSession(func(session core.Session) error {
				date := session.TimeProvider.Now().Time.In(loc)
				if len(args) > 0 {
					rawDate := strings.TrimSpace(args[0])

					if rawDate != "" {
						date, _ = time.ParseInLocation(time.DateOnly, rawDate, loc)
					}
				}

				durations := map[string]time.Duration{}
				overallDuration := time.Duration(0)

				for ctxId, _ := range session.State.Contexts {
					d, err := session.GetIntervalDurationsByDate(ctxId, ctxtime.ZonedTime{Time: date, Timezone: loc.String()})
					util.Checkm(err, "Unable to get interval durations for context "+ctxId)
					durations[ctxId] = roundDuration(d, roundUnit)
				}

				sortedIds := make([]string, 0, len(durations))
				for k := range durations {
					sortedIds = append(sortedIds, k)
				}
				sort.Strings(sortedIds)

				for _, c := range sortedIds {
					d := durations[c]
					if d > 0 {
						overallDuration += d
					}
				}

				if f, _ := cmd.Flags().GetBool("json"); f {

					outputContexts := []core.Context{}
					summary := DaySummary{}
					for _, c := range sortedIds {
						d := durations[c]
						ctx, _ := session.GetCtx(c)

						output := core.Context{
							Id:          c,
							Description: ctx.Description,
							Duration:    roundDuration(d, roundUnit),
						}

						if d > 0 {
							for _, interval := range session.GetIntervalsByDate(c, ctxtime.ZonedTime{Time: date, Timezone: loc.String()}) {
								output.Intervals[interval.Id] = core.Interval{
									Start:    interval.Start,
									End:      interval.End,
									Duration: roundDuration(interval.End.Time.Sub(interval.Start.Time), roundUnit),
								}
							}
						}

						outputContexts = append(outputContexts, output)
					}
					summary.Contexts = outputContexts
					summary.Duration = roundDuration(overallDuration, roundUnit)
					j, _ := json.Marshal(summary)

					fmt.Printf("%s", string(j))

				} else {
					for _, c := range sortedIds {
						d := durations[c]
						ctx, _ := session.GetCtx(c)
						if d > 0 {
							fmt.Printf("- %s: %s\n", ctx.Description, d)
							if f, _ := cmd.Flags().GetBool("verbose"); f {
								intervals := session.GetIntervalsByDate(c, ctxtime.ZonedTime{Time: date, Timezone: loc.String()})
								for _, interval := range ctx.Intervals {
									if containsInterval(intervals, interval.Id) {
										fmt.Printf("\t[%s] %s - %s\n", interval.Id, interval.Start.Time.Format(time.DateTime), interval.End.Time.Format(time.DateTime))
									}
								}
							}
						}
					}

					fmt.Printf("Overall: %s\n", overallDuration)
				}

				return nil
			})
		},
	}

}

func containsInterval(intervals []core.Interval, id string) bool {
	for _, i := range intervals {
		if i.Id == id {
			return true
		}
	}
	return false
}

func init() {
	summarizeDayCmd := newSummarizeDayCmd(bootstrap.CreateManager())
	summarizeDayCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	summarizeDayCmd.Flags().StringP("round", "r", "nanosecond", "Round to the nearest nanosecond, microsecond, millisecond, second, minute, hour")
	summarizeDayCmd.Flags().BoolP("json", "j", false, "Json output")
	summarizeCmd.AddCommand(summarizeDayCmd)
}
