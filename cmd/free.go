package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)


func NewFreeCmd(manager *core.ContextManager) *cobra.Command {
return &cobra.Command{
	Use:     "free",
	Aliases: []string{"f"},
	Short:   "Stop current context",
	Run: func(cmd *cobra.Command, args []string) {
	//	util.Check(manager.Free())
	},
}
}


func init() {
	cmd := NewFreeCmd(bootstrap.CreateManager())
	rootCmd.AddCommand(cmd)
}
