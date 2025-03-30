/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:     "switch",
	Aliases: []string{"sw", "s"},
	Short:   "Switch context",
	Long: `Switch context:
	- switch description, created if not exists
	- switch -i id"`,
	Run: func(cmd *cobra.Command, args []string) {
		description := strings.TrimSpace(args[0])
		byId, _ := cmd.Flags().GetBool("id")
		id, err := util.Id(description, byId)
		util.Checkm(err, "Unable to process id "+description)

		manager := ctx.CreateManager()
		if byId {
			util.Check(manager.Switch(id))
		} else {
			util.Check(manager.CreateIfNotExistsAndSwitch(id, description))
		}

	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().BoolP("id", "i", false, "stop by description")
}
