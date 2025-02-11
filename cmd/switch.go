/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strings"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/events"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if strings.TrimSpace(args[0]) == "" {
			return
		}

		eventsRegistry := events.Load()
		state := ctx.Load()

		isDescription, _ := cmd.Flags().GetBool("description")
		createIfNotFound, _ := cmd.Flags().GetBool("create")

		if isDescription {
			id = util.GenerateId(id)
		}

		err := ctx.Switch(id, &state, &eventsRegistry)

		if isDescription && createIfNotFound && err != nil {
			state.Contexts[id] = ctx.Context{
				Id:          id,
				Description: strings.TrimSpace(args[0]),
				State:       ctx.ACTIVE,
				Intervals:   []ctx.Interval{},
			}

			ctx.Switch(id, &state, &eventsRegistry)
		}

		log.Print(state)
		ctx.Save(&state)
		events.Save(&eventsRegistry)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().BoolP("create", "c", false, "create if not exists")
	switchCmd.Flags().BoolP("description", "d", false, "stop by description")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// switchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// switchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
