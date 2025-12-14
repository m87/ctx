package cmd

import (
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	return &cobra.Command{
		Use: "edit",
		Short: "Edit given resource",
	}
}

var editCmd = newEditCmd()

func init() {
	rootCmd.AddCommand(editCmd)
}
