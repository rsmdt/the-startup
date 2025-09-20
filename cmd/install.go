package cmd

import (
	"embed"

	"github.com/rsmdt/the-startup/internal/ui"
	"github.com/spf13/cobra"
)

// NewInstallCommand creates the install command
func NewInstallCommand(claudeAssets, startupAssets *embed.FS) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install The Startup agent system with interactive component selection",
		Long:  `Install agents, commands, and tools for development with interactive selection.

Choose from 58+ specialized agents across 9 development roles, 4 powerful commands,
and output styles. Select role-based presets or create custom combinations.

The installer provides:
• Interactive component selection using charm/huh forms
• Role-based presets (Frontend, Backend, Full Stack, etc.)
• Custom selection with category browsing
• Preview and confirmation before installation

Use --yes to skip interactive selection and install everything.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			local, _ := cmd.Flags().GetBool("local")
			yes, _ := cmd.Flags().GetBool("yes")

			// Run the bubbletea installer UI with enhanced file selection
			return ui.RunMainInstallerWithFlags(claudeAssets, startupAssets, local, yes)
		},
	}

	// Add flags
	cmd.Flags().BoolP("local", "l", false, "Use local installation paths for both directories (skip path selection screens)")
	cmd.Flags().BoolP("yes", "y", false, "Auto-confirm installation with recommended (global) paths")

	return cmd
}
