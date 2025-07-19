package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func NewListCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List contexts",
		Run: func(cmd *cobra.Command, args []string) {
			json, err := flags.ResolveJsonFlag(cmd)
			util.Check(err)
			verbose, err := flags.ResolveVerboseFlag(cmd)
			util.Check(err)

			if json {
				manager.ListJson()
			} else if verbose {
				manager.ListFull()
			} else {
				manager.List()
			}
		},
	}
}

func init() {
	cmd := NewListCmd(bootstrap.CreateManager())
	flags.AddVerboseFlag(cmd)
	flags.AddJsonFlag(cmd)
	rootCmd.AddCommand(cmd)
}
