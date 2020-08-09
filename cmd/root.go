package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var keychainName string
var defaultKeychaingName string

var rootCmd = &cobra.Command{
	Use:   "gopswdcli",
	Short: "Small tool for storing and retrieving secrets using OS's keychain",
	Long:  ``,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	if runtime.GOOS == "linux" {
		defaultKeychaingName = "Login"
	} else if runtime.GOOS == "darwin" {
		defaultKeychaingName = "login"
	} else {
		fmt.Printf("unsupported OS %s\n", runtime.GOOS)
	}
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gopswdcli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&keychainName, "keychain", "k", defaultKeychaingName, "name of the keychain")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gopswdcli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gopswdcli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
