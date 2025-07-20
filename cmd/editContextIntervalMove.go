package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func NewEditContextIntervalMoveCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "move",
		Aliases: []string{"mv"},
		Short:   "mv interval",
		Run: func(cmd *cobra.Command, args []string) {
			descriptionSrc := strings.TrimSpace(args[0])
			descriptionTarget := strings.TrimSpace(args[1])
			idSrc, err := util.Id(descriptionSrc, false)
			util.Checkm(err, "Unable to process id src "+descriptionSrc)

			idTarget, err := util.Id(descriptionTarget, false)
			util.Checkm(err, "Unable to process id target "+descriptionTarget)

			manager.WithSession(func(session core.Session) error {

				ctxSrc, err := session.GetCtx(idSrc)

				if err != nil {
					panic("Context not found: " + idSrc)
				}
				_, err = session.GetCtx(idTarget)
				if err != nil {
					panic("Context not found: " + idTarget)
				}

				if len(args) > 2 {
					var err error
					intervalId := args[1]
					util.Checkm(err, "Unable to parse id")

					manager.MoveIntervalById(idSrc, idTarget, intervalId)
				} else {
					for _, interval := range ctxSrc.Intervals {
						fmt.Printf("[%s] %s - %s\n", interval.Id, interval.Start.Time.Format(time.RFC3339), interval.End.Time.Format(time.RFC3339))
					}
				}

				return nil
			})

		},
	}

}

func init() {
	cmd := NewEditContextIntervalMoveCmd(bootstrap.CreateManager())
	editContextIntervalCmd.AddCommand(cmd)
}
