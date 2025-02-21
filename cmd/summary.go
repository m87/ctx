/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func contextSummary(state *ctx_model.State, input string, isRaw bool) {
	id, err := util.Id(input, isRaw)
	util.Check(err, "Unable to process id "+input)

	ctx := state.Contexts[id]
	fmt.Println(ctx.Description)
	fmt.Println("------------------------")
	fmt.Printf("duration: %s\n", ctx.Duration)
	fmt.Println("comments:")

	for _, v := range ctx.Comments {
		fmt.Printf("\t- %s\n", v)
	}
	fmt.Println("intervals:")

	for _, v := range ctx.Intervals {
		fmt.Printf("\t [%s-%s] %s\n", v.Start, v.End, v.Duration)
	}
}

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
		util.ApplyPatch(func(state *ctx_model.State) {
			isRaw, _ := cmd.Flags().GetBool("raw")
			contextSummary(state, args[0], isRaw)
		})

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
