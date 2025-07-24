
# CurseForge Auto-Update CLI (Golang)

A modern CLI tool to automatically manage Minecraft modpack updates on servers, with full lifecycle management (backup, update, restore) and notification capabilities.

## Vision

The CurseForge Auto-Update CLI aims to be the definitive solution for automated Minecraft modpack management on servers. It is designed to be robust, modular, and user-friendly, empowering server administrators to:


Development is iterative, with a focus on stability, extensibility, and user experience. There are no fixed deadlines—quality and adaptability are the primary goals.

## Features


## Directory Structure

```text
golang/
├── cmd/cli/         # CLI entry and commands
├── internal/api/    # CurseForge API client
├── internal/server/ # Server/backup logic
├── internal/config/ # Config types/templates
├── internal/notification/ # Notification system
├── helper/          # Env, filesystem, version helpers
└── templates/       # Config templates
```

## Quickstart

1. Install Go (1.20+ recommended)
2. Clone this repository and enter the `golang/` directory:

   ```bash
   cd golang
   go mod tidy
   ```

3. Run the CLI:

   ```bash
   go run ./cmd/cli/ --help
   ```

## Example Usage

_Ensure you are in this directory when running the commands below!_

```bash
# Scaffold a config file (recommended)
go run ./cmd/cli/ --init toml

# OR: Create a default config.toml
go run ./cmd/cli/ create-config

# Check if a mod exists (using config/env)
go run ./cmd/cli/ check

# (Stub) Update modpack (not yet implemented)
go run ./cmd/cli/ update
```

## Configuration

Configuration is managed via TOML, YAML, JSON, or .env files. See the `templates/` directory for examples.

## Roadmap

See [PLAN.md](./PLAN.md) for a detailed development plan, including architecture, features, and future enhancements.

## License

[MIT License](../LICENSE)

[trailing newline]
## Configuration
