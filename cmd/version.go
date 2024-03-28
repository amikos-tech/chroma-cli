package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var VersionCommand = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Get the version of the Chroma Server. If alias is not specified the currently active server is used.",
	Run: func(cmd *cobra.Command, args []string) {
		activeAlias := viper.GetString("active_server")
		alias, err := getStringFlagIfChangedWithDefault(cmd, "alias", &activeAlias)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		client, err := getClient(*alias)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		version, err := client.Version(context.TODO())
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		cmd.Printf("Chroma Server Version: %v\n", version)
	},
}

func init() {
	VersionCommand.Flags().StringP("alias", "s", "", "Server alias")
	RootCmd.AddCommand(VersionCommand)
}
