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
		Short: "Install The Startup agent system",
		Long:  `Install agents, hooks, and commands for development tools with an interactive TUI`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			local, _ := cmd.Flags().GetBool("local")
			yes, _ := cmd.Flags().GetBool("yes")
			
			// Run the new bubbletea installer UI with composable views
			return ui.RunMainInstallerWithFlags(claudeAssets, startupAssets, local, yes)
		},
	}

	// Add flags
	cmd.Flags().BoolP("local", "l", false, "Use local installation paths for both directories (skip path selection screens)")
	cmd.Flags().BoolP("yes", "y", false, "Auto-confirm installation with recommended (global) paths")

	return cmd
}
