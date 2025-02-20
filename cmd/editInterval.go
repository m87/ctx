/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strconv"
	"time"

	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

// editIntervalCmd represents the editInterval command
var editIntervalCmd = &cobra.Command{
	Use:   "interval",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		index, _ := strconv.Atoi(args[1])
		value := args[2]
		changeEnd, _ := cmd.Flags().GetBool("end")

		util.ApplyPatch(func(state *ctx_model.State) {
			context := state.Contexts[id]
			interval := context.Intervals[index]

			if context.Id == state.CurrentId {
				log.Fatalln("Cannont modify active context")
			}

			if changeEnd {
				newTime, _ := time.ParseInLocation(time.DateTime, value, time.Local)
				newDuration := newTime.Sub(interval.Start)
				interval.End = newTime
				diff := interval.Duration - newDuration
				interval.Duration = newDuration
				context.Duration = context.Duration - diff
				context.Intervals[index] = interval
				state.Contexts[id] = context
			}
		})
	},
}

func init() {
	editCmd.AddCommand(editIntervalCmd)
	editIntervalCmd.Flags().BoolP("end", "e", false, "edit interval end")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editIntervalCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editIntervalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
