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
				date, _ = time.ParseInLocation(time.DateOnly, rawDate, time.Local)
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
				overallDuration += d
				if f, _ := cmd.Flags().GetBool("verbose"); f {
					mgr.ContextStore.Read(func(s *ctx_model.State) error {
						for _, interval := range mgr.GetIntervalsByDate(s, c, date) {
							fmt.Printf("\t- %s - %s\n", interval.Start.Format(time.RFC3339Nano), interval.End.Format(time.RFC3339Nano))
						}
						return nil
					})
				}
			}
		}

		fmt.Printf("Overall: %s\n", overallDuration)
	},
}

func init() {
	summarizeCmd.AddCommand(summarizeDayCmd)
	summarizeDayCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
