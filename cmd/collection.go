package cmd

import (
	"chroma/utils"
	"context"
	"fmt"
	"github.com/amikos-tech/chroma-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func getClient(serverAlias string) (*chroma.Client, error) {
	var serverConfig map[string]interface{}
	var err error
	if serverAlias == "" {
		serverConfig, err = utils.GetServer(viper.GetString("active_server"))
	} else {
		serverConfig, err = utils.GetServer(serverAlias)
	}
	if err != nil {
		return nil, err
	}
	var scheme string
	if serverConfig["secure"].(bool) {
		scheme = "https"
	} else {
		scheme = "http"
	}
	client, err := chroma.NewClient(fmt.Sprintf("%v://%v:%v", scheme, serverConfig["host"], serverConfig["port"]))
	if err != nil {
		return nil, err
	}
	return client, nil
}

var ListCollectionsCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all available collections",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient("")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		col, err := client.ListCollections(context.TODO())
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		for _, collection := range col {
			fmt.Printf("%v\n", collection)
		}
	},
}

var collectionCommand = &cobra.Command{
	Use:     "collection",
	Aliases: []string{"c"},
	Short:   "Manage Chroma servers",
	Long:    ``,
}

func init() {
	rootCmd.AddCommand(collectionCommand)
	rootCmd.AddCommand(ListCollectionsCommand)
	collectionCommand.AddCommand(ListCollectionsCommand)
}
