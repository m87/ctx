package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/spf13/cobra"
)

var summarizeDayCmd = &cobra.Command{
	Use:     "day",
	Aliases: []string{"d", "day"},
	Short:   "Summarize day",
	Run: func(cmd *cobra.Command, args []string) {
		date := time.Now().Local()
		if len(args) > 0 {
			rawDate := strings.TrimSpace(args[0])

			if rawDate != "" {
				date, _ = time.Parse(time.DateOnly, rawDate)
			}
		}
		mgr := ctx.CreateManager()

		durations := map[string]time.Duration{}
		overallDuration := time.Duration(0)

		mgr.ContextStore.Read(func(s *ctx_model.State) error {
			for ctxId, _ := range s.Contexts {
				durations[ctxId], _ = mgr.GetIntervalDurationsByDate(s, ctxId, date)
			}
			return nil
		})

		for c, d := range durations {
			ctx, _ := mgr.Ctx(c)
			if d > 0 {
				fmt.Printf("- %s: %s\n", ctx.Description, d)
			}
		}

		fmt.Printf("Overall: %s\n", overallDuration)
	},
}

func init() {
	summarizeCmd.AddCommand(summarizeDayCmd)
}
