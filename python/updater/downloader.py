import requests
import json
import sys
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, Optional

def download_file(
    file_info: Dict[str, Any],
    api_key: str,
    download_path: Any
) -> bool:
    """
    Download a file with progress tracking.
    Returns True on success, False on failure.
    """
    download_url = file_info.get("downloadUrl")
    file_name = file_info.get("fileName")
    file_length = file_info.get("fileLength", 0)
    
    if not download_url:
        print("❌ No download URL available")
        return False
    
    # Ensure download directory exists
    download_path = Path(download_path)
    download_path.mkdir(parents=True, exist_ok=True)
    file_path = download_path / file_name
    
    headers = {
        "x-api-key": api_key,
        "User-Agent": "CurseForge Auto-Updater/1.0"
    }
    
    try:
        print(f"📥 Downloading {file_name}...")
        response = requests.get(download_url, headers=headers, stream=True, timeout=30)
        response.raise_for_status()
        
        # Download with progress
        downloaded = 0
        with open(file_path, 'wb') as f:
            for chunk in response.iter_content(chunk_size=8192):
                if chunk:
                    f.write(chunk)
                    downloaded += len(chunk)
                    
                    # Show progress
                    if file_length > 0:
                        progress = (downloaded / file_length) * 100
                        print(f"\r   Progress: {progress:.1f}% ({downloaded:,}/{file_length:,} bytes)", end='')
                    else:
                        print(f"\r   Downloaded: {downloaded:,} bytes", end='')
        
        print()  # New line after progress
        print(f"✅ Successfully downloaded: {file_path}")
        return True
        
    except requests.exceptions.RequestException as e:
        print(f"\n❌ Download failed: {e}")
        # Clean up partial download
        if file_path.exists():
            file_path.unlink()
        return False
    except Exception as e:
        print(f"\n❌ Unexpected download error: {e}")
        # Clean up partial download
        if file_path.exists():
            file_path.unlink()
        return False

def load_metadata(download_path: Any) -> Dict[str, Any]:
    """
    Load download metadata from JSON file.
    Returns a dictionary of metadata.
    """
    download_path = Path(download_path)
    metadata_file = download_path / "download_metadata.json"
    
    if metadata_file.exists():
        try:
            with open(metadata_file, 'r') as f:
                return json.load(f)
        except (json.JSONDecodeError, IOError) as e:
            print(f"⚠️  Warning: Could not load metadata: {e}")
            return {}
    return {}

def save_metadata(download_path: Any, metadata: Dict[str, Any]) -> None:
    """
    Save download metadata to JSON file.
    """
    download_path = Path(download_path)
    download_path.mkdir(parents=True, exist_ok=True)
    metadata_file = download_path / "download_metadata.json"
    
    try:
        with open(metadata_file, 'w') as f:
            json.dump(metadata, f, indent=2)
    except IOError as e:
        print(f"⚠️  Warning: Could not save metadata: {e}")

def record_download(
    file_info: Dict[str, Any],
    download_path: Any,
    metadata: Dict[str, Any]
) -> None:
    """
    Record a successful download in metadata.
    """
    file_id = str(file_info.get("id"))
    file_name = file_info.get("fileName")
    
    # Extract file hash if available
    file_hash = None
    for hash_info in file_info.get("hashes", []):
        if hash_info.get("algo") == 1:  # SHA-1
            file_hash = hash_info.get("value")
            break
    
    metadata[file_id] = {
        "fileName": file_name,
        "fileDate": file_info.get("fileDate"),
        "downloadedAt": datetime.now().isoformat(),
        "fileLength": file_info.get("fileLength"),
        "hash": file_hash,
        "displayName": file_info.get("displayName")
    }
    
    save_metadata(download_path, metadata)
    print(f"📝 Recorded download metadata for {file_name}")

def cleanup_old_downloads(download_path: Any, keep_count: int = 5) -> None:
    """
    Clean up old downloaded files, keeping only the most recent ones.
    """
    download_path = Path(download_path)
    if not download_path.exists():
        return
    
    # Get all downloaded files sorted by modification time
    files = []
    for file_path in download_path.iterdir():
        if file_path.is_file() and file_path.suffix in ['.jar', '.zip', '.mrpack']:
            files.append((file_path.stat().st_mtime, file_path))
    
    files.sort(reverse=True)  # Most recent first
    
    # Remove old files
    for _, file_path in files[keep_count:]:
        try:
            file_path.unlink()
            print(f"🗑️  Removed old file: {file_path.name}")
        except OSError as e:
            print(f"⚠️  Could not remove {file_path.name}: {e}")
