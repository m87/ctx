/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		util.Check(err, "Unable to get user home dir")

		os.RemoveAll(viper.GetString("ctxPath"))
		os.Remove(filepath.Join(home, ".ctx"))
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
