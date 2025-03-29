package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var summarizeContextCmd = &cobra.Command{
	Use:     "context",
	Aliases: []string{"ctx", "c"},
	Short:   "Summarize context",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("summarizeContext called")
	},
}

func init() {
	summarizeCmd.AddCommand(summarizeContextCmd)
}
