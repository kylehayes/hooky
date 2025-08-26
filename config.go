package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type HookScript struct {
	Name        string `yaml:"name"`
	Script      string `yaml:"script,omitempty"`
	Command     string `yaml:"command,omitempty"`
	Description string `yaml:"description"`
}

type Settings struct {
	AutoExecutable  bool   `yaml:"auto_executable"`
	BackupExisting  bool   `yaml:"backup_existing"`
	BackupDirectory string `yaml:"backup_directory"`
	Verbose         bool   `yaml:"verbose"`
}

type Config struct {
	Hooks    map[string][]HookScript `yaml:"hooks"`
	Settings Settings              `yaml:"settings"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	
	// Set defaults
	config.Settings = Settings{
		AutoExecutable:  true,
		BackupExisting:  true,
		BackupDirectory: ".hooky-backup",
		Verbose:         false,
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate hook script configuration
	if err := validateHookScripts(config.Hooks); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateHookScripts(hooks map[string][]HookScript) error {
	for hookName, scripts := range hooks {
		for i, script := range scripts {
			hasScript := script.Script != ""
			hasCommand := script.Command != ""
			
			if !hasScript && !hasCommand {
				return fmt.Errorf("hook %s[%d] (%s): must specify either 'script' or 'command'", hookName, i, script.Name)
			}
			
			if hasScript && hasCommand {
				return fmt.Errorf("hook %s[%d] (%s): cannot specify both 'script' and 'command', use only one", hookName, i, script.Name)
			}
		}
	}
	return nil
}

// GetSupportedHooks returns all git hooks that can be managed
func GetSupportedHooks() []string {
	return []string{
		"applypatch-msg",
		"pre-applypatch",
		"post-applypatch",
		"pre-commit",
		"prepare-commit-msg",
		"commit-msg",
		"post-commit",
		"pre-rebase",
		"post-checkout",
		"post-merge",
		"pre-receive",
		"update",
		"post-receive",
		"post-update",
		"pre-auto-gc",
		"post-rewrite",
		"pre-push",
		"push-to-checkout",
	}
}