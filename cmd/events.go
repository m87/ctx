package cmd

import (
	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:     "events",
	Aliases: []string{"event", "e"},
}

func init() {
	rootCmd.AddCommand(eventsCmd)
}
