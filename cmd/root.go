package cmd

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chroma",
	Short: "Chroma",
	Long:  ``,
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
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home + "/.chroma")
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			//create config file
			err := os.MkdirAll(home+"/.chroma", 0700)
			if err != nil {
				// Unable to create directory
				log.Fatal(err)
			}
			_, err = os.Create(home + "/.chroma" + "/config.yaml")
			if err != nil {
				//Unable to create file
				log.Fatal(err)
			}
			err = viper.WriteConfig()
			if err != nil {
				fmt.Println("Can't initialize config:", err)
				os.Exit(1)
			}
		}

	}
}