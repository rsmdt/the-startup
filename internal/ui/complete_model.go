package ui

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/rsmdt/the-startup/internal/installer"
)

type CompleteModel struct {
	styles        Styles
	installer     *installer.Installer
	selectedTool  string
	ready         bool
	mode          OperationMode
	claudeAssets  *embed.FS
	startupAssets *embed.FS
	selectedFiles []string
}

// autoExitMsg is sent after a delay to automatically exit
type autoExitMsg struct{}

func NewCompleteModel(selectedTool string, installer *installer.Installer) CompleteModel {
	return NewCompleteModelWithMode(selectedTool, installer, ModeInstall)
}

func NewCompleteModelWithMode(selectedTool string, installer *installer.Installer, mode OperationMode) CompleteModel {
	return CompleteModel{
		styles:       GetStyles(),
		installer:    installer,
		selectedTool: selectedTool,
		ready:        false,
		mode:         mode,
	}
}

func NewCompleteModelWithAssets(selectedTool string, installer *installer.Installer, mode OperationMode, claudeAssets, startupAssets *embed.FS, selectedFiles []string) CompleteModel {
	return CompleteModel{
		styles:        GetStyles(),
		installer:     installer,
		selectedTool:  selectedTool,
		ready:         false,
		mode:          mode,
		claudeAssets:  claudeAssets,
		startupAssets: startupAssets,
		selectedFiles: selectedFiles,
	}
}

func (m CompleteModel) Init() tea.Cmd {
	// Exit immediately
	return func() tea.Msg {
		return autoExitMsg{}
	}
}

func (m CompleteModel) Update(msg tea.Msg) (CompleteModel, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		// Allow immediate exit on any key press
		m.ready = true
	case autoExitMsg:
		// Auto-exit after delay
		m.ready = true
	}
	return m, nil
}

func (m CompleteModel) View() string {
	var s strings.Builder

	// Banner
	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")

	// Success message
	if m.mode == ModeUninstall {
		s.WriteString(m.styles.Success.Render("✅ Uninstallation Complete!"))
	} else {
		s.WriteString(m.styles.Success.Render("✅ Installation Complete!"))
	}
	s.WriteString("\n\n")

	// Installation locations
	claudePath := m.installer.GetClaudePath()
	startupPath := m.installer.GetInstallPath()

	// Simplify paths for display
	displayClaudePath := claudePath
	displayStartupPath := startupPath

	home := os.Getenv("HOME")
	if home != "" {
		if strings.HasPrefix(claudePath, home) {
			displayClaudePath = strings.Replace(claudePath, home, "~", 1)
		}
		if strings.HasPrefix(startupPath, home) {
			displayStartupPath = strings.Replace(startupPath, home, "~", 1)
		}
	}

	if m.mode == ModeUninstall {
		s.WriteString(m.styles.Normal.Render("Uninstallation locations:"))
		s.WriteString("\n")
		s.WriteString(m.styles.Warning.Render("  Claude files: " + displayClaudePath))
		s.WriteString("\n")
		s.WriteString(m.styles.Warning.Render("  Startup files: " + displayStartupPath))
		s.WriteString("\n\n")
	} else {
		s.WriteString(m.styles.Normal.Render("Installation locations:"))
		s.WriteString("\n")
		s.WriteString(m.styles.Info.Render("  Claude files: " + displayClaudePath))
		s.WriteString("\n")
		s.WriteString(m.styles.Info.Render("  Startup files: " + displayStartupPath))
		s.WriteString("\n\n")
	}

	// Display the tree of installed/removed files - same as during selection
	if m.mode == ModeUninstall {
		s.WriteString(m.styles.Warning.Render("Files removed:"))
	} else {
		s.WriteString(m.styles.Info.Render("Files installed:"))
	}
	s.WriteString("\n\n")
	
	// Build and display the tree (same structure as file selection)
	treeStr := m.buildCompletionTree()
	s.WriteString(treeStr)
	s.WriteString("\n")

	// Display removed files if any
	removedFiles := m.installer.GetDeprecatedFilesList()
	if len(removedFiles) > 0 {
		sort.Strings(removedFiles)
		s.WriteString(m.styles.Error.Render("  Removed deprecated files:"))
		s.WriteString("\n")
		for _, file := range removedFiles {
			s.WriteString(m.styles.Normal.Render("    ✗ " + file))
			s.WriteString("\n")
		}
		s.WriteString("\n")
	}

	// Add repository link
	s.WriteString(m.styles.Info.Render("See https://github.com/rsmdt/the-startup for details"))
	s.WriteString("\n")

	return s.String()
}

func (m CompleteModel) Ready() bool {
	return m.ready
}

// buildCompletionTree builds the same tree structure shown during file selection
func (m CompleteModel) buildCompletionTree() string {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	updateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange for updates
	
	// For install mode, build the tree showing what was installed
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
		
		// Add files from selectedFiles list that match this prefix
		for _, file := range m.selectedFiles {
			if strings.HasPrefix(file, prefix) {
				relPath := strings.TrimPrefix(file, prefix)
				filesSeen[relPath] = true
				
				// Build the tree structure
				parts := strings.Split(relPath, "/")
				current := root
				
				for i, part := range parts {
					if _, exists := current.children[part]; !exists {
						current.children[part] = &fileNode{
							name:     part,
							children: make(map[string]*fileNode),
							isFile:   i == len(parts)-1,
							fullPath: file,
							exists:   false, // For completion, we show them as new
						}
					}
					current = current.children[part]
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
					
					if m.mode == ModeUninstall {
						items = append(items, itemStyle.Render("✗ " + displayName))
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
	
	// Build trees for each category
	agentItems := buildSubtree(m.claudeAssets, "assets/claude/agents", "agents/")
	commandItems := buildSubtree(m.claudeAssets, "assets/claude/commands", "commands/")
	outputStyleItems := buildSubtree(m.claudeAssets, "assets/claude/output-styles", "output-styles/")
	
	// Check if settings.json was updated
	settingsExists := m.installer.CheckSettingsExists()
	
	claudePath := m.installer.GetClaudePath()
	displayPath := claudePath
	if strings.HasPrefix(claudePath, os.Getenv("HOME")) {
		displayPath = strings.Replace(claudePath, os.Getenv("HOME"), "~", 1)
	}
	
	// Add settings.json with appropriate styling
	settingsItem := "settings.json"
	if m.mode == ModeUninstall {
		settingsItem = itemStyle.Render("✗ settings.json")
	} else if settingsExists {
		settingsItem = updateStyle.Render("settings.json (updated)")
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
	
	// Add settings.local.json if it exists in selected files
	for _, file := range m.selectedFiles {
		if file == "settings.local.json" {
			localSettingsPath := filepath.Join(claudePath, "settings.local.json")
			localSettingsItem := "settings.local.json"
			if m.mode == ModeUninstall {
				localSettingsItem = itemStyle.Render("✗ settings.local.json")
			} else if _, err := os.Stat(localSettingsPath); err == nil {
				localSettingsItem = updateStyle.Render("settings.local.json (updated)")
			} else {
				localSettingsItem = itemStyle.Render("settings.local.json")
			}
			children = append(children, localSettingsItem)
			break
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
