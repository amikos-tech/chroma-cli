package cmd

import (
	"fmt"

	"chroma/utils"
	"github.com/spf13/viper"

	"github.com/amikos-tech/chroma-go"
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
	client, err := chroma.NewClient(fmt.Sprintf("%v://%v:%v", scheme, serverConfig["host"], serverConfig["port"]), chroma.WithDebug(false))
	if err != nil {
		return nil, err
	}
	return client, nil
}
