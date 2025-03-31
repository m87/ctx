package cmd

import (
	"fmt"
	"os"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "ctx",
	Run: func(cmd *cobra.Command, args []string) {
    mgr := ctx.CreateManager()

    mgr.ContextStore.Read(func(s *ctx_model.State) error {
      if s.CurrentId != "" {
        fmt.Println(s.Contexts[s.CurrentId].Description)
      }
      return nil
    })
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
