package cmd

import (
	"github.com/spf13/cobra"
)

var labelIntervalCmd = &cobra.Command{
	Use:     "labelInterval",
	Aliases: []string{"li"},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	intervalCmd.AddCommand(labelIntervalCmd)
}
