package cmd

import (
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func NewDeleteContextCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "d", "rm"},
		Short:   "Delete context",
		Run: func(cmd *cobra.Command, args []string) {
			description := strings.TrimSpace(args[0])
			id, err := util.Id(description, false)
			util.Checkm(err, "Unable to process id "+description)

			util.Check(manager.Delete(id))
		},
	}

}

var deleteCmd = NewDeleteContextCmd(bootstrap.CreateManager())

func init() {
	rootCmd.AddCommand(deleteCmd)
}
