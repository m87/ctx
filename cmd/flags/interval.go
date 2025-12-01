package flags

import (
	"errors"

	"github.com/spf13/cobra"
)

func ResolveIntervalId(id string) (string, error) {
	if id != "" {
		return id, nil
	}

	return "", errors.New("--interval must be provided")
}

func AddIntervalFlag(cmd *cobra.Command, intervalId *string) {
	cmd.Flags().StringVar(intervalId, "interval-id", "", "Interval id")
}
