package ui

import (
	"strings"
	"testing"
)

func TestClaudePathModelAlwaysShowsBothOptions(t *testing.T) {
	// Test that both options are always shown regardless of startup path
	testCases := []struct {
		name        string
		startupPath string
	}{
		{"global startup path", "/Users/test/.config/the-startup"},
		{"local startup path", "/Users/test/project/.the-startup"},
		{"custom startup path", "/opt/the-startup"},
		{"path with .config in it", "/Users/test/.config/project/.the-startup"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := NewClaudePathModel(tc.startupPath)

			// Should always have exactly 4 choices
			if len(model.choices) != 4 {
				t.Errorf("Expected 4 choices, got %d", len(model.choices))
				t.Logf("Choices: %v", model.choices)
			}

			// Check for specific options
			expectedChoices := []string{
				"~/.claude (recommended)",
				".claude (local)",
				"Custom location",
				"Cancel",
			}

			for i, expected := range expectedChoices {
				if i >= len(model.choices) {
					t.Errorf("Missing choice at index %d: expected %s", i, expected)
					continue
				}
				if model.choices[i] != expected {
					t.Errorf("Choice at index %d: expected %s, got %s", i, expected, model.choices[i])
				}
			}
		})
	}
}

func TestClaudePathModelView(t *testing.T) {
	model := NewClaudePathModel("/any/path")
	view := model.View()

	// Check that view contains both options
	if !strings.Contains(view, "~/.claude (recommended)") {
		t.Error("View should contain '~/.claude (recommended)' option")
	}
	if !strings.Contains(view, ".claude (local)") {
		t.Error("View should contain '.claude (local)' option")
	}
	if !strings.Contains(view, "Custom location") {
		t.Error("View should contain 'Custom location' option")
	}
	if !strings.Contains(view, "Cancel") {
		t.Error("View should contain 'Cancel' option")
	}
}
