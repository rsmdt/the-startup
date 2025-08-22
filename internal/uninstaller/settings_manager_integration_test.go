package uninstaller

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestSettingsManagerWithInstallerPatterns tests the settings manager with settings files
// that match the patterns created by the installer
func TestSettingsManagerWithInstallerPatterns(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "settings_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	claudePath := filepath.Join(tempDir, "claude")
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude dir: %v", err)
	}

	// Create settings files that match installer patterns
	settingsContent := createInstallerStyleSettings()
	localSettingsContent := createInstallerStyleLocalSettings()

	settingsPath := filepath.Join(claudePath, "settings.json")
	localSettingsPath := filepath.Join(claudePath, "settings.local.json")

	if err := os.WriteFile(settingsPath, []byte(settingsContent), 0644); err != nil {
		t.Fatalf("Failed to write settings.json: %v", err)
	}
	if err := os.WriteFile(localSettingsPath, []byte(localSettingsContent), 0644); err != nil {
		t.Fatalf("Failed to write settings.local.json: %v", err)
	}

	// Initialize settings manager
	sm := NewDefaultSettingsManager()
	sm.SetVerbose(true)

	t.Run("BackupAndCleanCycle", func(t *testing.T) {
		// Create backup
		backupPath, err := sm.BackupSettings(claudePath)
		if err != nil {
			t.Fatalf("Backup failed: %v", err)
		}

		// Verify backup was created
		if backupPath == "" {
			t.Fatal("Backup path is empty")
		}

		// Clean settings
		cleanedFiles, err := sm.CleanSettings(claudePath)
		if err != nil {
			t.Fatalf("Clean failed: %v", err)
		}

		if len(cleanedFiles) != 2 {
			t.Errorf("Expected 2 cleaned files, got %d", len(cleanedFiles))
		}

		// Verify startup content was removed
		verifyStartupContentRemovalFromFiles(t, cleanedFiles)

		// Restore from backup
		if err := sm.RestoreSettings(claudePath, backupPath); err != nil {
			t.Fatalf("Restore failed: %v", err)
		}

		// Verify startup content is back
		verifyStartupContentRestored(t, settingsPath, localSettingsPath)
	})

	t.Run("PreviewChanges", func(t *testing.T) {
		preview, err := sm.PreviewChanges(claudePath)
		if err != nil {
			t.Fatalf("Preview failed: %v", err)
		}

		// Should preview changes for both files
		if len(preview) != 2 {
			t.Errorf("Expected preview for 2 files, got %d", len(preview))
		}

		for filename, filePreview := range preview {
			previewMap, ok := filePreview.(map[string]interface{})
			if !ok {
				t.Errorf("Preview for %s is not a map", filename)
				continue
			}

			entriesRemoved, ok := previewMap["entries_to_remove"].(int)
			if !ok || entriesRemoved == 0 {
				t.Errorf("Expected entries to be removed for %s", filename)
			}
		}
	})

	t.Run("PartialStartupContent", func(t *testing.T) {
		// Create settings with mixed content (some startup, some user)
		mixedContent := createMixedContentSettings()
		mixedPath := filepath.Join(claudePath, "mixed.json")
		if err := os.WriteFile(mixedPath, []byte(mixedContent), 0644); err != nil {
			t.Fatalf("Failed to write mixed settings: %v", err)
		}

		// Clean the mixed settings
		modified, err := sm.cleanSettingsFile(mixedPath)
		if err != nil {
			t.Fatalf("Failed to clean mixed settings: %v", err)
		}

		if !modified {
			t.Error("Expected mixed settings to be modified")
		}

		// Verify user content was preserved and startup content was removed
		data, err := os.ReadFile(mixedPath)
		if err != nil {
			t.Fatalf("Failed to read cleaned mixed settings: %v", err)
		}

		var cleanedSettings map[string]interface{}
		if err := json.Unmarshal(data, &cleanedSettings); err != nil {
			t.Fatalf("Cleaned settings is not valid JSON: %v", err)
		}

		// Should preserve user settings
		if userSetting, exists := cleanedSettings["userSetting"]; !exists {
			t.Error("User setting was removed")
		} else if userStr, ok := userSetting.(string); !ok || userStr != "preserve this" {
			t.Error("User setting was modified")
		}

		// Should remove startup content
		verifyStartupContentRemoved(t, cleanedSettings)
	})
}

// TestSettingsManagerErrorRecovery tests error recovery scenarios
func TestSettingsManagerErrorRecovery(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "settings_error_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	claudePath := filepath.Join(tempDir, "claude")
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude dir: %v", err)
	}

	sm := NewDefaultSettingsManager()

	t.Run("CorruptedJSONRecovery", func(t *testing.T) {
		corruptedPath := filepath.Join(claudePath, "corrupted.json")
		corruptedContent := `{
			"hooks": {
				"PreToolUse": [
					{
						"matcher": "Task",
						"hooks": [
							{
								"command": "the-startup log --assistant",
								"_source": "the-startup"
							}
						]
					}
				]
			},
			"invalidField": { // Missing closing brace and value
		}`

		if err := os.WriteFile(corruptedPath, []byte(corruptedContent), 0644); err != nil {
			t.Fatalf("Failed to write corrupted settings: %v", err)
		}

		// Should fail gracefully
		_, err := sm.cleanSettingsFile(corruptedPath)
		if err == nil {
			t.Error("Expected error when cleaning corrupted JSON")
		}

		// Original file should remain untouched
		data, err := os.ReadFile(corruptedPath)
		if err != nil {
			t.Fatalf("Failed to read original file: %v", err)
		}

		if string(data) != corruptedContent {
			t.Error("Original corrupted file was modified despite error")
		}
	})

	t.Run("BackupFailureRecovery", func(t *testing.T) {
		// Test backup with nonexistent source files
		nonexistentClaudePath := filepath.Join(tempDir, "nonexistent")
		
		// Should fail gracefully when no settings files exist
		_, err := sm.BackupSettings(nonexistentClaudePath)
		if err == nil {
			t.Error("Expected backup to fail with nonexistent directory")
		}

		// Test backup creation failure by providing an invalid path
		invalidPath := "/root/restricted/path/that/should/not/exist"
		_, err = sm.BackupSettings(invalidPath)
		if err == nil {
			t.Error("Expected backup to fail with restricted path")
		}
	})
}

// Helper functions for integration tests

func createInstallerStyleSettings() string {
	// This matches the exact structure created by installer/installer.go
	return `{
  "permissions": {
    "additionalDirectories": [
      "~/.the-startup",
      "/some/other/user/directory"
    ]
  },
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Task",
        "hooks": [
          {
            "type": "command",
            "command": "~/.the-startup/bin/the-startup log --assistant",
            "_source": "the-startup"
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Task",
        "hooks": [
          {
            "type": "command",
            "command": "~/.the-startup/bin/the-startup log --user",
            "_source": "the-startup"
          }
        ]
      }
    ]
  },
  "statusLine": {
    "type": "command",
    "command": "~/.the-startup/bin/the-startup statusline"
  },
  "userCustomSetting": "this should be preserved",
  "outputStyle": "default"
}`
}

func createInstallerStyleLocalSettings() string {
	return `{
  "outputStyle": "The Startup",
  "userLocalCustomSetting": "this should also be preserved"
}`
}

func createMixedContentSettings() string {
	return `{
  "permissions": {
    "additionalDirectories": [
      "~/.the-startup",
      "/user/custom/directory"
    ]
  },
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Task",
        "hooks": [
          {
            "type": "command",
            "command": "~/.the-startup/bin/the-startup log --assistant",
            "_source": "the-startup"
          },
          {
            "type": "command",
            "command": "user-custom-hook",
            "_source": "user"
          }
        ]
      }
    ]
  },
  "statusLine": {
    "type": "command",
    "command": "~/.the-startup/bin/the-startup statusline"
  },
  "userSetting": "preserve this",
  "customHooks": {
    "someUserHook": "keep this too"
  }
}`
}

func verifyStartupContentRemovalFromFiles(t *testing.T, files []string) {
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read cleaned file %s: %v", file, err)
			continue
		}

		var settings map[string]interface{}
		if err := json.Unmarshal(data, &settings); err != nil {
			t.Errorf("Cleaned file %s contains invalid JSON: %v", file, err)
			continue
		}

		// Verify no startup content remains
		verifyStartupContentRemoved(t, settings)
	}
}

func verifyStartupContentRestored(t *testing.T, settingsPath, localSettingsPath string) {
	// Check main settings
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("Failed to read restored settings: %v", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatalf("Restored settings is not valid JSON: %v", err)
	}

	// Should have startup content back
	if !hasStartupContent(settings) {
		t.Error("Startup content was not restored to main settings")
	}

	// Check local settings
	data, err = os.ReadFile(localSettingsPath)
	if err != nil {
		t.Fatalf("Failed to read restored local settings: %v", err)
	}

	var localSettings map[string]interface{}
	if err := json.Unmarshal(data, &localSettings); err != nil {
		t.Fatalf("Restored local settings is not valid JSON: %v", err)
	}

	// Should have "The Startup" output style back
	if outputStyle, exists := localSettings["outputStyle"]; !exists {
		t.Error("outputStyle was not restored to local settings")
	} else if styleStr, ok := outputStyle.(string); !ok || styleStr != "The Startup" {
		t.Errorf("Expected outputStyle to be 'The Startup', got %v", outputStyle)
	}
}