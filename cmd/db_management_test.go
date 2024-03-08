package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTenant(t *testing.T) {
	command := rootCmd

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
