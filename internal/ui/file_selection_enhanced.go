package ui

import (
	"fmt"
	"io/fs"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// FileSelectionInterface defines the common interface for file selection models
type FileSelectionInterface interface {
	Update(tea.Msg) (FileSelectionInterface, tea.Cmd)
	View() string
	Ready() bool
	Confirmed() bool
	GetSelectedFiles() []string
	GetChoices() []string // For compatibility with tests
}

// FileSelectionWrapper wraps the original FileSelectionModel to implement the interface
type FileSelectionWrapper struct {
	FileSelectionModel
}

// Update wraps the original Update method to return the interface
func (w FileSelectionWrapper) Update(msg tea.Msg) (FileSelectionInterface, tea.Cmd) {
	original, cmd := w.FileSelectionModel.Update(msg)
	w.FileSelectionModel = original
	return w, cmd
}

// GetSelectedFiles returns the selected files
func (w FileSelectionWrapper) GetSelectedFiles() []string {
	return w.FileSelectionModel.selectedFiles
}

// GetChoices returns the file choices for test compatibility
func (w FileSelectionWrapper) GetChoices() []string {
	return w.FileSelectionModel.choices
}

// TreeItem represents a selectable item in the file tree
type TreeItem struct {
	Name        string // Display name (e.g., "the-analyst (5 specialized activities)")
	Path        string // Relative file path (e.g., "agents/the-analyst")
	Selected    bool   // Whether this item is selected
	IsFile      bool   // True if it's a file, false if it's a folder/category
	IsSettings  bool   // True if it's a settings file
	IsUpdate    bool   // True if this file already exists and will be updated
}

// EnhancedFileSelectionModel extends the existing FileSelectionModel with interactive tree selection
type EnhancedFileSelectionModel struct {
	FileSelectionModel

	// Tree selection state
	treeItems       []TreeItem
	cursor          int
	showingTree     bool
	showingSummary  bool
}


// NewEnhancedFileSelectionModelFromExisting creates enhanced model from existing model
func NewEnhancedFileSelectionModelFromExisting(existing FileSelectionModel) EnhancedFileSelectionModel {
	enhanced := EnhancedFileSelectionModel{
		FileSelectionModel: existing,
		cursor:             0,
		showingTree:        true, // Start directly in tree mode
	}

	// Build tree items using the same logic as the original tree
	enhanced.treeItems = enhanced.buildTreeItems()

	// Update selected files based on tree selections
	enhanced.updateSelectedFilesFromTree()

	return enhanced
}

// buildTreeItems creates the interactive tree structure using filesystem data
func (m EnhancedFileSelectionModel) buildTreeItems() []TreeItem {
	var items []TreeItem

	// Add directory headers with proper paths for selection tracking (Claude files only)
	items = append(items, TreeItem{Name: ".claude/", Path: "__DIR__.claude", Selected: true, IsFile: false})
	items = append(items, TreeItem{Name: "agents/", Path: "__DIR__agents", Selected: true, IsFile: false})
	items = append(items, m.buildAgentItems()...)

	items = append(items, TreeItem{Name: "commands/", Path: "__DIR__commands", Selected: true, IsFile: false})
	items = append(items, m.buildCommandItems()...)

	items = append(items, TreeItem{Name: "output-styles/", Path: "__DIR__output-styles", Selected: true, IsFile: false})
	items = append(items, m.buildOutputStyleItems()...)

	// Settings files - check if settings.json exists
	settingsExists := m.installer.CheckSettingsExists()
	items = append(items, TreeItem{Name: "settings.json", Path: "settings.json", Selected: true, IsFile: true, IsSettings: true, IsUpdate: settingsExists})

	// Add settings.local.json if it exists
	if m.claudeAssets != nil {
		if _, err := m.claudeAssets.ReadFile("assets/claude/settings.local.json"); err == nil {
			items = append(items, TreeItem{Name: "settings.local.json", Path: "settings.local.json", Selected: true, IsFile: true, IsSettings: true})
		}
	}

	// Add Continue option at the bottom
	items = append(items, TreeItem{Name: "Continue with selected files", Path: "__CONTINUE__", Selected: false, IsFile: false})

	return items
}

// Update handles enhanced file selection input and returns FileSelectionInterface
func (m EnhancedFileSelectionModel) Update(msg tea.Msg) (FileSelectionInterface, tea.Cmd) {
	// Handle enhanced selection mode
	if m.mode == ModeInstall && !m.ready {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if m.showingSummary {
					// Handle final confirmation choices
					if m.cursor == 0 {
						// Yes, install
						m.confirmed = true
						m.ready = true
						return m, nil
					} else {
						// No, go back to tree selection
						m.showingSummary = false
						m.showingTree = true
						m.cursor = len(m.treeItems) - 1 // Go back to Continue option
						return m, nil
					}
				} else if m.showingTree {
					// Check if we're on the Continue option
					if m.cursor < len(m.treeItems) && m.treeItems[m.cursor].Path == "__CONTINUE__" {
						// Update selected files and go to summary
						m.updateSelectedFilesFromTree()
						m.showingTree = false
						m.showingSummary = true
						m.cursor = 0 // Reset cursor for summary choices
						return m, nil
					}
				}
			case "esc":
				if m.showingSummary {
					// Go back to tree selection
					m.showingSummary = false
					m.showingTree = true
					m.cursor = len(m.treeItems) - 1 // Go back to Continue option
					return m, nil
				}
				// For tree view, fall through to normal handling (go back to previous screen)
			case "up", "k":
				if m.showingSummary {
					// Navigate summary choices
					if m.cursor > 0 {
						m.cursor--
					} else {
						m.cursor = 1 // Wrap to "No" option
					}
					return m, nil
				} else if m.showingTree {
					enhanced := m.handleTreeNavigation(-1)
					return enhanced, nil
				}
			case "down", "j":
				if m.showingSummary {
					// Navigate summary choices
					if m.cursor < 1 {
						m.cursor++
					} else {
						m.cursor = 0 // Wrap to "Yes" option
					}
					return m, nil
				} else if m.showingTree {
					enhanced := m.handleTreeNavigation(1)
					return enhanced, nil
				}
			case " ": // Space bar
				if m.showingTree {
					enhanced := m.toggleCurrentTreeItem()
					return enhanced, nil
				}
			}
		}
	}

	// Fall back to original logic for uninstall mode or when ready
	original, cmd := m.FileSelectionModel.Update(msg)
	m.FileSelectionModel = original
	return m, cmd
}

// handleTreeNavigation moves cursor up/down through tree items
func (m EnhancedFileSelectionModel) handleTreeNavigation(direction int) EnhancedFileSelectionModel {
	totalItems := len(m.treeItems)
	if totalItems == 0 {
		return m
	}

	m.cursor += direction

	// Wrap around
	if m.cursor < 0 {
		m.cursor = totalItems - 1
	} else if m.cursor >= totalItems {
		m.cursor = 0
	}

	return m
}

// toggleCurrentTreeItem toggles selection of current tree item and handles directory select/unselect all
func (m EnhancedFileSelectionModel) toggleCurrentTreeItem() EnhancedFileSelectionModel {
	if m.cursor >= 0 && m.cursor < len(m.treeItems) {
		item := &m.treeItems[m.cursor]

		if item.Path == "__CONTINUE__" {
			// Don't toggle continue option
			return m
		}

		if item.IsFile {
			// Toggle individual file
			item.Selected = !item.Selected
		} else {
			// Directory - toggle all children
			item.Selected = !item.Selected

			// Determine which children to toggle based on directory type
			var pathPrefix string
			switch item.Path {
			case "__DIR__.claude":
				// Toggle all agents, commands, output-styles, and settings
				for i := range m.treeItems {
					child := &m.treeItems[i]
					if child.IsFile && (strings.HasPrefix(child.Path, "agents/") ||
						strings.HasPrefix(child.Path, "commands/") ||
						strings.HasPrefix(child.Path, "output-styles/") ||
						child.IsSettings) {
						child.Selected = item.Selected
					}
				}
				// Also toggle subdirectories
				for i := range m.treeItems {
					child := &m.treeItems[i]
					if !child.IsFile && (child.Path == "__DIR__agents" ||
						child.Path == "__DIR__commands" ||
						child.Path == "__DIR__output-styles") {
						child.Selected = item.Selected
					}
				}
			case "__DIR__agents":
				pathPrefix = "agents/"
			case "__DIR__commands":
				pathPrefix = "commands/"
			case "__DIR__output-styles":
				pathPrefix = "output-styles/"
			}

			// Toggle children with matching prefix
			if pathPrefix != "" {
				for i := range m.treeItems {
					child := &m.treeItems[i]
					if child.IsFile && strings.HasPrefix(child.Path, pathPrefix) {
						child.Selected = item.Selected
					}
				}
			}
		}
	}
	return m
}







// View renders the enhanced file selection
func (m EnhancedFileSelectionModel) View() string {
	if m.mode == ModeInstall && !m.ready {
		if m.showingSummary {
			return m.renderSummary()
		} else {
			return m.renderTreeSelection()
		}
	}

	// Fall back to original view for uninstall or when ready
	return m.FileSelectionModel.View()
}

// renderInitialView shows the initial prompt to enter interactive selection
func (m EnhancedFileSelectionModel) renderInitialView() string {
	var s strings.Builder

	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")

	// Show both paths in the header using standardized display
	startupPath := m.installer.GetInstallPath()
	claudePath := m.installer.GetClaudePath()

	renderer := NewProgressiveDisclosureRenderer()
	s.WriteString(renderer.RenderSelectedPaths(claudePath, startupPath, m.mode))

	// Show component summary
	s.WriteString(m.styles.Title.Render("🚀 Interactive Component Selection"))
	s.WriteString("\n\n")

	// Count selected items
	selectedCount := 0
	for _, item := range m.treeItems {
		if item.Selected && item.IsFile {
			selectedCount++
		}
	}

	s.WriteString(m.renderer.RenderTitle("Files to be installed to .claude"))

	s.WriteString(m.styles.Info.Render("The following files will be installed to your Claude directory:"))
	s.WriteString("\n\n")

	// Show the original file tree using the same logic as the original
	s.WriteString(m.FileSelectionModel.buildStaticTree())
	s.WriteString("\n")

	s.WriteString("\n")
	s.WriteString(m.styles.Info.Render(fmt.Sprintf("Ready to install %d components from The (Agentic) Startup", selectedCount)))
	s.WriteString("\n")

	s.WriteString("\n")
	s.WriteString(m.styles.Selected.Render("> Press Enter to customize your selection"))
	s.WriteString("\n")
	s.WriteString(m.styles.Normal.Render("  Press Escape to go back"))
	s.WriteString("\n\n")
	s.WriteString(m.styles.Help.Render("Enter to continue • Esc to go back"))

	return s.String()
}

// renderTreeSelection shows the interactive tree selection interface
func (m EnhancedFileSelectionModel) renderTreeSelection() string {
	var s strings.Builder

	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")

	// Show both paths using standardized renderer
	startupPath := m.installer.GetInstallPath()
	claudePath := m.installer.GetClaudePath()

	renderer := NewProgressiveDisclosureRenderer()
	s.WriteString(renderer.RenderSelectedPaths(claudePath, startupPath, m.mode))

	s.WriteString(m.styles.Title.Render("Interactive Component Selection"))
	s.WriteString("\n\n")

	// Show selection summary
	selectedCount := 0
	for _, item := range m.treeItems {
		if item.Selected && item.IsFile {
			selectedCount++
		}
	}
	s.WriteString(m.styles.Info.Render(fmt.Sprintf("Selected: %d components", selectedCount)))
	s.WriteString("\n\n")

	// Render interactive tree using file tree format
	s.WriteString(m.buildInteractiveTree())

	// Add Continue option with visual separation
	s.WriteString("\n")
	continueStyle := m.styles.Normal
	if m.cursor == len(m.treeItems)-1 { // Continue is always last
		continueStyle = m.styles.Selected
	}
	s.WriteString(continueStyle.Render("▶ Continue with selected items"))
	s.WriteString("\n")

	s.WriteString("\n\n")
	s.WriteString(m.styles.Help.Render("Space: toggle • ↑↓: navigate • Enter: continue"))
	s.WriteString("\n")
	s.WriteString(m.styles.Help.Render("Navigate to 'Continue' at bottom to proceed • Esc: go back"))

	return s.String()
}

// GetSelectedFiles returns the selected files from the enhanced model
func (m EnhancedFileSelectionModel) GetSelectedFiles() []string {
	return m.selectedFiles
}

// GetChoices returns the file choices from the embedded model for test compatibility
func (m EnhancedFileSelectionModel) GetChoices() []string {
	return m.FileSelectionModel.choices
}

// buildAgentItems creates tree items for agents using filesystem data
func (m EnhancedFileSelectionModel) buildAgentItems() []TreeItem {
	var items []TreeItem

	if m.claudeAssets == nil {
		return items
	}

	// Get all available files first to check which ones exist
	allFiles := m.getAllAvailableFiles()
	existingFiles := m.installer.GetExistingFiles(allFiles)
	existingSet := make(map[string]bool)
	for _, file := range existingFiles {
		existingSet[file] = true
	}

	// Use the same logic as buildSubtree for agents but limit to depth 1
	fs.WalkDir(m.claudeAssets, "assets/claude/agents", func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == "assets/claude/agents" {
			return nil
		}

		// Extract relative path
		relPath := strings.TrimPrefix(path, "assets/claude/agents/")

		// Limit to depth 1 only (directly in agents directory)
		pathDepth := strings.Count(relPath, "/")
		if pathDepth > 0 {
			// Skip files/directories deeper than level 1
			return nil
		}

		if d.IsDir() {
			// Count the .md files for specialized activities
			fileCount := 0
			fs.WalkDir(m.claudeAssets, path, func(subPath string, subD fs.DirEntry, subErr error) error {
				if subErr == nil && !subD.IsDir() && strings.HasSuffix(subPath, ".md") {
					fileCount++
				}
				return nil
			})

			// Format: "the-analyst (5 specialized activities)"
			displayName := d.Name()
			if fileCount > 0 {
				activityLabel := "specialized activity"
				if fileCount > 1 {
					activityLabel = "specialized activities"
				}
				displayName = fmt.Sprintf("%s (%d %s)", d.Name(), fileCount, activityLabel)
			}

			// Add as selectable item, pre-selected
			items = append(items, TreeItem{
				Name:     "  " + displayName,
				Path:     "agents/" + relPath,
				Selected: true, // Pre-select all agents
				IsFile:   true,
				IsUpdate: existingSet["agents/"+relPath],
			})
		} else if strings.HasSuffix(d.Name(), ".md") {
			// Handle leaf files directly in the agents directory
			displayName := strings.TrimSuffix(d.Name(), ".md")

			// Add as selectable item, pre-selected
			items = append(items, TreeItem{
				Name:     "  " + displayName,
				Path:     "agents/" + relPath,
				Selected: true, // Pre-select all agents
				IsFile:   true,
				IsUpdate: existingSet["agents/"+relPath],
			})
		}
		return nil
	})

	return items
}

// buildCommandItems creates tree items for commands using filesystem data
func (m EnhancedFileSelectionModel) buildCommandItems() []TreeItem {
	var items []TreeItem

	if m.claudeAssets == nil {
		return items
	}

	// Get all available files first to check which ones exist
	allFiles := m.getAllAvailableFiles()
	existingFiles := m.installer.GetExistingFiles(allFiles)
	existingSet := make(map[string]bool)
	for _, file := range existingFiles {
		existingSet[file] = true
	}

	// Walk through commands directory (limit to depth 2 relative to claude)
	fs.WalkDir(m.claudeAssets, "assets/claude/commands", func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == "assets/claude/commands" {
			return nil
		}

		relPath := strings.TrimPrefix(path, "assets/claude/commands/")

		// Limit to depth 2 relative to .claude (so commands/subdir/file.md is allowed)
		pathDepth := strings.Count(relPath, "/")
		if pathDepth > 1 {
			// Skip files/directories deeper than level 2
			return nil
		}

		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			// For files, format display name
			displayName := strings.TrimSuffix(d.Name(), ".md")

			// Special formatting for s/ commands
			if strings.Contains(relPath, "/s/") {
				displayName = "/s:" + displayName
			}

			items = append(items, TreeItem{
				Name:     "  " + displayName,
				Path:     "commands/" + relPath,
				Selected: true, // Pre-select all commands
				IsFile:   true,
				IsUpdate: existingSet["commands/"+relPath],
			})
		}
		return nil
	})

	return items
}

// buildOutputStyleItems creates tree items for output-styles using filesystem data
func (m EnhancedFileSelectionModel) buildOutputStyleItems() []TreeItem {
	var items []TreeItem

	if m.claudeAssets == nil {
		return items
	}

	// Get all available files first to check which ones exist
	allFiles := m.getAllAvailableFiles()
	existingFiles := m.installer.GetExistingFiles(allFiles)
	existingSet := make(map[string]bool)
	for _, file := range existingFiles {
		existingSet[file] = true
	}

	// Walk through output-styles directory (limit to depth 2 relative to claude)
	fs.WalkDir(m.claudeAssets, "assets/claude/output-styles", func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == "assets/claude/output-styles" {
			return nil
		}

		relPath := strings.TrimPrefix(path, "assets/claude/output-styles/")

		// Limit to depth 2 relative to .claude (so output-styles/subdir/file.json is allowed)
		pathDepth := strings.Count(relPath, "/")
		if pathDepth > 1 {
			// Skip files/directories deeper than level 2
			return nil
		}

		if !d.IsDir() && (strings.HasSuffix(path, ".json") || strings.HasSuffix(path, ".md")) {
			displayName := strings.TrimSuffix(d.Name(), ".json")
			displayName = strings.TrimSuffix(displayName, ".md")

			items = append(items, TreeItem{
				Name:     "  " + displayName,
				Path:     "output-styles/" + relPath,
				Selected: true, // Pre-select all output styles
				IsFile:   true,
				IsUpdate: existingSet["output-styles/"+relPath],
			})
		}
		return nil
	})

	return items
}

// updateSelectedFilesFromTree updates the selectedFiles based on tree selections
func (m *EnhancedFileSelectionModel) updateSelectedFilesFromTree() {
	var selectedFiles []string

	for _, item := range m.treeItems {
		if item.Selected && item.IsFile && item.Path != "__CONTINUE__" {
			selectedFiles = append(selectedFiles, item.Path)
		}
	}

	m.selectedFiles = selectedFiles
	if m.installer != nil {
		m.installer.SetSelectedFiles(selectedFiles)
	}
}

// renderSummary shows the final confirmation screen with selected paths and files
func (m EnhancedFileSelectionModel) renderSummary() string {
	var s strings.Builder

	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")

	// Show both paths using standardized renderer
	startupPath := m.installer.GetInstallPath()
	claudePath := m.installer.GetClaudePath()

	renderer := NewProgressiveDisclosureRenderer()
	s.WriteString(renderer.RenderSelectedPaths(claudePath, startupPath, m.mode))

	s.WriteString(m.styles.Title.Render("📋 Installation Summary"))
	s.WriteString("\n\n")

	s.WriteString(m.styles.Info.Render("The following components will be installed:"))
	s.WriteString("\n\n")

	// Show the tree format for selected files (non-interactive)
	s.WriteString(m.buildNonInteractiveTree())

	// Count total selected for summary
	totalSelected := 0
	for _, item := range m.treeItems {
		if item.Selected && item.IsFile {
			totalSelected++
		}
	}
	s.WriteString("\n")
	s.WriteString(m.styles.Info.Render(fmt.Sprintf("Total: %d components will be installed", totalSelected)))
	s.WriteString("\n\n")

	s.WriteString(m.styles.Title.Render("Ready to install?"))
	s.WriteString("\n\n")

	// Show confirmation choices
	choices := []string{
		"Yes, install these components",
		"No, go back to selection",
	}

	for i, choice := range choices {
		if i == m.cursor {
			s.WriteString(m.styles.Selected.Render("> " + choice))
		} else {
			s.WriteString(m.styles.Normal.Render("  " + choice))
		}
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(m.styles.Help.Render("Enter to confirm • ↑↓ to navigate • Esc to go back"))

	return s.String()
}

// buildInteractiveTree creates a file tree format with left-aligned selection indicators
func (m EnhancedFileSelectionModel) buildInteractiveTree() string {
	var s strings.Builder

	// Use themed styles for consistency (matching selection views)
	styles := GetStyles()
	rootStyle := styles.Title
	checkboxStyle := styles.Cursor                                       // Purple for checkboxes (same as selection indicators)
	normalStyle := styles.Normal                                         // Default text color (light gray)
	greenTextStyle := styles.CursorLine                                  // Green text for cursor highlighting (no bold)
	updateStyle := styles.Warning                                        // Orange/Peach for update indicators

	// Title
	s.WriteString(rootStyle.Render("Files to be installed"))
	s.WriteString("\n")

	// Get directory items for checkbox display (Claude files only)
	var claudeDir, agentsDir, commandsDir *TreeItem
	for i := range m.treeItems {
		item := &m.treeItems[i]
		switch item.Path {
		case "__DIR__.claude":
			claudeDir = item
		case "__DIR__agents":
			agentsDir = item
		case "__DIR__commands":
			commandsDir = item
		}
	}

	// Claude section with checkbox
	claudeCheckbox := "□"
	if claudeDir != nil && claudeDir.Selected {
		claudeCheckbox = "✓"
	}
	claudeHighlighted := false
	// Find claude directory index for cursor highlighting
	for i, item := range m.treeItems {
		if item.Path == "__DIR__.claude" && i == m.cursor {
			claudeHighlighted = true
			break
		}
	}
	var claudeLine string
	if claudeHighlighted {
		// For cursor highlighting: purple checkbox, green text
		claudeLine = fmt.Sprintf("%s %s", checkboxStyle.Render(claudeCheckbox), greenTextStyle.Render(".claude/"))
	} else {
		// Default styling: default checkbox, default text color
		claudeLine = fmt.Sprintf("%s %s", normalStyle.Render(claudeCheckbox), normalStyle.Render(".claude/"))
	}
	s.WriteString(claudeLine)
	s.WriteString("\n")

	// Agents subsection with checkbox
	hasAgents := false
	for _, item := range m.treeItems {
		if strings.HasPrefix(item.Path, "agents/") && item.IsFile {
			hasAgents = true
			break
		}
	}

	if hasAgents {
		agentsCheckbox := "□"
		if agentsDir != nil && agentsDir.Selected {
			agentsCheckbox = "✓"
		}
		agentsHighlighted := false
		// Find agents directory index for cursor highlighting
		for i, item := range m.treeItems {
			if item.Path == "__DIR__agents" && i == m.cursor {
				agentsHighlighted = true
				break
			}
		}
		var agentsLine string
		if agentsHighlighted {
			// For cursor highlighting: purple checkbox, green text
			agentsLine = fmt.Sprintf("%s %s%s", checkboxStyle.Render(agentsCheckbox), normalStyle.Render("├── "), greenTextStyle.Render("agents/"))
		} else {
			// Default styling: default checkbox, default text color
			agentsLine = fmt.Sprintf("%s %s", normalStyle.Render(agentsCheckbox), normalStyle.Render("├── agents/"))
		}
		s.WriteString(agentsLine)
		s.WriteString("\n")

		for i, item := range m.treeItems {
			if strings.HasPrefix(item.Path, "agents/") && item.IsFile {
				checkbox := "□"
				if item.Selected {
					checkbox = "✓"
				}

				displayName := strings.TrimPrefix(item.Name, "  ")
				updateIndicator := ""
				if item.IsUpdate {
					updateIndicator = updateStyle.Render(" (will update)")
				}

				if i == m.cursor {
					// For cursor highlighting: purple checkbox, green text
					line := fmt.Sprintf("%s %s%s%s", checkboxStyle.Render(checkbox), normalStyle.Render("│   ├── "), greenTextStyle.Render(displayName), updateIndicator)
					s.WriteString(line)
				} else {
					// Default styling: default checkbox, default text color
					line := fmt.Sprintf("%s %s%s", normalStyle.Render(checkbox), normalStyle.Render("│   ├── "+displayName), updateIndicator)
					s.WriteString(line)
				}
				s.WriteString("\n")
			}
		}
	}

	// Commands subsection
	hasCommands := false
	for _, item := range m.treeItems {
		if strings.HasPrefix(item.Path, "commands/") && item.IsFile {
			hasCommands = true
			break
		}
	}

	if hasCommands {
		commandsCheckbox := "□"
		if commandsDir != nil && commandsDir.Selected {
			commandsCheckbox = "✓"
		}
		commandsHighlighted := false
		// Find commands directory index for cursor highlighting
		for i, item := range m.treeItems {
			if item.Path == "__DIR__commands" && i == m.cursor {
				commandsHighlighted = true
				break
			}
		}
		var commandsLine string
		if commandsHighlighted {
			// For cursor highlighting: purple checkbox, green text
			commandsLine = fmt.Sprintf("%s %s%s", checkboxStyle.Render(commandsCheckbox), normalStyle.Render("├── "), greenTextStyle.Render("commands/"))
		} else {
			// Default styling: default checkbox, default text color
			commandsLine = fmt.Sprintf("%s %s", normalStyle.Render(commandsCheckbox), normalStyle.Render("├── commands/"))
		}
		s.WriteString(commandsLine)
		s.WriteString("\n")

		for i, item := range m.treeItems {
			if strings.HasPrefix(item.Path, "commands/") && item.IsFile {
				checkbox := "□"
				if item.Selected {
					checkbox = "✓"
				}

				displayName := strings.TrimPrefix(item.Name, "  ")
				updateIndicator := ""
				if item.IsUpdate {
					updateIndicator = updateStyle.Render(" (will update)")
				}

				if i == m.cursor {
					// For cursor highlighting: purple checkbox, green text
					line := fmt.Sprintf("%s %s%s%s", checkboxStyle.Render(checkbox), normalStyle.Render("│   ├── "), greenTextStyle.Render(displayName), updateIndicator)
					s.WriteString(line)
				} else {
					// Default styling: default checkbox, default text color
					line := fmt.Sprintf("%s %s%s", normalStyle.Render(checkbox), normalStyle.Render("│   ├── "+displayName), updateIndicator)
					s.WriteString(line)
				}
				s.WriteString("\n")
			}
		}
	}

	// Output-styles subsection
	hasOutputStyles := false
	var outputStylesDir *TreeItem
	for i := range m.treeItems {
		item := &m.treeItems[i]
		if item.Path == "__DIR__output-styles" {
			outputStylesDir = item
		}
		if strings.HasPrefix(item.Path, "output-styles/") && item.IsFile {
			hasOutputStyles = true
		}
	}

	// Check if we have settings files
	hasSettings := false
	for _, item := range m.treeItems {
		if (item.Path == "settings.json" || item.Path == "settings.local.json") && item.IsFile {
			hasSettings = true
			break
		}
	}

	// Determine the tree branch character for output-styles
	outputStylesBranch := "└── "
	if hasSettings {
		outputStylesBranch = "├── "
	}

	if hasOutputStyles {
		outputStylesCheckbox := "□"
		if outputStylesDir != nil && outputStylesDir.Selected {
			outputStylesCheckbox = "✓"
		}
		outputStylesHighlighted := false
		// Find output-styles directory index for cursor highlighting
		for i, item := range m.treeItems {
			if item.Path == "__DIR__output-styles" && i == m.cursor {
				outputStylesHighlighted = true
				break
			}
		}
		var outputStylesLine string
		if outputStylesHighlighted {
			// For cursor highlighting: purple checkbox, green text
			outputStylesLine = fmt.Sprintf("%s %s%s", checkboxStyle.Render(outputStylesCheckbox), normalStyle.Render(outputStylesBranch), greenTextStyle.Render("output-styles/"))
		} else {
			// Default styling: default checkbox, default text color
			outputStylesLine = fmt.Sprintf("%s %s", normalStyle.Render(outputStylesCheckbox), normalStyle.Render(outputStylesBranch+"output-styles/"))
		}
		s.WriteString(outputStylesLine)
		s.WriteString("\n")

		for i, item := range m.treeItems {
			if strings.HasPrefix(item.Path, "output-styles/") && item.IsFile {
				checkbox := "□"
				if item.Selected {
					checkbox = "✓"
				}

				displayName := strings.TrimPrefix(item.Name, "  ")
				updateIndicator := ""
				if item.IsUpdate {
					updateIndicator = updateStyle.Render(" (will update)")
				}

				if i == m.cursor {
					// For cursor highlighting: purple checkbox, green text
					line := fmt.Sprintf("%s %s%s%s", checkboxStyle.Render(checkbox), normalStyle.Render("│   ├── "), greenTextStyle.Render(displayName), updateIndicator)
					s.WriteString(line)
				} else {
					// Default styling: default checkbox, default text color
					line := fmt.Sprintf("%s %s%s", normalStyle.Render(checkbox), normalStyle.Render("│   ├── "+displayName), updateIndicator)
					s.WriteString(line)
				}
				s.WriteString("\n")
			}
		}
	}

	// Settings files (inside .claude directory)
	if hasSettings {
		settingsCount := 0
		var settingsFiles []struct {
			item  TreeItem
			index int
		}

		// Collect settings files
		for i, item := range m.treeItems {
			if (item.Path == "settings.json" || item.Path == "settings.local.json") && item.IsFile {
				settingsFiles = append(settingsFiles, struct {
					item  TreeItem
					index int
				}{item, i})
				settingsCount++
			}
		}

		// Render settings files with proper tree structure
		for idx, sf := range settingsFiles {
			checkbox := "□"
			if sf.item.Selected {
				checkbox = "✓"
			}

			updateIndicator := ""
			if sf.item.IsUpdate {
				updateIndicator = updateStyle.Render(" (will update)")
			}

			// Use └── for the last settings file, ├── for others
			branch := "├── "
			if idx == len(settingsFiles)-1 {
				branch = "└── "
			}

			if sf.index == m.cursor {
				// For cursor highlighting: purple checkbox, green text
				line := fmt.Sprintf("%s %s%s%s", checkboxStyle.Render(checkbox), normalStyle.Render(branch), greenTextStyle.Render(sf.item.Name), updateIndicator)
				s.WriteString(line)
			} else {
				// Default styling: default checkbox, default text color
				line := fmt.Sprintf("%s %s%s", normalStyle.Render(checkbox), normalStyle.Render(branch+sf.item.Name), updateIndicator)
				s.WriteString(line)
			}
			s.WriteString("\n")
		}
	}

	return s.String()
}

// buildNonInteractiveTree creates a file tree format showing only selected files (for summary)
func (m EnhancedFileSelectionModel) buildNonInteractiveTree() string {
	var s strings.Builder

	// Use themed styles for consistency (matching selection views)
	styles := GetStyles()
	normalStyle := styles.Normal                                         // Default text color (light gray)
	updateStyle := styles.Warning                                        // Orange/Peach for update indicators

	// Root: .claude/ with left-aligned checkmark
	s.WriteString(normalStyle.Render("✓ .claude/"))
	s.WriteString("\n")

	// Agents subsection (only show selected)
	hasSelectedAgents := false
	for _, item := range m.treeItems {
		if strings.HasPrefix(item.Path, "agents/") && item.IsFile && item.Selected {
			hasSelectedAgents = true
			break
		}
	}

	// Commands subsection (only show selected)
	hasSelectedCommands := false
	for _, item := range m.treeItems {
		if strings.HasPrefix(item.Path, "commands/") && item.IsFile && item.Selected {
			hasSelectedCommands = true
			break
		}
	}

	// Output-styles subsection (only show selected)
	hasSelectedOutputStyles := false
	for _, item := range m.treeItems {
		if strings.HasPrefix(item.Path, "output-styles/") && item.IsFile && item.Selected {
			hasSelectedOutputStyles = true
			break
		}
	}

	// Check if we have selected settings files
	hasSelectedSettings := false
	for _, item := range m.treeItems {
		if (item.Path == "settings.json" || item.Path == "settings.local.json") && item.IsFile && item.Selected {
			hasSelectedSettings = true
			break
		}
	}

	if hasSelectedAgents {
		agentsBranch := "├── "
		if !hasSelectedCommands && !hasSelectedOutputStyles && !hasSelectedSettings {
			agentsBranch = "└── "
		}
		s.WriteString(normalStyle.Render("✓ " + agentsBranch + "agents/"))
		s.WriteString("\n")

		for _, item := range m.treeItems {
			if strings.HasPrefix(item.Path, "agents/") && item.IsFile && item.Selected {
				displayName := strings.TrimPrefix(item.Name, "  ")
				updateIndicator := ""
				if item.IsUpdate {
					updateIndicator = updateStyle.Render(" (will update)")
				}

				treeBranch := "│   ├── "
				if !hasSelectedCommands && !hasSelectedOutputStyles && !hasSelectedSettings {
					treeBranch = "    ├── "
				}
				s.WriteString(normalStyle.Render("✓ " + treeBranch + displayName + updateIndicator))
				s.WriteString("\n")
			}
		}
	}

	if hasSelectedCommands {
		commandsBranch := "├── "
		if !hasSelectedOutputStyles && !hasSelectedSettings {
			commandsBranch = "└── "
		}
		s.WriteString(normalStyle.Render("✓ " + commandsBranch + "commands/"))
		s.WriteString("\n")

		for _, item := range m.treeItems {
			if strings.HasPrefix(item.Path, "commands/") && item.IsFile && item.Selected {
				displayName := strings.TrimPrefix(item.Name, "  ")
				updateIndicator := ""
				if item.IsUpdate {
					updateIndicator = updateStyle.Render(" (will update)")
				}

				treeBranch := "│   ├── "
				if !hasSelectedOutputStyles && !hasSelectedSettings {
					treeBranch = "    ├── "
				}
				s.WriteString(normalStyle.Render("✓ " + treeBranch + displayName + updateIndicator))
				s.WriteString("\n")
			}
		}
	}

	if hasSelectedOutputStyles {
		outputStylesBranch := "├── "
		if !hasSelectedSettings {
			outputStylesBranch = "└── "
		}
		s.WriteString(normalStyle.Render("✓ " + outputStylesBranch + "output-styles/"))
		s.WriteString("\n")

		for _, item := range m.treeItems {
			if strings.HasPrefix(item.Path, "output-styles/") && item.IsFile && item.Selected {
				displayName := strings.TrimPrefix(item.Name, "  ")
				updateIndicator := ""
				if item.IsUpdate {
					updateIndicator = updateStyle.Render(" (will update)")
				}

				treeBranch := "│   ├── "
				if !hasSelectedSettings {
					treeBranch = "    ├── "
				}
				s.WriteString(normalStyle.Render("✓ " + treeBranch + displayName + updateIndicator))
				s.WriteString("\n")
			}
		}
	}

	// Settings files (inside .claude directory, only if selected)
	if hasSelectedSettings {
		var selectedSettingsFiles []TreeItem

		// Collect selected settings files
		for _, item := range m.treeItems {
			if (item.Path == "settings.json" || item.Path == "settings.local.json") && item.IsFile && item.Selected {
				selectedSettingsFiles = append(selectedSettingsFiles, item)
			}
		}

		// Render selected settings files with proper tree structure
		for idx, item := range selectedSettingsFiles {
			updateIndicator := ""
			if item.IsUpdate {
				updateIndicator = updateStyle.Render(" (will update)")
			}

			// Use └── for the last settings file, ├── for others
			branch := "├── "
			if idx == len(selectedSettingsFiles)-1 {
				branch = "└── "
			}

			s.WriteString(normalStyle.Render("✓ " + branch + item.Name + updateIndicator))
			s.WriteString("\n")
		}
	}

	return s.String()
}


