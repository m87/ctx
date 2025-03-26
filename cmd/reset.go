package cmd

import (
	"os"
	"path/filepath"

	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Deletes all app files",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		util.Checkm(err, "Unable to get user home dir")

		os.RemoveAll(viper.GetString("storePath"))
		os.Remove(filepath.Join(home, ".ctx"))
	},
}

func init() {
	admCmd.AddCommand(resetCmd)
}
