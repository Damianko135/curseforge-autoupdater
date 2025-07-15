package server

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/damianko135/curseforge-autoupdate/golang/helper/filesystem"
)

// MinecraftServer represents a Minecraft server instance
type MinecraftServer struct {
	serverPath string
	jarName    string
	process    *exec.Cmd
	isRunning  bool
	mu         sync.RWMutex
	stopChan   chan struct{}
	logChan    chan string
	errorChan  chan error
	startTime  time.Time
}

// NewMinecraftServer creates a new Minecraft server instance
func NewMinecraftServer(serverPath, jarName string) *MinecraftServer {
	return &MinecraftServer{
		serverPath: serverPath,
		jarName:    jarName,
		stopChan:   make(chan struct{}),
		logChan:    make(chan string, 100),
		errorChan:  make(chan error, 10),
	}
}

// Start starts the Minecraft server
func (s *MinecraftServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("server is already running")
	}

	// Check if server JAR exists
	jarPath := filepath.Join(s.serverPath, s.jarName)
	if !filesystem.FileExists(jarPath) {
		return fmt.Errorf("server JAR not found: %s", jarPath)
	}

	// Check if server directory exists
	if !filesystem.DirExists(s.serverPath) {
		return fmt.Errorf("server directory not found: %s", s.serverPath)
	}

	// Create start command
	args := []string{
		"-Xmx2G",
		"-Xms1G",
		"-jar",
		s.jarName,
		"nogui",
	}

	// #nosec G204 -- args are validated elsewhere
	s.process = exec.Command("java", args...)
	s.process.Dir = s.serverPath

	// Set up pipes for stdout and stderr
	stdout, err := s.process.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := s.process.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if _, err := s.process.StdinPipe(); err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	// Do not assign stdin to s.process.Stdin (types are incompatible)

	// Start the process
	if err := s.process.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.isRunning = true
	s.startTime = time.Now()

	// Start log monitoring goroutines
	if file, ok := stdout.(*os.File); ok {
		go s.monitorOutput(file, "stdout")
	}
	if file, ok := stderr.(*os.File); ok {
		go s.monitorOutput(file, "stderr")
	}

	// Start process monitoring
	go s.monitorProcess()

	// No need to store stdin again; already available via s.process.Stdin
	return nil
}

// Stop stops the Minecraft server gracefully
func (s *MinecraftServer) Stop(timeout time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("server is not running")
	}

	// Send stop command
	if err := s.sendCommand("stop"); err != nil {
		return fmt.Errorf("failed to send stop command: %w", err)
	}

	// Wait for graceful shutdown with timeout
	done := make(chan error, 1)
	go func() {
		done <- s.process.Wait()
	}()

	select {
	case err := <-done:
		s.isRunning = false
		if err != nil {
			return fmt.Errorf("server stopped with error: %w", err)
		}
		return nil
	case <-time.After(timeout):
		// Force kill if timeout reached
		if err := s.process.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill server process: %w", err)
		}
		s.isRunning = false
		return fmt.Errorf("server did not stop gracefully within timeout, killed")
	}
}

// IsRunning returns whether the server is currently running
func (s *MinecraftServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}

// SendCommand sends a command to the server
func (s *MinecraftServer) SendCommand(command string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isRunning {
		return fmt.Errorf("server is not running")
	}

	return s.sendCommand(command)
}

// sendCommand sends a command to the server (internal method)
func (s *MinecraftServer) sendCommand(command string) error {
	if s.process == nil || s.process.Stdin == nil {
		return fmt.Errorf("server process or stdin is not available")
	}

	// Use the original stdin pipe for writing commands
	if s.process == nil {
		return fmt.Errorf("server process is not available")
	}
	// Try to get the stdin pipe from the process
	if stdin, err := s.process.StdinPipe(); err == nil {
		defer stdin.Close()
		_, err := fmt.Fprintf(stdin, "%s\n", command)
		if err != nil {
			return fmt.Errorf("failed to write command to stdin: %w", err)
		}
		return nil
	} else {
		return fmt.Errorf("stdin pipe is not available: %w", err)
	}
}

// GetUptime returns the server uptime
func (s *MinecraftServer) GetUptime() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isRunning {
		return 0
	}

	return time.Since(s.startTime)
}

// GetLogChannel returns the log channel
func (s *MinecraftServer) GetLogChannel() <-chan string {
	return s.logChan
}

// GetErrorChannel returns the error channel
func (s *MinecraftServer) GetErrorChannel() <-chan error {
	return s.errorChan
}

// monitorOutput monitors server output
func (s *MinecraftServer) monitorOutput(pipe *os.File, source string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		logEntry := fmt.Sprintf("[%s] %s", source, line)

		select {
		case s.logChan <- logEntry:
		default:
			// Log channel is full, discard oldest entry and add new
			select {
			case <-s.logChan:
			default:
			}
			// Try again to send logEntry (should succeed now)
			select {
			case s.logChan <- logEntry:
			default:
			}
		}
	}
}

// monitorProcess monitors the server process
func (s *MinecraftServer) monitorProcess() {
	err := s.process.Wait()

	s.mu.Lock()
	s.isRunning = false
	s.mu.Unlock()

	if err != nil {
		select {
		case s.errorChan <- fmt.Errorf("server process exited with error: %w", err):
		default:
		}
	}

	close(s.stopChan)
}

// WaitForShutdown waits for the server to shut down
func (s *MinecraftServer) WaitForShutdown() {
	<-s.stopChan
}

// GetServerInfo returns basic server information
func (s *MinecraftServer) GetServerInfo() ServerInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info := ServerInfo{
		ServerPath: s.serverPath,
		JarName:    s.jarName,
		IsRunning:  s.isRunning,
		Uptime:     0,
	}

	if s.isRunning {
		info.Uptime = time.Since(s.startTime)
	}

	return info
}

// ServerInfo represents basic server information
type ServerInfo struct {
	ServerPath string
	JarName    string
	IsRunning  bool
	Uptime     time.Duration
}

// BroadcastMessage sends a message to all players
func (s *MinecraftServer) BroadcastMessage(message string) error {
	return s.SendCommand(fmt.Sprintf("say %s", message))
}

// KickAllPlayers kicks all players from the server
func (s *MinecraftServer) KickAllPlayers(reason string) error {
	if reason == "" {
		reason = "Server maintenance"
	}
	return s.SendCommand(fmt.Sprintf("kick @a %s", reason))
}

// SetWorldTime sets the world time
func (s *MinecraftServer) SetWorldTime(time string) error {
	return s.SendCommand(fmt.Sprintf("time set %s", time))
}

// SetWeather sets the weather
func (s *MinecraftServer) SetWeather(weather string) error {
	return s.SendCommand(fmt.Sprintf("weather %s", weather))
}

// SaveWorld saves the world
func (s *MinecraftServer) SaveWorld() error {
	return s.SendCommand("save-all")
}

// ReloadServer reloads the server
func (s *MinecraftServer) ReloadServer() error {
	return s.SendCommand("reload")
}

// GetOnlinePlayers gets the list of online players
func (s *MinecraftServer) GetOnlinePlayers() error {
	return s.SendCommand("list")
}

// NotifyPlayersBeforeShutdown notifies players before server shutdown
func (s *MinecraftServer) NotifyPlayersBeforeShutdown(countdown int) error {
	messages := []string{
		"Server will be shutting down for maintenance in %d minutes",
		"Server will be shutting down in %d minutes",
		"Server shutdown in %d minutes",
	}

	for i := countdown; i > 0; i-- {
		var message string
		if i <= len(messages) {
			message = fmt.Sprintf(messages[i-1], i)
		} else {
			message = fmt.Sprintf("Server will be shutting down in %d minutes", i)
		}

		if err := s.BroadcastMessage(message); err != nil {
			return fmt.Errorf("failed to broadcast countdown message: %w", err)
		}

		if i > 1 {
			time.Sleep(1 * time.Minute)
		}
	}

	return s.BroadcastMessage("Server is shutting down now for maintenance")
}

// CheckServerHealth checks if the server is healthy
func (s *MinecraftServer) CheckServerHealth() error {
	if !s.IsRunning() {
		return fmt.Errorf("server is not running")
	}

	// Check if server directory exists
	if _, err := os.Stat(s.serverPath); err != nil {
		return fmt.Errorf("server directory not accessible: %w", err)
	}

	// Check if JAR file exists
	jarPath := filepath.Join(s.serverPath, s.jarName)
	if _, err := os.Stat(jarPath); err != nil {
		return fmt.Errorf("server JAR not accessible: %w", err)
	}

	return nil
}

// GetServerProperties reads server properties
func (s *MinecraftServer) GetServerProperties() (map[string]string, error) {
	propertiesPath := filepath.Join(s.serverPath, "server.properties")

	// #nosec G304 -- propertiesPath is constructed internally
	file, err := os.Open(propertiesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open server.properties: %w", err)
	}
	defer file.Close()

	properties := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			properties[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading server.properties: %w", err)
	}

	return properties, nil
}

// UpdateServerProperties updates server properties
func (s *MinecraftServer) UpdateServerProperties(properties map[string]string) error {
	propertiesPath := filepath.Join(s.serverPath, "server.properties")

	// Read existing properties
	existing, err := s.GetServerProperties()
	if err != nil {
		existing = make(map[string]string)
	}

	// Merge with new properties
	for key, value := range properties {
		existing[key] = value
	}

	// Write back to file
	// #nosec G304 -- propertiesPath is constructed internally
	file, err := os.Create(propertiesPath)
	if err != nil {
		return fmt.Errorf("failed to create server.properties: %w", err)
	}
	defer file.Close()

	for key, value := range existing {
		_, err := fmt.Fprintf(file, "%s=%s\n", key, value)
		if err != nil {
			return fmt.Errorf("failed to write property %s: %w", key, err)
		}
	}

	return nil
}

// Restart restarts the server
func (s *MinecraftServer) Restart(timeout time.Duration) error {
	if s.IsRunning() {
		if err := s.Stop(timeout); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
		// Wait a moment before restarting
		time.Sleep(2 * time.Second)
	}
	return s.Start()
}

// GetServerVersion attempts to get the server version
func (s *MinecraftServer) GetServerVersion() (string, error) {
	// Try to read from server.properties
	properties, err := s.GetServerProperties()
	if err != nil {
		return "", fmt.Errorf("failed to read server properties: %w", err)
	}

	if version, exists := properties["version"]; exists {
		return version, nil
	}

	// If not found in properties, try to determine from JAR name
	if strings.Contains(s.jarName, "server") {
		// Extract version from JAR name if possible
		parts := strings.Split(s.jarName, "-")
		for _, part := range parts {
			if len(part) > 0 && (part[0] >= '0' && part[0] <= '9') {
				return part, nil
			}
		}
	}

	return "unknown", nil
}
