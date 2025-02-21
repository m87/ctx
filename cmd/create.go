/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func createContext(state *ctx_model.State, input string, isRaw bool) {
	id, err := util.Id(input, isRaw)
	util.Check(err, "Unable to process id "+input)

	state.Contexts[id] = ctx_model.Context{
		Id:          id,
		Description: strings.TrimSpace(input),
		State:       ctx_model.ACTIVE,
		Intervals:   []ctx_model.Interval{},
	}
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.ApplyPatch(func(state *ctx_model.State) {
			isRaw, _ := cmd.Flags().GetBool("raw")
			createContext(state, args[0], isRaw)
		})
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
