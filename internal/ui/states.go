package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// InstallerState represents the current state of the installer
type InstallerState int

const (
	StateStartupPath InstallerState = iota
	StateClaudePath
	StateFileSelection
	StateComplete
	StateError
)

// String returns the string representation of the state
func (s InstallerState) String() string {
	switch s {
	case StateStartupPath:
		return "Startup Path Selection"
	case StateClaudePath:
		return "Claude Path Selection"
	case StateFileSelection:
		return "File Selection"
	case StateComplete:
		return "Complete"
	case StateError:
		return "Error"
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
	// Forward transitions
	{StateStartupPath, StateClaudePath}:   true,
	{StateClaudePath, StateFileSelection}: true,
	{StateFileSelection, StateComplete}:   true,
	{StateFileSelection, StateError}:      true,

	// Backward transitions (ESC)
	{StateClaudePath, StateStartupPath}:   true,
	{StateFileSelection, StateClaudePath}: true,

	// Error recovery
	{StateError, StateStartupPath}: true,
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
	if tool == "" {
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
		sections = append(sections, fmt.Sprintf("Files: %d selected", filesCount))
	}

	if len(sections) == 0 {
		return ""
	}

	// Create bordered box with selections
	content := strings.Join(sections, " • ")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.styles.Info.GetForeground()).
		Padding(0, 1).
		MarginBottom(1).
		Render("The (Agentic) Startup Installation\n"+content) + "\n"
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
