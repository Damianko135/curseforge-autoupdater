package config

import (
	"bytes"
	"fmt"
	"text/template"
)

// DefaultConfigTemplate is the default configuration template
const DefaultConfigTemplate = `# CurseForge Auto-Update Configuration
# This file contains the main configuration for the CurseForge Auto-Update CLI tool

# ============================================================================
# API Configuration
# ============================================================================
# Your CurseForge API key (required)
# Get yours at: https://console.curseforge.com/
api_key = "{{.APIKey}}"

# ============================================================================
# Modpack Configuration
# ============================================================================
# The CurseForge modpack ID to track
modpack_id = {{.ModpackID}}

# Target Minecraft version
game_version = "{{.GameVersion}}"

# ============================================================================
# Server Configuration
# ============================================================================
# Path to your Minecraft server directory
server_path = "{{.ServerPath}}"

# Path where backups will be stored
backup_path = "{{.BackupPath}}"

# Name of the server JAR file
server_jar_name = "{{.ServerJarName}}"

# ============================================================================
# Update Configuration
# ============================================================================
# Enable automatic updates (be careful with this!)
auto_update = {{.AutoUpdate}}

# Update channel: stable, beta, alpha
update_channel = "{{.UpdateChannel}}"

# ============================================================================
# Logging Configuration
# ============================================================================
# Log level: debug, info, warn, error
log_level = "{{.LogLevel}}"

# Log file path (empty for stdout only)
log_file = "{{.LogFile}}"

# ============================================================================
# Notification Configuration
# ============================================================================
[notifications.discord]
# Enable Discord notifications
enabled = {{.Notifications.Discord.Enabled}}

# Discord webhook URL
webhook_url = "{{.Notifications.Discord.WebhookURL}}"

# Discord channel ID (optional)
channel_id = "{{.Notifications.Discord.ChannelID}}"

# Bot username for notifications
username = "{{.Notifications.Discord.Username}}"

# Bot avatar URL (optional)
avatar_url = "{{.Notifications.Discord.AvatarURL}}"

[notifications.webhook]
# Enable generic webhook notifications
enabled = {{.Notifications.Webhook.Enabled}}

# Webhook URL
url = "{{.Notifications.Webhook.URL}}"

# HTTP method (GET, POST, PUT, etc.)
method = "{{.Notifications.Webhook.Method}}"

# Content type
content_type = "{{.Notifications.Webhook.ContentType}}"

# Request timeout
timeout = "{{.Notifications.Webhook.Timeout}}"

# Custom headers (optional)
{{if .Notifications.Webhook.Headers}}
[notifications.webhook.headers]
{{range $key, $value := .Notifications.Webhook.Headers}}
"{{$key}}" = "{{$value}}"
{{end}}
{{end}}
`

// ServerConfigTemplate is the server-specific configuration template
const ServerConfigTemplate = `# Server-specific Configuration
# This file contains server-specific settings

# ============================================================================
# Server Settings
# ============================================================================
[server]
# Server name (for display purposes)
name = "{{.Name}}"

# Server port
port = {{.Port}}

# Maximum number of players
max_players = {{.MaxPlayers}}

# Shutdown timeout (how long to wait for graceful shutdown)
shutdown_timeout = "{{.ShutdownTimeout}}"

# Custom start command (optional)
start_command = "{{.StartCommand}}"

# Custom stop command (optional)
stop_command = "{{.StopCommand}}"

# ============================================================================
# Backup Settings
# ============================================================================
[backup]
# Number of days to retain backups
retention_days = {{.RetentionDays}}

# Enable compression for backups
compression = {{.Compression}}

# Enable incremental backups
incremental = {{.Incremental}}

# ============================================================================
# Maintenance Window
# ============================================================================
[maintenance]
# Maintenance window start time (HH:MM format)
window_start = "{{.WindowStart}}"

# Maintenance window end time (HH:MM format)
window_end = "{{.WindowEnd}}"

# Timezone for maintenance window
timezone = "{{.Timezone}}"
`

// GenerateDefaultConfig generates a default configuration with provided values
func GenerateDefaultConfig(config *Config) (string, error) {
	tmpl, err := template.New("config").Parse(DefaultConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse config template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return "", fmt.Errorf("failed to execute config template: %w", err)
	}

	return buf.String(), nil
}

// GenerateServerConfig generates a server configuration with provided values
func GenerateServerConfig(serverConfig *ServerConfig, backupConfig *BackupConfig, maintenanceConfig *MaintenanceConfig) (string, error) {
	data := struct {
		*ServerConfig
		*BackupConfig
		*MaintenanceConfig
	}{
		ServerConfig:      serverConfig,
		BackupConfig:      backupConfig,
		MaintenanceConfig: maintenanceConfig,
	}

	tmpl, err := template.New("server").Parse(ServerConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse server config template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute server config template: %w", err)
	}

	return buf.String(), nil
}

// GetDefaultConfig returns a default configuration with sensible defaults
func GetDefaultConfig() *Config {
	return &Config{
		APIKey:        "your-api-key-here",
		ModpackID:     0,
		GameVersion:   "1.20.1",
		ServerPath:    "./server",
		BackupPath:    "./backups",
		ServerJarName: "server.jar",
		AutoUpdate:    false,
		UpdateChannel: "stable",
		LogLevel:      "info",
		LogFile:       "",
		Notifications: NotificationConfig{
			Discord: DiscordConfig{
				Enabled:   false,
				Username:  "CurseForge Auto-Updater",
				AvatarURL: "",
			},
			Webhook: WebhookConfig{
				Enabled:     false,
				Method:      "POST",
				ContentType: "application/json",
				Timeout:     30000000000, // 30 seconds in nanoseconds
			},
		},
	}
}

// GetDefaultServerConfig returns a default server configuration
func GetDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Name:            "Minecraft Server",
		Port:            25565,
		MaxPlayers:      20,
		ShutdownTimeout: 30000000000, // 30 seconds in nanoseconds
		StartCommand:    "",
		StopCommand:     "",
	}
}

// GetDefaultBackupConfig returns a default backup configuration
func GetDefaultBackupConfig() *BackupConfig {
	return &BackupConfig{
		RetentionDays: 30,
		Compression:   true,
		Incremental:   true,
	}
}

// GetDefaultMaintenanceConfig returns a default maintenance configuration
func GetDefaultMaintenanceConfig() *MaintenanceConfig {
	return &MaintenanceConfig{
		WindowStart: "02:00",
		WindowEnd:   "04:00",
		Timezone:    "UTC",
	}
}
