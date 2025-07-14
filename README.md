
---

# CurseForge AutoUpdate

An automated tool for checking, downloading, and managing updates for CurseForge mods or modpacks. Supports both a Python proof-of-concept and a modern Golang CLI for advanced server automation.

---

## Project Overview

This repository contains:

- **Python PoC**: A proof-of-concept for CurseForge API integration and mod downloading.
- **Golang CLI**: A new, extensible CLI tool for server admins, with planned features for backup, update, restore, and notifications.

---

## Python Proof-of-Concept

Located in [`python/`](python/). See [`python/README.md`](python/README.md) for full details.

**Features:**
- Connects to CurseForge Core API v1
- Fetches mod info, lists files, downloads latest file
- Uses `.env` for config (API key, mod ID, download path)
- Handles errors and API fallbacks

**Quickstart:**
```bash
cd python
pip install -r requirements.txt
cp .env.example .env  # or create .env manually
python poc.py
```

See [`python/README.md`](python/README.md) for CLI usage, library usage, and config options.

---

## Golang CLI (In Development)

Located in [`golang/`](golang/). See [`golang/PLAN.md`](golang/PLAN.md) for full roadmap.

**Current Status:**
- CLI skeleton with commands: `init`, `check`, `update`, `backup`, `restore`, `notify`, `list`, `version`
- Config management: supports TOML, YAML, JSON, and .env templates
- Embedded config templates and interactive config creation
- Modular API client for CurseForge (mod info, file listing, download, search)

**Planned Features:**
- Full modpack update automation (backup, update, restore)
- Discord/webhook notifications
- Minecraft server integration (start/stop, notify players)
- Multi-modpack and multi-server support
- Scheduling, retention, rollback, and more

**Example CLI Usage:**
```bash
cd golang
go run ./cmd/cli/main.go --init toml   # Scaffold config
go run ./cmd/cli/main.go check          # Check mod exists
go run ./cmd/cli/main.go update         # (Planned) Update modpack
```

---

## Directory Structure

```
curseforge-autoupdate/
├── golang/      # Golang CLI implementation
│   ├── cmd/cli/         # CLI entry and commands
│   ├── internal/api/    # CurseForge API client
│   ├── internal/server/ # Server/backup logic
│   ├── internal/config/ # Config types/templates
│   ├── helper/          # Env, filesystem, version helpers
│   └── templates/       # Config templates
├── python/      # Python PoC and library
│   ├── updater/         # Main package code
│   ├── cli.py           # CLI
│   ├── poc.py           # Legacy PoC
│   └── downloads/       # Downloaded files
├── LICENSE
├── README.md    # This file
```

---

## Development Plan (Golang)

See [`golang/DEVELOPMENT_PLAN.md`](golang/DEVELOPMENT_PLAN.md) for a detailed roadmap, including:
- Enhanced config system (multi-format, templates)
- Modular API client
- Command structure: `init`, `check`, `update`, `backup`, `restore`, `notify`, `list`, `version`
- Server management, backup/restore, notification, scheduling

---

## Configuration Examples

**Python (.env):**
```env
CURSEFORGE_API_KEY=your_api_key_here
MOD_ID=123456
DOWNLOAD_PATH=./downloads
```

**Golang (config.toml):**
```toml
api_key = "your_api_key_here"
mod_id = 123456
# ...see templates for more fields
```

---

## License
[MIT](LICENSE)

---
   python poc.py
   ```

## How the Python PoC Works

The Python PoC follows these steps:

1. **Configuration Loading**: Reads API key, mod ID, and download path from `.env` file
2. **Mod Validation**: Queries the CurseForge API to verify the mod exists and get basic information
3. **File Discovery**: Fetches all available files for the specified mod
4. **Latest File Selection**: Identifies the most recent file based on upload date
5. **File Download**: Downloads the latest file to the specified directory with progress feedback

### Current Implementation Details

**API Integration:**
- Uses CurseForge Core API v1 endpoints
- Requires valid API key for authentication
- Implements multiple fallback strategies if initial requests fail
- Provides detailed debugging output for troubleshooting

**File Handling:**
- Downloads files using streaming for memory efficiency
- Preserves original filenames from CurseForge
- Creates download directories automatically

**Error Handling:**
- Validates API responses and provides helpful error messages
- Includes troubleshooting suggestions for common issues
- Gracefully handles network errors and API limitations

## File Structure

After running the Python PoC, you'll have:
```
curseforge-autoupdate/
├── python/                 # Python PoC directory
│   ├── .env               # Your configuration
│   ├── poc.py             # Main PoC script
│   ├── requirements.txt   # Python dependencies
│   ├── README.md          # Python-specific documentation
│   └── downloads/         # Downloaded files (auto-generated)
│       └── [mod-files]    # Downloaded mod/modpack files
├── .gitignore
├── LICENSE
└── README.md              # This file
```

## Configuration Examples

### Basic Mod Download
```env
CURSEFORGE_API_KEY=your_api_key_here
MOD_ID=123456
DOWNLOAD_PATH=./downloads
```

### Custom Download Location
```env
CURSEFORGE_API_KEY=your_api_key_here
MOD_ID=123456
DOWNLOAD_PATH=C:/Server/mods
```

## Finding Mod IDs

To find a CurseForge mod ID:
1. Go to the mod's CurseForge page
2. Look at the URL: `https://www.curseforge.com/minecraft/modpacks/[mod-name]/files`
3. Click on any file and check the URL: `https://www.curseforge.com/minecraft/modpacks/[mod-name]/files/[file-id]`
4. Or use the CurseForge API to search by name

## API Key Setup

An API key is **required** for the Python PoC:
1. Visit [CurseForge Core API Console](https://console.curseforge.com/)
2. Sign up or log in with your CurseForge account
3. Create a new API key
4. Add it to your `.env` file as `CURSEFORGE_API_KEY=your_key_here`

## Python PoC Limitations

The current Python implementation is a proof-of-concept with the following limitations:

- **No Version Tracking**: Downloads the latest file each time, doesn't check if it's already downloaded
- **No Filtering**: Cannot filter by Minecraft version, mod loader, or file type
- **No Manifest Extraction**: Downloads files as-is without extracting metadata
- **Single Mod Support**: Only handles one mod at a time
- **Basic Error Handling**: Limited recovery options for failed downloads

## Troubleshooting

### Common Issues

**"No API key found"**
- Ensure you have a `.env` file in the `python/` directory
- Verify the API key is correctly formatted: `CURSEFORGE_API_KEY=your_key_here`

**"No files found"**
- Check if the `MOD_ID` is correct (find it in the mod's CurseForge URL)
- Verify the mod has public files available
- Some mods may have restricted downloads

**"Request failed" or API errors**
- Verify your API key is valid and active
- Check your internet connection
- CurseForge API may have rate limits or temporary issues

**Download failures**
- Ensure the download directory exists and is writable
- Check available disk space
- Some files may have download restrictions

## Use Cases

**Current (Python PoC):**
- Testing CurseForge API integration
- Manual mod/modpack downloads
- API key validation and troubleshooting
- Learning the CurseForge API structure

**Planned (Golang Implementation):**
- **Server Administrators**: Automate modpack updates for Minecraft servers
- **Development**: Keep development environments up-to-date during modpack creation
- **CI/CD Integration**: Integrate into deployment pipelines for automated server updates

## Planned Features (Golang Implementation)

The planned Golang rewrite will include:

- **Smart Version Detection**: Compare current vs. latest file versions using CurseForge file IDs
- **Manifest-Based Tracking**: Use modpack `manifest.json` files for reliable version tracking
- **Automatic Extraction**: Extract and analyze ZIP contents including manifest files
- **Filtering Support**: Filter by Minecraft version, mod loader (Forge, Fabric, Quilt)
- **Selective Updates**: Compare individual mod versions to detect specific mod updates
- **Rollback Capability**: Keep previous versions for easy rollback
- **Multi-modpack Support**: Monitor and update multiple modpacks simultaneously
- **Integration Hooks**: Pre/post-download scripts for custom deployment logic
- **Configuration Management**: Enhanced configuration with validation and templates

## Development Notes

* **Current Status**: This repository contains a Python proof-of-concept demonstrating basic CurseForge API integration
* **Target Implementation**: The goal is to rewrite this functionality in **Golang** with comprehensive features
* **Focus**: Currently focused on Minecraft, but could be expanded to other games supported by CurseForge
* **Purpose**: Designed to reduce manual server update steps during modpack development
* **Scope**: The Python PoC only downloads files; installation/deployment logic will be added in the Golang version

## Dependencies (Python PoC)

- `requests>=2.28.0` - HTTP library for CurseForge API calls
- `python-dotenv>=1.0.0` - Environment variable management from `.env` files
- Built-in libraries: `pathlib`, `json`, `zipfile`, `os`


<!-- vim: set ft=markdown : -->
## License
[MIT](LICENSE)

---