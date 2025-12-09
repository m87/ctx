package cmd

import (
	"github.com/spf13/cobra"
)

func newAdminCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "admin",
		Short: "Admin command",
	}
}

var admCmd = newAdminCmd()

func init() {
	rootCmd.AddCommand(admCmd)
}
