# CurseForge Auto-Updater

A Go CLI tool for automatically checking and downloading the latest versions of CurseForge mods.

## Features

- Check for mod updates without downloading
- Download latest mod files automatically
- Track download metadata to avoid unnecessary downloads
- Support for configuration via files, environment variables, or command-line flags
- Structured logging with configurable levels

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd golang

# Install Mage (build tool)
go install github.com/magefile/mage@latest

# Install dependencies and build
mage deps build
```

## Build Commands

This project uses [Mage](https://magefile.org/) for build automation:

```bash
# Build for current platform
mage build

# Build for all platforms (Windows, Linux, macOS)
mage buildAll

# Development with live reload (optional)
mage dev

# Complete build pipeline
mage all

# View all available commands
mage -l
```

For more details, see [MAGE.md](MAGE.md).

## Configuration

The tool can be configured through:

1. **Command line flags**
2. **Environment variables** (prefixed with `CURSEFORGE_`)
3. **Configuration file** (YAML format)

### Required Configuration

- `api-key` / `CURSEFORGE_API_KEY`: Your CurseForge API key
- `mod-id` / `CURSEFORGE_MOD_ID`: The mod ID to check for updates

### Optional Configuration

- `download-path` / `CURSEFORGE_DOWNLOAD_PATH`: Download directory (default: `./downloads`)
- `game-id` / `CURSEFORGE_GAME_ID`: Game ID (default: `432` for Minecraft)
- `log-level` / `CURSEFORGE_LOG_LEVEL`: Log level (default: `info`)

### Example Environment File (.env)

```env
CURSEFORGE_API_KEY=your-api-key-here
CURSEFORGE_MOD_ID=123456
CURSEFORGE_DOWNLOAD_PATH=./downloads
CURSEFORGE_GAME_ID=432
CURSEFORGE_LOG_LEVEL=info
```

## Usage

### Basic Update Check and Download

```bash
# Check and download if needed
./curseforge-updater --api-key="your-key" --mod-id="123456"

# Using environment variables
export CURSEFORGE_API_KEY="your-key"
export CURSEFORGE_MOD_ID="123456"
./curseforge-updater
```

### Check Only (No Download)

```bash
./curseforge-updater check --api-key="your-key" --mod-id="123456"
```

### Force Download

```bash
./curseforge-updater download --api-key="your-key" --mod-id="123456"
```

### Show Mod Information

```bash
./curseforge-updater info --api-key="your-key" --mod-id="123456"
```

## Commands

- `curseforge-updater` - Check for updates and download if needed (default command)
- `curseforge-updater check` - Check for updates without downloading
- `curseforge-updater download` - Force download the latest file
- `curseforge-updater info` - Show detailed mod information

## API Key

You need a CurseForge API key to use this tool. You can get one from:
https://docs.curseforge.com/#getting-started

## License

See LICENSE file.
