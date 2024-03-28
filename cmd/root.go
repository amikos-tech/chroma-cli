package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "chroma",
	Short:   "Chroma Command Line Interface.",
	Long:    `Utility to manage local and remote Chroma servers.`,
	Version: "0.0.0",
}

type HomeDirProvider interface {
	GetHomeDir() (string, error)
}

type DefaultHomeDirProvider struct{}

func (d DefaultHomeDirProvider) GetHomeDir() (string, error) {
	return homedir.Dir()
}

func init() {
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
