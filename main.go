package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/the-startup/the-startup/cmd"
)

var (
	// Version info, set by build flags
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// Embedded assets
//
//go:embed agents/*.md
var agentFiles embed.FS

//go:embed commands/*.md
var commandFiles embed.FS

//go:embed hooks/*.py
var hookFiles embed.FS

//go:embed rules/*.md
var ruleFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

func main() {
	rootCmd := &cobra.Command{
		Use:   "the-startup",
		Short: "Agent system for development tools",
		Long: `The Startup - A comprehensive agent system for enhancing development 
workflows with specialized AI agents, hooks, and commands.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, GitCommit, BuildDate),
	}

	// Add commands
	rootCmd.AddCommand(cmd.NewInstallCommand(&agentFiles, &commandFiles, &hookFiles, &ruleFiles, &templateFiles))
	rootCmd.AddCommand(cmd.NewUpdateCommand())
	rootCmd.AddCommand(cmd.NewValidateCommand())
	rootCmd.AddCommand(cmd.NewHooksCommand())

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}