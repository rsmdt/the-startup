package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

/*
State Machine Architecture

This package implements a dual-workflow state machine supporting both install and uninstall operations:

Install Workflow:
  StateClaudePath → StateStartupPath → StateFileSelection → StateComplete
                ↓         ↓                ↓
            StateError ← StateError ← StateError

Uninstall Workflow:
  StateUninstallStartupPath → StateUninstallClaudePath → StateUninstallFileSelection → StateUninstallComplete
                         ↓                          ↓                            ↓
                     StateError ← StateError ← StateError

Navigation:
- ESC key allows backward navigation in both workflows
- Error states can return to appropriate starting points
- Both workflows support quit from initial states

The MainModel orchestrates the state machine using OperationMode to determine workflow behavior.
Sub-models handle individual state UI rendering and user interaction with mode-specific adaptations.
*/

// InstallerState represents the current state of the installer or uninstaller
type InstallerState int

const (
	// Install workflow states
	StateClaudePath InstallerState = iota
	StateStartupPath
	StateFileSelection
	StateComplete
	StateError

	// Uninstall workflow states
	StateUninstallStartupPath
	StateUninstallClaudePath
	StateUninstallFileSelection
	StateUninstallComplete
)

// OperationMode determines whether we're installing or uninstalling
type OperationMode int

const (
	ModeInstall OperationMode = iota
	ModeUninstall
)

// String returns the string representation of the operation mode
func (m OperationMode) String() string {
	switch m {
	case ModeInstall:
		return "Install"
	case ModeUninstall:
		return "Uninstall"
	default:
		return "Unknown"
	}
}

// String returns the string representation of the state
func (s InstallerState) String() string {
	switch s {
	// Install workflow states
	case StateClaudePath:
		return "Claude Path Selection"
	case StateStartupPath:
		return "Startup Path Selection"
	case StateFileSelection:
		return "File Selection"
	case StateComplete:
		return "Complete"
	case StateError:
		return "Error"

	// Uninstall workflow states
	case StateUninstallStartupPath:
		return "Uninstall Startup Path Selection"
	case StateUninstallClaudePath:
		return "Uninstall Claude Path Selection"
	case StateUninstallFileSelection:
		return "Uninstall File Selection"
	case StateUninstallComplete:
		return "Uninstall Complete"
	default:
		return "Unknown"
	}
}

// StateTransition handles transitions between states
type StateTransition struct {
	From InstallerState
	To   InstallerState
}

// ValidTransitions defines allowed state transitions
var ValidTransitions = map[StateTransition]bool{
	// Install workflow - Forward transitions
	{StateClaudePath, StateStartupPath}:   true,
	{StateStartupPath, StateFileSelection}: true,
	{StateFileSelection, StateComplete}:   true,
	{StateFileSelection, StateError}:      true,

	// Install workflow - Backward transitions (ESC)
	{StateStartupPath, StateClaudePath}:   true,
	{StateFileSelection, StateStartupPath}: true,

	// Uninstall workflow - Forward transitions
	{StateUninstallStartupPath, StateUninstallClaudePath}:   true,
	{StateUninstallClaudePath, StateUninstallFileSelection}: true,
	{StateUninstallFileSelection, StateUninstallComplete}:   true,
	{StateUninstallStartupPath, StateError}:                 true,
	{StateUninstallClaudePath, StateError}:                  true,
	{StateUninstallFileSelection, StateError}:               true,

	// Uninstall workflow - Backward transitions (ESC)
	{StateUninstallClaudePath, StateUninstallStartupPath}:     true,
	{StateUninstallFileSelection, StateUninstallClaudePath}:   true,

	// Error recovery - can return to appropriate starting points
	{StateError, StateClaudePath}:             true,
	{StateError, StateUninstallStartupPath}:    true,
}

// IsValidTransition checks if a state transition is allowed
func IsValidTransition(from, to InstallerState) bool {
	return ValidTransitions[StateTransition{from, to}]
}

// ProgressiveDisclosureRenderer creates a consistent header showing previous selections
type ProgressiveDisclosureRenderer struct {
	styles Styles
}

// NewProgressiveDisclosureRenderer creates a new renderer
func NewProgressiveDisclosureRenderer() *ProgressiveDisclosureRenderer {
	return &ProgressiveDisclosureRenderer{
		styles: GetStyles(),
	}
}

// RenderSelections creates the progressive disclosure header
func (p *ProgressiveDisclosureRenderer) RenderSelections(
	tool, path string,
	filesCount int,
) string {
	return p.RenderSelectionsWithMode(tool, path, filesCount, ModeInstall)
}

// RenderSelectionsWithMode creates the progressive disclosure header with operation mode
func (p *ProgressiveDisclosureRenderer) RenderSelectionsWithMode(
	tool, path string,
	filesCount int,
	mode OperationMode,
) string {
	if tool == "" && path == "" {
		return ""
	}

	var sections []string

	// Tool selection
	if tool != "" {
		sections = append(sections, fmt.Sprintf("Tool: %s", tool))
	}

	// Path selection
	if path != "" {
		// Show abbreviated path for display
		displayPath := path
		if strings.Contains(path, ".config/the-startup") {
			displayPath = "~/.config/the-startup (global)"
		} else if strings.Contains(path, ".the-startup") {
			displayPath = ".the-startup (local)"
		}
		sections = append(sections, fmt.Sprintf("Path: %s", displayPath))
	}

	// Files selection
	if filesCount > 0 {
		if mode == ModeUninstall {
			sections = append(sections, fmt.Sprintf("Files: %d to remove", filesCount))
		} else {
			sections = append(sections, fmt.Sprintf("Files: %d selected", filesCount))
		}
	}

	if len(sections) == 0 {
		return ""
	}

	// Create bordered box with selections
	content := strings.Join(sections, " • ")

	var title string
	if mode == ModeUninstall {
		title = "The (Agentic) Startup Uninstallation"
	} else {
		title = "The (Agentic) Startup Installation"
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.styles.Info.GetForeground()).
		Padding(0, 1).
		MarginBottom(1).
		Render(title+"\n"+content) + "\n"
}

// RenderTitle creates a consistent title bar
func (p *ProgressiveDisclosureRenderer) RenderTitle(title string) string {
	return p.styles.Title.Render(title) + "\n\n"
}

// RenderHelp creates consistent help text
func (p *ProgressiveDisclosureRenderer) RenderHelp(helpText string) string {
	return "\n" + p.styles.Help.Render(helpText)
}

// RenderError creates consistent error display
func (p *ProgressiveDisclosureRenderer) RenderError(err error, context string) string {
	return p.styles.Error.Render(fmt.Sprintf("❌ Error %s: %v", context, err)) + "\n"
}

// GetChoiceStyle returns the appropriate style for a choice
func (p *ProgressiveDisclosureRenderer) GetChoiceStyle(isCursor bool, isSelected bool) lipgloss.Style {
	if isCursor {
		if isSelected {
			return p.styles.CursorLine.Copy().Background(p.styles.Info.GetForeground())
		}
		return p.styles.CursorLine
	}

	if isSelected {
		return p.styles.Info
	}

	return p.styles.Normal
}

// RenderChoice formats a choice with cursor and selection indicators
func (p *ProgressiveDisclosureRenderer) RenderChoice(text string, isCursor, isSelected bool) string {
	return p.RenderChoiceWithMultiSelect(text, isCursor, isSelected, true) // Default to multi-select (backward compatibility)
}

// RenderChoiceWithMultiSelect formats a choice with optional selection indicators
func (p *ProgressiveDisclosureRenderer) RenderChoiceWithMultiSelect(text string, isCursor, isSelected, isMultiSelect bool) string {
	var prefix string
	if isCursor {
		prefix = p.styles.Cursor.Render("> ")
	} else {
		prefix = "  "
	}

	style := p.GetChoiceStyle(isCursor, isSelected)

	// Only show circles for multi-select contexts
	if isMultiSelect {
		if isSelected {
			return prefix + style.Render("● "+text)
		}
		return prefix + style.Render("○ "+text)
	} else {
		// Single-select: just show text without circles
		return prefix + style.Render(text)
	}
}

// RenderInstallationSummary creates a summary of what will be installed
func (p *ProgressiveDisclosureRenderer) RenderInstallationSummary(
	tool, path string,
	files []string,
	updatingFiles []string,
) string {
	var summary strings.Builder

	summary.WriteString(p.styles.Info.Render("📋 Installation Summary"))
	summary.WriteString("\n\n")

	// Basic configuration
	summary.WriteString(fmt.Sprintf("Tool: %s\n", tool))
	summary.WriteString(fmt.Sprintf("Path: %s\n", path))

	summary.WriteString(fmt.Sprintf("Files to install: %d\n", len(files)))

	// Group files by component
	filesByComponent := make(map[string][]string)
	for _, file := range files {
		parts := strings.Split(file, "/")
		if len(parts) >= 2 {
			component := parts[0]
			fileName := parts[1]
			filesByComponent[component] = append(filesByComponent[component], fileName)
		}
	}

	summary.WriteString("\n")
	for _, component := range []string{"agents", "commands", "hooks", "templates"} {
		if componentFiles, ok := filesByComponent[component]; ok && len(componentFiles) > 0 {
			summary.WriteString(fmt.Sprintf("%s/ (%d files)\n", component, len(componentFiles)))
			for i, file := range componentFiles {
				if i < 3 { // Show first 3 files
					summary.WriteString(fmt.Sprintf("  • %s\n", file))
				} else if i == 3 && len(componentFiles) > 3 {
					summary.WriteString(fmt.Sprintf("  ... and %d more\n", len(componentFiles)-3))
					break
				}
			}
		}
	}

	// File updates
	if len(updatingFiles) > 0 {
		summary.WriteString(fmt.Sprintf("\n%s Files to update: %d\n",
			p.styles.Warning.Render("⚠"), len(updatingFiles)))
	}

	return summary.String() + "\n"
}

// RenderUninstallSummary creates a summary of what will be removed
func (p *ProgressiveDisclosureRenderer) RenderUninstallSummary(
	installPath, claudePath string,
	filesToRemove []string,
	componentsToRemove map[string][]string,
) string {
	var summary strings.Builder

	summary.WriteString(p.styles.Warning.Render("🗑️  Uninstall Preview"))
	summary.WriteString("\n\n")

	// Installation paths
	summary.WriteString(fmt.Sprintf("Installation Path: %s\n", installPath))
	summary.WriteString(fmt.Sprintf("Claude Path: %s\n", claudePath))
	summary.WriteString(fmt.Sprintf("Total files to remove: %d\n", len(filesToRemove)))

	// Group files by component
	summary.WriteString("\n")
	for _, component := range []string{"agents", "commands", "output-styles", "templates", "rules", "bin"} {
		if componentFiles, ok := componentsToRemove[component]; ok && len(componentFiles) > 0 {
			summary.WriteString(fmt.Sprintf("%s/ (%d files)\n", component, len(componentFiles)))
			for i, file := range componentFiles {
				if i < 3 { // Show first 3 files
					summary.WriteString(fmt.Sprintf("  • %s\n", file))
				} else if i == 3 && len(componentFiles) > 3 {
					summary.WriteString(fmt.Sprintf("  ... and %d more\n", len(componentFiles)-3))
					break
				}
			}
		}
	}

	return summary.String() + "\n"
}

// RenderUninstallWarning creates a warning about destructive operations
func (p *ProgressiveDisclosureRenderer) RenderUninstallWarning() string {
	warning := p.styles.Error.Render("⚠️  WARNING: This action cannot be undone!")
	explanation := p.styles.Normal.Render("This will permanently remove all The (Agentic) Startup files from your system.\nYour Claude configuration will be updated to remove hooks and settings.\nYou can reinstall later, but any customizations will be lost.")
	
	return fmt.Sprintf("%s\n\n%s\n", warning, explanation)
}

// RenderUninstallHeader creates the uninstall workflow header
func (p *ProgressiveDisclosureRenderer) RenderUninstallHeader(currentSelections string) string {
	if currentSelections == "" {
		currentSelections = "Uninstalling The (Agentic) Startup"
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.styles.Warning.GetForeground()).
		Padding(0, 1).
		MarginBottom(1).
		Render(currentSelections) + "\n"
}

// RenderSelectedPaths creates a consistent display of selected paths across all screens
func (p *ProgressiveDisclosureRenderer) RenderSelectedPaths(claudePath, startupPath string, mode OperationMode) string {
	if claudePath == "" && startupPath == "" {
		return ""
	}

	var s strings.Builder

	// Title
	if mode == ModeUninstall {
		s.WriteString(p.styles.Warning.Render("📂 Selected Paths:"))
	} else {
		s.WriteString(p.styles.Info.Render("📂 Selected Paths:"))
	}
	s.WriteString("\n")

	// Format paths consistently
	home := os.Getenv("HOME")

	if claudePath != "" {
		displayPath := claudePath
		if home != "" && strings.HasPrefix(claudePath, home) {
			displayPath = "~" + strings.TrimPrefix(claudePath, home)
		}
		s.WriteString(p.styles.Normal.Render(fmt.Sprintf("  🔧 Claude:  %s", displayPath)))
		s.WriteString("\n")
	}

	if startupPath != "" {
		displayPath := startupPath
		if home != "" && strings.HasPrefix(startupPath, home) {
			displayPath = "~" + strings.TrimPrefix(startupPath, home)
		}
		s.WriteString(p.styles.Normal.Render(fmt.Sprintf("  🚀 Startup: %s", displayPath)))
		s.WriteString("\n")
	}

	return s.String() + "\n"
}
