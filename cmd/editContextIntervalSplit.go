/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newEditContextIntervalSplitCmd(manager *core.ContextManager) *cobra.Command {

	return &cobra.Command{
		Use:     "split",
		Aliases: []string{"s"},
		Short:   "Split interval",
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
					util.Checkm(err, "Unable to parse interval index")

					h, _ := strconv.Atoi(args[2])
					m, _ := strconv.Atoi(args[3])
					s, _ := strconv.Atoi(args[4])

					session.SplitContextIntervalById(id, intervalId, h, m ,s)
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
func init() {
	cmd := newEditContextIntervalSplitCmd(bootstrap.CreateManager())
	editContextIntervalCmd.AddCommand(cmd)
}
