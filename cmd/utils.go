package cmd

import (
	"context"
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

func collectionExists(client *chroma.Client, collectionName string) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("client is nil")
	}
	if collectionName == "" {
		return false, fmt.Errorf("collectionName is empty")
	}
	collections, err := client.ListCollections(context.TODO())
	if err != nil {
		return false, err
	}
	for _, collection := range collections {
		if collection.Name == collectionName {
			return true, nil
		}
	}
	return false, nil
}

func getCollection(client *chroma.Client, collectionName string) (*chroma.Collection, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if collectionName == "" {
		return nil, fmt.Errorf("collectionName is empty")
	}
	collections, err := client.ListCollections(context.TODO())
	if err != nil {
		return nil, err
	}
	for _, collection := range collections {
		if collection.Name == collectionName {
			return collection, nil
		}
	}
	return nil, fmt.Errorf("collection not found")
}
