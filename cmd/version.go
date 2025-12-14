package cmd

import (
	"fmt"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version:", core.Release)
			fmt.Println("Commit:", core.Commit)
			fmt.Println("Date:", core.Date)

		},
	}
}

func init() {
	rootCmd.AddCommand(newVersionCmd())
}
