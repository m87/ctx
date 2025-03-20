package cmd

import (
	"github.com/spf13/cobra"
)

var listEventsCmd = &cobra.Command{
	Use:   "listEvents",
	Short: "List events",
	Long:  `List events.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	eventsCmd.AddCommand(listEventsCmd)
}
