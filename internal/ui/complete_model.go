package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/installer"
)

type CompleteModel struct {
	styles        Styles
	installer     *installer.Installer
	selectedTool  string
	ready         bool
	mode          OperationMode
	claudeAssets  *embed.FS
	startupAssets *embed.FS
	selectedFiles []string
}

// autoExitMsg is sent after a delay to automatically exit
type autoExitMsg struct{}

func NewCompleteModel(selectedTool string, installer *installer.Installer) CompleteModel {
	return NewCompleteModelWithMode(selectedTool, installer, ModeInstall)
}

func NewCompleteModelWithMode(selectedTool string, installer *installer.Installer, mode OperationMode) CompleteModel {
	return CompleteModel{
		styles:       GetStyles(),
		installer:    installer,
		selectedTool: selectedTool,
		ready:        false,
		mode:         mode,
	}
}

func NewCompleteModelWithAssets(selectedTool string, installer *installer.Installer, mode OperationMode, claudeAssets, startupAssets *embed.FS, selectedFiles []string) CompleteModel {
	return CompleteModel{
		styles:        GetStyles(),
		installer:     installer,
		selectedTool:  selectedTool,
		ready:         false,
		mode:          mode,
		claudeAssets:  claudeAssets,
		startupAssets: startupAssets,
		selectedFiles: selectedFiles,
	}
}

func (m CompleteModel) Init() tea.Cmd {
	// Exit immediately
	return func() tea.Msg {
		return autoExitMsg{}
	}
}

func (m CompleteModel) Update(msg tea.Msg) (CompleteModel, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		// Allow immediate exit on any key press
		m.ready = true
	case autoExitMsg:
		// Auto-exit after delay
		m.ready = true
	}
	return m, nil
}

func (m CompleteModel) View() string {
	var s strings.Builder

	// Banner
	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")

	// Success message
	if m.mode == ModeUninstall {
		s.WriteString(m.styles.Success.Render("✅ Uninstallation Complete!"))
	} else {
		s.WriteString(m.styles.Success.Render("✅ Installation Complete!"))
	}
	s.WriteString("\n\n")

	// Use standardized path renderer for consistency
	claudePath := m.installer.GetClaudePath()
	startupPath := m.installer.GetInstallPath()

	renderer := NewProgressiveDisclosureRenderer()
	s.WriteString(renderer.RenderSelectedPaths(claudePath, startupPath, m.mode))

	// Display the tree of installed/removed files - same as during selection
	if m.mode == ModeUninstall {
		s.WriteString(m.styles.Warning.Render("Files removed:"))
	} else {
		s.WriteString(m.styles.Info.Render("Files installed:"))
	}
	s.WriteString("\n\n")
	
	// Build and display the tree (same structure as file selection)
	treeStr := m.buildCompletionTree()
	s.WriteString(treeStr)
	s.WriteString("\n")

	// Display removed files if any
	removedFiles := m.installer.GetDeprecatedFilesList()
	if len(removedFiles) > 0 {
		sort.Strings(removedFiles)
		s.WriteString(m.styles.Error.Render("  Removed deprecated files:"))
		s.WriteString("\n")
		for _, file := range removedFiles {
			s.WriteString(m.styles.Normal.Render("    ✗ " + file))
			s.WriteString("\n")
		}
		s.WriteString("\n")
	}

	// Add repository link
	s.WriteString(m.styles.Info.Render("See https://github.com/rsmdt/the-startup for details"))
	s.WriteString("\n")

	return s.String()
}

func (m CompleteModel) Ready() bool {
	return m.ready
}

// buildCompletionTree builds the same tree structure shown during file selection
func (m CompleteModel) buildCompletionTree() string {
	var s strings.Builder

	// Use themed styles for consistency (matching selection views)
	styles := GetStyles()
	rootStyle := styles.Title
	normalStyle := styles.Normal     // Default text color (light gray)
	updateStyle := styles.Warning    // Orange/Peach for update indicators

	claudePath := m.installer.GetClaudePath()
	displayPath := claudePath
	if strings.HasPrefix(claudePath, os.Getenv("HOME")) {
		displayPath = strings.Replace(claudePath, os.Getenv("HOME"), "~", 1)
	}

	// Check if settings.json was updated
	settingsExists := m.installer.CheckSettingsExists()

	// Root path
	s.WriteString(rootStyle.Render("⁜ " + displayPath))
	s.WriteString("\n")

	// Claude section with checkmark (always selected in completion)
	s.WriteString(normalStyle.Render("✓ .claude/"))
	s.WriteString("\n")

	// Check if we have agents files
	hasAgents := false
	for _, file := range m.selectedFiles {
		if strings.HasPrefix(file, "agents/") {
			hasAgents = true
			break
		}
	}

	// Check if we have commands files
	hasCommands := false
	for _, file := range m.selectedFiles {
		if strings.HasPrefix(file, "commands/") {
			hasCommands = true
			break
		}
	}

	// Check if we have output-styles files
	hasOutputStyles := false
	for _, file := range m.selectedFiles {
		if strings.HasPrefix(file, "output-styles/") {
			hasOutputStyles = true
			break
		}
	}

	// Check if we have settings files
	hasSettings := false
	for _, file := range m.selectedFiles {
		if file == "settings.json" || file == "settings.local.json" {
			hasSettings = true
			break
		}
	}

	// Agents subsection
	if hasAgents {
		s.WriteString(normalStyle.Render("✓ ├── agents/"))
		s.WriteString("\n")

		// Group and display agent files
		agentFiles := make(map[string]bool)

		// Process agent files based on debug output format
		for _, file := range m.selectedFiles {
			if strings.HasPrefix(file, "agents/") {
				relPath := strings.TrimPrefix(file, "agents/")

				if strings.HasSuffix(relPath, ".md") {
					// It's a direct .md file (like "the-chief.md", "the-meta-agent.md")
					agentName := strings.TrimSuffix(relPath, ".md")
					agentFiles[agentName] = true
				} else if !strings.Contains(relPath, "/") {
					// It's a bare directory name (like "the-analyst", "the-designer")
					// Count specialized activities by checking the embedded filesystem
					fileCount := 0
					if m.claudeAssets != nil {
						agentPath := "assets/claude/agents/" + relPath
						fs.WalkDir(m.claudeAssets, agentPath, func(path string, d fs.DirEntry, err error) error {
							if err == nil && !d.IsDir() && strings.HasSuffix(path, ".md") {
								fileCount++
							}
							return nil
						})
					}

					if fileCount > 0 {
						activityLabel := "specialized activity"
						if fileCount > 1 {
							activityLabel = "specialized activities"
						}
						displayName := fmt.Sprintf("%s (%d %s)", relPath, fileCount, activityLabel)
						agentFiles[displayName] = true
					} else {
						// If no specialized activities found, just show the agent name
						agentFiles[relPath] = true
					}
				}
			}
		}

		// Sort and display agent files
		var agentNames []string
		for name := range agentFiles {
			agentNames = append(agentNames, name)
		}
		sort.Strings(agentNames)

		for _, name := range agentNames {
			if m.mode == ModeUninstall {
				s.WriteString(normalStyle.Render("✗ │   ├── " + name))
			} else {
				s.WriteString(normalStyle.Render("✓ │   ├── " + name))
			}
			s.WriteString("\n")
		}
	}

	// Commands subsection
	if hasCommands {
		commandsBranch := "├── "
		if !hasOutputStyles && !hasSettings {
			commandsBranch = "└── "
		}
		s.WriteString(normalStyle.Render("✓ " + commandsBranch + "commands/"))
		s.WriteString("\n")

		// Group and display command files
		commandFiles := make([]string, 0)
		for _, file := range m.selectedFiles {
			if strings.HasPrefix(file, "commands/") {
				// Extract relative path
				relPath := strings.TrimPrefix(file, "commands/")

				// Format display name
				displayName := strings.TrimSuffix(filepath.Base(relPath), ".md")

				// Special formatting for s/ commands
				if strings.Contains(relPath, "/s/") {
					displayName = "/s:" + displayName
				}

				commandFiles = append(commandFiles, displayName)
			}
		}

		// Sort and display command files
		sort.Strings(commandFiles)

		for _, name := range commandFiles {
			treeBranch := "│   ├── "
			if !hasOutputStyles && !hasSettings {
				treeBranch = "    ├── "
			}
			if m.mode == ModeUninstall {
				s.WriteString(normalStyle.Render("✗ " + treeBranch + name))
			} else {
				s.WriteString(normalStyle.Render("✓ " + treeBranch + name))
			}
			s.WriteString("\n")
		}
	}

	// Output-styles subsection
	if hasOutputStyles {
		outputStylesBranch := "├── "
		if !hasSettings {
			outputStylesBranch = "└── "
		}
		s.WriteString(normalStyle.Render("✓ " + outputStylesBranch + "output-styles/"))
		s.WriteString("\n")

		// Group and display output-style files
		outputStyleFiles := make([]string, 0)
		for _, file := range m.selectedFiles {
			if strings.HasPrefix(file, "output-styles/") {
				// Extract relative path
				relPath := strings.TrimPrefix(file, "output-styles/")

				// Format display name
				displayName := strings.TrimSuffix(filepath.Base(relPath), ".json")
				displayName = strings.TrimSuffix(displayName, ".md")

				outputStyleFiles = append(outputStyleFiles, displayName)
			}
		}

		// Sort and display output-style files
		sort.Strings(outputStyleFiles)

		for _, name := range outputStyleFiles {
			treeBranch := "│   ├── "
			if !hasSettings {
				treeBranch = "    ├── "
			}
			if m.mode == ModeUninstall {
				s.WriteString(normalStyle.Render("✗ " + treeBranch + name))
			} else {
				s.WriteString(normalStyle.Render("✓ " + treeBranch + name))
			}
			s.WriteString("\n")
		}
	}

	// Settings files
	if hasSettings {
		for _, file := range m.selectedFiles {
			if file == "settings.json" {
				var settingsLine string
				if m.mode == ModeUninstall {
					settingsLine = "✗ └── settings.json"
				} else if settingsExists {
					settingsLine = "✓ └── " + updateStyle.Render("settings.json (updated)")
				} else {
					settingsLine = "✓ └── settings.json"
				}
				s.WriteString(normalStyle.Render(settingsLine))
				s.WriteString("\n")
				break
			}
		}

		// Handle settings.local.json if it exists
		for _, file := range m.selectedFiles {
			if file == "settings.local.json" {
				localSettingsPath := filepath.Join(claudePath, "settings.local.json")
				var localSettingsLine string
				if m.mode == ModeUninstall {
					localSettingsLine = "✗ └── settings.local.json"
				} else if _, err := os.Stat(localSettingsPath); err == nil {
					localSettingsLine = "✓ └── " + updateStyle.Render("settings.local.json (updated)")
				} else {
					localSettingsLine = "✓ └── settings.local.json"
				}
				s.WriteString(normalStyle.Render(localSettingsLine))
				s.WriteString("\n")
				break
			}
		}
	}

	return s.String()
}
