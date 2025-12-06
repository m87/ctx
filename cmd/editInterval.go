package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	ctxtime "github.com/m87/ctx/time"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newEditContextIntervalCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:   "interval",
		Short: "Edit selected interval",
		Run: func(cmd *cobra.Command, args []string) {
			description := strings.TrimSpace(args[0])
			id, err := util.Id(description, false)
			util.Checkm(err, "Unable to process id "+description)

			manager.WithSession(func(session core.Session) error {

				ctx, err := session.GetCtx(id)

				if err != nil {
					panic("Context not found: " + id)
				}

				if len(args) > 1 {
					var err error
					intervalId := args[1]
					util.Checkm(err, "Unable to parse id")

					loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
					if err != nil {
						loc = time.UTC
					}
					startDT, err := time.ParseInLocation(time.DateTime, strings.TrimSpace(args[2]), loc)
					util.Checkm(err, "Unable to parse start datetime")
					endDT, err := time.ParseInLocation(time.DateTime, strings.TrimSpace(args[3]), loc)
					util.Checkm(err, "Unable to parse end datetime")

					session.EditContextIntervalById(id, intervalId, ctxtime.ZonedTime{Time: startDT, Timezone: loc.String()}, ctxtime.ZonedTime{Time: endDT, Timezone: loc.String()})
				} else {
					for _, interval := range ctx.Intervals {
						fmt.Printf("[%s] %s - %s\n", interval.Id, interval.Start.Time.Format(time.RFC3339), interval.End.Time.Format(time.RFC3339))
					}
				}

				return nil
			})

		},
	}

}

var editIntervalCmd = newEditContextIntervalCmd(bootstrap.CreateManager())

func init() {
	editCmd.AddCommand(editIntervalCmd)
}
