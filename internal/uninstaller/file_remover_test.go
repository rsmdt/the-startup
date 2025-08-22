package uninstaller

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rsmdt/the-startup/internal/config"
)

func setupTestRemover(t *testing.T) (*SafeFileRemover, string, string, *config.LockFile) {
	tmpDir := t.TempDir()
	installPath := filepath.Join(tmpDir, ".the-startup")
	claudePath := filepath.Join(tmpDir, ".claude")
	
	// Create directories
	if err := os.MkdirAll(installPath, 0755); err != nil {
		t.Fatalf("Failed to create install dir: %v", err)
	}
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude dir: %v", err)
	}

	// Create a test lock file
	lockFile := &config.LockFile{
		Version:     "1.0.0",
		InstallDate: time.Now().Format(time.RFC3339),
		InstallPath: installPath,
		ClaudePath:  claudePath,
		Tool:        "claude-code",
		Components:  []string{"agents", "commands", "templates"},
		Files:       make(map[string]config.FileInfo),
	}

	remover := NewSafeFileRemover(installPath, claudePath, lockFile)
	return remover, installPath, claudePath, lockFile
}

func createTestFile(t *testing.T, path, content string) string {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create directory for %s: %v", path, err)
	}
	
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file %s: %v", path, err)
	}

	// Calculate checksum
	hasher := sha256.New()
	hasher.Write([]byte(content))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getFileInfo(t *testing.T, path string, checksum string) config.FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Failed to get file info for %s: %v", path, err)
	}
	
	return config.FileInfo{
		Size:         info.Size(),
		LastModified: info.ModTime().Format(time.RFC3339),
		Checksum:     checksum,
	}
}

func TestNewSafeFileRemover(t *testing.T) {
	installPath := "/test/install"
	claudePath := "/test/claude"
	lockFile := &config.LockFile{}
	
	remover := NewSafeFileRemover(installPath, claudePath, lockFile)
	
	if remover == nil {
		t.Fatal("NewSafeFileRemover returned nil")
	}
	
	if remover.installPath != installPath {
		t.Errorf("Expected install path %s, got %s", installPath, remover.installPath)
	}
	
	if remover.claudePath != claudePath {
		t.Errorf("Expected claude path %s, got %s", claudePath, remover.claudePath)
	}
	
	if remover.lockFile != lockFile {
		t.Error("Lock file reference not set correctly")
	}
	
	if remover.dryRun {
		t.Error("Expected dry-run to be false by default")
	}
}

func TestSetDryRun(t *testing.T) {
	remover, _, _, _ := setupTestRemover(t)
	
	if remover.dryRun {
		t.Error("Expected initial dry-run to be false")
	}
	
	remover.SetDryRun(true)
	if !remover.dryRun {
		t.Error("Expected dry-run to be true after SetDryRun(true)")
	}
	
	remover.SetDryRun(false)
	if remover.dryRun {
		t.Error("Expected dry-run to be false after SetDryRun(false)")
	}
}

func TestRemoveFiles_Interface(t *testing.T) {
	remover, installPath, _, _ := setupTestRemover(t)
	
	// Create test files
	testFile1 := filepath.Join(installPath, "test1.txt")
	createTestFile(t, testFile1, "content1")
	
	testFile2 := filepath.Join(installPath, "test2.txt")
	createTestFile(t, testFile2, "content2")
	
	// Test the FileRemover interface implementation
	errors := remover.RemoveFiles([]string{testFile1, testFile2})
	
	if len(errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(errors))
	}
	
	// Verify files are removed
	if _, err := os.Stat(testFile1); !os.IsNotExist(err) {
		t.Error("Expected test file 1 to be removed")
	}
	
	if _, err := os.Stat(testFile2); !os.IsNotExist(err) {
		t.Error("Expected test file 2 to be removed")
	}
}

func TestRemoveDirectory_Interface(t *testing.T) {
	remover, installPath, _, _ := setupTestRemover(t)
	
	// Create empty test directory
	testDir := filepath.Join(installPath, "empty_dir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	
	// Test removing empty directory
	err := remover.RemoveDirectory(testDir)
	if err != nil {
		t.Errorf("Expected no error removing empty directory, got: %v", err)
	}
	
	// Verify directory is removed
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		t.Error("Expected empty directory to be removed")
	}
}

func TestRemoveDirectory_NotEmpty(t *testing.T) {
	remover, installPath, _, _ := setupTestRemover(t)
	
	// Create directory with content
	testDir := filepath.Join(installPath, "nonempty_dir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	
	// Add a file to make it non-empty
	testFile := filepath.Join(testDir, "content.txt")
	createTestFile(t, testFile, "content")
	
	// Test removing non-empty directory
	err := remover.RemoveDirectory(testDir)
	if err == nil {
		t.Error("Expected error removing non-empty directory")
	}
	
	if !strings.Contains(err.Error(), "not empty") {
		t.Errorf("Expected 'not empty' error, got: %v", err)
	}
	
	// Verify directory still exists
	if _, err := os.Stat(testDir); err != nil {
		t.Error("Expected non-empty directory to remain")
	}
}

func TestRemoveAllFiles_EmptyLockFile(t *testing.T) {
	remover, _, _, _ := setupTestRemover(t)
	
	summary, err := remover.RemoveAllFiles()
	if err != nil {
		t.Errorf("Expected no error with empty lock file, got: %v", err)
	}
	
	if summary.FilesProcessed != 0 {
		t.Errorf("Expected 0 files processed, got %d", summary.FilesProcessed)
	}
	
	if len(summary.FileResults) != 0 {
		t.Errorf("Expected 0 file results, got %d", len(summary.FileResults))
	}
}

func TestRemoveAllFiles_NoLockFile(t *testing.T) {
	tmpDir := t.TempDir()
	remover := NewSafeFileRemover(tmpDir, tmpDir, nil)
	
	summary, err := remover.RemoveAllFiles()
	if err == nil {
		t.Error("Expected error when no lock file provided")
	}
	
	if summary != nil {
		t.Error("Expected nil summary when error occurs")
	}
}

func TestRemoveAllFiles_SuccessfulRemoval(t *testing.T) {
	remover, installPath, claudePath, lockFile := setupTestRemover(t)
	
	// Create test files
	agentPath := filepath.Join(claudePath, "agents", "test-agent.md")
	agentContent := "# Test Agent"
	agentChecksum := createTestFile(t, agentPath, agentContent)
	
	templatePath := filepath.Join(installPath, "templates", "test.md")
	templateContent := "# Test Template"
	templateChecksum := createTestFile(t, templatePath, templateContent)
	
	binaryPath := filepath.Join(installPath, "bin", "the-startup")
	binaryContent := "binary content"
	binaryChecksum := createTestFile(t, binaryPath, binaryContent)
	
	// Add files to lock file
	lockFile.Files["agents/test-agent.md"] = getFileInfo(t, agentPath, agentChecksum)
	lockFile.Files["startup/templates/test.md"] = getFileInfo(t, templatePath, templateChecksum)
	lockFile.Files["bin/the-startup"] = getFileInfo(t, binaryPath, binaryChecksum)
	
	// Remove files
	summary, err := remover.RemoveAllFiles()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Check summary
	if summary.FilesProcessed != 3 {
		t.Errorf("Expected 3 files processed, got %d", summary.FilesProcessed)
	}
	
	if summary.FilesRemoved != 3 {
		t.Errorf("Expected 3 files removed, got %d", summary.FilesRemoved)
	}
	
	if summary.FilesFailed != 0 {
		t.Errorf("Expected 0 files failed, got %d", summary.FilesFailed)
	}
	
	// Verify files are actually removed
	if _, err := os.Stat(agentPath); !os.IsNotExist(err) {
		t.Errorf("Expected agent file to be removed")
	}
	
	if _, err := os.Stat(templatePath); !os.IsNotExist(err) {
		t.Errorf("Expected template file to be removed")
	}
	
	if _, err := os.Stat(binaryPath); !os.IsNotExist(err) {
		t.Errorf("Expected binary file to be removed")
	}
}

func TestRemoveAllFiles_DryRun(t *testing.T) {
	remover, _, claudePath, lockFile := setupTestRemover(t)
	remover.SetDryRun(true)
	
	// Create test file
	agentPath := filepath.Join(claudePath, "agents", "test-agent.md")
	agentContent := "# Test Agent"
	agentChecksum := createTestFile(t, agentPath, agentContent)
	
	// Add file to lock file
	lockFile.Files["agents/test-agent.md"] = getFileInfo(t, agentPath, agentChecksum)
	
	// Remove files (dry run)
	summary, err := remover.RemoveAllFiles()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Check summary
	if summary.FilesRemoved != 1 {
		t.Errorf("Expected 1 file removed (dry run), got %d", summary.FilesRemoved)
	}
	
	// Verify file still exists
	if _, err := os.Stat(agentPath); err != nil {
		t.Errorf("Expected file to still exist in dry run mode: %v", err)
	}
	
	// Check for dry run warning
	found := false
	for _, result := range summary.FileResults {
		if strings.Contains(result.Warning, "dry-run mode") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected dry-run warning in results")
	}
}

func TestRemoveAllFiles_ModifiedFile(t *testing.T) {
	remover, _, claudePath, lockFile := setupTestRemover(t)
	
	// Create test file
	agentPath := filepath.Join(claudePath, "agents", "test-agent.md")
	originalContent := "# Test Agent"
	originalChecksum := createTestFile(t, agentPath, originalContent)
	
	// Add to lock file with original checksum
	lockFile.Files["agents/test-agent.md"] = getFileInfo(t, agentPath, originalChecksum)
	
	// Modify the file
	modifiedContent := "# Modified Test Agent"
	if err := os.WriteFile(agentPath, []byte(modifiedContent), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	
	// Remove files
	summary, err := remover.RemoveAllFiles()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Check that modification was detected
	if summary.FilesModified != 1 {
		t.Errorf("Expected 1 modified file, got %d", summary.FilesModified)
	}
	
	// Check for modification warning
	found := false
	for _, result := range summary.FileResults {
		if result.Modified && strings.Contains(result.Warning, "modified since installation") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected modification warning in results")
	}
	
	// File should still be removed despite being modified
	if _, err := os.Stat(agentPath); !os.IsNotExist(err) {
		t.Error("Expected modified file to be removed")
	}
}

func TestVerifyFiles(t *testing.T) {
	remover, _, claudePath, lockFile := setupTestRemover(t)
	
	// Create test file
	agentPath := filepath.Join(claudePath, "agents", "test-agent.md")
	agentContent := "# Test Agent"
	agentChecksum := createTestFile(t, agentPath, agentContent)
	
	// Add to lock file
	lockFile.Files["agents/test-agent.md"] = getFileInfo(t, agentPath, agentChecksum)
	
	// Verify files (should not remove them)
	summary, err := remover.VerifyFiles()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Check that verification succeeded
	if summary.FilesProcessed != 1 {
		t.Errorf("Expected 1 file processed, got %d", summary.FilesProcessed)
	}
	
	// File should still exist after verification
	if _, err := os.Stat(agentPath); err != nil {
		t.Errorf("Expected file to exist after verification: %v", err)
	}
	
	// Should indicate dry-run mode
	found := false
	for _, result := range summary.FileResults {
		if strings.Contains(result.Warning, "dry-run mode") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected dry-run indication in verify mode")
	}
}

func TestGetRemovalGuidance(t *testing.T) {
	remover, _, _, _ := setupTestRemover(t)
	
	tests := []struct {
		name     string
		summary  *RemovalSummary
		expected []string
	}{
		{
			name: "successful removal",
			summary: &RemovalSummary{
				FilesRemoved:     3,
				FilesProcessed:   3,
				FilesFailed:      0,
				FilesModified:    0,
				DirectoriesFailed: 0,
				Errors:           []string{},
				Warnings:         []string{},
			},
			expected: []string{"✅ All files removed successfully!"},
		},
		{
			name: "failed files",
			summary: &RemovalSummary{
				FilesRemoved:     2,
				FilesProcessed:   3,
				FilesFailed:      1,
				FilesModified:    0,
				DirectoriesFailed: 0,
				Errors:           []string{"agents/test.md: permission denied"},
				Warnings:         []string{},
			},
			expected: []string{
				"❌ 1 files failed to remove. Check file permissions and ensure files are not in use.",
				"",
				"Detailed errors:",
				"  • agents/test.md: permission denied",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			guidance := remover.GetRemovalGuidance(tt.summary)
			
			if len(guidance) != len(tt.expected) {
				t.Errorf("Expected %d guidance lines, got %d", len(tt.expected), len(guidance))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got: %v", guidance)
				return
			}
			
			for i, expected := range tt.expected {
				if guidance[i] != expected {
					t.Errorf("Line %d: expected %q, got %q", i, expected, guidance[i])
				}
			}
		})
	}
}

func TestCalculateChecksum(t *testing.T) {
	remover, _, _, _ := setupTestRemover(t)
	
	// Create test file
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	content := "test content for checksum"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Calculate checksum using remover
	checksum, err := remover.calculateChecksum(tmpFile)
	if err != nil {
		t.Errorf("Expected no error calculating checksum, got: %v", err)
	}
	
	// Calculate expected checksum
	hasher := sha256.New()
	hasher.Write([]byte(content))
	expected := hex.EncodeToString(hasher.Sum(nil))
	
	if checksum != expected {
		t.Errorf("Expected checksum %s, got %s", expected, checksum)
	}
}

// Benchmark tests for performance
func BenchmarkCalculateChecksum(b *testing.B) {
	tmpDir := b.TempDir()
	remover := NewSafeFileRemover(tmpDir, tmpDir, &config.LockFile{})
	
	// Create a 1MB test file
	testFile := filepath.Join(tmpDir, "large_file.bin")
	largeContent := make([]byte, 1024*1024) // 1MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	
	if err := os.WriteFile(testFile, largeContent, 0644); err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := remover.calculateChecksum(testFile)
		if err != nil {
			b.Errorf("Checksum calculation failed: %v", err)
		}
	}
}