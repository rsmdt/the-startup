package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/the-startup/the-startup/internal/installer"
	"github.com/the-startup/the-startup/internal/ui"
)

// Style definitions
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF06B7")).
			MarginTop(1).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C7EFF"))
)

// buildFileTree creates a tree structure from embedded files
func buildFileTree(agents, commands, hooks, templates *embed.FS) *ui.TreeNode {
	root := &ui.TreeNode{
		Name:     "Components",
		IsDir:    true,
		Selected: false,
		Children: []*ui.TreeNode{},
	}
	
	// Add agents
	agentsNode := &ui.TreeNode{
		Name:     "agents",
		Path:     "agents",
		IsDir:    true,
		Selected: false,
		Children: []*ui.TreeNode{},
		Parent:   root,
	}
	if files, err := fs.Glob(agents, "assets/agents/*.md"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			agentsNode.Children = append(agentsNode.Children, &ui.TreeNode{
				Name:     fileName,
				Path:     "agents/" + fileName,
				IsDir:    false,
				Selected: true, // Default to selected
				Parent:   agentsNode,
			})
		}
	}
	root.Children = append(root.Children, agentsNode)
	
	// Add commands
	commandsNode := &ui.TreeNode{
		Name:     "commands",
		Path:     "commands",
		IsDir:    true,
		Selected: false,
		Children: []*ui.TreeNode{},
		Parent:   root,
	}
	if files, err := fs.Glob(commands, "assets/commands/*.md"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			commandsNode.Children = append(commandsNode.Children, &ui.TreeNode{
				Name:     fileName,
				Path:     "commands/" + fileName,
				IsDir:    false,
				Selected: true, // Default to selected
				Parent:   commandsNode,
			})
		}
	}
	root.Children = append(root.Children, commandsNode)
	
	// Add hooks
	hooksNode := &ui.TreeNode{
		Name:     "hooks",
		Path:     "hooks",
		IsDir:    true,
		Selected: false,
		Children: []*ui.TreeNode{},
		Parent:   root,
	}
	if files, err := fs.Glob(hooks, "assets/hooks/*.py"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			hooksNode.Children = append(hooksNode.Children, &ui.TreeNode{
				Name:     fileName,
				Path:     "hooks/" + fileName,
				IsDir:    false,
				Selected: true, // Default to selected
				Parent:   hooksNode,
			})
		}
	}
	root.Children = append(root.Children, hooksNode)
	
	// Add templates
	templatesNode := &ui.TreeNode{
		Name:     "templates",
		Path:     "templates",
		IsDir:    true,
		Selected: false,
		Children: []*ui.TreeNode{},
		Parent:   root,
	}
	if files, err := fs.Glob(templates, "assets/templates/*"); err == nil {
		for _, file := range files {
			fileName := filepath.Base(file)
			templatesNode.Children = append(templatesNode.Children, &ui.TreeNode{
				Name:     fileName,
				Path:     "templates/" + fileName,
				IsDir:    false,
				Selected: true, // Default to selected
				Parent:   templatesNode,
			})
		}
	}
	root.Children = append(root.Children, templatesNode)
	
	return root
}

// confirmQuit asks the user to confirm they want to quit
func confirmQuit() bool {
	var shouldQuit bool
	quitForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you sure you want to quit the installer?").
				Value(&shouldQuit).
				Affirmative("Yes, quit").
				Negative("No, continue"),
		),
	)
	
	if err := quitForm.Run(); err != nil {
		// If the quit confirmation itself is cancelled, stay in installer
		return false
	}
	
	return shouldQuit
}

// NewInstallCommand creates the install command
func NewInstallCommand(agents, commands, hooks, templates *embed.FS) *cobra.Command {
	var (
		installPath string
		toolType    string
		nonInteractive bool
		treeMode    bool
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install The Startup agent system",
		Long:  `Install agents, hooks, and commands for development tools`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Print welcome message
			fmt.Println(titleStyle.Render("🚀 The Startup Installer"))
			fmt.Println(infoStyle.Render("Agent system for development tools"))
			fmt.Println()

			// Create installer
			inst := installer.New(agents, commands, hooks, templates)

			// Interactive mode
			if !nonInteractive {
				// Select tool type
			toolSelectionLoop:
				for {
					var selectedTool string
					toolForm := huh.NewForm(
						huh.NewGroup(
							huh.NewSelect[string]().
								Title("Select your development tool").
								Options(
									huh.NewOption("Claude Code", "claude-code"),
									huh.NewOption("Cancel", "cancel"),
								).
								Value(&selectedTool),
						),
					)

					if err := toolForm.Run(); err != nil {
						// ESC was pressed on first screen - show quit confirmation
						if confirmQuit() {
							fmt.Println(infoStyle.Render("Installation cancelled"))
							return nil
						}
						// User chose to continue, show the form again
						continue
					}

					if selectedTool == "cancel" {
						fmt.Println(infoStyle.Render("Installation cancelled"))
						return nil
					}
					toolType = selectedTool
					break toolSelectionLoop
				}

				// Select installation path
			pathSelectionLoop:
				for {
					var pathChoice string
					cwd, _ := os.Getwd()
					localPath := fmt.Sprintf(".the-startup (local to %s)", cwd)
					
					pathForm := huh.NewForm(
						huh.NewGroup(
							huh.NewSelect[string]().
								Title("Select installation location").
								Options(
									huh.NewOption("~/.config/the-startup (recommended)", "default"),
									huh.NewOption(localPath, "local"),
									huh.NewOption("Custom location", "custom"),
									huh.NewOption("Cancel", "cancel"),
								).
								Value(&pathChoice),
						),
					)

					if err := pathForm.Run(); err != nil {
						// ESC was pressed - go back to tool selection
						goto toolSelectionLoop
					}

					switch pathChoice {
					case "cancel":
						fmt.Println(infoStyle.Render("Installation cancelled"))
						return nil
					case "default":
						installPath = ""
						break pathSelectionLoop
					case "local":
						installPath = filepath.Join(cwd, ".the-startup")
						break pathSelectionLoop
					case "custom":
						var customPath string
						customForm := huh.NewForm(
							huh.NewGroup(
								huh.NewInput().
									Title("Enter custom installation path").
									Placeholder("~/.config/the-startup").
									Value(&customPath),
							),
						)
						if err := customForm.Run(); err != nil {
							// ESC in custom path input - go back to path selection
							continue pathSelectionLoop
						}
						installPath = customPath
						break pathSelectionLoop
					}
				}

				// Select components
			componentSelectionLoop:
				for {
					var selectedComponents []string
					var selectedFiles []string
					
					// Ask if user wants detailed file selection
					var useTreeMode bool
					treeModeForm := huh.NewForm(
						huh.NewGroup(
							huh.NewConfirm().
								Title("Would you like to select individual files? (Advanced)").
								Description("Choose specific files to install, or use simple component selection").
								Value(&useTreeMode).
								Affirmative("Yes, show all files").
								Negative("No, use simple mode"),
						),
					)
					
					if err := treeModeForm.Run(); err != nil {
						// ESC was pressed - go back to path selection
						goto pathSelectionLoop
					}
					
					if useTreeMode || treeMode {
						// Use tree selector for individual file selection
						root := buildFileTree(agents, commands, hooks, templates)
						selectedPaths, err := ui.RunTreeSelector(
							"Select files to install (use 'a' to toggle all, space to toggle selection)",
							root,
						)
						if err != nil {
							// ESC was pressed in tree selector - go back to mode selection
							continue componentSelectionLoop
						}
					
						// Convert selected files to components and individual files
						componentsMap := make(map[string]bool)
						for _, path := range selectedPaths {
							parts := strings.Split(path, "/")
							if len(parts) >= 1 {
								componentsMap[parts[0]] = true
								selectedFiles = append(selectedFiles, path)
							}
						}
						
						// Extract unique components
						for component := range componentsMap {
							selectedComponents = append(selectedComponents, component)
						}
						
						// Set both components and individual files
						inst.SetComponents(selectedComponents)
						inst.SetSelectedFiles(selectedFiles)
						break componentSelectionLoop
					} else {
						// Use simple multi-select for components
						componentForm := huh.NewForm(
							huh.NewGroup(
								huh.NewMultiSelect[string]().
									Title("Select components to install").
									Options(
										huh.NewOption("Agents (12 specialized agents)", "agents").Selected(true),
										huh.NewOption("Hooks (logging and tracking)", "hooks").Selected(true),
										huh.NewOption("Commands (develop, start)", "commands").Selected(true),
										huh.NewOption("Templates (document templates)", "templates").Selected(true),
									).
									Value(&selectedComponents),
							),
						)

						if err := componentForm.Run(); err != nil {
							// ESC was pressed - go back to mode selection
							continue componentSelectionLoop
						}
						
						// Set installer options
						inst.SetComponents(selectedComponents)
						break componentSelectionLoop
					}
				}
			}

			// Set configuration
			inst.SetTool(toolType)
			if installPath != "" {
				inst.SetInstallPath(installPath)
			}

			// Check for existing installation
			if inst.IsInstalled() {
				fmt.Println(infoStyle.Render("Found existing installation"))
				
				if !nonInteractive {
					var shouldUpdate bool
					updateForm := huh.NewForm(
						huh.NewGroup(
							huh.NewConfirm().
								Title("Update existing installation?").
								Value(&shouldUpdate),
						),
					)

					if err := updateForm.Run(); err != nil {
						return fmt.Errorf("update confirmation cancelled: %w", err)
					}

					if !shouldUpdate {
						fmt.Println(infoStyle.Render("Installation cancelled"))
						return nil
					}
				}
			}

			// Run installation
			fmt.Println()
			fmt.Println(infoStyle.Render("Installing components..."))
			
			if err := inst.Install(); err != nil {
				fmt.Println(errorStyle.Render(fmt.Sprintf("✗ Installation failed: %v", err)))
				return err
			}

			// Success message
			fmt.Println()
			fmt.Println(successStyle.Render("✓ Installation complete!"))
			fmt.Println()
			
			installLocation := inst.GetInstallPath()
			fmt.Println("The Startup has been installed to:", installLocation)
			
			// Determine if this is a local or global install
			isLocal := strings.Contains(installLocation, ".the-startup") && !strings.Contains(installLocation, ".config")
			if isLocal {
				fmt.Println(infoStyle.Render("Type: Local installation (project-specific)"))
			} else {
				fmt.Println(infoStyle.Render("Type: Global installation (all projects)"))
			}
			
			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Println("  1. Restart your", toolType, "session")
			if isLocal {
				fmt.Println("  2. Use /develop command in this project")
			} else {
				fmt.Println("  2. Use /develop command in any project")
			}
			fmt.Println("  3. Check logs in ~/.the-startup/")
			fmt.Println()
			fmt.Println(infoStyle.Render("Run 'the-startup help' for more commands"))

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&installPath, "path", "p", "", "Installation path (default: ~/.config/the-startup)")
	cmd.Flags().StringVarP(&toolType, "tool", "t", "claude-code", "Tool type (claude-code)")
	cmd.Flags().BoolVarP(&nonInteractive, "yes", "y", false, "Non-interactive mode, accept defaults")
	cmd.Flags().BoolVar(&treeMode, "tree", false, "Use tree mode for individual file selection")

	return cmd
}
