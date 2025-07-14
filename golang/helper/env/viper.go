package env

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(configPath string) error {
	ext := strings.TrimPrefix(filepath.Ext(configPath), ".")

	// If configPath has an extension, treat it as a full file path
	if ext != "" {
		viper.SetConfigFile(configPath)
	} else {
		// Default fallback: assume it's a name like "config" with .toml
		viper.SetConfigName(configPath)
		ext = "toml"
	}

	viper.SetConfigType(ext)

	// Add common search paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/curseforge-autoupdater")
	viper.AddConfigPath("$HOME/.curseforge-autoupdater")

	if err := viper.ReadInConfig(); err != nil {
		// Not fatal, caller decides what to do
		return fmt.Errorf("failed to read config (%s): %w", configPath, err)
	}

	log.Printf("✅ Loaded config: %s", viper.ConfigFileUsed())
	return nil
}
