
import requests
import json
from typing import Dict, List, Optional, Any

USER_AGENT = "CurseForge Auto-Updater/1.0"
BASE_URL = "https://api.curseforge.com/v1"


class CurseForgeAPIError(Exception):
    """Custom exception for CurseForge API errors."""
    pass


def _make_request(url: str, api_key: str, params: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
    """
    Make a request to the CurseForge API with error handling.
    Returns the parsed JSON response as a dictionary.
    Raises CurseForgeAPIError on error.
    """
    headers = {
        "Accept": "application/json",
        "x-api-key": api_key,
        "User-Agent": USER_AGENT
    }
    try:
        response = requests.get(url, headers=headers, params=params, timeout=30)
        if response.status_code == 401:
            print(f"[API] Invalid API key for {url}")
            raise CurseForgeAPIError("Invalid API key")
        elif response.status_code == 403:
            print(f"[API] Access forbidden for {url}")
            raise CurseForgeAPIError("API access forbidden")
        elif response.status_code == 404:
            print(f"[API] Resource not found: {url}")
            raise CurseForgeAPIError("Resource not found")
        elif response.status_code == 429:
            print(f"[API] Rate limit exceeded for {url}")
            raise CurseForgeAPIError("Rate limit exceeded")
        response.raise_for_status()
        return response.json()
    except requests.exceptions.Timeout:
        print(f"[API] Request timed out: {url}")
        raise CurseForgeAPIError("Request timed out")
    except requests.exceptions.ConnectionError:
        print(f"[API] Connection error: {url}")
        raise CurseForgeAPIError("Connection error")
    except requests.exceptions.RequestException as e:
        print(f"[API] Request failed: {e}")
        raise CurseForgeAPIError(f"Request failed: {e}")
    except json.JSONDecodeError:
        print(f"[API] Invalid JSON response from {url}")
        raise CurseForgeAPIError("Invalid JSON response")


def get_mod_info(api_key: str, mod_id: str) -> Dict[str, Any]:
    """Get information about a specific mod."""
    url = f"{BASE_URL}/mods/{mod_id}"
    data = _make_request(url, api_key)
    return data.get("data", {})


def get_mod_files(
    api_key: str,
    mod_id: str,
    game_version: Optional[str] = None,
    mod_loader: Optional[str] = None
) -> List[Dict[str, Any]]:
    """Get files for a specific mod with optional filtering."""
    url = f"{BASE_URL}/mods/{mod_id}/files"
    params = {}
    if game_version:
        params["gameVersion"] = game_version
    if mod_loader:
        params["modLoaderType"] = mod_loader
    data = _make_request(url, api_key, params)
    return data.get("data", [])


def get_mod_file_info(api_key: str, mod_id: str, file_id: str) -> Dict[str, Any]:
    """Get information about a specific mod file."""
    url = f"{BASE_URL}/mods/{mod_id}/files/{file_id}"
    data = _make_request(url, api_key)
    return data.get("data", {})


def search_mods(
    api_key: str,
    search_filter: str,
    game_id: int = 432,
    class_id: Optional[int] = None,
    sort_field: int = 2
) -> List[Dict[str, Any]]:
    """Search for mods on CurseForge."""
    url = f"{BASE_URL}/mods/search"
    params = {
        "gameId": game_id,
        "searchFilter": search_filter,
        "sortField": sort_field,  # 2 = Popularity
        "sortOrder": "desc"
    }
    if class_id:
        params["classId"] = class_id
    data = _make_request(url, api_key, params)
    return data.get("data", [])


def get_game_info(api_key: str, game_id: int = 432) -> Dict[str, Any]:
    """Get information about a game (default: Minecraft)."""
    url = f"{BASE_URL}/games/{game_id}"
    data = _make_request(url, api_key)
    return data.get("data", {})


def get_mod_categories(api_key: str, game_id: int = 432) -> List[Dict[str, Any]]:
    """Get available mod categories for a game."""
    url = f"{BASE_URL}/categories"
    params = {"gameId": game_id}
    data = _make_request(url, api_key, params)
    return data.get("data", [])


def validate_api_key(api_key: str) -> bool:
    """Validate that the API key works by making a simple request."""
    try:
        get_game_info(api_key)
        return True
    except CurseForgeAPIError:
        return False
