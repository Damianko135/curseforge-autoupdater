package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config represents the main configuration structure
type Config struct {
	// API Configuration
	APIKey string `mapstructure:"api_key"`

	// Modpack Configuration
	ModpackID   int    `mapstructure:"modpack_id"`
	GameVersion string `mapstructure:"game_version"`

	// Server Configuration
	ServerPath    string `mapstructure:"server_path"`
	BackupPath    string `mapstructure:"backup_path"`
	ServerJarName string `mapstructure:"server_jar_name"`

	// Notification Configuration
	Notifications NotificationConfig `mapstructure:"notifications"`

	// Update Configuration
	AutoUpdate    bool   `mapstructure:"auto_update"`
	UpdateChannel string `mapstructure:"update_channel"` // stable, beta, alpha

	// Logging Configuration
	LogLevel string `mapstructure:"log_level"`
	LogFile  string `mapstructure:"log_file"`
}

// NotificationConfig holds all notification settings
type NotificationConfig struct {
	Discord DiscordConfig `mapstructure:"discord"`
	Webhook WebhookConfig `mapstructure:"webhook"`
}

// DiscordConfig holds Discord-specific notification settings
type DiscordConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	WebhookURL string `mapstructure:"webhook_url"`
	ChannelID  string `mapstructure:"channel_id"`
	Username   string `mapstructure:"username"`
	AvatarURL  string `mapstructure:"avatar_url"`
}

// WebhookConfig holds generic webhook settings
type WebhookConfig struct {
	Enabled     bool              `mapstructure:"enabled"`
	URL         string            `mapstructure:"url"`
	Headers     map[string]string `mapstructure:"headers"`
	ContentType string            `mapstructure:"content_type"`
	Method      string            `mapstructure:"method"`
	Timeout     time.Duration     `mapstructure:"timeout"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Name            string        `mapstructure:"name"`
	Port            int           `mapstructure:"port"`
	MaxPlayers      int           `mapstructure:"max_players"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	StartCommand    string        `mapstructure:"start_command"`
	StopCommand     string        `mapstructure:"stop_command"`
}

// BackupConfig holds backup-specific configuration
type BackupConfig struct {
	RetentionDays int  `mapstructure:"retention_days"`
	Compression   bool `mapstructure:"compression"`
	Incremental   bool `mapstructure:"incremental"`
}

// MaintenanceConfig holds maintenance window configuration
type MaintenanceConfig struct {
	WindowStart string `mapstructure:"window_start"`
	WindowEnd   string `mapstructure:"window_end"`
	Timezone    string `mapstructure:"timezone"`
}

// LoadConfig loads configuration from file
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Set config file path
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	v.SetConfigFile(configPath)
	v.SetConfigType("toml")

	// Read environment variables
	v.AutomaticEnv()

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %s", configPath)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// API defaults
	v.SetDefault("api_key", "")

	// Modpack defaults
	v.SetDefault("modpack_id", 0)
	v.SetDefault("game_version", "1.20.1")

	// Server defaults
	v.SetDefault("server_path", "./server")
	v.SetDefault("backup_path", "./backups")
	v.SetDefault("server_jar_name", "server.jar")

	// Update defaults
	v.SetDefault("auto_update", false)
	v.SetDefault("update_channel", "stable")

	// Logging defaults
	v.SetDefault("log_level", "info")
	v.SetDefault("log_file", "")

	// Notification defaults
	v.SetDefault("notifications.discord.enabled", false)
	v.SetDefault("notifications.discord.username", "CurseForge Auto-Updater")
	v.SetDefault("notifications.webhook.enabled", false)
	v.SetDefault("notifications.webhook.method", "POST")
	v.SetDefault("notifications.webhook.content_type", "application/json")
	v.SetDefault("notifications.webhook.timeout", "30s")
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate API key
	if config.APIKey == "" {
		return fmt.Errorf("api_key is required")
	}

	// Validate modpack ID
	if config.ModpackID <= 0 {
		return fmt.Errorf("modpack_id must be greater than 0")
	}

	// Validate paths
	if config.ServerPath == "" {
		return fmt.Errorf("server_path is required")
	}

	if config.BackupPath == "" {
		return fmt.Errorf("backup_path is required")
	}

	// Validate update channel
	validChannels := []string{"stable", "beta", "alpha"}
	isValid := false
	for _, channel := range validChannels {
		if config.UpdateChannel == channel {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("update_channel must be one of: stable, beta, alpha")
	}

	// Validate Discord config if enabled
	if config.Notifications.Discord.Enabled {
		if config.Notifications.Discord.WebhookURL == "" {
			return fmt.Errorf("discord webhook_url is required when discord notifications are enabled")
		}
	}

	// Validate webhook config if enabled
	if config.Notifications.Webhook.Enabled {
		if config.Notifications.Webhook.URL == "" {
			return fmt.Errorf("webhook url is required when webhook notifications are enabled")
		}
	}

	return nil
}

// getDefaultConfigPath returns the default configuration file path
func getDefaultConfigPath() string {
	// Look for config file in current directory first
	configFiles := []string{
		"config.toml",
		"curseforge-autoupdate.toml",
		".curseforge-autoupdate.toml",
	}

	for _, file := range configFiles {
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}

	// If not found, return default name
	return "config.toml"
}

// SaveConfig saves configuration to file
func SaveConfig(config *Config, configPath string) error {
	v := viper.New()

	// Set values from config struct
	v.Set("api_key", config.APIKey)
	v.Set("modpack_id", config.ModpackID)
	v.Set("game_version", config.GameVersion)
	v.Set("server_path", config.ServerPath)
	v.Set("backup_path", config.BackupPath)
	v.Set("server_jar_name", config.ServerJarName)
	v.Set("auto_update", config.AutoUpdate)
	v.Set("update_channel", config.UpdateChannel)
	v.Set("log_level", config.LogLevel)
	v.Set("log_file", config.LogFile)

	// Set notification config
	v.Set("notifications.discord.enabled", config.Notifications.Discord.Enabled)
	v.Set("notifications.discord.webhook_url", config.Notifications.Discord.WebhookURL)
	v.Set("notifications.discord.channel_id", config.Notifications.Discord.ChannelID)
	v.Set("notifications.discord.username", config.Notifications.Discord.Username)
	v.Set("notifications.discord.avatar_url", config.Notifications.Discord.AvatarURL)

	v.Set("notifications.webhook.enabled", config.Notifications.Webhook.Enabled)
	v.Set("notifications.webhook.url", config.Notifications.Webhook.URL)
	v.Set("notifications.webhook.headers", config.Notifications.Webhook.Headers)
	v.Set("notifications.webhook.content_type", config.Notifications.Webhook.ContentType)
	v.Set("notifications.webhook.method", config.Notifications.Webhook.Method)
	v.Set("notifications.webhook.timeout", config.Notifications.Webhook.Timeout)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
