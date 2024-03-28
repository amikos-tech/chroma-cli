package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	// Create a new command
	command := RootCmd

	t.Run("Version", func(t *testing.T) {
		client := setup()
		defer tearDown(client)
		buf := new(bytes.Buffer)
		command.SetOut(buf)
		command.SetErr(buf)
		command.SetArgs([]string{"version"})
		_, err := command.ExecuteC()
		output := buf.String()
		fmt.Printf("Output: %v\n", output)
		assert.Contains(t, output, "Chroma Server Version: ")
		assert.NoError(t, err)
	})
}
