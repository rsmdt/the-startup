package cmd

import (
	"embed"

	"github.com/spf13/cobra"
	"github.com/the-startup/the-startup/internal/ui"
)

// NewInstallCommand creates the install command
func NewInstallCommand(agents, commands, hooks, templates *embed.FS) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install The Startup agent system",
		Long:  `Install agents, hooks, and commands for development tools with an interactive TUI`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Run the new bubbletea installer UI with composable views
			return ui.RunMainInstaller(agents, commands, hooks, templates)
		},
	}

	return cmd
}
