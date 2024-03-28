package cmd

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

func setup() *chroma.Client {
	client, err := chroma.NewClient("http://localhost:8000", chroma.WithDebug(false))
	if err != nil {
		panic(err)
	}
	_, err = client.Reset(context.TODO())
	if err != nil {
		panic(err)
	}
	return client
}

func resetCloneCommandFlags() {
	_ = CloneCollectionCommand.Flag("clone-batch-size").Value.Set("1000")
	CloneCollectionCommand.Flag("clone-batch-size").Changed = false
	_ = CloneCollectionCommand.Flag("alias").Value.Set("")
	CloneCollectionCommand.Flag("alias").Changed = false
	_ = CloneCollectionCommand.Flag("space").Value.Set(string(types.L2))
	CloneCollectionCommand.Flag("space").Changed = false
	_ = CloneCollectionCommand.Flag("m").Value.Set("16")
	CloneCollectionCommand.Flag("m").Changed = false
	_ = CloneCollectionCommand.Flag("construction-ef").Value.Set("100")
	CloneCollectionCommand.Flag("construction-ef").Changed = false
	_ = CloneCollectionCommand.Flag("search-ef").Value.Set("10")
	CloneCollectionCommand.Flag("search-ef").Changed = false
	_ = CloneCollectionCommand.Flag("batch-size").Value.Set("100")
	CloneCollectionCommand.Flag("batch-size").Changed = false
	_ = CloneCollectionCommand.Flag("sync-threshold").Value.Set("1000")
	CloneCollectionCommand.Flag("sync-threshold").Changed = false
	_ = CloneCollectionCommand.Flag("threads").Value.Set("-1")
	CloneCollectionCommand.Flag("threads").Changed = false
	_ = CloneCollectionCommand.Flag("resize-factor").Value.Set("1.2")
	CloneCollectionCommand.Flag("resize-factor").Changed = false
	_ = CloneCollectionCommand.Flags().Set("meta", "")
	CloneCollectionCommand.Flag("meta").Changed = false
	metaSlice = []string{}
}

func tearDown(client *chroma.Client) {
	_, err := client.Reset(context.TODO())
	if err != nil {
		panic(err)
	}
}

func assertCollectionExists(t *testing.T, client *chroma.Client, collectionName string) *chroma.Collection {
	col, err := client.GetCollection(context.TODO(), collectionName, nil)
	require.NoError(t, err)
	require.NotNil(t, col)
	require.Equal(t, collectionName, col.Name)
	return col
}

func assertCollectionHasMetadataAttr(t *testing.T, client *chroma.Client, collectionName string, key string, value interface{}) {
	col := assertCollectionExists(t, client, collectionName)
	require.Contains(t, col.Metadata, key)
	require.Equal(t, value, col.Metadata[key])
}

func helperCreateCollection(t *testing.T, client *chroma.Client, collectionName string) *chroma.Collection {
	col, err := client.CreateCollection(context.TODO(), collectionName, nil, false, nil, types.L2)
	require.NoError(t, err)
	return col
}
func helperCreateCollectionWithMetadata(t *testing.T, client *chroma.Client, collectionName string, metadata map[string]interface{}) *chroma.Collection {
	col, err := client.CreateCollection(context.TODO(), collectionName, metadata, false, nil, types.L2)
	require.NoError(t, err)
	return col
}

func helperCreateCollectionWithMetadataAndDF(t *testing.T, client *chroma.Client, collectionName string, metadata map[string]interface{}, distanceFunction types.DistanceFunction) *chroma.Collection {
	col, err := client.CreateCollection(context.TODO(), collectionName, metadata, false, nil, distanceFunction)
	require.NoError(t, err)
	return col
}

func addDummyRecordsToCollection(t *testing.T, client *chroma.Client, collectionName string, numRecords int) {
	col, err := client.GetCollection(context.TODO(), collectionName, types.NewConsistentHashEmbeddingFunction())
	require.NoError(t, err)
	var documents = make([]string, numRecords)
	var ids = make([]string, numRecords)
	for i := 0; i < numRecords; i++ {
		documents[i] = fmt.Sprintf("record-%v", i)
		ids[i] = fmt.Sprintf("id-%v", i)
	}
	_, err = col.Add(context.TODO(), nil, nil, documents, ids)
	require.NoError(t, err)
}
func TestCreateCollectionCommand(t *testing.T) {
	command := RootCmd

	t.Run("Create Collection basic", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		assertCollectionExists(t, client, collectionName)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Distance Function (space) full flag", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		var distanceFunction = string(types.COSINE)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--space", distanceFunction})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		assertCollectionExists(t, client, collectionName)
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWSpace, distanceFunction)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Distance Function (space) shorthand flag", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		var distanceFunction = string(types.COSINE)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-p", distanceFunction})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		assertCollectionExists(t, client, collectionName)
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWSpace, distanceFunction)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Ensure flag long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--ensure"})
		c, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		require.Equal(t, expectedOutput, output)
		buf.Reset()
		_, err = c.ExecuteC() // we execute the same command again, result is idempotent
		require.NoError(t, err)
		assertCollectionExists(t, client, collectionName)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with M flag long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		var value int32 = 10
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--m", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWM, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with M flag short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		var value int32 = 11
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-m", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWM, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with ConstructionEF flag long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		var value int32 = 360
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--construction-ef", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWConstructionEF, value)
		require.Equal(t, expectedOutput, output)
	})
	t.Run("Create Collection with ConstructionEF flag short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var value int32 = 330
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-u", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWConstructionEF, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with SearchEF flag long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var value int32 = 1000
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--search-ef", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWSearchEF, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with SearchEF flag short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var value int32 = 1001
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-f", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWSearchEF, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Batch Size flag long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var value int32 = 10000
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--batch-size", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWBatchSize, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Batch Size flag short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var value int32 = 10010
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-b", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWBatchSize, value)
		require.Equal(t, expectedOutput, output)
	})
	t.Run("Create Collection with Sync Threshold flag long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var value int32 = 100000
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--sync-threshold", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWSyncThreshold, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Sync Threshold flag short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var value int32 = 90010
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-k", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWSyncThreshold, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Number of Threads flag long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var value int32 = 100000
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--threads", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWNumThreads, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Number of Threads flag short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var value int32 = 90010
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-n", strconv.Itoa(int(value))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWNumThreads, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Resize Factor flag long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var value float32 = 2.5
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--resize-factor", strconv.FormatFloat(float64(value), 'f', -1, 32)})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()

		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWResizeFactor, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Resize Factor flag short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var value float32 = 3.1
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-r", strconv.FormatFloat(float64(value), 'f', -1, 32)})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()

		assertCollectionHasMetadataAttr(t, client, collectionName, types.HNSWResizeFactor, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Custom Metadata string long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var key = "my-key"
		var value = "my-value"
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--meta", fmt.Sprintf("%v=%v", key, value)})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()

		assertCollectionHasMetadataAttr(t, client, collectionName, key, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Custom Metadata string short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var key = "my-key"
		var value = "my-value"
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-a", fmt.Sprintf("%v=%v", key, value)})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()

		assertCollectionHasMetadataAttr(t, client, collectionName, key, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Custom Metadata int long", func(t *testing.T) {
		client := setup()
		// defer tearDown(client)
		var key = "my-key"
		var value int32 = 100
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--meta", fmt.Sprintf("%v=%v", key, value)})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()

		assertCollectionHasMetadataAttr(t, client, collectionName, key, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Custom Metadata int short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var key = "my-key"
		var value int32 = 200
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-a", fmt.Sprintf("%v=%v", key, value)})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()

		assertCollectionHasMetadataAttr(t, client, collectionName, key, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Custom Metadata float long", func(t *testing.T) {
		client := setup()
		// defer tearDown(client)
		var key = "my-key"
		var value float32 = 10.123
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--meta", fmt.Sprintf("%v=%v", key, value)})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()

		assertCollectionHasMetadataAttr(t, client, collectionName, key, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Custom Metadata float short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var key = "my-key"
		var value float32 = 200.33
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-a", fmt.Sprintf("%v=%v", key, value)})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, key, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Custom Metadata boolean long", func(t *testing.T) {
		t.Skip("Skipping until Chroma collection metadata, issue is fixed")
		client := setup()
		defer tearDown(client)
		var key = "my-key"
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "--meta", fmt.Sprintf("%v=%v", key, strconv.FormatBool(true))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionHasMetadataAttr(t, client, collectionName, key, true)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Custom Metadata boolean short", func(t *testing.T) {
		t.Skip("Skipping until Chroma collection metadata, issue is fixed")
		client := setup()
		defer tearDown(client)
		var key = "my-key"
		var value = false
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-a", fmt.Sprintf("%v=%v", key, value)})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()

		assertCollectionHasMetadataAttr(t, client, collectionName, key, value)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Custom Metadata multiple entries long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		var metadata = map[string]interface{}{
			"my-string-key": "my-value",
			"my-int-key":    int32(100),
			"my-float-key":  float32(10.123),
		}
		var cmdLine = make([]string, 0)
		cmdLine = append(cmdLine, "create", collectionName)
		for k, v := range metadata {
			cmdLine = append(cmdLine, "--meta", fmt.Sprintf("%v=%v", k, v))
		}
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs(cmdLine)
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		for k, v := range metadata {
			assertCollectionHasMetadataAttr(t, client, collectionName, k, v)
		}
		require.Equal(t, expectedOutput, output)
	})

	t.Run("Create Collection with Custom Metadata multiple entries short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		var metadata = map[string]interface{}{
			"my-string-key": "my-value",
			"my-int-key":    int32(100),
			"my-float-key":  float32(10.123),
		}
		var cmdLine = make([]string, 0)
		cmdLine = append(cmdLine, "create", collectionName)
		for k, v := range metadata {
			cmdLine = append(cmdLine, "-a", fmt.Sprintf("%v=%v", k, v))
		}
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs(cmdLine)
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()

		for k, v := range metadata {
			assertCollectionHasMetadataAttr(t, client, collectionName, k, v)
		}
		require.Equal(t, expectedOutput, output)
	})
}

func TestListCollectionsCommand(t *testing.T) {
	command := RootCmd

	t.Run("List Collections long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		helperCreateCollection(t, client, collectionName)
		assertCollectionExists(t, client, collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"list"})
		_, err := command.ExecuteC()
		assert.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, collectionName)
	})

	t.Run("List Collections short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		helperCreateCollection(t, client, collectionName)
		assertCollectionExists(t, client, collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"ls"})
		_, err := command.ExecuteC()
		assert.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, collectionName)
	})
}

func TestDeleteCollectionCommand(t *testing.T) {
	command := RootCmd

	t.Run("Delete Collection long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		helperCreateCollection(t, client, collectionName)
		assertCollectionExists(t, client, collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"delete", collectionName})
		_, err := command.ExecuteC()
		assert.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, collectionName)
	})

	t.Run("Delete Collection short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		helperCreateCollection(t, client, collectionName)
		assertCollectionExists(t, client, collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"rm", collectionName})
		_, err := command.ExecuteC()
		assert.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, collectionName)
	})
}

func TestCloneCollectionCommand(t *testing.T) {
	command := RootCmd

	t.Run("Clone Collection long", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection short", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 12)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"cp", sourceCollectionName, targetCollectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "12")
	})

	t.Run("Clone Collection with m", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)

		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-m", "128"})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWM, int32(128))
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with construction_ef", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-u", "512"})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWConstructionEF, int32(512))
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with space", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-p", "cosine"})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWSpace, string(types.COSINE))
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})
	t.Run("Clone Collection with batch-size", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-b", "1000"})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWBatchSize, int32(1000))
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})
	t.Run("Clone Collection with search-ef", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-f", "99"})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWSearchEF, int32(99))
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with sync-threshold", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-k", "9991"})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWSyncThreshold, int32(9991))
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with threads", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-n", "24"})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWNumThreads, int32(24))
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with resize-factor", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-r", "2.5"})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWResizeFactor, float32(2.5))
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with metadata", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-a", "k=100", "-a", "my-key=my-value"})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, "k", int32(100))
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, "my-key", "my-value")
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with openai ef", func(t *testing.T) {
		_ = godotenv.Load("../.env")
		if os.Getenv("OPENAI_API_KEY") == "" {
			t.Skip("Skipping test as OPENAI_API_KEY is not set")
		}
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		helperCreateCollection(t, client, sourceCollectionName)
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-e", "openai"})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})
}
func TestCloneCollectionCommandWithSource(t *testing.T) {
	command := RootCmd
	t.Run("Clone Collection with source with m", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var m = int32(128)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWM: m})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWM, m)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with space", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var spaceVar = string(types.IP)
		helperCreateCollectionWithMetadataAndDF(t, client, sourceCollectionName, map[string]interface{}{types.HNSWSpace: spaceVar}, types.DistanceFunction(spaceVar))
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWSpace, spaceVar)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with construction-ef", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var cef = int32(512)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWConstructionEF: cef})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWConstructionEF, cef)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with search-ef", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var sef = int32(512)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWSearchEF: sef})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWSearchEF, sef)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with batch-size", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var bs = int32(512)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWBatchSize: bs})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWBatchSize, bs)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with sync-threshold", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var st = int32(1001)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWSyncThreshold: st})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWSyncThreshold, st)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with threads", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var threadsVar = int32(45)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWNumThreads: threadsVar})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWNumThreads, threadsVar)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with resize-factor", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var rf = float32(3.1)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWResizeFactor: rf})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWResizeFactor, rf)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with meta", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var kVal = int32(101)
		var sVal = "my-value"
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{"k": kVal, "strval": sVal})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, "k", kVal)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, "strval", sVal)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})
}

func TestCloneCollectionCommandWithSourceOverride(t *testing.T) {
	command := RootCmd
	t.Run("Clone Collection with source with m override in target", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var sourceM = int32(128)
		var targetM = int32(256)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWM: sourceM})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-m", strconv.Itoa(int(targetM))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWM, targetM)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with space override in target", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var sourceSpace = string(types.IP)
		var targetSpace = string(types.COSINE)
		helperCreateCollectionWithMetadataAndDF(t, client, sourceCollectionName, map[string]interface{}{types.HNSWSpace: sourceSpace}, types.DistanceFunction(sourceSpace))
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-p", targetSpace})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWSpace, targetSpace)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with construction-ef override in target", func(t *testing.T) {
		resetCloneCommandFlags()
		client := setup()
		defer tearDown(client)
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var sourceCef = int32(512)
		var targetCef = int32(1024)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWConstructionEF: sourceCef})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-u", strconv.Itoa(int(targetCef))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWConstructionEF, targetCef)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with search-ef override in target", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var sourceSef = int32(512)
		var targetSef = int32(512)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWSearchEF: sourceSef})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-f", strconv.Itoa(int(targetSef))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWSearchEF, targetSef)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with batch-size override in target", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var sourceBs = int32(512)
		var targetBs = int32(256)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWBatchSize: sourceBs})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-b", strconv.Itoa(int(targetBs))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWBatchSize, targetBs)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with sync-threshold override in target", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var sourceSt = int32(1001)
		var targetSt = int32(2001)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWSyncThreshold: sourceSt})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-k", strconv.Itoa(int(targetSt))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWSyncThreshold, targetSt)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with threads override in target", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var sourceThreads = int32(45)
		var targetThreads = int32(80)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWNumThreads: sourceThreads})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-n", strconv.Itoa(int(targetThreads))})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWNumThreads, targetThreads)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with resize-factor override in target", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var sourceRf = float32(3.1)
		var targetRf = float32(5.3)
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{types.HNSWResizeFactor: sourceRf})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-r", strconv.FormatFloat(float64(targetRf), 'f', -1, 32)})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, types.HNSWResizeFactor, targetRf)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})

	t.Run("Clone Collection with source with meta override in target", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		resetCloneCommandFlags()
		var sourceCollectionName = "my-new-collection" + strconv.Itoa(rand.Int())
		var targetCollectionName = "my-new-collection-copy" + strconv.Itoa(rand.Int())
		var skVal = int32(101)
		var ssVal = "my-value"
		var tkVal = int32(202)
		var tsVal = "other-value"
		helperCreateCollectionWithMetadata(t, client, sourceCollectionName, map[string]interface{}{"k": skVal, "strval": ssVal})
		addDummyRecordsToCollection(t, client, sourceCollectionName, 10)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"clone", sourceCollectionName, targetCollectionName, "-a", "k=" + strconv.Itoa(int(tkVal)), "-a", "strval=" + tsVal})
		_, err := command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		assertCollectionExists(t, client, sourceCollectionName)
		assertCollectionExists(t, client, targetCollectionName)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, "k", tkVal)
		assertCollectionHasMetadataAttr(t, client, targetCollectionName, "strval", tsVal)
		require.Contains(t, output, "successfully cloned")
		require.Contains(t, output, sourceCollectionName)
		require.Contains(t, output, targetCollectionName)
		require.Contains(t, output, "10")
	})
}
