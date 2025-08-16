package cmd

import (
	"embed"

	"github.com/rsmdt/the-startup/internal/ui"
	"github.com/spf13/cobra"
)

// NewInstallCommand creates the install command
func NewInstallCommand(agents, commands, templates, settings *embed.FS) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install The Startup agent system",
		Long:  `Install agents, hooks, and commands for development tools with an interactive TUI`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Run the new bubbletea installer UI with composable views
			return ui.RunMainInstaller(agents, commands, templates, settings)
		},
	}

	return cmd
}
