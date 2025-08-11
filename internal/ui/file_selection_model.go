package ui

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

// FileSelectionModel handles file selection display
type FileSelectionModel struct {
	context *Context
	cursor  int
	choices []string
}

// NewFileSelectionModel creates a new file selection model
func NewFileSelectionModel(context *Context) *FileSelectionModel {
	return &FileSelectionModel{
		context: context,
		cursor:  0,
		choices: []string{"Continue"},
	}
}

// Init initializes the file selection model
func (m *FileSelectionModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for file selection
func (m *FileSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		
		case "esc":
			// Go back to path selection
			return m, func() tea.Msg {
				return ViewTransitionMsg{NextView: StatePathSelection}
			}
		
		case "enter":
			// Proceed to huh confirmation
			return m, func() tea.Msg {
				return ViewTransitionMsg{NextView: StateHuhConfirmation}
			}
		}
	}
	
	return m, nil
}

// buildStaticTree creates a lipgloss tree display for file preview
func (m *FileSelectionModel) buildStaticTree() string {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	
	// Build agents subtree
	agentsTree := tree.New()
	if files, err := fs.Glob(m.context.AgentFiles, "assets/agents/*.md"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			agentsTree = agentsTree.Child(fileName)
		}
	}
	
	// Build commands subtree
	commandsTree := tree.New()
	if files, err := fs.Glob(m.context.CommandFiles, "assets/commands/*.md"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			commandsTree = commandsTree.Child(fileName)
		}
	}
	
	// Build hooks subtree
	hooksTree := tree.New()
	if files, err := fs.Glob(m.context.HookFiles, "assets/hooks/*.py"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			hooksTree = hooksTree.Child(fileName)
		}
	}
	
	// Build templates subtree
	templatesTree := tree.New()
	if files, err := fs.Glob(m.context.TemplateFiles, "assets/templates/*"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			templatesTree = templatesTree.Child(fileName)
		}
	}
	
	// Get the actual Claude path that will be used
	claudePath := m.context.Installer.GetClaudePath()
	
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

// View renders the file selection screen
func (m *FileSelectionModel) View() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.context.Styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Progressive disclosure header
	displayPath := m.context.SelectedPath
	if displayPath == "" {
		displayPath = "~/.config/the-startup"
	}
	s.WriteString(m.context.Renderer.RenderSelections(m.context.SelectedTool, displayPath, len(m.context.SelectedFiles)))
	
	// Title
	s.WriteString(m.context.Renderer.RenderTitle("Files to be moved to your .claude directory"))
	
	// Informative message
	s.WriteString(m.context.Styles.Info.Render("The following files will be moved to your selected .claude directory:"))
	s.WriteString("\n\n")
	
	// Show static tree of files that will be installed
	s.WriteString(m.buildStaticTree())
	s.WriteString("\n")
	
	// Choices
	for i, choice := range m.choices {
		s.WriteString(m.context.Renderer.RenderChoiceWithMultiSelect(choice, i == m.cursor, false, false))
		s.WriteString("\n")
	}
	
	// Help
	s.WriteString(m.context.Renderer.RenderHelp("Enter: continue to installation • Escape: back"))
	
	return s.String()
}