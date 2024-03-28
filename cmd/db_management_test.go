package cmd

import (
	"bytes"
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chroma "github.com/amikos-tech/chroma-go"
)

func helperCreateTenant(t *testing.T, client *chroma.Client, tenantName string) {
	_, err := client.CreateTenant(context.TODO(), tenantName)
	assert.NoError(t, err)
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func getRandomName(prefix string) string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return prefix + string(b)
}

func TestCreateTenant(t *testing.T) {
	command := RootCmd

	t.Run("Create tenant long", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var tenantName = "test-tenant"
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"tenant", "create", tenantName})
		_, err := command.ExecuteC()
		assert.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, tenantName)
	})

	t.Run("Create tenant short", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		var tenantName = "test-tenant"
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"t", "c", tenantName})
		_, err := command.ExecuteC()
		assert.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, tenantName)
	})
}

func TestCreateDatabase(t *testing.T) {
	command := RootCmd
	t.Run("Create db default-tenant long", func(t *testing.T) {
		err := CreateDatabaseCommand.Flag("tenant").Value.Set(DefaultTenant)
		require.NoError(t, err)
		client := setup()
		defer tearDown(client)
		var dbName = getRandomName("test-db")
		buf := new(bytes.Buffer)
		command.ResetFlags()
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"database", "create", dbName})
		_, err = command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, dbName)
		require.Contains(t, output, DefaultTenant)
	})

	t.Run("Create db custom-tenant long", func(t *testing.T) {
		err := CreateDatabaseCommand.Flag("tenant").Value.Set(DefaultTenant)
		require.NoError(t, err)
		client := setup()
		defer tearDown(client)
		var tenantName = getRandomName("test-tenant")
		helperCreateTenant(t, client, tenantName)
		var dbName = getRandomName("test-db")
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"database", "create", dbName, "--tenant", tenantName})
		_, err = command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, dbName)
		require.Contains(t, output, tenantName)
	})

	t.Run("Create db default-tenant short", func(t *testing.T) {
		err := CreateDatabaseCommand.Flag("tenant").Value.Set(DefaultTenant)
		require.NoError(t, err)
		client := setup()
		defer tearDown(client)
		var dbName = getRandomName("test-db")
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"db", "c", dbName})
		_, err = command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, dbName)
		require.Contains(t, output, DefaultTenant)
	})

	t.Run("Create db custom-tenant short", func(t *testing.T) {
		err := CreateDatabaseCommand.Flag("tenant").Value.Set(DefaultTenant)
		require.NoError(t, err)
		client := setup()
		defer tearDown(client)
		var tenantName = getRandomName("test-tenant")
		helperCreateTenant(t, client, tenantName)
		var dbName = getRandomName("test-db")
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"db", "c", dbName, "-t", tenantName})
		_, err = command.ExecuteC()
		require.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, dbName)
		require.Contains(t, output, tenantName)
	})
}
