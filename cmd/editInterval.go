package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewEditIntervalCmd(manager *core.ContextManager) *cobra.Command {
	var (
		id        string
		contextID string
		startRaw  string
		endRaw    string
		status    string
	)

	cmd := &cobra.Command{
		Use:   "interval",
		Short: "Edit an existing interval",
		RunE: func(cmd *cobra.Command, args []string) error {
			intervalID := strings.TrimSpace(id)
			if intervalID == "" {
				return fmt.Errorf("id is required")
			}

			interval := &core.Interval{Id: intervalID}
			if resolveRemoteAddr() == "" {
				existing, err := manager.IntervalRepository.GetById(intervalID)
				if err != nil {
					return err
				}
				if existing == nil {
					return fmt.Errorf("interval not found")
				}
				interval = existing
			}

			if strings.TrimSpace(contextID) != "" {
				interval.ContextId = strings.TrimSpace(contextID)
			}
			if strings.TrimSpace(startRaw) != "" {
				start, err := parseDateTime(startRaw)
				if err != nil {
					return err
				}
				interval.Start = start
			}
			if strings.TrimSpace(endRaw) != "" {
				end, err := parseDateTime(endRaw)
				if err != nil {
					return err
				}
				interval.End = end
			}
			if strings.TrimSpace(status) != "" {
				interval.Status = strings.TrimSpace(status)
			}

			if !interval.End.Time.IsZero() {
				if interval.End.Time.Before(interval.Start.Time) {
					return fmt.Errorf("end must be after start")
				}
				interval.Duration = interval.End.Time.Sub(interval.Start.Time)
			} else {
				interval.End = core.ZonedTime{Time: time.Time{}, Timezone: interval.Start.Timezone, IsZero: true}
			}

			if resolveRemoteAddr() != "" {
				if err := remoteUpdateInterval(interval); err != nil {
					return err
				}
			} else {
				if _, err := manager.IntervalRepository.Save(interval); err != nil {
					return err
				}
			}

			return printOutput(cmd, interval, func() string {
				return "Interval updated successfully"
			}, nil)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID of the interval to edit")
	cmd.Flags().StringVar(&contextID, "context-id", "", "Context ID")
	cmd.Flags().StringVar(&startRaw, "start", "", "Start datetime in format 'YYYY-MM-DD HH:MM:SS' or RFC3339")
	cmd.Flags().StringVar(&endRaw, "end", "", "End datetime in format 'YYYY-MM-DD HH:MM:SS' or RFC3339")
	cmd.Flags().StringVar(&status, "status", "", "Interval status")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func init() {
	editCmd.AddCommand(NewEditIntervalCmd(bootstrap.CreateManager()))
}
