package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/amikos-tech/chroma-go/collection"
	"github.com/amikos-tech/chroma-go/types"
)

func listCollections(cmd *cobra.Command, args []string) error {
	activeAlias := viper.GetString("active_server")
	alias, err := getStringFlagIfChangedWithDefault(cmd, "alias", &activeAlias)
	if err != nil {
		return err
	}
	client, err := getClient(*alias)
	if err != nil {
		cmd.Printf("%v\n", err)
		return err
	}
	colList, err := client.ListCollections(context.TODO())
	if err != nil {
		cmd.Printf("%v\n", err)
		return err
	}
	for _, col := range colList {
		cmd.Printf("%v\n", col) // TODO make into a table?
	}
	return nil
}

var ListCollectionsCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all available collections",
	Run: func(cmd *cobra.Command, args []string) {
		err := listCollections(cmd, args)
		if err != nil {
			os.Exit(1)
		}
	},
}

func createCollection(cmd *cobra.Command, args []string) error {
	collectionName := args[0]
	activeAlias := viper.GetString("active_server")
	alias, err := getStringFlagIfChangedWithDefault(cmd, "alias", &activeAlias)
	if err != nil {
		return err
	}
	client, err := getClient(*alias)
	if err != nil {
		cmd.Printf("%v\n", err)
		return err
	}

	var options = make([]collection.Option, 0)
	options = append(options, collection.WithName(collectionName))

	if mVal, err := getIntFlagIfChangedWithDefault(cmd, "m", nil); err != nil {
		cmd.Printf("invalid m: %v\n", err)
		return err
	} else if mVal != nil {
		options = append(options, collection.WithHNSWM(int32(*mVal)))
	}
	if constructionEfVal, err := getIntFlagIfChangedWithDefault(cmd, "construction-ef", nil); err != nil {
		cmd.Printf("invalid construction-ef: %v\n", err)
		return err
	} else if constructionEfVal != nil {
		options = append(options, collection.WithHNSWConstructionEf(int32(*constructionEfVal)))
	}

	if searchEfVal, err := getIntFlagIfChangedWithDefault(cmd, "search-ef", nil); err != nil {
		cmd.Printf("invalid search-ef: %v\n", err)
		return err
	} else if searchEfVal != nil {
		options = append(options, collection.WithHNSWSearchEf(int32(*searchEfVal)))
	}

	if batchSizeVal, err := getIntFlagIfChangedWithDefault(cmd, "batch-size", nil); err != nil {
		cmd.Printf("invalid batch-size: %v\n", err)
		return err
	} else if batchSizeVal != nil {
		options = append(options, collection.WithHNSWBatchSize(int32(*batchSizeVal)))
	}

	if syncThresholdVal, err := getIntFlagIfChangedWithDefault(cmd, "sync-threshold", nil); err != nil {
		cmd.Printf("invalid sync-threshold: %v\n", err)
		return err
	} else if syncThresholdVal != nil {
		options = append(options, collection.WithHNSWSyncThreshold(int32(*syncThresholdVal)))
	}

	if threadsVal, err := getIntFlagIfChangedWithDefault(cmd, "threads", nil); err != nil {
		cmd.Printf("invalid threads: %v\n", err)
	} else if threadsVal != nil && *threadsVal > 0 {
		options = append(options, collection.WithHNSWNumThreads(int32(*threadsVal)))
	}

	if ensure, err := cmd.Flags().GetBool("ensure"); err != nil {
		cmd.Printf("invalid ensure: %v\n", err)
		return err
	} else if ensure {
		options = append(options, collection.WithCreateIfNotExist(ensure))
	}

	if resizeFactorVal, err := getFloatFlagIfChangedWithDefault(cmd, "resize-factor", nil); err != nil {
		cmd.Printf("invalid resize-factor: %v\n", err)
		return err
	} else if resizeFactorVal != nil {
		options = append(options, collection.WithHNSWResizeFactor(*resizeFactorVal))
	}

	if spaceVar, err := getStringFlagIfChangedWithDefault(cmd, "space", nil); err != nil {
		cmd.Printf("invalid space: %v\n", err)
		return err
	} else if spaceVar != nil {
		df, err := types.ToDistanceFunction(*spaceVar)
		if err != nil {
			cmd.Printf("invalid distance function: %v\n", err)
			return err
		}
		options = append(options, collection.WithHNSWDistanceFunction(df))
	}

	metadatasVar, err := getStringSliceFlagIfChangedWithDefault(cmd, "meta", &[]string{})
	if err != nil {
		cmd.Printf("invalid meta: %v\n", err)
		return err
	}
	if cmd.Flag("meta").Changed {
		metadata := make(map[string]interface{})
		for _, meta := range *metadatasVar {
			kvPair := strings.Split(meta, "=")
			if len(kvPair) != 2 {
				cmd.Printf("invalid metadata format: %v. should be key=value.", meta)
				os.Exit(1)
			}
			if b, err := strconv.ParseBool(kvPair[1]); err == nil {
				metadata[kvPair[0]] = b
			} else if f, err := strconv.ParseFloat(kvPair[1], 32); strings.Contains(kvPair[1], ".") && err == nil {
				metadata[kvPair[0]] = float32(f)
			} else if i, err := strconv.ParseInt(kvPair[1], 10, 32); err == nil {
				metadata[kvPair[0]] = i
			} else {
				metadata[kvPair[0]] = kvPair[1]
			}
		}
		options = append(options, collection.WithMetadatas(metadata))
	}

	_, err = client.NewCollection(
		context.Background(),
		options...,
	)
	if err != nil {
		cmd.Printf("failed to create collection: %v\n", err)
		return err
	}
	cmd.Printf("Collection created: %v\n", collectionName)

	return nil
}

var CreateCollectionCommand = &cobra.Command{
	Use:       "create",
	Aliases:   []string{"c"},
	Short:     "Create a new collection",
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: []string{"name"},
	Run: func(cmd *cobra.Command, args []string) {
		err := createCollection(cmd, args)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	},
}

func deleteCollection(cmd *cobra.Command, args []string) error {
	collectionName := args[0]
	activeAlias := viper.GetString("active_server")
	alias, err := getStringFlagIfChangedWithDefault(cmd, "alias", &activeAlias)
	if err != nil {
		return err
	}
	client, err := getClient(*alias)
	if err != nil {
		cmd.Printf("%v\n", err)
		return err
	}
	_, err = client.DeleteCollection(context.TODO(), collectionName)
	if err != nil {
		cmd.Printf("%v\n", err)
		return err
	}
	cmd.Printf("Collection deleted: %v\n", collectionName)
	return nil
}

var DeleteCollectionCommand = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm"},
	Short:   "Delete a collection",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := deleteCollection(cmd, args)
		if err != nil {
			os.Exit(1)
		}
	},
}

func cloneCollection(cmd *cobra.Command, args []string) error {
	sourceCollectionName := args[0]
	destinationCollectionName := args[1]
	activeAlias := viper.GetString("active_server")
	alias, err := getStringFlagIfChangedWithDefault(cmd, "alias", &activeAlias)
	if err != nil {
		return err
	}
	client, err := getClient(*alias)
	if err != nil {
		cmd.Printf("%v\n", err)
		return err
	}
	sourceExists, err := collectionExists(client, sourceCollectionName)
	if err != nil {
		cmd.Printf("%v\n", err)
		return err
	}
	if !sourceExists {
		cmd.Printf("source collection %v does not exist\n", sourceCollectionName)
		return err
	}
	destinationExists, err := collectionExists(client, destinationCollectionName)
	if err != nil {
		cmd.Printf("%v\n", err)
		return err
	}
	if destinationExists {
		cmd.Printf("destination collection %v already exists\n", destinationCollectionName)
		return err
	}
	sourceCollection, err := getCollection(client, sourceCollectionName)
	if err != nil {
		cmd.Printf("%v\n", err)
		return err
	}
	count, err := sourceCollection.Count(context.TODO())
	if err != nil {
		cmd.Printf("%v\n", err)
		return err
	} else if count == 0 {
		cmd.Printf("source collection %v is empty\n", sourceCollectionName)
		return err
	}
	cloneBatchSize := 100
	cbv, err := getIntFlagIfChangedWithDefault(cmd, "clone-batch-size", &cloneBatchSize)
	if err != nil {
		cmd.Printf("invalid clone-batch-size: %v\n", err)
		return err
	} else if cbv != nil {
		cloneBatchSize = *cbv
	}

	spaceVal, err := getStringFlagIfChangedWithDefault(cmd, "space", getMetadataStringValue(sourceCollection.Metadata, types.HNSWSpace))
	if err != nil {
		cmd.Printf("invalid space: %v\n", err)
		return err
	}

	mVal, err := getIntFlagIfChangedWithDefault(cmd, "m", getMetadataIntValue(sourceCollection.Metadata, types.HNSWM))
	if err != nil {
		cmd.Printf("invalid m: %v\n", err)
		return err
	}

	constructionEfVal, err := getIntFlagIfChangedWithDefault(cmd, "construction-ef", getMetadataIntValue(sourceCollection.Metadata, types.HNSWConstructionEF))
	if err != nil {
		cmd.Printf("invalid construction-ef: %v\n", err)
		return err
	}
	searchEfVal, err := getIntFlagIfChangedWithDefault(cmd, "search-ef", getMetadataIntValue(sourceCollection.Metadata, types.HNSWSearchEF))
	if err != nil {
		cmd.Printf("invalid search-ef: %v\n", err)
		return err
	}
	batchSizeVal, err := getIntFlagIfChangedWithDefault(cmd, "batch-size", getMetadataIntValue(sourceCollection.Metadata, types.HNSWBatchSize))
	if err != nil {
		cmd.Printf("invalid batch-size: %v\n", err)
		return err
	}
	syncThresholdVal, err := getIntFlagIfChangedWithDefault(cmd, "sync-threshold", getMetadataIntValue(sourceCollection.Metadata, types.HNSWSyncThreshold))
	if err != nil {
		cmd.Printf("invalid sync-threshold: %v\n", err)
		return err
	}
	threadsVal, err := getIntFlagIfChangedWithDefault(cmd, "threads", getMetadataIntValue(sourceCollection.Metadata, types.HNSWNumThreads))
	if err != nil {
		cmd.Printf("invalid threads: %v\n", err)
		return err
	}
	resizeFactorVal, err := getFloatFlagIfChangedWithDefault(cmd, "resize-factor", getMetadataFloatValue(sourceCollection.Metadata, types.HNSWResizeFactor))
	if err != nil {
		cmd.Printf("invalid resize-factor: %v\n", err)
		return err
	}
	metadatasVar, err := getStringSliceFlagIfChangedWithDefault(cmd, "meta", &[]string{})
	if err != nil {
		cmd.Printf("invalid meta: %v\n", err)
		return err
	}
	var metadatasVal = make(map[string]interface{})
	for k, v := range sourceCollection.Metadata {
		if k == types.HNSWSpace || k == types.HNSWM || k == types.HNSWConstructionEF || k == types.HNSWSearchEF || k == types.HNSWBatchSize || k == types.HNSWSyncThreshold || k == types.HNSWNumThreads || k == types.HNSWResizeFactor {
			continue
		}
		metadatasVal[k] = v
	}

	if cmd.Flag("meta").Changed {
		for _, meta := range *metadatasVar {
			kvPair := strings.Split(meta, "=")
			if len(kvPair) != 2 {
				cmd.Printf("invalid metadata format: %v. should be key=value.", meta)
				return err
			}
			if b, err := strconv.ParseBool(kvPair[1]); err == nil {
				metadatasVal[kvPair[0]] = b
			} else if f, err := strconv.ParseFloat(kvPair[1], 32); strings.Contains(kvPair[1], ".") && err == nil {
				metadatasVal[kvPair[0]] = float32(f)
			} else if i, err := strconv.ParseInt(kvPair[1], 10, 32); err == nil {
				metadatasVal[kvPair[0]] = i
			} else {
				metadatasVal[kvPair[0]] = kvPair[1]
			}
		}
	}
	var collectionOptions = make([]collection.Option, 0)
	collectionOptions = append(collectionOptions, collection.WithName(destinationCollectionName))
	if df, err := types.ToDistanceFunction(*spaceVal); err != nil {
		cmd.Printf("invalid distance function: %v\n", err)
		return err
	} else {
		collectionOptions = append(collectionOptions, collection.WithHNSWDistanceFunction(df))
	}
	var hasEf = false
	if efVal, err := embeddingFunctionForString(cmd.Flags().GetString("embedding-function")); err != nil {
		cmd.Printf("invalid embedding-function: %v\n", err)
		return err
	} else if efVal != nil {
		hasEf = true
		collectionOptions = append(collectionOptions, collection.WithEmbeddingFunction(efVal))
	}

	if mVal != nil {
		collectionOptions = append(collectionOptions, collection.WithHNSWM(int32(*mVal)))
	}
	if constructionEfVal != nil {
		collectionOptions = append(collectionOptions, collection.WithHNSWConstructionEf(int32(*constructionEfVal)))
	}
	if searchEfVal != nil {
		collectionOptions = append(collectionOptions, collection.WithHNSWSearchEf(int32(*searchEfVal)))
	}
	if batchSizeVal != nil {
		collectionOptions = append(collectionOptions, collection.WithHNSWBatchSize(int32(*batchSizeVal)))
	}
	if syncThresholdVal != nil {
		collectionOptions = append(collectionOptions, collection.WithHNSWSyncThreshold(int32(*syncThresholdVal)))
	}
	if threadsVal != nil && *threadsVal > 0 {
		collectionOptions = append(collectionOptions, collection.WithHNSWNumThreads(int32(*threadsVal)))
	}
	if resizeFactorVal != nil {
		collectionOptions = append(collectionOptions, collection.WithHNSWResizeFactor(*resizeFactorVal))
	}
	if len(metadatasVal) > 0 {
		collectionOptions = append(collectionOptions, collection.WithMetadatas(metadatasVal))
	}
	targetCollection, err := client.NewCollection(context.TODO(),
		collectionOptions...,
	)
	if err != nil {
		cmd.Printf("%v\n", err)
		return err
	}

	var totalNumberOfRecordsCopied = 0

	for start := 0; start < int(count); start += cloneBatchSize {
		end := start + cloneBatchSize
		if end > int(count) {
			end = int(count)
		}

		result, err := sourceCollection.GetWithOptions(
			context.TODO(),
			types.WithOffset(int32(start)),
			types.WithLimit(int32(end)),
			types.WithInclude(types.IMetadatas, types.IDocuments, types.IEmbeddings),
		)
		if err != nil {
			cmd.Printf("%v\n", err)
			return err
		}
		var _embeddings []*types.Embedding
		if !hasEf {
			_embeddings = result.Embeddings
		}
		_, err = targetCollection.Add(context.TODO(), _embeddings, result.Metadatas, result.Documents, result.Ids)
		if err != nil { // TODO not great to exit on first error but for now that will do. Consider rollback?
			cmd.Printf("%v\n", err)
			return err
		}
		totalNumberOfRecordsCopied += len(result.Ids)
	}
	cmd.Printf("successfully cloned %v to %v. copied records: %v\n", sourceCollection.Name, targetCollection.Name, totalNumberOfRecordsCopied)
	return nil
}

var CloneCollectionCommand = &cobra.Command{
	Use:     "clone",
	Aliases: []string{"cp"},
	Short:   "Clone a collection",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := cloneCollection(cmd, args)
		if err != nil {
			os.Exit(1)
		}
	},
}

func getIntFlagIfChangedWithDefault(cmd *cobra.Command, flag string, defaultValue *int) (*int, error) {
	if cmd.Flag(flag).Changed {
		flagValue, err := cmd.Flags().GetInt(flag)
		if err != nil {
			return &flagValue, err
		}
		return &flagValue, nil
	}
	if defaultValue != nil {
		return defaultValue, nil
	}
	return nil, nil
}

func getStringFlagIfChangedWithDefault(cmd *cobra.Command, flag string, defaultValue *string) (*string, error) {
	if cmd.Flag(flag).Changed {
		flagValue, err := cmd.Flags().GetString(flag)
		if err != nil {
			return &flagValue, err
		}
		return &flagValue, nil
	}
	if defaultValue != nil {
		return defaultValue, nil
	}
	return nil, nil
}

func getFloatFlagIfChangedWithDefault(cmd *cobra.Command, flag string, defaultValue *float32) (*float32, error) {
	if cmd.Flag(flag).Changed {
		flagValue, err := cmd.Flags().GetFloat32(flag)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return nil, err
		}
		return &flagValue, nil
	}
	if defaultValue != nil {
		return defaultValue, nil
	}
	return nil, nil
}

func getMetadataIntValue(metadata map[string]interface{}, key string) *int {
	if val, ok := metadata[key].(int); ok {
		return &val
	} else if val, ok := metadata[key].(int32); ok {
		v := int(val)
		return &v
	} else if val, ok := metadata[key].(int64); ok {
		v := int(val)
		return &v
	}
	return nil
}

func getMetadataFloatValue(metadata map[string]interface{}, key string) *float32 {
	if metadata[key] != nil {
		v := metadata[key].(float32)
		return &v
	}
	return nil
}

func getMetadataStringValue(metadata map[string]interface{}, key string) *string {
	if metadata[key] != nil {
		v := metadata[key].(string)
		return &v
	}
	return nil
}

func getStringSliceFlagIfChangedWithDefault(cmd *cobra.Command, flag string, defaultValue *[]string) (*[]string, error) {
	if cmd.Flag(flag).Changed {
		flagValue, err := cmd.Flags().GetStringSlice(flag)
		if err != nil {
			return &flagValue, err
		}
		return &flagValue, nil
	}
	if defaultValue != nil {
		return defaultValue, nil
	}
	return nil, nil
}

var collectionCommand = &cobra.Command{
	Use:     "collection",
	Aliases: []string{"c"},
	Short:   "Manage Chroma servers",
	Long:    ``,
}

var metaSlice = make([]string, 0)

func init() {
	RootCmd.AddCommand(collectionCommand)
	RootCmd.AddCommand(ListCollectionsCommand)
	collectionCommand.AddCommand(ListCollectionsCommand)
	ListCollectionsCommand.Flags().StringP("alias", "s", "", "Server alias name. If not provided, the active server will be used.")
	CreateCollectionCommand.Flags().String("name", "", "Name of the collection")
	CreateCollectionCommand.Flags().StringP("alias", "s", "", "Server alias name. If not provided, the active server will be used.")
	CreateCollectionCommand.Flags().Bool("ensure", false, "Create collection only if it doesn't exist. Chroma will be queried before sending create, if the collection exists, exit with 0. The metadata will be overwritten.")
	CreateCollectionCommand.Flags().StringP("space", "p", string(types.L2), "Distance metric to use for the collection")
	CreateCollectionCommand.Flags().IntP("m", "m", 16, "hnsw:m - The maximum number of outgoing connections (links) for a single node within the HNSW graph.")
	CreateCollectionCommand.Flags().IntP("construction-ef", "u", 100, "hnsw:construction_ef - This parameter influences the size of the dynamic list used during the graph construction phase.")
	CreateCollectionCommand.Flags().IntP("search-ef", "f", 10, "hnsw:search_ef - The size of the dynamic list employed during the search phase.")
	CreateCollectionCommand.Flags().IntP("batch-size", "b", 100, "hnsw:batch_size - The number of elements held in brute force index (in-memory), before adding them to the HNSW index.")
	CreateCollectionCommand.Flags().IntP("sync-threshold", "k", 1000, "hnsw:sync_threshold - The number of elements added to the HNSW index before the index is synced to disk.")
	CreateCollectionCommand.Flags().IntP("threads", "n", -1, "hnsw:threads - The number of threads to use during index construction and searches. Defaults to the number of logical cores on the machine.")
	CreateCollectionCommand.Flags().Float32P("resize-factor", "r", 1.2, "hnsw:resize_factor - This parameter is used by HNSW's hierarchical layers during insertion..")
	CreateCollectionCommand.Flags().StringSliceVarP(&metaSlice, "meta", "a", []string{}, "Defines a single key-value attribute (KVP) to added to collection metadata.")
	collectionCommand.AddCommand(CreateCollectionCommand)
	RootCmd.AddCommand(CreateCollectionCommand)
	DeleteCollectionCommand.Flags().StringP("alias", "s", "", "Server alias name. If not provided, the active server will be used.")
	collectionCommand.AddCommand(DeleteCollectionCommand)
	RootCmd.AddCommand(DeleteCollectionCommand)
	CloneCollectionCommand.Flags().IntP("clone-batch-size", "z", 100, "The batch size for cloning from one collection to another.")
	CloneCollectionCommand.Flags().StringP("alias", "s", "", "Server alias name. If not provided, the active server will be used.")
	CloneCollectionCommand.Flags().StringP("space", "p", string(types.L2), "Distance metric to use for the collection")
	CloneCollectionCommand.Flags().IntP("m", "m", 16, "hnsw:m - The maximum number of outgoing connections (links) for a single node within the HNSW graph.")
	CloneCollectionCommand.Flags().IntP("construction-ef", "u", 100, "hnsw:construction_ef - This parameter influences the size of the dynamic list used during the graph construction phase.")
	CloneCollectionCommand.Flags().IntP("search-ef", "f", 10, "hnsw:search_ef - The size of the dynamic list employed during the search phase.")
	CloneCollectionCommand.Flags().IntP("batch-size", "b", 100, "hnsw:batch_size - The number of elements held in brute force index (in-memory), before adding them to the HNSW index.")
	CloneCollectionCommand.Flags().IntP("sync-threshold", "k", 1000, "hnsw:sync_threshold - The number of elements added to the HNSW index before the index is synced to disk.")
	CloneCollectionCommand.Flags().IntP("threads", "n", -1, "hnsw:threads - The number of threads to use during index construction and searches. Defaults to the number of logical cores on the machine.")
	CloneCollectionCommand.Flags().Float32P("resize-factor", "r", 1.2, "hnsw:resize_factor - This parameter is used by HNSW's hierarchical layers during insertion..")
	CloneCollectionCommand.Flags().StringP("embedding-function", "e", "", "The name of the embedding function to use for the target collection")
	CloneCollectionCommand.Flags().StringSliceVarP(&metaSlice, "meta", "a", []string{}, "Defines a single key-value attribute (KVP) to added to collection metadata.")
	RootCmd.AddCommand(CloneCollectionCommand)
}
