"""
CurseForge Auto-Updater

A Python package for automatically downloading and updating CurseForge mods.
"""

__version__ = "1.0.0"
__author__ = "Damian Korver"

from .main import main
from .config import get_config
from . import api, downloader, utils

__all__ = ["main", "get_config", "api", "downloader", "utils"]
