package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewAdmInitCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize environment",
		Long:    `Creates default directories and config file`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(viper.ConfigFileUsed()) > 0 {
				path := viper.GetString("storePath")

				os.Mkdir(path, 0777)
				os.WriteFile(filepath.Join(path, "state"), []byte("{\"contexts\": {}}"), 0777)
			} else {
				home, err := os.UserHomeDir()
				util.Checkm(err, "Unable to get user home dir")

				os.Mkdir(filepath.Join(home, ".ctx.d"), 0777)
				os.WriteFile(filepath.Join(home, ".ctx"), []byte(fmt.Sprintf(`version: %s
storePath: %s`, core.Version, filepath.Join(home, ".ctx.d"))), 0777)
				os.WriteFile(filepath.Join(home, ".ctx.d", "state"), []byte("{\"contexts\": {}}"), 0777)
			}

		},
	}

}

func init() {
	admCmd.AddCommand(NewAdmInitCmd(bootstrap.CreateManager()))
}
