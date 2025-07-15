package env

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(configPath string) error {
	// If configPath has an extension, treat it as a file path (absolute or relative)
	ext := strings.TrimPrefix(filepath.Ext(configPath), ".")
	if ext != "" {
		// If the path is not absolute, make it relative to the current working directory
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			return fmt.Errorf("could not resolve config path: %w", err)
		}
		viper.SetConfigFile(absPath)
		viper.SetConfigType(ext)
	} else {
		// Treat as config name (no extension), search in current and standard locations
		viper.SetConfigName(configPath)
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/curseforge-autoupdater")
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(filepath.Join(home, ".curseforge-autoupdater"))
		} else {
			log.Printf("⚠️ Could not resolve user home directory: %v", err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		// Not fatal, caller decides what to do
		return fmt.Errorf("failed to read config (%s): %w", configPath, err)
	}

	log.Printf("✅ Loaded config: %s", viper.ConfigFileUsed())
	return nil
}
