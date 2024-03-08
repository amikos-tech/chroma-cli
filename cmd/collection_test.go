package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

func setup() *chroma.Client {
	client, err := chroma.NewClient("http://localhost:8000", chroma.WithDebug(true))
	if err != nil {
		panic(err)
	}
	// _, err = client.Reset(context.TODO())
	// if err != nil {
	//	panic(err)
	//}
	return client
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

func TestCreateCollectionCommand(t *testing.T) {
	command := rootCmd

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

	t.Run("Create Collection with Ensure flag short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var collectionName = "my-new-collection"
		expectedOutput := fmt.Sprintf("Collection created: %v\n", collectionName)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"create", collectionName, "-x"})
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
	command := rootCmd

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
	command := rootCmd

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
