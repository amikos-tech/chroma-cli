package cmd

import (
	"context"
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
		cmd.Printf("Tenant created: %v\n", tenantName)
	},
}

var TenantCommand = &cobra.Command{
	Use:     "tenant",
	Aliases: []string{"t"},
	Short:   "Tenant management",
}

func init() {
	CreateTenantCommand.Flags().StringVarP(&alias, "alias", "s", "", "Server alias")
	CreateTenantCommand.ValidArgs = []string{"tenant"}
	TenantCommand.AddCommand(CreateTenantCommand)
	rootCmd.AddCommand(TenantCommand)
}
