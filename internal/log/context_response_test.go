package log

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadContext(t *testing.T) {
	tempDir := t.TempDir()
	originalClaudeDir := os.Getenv("CLAUDE_PROJECT_DIR")
	defer os.Setenv("CLAUDE_PROJECT_DIR", originalClaudeDir)
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
	
	// Create the .the-startup directory
	startupDir := filepath.Join(tempDir, ".the-startup")
	require.NoError(t, os.MkdirAll(startupDir, 0755))
	
	sessionID := "dev-20250812-143022"
	agentID := "arch-001"

	// Create test data using the working pattern
	sessionDir := filepath.Join(startupDir, sessionID)
	require.NoError(t, os.MkdirAll(sessionDir, 0755))
	
	contextFile := filepath.Join(sessionDir, agentID+".jsonl")
	testData := []string{
		`{"event":"agent_start","agent_type":"the-architect","agent_id":"arch-001","description":"Design authentication system","session_id":"dev-20250812-143022","timestamp":"2025-08-12T14:00:00.000Z"}`,
		`{"event":"agent_complete","agent_type":"the-architect","agent_id":"arch-001","description":"Authentication design completed","session_id":"dev-20250812-143022","timestamp":"2025-08-12T14:05:00.000Z"}`,
		`{"event":"agent_start","agent_type":"the-architect","agent_id":"arch-001","description":"Implement security features","session_id":"dev-20250812-143022","timestamp":"2025-08-12T14:10:00.000Z"}`,
	}
	
	require.NoError(t, os.WriteFile(contextFile, []byte(strings.Join(testData, "\n")+"\n"), 0644))

	testCases := []struct {
		name        string
		query       ContextQuery
		expectError bool
		checkResult func(t *testing.T, result string, err error)
	}{
		{
			name: "Valid JSON format request",
			query: ContextQuery{
				AgentID:         agentID,
				SessionID:       sessionID,
				MaxLines:        10,
				Format:          "json",
				IncludeMetadata: true,
			},
			expectError: false,
			checkResult: func(t *testing.T, result string, err error) {
				assert.NoError(t, err)
				
				var response ContextResponse
				err = json.Unmarshal([]byte(result), &response)
				assert.NoError(t, err)
				
				assert.Equal(t, agentID, response.Metadata.AgentID)
				assert.Equal(t, sessionID, response.Metadata.SessionID)
				assert.Equal(t, 3, response.Metadata.TotalEntries)
				assert.Len(t, response.Entries, 3)
			},
		},
		{
			name: "Valid text format request",
			query: ContextQuery{
				AgentID:   agentID,
				SessionID: sessionID,
				MaxLines:  10,
				Format:    "text",
			},
			expectError: false,
			checkResult: func(t *testing.T, result string, err error) {
				assert.NoError(t, err)
				
				lines := strings.Split(result, "\n")
				assert.Contains(t, lines[0], "Agent: arch-001")
				assert.Contains(t, lines[0], "Session: dev-20250812-143022")
				assert.Contains(t, lines[0], "Entries: 3")
				assert.Equal(t, "---", lines[1])
				
				// Check entries format
				assert.Contains(t, lines[2], "[agent_start]")
				assert.Contains(t, lines[3], "[agent_complete]")
				assert.Contains(t, lines[4], "[agent_start]")
			},
		},
		{
			name: "Request with limited lines",
			query: ContextQuery{
				AgentID:   agentID,
				SessionID: sessionID,
				MaxLines:  2,
				Format:    "json",
			},
			expectError: false,
			checkResult: func(t *testing.T, result string, err error) {
				assert.NoError(t, err)
				
				var response ContextResponse
				err = json.Unmarshal([]byte(result), &response)
				assert.NoError(t, err)
				
				assert.Len(t, response.Entries, 2)
				assert.Equal(t, 2, response.Metadata.TotalEntries)
			},
		},
		{
			name: "Non-existent agent",
			query: ContextQuery{
				AgentID:   "nonexistent",
				SessionID: sessionID,
				MaxLines:  10,
				Format:    "json",
			},
			expectError: false, // Should return empty context, not error
			checkResult: func(t *testing.T, result string, err error) {
				assert.NoError(t, err)
				
				var response ContextResponse
				err = json.Unmarshal([]byte(result), &response)
				assert.NoError(t, err)
				
				assert.Empty(t, response.Entries)
				assert.Equal(t, 0, response.Metadata.TotalEntries)
			},
		},
		{
			name: "Empty agent ID",
			query: ContextQuery{
				AgentID:   "",
				SessionID: sessionID,
				MaxLines:  10,
				Format:    "json",
			},
			expectError: true,
			checkResult: func(t *testing.T, result string, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "agentID is required")
			},
		},
		{
			name: "Default parameters",
			query: ContextQuery{
				AgentID: agentID,
			},
			expectError: false,
			checkResult: func(t *testing.T, result string, err error) {
				assert.NoError(t, err)
				
				var response ContextResponse
				err = json.Unmarshal([]byte(result), &response)
				assert.NoError(t, err)
				
				// Should use defaults: maxLines=50, format=json
				assert.Len(t, response.Entries, 3) // All entries since less than 50
			},
		},
		{
			name: "Invalid format falls back to JSON",
			query: ContextQuery{
				AgentID:   agentID,
				SessionID: sessionID,
				Format:    "xml",
			},
			expectError: false,
			checkResult: func(t *testing.T, result string, err error) {
				assert.NoError(t, err)
				
				// Should be valid JSON despite invalid format request
				var response ContextResponse
				err = json.Unmarshal([]byte(result), &response)
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ReadContext(tc.query)
			tc.checkResult(t, result, err)
		})
	}
}

func TestBuildMetadata(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	sessionID := "dev-20250812-143022"
	agentID := "arch-001"

	// Setup project directory to use temp directory
	oldProjectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	startupDir := filepath.Join(tempDir, ".the-startup")
	require.NoError(t, os.MkdirAll(startupDir, 0755))
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
	defer os.Setenv("CLAUDE_PROJECT_DIR", oldProjectDir)

	// Create test context file
	setupTestContext(t, startupDir, sessionID, agentID)
	
	// Verify the test file was created
	expectedFile := filepath.Join(startupDir, sessionID, agentID+".jsonl")
	require.FileExists(t, expectedFile)

	testCases := []struct {
		name      string
		agentID   string
		sessionID string
		checkFunc func(t *testing.T, metadata ContextMetadata, err error)
	}{
		{
			name:      "Valid metadata build",
			agentID:   agentID,
			sessionID: sessionID,
			checkFunc: func(t *testing.T, metadata ContextMetadata, err error) {
				assert.NoError(t, err)
				assert.Equal(t, agentID, metadata.AgentID)
				assert.Equal(t, sessionID, metadata.SessionID)
				assert.Greater(t, metadata.FileSizeBytes, int64(0))
				assert.False(t, metadata.LastModified.IsZero())
				assert.Contains(t, metadata.ContextFile, agentID+".jsonl")
			},
		},
		{
			name:      "Non-existent file",
			agentID:   "nonexistent",
			sessionID: sessionID,
			checkFunc: func(t *testing.T, metadata ContextMetadata, err error) {
				assert.NoError(t, err) // Should not error, just empty metadata
				assert.Equal(t, "nonexistent", metadata.AgentID)
				assert.Equal(t, sessionID, metadata.SessionID)
				assert.Equal(t, int64(0), metadata.FileSizeBytes)
				assert.True(t, metadata.LastModified.IsZero())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metadata, err := buildMetadata(tc.agentID, tc.sessionID)
			tc.checkFunc(t, metadata, err)
		})
	}
}

func TestFormatTextResponse(t *testing.T) {
	response := ContextResponse{
		Metadata: ContextMetadata{
			AgentID:      "arch-001",
			SessionID:    "dev-20250812-143022",
			TotalEntries: 2,
		},
		Entries: []*HookData{
			{
				Event:       "agent_start",
				AgentID:     "arch-001",
				Description: "Design authentication system",
				Timestamp:   "2025-08-12T14:00:00.000Z",
			},
			{
				Event:       "agent_complete",
				AgentID:     "arch-001",
				Description: "Authentication system completed",
				Timestamp:   "2025-08-12T14:05:00.000Z",
			},
		},
	}

	result := formatTextResponse(response)
	
	lines := strings.Split(result, "\n")
	assert.Contains(t, lines[0], "Agent: arch-001")
	assert.Contains(t, lines[0], "Session: dev-20250812-143022")
	assert.Contains(t, lines[0], "Entries: 2")
	assert.Equal(t, "---", lines[1])
	assert.Contains(t, lines[2], "2025-08-12T14:00:00.000Z [agent_start] Design authentication system")
	assert.Contains(t, lines[3], "2025-08-12T14:05:00.000Z [agent_complete] Authentication system completed")
}

func TestFormatJSONResponse(t *testing.T) {
	response := ContextResponse{
		Metadata: ContextMetadata{
			AgentID:      "arch-001",
			SessionID:    "dev-20250812-143022",
			TotalEntries: 1,
		},
		Entries: []*HookData{
			{
				Event:       "agent_start",
				AgentID:     "arch-001",
				Description: "Test task",
				Timestamp:   "2025-08-12T14:00:00.000Z",
			},
		},
	}

	result, err := formatJSONResponse(response)
	assert.NoError(t, err)
	
	// Verify it's valid JSON
	var parsed ContextResponse
	err = json.Unmarshal([]byte(result), &parsed)
	assert.NoError(t, err)
	
	assert.Equal(t, response.Metadata.AgentID, parsed.Metadata.AgentID)
	assert.Equal(t, response.Metadata.SessionID, parsed.Metadata.SessionID)
	assert.Len(t, parsed.Entries, 1)
}

// setupTestContext creates test context files for testing
func setupTestContext(t *testing.T, startupDir, sessionID, agentID string) {
	// Create session directory
	sessionDir := filepath.Join(startupDir, sessionID)
	require.NoError(t, os.MkdirAll(sessionDir, 0755))

	// Create test context file using the same format as existing tests
	contextFile := filepath.Join(sessionDir, agentID+".jsonl")
	testData := []string{
		`{"event":"agent_start","agent_type":"the-architect","agent_id":"arch-001","description":"Design authentication system","session_id":"dev-20250812-143022","timestamp":"2025-08-12T14:00:00.000Z"}`,
		`{"event":"agent_complete","agent_type":"the-architect","agent_id":"arch-001","description":"Authentication design completed","session_id":"dev-20250812-143022","timestamp":"2025-08-12T14:05:00.000Z"}`,
		`{"event":"agent_start","agent_type":"the-architect","agent_id":"arch-001","description":"Implement security features","session_id":"dev-20250812-143022","timestamp":"2025-08-12T14:10:00.000Z"}`,
	}

	require.NoError(t, os.WriteFile(contextFile, []byte(strings.Join(testData, "\n")+"\n"), 0644))
}