package server

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

const (
	DefaultHost = "localhost"
	DefaultPort = "8000"
)

// chroma server add -h localhost -p 8080

var AddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a new server",

	Run: func(cmd *cobra.Command, args []string) {
		iHost := huh.NewInput()
		host := ""
		port := DefaultPort
		iHost.Title("Host").Placeholder(host).CharLimit(256).Value(&host).Placeholder(DefaultHost)
		err := iHost.Run()
		if err != nil {
			return
		}
		if host == "" {
			host = DefaultHost
		}
		iPort := huh.NewInput().Value(&port).Title("Port").Placeholder(DefaultPort)
		err = iPort.Run()
		if err != nil {
			return
		}
		if port == "" {
			port = DefaultPort
		}

		fmt.Printf("Adding server: %v:%v \n", host, port)
		confirm := false
		err = huh.NewConfirm().
			Title("Is the above correct?").
			Affirmative("Yes!").
			Negative("No.").
			Value(&confirm).Run()
		if err != nil {
			return
		}
	},
}
