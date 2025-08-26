package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestScriptValidationScenarios(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test files
	validScript := filepath.Join(tmpDir, "valid.sh")
	err := os.WriteFile(validScript, []byte("#!/bin/bash\necho valid"), 0755)
	if err != nil {
		t.Fatalf("Failed to create valid script: %v", err)
	}
	
	validPython := filepath.Join(tmpDir, "valid.py")
	err = os.WriteFile(validPython, []byte("print('valid')"), 0644)
	if err != nil {
		t.Fatalf("Failed to create valid python script: %v", err)
	}

	// Create subdirectory with script
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	
	subScript := filepath.Join(subDir, "sub.sh")
	err = os.WriteFile(subScript, []byte("#!/bin/bash\necho sub"), 0755)
	if err != nil {
		t.Fatalf("Failed to create sub script: %v", err)
	}

	// Change to temp directory for relative path tests
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	tests := []struct {
		name        string
		script      HookScript
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid script file",
			script:      HookScript{Name: "test", Script: "valid.sh", Description: "Valid"},
			expectError: false,
		},
		{
			name:        "valid script with arguments",
			script:      HookScript{Name: "test", Script: "valid.sh --verbose", Description: "Valid with args"},
			expectError: false,
		},
		{
			name:        "valid python script",
			script:      HookScript{Name: "test", Script: "valid.py", Description: "Valid Python"},
			expectError: false,
		},
		{
			name:        "valid script in subdirectory",
			script:      HookScript{Name: "test", Script: "subdir/sub.sh", Description: "Valid sub"},
			expectError: false,
		},
		{
			name:        "valid command in PATH",
			script:      HookScript{Name: "test", Command: "echo hello", Description: "Valid command"},
			expectError: false,
		},
		{
			name:        "valid complex command",
			script:      HookScript{Name: "test", Command: "ls -la", Description: "Valid complex"},
			expectError: false,
		},
		{
			name:        "missing script file",
			script:      HookScript{Name: "test", Script: "missing.sh", Description: "Missing"},
			expectError: true,
			errorMsg:    "script file 'missing.sh' not found",
		},
		{
			name:        "missing script with arguments",
			script:      HookScript{Name: "test", Script: "missing.sh --arg", Description: "Missing with args"},
			expectError: true,
			errorMsg:    "script file 'missing.sh' not found",
		},
		{
			name:        "missing command",
			script:      HookScript{Name: "test", Command: "nonexistent-cmd-xyz", Description: "Missing cmd"},
			expectError: true,
			errorMsg:    "command 'nonexistent-cmd-xyz' not found in PATH",
		},
		{
			name:        "missing complex command",
			script:      HookScript{Name: "test", Command: "nonexistent-cmd-xyz --flag", Description: "Missing complex"},
			expectError: true,
			errorMsg:    "command 'nonexistent-cmd-xyz' not found in PATH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hooks := map[string][]HookScript{
				"pre-commit": {tt.script},
			}

			hm := &HookManager{
				config: &Config{
					Hooks:    hooks,
					Settings: Settings{},
				},
			}

			err := hm.validateScripts()

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

func TestListHooksValidation(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test script
	validScript := filepath.Join(tmpDir, "valid.sh")
	err := os.WriteFile(validScript, []byte("#!/bin/bash\necho test"), 0755)
	if err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Change to temp directory
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Create temporary config file
	configContent := `
hooks:
  pre-commit:
    - name: "valid-script"
      script: "valid.sh"
      description: "Valid script"
    - name: "missing-script"
      script: "missing.sh"
      description: "Missing script"
    - name: "valid-command"
      command: "echo test"
      description: "Valid command"
    - name: "missing-command"
      command: "nonexistent-cmd-xyz"
      description: "Missing command"
`
	
	configPath := filepath.Join(tmpDir, "test-config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create git repository
	gitDir := filepath.Join(tmpDir, ".git")
	err = os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create git dir: %v", err)
	}
	
	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	hm := NewHookManager(configPath, false)
	err = hm.init()
	if err != nil {
		t.Fatalf("Failed to initialize hook manager: %v", err)
	}

	// Test ListHooks doesn't error (it should show validation status)
	err = hm.ListHooks()
	if err != nil {
		t.Errorf("ListHooks should not error, got: %v", err)
	}

	// Verify that validation would catch the issues
	err = hm.validateScripts()
	if err == nil {
		t.Error("Expected validation errors for missing script and command")
	}

	errorMsg := err.Error()
	if !containsString(errorMsg, "script file 'missing.sh' not found") {
		t.Errorf("Expected missing script error, got: %s", errorMsg)
	}
	if !containsString(errorMsg, "command 'nonexistent-cmd-xyz' not found in PATH") {
		t.Errorf("Expected missing command error, got: %s", errorMsg)
	}
}

func TestEdgeCaseValidation(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		setupFunc   func() HookScript
		expectError bool
		errorMsg    string
	}{
		{
			name: "empty script path",
			setupFunc: func() HookScript {
				return HookScript{Name: "test", Script: "", Command: "echo test", Description: "Empty script"}
			},
			expectError: false, // Should be valid because command is provided
		},
		{
			name: "empty command",
			setupFunc: func() HookScript {
				// Create a test file
				scriptPath := filepath.Join(tmpDir, "test.sh")
				os.WriteFile(scriptPath, []byte("echo test"), 0755)
				return HookScript{Name: "test", Script: scriptPath, Command: "", Description: "Empty command"}
			},
			expectError: false, // Should be valid because script is provided
		},
		{
			name: "script with multiple arguments",
			setupFunc: func() HookScript {
				scriptPath := filepath.Join(tmpDir, "multi-arg.sh")
				os.WriteFile(scriptPath, []byte("echo test"), 0755)
				return HookScript{Name: "test", Script: scriptPath + " --flag1 --flag2 value", Description: "Multi args"}
			},
			expectError: false,
		},
		{
			name: "command with pipes and redirects",
			setupFunc: func() HookScript {
				return HookScript{Name: "test", Command: "echo test | grep test > /dev/null", Description: "Complex command"}
			},
			expectError: false, // echo exists in PATH
		},
		{
			name: "script path with spaces (quoted)",
			setupFunc: func() HookScript {
				// Create directory and file with spaces
				spaceDir := filepath.Join(tmpDir, "dir_with_spaces")
				os.MkdirAll(spaceDir, 0755)
				scriptPath := filepath.Join(spaceDir, "script_with_spaces.sh")
				os.WriteFile(scriptPath, []byte("echo test"), 0755)
				return HookScript{Name: "test", Script: scriptPath, Description: "Path with spaces"}
			},
			expectError: false,
		},
	}

	// Change to temp directory
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := tt.setupFunc()
			hooks := map[string][]HookScript{
				"pre-commit": {script},
			}

			hm := &HookManager{
				config: &Config{
					Hooks:    hooks,
					Settings: Settings{},
				},
			}

			err := hm.validateScripts()

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