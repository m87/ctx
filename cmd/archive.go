package cmd

import (
	"github.com/m87/ctx/archive"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/events_model"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive active contexts",
	Long: `Archive single context: ctx archive test
	Archvie all active contexts: 
		ctx archive --all
		ctx archive -a

	`,
	Run: func(cmd *cobra.Command, args []string) {
		isRaw, _ := cmd.Flags().GetBool("raw")
		archiveAll, _ := cmd.Flags().GetBool("all")
		id, err := util.Id(args[0], isRaw)
		util.Check(err, "Unable to process id "+args[0])

		util.ApplyPatch(func(state *ctx_model.State) {
			util.ApplyEventsPatch(func(eventsRegistry *events_model.EventRegistry) {
				if archiveAll {
					archive.ArchiveAll(state, eventsRegistry)
				} else {
					err := archive.Archive(id, state, eventsRegistry)
					util.Warn(err, id+" - context is active")
				}
			})
		})
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().BoolP("all", "a", false, "Archive all active contexts")
}
