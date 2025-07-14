package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/damianko135/curseforge-autoupdate/golang/internal/config"
)

// DiscordNotifier handles Discord webhook notifications
type DiscordNotifier struct {
	config *config.DiscordConfig
	client *http.Client
}

// NewDiscordNotifier creates a new Discord notifier
func NewDiscordNotifier(config *config.DiscordConfig) *DiscordNotifier {
	return &DiscordNotifier{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// DiscordWebhookPayload represents the payload for Discord webhook
type DiscordWebhookPayload struct {
	Username  string         `json:"username,omitempty"`
	AvatarURL string         `json:"avatar_url,omitempty"`
	Content   string         `json:"content,omitempty"`
	Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

// DiscordEmbed represents a Discord embed
type DiscordEmbed struct {
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Color       int                 `json:"color,omitempty"`
	Fields      []DiscordEmbedField `json:"fields,omitempty"`
	Footer      *DiscordEmbedFooter `json:"footer,omitempty"`
	Timestamp   string              `json:"timestamp,omitempty"`
	URL         string              `json:"url,omitempty"`
	Thumbnail   *DiscordEmbedImage  `json:"thumbnail,omitempty"`
	Image       *DiscordEmbedImage  `json:"image,omitempty"`
}

// DiscordEmbedField represents a field in a Discord embed
type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// DiscordEmbedFooter represents the footer of a Discord embed
type DiscordEmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// DiscordEmbedImage represents an image in a Discord embed
type DiscordEmbedImage struct {
	URL string `json:"url"`
}

// Discord color constants
const (
	ColorSuccess int = 0x00FF00 // Green
	ColorWarning int = 0xFFFF00 // Yellow
	ColorError   int = 0xFF0000 // Red
	ColorInfo    int = 0x0099FF // Blue
	ColorUpdate  int = 0xFF9900 // Orange
)

// SendMessage sends a simple message to Discord
func (d *DiscordNotifier) SendMessage(message string) error {
	if !d.config.Enabled {
		return nil // Skip if not enabled
	}

	payload := DiscordWebhookPayload{
		Username:  d.config.Username,
		AvatarURL: d.config.AvatarURL,
		Content:   message,
	}

	return d.sendWebhook(payload)
}

// SendEmbed sends an embed message to Discord
func (d *DiscordNotifier) SendEmbed(embed DiscordEmbed) error {
	if !d.config.Enabled {
		return nil // Skip if not enabled
	}

	payload := DiscordWebhookPayload{
		Username:  d.config.Username,
		AvatarURL: d.config.AvatarURL,
		Embeds:    []DiscordEmbed{embed},
	}

	return d.sendWebhook(payload)
}

// SendUpdateNotification sends a modpack update notification
func (d *DiscordNotifier) SendUpdateNotification(modpackName, currentVersion, newVersion, changelog string) error {
	embed := DiscordEmbed{
		Title:       fmt.Sprintf("🔄 Modpack Update Available: %s", modpackName),
		Description: fmt.Sprintf("A new version of **%s** is available!", modpackName),
		Color:       ColorUpdate,
		Fields: []DiscordEmbedField{
			{
				Name:   "Current Version",
				Value:  currentVersion,
				Inline: true,
			},
			{
				Name:   "New Version",
				Value:  newVersion,
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  "🟡 Ready to Update",
				Inline: true,
			},
		},
		Footer: &DiscordEmbedFooter{
			Text: "CurseForge Auto-Updater",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if changelog != "" {
		embed.Fields = append(embed.Fields, DiscordEmbedField{
			Name:   "Changelog",
			Value:  truncateString(changelog, 1024),
			Inline: false,
		})
	}

	return d.SendEmbed(embed)
}

// SendUpdateStartNotification sends a notification when update starts
func (d *DiscordNotifier) SendUpdateStartNotification(modpackName, version string) error {
	embed := DiscordEmbed{
		Title:       fmt.Sprintf("⚙️ Starting Update: %s", modpackName),
		Description: fmt.Sprintf("Beginning update process for **%s** to version **%s**", modpackName, version),
		Color:       ColorInfo,
		Fields: []DiscordEmbedField{
			{
				Name:   "Status",
				Value:  "🔄 Updating...",
				Inline: true,
			},
			{
				Name:   "Action",
				Value:  "Server will be temporarily unavailable",
				Inline: true,
			},
		},
		Footer: &DiscordEmbedFooter{
			Text: "CurseForge Auto-Updater",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return d.SendEmbed(embed)
}

// SendUpdateSuccessNotification sends a notification when update succeeds
func (d *DiscordNotifier) SendUpdateSuccessNotification(modpackName, version string, duration time.Duration) error {
	embed := DiscordEmbed{
		Title:       fmt.Sprintf("✅ Update Completed: %s", modpackName),
		Description: fmt.Sprintf("**%s** has been successfully updated to version **%s**", modpackName, version),
		Color:       ColorSuccess,
		Fields: []DiscordEmbedField{
			{
				Name:   "Status",
				Value:  "🟢 Online",
				Inline: true,
			},
			{
				Name:   "Duration",
				Value:  duration.String(),
				Inline: true,
			},
			{
				Name:   "Action",
				Value:  "Server is now available",
				Inline: true,
			},
		},
		Footer: &DiscordEmbedFooter{
			Text: "CurseForge Auto-Updater",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return d.SendEmbed(embed)
}

// SendUpdateFailureNotification sends a notification when update fails
func (d *DiscordNotifier) SendUpdateFailureNotification(modpackName, version string, errorMsg string) error {
	embed := DiscordEmbed{
		Title:       fmt.Sprintf("❌ Update Failed: %s", modpackName),
		Description: fmt.Sprintf("Failed to update **%s** to version **%s**", modpackName, version),
		Color:       ColorError,
		Fields: []DiscordEmbedField{
			{
				Name:   "Status",
				Value:  "🔴 Failed",
				Inline: true,
			},
			{
				Name:   "Error",
				Value:  truncateString(errorMsg, 1024),
				Inline: false,
			},
		},
		Footer: &DiscordEmbedFooter{
			Text: "CurseForge Auto-Updater",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return d.SendEmbed(embed)
}

// SendBackupNotification sends a backup notification
func (d *DiscordNotifier) SendBackupNotification(action, backupName string, size int64) error {
	var title, description string
	var color int

	switch action {
	case "created":
		title = "💾 Backup Created"
		description = fmt.Sprintf("Backup **%s** has been created successfully", backupName)
		color = ColorSuccess
	case "restored":
		title = "🔄 Backup Restored"
		description = fmt.Sprintf("Backup **%s** has been restored successfully", backupName)
		color = ColorInfo
	case "failed":
		title = "❌ Backup Failed"
		description = fmt.Sprintf("Failed to create backup **%s**", backupName)
		color = ColorError
	default:
		title = "💾 Backup Operation"
		description = fmt.Sprintf("Backup operation for **%s**", backupName)
		color = ColorInfo
	}

	embed := DiscordEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Fields: []DiscordEmbedField{
			{
				Name:   "Backup Name",
				Value:  backupName,
				Inline: true,
			},
		},
		Footer: &DiscordEmbedFooter{
			Text: "CurseForge Auto-Updater",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if size > 0 {
		embed.Fields = append(embed.Fields, DiscordEmbedField{
			Name:   "Size",
			Value:  formatSize(size),
			Inline: true,
		})
	}

	return d.SendEmbed(embed)
}

// SendServerStatusNotification sends a server status notification
func (d *DiscordNotifier) SendServerStatusNotification(status, message string) error {
	var title string
	var color int

	switch status {
	case "starting":
		title = "🟡 Server Starting"
		color = ColorWarning
	case "online":
		title = "🟢 Server Online"
		color = ColorSuccess
	case "stopping":
		title = "🟡 Server Stopping"
		color = ColorWarning
	case "offline":
		title = "🔴 Server Offline"
		color = ColorError
	default:
		title = "ℹ️ Server Status"
		color = ColorInfo
	}

	embed := DiscordEmbed{
		Title:       title,
		Description: message,
		Color:       color,
		Footer: &DiscordEmbedFooter{
			Text: "CurseForge Auto-Updater",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return d.SendEmbed(embed)
}

// sendWebhook sends a webhook payload to Discord
func (d *DiscordNotifier) sendWebhook(payload DiscordWebhookPayload) error {
	if d.config.WebhookURL == "" {
		return fmt.Errorf("Discord webhook URL is not configured")
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord payload: %w", err)
	}

	resp, err := d.client.Post(d.config.WebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send Discord webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Discord webhook returned status code: %d", resp.StatusCode)
	}

	return nil
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// formatSize formats a byte size into human-readable format
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// TestConnection tests the Discord webhook connection
func (d *DiscordNotifier) TestConnection() error {
	if !d.config.Enabled {
		return fmt.Errorf("Discord notifications are not enabled")
	}

	testEmbed := DiscordEmbed{
		Title:       "🧪 Test Notification",
		Description: "This is a test notification from CurseForge Auto-Updater",
		Color:       ColorInfo,
		Fields: []DiscordEmbedField{
			{
				Name:   "Status",
				Value:  "✅ Connection Successful",
				Inline: true,
			},
		},
		Footer: &DiscordEmbedFooter{
			Text: "CurseForge Auto-Updater",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return d.SendEmbed(testEmbed)
}
