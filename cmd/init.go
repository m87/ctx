package cmd

import (
	"os"
	"path/filepath"

	"github.com/m87/ctx/bootstrap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newAdmInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize environment",
		Long: `Initialize environment. For example:
		ctx adm init
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(viper.ConfigFileUsed()) > 0 {
				path := viper.GetString("storePath")

				os.Mkdir(path, 0777)
				os.Mkdir(filepath.Join(path, "archive"), 0777)
				os.WriteFile(filepath.Join(path, "state"), []byte("{\"contexts\": {}}"), 0777)
			} else {
				bootstrap.InitDefault()
			}

		},
	}

}

func init() {
	admCmd.AddCommand(newAdmInitCmd())
}
