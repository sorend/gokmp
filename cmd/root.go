package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	// jww "github.com/spf13/jwalterweatherman"
)

var rootCmd = &cobra.Command{
	Use:   "gokmp",
	Short: "gokmp: Backup your Flickr pictures",
	Long: `
                    ___
     _______ ______/   / ___ _  ___________
    / __   // _   //  / /  /  \/ __\    __ \
   /  |/   /  /   /  /_/  /   ______\   |/  \
  /   /   /  /   /     __/       ___/   /   /
  \__    /______/__/\   \ _/\__/  _/   ____/  Go Keep My Photos
==__/   / ========== \___\ == /___/   / =========================
 /_____/                         /___/

gokmp: Commandline tool for backing up your Flickr pictures locally.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

	// jww.SetLogThreshold(jww.LevelTrace)
	// jww.SetStdoutThreshold(jww.LevelInfo)
	viper.SetConfigName("gokmp")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/gokmp/")
	viper.AddConfigPath("$HOME/.gokmp")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("GOKMP")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		// do nothing
	}
}
