# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.3.0] - 2024-08-26

### Changed
- **BREAKING**: Replaced single `script` property with separate `script` and `command` properties
- **script**: Now only for executable files (with optional arguments)
- **command**: New property for direct commands to run

### Added
- Configuration validation: Each hook entry must have either `script` OR `command` (not both, not neither)
- Clear distinction between file scripts and direct commands in `--list` output
- Better error messages specify whether script file or command is missing

### Improved
- Eliminated complex auto-detection logic in favor of explicit configuration
- Cleaner, more predictable behavior
- Better validation and error reporting

### Migration
**Old format (v1.2.x):**
```yaml
hooks:
  pre-commit:
    - script: "hooks/format.sh"           # File script
    - script: "go test ./..."            # Direct command (auto-detected)
```

**New format (v1.3.0):**
```yaml
hooks:
  pre-commit:
    - script: "hooks/format.sh"          # Explicit file script
    - command: "go test ./..."           # Explicit command
```

## [1.2.0] - 2024-08-26

### Added
- **Command support**: Scripts can now be direct commands like `go test ./...`, `npm run lint`, `fvm dart format --set-exit-if-changed lib`
- **Smart validation**: Automatically distinguishes between file scripts and commands
- **Command validation**: Validates that commands are available in PATH before installation
- **Script arguments support**: File scripts can include arguments like `hooks/lint.sh --verbose`
- Enhanced `--list` output shows `[file]` vs `[cmd]` indicators

### Improved
- Better error messages distinguish between missing files and missing commands
- Validation handles complex command-line arguments properly
- Support for script files with arguments (validates just the file part)

### Examples
```yaml
hooks:
  pre-commit:
    - script: "hooks/format.sh"                    # File script
    - script: "hooks/lint.sh --verbose"            # File with args  
    - script: "go test ./..."                      # Direct command
    - script: "fvm dart format --set-exit-if-changed lib" # Complex command
```

## [1.1.0] - 2024-08-26

### Changed
- **BREAKING**: Removed `hooks_directory` configuration option
- Script paths are now specified directly in the `script` property as relative paths
- More flexible script organization - scripts can be in any directory structure

### Improved  
- Simplified configuration format
- Better support for different project layouts
- Scripts can be anywhere: `hooks/script.sh`, `tools/check.py`, `./validate.sh`, or even direct commands

### Migration
- Old format: `hooks_directory: "hooks"` + `script: "format.sh"` 
- New format: `script: "hooks/format.sh"`

## [1.0.1] - 2024-08-26

### Added
- Script validation: Hooky now validates that all referenced scripts exist before installing hooks
- Enhanced `--list` command shows ✅ for existing scripts and ❌ for missing scripts
- Clear error messages when referenced scripts are missing
- Warning display when using `--list` if there are missing scripts

### Improved
- Better error reporting with specific script names and associated hooks
- Installation fails early if any scripts are missing, preventing partial installations

## [1.0.0] - 2024-08-26

### Added
- Initial release of Hooky git hooks manager
- Cross-platform support (Linux, macOS, Windows, FreeBSD)
- YAML-based configuration system
- Support for all standard git hooks
- Multiple scripts per hook with sequential execution
- Automatic backup of existing hooks
- Verbose logging and clear error messages
- Example hook scripts for common tasks:
  - Code formatting checks (Go, JavaScript, Python)
  - Linting (multiple languages)
  - Testing (unit and integration)
  - Commit message validation
  - Pre-rebase safety checks
  - Dependency checking
  - Post-merge cleanup
- Build system with cross-platform compilation
- GitHub Actions CI/CD pipeline
- Comprehensive documentation and examples
- MIT license
- Contributing guidelines

### Features
- `--install` - Install hooks from configuration
- `--uninstall` - Remove hooky-managed hooks
- `--list` - Display configured hooks
- `--config` - Use custom configuration file
- `--verbose` - Enable detailed output
- `--version` - Show version information

### Configuration Options
- `hooks_directory` - Specify where hook scripts are stored
- `settings.auto_executable` - Automatically make scripts executable
- `settings.backup_existing` - Backup existing hooks before installation
- `settings.backup_directory` - Custom backup location
- `settings.verbose` - Enable verbose output by default

### Supported Platforms
- Linux (amd64, arm64)
- macOS (amd64, arm64) 
- Windows (amd64, arm64)
- FreeBSD (amd64)

[Unreleased]: https://github.com/yourusername/hooky/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/yourusername/hooky/releases/tag/v1.0.0