#!/bin/bash

# CurseForge Auto-Updater Setup and Run Script
# This script sets up a Python virtual environment and runs the PoC

set -e  # Exit on any error

echo "CurseForge Auto-Updater - Setup and Run"
echo "======================================="

# Detect platform-specific venv activation script
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    ACTIVATE_PATH="venv/Scripts/activate"
    PYTHON_CMD="python"
else
    ACTIVATE_PATH="venv/bin/activate"
    PYTHON_CMD="python3"
fi

# Check for Python
if ! command -v $PYTHON_CMD &>/dev/null; then
    echo "❌ $PYTHON_CMD not found. Please install Python."
    exit 1
fi

echo "✓ Using Python: $PYTHON_CMD"

# Check for valid venv by looking for activation script
if [ -f "$ACTIVATE_PATH" ]; then
    echo "✓ Virtual environment already exists and looks valid. Skipping creation."
else
    echo "� Creating virtual environment..."
    $PYTHON_CMD -m venv venv
    echo "✓ Virtual environment created"
fi

# Activate venv
echo "🔄 Activating virtual environment..."
source "$ACTIVATE_PATH"

# Verify activation
if [[ "$VIRTUAL_ENV" != "" ]]; then
    echo "✓ Virtual environment activated: $(basename $VIRTUAL_ENV)"
else
    echo "❌ Failed to activate virtual environment"
    exit 1
fi

# Install dependencies
if [ -f "requirements.txt" ]; then
    echo "� Installing dependencies from requirements.txt..."
    pip install --upgrade pip
    pip install -r requirements.txt
    echo "✓ Dependencies installed"
else
    echo "⚠️  requirements.txt not found, installing basic dependencies..."
    pip install --upgrade pip
    pip install requests python-dotenv
    echo "✓ Basic dependencies installed"
fi

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo ""
    echo "⚠️  WARNING: .env file not found!"
    echo "   Please create a .env file with your CurseForge API key:"
    echo "   CURSEFORGE_API_KEY=your_api_key_here"
    echo "   MOD_ID=1300837"
    echo "   DOWNLOAD_PATH=./downloads"
    echo ""
fi

# Run PoC
if [ -f "poc.py" ]; then
    echo ""
    echo "🚀 Running CurseForge Auto-Updater PoC..."
    echo "----------------------------------------"
    python poc.py
elif [ -f "PoC.py" ]; then
    echo ""
    echo "🚀 Running CurseForge Auto-Updater PoC..."
    echo "----------------------------------------"
    python PoC.py
else
    echo "❌ Python file poc.py/PoC.py not found."
    echo "   Please ensure the PoC script exists in the current directory."
    exit 1
fi

echo ""
echo "✅ Script completed!"
