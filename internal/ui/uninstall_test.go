package ui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/rsmdt/the-startup/internal/config"
	"github.com/rsmdt/the-startup/internal/installer"
)

func TestUninstallFlow(t *testing.T) {
	// Create temporary directories for test
	tempDir := t.TempDir()
	startupPath := filepath.Join(tempDir, "startup")
	claudePath := filepath.Join(tempDir, "claude")

	// Create directories
	if err := os.MkdirAll(startupPath, 0755); err != nil {
		t.Fatalf("Failed to create startup directory: %v", err)
	}
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude directory: %v", err)
	}

	// Create test files that would be installed
	testFiles := map[string]string{
		// Claude files
		filepath.Join(claudePath, "agents", "the-test.md"):         "# Test Agent",
		filepath.Join(claudePath, "commands", "s", "test.md"):      "# Test Command",
		filepath.Join(claudePath, "output-styles", "test.md"):     "# Test Style",
		// Startup files  
		filepath.Join(startupPath, "templates", "test.md"):        "# Test Template",
		filepath.Join(startupPath, "rules", "test.md"):            "# Test Rule",
		filepath.Join(startupPath, "bin", "the-startup"):          "#!/bin/bash\necho test",
		filepath.Join(startupPath, "the-startup.lock"):            "", // Will be created below
	}

	// Create directories and files
	for filePath, content := range testFiles {
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("Failed to create directory for %s: %v", filePath, err)
		}
		if filepath.Base(filePath) != "the-startup.lock" {
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to create test file %s: %v", filePath, err)
			}
		}
	}

	// Create a lockfile that references these files
	lockFile := config.LockFile{
		InstallPath: startupPath,
		ClaudePath:  claudePath,
		Files: map[string]config.FileInfo{
			"agents/the-test.md":         {Size: 12, Checksum: "abc123"},
			"commands/s/test.md":         {Size: 15, Checksum: "def456"},
			"output-styles/test.md":      {Size: 13, Checksum: "ghi789"},
			"startup/templates/test.md":  {Size: 16, Checksum: "jkl012"},
			"startup/rules/test.md":      {Size: 12, Checksum: "mno345"},
			"bin/the-startup":            {Size: 20, Checksum: "pqr678"},
		},
	}

	lockFileData, err := json.Marshal(lockFile)
	if err != nil {
		t.Fatalf("Failed to marshal lockfile: %v", err)
	}

	lockFilePath := filepath.Join(startupPath, "the-startup.lock")
	if err := os.WriteFile(lockFilePath, lockFileData, 0644); err != nil {
		t.Fatalf("Failed to create lockfile: %v", err)
	}

	// Verify all files exist before uninstall
	for filePath := range testFiles {
		if _, err := os.Stat(filePath); err != nil {
			t.Errorf("Test file %s should exist before uninstall: %v", filePath, err)
		}
	}
	if _, err := os.Stat(lockFilePath); err != nil {
		t.Errorf("Lockfile should exist: %v", err)
	}

	// Test the uninstall functionality
	t.Run("FileSelectionModel gets files from lockfile", func(t *testing.T) {
		inst := installer.New(nil, nil)
		inst.SetInstallPath(startupPath)
		inst.SetClaudePath(claudePath)

		model := NewFileSelectionModelWithMode("claude-code", claudePath, inst, nil, nil, ModeUninstall)
		
		// Should have found the files
		if len(model.selectedFiles) == 0 {
			t.Error("Expected files to be found from lockfile, got none")
		}

		// Check that only existing files are included
		for _, filePath := range model.selectedFiles {
			if _, err := os.Stat(filePath); err != nil {
				t.Errorf("Selected file %s does not exist: %v", filePath, err)
			}
		}
	})

	t.Run("performUninstall actually removes files", func(t *testing.T) {
		// Create MainModel with uninstall setup
		inst := installer.New(nil, nil)
		inst.SetInstallPath(startupPath)
		inst.SetClaudePath(claudePath)

		model := &MainModel{
			installer: inst,
			mode:      ModeUninstall,
		}
		model.fileSelectionModel = NewFileSelectionModelWithMode("claude-code", claudePath, inst, nil, nil, ModeUninstall)

		// Perform uninstall
		err := model.performUninstall()
		if err != nil {
			t.Fatalf("performUninstall failed: %v", err)
		}

		// Verify files were removed
		for filePath := range testFiles {
			if _, err := os.Stat(filePath); err == nil {
				t.Errorf("File %s should have been removed but still exists", filePath)
			}
		}
		if _, err := os.Stat(lockFilePath); err == nil {
			t.Error("Lockfile should have been removed but still exists")
		}
	})
}

func TestUninstallWithMissingFiles(t *testing.T) {
	// Test scenario where some files in lockfile don't exist
	tempDir := t.TempDir()
	startupPath := filepath.Join(tempDir, "startup")
	claudePath := filepath.Join(tempDir, "claude")

	// Create directories
	if err := os.MkdirAll(startupPath, 0755); err != nil {
		t.Fatalf("Failed to create startup directory: %v", err)
	}
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude directory: %v", err)
	}

	// Create only some of the files
	existingFile := filepath.Join(claudePath, "agents", "existing.md")
	if err := os.MkdirAll(filepath.Dir(existingFile), 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.WriteFile(existingFile, []byte("# Existing"), 0644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	// Create lockfile that references both existing and missing files
	lockFile := config.LockFile{
		InstallPath: startupPath,
		ClaudePath:  claudePath,
		Files: map[string]config.FileInfo{
			"agents/existing.md": {Size: 10, Checksum: "abc123"},
			"agents/missing.md":  {Size: 10, Checksum: "def456"}, // This file doesn't exist
		},
	}

	lockFileData, err := json.Marshal(lockFile)
	if err != nil {
		t.Fatalf("Failed to marshal lockfile: %v", err)
	}

	lockFilePath := filepath.Join(startupPath, "the-startup.lock")
	if err := os.WriteFile(lockFilePath, lockFileData, 0644); err != nil {
		t.Fatalf("Failed to create lockfile: %v", err)
	}

	// Test file selection only includes existing files
	inst := installer.New(nil, nil)
	inst.SetInstallPath(startupPath)
	inst.SetClaudePath(claudePath)

	model := NewFileSelectionModelWithMode("claude-code", claudePath, inst, nil, nil, ModeUninstall)
	
	// Should only include existing files
	foundExisting := false
	foundMissing := false
	
	for _, filePath := range model.selectedFiles {
		if filePath == existingFile {
			foundExisting = true
		}
		if filepath.Base(filePath) == "missing.md" {
			foundMissing = true
		}
	}

	if !foundExisting {
		t.Error("Should have found the existing file")
	}
	if foundMissing {
		t.Error("Should not have included the missing file")
	}
}

func TestUninstallWithNoLockfile(t *testing.T) {
	// Test scenario where no lockfile exists
	tempDir := t.TempDir()
	startupPath := filepath.Join(tempDir, "startup")
	claudePath := filepath.Join(tempDir, "claude")

	// Create directories but no lockfile
	if err := os.MkdirAll(startupPath, 0755); err != nil {
		t.Fatalf("Failed to create startup directory: %v", err)
	}
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude directory: %v", err)
	}

	// Test file selection with no lockfile
	inst := installer.New(nil, nil)
	inst.SetInstallPath(startupPath)
	inst.SetClaudePath(claudePath)

	model := NewFileSelectionModelWithMode("claude-code", claudePath, inst, nil, nil, ModeUninstall)
	
	// Should have no files to remove
	if len(model.selectedFiles) != 0 {
		t.Errorf("Expected no files when lockfile missing, got %d files", len(model.selectedFiles))
	}

	// Should show "Go back to path selection" option
	if len(model.choices) != 1 || model.choices[0] != "Go back to path selection" {
		t.Errorf("Expected 'Go back to path selection' option when no files found, got %v", model.choices)
	}
}

func TestRemoveStartupDirectoryIfEmpty(t *testing.T) {
	tempDir := t.TempDir()
	startupPath := filepath.Join(tempDir, "startup")

	// Create empty startup directory
	if err := os.MkdirAll(startupPath, 0755); err != nil {
		t.Fatalf("Failed to create startup directory: %v", err)
	}

	// Create model and test directory removal
	inst := installer.New(nil, nil)
	inst.SetInstallPath(startupPath)

	model := &MainModel{
		installer: inst,
	}

	// Directory should exist before cleanup
	if _, err := os.Stat(startupPath); err != nil {
		t.Fatalf("Startup directory should exist: %v", err)
	}

	// Remove empty directory
	model.removeStartupDirectoryIfEmpty()

	// Directory should be removed
	if _, err := os.Stat(startupPath); err == nil {
		t.Error("Empty startup directory should have been removed")
	}
}

func TestRemoveStartupDirectoryIfEmptyWithFiles(t *testing.T) {
	tempDir := t.TempDir()
	startupPath := filepath.Join(tempDir, "startup")

	// Create startup directory with a file
	if err := os.MkdirAll(startupPath, 0755); err != nil {
		t.Fatalf("Failed to create startup directory: %v", err)
	}

	testFile := filepath.Join(startupPath, "remaining-file.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create model and test directory removal
	inst := installer.New(nil, nil)
	inst.SetInstallPath(startupPath)

	model := &MainModel{
		installer: inst,
	}

	// Try to remove directory (should not remove it because it has files)
	model.removeStartupDirectoryIfEmpty()

	// Directory should still exist because it's not empty
	if _, err := os.Stat(startupPath); err != nil {
		t.Error("Non-empty startup directory should not have been removed")
	}

	// File should still exist
	if _, err := os.Stat(testFile); err != nil {
		t.Error("File in startup directory should still exist")
	}
}