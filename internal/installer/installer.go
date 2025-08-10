package installer

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/the-startup/the-startup/internal/config"
)

// Installer handles the installation process
type Installer struct {
	agentFiles    *embed.FS
	commandFiles  *embed.FS
	hookFiles     *embed.FS
	ruleFiles     *embed.FS
	templateFiles *embed.FS
	
	installPath  string
	claudePath   string
	tool         string
	components   []string
	lockFile     *config.LockFile
}

// New creates a new installer
func New(agents, commands, hooks, rules, templates *embed.FS) *Installer {
	homeDir, _ := os.UserHomeDir()
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = filepath.Join(homeDir, ".config")
	}
	
	return &Installer{
		agentFiles:    agents,
		commandFiles:  commands,
		hookFiles:     hooks,
		ruleFiles:     rules,
		templateFiles: templates,
		installPath:   filepath.Join(configDir, "the-startup"),
		claudePath:    filepath.Join(homeDir, ".claude"),
		tool:          "claude-code",
		components:    []string{"agents", "hooks", "commands", "rules"},
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

// SetTool sets the tool type
func (i *Installer) SetTool(tool string) {
	i.tool = tool
}

// SetComponents sets which components to install
func (i *Installer) SetComponents(components []string) {
	i.components = components
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
	// Create installation directory
	if err := os.MkdirAll(i.installPath, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}
	
	// Install each component
	for _, component := range i.components {
		if err := i.installComponent(component); err != nil {
			return fmt.Errorf("failed to install %s: %w", component, err)
		}
	}
	
	// Create references in Claude directory
	if i.tool == "claude-code" {
		if err := i.createClaudeReferences(); err != nil {
			return fmt.Errorf("failed to create Claude references: %w", err)
		}
		
		// Configure hooks
		if contains(i.components, "hooks") {
			if err := i.configureHooks(); err != nil {
				return fmt.Errorf("failed to configure hooks: %w", err)
			}
		}
	}
	
	// Create lock file
	if err := i.createLockFile(); err != nil {
		return fmt.Errorf("failed to create lock file: %w", err)
	}
	
	return nil
}

// installComponent installs a specific component
func (i *Installer) installComponent(component string) error {
	var sourceFS *embed.FS
	var pattern string
	
	switch component {
	case "agents":
		sourceFS = i.agentFiles
		pattern = "agents/*.md"
	case "commands":
		sourceFS = i.commandFiles
		pattern = "commands/*.md"
	case "hooks":
		sourceFS = i.hookFiles
		pattern = "hooks/*.py"
	case "rules":
		sourceFS = i.ruleFiles
		pattern = "rules/*.md"
	case "templates":
		sourceFS = i.templateFiles
		pattern = "templates/*"
	default:
		return fmt.Errorf("unknown component: %s", component)
	}
	
	// Create component directory
	componentPath := filepath.Join(i.installPath, component)
	if err := os.MkdirAll(componentPath, 0755); err != nil {
		return err
	}
	
	// Copy files
	files, err := fs.Glob(sourceFS, pattern)
	if err != nil {
		return err
	}
	
	for _, file := range files {
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

// createClaudeReferences creates reference files in ~/.claude
func (i *Installer) createClaudeReferences() error {
	// Create directories
	claudeAgents := filepath.Join(i.claudePath, "agents")
	claudeCommands := filepath.Join(i.claudePath, "commands")
	claudeHooks := filepath.Join(i.claudePath, "hooks")
	claudeRules := filepath.Join(i.claudePath, "rules")
	
	for _, dir := range []string{claudeAgents, claudeCommands, claudeHooks, claudeRules} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	
	// Create reference files for agents
	if contains(i.components, "agents") {
		agentsDir := filepath.Join(i.installPath, "agents")
		entries, _ := os.ReadDir(agentsDir)
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".md") {
				name := strings.TrimSuffix(entry.Name(), ".md")
				refPath := filepath.Join(claudeAgents, name)
				refContent := fmt.Sprintf("@%s/agents/%s", i.installPath, entry.Name())
				if err := os.WriteFile(refPath, []byte(refContent), 0644); err != nil {
					return err
				}
			}
		}
	}
	
	// Create reference files for commands
	if contains(i.components, "commands") {
		commandsDir := filepath.Join(i.installPath, "commands")
		entries, _ := os.ReadDir(commandsDir)
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".md") {
				name := strings.TrimSuffix(entry.Name(), ".md")
				refPath := filepath.Join(claudeCommands, name+".md")
				
				// Commands are copied directly (not references)
				srcPath := filepath.Join(commandsDir, entry.Name())
				data, err := os.ReadFile(srcPath)
				if err != nil {
					return err
				}
				if err := os.WriteFile(refPath, data, 0644); err != nil {
					return err
				}
			}
		}
	}
	
	// Copy hooks directly (they need to be executable)
	if contains(i.components, "hooks") {
		hooksDir := filepath.Join(i.installPath, "hooks")
		entries, _ := os.ReadDir(hooksDir)
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".py") {
				srcPath := filepath.Join(hooksDir, entry.Name())
				destPath := filepath.Join(claudeHooks, entry.Name())
				
				// Copy file
				data, err := os.ReadFile(srcPath)
				if err != nil {
					return err
				}
				if err := os.WriteFile(destPath, data, 0755); err != nil {
					return err
				}
			}
		}
	}
	
	// Copy rules
	if contains(i.components, "rules") {
		rulesDir := filepath.Join(i.installPath, "rules")
		entries, _ := os.ReadDir(rulesDir)
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".md") {
				srcPath := filepath.Join(rulesDir, entry.Name())
				destPath := filepath.Join(claudeRules, entry.Name())
				
				data, err := os.ReadFile(srcPath)
				if err != nil {
					return err
				}
				if err := os.WriteFile(destPath, data, 0644); err != nil {
					return err
				}
			}
		}
	}
	
	return nil
}

// configureHooks updates settings.json to include hooks
func (i *Installer) configureHooks() error {
	settingsPath := filepath.Join(i.claudePath, "settings.json")
	
	// Read existing settings or create new
	var settings map[string]interface{}
	if data, err := os.ReadFile(settingsPath); err == nil {
		json.Unmarshal(data, &settings)
	} else {
		settings = make(map[string]interface{})
	}
	
	// Ensure hooks section exists
	if _, ok := settings["hooks"]; !ok {
		settings["hooks"] = make(map[string]interface{})
	}
	hooks := settings["hooks"].(map[string]interface{})
	
	// Add PreToolUse hook
	preToolUse := []map[string]interface{}{
		{
			"matcher": "Task",
			"hooks": []map[string]interface{}{
				{
					"type":    "command",
					"command": "uv run $CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_start.py",
					"_source": "the-startup",
				},
			},
		},
	}
	hooks["PreToolUse"] = preToolUse
	
	// Add PostToolUse hook
	postToolUse := []map[string]interface{}{
		{
			"matcher": "Task",
			"hooks": []map[string]interface{}{
				{
					"type":    "command",
					"command": "uv run $CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_complete.py",
					"_source": "the-startup",
				},
			},
		},
	}
	hooks["PostToolUse"] = postToolUse
	
	// Write updated settings
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(settingsPath, data, 0644)
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
	
	// Record installed files
	for _, component := range i.components {
		componentPath := filepath.Join(i.installPath, component)
		err := filepath.Walk(componentPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			
			relPath, _ := filepath.Rel(i.installPath, path)
			lockFile.Files[relPath] = config.FileInfo{
				Size:         info.Size(),
				LastModified: info.ModTime().Format(time.RFC3339),
			}
			return nil
		})
		if err != nil {
			return err
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

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}