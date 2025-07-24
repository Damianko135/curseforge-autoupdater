#!/usr/bin/env python3
"""
Setup script for CurseForge Auto-Updater
"""

from pathlib import Path

from setuptools import find_packages, setup

# Read the README file
this_directory = Path(__file__).parent
long_description = (this_directory / "README.md").read_text(encoding="utf-8")

# Read requirements
requirements = []
with open("requirements.txt", "r", encoding="utf-8") as f:
    for line in f:
        line = line.strip()
        if line and not line.startswith("#"):
            requirements.append(line)

setup(
    name="curseforge-autoupdate",
    version="1.0.0",
    description="Automatically download and update CurseForge mods",
    long_description=long_description,
    long_description_content_type="text/markdown",
    author="Damian Korver",
    author_email="your.email@example.com",
    url="https://github.com/Damianko135/curseforge-autoupdate",
    packages=find_packages(include=["updater", "updater.*"]),
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: End Users/Desktop",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "Programming Language :: Python :: 3.13",
        "Topic :: Games/Entertainment",
        "Topic :: Utilities",
    ],
    python_requires=">=3.8",
    install_requires=requirements,
    entry_points={
        "console_scripts": [
            "curseforge-update=cli:cli_main",
            "cf-update=cli:cli_main",
        ],
    },
    include_package_data=True,
    zip_safe=False,
)
