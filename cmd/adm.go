package cmd

import (
	"github.com/spf13/cobra"
)

var admCmd = &cobra.Command{
	Use:     "admin",
	Aliases: []string{"adm"},
	Short:   "Admin command",
}

func init() {
	rootCmd.AddCommand(admCmd)
}
