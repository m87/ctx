package cmd

import (
	"fmt"

	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v", "ver"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(core.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
