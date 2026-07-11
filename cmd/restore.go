package cmd

import "github.com/spf13/cobra"

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a context",
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
