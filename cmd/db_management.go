package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateTenantCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a tenant",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tenantName := args[0]
		activeAlias := viper.GetString("active_server")
		alias, err := getStringFlagIfChangedWithDefault(cmd, "alias", &activeAlias)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		client, err := getClient(*alias)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		_, err = client.CreateTenant(context.TODO(), tenantName)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		cmd.Printf("Tenant '%v' created\n", tenantName)
	},
}

var tenant string // Tenant name
var CreateDatabaseCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a db for a tenant, if no tenant is specified with --tenant/-t, the default_tenant is used.",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbName := args[0]
		activeAlias := viper.GetString("active_server")
		alias, err := getStringFlagIfChangedWithDefault(cmd, "alias", &activeAlias)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		client, err := getClient(*alias)
		fmt.Printf("tenant: %v\n", tenant)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		_, err = client.CreateDatabase(context.TODO(), dbName, &tenant)
		if err != nil {
			cmd.Printf("%v\n", err)
			os.Exit(1)
		}
		cmd.Printf("Database '%v' created in tenant '%v'\n", dbName, tenant)
	},
}

var TenantCommand = &cobra.Command{
	Use:     "tenant",
	Aliases: []string{"t"},
	Short:   "Tenant management",
}

var DBCommand = &cobra.Command{
	Use:     "database",
	Aliases: []string{"db"},
	Short:   "Database management",
}

func init() {
	CreateTenantCommand.Flags().StringP("alias", "s", "", "Server alias")
	CreateTenantCommand.ValidArgs = []string{"tenant"}
	CreateDatabaseCommand.Flags().StringP("alias", "s", "", "Server alias")
	CreateDatabaseCommand.Flags().StringVarP(&tenant, "tenant", "t", DefaultTenant, "Tenant name")
	CreateDatabaseCommand.ValidArgs = []string{"db"}
	TenantCommand.AddCommand(CreateTenantCommand)
	DBCommand.AddCommand(CreateDatabaseCommand)
	RootCmd.AddCommand(TenantCommand)
	RootCmd.AddCommand(DBCommand)
}
