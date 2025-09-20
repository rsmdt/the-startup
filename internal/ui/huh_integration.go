package ui

import (
	"embed"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/installer"
)

// RunHuhSelection runs the Huh-based selection process outside of BubbleTea
// This is called before launching the main BubbleTea interface
func RunHuhSelection(claudeAssets, startupAssets *embed.FS, flags InstallFlags) ([]string, error) {
	// Skip interactive selection for certain flags
	if flags.Yes || flags.Local {
		// Use all files for automatic installation
		return getAllAvailableFiles(claudeAssets, startupAssets), nil
	}

	// Create discovery and selector
	discovery := NewAssetDiscovery(claudeAssets, startupAssets)
	selector := NewInteractiveSelector(discovery)

	// Run interactive selection
	selectedItems, err := selector.RunInteractiveSelection()
	if err != nil {
		return nil, fmt.Errorf("interactive selection failed: %w", err)
	}

	// Convert to file paths
	filePaths := discovery.GetSelectedFilePaths(selectedItems)
	return filePaths, nil
}

// getAllAvailableFiles returns all available files from both embedded filesystems
func getAllAvailableFiles(claudeAssets, startupAssets *embed.FS) []string {
	var allFiles []string

	// Discover all items and return their file paths
	discovery := NewAssetDiscovery(claudeAssets, startupAssets)
	items, err := discovery.DiscoverAllItems()
	if err != nil {
		return allFiles // Return empty slice on error
	}

	for _, item := range items {
		allFiles = append(allFiles, item.FilePath)
	}

	return allFiles
}

// PreselectedFileSelectionModel is a simplified model that shows already-selected files
// and just asks for confirmation
type PreselectedFileSelectionModel struct {
	styles        Styles
	installer     *installer.Installer
	claudeAssets  *embed.FS
	startupAssets *embed.FS
	selectedFiles []string
	cursor        int
	choices       []string
	ready         bool
	confirmed     bool
	mode          OperationMode
	customMessage string
}

// NewPreselectedFileSelectionModel creates a model with pre-selected files
func NewPreselectedFileSelectionModel(installer *installer.Installer, claudeAssets, startupAssets *embed.FS, selectedFiles []string, mode OperationMode) PreselectedFileSelectionModel {
	var choices []string
	var message string

	if mode == ModeUninstall {
		choices = []string{
			"Yes, remove everything",
			"No, keep everything as-is",
		}
		message = "Ready to uninstall?"
	} else {
		choices = []string{
			"Yes, install selected components",
			"No, cancel installation",
		}
		if len(selectedFiles) > 0 {
			message = fmt.Sprintf("Ready to install %d selected components?", len(selectedFiles))
		} else {
			message = "No components selected for installation"
			choices = []string{"Go back to selection", "Cancel"}
		}
	}

	return PreselectedFileSelectionModel{
		styles:        GetStyles(),
		installer:     installer,
		claudeAssets:  claudeAssets,
		startupAssets: startupAssets,
		selectedFiles: selectedFiles,
		cursor:        0,
		choices:       choices,
		ready:         false,
		confirmed:     false,
		mode:          mode,
		customMessage: message,
	}
}

// Init initializes the model
func (m PreselectedFileSelectionModel) Init() tea.Cmd {
	return nil
}

// Update handles user input
func (m PreselectedFileSelectionModel) Update(msg tea.Msg) (PreselectedFileSelectionModel, tea.Cmd) {
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
					(m.mode == ModeInstall && choice == "Yes, install selected components") {
					m.confirmed = true
				} else {
					m.confirmed = false
				}
				m.ready = true
			}
		}
	}
	return m, nil
}

// View renders the confirmation interface
func (m PreselectedFileSelectionModel) View() string {
	var s strings.Builder

	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")

	// Show paths
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
	} else {
		s.WriteString(m.styles.Info.Render("Installation Paths:"))
	}
	s.WriteString("\n")
	s.WriteString(m.styles.Normal.Render(fmt.Sprintf("  Startup: %s", startupPath)))
	s.WriteString("\n")
	s.WriteString(m.styles.Normal.Render(fmt.Sprintf("  Claude:  %s", claudePath)))
	s.WriteString("\n\n")

	// Show selection summary if in install mode
	if m.mode == ModeInstall && len(m.selectedFiles) > 0 {
		s.WriteString(m.styles.Info.Render("📋 Selected Components Summary:"))
		s.WriteString("\n\n")

		// Group by type for display
		s.WriteString(m.renderSelectionSummary())
		s.WriteString("\n")
	}

	// Show confirmation message
	s.WriteString(m.styles.Title.Render(m.customMessage))
	s.WriteString("\n\n")

	// Show choices
	for i, option := range m.choices {
		if i == m.cursor {
			s.WriteString(m.styles.Selected.Render("> " + option))
		} else {
			s.WriteString(m.styles.Normal.Render("  " + option))
		}
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(m.styles.Help.Render("Press Enter to confirm • ↑↓ to navigate • Escape to go back"))

	return s.String()
}

// renderSelectionSummary creates a summary of selected components
func (m PreselectedFileSelectionModel) renderSelectionSummary() string {
	discovery := NewAssetDiscovery(m.claudeAssets, m.startupAssets)
	items, err := discovery.DiscoverAllItems()
	if err != nil {
		return "Unable to load component details"
	}

	// Create a map of selected files
	selectedSet := make(map[string]bool)
	for _, file := range m.selectedFiles {
		selectedSet[file] = true
	}

	// Group selected items by type
	categories := make(map[string][]SelectableItem)
	for _, item := range items {
		if selectedSet[item.FilePath] {
			categories[item.Category] = append(categories[item.Category], item)
		}
	}

	var result []string
	for category, categoryItems := range categories {
		if len(categoryItems) == 0 {
			continue
		}

		result = append(result, fmt.Sprintf("▶ %s (%d items)", category, len(categoryItems)))

		// Show first few items
		for i, item := range categoryItems {
			if i < 3 {
				result = append(result, fmt.Sprintf("  • %s", item.Name))
			} else {
				remaining := len(categoryItems) - 3
				result = append(result, fmt.Sprintf("  • ... and %d more", remaining))
				break
			}
		}
	}

	return strings.Join(result, "\n")
}

// Ready returns whether the model is ready to proceed
func (m PreselectedFileSelectionModel) Ready() bool {
	return m.ready
}

// Confirmed returns whether the user confirmed the action
func (m PreselectedFileSelectionModel) Confirmed() bool {
	return m.confirmed
}

// Reset resets the model state
func (m PreselectedFileSelectionModel) Reset() PreselectedFileSelectionModel {
	m.ready = false
	m.cursor = 0
	return m
}