package notification

import (
	"fmt"
	"sync"
	"time"

	"github.com/damianko135/curseforge-autoupdate/golang/internal/config"
)

// Manager handles all notification channels
type Manager struct {
	discord *DiscordNotifier
	webhook *WebhookNotifier
	enabled bool
	mu      sync.RWMutex
}

// NewManager creates a new notification manager
func NewManager(config *config.NotificationConfig) *Manager {
	var discord *DiscordNotifier
	var webhook *WebhookNotifier

	if config.Discord.Enabled {
		discord = NewDiscordNotifier(&config.Discord)
	}

	if config.Webhook.Enabled {
		webhook = NewWebhookNotifier(&config.Webhook)
	}

	enabled := config.Discord.Enabled || config.Webhook.Enabled

	return &Manager{
		discord: discord,
		webhook: webhook,
		enabled: enabled,
	}
}

// IsEnabled returns whether notifications are enabled
func (m *Manager) IsEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enabled
}

// SendMessage sends a simple message to all enabled channels
func (m *Manager) SendMessage(message string) error {
	if !m.IsEnabled() {
		return nil
	}

	var errors []error

	// Send to Discord
	if m.discord != nil {
		if err := m.discord.SendMessage(message); err != nil {
			errors = append(errors, fmt.Errorf("Discord: %w", err))
		}
	}

	// Send to webhook
	if m.webhook != nil {
		if err := m.webhook.SendNotification("message", message, nil); err != nil {
			errors = append(errors, fmt.Errorf("Webhook: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %v", errors)
	}

	return nil
}

// SendUpdateNotification sends an update notification to all enabled channels
func (m *Manager) SendUpdateNotification(modpackName, currentVersion, newVersion, changelog string) error {
	if !m.IsEnabled() {
		return nil
	}

	var errors []error

	// Send to Discord
	if m.discord != nil {
		if err := m.discord.SendUpdateNotification(modpackName, currentVersion, newVersion, changelog); err != nil {
			errors = append(errors, fmt.Errorf("Discord: %w", err))
		}
	}

	// Send to webhook
	if m.webhook != nil {
		if err := m.webhook.SendUpdateNotification(modpackName, currentVersion, newVersion, changelog); err != nil {
			errors = append(errors, fmt.Errorf("Webhook: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %v", errors)
	}

	return nil
}

// SendUpdateStartNotification sends a notification when update starts
func (m *Manager) SendUpdateStartNotification(modpackName, version string) error {
	if !m.IsEnabled() {
		return nil
	}

	var errors []error

	// Send to Discord
	if m.discord != nil {
		if err := m.discord.SendUpdateStartNotification(modpackName, version); err != nil {
			errors = append(errors, fmt.Errorf("Discord: %w", err))
		}
	}

	// Send to webhook
	if m.webhook != nil {
		if err := m.webhook.SendUpdateStartNotification(modpackName, version); err != nil {
			errors = append(errors, fmt.Errorf("Webhook: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %v", errors)
	}

	return nil
}

// SendUpdateSuccessNotification sends a notification when update succeeds
func (m *Manager) SendUpdateSuccessNotification(modpackName, version string, duration time.Duration) error {
	if !m.IsEnabled() {
		return nil
	}

	var errors []error

	// Send to Discord
	if m.discord != nil {
		if err := m.discord.SendUpdateSuccessNotification(modpackName, version, duration); err != nil {
			errors = append(errors, fmt.Errorf("Discord: %w", err))
		}
	}

	// Send to webhook
	if m.webhook != nil {
		if err := m.webhook.SendUpdateSuccessNotification(modpackName, version, duration); err != nil {
			errors = append(errors, fmt.Errorf("Webhook: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %v", errors)
	}

	return nil
}

// SendUpdateFailureNotification sends a notification when update fails
func (m *Manager) SendUpdateFailureNotification(modpackName, version string, errorMsg string) error {
	if !m.IsEnabled() {
		return nil
	}

	var errors []error

	// Send to Discord
	if m.discord != nil {
		if err := m.discord.SendUpdateFailureNotification(modpackName, version, errorMsg); err != nil {
			errors = append(errors, fmt.Errorf("Discord: %w", err))
		}
	}

	// Send to webhook
	if m.webhook != nil {
		if err := m.webhook.SendUpdateFailureNotification(modpackName, version, errorMsg); err != nil {
			errors = append(errors, fmt.Errorf("Webhook: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %v", errors)
	}

	return nil
}

// SendBackupNotification sends a backup notification
func (m *Manager) SendBackupNotification(action, backupName string, size int64) error {
	if !m.IsEnabled() {
		return nil
	}

	var errors []error

	// Send to Discord
	if m.discord != nil {
		if err := m.discord.SendBackupNotification(action, backupName, size); err != nil {
			errors = append(errors, fmt.Errorf("Discord: %w", err))
		}
	}

	// Send to webhook
	if m.webhook != nil {
		if err := m.webhook.SendBackupNotification(action, backupName, size); err != nil {
			errors = append(errors, fmt.Errorf("Webhook: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %v", errors)
	}

	return nil
}

// SendServerStatusNotification sends a server status notification
func (m *Manager) SendServerStatusNotification(status, message string) error {
	if !m.IsEnabled() {
		return nil
	}

	var errors []error

	// Send to Discord
	if m.discord != nil {
		if err := m.discord.SendServerStatusNotification(status, message); err != nil {
			errors = append(errors, fmt.Errorf("Discord: %w", err))
		}
	}

	// Send to webhook
	if m.webhook != nil {
		if err := m.webhook.SendServerStatusNotification(status, message); err != nil {
			errors = append(errors, fmt.Errorf("Webhook: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %v", errors)
	}

	return nil
}

// TestConnections tests all notification channels
func (m *Manager) TestConnections() error {
	if !m.IsEnabled() {
		return fmt.Errorf("notifications are not enabled")
	}

	var errors []error

	// Test Discord
	if m.discord != nil {
		if err := m.discord.TestConnection(); err != nil {
			errors = append(errors, fmt.Errorf("Discord test failed: %w", err))
		}
	}

	// Test webhook
	if m.webhook != nil {
		if err := m.webhook.TestConnection(); err != nil {
			errors = append(errors, fmt.Errorf("Webhook test failed: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification test errors: %v", errors)
	}

	return nil
}

// GetDiscordNotifier returns the Discord notifier (if enabled)
func (m *Manager) GetDiscordNotifier() *DiscordNotifier {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.discord
}

// GetWebhookNotifier returns the webhook notifier (if enabled)
func (m *Manager) GetWebhookNotifier() *WebhookNotifier {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.webhook
}

// UpdateConfig updates the notification configuration
func (m *Manager) UpdateConfig(config *config.NotificationConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Update Discord notifier
	if config.Discord.Enabled {
		m.discord = NewDiscordNotifier(&config.Discord)
	} else {
		m.discord = nil
	}

	// Update webhook notifier
	if config.Webhook.Enabled {
		m.webhook = NewWebhookNotifier(&config.Webhook)
	} else {
		m.webhook = nil
	}

	// Update enabled status
	m.enabled = config.Discord.Enabled || config.Webhook.Enabled
}

// SendCustomNotification sends a custom notification to specific channels
func (m *Manager) SendCustomNotification(message string, channels []string) error {
	if !m.IsEnabled() {
		return nil
	}

	var errors []error

	for _, channel := range channels {
		switch channel {
		case "discord":
			if m.discord != nil {
				if err := m.discord.SendMessage(message); err != nil {
					errors = append(errors, fmt.Errorf("Discord: %w", err))
				}
			}
		case "webhook":
			if m.webhook != nil {
				if err := m.webhook.SendNotification("custom", message, nil); err != nil {
					errors = append(errors, fmt.Errorf("Webhook: %w", err))
				}
			}
		default:
			errors = append(errors, fmt.Errorf("unknown channel: %s", channel))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %v", errors)
	}

	return nil
}

// GetEnabledChannels returns a list of enabled notification channels
func (m *Manager) GetEnabledChannels() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var channels []string

	if m.discord != nil {
		channels = append(channels, "discord")
	}

	if m.webhook != nil {
		channels = append(channels, "webhook")
	}

	return channels
}

// GetStatus returns the status of all notification channels
func (m *Manager) GetStatus() map[string]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[string]bool)

	status["discord"] = m.discord != nil
	status["webhook"] = m.webhook != nil
	status["enabled"] = m.enabled

	return status
}

// Disable disables all notifications
func (m *Manager) Disable() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.enabled = false
	m.discord = nil
	m.webhook = nil
}

// Enable enables notifications with the given configuration
func (m *Manager) Enable(config *config.NotificationConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if config.Discord.Enabled {
		m.discord = NewDiscordNotifier(&config.Discord)
	}

	if config.Webhook.Enabled {
		m.webhook = NewWebhookNotifier(&config.Webhook)
	}

	m.enabled = config.Discord.Enabled || config.Webhook.Enabled
}
