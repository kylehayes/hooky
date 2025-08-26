package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configYAML  string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config with scripts",
			configYAML: `
hooks:
  pre-commit:
    - name: "test"
      script: "test.sh"
      description: "Test script"
settings:
  verbose: true
`,
			expectError: false,
		},
		{
			name: "valid config with commands",
			configYAML: `
hooks:
  pre-commit:
    - name: "test"
      command: "go test ./..."
      description: "Test command"
`,
			expectError: false,
		},
		{
			name: "valid config with mixed script and command",
			configYAML: `
hooks:
  pre-commit:
    - name: "script-test"
      script: "test.sh"
      description: "Script test"
    - name: "command-test"
      command: "go test ./..."
      description: "Command test"
`,
			expectError: false,
		},
		{
			name: "invalid config - both script and command",
			configYAML: `
hooks:
  pre-commit:
    - name: "test"
      script: "test.sh"
      command: "go test ./..."
      description: "Invalid"
`,
			expectError: true,
			errorMsg:    "cannot specify both 'script' and 'command'",
		},
		{
			name: "invalid config - neither script nor command",
			configYAML: `
hooks:
  pre-commit:
    - name: "test"
      description: "Invalid"
`,
			expectError: true,
			errorMsg:    "must specify either 'script' or 'command'",
		},
		{
			name: "invalid YAML",
			configYAML: `
hooks:
  pre-commit:
    - name: "test"
      script: "test.sh
      description: "Unclosed quote
`,
			expectError: true,
			errorMsg:    "failed to parse config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "test-config.yaml")
			
			err := os.WriteFile(configPath, []byte(tt.configYAML), 0644)
			if err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			// Test LoadConfig
			config, err := LoadConfig(configPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if config == nil {
				t.Errorf("Config is nil")
				return
			}

			// Verify defaults are set
			if config.Settings.BackupDirectory != ".hooky-backup" {
				t.Errorf("Expected default backup directory '.hooky-backup', got '%s'", config.Settings.BackupDirectory)
			}
		})
	}
}

func TestLoadConfigNonexistentFile(t *testing.T) {
	_, err := LoadConfig("nonexistent-file.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
	if !containsString(err.Error(), "failed to read config file") {
		t.Errorf("Expected 'failed to read config file' error, got: %s", err.Error())
	}
}

func TestValidateHookScripts(t *testing.T) {
	tests := []struct {
		name        string
		hooks       map[string][]HookScript
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid hooks with scripts",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "test1", Script: "test.sh", Description: "Test"},
				},
			},
			expectError: false,
		},
		{
			name: "valid hooks with commands",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "test1", Command: "go test", Description: "Test"},
				},
			},
			expectError: false,
		},
		{
			name: "valid mixed hooks",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "script", Script: "test.sh", Description: "Script"},
					{Name: "command", Command: "go test", Description: "Command"},
				},
			},
			expectError: false,
		},
		{
			name: "invalid - both script and command",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "invalid", Script: "test.sh", Command: "go test", Description: "Invalid"},
				},
			},
			expectError: true,
			errorMsg:    "cannot specify both 'script' and 'command'",
		},
		{
			name: "invalid - neither script nor command",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "invalid", Description: "Invalid"},
				},
			},
			expectError: true,
			errorMsg:    "must specify either 'script' or 'command'",
		},
		{
			name: "multiple hooks with mixed validity",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "valid", Script: "test.sh", Description: "Valid"},
					{Name: "invalid", Description: "Invalid"},
				},
			},
			expectError: true,
			errorMsg:    "must specify either 'script' or 'command'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHookScripts(tt.hooks)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetSupportedHooks(t *testing.T) {
	hooks := GetSupportedHooks()
	
	if len(hooks) == 0 {
		t.Error("Expected non-empty list of supported hooks")
	}

	// Check for some known hooks
	expectedHooks := []string{"pre-commit", "post-commit", "pre-push", "commit-msg"}
	for _, expected := range expectedHooks {
		found := false
		for _, hook := range hooks {
			if hook == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find hook '%s' in supported hooks", expected)
		}
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    s[:len(substr)] == substr || 
		    s[len(s)-len(substr):] == substr ||
		    func() bool {
		    	for i := 1; i <= len(s)-len(substr); i++ {
		    		if s[i:i+len(substr)] == substr {
		    			return true
		    		}
		    	}
		    	return false
		    }())
}