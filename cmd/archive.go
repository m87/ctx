package cmd

import (
	"log"

	"github.com/m87/ctx/archive"
	"github.com/m87/ctx/archive_store"
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

		util.ApplyPatch(func(state *ctx_model.State) {
			util.ApplyEventsPatch(func(eventsRegistry *events_model.EventRegistry) {
				if archiveAll {
					ArchiveAll(state, eventsRegistry)
				} else {
					id, err := util.Id(args[0], isRaw)
					util.Check(err, "Unable to process id "+args[0])
					entry := archive_store.LoadArchive(id)
					events, err := archive.Archive(id, state, eventsRegistry, &entry)
					archive_store.SaveArchive(&entry)
					archive_store.SaveEventsArchive(events)
					util.Warn(err, id+" - context is active")
				}
			})
		})
	},
}

func ArchiveAll(state *ctx_model.State, eventsRegistry *events_model.EventRegistry) {
	for _, v := range state.Contexts {
		entry := archive_store.LoadArchive(v.Id)
		events, err := archive.Archive(v.Id, state, eventsRegistry, &entry)
		if err != nil {
			log.Printf("Active context %s, skipping\n", v.Id)
		}
		archive_store.SaveArchive(&entry)
		archive_store.SaveEventsArchive(events)
	}
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().BoolP("all", "a", false, "Archive all active contexts")
}
