package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/rsmdt/the-startup/cmd"
	"github.com/spf13/cobra"
)

var (
	// Version info, set by build flags
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// Embedded assets - Claude configuration
//
//go:embed assets/claude
var claudeAssets embed.FS

// Embedded assets - The Startup configuration
//
//go:embed assets/the-startup
var startupAssets embed.FS

func main() {
	rootCmd := &cobra.Command{
		Use:   "the-startup",
		Short: "Agent system for development tools",
		Long: `The Startup - A comprehensive agent system for enhancing development 
workflows with specialized AI agents, hooks, and commands.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, GitCommit, BuildDate),
	}

	// Add commands
	rootCmd.AddCommand(cmd.NewInstallCommand(&claudeAssets, &startupAssets))
	rootCmd.AddCommand(cmd.NewUninstallCommand(&claudeAssets, &startupAssets))
	rootCmd.AddCommand(cmd.NewInitCommand(&startupAssets))
	rootCmd.AddCommand(cmd.NewStatsCommand())
	rootCmd.AddCommand(cmd.NewStatuslineCommand())
	rootCmd.AddCommand(cmd.NewSpecCommand(&startupAssets))

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
