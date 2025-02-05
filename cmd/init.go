/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(viper.ConfigFileUsed()) > 0 {
			log.Print(viper.ConfigFileUsed())
		} else {
			home, err := os.UserHomeDir()
			util.Check(err, "Unable to get user home dir")

			os.Mkdir(filepath.Join(home, ".ctx.d"), 0777)
			os.WriteFile(filepath.Join(home, ".ctx"), []byte(fmt.Sprintf(`ctxPath: %s`, filepath.Join(home, ".ctx.d"))), 0777)
			os.WriteFile(filepath.Join(home, ".ctx.d", "state"), []byte("{\"Contexts\": {}}"), 0777)
			os.WriteFile(filepath.Join(home, ".ctx.d", "events"), []byte("{\"Events\": []}"), 0777)
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
