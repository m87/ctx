package bootstrap

import (
	"os"

	"github.com/m87/ctx/core"
	localstorage "github.com/m87/ctx/storage/local"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func CreateManager() *core.ContextManager {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ctx")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
	}
	return localstorage.CreateManager()
}
