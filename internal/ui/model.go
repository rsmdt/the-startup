package ui

import (
	"embed"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/installer"
)

// MainModel manages the overall installer flow by composing sub-models
type MainModel struct {
	// State management
	state InstallerState
	mode  OperationMode

	// Dependencies
	installer     *installer.Installer
	claudeAssets  *embed.FS
	startupAssets *embed.FS

	// User selections (shared state)
	startupPath   string
	claudePath    string
	selectedFiles []string

	// Sub-models
	startupPathModel   StartupPathModel
	claudePathModel    ClaudePathModel
	fileSelectionModel FileSelectionModel
	completeModel      CompleteModel
	errorModel         ErrorModel

	// UI state
	width  int
	height int
}

// NewMainModel creates a new main model with composed sub-models for installation
func NewMainModel(claudeAssets, startupAssets *embed.FS) *MainModel {
	return NewMainModelWithMode(claudeAssets, startupAssets, ModeInstall)
}

// NewMainUninstallModel creates a new main model for uninstallation
func NewMainUninstallModel(claudeAssets, startupAssets *embed.FS) *MainModel {
	return NewMainModelWithMode(claudeAssets, startupAssets, ModeUninstall)
}

// NewMainModelWithMode creates a new main model with composed sub-models for specified mode
func NewMainModelWithMode(claudeAssets, startupAssets *embed.FS, mode OperationMode) *MainModel {
	installerInstance := installer.New(claudeAssets, startupAssets)

	var initialState InstallerState
	if mode == ModeUninstall {
		initialState = StateUninstallStartupPath
	} else {
		initialState = StateStartupPath
	}

	m := &MainModel{
		state:            initialState,
		mode:             mode,
		installer:        installerInstance,
		claudeAssets:     claudeAssets,
		startupAssets:    startupAssets,
		width:            80,
		height:           24,
		startupPathModel: NewStartupPathModelWithMode(mode),
		errorModel:       NewErrorModel(nil, ""),
	}

	return m
}

// Init initializes the main model
func (m *MainModel) Init() tea.Cmd {
	return nil
}

// Update handles all messages and delegates to sub-models
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global keys first
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c", "q":
			if m.state == StateStartupPath || m.state == StateUninstallStartupPath {
				return m, tea.Quit
			}
			// In other states, treat as ESC (go back)
			return m.handleBack()

		case "esc":
			return m.handleBack()
		}
	}

	// Handle window size
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = sizeMsg.Width
		m.height = sizeMsg.Height
		return m, nil
	}

	// Delegate to appropriate sub-model based on state
	switch m.state {
	case StateStartupPath, StateUninstallStartupPath:
		newModel, cmd := m.startupPathModel.Update(msg)
		m.startupPathModel = newModel
		if m.startupPathModel.Ready() {
			path := m.startupPathModel.SelectedPath()
			if path == "CANCEL" {
				return m, tea.Quit
			}
			m.startupPath = path
			m.installer.SetInstallPath(path)
			if m.mode == ModeUninstall {
				m.transitionToState(StateUninstallClaudePath)
			} else {
				m.transitionToState(StateClaudePath)
			}
		}
		return m, cmd

	case StateClaudePath, StateUninstallClaudePath:
		newModel, cmd := m.claudePathModel.Update(msg)
		m.claudePathModel = newModel
		if m.claudePathModel.Ready() {
			path := m.claudePathModel.SelectedPath()
			if path == "CANCEL" {
				return m, tea.Quit
			}
			m.claudePath = path
			// Set Claude path in installer
			if m.installer != nil {
				// We need to add a method to set Claude path
				m.installer.SetClaudePath(path)
			}
			m.installer.SetTool("claude-code") // Always use claude-code
			if m.mode == ModeUninstall {
				m.transitionToState(StateUninstallFileSelection)
			} else {
				m.transitionToState(StateFileSelection)
			}
		}
		return m, cmd

	case StateFileSelection, StateUninstallFileSelection:
		newModel, cmd := m.fileSelectionModel.Update(msg)
		m.fileSelectionModel = newModel
		if m.fileSelectionModel.Ready() {
			if m.fileSelectionModel.Confirmed() {
				if m.mode == ModeUninstall {
					// Perform uninstallation synchronously
					err := m.performUninstall()
					if err != nil {
						m.errorModel = NewErrorModel(err, "during uninstallation")
						m.transitionToState(StateError)
					} else {
						m.transitionToState(StateUninstallComplete)
						return m, m.completeModel.Init()
					}
				} else {
					// Perform installation synchronously
					err := m.installer.Install()
					if err != nil {
						m.errorModel = NewErrorModel(err, "during installation")
						m.transitionToState(StateError)
					} else {
						m.transitionToState(StateComplete)
						// Return the Init command from the CompleteModel to start the auto-exit timer
						return m, m.completeModel.Init()
					}
				}
			} else {
				if m.mode == ModeUninstall {
					m.transitionToState(StateUninstallClaudePath)
				} else {
					m.transitionToState(StateClaudePath)
				}
			}
		}
		return m, cmd

	case StateComplete, StateUninstallComplete:
		newModel, cmd := m.completeModel.Update(msg)
		m.completeModel = newModel
		if m.completeModel.Ready() {
			return m, tea.Quit
		}
		return m, cmd

	case StateError:
		newModel, cmd := m.errorModel.Update(msg)
		m.errorModel = newModel
		if m.errorModel.Ready() {
			if m.mode == ModeUninstall {
				m.transitionToState(StateUninstallStartupPath)
			} else {
				m.transitionToState(StateStartupPath)
			}
		}
		return m, cmd
	}

	return m, nil
}

// transitionToState handles state transitions and initializes sub-models
func (m *MainModel) transitionToState(newState InstallerState) {
	m.state = newState

	// Initialize or reset sub-models as needed
	switch newState {
	case StateStartupPath, StateUninstallStartupPath:
		m.startupPathModel = m.startupPathModel.Reset()

	case StateClaudePath, StateUninstallClaudePath:
		m.claudePathModel = NewClaudePathModelWithMode(m.startupPath, m.mode)

	case StateFileSelection, StateUninstallFileSelection:
		m.fileSelectionModel = NewFileSelectionModelWithMode(
			"claude-code", // Always use claude-code
			m.claudePath,
			m.installer,
			m.claudeAssets,
			m.startupAssets,
			m.mode,
		)
		m.selectedFiles = m.fileSelectionModel.selectedFiles

	case StateComplete, StateUninstallComplete:
		m.completeModel = NewCompleteModelWithMode("claude-code", m.installer, m.mode)

	case StateError:
		// Error model is set before transition
	}
}

// View delegates rendering to the appropriate sub-model
func (m *MainModel) View() string {
	switch m.state {
	case StateStartupPath, StateUninstallStartupPath:
		return m.startupPathModel.View()
	case StateClaudePath, StateUninstallClaudePath:
		return m.claudePathModel.View()
	case StateFileSelection, StateUninstallFileSelection:
		return m.fileSelectionModel.View()
	case StateComplete, StateUninstallComplete:
		return m.completeModel.View()
	case StateError:
		return m.errorModel.View()
	default:
		return "Unknown state"
	}
}

// handleBack processes back navigation (ESC key)
func (m *MainModel) handleBack() (tea.Model, tea.Cmd) {
	switch m.state {
	case StateStartupPath, StateUninstallStartupPath:
		return m, tea.Quit
	case StateClaudePath:
		m.transitionToState(StateStartupPath)
	case StateUninstallClaudePath:
		m.transitionToState(StateUninstallStartupPath)
	case StateFileSelection:
		m.transitionToState(StateClaudePath)
	case StateUninstallFileSelection:
		m.transitionToState(StateUninstallClaudePath)
	case StateComplete, StateUninstallComplete:
		return m, tea.Quit
	case StateError:
		if m.mode == ModeUninstall {
			m.transitionToState(StateUninstallStartupPath)
		} else {
			m.transitionToState(StateStartupPath)
		}
	default:
		return m, nil
	}
	return m, nil
}

// getAllAvailableFiles returns all files that would be installed (for testing)
func (m *MainModel) getAllAvailableFiles() []string {
	if m.fileSelectionModel.selectedFiles != nil {
		return m.fileSelectionModel.selectedFiles
	}
	// Create a temporary file selection model to get the files
	tempModel := NewFileSelectionModel("claude-code", m.claudePath, m.installer, m.claudeAssets, m.startupAssets)
	return tempModel.getAllAvailableFiles()
}

// RunMainInstaller starts the installation UI using the MainModel
func RunMainInstaller(claudeAssets, startupAssets *embed.FS) error {
	model := NewMainModel(claudeAssets, startupAssets)

	program := tea.NewProgram(model)

	_, err := program.Run()
	return err
}

// RunMainUninstaller starts the uninstallation UI using the MainModel
func RunMainUninstaller(claudeAssets, startupAssets *embed.FS) error {
	model := NewMainUninstallModel(claudeAssets, startupAssets)

	program := tea.NewProgram(model)

	_, err := program.Run()
	return err
}
