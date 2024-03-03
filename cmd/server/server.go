package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultHost = "localhost"
	DefaultPort = "8000"
)

func validateHost(host string) error {
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	v := validator.New()
	hostErr := v.Var(host, "hostname")
	ipErr := v.Var(host, "ip4_addr")
	if hostErr != nil && ipErr != nil {
		return fmt.Errorf("invalid host: %v", host)
	}

	return nil
}

func getPort(changed bool) (int, error) {
	var port string
	if changed {
		port = Port
	} else {
		port = DefaultPort
	}
	var actualPort int
	// if port == "" {
	//	iPort := huh.NewInput().Value(&port).Title("Port").Placeholder(DefaultPort)
	//	err := iPort.Run()
	//	if err != nil {
	//		return -1, fmt.Errorf("unable to get port: %v", err)
	//	}
	//	if port == "" {
	//		port = DefaultPort
	//	}
	//}

	actualPort, err := strconv.Atoi(port)
	if err != nil {
		return -1, fmt.Errorf("invalid port!, must be a number! %v was given", port)
	}
	return actualPort, nil
}
func getHost(changed bool) (string, error) {
	var host string
	if changed {
		host = Host
	} else {
		host = DefaultHost
	}

	err := validateHost(host)
	if err != nil {
		return "", fmt.Errorf("invalid host: %v", err)
	}

	// if host == "" {
	//	iHost := huh.NewInput().Value(&host).Title("Host").Placeholder(DefaultHost)
	//	err := iHost.Run()
	//	if err != nil {
	//		return "", fmt.Errorf("unable to get host: %v", err)
	//	}
	// if host == "" {
	//		host = DefaultHost
	//	}
	//}

	return host, nil
}

var Host string
var Port string
var Overwrite bool
var Secure bool

var AddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add new or Update existing Chroma server",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// get the first argument tht is our alias
		alias := args[0]
		hostChanged := cmd.Flags().Changed("host")
		host, hostErr := getHost(hostChanged)
		if hostErr != nil {
			fmt.Printf("%v\n", hostErr)
			os.Exit(1)
		}
		portChanged := cmd.Flags().Changed("port")
		var actualPort, portErr = getPort(portChanged)
		if portErr != nil {
			fmt.Printf("%v\n", portErr)
			os.Exit(1)
		}
		if !hostChanged && !portChanged {
			fmt.Printf("You must specify either host or port\n")
			os.Exit(1)
		}
		if !hostChanged {
			fmt.Printf("Using default host: %v\n", DefaultHost)
		}
		if !portChanged {
			fmt.Printf("Using default port: %v\n", DefaultPort)
		}
		// confirm := false
		// if Host != "" || Port != "" {
		//	confirm = true
		//} else {
		//	err := huh.NewConfirm().
		//		Title("Is the above correct?").
		//		Affirmative("Yes!").
		//		Negative("No.").
		//		Value(&confirm).Run()
		//	if err != nil {
		//		fmt.Printf("unable to get confirmation: %v\n", err)
		//		os.Exit(1)
		//	}
		//}
		//if confirm {
		var servers = viper.GetStringMap("servers")
		if servers == nil {
			servers = make(map[string]interface{})
		}
		if !Overwrite {
			if _, ok := servers[alias]; ok {
				fmt.Printf("Server with alias %v already exists! \n", alias)
				os.Exit(1)
			}
		}
		servers[alias] = map[string]interface{}{"host": host, "port": actualPort, "secure": Secure}
		viper.Set("servers", servers)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Printf("unable to write to config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Server '%v:%v' (secure=%v) successfully added!\n", host, actualPort, Secure)
		//}
	},
}

var ForceDelete bool
var RmCommand = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Add new or Update existing Chroma server",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		var servers = viper.GetStringMap("servers")
		if servers == nil {
			servers = make(map[string]interface{})
		}
		if _, ok := servers[alias]; ok {
			confirm := ForceDelete
			if !ForceDelete {
				err := huh.NewConfirm().
					Title("Are you sure you want to remove [" + alias + "]?").
					Affirmative("Yes!").
					Negative("No.").
					Value(&confirm).Run()
				if err != nil {
					fmt.Printf("unable to get confirmation: %v\n", err)
					os.Exit(1)
				}
			}
			if !confirm {
				fmt.Printf("Operation aborted!\n")
				os.Exit(0)
			}
			delete(servers, alias)
			viper.Set("servers", servers)
			err := viper.WriteConfig()
			if err != nil {
				fmt.Printf("unable to write to config file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Server '%v' successfully removed!\n", alias)
		} else {
			fmt.Printf("Server with alias %v does not exist! \n", alias)
			os.Exit(1)
		}
	},
}

var ListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all available Chroma servers",
	Run: func(cmd *cobra.Command, args []string) {
		var servers = viper.GetStringMap("servers")
		if servers == nil {
			servers = make(map[string]interface{})
		}
		fmt.Printf("Available servers: \n")
		for alias, server := range servers {
			fmt.Printf("%v: %v\n", alias, server)
		}
	},
}

func init() {
	AddCommand.Flags().StringVarP(&Host, "host", "H", "", "Chroma server host")
	AddCommand.Flags().StringVarP(&Port, "port", "p", "", "Chroma server port")
	AddCommand.Flags().BoolVarP(&Overwrite, "overwrite", "o", false, "Overwrite existing server with the same alias")
	AddCommand.Flags().BoolVarP(&Secure, "secure", "s", false, "Use secure connection (https).")
	// AddCommand.MarkFlagsRequiredTogether("host", "port")
	AddCommand.ValidArgs = []string{"alias"}
	RmCommand.ValidArgs = []string{"alias"}
	RmCommand.Flags().BoolVarP(&ForceDelete, "force", "f", false, "Force remove server without confirmation")
}
