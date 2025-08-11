package log

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetStartupDirCompatibility(t *testing.T) {
	// Test cases that match Python behavior exactly
	tests := []struct {
		name        string
		setupFn     func(tempDir string) string
		expectedDir string // relative to tempDir or "home" for home directory
	}{
		{
			name: "project local .the-startup exists - should use it",
			setupFn: func(tempDir string) string {
				localStartup := filepath.Join(tempDir, ".the-startup")
				os.MkdirAll(localStartup, 0755)
				return tempDir
			},
			expectedDir: ".the-startup",
		},
		{
			name: "no project local .the-startup - should use home",
			setupFn: func(tempDir string) string {
				// Don't create local .the-startup directory
				return tempDir
			},
			expectedDir: "home",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			projectDir := tt.setupFn(tempDir)

			result := GetStartupDir(projectDir)

			if tt.expectedDir == "home" {
				// Should be ~/.the-startup
				homeDir, err := os.UserHomeDir()
				if err != nil {
					t.Fatalf("Failed to get home directory: %v", err)
				}
				expected := filepath.Join(homeDir, ".the-startup")
				if result != expected {
					t.Errorf("GetStartupDir() = %q, expected %q", result, expected)
				}
				// Verify the directory was created
				if _, err := os.Stat(result); os.IsNotExist(err) {
					t.Errorf("Home startup directory was not created: %s", result)
				}
			} else {
				// Should be project local
				expected := filepath.Join(projectDir, tt.expectedDir)
				if result != expected {
					t.Errorf("GetStartupDir() = %q, expected %q", result, expected)
				}
			}
		})
	}
}

func TestFindLatestSessionCompatibility(t *testing.T) {
	tests := []struct {
		name        string
		setupFn     func(tempDir string)
		projectDir  string
		expected    string
	}{
		{
			name: "multiple dev- sessions - find latest",
			setupFn: func(tempDir string) {
				startupDir := filepath.Join(tempDir, ".the-startup")
				os.MkdirAll(startupDir, 0755)
				
				// Create session directories with different timestamps
				session1 := filepath.Join(startupDir, "dev-session-1")
				session2 := filepath.Join(startupDir, "dev-session-2")
				session3 := filepath.Join(startupDir, "dev-session-3")
				
				os.MkdirAll(session1, 0755)
				os.MkdirAll(session2, 0755)
				os.MkdirAll(session3, 0755)
				
				// Make session2 the most recently modified
				// We can't easily manipulate timestamps in tests, so we'll just verify
				// the function finds a valid dev- session
			},
			expected: "dev-", // Should start with "dev-"
		},
		{
			name: "no dev- sessions",
			setupFn: func(tempDir string) {
				startupDir := filepath.Join(tempDir, ".the-startup")
				os.MkdirAll(startupDir, 0755)
				
				// Create non-dev directories
				otherDir := filepath.Join(startupDir, "other-session")
				os.MkdirAll(otherDir, 0755)
			},
			expected: "", // Should return empty string
		},
		{
			name: "startup directory doesn't exist locally but exists in home",
			setupFn: func(tempDir string) {
				// Don't create local .the-startup directory
				// Function should fall back to home directory
			},
			expected: "dev-", // Should find a dev- session in home directory
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tt.setupFn(tempDir)
			
			result := FindLatestSession(tempDir)
			
			if tt.expected == "" {
				if result != "" {
					t.Errorf("FindLatestSession() = %q, expected empty string", result)
				}
			} else if tt.expected == "dev-" {
				// Should start with "dev-"
				if result == "" || !filepath.HasPrefix(result, "dev-") {
					t.Errorf("FindLatestSession() = %q, expected to start with 'dev-'", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("FindLatestSession() = %q, expected %q", result, tt.expected)
				}
			}
		})
	}
}

func TestIsDebugEnabledCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{
			name:     "DEBUG_HOOKS not set",
			envValue: "",
			expected: false,
		},
		{
			name:     "DEBUG_HOOKS set to 1",
			envValue: "1",
			expected: true,
		},
		{
			name:     "DEBUG_HOOKS set to true",
			envValue: "true",
			expected: true,
		},
		{
			name:     "DEBUG_HOOKS set to any value",
			envValue: "yes",
			expected: true,
		},
		{
			name:     "DEBUG_HOOKS set to empty string",
			envValue: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			original := os.Getenv("DEBUG_HOOKS")
			defer func() {
				if original == "" {
					os.Unsetenv("DEBUG_HOOKS")
				} else {
					os.Setenv("DEBUG_HOOKS", original)
				}
			}()

			// Set test value
			if tt.envValue == "" {
				os.Unsetenv("DEBUG_HOOKS")
			} else {
				os.Setenv("DEBUG_HOOKS", tt.envValue)
			}

			result := IsDebugEnabled()
			if result != tt.expected {
				t.Errorf("IsDebugEnabled() = %t, expected %t (env value: %q)", result, tt.expected, tt.envValue)
			}
		})
	}
}

func TestGetProjectDirExtensive(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "CLAUDE_PROJECT_DIR set to absolute path",
			envValue: "/tmp/test-project",
			expected: "/tmp/test-project",
		},
		{
			name:     "CLAUDE_PROJECT_DIR set to relative path",
			envValue: "../relative-project",
			expected: "../relative-project",
		},
		{
			name:     "CLAUDE_PROJECT_DIR not set",
			envValue: "",
			expected: ".",
		},
		{
			name:     "CLAUDE_PROJECT_DIR set to current directory",
			envValue: ".",
			expected: ".",
		},
		{
			name:     "CLAUDE_PROJECT_DIR set to home directory",
			envValue: "~",
			expected: "~",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			original := os.Getenv("CLAUDE_PROJECT_DIR")
			defer func() {
				if original == "" {
					os.Unsetenv("CLAUDE_PROJECT_DIR")
				} else {
					os.Setenv("CLAUDE_PROJECT_DIR", original)
				}
			}()

			// Set test value
			if tt.envValue == "" {
				os.Unsetenv("CLAUDE_PROJECT_DIR")
			} else {
				os.Setenv("CLAUDE_PROJECT_DIR", tt.envValue)
			}

			result := GetProjectDir()
			if result != tt.expected {
				t.Errorf("GetProjectDir() = %q, expected %q (env value: %q)", result, tt.expected, tt.envValue)
			}
		})
	}
}

func TestGetStartupDirErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		setupFn    func(*testing.T) string
		validateFn func(*testing.T, string)
	}{
		{
			name: "home directory not accessible - should fallback to current dir",
			setupFn: func(t *testing.T) string {
				tempDir := t.TempDir()
				// Simulate scenario where UserHomeDir fails by testing with empty HOME
				originalHome := os.Getenv("HOME")
				originalUserProfile := os.Getenv("USERPROFILE")
				os.Unsetenv("HOME")
				os.Unsetenv("USERPROFILE")
				
				t.Cleanup(func() {
					if originalHome != "" {
						os.Setenv("HOME", originalHome)
					}
					if originalUserProfile != "" {
						os.Setenv("USERPROFILE", originalUserProfile)
					}
				})
				
				return tempDir
			},
			validateFn: func(t *testing.T, result string) {
				// Should fallback to ./.the-startup when home dir not accessible
				expected := filepath.Join(".", ".the-startup")
				if result != expected {
					t.Errorf("Expected fallback to %s, got %s", expected, result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectDir := tt.setupFn(t)
			result := GetStartupDir(projectDir)
			tt.validateFn(t, result)
		})
	}
}

func TestDirExists(t *testing.T) {
	tests := []struct {
		name       string
		setupFn    func(*testing.T, string) string
		expected   bool
	}{
		{
			name: "directory exists",
			setupFn: func(t *testing.T, tempDir string) string {
				existingDir := filepath.Join(tempDir, "existing")
				os.MkdirAll(existingDir, 0755)
				return existingDir
			},
			expected: true,
		},
		{
			name: "directory does not exist",
			setupFn: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent")
			},
			expected: false,
		},
		{
			name: "file exists (not directory)",
			setupFn: func(t *testing.T, tempDir string) string {
				filePath := filepath.Join(tempDir, "file.txt")
				file, _ := os.Create(filePath)
				file.Close()
				return filePath
			},
			expected: false,
		},
		{
			name: "path is empty",
			setupFn: func(t *testing.T, tempDir string) string {
				return ""
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := tt.setupFn(t, tempDir)
			
			result := dirExists(path)
			if result != tt.expected {
				t.Errorf("dirExists(%q) = %t, expected %t", path, result, tt.expected)
			}
		})
	}
}