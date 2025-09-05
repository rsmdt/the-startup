package log

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAppendMetrics(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "startup-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test installation directory
	testStartupDir := filepath.Join(tempDir, ".the-startup")
	if err := os.MkdirAll(testStartupDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory for testing
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	// Create test metrics entry
	testTime := time.Date(2025, 9, 3, 10, 30, 0, 0, time.UTC)
	success := true
	entry := &MetricsEntry{
		ToolName:   "TestTool",
		SessionID:  "test-session-123",
		HookEvent:  "PreToolUse",
		Timestamp:  testTime,
		ToolInput:  json.RawMessage(`{"param":"value"}`),
		ToolID:     "tool-456",
		Success:    &success,
	}

	// Test appending metrics
	err = AppendMetrics(entry)
	if err != nil {
		t.Errorf("AppendMetrics() returned error: %v", err)
	}

	// Verify file was created
	expectedFile := filepath.Join(testStartupDir, "logs", "20250903.jsonl")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s was not created", expectedFile)
	}

	// Read and verify content
	data, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	// Unmarshal to verify JSON structure
	var readEntry MetricsEntry
	if err := json.Unmarshal(data, &readEntry); err != nil {
		t.Errorf("Failed to unmarshal written data: %v", err)
	}

	// Verify fields
	if readEntry.ToolName != entry.ToolName {
		t.Errorf("ToolName mismatch: got %s, want %s", readEntry.ToolName, entry.ToolName)
	}
	if readEntry.SessionID != entry.SessionID {
		t.Errorf("SessionID mismatch: got %s, want %s", readEntry.SessionID, entry.SessionID)
	}
	if readEntry.ToolID != entry.ToolID {
		t.Errorf("ToolID mismatch: got %s, want %s", readEntry.ToolID, entry.ToolID)
	}
}

func TestAppendMetrics_NilEntry(t *testing.T) {
	// Test that nil entry returns nil (no error)
	err := AppendMetrics(nil)
	if err != nil {
		t.Errorf("AppendMetrics(nil) returned error: %v, want nil", err)
	}
}

func TestAppendMetrics_MultipleEntries(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "startup-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test installation directory
	testStartupDir := filepath.Join(tempDir, ".the-startup")
	if err := os.MkdirAll(testStartupDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory for testing
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	testTime := time.Date(2025, 9, 3, 10, 30, 0, 0, time.UTC)

	// Append multiple entries
	for i := 0; i < 3; i++ {
		entry := &MetricsEntry{
			ToolName:  "TestTool",
			SessionID: "session-" + string(rune('0'+i)),
			HookEvent: "PreToolUse",
			Timestamp: testTime,
			ToolID:    "tool-" + string(rune('0'+i)),
		}
		if err := AppendMetrics(entry); err != nil {
			t.Errorf("AppendMetrics() entry %d returned error: %v", i, err)
		}
	}

	// Read file and count lines
	expectedFile := filepath.Join(testStartupDir, "logs", "20250903.jsonl")
	data, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	lines := 0
	for _, b := range data {
		if b == '\n' {
			lines++
		}
	}

	if lines != 3 {
		t.Errorf("Expected 3 lines in JSONL file, got %d", lines)
	}
}

func TestGetStartupPath(t *testing.T) {
	// Create a completely isolated temporary directory
	tempDir, err := os.MkdirTemp("", "startup-isolated-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a subdirectory to ensure we're not in any existing project
	testDir := filepath.Join(tempDir, "test-project")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Change to test directory
	originalDir, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalDir)

	// Override HOME to prevent finding user installations
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", filepath.Join(tempDir, "fake-home"))
	defer os.Setenv("HOME", originalHome)

	// Test 1: No existing installation, should create default
	path, err := getStartupPath()
	if err != nil {
		t.Errorf("getStartupPath() failed: %v", err)
	}
	if path != ".the-startup" {
		t.Errorf("Expected default path .the-startup, got %s", path)
	}

	// Verify directory was created
	if _, err := os.Stat(".the-startup"); os.IsNotExist(err) {
		t.Errorf("Default directory was not created")
	}

	// Test 2: Directory exists, should return it
	path, err = getStartupPath()
	if err != nil {
		t.Errorf("getStartupPath() with existing dir failed: %v", err)
	}
	if path != ".the-startup" {
		t.Errorf("Expected .the-startup with existing dir, got %s", path)
	}
}

func TestExpandHome(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Could not get home directory")
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde path",
			input:    "~/test/path",
			expected: filepath.Join(homeDir, "test/path"),
		},
		{
			name:     "absolute path",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "relative path",
			input:    "relative/path",
			expected: "relative/path",
		},
		{
			name:     "just tilde",
			input:    "~",
			expected: "~",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandHome(tt.input)
			if result != tt.expected {
				t.Errorf("expandHome(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}