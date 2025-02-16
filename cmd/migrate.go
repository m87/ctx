/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/m87/ctx/events"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		version := viper.GetString("version")

		if version == "" {
			migrateV02()
		} else {
			log.Println("Config up to date")
		}

	},
}

func migrateV02() {
	home, _ := os.UserHomeDir()
	os.Mkdir(filepath.Join(home, ".ctx.d", "archive"), 0777)

	eventsRegsitry := events.Load()

	changedEvents := []events.Event{}

	for _, ev := range eventsRegsitry.Events {
		ev.UUID = uuid.NewString()
		changedEvents = append(changedEvents, ev)
	}

	eventsRegsitry.Events = changedEvents
	events.Save(&eventsRegsitry)

	viper.Set("version", 0.2)
	viper.WriteConfig()
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
