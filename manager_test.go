package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewHookManager(t *testing.T) {
	hm := NewHookManager("test-config.yaml", true)
	
	if hm == nil {
		t.Error("Expected non-nil HookManager")
	}
	
	if hm.configPath != "test-config.yaml" {
		t.Errorf("Expected configPath 'test-config.yaml', got '%s'", hm.configPath)
	}
	
	if !hm.verbose {
		t.Error("Expected verbose to be true")
	}
}

func TestValidateScripts(t *testing.T) {
	// Create temporary directory with test files
	tmpDir := t.TempDir()
	
	// Create test script file
	scriptPath := filepath.Join(tmpDir, "test-script.sh")
	err := os.WriteFile(scriptPath, []byte("#!/bin/bash\necho test"), 0755)
	if err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	tests := []struct {
		name        string
		hooks       map[string][]HookScript
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid script file exists",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "test", Script: scriptPath, Description: "Test"},
				},
			},
			expectError: false,
		},
		{
			name: "valid script with arguments",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "test", Script: scriptPath + " --verbose", Description: "Test with args"},
				},
			},
			expectError: false,
		},
		{
			name: "valid command exists in PATH",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "test", Command: "echo hello", Description: "Test command"},
				},
			},
			expectError: false,
		},
		{
			name: "invalid script file missing",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "test", Script: filepath.Join(tmpDir, "missing.sh"), Description: "Missing script"},
				},
			},
			expectError: true,
			errorMsg:    "script file",
		},
		{
			name: "invalid command not in PATH",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "test", Command: "nonexistent-command-12345", Description: "Missing command"},
				},
			},
			expectError: true,
			errorMsg:    "command 'nonexistent-command-12345' not found in PATH",
		},
		{
			name: "mixed valid and invalid",
			hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "valid", Script: scriptPath, Description: "Valid"},
					{Name: "invalid", Script: filepath.Join(tmpDir, "missing.sh"), Description: "Invalid"},
				},
			},
			expectError: true,
			errorMsg:    "script file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HookManager with test configuration
			hm := &HookManager{
				config: &Config{
					Hooks:    tt.hooks,
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

func TestGenerateHookScript(t *testing.T) {
	tmpDir := t.TempDir()
	
	tests := []struct {
		name        string
		hookName    string
		scripts     []HookScript
		expectError bool
		contains    []string
	}{
		{
			name:     "single script",
			hookName: "pre-commit",
			scripts: []HookScript{
				{Name: "test", Script: "test.sh", Description: "Test script"},
			},
			expectError: false,
			contains:    []string{"#!/bin/sh", "test.sh", "Running: test", "Test script"},
		},
		{
			name:     "single command",
			hookName: "pre-commit",
			scripts: []HookScript{
				{Name: "test", Command: "go test", Description: "Test command"},
			},
			expectError: false,
			contains:    []string{"#!/bin/sh", "go test", "Running: test", "Test command"},
		},
		{
			name:     "mixed script and command",
			hookName: "pre-commit",
			scripts: []HookScript{
				{Name: "script-test", Script: "test.sh", Description: "Script test"},
				{Name: "command-test", Command: "go test", Description: "Command test"},
			},
			expectError: false,
			contains:    []string{"test.sh", "go test", "Running: script-test", "Running: command-test"},
		},
		{
			name:     "multiple scripts",
			hookName: "pre-push",
			scripts: []HookScript{
				{Name: "build", Script: "build.sh", Description: "Build"},
				{Name: "test", Command: "npm test", Description: "Test"},
			},
			expectError: false,
			contains:    []string{"build.sh", "npm test", "Running: build", "Running: test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary git directory
			gitDir := filepath.Join(tmpDir, ".git")
			err := os.MkdirAll(gitDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create git dir: %v", err)
			}

			hm := &HookManager{
				config: &Config{
					Hooks:    map[string][]HookScript{tt.hookName: tt.scripts},
					Settings: Settings{},
				},
				gitDir: gitDir,
			}

			// Change to temp directory for working directory
			oldDir, _ := os.Getwd()
			defer os.Chdir(oldDir)
			os.Chdir(tmpDir)

			content, err := hm.generateHookScript(tt.hookName, tt.scripts)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check that content contains expected strings
			for _, expected := range tt.contains {
				if !strings.Contains(content, expected) {
					t.Errorf("Expected content to contain '%s', got:\n%s", expected, content)
				}
			}

			// Check that it's a valid shell script
			if !strings.HasPrefix(content, "#!/bin/sh") {
				t.Error("Expected script to start with shebang")
			}

			if !strings.Contains(content, "Generated by hooky") {
				t.Error("Expected script to contain hooky signature")
			}
		})
	}
}

func TestFindGitDirectory(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(string) error
		expectError bool
	}{
		{
			name: "valid git repository",
			setupFunc: func(dir string) error {
				gitDir := filepath.Join(dir, ".git")
				if err := os.MkdirAll(gitDir, 0755); err != nil {
					return err
				}
				// Initialize git repo
				cmd := exec.Command("git", "init")
				cmd.Dir = dir
				return cmd.Run()
			},
			expectError: false,
		},
		{
			name: "no git repository",
			setupFunc: func(dir string) error {
				return nil // Don't create .git directory
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			
			if err := tt.setupFunc(tmpDir); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Change to temp directory
			oldDir, _ := os.Getwd()
			defer os.Chdir(oldDir)
			os.Chdir(tmpDir)

			hm := &HookManager{}
			gitDir, err := hm.findGitDirectory()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if gitDir == "" {
				t.Error("Expected non-empty git directory")
			}

			// Check that the directory exists
			if _, err := os.Stat(gitDir); os.IsNotExist(err) {
				t.Errorf("Git directory does not exist: %s", gitDir)
			}
		})
	}
}

func TestInstallHooks(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test script
	scriptPath := filepath.Join(tmpDir, "test.sh")
	err := os.WriteFile(scriptPath, []byte("#!/bin/bash\necho test"), 0755)
	if err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Create git repository
	gitDir := filepath.Join(tmpDir, ".git")
	hooksDir := filepath.Join(gitDir, "hooks")
	err = os.MkdirAll(hooksDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create git hooks dir: %v", err)
	}

	// Change to temp directory
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	hm := &HookManager{
		config: &Config{
			Hooks: map[string][]HookScript{
				"pre-commit": {
					{Name: "test-script", Script: "test.sh", Description: "Test script"},
					{Name: "test-command", Command: "echo test", Description: "Test command"},
				},
			},
			Settings: Settings{
				BackupExisting:  true,
				BackupDirectory: ".hooky-backup",
			},
		},
		gitDir: gitDir,
	}

	err = hm.InstallHooks()
	if err != nil {
		t.Errorf("InstallHooks failed: %v", err)
		return
	}

	// Check that hook file was created
	hookPath := filepath.Join(hooksDir, "pre-commit")
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		t.Error("Hook file was not created")
		return
	}

	// Check hook content
	content, err := os.ReadFile(hookPath)
	if err != nil {
		t.Errorf("Failed to read hook file: %v", err)
		return
	}

	hookContent := string(content)
	if !strings.Contains(hookContent, "test.sh") {
		t.Error("Hook content should contain script path")
	}
	if !strings.Contains(hookContent, "echo test") {
		t.Error("Hook content should contain command")
	}
	if !strings.Contains(hookContent, "Generated by hooky") {
		t.Error("Hook content should contain hooky signature")
	}
}

func TestUninstallHooks(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create git repository
	gitDir := filepath.Join(tmpDir, ".git")
	hooksDir := filepath.Join(gitDir, "hooks")
	err := os.MkdirAll(hooksDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create git hooks dir: %v", err)
	}

	// Create test hooks (one generated by hooky, one not)
	hookyHookPath := filepath.Join(hooksDir, "pre-commit")
	hookyHookContent := `#!/bin/sh
# Generated by hooky - Do not edit manually
echo "hooky generated hook"
`
	err = os.WriteFile(hookyHookPath, []byte(hookyHookContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create hooky hook: %v", err)
	}

	otherHookPath := filepath.Join(hooksDir, "pre-push")
	otherHookContent := `#!/bin/sh
echo "other hook"
`
	err = os.WriteFile(otherHookPath, []byte(otherHookContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create other hook: %v", err)
	}

	hm := &HookManager{
		config: &Config{
			Hooks: map[string][]HookScript{
				"pre-commit": {{Name: "test", Command: "echo test", Description: "Test"}},
				"pre-push":   {{Name: "test", Command: "echo test", Description: "Test"}},
			},
			Settings: Settings{},
		},
		gitDir: gitDir,
	}

	err = hm.UninstallHooks()
	if err != nil {
		t.Errorf("UninstallHooks failed: %v", err)
		return
	}

	// Check that hooky-generated hook was removed
	if _, err := os.Stat(hookyHookPath); !os.IsNotExist(err) {
		t.Error("Hooky-generated hook should have been removed")
	}

	// Check that other hook was not removed
	if _, err := os.Stat(otherHookPath); os.IsNotExist(err) {
		t.Error("Non-hooky hook should not have been removed")
	}
}

// Helper function to run a command and capture output
func runCommand(t *testing.T, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}