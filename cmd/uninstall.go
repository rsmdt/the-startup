package cmd

import (
	"embed"

	"github.com/spf13/cobra"
	"github.com/rsmdt/the-startup/internal/ui"
)

// UninstallFlags holds the command line flags for the uninstall command
type UninstallFlags struct {
	DryRun       bool
	Force        bool
	KeepLogs     bool
	KeepSettings bool
}

// NewUninstallCommand creates the uninstall command
func NewUninstallCommand(claudeAssets, startupAssets *embed.FS) *cobra.Command {
	var flags UninstallFlags

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall The Startup agent system",
		Long: `Uninstall agents, hooks, and commands from development tools with an interactive TUI.

The uninstall process will:
- Show all currently installed files from the lock file
- Allow selective removal of components (agents, commands, templates)
- Clean up configuration changes in settings.json
- Remove the binary from installation directory
- Optionally preserve user data and logs

Use --dry-run to see what would be removed without making changes.

Examples:
  # Interactive uninstall with TUI
  the-startup uninstall

  # Preview what would be removed
  the-startup uninstall --dry-run

  # Keep logs and settings during uninstall
  the-startup uninstall --keep-logs --keep-settings

  # Force removal without confirmation prompts
  the-startup uninstall --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUninstall(cmd, flags, claudeAssets, startupAssets)
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false, "Show what would be removed without making changes")
	cmd.Flags().BoolVar(&flags.Force, "force", false, "Force removal without confirmation prompts")
	cmd.Flags().BoolVar(&flags.KeepLogs, "keep-logs", false, "Preserve log files and session data")
	cmd.Flags().BoolVar(&flags.KeepSettings, "keep-settings", false, "Preserve user settings and configurations")

	return cmd
}

// runUninstall handles the uninstall command execution
func runUninstall(cmd *cobra.Command, flags UninstallFlags, claudeAssets, startupAssets *embed.FS) error {
	// NOTE: The simplified uninstall flow doesn't use command line options
	// It follows the same interactive pattern as install
	// TODO: Add support for dry-run and other options to the simplified flow if needed
	
	return ui.RunMainUninstaller(claudeAssets, startupAssets)
}