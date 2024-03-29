package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/amikos-tech/chroma-cli/chroma/utils"
	"github.com/charmbracelet/huh"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultHost     = "localhost"
	DefaultPort     = "8000"
	DefaultTenant   = "default_tenant"
	DefaultDatabase = "default_database"
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

func getTenant(changed bool) (string, error) {
	if !changed {
		return DefaultTenant, nil
	}
	// TODO validate tenant name
	return Tenant, nil
}

func getDatabase(changed bool) (string, error) {
	if !changed {
		return DefaultDatabase, nil
	}
	// TODO validate database name
	return Database, nil
}

const (
	EnvChromaAPIToken  = "CHROMA_API_TOKEN"
	EnvChromaXAPIToken = "CHROMA_X_API_TOKEN"
	EnvChromaBasicAuth = "CHROMA_BASIC_AUTH"
)

var Host string
var Port string
var Overwrite bool
var Secure bool
var Tenant string
var Database string
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
			cmd.Printf("%v\n", hostErr)
			os.Exit(1)
		}
		portChanged := cmd.Flags().Changed("port")
		var actualPort, portErr = getPort(portChanged)
		if portErr != nil {
			cmd.Printf("%v\n", portErr)
			os.Exit(1)
		}
		if !hostChanged && !portChanged {
			cmd.Printf("You must specify either host or port\n")
			os.Exit(1)
		}
		if !hostChanged {
			cmd.Printf("Using default host: %v\n", DefaultHost)
		}
		if !portChanged {
			cmd.Printf("Using default port: %v\n", DefaultPort)
		}
		var tenant, tenantErr = getTenant(cmd.Flags().Changed("tenant"))
		if tenantErr != nil {
			cmd.Printf("%v\n", tenantErr)
			os.Exit(1)
		}
		var database, databaseErr = getDatabase(cmd.Flags().Changed("database"))
		if databaseErr != nil {
			cmd.Printf("%v\n", databaseErr)
			os.Exit(1)
		}
		// confirm := false
		// if Host != "" || Port != "" {
		//	confirm = true
		// } else {
		//	err := huh.NewConfirm().
		//		Title("Is the above correct?").
		//		Affirmative("Yes!").
		//		Negative("No.").
		//		Value(&confirm).Run()
		//	if err != nil {
		//		cmd.Printf("unable to get confirmation: %v\n", err)
		//		os.Exit(1)
		//	}
		// }
		// if confirm {
		var servers = viper.GetStringMap("servers")
		var setActive = false
		if len(servers) == 0 {
			servers = make(map[string]interface{})
			setActive = true
		}
		if !Overwrite {
			if _, ok := servers[alias]; ok {
				cmd.Printf("Server with alias %v already exists! \n", alias)
				os.Exit(1)
			}
		}
		servers[alias] = map[string]interface{}{
			"host":     host,
			"port":     actualPort,
			"secure":   Secure,
			"tenant":   tenant,
			"database": database,
		}
		var envConfigProvided = false
		var _authType = AuthTypeNone
		var _authToken string
		if os.Getenv(EnvChromaAPIToken) != "" {
			envConfigProvided = true
			_authType = AuthTypeToken
			_authToken = os.Getenv(EnvChromaAPIToken)
		}
		if os.Getenv(EnvChromaXAPIToken) != "" {
			envConfigProvided = true
			_authType = AuthTypeXToken
			_authToken = os.Getenv(EnvChromaXAPIToken)
		}
		if os.Getenv(EnvChromaBasicAuth) != "" {
			envConfigProvided = true
			_authType = AuthTypeBasic
			_authToken = os.Getenv(EnvChromaBasicAuth)
		}
		fmt.Printf("login: %v\n", cmd.Flags().Changed("login"))
		fmt.Printf("envConfigProvided: %v\n", envConfigProvided)
		if !envConfigProvided && cmd.Flags().Changed("login") {
			err := huh.NewSelect[AuthType]().
				Title("Authorization Type").
				Options(
					huh.NewOption("Basic", AuthTypeBasic),
					huh.NewOption("Token (Authorization)", AuthTypeToken),
					huh.NewOption("Token (X-Chroma-Token)", AuthTypeXToken),
				).
				Value(&_authType).Run()
			if err != nil {
				return // TODO fix this
			}
			if _authType == AuthTypeBasic {
				err := huh.NewInput().Value(&_authToken).Title("Basic Auth").Placeholder("username:password").Run()
				if err != nil {
					cmd.Printf("unable to get basic auth: %v\n", err)
					os.Exit(1)
				}
			} else {
				err := huh.NewInput().Value(&_authToken).Title("Token").Placeholder("token").Run()
				if err != nil {
					cmd.Printf("unable to get token: %v\n", err)
					os.Exit(1)
				}
			}
		}
		if _authType != AuthTypeNone {
			var _authInfo = make(map[string]string)
			_authInfo["type"] = string(_authType)
			_authInfo["token"] = _authToken
			servers[alias].(map[string]interface{})["auth"] = _authInfo
		}
		viper.Set("servers", servers)
		err := viper.WriteConfig()
		if err != nil {
			cmd.Printf("unable to write to config file: %v\n", err)
			os.Exit(1)
		}
		if setActive {
			err := utils.SetActiveServer(alias)
			if err != nil {
				cmd.Printf("unable to write to config file: %v\n", err)
				os.Exit(1)
			}
		}
		cmd.Printf("Server '%v:%v' (secure=%v) successfully added!\n", host, actualPort, Secure)
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
					cmd.Printf("unable to get confirmation: %v\n", err)
					os.Exit(1)
				}
			}
			if !confirm {
				cmd.Printf("Operation aborted!\n")
				os.Exit(0)
			}
			delete(servers, alias)
			if viper.GetString("active") == alias {
				viper.Set("active", "")
				cmd.Println(alias, "was the active server. You will need to set a new active server.")
			}
			viper.Set("servers", servers)
			err := viper.WriteConfig()
			if err != nil {
				cmd.Printf("unable to write to config file: %v\n", err)
				os.Exit(1)
			}
			cmd.Printf("Server '%v' successfully removed!\n", alias)
		} else {
			cmd.Printf("Server with alias %v does not exist! \n", alias)
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
		cmd.Printf("Available servers: \n")
		for alias, server := range servers {
			cmd.Printf("%v: %v\n", alias, server)
		}
	},
}

var DBAndTenantDefaults bool

var UseCommand = &cobra.Command{
	Use:   "use",
	Short: "Set active server",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		err := utils.SetActiveServer(alias)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("tenant") {
			err := utils.SetActiveTenant(Tenant)
			if err != nil {
				cmd.Printf("%v\n", err)
				os.Exit(1)
			}
			cmd.Printf("Tenant '%v' set as active!\n", Tenant)
		} else if cmd.Flags().Changed("defaults") {
			getSrv, err := utils.GetServer(alias)
			if err != nil {
				cmd.Printf("%v\n", err)
				os.Exit(1)
			}
			if getSrv["tenant"] == nil {
				getSrv["tenant"] = DefaultTenant
			}
			if _, ok := getSrv["tenant"]; ok {
				err := utils.SetActiveTenant(getSrv["tenant"].(string))
				if err != nil {
					cmd.Printf("%v\n", err)
					os.Exit(1)
				}
				cmd.Printf("Tenant '%v' set as active!\n", getSrv["tenant"])
			}
		}
		if cmd.Flags().Changed("database") {
			err := utils.SetActiveDatabase(Database)
			if err != nil {
				cmd.Printf("%v\n", err)
				os.Exit(1)
			}
			cmd.Printf("Database '%v' set as active!\n", Database)
		} else if cmd.Flags().Changed("defaults") {
			getSrv, err := utils.GetServer(alias)
			if err != nil {
				cmd.Printf("%v\n", err)
				os.Exit(1)
			}
			if getSrv["database"] == nil {
				getSrv["database"] = DefaultDatabase
			}
			if _, ok := getSrv["database"]; ok {
				err := utils.SetActiveDatabase(getSrv["database"].(string))
				if err != nil {
					cmd.Printf("%v\n", err)
					os.Exit(1)
				}
			}
			cmd.Printf("Database '%v' set as active!\n", getSrv["database"])
		}
		cmd.Printf("Server '%v' set as active!\n", alias)
	},
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"s"},
	Short:   "Manage Chroma servers",
	Long:    ``,
}

type AuthType string

const (
	AuthTypeNone   AuthType = "none"
	AuthTypeBasic  AuthType = "basic"
	AuthTypeToken  AuthType = "token"
	AuthTypeXToken AuthType = "x-token"
)

func init() {
	AddCommand.Flags().StringVarP(&Host, "host", "H", "", "Chroma server host")
	AddCommand.Flags().StringVarP(&Port, "port", "p", "", "Chroma server port")
	AddCommand.Flags().BoolVarP(&Overwrite, "force", "f", false, "Overwrite existing server with the same alias")
	AddCommand.Flags().BoolVar(&Secure, "secure", false, "Use secure connection (https).")
	AddCommand.Flags().Bool("login", false, "Authenticate with the server. If the following env vars are not provided the user will be prompted to enter the login information - CHROMA_API_TOKEN or CHROMA_BASIC_AUTH.")
	AddCommand.Flags().StringVar(&Tenant, "tenant", DefaultTenant, "Default tenant for the server")
	AddCommand.Flags().StringVar(&Database, "database", DefaultDatabase, "Default database for the server")
	// AddCommand.MarkFlagsRequiredTogether("host", "port")
	AddCommand.ValidArgs = []string{"alias"}
	RmCommand.ValidArgs = []string{"alias"}
	RmCommand.Flags().BoolVarP(&ForceDelete, "force", "f", false, "Force remove server without confirmation")
	UseCommand.Flags().StringVarP(&Tenant, "tenant", "t", "", "Default tenant for the server")
	UseCommand.Flags().StringVarP(&Database, "database", "d", "", "Default database for the server")
	UseCommand.Flags().BoolVar(&DBAndTenantDefaults, "defaults", false, "Reset active tenant and database to defaults")
	UseCommand.MarkFlagsMutuallyExclusive("tenant", "defaults")
	UseCommand.MarkFlagsMutuallyExclusive("database", "defaults")
	RootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(AddCommand)
	serverCmd.AddCommand(ListCommand)
	serverCmd.AddCommand(RmCommand)
	RootCmd.AddCommand(UseCommand)
}
