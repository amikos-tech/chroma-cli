package cmd

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/amikos-tech/chroma-go/collection"
	"github.com/amikos-tech/chroma-go/types"
)

var ListCollectionsCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all available collections",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient(alias)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		col, err := client.ListCollections(context.TODO())
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		for _, collection := range col {
			cmd.Printf("%v\n", collection)
		}
	},
}

var createIfNotExist bool
var space string
var m int
var constructionEf int
var searchEf int
var batchSize int
var syncThreshold int
var threads int
var alias string
var metadatas []string
var resizeFactor float32
var CreateCollectionCommand = &cobra.Command{
	Use:       "create",
	Aliases:   []string{"c"},
	Short:     "Create a new collection",
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: []string{"name"},
	Run: func(cmd *cobra.Command, args []string) {
		collectionName := args[0]
		client, err := getClient(alias)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		var options = make([]collection.Option, 0)
		options = append(options, collection.WithName(collectionName))
		if cmd.Flag("space").Changed {
			options = append(options, collection.WithHNSWDistanceFunction(types.DistanceFunction(space)))
		}
		if cmd.Flag("m").Changed {
			options = append(options, collection.WithHNSWM(int32(m)))
		}
		if cmd.Flag("construction-ef").Changed {
			options = append(options, collection.WithHNSWConstructionEf(int32(constructionEf)))
		}
		if cmd.Flag("search-ef").Changed {
			options = append(options, collection.WithHNSWSearchEf(int32(searchEf)))
		}
		if cmd.Flag("batch-size").Changed {
			options = append(options, collection.WithHNSWBatchSize(int32(batchSize)))
		}
		if cmd.Flag("sync-threshold").Changed {
			options = append(options, collection.WithHNSWSyncThreshold(int32(syncThreshold)))
		}
		if cmd.Flag("threads").Changed {
			options = append(options, collection.WithHNSWNumThreads(int32(threads)))
		}
		if cmd.Flag("ensure").Changed {
			options = append(options, collection.WithCreateIfNotExist(createIfNotExist))
		}
		if cmd.Flag("resize-factor").Changed {
			options = append(options, collection.WithHNSWResizeFactor(resizeFactor))
		}
		if cmd.Flag("meta").Changed {
			metadata := make(map[string]interface{})
			for _, meta := range metadatas {
				kvPair := strings.Split(meta, "=")
				if len(kvPair) != 2 {
					cmd.Printf("invalid metadata format: %v. should be key=value.", meta)
					os.Exit(1)
				}
				if b, err := strconv.ParseBool(kvPair[1]); err == nil {
					cmd.Printf("bool: %v\n", b)
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
			os.Exit(1)
		}
		cmd.Printf("Collection created: %v\n", collectionName)
	},
}

var DeleteCollectionCommand = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm"},
	Short:   "Delete a collection",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		collectionName := args[0]
		client, err := getClient(alias)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		_, err = client.DeleteCollection(context.TODO(), collectionName)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		cmd.Printf("Collection deleted: %v\n", collectionName)
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

	CreateCollectionCommand.Flags().String("name", "", "Name of the collection")
	CreateCollectionCommand.Flags().StringVarP(&alias, "alias", "s", "", "Server alias name. If not provided, the active server will be used.")
	CreateCollectionCommand.Flags().BoolVarP(&createIfNotExist, "ensure", "x", false, "Create collection only if it doesn't exist. Chroma will be queried before sending create, if the collection exists, exit with 0. The metadata will be overwritten.")
	CreateCollectionCommand.Flags().StringVarP(&space, "space", "p", string(types.L2), "Server alias name. If not provided, the active server will be used.")
	CreateCollectionCommand.Flags().IntVarP(&m, "m", "m", 16, "Server alias name. If not provided, the active server will be used.")
	CreateCollectionCommand.Flags().IntVarP(&constructionEf, "construction-ef", "u", 100, "Server alias name. If not provided, the active server will be used.")
	CreateCollectionCommand.Flags().IntVarP(&searchEf, "search-ef", "f", 10, "Server alias name. If not provided, the active server will be used.")
	CreateCollectionCommand.Flags().IntVarP(&batchSize, "batch-size", "b", 100, "Server alias name. If not provided, the active server will be used.")
	CreateCollectionCommand.Flags().IntVarP(&syncThreshold, "sync-threshold", "k", 1000, "Server alias name. If not provided, the active server will be used.")
	CreateCollectionCommand.Flags().IntVarP(&threads, "threads", "n", -1, "Server alias name. If not provided, the active server will be used.")
	CreateCollectionCommand.Flags().Float32VarP(&resizeFactor, "resize-factor", "r", 1.2, "Resize factor")
	CreateCollectionCommand.Flags().StringSliceVarP(&metadatas, "meta", "a", []string{}, "Server alias name. If not provided, the active server will be used.")
	collectionCommand.AddCommand(CreateCollectionCommand)
	rootCmd.AddCommand(CreateCollectionCommand)

	DeleteCollectionCommand.Flags().StringVarP(&alias, "alias", "s", "", "Server alias name. If not provided, the active server will be used.")
	collectionCommand.AddCommand(DeleteCollectionCommand)
	rootCmd.AddCommand(DeleteCollectionCommand)
}
