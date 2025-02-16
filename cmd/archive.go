/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/m87/ctx/archive"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/events"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

// archiveCmd represents the archive command
var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if f, _ := cmd.Flags().GetBool("all"); f {
			// st := ctx_store.Load()
			// util.ApplyPatch(func(state *ctx_model.State) {
			// 	for _, v := range st.Contexts {
			// 		id, err := util.Id(v.Id, cmd)
			// 		util.Check(err, "Unable to process id "+v.Id)

			// 		eventsRegistry := events.Load()
			// 		archive.Archive(id, state, &eventsRegistry)
			// 		events.Save(&eventsRegistry)
			// 	}

			// })
		} else {
			util.ApplyPatch(func(state *ctx_model.State) {
				id, err := util.Id(args[0], cmd)
				util.Check(err, "Unable to process id "+args[0])

				eventsRegistry := events.Load()
				archive.Archive(id, state, &eventsRegistry)
				events.Save(&eventsRegistry)

			})
		}

	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().BoolP("all", "a", false, "Show full info")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// archiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// archiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
