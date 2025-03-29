package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var summarizeDayCmd = &cobra.Command{
	Use:     "day",
	Aliases: []string{"d", "day"},
	Short:   "Summarize day",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("summarizeDay called")
	},
}

func init() {
	summarizeCmd.AddCommand(summarizeDayCmd)
}
