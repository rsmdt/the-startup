package stats

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSanitizeProjectPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic path with leading slash",
			input:    "/Users/irudi/Code/personal/the-startup",
			expected: "-Users-irudi-Code-personal-the-startup",
		},
		{
			name:     "path without leading slash",
			input:    "Users/irudi/Code/personal/the-startup",
			expected: "Users-irudi-Code-personal-the-startup",
		},
		{
			name:     "root path",
			input:    "/",
			expected: "-",
		},
		{
			name:     "path with trailing slash",
			input:    "/Users/irudi/Code/personal/the-startup/",
			expected: "-Users-irudi-Code-personal-the-startup-",
		},
		{
			name:     "deep nested path",
			input:    "/Users/john/Documents/Projects/2024/january/my-project",
			expected: "-Users-john-Documents-Projects-2024-january-my-project",
		},
		{
			name:     "path with spaces",
			input:    "/Users/john/My Documents/Project Name",
			expected: "-Users-john-My Documents-Project Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeProjectPath(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeProjectPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractTimestampFromPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		filename     string
		shouldParse  bool
		expectedYear int
		expectedMonth time.Month
		expectedDay  int
	}{
		{
			name:         "date with session suffix",
			filename:     "2025-01-05-session.jsonl",
			shouldParse:  true,
			expectedYear: 2025,
			expectedMonth: time.January,
			expectedDay:  5,
		},
		{
			name:         "full timestamp",
			filename:     "2025-01-05-14-30-45.jsonl",
			shouldParse:  true,
			expectedYear: 2025,
			expectedMonth: time.January,
			expectedDay:  5,
		},
		{
			name:         "date only",
			filename:     "2025-01-05.jsonl",
			shouldParse:  true,
			expectedYear: 2025,
			expectedMonth: time.January,
			expectedDay:  5,
		},
		{
			name:         "invalid format",
			filename:     "session-abc.jsonl",
			shouldParse:  false,
		},
		{
			name:         "partial date",
			filename:     "2025-01.jsonl",
			shouldParse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy file for fallback testing
			filePath := filepath.Join(tempDir, tt.filename)
			if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			result := extractTimestampFromPath(filePath)
			
			if tt.shouldParse {
				if result == nil {
					t.Errorf("extractTimestampFromPath(%q) returned nil, expected timestamp", filePath)
					return
				}
				
				if result.Year() != tt.expectedYear {
					t.Errorf("Year = %d, want %d", result.Year(), tt.expectedYear)
				}
				if result.Month() != tt.expectedMonth {
					t.Errorf("Month = %v, want %v", result.Month(), tt.expectedMonth)
				}
				if result.Day() != tt.expectedDay {
					t.Errorf("Day = %d, want %d", result.Day(), tt.expectedDay)
				}
			} else {
				// For files that shouldn't parse, we still get modification time as fallback
				if result == nil {
					t.Errorf("extractTimestampFromPath(%q) returned nil, expected at least mod time", filePath)
				}
			}
		})
	}
}

func TestLogDiscovery_GetCurrentProject(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	
	// Create a test project directory structure
	actualTestPath := filepath.Join(tempDir, "Users", "test", "project")
	if err := os.MkdirAll(actualTestPath, 0755); err != nil {
		t.Fatalf("Failed to create test project directory: %v", err)
	}
	
	// Save current directory and restore it after test
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	
	// Change to the test directory first to get the actual absolute path
	if err := os.Chdir(actualTestPath); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	// Get the actual current directory (which may have /private prefix on macOS)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	
	// Get the sanitized project name based on the actual path
	sanitizedProjectName := sanitizeProjectPath(cwd)
	
	// Create a mock Claude projects directory with the correctly sanitized name
	claudeDir := filepath.Join(tempDir, ".claude", "projects")
	projectDir := filepath.Join(claudeDir, sanitizedProjectName)
	
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}

	// Create a dummy log file
	logFile := filepath.Join(projectDir, "2025-01-05-session.jsonl")
	if err := os.WriteFile(logFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	discovery := &LogDiscovery{
		HomeDir: tempDir,
	}

	// Test case 1: Current directory matches a Claude project
	result := discovery.GetCurrentProject()
	if result != sanitizedProjectName {
		t.Errorf("GetCurrentProject() = %q, want %q", result, sanitizedProjectName)
	}

	// Test case 2: Current directory doesn't match any Claude project
	nonExistentPath := filepath.Join(tempDir, "nonexistent")
	if err := os.MkdirAll(nonExistentPath, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.Chdir(nonExistentPath); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	result = discovery.GetCurrentProject()
	if result != "" {
		t.Errorf("GetCurrentProject() = %q, want empty string for non-existent project", result)
	}
}

func TestLogDiscovery_ValidateProjectPath(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	
	// Create mock Claude projects
	claudeDir := filepath.Join(tempDir, ".claude", "projects")
	
	// Project with logs
	projectWithLogs := "project-with-logs"
	projectWithLogsDir := filepath.Join(claudeDir, projectWithLogs)
	if err := os.MkdirAll(projectWithLogsDir, 0755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}
	logFile := filepath.Join(projectWithLogsDir, "2025-01-05-session.jsonl")
	if err := os.WriteFile(logFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	// Project without logs
	projectWithoutLogs := "project-without-logs"
	projectWithoutLogsDir := filepath.Join(claudeDir, projectWithoutLogs)
	if err := os.MkdirAll(projectWithoutLogsDir, 0755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}

	discovery := &LogDiscovery{
		HomeDir: tempDir,
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "sanitized project with logs",
			path:     projectWithLogs,
			expected: true,
		},
		{
			name:     "sanitized project without logs",
			path:     projectWithoutLogs,
			expected: false,
		},
		{
			name:     "full path that needs sanitization",
			path:     "/project/with/logs",
			expected: true,
		},
		{
			name:     "non-existent project",
			path:     "non-existent-project",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := discovery.ValidateProjectPath(tt.path)
			if result != tt.expected {
				t.Errorf("ValidateProjectPath(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestLogDiscovery_FindLogFiles(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	
	// Create mock Claude projects directory
	claudeDir := filepath.Join(tempDir, ".claude", "projects")
	projectName := "test-project"
	projectDir := filepath.Join(claudeDir, projectName)
	
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}

	// Create test log files with different timestamps
	testFiles := []string{
		"2025-01-03-session.jsonl",
		"2025-01-05-morning.jsonl",
		"2025-01-05-afternoon.jsonl",
		"2025-01-07-session.jsonl",
		"2025-01-10-session.jsonl",
		"not-a-log.txt", // Should be ignored
	}

	for _, filename := range testFiles {
		filePath := filepath.Join(projectDir, filename)
		if err := os.WriteFile(filePath, []byte("{}"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	discovery := &LogDiscovery{
		HomeDir: tempDir,
	}

	t.Run("find all files", func(t *testing.T) {
		files, err := discovery.FindLogFiles(projectName, FilterOptions{})
		if err != nil {
			t.Fatalf("FindLogFiles failed: %v", err)
		}

		// Should find 5 JSONL files (excluding the .txt file)
		if len(files) != 5 {
			t.Errorf("Found %d files, want 5", len(files))
		}

		// Check that files are sorted
		for i := 1; i < len(files); i++ {
			if files[i] < files[i-1] {
				t.Errorf("Files not sorted: %s comes after %s", files[i], files[i-1])
			}
		}
	})

	t.Run("filter by start time", func(t *testing.T) {
		// Filter for files from Jan 5, 2025 onwards
		startTime := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
		files, err := discovery.FindLogFiles(projectName, FilterOptions{
			StartTime: &startTime,
		})
		if err != nil {
			t.Fatalf("FindLogFiles failed: %v", err)
		}

		// Should find 4 files (Jan 5, 7, and 10)
		if len(files) != 4 {
			t.Errorf("Found %d files, want 4", len(files))
			for _, f := range files {
				t.Logf("File: %s", f)
			}
		}
	})

	t.Run("filter by end time", func(t *testing.T) {
		// Filter for files up to Jan 7, 2025
		endTime := time.Date(2025, 1, 7, 23, 59, 59, 0, time.UTC)
		files, err := discovery.FindLogFiles(projectName, FilterOptions{
			EndTime: &endTime,
		})
		if err != nil {
			t.Fatalf("FindLogFiles failed: %v", err)
		}

		// Should find 4 files (Jan 3, 5, and 7)
		if len(files) != 4 {
			t.Errorf("Found %d files, want 4", len(files))
		}
	})

	t.Run("filter by time range", func(t *testing.T) {
		// Filter for files between Jan 5 and Jan 7
		startTime := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2025, 1, 7, 23, 59, 59, 0, time.UTC)
		files, err := discovery.FindLogFiles(projectName, FilterOptions{
			StartTime: &startTime,
			EndTime:   &endTime,
		})
		if err != nil {
			t.Fatalf("FindLogFiles failed: %v", err)
		}

		// Should find 3 files (Jan 5 and 7)
		if len(files) != 3 {
			t.Errorf("Found %d files, want 3", len(files))
		}
	})

	t.Run("non-existent project", func(t *testing.T) {
		_, err := discovery.FindLogFiles("non-existent", FilterOptions{})
		if err == nil {
			t.Error("Expected error for non-existent project, got nil")
		}
	})

	t.Run("empty project path uses current", func(t *testing.T) {
		// Save and restore current directory
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)

		// Create a test directory
		testPath := filepath.Join(tempDir, "test", "project2")
		if err := os.MkdirAll(testPath, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
		
		// Change to the directory first to get the actual absolute path
		if err := os.Chdir(testPath); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		
		// Get the actual current directory path
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get current directory: %v", err)
		}
		
		// Get the sanitized name and create a Claude project for it
		sanitizedName := sanitizeProjectPath(cwd)
		claudeProjectDir := filepath.Join(claudeDir, sanitizedName)
		if err := os.MkdirAll(claudeProjectDir, 0755); err != nil {
			t.Fatalf("Failed to create Claude project directory: %v", err)
		}
		
		// Create test files in the new project directory
		for _, filename := range []string{
			"2025-01-03-session.jsonl",
			"2025-01-05-morning.jsonl",
			"2025-01-05-afternoon.jsonl",
			"2025-01-07-session.jsonl",
			"2025-01-10-session.jsonl",
		} {
			dst := filepath.Join(claudeProjectDir, filename)
			if err := os.WriteFile(dst, []byte("{}"), 0644); err != nil {
				t.Fatalf("Failed to create test file %s: %v", filename, err)
			}
		}

		files, err := discovery.FindLogFiles("", FilterOptions{})
		if err != nil {
			t.Fatalf("FindLogFiles with empty path failed: %v", err)
		}

		// Should find all 5 JSONL files
		if len(files) != 5 {
			t.Errorf("Found %d files, want 5", len(files))
		}
	})
}

func TestLogDiscovery_Integration(t *testing.T) {
	// This test simulates a real-world scenario with the actual directory structure
	// Skip if we can't access the real Claude directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot determine user home directory")
	}

	claudeProjectsDir := filepath.Join(homeDir, ".claude", "projects")
	if _, err := os.Stat(claudeProjectsDir); os.IsNotExist(err) {
		t.Skip("Claude projects directory does not exist")
	}

	discovery := NewLogDiscovery()

	t.Run("real project detection", func(t *testing.T) {
		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			t.Skip("Cannot get current working directory")
		}

		// Check if current project exists in Claude
		projectPath := discovery.GetCurrentProject()
		if projectPath != "" {
			t.Logf("Detected current project: %s", projectPath)
			
			// Validate the detected project
			if !discovery.ValidateProjectPath(projectPath) {
				t.Error("Detected project does not validate")
			}

			// Try to find log files
			files, err := discovery.FindLogFiles(projectPath, FilterOptions{})
			if err != nil {
				t.Logf("Could not find log files: %v", err)
			} else {
				t.Logf("Found %d log files in current project", len(files))
			}
		} else {
			t.Logf("Current directory %s is not a Claude project", cwd)
		}
	})
}

// Benchmark tests
func BenchmarkSanitizeProjectPath(b *testing.B) {
	path := "/Users/irudi/Code/personal/the-startup/with/very/deep/nesting/structure"
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = sanitizeProjectPath(path)
	}
}

func BenchmarkExtractTimestampFromPath(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "2025-01-05-14-30-45-session.jsonl")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = extractTimestampFromPath(testFile)
	}
}

func BenchmarkFindLogFiles(b *testing.B) {
	// Create a temporary directory with many log files
	tempDir := b.TempDir()
	projectDir := filepath.Join(tempDir, ".claude", "projects", "test-project")
	os.MkdirAll(projectDir, 0755)
	
	// Create 100 test files
	for i := 0; i < 100; i++ {
		filename := fmt.Sprintf("2025-01-%02d-session.jsonl", i%28+1)
		filePath := filepath.Join(projectDir, filename)
		os.WriteFile(filePath, []byte("{}"), 0644)
	}
	
	discovery := &LogDiscovery{HomeDir: tempDir}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, _ = discovery.FindLogFiles("test-project", FilterOptions{})
	}
}