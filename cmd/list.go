package cmd

import (
	"fmt"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/cmd/flags"
	"github.com/m87/ctx/cmd/tui"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newListCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List contexts",
		Run: func(cmd *cobra.Command, args []string) {
			json, err := flags.ResolveJsonFlag(cmd)
			util.Check(err)
			verbose, err := flags.ResolveVerboseFlag(cmd)
			util.Check(err)
			bash, err := flags.ResolveShellFlag(cmd)
			util.Check(err)

			manager.WithSession(func(session core.Session) error {
				if json {
					fmt.Println(tui.ListJson(session))
				} else if bash {
					fmt.Println(tui.ListBash(session))
				} else if verbose {
					fmt.Println(tui.ListFull(session))
				} else {
					fmt.Println(tui.List(session))
				}

				return nil
			})
		},
	}
}

func init() {
	cmd := newListCmd(bootstrap.CreateManager())
	flags.AddVerboseFlag(cmd)
	flags.AddShellFlag(cmd)
	flags.AddJsonFlag(cmd)
	rootCmd.AddCommand(cmd)
}
