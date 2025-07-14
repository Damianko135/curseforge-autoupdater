package env

import (
	"github.com/spf13/viper"
)

// LoadTOMLConfig loads configuration from a TOML file using viper.
func LoadTOMLConfig(configName string) error {
	viper.SetConfigType("toml")
	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/curseforge-autoupdater")
	viper.AddConfigPath("$HOME/.curseforge-autoupdater")
	return viper.ReadInConfig()
}

// SaveTOMLConfig saves the current viper config to a TOML file.
func SaveTOMLConfig(filename string) error {
	return viper.WriteConfigAs(filename)
}
