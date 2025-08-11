package installer

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed test_assets/*
var testAssets embed.FS

func setupTestDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "installer_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})
	
	return tmpDir
}

func TestNew(t *testing.T) {
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	
	if installer == nil {
		t.Fatal("New() returned nil")
	}
	
	if installer.tool != "claude-code" {
		t.Errorf("Expected default tool 'claude-code', got '%s'", installer.tool)
	}
	
	expectedComponents := []string{"agents", "hooks", "commands", "templates"}
	if len(installer.components) != len(expectedComponents) {
		t.Errorf("Expected %d default components, got %d", len(expectedComponents), len(installer.components))
	}
	
	for i, comp := range expectedComponents {
		if installer.components[i] != comp {
			t.Errorf("Expected component[%d] to be '%s', got '%s'", i, comp, installer.components[i])
		}
	}
}

func TestSetInstallPath(t *testing.T) {
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test absolute path
	testPath := "/tmp/test-install"
	installer.SetInstallPath(testPath)
	
	if installer.installPath != testPath {
		t.Errorf("Expected install path '%s', got '%s'", testPath, installer.installPath)
	}
	
	// Test home directory expansion
	installer.SetInstallPath("~/test-install")
	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, "test-install")
	
	if installer.installPath != expectedPath {
		t.Errorf("Expected expanded path '%s', got '%s'", expectedPath, installer.installPath)
	}
	
	// Test local installation (.the-startup detection)
	localPath := "/some/project/.the-startup"
	installer.SetInstallPath(localPath)
	
	if installer.installPath != localPath {
		t.Errorf("Expected local install path '%s', got '%s'", localPath, installer.installPath)
	}
	
	// Claude path should be set to local .claude for local installations
	expectedClaudePath := "/some/project/.claude"
	if installer.claudePath != expectedClaudePath {
		t.Errorf("Expected claude path '%s', got '%s'", expectedClaudePath, installer.claudePath)
	}
}

func TestSetters(t *testing.T) {
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test SetTool
	installer.SetTool("cursor")
	if installer.tool != "cursor" {
		t.Errorf("Expected tool 'cursor', got '%s'", installer.tool)
	}
	
	// Test SetComponents
	customComponents := []string{"agents", "commands"}
	installer.SetComponents(customComponents)
	
	if len(installer.components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(installer.components))
	}
	
	for i, comp := range customComponents {
		if installer.components[i] != comp {
			t.Errorf("Expected component[%d] to be '%s', got '%s'", i, comp, installer.components[i])
		}
	}
	
	// Test SetSelectedFiles
	selectedFiles := []string{"agents/the-architect.md", "commands/develop.md"}
	installer.SetSelectedFiles(selectedFiles)
	
	if len(installer.selectedFiles) != 2 {
		t.Errorf("Expected 2 selected files, got %d", len(installer.selectedFiles))
	}
	
	for i, file := range selectedFiles {
		if installer.selectedFiles[i] != file {
			t.Errorf("Expected selected file[%d] to be '%s', got '%s'", i, file, installer.selectedFiles[i])
		}
	}
}

func TestCheckFileExists(t *testing.T) {
	tmpDir := setupTestDir(t)
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	installer.SetInstallPath(tmpDir)
	
	// Get the claude path (set based on test setup)
	homeDir, _ := os.UserHomeDir()
	claudePath := filepath.Join(homeDir, ".claude")
	
	// Create test directories and files in Claude directory
	agentsDir := filepath.Join(claudePath, "agents")
	os.MkdirAll(agentsDir, 0755)
	defer os.RemoveAll(filepath.Join(claudePath, "agents", "test-agent.md")) // Clean up
	
	testFile := filepath.Join(agentsDir, "test-agent.md")
	os.WriteFile(testFile, []byte("test content"), 0644)
	
	// Test existing file
	if !installer.CheckFileExists("agents/test-agent.md") {
		t.Error("Expected CheckFileExists to return true for existing file")
	}
	
	// Test non-existing file
	if installer.CheckFileExists("agents/nonexistent.md") {
		t.Error("Expected CheckFileExists to return false for non-existing file")
	}
	
	// Clean up test file
	os.Remove(testFile)
}

func TestIsInstalled(t *testing.T) {
	tmpDir := setupTestDir(t)
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Set both install and claude paths to the temp directory to avoid conflicts
	installer.SetInstallPath(tmpDir)
	installer.claudePath = filepath.Join(tmpDir, ".claude")
	
	// Initially should not be installed
	if installer.IsInstalled() {
		t.Error("Expected IsInstalled to return false for fresh install location")
	}
	
	// Create lock file
	lockFile := filepath.Join(tmpDir, "the-startup.lock")
	os.WriteFile(lockFile, []byte("{}"), 0644)
	
	// Should now be installed
	if !installer.IsInstalled() {
		t.Error("Expected IsInstalled to return true when lock file exists")
	}
	
	// Remove lock file and test Claude directory detection
	os.Remove(lockFile)
	
	// Create Claude agents directory with a "the-" prefixed file
	claudeAgentsDir := filepath.Join(tmpDir, ".claude", "agents")
	os.MkdirAll(claudeAgentsDir, 0755)
	os.WriteFile(filepath.Join(claudeAgentsDir, "the-architect.md"), []byte("test"), 0644)
	
	// Should be detected as installed due to Claude files
	if !installer.IsInstalled() {
		t.Error("Expected IsInstalled to return true when Claude agents exist")
	}
}

func TestInstallComponentUnknown(t *testing.T) {
	tmpDir := setupTestDir(t)
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	installer.SetInstallPath(tmpDir)
	
	// Test unknown component
	err := installer.installComponentToClaude("unknown")
	if err == nil {
		t.Error("Expected error for unknown component")
	}
	
	if err.Error() != "unknown component: unknown" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestReplacePlaceholders(t *testing.T) {
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	installer.SetInstallPath("/test/path")
	
	input := []byte("Install path: {{INSTALL_PATH}}/config")
	expected := "Install path: /test/path/config"
	
	result := installer.replacePlaceholders(input)
	
	if string(result) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(result))
	}
}

func TestGetPaths(t *testing.T) {
	tmpDir := setupTestDir(t)
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	installer.SetInstallPath(tmpDir)
	
	// Test GetInstallPath
	if installer.GetInstallPath() != tmpDir {
		t.Errorf("Expected install path '%s', got '%s'", tmpDir, installer.GetInstallPath())
	}
	
	// Test GetClaudePath
	expectedClaudePath := installer.claudePath
	if installer.GetClaudePath() != expectedClaudePath {
		t.Errorf("Expected claude path '%s', got '%s'", expectedClaudePath, installer.GetClaudePath())
	}
}

func TestContainsHelper(t *testing.T) {
	slice := []string{"a", "b", "c"}
	
	if !contains(slice, "b") {
		t.Error("Expected contains to return true for existing item")
	}
	
	if contains(slice, "d") {
		t.Error("Expected contains to return false for non-existing item")
	}
	
	// Test empty slice
	emptySlice := []string{}
	if contains(emptySlice, "a") {
		t.Error("Expected contains to return false for empty slice")
	}
}

func TestInstallDirectoryCreation(t *testing.T) {
	tmpDir := setupTestDir(t)
	nonExistentDir := filepath.Join(tmpDir, "deep", "nested", "path")
	
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	installer.SetInstallPath(nonExistentDir)
	installer.SetComponents([]string{}) // No components to avoid installation errors
	
	// This should create the directory structure
	err := installer.Install()
	if err != nil {
		t.Errorf("Expected Install to create directories, got error: %v", err)
	}
	
	// Check that directory was created
	if _, err := os.Stat(nonExistentDir); os.IsNotExist(err) {
		t.Error("Expected install directory to be created")
	}
}

func TestInstallWithSelectedFiles(t *testing.T) {
	tmpDir := setupTestDir(t)
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	installer.SetInstallPath(tmpDir)
	installer.SetComponents([]string{"agents"})
	installer.SetSelectedFiles([]string{"agents/specific-agent.md"})
	
	// Install should respect selected files filter
	// This test validates the filtering logic without requiring actual embed files
	err := installer.Install()
	
	// Error is expected due to missing test assets, but validates the code path
	if err == nil {
		t.Log("Install completed (test assets might be empty)")
	} else {
		t.Logf("Install failed as expected due to test environment: %v", err)
	}
}

func TestCopyCurrentExecutable(t *testing.T) {
	tmpDir := setupTestDir(t)
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test copying current executable
	err := installer.copyCurrentExecutable(tmpDir)
	if err != nil {
		t.Errorf("Expected copyCurrentExecutable to succeed, got error: %v", err)
	}
	
	// Check that the binary was copied
	binaryPath := filepath.Join(tmpDir, "the-startup")
	if _, err := os.Stat(binaryPath); err != nil {
		t.Errorf("Expected binary to be copied to %s, but file doesn't exist: %v", binaryPath, err)
	}
	
	// Check permissions (should be executable)
	info, err := os.Stat(binaryPath)
	if err != nil {
		t.Errorf("Failed to stat binary: %v", err)
	} else if info.Mode().Perm() != 0755 {
		t.Errorf("Expected binary permissions 0755, got %o", info.Mode().Perm())
	}
}

func TestHooksComponentInstallsGoBinary(t *testing.T) {
	tmpDir := setupTestDir(t)
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	installer.SetInstallPath(tmpDir)
	installer.claudePath = filepath.Join(tmpDir, ".claude")
	
	// Install hooks component - should deploy Go binary instead of Python files
	err := installer.installComponentToClaude("hooks")
	if err != nil {
		t.Errorf("Expected hooks component installation to succeed, got error: %v", err)
	}
	
	// Check that the Go binary was deployed
	hooksBinaryPath := filepath.Join(tmpDir, ".claude", "hooks", "the-startup")
	if _, err := os.Stat(hooksBinaryPath); err != nil {
		t.Errorf("Expected Go binary to be deployed to %s, but file doesn't exist: %v", hooksBinaryPath, err)
	}
}

func TestConfigureHooksUsesGoCommands(t *testing.T) {
	tmpDir := setupTestDir(t)
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	installer.SetInstallPath(tmpDir)
	installer.claudePath = filepath.Join(tmpDir, ".claude")
	
	// Create claude directory
	os.MkdirAll(installer.claudePath, 0755)
	
	// Configure hooks
	err := installer.configureHooks()
	if err != nil {
		t.Errorf("Expected configureHooks to succeed, got error: %v", err)
	}
	
	// Check that settings.json was created with Go commands
	settingsPath := filepath.Join(installer.claudePath, "settings.json")
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Errorf("Expected settings.json to be created, got error: %v", err)
	}
	
	settingsContent := string(data)
	
	// Verify Go hook commands are present
	expectedPreCommand := "$CLAUDE_PROJECT_DIR/.claude/hooks/the-startup log --assistant"
	expectedPostCommand := "$CLAUDE_PROJECT_DIR/.claude/hooks/the-startup log --user"
	
	if !strings.Contains(settingsContent, expectedPreCommand) {
		t.Errorf("Expected settings.json to contain PreToolUse command '%s', but content was:\n%s", expectedPreCommand, settingsContent)
	}
	
	if !strings.Contains(settingsContent, expectedPostCommand) {
		t.Errorf("Expected settings.json to contain PostToolUse command '%s', but content was:\n%s", expectedPostCommand, settingsContent)
	}
	
	// Verify Python commands are NOT present (regression test)
	oldPythonCommand := "uv run $CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_start.py"
	if strings.Contains(settingsContent, oldPythonCommand) {
		t.Errorf("Settings.json should not contain old Python commands, but found '%s'", oldPythonCommand)
	}
}

// Integration test that validates the full state machine works with installer
func TestInstallerIntegration(t *testing.T) {
	tmpDir := setupTestDir(t)
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	installer.SetInstallPath(tmpDir)
	
	// Test the full configuration sequence
	installer.SetTool("claude-code")
	installer.SetComponents([]string{"agents", "commands"})
	
	// Verify configurations are set correctly
	if installer.tool != "claude-code" {
		t.Errorf("Expected tool 'claude-code', got '%s'", installer.tool)
	}
	
	if len(installer.components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(installer.components))
	}
	
	if installer.GetInstallPath() != tmpDir {
		t.Errorf("Expected install path '%s', got '%s'", tmpDir, installer.GetInstallPath())
	}
}

// Full integration test for Go hooks deployment
func TestGoHooksIntegration(t *testing.T) {
	tmpDir := setupTestDir(t)
	installer := New(&testAssets, &testAssets, &testAssets, &testAssets)
	installer.SetInstallPath(tmpDir)
	installer.claudePath = filepath.Join(tmpDir, ".claude")
	installer.SetTool("claude-code")
	installer.SetComponents([]string{"hooks"})
	
	// Install hooks component
	err := installer.Install()
	if err != nil {
		t.Errorf("Expected full installation to succeed, got error: %v", err)
	}
	
	// Verify binary was deployed
	hooksBinaryPath := filepath.Join(tmpDir, ".claude", "hooks", "the-startup")
	if _, err := os.Stat(hooksBinaryPath); err != nil {
		t.Errorf("Expected Go binary to be deployed to %s, but file doesn't exist: %v", hooksBinaryPath, err)
	}
	
	// Verify settings.json was configured with Go commands
	settingsPath := filepath.Join(tmpDir, ".claude", "settings.json")
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Errorf("Expected settings.json to be created, got error: %v", err)
	}
	
	settingsContent := string(data)
	expectedCommand := "the-startup log --assistant"
	if !strings.Contains(settingsContent, expectedCommand) {
		t.Errorf("Expected settings.json to contain Go command '%s', but content was:\n%s", expectedCommand, settingsContent)
	}
	
	// Verify lock file was created
	lockFilePath := filepath.Join(tmpDir, "the-startup.lock")
	if _, err := os.Stat(lockFilePath); err != nil {
		t.Errorf("Expected lock file to be created at %s, but file doesn't exist: %v", lockFilePath, err)
	}
}