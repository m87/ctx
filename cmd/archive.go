package cmd

import (
	"fmt"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func CreateArchiveCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "archive",
		Aliases: []string{"a"},
		Short:   "Archive contexts",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("archive called")
		},
	}
}

func init() {
	cmd := CreateArchiveCmd(bootstrap.CreateManager())
	flags.AddContxtFlag(cmd)
	cmd.Flags().BoolP("all", "a", false, "Archive all contexts")
	rootCmd.AddCommand(cmd)
}
