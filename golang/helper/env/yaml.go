package env

import (
	"github.com/spf13/viper"
)

// LoadYAMLConfig loads configuration from a YAML file using viper.
func LoadYAMLConfig(configName string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/curseforge-autoupdater")
	viper.AddConfigPath("$HOME/.curseforge-autoupdater")
	return viper.ReadInConfig()
}

// SaveYAMLConfig saves the current viper config to a YAML file.
func SaveYAMLConfig(filename string) error {
	return viper.WriteConfigAs(filename)
}
