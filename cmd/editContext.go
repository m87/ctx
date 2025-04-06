package cmd

import (
	"github.com/spf13/cobra"
)

var editContextCmd = &cobra.Command{
	Use:   "context",
}

func init() {
	editCmd.AddCommand(editContextCmd)
}
