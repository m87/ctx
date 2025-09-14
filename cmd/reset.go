package cmd

import (
	"os"
	"path/filepath"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newAdmResetCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "reset",
		Aliases: []string{"r"},
		Short:   "Deletes all app files",
		Run: func(cmd *cobra.Command, args []string) {
			home, err := os.UserHomeDir()
			util.Checkm(err, "Unable to get user home dir")

			os.RemoveAll(viper.GetString("storePath"))
			os.Remove(filepath.Join(home, ".ctx"))
		},
	}

}

func init() {
	admCmd.AddCommand(newAdmResetCmd(bootstrap.CreateManager()))
}
