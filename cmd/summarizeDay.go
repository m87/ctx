package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/m87/ctx/core"
	localstorage "github.com/m87/ctx/storage/local"
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

var summarizeDayCmd = &cobra.Command{
	Use:     "day",
	Aliases: []string{"d", "day"},
	Short:   "Summarize day",
	Run: func(cmd *cobra.Command, args []string) {
		roundUnit, _ := cmd.Flags().GetString("round")
		loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
		if err != nil {
			loc = time.UTC
		}
		mgr := localstorage.CreateManager()
		date := mgr.TimeProvider.Now().Time.In(loc)
		if len(args) > 0 {
			rawDate := strings.TrimSpace(args[0])

			if rawDate != "" {
				date, _ = time.ParseInLocation(time.DateOnly, rawDate, loc)
			}
		}

		durations := map[string]time.Duration{}
		overallDuration := time.Duration(0)

		mgr.ContextStore.Read(func(s *core.State) error {
			for ctxId, _ := range s.Contexts {
				d, err := mgr.GetIntervalDurationsByDate(s, ctxId, ctxtime.ZonedTime{Time: date, Timezone: loc.String()})
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
			if d > 0 {
				overallDuration += d
			}
		}

		if f, _ := cmd.Flags().GetBool("json"); f {

			outputContexts := []core.Context{}
			summary := DaySummary{}
			mgr.ContextStore.Read(func(s *core.State) error {
				for _, c := range sortedIds {
					d := durations[c]
					ctx, _ := mgr.Ctx(c)

					output := core.Context{
						Id:          c,
						Description: ctx.Description,
						Duration:    roundDuration(d, roundUnit),
					}

					if d > 0 {
						for _, interval := range mgr.GetIntervalsByDate(s, c, ctxtime.ZonedTime{Time: date, Timezone: loc.String()}) {
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
				return nil
			})

		} else {
			for _, c := range sortedIds {
				d := durations[c]
				ctx, _ := mgr.Ctx(c)
				if d > 0 {
					fmt.Printf("- %s: %s\n", ctx.Description, d)
					if f, _ := cmd.Flags().GetBool("verbose"); f {
						mgr.ContextStore.Read(func(s *core.State) error {
							intervals := mgr.GetIntervalsByDate(s, c, ctxtime.ZonedTime{Time: date, Timezone: loc.String()})
							for _, interval := range ctx.Intervals {
								if containsInterval(intervals, interval.Id) {
									fmt.Printf("\t[%s] %s - %s\n", interval.Id, interval.Start.Time.Format(time.DateTime), interval.End.Time.Format(time.DateTime))
								}
							}
							return nil
						})
					}
				}
			}

			fmt.Printf("Overall: %s\n", overallDuration)
		}
	},
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
	summarizeCmd.AddCommand(summarizeDayCmd)
	summarizeDayCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	summarizeDayCmd.Flags().StringP("round", "r", "nanosecond", "Round to the nearest nanosecond, microsecond, millisecond, second, minute, hour")
	summarizeDayCmd.Flags().BoolP("json", "j", false, "Json output")
}
