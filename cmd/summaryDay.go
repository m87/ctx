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
				stats, err := remoteSummaryDay(dayStr)
				if err != nil {
					return err
				}
				return printOutput(cmd, stats, func() string {
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
		},
	}

	cmd.Flags().StringVar(&dayRaw, "day", "", "Day in YYYY-MM-DD, default today")
	return cmd
}

func init() {
	summaryCmd.AddCommand(NewSummaryDayCmd(bootstrap.CreateManager()))
}
