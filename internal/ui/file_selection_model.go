package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/rsmdt/the-startup/internal/installer"
)

type FileSelectionModel struct {
	styles        Styles
	renderer      *ProgressiveDisclosureRenderer
	installer     *installer.Installer
	claudeAssets  *embed.FS
	startupAssets *embed.FS
	selectedTool  string
	selectedPath  string
	selectedFiles []string
	cursor        int
	choices       []string
	ready         bool
	confirmed     bool
}

func NewFileSelectionModel(selectedTool, selectedPath string, installer *installer.Installer, claudeAssets, startupAssets *embed.FS) FileSelectionModel {
	m := FileSelectionModel{
		styles:        GetStyles(),
		renderer:      NewProgressiveDisclosureRenderer(),
		installer:     installer,
		claudeAssets:  claudeAssets,
		startupAssets: startupAssets,
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

	// Show both paths in the header
	startupPath := m.installer.GetInstallPath()
	claudePath := m.installer.GetClaudePath()

	// Format paths for display
	home := os.Getenv("HOME")
	if home != "" {
		if strings.HasPrefix(startupPath, home) {
			startupPath = "~" + strings.TrimPrefix(startupPath, home)
		}
		if strings.HasPrefix(claudePath, home) {
			claudePath = "~" + strings.TrimPrefix(claudePath, home)
		}
	}

	s.WriteString(m.styles.Info.Render("Installation Paths:"))
	s.WriteString("\n")
	s.WriteString(m.styles.Normal.Render(fmt.Sprintf("  Startup: %s", startupPath)))
	s.WriteString("\n")
	s.WriteString(m.styles.Normal.Render(fmt.Sprintf("  Claude:  %s", claudePath)))
	s.WriteString("\n\n")

	s.WriteString(m.renderer.RenderTitle("Files to be installed to .claude"))

	s.WriteString(m.styles.Info.Render("The following files will be installed to your Claude directory:"))
	s.WriteString("\n\n")

	s.WriteString(m.buildStaticTree())
	s.WriteString("\n")

	s.WriteString("\n\n")
	s.WriteString(m.styles.Title.Render("Ready to install?"))
	s.WriteString("\n")
	s.WriteString(m.styles.Info.Render("This will install The (Agentic) Startup to the selected directories."))
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

	// Walk through all files in Claude assets
	fs.WalkDir(m.claudeAssets, "assets/claude", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		// Create relative path for display/selection
		relPath := strings.TrimPrefix(path, "assets/claude/")
		allFiles = append(allFiles, relPath)
		return nil
	})

	// Walk through all files in Startup assets  
	fs.WalkDir(m.startupAssets, "assets/the-startup", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		// Create relative path for display/selection
		relPath := strings.TrimPrefix(path, "assets/the-startup/")
		allFiles = append(allFiles, relPath)
		return nil
	})

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
					// Extract relative path from assets/claude/[type]/
					relPath := strings.TrimPrefix(file, "assets/claude/")
					relPath = strings.TrimPrefix(relPath, prefix)
					filePath := prefix + relPath

					// Format display name (preserve namespace for commands)
					displayName := relPath

					// Apply orange color if file will be updated
					if existingFiles[filePath] {
						items = append(items, updateStyle.Render(displayName+" (will update)"))
					} else {
						items = append(items, itemStyle.Render(displayName))
					}
				}
			}
		}
		return items
	}

	// Only show files that go to .claude directory (agents and commands)
	agentItems := buildSubtree(m.claudeAssets, []string{"assets/claude/agents/**/*.md", "assets/claude/agents/*.md"}, "agents/")
	commandItems := buildSubtree(m.claudeAssets, []string{"assets/claude/commands/**/*.md"}, "commands/")
	outputStyleItems := buildSubtree(m.claudeAssets, []string{"assets/claude/output-styles/*.md"}, "output-styles/")
	// Don't show hooks and templates as they go to .the-startup, not .claude

	// Check if settings.json exists using installer's method
	settingsExists := m.installer.CheckSettingsExists()

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

	outputStylesTree := tree.New()
	for _, item := range outputStyleItems {
		outputStylesTree = outputStylesTree.Child(item)
	}

	// Add settings.json with appropriate styling
	settingsItem := "settings.json"
	if settingsExists {
		settingsItem = updateStyle.Render("settings.json (will update)")
	} else {
		settingsItem = itemStyle.Render("settings.json")
	}
	
	// Build children list for the tree
	children := []any{
		"agents",
		agentsTree,
		"commands",
		commandsTree,
		"output-styles",
		outputStylesTree,
		settingsItem,
	}
	
	// Add settings.local.json if it exists in assets
	if _, err := m.claudeAssets.ReadFile("assets/claude/settings.local.json"); err == nil {
		localSettingsPath := filepath.Join(claudePath, "settings.local.json")
		localSettingsItem := "settings.local.json"
		if _, err := os.Stat(localSettingsPath); err == nil {
			localSettingsItem = updateStyle.Render("settings.local.json (will update)")
		} else {
			localSettingsItem = itemStyle.Render("settings.local.json")
		}
		children = append(children, localSettingsItem)
	}

	t := tree.
		Root("⁜ "+displayPath).
		Child(children...).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle)

	return t.String()
}
