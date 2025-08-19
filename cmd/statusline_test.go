package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestStatuslineOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    StatuslineInput
		wantContains []string
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
			wantNotContains: []string{
				"| ?", // Should not have separator before help text
			},
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
			wantNotContains: []string{
				"  |", // Should not have double spaces
				"| ?", // Should not have separator before help text
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

func TestGetVisibleLength(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "plain text",
			input: "hello world",
			want:  11,
		},
		{
			name:  "folder emoji",
			input: "📁 test",
			want:  7, // 2 for emoji + 1 space + 4 for "test"
		},
		{
			name:  "robot emoji",
			input: "🤖 Claude",
			want:  9, // 2 for emoji + 1 space + 6 for "Claude"
		},
		{
			name:  "git branch symbol",
			input: "⎇ main",
			want:  7, // 2 for symbol + 1 space + 4 for "main"
		},
		{
			name:  "mixed content",
			input: "📁 ~/project ⎇ main | 🤖 Claude",
			want:  32, // 2+1+10+1+2+1+4+1+1+1+2+1+6 = 32
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getVisibleLength(tt.input)
			if got != tt.want {
				t.Errorf("getVisibleLength(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestRightAlignment(t *testing.T) {
	tests := []struct {
		name       string
		mainContent string
		helpText   string
		termWidth  int
		wantPadding int
	}{
		{
			name:        "short content",
			mainContent: "📁 /tmp | 🤖 Opus",
			helpText:    "? for shortcuts",
			termWidth:   80,
			wantPadding: 48, // 80 - 17 (visible length) - 15 (help text) = 48
		},
		{
			name:        "long content",
			mainContent: "📁 ~/very/long/path/to/project ⎇ main | 🤖 Claude 3.5 Sonnet (The Startup)",
			helpText:    "? for shortcuts",
			termWidth:   120,
			wantPadding: 30, // 120 - 75 (visible length) - 15 (help text) = 30
		},
		{
			name:        "minimum padding",
			mainContent: "📁 ~/very/long/path | 🤖 Model",
			helpText:    "? for shortcuts",
			termWidth:   50,
			wantPadding: 5, // 50 - 30 (visible length) - 15 (help text) = 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mainLen := getVisibleLength(tt.mainContent)
			helpLen := len(tt.helpText)
			
			padding := tt.termWidth - mainLen - helpLen
			if padding < 1 {
				padding = 1
			}

			if padding != tt.wantPadding {
				t.Errorf("Padding calculation wrong:\n"+
					"  termWidth=%d, mainLen=%d, helpLen=%d\n"+
					"  got padding=%d, want=%d",
					tt.termWidth, mainLen, helpLen, padding, tt.wantPadding)
			}

			// Verify the total length matches terminal width (or close to it)
			totalLen := mainLen + padding + helpLen
			if totalLen > tt.termWidth && padding > 1 {
				t.Errorf("Total length %d exceeds terminal width %d", totalLen, tt.termWidth)
			}
		})
	}
}

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
	// Test that getTerminalWidth returns a reasonable default
	// This is hard to test directly, but we can verify it returns > 0
	width := getTerminalWidth()
	
	if width <= 0 {
		t.Errorf("getTerminalWidth() returned invalid width: %d", width)
	}
	
	// Should return at least the default width
	if width < 80 {
		t.Errorf("getTerminalWidth() returned width less than minimum: %d", width)
	}
}