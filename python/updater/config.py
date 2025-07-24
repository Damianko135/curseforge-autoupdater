
import os
import sys
from pathlib import Path
from dotenv import load_dotenv
from typing import Dict, Any, Optional

# Load environment variables
load_dotenv()


def get_config() -> Dict[str, Any]:
    """
    Load and validate configuration from environment variables.
    Returns a dictionary with all config values.
    """
    config = {
        "api_key": os.getenv("CURSEFORGE_API_KEY"),
        "mod_id": os.getenv("MOD_ID", "1300837"),
        "download_path": Path(os.getenv("DOWNLOAD_PATH", "./downloads")),
        "game_id": int(os.getenv("GAME_ID", "432")),
        "mod_loader": os.getenv("MOD_LOADER"),
        "minecraft_version": os.getenv("MINECRAFT_VERSION"),
        "extract_path": Path(os.getenv("EXTRACT_PATH", "./extracted")),
    }

    # Validate mod_id is numeric
    try:
        config["mod_id"] = str(int(config["mod_id"]))
    except ValueError:
        print(f"\u274c Invalid MOD_ID: {config['mod_id']}. Must be a number.")
        sys.exit(1)

    return config


def validate_config(config: Dict[str, Any]) -> bool:
    """
    Validate the loaded configuration. Returns True if valid, False otherwise.
    """
    errors = []

    if not config["api_key"]:
        errors.append("CURSEFORGE_API_KEY is required")

    if not config["mod_id"]:
        errors.append("MOD_ID is required")

    try:
        int(config["mod_id"])
    except ValueError:
        errors.append("MOD_ID must be a valid number")

    if errors:
        print("\u274c Configuration errors:")
        for error in errors:
            print(f"   - {error}")
        return False

    return True


def print_config(config: Dict[str, Any]) -> None:
    """
    Print the current configuration (with API key masked).
    """
    print("\U0001F4CB Current Configuration:")

    api_key = config["api_key"]
    if api_key:
        masked_key = f"{'*' * (len(api_key) - 4)}{api_key[-4:]}"
    else:
        masked_key = "\u274c Not set"

    print(f"   API Key: {masked_key}")
    print(f"   Mod ID: {config['mod_id']}")
    print(f"   Game ID: {config['game_id']}")
    print(f"   Download Path: {config['download_path']}")
    print(f"   Extract Path: {config['extract_path']}")

    if config["mod_loader"]:
        print(f"   Mod Loader: {config['mod_loader']}")
    if config["minecraft_version"]:
        print(f"   Minecraft Version: {config['minecraft_version']}")


def create_example_env() -> None:
    """
    Create an example .env file if it doesn't exist.
    """
    env_example_path = Path(".env.example")
    env_path = Path(".env")

    if not env_path.exists() and env_example_path.exists():
        print("\U0001F4DD Creating .env file from .env.example...")
        try:
            with open(env_example_path, 'r') as src, open(env_path, 'w') as dst:
                dst.write(src.read())
            print("\u2705 Created .env file. Please edit it with your API key.")
        except IOError as e:
            print(f"\u274c Could not create .env file: {e}")
    elif not env_example_path.exists():
        print("\u26A0\uFE0F  No .env.example file found to copy from.")


def require_env(var: str, example: Optional[str] = None) -> str:
    """
    Helper to require an environment variable, with an optional example value.
    """
    value = os.getenv(var)
    if not value:
        msg = f"\u274c Required environment variable '{var}' is missing."
        if example:
            msg += f" Example: {var}={example}"
        print(msg)
        sys.exit(1)
    return value
