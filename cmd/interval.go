package cmd

import (
	"github.com/spf13/cobra"
)

var intervalCmd = &cobra.Command{
	Use: "interval",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(intervalCmd)
}
