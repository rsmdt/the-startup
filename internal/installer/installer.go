package installer

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rsmdt/the-startup/internal/config"
)

// Installer handles the installation process
type Installer struct {
	claudeAssets  *embed.FS
	startupAssets *embed.FS

	installPath     string
	claudePath      string
	tool            string
	components      []string
	selectedFiles   []string // Specific files to install (optional)
	lockFile        *config.LockFile
	existingLock    *config.LockFile // Previously installed files
	deprecatedFiles []string          // Files to be removed
}

// New creates a new installer
func New(claudeAssets, startupAssets *embed.FS) *Installer {
	homeDir, _ := os.UserHomeDir()

	// Default to project-local .the-startup directory
	installPath := ".the-startup"

	return &Installer{
		claudeAssets:  claudeAssets,
		startupAssets: startupAssets,
		installPath:   installPath,
		claudePath:    filepath.Join(homeDir, ".claude"),
		tool:          "claude-code",
		components:    []string{"agents", "commands", "templates", "rules"},
	}
}

// SetInstallPath sets the installation path
func (i *Installer) SetInstallPath(path string) {
	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		path = filepath.Join(homeDir, path[2:])
	}
	i.installPath = path
}

// SetClaudePath sets the Claude directory path
func (i *Installer) SetClaudePath(path string) {
	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		path = filepath.Join(homeDir, path[2:])
	}
	i.claudePath = path
}

// GetInstallPath returns the installation path
func (i *Installer) GetInstallPath() string {
	return i.installPath
}

// GetClaudePath returns the claude path
func (i *Installer) GetClaudePath() string {
	return i.claudePath
}

// GetExistingFiles returns a list of files that already exist and will be updated
func (i *Installer) GetExistingFiles(files []string) []string {
	var existing []string
	for _, file := range files {
		// Only check if file exists in the actual Claude directory where it will be installed
		if i.checkFileExistsInClaude(file) {
			existing = append(existing, file)
		}
	}
	return existing
}

// CheckSettingsExists checks if settings.json exists in the Claude directory
func (i *Installer) CheckSettingsExists() bool {
	settingsPath := filepath.Join(i.claudePath, "settings.json")
	if _, err := os.Stat(settingsPath); err == nil {
		return true
	}
	return false
}

// checkFileExistsInClaude checks if a file exists specifically in the Claude directory
func (i *Installer) checkFileExistsInClaude(componentPath string) bool {
	// Dynamically determine destination path based on component type
	// This mirrors the exact structure: assets/X/... -> claudePath/X/...

	parts := strings.Split(componentPath, "/")
	if len(parts) < 2 {
		return false
	}

	component := parts[0]

	// Build the rest of the path after the component
	restPath := strings.Join(parts[1:], "/")

	var filePath string
	switch component {
	case "agents", "commands":
		// These go to .claude directory, preserving full path structure
		filePath = filepath.Join(i.claudePath, component, restPath)
	case "hooks":
		// Hooks are now the binary in STARTUP_PATH/bin
		// This shouldn't be checked as part of .claude files
		return false
	case "templates":
		// Templates go to STARTUP_PATH/templates
		// This shouldn't be checked as part of .claude files
		return false
	case "rules":
		// Rules go to STARTUP_PATH/rules
		// This shouldn't be checked as part of .claude files
		return false
	default:
		// Unknown component type
		return false
	}

	// Check if the file exists
	if filePath != "" {
		if _, err := os.Stat(filePath); err == nil {
			return true
		}
	}

	return false
}

// CheckFileExists checks if a component file exists
func (i *Installer) CheckFileExists(componentPath string) bool {
	// Use the same logic as checkFileExistsInClaude for consistency
	return i.checkFileExistsInClaude(componentPath)
}

// SetTool sets the tool type
func (i *Installer) SetTool(tool string) {
	i.tool = tool
}

// SetComponents sets which components to install
func (i *Installer) SetComponents(components []string) {
	i.components = components
}

// SetSelectedFiles sets specific files to install
func (i *Installer) SetSelectedFiles(files []string) {
	i.selectedFiles = files
}

// IsInstalled checks if already installed
func (i *Installer) IsInstalled() bool {
	lockFilePath := filepath.Join(i.installPath, "the-startup.lock")
	if _, err := os.Stat(lockFilePath); err == nil {
		return true
	}

	// Also check if any claude files exist
	claudeAgentsPath := filepath.Join(i.claudePath, "agents")
	if _, err := os.Stat(claudeAgentsPath); err == nil {
		// Check if any of our agents exist
		entries, _ := os.ReadDir(claudeAgentsPath)
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "the-") {
				return true
			}
		}
	}

	return false
}

// LoadExistingLockFile loads the previous installation's lock file if it exists
func (i *Installer) LoadExistingLockFile() error {
	lockFilePath := filepath.Join(i.installPath, "the-startup.lock")
	
	// Check if lock file exists
	data, err := os.ReadFile(lockFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// No existing installation, this is fine
			return nil
		}
		return fmt.Errorf("failed to read lock file: %w", err)
	}
	
	// Parse the lock file
	var lockFile config.LockFile
	if err := json.Unmarshal(data, &lockFile); err != nil {
		return fmt.Errorf("failed to parse lock file: %w", err)
	}
	
	i.existingLock = &lockFile
	return nil
}

// GetDeprecatedFiles returns a list of files that were previously installed but are no longer in the embedded assets
func (i *Installer) GetDeprecatedFiles() []string {
	if i.existingLock == nil {
		return nil
	}
	
	// Build a set of current files from embedded assets
	currentFiles := make(map[string]bool)
	
	// Add files from claude assets
	if i.claudeAssets != nil {
		fs.WalkDir(i.claudeAssets, "assets/claude", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			// Get relative path that would be in lock file
			relPath := strings.TrimPrefix(path, "assets/claude/")
			currentFiles[relPath] = true
			return nil
		})
	}
	
	// Check which files from the lock file no longer exist in current assets
	// The lock file only contains files WE installed, so anything in it is ours to manage
	var deprecated []string
	for filePath := range i.existingLock.Files {
		// Skip non-claude files (like bin/the-startup, templates, etc)
		if !strings.HasPrefix(filePath, "agents/") && !strings.HasPrefix(filePath, "commands/") {
			continue
		}
		
		// Check if this file still exists in current assets
		if !currentFiles[filePath] {
			deprecated = append(deprecated, filePath)
		}
	}
	
	i.deprecatedFiles = deprecated
	return deprecated
}

// RemoveDeprecatedFiles removes files that are no longer part of the distribution
func (i *Installer) RemoveDeprecatedFiles() error {
	if len(i.deprecatedFiles) == 0 {
		return nil
	}
	
	fmt.Printf("Removing %d deprecated files...\n", len(i.deprecatedFiles))
	
	for _, relPath := range i.deprecatedFiles {
		// Construct full path
		fullPath := filepath.Join(i.claudePath, relPath)
		
		// Remove the file
		if err := os.Remove(fullPath); err != nil {
			// File might already be gone, that's ok
			if !os.IsNotExist(err) {
				fmt.Printf("  Warning: Failed to remove %s: %v\n", relPath, err)
			}
		} else {
			fmt.Printf("  ✗ Removed: %s\n", relPath)
		}
	}
	
	fmt.Printf("✓ Deprecated files removed\n")
	return nil
}

// Install performs the installation
func (i *Installer) Install() error {
	// Create installation directory for lock file
	if err := os.MkdirAll(i.installPath, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Load existing lock file to detect deprecated files
	if err := i.LoadExistingLockFile(); err != nil {
		fmt.Printf("Warning: Failed to load existing lock file: %v\n", err)
		// Continue anyway, this just means we won't remove deprecated files
	}

	// Detect deprecated files
	if i.existingLock != nil {
		i.GetDeprecatedFiles()
	}

	// Remove deprecated files first (before installing new ones)
	if len(i.deprecatedFiles) > 0 {
		if err := i.RemoveDeprecatedFiles(); err != nil {
			fmt.Printf("Warning: Failed to remove some deprecated files: %v\n", err)
			// Continue with installation
		}
	}

	// Install components directly to Claude directory
	if i.tool == "claude-code" {
		// Install Claude assets to CLAUDE_PATH
		fmt.Println("Installing Claude assets...")
		if err := i.installClaudeAssets(); err != nil {
			fmt.Printf("✗ Error installing Claude assets: %v\n", err)
			return fmt.Errorf("failed to install Claude assets: %w", err)
		}
		fmt.Println("✓ Claude assets installed")

		// Install Startup assets to STARTUP_PATH
		fmt.Println("Installing Startup assets...")
		if err := i.installStartupAssets(); err != nil {
			fmt.Printf("✗ Error installing Startup assets: %v\n", err)
			return fmt.Errorf("failed to install Startup assets: %w", err)
		}
		fmt.Println("✓ Startup assets installed")

		// Always install the binary to STARTUP_PATH/bin and configure hooks
		fmt.Println("Installing the-startup binary...")
		if err := i.installBinary(); err != nil {
			fmt.Printf("✗ Error installing binary: %v\n", err)
			return fmt.Errorf("failed to install binary: %w", err)
		}
		fmt.Printf("✓ Binary installed to %s/bin\n", i.installPath)

		// Configure hooks in settings.json
		if err := i.configureHooks(); err != nil {
			return fmt.Errorf("failed to configure hooks: %w", err)
		}
	}

	// Create lock file
	fmt.Println("Creating lock file...")

	if err := i.createLockFile(); err != nil {
		fmt.Printf("✗ Error creating lock file: %v\n", err)
		return fmt.Errorf("failed to create lock file: %w", err)
	}

	fmt.Println("✓ Lock file created")

	return nil
}

// installClaudeAssets installs all assets from the claude directory to CLAUDE_PATH
func (i *Installer) installClaudeAssets() error {
	// Walk through all files in assets/claude and copy them to CLAUDE_PATH
	return fs.WalkDir(i.claudeAssets, "assets/claude", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Get relative path from assets/claude
		relPath := strings.TrimPrefix(path, "assets/claude/")

		// Check if this file should be installed based on selected files
		if i.selectedFiles != nil && len(i.selectedFiles) > 0 {
			found := false
			for _, selected := range i.selectedFiles {
				if selected == relPath {
					found = true
					break
				}
			}
			if !found {
				return nil // Skip this file
			}
		}

		// Create destination path in CLAUDE_PATH
		destPath := filepath.Join(i.claudePath, relPath)

		// Create destination directory if needed
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}

		// Read file content
		data, err := i.claudeAssets.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// Apply placeholder replacement to ALL files
		data = i.replacePlaceholders(data)

		// Skip writing settings.json and settings.local.json here as they're handled by configureHooks
		if filepath.Base(path) == "settings.json" || filepath.Base(path) == "settings.local.json" {
			return nil
		}

		// Write file to destination
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", destPath, err)
		}

		return nil
	})
}

// installStartupAssets installs all assets from the-startup directory to STARTUP_PATH
func (i *Installer) installStartupAssets() error {
	// Walk through all files in assets/the-startup and copy them to STARTUP_PATH
	return fs.WalkDir(i.startupAssets, "assets/the-startup", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Get relative path from assets/the-startup
		relPath := strings.TrimPrefix(path, "assets/the-startup/")

		// Check if this file should be installed based on selected files
		if i.selectedFiles != nil && len(i.selectedFiles) > 0 {
			found := false
			for _, selected := range i.selectedFiles {
				if selected == relPath {
					found = true
					break
				}
			}
			if !found {
				return nil // Skip this file
			}
		}

		// Create destination path in STARTUP_PATH
		destPath := filepath.Join(i.installPath, relPath)

		// Create destination directory if needed
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}

		// Read file content
		data, err := i.startupAssets.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// Apply placeholder replacement for template files
		data = i.replacePlaceholders(data)

		// Write file to destination
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", destPath, err)
		}

		return nil
	})
}

// configureHooks updates settings.json and settings.local.json to include hooks and permissions
func (i *Installer) configureHooks() error {
	// Handle settings.json
	settingsPath := filepath.Join(i.claudePath, "settings.json")

	// Read the template settings
	var templateSettings map[string]interface{}
	if i.claudeAssets != nil {
		templateData, err := i.claudeAssets.ReadFile("assets/claude/settings.json")
		if err == nil {
			// Replace placeholders in template
			templateData = i.replacePlaceholders(templateData)
			json.Unmarshal(templateData, &templateSettings)
		}
	}

	// If template loading failed, fall back to hardcoded
	if templateSettings == nil {
		templateSettings = i.createDefaultSettings()
	}

	// Read existing settings if present
	var existingSettings map[string]interface{}
	if data, err := os.ReadFile(settingsPath); err == nil {
		json.Unmarshal(data, &existingSettings)
	}

	// Merge settings (template takes precedence for our managed sections)
	settings := i.mergeSettings(existingSettings, templateSettings)

	// Write updated settings
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return err
	}

	// Handle settings.local.json if template exists
	localSettingsPath := filepath.Join(i.claudePath, "settings.local.json")

	if i.claudeAssets != nil {
		if templateLocalData, err := i.claudeAssets.ReadFile("assets/claude/settings.local.json"); err == nil {
			// Replace placeholders in template
			templateLocalData = i.replacePlaceholders(templateLocalData)

			var templateLocalSettings map[string]interface{}
			if err := json.Unmarshal(templateLocalData, &templateLocalSettings); err != nil {
				// Log the error but continue - invalid JSON in template
				fmt.Printf("Warning: Failed to parse settings.local.json template: %v\n", err)
			} else {
				// Read existing local settings if present
				var existingLocalSettings map[string]interface{}
				if data, err := os.ReadFile(localSettingsPath); err == nil {
					if err := json.Unmarshal(data, &existingLocalSettings); err != nil {
						fmt.Printf("Warning: Failed to parse existing settings.local.json: %v\n", err)
						// Continue with just the template settings
					}
				}

				// Merge local settings (template takes precedence for our managed sections)
				localSettings := i.mergeSettings(existingLocalSettings, templateLocalSettings)

				// Write updated local settings
				localData, err := json.MarshalIndent(localSettings, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal settings.local.json: %w", err)
				}

				if err := os.WriteFile(localSettingsPath, localData, 0644); err != nil {
					return fmt.Errorf("failed to write settings.local.json: %w", err)
				}

				fmt.Println("✓ Settings.local.json configured")
			}
		}
		// If settings.local.json doesn't exist in assets, that's OK - it's optional
	}

	return nil
}

// mergeSettings merges template settings with existing settings
func (i *Installer) mergeSettings(existing, template map[string]interface{}) map[string]interface{} {
	if existing == nil {
		return template
	}
	if template == nil {
		return existing
	}

	result := make(map[string]interface{})

	// Copy all existing settings first
	for k, v := range existing {
		result[k] = v
	}

	// Now merge template settings into result
	for key, templateValue := range template {
		existingValue, exists := result[key]

		if !exists {
			// Key doesn't exist in existing settings, add it
			result[key] = templateValue
		} else {
			// Key exists, merge based on type
			result[key] = i.mergeValues(existingValue, templateValue)
		}
	}

	return result
}

// mergeValues recursively merges two values
func (i *Installer) mergeValues(existing, template interface{}) interface{} {
	// If both are maps, merge them recursively
	if existingMap, existingIsMap := existing.(map[string]interface{}); existingIsMap {
		if templateMap, templateIsMap := template.(map[string]interface{}); templateIsMap {
			merged := make(map[string]interface{})

			// Copy all existing entries
			for k, v := range existingMap {
				merged[k] = v
			}

			// Merge in template entries
			for k, templateVal := range templateMap {
				if existingVal, exists := merged[k]; exists {
					// Special handling for specific keys
					if k == "hooks" || k == "additionalDirectories" {
						// For these arrays, we want to deduplicate
						merged[k] = i.mergeAndDeduplicate(existingVal, templateVal)
					} else {
						// Recursively merge for other keys
						merged[k] = i.mergeValues(existingVal, templateVal)
					}
				} else {
					// Add new key from template
					merged[k] = templateVal
				}
			}

			return merged
		}
	}

	// If both are slices, merge them
	if existingSlice, existingIsSlice := existing.([]interface{}); existingIsSlice {
		if templateSlice, templateIsSlice := template.([]interface{}); templateIsSlice {
			return i.mergeSlices(existingSlice, templateSlice)
		}
	}

	// For other types or mismatched types, template takes precedence
	// This includes strings, numbers, booleans
	return template
}

// mergeAndDeduplicate merges arrays and removes duplicates
func (i *Installer) mergeAndDeduplicate(existing, template interface{}) interface{} {
	// Handle hooks array (array of hook entries)
	if existingSlice, existingIsSlice := existing.([]interface{}); existingIsSlice {
		if templateSlice, templateIsSlice := template.([]interface{}); templateIsSlice {
			// Check if this is a hooks array (has command field)
			if len(templateSlice) > 0 {
				if hookEntry, ok := templateSlice[0].(map[string]interface{}); ok {
					if _, hasCommand := hookEntry["command"]; hasCommand {
						// This is a hooks array, deduplicate by command
						seen := make(map[string]bool)
						var result []interface{}

						// Use template hooks (they have the right paths)
						for _, item := range templateSlice {
							if hook, ok := item.(map[string]interface{}); ok {
								if cmd, hasCmd := hook["command"].(string); hasCmd {
									if !seen[cmd] {
										seen[cmd] = true
										result = append(result, item)
									}
								}
							}
						}

						// Add any existing hooks that aren't in template
						for _, item := range existingSlice {
							if hook, ok := item.(map[string]interface{}); ok {
								if cmd, hasCmd := hook["command"].(string); hasCmd {
									// Only add if it's not a startup command (we manage those)
									if !strings.Contains(cmd, "the-startup") && !seen[cmd] {
										seen[cmd] = true
										result = append(result, item)
									}
								}
							}
						}

						return result
					}
				}
			}

			// For additionalDirectories or other string arrays
			if len(templateSlice) > 0 {
				if _, isString := templateSlice[0].(string); isString {
					seen := make(map[string]bool)
					var result []interface{}

					// Use template directories (they have the right paths)
					for _, item := range templateSlice {
						if str, ok := item.(string); ok {
							if !seen[str] {
								seen[str] = true
								result = append(result, str)
							}
						}
					}

					// Add any existing directories that aren't startup-related
					for _, item := range existingSlice {
						if str, ok := item.(string); ok {
							// Only add if it's not a startup directory
							if !strings.Contains(str, "the-startup") && !seen[str] {
								seen[str] = true
								result = append(result, str)
							}
						}
					}

					return result
				}
			}
		}
	}

	// Fallback to template
	return template
}

// mergeSlices merges two slices, typically used for hooks
func (i *Installer) mergeSlices(existing, template []interface{}) []interface{} {
	// For hooks arrays, we want to merge based on matcher
	// Check if these are hook entries (have "matcher" field)
	isHookArray := false
	if len(template) > 0 {
		if hookEntry, ok := template[0].(map[string]interface{}); ok {
			if _, hasMatcher := hookEntry["matcher"]; hasMatcher {
				isHookArray = true
			}
		}
	}

	if isHookArray {
		// Create a map of existing hooks by matcher
		existingByMatcher := make(map[string]interface{})
		for _, item := range existing {
			if hookEntry, ok := item.(map[string]interface{}); ok {
				if matcher, hasMatcher := hookEntry["matcher"].(string); hasMatcher {
					existingByMatcher[matcher] = item
				}
			}
		}

		// Merge template hooks
		var result []interface{}
		for _, templateItem := range template {
			if templateHook, ok := templateItem.(map[string]interface{}); ok {
				if matcher, hasMatcher := templateHook["matcher"].(string); hasMatcher {
					if existingHook, exists := existingByMatcher[matcher]; exists {
						// Merge the hooks for this matcher
						merged := i.mergeValues(existingHook, templateItem)
						result = append(result, merged)
						delete(existingByMatcher, matcher)
					} else {
						// New matcher from template
						result = append(result, templateItem)
					}
				}
			}
		}

		// Add remaining existing hooks that weren't in template
		for _, item := range existingByMatcher {
			result = append(result, item)
		}

		return result
	}

	// For non-hook arrays, just return template
	// (we handle deduplication in mergeAndDeduplicate)
	return template
}

// createDefaultSettings creates default settings if template is not available
func (i *Installer) createDefaultSettings() map[string]interface{} {
	// Use the binary path in STARTUP_PATH/bin
	binaryPath := filepath.Join(i.installPath, "bin", "the-startup")
	return map[string]interface{}{
		"permissions": map[string]interface{}{
			"additionalDirectories": []string{i.installPath},
		},
		"hooks": map[string]interface{}{
			"PreToolUse": []map[string]interface{}{
				{
					"matcher": "Task",
					"hooks": []map[string]interface{}{
						{
							"type":    "command",
							"command": fmt.Sprintf("%s log --assistant", binaryPath),
							"_source": "the-startup",
						},
					},
				},
			},
			"PostToolUse": []map[string]interface{}{
				{
					"matcher": "Task",
					"hooks": []map[string]interface{}{
						{
							"type":    "command",
							"command": fmt.Sprintf("%s log --user", binaryPath),
							"_source": "the-startup",
						},
					},
				},
			},
		},
		"statusLine": map[string]interface{}{
			"type":    "command",
			"command": fmt.Sprintf("%s statusline", binaryPath),
		},
	}
}

// calculateFileChecksum computes the SHA256 checksum of a file
func (i *Installer) calculateFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// createLockFile creates the lock file
func (i *Installer) createLockFile() error {
	lockFile := &config.LockFile{
		Version:     "1.0.0",
		InstallDate: time.Now().Format(time.RFC3339),
		InstallPath: i.installPath,
		ClaudePath:  i.claudePath,
		Tool:        i.tool,
		Components:  i.components,
		Files:       make(map[string]config.FileInfo),
	}

	// Record only the files WE installed from our embedded assets
	// Walk through Claude assets and record what we actually installed
	if i.claudeAssets != nil {
		fs.WalkDir(i.claudeAssets, "assets/claude", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			
			// Get relative path from assets/claude
			relPath := strings.TrimPrefix(path, "assets/claude/")
			
			// Check if this file should be installed based on selected files
			if i.selectedFiles != nil && len(i.selectedFiles) > 0 {
				found := false
				for _, selected := range i.selectedFiles {
					if selected == relPath {
						found = true
						break
					}
				}
				if !found {
					return nil // Skip this file
				}
			}
			
			// Skip settings files as they're handled specially
			if filepath.Base(path) == "settings.json" || filepath.Base(path) == "settings.local.json" {
				return nil
			}
			
			// Check if the file exists in the destination
			destPath := filepath.Join(i.claudePath, relPath)
			if info, err := os.Stat(destPath); err == nil {
				// Calculate checksum for installed file
				checksum, checksumErr := i.calculateFileChecksum(destPath)
				if checksumErr != nil {
					// Log warning but don't fail - checksum is for integrity checking
					fmt.Printf("Warning: Failed to calculate checksum for %s: %v\n", relPath, checksumErr)
					checksum = "" // Empty checksum means no verification possible
				}
				
				lockFile.Files[relPath] = config.FileInfo{
					Size:         info.Size(),
					LastModified: info.ModTime().Format(time.RFC3339),
					Checksum:     checksum,
				}
			}
			
			return nil
		})
	}
	
	// Also record files from startup assets that go to STARTUP_PATH
	if i.startupAssets != nil {
		fs.WalkDir(i.startupAssets, "assets/the-startup", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			
			// Get relative path from assets/the-startup
			relPath := strings.TrimPrefix(path, "assets/the-startup/")
			
			// Check if this file should be installed based on selected files
			if i.selectedFiles != nil && len(i.selectedFiles) > 0 {
				found := false
				for _, selected := range i.selectedFiles {
					if selected == relPath {
						found = true
						break
					}
				}
				if !found {
					return nil // Skip this file
				}
			}
			
			// Check if the file exists in the destination
			destPath := filepath.Join(i.installPath, relPath)
			if info, err := os.Stat(destPath); err == nil {
				// Calculate checksum for installed file
				checksum, checksumErr := i.calculateFileChecksum(destPath)
				if checksumErr != nil {
					// Log warning but don't fail - checksum is for integrity checking
					fmt.Printf("Warning: Failed to calculate checksum for startup/%s: %v\n", relPath, checksumErr)
					checksum = "" // Empty checksum means no verification possible
				}
				
				// Store with a prefix to distinguish from Claude files
				lockFile.Files["startup/"+relPath] = config.FileInfo{
					Size:         info.Size(),
					LastModified: info.ModTime().Format(time.RFC3339),
					Checksum:     checksum,
				}
			}
			
			return nil
		})
	}

	// Record the binary in STARTUP_PATH/bin
	binPath := filepath.Join(i.installPath, "bin", "the-startup")
	if info, err := os.Stat(binPath); err == nil {
		// Calculate checksum for binary
		checksum, checksumErr := i.calculateFileChecksum(binPath)
		if checksumErr != nil {
			// Log warning but don't fail - checksum is for integrity checking
			fmt.Printf("Warning: Failed to calculate checksum for binary: %v\n", checksumErr)
			checksum = "" // Empty checksum means no verification possible
		}
		
		relPath := filepath.Join("bin", "the-startup")
		lockFile.Files[relPath] = config.FileInfo{
			Size:         info.Size(),
			LastModified: info.ModTime().Format(time.RFC3339),
			Checksum:     checksum,
		}
	}

	// Write lock file
	lockFilePath := filepath.Join(i.installPath, "the-startup.lock")
	data, err := json.MarshalIndent(lockFile, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(lockFilePath, data, 0644)
}

// installBinary installs the executable to STARTUP_PATH/bin
func (i *Installer) installBinary() error {
	binDir := filepath.Join(i.installPath, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	return i.copyCurrentExecutable(binDir)
}

// copyCurrentExecutable copies the current executable to the destination directory
func (i *Installer) copyCurrentExecutable(destDir string) error {
	// Get the path of the current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Open source file
	src, err := os.Open(execPath)
	if err != nil {
		return fmt.Errorf("failed to open source executable: %w", err)
	}
	defer src.Close()

	// Create destination path
	destPath := filepath.Join(destDir, "the-startup")

	// Create destination file
	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the file
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy executable: %w", err)
	}

	// Make it executable
	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("failed to make file executable: %w", err)
	}

	return nil
}

// GetInstalledAgents returns the list of installed agent files
func (i *Installer) GetInstalledAgents() []string {
	var agents []string

	// If specific files were selected, filter for agents
	if len(i.selectedFiles) > 0 {
		for _, file := range i.selectedFiles {
			if strings.HasPrefix(file, "agents/") && strings.HasSuffix(file, ".md") {
				// Extract agent name from path (e.g., "agents/the-architect.md" -> "the-architect")
				name := filepath.Base(file)
				name = strings.TrimSuffix(name, ".md")
				agents = append(agents, name)
			}
		}
	} else {
		// Get all agent files from embedded FS
		pattern := "assets/claude/agents/*.md"
		files, _ := fs.Glob(i.claudeAssets, pattern)
		for _, file := range files {
			name := filepath.Base(file)
			name = strings.TrimSuffix(name, ".md")
			agents = append(agents, name)
		}
	}

	return agents
}

// GetDeprecatedFilesList returns the list of deprecated files (for UI display)
func (i *Installer) GetDeprecatedFilesList() []string {
	return i.deprecatedFiles
}

// GetInstalledCommands returns the list of installed command files
func (i *Installer) GetInstalledCommands() []string {
	var commands []string

	// If specific files were selected, filter for commands
	if len(i.selectedFiles) > 0 {
		for _, file := range i.selectedFiles {
			if strings.HasPrefix(file, "commands/") && strings.HasSuffix(file, ".md") {
				// Extract command path (e.g., "commands/s/specify.md" -> "/s:specify")
				relPath := strings.TrimPrefix(file, "commands/")
				relPath = strings.TrimSuffix(relPath, ".md")
				// Convert path separators to colons
				parts := strings.Split(relPath, "/")
				command := "/" + strings.Join(parts, ":")
				commands = append(commands, command)
			}
		}
	} else {
		// Get all command files from embedded FS
		pattern := "assets/claude/commands/**/*.md"
		files, _ := fs.Glob(i.claudeAssets, pattern)
		for _, file := range files {
			// Extract command path (e.g., "assets/claude/commands/s/specify.md" -> "/s:specify")
			relPath := strings.TrimPrefix(file, "assets/claude/commands/")
			relPath = strings.TrimSuffix(relPath, ".md")
			// Convert path separators to colons
			parts := strings.Split(relPath, "/")
			command := "/" + strings.Join(parts, ":")
			commands = append(commands, command)
		}
	}

	return commands
}

// toTildePath converts an absolute path to use ~ for the home directory
func toTildePath(path string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path // Return original if we can't get home dir
	}

	// If path starts with home directory, replace it with ~
	if strings.HasPrefix(path, homeDir) {
		return "~" + strings.TrimPrefix(path, homeDir)
	}

	return path
}

// replacePlaceholders replaces template variables with actual paths.
// This function is called for ALL asset files (agents, commands, templates, rules)
// during installation to ensure placeholders work across all content types.
//
// Supported placeholders:
//   - {{STARTUP_PATH}}: The user-selected installation directory (e.g., ~/.the-startup)
//   - {{CLAUDE_PATH}}: The Claude configuration directory (e.g., ~/.claude)
//
// This generic approach ensures that any future agents, commands, or other assets
// can use these placeholders without requiring code changes.
func (i *Installer) replacePlaceholders(data []byte) []byte {
	// Convert paths to ~ format for better readability
	startupPath := toTildePath(i.installPath)
	claudePath := toTildePath(i.claudePath)

	// Replace {{STARTUP_PATH}} with the installation path (using ~ format)
	data = bytes.ReplaceAll(data, []byte("{{STARTUP_PATH}}"), []byte(startupPath))

	// Replace {{CLAUDE_PATH}} with ~/.claude (using ~ format)
	data = bytes.ReplaceAll(data, []byte("{{CLAUDE_PATH}}"), []byte(claudePath))

	// Future placeholders can be added here without modifying the calling code
	// Example: data = bytes.ReplaceAll(data, []byte("{{USER_HOME}}"), []byte("~"))

	return data
}
