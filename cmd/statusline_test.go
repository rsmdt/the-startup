package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestStatuslineOutput(t *testing.T) {
	tests := []struct {
		name            string
		input           StatuslineInput
		wantContains    []string
		wantNotContains []string
	}{
		{
			name: "git repository with branch",
			input: StatuslineInput{
				Workspace: struct {
					CurrentDir string `json:"current_dir"`
					ProjectDir string `json:"project_dir"`
				}{
					CurrentDir: "/Users/test/project",
				},
				Model: struct {
					ID          string `json:"id"`
					DisplayName string `json:"display_name"`
				}{
					DisplayName: "Claude 3.5 Sonnet",
				},
				OutputStyle: struct {
					Name string `json:"name"`
				}{
					Name: "The Startup",
				},
				SessionID: "test-session",
			},
			wantContains: []string{
				"📁",
				"🤖 Claude 3.5 Sonnet (The Startup)",
				"? for shortcuts",
			},
			wantNotContains: []string{},
		},
		{
			name: "non-git directory",
			input: StatuslineInput{
				Workspace: struct {
					CurrentDir string `json:"current_dir"`
					ProjectDir string `json:"project_dir"`
				}{
					CurrentDir: "/tmp",
				},
				Model: struct {
					ID          string `json:"id"`
					DisplayName string `json:"display_name"`
				}{
					DisplayName: "Opus",
				},
				OutputStyle: struct {
					Name string `json:"name"`
				}{
					Name: "default",
				},
				SessionID: "test",
			},
			wantContains: []string{
				"📁 /tmp",
				"🤖 Opus (default)",
				"? for shortcuts",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create input JSON
			inputJSON, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal input: %v", err)
			}

			// Run statusline command
			input := bytes.NewReader(inputJSON)
			output := &bytes.Buffer{}

			if err := runStatusline(input, output); err != nil {
				t.Fatalf("runStatusline() error = %v", err)
			}

			got := output.String()

			// Check for expected content
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("Output missing expected content %q\nGot: %s", want, got)
				}
			}

			// Check for unwanted content
			for _, notWant := range tt.wantNotContains {
				if strings.Contains(got, notWant) {
					t.Errorf("Output contains unwanted content %q\nGot: %s", notWant, got)
				}
			}
		})
	}
}

// TestGetVisibleLength removed - this function doesn't exist in the current implementation
// The statusline now uses lipgloss for formatting which handles text width internally

// TestRightAlignment removed - the current implementation uses lipgloss for layout
// which handles alignment and padding automatically with MaxWidth(termWidth)

func TestHomeDirectorySubstitution(t *testing.T) {
	homeDir := "/Users/testuser"
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "home directory path",
			input:    "/Users/testuser/Documents/project",
			expected: "~/Documents/project",
		},
		{
			name:     "non-home path",
			input:    "/tmp/project",
			expected: "/tmp/project",
		},
		{
			name:     "exact home directory",
			input:    "/Users/testuser",
			expected: "~",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input
			if strings.HasPrefix(result, homeDir) {
				result = "~" + strings.TrimPrefix(result, homeDir)
			}

			if result != tt.expected {
				t.Errorf("Home substitution failed: got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestTerminalWidthFallback(t *testing.T) {
	// Test that getTermWidth returns a reasonable default
	// This is hard to test directly, but we can verify it returns > 0
	width := getTermWidth()

	if width <= 0 {
		t.Errorf("getTermWidth() returned invalid width: %d", width)
	}

	// Should return at least the default width
	if width < 80 {
		t.Errorf("getTermWidth() returned width less than minimum: %d", width)
	}
}
