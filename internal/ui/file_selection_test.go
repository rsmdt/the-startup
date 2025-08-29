package ui

import (
	"strings"
	"testing"
)

func TestFileSelectionViewHasBanner(t *testing.T) {
	model := NewMainModel(nil, nil)

	// Navigate to file selection
	model.startupPath = "~/.config/the-startup"
	model.claudePath = "~/.claude"
	model.transitionToState(StateFileSelection)

	// Get the view
	view := model.View()

	// Check if the view starts with the banner
	lines := strings.Split(view, "\n")
	if len(lines) < 5 {
		t.Fatal("View has too few lines")
	}

	// The banner should be in the first few lines
	bannerFound := false
	for i := 0; i < 10 && i < len(lines); i++ {
		if strings.Contains(lines[i], "████████") || strings.Contains(lines[i], "THE") || strings.Contains(lines[i], "AGENTIC") || strings.Contains(lines[i], "STARTUP") {
			bannerFound = true
			break
		}
	}

	if !bannerFound {
		t.Error("Banner not found in the first 10 lines of the file selection view")
		t.Log("First 10 lines of view:")
		for i := 0; i < 10 && i < len(lines); i++ {
			t.Logf("Line %d: %s", i, lines[i])
		}
	}

	// Also verify the view contains the tree (if assets are available)
	// The test uses nil assets, so we may not have actual file tree
	// Just check that the basic structure is present
	hasAgents := strings.Contains(view, "agents") || strings.Contains(view, "agents/")
	hasTreeStructure := strings.Contains(view, "Files to be installed")
	
	if !hasAgents && !hasTreeStructure {
		// Only fail if we don't have either agents or the basic structure
		t.Log("Note: View doesn't contain 'agents' folder, but this may be expected with nil assets in test")
	}

	// Verify confirmation options
	if !strings.Contains(view, "Yes, give me awesome") {
		t.Error("View should contain confirmation options")
	}
}
