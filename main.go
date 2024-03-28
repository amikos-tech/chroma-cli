package main

import (
	"errors"
	"fmt"
	"github.com/amikos-tech/chroma-cli/chroma/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	Version   = "0.0.0-development-build" // Replaced at build time
	BuildDate = "9999-12-31"              // Replace with the actual build date
)

type ChromaCLI struct {
	rootCmd         *cobra.Command
	homeDirProvider cmd.HomeDirProvider
}
type CliOption func(*ChromaCLI) error

func WithHomeDirProvider(provider cmd.HomeDirProvider) CliOption {
	return func(c *ChromaCLI) error {
		c.homeDirProvider = provider
		return nil
	}
}

func (c *ChromaCLI) Initialize(options ...CliOption) error {
	for _, option := range options {
		err := option(c)
		if err != nil {
			return err
		}
	}
	c.rootCmd = cmd.RootCmd
	c.rootCmd.SetVersionTemplate(fmt.Sprintf("Chroma version %s, build date %s\n", Version, BuildDate))
	cobra.OnInitialize(func() {
		home, err := c.homeDirProvider.GetHomeDir()
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
				// create config file
				err := os.MkdirAll(home+"/.chroma", 0700)
				if err != nil {
					// Unable to create directory
					log.Fatal(err)
				}
				_, err = os.Create(home + "/.chroma" + "/config.yaml")
				if err != nil {
					// Unable to create file
					log.Fatal(err)
				}
				err = viper.WriteConfig()
				if err != nil {
					fmt.Println("Can't initialize config:", err)
					os.Exit(1)
				}
			}
		}
	})
	err := c.rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	cli := &ChromaCLI{}
	err := cli.Initialize(WithHomeDirProvider(cmd.DefaultHomeDirProvider{}))
	if err != nil {
		fmt.Printf("Error initializing CLI: %s", err)
		os.Exit(1)
	}
}
