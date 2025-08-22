package uninstaller

import (
	"fmt"
	"os"
	"path/filepath"
)

// ExampleDefaultSettingsManager demonstrates how to use the settings manager
func ExampleDefaultSettingsManager() {
	// Create a temporary directory for this example
	tempDir, _ := os.MkdirTemp("", "settings_example")
	defer os.RemoveAll(tempDir)

	claudePath := filepath.Join(tempDir, ".claude")
	os.MkdirAll(claudePath, 0755)

	// Create sample settings files with startup content
	settingsContent := `{
  "permissions": {
    "additionalDirectories": ["~/.the-startup"]
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
    ]
  },
  "statusLine": {
    "type": "command",
    "command": "~/.the-startup/bin/the-startup statusline"
  },
  "userSetting": "this will be preserved"
}`

	settingsPath := filepath.Join(claudePath, "settings.json")
	os.WriteFile(settingsPath, []byte(settingsContent), 0644)

	// Initialize the settings manager
	sm := NewDefaultSettingsManager()
	sm.SetVerbose(false) // Reduce output for example

	// Preview what changes would be made
	preview, err := sm.PreviewChanges(claudePath)
	if err == nil && len(preview) > 0 {
		fmt.Println("Preview shows changes would be made to settings files")
	}

	// Create a backup before making changes
	backupPath, err := sm.BackupSettings(claudePath)
	if err == nil && backupPath != "" {
		fmt.Println("Backup created successfully")
	}

	// Clean the startup-managed settings
	cleanedFiles, err := sm.CleanSettings(claudePath)
	if err == nil && len(cleanedFiles) > 0 {
		fmt.Println("Startup settings cleaned successfully")
	}

	// Restore from backup if needed
	err = sm.RestoreSettings(claudePath, backupPath)
	if err == nil {
		fmt.Println("Settings restored from backup")
	}

	// Output:
	// Preview shows changes would be made to settings files
	// Backup created successfully
	// Startup settings cleaned successfully
	// Settings restored from backup
}

// ExampleDefaultSettingsManager_dryRun demonstrates dry-run mode
func ExampleDefaultSettingsManager_dryRun() {
	// Create a temporary directory for this example
	tempDir, _ := os.MkdirTemp("", "settings_dryrun_example")
	defer os.RemoveAll(tempDir)

	claudePath := filepath.Join(tempDir, ".claude")
	os.MkdirAll(claudePath, 0755)

	// Create sample settings
	settingsPath := filepath.Join(claudePath, "settings.json")
	os.WriteFile(settingsPath, []byte(`{"hooks": {"PreToolUse": [{"matcher": "Task", "hooks": [{"command": "the-startup log --assistant", "_source": "the-startup"}]}]}}`), 0644)

	// Initialize settings manager in dry-run mode
	sm := NewDefaultSettingsManager()
	sm.SetDryRun(true)
	sm.SetVerbose(false)

	// Operations will show what they would do without actually doing it
	_, err := sm.BackupSettings(claudePath)
	if err == nil {
		fmt.Println("Dry-run backup completed")
	}

	cleanedFiles, err := sm.CleanSettings(claudePath)
	if err == nil && len(cleanedFiles) > 0 {
		fmt.Println("Dry-run cleaning completed")
	}

	// Output:
	// Dry-run backup completed
	// Dry-run cleaning completed
}

// ExampleNewWithDefaults demonstrates integration with the full uninstaller
func ExampleNewWithDefaults() {
	// Create uninstaller with default settings manager
	uninstaller := NewWithDefaults()
	
	// Configure options - this will also configure the settings manager
	uninstaller.SetOptions(
		false, // dryRun
		false, // forceRemove
		true,  // createBackup
		false, // verbose
	)

	fmt.Println("Uninstaller created with settings manager")
	
	// The settings manager is now ready to use as part of the uninstall process
	// When ExecuteUninstallPlan is called, it will:
	// 1. Create backups if requested
	// 2. Clean settings files 
	// 3. Remove tracked files
	// 4. Clean up directories

	// Output:
	// Uninstaller created with settings manager
}