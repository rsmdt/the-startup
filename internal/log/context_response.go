package log

import (
	"strings"
)

// ContextQuery represents a simplified context read request
type ContextQuery struct {
	AgentID         string // Required: agent identifier
	SessionID       string // Optional: specific session (empty = latest)
	MaxLines        int    // Number of entries to return (1-1000)
	Format          string // Kept for compatibility but ignored - always returns raw JSONL
	IncludeMetadata bool   // Kept for compatibility but ignored
}

// ReadContext reads agent context and returns raw JSONL lines
func ReadContext(query ContextQuery) (string, error) {
	// Set defaults
	if query.MaxLines == 0 {
		query.MaxLines = 50
	}

	// Validate query
	if query.AgentID == "" {
		return "", nil // Silent error handling for hook compatibility
	}
	if query.MaxLines < 1 || query.MaxLines > 1000 {
		query.MaxLines = 50 // Use default if invalid
	}

	// Read raw JSONL lines from the agent's context file
	lines, err := ReadAgentContextRaw(query.SessionID, query.AgentID, query.MaxLines)
	if err != nil {
		// Silent error handling - return empty
		return "", nil
	}

	// Join the lines with newlines to return raw JSONL
	return strings.Join(lines, "\n"), nil
}
