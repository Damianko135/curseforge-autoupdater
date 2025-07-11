#!/usr/bin/env python3
"""
Command-line interface for CurseForge Auto-Updater
"""

import sys
import argparse
from pathlib import Path

# Add the package to the path
sys.path.insert(0, str(Path(__file__).parent))

from updater import main, get_config, api
from updater.config import validate_config, print_config, create_example_env
from updater.utils import validate_mod_id

def cli_main():
    """Command-line interface main function."""
    parser = argparse.ArgumentParser(
        description="CurseForge Auto-Updater - Download and update CurseForge mods automatically",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s                          # Run with default configuration
  %(prog)s --mod-id 123456          # Override mod ID
  %(prog)s --config                 # Show current configuration
  %(prog)s --validate-key           # Validate API key
  %(prog)s --create-env             # Create .env file from example
        """
    )
    
    parser.add_argument(
        "--mod-id", 
        type=str,
        help="Override the mod ID from environment"
    )
    
    parser.add_argument(
        "--download-path",
        type=Path,
        help="Override the download path from environment"
    )
    
    parser.add_argument(
        "--config",
        action="store_true",
        help="Show current configuration and exit"
    )
    
    parser.add_argument(
        "--validate-key",
        action="store_true", 
        help="Validate API key and exit"
    )
    
    parser.add_argument(
        "--create-env",
        action="store_true",
        help="Create .env file from .env.example and exit"
    )
    
    parser.add_argument(
        "--version",
        action="version",
        version="CurseForge Auto-Updater 1.0.0"
    )
    
    args = parser.parse_args()
    
    # Handle special commands
    if args.create_env:
        create_example_env()
        return 0
    
    # Load configuration
    config = get_config()
    
    # Override with command-line arguments
    if args.mod_id:
        try:
            config["mod_id"] = str(validate_mod_id(args.mod_id))
        except ValueError as e:
            print(f"❌ {e}")
            return 1
    
    if args.download_path:
        config["download_path"] = args.download_path
    
    if args.config:
        print_config(config)
        return 0
    
    if args.validate_key:
        if not config["api_key"]:
            print("❌ No API key found in configuration")
            return 1
        
        print("🔑 Validating API key...")
        try:
            if api.validate_api_key(config["api_key"]):
                print("✅ API key is valid")
                return 0
            else:
                print("❌ API key is invalid")
                return 1
        except Exception as e:
            print(f"❌ Error validating API key: {e}")
            return 1
    
    # Validate configuration
    if not validate_config(config):
        return 1
    
    # Update environment variables with overrides
    import os
    if args.mod_id:
        os.environ["MOD_ID"] = config["mod_id"]
    if args.download_path:
        os.environ["DOWNLOAD_PATH"] = str(config["download_path"])
    
    # Run the main application
    return main()

if __name__ == "__main__":
    sys.exit(cli_main())
