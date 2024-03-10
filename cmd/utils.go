package cmd

import (
	"context"
	"fmt"
	"os"

	"chroma/utils"
	"github.com/spf13/viper"

	"github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/cohere"
	"github.com/amikos-tech/chroma-go/hf"
	"github.com/amikos-tech/chroma-go/openai"
	"github.com/amikos-tech/chroma-go/types"
)

func embeddingFunctionForString(embedder string, err error) (types.EmbeddingFunction, error) {
	if err != nil {
		return nil, err
	}
	if embedder == "" {
		return nil, nil // TODO this is rather a hack
	}
	switch embedder {
	case "openai":
		return openai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
	case "cohere":
		return cohere.NewCohereEmbeddingFunction(os.Getenv("COHERE_API_KEY")), nil
	case "hf":
		return hf.NewHuggingFaceEmbeddingFunction(os.Getenv("HF_API_KEY"), os.Getenv("HF_MODEL")), nil
	case "hash": // dummy embedding function
		return types.NewConsistentHashEmbeddingFunction(), nil
	default:
		return nil, fmt.Errorf("embedding function not found")
	}
}
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
