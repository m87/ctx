package cmd

import (
	"fmt"
	"os"
	"strings"

	ctxlog "github.com/m87/ctx/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var RemoteAddr string
var OutputFormat string
var Verbose bool

var rootCmd = &cobra.Command{
	Use: "ctx",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ctxlog.SetupLogger(viper.GetString("log_level"))
		format := strings.ToLower(strings.TrimSpace(OutputFormat))
		switch format {
		case "", "text", "json", "yaml", "shell":
			if format == "" {
				OutputFormat = "text"
			} else {
				OutputFormat = format
			}
			return nil
		default:
			return fmt.Errorf("unsupported output format: %s", OutputFormat)
		}
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

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&RemoteAddr, "remote", "r", "", "Remote server address")
	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "text", "Output format: text|json|yaml|shell")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output (include detailed fields and intervals)")
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

	if RemoteAddr == "" {
		RemoteAddr = viper.GetString("remote")
	}
}
