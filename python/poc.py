"""
Simple CurseForge Auto-Updater PoC
"""

import requests
import os
import json
import zipfile
from pathlib import Path
from dotenv import load_dotenv
from datetime import datetime
import traceback

# Load environment variables
load_dotenv()

def get_mod_files(api_key, mod_id):
    """Get mod files from CurseForge API."""
    url = f"https://api.curseforge.com/v1/mods/{mod_id}/files"
    headers = {
        "Accept": "application/json",
        "x-api-key": api_key,
        "User-Agent": "CurseForge Auto-Updater PoC/1.0"
    }
    
    print(f"Making API request to: {url}")
    print(f"Headers: {dict(headers)}")
    
    try:
        response = requests.get(url, headers=headers)
        print(f"Response status code: {response.status_code}")
        print(f"Response headers: {dict(response.headers)}")
        
        response.raise_for_status()
        
        data = response.json()
        print(f"Response JSON keys: {list(data.keys())}")
        
        # Check pagination info
        pagination = data.get("pagination", {})
        if pagination:
            print(f"Pagination: {pagination}")
        
        files = data.get("data", [])
        print(f"Number of files found: {len(files)}")
        
        # Print the full response for debugging (first time)
        if len(files) == 0:
            print("Full API response:")
            print("Writing to 'full_response.json'")
            with open('full_response.json', 'w') as f:
                json.dump(data, f, indent=2)
        else:
            print("Full API response (truncated):")
            print(json.dumps(data, indent=2)[:500] + "...")
            print("Writing to 'full_response.json'")
            with open('full_response.json', 'w') as f:
                json.dump(data, f, indent=2) 
        
        if files:
            print("Sample file info:")
            for i, file in enumerate(files[:3]):  # Show first 3 files
                print(f"  File {i+1}: {file.get('fileName')} (ID: {file.get('id')}, Date: {file.get('fileDate')})")
        
        return files
        
    except requests.exceptions.RequestException as e:
        print(f"Request failed: {e}")
        if hasattr(e, 'response') and e.response is not None:
            print(f"Error response: {e.response.text}")
        return []

def get_latest_file(files):
    """Get the latest file from the list."""
    if not files:
        return None
    return max(files, key=lambda x: x.get("fileDate", ""))

def download_file(file_info, api_key, download_path):
    """Download the file."""
    download_url = file_info.get("downloadUrl")
    file_name = file_info.get("fileName")
    
    if not download_url:
        print("No download URL available")
        return False
    
    download_path.mkdir(parents=True, exist_ok=True)
    file_path = download_path / file_name
    
    headers = {
        "x-api-key": api_key,
        "User-Agent": "CurseForge Auto-Updater PoC/1.0"
    }
    response = requests.get(download_url, headers=headers, stream=True, timeout=60)
    response.raise_for_status()
    
    with open(file_path, 'wb') as f:
        for chunk in response.iter_content(chunk_size=8192):
            if chunk:
                f.write(chunk)
    
    print(f"Downloaded: {file_path}")
    return True

def get_mod_files_with_params(api_key, mod_id, params):
    """Get mod files with additional parameters."""
    url = f"https://api.curseforge.com/v1/mods/{mod_id}/files"
    headers = {"Accept": "application/json", "x-api-key": api_key}
    
    print(f"Making API request with params: {params}")
    
    try:
        response = requests.get(url, headers=headers, params=params)
        print(f"Response status: {response.status_code}")
        
        if response.status_code == 200:
            data = response.json()
            files = data.get("data", [])
            print(f"Files found with params: {len(files)}")
            return files
        else:
            print(f"Request failed: {response.text}")
            return []
            
    except requests.exceptions.RequestException as e:
        print(f"Error with params: {e}")
        return []

def get_mod_files_raw(api_key, mod_id):
    """Get mod files with minimal processing."""
    url = f"https://api.curseforge.com/v1/mods/{mod_id}/files"
    headers = {"Accept": "application/json", "x-api-key": api_key}
    
    print(f"Making raw API request...")
    
    try:
        response = requests.get(url, headers=headers)
        print(f"Raw response status: {response.status_code}")
        print(f"Raw response text: {response.text[:500]}...")  # First 500 chars
        
        if response.status_code == 200:
            data = response.json()
            return data.get("data", [])
        else:
            return []
            
    except Exception as e:
        print(f"Raw request error: {e}")
        return []

def load_download_metadata(download_path):
    """Load metadata about previously downloaded files."""
    metadata_file = download_path / "download_metadata.json"
    if metadata_file.exists():
        with open(metadata_file, 'r') as f:
            return json.load(f)
    return {}

def save_download_metadata(download_path, metadata):
    """Save metadata about downloaded files."""
    metadata_file = download_path / "download_metadata.json"
    with open(metadata_file, 'w') as f:
        json.dump(metadata, f, indent=2)

def is_download_needed(file_info, download_path, metadata):
    """Check if a file needs to be downloaded."""
    file_name = file_info.get("fileName")
    file_id = file_info.get("id")
    file_date = file_info.get("fileDate")
    
    # Check if file exists locally
    local_file_path = download_path / file_name
    if not local_file_path.exists():
        print(f"  ➤ File not found locally: {file_name}")
        return True, "File not downloaded yet"
    
    # Check metadata
    if str(file_id) not in metadata:
        print(f"  ➤ No metadata found for file ID {file_id}")
        return True, "No metadata for this file"
    
    local_metadata = metadata[str(file_id)]
    local_date = local_metadata.get("fileDate")
    
    if local_date != file_date:
        print(f"  ➤ Date mismatch - Local: {local_date}, Remote: {file_date}")
        return True, f"File updated (was: {local_date}, now: {file_date})"
    
    # Check file hash if available
    remote_hash = None
    for hash_info in file_info.get("hashes", []):
        if hash_info.get("algo") == 1:  # SHA-1
            remote_hash = hash_info.get("value")
            break
    
    if remote_hash and local_metadata.get("hash") != remote_hash:
        print(f"  ➤ Hash mismatch - Local: {local_metadata.get('hash')}, Remote: {remote_hash}")
        return True, "File hash changed"
    
    print(f"  ✓ File up to date: {file_name}")
    return False, "File is current"

def record_download(file_info, download_path, metadata):
    """Record a successful download in metadata."""
    file_id = str(file_info.get("id"))
    file_name = file_info.get("fileName")
    
    # Get hash
    file_hash = None
    for hash_info in file_info.get("hashes", []):
        if hash_info.get("algo") == 1:  # SHA-1
            file_hash = hash_info.get("value")
            break
    
    metadata[file_id] = {
        "fileName": file_name,
        "fileDate": file_info.get("fileDate"),
        "downloadedAt": datetime.now().isoformat(),
        "hash": file_hash,
        "fileLength": file_info.get("fileLength")
    }
    
    save_download_metadata(download_path, metadata)
    print(f"  ✓ Recorded download metadata for {file_name}")

def main():
    """Main function for the CurseForge updater PoC."""
    print("CurseForge Auto-Updater PoC")
    print("=" * 40)
    
    # Get configuration from environment
    api_key = os.getenv('CURSEFORGE_API_KEY')
    mod_id = os.getenv('MOD_ID', '1300837')  # Default to some mod
    download_path = Path(os.getenv('DOWNLOAD_PATH', './downloads'))
    
    print(f"Configuration:")
    if api_key:
        print(f"  API key: {'*' * (len(api_key) - 4)}{api_key[-4:]}")
    else:
        print(f"  API key: None")
    print(f"  Mod ID: {mod_id}")
    print(f"  Download path: {download_path}")
    print()
    
    if not api_key:
        print("❌ No API key found. Create a .env file with CURSEFORGE_API_KEY=your_key")
        return
    
    # First, let's test if we can get basic mod info
    print("Step 1: Testing mod info endpoint...")
    try:
        mod_info_url = f"https://api.curseforge.com/v1/mods/{mod_id}"
        headers = {"Accept": "application/json", "x-api-key": api_key, "User-Agent": "CurseForge Auto-Updater PoC/1.0"}
        
        response = requests.get(mod_info_url, headers=headers)
        print(f"Mod info response: {response.status_code}")
        
        if response.status_code == 200:
            mod_data = response.json().get("data", {})
            print(f"✓ Mod found: {mod_data.get('name', 'Unknown')} by {mod_data.get('authors', [{}])[0].get('name', 'Unknown') if mod_data.get('authors') else 'Unknown'}")
            print(f"  Game ID: {mod_data.get('gameId')}")
            print(f"  Category: {mod_data.get('classId')}")
        else:
            print(f"❌ Failed to get mod info: {response.text}")
            return
            
    except Exception as e:
        print(f"❌ Error getting mod info: {e}")
        return
    
    print()
    print("Step 2: Fetching mod files...")
    
    try:
        files = get_mod_files(api_key, mod_id)
        
        if not files:
            print("❌ No files found")
            print("This could mean:")
            print("  - The mod has no public files")
            print("  - The mod ID is incorrect") 
            print("  - API permissions issue")
            print("  - Files might be in a different game/category")
            print()
            print("Let's try some alternative approaches...")
            
            # Try with different parameters
            print("Trying with gameId parameter...")
            game_files = get_mod_files_with_params(api_key, mod_id, {"gameId": 432})
            if game_files:
                files = game_files
            else:
                print("Still no files with gameId parameter")
                
            if not files:
                print("Trying to get ALL files (no filters)...")
                all_files = get_mod_files_raw(api_key, mod_id)
                if all_files:
                    files = all_files
                else:
                    print("No files found even with no filters")
                    return
        
        print(f"✓ Found {len(files)} files")
        
        # Get latest file
        latest_file = get_latest_file(files)
        if not latest_file:
            print("❌ No latest file found")
            return
        
        print()
        print("Step 3: Latest file info:")
        print(f"  Name: {latest_file.get('fileName')}")
        print(f"  Display Name: {latest_file.get('displayName')}")
        print(f"  Date: {latest_file.get('fileDate')}")
        print(f"  Size: {latest_file.get('fileLength', 0)} bytes")
        print(f"  Download URL: {'Available' if latest_file.get('downloadUrl') else 'Not available'}")
        
        # Check if download is needed
        print()
        print("Step 4: Checking if download is needed...")
        metadata = load_download_metadata(download_path)
        print(f"Found metadata for {len(metadata)} previously downloaded files")
        
        needs_download, reason = is_download_needed(latest_file, download_path, metadata)
        
        if needs_download:
            print(f"📥 Download needed: {reason}")
            print()
            print("Step 5: Downloading...")
            if download_file(latest_file, api_key, download_path):
                print("✓ Download completed!")
                record_download(latest_file, download_path, metadata)
                print("✓ PoC completed successfully!")
        else:
            print("✓ PoC completed - everything up to date!")
        
    except Exception as e:
        print(f"❌ Error: {e}")
        traceback.print_exc()

if __name__ == "__main__":
    main()