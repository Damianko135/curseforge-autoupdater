from pathlib import Path
from datetime import datetime
from typing import Any, Dict, List, Optional, Tuple

def get_latest_file(files: List[Dict[str, Any]]) -> Optional[Dict[str, Any]]:
    """Get the latest file from a list of mod files."""
    if not files:
        return None
    
    # Sort by file date (most recent first)
    return max(files, key=lambda x: x.get("fileDate", ""))

def is_download_needed(
    file_info: Dict[str, Any],
    download_path: Any,
    metadata: Dict[str, Any]
) -> Tuple[bool, str]:
    """
    Check if a file needs to be downloaded.
    
    Returns:
        tuple: (needs_download: bool, reason: str)
    """
    file_name = file_info.get("fileName")
    if not file_name:
        return True, "File info missing fileName"
    file_id = str(file_info.get("id"))
    file_date = file_info.get("fileDate")
    file_length = file_info.get("fileLength", 0)
    # Ensure download_path is a Path
    if not isinstance(download_path, Path):
        download_path = Path(download_path)
    file_path = download_path / file_name
    
    # Check if file exists locally
    if not file_path.exists():
        return True, "File not found locally"
    
    # Check if we have metadata for this file
    if file_id not in metadata:
        return True, "No metadata found for this file"
    
    local_metadata = metadata[file_id]
    
    # Check file date
    local_date = local_metadata.get("fileDate")
    if local_date != file_date:
        return True, f"File updated (local: {local_date}, remote: {file_date})"
    
    # Check file size
    try:
        local_size = file_path.stat().st_size
        if local_size != file_length:
            return True, f"File size mismatch (local: {local_size}, remote: {file_length})"
    except OSError:
        return True, "Could not check local file size"
    
    # Check hash if available
    remote_hash = None
    for hash_info in file_info.get("hashes", []):
        if hash_info.get("algo") == 1:  # SHA-1
            remote_hash = hash_info.get("value")
            break
    
    if remote_hash and local_metadata.get("hash"):
        if local_metadata.get("hash") != remote_hash:
            return True, "File hash mismatch"
    
    return False, "File is up to date"

def format_file_size(size_bytes: int) -> str:
    """Format file size in human readable format."""
    if size_bytes == 0:
        return "0 B"
    
    size_names = ["B", "KB", "MB", "GB"]
    size_index = 0
    size = float(size_bytes)
    
    while size >= 1024.0 and size_index < len(size_names) - 1:
        size /= 1024.0
        size_index += 1
    
    return f"{size:.1f} {size_names[size_index]}"

def format_date(date_string: str) -> str:
    """Format ISO date string to human readable format."""
    try:
        dt = datetime.fromisoformat(date_string.replace('Z', '+00:00'))
        return dt.strftime("%Y-%m-%d %H:%M:%S UTC")
    except (ValueError, AttributeError):
        return date_string

def validate_mod_id(mod_id: Any) -> int:
    """Validate that mod_id is a valid integer."""
    try:
        return int(mod_id)
    except (ValueError, TypeError):
        raise ValueError(f"Invalid mod ID: {mod_id}. Must be a number.")

def get_file_extension_from_url(url: str) -> str:
    """Extract file extension from download URL."""
    if not url:
        return ""
    
    # Remove query parameters and get the last part
    path = url.split('?')[0]
    parts = path.split('.')
    
    if len(parts) > 1:
        return f".{parts[-1]}"
    
    return ""
