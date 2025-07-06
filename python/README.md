# CurseForge Auto-Updater PoC

Simple proof-of-concept for downloading the latest CurseForge mod file.

## Quick Setup

1. `pip install -r requirements.txt`
2. `cp .env.example .env` and add your CurseForge API key
3. `python poc.py`

Get API key from: https://console.curseforge.com/

## What it does

```bash
python poc.py
```

## Configuration

- Connects to CurseForge API
- Fetches latest file for a mod  
- Downloads the file

## Configuration

Set in `.env` file:
- `CURSEFORGE_API_KEY` - Your API key (required)
- `MOD_ID` - Mod ID to check (optional)
- `DOWNLOAD_PATH` - Download location (optional)