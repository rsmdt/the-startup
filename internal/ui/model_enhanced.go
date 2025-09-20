package ui

import (
	"embed"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/installer"
)

// RunMainInstallerWithInteractiveSelection runs the installer with Huh-based component selection
func RunMainInstallerWithInteractiveSelection(claudeAssets, startupAssets *embed.FS, local, yes bool) error {
	flags := InstallFlags{
		Local: local,
		Yes:   yes,
	}

	// Step 1: Run Huh-based selection (if interactive mode)
	var selectedFiles []string
	var err error

	if !yes && !local {
		// Show the banner first
		fmt.Print(AppBanner)
		fmt.Println()
		fmt.Println("🚀 Welcome to The (Agentic) Startup Interactive Installer!")

		selectedFiles, err = RunHuhSelection(claudeAssets, startupAssets, flags)
		if err != nil {
			return fmt.Errorf("interactive selection failed: %w", err)
		}

		if len(selectedFiles) == 0 {
			fmt.Println("No components selected. Installation cancelled.")
			return nil
		}

		fmt.Printf("✓ Selected %d components for installation\n", len(selectedFiles))
		fmt.Println()
	} else {
		// For --yes or --local flags, select all files
		selectedFiles = getAllAvailableFiles(claudeAssets, startupAssets)
	}

	// Step 2: Run modified BubbleTea UI for path selection and confirmation
	model := NewMainModelWithSelectedFiles(claudeAssets, startupAssets, flags, selectedFiles)

	program := tea.NewProgram(model)
	_, err = program.Run()
	return err
}

// NewMainModelWithSelectedFiles creates a main model with pre-selected files
func NewMainModelWithSelectedFiles(claudeAssets, startupAssets *embed.FS, flags InstallFlags, selectedFiles []string) *MainModel {
	installerInstance := installer.New(claudeAssets, startupAssets)

	// Set selected files in installer
	if len(selectedFiles) > 0 {
		installerInstance.SetSelectedFiles(selectedFiles)
	}

	var initialState InstallerState
	initialState = StateStartupPath

	m := &MainModel{
		state:            initialState,
		mode:             ModeInstall,
		installer:        installerInstance,
		claudeAssets:     claudeAssets,
		startupAssets:    startupAssets,
		width:            80,
		height:           24,
		flags:            flags,
		selectedFiles:    selectedFiles,
		startupPathModel: NewStartupPathModelWithMode("", ModeInstall),
		errorModel:       NewErrorModel(nil, ""),
	}

	return m
}

// Enhanced transitionToState that uses preselected files
func (m *MainModel) transitionToStateEnhanced(newState InstallerState) {
	m.state = newState

	// Initialize or reset sub-models as needed
	switch newState {
	case StateStartupPath:
		m.startupPathModel = m.startupPathModel.Reset()

	case StateClaudePath:
		m.claudePathModel = NewClaudePathModelWithMode(m.startupPath, m.mode)

	case StateFileSelection:
		// Use preselected files if available
		if len(m.selectedFiles) > 0 {
			// Create a regular model but with the preselected files
			original := NewFileSelectionModelWithMode(
				"claude-code",
				m.claudePath,
				m.installer,
				m.claudeAssets,
				m.startupAssets,
				m.mode,
			)
			// Override the selected files
			original.selectedFiles = m.selectedFiles
			m.fileSelectionModel = FileSelectionWrapper{FileSelectionModel: original}
		} else {
			// Use original logic for classic mode
			original := NewFileSelectionModelWithMode(
				"claude-code",
				m.claudePath,
				m.installer,
				m.claudeAssets,
				m.startupAssets,
				m.mode,
			)
			m.fileSelectionModel = FileSelectionWrapper{FileSelectionModel: original}
		}

	case StateComplete:
		m.completeModel = NewCompleteModelWithAssets("claude-code", m.installer, m.mode, m.claudeAssets, m.startupAssets, m.selectedFiles)

	case StateError:
		// Error model is set before transition
	}
}

// RunEnhancedInstaller is the new entry point for the enhanced installer
func RunEnhancedInstaller(claudeAssets, startupAssets *embed.FS, local, yes bool) error {
	// Check if user wants interactive mode
	if !yes {
		return RunMainInstallerWithInteractiveSelection(claudeAssets, startupAssets, local, yes)
	}

	// Fall back to original installer for --yes flag
	return RunMainInstallerWithFlags(claudeAssets, startupAssets, local, yes)
}

// UpdateFileSelectionToShowPreselected modifies the file selection to show preselected items
func (m *MainModel) UpdateFileSelectionToShowPreselected() {
	if len(m.selectedFiles) > 0 && m.state == StateFileSelection {
		// This functionality is now handled by the enhanced model creation
		// The selected files are already set when creating the model
	}
}
