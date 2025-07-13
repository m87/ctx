package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	localstorage "github.com/m87/ctx/storage/local"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

var editContextIntervalMoveCmd = &cobra.Command{
	Use:     "move",
	Aliases: []string{"mv"},
	Short:   "mv interval",
	Run: func(cmd *cobra.Command, args []string) {
		descriptionSrc := strings.TrimSpace(args[0])
		descriptionTarget := strings.TrimSpace(args[1])
		intervalIndex := -1
		idSrc, err := util.Id(descriptionSrc, false)
		util.Checkm(err, "Unable to process id src "+descriptionSrc)

		idTarget, err := util.Id(descriptionTarget, false)
		util.Checkm(err, "Unable to process id target "+descriptionTarget)

		mgr := localstorage.CreateManager()

		ctxSrc, err := mgr.Ctx(idSrc)

		if err != nil {
			panic("Context not found: " + idSrc)
		}
		_, err = mgr.Ctx(idTarget)
		if err != nil {
			panic("Context not found: " + idTarget)
		}

		if len(args) > 2 {
			var err error
			intervalIndex, err = strconv.Atoi(args[1])
			util.Checkm(err, "Unable to parse id")

			if intervalIndex < 0 || intervalIndex > len(ctxSrc.Intervals)-1 {
				panic("interval index out of range")
			}

			mgr.MoveIntervalByIndex(idSrc, idTarget, intervalIndex)
		} else {
			for index, interval := range ctxSrc.Intervals {
				fmt.Printf("[%d] %s - %s\n", index, interval.Start.Time.Format(time.RFC3339), interval.End.Time.Format(time.RFC3339))
			}
		}

	},
}

func init() {
	editContextIntervalCmd.AddCommand(editContextIntervalMoveCmd)
}
