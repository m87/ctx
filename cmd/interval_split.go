package cmd

import (
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/spf13/cobra"
)


func paresLocalDateTime(timeStr string) (*time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return nil, err
	}
	return &parsedTime, nil
}

func NewIntervalSplitCmd() *cobra.Command {
	var (
		id string
		splitTime string
	)

	cmd := &cobra.Command{
		Use:   "split",
		Short: "Split an interval into two at the specified time",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := bootstrap.CreateManager()
			if err != nil {
				return err
			}


			parsedTime, err := paresLocalDateTime(splitTime)
			if err != nil {
				return err
			}

			err = manager.SplitInterval(id, *parsedTime)
			if err != nil {
				return err
			}

			cmd.Printf("Interval %q split at %s\n", id, parsedTime.Format(time.RFC3339))

			return nil
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID of the interval to split")
	cmd.Flags().StringVarP(&splitTime, "time", "t", "", "Time to split the interval (Local format: YYYY-MM-DD HH:MM:SS)")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagsOneRequired("time")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewIntervalSplitCmd())
}
