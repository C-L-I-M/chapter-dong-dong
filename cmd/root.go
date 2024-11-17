package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "chapter-dong-dong",
	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(prerun)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().String("config", "", "config file (default is config.yaml)")
	rootCmd.PersistentFlags().String("discord_token", "", "Bot access token")
	rootCmd.PersistentFlags().String("server_id", "", "Server id")

	viper.BindPFlag("discord_token", rootCmd.PersistentFlags().Lookup("discord_token"))
	viper.BindPFlag("server_id", rootCmd.PersistentFlags().Lookup("server_id"))
}

// prerun reads in config file and ENV variables if set.
func prerun() {
	// Reading .env config
	cwd, err := os.Getwd()
	cobra.CheckErr(err)
	viper.AddConfigPath(cwd)

	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Debug("No .env file found", err)
	} else {
		log.Info("Using .env file:", viper.ConfigFileUsed())
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.MergeInConfig(); err != nil {
		log.Debug("No config.yaml file found", err)
	} else {
		log.Info("Merging config file into it:", viper.ConfigFileUsed())
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		log.Info("Merging config file into it:", cfgFile)
	} else {
		cfgFile = viper.ConfigFileUsed()
	}
}
