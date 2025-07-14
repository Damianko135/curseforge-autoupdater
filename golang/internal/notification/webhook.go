package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/damianko135/curseforge-autoupdate/golang/internal/config"
)

// WebhookNotifier handles generic webhook notifications
type WebhookNotifier struct {
	config *config.WebhookConfig
	client *http.Client
}

// NewWebhookNotifier creates a new webhook notifier
func NewWebhookNotifier(config *config.WebhookConfig) *WebhookNotifier {
	return &WebhookNotifier{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// WebhookPayload represents a generic webhook payload
type WebhookPayload struct {
	Event     string                 `json:"event"`
	Timestamp string                 `json:"timestamp"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// SendNotification sends a generic notification via webhook
func (w *WebhookNotifier) SendNotification(event, message string, data map[string]interface{}) error {
	if !w.config.Enabled {
		return nil // Skip if not enabled
	}

	payload := WebhookPayload{
		Event:     event,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   message,
		Data:      data,
	}

	return w.sendWebhook(payload)
}

// SendUpdateNotification sends an update notification via webhook
func (w *WebhookNotifier) SendUpdateNotification(modpackName, currentVersion, newVersion, changelog string) error {
	data := map[string]interface{}{
		"modpack_name":    modpackName,
		"current_version": currentVersion,
		"new_version":     newVersion,
		"changelog":       changelog,
	}

	message := fmt.Sprintf("Modpack update available: %s (%s -> %s)", modpackName, currentVersion, newVersion)
	return w.SendNotification("update_available", message, data)
}

// SendUpdateStartNotification sends a notification when update starts
func (w *WebhookNotifier) SendUpdateStartNotification(modpackName, version string) error {
	data := map[string]interface{}{
		"modpack_name": modpackName,
		"version":      version,
	}

	message := fmt.Sprintf("Starting update: %s to version %s", modpackName, version)
	return w.SendNotification("update_started", message, data)
}

// SendUpdateSuccessNotification sends a notification when update succeeds
func (w *WebhookNotifier) SendUpdateSuccessNotification(modpackName, version string, duration time.Duration) error {
	data := map[string]interface{}{
		"modpack_name": modpackName,
		"version":      version,
		"duration":     duration.String(),
	}

	message := fmt.Sprintf("Update completed successfully: %s updated to version %s", modpackName, version)
	return w.SendNotification("update_success", message, data)
}

// SendUpdateFailureNotification sends a notification when update fails
func (w *WebhookNotifier) SendUpdateFailureNotification(modpackName, version string, errorMsg string) error {
	data := map[string]interface{}{
		"modpack_name": modpackName,
		"version":      version,
		"error":        errorMsg,
	}

	message := fmt.Sprintf("Update failed: %s to version %s - %s", modpackName, version, errorMsg)
	return w.SendNotification("update_failed", message, data)
}

// SendBackupNotification sends a backup notification
func (w *WebhookNotifier) SendBackupNotification(action, backupName string, size int64) error {
	data := map[string]interface{}{
		"action":      action,
		"backup_name": backupName,
		"size":        size,
	}

	message := fmt.Sprintf("Backup %s: %s", action, backupName)
	return w.SendNotification("backup_"+action, message, data)
}

// SendServerStatusNotification sends a server status notification
func (w *WebhookNotifier) SendServerStatusNotification(status, message string) error {
	data := map[string]interface{}{
		"status": status,
	}

	return w.SendNotification("server_status", message, data)
}

// sendWebhook sends a webhook payload
func (w *WebhookNotifier) sendWebhook(payload WebhookPayload) error {
	if w.config.URL == "" {
		return fmt.Errorf("webhook URL is not configured")
	}

	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Create request
	req, err := http.NewRequest(w.config.Method, w.config.URL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", w.config.ContentType)
	req.Header.Set("User-Agent", "CurseForge Auto-Updater/1.0")

	// Set custom headers
	for key, value := range w.config.Headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status code: %d", resp.StatusCode)
	}

	return nil
}

// TestConnection tests the webhook connection
func (w *WebhookNotifier) TestConnection() error {
	if !w.config.Enabled {
		return fmt.Errorf("webhook notifications are not enabled")
	}

	testData := map[string]interface{}{
		"test": true,
	}

	return w.SendNotification("test", "This is a test notification from CurseForge Auto-Updater", testData)
}

// SendCustomNotification sends a custom notification with full control over the payload
func (w *WebhookNotifier) SendCustomNotification(payload interface{}) error {
	if !w.config.Enabled {
		return nil // Skip if not enabled
	}

	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal custom webhook payload: %w", err)
	}

	// Create request
	req, err := http.NewRequest(w.config.Method, w.config.URL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create custom webhook request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", w.config.ContentType)
	req.Header.Set("User-Agent", "CurseForge Auto-Updater/1.0")

	// Set custom headers
	for key, value := range w.config.Headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send custom webhook: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("custom webhook returned status code: %d", resp.StatusCode)
	}

	return nil
}

// SendSlackNotification sends a Slack-compatible notification
func (w *WebhookNotifier) SendSlackNotification(text, channel, username, iconEmoji string, attachments []SlackAttachment) error {
	payload := SlackPayload{
		Text:        text,
		Channel:     channel,
		Username:    username,
		IconEmoji:   iconEmoji,
		Attachments: attachments,
	}

	return w.SendCustomNotification(payload)
}

// SlackPayload represents a Slack webhook payload
type SlackPayload struct {
	Text        string            `json:"text"`
	Channel     string            `json:"channel,omitempty"`
	Username    string            `json:"username,omitempty"`
	IconEmoji   string            `json:"icon_emoji,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

// SlackAttachment represents a Slack attachment
type SlackAttachment struct {
	Color      string       `json:"color,omitempty"`
	Title      string       `json:"title,omitempty"`
	Text       string       `json:"text,omitempty"`
	Fields     []SlackField `json:"fields,omitempty"`
	Footer     string       `json:"footer,omitempty"`
	Timestamp  int64        `json:"ts,omitempty"`
	MarkdownIn []string     `json:"mrkdwn_in,omitempty"`
}

// SlackField represents a field in a Slack attachment
type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

// SendTeamsNotification sends a Microsoft Teams-compatible notification
func (w *WebhookNotifier) SendTeamsNotification(title, text, themeColor string, sections []TeamsSection) error {
	payload := TeamsPayload{
		Type:       "MessageCard",
		Context:    "https://schema.org/extensions",
		Title:      title,
		Text:       text,
		ThemeColor: themeColor,
		Sections:   sections,
	}

	return w.SendCustomNotification(payload)
}

// TeamsPayload represents a Microsoft Teams webhook payload
type TeamsPayload struct {
	Type       string         `json:"@type"`
	Context    string         `json:"@context"`
	Title      string         `json:"title"`
	Text       string         `json:"text"`
	ThemeColor string         `json:"themeColor,omitempty"`
	Sections   []TeamsSection `json:"sections,omitempty"`
}

// TeamsSection represents a section in a Teams message
type TeamsSection struct {
	ActivityTitle    string      `json:"activityTitle,omitempty"`
	ActivitySubtitle string      `json:"activitySubtitle,omitempty"`
	ActivityImage    string      `json:"activityImage,omitempty"`
	Facts            []TeamsFact `json:"facts,omitempty"`
	Text             string      `json:"text,omitempty"`
	Markdown         bool        `json:"markdown,omitempty"`
}

// TeamsFact represents a fact in a Teams section
type TeamsFact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// BuildSlackUpdateNotification builds a Slack notification for updates
func BuildSlackUpdateNotification(modpackName, currentVersion, newVersion, changelog string) SlackPayload {
	attachment := SlackAttachment{
		Color: "warning",
		Title: fmt.Sprintf("🔄 Modpack Update Available: %s", modpackName),
		Text:  fmt.Sprintf("A new version of *%s* is available!", modpackName),
		Fields: []SlackField{
			{
				Title: "Current Version",
				Value: currentVersion,
				Short: true,
			},
			{
				Title: "New Version",
				Value: newVersion,
				Short: true,
			},
		},
		Footer:    "CurseForge Auto-Updater",
		Timestamp: time.Now().Unix(),
	}

	if changelog != "" {
		attachment.Fields = append(attachment.Fields, SlackField{
			Title: "Changelog",
			Value: changelog,
			Short: false,
		})
	}

	return SlackPayload{
		Text:        fmt.Sprintf("Modpack update available: %s", modpackName),
		Username:    "CurseForge Auto-Updater",
		IconEmoji:   ":arrow_up:",
		Attachments: []SlackAttachment{attachment},
	}
}

// BuildTeamsUpdateNotification builds a Teams notification for updates
func BuildTeamsUpdateNotification(modpackName, currentVersion, newVersion, changelog string) TeamsPayload {
	section := TeamsSection{
		ActivityTitle:    fmt.Sprintf("🔄 Modpack Update Available: %s", modpackName),
		ActivitySubtitle: fmt.Sprintf("A new version of %s is available!", modpackName),
		Facts: []TeamsFact{
			{
				Name:  "Current Version",
				Value: currentVersion,
			},
			{
				Name:  "New Version",
				Value: newVersion,
			},
		},
		Markdown: true,
	}

	if changelog != "" {
		section.Facts = append(section.Facts, TeamsFact{
			Name:  "Changelog",
			Value: changelog,
		})
	}

	return TeamsPayload{
		Type:       "MessageCard",
		Context:    "https://schema.org/extensions",
		Title:      "CurseForge Auto-Updater",
		Text:       fmt.Sprintf("Modpack update available: %s", modpackName),
		ThemeColor: "FF9900",
		Sections:   []TeamsSection{section},
	}
}

// ValidateWebhookConfig validates the webhook configuration
func ValidateWebhookConfig(config *config.WebhookConfig) error {
	if !config.Enabled {
		return nil // Skip validation if not enabled
	}

	if config.URL == "" {
		return fmt.Errorf("webhook URL is required")
	}

	if config.Method == "" {
		return fmt.Errorf("webhook method is required")
	}

	// Validate HTTP method
	validMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	isValidMethod := false
	for _, method := range validMethods {
		if strings.ToUpper(config.Method) == method {
			isValidMethod = true
			break
		}
	}

	if !isValidMethod {
		return fmt.Errorf("invalid webhook method: %s", config.Method)
	}

	if config.ContentType == "" {
		return fmt.Errorf("webhook content type is required")
	}

	if config.Timeout <= 0 {
		return fmt.Errorf("webhook timeout must be positive")
	}

	return nil
}
