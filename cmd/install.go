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

// buildFileTree creates a tree structure from embedded files and checks for existing files
func buildFileTree(agents, commands, hooks, templates *embed.FS, inst *installer.Installer) *ui.TreeNode {
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
			filePath := "agents/" + fileName
			agentsNode.Children = append(agentsNode.Children, &ui.TreeNode{
				Name:     fileName,
				Path:     filePath,
				IsDir:    false,
				Selected: true, // Default to selected
				Exists:   inst.CheckFileExists(filePath),
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
			filePath := "commands/" + fileName
			commandsNode.Children = append(commandsNode.Children, &ui.TreeNode{
				Name:     fileName,
				Path:     filePath,
				IsDir:    false,
				Selected: true, // Default to selected
				Exists:   inst.CheckFileExists(filePath),
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
			filePath := "hooks/" + fileName
			hooksNode.Children = append(hooksNode.Children, &ui.TreeNode{
				Name:     fileName,
				Path:     filePath,
				IsDir:    false,
				Selected: true, // Default to selected
				Exists:   inst.CheckFileExists(filePath),
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
			filePath := "templates/" + fileName
			templatesNode.Children = append(templatesNode.Children, &ui.TreeNode{
				Name:     fileName,
				Path:     filePath,
				IsDir:    false,
				Selected: true, // Default to selected
				Exists:   inst.CheckFileExists(filePath),
				Parent:   templatesNode,
			})
		}
	}
	root.Children = append(root.Children, templatesNode)
	
	return root
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
			var updatesAlreadyConfirmed bool
			if !nonInteractive {
				// Select tool type
			toolSelectionLoop:
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
					// ESC was pressed - exit immediately
					fmt.Println(infoStyle.Render("Installation cancelled"))
					return nil
				}

				if selectedTool == "cancel" {
					fmt.Println(infoStyle.Render("Installation cancelled"))
					return nil
				}
				toolType = selectedTool

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

				// Set the installer configuration NOW so file existence checks work
				inst.SetTool(toolType)
				if installPath != "" {
					inst.SetInstallPath(installPath)
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
						root := buildFileTree(agents, commands, hooks, templates, inst)
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
						
						// Check for existing files that will be updated
						var updatingFiles []string
						for _, path := range selectedFiles {
							if inst.CheckFileExists(path) {
								updatingFiles = append(updatingFiles, path)
							}
						}
						
						// If there are files to update, show confirmation
						if len(updatingFiles) > 0 {
							var confirmUpdate bool
							
							// Build tree-structured update message
							updateMsg := fmt.Sprintf("%d file(s) will be updated:\n\n", len(updatingFiles))
							
							// Group files by component for tree display
							filesByComponent := make(map[string][]string)
							for _, file := range updatingFiles {
								parts := strings.Split(file, "/")
								if len(parts) >= 2 {
									component := parts[0]
									fileName := parts[1]
									filesByComponent[component] = append(filesByComponent[component], fileName)
								}
							}
							
							// Display as tree
							for _, component := range []string{"agents", "commands", "hooks", "templates"} {
								if files, ok := filesByComponent[component]; ok && len(files) > 0 {
									updateMsg += fmt.Sprintf("%s/\n", component)
									for _, file := range files {
										updateMsg += fmt.Sprintf("  • %s\n", file)
									}
								}
							}
							updateMsg += "\nDo you want to continue?"
							
							updateConfirm := huh.NewForm(
								huh.NewGroup(
									huh.NewConfirm().
										Title("Files will be updated").
										Description(updateMsg).
										Value(&confirmUpdate).
										Affirmative("Yes, update files").
										Negative("No, go back"),
								),
							)
							
							if err := updateConfirm.Run(); err != nil || !confirmUpdate {
								// User cancelled or said no - go back to file selection
								continue componentSelectionLoop
							}
							// User confirmed updates
							updatesAlreadyConfirmed = true
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
						
						// Check for existing files in selected components
						var updatingFiles []string
						for _, component := range selectedComponents {
							// Check all files in this component
							switch component {
							case "agents":
								if files, err := fs.Glob(agents, "assets/agents/*.md"); err == nil {
									for _, file := range files {
										fileName := filepath.Base(file)
										filePath := "agents/" + fileName
										if inst.CheckFileExists(filePath) {
											updatingFiles = append(updatingFiles, filePath)
										}
									}
								}
							case "commands":
								if files, err := fs.Glob(commands, "assets/commands/*.md"); err == nil {
									for _, file := range files {
										fileName := filepath.Base(file)
										filePath := "commands/" + fileName
										if inst.CheckFileExists(filePath) {
											updatingFiles = append(updatingFiles, filePath)
										}
									}
								}
							case "hooks":
								if files, err := fs.Glob(hooks, "assets/hooks/*.py"); err == nil {
									for _, file := range files {
										fileName := filepath.Base(file)
										filePath := "hooks/" + fileName
										if inst.CheckFileExists(filePath) {
											updatingFiles = append(updatingFiles, filePath)
										}
									}
								}
							case "templates":
								if files, err := fs.Glob(templates, "assets/templates/*"); err == nil {
									for _, file := range files {
										fileName := filepath.Base(file)
										filePath := "templates/" + fileName
										if inst.CheckFileExists(filePath) {
											updatingFiles = append(updatingFiles, filePath)
										}
									}
								}
							}
						}
						
						// If there are files to update, show confirmation
						if len(updatingFiles) > 0 {
							var confirmUpdate bool
							
							// Build tree-structured update message
							updateMsg := fmt.Sprintf("%d file(s) will be updated:\n\n", len(updatingFiles))
							
							// Group files by component for tree display
							filesByComponent := make(map[string][]string)
							for _, file := range updatingFiles {
								parts := strings.Split(file, "/")
								if len(parts) >= 2 {
									component := parts[0]
									fileName := parts[1]
									filesByComponent[component] = append(filesByComponent[component], fileName)
								}
							}
							
							// Display as tree (limit to 10 files per component for readability)
							totalShown := 0
							for _, component := range []string{"agents", "commands", "hooks", "templates"} {
								if files, ok := filesByComponent[component]; ok && len(files) > 0 {
									updateMsg += fmt.Sprintf("%s/\n", component)
									showCount := len(files)
									if showCount > 10 {
										showCount = 10
									}
									for i := 0; i < showCount; i++ {
										updateMsg += fmt.Sprintf("  • %s\n", files[i])
										totalShown++
									}
									if len(files) > 10 {
										updateMsg += fmt.Sprintf("  ... and %d more\n", len(files)-10)
									}
								}
							}
							updateMsg += "\nDo you want to continue?"
							
							updateConfirm := huh.NewForm(
								huh.NewGroup(
									huh.NewConfirm().
										Title("Files will be updated").
										Description(updateMsg).
										Value(&confirmUpdate).
										Affirmative("Yes, update files").
										Negative("No, go back"),
								),
							)
							
							if err := updateConfirm.Run(); err != nil || !confirmUpdate {
								// User cancelled or said no - go back to component selection
								continue componentSelectionLoop
							}
							// User confirmed updates
							updatesAlreadyConfirmed = true
						}
						
						// Set installer options
						inst.SetComponents(selectedComponents)
						break componentSelectionLoop
					}
				}
			}

			// Configuration already set above

			// Check for existing installation (skip if updates were already confirmed)
			if !updatesAlreadyConfirmed && inst.IsInstalled() {
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
						// ESC was pressed - for now just return error
						// (can't easily go back to component selection from here due to scope)
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
