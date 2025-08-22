package uninstaller

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultSettingsManager(t *testing.T) {
	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "settings_manager_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewDefaultSettingsManager()
	sm.SetVerbose(false) // Reduce noise in tests

	t.Run("NewDefaultSettingsManager", func(t *testing.T) {
		if sm == nil {
			t.Fatal("NewDefaultSettingsManager returned nil")
		}
		if sm.maxBackupRetention != 10 {
			t.Errorf("Expected default max backup retention to be 10, got %d", sm.maxBackupRetention)
		}
	})

	t.Run("SetOptions", func(t *testing.T) {
		sm.SetDryRun(true)
		sm.SetVerbose(true)
		sm.SetMaxBackupRetention(5)

		if !sm.dryRun {
			t.Error("SetDryRun(true) did not set dryRun to true")
		}
		if !sm.verbose {
			t.Error("SetVerbose(true) did not set verbose to true")
		}
		if sm.maxBackupRetention != 5 {
			t.Error("SetMaxBackupRetention(5) did not set maxBackupRetention to 5")
		}

		// Reset for other tests
		sm.SetDryRun(false)
		sm.SetVerbose(false)
	})

	t.Run("BackupSettings", func(t *testing.T) {
		testBackupSettings(t, tempDir, sm)
	})

	t.Run("CleanSettings", func(t *testing.T) {
		testCleanSettings(t, tempDir, sm)
	})

	t.Run("RestoreSettings", func(t *testing.T) {
		testRestoreSettings(t, tempDir, sm)
	})

	t.Run("PreviewChanges", func(t *testing.T) {
		testPreviewChanges(t, tempDir, sm)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		testErrorHandling(t, tempDir, sm)
	})

	t.Run("EdgeCases", func(t *testing.T) {
		testEdgeCases(t, tempDir, sm)
	})
}

func testBackupSettings(t *testing.T, tempDir string, sm *DefaultSettingsManager) {
	claudePath := filepath.Join(tempDir, "claude_backup_test")
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude dir: %v", err)
	}

	// Create test settings files
	settingsJSON := createTestSettingsJSON()
	settingsLocalJSON := createTestSettingsLocalJSON()

	settingsPath := filepath.Join(claudePath, "settings.json")
	localSettingsPath := filepath.Join(claudePath, "settings.local.json")

	if err := os.WriteFile(settingsPath, []byte(settingsJSON), 0644); err != nil {
		t.Fatalf("Failed to write settings.json: %v", err)
	}
	if err := os.WriteFile(localSettingsPath, []byte(settingsLocalJSON), 0644); err != nil {
		t.Fatalf("Failed to write settings.local.json: %v", err)
	}

	// Test backup creation
	backupDir, err := sm.BackupSettings(claudePath)
	if err != nil {
		t.Fatalf("BackupSettings failed: %v", err)
	}

	if backupDir == "" {
		t.Fatal("BackupSettings returned empty backup path")
	}

	// Verify backup files exist
	backupFiles, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("Failed to read backup directory: %v", err)
	}

	expectedBackups := []string{"settings.json.backup", "settings.local.json.backup"}
	found := make(map[string]bool)

	for _, file := range backupFiles {
		name := file.Name()
		for _, expected := range expectedBackups {
			if strings.HasPrefix(name, expected) {
				found[expected] = true
			}
		}
	}

	for _, expected := range expectedBackups {
		if !found[expected] {
			t.Errorf("Expected backup file %s not found", expected)
		}
	}

	// Test dry run mode
	sm.SetDryRun(true)
	dryRunBackupDir, err := sm.BackupSettings(claudePath)
	if err != nil {
		t.Errorf("Dry run backup failed: %v", err)
	}
	if dryRunBackupDir != "" {
		t.Error("Dry run should return empty backup path")
	}
	sm.SetDryRun(false)

	// Test backup cleanup (create multiple backups)
	sm.SetMaxBackupRetention(2)
	time.Sleep(time.Millisecond * 10) // Ensure different timestamps
	sm.BackupSettings(claudePath)
	time.Sleep(time.Millisecond * 10)
	sm.BackupSettings(claudePath)

	// Check that old backups are cleaned up
	backupFiles, _ = os.ReadDir(backupDir)
	backupCount := 0
	for _, file := range backupFiles {
		if strings.Contains(file.Name(), ".backup.") {
			backupCount++
		}
	}

	// Should have at most 4 files (2 retention * 2 files)
	if backupCount > 4 {
		t.Errorf("Expected at most 4 backup files after cleanup, got %d", backupCount)
	}
}

func testCleanSettings(t *testing.T, tempDir string, sm *DefaultSettingsManager) {
	claudePath := filepath.Join(tempDir, "claude_clean_test")
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude dir: %v", err)
	}

	// Test cases for different settings structures
	testCases := []struct {
		name           string
		input          string
		expectModified bool
		expectError    bool
	}{
		{
			name:           "settings_with_startup_hooks",
			input:          createTestSettingsJSON(),
			expectModified: true,
			expectError:    false,
		},
		{
			name:           "settings_without_startup_content",
			input:          createCleanSettingsJSON(),
			expectModified: false,
			expectError:    false,
		},
		{
			name:           "empty_settings",
			input:          "{}",
			expectModified: false,
			expectError:    false,
		},
		{
			name:           "invalid_json",
			input:          "{invalid json",
			expectModified: false,
			expectError:    true,
		},
		{
			name:           "empty_file",
			input:          "",
			expectModified: false,
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settingsPath := filepath.Join(claudePath, fmt.Sprintf("settings_%s.json", tc.name))
			
			if err := os.WriteFile(settingsPath, []byte(tc.input), 0644); err != nil {
				t.Fatalf("Failed to write test settings: %v", err)
			}

			modified, err := sm.cleanSettingsFile(settingsPath)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tc.expectModified != modified {
				t.Errorf("Expected modified=%v, got %v", tc.expectModified, modified)
			}

			// If successful and modified, verify the content
			if !tc.expectError && tc.expectModified {
				data, err := os.ReadFile(settingsPath)
				if err != nil {
					t.Errorf("Failed to read cleaned settings: %v", err)
				} else {
					var settings map[string]interface{}
					if err := json.Unmarshal(data, &settings); err != nil {
						t.Errorf("Cleaned settings is not valid JSON: %v", err)
					} else {
						// Verify that startup content was removed
						verifyStartupContentRemoved(t, settings)
					}
				}
			}
		})
	}

	// Test dry run mode
	settingsPath := filepath.Join(claudePath, "settings_dryrun.json")
	if err := os.WriteFile(settingsPath, []byte(createTestSettingsJSON()), 0644); err != nil {
		t.Fatalf("Failed to write dry run test settings: %v", err)
	}

	sm.SetDryRun(true)
	modified, err := sm.cleanSettingsFile(settingsPath)
	if err != nil {
		t.Errorf("Dry run clean failed: %v", err)
	}
	if !modified {
		t.Error("Dry run should report modifications would be made")
	}

	// Verify file wasn't actually modified
	data, _ := os.ReadFile(settingsPath)
	var settings map[string]interface{}
	json.Unmarshal(data, &settings)
	if !hasStartupContent(settings) {
		t.Error("Dry run should not have actually modified the file")
	}
	sm.SetDryRun(false)
}

func testRestoreSettings(t *testing.T, tempDir string, sm *DefaultSettingsManager) {
	claudePath := filepath.Join(tempDir, "claude_restore_test")
	backupDir := filepath.Join(claudePath, "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	// Create original settings
	settingsPath := filepath.Join(claudePath, "settings.json")
	originalSettings := createTestSettingsJSON()
	if err := os.WriteFile(settingsPath, []byte(originalSettings), 0644); err != nil {
		t.Fatalf("Failed to write original settings: %v", err)
	}

	// Create backup
	backupPath, err := sm.BackupSettings(claudePath)
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	// Modify the settings file
	modifiedSettings := createCleanSettingsJSON()
	if err := os.WriteFile(settingsPath, []byte(modifiedSettings), 0644); err != nil {
		t.Fatalf("Failed to write modified settings: %v", err)
	}

	// Restore from backup
	if err := sm.RestoreSettings(claudePath, backupPath); err != nil {
		t.Fatalf("Failed to restore settings: %v", err)
	}

	// Verify restoration
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("Failed to read restored settings: %v", err)
	}

	var restoredSettings map[string]interface{}
	if err := json.Unmarshal(data, &restoredSettings); err != nil {
		t.Fatalf("Restored settings is not valid JSON: %v", err)
	}

	if !hasStartupContent(restoredSettings) {
		t.Error("Restored settings does not contain expected startup content")
	}

	// Test dry run restore
	sm.SetDryRun(true)
	if err := sm.RestoreSettings(claudePath, backupPath); err != nil {
		t.Errorf("Dry run restore failed: %v", err)
	}
	sm.SetDryRun(false)

	// Test restore with invalid backup path
	if err := sm.RestoreSettings(claudePath, "/nonexistent/path"); err == nil {
		t.Error("Expected error when restoring from nonexistent path")
	}
}

func testPreviewChanges(t *testing.T, tempDir string, sm *DefaultSettingsManager) {
	claudePath := filepath.Join(tempDir, "claude_preview_test")
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude dir: %v", err)
	}

	// Create test settings
	settingsPath := filepath.Join(claudePath, "settings.json")
	if err := os.WriteFile(settingsPath, []byte(createTestSettingsJSON()), 0644); err != nil {
		t.Fatalf("Failed to write settings: %v", err)
	}

	preview, err := sm.PreviewChanges(claudePath)
	if err != nil {
		t.Fatalf("PreviewChanges failed: %v", err)
	}

	if len(preview) == 0 {
		t.Error("Expected preview to contain changes")
	}

	settingsPreview, exists := preview["settings.json"]
	if !exists {
		t.Error("Expected preview for settings.json")
	}

	previewMap, ok := settingsPreview.(map[string]interface{})
	if !ok {
		t.Error("Expected preview to be a map")
	}

	// Verify preview contains expected fields
	expectedFields := []string{"original_entries", "remaining_entries", "entries_to_remove", "removed_keys"}
	for _, field := range expectedFields {
		if _, exists := previewMap[field]; !exists {
			t.Errorf("Expected preview field %s not found", field)
		}
	}

	// Verify that entries would be removed
	entriesRemoved, ok := previewMap["entries_to_remove"].(int)
	if !ok || entriesRemoved == 0 {
		t.Error("Expected some entries to be marked for removal")
	}
}

func testErrorHandling(t *testing.T, tempDir string, sm *DefaultSettingsManager) {
	// Test with nonexistent directory
	nonexistentPath := filepath.Join(tempDir, "nonexistent")
	_, err := sm.BackupSettings(nonexistentPath)
	if err == nil {
		t.Error("Expected error when backing up nonexistent directory")
	}

	_, err = sm.CleanSettings(nonexistentPath)
	if err != nil {
		t.Errorf("CleanSettings should handle nonexistent directory gracefully: %v", err)
	}

	// Test with insufficient permissions (create a directory we can't write to)
	restrictedPath := filepath.Join(tempDir, "restricted")
	if err := os.MkdirAll(restrictedPath, 0000); err != nil {
		t.Fatalf("Failed to create restricted directory: %v", err)
	}
	defer os.Chmod(restrictedPath, 0755) // Cleanup

	_, err = sm.BackupSettings(restrictedPath)
	if err == nil {
		t.Error("Expected error when backing up to restricted directory")
	}
}

func testEdgeCases(t *testing.T, tempDir string, sm *DefaultSettingsManager) {
	claudePath := filepath.Join(tempDir, "claude_edge_test")
	if err := os.MkdirAll(claudePath, 0755); err != nil {
		t.Fatalf("Failed to create claude dir: %v", err)
	}

	// Test with malformed JSON that has startup content
	malformedJSON := `{
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
		"invalidField": {unclosed object
	}`

	settingsPath := filepath.Join(claudePath, "malformed.json")
	if err := os.WriteFile(settingsPath, []byte(malformedJSON), 0644); err != nil {
		t.Fatalf("Failed to write malformed settings: %v", err)
	}

	_, err := sm.cleanSettingsFile(settingsPath)
	if err == nil {
		t.Error("Expected error when cleaning malformed JSON")
	}

	// Test with very large settings file (within reason)
	largeSettings := createLargeSettingsJSON()
	largeSettingsPath := filepath.Join(claudePath, "large.json")
	if err := os.WriteFile(largeSettingsPath, []byte(largeSettings), 0644); err != nil {
		t.Fatalf("Failed to write large settings: %v", err)
	}

	modified, err := sm.cleanSettingsFile(largeSettingsPath)
	if err != nil {
		t.Errorf("Failed to clean large settings file: %v", err)
	}
	if !modified {
		t.Error("Expected large settings file to be modified")
	}

	// Test cleaning file with only startup content (should result in empty object)
	startupOnlyJSON := `{
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
		"permissions": {
			"additionalDirectories": ["~/.the-startup"]
		},
		"statusLine": {
			"command": "the-startup statusline"
		}
	}`

	startupOnlyPath := filepath.Join(claudePath, "startup_only.json")
	if err := os.WriteFile(startupOnlyPath, []byte(startupOnlyJSON), 0644); err != nil {
		t.Fatalf("Failed to write startup-only settings: %v", err)
	}

	modified, err = sm.cleanSettingsFile(startupOnlyPath)
	if err != nil {
		t.Errorf("Failed to clean startup-only settings: %v", err)
	}
	if !modified {
		t.Error("Expected startup-only settings to be modified")
	}

	// Verify file is now empty or contains only empty object
	data, err := os.ReadFile(startupOnlyPath)
	if err != nil {
		t.Fatalf("Failed to read cleaned startup-only settings: %v", err)
	}

	var cleanedSettings map[string]interface{}
	if err := json.Unmarshal(data, &cleanedSettings); err != nil {
		t.Errorf("Cleaned settings is not valid JSON: %v", err)
	}

	if len(cleanedSettings) != 0 {
		t.Errorf("Expected cleaned settings to be empty, got %d entries", len(cleanedSettings))
	}
}

// Helper functions for creating test data

func createTestSettingsJSON() string {
	return `{
  "permissions": {
    "additionalDirectories": [
      "~/.the-startup",
      "/some/other/dir"
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
            "command": "some-other-command",
            "_source": "other"
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
  "userSetting": "should be preserved",
  "outputStyle": "default"
}`
}

func createTestSettingsLocalJSON() string {
	return `{
  "outputStyle": "The Startup",
  "userLocalSetting": "should be preserved"
}`
}

func createCleanSettingsJSON() string {
	return `{
  "permissions": {
    "additionalDirectories": [
      "/some/other/dir"
    ]
  },
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Other",
        "hooks": [
          {
            "type": "command",
            "command": "some-other-command",
            "_source": "other"
          }
        ]
      }
    ]
  },
  "userSetting": "should be preserved",
  "outputStyle": "default"
}`
}

func createLargeSettingsJSON() string {
	// Create a settings file with many entries including startup content
	base := `{
  "permissions": {
    "additionalDirectories": [
      "~/.the-startup"`

	// Add many additional directories
	for i := 0; i < 100; i++ {
		base += fmt.Sprintf(`,
      "/path/to/dir%d"`, i)
	}

	base += `
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
          }`

	// Add many other hooks
	for i := 0; i < 50; i++ {
		base += fmt.Sprintf(`,
          {
            "type": "command",
            "command": "command%d",
            "_source": "other"
          }`, i)
	}

	base += `
        ]
      }
    ]
  },
  "statusLine": {
    "type": "command",
    "command": "~/.the-startup/bin/the-startup statusline"
  }`

	// Add many user settings
	for i := 0; i < 50; i++ {
		base += fmt.Sprintf(`,
  "userSetting%d": "value%d"`, i, i)
	}

	base += `
}`

	return base
}

func verifyStartupContentRemoved(t *testing.T, settings map[string]interface{}) {
	// Check that startup-managed content was removed
	if permissions, exists := settings["permissions"]; exists {
		if permMap, ok := permissions.(map[string]interface{}); ok {
			if additionalDirs, exists := permMap["additionalDirectories"]; exists {
				if dirs, ok := additionalDirs.([]interface{}); ok {
					for _, dir := range dirs {
						if dirStr, ok := dir.(string); ok {
							if strings.Contains(dirStr, "the-startup") {
								t.Error("Found startup directory that should have been removed")
							}
						}
					}
				}
			}
		}
	}

	if hooks, exists := settings["hooks"]; exists {
		if hooksMap, ok := hooks.(map[string]interface{}); ok {
			for _, hookList := range hooksMap {
				if hooks, ok := hookList.([]interface{}); ok {
					for _, hook := range hooks {
						if hookMap, ok := hook.(map[string]interface{}); ok {
							if hooksList, exists := hookMap["hooks"]; exists {
								if entries, ok := hooksList.([]interface{}); ok {
									for _, entry := range entries {
										if entryMap, ok := entry.(map[string]interface{}); ok {
											if source, exists := entryMap["_source"]; exists {
												if sourceStr, ok := source.(string); ok && sourceStr == "the-startup" {
													t.Error("Found startup hook that should have been removed")
												}
											}
											if command, exists := entryMap["command"]; exists {
												if cmdStr, ok := command.(string); ok && strings.Contains(cmdStr, "the-startup") {
													t.Error("Found startup command that should have been removed")
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if statusLine, exists := settings["statusLine"]; exists {
		if statusMap, ok := statusLine.(map[string]interface{}); ok {
			if command, exists := statusMap["command"]; exists {
				if cmdStr, ok := command.(string); ok && strings.Contains(cmdStr, "the-startup") {
					t.Error("Found startup statusLine that should have been removed")
				}
			}
		}
	}

	if outputStyle, exists := settings["outputStyle"]; exists {
		if styleStr, ok := outputStyle.(string); ok && styleStr == "The Startup" {
			t.Error("Found startup outputStyle that should have been removed")
		}
	}
}

func hasStartupContent(settings map[string]interface{}) bool {
	// Check if settings contain startup-managed content
	if permissions, exists := settings["permissions"]; exists {
		if permMap, ok := permissions.(map[string]interface{}); ok {
			if additionalDirs, exists := permMap["additionalDirectories"]; exists {
				if dirs, ok := additionalDirs.([]interface{}); ok {
					for _, dir := range dirs {
						if dirStr, ok := dir.(string); ok {
							if strings.Contains(dirStr, "the-startup") {
								return true
							}
						}
					}
				}
			}
		}
	}

	if hooks, exists := settings["hooks"]; exists {
		if hooksMap, ok := hooks.(map[string]interface{}); ok {
			for _, hookList := range hooksMap {
				if hooks, ok := hookList.([]interface{}); ok {
					for _, hook := range hooks {
						if hookMap, ok := hook.(map[string]interface{}); ok {
							if hooksList, exists := hookMap["hooks"]; exists {
								if entries, ok := hooksList.([]interface{}); ok {
									for _, entry := range entries {
										if entryMap, ok := entry.(map[string]interface{}); ok {
											if source, exists := entryMap["_source"]; exists {
												if sourceStr, ok := source.(string); ok && sourceStr == "the-startup" {
													return true
												}
											}
											if command, exists := entryMap["command"]; exists {
												if cmdStr, ok := command.(string); ok && strings.Contains(cmdStr, "the-startup") {
													return true
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if statusLine, exists := settings["statusLine"]; exists {
		if statusMap, ok := statusLine.(map[string]interface{}); ok {
			if command, exists := statusMap["command"]; exists {
				if cmdStr, ok := command.(string); ok && strings.Contains(cmdStr, "the-startup") {
					return true
				}
			}
		}
	}

	if outputStyle, exists := settings["outputStyle"]; exists {
		if styleStr, ok := outputStyle.(string); ok && styleStr == "The Startup" {
			return true
		}
	}

	return false
}