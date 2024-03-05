package cmd

import (
	"chroma/cmd/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"s"},
	Short:   "Manage Chroma servers",
	Long:    ``,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(server.AddCommand)
	serverCmd.AddCommand(server.ListCommand)
	serverCmd.AddCommand(server.RmCommand)
	rootCmd.AddCommand(server.SwitchCommand)
}
