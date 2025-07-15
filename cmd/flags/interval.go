package flags

import (
	"errors"

	"github.com/spf13/cobra"
)

func ResolveIntervalId(cmd *cobra.Command) (string, error) {
	flags := cmd.Flags()

	id, err := flags.GetString("interval")
	if err != nil {
		return "", err
	}

	if !flags.Changed("interval") {
		id = ""
	}

	if id != "" {
		return id, nil
	}

	return "", errors.New("--interval must be provided")
}

func AddIntervalFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("interval", "i", "", "Interval id")
}
