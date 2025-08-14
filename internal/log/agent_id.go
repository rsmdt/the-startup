package log

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ExtractAgentIDFromPrompt extracts agent ID from prompt text using the specified regex pattern.
// Returns empty string if no valid agent ID is found.
//
// The function implements the specification from SDD section 3.2.1:
// - Uses case-insensitive regex pattern: (?i)\bAgentId\s*:\s*([a-zA-Z0-9][a-zA-Z0-9\-_]{0,62}[a-zA-Z0-9])
// - Validates extracted ID against format rules (2-64 chars, alphanumeric with hyphens/underscores)
// - Checks against reserved words: ["main", "global", "system"]
// - Converts result to lowercase as specified
func ExtractAgentIDFromPrompt(prompt string) string {
	// Find AgentId: pattern and capture the complete line after it
	pattern := `(?i)\bAgentId\s*:\s*([^\n\r]*)`
	re := regexp.MustCompile(pattern)
	
	// Find the first match
	matches := re.FindStringSubmatch(prompt)
	if len(matches) < 2 {
		return "" // No match found
	}
	
	// Extract the captured line and trim whitespace
	line := strings.TrimSpace(matches[1])
	if line == "" {
		return "" // Empty or whitespace-only value
	}
	
	// Split by whitespace to get the first token (the agent ID)
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return ""
	}
	
	// If there's more than one token, it means there are additional words which makes it invalid
	if len(tokens) > 1 {
		return "" // Invalid: contains spaces separating multiple tokens
	}
	
	extracted := tokens[0]
	
	// Convert to lowercase as specified
	extracted = strings.ToLower(extracted)
	
	// Validate the extracted ID against the exact specification pattern
	if !isValidAgentID(extracted) {
		return ""
	}
	
	// Check if it's a reserved word
	if isReservedWord(extracted) {
		return ""
	}
	
	return extracted
}

// isValidAgentID validates agent ID format according to specification:
// - Length: 2-64 characters
// - Characters: alphanumeric, hyphens, underscores
// - Must start and end with alphanumeric characters
func isValidAgentID(id string) bool {
	// Check length requirements
	if len(id) < 2 || len(id) > 64 {
		return false
	}
	
	// Must start and end with alphanumeric
	if !isAlphanumeric(id[0]) || !isAlphanumeric(id[len(id)-1]) {
		return false
	}
	
	// Check all characters are valid (alphanumeric, hyphen, or underscore)
	for _, char := range id {
		if !isAlphanumeric(byte(char)) && char != '-' && char != '_' {
			return false
		}
	}
	
	return true
}

// isReservedWord checks if the ID is a reserved word (case-insensitive)
// Reserved words: ["main", "global", "system"]
func isReservedWord(id string) bool {
	reserved := []string{"main", "global", "system"}
	lower := strings.ToLower(id)
	
	for _, word := range reserved {
		if lower == word {
			return true
		}
	}
	
	return false
}

// isAlphanumeric checks if a byte is alphanumeric (a-z, A-Z, 0-9)
func isAlphanumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// GenerateAgentID generates a simple AgentID with format "{shortId}-{agentType}"
// Example: "ppYh6tV7-the-architect"
//
// The function generates:
// - An 8-character alphanumeric shortId
// - Format: "{shortId}-{agentType}"
// - Returns a valid AgentID that is easy to read and type
func GenerateAgentID(sessionID, agentType, promptExcerpt string) string {
	// Handle empty or invalid agent type
	if agentType == "" {
		agentType = "unknown"
	}
	
	// Generate 8-character alphanumeric shortId
	shortId := generateShortId()
	
	// Format: "{shortId}-{agentType}"
	baseID := fmt.Sprintf("%s-%s", shortId, agentType)
	
	return baseID
}

// generateShortId generates an 8-character alphanumeric ID
func generateShortId() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based if random fails
		return fmt.Sprintf("%08x", time.Now().Unix())
	}
	for i := range b {
		b[i] = chars[b[i]%byte(len(chars))]
	}
	return string(b)
}


// ExtractOrGenerateAgentID implements the complete AgentID extraction/generation flow
// as specified in SDD section 3.2.3. First attempts regex extraction, then falls back
// to deterministic generation with uniqueness checking.
func ExtractOrGenerateAgentID(prompt, agentType, sessionID string) string {
	// Step 1: Try regex extraction
	extracted := ExtractAgentIDFromPrompt(prompt)
	if extracted != "" {
		return extracted
	}
	
	// Step 2: Generate deterministic fallback
	baseID := GenerateAgentID(sessionID, agentType, prompt)
	
	// Step 3: Ensure uniqueness (using default base directory ".the-startup")
	return ensureUniqueness(baseID, sessionID, ".the-startup")
}

// truncateString truncates a string to the specified maximum length, handling Unicode correctly.
// Used to ensure prompt excerpts are exactly 100 characters for consistent hash generation.
func truncateString(s string, maxLen int) string {
	// Handle Unicode by converting to runes first
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	
	return string(runes[:maxLen])
}

// ensureUniqueness checks for existing agent context files in the session directory
// and applies disambiguation with "-N" suffix (N=1-999) if needed. Uses ultimate
// fallback with 8-character UUID segment if all numeric suffixes are exhausted.
//
// Parameters:
//   - baseID: The base agent ID to check for uniqueness
//   - sessionID: The session identifier
//   - baseDir: The base directory (typically ".the-startup")
//
// Returns a unique agent ID that doesn't conflict with existing files.
func ensureUniqueness(baseID, sessionID, baseDir string) string {
	// If session directory doesn't exist, baseID is unique
	sessionDir := filepath.Join(baseDir, sessionID)
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		return baseID
	}
	
	// Check if base ID is available
	baseFilePath := filepath.Join(sessionDir, baseID+".jsonl")
	if _, err := os.Stat(baseFilePath); os.IsNotExist(err) {
		return baseID // Base ID is available
	}
	
	// Try numeric suffixes 1-999
	for i := 1; i <= 999; i++ {
		candidateID := fmt.Sprintf("%s-%d", baseID, i)
		candidateFilePath := filepath.Join(sessionDir, candidateID+".jsonl")
		if _, err := os.Stat(candidateFilePath); os.IsNotExist(err) {
			return candidateID // Found available suffix
		}
	}
	
	// Ultimate fallback: append 8-character UUID segment
	uuidSegment := generateRandomHex(8)
	fallbackID := fmt.Sprintf("%s-uuid%s", baseID, uuidSegment)
	
	return fallbackID
}

// generateRandomHex generates a random hex string of the specified length.
// Used for ultimate fallback when all numeric disambiguation suffixes are exhausted.
func generateRandomHex(length int) string {
	bytes := make([]byte, (length+1)/2) // Round up for odd lengths
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based generation if crypto/rand fails
		timestamp := time.Now().UnixNano()
		return fmt.Sprintf("%x", timestamp)[:length]
	}
	
	hexStr := hex.EncodeToString(bytes)
	if len(hexStr) > length {
		hexStr = hexStr[:length]
	}
	
	return hexStr
}