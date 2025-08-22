package uninstaller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DefaultSettingsManager implements the SettingsManager interface
type DefaultSettingsManager struct {
	dryRun                bool
	maxBackupRetention    int // Number of backup files to keep
	verbose               bool
}

// NewDefaultSettingsManager creates a new DefaultSettingsManager
func NewDefaultSettingsManager() *DefaultSettingsManager {
	return &DefaultSettingsManager{
		dryRun:             false,
		maxBackupRetention: 10, // Keep last 10 backups by default
		verbose:            false,
	}
}

// SetDryRun enables or disables dry-run mode
func (sm *DefaultSettingsManager) SetDryRun(dryRun bool) {
	sm.dryRun = dryRun
}

// SetVerbose enables or disables verbose logging
func (sm *DefaultSettingsManager) SetVerbose(verbose bool) {
	sm.verbose = verbose
}

// SetMaxBackupRetention sets the maximum number of backup files to keep
func (sm *DefaultSettingsManager) SetMaxBackupRetention(count int) {
	sm.maxBackupRetention = count
}

// BackupSettings creates a timestamped backup of settings files before modification
func (sm *DefaultSettingsManager) BackupSettings(claudePath string) (string, error) {
	if sm.dryRun {
		if sm.verbose {
			fmt.Println("DRY RUN: Would create settings backup")
		}
		return "", nil
	}

	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(claudePath, "backups")
	
	// Ensure backup directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	var backupFiles []string
	settingsFiles := []string{"settings.json", "settings.local.json"}

	for _, filename := range settingsFiles {
		sourcePath := filepath.Join(claudePath, filename)
		
		// Check if source file exists
		sourceInfo, err := os.Stat(sourcePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Skip files that don't exist
			}
			return "", fmt.Errorf("failed to stat %s: %w", sourcePath, err)
		}

		// Create backup filename with timestamp
		backupFilename := fmt.Sprintf("%s.backup.%s", filename, timestamp)
		backupPath := filepath.Join(backupDir, backupFilename)

		// Perform atomic backup operation
		if err := sm.atomicFileCopy(sourcePath, backupPath); err != nil {
			return "", fmt.Errorf("failed to backup %s: %w", filename, err)
		}

		// Verify backup integrity
		if err := sm.verifyBackupIntegrity(sourcePath, backupPath, sourceInfo.Size()); err != nil {
			// Clean up failed backup
			os.Remove(backupPath)
			return "", fmt.Errorf("backup verification failed for %s: %w", filename, err)
		}

		backupFiles = append(backupFiles, backupPath)
		
		if sm.verbose {
			fmt.Printf("✓ Backed up %s to %s\n", filename, backupFilename)
		}
	}

	if len(backupFiles) == 0 {
		return "", fmt.Errorf("no settings files found to backup")
	}

	// Clean up old backups
	if err := sm.cleanupOldBackups(backupDir); err != nil {
		// Log warning but don't fail the operation
		if sm.verbose {
			fmt.Printf("Warning: Failed to cleanup old backups: %v\n", err)
		}
	}

	if sm.verbose {
		fmt.Printf("✓ Created backup with %d files in %s\n", len(backupFiles), backupDir)
	}

	return backupDir, nil
}

// CleanSettings removes the-startup related settings from Claude configuration files
func (sm *DefaultSettingsManager) CleanSettings(claudePath string) ([]string, error) {
	var cleanedSettings []string
	settingsFiles := []string{"settings.json", "settings.local.json"}

	for _, filename := range settingsFiles {
		settingsPath := filepath.Join(claudePath, filename)
		
		// Check if file exists
		if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
			continue // Skip files that don't exist
		}

		cleaned, err := sm.cleanSettingsFile(settingsPath)
		if err != nil {
			return cleanedSettings, fmt.Errorf("failed to clean %s: %w", filename, err)
		}

		if cleaned {
			cleanedSettings = append(cleanedSettings, settingsPath)
			if sm.verbose {
				fmt.Printf("✓ Cleaned startup-managed settings from %s\n", filename)
			}
		}
	}

	return cleanedSettings, nil
}

// RestoreSettings restores settings from a backup directory
func (sm *DefaultSettingsManager) RestoreSettings(claudePath, backupPath string) error {
	if sm.dryRun {
		if sm.verbose {
			fmt.Printf("DRY RUN: Would restore settings from %s\n", backupPath)
		}
		return nil
	}

	// Find all backup files in the backup directory
	entries, err := os.ReadDir(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	var restoredFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".backup") && !strings.Contains(filename, ".backup.") {
			continue
		}

		// Extract original filename (e.g., "settings.json.backup.20240101_120000" -> "settings.json")
		originalName := strings.Split(filename, ".backup")[0]
		
		backupFilePath := filepath.Join(backupPath, filename)
		restorePath := filepath.Join(claudePath, originalName)

		// Perform atomic restore
		if err := sm.atomicFileCopy(backupFilePath, restorePath); err != nil {
			return fmt.Errorf("failed to restore %s: %w", originalName, err)
		}

		restoredFiles = append(restoredFiles, originalName)
		
		if sm.verbose {
			fmt.Printf("✓ Restored %s from backup\n", originalName)
		}
	}

	if len(restoredFiles) == 0 {
		return fmt.Errorf("no backup files found in %s", backupPath)
	}

	if sm.verbose {
		fmt.Printf("✓ Successfully restored %d files from backup\n", len(restoredFiles))
	}

	return nil
}

// PreviewChanges shows what changes would be made without actually applying them
func (sm *DefaultSettingsManager) PreviewChanges(claudePath string) (map[string]interface{}, error) {
	preview := make(map[string]interface{})
	settingsFiles := []string{"settings.json", "settings.local.json"}

	for _, filename := range settingsFiles {
		settingsPath := filepath.Join(claudePath, filename)
		
		if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
			continue
		}

		filePreview, err := sm.previewSettingsFileChanges(settingsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to preview changes for %s: %w", filename, err)
		}

		if len(filePreview) > 0 {
			preview[filename] = filePreview
		}
	}

	return preview, nil
}

// Private helper methods

// cleanSettingsFile removes the-startup managed entries from a single settings file
func (sm *DefaultSettingsManager) cleanSettingsFile(settingsPath string) (bool, error) {
	if sm.dryRun {
		if sm.verbose {
			fmt.Printf("DRY RUN: Would clean settings in %s\n", settingsPath)
		}
		return true, nil
	}

	// Read and parse the settings file
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return false, fmt.Errorf("failed to read settings file: %w", err)
	}

	// Handle empty files
	if len(bytes.TrimSpace(data)) == 0 {
		if sm.verbose {
			fmt.Printf("Settings file %s is empty, skipping\n", filepath.Base(settingsPath))
		}
		return false, nil
	}

	// Parse JSON with error recovery
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return false, fmt.Errorf("failed to parse settings JSON (file may be corrupted): %w", err)
	}

	// Track if any changes were made
	modified := false

	// Remove startup-managed sections
	modified = sm.removeStartupManagedEntries(settings) || modified

	// If no modifications were made, don't rewrite the file
	if !modified {
		if sm.verbose {
			fmt.Printf("No startup-managed settings found in %s\n", filepath.Base(settingsPath))
		}
		return false, nil
	}

	// Write the cleaned settings back to file atomically
	if err := sm.atomicJSONWrite(settingsPath, settings); err != nil {
		return false, fmt.Errorf("failed to write cleaned settings: %w", err)
	}

	return true, nil
}

// removeStartupManagedEntries removes entries that were managed by the-startup
func (sm *DefaultSettingsManager) removeStartupManagedEntries(settings map[string]interface{}) bool {
	modified := false

	// Remove hooks with "_source": "the-startup"
	if hooksSection, exists := settings["hooks"]; exists {
		if hooksMap, ok := hooksSection.(map[string]interface{}); ok {
			for hookType, hookList := range hooksMap {
				if hooks, ok := hookList.([]interface{}); ok {
					cleanedHooks := sm.removeStartupHooks(hooks)
					if len(cleanedHooks) != len(hooks) {
						modified = true
						if len(cleanedHooks) == 0 {
							delete(hooksMap, hookType)
						} else {
							hooksMap[hookType] = cleanedHooks
						}
					}
				}
			}
			
			// Remove empty hooks section
			if len(hooksMap) == 0 {
				delete(settings, "hooks")
			}
		}
	}

	// Remove startup paths from additionalDirectories
	if permissionsSection, exists := settings["permissions"]; exists {
		if permissions, ok := permissionsSection.(map[string]interface{}); ok {
			if additionalDirs, exists := permissions["additionalDirectories"]; exists {
				if dirs, ok := additionalDirs.([]interface{}); ok {
					cleanedDirs := sm.removeStartupDirectories(dirs)
					if len(cleanedDirs) != len(dirs) {
						modified = true
						if len(cleanedDirs) == 0 {
							delete(permissions, "additionalDirectories")
						} else {
							permissions["additionalDirectories"] = cleanedDirs
						}
					}
				}
			}
			
			// Remove empty permissions section
			if len(permissions) == 0 {
				delete(settings, "permissions")
			}
		}
	}

	// Remove startup-managed statusLine
	if statusLineSection, exists := settings["statusLine"]; exists {
		if statusLine, ok := statusLineSection.(map[string]interface{}); ok {
			if command, exists := statusLine["command"]; exists {
				if cmdStr, ok := command.(string); ok {
					if strings.Contains(cmdStr, "the-startup") {
						delete(settings, "statusLine")
						modified = true
					}
				}
			}
		}
	}

	// Remove startup-managed outputStyle (only if it's exactly "The Startup")
	if outputStyle, exists := settings["outputStyle"]; exists {
		if styleStr, ok := outputStyle.(string); ok {
			if styleStr == "The Startup" {
				delete(settings, "outputStyle")
				modified = true
			}
		}
	}

	return modified
}

// removeStartupHooks filters out hooks that contain "the-startup" in their command
func (sm *DefaultSettingsManager) removeStartupHooks(hooks []interface{}) []interface{} {
	var cleaned []interface{}

	for _, hook := range hooks {
		if hookMap, ok := hook.(map[string]interface{}); ok {
			if hooksList, exists := hookMap["hooks"]; exists {
				if hookEntries, ok := hooksList.([]interface{}); ok {
					var cleanedEntries []interface{}
					
					for _, entry := range hookEntries {
						if entryMap, ok := entry.(map[string]interface{}); ok {
							// Check if this is a startup hook by looking for "_source": "the-startup"
							// or "the-startup" in the command
							isStartupHook := false
							
							if source, exists := entryMap["_source"]; exists {
								if sourceStr, ok := source.(string); ok && sourceStr == "the-startup" {
									isStartupHook = true
								}
							}
							
							if command, exists := entryMap["command"]; exists && !isStartupHook {
								if cmdStr, ok := command.(string); ok && strings.Contains(cmdStr, "the-startup") {
									isStartupHook = true
								}
							}
							
							if !isStartupHook {
								cleanedEntries = append(cleanedEntries, entry)
							}
						}
					}
					
					// Only keep this hook group if it has remaining entries
					if len(cleanedEntries) > 0 {
						hookMap["hooks"] = cleanedEntries
						cleaned = append(cleaned, hook)
					}
				} else {
					// Keep hooks that don't match the expected structure
					cleaned = append(cleaned, hook)
				}
			} else {
				// Keep hooks that don't have a "hooks" sub-array
				cleaned = append(cleaned, hook)
			}
		} else {
			// Keep hooks that aren't objects
			cleaned = append(cleaned, hook)
		}
	}

	return cleaned
}

// removeStartupDirectories filters out directories that contain "the-startup"
func (sm *DefaultSettingsManager) removeStartupDirectories(dirs []interface{}) []interface{} {
	var cleaned []interface{}

	for _, dir := range dirs {
		if dirStr, ok := dir.(string); ok {
			// Remove directories that contain "the-startup" or ".the-startup"
			if !strings.Contains(dirStr, "the-startup") {
				cleaned = append(cleaned, dir)
			}
		} else {
			// Keep non-string entries as-is
			cleaned = append(cleaned, dir)
		}
	}

	return cleaned
}

// previewSettingsFileChanges shows what would be removed from a settings file
func (sm *DefaultSettingsManager) previewSettingsFileChanges(settingsPath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	if len(bytes.TrimSpace(data)) == 0 {
		return make(map[string]interface{}), nil
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings JSON: %w", err)
	}

	// Create a copy and track what would be removed
	originalSize := len(settings)
	settingsCopy := make(map[string]interface{})
	for k, v := range settings {
		settingsCopy[k] = v
	}

	sm.removeStartupManagedEntries(settingsCopy)

	preview := make(map[string]interface{})
	preview["original_entries"] = originalSize
	preview["remaining_entries"] = len(settingsCopy)
	preview["entries_to_remove"] = originalSize - len(settingsCopy)

	// Identify specific entries that would be removed
	var removedKeys []string
	for key := range settings {
		if _, exists := settingsCopy[key]; !exists {
			removedKeys = append(removedKeys, key)
		}
	}
	preview["removed_keys"] = removedKeys

	return preview, nil
}

// atomicFileCopy performs an atomic file copy operation
func (sm *DefaultSettingsManager) atomicFileCopy(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create temporary destination file
	tmpDst := dst + ".tmp"
	dstFile, err := os.Create(tmpDst)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer func() {
		dstFile.Close()
		// Clean up temp file if something goes wrong
		if err != nil {
			os.Remove(tmpDst)
		}
	}()

	// Copy the content
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Sync to ensure data is written to disk
	if err = dstFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	// Close the file before renaming
	if err = dstFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Atomically move temp file to final destination
	if err = os.Rename(tmpDst, dst); err != nil {
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// atomicJSONWrite writes JSON data to a file atomically
func (sm *DefaultSettingsManager) atomicJSONWrite(path string, data interface{}) error {
	// Create temporary file
	tmpPath := path + ".tmp"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer func() {
		tmpFile.Close()
		// Clean up temp file if something goes wrong
		if err != nil {
			os.Remove(tmpPath)
		}
	}()

	// Write JSON with proper formatting
	encoder := json.NewEncoder(tmpFile)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	// Sync to ensure data is written to disk
	if err = tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	// Close the file before renaming
	if err = tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Atomically move temp file to final destination
	if err = os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// verifyBackupIntegrity checks that a backup file was created correctly
func (sm *DefaultSettingsManager) verifyBackupIntegrity(originalPath, backupPath string, expectedSize int64) error {
	backupInfo, err := os.Stat(backupPath)
	if err != nil {
		return fmt.Errorf("backup file does not exist: %w", err)
	}

	if backupInfo.Size() != expectedSize {
		return fmt.Errorf("backup file size mismatch: expected %d bytes, got %d bytes", expectedSize, backupInfo.Size())
	}

	// Verify that the backup file contains valid JSON (for .json files)
	if strings.HasSuffix(backupPath, ".json") {
		data, err := os.ReadFile(backupPath)
		if err != nil {
			return fmt.Errorf("failed to read backup file for verification: %w", err)
		}

		var jsonData interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			return fmt.Errorf("backup contains invalid JSON: %w", err)
		}
	}

	return nil
}

// cleanupOldBackups removes old backup files to prevent disk space accumulation
func (sm *DefaultSettingsManager) cleanupOldBackups(backupDir string) error {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}

	// Group backup files by their base name (e.g., "settings.json")
	backupGroups := make(map[string][]os.DirEntry)
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		name := entry.Name()
		if !strings.Contains(name, ".backup.") {
			continue
		}
		
		// Extract base name (e.g., "settings.json.backup.20240101_120000" -> "settings.json")
		baseName := strings.Split(name, ".backup.")[0]
		backupGroups[baseName] = append(backupGroups[baseName], entry)
	}

	// For each group, keep only the most recent backups
	for _, backups := range backupGroups {
		if len(backups) <= sm.maxBackupRetention {
			continue
		}

		// Sort by modification time (newest first)
		sort.Slice(backups, func(i, j int) bool {
			info1, _ := backups[i].Info()
			info2, _ := backups[j].Info()
			return info1.ModTime().After(info2.ModTime())
		})

		// Remove old backups
		for i := sm.maxBackupRetention; i < len(backups); i++ {
			backupPath := filepath.Join(backupDir, backups[i].Name())
			if err := os.Remove(backupPath); err != nil {
				if sm.verbose {
					fmt.Printf("Warning: Failed to remove old backup %s: %v\n", backups[i].Name(), err)
				}
			} else if sm.verbose {
				fmt.Printf("✓ Removed old backup: %s\n", backups[i].Name())
			}
		}
	}

	return nil
}