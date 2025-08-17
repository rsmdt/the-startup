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

// Embedded assets
//
//go:embed assets/agents/*.md
var agentFiles embed.FS

//go:embed assets/commands/**/*.md
var commandFiles embed.FS

//go:embed assets/templates/*
var templateFiles embed.FS

//go:embed assets/rules/*.md
var rulesFiles embed.FS

//go:embed assets/settings.json
var settingsTemplate embed.FS

func main() {
	rootCmd := &cobra.Command{
		Use:   "the-startup",
		Short: "Agent system for development tools",
		Long: `The Startup - A comprehensive agent system for enhancing development 
workflows with specialized AI agents, hooks, and commands.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, GitCommit, BuildDate),
	}

	// Add commands
	rootCmd.AddCommand(cmd.NewInstallCommand(&agentFiles, &commandFiles, &templateFiles, &rulesFiles, &settingsTemplate))
	rootCmd.AddCommand(cmd.NewUpdateCommand())
	rootCmd.AddCommand(cmd.NewValidateCommand())
	rootCmd.AddCommand(cmd.NewHooksCommand())
	rootCmd.AddCommand(cmd.NewLogCommand())

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
