package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetServer(alias string) (map[string]interface{}, error) {
	var servers = viper.GetStringMap("servers")
	if servers == nil {
		servers = make(map[string]interface{})
	}
	if server, ok := servers[alias]; ok {
		return server.(map[string]interface{}), nil
	}
	return nil, fmt.Errorf("server with alias %v does not exist", alias)
}

// SetActiveServer sets the active server to the one with the given alias
func SetActiveServer(alias string) error {
	var servers = viper.GetStringMap("servers")
	if servers == nil {
		servers = make(map[string]interface{})
	}
	if _, ok := servers[alias]; ok {
		viper.Set("active_server", alias)
		err := viper.WriteConfig()
		if err != nil {
			return fmt.Errorf("unable to write to config file: %v", err)
		}
	} else {
		return fmt.Errorf("server with alias %v does not exist", alias)
	}
	return nil
}

// SetActiveDatabase sets the active database to the one with the given name
func SetActiveDatabase(database string) error {
	viper.Set("active_db", database)
	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("unable to write to config file: %v", err)
	}
	return nil
}

// SetActiveTenant sets the active tenant to the one with the given name
func SetActiveTenant(tenant string) error {
	viper.Set("active_tenant", tenant)
	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("unable to write to config file: %v", err)
	}
	return nil
}
