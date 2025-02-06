/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/m87/ctx/events"
	"github.com/spf13/cobra"
)

// listEventsCmd represents the listEvents command
var listEventsCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		eventsRegistry := events.Load()

		for _, v := range eventsRegistry.Events {
			fmt.Printf("[%s] %s\n", v.DateTime.Local().Format(time.DateTime), v.Data["description"])
		}
	},
}

func init() {
	eventsCmd.AddCommand(listEventsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listEventsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listEventsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
