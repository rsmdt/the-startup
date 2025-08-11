package installer

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/the-startup/the-startup/internal/config"
)

// Installer handles the installation process
type Installer struct {
	agentFiles      *embed.FS
	commandFiles    *embed.FS
	hookFiles       *embed.FS
	templateFiles   *embed.FS
	settingsFile    *embed.FS
	
	installPath    string
	claudePath     string
	tool           string
	components     []string
	selectedFiles  []string // Specific files to install (optional)
	lockFile       *config.LockFile
}

// New creates a new installer
func New(agents, commands, hooks, templates, settings *embed.FS) *Installer {
	homeDir, _ := os.UserHomeDir()
	
	// Default to project-local .the-startup directory
	installPath := ".the-startup"
	
	return &Installer{
		agentFiles:      agents,
		commandFiles:    commands,
		hookFiles:       hooks,
		templateFiles:   templates,
		settingsFile:    settings,
		installPath:     installPath,
		claudePath:      filepath.Join(homeDir, ".claude"),
		tool:            "claude-code",
		components:      []string{"agents", "hooks", "commands", "templates"},
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
	
	// If this is a local installation (.the-startup in current dir),
	// also set Claude path to local .claude
	if strings.Contains(path, ".the-startup") && !strings.Contains(path, ".config") {
		// This appears to be a local installation
		dir := filepath.Dir(path)
		i.claudePath = filepath.Join(dir, ".claude")
	}
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

// checkFileExistsInClaude checks if a file exists specifically in the Claude directory
func (i *Installer) checkFileExistsInClaude(componentPath string) bool {
	parts := strings.Split(componentPath, "/")
	if len(parts) >= 2 {
		component := parts[0]
		fileName := parts[1]
		
		var filePath string
		switch component {
		case "agents":
			filePath = filepath.Join(i.claudePath, "agents", fileName)
		case "commands":
			filePath = filepath.Join(i.claudePath, "commands", fileName)
		case "hooks":
			// Hooks are now the binary in STARTUP_PATH/bin
			if fileName == "the-startup" {
				filePath = filepath.Join(i.installPath, "bin", "the-startup")
			}
		case "templates":
			// Templates are in STARTUP_PATH/templates
			filePath = filepath.Join(i.installPath, "templates", fileName)
		}
		
		if filePath != "" {
			if _, err := os.Stat(filePath); err == nil {
				return true
			}
		}
	}
	
	return false
}

// CheckFileExists checks if a component file exists
func (i *Installer) CheckFileExists(componentPath string) bool {
	parts := strings.Split(componentPath, "/")
	if len(parts) >= 2 {
		component := parts[0]
		fileName := parts[1]
		
		var filePath string
		switch component {
		case "agents":
			filePath = filepath.Join(i.claudePath, "agents", fileName)
		case "commands":
			filePath = filepath.Join(i.claudePath, "commands", fileName)
		case "hooks":
			// Hooks are now the binary in STARTUP_PATH/bin
			if fileName == "the-startup" {
				filePath = filepath.Join(i.installPath, "bin", "the-startup")
			}
		case "templates":
			// Templates are in STARTUP_PATH/templates
			filePath = filepath.Join(i.installPath, "templates", fileName)
		}
		
		if filePath != "" {
			if _, err := os.Stat(filePath); err == nil {
				return true
			}
		}
	}
	
	return false
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

// Install performs the installation
func (i *Installer) Install() error {
	// Create installation directory for lock file
	if err := os.MkdirAll(i.installPath, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}
	
	// Install components directly to Claude directory
	if i.tool == "claude-code" {
		for _, component := range i.components {
			fmt.Printf("Installing %s...\n", component)
			
			if err := i.installComponentToClaude(component); err != nil {
				fmt.Printf("✗ Error installing %s: %v\n", component, err)
				return fmt.Errorf("failed to install %s: %w", component, err)
			}
			
			fmt.Printf("✓ %s installed\n", component)
		}
		
		// Install the binary to STARTUP_PATH/bin
		if contains(i.components, "hooks") {
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

// installComponentToClaude installs a specific component directly to Claude directory
func (i *Installer) installComponentToClaude(component string) error {
	// Special handling for templates - install to STARTUP_PATH instead of CLAUDE_PATH
	if component == "templates" {
		templatePath := filepath.Join(i.installPath, "templates")
		if err := os.MkdirAll(templatePath, 0755); err != nil {
			return fmt.Errorf("failed to create template directory: %w", err)
		}
		
		// Copy template files to STARTUP_PATH/templates
		files, err := fs.Glob(i.templateFiles, "assets/templates/*")
		if err != nil {
			return err
		}
		
		for _, file := range files {
			fileName := filepath.Base(file)
			destPath := filepath.Join(templatePath, fileName)
			
			data, err := i.templateFiles.ReadFile(file)
			if err != nil {
				return err
			}
			
			if err := os.WriteFile(destPath, data, 0644); err != nil {
				return err
			}
		}
		return nil
	}
	
	// Create component directory in Claude path
	componentPath := filepath.Join(i.claudePath, component)
	if err := os.MkdirAll(componentPath, 0755); err != nil {
		return err
	}
	
	// Skip hooks component - it's now handled separately
	if component == "hooks" {
		return nil
	}
	
	// For other components, use the original logic
	var sourceFS *embed.FS
	var pattern string
	
	switch component {
	case "agents":
		sourceFS = i.agentFiles
		pattern = "assets/agents/*.md"
	case "commands":
		sourceFS = i.commandFiles
		pattern = "assets/commands/*.md"
	default:
		return fmt.Errorf("unknown component: %s", component)
	}
	
	// Copy files
	files, err := fs.Glob(sourceFS, pattern)
	if err != nil {
		return err
	}
	
	for _, file := range files {
		// If specific files are selected, check if this file should be installed
		if len(i.selectedFiles) > 0 {
			fileName := filepath.Base(file)
			filePath := component + "/" + fileName
			shouldInstall := false
			for _, selected := range i.selectedFiles {
				if selected == filePath {
					shouldInstall = true
					break
				}
			}
			if !shouldInstall {
				continue // Skip this file
			}
		}
		
		if err := i.copyFile(sourceFS, file, componentPath); err != nil {
			return fmt.Errorf("failed to copy %s: %w", file, err)
		}
	}
	
	return nil
}

// copyFile copies a single file from embed.FS to disk
func (i *Installer) copyFile(sourceFS *embed.FS, sourcePath, destDir string) error {
	// Read source file
	data, err := sourceFS.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	
	// Replace placeholders in .md files
	if strings.HasSuffix(sourcePath, ".md") {
		data = i.replacePlaceholders(data)
	}
	
	// Determine destination path
	fileName := filepath.Base(sourcePath)
	destPath := filepath.Join(destDir, fileName)
	
	// Write file
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return err
	}
	
	// Make Python files executable
	if strings.HasSuffix(destPath, ".py") {
		os.Chmod(destPath, 0755)
	}
	
	return nil
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
	
	// Copy file contents
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy executable: %w", err)
	}
	
	// Make executable
	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permissions: %w", err)
	}
	
	return nil
}

// configureHooks updates settings.json to include hooks and permissions
func (i *Installer) configureHooks() error {
	settingsPath := filepath.Join(i.claudePath, "settings.json")
	
	// Read the template settings
	var templateSettings map[string]interface{}
	if i.settingsFile != nil {
		templateData, err := i.settingsFile.ReadFile("assets/settings.json")
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
	
	return os.WriteFile(settingsPath, data, 0644)
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
	
	// Merge permissions
	if templatePerms, ok := template["permissions"].(map[string]interface{}); ok {
		if existingPerms, ok := existing["permissions"].(map[string]interface{}); ok {
			// Merge additionalDirectories
			if templateDirs, ok := templatePerms["additionalDirectories"].([]interface{}); ok {
				existingPerms["additionalDirectories"] = templateDirs
			}
			result["permissions"] = existingPerms
		} else {
			result["permissions"] = templatePerms
		}
	}
	
	// Replace hooks entirely with our template (we manage these)
	if templateHooks, ok := template["hooks"].(map[string]interface{}); ok {
		result["hooks"] = templateHooks
	}
	
	return result
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
	}
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
	
	// Record installed files from Claude directory
	for _, component := range i.components {
		// Skip hooks component as it's not in Claude directory anymore
		if component == "hooks" {
			// Record the binary instead
			binPath := filepath.Join(i.installPath, "bin", "the-startup")
			if info, err := os.Stat(binPath); err == nil {
				relPath := filepath.Join("bin", "the-startup")
				lockFile.Files[relPath] = config.FileInfo{
					Size:         info.Size(),
					LastModified: info.ModTime().Format(time.RFC3339),
				}
			}
			continue
		}
		
		componentPath := filepath.Join(i.claudePath, component)
		err := filepath.Walk(componentPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			
			relPath, _ := filepath.Rel(i.claudePath, path)
			lockFile.Files[relPath] = config.FileInfo{
				Size:         info.Size(),
				LastModified: info.ModTime().Format(time.RFC3339),
			}
			return nil
		})
		if err != nil {
			// Component might not exist if not selected, that's ok
			if !os.IsNotExist(err) {
				return err
			}
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

// replacePlaceholders replaces template variables with actual paths
func (i *Installer) replacePlaceholders(data []byte) []byte {
	// Replace {{STARTUP_PATH}} with the installation path
	data = bytes.ReplaceAll(data, []byte("{{STARTUP_PATH}}"), []byte(i.installPath))
	
	// Keep {{INSTALL_PATH}} for backward compatibility (same as STARTUP_PATH)
	data = bytes.ReplaceAll(data, []byte("{{INSTALL_PATH}}"), []byte(i.installPath))
	
	// Replace {{CLAUDE_PATH}} with ~/.claude
	data = bytes.ReplaceAll(data, []byte("{{CLAUDE_PATH}}"), []byte(i.claudePath))
	
	return data
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}