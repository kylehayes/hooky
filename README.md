# 🪝 Hooky

A simple, cross-platform git hooks manager that works without package managers - just like Husky and pre-commit, but even simpler.

[![CI](https://github.com/yourusername/hooky/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/hooky/actions/workflows/ci.yml)
[![Release](https://github.com/yourusername/hooky/actions/workflows/release.yml/badge.svg)](https://github.com/yourusername/hooky/actions/workflows/release.yml)

## ✨ Features

- 🚀 **Zero dependencies** - Single binary, no package manager required
- 🌍 **Cross-platform** - Works on Linux, macOS, Windows, and FreeBSD
- ⚙️ **YAML configuration** - Simple and readable configuration format
- 📁 **Conventional directory structure** - Organize hooks in a `hooks/` directory
- 🔄 **Multiple scripts per hook** - Run multiple scripts in sequence for each git hook
- 🛡️ **Safe installation** - Automatically backs up existing hooks
- 🎯 **Selective activation** - Enable only the hooks you need
- ✅ **Script validation** - Ensures all referenced scripts exist before installation
- 📝 **Rich logging** - Clear output showing what's running and why

## 🚀 Quick Start

### Installation

1. **Download the appropriate binary for your platform:**
   ```bash
   # Linux (amd64)
   curl -L -o hooky.tar.gz https://github.com/yourusername/hooky/releases/latest/download/hooky-1.0.0-linux-amd64.tar.gz
   tar -xzf hooky.tar.gz && cd hooky-*
   
   # macOS (arm64 for Apple Silicon)
   curl -L -o hooky.tar.gz https://github.com/yourusername/hooky/releases/latest/download/hooky-1.0.0-darwin-arm64.tar.gz
   tar -xzf hooky.tar.gz && cd hooky-*
   
   # Windows (amd64)
   # Download hooky-1.0.0-windows-amd64.zip from releases page
   ```

2. **Move the binary to your PATH:**
   ```bash
   # Linux/macOS
   sudo mv hooky /usr/local/bin/
   
   # Or add to your project (no sudo required)
   mv hooky ./hooky
   ```

3. **Initialize in your git repository:**
   ```bash
   # Copy the configuration and hooks directory to your project
   cp hooky.yaml your-project/
   cp -r hooks your-project/
   cd your-project
   
   # Install the hooks
   ./hooky --install
   ```

## 📖 Usage

### Basic Commands

```bash
# Install hooks from configuration
hooky --install

# Uninstall hooks (removes only hooky-generated hooks)
hooky --uninstall

# List configured hooks (shows ✅ for existing scripts/commands, ❌ for missing)
# Also shows [file] vs [cmd] to indicate script type
hooky --list

# Use custom configuration file
hooky --config custom-hooks.yaml --install

# Enable verbose output
hooky --verbose --install

# Show version
hooky --version
```

### Configuration

Create a `hooky.yaml` file in your repository root:

```yaml
# Git hooks configuration
# Use "script" for executable files or "command" for direct commands
hooks:
  pre-commit:
    - name: "format-check"
      script: "hooks/format.sh"
      description: "Check code formatting with script file"
    - name: "lint"
      script: "hooks/lint.sh --verbose" 
      description: "Run linting checks (script with arguments)"
    - name: "go-test"
      command: "go test ./..."
      description: "Run Go tests directly"

  pre-push:
    - name: "dart-format"
      command: "fvm dart format --set-exit-if-changed lib packages test"
      description: "Check Dart code formatting"
    - name: "integration-tests"
      script: "tools/integration-tests.py"
      description: "Integration tests (Python script)"

  commit-msg:
    - name: "conventional-commits"
      script: "check-commit-msg.sh"
      description: "Script in project root"

# Global settings
settings:
  auto_executable: true        # Make hook scripts executable automatically
  backup_existing: true        # Backup existing hooks before installing
  backup_directory: ".hooky-backup"  # Where to store backups
  verbose: false              # Show detailed output
```

**Key features:**
- **Separate script and command properties**: Clear distinction between executable files and direct commands
- **Script files**: Use `script` for executable files - `hooks/script.sh`, `tools/check.py --verbose`
- **Direct commands**: Use `command` for direct commands - `go test ./...`, `fvm dart format --set-exit-if-changed lib`
- **Multiple languages**: Use any executable script (shell, Python, Node.js, Dart, etc.)
- **Custom organization**: Organize scripts however you want in your project

### Directory Structure

```
your-project/
├── .git/
├── hooky.yaml          # Configuration file
├── hooks/              # Common location for hook scripts
│   ├── format.sh
│   ├── lint.sh
│   └── ...
├── scripts/            # Alternative location
│   ├── test.py
│   └── build.sh
├── tools/              # Another alternative
│   └── integration-test.js
├── check-root.sh       # Scripts can be anywhere
└── ... (your project files)
```

With this structure, your configuration might look like:
```yaml
hooks:
  pre-commit:
    - script: "hooks/format.sh"                    # File in hooks/ directory
    - script: "scripts/test.py --verbose"          # Python script with args
    - command: "go test ./..."                     # Direct Go command
    - command: "npm run lint"                      # npm script
    - command: "fvm dart format --set-exit-if-changed lib" # Dart command
```

**Property Types:**
- **script**: Executable files (with optional arguments)
  - `hooks/format.sh`
  - `scripts/test.py --coverage`
  - `tools/lint.js --fix`
- **command**: Direct commands to run
  - `go test ./...`
  - `npm run build`
  - `fvm dart format --set-exit-if-changed lib packages test`
  - `make test`

### Script Validation

Hooky validates both script files and commands before installing:

```bash
# Use --list to see script/command status
$ hooky --list
Hook: pre-commit
  1. format-check (hooks/format.sh) [script] ✅
     Check code formatting
  2. go-test (go test ./...) [command] ✅
     Run Go tests
  3. dart-format (fvm dart format --set-exit-if-changed lib) [command] ✅
     Format Dart code  
  4. missing-script (scripts/missing.py) [script] ❌ MISSING
     Script that doesn't exist
  5. missing-command (nonexistent-cmd --flag) [command] ❌ MISSING
     Command not in PATH

⚠️  Missing scripts that would prevent installation:
  script file 'scripts/missing.py' not found (from: scripts/missing.py, hook: pre-commit)
  command 'nonexistent-cmd' not found in PATH (from: nonexistent-cmd --flag, hook: pre-commit)

# Installation fails if anything is missing
$ hooky --install
Error installing hooks: missing scripts/commands:
  script file 'scripts/missing.py' not found (from: scripts/missing.py, hook: pre-commit)
  command 'nonexistent-cmd' not found in PATH (from: nonexistent-cmd --flag, hook: pre-commit)
```

**Configuration Validation:**
- Each hook entry must have either `script` OR `command` (not both, not neither)
- **script**: Validates the file exists (ignores arguments after first space)
- **command**: Validates the command is available in PATH

**Configuration Errors:**
```yaml
# ❌ Invalid - has both script and command
- name: "invalid"
  script: "hooks/test.sh"
  command: "go test"

# ❌ Invalid - has neither script nor command  
- name: "invalid"
  description: "Missing script or command"

# ✅ Valid - has script only
- name: "valid"
  script: "hooks/test.sh"
  
# ✅ Valid - has command only
- name: "valid"
  command: "go test ./..."
```

## 🔧 Hook Scripts

Hook scripts should be executable shell scripts. Here's a simple example:

```bash
#!/bin/bash
# hooks/format.sh

echo "🎨 Checking code formatting..."

# Check Go files
if find . -name "*.go" | grep -q .; then
    unformatted=$(gofmt -l . | grep -v vendor || true)
    if [ -n "$unformatted" ]; then
        echo "❌ The following files need formatting:"
        echo "$unformatted"
        echo "Run: gofmt -w ."
        exit 1
    fi
fi

echo "✅ Code formatting looks good!"
exit 0
```

### Supported Git Hooks

Hooky supports all standard git hooks:

- `pre-commit` - Before each commit
- `prepare-commit-msg` - Prepare commit message
- `commit-msg` - Validate commit message
- `post-commit` - After each commit
- `pre-rebase` - Before rebase
- `post-checkout` - After checkout
- `post-merge` - After merge
- `pre-push` - Before push
- `pre-receive` - Before receive (server-side)
- `update` - Update reference (server-side)
- `post-receive` - After receive (server-side)
- `post-update` - After update (server-side)
- `pre-auto-gc` - Before auto garbage collection
- `post-rewrite` - After rewrite (rebase/amend)
- `push-to-checkout` - Push to checkout
- `applypatch-msg` - Apply patch message
- `pre-applypatch` - Before apply patch
- `post-applypatch` - After apply patch

## 🎯 Examples

### Node.js Project

```yaml
hooks:
  pre-commit:
    - name: "prettier"
      command: "prettier --check ."
      description: "Format code with Prettier"
    - name: "eslint" 
      command: "eslint . --ext .js,.ts,.jsx,.tsx"
      description: "Lint JavaScript/TypeScript"
    - name: "test"
      command: "npm run test"
      description: "Run unit tests using npm"
    - name: "custom-check"
      script: "scripts/custom-validation.sh"
      description: "Custom validation script"

  pre-push:
    - name: "build"
      command: "npm run build"
      description: "Ensure project builds"
```

### Python Project

```yaml
hooks:
  pre-commit:
    - name: "black"
      command: "black --check ."
      description: "Format Python code with Black"
    - name: "flake8"
      command: "flake8 ." 
      description: "Lint Python code"
    - name: "pytest"
      command: "pytest"
      description: "Run Python tests"
    - name: "custom-check"
      script: "dev-tools/custom-python-check.py --strict"
      description: "Custom Python validation script"

  commit-msg:
    - name: "conventional"
      script: "tools/check-commit-format.sh"
      description: "Validate conventional commit format"
```

### Multi-language Project

```yaml
hooks:
  pre-commit:
    - name: "format-go"
      command: "gofmt -l ."
      description: "Format Go code"
    - name: "format-js"
      command: "prettier --check frontend/"
      description: "Format JavaScript code"
    - name: "lint-all"
      script: "tools/lint-all-languages.py"
      description: "Lint all supported languages"
    - name: "test-backend"
      command: "go test ./..."
      description: "Run Go tests"
    - name: "test-frontend" 
      command: "npm run test --prefix frontend"
      description: "Run JavaScript tests"

  pre-push:
    - name: "integration-tests"
      script: "ci/integration-test.sh"
      description: "Run full integration test suite"
    - name: "security-scan"
      script: "security/scan.py --all"
      description: "Run security vulnerability scan"
    - name: "docker-build"
      command: "docker build -t myapp ."
      description: "Test Docker build"

### Dart/Flutter Project

```yaml
hooks:
  pre-commit:
    - name: "dart-format"
      command: "fvm dart format --set-exit-if-changed lib packages test"
      description: "Check Dart code formatting"
    - name: "dart-analyze"
      command: "fvm dart analyze"
      description: "Run Dart analyzer"
    - name: "flutter-test"
      command: "fvm flutter test"
      description: "Run Flutter tests"
    - name: "custom-checks"
      script: "scripts/flutter-custom-checks.sh"
      description: "Custom Flutter validation"

settings:
  verbose: true
  backup_existing: true
```

## 🔨 Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/hooky.git
cd hooky

# Install dependencies
make deps

# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Development workflow (format, lint, test, build)
make dev
```

### Development Setup

```bash
# Install development dependencies
go mod download

# Run tests
go test ./...

# Install locally for testing
make install

# Test with a sample project
cd /path/to/your/project
hooky --install
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Contribution Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Ensure all tests pass (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Workflow

- `make dev` - Run the full development workflow
- `make test` - Run tests
- `make lint` - Run linter
- `make fmt` - Format code

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [Husky](https://typicode.github.io/husky/) and [pre-commit](https://pre-commit.com/)
- Built with ❤️ in Go for maximum portability and performance

## 📞 Support

- 🐛 [Report bugs](https://github.com/yourusername/hooky/issues)
- 💡 [Request features](https://github.com/yourusername/hooky/issues)
- 📖 [Documentation](https://github.com/yourusername/hooky/wiki)
- 💬 [Discussions](https://github.com/yourusername/hooky/discussions)

---

**Made with 🪝 by the Hooky team**