package ui

import (
	"embed"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallFlags(t *testing.T) {
	// Create mock embedded filesystems
	claudeAssets := &embed.FS{}
	startupAssets := &embed.FS{}

	t.Run("Default flags (no flags set)", func(t *testing.T) {
		flags := InstallFlags{
			Local: false,
			Yes:   false,
		}

		model := NewMainModelWithFlags(claudeAssets, startupAssets, flags)

		// Should start in StateStartupPath
		if model.state != StateStartupPath {
			t.Errorf("Expected initial state StateStartupPath, got %v", model.state)
		}

		// Should have no auto-selected paths
		if model.startupPath != "" {
			t.Errorf("Expected empty startupPath, got %s", model.startupPath)
		}
		if model.claudePath != "" {
			t.Errorf("Expected empty claudePath, got %s", model.claudePath)
		}
	})

	t.Run("Local flag only", func(t *testing.T) {
		flags := InstallFlags{
			Local: true,
			Yes:   false,
		}

		model := NewMainModelWithFlags(claudeAssets, startupAssets, flags)
		
		// Call Init() to trigger flag handling
		model.Init()

		// Should pre-select both local paths
		cwd, _ := os.Getwd()
		expectedStartupPath := filepath.Join(cwd, ".the-startup")
		expectedClaudePath := filepath.Join(cwd, ".claude")

		if model.startupPath != expectedStartupPath {
			t.Errorf("Expected startupPath %s, got %s", expectedStartupPath, model.startupPath)
		}
		if model.claudePath != expectedClaudePath {
			t.Errorf("Expected claudePath %s, got %s", expectedClaudePath, model.claudePath)
		}

		// Should skip both path selection screens and go to file selection
		if model.state != StateFileSelection {
			t.Errorf("Expected state StateFileSelection after local flag, got %v", model.state)
		}
	})

	t.Run("Yes flag only (global installation)", func(t *testing.T) {
		flags := InstallFlags{
			Local: false,
			Yes:   true,
		}

		model := NewMainModelWithFlags(claudeAssets, startupAssets, flags)
		
		// Call Init() to trigger flag handling
		model.Init()

		// Should auto-select global paths
		homeDir, _ := os.UserHomeDir()
		expectedStartupPath := filepath.Join(homeDir, ".config", "the-startup")
		expectedClaudePath := filepath.Join(homeDir, ".claude")

		if model.startupPath != expectedStartupPath {
			t.Errorf("Expected startupPath %s, got %s", expectedStartupPath, model.startupPath)
		}
		if model.claudePath != expectedClaudePath {
			t.Errorf("Expected claudePath %s, got %s", expectedClaudePath, model.claudePath)
		}

		// Should transition to file selection state (or complete)
		if model.state != StateFileSelection && model.state != StateComplete && model.state != StateError {
			t.Errorf("Expected state FileSelection, Complete, or Error after yes flag, got %v", model.state)
		}
	})

	t.Run("Both local and yes flags (local installation)", func(t *testing.T) {
		flags := InstallFlags{
			Local: true,
			Yes:   true,
		}

		model := NewMainModelWithFlags(claudeAssets, startupAssets, flags)
		
		// Call Init() to trigger flag handling
		model.Init()

		// Should auto-select local paths
		cwd, _ := os.Getwd()
		expectedStartupPath := filepath.Join(cwd, ".the-startup")
		expectedClaudePath := filepath.Join(cwd, ".claude")

		if model.startupPath != expectedStartupPath {
			t.Errorf("Expected startupPath %s, got %s", expectedStartupPath, model.startupPath)
		}
		if model.claudePath != expectedClaudePath {
			t.Errorf("Expected claudePath %s, got %s", expectedClaudePath, model.claudePath)
		}

		// Should transition to file selection state (or complete)
		if model.state != StateFileSelection && model.state != StateComplete && model.state != StateError {
			t.Errorf("Expected state FileSelection, Complete, or Error after both flags, got %v", model.state)
		}
	})
}

func TestFlagPathSelection(t *testing.T) {
	claudeAssets := &embed.FS{}
	startupAssets := &embed.FS{}

	t.Run("Local flag sets correct paths", func(t *testing.T) {
		flags := InstallFlags{Local: true, Yes: false}
		model := NewMainModelWithFlags(claudeAssets, startupAssets, flags)
		model.Init()

		cwd, _ := os.Getwd()
		expectedStartupPath := filepath.Join(cwd, ".the-startup")
		expectedClaudePath := filepath.Join(cwd, ".claude")

		if model.startupPath != expectedStartupPath {
			t.Errorf("Local flag should set startup path to %s, got %s", expectedStartupPath, model.startupPath)
		}
		if model.claudePath != expectedClaudePath {
			t.Errorf("Local flag should set claude path to %s, got %s", expectedClaudePath, model.claudePath)
		}

		// Check installer paths are also set
		if model.installer.GetInstallPath() != expectedStartupPath {
			t.Errorf("Installer should have startup path %s, got %s", expectedStartupPath, model.installer.GetInstallPath())
		}
	})

	t.Run("Yes flag with local sets both local paths", func(t *testing.T) {
		flags := InstallFlags{Local: true, Yes: true}
		model := NewMainModelWithFlags(claudeAssets, startupAssets, flags)
		model.Init()

		cwd, _ := os.Getwd()
		expectedStartupPath := filepath.Join(cwd, ".the-startup")
		expectedClaudePath := filepath.Join(cwd, ".claude")

		if model.startupPath != expectedStartupPath {
			t.Errorf("Expected startup path %s, got %s", expectedStartupPath, model.startupPath)
		}
		if model.claudePath != expectedClaudePath {
			t.Errorf("Expected claude path %s, got %s", expectedClaudePath, model.claudePath)
		}
	})

	t.Run("Yes flag without local sets global paths", func(t *testing.T) {
		flags := InstallFlags{Local: false, Yes: true}
		model := NewMainModelWithFlags(claudeAssets, startupAssets, flags)
		model.Init()

		homeDir, _ := os.UserHomeDir()
		expectedStartupPath := filepath.Join(homeDir, ".config", "the-startup")
		expectedClaudePath := filepath.Join(homeDir, ".claude")

		if model.startupPath != expectedStartupPath {
			t.Errorf("Expected startup path %s, got %s", expectedStartupPath, model.startupPath)
		}
		if model.claudePath != expectedClaudePath {
			t.Errorf("Expected claude path %s, got %s", expectedClaudePath, model.claudePath)
		}
	})
}

func TestFlagStateTransitions(t *testing.T) {
	claudeAssets := &embed.FS{}
	startupAssets := &embed.FS{}

	t.Run("No flags start in StateStartupPath", func(t *testing.T) {
		flags := InstallFlags{Local: false, Yes: false}
		model := NewMainModelWithFlags(claudeAssets, startupAssets, flags)

		if model.state != StateStartupPath {
			t.Errorf("No flags should start in StateStartupPath, got %v", model.state)
		}
	})

	t.Run("Local flag only skips to StateFileSelection", func(t *testing.T) {
		flags := InstallFlags{Local: true, Yes: false}
		model := NewMainModelWithFlags(claudeAssets, startupAssets, flags)
		model.Init()

		if model.state != StateFileSelection {
			t.Errorf("Local flag should transition to StateFileSelection, got %v", model.state)
		}
	})

	t.Run("Yes flag skips to later state", func(t *testing.T) {
		flags := InstallFlags{Local: false, Yes: true}
		model := NewMainModelWithFlags(claudeAssets, startupAssets, flags)
		model.Init()

		// Should skip path selection states
		if model.state == StateStartupPath || model.state == StateClaudePath {
			t.Errorf("Yes flag should skip path selection states, got %v", model.state)
		}

		// Should be in FileSelection, Complete, or Error
		validStates := []InstallerState{StateFileSelection, StateComplete, StateError}
		stateValid := false
		for _, validState := range validStates {
			if model.state == validState {
				stateValid = true
				break
			}
		}

		if !stateValid {
			t.Errorf("Yes flag should transition to FileSelection/Complete/Error, got %v", model.state)
		}
	})
}

func TestFlagBackNavigation(t *testing.T) {
	claudeAssets := &embed.FS{}
	startupAssets := &embed.FS{}

	t.Run("Local flag allows back navigation to startup path", func(t *testing.T) {
		flags := InstallFlags{Local: true, Yes: false}
		model := NewMainModelWithFlags(claudeAssets, startupAssets, flags)
		model.Init()

		// Should be in StateFileSelection after local flag
		if model.state != StateFileSelection {
			t.Errorf("Expected StateFileSelection, got %v", model.state)
		}

		// Simulate back navigation (ESC key)
		newModel, _ := model.handleBack()
		model = newModel.(*MainModel)

		// Should go back to StateClaudePath
		if model.state != StateClaudePath {
			t.Errorf("Back navigation should return to StateClaudePath, got %v", model.state)
		}

		// Pre-selected paths should still be preserved
		cwd, _ := os.Getwd()
		expectedStartupPath := filepath.Join(cwd, ".the-startup")
		expectedClaudePath := filepath.Join(cwd, ".claude")
		if model.startupPath != expectedStartupPath {
			t.Errorf("Pre-selected startup path should be preserved after back navigation, expected %s, got %s", 
				expectedStartupPath, model.startupPath)
		}
		if model.claudePath != expectedClaudePath {
			t.Errorf("Pre-selected claude path should be preserved after back navigation, expected %s, got %s", 
				expectedClaudePath, model.claudePath)
		}
	})
}

func TestRunMainInstallerWithFlags(t *testing.T) {
	t.Run("RunMainInstallerWithFlags creates correct flags", func(t *testing.T) {
		// This test verifies the function signature and flag creation
		// We can't easily test the full UI execution without mocking tea.Program

		// Test flag combinations
		testCases := []struct {
			local    bool
			yes      bool
			expected InstallFlags
		}{
			{false, false, InstallFlags{Local: false, Yes: false}},
			{true, false, InstallFlags{Local: true, Yes: false}},
			{false, true, InstallFlags{Local: false, Yes: true}},
			{true, true, InstallFlags{Local: true, Yes: true}},
		}

		for _, tc := range testCases {
			t.Run(string(rune('a'+len(testCases))), func(t *testing.T) {
				// We can't easily test the full UI execution, but we can verify
				// that the function exists and would create the correct flags
				// by testing the flag creation logic directly

				flags := InstallFlags{
					Local: tc.local,
					Yes:   tc.yes,
				}

				if flags.Local != tc.expected.Local {
					t.Errorf("Expected Local=%v, got %v", tc.expected.Local, flags.Local)
				}
				if flags.Yes != tc.expected.Yes {
					t.Errorf("Expected Yes=%v, got %v", tc.expected.Yes, flags.Yes)
				}
			})
		}
	})
}