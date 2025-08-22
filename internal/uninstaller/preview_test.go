package uninstaller

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rsmdt/the-startup/internal/config"
)

// MockPathDiscoverer for testing
type MockPathDiscoverer struct {
	installPath     string
	claudePath      string
	discoverySource DiscoverySource
	shouldError     bool
}

func (m *MockPathDiscoverer) DiscoverPaths() (string, string, DiscoverySource, error) {
	if m.shouldError {
		return "", "", DiscoverySourceUserInput, fmt.Errorf("mock discovery error")
	}
	return m.installPath, m.claudePath, m.discoverySource, nil
}

func TestFileCategory_String(t *testing.T) {
	tests := []struct {
		category FileCategory
		expected string
	}{
		{CategoryAgent, "agent"},
		{CategoryCommand, "command"},
		{CategoryTemplate, "template"},
		{CategoryRule, "rule"},
		{CategoryBinary, "binary"},
		{CategoryLog, "log"},
		{CategorySettings, "settings"},
		{CategoryOther, "other"},
		{CategoryUntracked, "untracked"},
		{FileCategory(999), "unknown"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if got := test.category.String(); got != test.expected {
				t.Errorf("expected %s, got %s", test.expected, got)
			}
		})
	}
}

func TestPreviewGenerator_GeneratePreview_PathDiscoveryError(t *testing.T) {
	mockDiscoverer := &MockPathDiscoverer{
		shouldError: true,
	}

	generator := NewPreviewGenerator(mockDiscoverer)
	preview, err := generator.GeneratePreview()

	if err == nil {
		t.Fatal("expected error from path discovery failure")
	}

	if preview != nil {
		t.Error("expected nil preview when path discovery fails")
	}

	if !containsString(err.Error(), "failed to discover installation paths") {
		t.Errorf("expected path discovery error, got: %s", err.Error())
	}
}

func TestPreviewGenerator_GeneratePreview_NonExistentDirectories(t *testing.T) {
	// Use non-existent directories
	installPath := "/tmp/non-existent-install"
	claudePath := "/tmp/non-existent-claude"

	mockDiscoverer := &MockPathDiscoverer{
		installPath:     installPath,
		claudePath:      claudePath,
		discoverySource: DiscoverySourceAutoDetect,
	}

	generator := NewPreviewGenerator(mockDiscoverer)
	generator.SetVerbose(true)

	preview, err := generator.GeneratePreview()

	// Should not error - non-existent directories are handled gracefully
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if preview == nil {
		t.Fatal("expected non-nil preview")
	}

	// Validate basic structure
	if preview.InstallPath != installPath {
		t.Errorf("expected install path %s, got %s", installPath, preview.InstallPath)
	}

	if preview.ClaudePath != claudePath {
		t.Errorf("expected claude path %s, got %s", claudePath, preview.ClaudePath)
	}

	if preview.DiscoverySource != DiscoverySourceAutoDetect {
		t.Errorf("expected discovery source %v, got %v", DiscoverySourceAutoDetect, preview.DiscoverySource)
	}

	// Should have empty lists
	if len(preview.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(preview.Files))
	}
}

func TestPreviewGenerator_GeneratePreview_WithMockInstallation(t *testing.T) {
	// Create temporary directories for testing
	tempDir := t.TempDir()
	installPath := filepath.Join(tempDir, "install")
	claudePath := filepath.Join(tempDir, "claude")

	// Create directory structure
	if err := os.MkdirAll(installPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatal(err)
	}

	// Create mock files
	testFiles := map[string]struct {
		path        string
		content     string
		category    FileCategory
	}{
		"lockfile": {
			path:     filepath.Join(installPath, "the-startup.lock"),
			content:  createMockLockFile(installPath, claudePath),
			category: CategoryOther,
		},
		"binary": {
			path:     filepath.Join(installPath, "bin", "the-startup"),
			content:  "fake binary content",
			category: CategoryBinary,
		},
		"template": {
			path:     filepath.Join(installPath, "templates", "test.md"),
			content:  "# Test Template",
			category: CategoryTemplate,
		},
		"rule": {
			path:     filepath.Join(installPath, "rules", "test.md"),
			content:  "# Test Rule",
			category: CategoryRule,
		},
		"agent": {
			path:     filepath.Join(claudePath, "agents", "the-test-agent.md"),
			content:  "# Test Agent",
			category: CategoryAgent,
		},
		"command": {
			path:     filepath.Join(claudePath, "commands", "s", "test.md"),
			content:  "# Test Command",
			category: CategoryCommand,
		},
		"settings": {
			path:     filepath.Join(claudePath, "settings.json"),
			content:  `{"test": "settings"}`,
			category: CategorySettings,
		},
	}

	for name, file := range testFiles {
		if err := os.MkdirAll(filepath.Dir(file.path), 0755); err != nil {
			t.Fatalf("failed to create directory for %s: %v", name, err)
		}
		if err := os.WriteFile(file.path, []byte(file.content), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", name, err)
		}
	}

	// Create preview generator
	mockDiscoverer := &MockPathDiscoverer{
		installPath:     installPath,
		claudePath:      claudePath,
		discoverySource: DiscoverySourceLockfile,
	}

	generator := NewPreviewGenerator(mockDiscoverer)
	generator.SetVerbose(true)

	// Generate preview
	preview, err := generator.GeneratePreview()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate results
	if preview == nil {
		t.Fatal("expected non-nil preview")
	}

	// Should have found files
	if len(preview.Files) == 0 {
		t.Error("expected to find some files")
	}

	// Check that lockfile was loaded
	if preview.LockFile == nil {
		t.Error("expected lockfile to be loaded")
	}

	// Validate file categorization
	categoryCount := make(map[FileCategory]int)
	for _, file := range preview.Files {
		categoryCount[file.Category]++
	}

	// Should have found each type of file
	expectedCategories := []FileCategory{
		CategoryOther,   // lockfile
		CategoryBinary,  // the-startup binary
		CategoryTemplate, CategoryRule, CategoryAgent, CategoryCommand, CategorySettings,
	}

	for _, expectedCategory := range expectedCategories {
		if count, found := categoryCount[expectedCategory]; !found || count == 0 {
			t.Errorf("expected to find files of category %s", expectedCategory.String())
		}
	}

	// Validate category summaries
	if len(preview.CategorySummary) == 0 {
		t.Error("expected category summaries to be generated")
	}

	// Check total counts
	if preview.TotalFiles != len(preview.Files) {
		t.Errorf("total files mismatch: expected %d, got %d", len(preview.Files), preview.TotalFiles)
	}

	var expectedTotalSize int64
	for _, file := range preview.Files {
		expectedTotalSize += file.Size
	}

	if preview.TotalSize != expectedTotalSize {
		t.Errorf("total size mismatch: expected %d, got %d", expectedTotalSize, preview.TotalSize)
	}
}

func TestPreviewGenerator_CategorizeFile(t *testing.T) {
	generator := NewPreviewGenerator(nil)

	tests := []struct {
		fullPath    string
		relPath     string
		dirType     string
		expected    FileCategory
	}{
		{"/test/agents/the-test.md", "agents/the-test.md", "claude", CategoryAgent},
		{"/test/commands/s/test.md", "commands/s/test.md", "claude", CategoryCommand},
		{"/test/templates/test.md", "templates/test.md", "install", CategoryTemplate},
		{"/test/rules/test.md", "rules/test.md", "install", CategoryRule},
		{"/test/bin/the-startup", "bin/the-startup", "install", CategoryBinary},
		{"/test/logs/session.jsonl", "logs/session.jsonl", "install", CategoryLog},
		{"/test/settings.json", "settings.json", "claude", CategorySettings},
		{"/test/settings.local.json", "settings.local.json", "claude", CategorySettings},
		{"/test/the-startup.lock", "the-startup.lock", "install", CategoryOther},
		{"/test/unknown.txt", "unknown.txt", "install", CategoryOther},
	}

	for _, test := range tests {
		t.Run(test.relPath, func(t *testing.T) {
			result := generator.categorizeFile(test.fullPath, test.relPath, test.dirType)
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected.String(), result.String())
			}
		})
	}
}

func TestPreviewGenerator_IsStartupRelated(t *testing.T) {
	generator := NewPreviewGenerator(nil)

	tests := []struct {
		name     string
		fullPath string
		relPath  string
		dirType  string
		fileName string
		expected bool
	}{
		{"lockfile", "/test/the-startup.lock", "the-startup.lock", "install", "the-startup.lock", true},
		{"agent", "/test/agents/the-test.md", "agents/the-test.md", "claude", "the-test.md", true},
		{"command", "/test/commands/s/test.md", "commands/s/test.md", "claude", "test.md", true},
		{"settings", "/test/settings.json", "settings.json", "claude", "settings.json", true},
		{"binary", "/test/bin/the-startup", "bin/the-startup", "install", "the-startup", true},
		{"template", "/test/templates/test.md", "templates/test.md", "install", "test.md", true},
		{"log", "/test/logs/session.jsonl", "logs/session.jsonl", "install", "session.jsonl", true},
		{"unrelated", "/test/random.txt", "random.txt", "install", "random.txt", false},
		{"claude_unrelated", "/test/random.md", "random.md", "claude", "random.md", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a mock DirEntry
			mockDirEntry := &mockDirEntry{name: test.fileName}
			
			result := generator.isStartupRelated(test.fullPath, test.relPath, test.dirType, mockDirEntry)
			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestPreviewGenerator_GenerateLockfilePaths(t *testing.T) {
	generator := NewPreviewGenerator(nil)

	preview := &RemovalPreview{
		InstallPath: "/test/install",
		ClaudePath:  "/test/claude",
	}

	tests := []struct {
		name     string
		file     FileInfo
		expected []string
	}{
		{
			name: "install_template",
			file: FileInfo{
				Path:         "/test/install/templates/test.md",
				RelativePath: "templates/test.md",
			},
			expected: []string{"startup/templates/test.md", "templates/test.md"},
		},
		{
			name: "install_binary",
			file: FileInfo{
				Path:         "/test/install/bin/the-startup",
				RelativePath: "bin/the-startup",
			},
			expected: []string{"bin/the-startup", "bin/the-startup"},
		},
		{
			name: "claude_agent",
			file: FileInfo{
				Path:         "/test/claude/agents/the-test.md",
				RelativePath: "agents/the-test.md",
			},
			expected: []string{"agents/the-test.md", "agents/the-test.md"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := generator.generateLockfilePaths(&test.file, preview)
			
			if len(result) != len(test.expected) {
				t.Errorf("expected %d paths, got %d", len(test.expected), len(result))
				return
			}

			for i, expected := range test.expected {
				if result[i] != expected {
					t.Errorf("path %d: expected %s, got %s", i, expected, result[i])
				}
			}
		})
	}
}

func TestPreviewGenerator_SecurityAnalysis(t *testing.T) {
	generator := NewPreviewGenerator(nil)

	preview := &RemovalPreview{
		InstallPath: "/safe/install",
		ClaudePath:  "/safe/claude",
		Files: []FileInfo{
			{Path: "/safe/install/normal.txt", Size: 1024},
			{Path: "/safe/install/../evil.txt", Size: 1024}, // Path traversal
			{Path: "/unsafe/location/file.txt", Size: 1024}, // Outside safe dirs
			{Path: "/safe/install/huge.bin", Size: 200 * 1024 * 1024}, // Large file
			{
				Path:          "/safe/install/link.txt",
				Size:          1024,
				IsSymlink:     true,
				SymlinkTarget: "/unsafe/target.txt", // Unsafe symlink
			},
		},
	}

	generator.performSecurityAnalysis(preview)

	// Should have identified security issues
	if len(preview.SecurityIssues) == 0 {
		t.Error("expected security issues to be identified")
	}

	// Check for specific issue types
	issueTypes := make(map[string]bool)
	for _, issue := range preview.SecurityIssues {
		issueTypes[issue.Type] = true
	}

	expectedIssueTypes := []string{"path_traversal", "unexpected_location", "large_file", "suspicious_symlink"}
	for _, expectedType := range expectedIssueTypes {
		if !issueTypes[expectedType] {
			t.Errorf("expected to find security issue type: %s", expectedType)
		}
	}
}

func TestPreviewGenerator_ValidationErrors(t *testing.T) {
	// Create temporary file to test validation
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.bin")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	generator := NewPreviewGenerator(nil)

	preview := &RemovalPreview{
		InstallPath: tempDir,
		ClaudePath:  tempDir,
		Files: []FileInfo{
			{
				Path:            testFile,
				Category:        CategoryBinary,
				PermissionIssue: true,
				PermissionError: "test permission error",
			},
		},
	}

	generator.validateRemovalFeasibility(preview)

	// Should have identified validation errors
	if len(preview.ValidationErrors) == 0 {
		t.Error("expected validation errors to be identified")
	}

	// Check for permission error
	found := false
	for _, err := range preview.ValidationErrors {
		if err.Type == "permission_denied" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected to find permission_denied validation error")
	}
}

// Helper functions and mocks

type mockDirEntry struct {
	name string
	size int64
	mode os.FileMode
}

func (m *mockDirEntry) Name() string               { return m.name }
func (m *mockDirEntry) IsDir() bool                { return m.mode.IsDir() }
func (m *mockDirEntry) Type() os.FileMode          { return m.mode.Type() }
func (m *mockDirEntry) Info() (os.FileInfo, error) {
	return &mockFileInfo{
		name: m.name,
		size: m.size,
		mode: m.mode,
	}, nil
}

type mockFileInfo struct {
	name string
	size int64
	mode os.FileMode
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m *mockFileInfo) ModTime() time.Time { return time.Now() }
func (m *mockFileInfo) IsDir() bool        { return m.mode.IsDir() }
func (m *mockFileInfo) Sys() interface{}   { return nil }

func createMockLockFile(installPath, claudePath string) string {
	lockFile := config.LockFile{
		Version:     "1.0.0",
		InstallDate: time.Now().Format(time.RFC3339),
		InstallPath: installPath,
		ClaudePath:  claudePath,
		Tool:        "claude-code",
		Components:  []string{"agents", "commands", "templates", "rules"},
		Files: map[string]config.FileInfo{
			"agents/the-test-agent.md": {
				Size:         100,
				LastModified: time.Now().Add(-time.Hour).Format(time.RFC3339),
			},
			"commands/s/test.md": {
				Size:         200,
				LastModified: time.Now().Add(-time.Hour).Format(time.RFC3339),
			},
			"startup/templates/test.md": {
				Size:         150,
				LastModified: time.Now().Add(-time.Hour).Format(time.RFC3339),
			},
			"startup/rules/test.md": {
				Size:         180,
				LastModified: time.Now().Add(-time.Hour).Format(time.RFC3339),
			},
			"bin/the-startup": {
				Size:         5000,
				LastModified: time.Now().Add(-time.Hour).Format(time.RFC3339),
			},
		},
	}

	data, _ := json.MarshalIndent(lockFile, "", "  ")
	return string(data)
}

func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Benchmark tests
func BenchmarkPreviewGenerator_CategorizeFile(b *testing.B) {
	generator := NewPreviewGenerator(nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.categorizeFile("/test/agents/the-test.md", "agents/the-test.md", "claude")
	}
}

func BenchmarkPreviewGenerator_IsStartupRelated(b *testing.B) {
	generator := NewPreviewGenerator(nil)
	mockDirEntry := &mockDirEntry{name: "the-test.md"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.isStartupRelated("/test/agents/the-test.md", "agents/the-test.md", "claude", mockDirEntry)
	}
}