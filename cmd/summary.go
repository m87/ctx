package cmd

import "github.com/spf13/cobra"

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show summaries",
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}
