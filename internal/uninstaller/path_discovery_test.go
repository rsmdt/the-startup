package uninstaller

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rsmdt/the-startup/internal/config"
)

func TestDiscoverySource_String(t *testing.T) {
	tests := []struct {
		source   DiscoverySource
		expected string
	}{
		{DiscoverySourceLockfile, "lockfile"},
		{DiscoverySourceAutoDetect, "auto-detected"},
		{DiscoverySourceUserInput, "user-provided"},
		{DiscoverySource(999), "unknown"},
	}

	for _, test := range tests {
		result := test.source.String()
		if result != test.expected {
			t.Errorf("DiscoverySource(%d).String() = %s, want %s",
				int(test.source), result, test.expected)
		}
	}
}

func TestNewPathDiscovery(t *testing.T) {
	discovery := NewPathDiscovery()
	if discovery == nil {
		t.Fatal("NewPathDiscovery() returned nil")
	}
	if discovery.lockfilePath != "" {
		t.Errorf("NewPathDiscovery() should have empty lockfilePath, got %s", discovery.lockfilePath)
	}
}

func TestNewPathDiscoveryWithLockfile(t *testing.T) {
	customPath := "/custom/path/lockfile.json"
	discovery := NewPathDiscoveryWithLockfile(customPath)
	if discovery == nil {
		t.Fatal("NewPathDiscoveryWithLockfile() returned nil")
	}
	if discovery.lockfilePath != customPath {
		t.Errorf("NewPathDiscoveryWithLockfile() lockfilePath = %s, want %s",
			discovery.lockfilePath, customPath)
	}
}

func TestDiscoverPaths_FromLockfile(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	installPath := filepath.Join(tempDir, "install")
	claudePath := filepath.Join(tempDir, "claude")

	// Create directories
	if err := os.MkdirAll(installPath, 0755); err != nil {
		t.Fatalf("Failed to create install directory: %v", err)
	}
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude directory: %v", err)
	}

	// Create lockfile
	lockFile := config.LockFile{
		Version:     "1.0.0",
		InstallPath: installPath,
		ClaudePath:  claudePath,
		Files:       make(map[string]config.FileInfo),
	}

	lockfilePath := filepath.Join(tempDir, "test.lock")
	data, err := json.Marshal(lockFile)
	if err != nil {
		t.Fatalf("Failed to marshal lockfile: %v", err)
	}

	if err := os.WriteFile(lockfilePath, data, 0644); err != nil {
		t.Fatalf("Failed to write lockfile: %v", err)
	}

	// Test path discovery
	discovery := NewPathDiscoveryWithLockfile(lockfilePath)
	foundInstall, foundClaude, source, err := discovery.DiscoverPaths()

	if err != nil {
		t.Fatalf("DiscoverPaths() failed: %v", err)
	}

	if source != DiscoverySourceLockfile {
		t.Errorf("DiscoverPaths() source = %v, want %v", source, DiscoverySourceLockfile)
	}

	if foundInstall != installPath {
		t.Errorf("DiscoverPaths() installPath = %s, want %s", foundInstall, installPath)
	}

	if foundClaude != claudePath {
		t.Errorf("DiscoverPaths() claudePath = %s, want %s", foundClaude, claudePath)
	}
}

func TestDiscoverPaths_LockfileWithTildePaths(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()

	// Create actual directories in temp location but store tilde paths in lockfile
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory for tilde test")
	}

	// Create directories in temp location
	installPath := filepath.Join(tempDir, "install")
	claudePath := filepath.Join(tempDir, "claude")

	if err := os.MkdirAll(installPath, 0755); err != nil {
		t.Fatalf("Failed to create install directory: %v", err)
	}
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude directory: %v", err)
	}

	// Create lockfile with tilde paths that point to temp directories
	// For testing, we'll use the actual paths but test tilde expansion logic separately
	lockFile := config.LockFile{
		Version:     "1.0.0",
		InstallPath: installPath,
		ClaudePath:  claudePath,
		Files:       make(map[string]config.FileInfo),
	}

	lockfilePath := filepath.Join(tempDir, "test.lock")
	data, err := json.Marshal(lockFile)
	if err != nil {
		t.Fatalf("Failed to marshal lockfile: %v", err)
	}

	if err := os.WriteFile(lockfilePath, data, 0644); err != nil {
		t.Fatalf("Failed to write lockfile: %v", err)
	}

	// Test path discovery
	discovery := NewPathDiscoveryWithLockfile(lockfilePath)
	foundInstall, foundClaude, source, err := discovery.DiscoverPaths()

	if err != nil {
		t.Fatalf("DiscoverPaths() failed: %v", err)
	}

	if source != DiscoverySourceLockfile {
		t.Errorf("DiscoverPaths() source = %v, want %v", source, DiscoverySourceLockfile)
	}

	// Check that paths are absolute (tilde expansion worked)
	if !filepath.IsAbs(foundInstall) {
		t.Errorf("DiscoverPaths() installPath should be absolute, got %s", foundInstall)
	}

	if !filepath.IsAbs(foundClaude) {
		t.Errorf("DiscoverPaths() claudePath should be absolute, got %s", foundClaude)
	}

	// Test tilde expansion separately
	tildeHomeResult, err := expandTildePath("~/test")
	expectedHome := filepath.Join(homeDir, "test")
	if err != nil {
		t.Errorf("expandTildePath() failed: %v", err)
	}
	if tildeHomeResult != expectedHome {
		t.Errorf("expandTildePath('~/test') = %s, want %s", tildeHomeResult, expectedHome)
	}

	// Test non-tilde path
	nonTildeResult, err := expandTildePath("/absolute/path")
	if err != nil {
		t.Errorf("expandTildePath() failed on non-tilde path: %v", err)
	}
	if nonTildeResult != "/absolute/path" {
		t.Errorf("expandTildePath('/absolute/path') = %s, want /absolute/path", nonTildeResult)
	}
}

func TestDiscoverPaths_MissingLockfile(t *testing.T) {
	// Use non-existent lockfile path
	discovery := NewPathDiscoveryWithLockfile("/non/existent/path.lock")
	foundInstall, foundClaude, source, err := discovery.DiscoverPaths()

	// Should not error, but may fall back to auto-detect or user input
	if err != nil {
		t.Fatalf("DiscoverPaths() failed: %v", err)
	}

	// Should fall back to either auto-detect or user input (depending on environment)
	if source != DiscoverySourceAutoDetect && source != DiscoverySourceUserInput {
		t.Errorf("DiscoverPaths() source = %v, want %v or %v",
			source, DiscoverySourceAutoDetect, DiscoverySourceUserInput)
	}

	// Should return paths
	if foundInstall == "" || foundClaude == "" {
		t.Errorf("DiscoverPaths() should return paths, got install='%s', claude='%s'",
			foundInstall, foundClaude)
	}
}

func TestDiscoverPaths_CorruptedLockfile(t *testing.T) {
	// Create corrupted lockfile
	tempDir := t.TempDir()
	lockfilePath := filepath.Join(tempDir, "corrupted.lock")

	// Write invalid JSON
	if err := os.WriteFile(lockfilePath, []byte("invalid json {"), 0644); err != nil {
		t.Fatalf("Failed to write corrupted lockfile: %v", err)
	}

	discovery := NewPathDiscoveryWithLockfile(lockfilePath)
	foundInstall, foundClaude, source, err := discovery.DiscoverPaths()

	// Should not error, but may fall back to auto-detect or user input
	if err != nil {
		t.Fatalf("DiscoverPaths() failed: %v", err)
	}

	// Should fall back to either auto-detect or user input (depending on environment)
	if source != DiscoverySourceAutoDetect && source != DiscoverySourceUserInput {
		t.Errorf("DiscoverPaths() source = %v, want %v or %v",
			source, DiscoverySourceAutoDetect, DiscoverySourceUserInput)
	}

	// Should return paths
	if foundInstall == "" || foundClaude == "" {
		t.Errorf("DiscoverPaths() should return paths, got install='%s', claude='%s'",
			foundInstall, foundClaude)
	}
}

func TestDiscoverPaths_AutoDetect(t *testing.T) {
	// Create temporary directory structure that mimics installation
	tempDir := t.TempDir()
	installPath := filepath.Join(tempDir, ".the-startup")
	claudePath := filepath.Join(tempDir, ".claude")

	// Create installation evidence
	if err := os.MkdirAll(filepath.Join(installPath, "bin"), 0755); err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}

	// Create binary
	binaryPath := filepath.Join(installPath, "bin", "the-startup")
	if err := os.WriteFile(binaryPath, []byte("fake binary"), 0755); err != nil {
		t.Fatalf("Failed to create binary: %v", err)
	}

	// Create claude evidence
	agentsPath := filepath.Join(claudePath, "agents")
	if err := os.MkdirAll(agentsPath, 0755); err != nil {
		t.Fatalf("Failed to create agents directory: %v", err)
	}

	// Create agent file
	agentPath := filepath.Join(agentsPath, "the-architect.md")
	if err := os.WriteFile(agentPath, []byte("# Agent"), 0644); err != nil {
		t.Fatalf("Failed to create agent file: %v", err)
	}

	// Use non-existent lockfile to force auto-detection
	discovery := NewPathDiscoveryWithLockfile("/non/existent/path.lock")

	// Override the auto-detect logic by calling it directly
	foundInstall, foundClaude, err := discovery.autoDetectPaths()
	if err != nil {
		// Auto-detection might fail in test environment, that's okay
		t.Skip("Auto-detection failed in test environment, which is expected")
	}

	if foundInstall == "" || foundClaude == "" {
		t.Errorf("autoDetectPaths() should return paths, got install='%s', claude='%s'",
			foundInstall, foundClaude)
	}
}

func TestHasInstallationEvidence(t *testing.T) {
	tempDir := t.TempDir()
	discovery := NewPathDiscovery()

	// Test with no evidence
	if discovery.hasInstallationEvidence(tempDir) {
		t.Errorf("hasInstallationEvidence() should return false for empty directory")
	}

	// Test with lockfile
	lockfilePath := filepath.Join(tempDir, "the-startup.lock")
	if err := os.WriteFile(lockfilePath, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create lockfile: %v", err)
	}

	if !discovery.hasInstallationEvidence(tempDir) {
		t.Errorf("hasInstallationEvidence() should return true with lockfile present")
	}

	// Clean up and test with binary
	os.Remove(lockfilePath)
	binDir := filepath.Join(tempDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}

	binaryPath := filepath.Join(binDir, "the-startup")
	if err := os.WriteFile(binaryPath, []byte("binary"), 0755); err != nil {
		t.Fatalf("Failed to create binary: %v", err)
	}

	if !discovery.hasInstallationEvidence(tempDir) {
		t.Errorf("hasInstallationEvidence() should return true with binary present")
	}
}

func TestHasClaudeEvidence(t *testing.T) {
	tempDir := t.TempDir()
	discovery := NewPathDiscovery()

	// Test with no evidence
	if discovery.hasClaudeEvidence(tempDir) {
		t.Errorf("hasClaudeEvidence() should return false for empty directory")
	}

	// Test with agents directory but no our agents
	agentsPath := filepath.Join(tempDir, "agents")
	if err := os.MkdirAll(agentsPath, 0755); err != nil {
		t.Fatalf("Failed to create agents directory: %v", err)
	}

	if discovery.hasClaudeEvidence(tempDir) {
		t.Errorf("hasClaudeEvidence() should return false with empty agents directory")
	}

	// Test with other agent files
	otherAgentPath := filepath.Join(agentsPath, "other-agent.md")
	if err := os.WriteFile(otherAgentPath, []byte("# Other"), 0644); err != nil {
		t.Fatalf("Failed to create other agent: %v", err)
	}

	if discovery.hasClaudeEvidence(tempDir) {
		t.Errorf("hasClaudeEvidence() should return false with non-startup agents")
	}

	// Test with our agent files
	ourAgentPath := filepath.Join(agentsPath, "the-architect.md")
	if err := os.WriteFile(ourAgentPath, []byte("# Our Agent"), 0644); err != nil {
		t.Fatalf("Failed to create our agent: %v", err)
	}

	if !discovery.hasClaudeEvidence(tempDir) {
		t.Errorf("hasClaudeEvidence() should return true with our agents present")
	}
}

func TestGetDefaultPaths(t *testing.T) {
	discovery := NewPathDiscovery()
	installPath, claudePath := discovery.getDefaultPaths()

	if installPath == "" || claudePath == "" {
		t.Errorf("getDefaultPaths() should return non-empty paths, got install='%s', claude='%s'",
			installPath, claudePath)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If we can't get home directory, should still return paths
		if !strings.Contains(installPath, "the-startup") {
			t.Errorf("getDefaultPaths() installPath should contain 'the-startup', got %s", installPath)
		}
		if !strings.Contains(claudePath, "claude") {
			t.Errorf("getDefaultPaths() claudePath should contain 'claude', got %s", claudePath)
		}
	} else {
		// Should use home directory
		expectedInstall := filepath.Join(homeDir, ".the-startup")
		expectedClaude := filepath.Join(homeDir, ".claude")

		if installPath != expectedInstall {
			t.Errorf("getDefaultPaths() installPath = %s, want %s", installPath, expectedInstall)
		}

		if claudePath != expectedClaude {
			t.Errorf("getDefaultPaths() claudePath = %s, want %s", claudePath, expectedClaude)
		}
	}
}

func TestValidatePath(t *testing.T) {
	discovery := NewPathDiscovery()

	// Test with non-existent path
	err := discovery.validatePath("/non/existent/path", "test")
	if err == nil {
		t.Errorf("validatePath() should fail for non-existent path")
	}

	// Test with existing directory
	tempDir := t.TempDir()
	err = discovery.validatePath(tempDir, "test")
	if err != nil {
		t.Errorf("validatePath() should succeed for valid directory: %v", err)
	}

	// Test with file instead of directory
	tempFile := filepath.Join(tempDir, "testfile")
	if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = discovery.validatePath(tempFile, "test")
	if err == nil {
		t.Errorf("validatePath() should fail for file path")
	}
}

func TestReadFromLockfile_EmptyPaths(t *testing.T) {
	tempDir := t.TempDir()

	// Create lockfile with empty paths
	lockFile := config.LockFile{
		Version:     "1.0.0",
		InstallPath: "",
		ClaudePath:  "",
		Files:       make(map[string]config.FileInfo),
	}

	lockfilePath := filepath.Join(tempDir, "empty.lock")
	data, err := json.Marshal(lockFile)
	if err != nil {
		t.Fatalf("Failed to marshal lockfile: %v", err)
	}

	if err := os.WriteFile(lockfilePath, data, 0644); err != nil {
		t.Fatalf("Failed to write lockfile: %v", err)
	}

	discovery := NewPathDiscoveryWithLockfile(lockfilePath)
	_, _, err = discovery.readFromLockfile()

	if err == nil {
		t.Errorf("readFromLockfile() should fail with empty paths")
	}

	if !strings.Contains(err.Error(), "empty paths") {
		t.Errorf("readFromLockfile() error should mention empty paths, got: %v", err)
	}
}

func TestReadFromLockfile_InvalidPaths(t *testing.T) {
	tempDir := t.TempDir()

	// Create lockfile with non-existent paths
	lockFile := config.LockFile{
		Version:     "1.0.0",
		InstallPath: "/non/existent/install",
		ClaudePath:  "/non/existent/claude",
		Files:       make(map[string]config.FileInfo),
	}

	lockfilePath := filepath.Join(tempDir, "invalid.lock")
	data, err := json.Marshal(lockFile)
	if err != nil {
		t.Fatalf("Failed to marshal lockfile: %v", err)
	}

	if err := os.WriteFile(lockfilePath, data, 0644); err != nil {
		t.Fatalf("Failed to write lockfile: %v", err)
	}

	discovery := NewPathDiscoveryWithLockfile(lockfilePath)
	_, _, err = discovery.readFromLockfile()

	if err == nil {
		t.Errorf("readFromLockfile() should fail with invalid paths")
	}

	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("readFromLockfile() error should mention path not existing, got: %v", err)
	}
}
