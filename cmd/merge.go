package cmd

import "github.com/spf13/cobra"

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge resources",
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}
