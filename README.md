
---

# CurseForge AutoUpdate

An automated tool for checking and downloading updates for CurseForge mods or modpacks using the CurseForge API. This script intelligently compares your current version with the latest available version and downloads updates when needed.

## Features

- **Smart Version Detection**: Compares current vs. latest file versions using CurseForge file IDs
- **Automatic Downloads**: Downloads new versions with progress indication
- **Version Tracking**: Maintains local version state between runs
- **Flexible Filtering**: Filter by Minecraft version and mod loader (Forge, Fabric, Quilt)
- **Environment Configuration**: Uses `.env` files for easy configuration management
- **API Key Support**: Optional API key support for higher rate limits
- **Error Handling**: Comprehensive error handling and user feedback

## Setup

1. **Clone and navigate to the project:**
   ```bash
   git clone <repository-url>
   cd curseforge-autoupdate
   ```

2. **Copy `.env.example` to `.env` and configure:**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` with your configuration:
   - `CURSEFORGE_API_KEY` - Optional API key for higher rate limits
   - `MOD_ID` - The CurseForge mod/modpack ID (required)
   - `GAME_ID` - Game ID (432 for Minecraft, default)
   - `DOWNLOAD_PATH` - Download directory (default: ./downloads)
   - `EXTRACT_PATH` - Directory to extract manifest.json (default: ./extracted)
   - `MINECRAFT_VERSION` - Filter by MC version (optional)
   - `MOD_LOADER` - Filter by loader: forge, fabric, quilt (optional)

3. **Install Python dependencies:**
   ```bash
   pip install -r requirements.txt
   ```

4. **Run the script:**
   ```bash
   python poc.py
   ```

## How It Works

1. **Current Version Check**: Looks for existing `manifest.json` in the extracted directory or extracts it from downloaded ZIP files
2. **API Query**: Fetches latest files from CurseForge API for the specified mod
3. **Manifest Comparison**: Compares current manifest data (version, file date) with latest available file
4. **Download**: If versions differ, downloads the new version with progress tracking
5. **Manifest Extraction**: Automatically extracts and saves `manifest.json` from the downloaded ZIP file for future comparisons

### Manifest-Based Version Detection

The script now uses the modpack's own `manifest.json` file for version tracking, which provides:

- **Modpack version** - The actual version number from the modpack creator
- **File metadata** - Creation date, file ID, and other CurseForge data
- **Mod list** - Complete list of mods with their project and file IDs
- **Minecraft version** - Target Minecraft version for the modpack
- **Mod loader info** - Forge, Fabric, or Quilt version information

This approach is more reliable than external tracking files and uses the modpack's native metadata.

## File Structure

After running the script, you'll have:
```
curseforge-autoupdate/
├── .env                    # Your configuration (create from .env.example)
├── .env.example           # Configuration template
├── poc.py                 # Main script
├── requirements.txt       # Python dependencies
├── extracted/             # Extracted manifest files (auto-generated)
│   └── manifest.json     # Current modpack manifest
└── downloads/            # Downloaded files (auto-generated)
    └── [modpack-files.zip]
```

The script now uses the modpack's native `manifest.json` file (extracted from ZIP files) for version tracking, eliminating the need for separate tracking files.

## Configuration Examples

### Basic Modpack Setup
```env
MOD_ID=123456
DOWNLOAD_PATH=./modpacks
```

### Minecraft 1.20.1 with Forge
```env
MOD_ID=123456
MINECRAFT_VERSION=1.20.1
MOD_LOADER=forge
DOWNLOAD_PATH=./server-files
```

### With API Key (Recommended)
```env
CURSEFORGE_API_KEY=your-api-key-here
MOD_ID=123456
```

## Finding Mod IDs

To find a CurseForge mod ID:
1. Go to the mod's CurseForge page
2. Look at the URL: `https://www.curseforge.com/minecraft/modpacks/[mod-name]/files`
3. Click on any file and check the URL: `https://www.curseforge.com/minecraft/modpacks/[mod-name]/files/[file-id]`
4. Or use the CurseForge API to search by name

## API Key Setup

While not required, an API key provides higher rate limits:
1. Visit [CurseForge Core API](https://docs.curseforge.com/)
2. Register for an API key
3. Add it to your `.env` file as `CURSEFORGE_API_KEY`

## Troubleshooting

- **No files found**: Check if `MOD_ID` is correct and the mod has files for your specified filters
- **Download fails**: Verify internet connection and that the file isn't restricted
- **Import errors**: Ensure all dependencies are installed with `pip install -r requirements.txt`
- **Permission errors**: Check write permissions for download directory

## Use Cases

- **Server Administrators**: Automate modpack updates for Minecraft servers
- **Development**: Keep development environments up-to-date during modpack creation
- **CI/CD Integration**: Integrate into deployment pipelines for automated server updates

## Future Enhancements

- **Selective Mod Updates**: Compare individual mod versions from manifest files to detect specific mod updates
- **Rollback Capability**: Keep previous manifests for easy rollback to earlier versions
- **Integration Hooks**: Add pre/post-download scripts for custom deployment logic
- **Multi-modpack Support**: Monitor and update multiple modpacks simultaneously
- **Automatic Extraction**: Optionally extract the entire modpack, not just the manifest

## Notes

* This PoC is written in Python. The goal is to rework it in **Golang** later.
* Currently focused on Minecraft, but could be expanded to other games supported by CurseForge.
* Designed to reduce manual server update steps during modpack development.
* The script only downloads; installation/deployment logic can be added as needed.

## Dependencies

- `requests` - HTTP library for API calls
- `python-dotenv` - Environment variable management
- `pathlib` - Modern path handling (built-in)
- `json` - JSON parsing (built-in)


<!-- vim: set ft=markdown : -->
## License
[MIT](LICENSE)

---