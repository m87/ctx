package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initialize environment",
	Long:    `Creates default directories and config file`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(viper.ConfigFileUsed()) > 0 {
		} else {
			home, err := os.UserHomeDir()
			util.Checkm(err, "Unable to get user home dir")

			os.Mkdir(filepath.Join(home, ".ctx.d"), 0777)
			os.Mkdir(filepath.Join(home, ".ctx.d", "archive"), 0777)
			os.WriteFile(filepath.Join(home, ".ctx"), []byte(fmt.Sprintf(`storePath: %s`, filepath.Join(home, ".ctx.d"))), 0777)
			os.WriteFile(filepath.Join(home, ".ctx.d", "state"), []byte("{\"contexts\": {}}"), 0777)
			os.WriteFile(filepath.Join(home, ".ctx.d", "events"), []byte("{\"events\": []}"), 0777)
		}

	},
}

func init() {
	admCmd.AddCommand(initCmd)
}
