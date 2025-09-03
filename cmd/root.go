package cmd

import (
	"fmt"
	"os"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var rootCmd = cobra.Command{
	Use:     "ctx",
	Version: core.Version,
	Run: func(cmd *cobra.Command, args []string) {
		manager := bootstrap.CreateManager()
		util.Check(manager.WithSession(func(session core.Session) error {
			ctx, err := session.GetActiveCtx()
			if err != nil {
				return err
			} else {
				fmt.Println(ctx.Description)
			}
			return nil
		}))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ctx.yaml)")
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

func initConfig() {
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
}
