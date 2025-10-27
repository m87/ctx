package cmd

import (
	"fmt"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func newVersionCmd(manager *core.ContextManager) *cobra.Command {

	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v", "ver"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(core.Release)
		},
	}
}

func init() {
	rootCmd.AddCommand(newVersionCmd(bootstrap.CreateManager()))
}
