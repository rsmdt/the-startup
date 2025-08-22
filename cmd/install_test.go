package cmd

import (
	"embed"
	"testing"
)

func TestInstallCommandFlags(t *testing.T) {
	// Create mock embedded filesystems
	claudeAssets := &embed.FS{}
	startupAssets := &embed.FS{}

	// Create the install command
	installCmd := NewInstallCommand(claudeAssets, startupAssets)

	t.Run("Install command has correct flags", func(t *testing.T) {
		// Check that the command has the expected flags
		localFlag := installCmd.Flags().Lookup("local")
		if localFlag == nil {
			t.Error("Expected --local flag to be defined")
		} else {
			if localFlag.Shorthand != "l" {
				t.Errorf("Expected --local flag shorthand to be 'l', got '%s'", localFlag.Shorthand)
			}
			if localFlag.DefValue != "false" {
				t.Errorf("Expected --local flag default value to be 'false', got '%s'", localFlag.DefValue)
			}
		}

		yesFlag := installCmd.Flags().Lookup("yes")
		if yesFlag == nil {
			t.Error("Expected --yes flag to be defined")
		} else {
			if yesFlag.Shorthand != "y" {
				t.Errorf("Expected --yes flag shorthand to be 'y', got '%s'", yesFlag.Shorthand)
			}
			if yesFlag.DefValue != "false" {
				t.Errorf("Expected --yes flag default value to be 'false', got '%s'", yesFlag.DefValue)
			}
			if yesFlag.Usage != "Auto-confirm installation with recommended (global) paths" {
				t.Errorf("Expected --yes flag usage to mention global paths, got '%s'", yesFlag.Usage)
			}
		}
	})

	t.Run("Install command has correct usage", func(t *testing.T) {
		if installCmd.Use != "install" {
			t.Errorf("Expected command use to be 'install', got '%s'", installCmd.Use)
		}
		
		if installCmd.Short != "Install The Startup agent system" {
			t.Errorf("Expected short description to be 'Install The Startup agent system', got '%s'", installCmd.Short)
		}
	})

	t.Run("Flag parsing works correctly", func(t *testing.T) {
		// Test cases for different flag combinations
		testCases := []struct {
			name     string
			args     []string
			expected map[string]bool
		}{
			{
				name: "no flags",
				args: []string{},
				expected: map[string]bool{
					"local": false,
					"yes":   false,
				},
			},
			{
				name: "local flag short",
				args: []string{"-l"},
				expected: map[string]bool{
					"local": true,
					"yes":   false,
				},
			},
			{
				name: "local flag long",
				args: []string{"--local"},
				expected: map[string]bool{
					"local": true,
					"yes":   false,
				},
			},
			{
				name: "yes flag short",
				args: []string{"-y"},
				expected: map[string]bool{
					"local": false,
					"yes":   true,
				},
			},
			{
				name: "yes flag long",
				args: []string{"--yes"},
				expected: map[string]bool{
					"local": false,
					"yes":   true,
				},
			},
			{
				name: "both flags",
				args: []string{"-l", "-y"},
				expected: map[string]bool{
					"local": true,
					"yes":   true,
				},
			},
			{
				name: "both flags combined",
				args: []string{"-ly"},
				expected: map[string]bool{
					"local": true,
					"yes":   true,
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create a new command instance for each test
				cmd := NewInstallCommand(claudeAssets, startupAssets)
				
				// Parse the arguments
				cmd.SetArgs(tc.args)
				err := cmd.ParseFlags(tc.args)
				if err != nil {
					t.Fatalf("Failed to parse flags: %v", err)
				}

				// Check flag values
				localVal, err := cmd.Flags().GetBool("local")
				if err != nil {
					t.Fatalf("Failed to get local flag value: %v", err)
				}
				if localVal != tc.expected["local"] {
					t.Errorf("Expected local flag to be %v, got %v", tc.expected["local"], localVal)
				}

				yesVal, err := cmd.Flags().GetBool("yes")
				if err != nil {
					t.Fatalf("Failed to get yes flag value: %v", err)
				}
				if yesVal != tc.expected["yes"] {
					t.Errorf("Expected yes flag to be %v, got %v", tc.expected["yes"], yesVal)
				}
			})
		}
	})
}

func TestInstallCommandHelp(t *testing.T) {
	claudeAssets := &embed.FS{}
	startupAssets := &embed.FS{}

	cmd := NewInstallCommand(claudeAssets, startupAssets)

	t.Run("Help text contains flag descriptions", func(t *testing.T) {
		// Get help text
		helpText := cmd.Long + "\n" + cmd.Flags().FlagUsages()

		// Check that help contains information about both flags
		if !containsAnySubstring(helpText, []string{"local", "-l"}) {
			t.Error("Help text should mention the --local/-l flag")
		}
		if !containsAnySubstring(helpText, []string{"yes", "-y"}) {
			t.Error("Help text should mention the --yes/-y flag")
		}
	})
}

// Helper function to check if a string contains any of the given substrings
func containsAnySubstring(text string, substrings []string) bool {
	for _, substr := range substrings {
		if len(text) >= len(substr) {
			for i := 0; i <= len(text)-len(substr); i++ {
				if text[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}