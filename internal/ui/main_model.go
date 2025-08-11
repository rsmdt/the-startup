package ui

import (
	"embed"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/the-startup/the-startup/internal/installer"
)

// MainModel manages the overall installer flow using composable views
type MainModel struct {
	context     *Context
	currentView ViewModel
	state       InstallerState
}

// NewMainModel creates a new main model with composable views
func NewMainModel(agents, commands, hooks, templates *embed.FS) *MainModel {
	context := &Context{
		Installer:     installer.New(agents, commands, hooks, templates),
		AgentFiles:    agents,
		CommandFiles:  commands,
		HookFiles:     hooks,
		TemplateFiles: templates,
		Styles:        GetStyles(),
		Renderer:      NewProgressiveDisclosureRenderer(),
		Width:         80,
		Height:        24,
	}

	main := &MainModel{
		context: context,
		state:   StateWelcome,
	}

	// Initialize with welcome view
	main.currentView = NewWelcomeModel(context)
	return main
}

// Init initializes the main model
func (m *MainModel) Init() tea.Cmd {
	return m.currentView.Init()
}

// Update handles messages and routes them to appropriate views
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update context dimensions
		m.context.Width = msg.Width
		m.context.Height = msg.Height
		
		// Pass to current view
		var cmd tea.Cmd
		m.currentView, cmd = m.currentView.Update(msg)
		return m, cmd

	case ViewTransitionMsg:
		return m.transitionToView(msg.NextView, msg.Data)

	case SelectionMadeMsg:
		// Update context with user selections
		if msg.Tool != "" {
			m.context.SelectedTool = msg.Tool
			m.context.Installer.SetTool(msg.Tool)
		}
		if msg.Path != "" {
			m.context.SelectedPath = msg.Path
			m.context.Installer.SetInstallPath(msg.Path)
		}
		if msg.Files != nil {
			m.context.SelectedFiles = msg.Files
			m.context.Installer.SetSelectedFiles(msg.Files)
		}
		return m, nil

	default:
		// Pass all other messages to current view
		var cmd tea.Cmd
		m.currentView, cmd = m.currentView.Update(msg)
		return m, cmd
	}
}

// transitionToView switches to a new view based on state
func (m *MainModel) transitionToView(newState InstallerState, data interface{}) (tea.Model, tea.Cmd) {
	if !IsValidTransition(m.state, newState) {
		// Transition to error view with invalid transition error
		errorModel := NewErrorModel(m.context)
		errorModel.SetError(fmt.Errorf("invalid state transition from %s to %s", m.state, newState), "changing state")
		m.currentView = errorModel
		m.state = StateError
		return m, errorModel.Init()
	}

	m.state = newState
	var newView ViewModel
	
	switch newState {
	case StateWelcome:
		newView = NewWelcomeModel(m.context)
	case StateToolSelection:
		newView = NewToolSelectionModel(m.context)
	case StatePathSelection:
		newView = NewPathSelectionModel(m.context)
	case StateFileSelection:
		newView = NewFileSelectionModel(m.context)
	case StateHuhConfirmation:
		newView = NewHuhConfirmationModel(m.context)
	case StateConfirmation:
		newView = NewConfirmationModel(m.context)
	case StateInstalling:
		newView = NewInstallingModel(m.context)
	case StateComplete:
		newView = NewCompleteModel(m.context)
	case StateError:
		errorModel := NewErrorModel(m.context)
		// Check if error data was provided
		if data != nil {
			if errorData, ok := data.(map[string]interface{}); ok {
				if err, hasErr := errorData["error"].(error); hasErr {
					context, _ := errorData["context"].(string)
					errorModel.SetError(err, context)
				}
			}
		}
		newView = errorModel
	default:
		// Fallback to error view
		errorModel := NewErrorModel(m.context)
		errorModel.SetError(fmt.Errorf("unknown state: %s", newState), "transitioning view")
		newView = errorModel
		m.state = StateError
	}

	m.currentView = newView
	return m, newView.Init()
}

// View renders the current view
func (m *MainModel) View() string {
	return m.currentView.View()
}

// RunMainInstaller starts the installation UI using the MainModel
func RunMainInstaller(agents, commands, hooks, templates *embed.FS) error {
	model := NewMainModel(agents, commands, hooks, templates)
	
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	
	_, err := program.Run()
	return err
}