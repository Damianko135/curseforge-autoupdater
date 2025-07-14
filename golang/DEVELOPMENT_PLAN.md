# **CurseForge Auto-Update CLI Tool - Detailed Development Plan**

## **Project Overview**
A CLI tool to automatically manage Minecraft modpack updates on servers with full lifecycle management (backup, update, restore) and notification capabilities.

## **Core Architecture**

### **Current State Analysis**
✅ **Existing Components:**
- Basic CLI structure with Cobra
- Config management (TOML, .env, YAML, JSON)
- CurseForge API client with mod existence checking
- Basic command structure

### **Proposed Directory Structure**
```
golang/
├── cmd/cli/main.go              # Main CLI entry point
├── internal/
│   ├── api/                     # CurseForge API client
│   │   ├── client.go           # Enhanced API client
│   │   ├── modpack.go          # Modpack-specific operations
│   │   └── types.go            # API response types
│   ├── server/                  # Server management
│   │   ├── minecraft.go        # Minecraft server operations
│   │   └── backup.go           # Backup/restore functionality
│   ├── notification/            # Notification system
│   │   ├── discord.go          # Discord notifications
│   │   └── webhook.go          # Generic webhook support
│   └── config/                  # Configuration management
│       ├── types.go            # Config structures
│       └── templates.go        # Config templates
├── helper/
│   ├── env/                    # Environment helpers (existing)
│   ├── filesystem/             # File operations
│   └── version/                # Version comparison
└── templates/                   # Config templates
    ├── init.toml               # Initial setup template
    └── server.toml             # Server configuration template
```

## **Phase 1: Core Infrastructure Enhancement**

### **1.1 Enhanced Configuration System**
**New Config Structure:**
```go
type Config struct {
    // API Configuration
    APIKey string `mapstructure:"api_key"`
    
    // Modpack Configuration
    ModpackID   int    `mapstructure:"modpack_id"`
    GameVersion string `mapstructure:"game_version"`
    
    // Server Configuration
    ServerPath     string `mapstructure:"server_path"`
    BackupPath     string `mapstructure:"backup_path"`
    ServerJarName  string `mapstructure:"server_jar_name"`
    
    // Notification Configuration
    Notifications NotificationConfig `mapstructure:"notifications"`
    
    // Update Configuration
    AutoUpdate    bool   `mapstructure:"auto_update"`
    UpdateChannel string `mapstructure:"update_channel"` // stable, beta, alpha
}

type NotificationConfig struct {
    Discord DiscordConfig `mapstructure:"discord"`
    Webhook WebhookConfig `mapstructure:"webhook"`
}
```

### **1.2 Enhanced API Client**
**New Features:**
- Modpack version checking
- File download capabilities
- Version comparison
- Changelog retrieval

## **Phase 2: Command Implementation**

### **2.1 Command Structure**
```
curseforge-autoupdate
├── init           # Initialize new project
├── check          # Check for updates (enhanced)
├── update         # Perform update process
├── backup         # Manual backup operations
├── restore        # Restore from backup
├── notify         # Send notifications
├── list           # List available commands/info
├── help           # Detailed help system
└── version        # Show version info
```

### **2.2 Command Details**

#### **`init` Command**
- **Purpose**: Initialize a new project with configuration templates
- **Features**:
  - Interactive setup wizard
  - Generate config templates
  - Validate server paths
  - Test API connectivity
  - Create initial directory structure

#### **`check` Command (Enhanced)**
- **Purpose**: Check for modpack updates
- **Features**:
  - Compare current vs latest version
  - Show changelog if available
  - Validate compatibility
  - Dry-run mode

#### **`update` Command**
- **Purpose**: Perform the full update process
- **Features**:
  - Pre-update validation
  - Automatic backup creation
  - Download new modpack
  - Server shutdown/startup
  - Rollback on failure
  - Post-update notifications

#### **`backup` Command**
- **Purpose**: Manual backup operations
- **Features**:
  - Create named backups
  - List existing backups
  - Compress/decompress
  - Cleanup old backups

#### **`restore` Command**
- **Purpose**: Restore from backup
- **Features**:
  - List available backups
  - Selective restore
  - Validation before restore

#### **`notify` Command**
- **Purpose**: Send notifications manually
- **Features**:
  - Test notification systems
  - Send custom messages
  - Verify webhook endpoints

## **Phase 3: Server Management**

### **3.1 Minecraft Server Integration**
**Features:**
- Server process management (start/stop)
- Player notification before shutdown
- Server status monitoring
- Configuration file management

### **3.2 Backup System**
**Features:**
- Incremental backups
- Compression support
- Retention policies
- Backup verification

## **Phase 4: Notification System**

### **4.1 Discord Integration**
**Features:**
- Rich embed messages
- Update notifications
- Error alerts
- Status updates

### **4.2 Generic Webhook Support**
**Features:**
- Configurable webhook endpoints
- Custom message templates
- Retry mechanisms
- Error handling

## **Phase 5: Advanced Features**

### **5.1 Scheduling & Automation**
**Features:**
- Cron-like scheduling
- Automatic update checks
- Maintenance windows
- Player activity monitoring

### **5.2 Multi-Server Support**
**Features:**
- Multiple server configurations
- Batch operations
- Server groups
- Centralized management

## **Implementation Timeline**

### **Week 1-2: Core Infrastructure**
- [ ] Enhance configuration system
- [ ] Restructure API client
- [ ] Implement basic backup functionality
- [ ] Create config templates

### **Week 3-4: Command Implementation**
- [ ] Implement `init` command
- [ ] Enhance `check` command
- [ ] Implement `update` command
- [ ] Add `backup` and `restore` commands

### **Week 5-6: Server Management**
- [ ] Minecraft server integration
- [ ] Process management
- [ ] File system operations
- [ ] Error handling

### **Week 7-8: Notification System**
- [ ] Discord integration
- [ ] Webhook support
- [ ] Message templates
- [ ] Testing framework

### **Week 9-10: Polish & Testing**
- [ ] Comprehensive testing
- [ ] Documentation
- [ ] Error handling improvements
- [ ] Performance optimization

## **Dependencies to Add**
```go
// Additional dependencies needed
"github.com/robfig/cron/v3"           // Scheduling
"github.com/bwmarrin/discordgo"       // Discord integration
"github.com/klauspost/compress/zip"   // Compression
"github.com/shirou/gopsutil/v3"      // System monitoring
```

## **Key Design Decisions**

1. **Modular Architecture**: Each major feature in separate packages
2. **Configuration-Driven**: All behavior configurable via TOML files
3. **Fail-Safe Operations**: Always backup before making changes
4. **Extensible Notifications**: Plugin-like notification system
5. **Comprehensive Logging**: Detailed logging for troubleshooting
6. **Rollback Capability**: Ability to undo changes if something goes wrong

## **Technical Requirements**

### **API Integration**
- CurseForge API v1 compatibility
- Rate limiting and error handling
- Modpack-specific endpoints focus
- File download with progress tracking

### **File System Operations**
- Cross-platform file operations
- Atomic file operations where possible
- Permission handling
- Disk space monitoring

### **Process Management**
- Graceful server shutdown
- Process monitoring and health checks
- Signal handling
- Timeout management

### **Error Handling**
- Comprehensive error types
- Graceful degradation
- Retry mechanisms
- User-friendly error messages

## **Configuration Templates**

### **Main Config Template (init.toml)**
```toml
# CurseForge API Configuration
api_key = "your-api-key-here"

# Modpack Configuration
modpack_id = 0
game_version = "1.20.1"

# Server Configuration
server_path = "/path/to/server"
backup_path = "/path/to/backups"
server_jar_name = "server.jar"

# Update Configuration
auto_update = false
update_channel = "stable"

# Notification Configuration
[notifications.discord]
enabled = false
webhook_url = ""
channel_id = ""

[notifications.webhook]
enabled = false
url = ""
```

### **Server Config Template (server.toml)**
```toml
# Server-specific configuration
[server]
name = "Minecraft Server"
port = 25565
max_players = 20
shutdown_timeout = 30

[backup]
retention_days = 30
compression = true
incremental = true

[maintenance]
window_start = "02:00"
window_end = "04:00"
timezone = "UTC"
```

## **Security Considerations**

1. **API Key Management**: Secure storage and handling of API keys
2. **File System Security**: Proper permission checks and path validation
3. **Process Security**: Safe process execution and signal handling
4. **Network Security**: HTTPS enforcement and certificate validation
5. **Input Validation**: Sanitization of all user inputs and configuration values

## **Testing Strategy**

### **Unit Tests**
- API client functionality
- Configuration parsing
- File system operations
- Notification systems

### **Integration Tests**
- End-to-end command execution
- Server management workflows
- Backup and restore processes
- Notification delivery

### **Manual Tests**
- Real server update scenarios
- Error recovery testing
- Performance under load
- User experience validation

## **Documentation Plan**

1. **README.md**: Quick start guide and overview
2. **INSTALLATION.md**: Installation instructions
3. **CONFIGURATION.md**: Detailed configuration guide
4. **COMMANDS.md**: Command reference
5. **TROUBLESHOOTING.md**: Common issues and solutions
6. **API.md**: API integration details
7. **CONTRIBUTING.md**: Development guidelines

## **Future Enhancements**

### **Phase 6: Advanced Features**
- Web dashboard for monitoring
- Multiple game support
- Plugin/mod management
- Performance metrics
- Community features

### **Phase 7: Enterprise Features**
- Multi-server management
- Role-based access control
- Audit logging
- API for external integrations
- Clustering support

---

*This document serves as the master plan for the CurseForge Auto-Update CLI tool development. It should be updated as requirements change and implementation progresses.*