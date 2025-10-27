package bootstrap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/m87/ctx/core"
	localstorage "github.com/m87/ctx/storage/local"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func CreateManager() *core.ContextManager {
	l, err := core.LockWithTimeout()
	if err != nil {
		panic(err)
	}
	defer l.Unlock()

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
	viper.ReadInConfig()

	if len(viper.ConfigFileUsed()) <= 0 {
		log.Println("Ctx is not initialized. Iinitializing with defaults")
		InitDefault()

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ctx")
	}

	viper.ReadInConfig()

	return localstorage.CreateManager()
}

func InitDefault() {
	home, err := os.UserHomeDir()
	util.Checkm(err, "Unable to get user home dir")

	os.Mkdir(filepath.Join(home, ".ctx.d"), 0777)
	os.Mkdir(filepath.Join(home, ".ctx.d", "archive"), 0777)
	os.WriteFile(filepath.Join(home, ".ctx"), []byte(fmt.Sprintf(`version: %s
storePath: %s`, core.Release, filepath.Join(home, ".ctx.d"))), 0777)
	os.WriteFile(filepath.Join(home, ".ctx.d", "state"), []byte("{\"contexts\": {}}"), 0777)
}
