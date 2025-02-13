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

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := strings.TrimSpace(args[0])
		comment := strings.TrimSpace(args[1])
		if id == "" || comment == "" {
			return
		}
		state := ctx.Load()

		isDescription, _ := cmd.Flags().GetBool("description")

		if isDescription {
			id = util.GenerateId(id)
		}

		ctx.Comment(id, comment, &state)
		ctx.Save(&state)
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)
	commentCmd.Flags().BoolP("description", "d", false, "stop by description")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
