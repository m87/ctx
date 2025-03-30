package cmd

import (
	"github.com/spf13/cobra"
)

var summarizeCmd = &cobra.Command{
	Use:     "summarize",
	Aliases: []string{"sum", "s"},
}

func init() {
	rootCmd.AddCommand(summarizeCmd)
}
