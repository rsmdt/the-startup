package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAgentFileManager tests the AgentFileManager structure and its methods
func TestAgentFileManager(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("GetAgentFilePath", func(t *testing.T) {
		afm := &AgentFileManager{
			SessionID: "dev-20250812-143022",
			AgentID:   "arch-001",
			BaseDir:   tempDir,
		}

		expected := filepath.Join(tempDir, "dev-20250812-143022", "arch-001.jsonl")
		actual := afm.GetAgentFilePath()
		assert.Equal(t, expected, actual)
	})

	t.Run("GetSessionDir", func(t *testing.T) {
		afm := &AgentFileManager{
			SessionID: "dev-20250812-143022",
			AgentID:   "arch-001",
			BaseDir:   tempDir,
		}

		expected := filepath.Join(tempDir, "dev-20250812-143022")
		actual := afm.GetSessionDir()
		assert.Equal(t, expected, actual)
	})

	t.Run("EnsureSessionDir", func(t *testing.T) {
		afm := &AgentFileManager{
			SessionID: "dev-20250812-143022",
			AgentID:   "arch-001",
			BaseDir:   tempDir,
		}

		err := afm.EnsureSessionDir()
		assert.NoError(t, err)

		// Verify directory was created with proper permissions
		sessionDir := filepath.Join(tempDir, "dev-20250812-143022")
		info, err := os.Stat(sessionDir)
		assert.NoError(t, err)
		assert.True(t, info.IsDir())
		assert.Equal(t, os.FileMode(0755), info.Mode().Perm())
	})
}

// TestWriteAgentContext tests the WriteAgentContext function
func TestWriteAgentContext(t *testing.T) {
	// Set up temporary directory for testing
	tempDir := t.TempDir()
	// Use environment variable to control GetProjectDir behavior
	originalClaudeDir := os.Getenv("CLAUDE_PROJECT_DIR")
	defer os.Setenv("CLAUDE_PROJECT_DIR", originalClaudeDir)
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)

	// Create the .the-startup directory so GetStartupDir will use project-local directory
	startupDir := filepath.Join(tempDir, ".the-startup")
	require.NoError(t, os.MkdirAll(startupDir, 0755))

	sessionID := "dev-test-session"
	agentID := "arch-001"

	hookData := &HookData{
		Role:      "user",
		Content:   "Test task content",
		Timestamp: "2025-08-12T14:00:00.000Z",
		SessionID: sessionID,
		AgentID:   agentID,
	}

	t.Run("Successful write", func(t *testing.T) {
		err := WriteAgentContext(sessionID, agentID, hookData)
		assert.NoError(t, err)

		// Verify file exists
		expectedPath := filepath.Join(tempDir, ".the-startup", sessionID, agentID+".jsonl")
		assert.FileExists(t, expectedPath)

		// Verify content
		content, err := os.ReadFile(expectedPath)
		assert.NoError(t, err)

		// Parse the JSONL line (remove trailing newline if present)
		jsonLine := strings.TrimSpace(string(content))
		var parsedData HookData
		err = json.Unmarshal([]byte(jsonLine), &parsedData)
		assert.NoError(t, err)
		// SessionID and AgentID are not serialized (marked with json:"-")
		assert.Equal(t, hookData.Role, parsedData.Role)
		assert.Equal(t, hookData.Content, parsedData.Content)
		assert.Equal(t, hookData.Timestamp, parsedData.Timestamp)
	})

	t.Run("Append operation", func(t *testing.T) {
		// Write second entry
		hookData2 := &HookData{
			Role:      "assistant",
			Content:   "Task completed content",
			Timestamp: "2025-08-12T14:05:00.000Z",
			SessionID: sessionID,
			AgentID:   agentID,
		}

		err := WriteAgentContext(sessionID, agentID, hookData2)
		assert.NoError(t, err)

		// Verify file contains both entries
		expectedPath := filepath.Join(tempDir, ".the-startup", sessionID, agentID+".jsonl")
		content, err := os.ReadFile(expectedPath)
		assert.NoError(t, err)
		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		assert.Len(t, lines, 2)
	})

	t.Run("Parameter validation", func(t *testing.T) {
		// Test empty sessionID
		err := WriteAgentContext("", agentID, hookData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sessionID, agentID, and data required")

		// Test empty agentID
		err = WriteAgentContext(sessionID, "", hookData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sessionID, agentID, and data required")

		// Test nil data
		err = WriteAgentContext(sessionID, agentID, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sessionID, agentID, and data required")
	})
}

// TestReadAgentContextRaw tests the ReadAgentContextRaw function
func TestReadAgentContextRaw(t *testing.T) {
	tempDir := t.TempDir()
	originalClaudeDir := os.Getenv("CLAUDE_PROJECT_DIR")
	defer os.Setenv("CLAUDE_PROJECT_DIR", originalClaudeDir)
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)

	// Create the .the-startup directory so GetStartupDir will use project-local directory
	startupDir := filepath.Join(tempDir, ".the-startup")
	require.NoError(t, os.MkdirAll(startupDir, 0755))

	sessionID := "dev-test-session"
	agentID := "arch-001"

	t.Run("Read existing context", func(t *testing.T) {
		// Create test context file
		sessionDir := filepath.Join(tempDir, ".the-startup", sessionID)
		require.NoError(t, os.MkdirAll(sessionDir, 0755))

		contextFile := filepath.Join(sessionDir, agentID+".jsonl")
		testData := []string{
			`{"role":"user","content":"Task 1","timestamp":"2025-08-12T14:00:00.000Z"}`,
			`{"role":"assistant","content":"Task 1 complete","timestamp":"2025-08-12T14:05:00.000Z"}`,
			`{"role":"user","content":"Task 2","timestamp":"2025-08-12T14:10:00.000Z"}`,
		}

		err := os.WriteFile(contextFile, []byte(strings.Join(testData, "\n")+"\n"), 0644)
		require.NoError(t, err)

		// Test reading with limit
		lines, err := ReadAgentContextRaw(sessionID, agentID, 2)
		assert.NoError(t, err)
		assert.Len(t, lines, 2)

		// Should return most recent lines (last 2 lines)
		assert.Contains(t, lines[0], "Task 1 complete")
		assert.Contains(t, lines[1], "Task 2")
	})

	t.Run("Non-existent agent returns empty", func(t *testing.T) {
		lines, err := ReadAgentContextRaw("nonexistent-session", "nonexistent-agent", 10)
		assert.NoError(t, err)
		assert.Empty(t, lines)
	})

	t.Run("Parameter validation", func(t *testing.T) {
		// Test empty agentID
		lines, err := ReadAgentContextRaw(sessionID, "", 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "agentID is required")
		assert.Nil(t, lines)

		// Test invalid maxLines - should default to 50
		lines2, err2 := ReadAgentContextRaw(sessionID, "some-agent", 0)
		assert.NoError(t, err2)  // Should not error but use default
		assert.NotNil(t, lines2) // Should return empty slice, not nil

		lines3, err3 := ReadAgentContextRaw(sessionID, "some-agent", 2000)
		assert.NoError(t, err3)  // Should not error but use default
		assert.NotNil(t, lines3) // Should return empty slice, not nil
	})

	t.Run("Corrupted JSONL handling", func(t *testing.T) {
		sessionDir := filepath.Join(tempDir, ".the-startup", "corrupt-session")
		require.NoError(t, os.MkdirAll(sessionDir, 0755))

		contextFile := filepath.Join(sessionDir, "arch-corrupt.jsonl")
		corruptData := []string{
			`{"role":"user","content":"start","timestamp":"2025-08-12T14:00:00.000Z"}`,
			`invalid json line`,
			`{"role":"assistant","content":"complete","timestamp":"2025-08-12T14:05:00.000Z"}`,
		}

		err := os.WriteFile(contextFile, []byte(strings.Join(corruptData, "\n")+"\n"), 0644)
		require.NoError(t, err)

		lines, err := ReadAgentContextRaw("corrupt-session", "arch-corrupt", 10)
		assert.NoError(t, err)
		assert.Len(t, lines, 3) // Should return all lines including corrupted one
	})
}

// TestFindLatestSessionWithAgent tests the helper function for finding sessions
func TestFindLatestSessionWithAgent(t *testing.T) {
	tempDir := t.TempDir()
	originalClaudeDir := os.Getenv("CLAUDE_PROJECT_DIR")
	defer os.Setenv("CLAUDE_PROJECT_DIR", originalClaudeDir)
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)

	baseDir := filepath.Join(tempDir, ".the-startup")
	require.NoError(t, os.MkdirAll(baseDir, 0755))
	agentID := "arch-001"

	t.Run("Find latest session with agent", func(t *testing.T) {
		// Create multiple sessions with the agent
		sessions := []string{
			"dev-20250812-140000",
			"dev-20250812-150000", // This should be latest
			"dev-20250812-130000",
		}

		baseTime := time.Now()
		for _, session := range sessions {
			sessionDir := filepath.Join(baseDir, session)
			require.NoError(t, os.MkdirAll(sessionDir, 0755))

			agentFile := filepath.Join(sessionDir, agentID+".jsonl")
			testData := fmt.Sprintf(`{"role":"user","content":"test content","timestamp":"2025-08-12T14:00:00.000Z"}`)
			require.NoError(t, os.WriteFile(agentFile, []byte(testData), 0644))

			// Set specific modification times - make dev-20250812-150000 (index 1) the latest
			var modTime time.Time
			switch session {
			case "dev-20250812-140000":
				modTime = baseTime.Add(-2 * time.Hour) // Oldest
			case "dev-20250812-150000":
				modTime = baseTime // Most recent (latest)
			case "dev-20250812-130000":
				modTime = baseTime.Add(-3 * time.Hour) // Middle
			}

			require.NoError(t, os.Chtimes(agentFile, modTime, modTime))
		}

		latestSession, err := findLatestSessionWithAgent(agentID)
		assert.NoError(t, err)
		assert.Equal(t, "dev-20250812-150000", latestSession)
	})

	t.Run("No session found", func(t *testing.T) {
		session, err := findLatestSessionWithAgent("nonexistent-agent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no session found for agent")
		assert.Empty(t, session)
	})
}

// TestReadLastNLines tests the efficient file reading helper
func TestReadLastNLines(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("Small file - read all", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "small.jsonl")
		testLines := []string{
			"line 1",
			"line 2",
			"line 3",
			"line 4",
			"line 5",
		}

		content := strings.Join(testLines, "\n") + "\n"
		require.NoError(t, os.WriteFile(testFile, []byte(content), 0644))

		lines, err := readLastNLines(testFile, 3)
		assert.NoError(t, err)
		assert.Len(t, lines, 3)
		assert.Equal(t, "line 3", lines[0])
		assert.Equal(t, "line 4", lines[1])
		assert.Equal(t, "line 5", lines[2])
	})

	t.Run("Empty file", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "empty.jsonl")
		require.NoError(t, os.WriteFile(testFile, []byte(""), 0644))

		lines, err := readLastNLines(testFile, 10)
		assert.NoError(t, err)
		assert.Empty(t, lines)
	})

	t.Run("File not exists", func(t *testing.T) {
		lines, err := readLastNLines("/nonexistent/file.jsonl", 10)
		assert.Error(t, err)
		assert.Nil(t, lines)
	})

	t.Run("Large file - use buffered reading", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "large.jsonl")

		// Create a file larger than 1MB threshold
		var lines []string
		for i := 0; i < 50000; i++ {
			// Each line is about 50 bytes, so 50000 lines ~ 2.5MB
			line := fmt.Sprintf(`{"role":"user","content":"Line %05d with some extra padding","timestamp":"2025-08-12T14:00:00.000Z"}`, i)
			lines = append(lines, line)
		}

		content := strings.Join(lines, "\n") + "\n"
		require.NoError(t, os.WriteFile(testFile, []byte(content), 0644))

		result, err := readLastNLines(testFile, 3)
		assert.NoError(t, err)
		assert.Len(t, result, 3)

		// Verify we got the last 3 lines
		expectedLines := lines[len(lines)-3:]
		assert.Equal(t, expectedLines, result)
	})
}

// TestReadTailWithBuffer tests the efficient tail reading for large files
func TestReadTailWithBuffer(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.jsonl")

	// Create test data
	testLines := []string{
		"first line",
		"second line",
		"third line",
		"fourth line",
		"fifth line",
	}

	content := strings.Join(testLines, "\n") + "\n"
	require.NoError(t, os.WriteFile(testFile, []byte(content), 0644))

	file, err := os.Open(testFile)
	require.NoError(t, err)
	defer file.Close()

	stat, err := file.Stat()
	require.NoError(t, err)

	t.Run("Read last N lines", func(t *testing.T) {
		lines, err := readTailWithBuffer(file, stat.Size(), 3)
		assert.NoError(t, err)
		assert.Len(t, lines, 3)

		expected := []string{"third line", "fourth line", "fifth line"}
		assert.Equal(t, expected, lines)
	})

	t.Run("Request more lines than available", func(t *testing.T) {
		lines, err := readTailWithBuffer(file, stat.Size(), 10)
		assert.NoError(t, err)
		assert.Len(t, lines, 5) // Should return all available lines
		assert.Equal(t, testLines, lines)
	})
}
