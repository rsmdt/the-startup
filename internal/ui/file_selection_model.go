package ui

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/the-startup/the-startup/internal/installer"
)

type FileSelectionModel struct {
	styles        Styles
	renderer      *ProgressiveDisclosureRenderer
	installer     *installer.Installer
	agentFiles    *embed.FS
	commandFiles  *embed.FS
	hookFiles     *embed.FS
	templateFiles *embed.FS
	selectedTool  string
	selectedPath  string
	selectedFiles []string
	cursor        int
	choices       []string
	ready         bool
	confirmed     bool
}

func NewFileSelectionModel(selectedTool, selectedPath string, installer *installer.Installer, agents, commands, hooks, templates *embed.FS) FileSelectionModel {
	m := FileSelectionModel{
		styles:        GetStyles(),
		renderer:      NewProgressiveDisclosureRenderer(),
		installer:     installer,
		agentFiles:    agents,
		commandFiles:  commands,
		hookFiles:     hooks,
		templateFiles: templates,
		selectedTool:  selectedTool,
		selectedPath:  selectedPath,
		choices: []string{
			"Yes, give me awesome",
			"Huh? I did not sign up for this",
		},
		cursor: 0,
		ready:  false,
	}
	
	m.selectedFiles = m.getAllAvailableFiles()
	if m.installer != nil {
		m.installer.SetSelectedFiles(m.selectedFiles)
	}
	
	return m
}

func (m FileSelectionModel) Init() tea.Cmd {
	return nil
}

func (m FileSelectionModel) Update(msg tea.Msg) (FileSelectionModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor < len(m.choices) {
				choice := m.choices[m.cursor]
				if choice == "Yes, give me awesome" {
					m.confirmed = true
					m.ready = true
				} else {
					m.confirmed = false
					m.ready = true
				}
			}
		}
	}
	return m, nil
}

func (m FileSelectionModel) View() string {
	var s strings.Builder
	
	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")
	
	displayPath := m.selectedPath
	if displayPath == "" {
		displayPath = "~/.config/the-startup"
	}
	s.WriteString(m.renderer.RenderSelections(m.selectedTool, displayPath, len(m.selectedFiles)))
	
	s.WriteString(m.renderer.RenderTitle("Files to be moved to your .claude directory"))
	
	s.WriteString(m.styles.Info.Render("The following files will be moved to your selected .claude directory:"))
	s.WriteString("\n\n")
	
	s.WriteString(m.buildStaticTree())
	s.WriteString("\n")
	
	s.WriteString("\n\n")
	s.WriteString(m.styles.Title.Render("Ready to install?"))
	s.WriteString("\n")
	s.WriteString(m.styles.Info.Render("This will install The (Agentic) Startup to your .claude directory."))
	s.WriteString("\n\n")
	
	for i, option := range m.choices {
		if i == m.cursor {
			s.WriteString(m.styles.Selected.Render("> " + option))
		} else {
			s.WriteString(m.styles.Normal.Render("  " + option))
		}
		s.WriteString("\n")
	}
	
	s.WriteString("\n")
	s.WriteString(m.styles.Help.Render("Press Enter to confirm • Escape to go back"))
	
	return s.String()
}

func (m FileSelectionModel) Ready() bool {
	return m.ready
}

func (m FileSelectionModel) Confirmed() bool {
	return m.confirmed
}

func (m FileSelectionModel) Reset() FileSelectionModel {
	m.ready = false
	m.cursor = 0
	return m
}

func (m FileSelectionModel) getAllAvailableFiles() []string {
	allFiles := make([]string, 0)
	
	addFiles := func(embedFS *embed.FS, pattern, prefix string) {
		if files, err := fs.Glob(embedFS, pattern); err == nil {
			for _, file := range files {
				fileName := filepath.Base(file)
				filePath := prefix + fileName
				allFiles = append(allFiles, filePath)
			}
		}
	}
	
	patterns := []string{"assets/agents/*.md", "test_assets/assets/agents/*.md"}
	for _, pattern := range patterns {
		addFiles(m.agentFiles, pattern, "agents/")
	}
	
	patterns = []string{"assets/commands/*.md", "test_assets/assets/commands/*.md"}
	for _, pattern := range patterns {
		addFiles(m.commandFiles, pattern, "commands/")
	}
	
	patterns = []string{"assets/hooks/*.py", "test_assets/assets/hooks/*.py"}
	for _, pattern := range patterns {
		addFiles(m.hookFiles, pattern, "hooks/")
	}
	
	patterns = []string{"assets/templates/*", "test_assets/assets/templates/*"}
	for _, pattern := range patterns {
		addFiles(m.templateFiles, pattern, "templates/")
	}
	
	return allFiles
}

func (m FileSelectionModel) buildStaticTree() string {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	updateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange for updates
	
	// Get list of existing files that will be updated
	existingFiles := make(map[string]bool)
	if m.installer != nil {
		for _, file := range m.installer.GetExistingFiles(m.selectedFiles) {
			existingFiles[file] = true
		}
	}
	
	buildSubtree := func(embedFS *embed.FS, patterns []string, prefix string) []string {
		var items []string
		for _, pattern := range patterns {
			if files, err := fs.Glob(embedFS, pattern); err == nil {
				for _, file := range files {
					fileName := filepath.Base(file)
					filePath := prefix + fileName
					
					// Apply orange color if file will be updated
					if existingFiles[filePath] {
						items = append(items, updateStyle.Render(fileName+" (will update)"))
					} else {
						items = append(items, itemStyle.Render(fileName))
					}
				}
			}
		}
		return items
	}
	
	agentItems := buildSubtree(m.agentFiles, []string{"assets/agents/*.md", "test_assets/assets/agents/*.md"}, "agents/")
	commandItems := buildSubtree(m.commandFiles, []string{"assets/commands/*.md", "test_assets/assets/commands/*.md"}, "commands/")
	hookItems := buildSubtree(m.hookFiles, []string{"assets/hooks/*.py", "test_assets/assets/hooks/*.py"}, "hooks/")
	templateItems := buildSubtree(m.templateFiles, []string{"assets/templates/*", "test_assets/assets/templates/*"}, "templates/")
	
	claudePath := m.installer.GetClaudePath()
	
	displayPath := claudePath
	if strings.HasPrefix(claudePath, os.Getenv("HOME")) {
		displayPath = strings.Replace(claudePath, os.Getenv("HOME"), "~", 1)
	}
	
	// Build the tree with colored items
	agentsTree := tree.New()
	for _, item := range agentItems {
		agentsTree = agentsTree.Child(item)
	}
	
	commandsTree := tree.New()
	for _, item := range commandItems {
		commandsTree = commandsTree.Child(item)
	}
	
	hooksTree := tree.New()
	for _, item := range hookItems {
		hooksTree = hooksTree.Child(item)
	}
	
	templatesTree := tree.New()
	for _, item := range templateItems {
		templatesTree = templatesTree.Child(item)
	}
	
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
		RootStyle(rootStyle)
	
	return t.String()
}