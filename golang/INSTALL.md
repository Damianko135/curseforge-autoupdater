# Installation Guide - CurseForge Auto-Updater

This guide walks you through installing all dependencies needed to build and develop the CurseForge Auto-Updater.

## Prerequisites

### 1. Go Programming Language
First, you need Go installed on your system.

**Windows:**
```powershell
# Using Chocolatey
choco install golang

# Or download from: https://golang.org/dl/
```

**macOS:**
```bash
# Using Homebrew
brew install go

# Or download from: https://golang.org/dl/
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# Or download from: https://golang.org/dl/
```

Verify installation:
```bash
go version
```

## Required Dependencies

### 2. Mage (Build Tool)
Mage is our primary build automation tool.

```bash
go install github.com/magefile/mage@latest
```

Verify installation:
```bash
mage -version
```

## Optional Development Dependencies

### 3. Air (Live Reload - Optional)
Only needed if you want to use `mage dev` for live reload during development.

```bash
go install github.com/air-verse/air@latest
```

Verify installation:
```bash
air -v
```

### 4. golangci-lint (Code Linting - Optional)
Only needed if you want to use `mage lint` for code quality checks.

**Windows:**
```powershell
# Using Chocolatey
choco install golangci-lint

# Or using Go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**macOS:**
```bash
# Using Homebrew
brew install golangci-lint

# Or using Go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**Linux:**
```bash
# Using Go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Or download binary from: https://github.com/golangci/golangci-lint/releases
```

Verify installation:
```bash
golangci-lint version
```

### 5. gosec (Security Scanner - Optional)
Only needed if you want to use `mage security` for security checks.

```bash
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
```

Verify installation:
```bash
gosec -version
```

## Project Setup

After installing the dependencies above:

### 1. Clone and Setup
```bash
# Clone the repository
git clone https://github.com/Damianko135/curseforge-autoupdate.git
cd curseforge-autoupdate/golang

# Install project dependencies
mage deps

# Build the application
mage build
```

### 2. Verify Everything Works
```bash
# Run all available mage commands
mage -l

# Test the complete build pipeline
mage all
```

## Quick Start Summary

**Minimal setup (required only):**
```bash
# 1. Install Go (see above for your OS)
# 2. Install Mage
go install github.com/magefile/mage@latest

# 3. Setup project
git clone https://github.com/Damianko135/curseforge-autoupdate.git
cd curseforge-autoupdate/golang
mage deps build
```

**Full development setup (with all tools):**
```bash
# After Go is installed
go install github.com/magefile/mage@latest
go install github.com/air-verse/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Setup project
git clone https://github.com/Damianko135/curseforge-autoupdate.git
cd curseforge-autoupdate/golang
mage deps build
```

## Available Mage Commands

After setup, you can use these commands:

**Core Commands:**
- `mage build` - Build for current platform
- `mage buildAll` - Build for all platforms
- `mage run` - Run the application
- `mage clean` - Clean build artifacts

**Development Commands:**
- `mage dev` - Live reload development (requires air)
- `mage fmt` - Format code
- `mage test` - Run tests
- `mage lint` - Run linter (requires golangci-lint)
- `mage security` - Security scan (requires gosec)

**Composite Commands:**
- `mage all` - Complete pipeline: deps → fmt → test → build
- `mage ci` - CI pipeline: deps → fmt → lint → test → buildAll

## Troubleshooting

**If `go install` fails:**
- Make sure `$GOPATH/bin` is in your `$PATH`
- On Windows, this is usually `%USERPROFILE%\go\bin`

**If mage commands fail:**
- Ensure you're in the `golang/` directory
- Run `mage deps` first to install project dependencies

**If air fails to start:**
- Make sure you have write permissions in the project directory
- Check that `tmp/` directory can be created
