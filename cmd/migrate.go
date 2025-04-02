package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrate called")
	},
}

func init() {
	admCmd.AddCommand(migrateCmd)
}
