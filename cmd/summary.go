/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
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
		isDescription, _ := cmd.Flags().GetBool("description")

		if isDescription {
			id = util.GenerateId(id)
		}

		state := ctx.Load()
		ctx := state.Contexts[id]
		fmt.Println(ctx.Description)
		fmt.Println("------------------------")
		fmt.Printf("duration: %s\n", ctx.Duration)
		fmt.Println("intervals:")

		for _, v := range ctx.Intervals {
			fmt.Printf("\t [%s-%s] %s\n", v.Start, v.End, v.Duration)
		}

	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)
	summaryCmd.Flags().BoolP("description", "d", false, "stop by description")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// summaryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// summaryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
