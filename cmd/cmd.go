package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wantedly/ev-cli/config"
	"os"
)

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ev/config.yml)")
}

// RootCmd represents the base command when called without any subcommands
var (
	RootCmd = &cobra.Command{
		Use:           "ev",
		Short:         "CLI tool for managing evaluation data",
		Long:          "CLI tool for managing evaluation data",
		SilenceErrors: true,
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// initConfig reads in config file
func initConfig() {
	v := viper.New()

	// Read config file
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName("config")
		v.AddConfigPath("$HOME/.ev")
	}
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Ignore ConfigFileNotFoundError
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	config.InitWithViper(v)
}
