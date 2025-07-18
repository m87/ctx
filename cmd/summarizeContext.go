package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m87/ctx/core"
	localstorage "github.com/m87/ctx/storage/local"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func getContextCreationTimeFromEvents(ctxId string) (string, error) {
	mgr := localstorage.CreateManager()
	var creationTime string
	err := mgr.EventsStore.Read(func(er *core.EventRegistry) error {
		for _, event := range er.Events {
			if event.Type == core.CREATE_CTX && event.CtxId == ctxId {
				creationTime = event.DateTime.Time.Format(time.RFC3339)
				return nil
			}
		}
		return nil
	})
	return creationTime, err
}

var summarizeContextCmd = &cobra.Command{
	Use:     "context",
	Aliases: []string{"ctx", "c"},
	Short:   "Summarize context",
	Run: func(cmd *cobra.Command, args []string) {
		id := ""
		if len(args) > 0 {
			description := strings.TrimSpace(args[0])
			selectedId, err := util.Id(description, false)
			util.Checkm(err, "Unable to process id "+description)
			id = selectedId
		}

		verbose, _ := cmd.Flags().GetBool("verbose")

		mgr := localstorage.CreateManager()
		if id == "" {
			mgr.ContextStore.Read(func(s *core.State) error {
				if s.CurrentId != "" {
					id = s.CurrentId
				} else {
					panic("No context selected")
				}
				return nil
			})
		}
		ctx, err := mgr.Ctx(id)
		util.Checkm(err, "Unable to find context "+id)

		fmt.Printf("Context: %s\n", ctx.Description)
		if verbose {
			fmt.Printf("Id: %s\n", ctx.Id)
		}
		ceationTime, _ := getContextCreationTimeFromEvents(ctx.Id)
		fmt.Printf("Created: %s\n", ceationTime)
		fmt.Printf("Duration: %s\n", ctx.Duration)

	},
}

func init() {
	summarizeCmd.AddCommand(summarizeContextCmd)
	summarizeContextCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
