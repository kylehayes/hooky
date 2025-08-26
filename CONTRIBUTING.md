# Contributing to Hooky

Thank you for your interest in contributing to Hooky! We welcome contributions from everyone, whether you're fixing bugs, adding features, improving documentation, or suggesting enhancements.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, but recommended)

### Development Setup

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/yourusername/hooky.git
   cd hooky
   ```

2. **Install dependencies:**
   ```bash
   make deps
   # or manually:
   go mod download && go mod tidy
   ```

3. **Build and test:**
   ```bash
   make dev  # Runs format, lint, test, and build
   ```

4. **Install locally for testing:**
   ```bash
   make install
   ```

## 🔄 Development Workflow

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-number
   ```

2. **Make your changes:**
   - Write code following the existing style
   - Add tests for new functionality
   - Update documentation if needed

3. **Test your changes:**
   ```bash
   make test      # Run tests
   make lint      # Run linter
   make build     # Build for current platform
   make build-all # Build for all platforms (optional)
   ```

4. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat: add amazing new feature"
   # or
   git commit -m "fix: resolve issue with hook installation"
   ```

5. **Push and create a Pull Request:**
   ```bash
   git push origin feature/your-feature-name
   ```

## 📝 Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or modifying tests
- `chore:` - Build process or auxiliary tool changes
- `perf:` - Performance improvements
- `ci:` - Continuous integration changes

Examples:
```bash
feat: add support for custom hook directories
fix: resolve Windows path handling issue
docs: update installation instructions
test: add unit tests for config parsing
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Writing Tests

- Write unit tests for all new functionality
- Place tests in `*_test.go` files alongside the code they test
- Use table-driven tests where appropriate
- Mock external dependencies

Example test structure:
```go
func TestHookManager_InstallHooks(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
    }{
        {
            name:    "valid config",
            config:  &Config{...},
            wantErr: false,
        },
        // more test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## 📖 Documentation

### Code Documentation

- Write clear, concise comments for public functions and types
- Use godoc-style comments
- Document complex logic and algorithms

### README and Documentation

- Update the README.md if you add new features
- Add examples for new functionality
- Update configuration documentation for new options

## 🏗️ Project Structure

```
hooky/
├── .github/           # GitHub Actions workflows
├── hooks/             # Example hook scripts
├── dist/              # Build artifacts (generated)
├── main.go            # Main entry point
├── config.go          # Configuration parsing
├── manager.go         # Core hook management logic
├── go.mod             # Go module definition
├── go.sum             # Go module checksums
├── hooky.yaml         # Example configuration
├── Makefile           # Build automation
├── build.sh           # Cross-platform build script
├── README.md          # Project documentation
├── CONTRIBUTING.md    # This file
└── LICENSE            # MIT license
```

## 🎯 Areas for Contribution

### High Priority
- [ ] Windows-specific improvements and testing
- [ ] Additional hook script examples for different languages
- [ ] Performance optimizations
- [ ] Better error handling and user feedback

### Medium Priority
- [ ] Hook template generation commands
- [ ] Integration with popular CI/CD systems
- [ ] Configuration validation improvements
- [ ] Parallel hook execution support

### Low Priority
- [ ] GUI or TUI interface
- [ ] Plugin system for custom hook types
- [ ] Advanced configuration options
- [ ] Metrics and reporting features

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Environment information:**
   - Operating system and version
   - Go version
   - Hooky version

2. **Steps to reproduce:**
   - Clear, numbered steps
   - Expected vs actual behavior
   - Any error messages

3. **Configuration:**
   - Your `hooky.yaml` configuration
   - Hook scripts (if relevant)
   - Git repository state

4. **Minimal reproduction:**
   - If possible, create a minimal example that reproduces the issue

## 💡 Feature Requests

For feature requests, please:

1. Check if the feature already exists or is planned
2. Describe the use case and why it would be valuable
3. Provide examples of how it should work
4. Consider backward compatibility

## 🔍 Code Review Process

1. **All changes require review:** Even small changes benefit from a second pair of eyes
2. **Be constructive:** Provide helpful feedback and suggestions
3. **Be responsive:** Address feedback promptly and professionally
4. **Test thoroughly:** Ensure all tests pass and new code is well-tested

### What Reviewers Look For

- Code quality and style consistency
- Test coverage for new functionality
- Documentation updates
- Backward compatibility
- Security considerations
- Performance impact

## 🎨 Code Style

- Follow standard Go conventions
- Use `gofmt` to format code
- Run `golangci-lint` or `go vet` to catch issues
- Write clear, descriptive variable and function names
- Keep functions focused and reasonably sized

### Example Style

```go
// InstallHooks installs all configured git hooks.
// It creates backup copies of existing hooks if backup_existing is enabled.
func (hm *HookManager) InstallHooks() error {
    if err := hm.init(); err != nil {
        return fmt.Errorf("failed to initialize hook manager: %w", err)
    }

    for hookName, scripts := range hm.config.Hooks {
        if err := hm.installHook(hookName, scripts); err != nil {
            return fmt.Errorf("failed to install hook %s: %w", hookName, err)
        }
    }

    return nil
}
```

## 🏷️ Release Process

Releases are automated through GitHub Actions when tags are pushed:

1. **Create a release branch:**
   ```bash
   git checkout -b release/v1.1.0
   ```

2. **Update version numbers:**
   - Update version in `main.go`
   - Update version in documentation
   - Update CHANGELOG.md

3. **Create a pull request and merge**

4. **Tag the release:**
   ```bash
   git tag v1.1.0
   git push origin v1.1.0
   ```

5. **GitHub Actions will automatically:**
   - Build binaries for all platforms
   - Create a GitHub release
   - Attach build artifacts

## 🤝 Community Guidelines

- Be respectful and inclusive
- Help others learn and grow
- Follow the [Go Community Code of Conduct](https://golang.org/conduct)
- When in doubt, ask questions

## 📞 Getting Help

- **GitHub Issues:** For bugs and feature requests
- **GitHub Discussions:** For questions and general discussion
- **Code Review:** For feedback on your contributions

## 🙏 Recognition

Contributors are recognized in:
- GitHub's contribution graphs
- Release notes for significant contributions
- The project's acknowledgments

Thank you for contributing to Hooky! 🪝