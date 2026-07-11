package cmd

import "github.com/spf13/cobra"

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive a context",
}

func init() {
	rootCmd.AddCommand(archiveCmd)
}
