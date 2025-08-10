package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/the-startup/the-startup/internal/installer"
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

// NewInstallCommand creates the install command
func NewInstallCommand(agents, commands, hooks, rules, templates *embed.FS) *cobra.Command {
	var (
		installPath string
		toolType    string
		nonInteractive bool
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
			inst := installer.New(agents, commands, hooks, rules, templates)

			// Interactive mode
			if !nonInteractive {
				// Select tool type
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
					return fmt.Errorf("tool selection cancelled: %w", err)
				}

				if selectedTool == "cancel" {
					fmt.Println(infoStyle.Render("Installation cancelled"))
					return nil
				}
				toolType = selectedTool

				// Select installation path
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
					return fmt.Errorf("path selection cancelled: %w", err)
				}

				switch pathChoice {
				case "cancel":
					fmt.Println(infoStyle.Render("Installation cancelled"))
					return nil
				case "default":
					installPath = ""
				case "local":
					installPath = filepath.Join(cwd, ".the-startup")
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
						return fmt.Errorf("path input cancelled: %w", err)
					}
					installPath = customPath
				}

				// Select components
				var selectedComponents []string
				componentForm := huh.NewForm(
					huh.NewGroup(
						huh.NewMultiSelect[string]().
							Title("Select components to install").
							Options(
								huh.NewOption("Agents (12 specialized agents)", "agents").Selected(true),
								huh.NewOption("Hooks (logging and tracking)", "hooks").Selected(true),
								huh.NewOption("Commands (develop, start)", "commands").Selected(true),
								huh.NewOption("Rules (context management)", "rules").Selected(true),
							).
							Value(&selectedComponents),
					),
				)

				if err := componentForm.Run(); err != nil {
					return fmt.Errorf("component selection cancelled: %w", err)
				}

				// Set installer options
				inst.SetComponents(selectedComponents)
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

	return cmd
}
