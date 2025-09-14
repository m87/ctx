package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newSummarizeContextCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
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

			manager.WithSession(func(session core.Session) error {

				if id == "" && session.State.CurrentId != "" {
					id = session.State.CurrentId
				} else {
					panic("No context selected")
				}

				ctx, err := session.GetCtx(id)
				util.Checkm(err, "Unable to find context "+id)

				fmt.Printf("Context: %s\n", ctx.Description)
				if verbose {
					fmt.Printf("Id: %s\n", ctx.Id)
				}
				fmt.Printf("Duration: %s\n", ctx.Duration)

				return nil
			})
		},
	}
}

func init() {
	cmd := newSummarizeContextCmd(bootstrap.CreateManager())
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	summarizeCmd.AddCommand(cmd)
}
