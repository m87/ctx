package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewCreateIntervalCmd(manager *core.ContextManager) *cobra.Command {
	var (
		contextID string
		startRaw  string
		endRaw    string
		status    string
	)

	cmd := &cobra.Command{
		Use:   "interval",
		Short: "Create a new interval",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(contextID) == "" {
				return fmt.Errorf("context-id is required")
			}

			var start core.ZonedTime
			if strings.TrimSpace(startRaw) == "" {
				start = core.NewTimer().Now()
			} else {
				parsed, err := parseDateTime(startRaw)
				if err != nil {
					return err
				}
				start = parsed
			}

			interval := &core.Interval{
				ContextId: strings.TrimSpace(contextID),
				Start:     start,
				Status:    strings.TrimSpace(status),
			}

			if interval.Status == "" {
				if strings.TrimSpace(endRaw) == "" {
					interval.Status = "active"
				} else {
					interval.Status = "completed"
				}
			}

			if strings.TrimSpace(endRaw) != "" {
				end, err := parseDateTime(endRaw)
				if err != nil {
					return err
				}
				interval.End = end
				if interval.End.Time.Before(interval.Start.Time) {
					return fmt.Errorf("end must be after start")
				}
				interval.Duration = interval.End.Time.Sub(interval.Start.Time)
			} else {
				interval.End = core.ZonedTime{Time: time.Time{}, Timezone: interval.Start.Timezone, IsZero: true}
			}

			if resolveRemoteAddr() != "" {
				if err := remoteCreateInterval(interval); err != nil {
					return err
				}
			} else {
				id, err := manager.IntervalRepository.Save(interval)
				if err != nil {
					return err
				}
				interval.Id = id
			}

			return printOutput(cmd, interval, func() string {
				return "Interval created successfully"
			}, nil)
		},
	}

	cmd.Flags().StringVar(&contextID, "context-id", "", "Context ID")
	cmd.Flags().StringVar(&startRaw, "start", "", "Start datetime in RFC3339")
	cmd.Flags().StringVar(&endRaw, "end", "", "End datetime in RFC3339")
	cmd.Flags().StringVar(&status, "status", "", "Interval status")
	_ = cmd.MarkFlagRequired("context-id")

	return cmd
}

func init() {
	createCmd.AddCommand(NewCreateIntervalCmd(bootstrap.CreateManager()))
}
