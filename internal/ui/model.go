package ui

import (
	"embed"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/the-startup/the-startup/internal/installer"
)


// MainModel manages the overall installer flow by composing sub-models
type MainModel struct {
	// State management
	state InstallerState
	
	// Dependencies
	installer     *installer.Installer
	agentFiles    *embed.FS
	commandFiles  *embed.FS
	hookFiles     *embed.FS
	templateFiles *embed.FS
	
	// User selections (shared state)
	selectedTool  string
	selectedPath  string
	selectedFiles []string
	
	// Sub-models
	toolSelectionModel ToolSelectionModel
	pathSelectionModel PathSelectionModel
	fileSelectionModel FileSelectionModel
	completeModel      CompleteModel
	errorModel         ErrorModel
	
	// UI state
	width  int
	height int
}

// NewMainModel creates a new main model with composed sub-models
func NewMainModel(agents, commands, hooks, templates *embed.FS) *MainModel {
	installerInstance := installer.New(agents, commands, hooks, templates)
	
	m := &MainModel{
		state:              StateToolSelection, // Start directly with tool selection
		installer:          installerInstance,
		agentFiles:         agents,
		commandFiles:       commands,
		hookFiles:          hooks,
		templateFiles:      templates,
		width:              80,
		height:             24,
		toolSelectionModel: NewToolSelectionModel(),
		errorModel:         NewErrorModel(nil, ""),
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
			if m.state == StateToolSelection {
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
	case StateToolSelection:
		newModel, cmd := m.toolSelectionModel.Update(msg)
		m.toolSelectionModel = newModel
		if m.toolSelectionModel.Ready() {
			if m.toolSelectionModel.Selected() == "Cancel" {
				return m, tea.Quit
			}
			m.selectedTool = m.toolSelectionModel.Selected()
			m.installer.SetTool(m.selectedTool)
			m.transitionToState(StatePathSelection)
		}
		return m, cmd
	
	case StatePathSelection:
		newModel, cmd := m.pathSelectionModel.Update(msg)
		m.pathSelectionModel = newModel
		if m.pathSelectionModel.Ready() {
			path := m.pathSelectionModel.SelectedPath()
			switch path {
			case "CANCEL":
				return m, tea.Quit
			case "CUSTOM":
				m.errorModel = NewErrorModel(fmt.Errorf("custom path input not yet implemented"), "selecting custom path")
				m.transitionToState(StateError)
			default:
				m.selectedPath = path
				if path != "" {
					m.installer.SetInstallPath(path)
				}
				m.transitionToState(StateFileSelection)
			}
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
				m.transitionToState(StatePathSelection)
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
			m.transitionToState(StateToolSelection)
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
	case StateToolSelection:
		m.toolSelectionModel = m.toolSelectionModel.Reset()
	
	case StatePathSelection:
		m.pathSelectionModel = NewPathSelectionModel(m.selectedTool)
	
	case StateFileSelection:
		m.fileSelectionModel = NewFileSelectionModel(
			m.selectedTool,
			m.selectedPath,
			m.installer,
			m.agentFiles,
			m.commandFiles,
			m.hookFiles,
			m.templateFiles,
		)
		m.selectedFiles = m.fileSelectionModel.selectedFiles
	
	case StateComplete:
		m.completeModel = NewCompleteModel(m.selectedTool, m.installer)
	
	case StateError:
		// Error model is set before transition
	}
}

// View delegates rendering to the appropriate sub-model
func (m *MainModel) View() string {
	switch m.state {
	case StateToolSelection:
		return m.toolSelectionModel.View()
	case StatePathSelection:
		return m.pathSelectionModel.View()
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
	case StateToolSelection:
		return m, tea.Quit
	case StatePathSelection:
		m.transitionToState(StateToolSelection)
	case StateFileSelection:
		m.transitionToState(StatePathSelection)
	case StateComplete:
		return m, tea.Quit
	case StateError:
		m.transitionToState(StateToolSelection) // Go back to tool selection instead of welcome
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
	tempModel := NewFileSelectionModel(m.selectedTool, m.selectedPath, m.installer, m.agentFiles, m.commandFiles, m.hookFiles, m.templateFiles)
	return tempModel.getAllAvailableFiles()
}

// RunMainInstaller starts the installation UI using the MainModel
func RunMainInstaller(agents, commands, hooks, templates *embed.FS) error {
	model := NewMainModel(agents, commands, hooks, templates)
	
	program := tea.NewProgram(model)
	
	_, err := program.Run()
	return err
}

