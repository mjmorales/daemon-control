## [1.0.3](https://github.com/mjmorales/daemon-control/compare/v1.0.2...v1.0.3) (2025-07-25)

### 🐛 Bug Fixes

* remove package.json to prevent dirty git state during release ([697b63b](https://github.com/mjmorales/daemon-control/commit/697b63b91dfe8ab0efb43da49964d656e087292e))

### ♻️ Code Refactoring

* rename project from mac-daemon-control to daemon-control ([c18243c](https://github.com/mjmorales/daemon-control/commit/c18243c9377c19682b51d131bc675c1583941917))

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of daemon-control
- Generic daemon management for macOS LaunchAgents
- YAML-based daemon configuration
- Automatic plist generation from YAML
- Comprehensive daemon control commands (install, uninstall, start, stop, restart, status)
- Log management with logs and tail commands
- Core configuration management system
- Editor integration for quick config editing
- CI/CD with GitHub Actions, semantic-release, and goreleaser
- Homebrew tap support
