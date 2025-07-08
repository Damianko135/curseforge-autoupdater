package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/damianko135/curseforge-autoupdate/golang/pkg/models"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

// LoadConfig loads configuration from various sources (flags, env, config file)
func LoadConfig() (*models.Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	k := koanf.New(".")

	// Define flags
	f := pflag.NewFlagSet("config", pflag.ContinueOnError)
	f.String("api-key", "", "CurseForge API key")
	f.String("mod-id", "", "Mod ID to check for updates")
	f.String("download-path", "./downloads", "Path to download files")
	f.Int("game-id", 432, "Game ID (default: 432 for Minecraft)")
	f.String("log-level", "info", "Log level (debug, info, warn, error)")
	f.String("config", "", "Path to config file")

	if err := f.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	// Load from config file if specified
	if configFile, _ := f.GetString("config"); configFile != "" {
		if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	// Load from environment variables
	if err := k.Load(env.Provider("CURSEFORGE_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "CURSEFORGE_"))
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Load from flags (highest priority)
	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		return nil, fmt.Errorf("failed to load flags: %w", err)
	}

	var config models.Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate required fields
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required (set CURSEFORGE_API_KEY or use --api-key)")
	}

	if config.ModID == "" {
		return nil, fmt.Errorf("Mod ID is required (set CURSEFORGE_MOD_ID or use --mod-id)")
	}

	return &config, nil
}

// GetRemainingArgs returns any remaining command line arguments after flags
func GetRemainingArgs() []string {
	f := pflag.NewFlagSet("config", pflag.ContinueOnError)
	f.String("api-key", "", "")
	f.String("mod-id", "", "")
	f.String("download-path", "", "")
	f.Int("game-id", 0, "")
	f.String("log-level", "", "")
	f.String("config", "", "")

	_ = f.Parse(os.Args[1:])
	return f.Args()
}
