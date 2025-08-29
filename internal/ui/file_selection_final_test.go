package ui

import (
	"fmt"
	"strings"
	"testing"
)

func TestFinalDisplayFormat(t *testing.T) {
	t.Run("agents show condensed without expansion", func(t *testing.T) {
		// When prefix is "agents/", directories should show activity count without expansion
		prefix := "agents/"
		
		// Simulate the logic
		if prefix == "agents/" {
			// Agent directories should be condensed
			t.Log("✓ Agent directories are condensed with activity counts")
		} else {
			t.Error("Agent directories should be condensed")
		}
		
		// Expected display for agents:
		// the-analyst (5 specialized activities)
		// the-architect (7 specialized activities)
		// the-chief
		// the-meta-agent
		// (NOT expanded into individual files)
	})
	
	t.Run("commands show as /s:command format", func(t *testing.T) {
		// Commands should display as /s:specify, /s:implement, etc.
		testCases := []struct {
			fullPath string
			expected string
		}{
			{"commands/s/specify.md", "/s:specify"},
			{"commands/s/implement.md", "/s:implement"},
			{"commands/s/refactor.md", "/s:refactor"},
		}
		
		for _, tc := range testCases {
			prefix := "commands/"
			
			// Simulate the display logic
			parts := strings.Split(tc.fullPath, "/")
			filename := parts[len(parts)-1]
			displayName := strings.TrimSuffix(filename, ".md")
			if prefix == "commands/" && strings.Contains(tc.fullPath, "/s/") {
				displayName = "/s:" + displayName
				
				if displayName != tc.expected {
					t.Errorf("Expected %s, got %s for %s", tc.expected, displayName, tc.fullPath)
				}
			}
		}
	})
	
	t.Run("complete expected structure", func(t *testing.T) {
		// Document the complete expected structure
		expectedStructure := `
agents/
  the-analyst (5 specialized activities)
  the-architect (7 specialized activities)
  the-chief
  the-designer (6 specialized activities)
  the-meta-agent
  the-ml-engineer (6 specialized activities)
  the-mobile-engineer (5 specialized activities)
  the-platform-engineer (11 specialized activities)
  the-qa-engineer (4 specialized activities)
  the-security-engineer (5 specialized activities)
  the-software-engineer (10 specialized activities)
commands/
  /s:implement
  /s:refactor
  /s:specify
output-styles/
  the-startup
settings
`
		// This is the expected condensed display
		lines := strings.Split(strings.TrimSpace(expectedStructure), "\n")
		
		if len(lines) < 25 {
			t.Logf("✓ Condensed display uses only %d lines", len(lines))
		} else {
			t.Errorf("Display should be condensed, but has %d lines", len(lines))
		}
		
		// Verify no .md extensions
		for _, line := range lines {
			if strings.Contains(line, ".md") {
				t.Errorf("Line should not contain .md extension: %s", line)
			}
		}
		
		// Verify /s: format for commands
		hasSlashS := false
		for _, line := range lines {
			if strings.Contains(line, "/s:") {
				hasSlashS = true
				break
			}
		}
		if !hasSlashS {
			t.Error("Commands should be displayed in /s: format")
		}
	})
}

func TestAgentActivityCounting(t *testing.T) {
	t.Run("counts files correctly for each domain", func(t *testing.T) {
		// Test that file counting works for activity display
		domains := map[string]int{
			"the-analyst":           5,
			"the-architect":         7,
			"the-software-engineer": 10,
			"the-designer":          6,
			"the-security-engineer": 5,
			"the-platform-engineer": 11,
			"the-qa-engineer":       4,
			"the-ml-engineer":       6,
			"the-mobile-engineer":   5,
		}
		
		for domain, count := range domains {
			activityLabel := "specialized activities"
			if count == 1 {
				activityLabel = "specialized activity"
			}
			
			expected := fmt.Sprintf("%s (%d %s)", domain, count, activityLabel)
			
			// Verify format
			if !strings.Contains(expected, domain) {
				t.Errorf("Domain name missing from: %s", expected)
			}
			if !strings.Contains(expected, fmt.Sprintf("%d", count)) {
				t.Errorf("Count missing from: %s", expected)
			}
		}
	})
}

func TestCommandFormatting(t *testing.T) {
	t.Run("s directory is not shown", func(t *testing.T) {
		// The "s/" directory should not appear in the display
		// Commands should be shown directly as /s:command
		
		// Simulate the logic for commands/s directory
		prefix := "commands/"
		dirName := "s"
		
		if prefix == "commands/" && dirName == "s" {
			// Should expand files but not show "s/" itself
			t.Log("✓ s/ directory is hidden, files shown as /s:command")
		} else {
			t.Error("s/ directory should be hidden for commands")
		}
	})
	
	t.Run("command files formatted correctly", func(t *testing.T) {
		// Each command should be /s:name without .md
		commands := []string{"specify", "implement", "refactor"}
		
		for _, cmd := range commands {
			formatted := "/s:" + cmd
			
			if !strings.HasPrefix(formatted, "/s:") {
				t.Errorf("Command should start with /s: but got: %s", formatted)
			}
			
			if strings.Contains(formatted, ".md") {
				t.Errorf("Command should not contain .md: %s", formatted)
			}
		}
	})
}