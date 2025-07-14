package env

import (
	"github.com/spf13/viper"
)

// LoadJSONConfig loads configuration from a JSON file using viper.
func LoadJSONConfig(configName string) error {
	viper.SetConfigType("json")
	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/curseforge-autoupdater")
	viper.AddConfigPath("$HOME/.curseforge-autoupdater")
	return viper.ReadInConfig()
}

// SaveJSONConfig saves the current viper config to a JSON file.
func SaveJSONConfig(filename string) error {
	return viper.WriteConfigAs(filename)
}
