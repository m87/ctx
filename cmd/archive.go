package cmd

import (
	"github.com/m87/ctx/archive"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/ctx_store"
	"github.com/m87/ctx/events_model"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func archiveContext(input string, isRaw bool, archiveAll bool) {
	if archiveAll {
		st := ctx_store.Load()
		util.ApplyPatch(func(state *ctx_model.State) {
			for _, v := range st.Contexts {
				util.ApplyEventsPatch(func(eventsRegistry *events_model.EventRegistry) {
					archive.Archive(v.Id, state, eventsRegistry)
				})
			}

		})
	} else {
		util.ApplyPatch(func(state *ctx_model.State) {
			id, err := util.Id(input, isRaw)
			util.Check(err, "Unable to process id "+input)

			util.ApplyEventsPatch(func(eventsRegistry *events_model.EventRegistry) {
				archive.Archive(id, state, eventsRegistry)
			})

		})
	}
}

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		isRaw, _ := cmd.Flags().GetBool("raw")
		archiveAll, _ := cmd.Flags().GetBool("all")
		archiveContext(args[0], isRaw, archiveAll)
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().BoolP("all", "a", false, "Show full info")
}
