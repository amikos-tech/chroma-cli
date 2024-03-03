package server

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateHost(t *testing.T) {
	t.Run("Valid Host", func(t *testing.T) {
		err := validateHost("localhost")
		require.NoError(t, err, "Unexpected error: %v", err)
	})
	t.Run("Invalid Host", func(t *testing.T) {
		err := validateHost("localhost:8080")
		require.Error(t, err)
	})

	t.Run("Valid FQDN Host", func(t *testing.T) {
		err := validateHost("api.trychroma.com")
		require.NoError(t, err, "Unexpected error: %v", err)
	})

	t.Run("Invalid FQDN Host", func(t *testing.T) {
		err := validateHost("api.trychroma.com:8080")
		require.Error(t, err)
	})

	t.Run("Empty Host", func(t *testing.T) {
		err := validateHost("")
		require.Error(t, err)
	})

	t.Run("Valid IPv4", func(t *testing.T) {
		err := validateHost("10.10.10.10")
		require.NoError(t, err, "Unexpected error: %v", err)
	})

	t.Run("Invalid IPv4", func(t *testing.T) {
		err := validateHost("10.10.10.256")
		require.Error(t, err)
	})
	t.Run("Invalid FQDN", func(t *testing.T) {
		err := validateHost("1231.com")
		require.Error(t, err)
	})
}
