# CurseForge Auto-Updater Setup and Run Script (PowerShell)
# This script sets up a Python virtual environment and runs the PoC

param(
    [switch]$Help
)

if ($Help) {
    Write-Host "CurseForge Auto-Updater Setup Script"
    Write-Host "Usage: .\run.ps1"
    Write-Host "This script will:"
    Write-Host "  1. Create a Python virtual environment (if not exists)"
    Write-Host "  2. Activate the virtual environment"
    Write-Host "  3. Install dependencies from requirements.txt"
    Write-Host "  4. Run the poc.py script"
    exit 0
}

Write-Host "CurseForge Auto-Updater - Setup and Run" -ForegroundColor Cyan
Write-Host "=======================================" -ForegroundColor Cyan

# Check if Python is available
$pythonCmd = $null
if (Get-Command python -ErrorAction SilentlyContinue) {
    $pythonCmd = "python"
} elseif (Get-Command python3 -ErrorAction SilentlyContinue) {
    $pythonCmd = "python3"
} else {
    Write-Host "❌ Error: Python is not installed or not in PATH" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Using Python: $pythonCmd" -ForegroundColor Green

# Check if virtual environment already exists
if (Test-Path "venv\Scripts\activate.ps1") {
    Write-Host "✓ Virtual environment already exists" -ForegroundColor Green
} else {
    Write-Host "📦 Creating virtual environment..." -ForegroundColor Yellow
    & $pythonCmd -m venv venv
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to create virtual environment" -ForegroundColor Red
        exit 1
    }
    Write-Host "✓ Virtual environment created" -ForegroundColor Green
}

# Activate virtual environment
Write-Host "🔄 Activating virtual environment..." -ForegroundColor Yellow
& "venv\Scripts\Activate.ps1"

# Verify activation
if ($env:VIRTUAL_ENV) {
    Write-Host "✓ Virtual environment activated: $(Split-Path $env:VIRTUAL_ENV -Leaf)" -ForegroundColor Green
} else {
    Write-Host "❌ Failed to activate virtual environment" -ForegroundColor Red
    exit 1
}

# Install dependencies
if (Test-Path "requirements.txt") {
    Write-Host "📋 Installing dependencies from requirements.txt..." -ForegroundColor Yellow
    pip install --upgrade pip
    pip install -r requirements.txt
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to install dependencies" -ForegroundColor Red
        exit 1
    }
    Write-Host "✓ Dependencies installed" -ForegroundColor Green
} else {
    Write-Host "⚠️  requirements.txt not found, installing basic dependencies..." -ForegroundColor Yellow
    pip install --upgrade pip
    pip install requests python-dotenv
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to install basic dependencies" -ForegroundColor Red
        exit 1
    }
    Write-Host "✓ Basic dependencies installed" -ForegroundColor Green
}

# Check if .env file exists
if (-not (Test-Path ".env")) {
    Write-Host ""
    Write-Host "⚠️  WARNING: .env file not found!" -ForegroundColor Yellow
    Write-Host "   Please create a .env file with your CurseForge API key:" -ForegroundColor Yellow
    Write-Host "   CURSEFORGE_API_KEY=your_api_key_here" -ForegroundColor Gray
    Write-Host "   MOD_ID=1300837" -ForegroundColor Gray
    Write-Host "   DOWNLOAD_PATH=./downloads" -ForegroundColor Gray
    Write-Host ""
}

# Run the PoC
if (Test-Path "poc.py") {
    Write-Host ""
    Write-Host "🚀 Running CurseForge Auto-Updater PoC..." -ForegroundColor Cyan
    Write-Host "----------------------------------------" -ForegroundColor Cyan
    python poc.py
} elseif (Test-Path "PoC.py") {
    Write-Host ""
    Write-Host "🚀 Running CurseForge Auto-Updater PoC..." -ForegroundColor Cyan
    Write-Host "----------------------------------------" -ForegroundColor Cyan
    python PoC.py
} else {
    Write-Host "❌ Python file poc.py/PoC.py not found." -ForegroundColor Red
    Write-Host "   Please ensure the PoC script exists in the current directory." -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "✅ Script completed!" -ForegroundColor Green
