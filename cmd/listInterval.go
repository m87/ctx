package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewListIntervalCmd(manager *core.ContextManager) *cobra.Command {
	var dayRaw string

	cmd := &cobra.Command{
		Use:   "interval",
		Short: "List intervals for a day",
		RunE: func(cmd *cobra.Command, args []string) error {
			day, err := parseDay(dayRaw)
			if err != nil {
				return err
			}
			dayStr := day.Format("2006-01-02")

			var report *DayReport
			if resolveRemoteAddr() != "" {
				report, err = remoteListIntervalsByDay(dayStr)
				if err != nil {
					return err
				}
			} else {
				intervals, err := manager.IntervalRepository.ListByDay(day)
				if err != nil {
					return err
				}
				report = &DayReport{Intervals: intervals}
			}

			return printOutput(cmd, report, func() string {
				if report == nil || len(report.Intervals) == 0 {
					return "No intervals found"
				}
				lines := make([]string, 0, len(report.Intervals)+1)
				lines = append(lines, fmt.Sprintf("Intervals for %s:", dayStr))
				for _, interval := range report.Intervals {
					if interval == nil {
						continue
					}
					lines = append(lines, fmt.Sprintf("- ID: %s, ContextID: %s, Start: %s, End: %s, Status: %s", interval.Id, interval.ContextId, interval.Start.Time.Format("2006-01-02T15:04:05Z07:00"), interval.End.Time.Format("2006-01-02T15:04:05Z07:00"), interval.Status))
				}
				return strings.Join(lines, "\n")
			}, nil)
		},
	}

	cmd.Flags().StringVar(&dayRaw, "day", "", "Day in YYYY-MM-DD, default today")
	return cmd
}

func init() {
	listCmd.AddCommand(NewListIntervalCmd(bootstrap.CreateManager()))
}
