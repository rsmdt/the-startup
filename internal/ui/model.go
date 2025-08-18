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

// NewMainModel creates a new main model with composed sub-models
func NewMainModel(claudeAssets, startupAssets *embed.FS) *MainModel {
	installerInstance := installer.New(claudeAssets, startupAssets)

	m := &MainModel{
		state:            StateStartupPath, // Start with startup path selection
		installer:        installerInstance,
		claudeAssets:     claudeAssets,
		startupAssets:    startupAssets,
		width:            80,
		height:           24,
		startupPathModel: NewStartupPathModel(),
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
			if m.state == StateStartupPath {
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
	case StateStartupPath:
		newModel, cmd := m.startupPathModel.Update(msg)
		m.startupPathModel = newModel
		if m.startupPathModel.Ready() {
			path := m.startupPathModel.SelectedPath()
			if path == "CANCEL" {
				return m, tea.Quit
			}
			m.startupPath = path
			m.installer.SetInstallPath(path)
			m.transitionToState(StateClaudePath)
		}
		return m, cmd

	case StateClaudePath:
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
			m.transitionToState(StateFileSelection)
		}
		return m, cmd

	case StateFileSelection:
		newModel, cmd := m.fileSelectionModel.Update(msg)
		m.fileSelectionModel = newModel
		if m.fileSelectionModel.Ready() {
			if m.fileSelectionModel.Confirmed() {
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
			} else {
				m.transitionToState(StateClaudePath)
			}
		}
		return m, cmd

	case StateComplete:
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
			m.transitionToState(StateStartupPath)
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
	case StateStartupPath:
		m.startupPathModel = m.startupPathModel.Reset()

	case StateClaudePath:
		m.claudePathModel = NewClaudePathModel(m.startupPath)

	case StateFileSelection:
		m.fileSelectionModel = NewFileSelectionModel(
			"claude-code", // Always use claude-code
			m.claudePath,
			m.installer,
			m.claudeAssets,
			m.startupAssets,
		)
		m.selectedFiles = m.fileSelectionModel.selectedFiles

	case StateComplete:
		m.completeModel = NewCompleteModel("claude-code", m.installer)

	case StateError:
		// Error model is set before transition
	}
}

// View delegates rendering to the appropriate sub-model
func (m *MainModel) View() string {
	switch m.state {
	case StateStartupPath:
		return m.startupPathModel.View()
	case StateClaudePath:
		return m.claudePathModel.View()
	case StateFileSelection:
		return m.fileSelectionModel.View()
	case StateComplete:
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
	case StateStartupPath:
		return m, tea.Quit
	case StateClaudePath:
		m.transitionToState(StateStartupPath)
	case StateFileSelection:
		m.transitionToState(StateClaudePath)
	case StateComplete:
		return m, tea.Quit
	case StateError:
		m.transitionToState(StateStartupPath)
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
