package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var VersionCommand = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Get the version of the Chroma Server",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient("")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		version, err := client.Version(context.TODO())
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		cmd.Printf("Chroma Server Version: %v\n", version)
	},
}

func init() {
	VersionCommand.PersistentFlags().StringP("alias", "s", "", "Server alias")
	rootCmd.AddCommand(VersionCommand)
}
