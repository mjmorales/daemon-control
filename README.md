# daemon-control

[![CI](https://github.com/mjmorales/mac-daemon-control/actions/workflows/ci.yml/badge.svg)](https://github.com/mjmorales/mac-daemon-control/actions/workflows/ci.yml)
[![Release](https://github.com/mjmorales/mac-daemon-control/actions/workflows/release.yml/badge.svg)](https://github.com/mjmorales/mac-daemon-control/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mjmorales/mac-daemon-control)](https://goreportcard.com/report/github.com/mjmorales/mac-daemon-control)
[![codecov](https://codecov.io/gh/mjmorales/mac-daemon-control/branch/main/graph/badge.svg)](https://codecov.io/gh/mjmorales/mac-daemon-control)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mjmorales/mac-daemon-control)](https://go.dev/)
[![Latest Release](https://img.shields.io/github/release/mjmorales/mac-daemon-control.svg)](https://github.com/mjmorales/mac-daemon-control/releases/latest)

A powerful and generic daemon control tool for managing macOS LaunchAgent daemons. Simplify the creation, installation, and management of background services on macOS.

## 🚀 Features

- **Generic Daemon Management**: Works with any plist file, not tied to specific services
- **YAML Configuration**: Define daemons in human-readable YAML format
- **Automatic Plist Generation**: Generate valid macOS plist files from YAML configs
- **Comprehensive Controls**: Install, uninstall, start, stop, restart, and monitor daemons
- **Configuration Management**: Centralized configuration with sensible defaults
- **Log Management**: Easy access to daemon logs with `logs` and `tail` commands
- **Editor Integration**: Quick config editing with automatic editor detection
- **Zero Dependencies**: Single binary with no runtime dependencies

## 📦 Installation

### Download Pre-built Binary

Download the latest release from the [releases page](https://github.com/mjmorales/mac-daemon-control/releases).

```bash
# Download and extract (example for macOS arm64)
curl -L https://github.com/mjmorales/mac-daemon-control/releases/latest/download/daemon-control_Darwin_arm64.tar.gz | tar xz
sudo mv daemon-control /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/mjmorales/mac-daemon-control.git
cd mac-daemon-control
make build
sudo make install
```

## 🏃 Quick Start

1. **Initialize configuration**:
   ```bash
   daemon-control config init
   ```

2. **Create your daemon configuration**:
   ```bash
   daemon-control edit
   ```

3. **Add a daemon to the YAML config**:
   ```yaml
   daemons:
     - name: my-service
       label: com.example.my-service
       program_arguments:
         - /usr/local/bin/my-program
         - --config
         - /etc/my-program/config.yml
       working_directory: /var/lib/my-program
       run_at_load: true
       keep_alive:
         successful_exit: false
         crashed: true
   ```

4. **Generate plist files**:
   ```bash
   daemon-control generate
   ```

5. **Install and start your daemon**:
   ```bash
   daemon-control install my-service
   daemon-control start my-service
   ```

## 📖 Usage

### Core Commands

```bash
# Configuration management
daemon-control config init          # Initialize configuration
daemon-control config show          # Show current configuration
daemon-control config set <key> <value>  # Set a config value
daemon-control edit                 # Edit daemon configuration
daemon-control edit --core          # Edit core configuration

# Daemon management
daemon-control list                 # List all available daemons
daemon-control generate             # Generate plist files from YAML
daemon-control install <daemon>     # Install a daemon
daemon-control uninstall <daemon>   # Uninstall a daemon
daemon-control start <daemon>       # Start a daemon
daemon-control stop <daemon>        # Stop a daemon
daemon-control restart <daemon>     # Restart a daemon
daemon-control status <daemon>      # Check daemon status

# Log management
daemon-control logs <daemon>        # Show recent logs
daemon-control tail <daemon>        # Tail logs in real-time
```

### Configuration

The tool uses two configuration files:

1. **Core Config** (`~/.daemon-control/core.config.yaml`): Controls daemon-control behavior
2. **Daemon Config** (customizable path): Defines your daemons

#### Core Configuration Options

- `daemon_config_path`: Path to your daemons YAML file
- `daemons_dir`: Directory for plist files
- `output_dir`: Output directory for generated plists
- `auto_generate_plists`: Auto-copy generated plists to daemons dir
- `log_level`: Logging level (debug, info, warn, error)
- `log_format`: Log format (console or json)

### Example Daemon Configurations

See [daemons.example.yaml](daemons.example.yaml) for comprehensive examples including:
- Simple background services
- Scheduled tasks with calendar intervals
- Socket-activated services
- Services with resource limits
- File watchers
- User-specific daemons

## 🔧 Development

### Prerequisites

- Go 1.21 or later
- macOS (for testing)
- Make

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linters
make lint
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes using [conventional commits](https://www.conventionalcommits.org/) (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

### Commit Message Format

This project uses [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `style:` Code style changes
- `refactor:` Code refactoring
- `perf:` Performance improvements
- `test:` Test additions or corrections
- `build:` Build system changes
- `ci:` CI/CD changes
- `chore:` Other changes

## 🚀 CI/CD

This project uses:
- **GitHub Actions** for CI/CD
- **semantic-release** for automated versioning
- **goreleaser** for building and releasing binaries
- **commitlint** for commit message validation

Releases are automatically created when commits are pushed to the `main` branch.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI management
- Uses [Viper](https://github.com/spf13/viper) for configuration
- Logging powered by [zerolog](https://github.com/rs/zerolog)

## 📮 Support

- Report bugs via [GitHub Issues](https://github.com/mjmorales/mac-daemon-control/issues)
- Request features via [GitHub Issues](https://github.com/mjmorales/mac-daemon-control/issues)
- Ask questions in [GitHub Discussions](https://github.com/mjmorales/mac-daemon-control/discussions)