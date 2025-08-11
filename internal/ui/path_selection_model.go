package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// PathSelectionModel handles path selection
type PathSelectionModel struct {
	context *Context
	cursor  int
	choices []string
}

// NewPathSelectionModel creates a new path selection model
func NewPathSelectionModel(context *Context) *PathSelectionModel {
	cwd, _ := os.Getwd()
	localPath := fmt.Sprintf(".the-startup (local to %s)", cwd)
	
	return &PathSelectionModel{
		context: context,
		cursor:  0,
		choices: []string{
			"~/.config/the-startup (recommended)",
			localPath,
			"Custom location",
			"Cancel",
		},
	}
}

// Init initializes the path selection model
func (m *PathSelectionModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for path selection
func (m *PathSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		
		case "esc":
			// Go back to tool selection
			return m, func() tea.Msg {
				return ViewTransitionMsg{NextView: StateToolSelection}
			}
		
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		
		case "enter":
			choice := m.choices[m.cursor]
			
			switch choice {
			case "Cancel":
				return m, tea.Quit
			
			case "~/.config/the-startup (recommended)":
				// Set empty path for recommended location, get all files
				allFiles := m.getAllAvailableFiles()
				return m, tea.Batch(
					func() tea.Msg {
						return SelectionMadeMsg{Path: "", Files: allFiles}
					},
					func() tea.Msg {
						return ViewTransitionMsg{NextView: StateFileSelection}
					},
				)
			
			case "Custom location":
				// TODO: Implement custom path input
				return m, func() tea.Msg {
					return ViewTransitionMsg{NextView: StateError, Data: map[string]interface{}{
						"error": fmt.Errorf("custom path input not yet implemented"),
						"context": "selecting custom path",
					}}
				}
			
			default:
				// Local path option
				cwd, _ := os.Getwd()
				selectedPath := filepath.Join(cwd, ".the-startup")
				allFiles := m.getAllAvailableFiles()
				
				return m, tea.Batch(
					func() tea.Msg {
						return SelectionMadeMsg{Path: selectedPath, Files: allFiles}
					},
					func() tea.Msg {
						return ViewTransitionMsg{NextView: StateFileSelection}
					},
				)
			}
		}
	}
	
	return m, nil
}

// getAllAvailableFiles returns all files that would be installed
func (m *PathSelectionModel) getAllAvailableFiles() []string {
	allFiles := make([]string, 0)
	
	// Helper to add files from embed.FS
	addFiles := func(embedFS *embed.FS, pattern, prefix string) {
		if files, err := fs.Glob(embedFS, pattern); err == nil {
			for _, file := range files {
				fileName := filepath.Base(file)
				filePath := prefix + fileName
				allFiles = append(allFiles, filePath)
			}
		}
	}
	
	// Add agent files
	addFiles(m.context.AgentFiles, "assets/agents/*.md", "agents/")
	
	// Add command files
	addFiles(m.context.CommandFiles, "assets/commands/*.md", "commands/")
	
	// Add hook files
	addFiles(m.context.HookFiles, "assets/hooks/*.py", "hooks/")
	
	// Add template files
	addFiles(m.context.TemplateFiles, "assets/templates/*", "templates/")
	
	return allFiles
}

// View renders the path selection screen
func (m *PathSelectionModel) View() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.context.Styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Progressive disclosure header
	s.WriteString(m.context.Renderer.RenderSelections(m.context.SelectedTool, "", 0))
	
	// Title
	s.WriteString(m.context.Renderer.RenderTitle("Select installation location"))
	
	// Choices
	for i, choice := range m.choices {
		s.WriteString(m.context.Renderer.RenderChoiceWithMultiSelect(choice, i == m.cursor, false, false))
		s.WriteString("\n")
	}
	
	// Help
	s.WriteString(m.context.Renderer.RenderHelp("↑↓ navigate • Enter: select • Escape: back"))
	
	return s.String()
}