# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

daemon-control is a macOS daemon management tool that simplifies creating, installing, and managing LaunchAgent daemons. It works with YAML configuration files to generate macOS plist files and provides comprehensive daemon lifecycle management.

## Build and Development Commands

### Building
```bash
# Build for current platform
make build

# Build for all platforms (darwin/amd64, darwin/arm64, linux/amd64, linux/arm64)
make build-all

# Build with debug symbols
make build-debug

# Install to /usr/local/bin
make install
```

### Testing
```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run a single test
go test -v -run TestFunctionName ./path/to/package
```

### Linting and Formatting
```bash
# Run golangci-lint (v2.3.0)
make lint
golangci-lint run

# Format code
make fmt
gofmt -w .
goimports -local github.com/mjmorales/daemon-control -w .

# Run go vet
make vet
```

### Dependencies
```bash
# Update dependencies
make deps

# Verify dependencies
make deps-verify
```

## Architecture and Code Structure

### Core Components

1. **Core Configuration Manager** (`internal/core/`)
   - Singleton pattern via `GetManager()` in `singleton.go`
   - Manages daemon-control's own configuration (not daemon configs)
   - Default config location: `~/.daemon-control/core.config.yaml`
   - Handles logging configuration and path resolution

2. **Daemon Configuration** (`internal/config/`)
   - `schema.go`: Defines the YAML structure for daemon definitions
   - `loader.go`: Loads and validates daemon configurations
   - Supports complex daemon setups including calendar intervals, keep-alive conditions, and socket activation

3. **Plist Generation** (`internal/plist/`)
   - `generator.go`: Converts YAML daemon configs to macOS plist XML
   - `types.go`: Defines plist XML structures
   - Handles all macOS-specific daemon features

4. **Command Structure** (`cmd/`)
   - Each command is a separate file (e.g., `install.go`, `start.go`)
   - All commands use the singleton core manager for configuration
   - Commands interact with launchctl for daemon management

5. **Utilities** (`internal/utils/`)
   - Common functions for plist path resolution
   - Wrappers around launchctl commands
   - File operations with proper error handling

### Key Design Patterns

1. **Configuration Hierarchy**:
   - Core config (daemon-control settings) → stored in `~/.daemon-control/`
   - Daemon config (user's daemons) → customizable path, default `./daemons.yaml`
   - Generated plists → configurable output directory

2. **Error Handling**:
   - All file operations check errors
   - Commands use `exec.CommandContext` with context.Background()
   - Proper file permission settings (0600 for files, 0700/0750 for directories)

3. **Logging**:
   - Uses zerolog throughout
   - Configurable via core config (level and format)
   - Contextual logging with structured fields

## CI/CD and Release Process

- Uses GitHub Actions for CI (`.github/workflows/ci.yml`)
- Semantic-release for automated versioning (config in `.releaserc.json`)
- GoReleaser for building release binaries (`.goreleaser.yml`)
- Conventional commits required (enforced by commitlint)
- Branch strategy: `master` (main), `beta`, `alpha`

## Important Configuration Notes

- The project uses golangci-lint v2.3.0 with custom exclusion rules in `.golangci.yml`
- GoReleaser v2 configuration format
- Repository URL: `https://github.com/mjmorales/daemon-control.git`
- Package import path: `github.com/mjmorales/daemon-control`

## Testing Daemon Operations

When testing daemon operations locally:
1. Use the `--dry-run` flag if available
2. Test with non-system daemons first
3. Check daemon status with `daemon-control status <daemon-name>`
4. View logs with `daemon-control logs <daemon-name>` or `tail <daemon-name>`

## Common Development Tasks

### Adding a New Command
1. Create new file in `cmd/` directory
2. Define cobra.Command with proper Use, Short, and Long descriptions
3. Add command to rootCmd in init() function
4. Use `core.GetManager()` for configuration access
5. Use `internal/utils` for common operations

### Modifying Daemon Schema
1. Update `internal/config/schema.go` with new fields
2. Update validation in `internal/config/loader.go`
3. Update plist generation in `internal/plist/generator.go`
4. Add example to `daemons.example.yaml`