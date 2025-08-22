package cmd

import (
	"embed"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewUninstallCommand(t *testing.T) {
	// Mock empty embedded filesystems
	var claudeAssets embed.FS
	var startupAssets embed.FS

	cmd := NewUninstallCommand(&claudeAssets, &startupAssets)

	// Test command properties
	if cmd.Use != "uninstall" {
		t.Errorf("Expected Use to be 'uninstall', got %s", cmd.Use)
	}

	if cmd.Short != "Uninstall The Startup agent system" {
		t.Errorf("Expected Short description, got %s", cmd.Short)
	}

	// Test that Long description contains expected content
	expectedContent := []string{
		"Uninstall agents, hooks, and commands",
		"interactive TUI",
		"dry-run",
		"lock file",
		"selective removal",
	}

	for _, content := range expectedContent {
		if !strings.Contains(cmd.Long, content) {
			t.Errorf("Expected Long description to contain '%s'", content)
		}
	}

	// Test flags are registered
	expectedFlags := []string{"dry-run", "force", "keep-logs", "keep-settings"}
	for _, flagName := range expectedFlags {
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to be registered", flagName)
		}
	}

	// Test flag defaults
	if flag := cmd.Flags().Lookup("dry-run"); flag.DefValue != "false" {
		t.Errorf("Expected dry-run default to be false, got %s", flag.DefValue)
	}
}

func TestUninstallCommandExecution(t *testing.T) {
	var claudeAssets embed.FS
	var startupAssets embed.FS

	cmd := NewUninstallCommand(&claudeAssets, &startupAssets)
	
	// Test basic execution
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Since we're using direct output to stdout/stderr, we can't easily capture it
	// But we can verify the command executed without error, which indicates proper setup
}

func TestUninstallCommandFlags(t *testing.T) {
	var claudeAssets embed.FS
	var startupAssets embed.FS

	cmd := NewUninstallCommand(&claudeAssets, &startupAssets)
	cmd.SetArgs([]string{"--dry-run", "--keep-logs"})

	// Test execution with flags
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test that flags are properly parsed by checking the flag objects
	dryRunFlag := cmd.Flags().Lookup("dry-run")
	if dryRunFlag == nil {
		t.Error("Expected dry-run flag to exist")
	}
	
	keepLogsFlag := cmd.Flags().Lookup("keep-logs")
	if keepLogsFlag == nil {
		t.Error("Expected keep-logs flag to exist")
	}
}

func TestUninstallCommandIntegration(t *testing.T) {
	// Test that the command integrates properly with Cobra
	var claudeAssets embed.FS
	var startupAssets embed.FS

	rootCmd := &cobra.Command{Use: "the-startup"}
	rootCmd.AddCommand(NewUninstallCommand(&claudeAssets, &startupAssets))

	// Find the uninstall command
	uninstallCmd, _, err := rootCmd.Find([]string{"uninstall"})
	if err != nil {
		t.Errorf("Expected to find uninstall command, got error: %v", err)
	}

	if uninstallCmd == nil {
		t.Error("Expected to find uninstall command, got nil")
	}

	if uninstallCmd.Use != "uninstall" {
		t.Errorf("Expected found command to have Use 'uninstall', got %s", uninstallCmd.Use)
	}
}