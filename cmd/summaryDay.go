package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewSummaryDayCmd(manager *core.ContextManager) *cobra.Command {
	var dayRaw string

	cmd := &cobra.Command{
		Use:   "day",
		Short: "Show day summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			day, err := parseDay(dayRaw)
			if err != nil {
				return err
			}
			dayStr := day.Format("2006-01-02")

			if resolveRemoteAddr() != "" {
				fmt.Printf("Fetching summary for %s from remote...\n", dayStr)
				stats, err := remoteSummaryDay(dayStr)
				if err != nil {
					return err
				}

				if !Verbose {
					minimal := map[string]any{"date": stats.Date, "contextStats": stats.ContextStats}
					return printOutput(cmd, minimal, func() string {
						if stats == nil || len(stats.ContextStats) == 0 {
							return "No summary data found"
						}
						lines := []string{fmt.Sprintf("Summary for %s:", stats.Date)}
						for _, stat := range stats.ContextStats {
							lines = append(lines, fmt.Sprintf("- ContextID: %s, Duration: %s, Intervals: %d, Share: %.2f%%", stat.ContextId, time.Duration(stat.Duration).String(), stat.IntervalCount, stat.Percentage))
						}
						return strings.Join(lines, "\n")
					}, nil)
				}

				textRenderer := func() string {
					if stats == nil || len(stats.ContextStats) == 0 {
						return "No summary data found"
					}
					lines := []string{fmt.Sprintf("Summary for %s:", stats.Date)}
					for _, ctx := range stats.Contexts {
						lines = append(lines, fmt.Sprintf("Context ID: %s, Name: %s, Status: %s", ctx.Id, ctx.Name, ctx.Status))
						ivs := stats.Intervals[ctx.Id]
						if len(ivs) == 0 {
							lines = append(lines, "  (no intervals)")
						} else {
							for _, iv := range ivs {
								endStr := "(ongoing)"
								if !iv.End.IsZero {
									endStr = iv.End.Time.In(iv.End.Time.Location()).Format(time.RFC3339)
								}
								lines = append(lines, fmt.Sprintf("  - ID: %s, Start: %s, End: %s, Status: %s", iv.Id, iv.Start.Time.In(iv.Start.Time.Location()).Format(time.RFC3339), endStr, iv.Status))
							}
						}
					}
					return strings.Join(lines, "\n")
				}

				return printOutput(cmd, stats, textRenderer, nil)
			}

			intervals, err := manager.IntervalRepository.ListByDay(day)
			if err != nil {
				return err
			}
			now := manager.TimeProvider.Now().Time.UTC()
			rangesByContext := map[string][]core.TimeRange{}
			countByContext := map[string]int{}
			for _, interval := range intervals {
				rng, ok := core.ClipIntervalRangeToDay(interval, day, now)
				if !ok {
					continue
				}
				rangesByContext[interval.ContextId] = append(rangesByContext[interval.ContextId], rng)
				countByContext[interval.ContextId]++
			}

			stats := make([]map[string]any, 0, len(rangesByContext))
			for contextID, ranges := range rangesByContext {
				duration := core.SumMergedRangesDuration(ranges)
				stats = append(stats, map[string]any{
					"contextId": contextID,
					"duration":  duration.String(),
					"intervals": countByContext[contextID],
				})
			}
			sort.Slice(stats, func(i, j int) bool {
				return fmt.Sprintf("%v", stats[i]["contextId"]) < fmt.Sprintf("%v", stats[j]["contextId"])
			})

			if !Verbose {
				result := map[string]any{"date": dayStr, "contextStats": stats}
				return printOutput(cmd, result, func() string {
					if len(stats) == 0 {
						return "No summary data found"
					}
					lines := []string{fmt.Sprintf("Summary for %s:", dayStr)}
					for _, item := range stats {
						lines = append(lines, fmt.Sprintf("- ContextID: %v, Duration: %v, Intervals: %v", item["contextId"], item["duration"], item["intervals"]))
					}
					return strings.Join(lines, "\n")
				}, nil)
			}

			intervalsByContext := map[string][]*core.Interval{}
			for _, iv := range intervals {
				intervalsByContext[iv.ContextId] = append(intervalsByContext[iv.ContextId], iv)
			}

			var totalDuration time.Duration
			ctxStats := make([]*DayContextStats, 0, len(stats))
			for _, item := range stats {
				durationStr := item["duration"].(string)
				d, _ := time.ParseDuration(durationStr)
				totalDuration += d
			}
			for _, item := range stats {
				cid := fmt.Sprintf("%v", item["contextId"])
				durationStr := item["duration"].(string)
				d, _ := time.ParseDuration(durationStr)
				pct := 0.0
				if totalDuration > 0 {
					pct = (float64(d) / float64(totalDuration)) * 100.0
				}
				ctxStats = append(ctxStats, &DayContextStats{ContextId: cid, Duration: int64(d), Percentage: pct, IntervalCount: item["intervals"].(int)})
			}

			contexts := make([]*core.Context, 0, len(ctxStats))
			intervalsMap := map[string][]*core.Interval{}
			for _, cs := range ctxStats {
				cObj, _ := manager.ContextRepository.GetById(cs.ContextId)
				if cObj == nil {
					cObj = &core.Context{Id: cs.ContextId}
				}
				contexts = append(contexts, cObj)
				intervalsMap[cs.ContextId] = intervalsByContext[cs.ContextId]
			}

			verboseResult := &DayStats{
				Date:         dayStr,
				ContextStats: ctxStats,
				Contexts:     contexts,
				Intervals:    intervalsMap,
			}

			textRenderer := func() string {
				if len(verboseResult.ContextStats) == 0 {
					return "No summary data found"
				}
				lines := []string{fmt.Sprintf("Summary for %s:", verboseResult.Date)}
				for _, ctx := range verboseResult.Contexts {
					lines = append(lines, fmt.Sprintf("Context ID: %s, Name: %s, Status: %s", ctx.Id, ctx.Name, ctx.Status))
					ivs := verboseResult.Intervals[ctx.Id]
					if len(ivs) == 0 {
						lines = append(lines, "  (no intervals)")
					} else {
						for _, iv := range ivs {
							endStr := "(ongoing)"
							if !iv.End.IsZero {
								endStr = iv.End.Time.In(iv.End.Time.Location()).Format(time.RFC3339)
							}
							lines = append(lines, fmt.Sprintf("  - ID: %s, Start: %s, End: %s, Status: %s", iv.Id, iv.Start.Time.In(iv.Start.Time.Location()).Format(time.RFC3339), endStr, iv.Status))
						}
					}
				}
				return strings.Join(lines, "\n")
			}

			return printOutput(cmd, verboseResult, textRenderer, nil)
		},
	}

	cmd.Flags().StringVar(&dayRaw, "day", "", "Day in YYYY-MM-DD, default today")
	return cmd
}

func init() {
	summaryCmd.AddCommand(NewSummaryDayCmd(bootstrap.CreateManager()))
}
