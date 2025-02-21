/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/events"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func switchContext(state *ctx_model.State, input string, isRawId bool, createIfNotFound bool) {
	id, err := util.Id(input, isRawId)
	util.Check(err, "Unable to process id "+id)

	eventsRegistry := events.Load()

	err = ctx.Switch(id, state, &eventsRegistry)

	if createIfNotFound && err != nil {
		state.Contexts[id] = ctx_model.Context{
			Id:          id,
			Description: strings.TrimSpace(input),
			State:       ctx_model.ACTIVE,
			Intervals:   []ctx_model.Interval{},
		}

		ctx.Switch(id, state, &eventsRegistry)
	}

	events.Save(&eventsRegistry)
}

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
		util.ApplyPatch(func(state *ctx_model.State) {
			isRaw, _ := cmd.Flags().GetBool("raw")
			createIfNotFound, _ := cmd.Flags().GetBool("create")
			switchContext(state, args[0], isRaw, createIfNotFound)
		})

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
