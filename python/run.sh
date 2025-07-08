#!/bin/bash
set -e

# Detect platform-specific venv activation script
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    ACTIVATE_PATH="venv/Scripts/activate"
else
    ACTIVATE_PATH="venv/bin/activate"
fi

# Check for Python 3
if ! command -v python3 &>/dev/null; then
    echo "❌ Python 3 not found."
    exit 1
fi

# Check for valid venv by looking for activation script
if [ -f "$ACTIVATE_PATH" ]; then
    echo "✔️ Virtual environment already exists and looks valid. Skipping creation."
else
    echo "🔧 Creating virtual environment..."
    python3 -m venv venv
fi

# Activate venv
source "$ACTIVATE_PATH"

# Install dependencies
if [ -f "requirements.txt" ]; then
    echo "📦 Installing requirements..."
    pip install --upgrade pip
    pip install -r requirements.txt
else
    echo "⚠️ No requirements.txt found. Skipping dependencies."
fi

# Run PoC
if [ -f "PoC.py" || -f "poc.py" ]; then
    echo "🚀 Running PoC.py..."

    python3 PoC.py || python3 poc.py
else
    echo "❌ Python file PoC.py/poc.py not found."
    echo "If named differently, please rename it in this script."
    exit 1
fi
