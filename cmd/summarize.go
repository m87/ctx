package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func newSummarizeCmd(manager *core.ContextManager) *cobra.Command {

	return &cobra.Command{
		Use:     "summarize",
		Aliases: []string{"sum", "s"},
	}

}

var summarizeCmd = newSummarizeCmd(bootstrap.CreateManager())
func init() {
	rootCmd.AddCommand(summarizeCmd)
}
