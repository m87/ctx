package cmd

import (
	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/spf13/cobra"
)

var listEventsCmd = &cobra.Command{
	Use:   "list",
	Short: "List events",
	Long:  `List events.`,
	Run: func(cmd *cobra.Command, args []string) {
		date, _ := cmd.Flags().GetString("date")
		types, _ := cmd.Flags().GetStringArray("types")

		filter := ctx_model.EventsFilter{
			Date:  date,
			Types: types,
		}

		mgr := ctx.CreateManager()
		if j, _ := cmd.Flags().GetBool("json"); j {
			mgr.ListEventsJson(filter)
		} else if f, _ := cmd.Flags().GetBool("full"); f {
			mgr.ListEventsFull(filter)
		} else {
			mgr.ListEvents(filter)
		}
	},
}

func init() {
	eventsCmd.AddCommand(listEventsCmd)
	listEventsCmd.Flags().StringP("date", "d", "", "show for date")
	listEventsCmd.Flags().StringArrayP("types", "t", []string{}, "show for date")
	listEventsCmd.Flags().BoolP("full", "f", false, "show full list")
	listEventsCmd.Flags().BoolP("json", "j", false, "show list as json")
}
