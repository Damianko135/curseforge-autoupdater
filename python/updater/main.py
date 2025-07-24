
import sys
import traceback
from typing import Optional
from .config import get_config
from . import api, downloader, utils

try:
    import logging
    logger = logging.getLogger("updater")
    logging.basicConfig(level=logging.INFO)
except ImportError:
    logger = None


def log(msg: str, level: str = "info"):
    if logger:
        getattr(logger, level, logger.info)(msg)
    else:
        print(msg)


def main() -> int:
    """
    Main entry point for the CurseForge Auto-Updater.
    Returns 0 on success, 1 on error.
    """
    log("🔧 CurseForge Auto-Updater v1.0.0", "info")
    log("=" * 50, "info")

    try:
        # Load configuration
        config = get_config()
        api_key: Optional[str] = config["api_key"]
        mod_id: str = config["mod_id"]
        download_path = config["download_path"]

        log(f"📋 Configuration:")
        if api_key:
            log(f"   API Key: {'*' * (len(api_key) - 4)}{api_key[-4:]}")
        else:
            log(f"   API Key: ❌ Missing")
        log(f"   Mod ID: {mod_id}")
        log(f"   Download Path: {download_path}")
        log("", "info")

        if not api_key:
            log("❌ API key missing. Please:", "error")
            log("   1. Copy .env.example to .env", "error")
            log("   2. Add your CurseForge API key", "error")
            log("   3. Get API key from: https://console.curseforge.com/", "error")
            return 1

        # Get mod information
        log("🔍 Fetching mod information...", "info")
        try:
            mod_info = api.get_mod_info(api_key, mod_id)
            mod_name = mod_info.get('name', 'Unknown')
            mod_authors = mod_info.get('authors', [])
            author_name = mod_authors[0].get('name', 'Unknown') if mod_authors else 'Unknown'

            log(f"✅ Found mod: {mod_name}")
            log(f"   Author: {author_name}")
            log(f"   Game ID: {mod_info.get('gameId', 'Unknown')}")
            log("", "info")
        except Exception as e:
            log(f"❌ Failed to fetch mod info: {e}", "error")
            return 1

        # Get mod files
        log("📂 Fetching mod files...", "info")
        try:
            files = api.get_mod_files(api_key, mod_id)
            log(f"✅ Found {len(files)} files")

            if not files:
                log("❌ No files found for this mod", "error")
                return 1

        except Exception as e:
            log(f"❌ Failed to fetch mod files: {e}", "error")
            return 1

        # Get latest file
        latest_file = utils.get_latest_file(files)
        if not latest_file:
            log("❌ No latest file found.", "error")
            return 1

        log(f"📄 Latest file: {latest_file.get('fileName')}")
        log(f"   Display Name: {latest_file.get('displayName', 'N/A')}")
        log(f"   Date: {latest_file.get('fileDate')}")
        log(f"   Size: {latest_file.get('fileLength', 0):,} bytes")
        log("", "info")

        # Check if download is needed
        log("🔄 Checking if update is needed...", "info")
        metadata = downloader.load_metadata(download_path)
        needs_download, reason = utils.is_download_needed(latest_file, download_path, metadata)

        if needs_download:
            log(f"📥 Download needed: {reason}")
            log("⬇️  Starting download...", "info")

            try:
                success = downloader.download_file(latest_file, api_key, download_path)
                if success:
                    downloader.record_download(latest_file, download_path, metadata)
                    log("✅ Download completed and recorded successfully!", "info")
                    return 0
                else:
                    log("❌ Download failed", "error")
                    return 1
            except Exception as e:
                log(f"❌ Download error: {e}", "error")
                return 1
        else:
            log(f"✅ File is up to date: {reason}", "info")
            return 0

    except KeyboardInterrupt:
        log("\n⚠️  Operation cancelled by user", "warning")
        return 1
    except Exception as e:
        log(f"\n❌ Unexpected error: {e}", "error")
        log("🐛 Full traceback:", "error")
        traceback.print_exc()
        return 1


if __name__ == "__main__":
    sys.exit(main())
