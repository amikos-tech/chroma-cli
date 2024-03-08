package cmd

import (
	"chroma/cmd/server"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var CreateTenantCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a tenant",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tenantName := args[0]
		client, err := getClient(alias)
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
	Short:   "Create a db",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbName := args[0]
		client, err := getClient(alias)
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
	CreateTenantCommand.Flags().StringVarP(&alias, "alias", "s", "", "Server alias")
	CreateTenantCommand.ValidArgs = []string{"tenant"}
	CreateDatabaseCommand.Flags().StringVarP(&alias, "alias", "s", "", "Server alias")
	CreateDatabaseCommand.Flags().StringVarP(&tenant, "tenant", "t", server.DefaultTenant, "Tenant name")
	CreateDatabaseCommand.ValidArgs = []string{"db"}
	TenantCommand.AddCommand(CreateTenantCommand)
	DBCommand.AddCommand(CreateDatabaseCommand)
	rootCmd.AddCommand(TenantCommand)
	rootCmd.AddCommand(DBCommand)
}
