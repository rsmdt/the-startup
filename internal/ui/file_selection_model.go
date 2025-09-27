package ui

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/rsmdt/the-startup/internal/installer"
)

type FileSelectionModel struct {
	styles        Styles
	renderer      *ProgressiveDisclosureRenderer
	installer     *installer.Installer
	claudeAssets  *embed.FS
	startupAssets *embed.FS
	selectedTool  string
	selectedPath  string
	selectedFiles []string
	cursor        int
	choices       []string
	ready         bool
	confirmed     bool
	mode          OperationMode
}

func NewFileSelectionModel(selectedTool, selectedPath string, installer *installer.Installer, claudeAssets, startupAssets *embed.FS) FileSelectionModel {
	return NewFileSelectionModelWithMode(selectedTool, selectedPath, installer, claudeAssets, startupAssets, ModeInstall)
}

func NewFileSelectionModelWithMode(selectedTool, selectedPath string, installer *installer.Installer, claudeAssets, startupAssets *embed.FS, mode OperationMode) FileSelectionModel {
	var choices []string
	if mode == ModeUninstall {
		choices = []string{
			"Yes, remove everything",
			"No, keep everything as-is",
		}
	} else {
		choices = []string{
			"Yes, give me awesome",
			"Huh? I did not sign up for this",
		}
	}

	m := FileSelectionModel{
		styles:        GetStyles(),
		renderer:      NewProgressiveDisclosureRenderer(),
		installer:     installer,
		claudeAssets:  claudeAssets,
		startupAssets: startupAssets,
		selectedTool:  selectedTool,
		selectedPath:  selectedPath,
		choices:       choices,
		cursor:        0,
		ready:         false,
		mode:          mode,
	}

	if mode == ModeUninstall {
		// Load files from lockfile for uninstall
		m.selectedFiles = m.getFilesFromLockfile()
		// Update choices if no files found
		if len(m.selectedFiles) == 0 {
			m.choices = []string{"Go back to path selection"}
		}
	} else {
		m.selectedFiles = m.getAllAvailableFiles()
		if m.installer != nil {
			m.installer.SetSelectedFiles(m.selectedFiles)
			// Load existing lock file to detect deprecated files
			m.installer.LoadExistingLockFile()
			m.installer.GetDeprecatedFiles()
		}
	}

	return m
}

func (m FileSelectionModel) Init() tea.Cmd {
	return nil
}

func (m FileSelectionModel) Update(msg tea.Msg) (FileSelectionModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor < len(m.choices) {
				choice := m.choices[m.cursor]
				if (m.mode == ModeUninstall && choice == "Yes, remove everything") ||
				   (m.mode == ModeInstall && choice == "Yes, give me awesome") {
					m.confirmed = true
					m.ready = true
				} else if m.mode == ModeUninstall && choice == "Go back to path selection" {
					// Special case for when no files found - signal to go back
					m.confirmed = false
					m.ready = true
				} else {
					m.confirmed = false
					m.ready = true
				}
			}
		}
	}
	return m, nil
}

func (m FileSelectionModel) View() string {
	var s strings.Builder

	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")

	// Show both paths in the header
	startupPath := m.installer.GetInstallPath()
	claudePath := m.installer.GetClaudePath()

	// Format paths for display
	home := os.Getenv("HOME")
	if home != "" {
		if strings.HasPrefix(startupPath, home) {
			startupPath = "~" + strings.TrimPrefix(startupPath, home)
		}
		if strings.HasPrefix(claudePath, home) {
			claudePath = "~" + strings.TrimPrefix(claudePath, home)
		}
	}

	if m.mode == ModeUninstall {
		s.WriteString(m.styles.Warning.Render("Uninstallation Paths:"))
		s.WriteString("\n")
		s.WriteString(m.styles.Normal.Render(fmt.Sprintf("  Startup: %s", startupPath)))
		s.WriteString("\n")
		s.WriteString(m.styles.Normal.Render(fmt.Sprintf("  Claude:  %s", claudePath)))
		s.WriteString("\n\n")

		// Check if we have any files to remove
		if len(m.selectedFiles) == 0 {
			s.WriteString(m.renderer.RenderTitle("No Installation Found"))
			s.WriteString(m.styles.Warning.Render("No installation found - nothing to uninstall"))
			s.WriteString("\n\n")
			s.WriteString(m.styles.Help.Render("The lockfile does not exist or no tracked files were found."))
			s.WriteString("\n")
			s.WriteString(m.styles.Help.Render("This could mean:"))
			s.WriteString("\n")
			s.WriteString(m.styles.Help.Render("  • The (Agentic) Startup is not installed in this location"))
			s.WriteString("\n")
			s.WriteString(m.styles.Help.Render("  • Files have already been manually removed"))
			s.WriteString("\n")
			s.WriteString(m.styles.Help.Render("  • Installation was corrupted"))
			s.WriteString("\n\n")
			s.WriteString(m.styles.Info.Render("Please select different paths or use install mode to set up The (Agentic) Startup."))
			s.WriteString("\n\n")
		} else {
			s.WriteString(m.renderer.RenderTitle("Files to be removed"))
			s.WriteString(m.styles.Warning.Render(fmt.Sprintf("%d files will be removed:", len(m.selectedFiles))))
			s.WriteString("\n\n")

			s.WriteString(m.buildStaticTree())
			s.WriteString("\n")
		}

		s.WriteString("\n\n")
		s.WriteString(m.styles.Title.Render("Ready to uninstall?"))
		s.WriteString("\n")
		s.WriteString(m.styles.Warning.Render("This will remove The (Agentic) Startup from the selected directories."))
		s.WriteString("\n\n")
	} else {
		s.WriteString(m.styles.Info.Render("Installation Paths:"))
		s.WriteString("\n")
		s.WriteString(m.styles.Normal.Render(fmt.Sprintf("  Startup: %s", startupPath)))
		s.WriteString("\n")
		s.WriteString(m.styles.Normal.Render(fmt.Sprintf("  Claude:  %s", claudePath)))
		s.WriteString("\n\n")

		s.WriteString(m.renderer.RenderTitle("Files to be installed to .claude"))

		s.WriteString(m.styles.Info.Render("The following files will be installed to your Claude directory:"))
		s.WriteString("\n\n")

		s.WriteString(m.buildStaticTree())
		s.WriteString("\n")

		s.WriteString("\n\n")
		s.WriteString(m.styles.Title.Render("Ready to install?"))
		s.WriteString("\n")
		s.WriteString(m.styles.Info.Render("This will install The (Agentic) Startup to the selected directories."))
		s.WriteString("\n\n")
	}

	for i, option := range m.choices {
		if i == m.cursor {
			s.WriteString(m.styles.Selected.Render("> " + option))
		} else {
			s.WriteString(m.styles.Normal.Render("  " + option))
		}
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(m.styles.Help.Render("Press Enter to confirm • Escape to go back"))

	return s.String()
}

func (m FileSelectionModel) Ready() bool {
	return m.ready
}

func (m FileSelectionModel) Confirmed() bool {
	return m.confirmed
}

func (m FileSelectionModel) Reset() FileSelectionModel {
	m.ready = false
	m.cursor = 0
	return m
}

func (m FileSelectionModel) getAllAvailableFiles() []string {
	allFiles := make([]string, 0)

	// Walk through all files in Claude assets (handle nil embed.FS gracefully)
	if m.claudeAssets != nil {
		fs.WalkDir(m.claudeAssets, "assets/claude", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			// Create relative path for display/selection
			relPath := strings.TrimPrefix(path, "assets/claude/")
			allFiles = append(allFiles, relPath)
			return nil
		})
	}

	// Walk through all files in Startup assets (handle nil embed.FS gracefully)
	if m.startupAssets != nil {
		fs.WalkDir(m.startupAssets, "assets/the-startup", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			// Create relative path for display/selection
			relPath := strings.TrimPrefix(path, "assets/the-startup/")
			allFiles = append(allFiles, relPath)
			return nil
		})
	}

	return allFiles
}

func (m FileSelectionModel) buildStaticTree() string {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	updateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange for updates
	removeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Strikethrough(true) // Red with strike-through for removal

	// For uninstall mode, show the actual files that will be removed in tree structure
	if m.mode == ModeUninstall {
		// Organize files by category for tree display
		claudeFiles := make(map[string][]string)
		startupFiles := make(map[string][]string)
		
		for _, filePath := range m.selectedFiles {
			if strings.Contains(filePath, m.installer.GetInstallPath()) {
				// Startup files
				relPath := strings.TrimPrefix(filePath, m.installer.GetInstallPath()+"/")
				if strings.HasPrefix(relPath, "templates/") {
					startupFiles["templates"] = append(startupFiles["templates"], strings.TrimPrefix(relPath, "templates/"))
				} else if strings.HasPrefix(relPath, "rules/") {
					startupFiles["rules"] = append(startupFiles["rules"], strings.TrimPrefix(relPath, "rules/"))
				} else if strings.HasPrefix(relPath, "bin/") {
					startupFiles["bin"] = append(startupFiles["bin"], strings.TrimPrefix(relPath, "bin/"))
				} else {
					startupFiles["other"] = append(startupFiles["other"], relPath)
				}
			} else if strings.Contains(filePath, m.installer.GetClaudePath()) {
				// Claude files
				relPath := strings.TrimPrefix(filePath, m.installer.GetClaudePath()+"/")
				if strings.HasPrefix(relPath, "agents/") {
					claudeFiles["agents"] = append(claudeFiles["agents"], strings.TrimPrefix(relPath, "agents/"))
				} else if strings.HasPrefix(relPath, "commands/") {
					claudeFiles["commands"] = append(claudeFiles["commands"], strings.TrimPrefix(relPath, "commands/"))
				} else if strings.HasPrefix(relPath, "output-styles/") {
					claudeFiles["output-styles"] = append(claudeFiles["output-styles"], strings.TrimPrefix(relPath, "output-styles/"))
				} else {
					claudeFiles["other"] = append(claudeFiles["other"], relPath)
				}
			}
		}

		// Build tree structure with proper children format
		var children []any
		
		// Add Claude files section
		if len(claudeFiles) > 0 {
			children = append(children, ".claude/")
			
			// Sort categories for consistent display
			categories := []string{"agents", "commands", "output-styles", "other"}
			for _, category := range categories {
				if files, exists := claudeFiles[category]; exists && len(files) > 0 {
					children = append(children, "  "+category+"/")
					categoryTree := tree.New()
					sort.Strings(files)
					for _, file := range files {
						categoryTree = categoryTree.Child(removeStyle.Render("✗ " + file))
					}
					children = append(children, categoryTree)
				}
			}
		}

		// Add startup files section
		if len(startupFiles) > 0 {
			children = append(children, ".the-startup/")
			
			// Sort categories for consistent display
			categories := []string{"templates", "rules", "bin", "other"}
			for _, category := range categories {
				if files, exists := startupFiles[category]; exists && len(files) > 0 {
					if category == "other" && len(files) == 1 && files[0] == "the-startup.lock" {
						// Show lockfile directly under startup
						children = append(children, removeStyle.Render("✗ " + files[0]))
					} else {
						children = append(children, "  "+category+"/")
						categoryTree := tree.New()
						sort.Strings(files)
						for _, file := range files {
							categoryTree = categoryTree.Child(removeStyle.Render("✗ " + file))
						}
						children = append(children, categoryTree)
					}
				}
			}
		}
		
		// Build final tree
		t := tree.Root("Files to be removed").
			Child(children...).
			Enumerator(tree.RoundedEnumerator).
			EnumeratorStyle(enumeratorStyle).
			RootStyle(rootStyle)

		return t.String()
	}

	// Get list of existing files that will be updated
	existingFiles := make(map[string]bool)
	if m.installer != nil {
		for _, file := range m.installer.GetExistingFiles(m.selectedFiles) {
			existingFiles[file] = true
		}
	}

	// Get list of deprecated files that will be removed
	deprecatedFiles := make(map[string]bool)
	if m.installer != nil {
		for _, file := range m.installer.GetDeprecatedFilesList() {
			deprecatedFiles[file] = true
		}
	}

	buildSubtree := func(embedFS *embed.FS, basePath string, prefix string) []any {
		filesSeen := make(map[string]bool)
		
		// Build a hierarchical structure for files
		type fileNode struct {
			name     string
			children map[string]*fileNode
			isFile   bool
			fullPath string
			exists   bool
		}
		
		root := &fileNode{
			children: make(map[string]*fileNode),
		}
		
		// First, add current files from embedded assets using WalkDir for recursive discovery
		if embedFS != nil {
			fs.WalkDir(embedFS, basePath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				
				// Skip the base directory itself
				if path == basePath {
					return nil
				}
				
				// For files, only process markdown files
				if !d.IsDir() && !strings.HasSuffix(path, ".md") && !strings.HasSuffix(path, ".json") {
					return nil
				}
				
				// Extract relative path from assets/claude/[type]/
				relPath := strings.TrimPrefix(path, basePath+"/")
				
				// For agents, validate naming convention
				// Agents at root level should start with "the-", but nested agents don't need to
				if prefix == "agents/" && !d.IsDir() {
					// Check if it's a root-level agent (no directory separator in path)
					if !strings.Contains(relPath, "/") {
						name := filepath.Base(path)
						name = strings.TrimSuffix(name, ".md")
						if !strings.HasPrefix(name, "the-") {
							return nil // Skip root agents that don't follow naming convention
						}
					}
					// Nested agents don't need "the-" prefix
				}
				filePath := prefix + relPath
				
				if !d.IsDir() {
					filesSeen[relPath] = true
				}
				
				// Build the tree structure
				parts := strings.Split(relPath, "/")
				current := root
				
				for i, part := range parts {
					if _, exists := current.children[part]; !exists {
						current.children[part] = &fileNode{
							name:     part,
							children: make(map[string]*fileNode),
							isFile:   i == len(parts)-1 && !d.IsDir(),
							fullPath: filePath,
							exists:   existingFiles[filePath],
						}
					}
					current = current.children[part]
				}
				
				return nil
			})
		}
		
		// Then, add deprecated files that will be removed (only for this prefix)
		for depFile := range deprecatedFiles {
			if strings.HasPrefix(depFile, prefix) {
				relPath := strings.TrimPrefix(depFile, prefix)
				// Only add if we haven't seen this file (it's truly deprecated)
				if !filesSeen[relPath] {
					// Build the tree structure for deprecated files
					parts := strings.Split(relPath, "/")
					current := root
					
					for i, part := range parts {
						if _, exists := current.children[part]; !exists {
							current.children[part] = &fileNode{
								name:     part,
								children: make(map[string]*fileNode),
								isFile:   i == len(parts)-1,
								fullPath: depFile,
								exists:   false, // Deprecated files don't exist in new structure
							}
						}
						current = current.children[part]
					}
				}
			}
		}
		
		// Convert the tree structure to lipgloss tree items
		var buildTreeItems func(node *fileNode, depth int) []any
		buildTreeItems = func(node *fileNode, depth int) []any {
			var items []any
			
			// Sort children for consistent display
			var sortedKeys []string
			for k := range node.children {
				sortedKeys = append(sortedKeys, k)
			}
			sort.Strings(sortedKeys)
			
			for _, key := range sortedKeys {
				child := node.children[key]
				
				if child.isFile {
					// It's a file - remove .md extension for display
					displayName := strings.TrimSuffix(child.name, ".md")
					displayName = strings.TrimSuffix(displayName, ".json")
					
					// For commands, format as /s:command
					if prefix == "commands/" && strings.Contains(child.fullPath, "/s/") {
						// Convert s/specify.md to /s:specify
						displayName = "/s:" + displayName
					}
					
					// Check if it's deprecated (not in filesSeen)
					if !filesSeen[strings.TrimPrefix(child.fullPath, prefix)] && deprecatedFiles[child.fullPath] {
						displayName = "✗ " + displayName + " (will remove)"
						items = append(items, removeStyle.Render(displayName))
					} else if child.exists {
						items = append(items, updateStyle.Render(displayName+" (will update)"))
					} else {
						items = append(items, itemStyle.Render(displayName))
					}
				} else {
					// It's a directory
					
					// For agents prefix, show condensed format with activity count
					if prefix == "agents/" {
						// Count the .md files for specialized activities
						fileCount := 0
						for _, grandchild := range child.children {
							if grandchild.isFile && strings.HasSuffix(grandchild.name, ".md") {
								fileCount++
							}
						}
						
						// Format: "the-analyst (5 specialized activities)"
						displayName := child.name
						if fileCount > 0 {
							activityLabel := "specialized activity"
							if fileCount > 1 {
								activityLabel = "specialized activities"
							}
							displayName = fmt.Sprintf("%s (%d %s)", child.name, fileCount, activityLabel)
						}
						
						// Don't expand agent directories, just show the count
						items = append(items, itemStyle.Render(displayName))
					} else if prefix == "commands/" && child.name == "s" {
						// For commands/s directory, expand but format files specially
						dirItems := buildTreeItems(child, depth+1)
						// Add the command items directly without showing "s/" directory
						items = append(items, dirItems...)
					} else {
						// For other directories, expand normally
						dirItems := buildTreeItems(child, depth+1)
						if len(dirItems) > 0 {
							// Create a subtree for this directory
							dirTree := tree.New()
							for _, item := range dirItems {
								dirTree = dirTree.Child(item)
							}
							items = append(items, child.name+"/")
							items = append(items, dirTree)
						}
					}
				}
			}
			
			return items
		}
		
		return buildTreeItems(root, 0)
	}

	// Only show files that go to .claude directory (agents and commands)
	agentItems := buildSubtree(m.claudeAssets, "assets/claude/agents", "agents/")
	commandItems := buildSubtree(m.claudeAssets, "assets/claude/commands", "commands/")
	outputStyleItems := buildSubtree(m.claudeAssets, "assets/claude/output-styles", "output-styles/")
	// Don't show hooks and templates as they go to .the-startup, not .claude

	// Check if settings.json exists using installer's method
	settingsExists := m.installer.CheckSettingsExists()

	claudePath := m.installer.GetClaudePath()

	displayPath := claudePath
	if strings.HasPrefix(claudePath, os.Getenv("HOME")) {
		displayPath = strings.Replace(claudePath, os.Getenv("HOME"), "~", 1)
	}

	// Add settings.json with appropriate styling
	settingsItem := "settings.json"
	if settingsExists {
		settingsItem = updateStyle.Render("settings.json (will update)")
	} else {
		settingsItem = itemStyle.Render("settings.json")
	}

	// Build children list for the tree with proper nesting
	children := []any{}
	
	// Add agents folder with its items as a subtree
	if len(agentItems) > 0 {
		agentsTree := tree.New()
		for _, item := range agentItems {
			agentsTree = agentsTree.Child(item)
		}
		children = append(children, "agents/")
		children = append(children, agentsTree)
	}
	
	// Add commands folder with its items as a subtree
	if len(commandItems) > 0 {
		commandsTree := tree.New()
		for _, item := range commandItems {
			commandsTree = commandsTree.Child(item)
		}
		children = append(children, "commands/")
		children = append(children, commandsTree)
	}
	
	// Add output-styles folder with its items as a subtree
	if len(outputStyleItems) > 0 {
		outputStylesTree := tree.New()
		for _, item := range outputStyleItems {
			outputStylesTree = outputStylesTree.Child(item)
		}
		children = append(children, "output-styles/")
		children = append(children, outputStylesTree)
	}
	
	// Settings files go at root level
	children = append(children, settingsItem)

	// Add settings.local.json if it exists in assets (handle nil embed.FS gracefully)
	if m.claudeAssets != nil {
		if _, err := m.claudeAssets.ReadFile("assets/claude/settings.local.json"); err == nil {
			localSettingsPath := filepath.Join(claudePath, "settings.local.json")
			localSettingsItem := "settings.local.json"
			if _, err := os.Stat(localSettingsPath); err == nil {
				localSettingsItem = updateStyle.Render("settings.local.json (will update)")
			} else {
				localSettingsItem = itemStyle.Render("settings.local.json")
			}
			children = append(children, localSettingsItem)
		}
	}

	t := tree.
		Root("⁜ " + displayPath).
		Child(children...).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle)

	return t.String()
}

// getFilesFromLockfile reads the lockfile to get the list of installed files for uninstall
// It validates that files actually exist on disk and handles missing lockfile gracefully
func (m FileSelectionModel) getFilesFromLockfile() []string {
	// Get paths from installer (set by the path selection steps)
	if m.installer == nil {
		return []string{}
	}
	
	startupPath := m.installer.GetInstallPath()
	claudePath := m.installer.GetClaudePath()
	
	// Try to read lockfile from the selected startup path
	lockfilePath := filepath.Join(startupPath, "the-startup.lock")
	if _, err := os.Stat(lockfilePath); err != nil {
		// Lockfile doesn't exist - nothing to uninstall
		return []string{}
	}
	
	// Read and parse lockfile
	lockfileData, err := os.ReadFile(lockfilePath)
	if err != nil {
		return []string{}
	}
	
	var lockFile map[string]interface{}
	if err := json.Unmarshal(lockfileData, &lockFile); err != nil {
		return []string{}
	}
	
	// Extract files from lockfile
	filesInterface, ok := lockFile["files"]
	if !ok {
		return []string{}
	}
	
	filesMap, ok := filesInterface.(map[string]interface{})
	if !ok {
		return []string{}
	}
	
	var existingFiles []string
	
	// Check each file in lockfile to see if it still exists
	for lockfilePath := range filesMap {
		var fullPath string
		
		// Determine full path based on file type
		if strings.HasPrefix(lockfilePath, "startup/") {
			// Startup files go to install path
			relPath := strings.TrimPrefix(lockfilePath, "startup/")
			fullPath = filepath.Join(startupPath, relPath)
		} else if strings.HasPrefix(lockfilePath, "bin/") {
			// Binary files go to install path
			fullPath = filepath.Join(startupPath, lockfilePath)
		} else {
			// Claude files (agents, commands, output-styles) go to claude path
			fullPath = filepath.Join(claudePath, lockfilePath)
		}
		
		// Only include files that actually exist on disk
		if _, err := os.Stat(fullPath); err == nil {
			existingFiles = append(existingFiles, fullPath)
		}
	}
	
	// Also check for lockfile itself
	if _, err := os.Stat(lockfilePath); err == nil {
		existingFiles = append(existingFiles, lockfilePath)
	}
	
	return existingFiles
}
