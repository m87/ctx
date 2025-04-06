package cmd

import (
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
}

func init() {
	rootCmd.AddCommand(editCmd)
}
