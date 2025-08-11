package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/the-startup/the-startup/internal/installer"
)

// InstallProgressMsg represents installation progress
type InstallProgressMsg struct {
	Stage   string
	Message string
	Percent float64
}

// InstallCompleteMsg indicates installation completion
type InstallCompleteMsg struct {
	Success bool
	Error   error
}

// HuhConfirmCompleteMsg indicates huh confirmation completion
type HuhConfirmCompleteMsg struct {
	Confirmed bool
}

// InstallerModel is the main bubbletea model for the installation UI
type InstallerModel struct {
	// State machine
	state InstallerState
	
	// Dependencies
	installer   *installer.Installer
	agentFiles  *embed.FS
	commandFiles *embed.FS
	hookFiles   *embed.FS
	templateFiles *embed.FS
	
	// User selections (progressive disclosure)
	selectedTool       string
	selectedPath       string
	selectedMode       string // "Simple (components)" or "Advanced (individual files)"
	selectedComponents []string
	selectedFiles      []string
	
	// UI state
	cursor     int
	choices    []string
	width      int
	height     int
	
	// Static tree display (no longer using interactive TreeSelector)
	
	// Huh confirmation form
	huhForm        *huh.Form
	huhConfirmed   bool
	
	// Error handling
	err        error
	errContext string
	
	// Installation state
	installProgress InstallProgressMsg
	installComplete bool
	
	// Styles and rendering
	styles   Styles
	renderer *ProgressiveDisclosureRenderer
	
	// Removed cumulative view - now rendering one screen at a time
}

// NewInstallerModel creates a new installer model
func NewInstallerModel(agents, commands, hooks, templates *embed.FS) *InstallerModel {
	return &InstallerModel{
		state:         StateWelcome,
		installer:     installer.New(agents, commands, hooks, templates),
		agentFiles:    agents,
		commandFiles:  commands,
		hookFiles:     hooks,
		templateFiles: templates,
		width:         80,
		height:        24,
		styles:        GetStyles(),
		renderer:      NewProgressiveDisclosureRenderer(),
	}
}

// Init initializes the model
func (m *InstallerModel) Init() tea.Cmd {
	// Auto-advance from welcome screen after 2 seconds
	if m.state == StateWelcome {
		return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
			return TickMsg{}
		})
	}
	return nil
}

// Update handles messages and updates the model
func (m *InstallerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	
	case TickMsg:
		// Handle auto-advance from welcome screen
		if m.state == StateWelcome {
			m.transitionToState(StateToolSelection)
		}
		return m, nil
	
	case InstallProgressMsg:
		m.installProgress = msg
		return m, nil
	
	case InstallCompleteMsg:
		m.installComplete = true
		if msg.Success {
			m.state = StateComplete
		} else {
			m.state = StateError
			m.err = msg.Error
			m.errContext = "during installation"
		}
		return m, nil
	
	case HuhConfirmCompleteMsg:
		if msg.Confirmed {
			m.transitionToState(StateConfirmation)
		} else {
			m.transitionToState(StateFileSelection)
		}
		return m, nil
	
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}
	
	// Handle huh form updates when in StateHuhConfirmation
	if m.state == StateHuhConfirmation && m.huhForm != nil {
		form, cmd := m.huhForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.huhForm = f
			
			// Check if form is complete
			if m.huhForm.State == huh.StateCompleted {
				// Process the confirmation result
				return m, func() tea.Msg {
					return HuhConfirmCompleteMsg{Confirmed: m.huhConfirmed}
				}
			}
		}
		return m, cmd
	}
	
	return m, nil
}

// handleKeyMsg processes keyboard input based on current state
func (m *InstallerModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if m.state == StateWelcome {
			return m, tea.Quit
		}
		// In other states, treat as ESC (go back)
		return m.handleBack()
	
	case "esc":
		return m.handleBack()
	
	case "enter":
		return m.handleEnter()
	
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	}
	
	return m, nil
}

// handleBack processes back navigation (ESC key)
func (m *InstallerModel) handleBack() (tea.Model, tea.Cmd) {
	switch m.state {
	case StateWelcome:
		return m, tea.Quit
	case StateToolSelection:
		m.transitionToState(StateWelcome)
	case StatePathSelection:
		m.transitionToState(StateToolSelection)
	case StateFileSelection:
		m.transitionToState(StatePathSelection)
	case StateHuhConfirmation:
		m.transitionToState(StateFileSelection)
	case StateConfirmation:
		m.transitionToState(StateHuhConfirmation)
	case StateError:
		m.transitionToState(StateWelcome)
	default:
		// States that don't support back navigation
		return m, nil
	}
	return m, nil
}

// handleEnter processes selection confirmation (Enter key)
func (m *InstallerModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case StateWelcome:
		m.transitionToState(StateToolSelection)
	
	case StateToolSelection:
		if m.cursor < len(m.choices) {
			choice := m.choices[m.cursor]
			if choice == "Cancel" {
				return m, tea.Quit
			}
			m.selectedTool = choice
			m.installer.SetTool(choice)
			m.transitionToState(StatePathSelection)
		}
	
	case StatePathSelection:
		if m.cursor < len(m.choices) {
			choice := m.choices[m.cursor]
			switch choice {
			case "Cancel":
				return m, tea.Quit
			case "~/.config/the-startup (recommended)":
				m.selectedPath = ""
				m.selectedFiles = m.getAllAvailableFiles()
				m.installer.SetSelectedFiles(m.selectedFiles)
				m.transitionToState(StateFileSelection)
			case "Custom location":
				// TODO: Implement custom path input
				m.err = fmt.Errorf("custom path input not yet implemented")
				m.errContext = "selecting custom path"
				m.transitionToState(StateError)
			default:
				// Local path option
				cwd, _ := os.Getwd()
				m.selectedPath = filepath.Join(cwd, ".the-startup")
				m.installer.SetInstallPath(m.selectedPath)
				m.selectedFiles = m.getAllAvailableFiles()
				m.installer.SetSelectedFiles(m.selectedFiles)
				m.transitionToState(StateFileSelection)
			}
		}
	
	case StateFileSelection:
		// Proceed with all selected files
		m.installer.SetSelectedFiles(m.selectedFiles)
		m.transitionToState(StateHuhConfirmation)
	
	case StateHuhConfirmation:
		// The huh form handles its own Enter key
		// Transition happens via HuhConfirmCompleteMsg
		return m, nil
	
	case StateConfirmation:
		if m.cursor == 0 { // "Install" option
			m.transitionToState(StateInstalling)
			return m, m.performInstallation()
		} else { // "Go back" option
			return m.handleBack()
		}
	
	default:
		// States that don't handle enter
		return m, nil
	}
	
	return m, nil
}


// transitionToState handles state transitions and updates UI choices
func (m *InstallerModel) transitionToState(newState InstallerState) {
	if !IsValidTransition(m.state, newState) {
		m.err = fmt.Errorf("invalid state transition from %s to %s", m.state, newState)
		m.errContext = "changing state"
		m.state = StateError
		return
	}
	
	m.state = newState
	m.cursor = 0
	m.err = nil
	m.errContext = ""
	
	// Update choices based on new state
	switch newState {
	case StateWelcome:
		m.choices = []string{}
	
	case StateToolSelection:
		m.choices = []string{"claude-code", "Cancel"}
	
	case StatePathSelection:
		cwd, _ := os.Getwd()
		localPath := fmt.Sprintf(".the-startup (local to %s)", cwd)
		m.choices = []string{
			"~/.config/the-startup (recommended)",
			localPath,
			"Custom location",
			"Cancel",
		}
	
	
	case StateFileSelection:
		m.choices = []string{"Continue"}
		// Auto-select all files if not already selected
		if len(m.selectedFiles) == 0 {
			m.selectedFiles = m.getAllAvailableFiles()
			m.installer.SetSelectedFiles(m.selectedFiles)
		}
	
	case StateHuhConfirmation:
		// Create huh confirmation form
		m.huhForm = huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Ready to install?").
					Description("This will install The (Agentic) Startup to your .claude directory.").
					Affirmative("Yes, give me awesome").
					Negative("Huh? I did not sign up for this").
					Value(&m.huhConfirmed),
			),
		)
		m.choices = []string{}
	
	case StateConfirmation:
		m.choices = []string{"Install", "Go back"}
	
	default:
		m.choices = []string{}
	}
}

// getAllAvailableFiles returns all files that would be installed in advanced mode
func (m *InstallerModel) getAllAvailableFiles() []string {
	allFiles := make([]string, 0)
	
	// Add agent files
	if files, err := fs.Glob(m.agentFiles, "assets/agents/*.md"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			filePath := "agents/" + fileName
			allFiles = append(allFiles, filePath)
		}
	}
	
	// Add command files
	if files, err := fs.Glob(m.commandFiles, "assets/commands/*.md"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			filePath := "commands/" + fileName
			allFiles = append(allFiles, filePath)
		}
	}
	
	// Add hook files
	if files, err := fs.Glob(m.hookFiles, "assets/hooks/*.py"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			filePath := "hooks/" + fileName
			allFiles = append(allFiles, filePath)
		}
	}
	
	// Add template files
	if files, err := fs.Glob(m.templateFiles, "assets/templates/*"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			filePath := "templates/" + fileName
			allFiles = append(allFiles, filePath)
		}
	}
	
	return allFiles
}

// buildStaticTree creates a lipgloss tree display for advanced mode preview
func (m *InstallerModel) buildStaticTree() string {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	
	// Build agents subtree
	agentsTree := tree.New()
	if files, err := fs.Glob(m.agentFiles, "assets/agents/*.md"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			agentsTree = agentsTree.Child(fileName)
		}
	}
	
	// Build commands subtree
	commandsTree := tree.New()
	if files, err := fs.Glob(m.commandFiles, "assets/commands/*.md"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			commandsTree = commandsTree.Child(fileName)
		}
	}
	
	// Build hooks subtree
	hooksTree := tree.New()
	if files, err := fs.Glob(m.hookFiles, "assets/hooks/*.py"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			hooksTree = hooksTree.Child(fileName)
		}
	}
	
	// Build templates subtree
	templatesTree := tree.New()
	if files, err := fs.Glob(m.templateFiles, "assets/templates/*"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			templatesTree = templatesTree.Child(fileName)
		}
	}
	
	// Get the actual Claude path that will be used
	claudePath := m.installer.GetClaudePath()
	
	// Abbreviate the path for display
	displayPath := claudePath
	if strings.HasPrefix(claudePath, os.Getenv("HOME")) {
		displayPath = strings.Replace(claudePath, os.Getenv("HOME"), "~", 1)
	}
	
	// Create main tree with actual Claude path
	t := tree.
		Root("⁜ " + displayPath).
		Child(
			"agents",
			agentsTree,
			"commands", 
			commandsTree,
			"hooks",
			hooksTree,
			"templates",
			templatesTree,
		).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle).
		ItemStyle(itemStyle)
	
	return t.String()
}


// performInstallation starts the installation process
func (m *InstallerModel) performInstallation() tea.Cmd {
	return func() tea.Msg {
		// Simulate installation progress (in a real implementation, this would be async)
		if err := m.installer.Install(); err != nil {
			return InstallCompleteMsg{Success: false, Error: err}
		}
		return InstallCompleteMsg{Success: true}
	}
}

// View renders the current state as a single complete screen
func (m *InstallerModel) View() string {
	switch m.state {
	case StateWelcome:
		return m.renderWelcome()
	case StateToolSelection:
		return m.renderToolSelection()
	case StatePathSelection:
		return m.renderPathSelection()
	case StateFileSelection:
		return m.renderFileSelection()
	case StateHuhConfirmation:
		return m.renderHuhConfirmation()
	case StateConfirmation:
		return m.renderConfirmation()
	case StateInstalling:
		return m.renderInstalling()
	case StateComplete:
		return m.renderComplete()
	case StateError:
		return m.renderError()
	default:
		return "Unknown state"
	}
}

// renderWelcome renders the welcome screen
func (m *InstallerModel) renderWelcome() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.styles.Title.Render(WelcomeBanner))
	s.WriteString("\n")
	
	// Welcome message
	s.WriteString(m.styles.Normal.Render(WelcomeMessage))
	s.WriteString("\n\n")
	
	// Instructions
	s.WriteString(m.styles.Help.Render("Press Enter to begin • Escape or Ctrl+C to exit"))
	
	return s.String()
}

// renderToolSelection renders the tool selection screen
func (m *InstallerModel) renderToolSelection() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Progressive disclosure header (empty for first step)
	s.WriteString(m.renderer.RenderSelections("", "", 0))
	
	// Title
	s.WriteString(m.renderer.RenderTitle("Select your development tool"))
	
	// Choices
	for i, choice := range m.choices {
		if choice == "Cancel" {
			s.WriteString(m.renderer.RenderChoiceWithMultiSelect(choice, i == m.cursor, false, false))
		} else {
			s.WriteString(m.renderer.RenderChoiceWithMultiSelect("Claude Code", i == m.cursor, false, false))
		}
		s.WriteString("\n")
	}
	
	// Help
	s.WriteString(m.renderer.RenderHelp("↑↓ navigate • Enter: select • Escape: back"))
	
	return s.String()
}

// renderPathSelection renders the path selection screen
func (m *InstallerModel) renderPathSelection() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Progressive disclosure header
	s.WriteString(m.renderer.RenderSelections(m.selectedTool, "", 0))
	
	// Title
	s.WriteString(m.renderer.RenderTitle("Select installation location"))
	
	// Choices
	for i, choice := range m.choices {
		s.WriteString(m.renderer.RenderChoiceWithMultiSelect(choice, i == m.cursor, false, false))
		s.WriteString("\n")
	}
	
	// Help
	s.WriteString(m.renderer.RenderHelp("↑↓ navigate • Enter: select • Escape: back"))
	
	return s.String()
}

// renderFileSelection renders the file selection screen showing .claude directory files
func (m *InstallerModel) renderFileSelection() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Progressive disclosure header
	displayPath := m.selectedPath
	if displayPath == "" {
		displayPath = "~/.config/the-startup"
	}
	s.WriteString(m.renderer.RenderSelections(m.selectedTool, displayPath, len(m.selectedFiles)))
	
	// Title
	s.WriteString(m.renderer.RenderTitle("Files to be moved to your .claude directory"))
	
	// Informative message
	s.WriteString(m.styles.Info.Render("The following files will be moved to your selected .claude directory:"))
	s.WriteString("\n\n")
	
	// Show static tree of files that will be installed
	s.WriteString(m.buildStaticTree())
	s.WriteString("\n")
	
	// Choices
	for i, choice := range m.choices {
		s.WriteString(m.renderer.RenderChoiceWithMultiSelect(choice, i == m.cursor, false, false))
		s.WriteString("\n")
	}
	
	// Help
	s.WriteString(m.renderer.RenderHelp("Enter: continue to installation • Escape: back"))
	
	return s.String()
}

// renderHuhConfirmation renders the huh confirmation dialog
func (m *InstallerModel) renderHuhConfirmation() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Progressive disclosure header
	displayPath := m.selectedPath
	if displayPath == "" {
		displayPath = "~/.config/the-startup"
	}
	s.WriteString(m.renderer.RenderSelections(m.selectedTool, displayPath, len(m.selectedFiles)))
	
	// Render the huh form
	if m.huhForm != nil {
		s.WriteString(m.huhForm.View())
	}
	
	return s.String()
}

// renderConfirmation renders the installation confirmation screen
func (m *InstallerModel) renderConfirmation() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Progressive disclosure header
	displayPath := m.selectedPath
	if displayPath == "" {
		displayPath = "~/.config/the-startup"
	}
	s.WriteString(m.renderer.RenderSelections(m.selectedTool, displayPath, len(m.selectedFiles)))
	
	// Installation summary
	updatingFiles := []string{}
	if len(m.selectedFiles) > 0 {
		for _, file := range m.selectedFiles {
			if m.installer.CheckFileExists(file) {
				updatingFiles = append(updatingFiles, file)
			}
		}
	}
	
	s.WriteString(m.renderer.RenderInstallationSummary(
		m.selectedTool, displayPath, m.selectedFiles, updatingFiles))
	
	// Confirmation choices
	for i, choice := range m.choices {
		s.WriteString(m.renderer.RenderChoiceWithMultiSelect(choice, i == m.cursor, false, false))
		s.WriteString("\n")
	}
	
	// Help
	s.WriteString(m.renderer.RenderHelp("↑↓ navigate • Enter: select • Escape: back"))
	
	return s.String()
}

// renderInstalling renders the installation progress screen
func (m *InstallerModel) renderInstalling() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Title
	s.WriteString(m.renderer.RenderTitle("Installing The (Agentic) Startup"))
	
	// Progress information
	if m.installProgress.Stage != "" {
		s.WriteString(m.styles.Info.Render(fmt.Sprintf("Stage: %s", m.installProgress.Stage)))
		s.WriteString("\n")
	}
	
	if m.installProgress.Message != "" {
		s.WriteString(m.styles.Normal.Render(m.installProgress.Message))
		s.WriteString("\n")
	}
	
	// Simple progress indicator (spinner could be added here)
	s.WriteString("\n")
	s.WriteString(m.styles.Info.Render("⚡ Installing components..."))
	s.WriteString("\n\n")
	
	// Progress bar could be added here based on m.installProgress.Percent
	
	return s.String()
}

// renderComplete renders the installation complete screen
func (m *InstallerModel) renderComplete() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Success message
	s.WriteString(m.styles.Success.Render("✅ Installation Complete!"))
	s.WriteString("\n\n")
	
	// Installation details
	installLocation := m.installer.GetInstallPath()
	s.WriteString(m.styles.Normal.Render(fmt.Sprintf("The Startup has been installed to: %s", installLocation)))
	s.WriteString("\n\n")
	
	// Next steps
	s.WriteString(m.styles.Info.Render("Next steps:"))
	s.WriteString("\n")
	s.WriteString(m.styles.Normal.Render("1. Restart your " + m.selectedTool + " session"))
	s.WriteString("\n")
	
	isLocal := strings.Contains(installLocation, ".the-startup") && !strings.Contains(installLocation, ".config")
	if isLocal {
		s.WriteString(m.styles.Normal.Render("2. Use /develop command in this project"))
	} else {
		s.WriteString(m.styles.Normal.Render("2. Use /develop command in any project"))
	}
	s.WriteString("\n")
	s.WriteString(m.styles.Normal.Render("3. Check logs in ~/.the-startup/"))
	s.WriteString("\n\n")
	
	// Exit instruction
	s.WriteString(m.styles.Help.Render("Press any key to exit"))
	
	return s.String()
}

// renderError renders the error screen
func (m *InstallerModel) renderError() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Error message
	s.WriteString(m.renderer.RenderError(m.err, m.errContext))
	s.WriteString("\n")
	
	// Recovery instructions
	s.WriteString(m.styles.Help.Render("Press Escape to return to welcome screen"))
	
	return s.String()
}

// RunInstaller starts the installation UI (legacy - now uses MainModel)
func RunInstaller(agents, commands, hooks, templates *embed.FS) error {
	// Delegate to the new MainModel implementation
	return RunMainInstaller(agents, commands, hooks, templates)
}