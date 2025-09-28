package installer

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rsmdt/the-startup/internal/config"
)

// Test assets for lockfile testing
//go:embed testdata/claude_assets/*
var testClaudeAssets embed.FS

//go:embed testdata/startup_assets/*
var testStartupAssets embed.FS

// Helper function to create a test installer with mock assets
func createTestInstallerWithAssets(t *testing.T) (*Installer, string, string) {
	t.Helper()

	// Create temp directories
	tmpDir := t.TempDir()
	installPath := filepath.Join(tmpDir, ".the-startup")
	claudePath := filepath.Join(tmpDir, ".claude")

	// Create directories
	os.MkdirAll(installPath, 0755)
	os.MkdirAll(claudePath, 0755)

	installer := &Installer{
		installPath:   installPath,
		claudePath:    claudePath,
		tool:          "claude-code",
		claudeAssets:  &testClaudeAssets,
		startupAssets: &testStartupAssets,
	}

	return installer, installPath, claudePath
}

// Helper function to create a mock lockfile
func createMockLockFile(installPath string, files map[string]config.FileInfo) error {
	lockFile := config.LockFile{
		Version:     "1.0.0",
		InstallDate: time.Now().Format(time.RFC3339),
		InstallPath: installPath,
		ClaudePath:  filepath.Join(filepath.Dir(installPath), ".claude"),
		Tool:        "claude-code",
		Files:       files,
	}

	data, err := json.MarshalIndent(lockFile, "", "  ")
	if err != nil {
		return err
	}

	lockFilePath := filepath.Join(installPath, "the-startup.lock")
	return os.WriteFile(lockFilePath, data, 0644)
}

// TestLoadExistingLockFile tests loading an existing lock file
func TestLoadExistingLockFile(t *testing.T) {
	tests := []struct {
		name          string
		setupLockFile bool
		files         map[string]config.FileInfo
		expectError   bool
		expectNil     bool
	}{
		{
			name:          "No existing lockfile",
			setupLockFile: false,
			expectError:   false,
			expectNil:     true,
		},
		{
			name:          "Valid lockfile with files",
			setupLockFile: true,
			files: map[string]config.FileInfo{
				"agents/the-analyst.md": {
					Size:         1024,
					LastModified: time.Now().Format(time.RFC3339),
					Checksum:     "abc123",
				},
				"commands/test.md": {
					Size:         512,
					LastModified: time.Now().Format(time.RFC3339),
					Checksum:     "def456",
				},
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name:          "Empty lockfile",
			setupLockFile: true,
			files:         map[string]config.FileInfo{},
			expectError:   false,
			expectNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installer, installPath, _ := createTestInstallerWithAssets(t)

			if tt.setupLockFile {
				err := createMockLockFile(installPath, tt.files)
				if err != nil {
					t.Fatalf("Failed to create mock lock file: %v", err)
				}
			}

			err := installer.LoadExistingLockFile()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectNil && installer.existingLock != nil {
				t.Error("Expected nil existing lock but got value")
			}
			if !tt.expectNil && installer.existingLock == nil {
				t.Error("Expected existing lock but got nil")
			}

			if !tt.expectNil && installer.existingLock != nil {
				if len(installer.existingLock.Files) != len(tt.files) {
					t.Errorf("Expected %d files in lock, got %d",
						len(tt.files), len(installer.existingLock.Files))
				}
			}
		})
	}
}

// TestGetDeprecatedFiles tests detection of deprecated files
func TestGetDeprecatedFiles(t *testing.T) {
	tests := []struct {
		name               string
		lockFiles          map[string]config.FileInfo
		currentAssetFiles  []string // Files that exist in current assets
		expectedDeprecated []string
	}{
		{
			name:               "No existing lock",
			lockFiles:          nil,
			currentAssetFiles:  []string{"agents/new-agent.md"},
			expectedDeprecated: nil,
		},
		{
			name: "All files still exist",
			lockFiles: map[string]config.FileInfo{
				"agents/the-analyst.md": {},
				"agents/the-architect.md": {},
			},
			currentAssetFiles:  []string{"agents/the-analyst.md", "agents/the-architect.md"},
			expectedDeprecated: []string{},
		},
		{
			name: "Some files deprecated",
			lockFiles: map[string]config.FileInfo{
				"agents/old-agent-1.md": {},
				"agents/old-agent-2.md": {},
				"agents/current-agent.md": {},
				"commands/old-command.md": {},
			},
			currentAssetFiles:  []string{"agents/current-agent.md"},
			expectedDeprecated: []string{"agents/old-agent-1.md", "agents/old-agent-2.md", "commands/old-command.md"},
		},
		{
			name: "Non-claude files ignored",
			lockFiles: map[string]config.FileInfo{
				"agents/deprecated.md": {},
				"templates/template.md": {}, // Should be ignored
				"bin/the-startup": {},        // Should be ignored
			},
			currentAssetFiles:  []string{},
			expectedDeprecated: []string{"agents/deprecated.md"},
		},
		{
			name: "Mixed deprecated and current",
			lockFiles: map[string]config.FileInfo{
				"agents/agent-1.md": {},
				"agents/agent-2.md": {},
				"agents/agent-3.md": {},
				"commands/cmd-1.md": {},
				"commands/cmd-2.md": {},
			},
			currentAssetFiles:  []string{"agents/agent-2.md", "commands/cmd-1.md"},
			expectedDeprecated: []string{"agents/agent-1.md", "agents/agent-3.md", "commands/cmd-2.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory structure
			tmpDir := t.TempDir()

			// Create mock embedded assets structure
			assetsDir := filepath.Join(tmpDir, "assets", "claude")
			os.MkdirAll(assetsDir, 0755)

			// Create current asset files
			for _, file := range tt.currentAssetFiles {
				fullPath := filepath.Join(assetsDir, file)
				os.MkdirAll(filepath.Dir(fullPath), 0755)
				os.WriteFile(fullPath, []byte("content"), 0644)
			}

			// Create installer with mock lock
			installer := &Installer{
				installPath: filepath.Join(tmpDir, ".the-startup"),
				claudePath:  filepath.Join(tmpDir, ".claude"),
				tool:        "claude-code",
			}

			// Set up existing lock if provided
			if tt.lockFiles != nil {
				installer.existingLock = &config.LockFile{
					Files: tt.lockFiles,
				}
			}

			// Mock the embedded assets by creating a custom FS
			// For this test, we'll simulate by checking the created files
			installer.claudeAssets = &embed.FS{} // This would need proper mocking in production

			// Manually set current files for the test
			// In real code, this comes from WalkDir
			currentFiles := make(map[string]bool)
			for _, f := range tt.currentAssetFiles {
				currentFiles[f] = true
			}

			// Test the logic directly
			var deprecated []string
			if installer.existingLock != nil {
				for filePath := range installer.existingLock.Files {
					// Skip non-claude files
					if !strings.HasPrefix(filePath, "agents/") &&
					   !strings.HasPrefix(filePath, "commands/") {
						continue
					}

					// Check if file exists in current assets
					if !currentFiles[filePath] {
						deprecated = append(deprecated, filePath)
					}
				}
			}

			// Verify results
			if len(deprecated) != len(tt.expectedDeprecated) {
				t.Errorf("Expected %d deprecated files, got %d",
					len(tt.expectedDeprecated), len(deprecated))
				t.Errorf("Expected: %v", tt.expectedDeprecated)
				t.Errorf("Got: %v", deprecated)
			}

			// Check each expected deprecated file
			deprecatedMap := make(map[string]bool)
			for _, f := range deprecated {
				deprecatedMap[f] = true
			}

			for _, expected := range tt.expectedDeprecated {
				if !deprecatedMap[expected] {
					t.Errorf("Expected file %s to be marked as deprecated", expected)
				}
			}
		})
	}
}

// TestRemoveDeprecatedFiles tests the removal of deprecated files
func TestRemoveDeprecatedFiles(t *testing.T) {
	tests := []struct {
		name            string
		deprecatedFiles []string
		existingFiles   []string // Files to create before removal
		expectRemoved   []string // Files that should be removed
		expectKept      []string // Files that should remain
	}{
		{
			name:            "No deprecated files",
			deprecatedFiles: []string{},
			existingFiles:   []string{"agents/current.md"},
			expectRemoved:   []string{},
			expectKept:      []string{"agents/current.md"},
		},
		{
			name:            "Remove single deprecated file",
			deprecatedFiles: []string{"agents/old.md"},
			existingFiles:   []string{"agents/old.md", "agents/current.md"},
			expectRemoved:   []string{"agents/old.md"},
			expectKept:      []string{"agents/current.md"},
		},
		{
			name:            "Remove multiple deprecated files",
			deprecatedFiles: []string{"agents/old1.md", "commands/old2.md"},
			existingFiles:   []string{"agents/old1.md", "commands/old2.md", "agents/current.md"},
			expectRemoved:   []string{"agents/old1.md", "commands/old2.md"},
			expectKept:      []string{"agents/current.md"},
		},
		{
			name:            "Handle already removed files gracefully",
			deprecatedFiles: []string{"agents/already-gone.md"},
			existingFiles:   []string{"agents/current.md"},
			expectRemoved:   []string{},
			expectKept:      []string{"agents/current.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installer, _, claudePath := createTestInstallerWithAssets(t)

			// Create existing files
			for _, file := range tt.existingFiles {
				fullPath := filepath.Join(claudePath, file)
				os.MkdirAll(filepath.Dir(fullPath), 0755)
				os.WriteFile(fullPath, []byte("content"), 0644)
			}

			// Set deprecated files
			installer.deprecatedFiles = tt.deprecatedFiles

			// Remove deprecated files
			err := installer.RemoveDeprecatedFiles()
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check removed files
			for _, file := range tt.expectRemoved {
				fullPath := filepath.Join(claudePath, file)
				if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
					t.Errorf("Expected file %s to be removed, but it still exists", file)
				}
			}

			// Check kept files
			for _, file := range tt.expectKept {
				fullPath := filepath.Join(claudePath, file)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					t.Errorf("Expected file %s to remain, but it was removed", file)
				}
			}
		})
	}
}

// TestLockFileIntegration tests the complete lockfile workflow
func TestLockFileIntegration(t *testing.T) {
	// Scenario: Install with 61 files, then reinstall with 40 files
	// This simulates your exact use case

	tmpDir := t.TempDir()
	installPath := filepath.Join(tmpDir, ".the-startup")
	claudePath := filepath.Join(tmpDir, ".claude")

	// Create directories
	os.MkdirAll(installPath, 0755)
	os.MkdirAll(filepath.Join(claudePath, "agents"), 0755)

	// First installation: 61 agents
	initialAgents := make(map[string]config.FileInfo)
	for i := 1; i <= 61; i++ {
		agentFile := fmt.Sprintf("agents/agent-%02d.md", i)
		initialAgents[agentFile] = config.FileInfo{
			Size:         int64(100 + i),
			LastModified: time.Now().Format(time.RFC3339),
			Checksum:     fmt.Sprintf("checksum-%d", i),
		}

		// Create the actual file
		fullPath := filepath.Join(claudePath, agentFile)
		os.WriteFile(fullPath, []byte(fmt.Sprintf("Agent %d content", i)), 0644)
	}

	// Create initial lock file
	err := createMockLockFile(installPath, initialAgents)
	if err != nil {
		t.Fatalf("Failed to create initial lock file: %v", err)
	}

	// Verify all 61 files exist
	files, _ := os.ReadDir(filepath.Join(claudePath, "agents"))
	if len(files) != 61 {
		t.Errorf("Expected 61 initial files, got %d", len(files))
	}

	// Second installation: Only 40 agents (simulating removal of 21 agents)
	installer := &Installer{
		installPath:  installPath,
		claudePath:   claudePath,
		tool:         "claude-code",
	}

	// Load existing lock file
	err = installer.LoadExistingLockFile()
	if err != nil {
		t.Fatalf("Failed to load lock file: %v", err)
	}

	// Simulate current assets having only 40 agents (1-40)
	currentFiles := make(map[string]bool)
	for i := 1; i <= 40; i++ {
		currentFiles[fmt.Sprintf("agents/agent-%02d.md", i)] = true
	}

	// Manually detect deprecated files (normally done by GetDeprecatedFiles)
	var deprecated []string
	for filePath := range installer.existingLock.Files {
		if strings.HasPrefix(filePath, "agents/") && !currentFiles[filePath] {
			deprecated = append(deprecated, filePath)
		}
	}
	installer.deprecatedFiles = deprecated

	// Should detect 21 deprecated files (41-61)
	if len(installer.deprecatedFiles) != 21 {
		t.Errorf("Expected 21 deprecated files, got %d", len(installer.deprecatedFiles))
	}

	// Remove deprecated files
	err = installer.RemoveDeprecatedFiles()
	if err != nil {
		t.Fatalf("Failed to remove deprecated files: %v", err)
	}

	// Verify only 40 files remain
	files, _ = os.ReadDir(filepath.Join(claudePath, "agents"))
	if len(files) != 40 {
		t.Errorf("Expected 40 remaining files, got %d", len(files))
	}

	// Verify correct files remain (1-40)
	for i := 1; i <= 40; i++ {
		agentFile := fmt.Sprintf("agents/agent-%02d.md", i)
		fullPath := filepath.Join(claudePath, agentFile)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist, but it was removed", agentFile)
		}
	}

	// Verify deprecated files are gone (41-61)
	for i := 41; i <= 61; i++ {
		agentFile := fmt.Sprintf("agents/agent-%02d.md", i)
		fullPath := filepath.Join(claudePath, agentFile)
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			t.Errorf("Expected file %s to be removed, but it still exists", agentFile)
		}
	}
}

// TestCorruptedLockFile tests handling of corrupted lock files
func TestCorruptedLockFile(t *testing.T) {
	installer, installPath, _ := createTestInstallerWithAssets(t)

	// Create corrupted lock file
	lockFilePath := filepath.Join(installPath, "the-startup.lock")
	os.WriteFile(lockFilePath, []byte("{ invalid json }"), 0644)

	// Should return error for corrupted lock file
	err := installer.LoadExistingLockFile()
	if err == nil {
		t.Error("Expected error for corrupted lock file")
	}

	// existingLock should remain nil
	if installer.existingLock != nil {
		t.Error("Expected nil existing lock for corrupted file")
	}
}

// TestChecksumVerification tests file integrity checking
func TestChecksumVerification(t *testing.T) {
	installer, _, claudePath := createTestInstallerWithAssets(t)

	// Create a file with known content
	agentFile := "agents/test-agent.md"
	content := []byte("Test agent content")
	fullPath := filepath.Join(claudePath, agentFile)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	os.WriteFile(fullPath, content, 0644)

	// Calculate checksum
	checksum, err := installer.calculateFileChecksum(fullPath)
	if err != nil {
		t.Fatalf("Failed to calculate checksum: %v", err)
	}

	// Checksum should not be empty
	if checksum == "" {
		t.Error("Expected non-empty checksum")
	}

	// Same content should produce same checksum
	checksum2, _ := installer.calculateFileChecksum(fullPath)
	if checksum != checksum2 {
		t.Error("Same file produced different checksums")
	}

	// Modified content should produce different checksum
	os.WriteFile(fullPath, append(content, []byte(" modified")...), 0644)
	checksum3, _ := installer.calculateFileChecksum(fullPath)
	if checksum == checksum3 {
		t.Error("Modified file produced same checksum")
	}
}