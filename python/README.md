---
# CurseForge Auto-Updater

A modern Python tool and library for automatically downloading and updating mods from CurseForge.

## Features

- Connects to the CurseForge API
- Fetches the latest file for a mod
- Downloads and records mod files
- CLI and library usage
- Metadata and update checks
- Configurable via `.env` or environment variables

## Quick Setup

1. Create and activate a virtual environment (optional but recommended):

   ```bash
   python3 -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   ```

2. Install dependencies:

   ```bash
   pip install -r requirements.txt
   ```

3. Copy the example environment file and add your CurseForge API key:

   ```bash
   cp .env.example .env
   # Then edit .env and set CURSEFORGE_API_KEY=your_key
   ```

4. (Optional) Install as a package:

   ```bash
   pip install .
   ```

Get your API key from: <https://console.curseforge.com/>

## Usage

### CLI

Run the updater using the CLI:

```bash
python cli.py --mod-id 123456
# or, if installed:
curseforge-update --mod-id 123456
```

See all options:

```bash
python cli.py --help
```

### As a Library

You can use the updater in your own Python scripts:

```python
from updater import main, get_config
config = get_config()
```

## Vision

The Python implementation serves as a proof-of-concept for CurseForge API integration and mod downloading. It demonstrates core update logic, error handling, and configuration management, and acts as a foundation for the more advanced Golang CLI. The goal is to provide a simple, reliable, and easily extensible tool for mod and modpack management.
   ```bash
   pip install -r requirements.txt
   ```
3. Copy the example environment file and add your CurseForge API key:
   ```bash
   cp .env.example .env
   # Then edit .env and set CURSEFORGE_API_KEY=your_key
   ```
4. (Optional) Install as a package:
   ```bash
   pip install .
   ```

Get your API key from: https://console.curseforge.com/

## Usage

### CLI

Run the updater using the CLI:

```bash
python cli.py --mod-id 123456
# or, if installed:
curseforge-update --mod-id 123456
```

See all options:
```bash
python cli.py --help
```

### As a Library

You can use the updater in your own Python scripts:

```python
from updater import main, get_config
config = get_config()
main()
```

### 'Legacy' PoC

The original proof-of-concept is still available:

```bash
python poc.py
```

## Configuration

Set in `.env` file:
- `CURSEFORGE_API_KEY` - Your API key (required)
- `MOD_ID` - Mod ID to check (optional)
- `DOWNLOAD_PATH` - Download location (optional)
- `GAME_ID` - Game ID (default: 432 for Minecraft)
- `MOD_LOADER` - Mod loader type (optional)
- `MINECRAFT_VERSION` - Minecraft version filter (optional)
- `EXTRACT_PATH` - Extraction path (optional)

## Project Structure

- `updater/` - Main package code
- `cli.py`   - Command-line interface
- `poc.py`   - 'Legacy' proof-of-concept script
- `test.py`  - Basic test script
- `requirements.txt` - Dependencies
- `setup.py` - Packaging script

## License

[MIT License](LICENSE) (C) 2025 by [Damianko135](https://github.com/Damianko135)