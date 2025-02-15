/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.ApplyPatch(func(state *ctx_model.State) {
			if j, _ := cmd.Flags().GetBool("json"); j {
				v := make([]ctx_model.Context, 0, len(state.Contexts))
				for _, c := range state.Contexts {
					v = append(v, c)
				}
				s, _ := json.Marshal(v)

				fmt.Printf("%s", string(s))
			} else {
				for _, v := range state.Contexts {

					if f, _ := cmd.Flags().GetBool("full"); f {
						fmt.Printf("- [%s] %s\n", v.Id, v.Description)
						for _, interval := range v.Intervals {
							fmt.Printf("\t- %s - %s\n", interval.Start.Local().Format(time.DateTime), interval.End.Local().Format(time.DateTime))
						}
					} else {
						fmt.Printf("- %s\n", v.Description)
					}

				}
			}
		})
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("full", "f", false, "show full list")
	listCmd.Flags().BoolP("json", "j", false, "show list as json")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
