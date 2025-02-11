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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := strings.TrimSpace(args[0])
		if id == "" {
			return
		}
		state := ctx.Load()

		isDescription, _ := cmd.Flags().GetBool("description")

		if isDescription {
			id = util.GenerateId(id)
		}

		ctx.Delete(id, &state)

		ctx.Save(&state)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("description", "d", false, "stop by description")
	//TODO flag for history
}
