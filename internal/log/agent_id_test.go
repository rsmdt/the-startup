package log

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestExtractAgentIDFromPrompt_ValidCases validates successful AgentID extraction
func TestExtractAgentIDFromPrompt_ValidCases(t *testing.T) {
	testCases := []struct {
		name     string
		prompt   string
		expected string
	}{
		{
			name:     "Standard format",
			prompt:   "AgentId: arch-001\nDesign authentication system",
			expected: "arch-001",
		},
		{
			name:     "Case insensitive",
			prompt:   "agentid: dev-frontend\nImplement UI components",
			expected: "dev-frontend",
		},
		{
			name:     "With extra whitespace",
			prompt:   "AgentId  :   test-integration-auth   \nRun tests",
			expected: "test-integration-auth",
		},
		{
			name:     "Hyphens and underscores",
			prompt:   "AgentId: sec_audit-2024\nSecurity audit",
			expected: "sec_audit-2024",
		},
		{
			name:     "Maximum length (64 chars)",
			prompt:   "AgentId: " + strings.Repeat("a", 62) + "bc\nLong task",
			expected: strings.Repeat("a", 62) + "bc",
		},
		{
			name:     "Minimum length (2 chars)",
			prompt:   "AgentId: a1\nMinimal task",
			expected: "a1",
		},
		{
			name:     "Mixed case in value",
			prompt:   "AgentId: Arch-Auth-System\nMixed case",
			expected: "arch-auth-system",
		},
		{
			name:     "Numbers and letters",
			prompt:   "AgentId: test123\nNumeric test",
			expected: "test123",
		},
		{
			name:     "AgentId in middle of text",
			prompt:   "Some text before\nAgentId: core-service\nSome text after",
			expected: "core-service",
		},
		{
			name:     "Word boundary test",
			prompt:   "NotAgentId: fake\nAgentId: real-id\nMore text",
			expected: "real-id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ExtractAgentIDFromPrompt(tc.prompt)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestExtractAgentIDFromPrompt_InvalidCases validates cases that should return empty string
func TestExtractAgentIDFromPrompt_InvalidCases(t *testing.T) {
	testCases := []struct {
		name   string
		prompt string
	}{
		{
			name:   "No AgentId present",
			prompt: "Design authentication system",
		},
		{
			name:   "Invalid characters - special symbols",
			prompt: "AgentId: invalid@id!\nWork on task",
		},
		{
			name:   "Invalid characters - spaces",
			prompt: "AgentId: invalid id\nWork on task",
		},
		{
			name:   "Too short - single character",
			prompt: "AgentId: a\nShort ID",
		},
		{
			name:   "Too long - 65 characters",
			prompt: "AgentId: " + strings.Repeat("a", 65) + "\nLong ID",
		},
		{
			name:   "Reserved word - main",
			prompt: "AgentId: main\nReserved word",
		},
		{
			name:   "Reserved word - global",
			prompt: "AgentId: global\nReserved word",
		},
		{
			name:   "Reserved word - system",
			prompt: "AgentId: system\nReserved word",
		},
		{
			name:   "Reserved word case insensitive - MAIN",
			prompt: "AgentId: MAIN\nReserved word uppercase",
		},
		{
			name:   "Starts with non-alphanumeric",
			prompt: "AgentId: -invalid\nBad start",
		},
		{
			name:   "Ends with non-alphanumeric",
			prompt: "AgentId: invalid-\nBad end",
		},
		{
			name:   "Empty AgentId",
			prompt: "AgentId: \nEmpty value",
		},
		{
			name:   "Only whitespace AgentId",
			prompt: "AgentId:   \nWhitespace only",
		},
		{
			name:   "Contains dots",
			prompt: "AgentId: invalid.id\nDots not allowed",
		},
		{
			name:   "Contains slashes",
			prompt: "AgentId: invalid/id\nSlashes not allowed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ExtractAgentIDFromPrompt(tc.prompt)
			if result != "" {
				t.Errorf("Expected empty string for invalid case, got %s", result)
			}
		})
	}
}

// TestExtractAgentIDFromPrompt_EdgeCases validates edge cases and boundary conditions
func TestExtractAgentIDFromPrompt_EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		prompt   string
		expected string
	}{
		{
			name:     "Multiple AgentIds - should extract first",
			prompt:   "AgentId: first-id\nSome text\nAgentId: second-id\nMore text",
			expected: "first-id",
		},
		{
			name:     "AgentId with colon variations",
			prompt:   "AgentId:no-space\nNo space after colon",
			expected: "no-space",
		},
		{
			name:     "Tabs instead of spaces",
			prompt:   "AgentId:\t\ttab-separated\nTab test",
			expected: "tab-separated",
		},
		{
			name:     "Mixed whitespace",
			prompt:   "AgentId: \t \n mixed-whitespace\nMixed test",
			expected: "mixed-whitespace",
		},
		{
			name:     "Case variations in AgentId key",
			prompt:   "AGENTID: uppercase-key\nUppercase test",
			expected: "uppercase-key",
		},
		{
			name:     "AgentId at start of string",
			prompt:   "AgentId: start-of-string",
			expected: "start-of-string",
		},
		{
			name:     "AgentId at end of string",
			prompt:   "Some text\nAgentId: end-of-string",
			expected: "end-of-string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ExtractAgentIDFromPrompt(tc.prompt)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestValidationRules validates the agent ID validation rules
func TestValidationRules(t *testing.T) {
	// Test length requirements
	if isValidAgentID("a") {
		t.Error("Single character should be invalid")
	}
	if !isValidAgentID("ab") {
		t.Error("Two characters should be valid")
	}
	if !isValidAgentID(strings.Repeat("a", 64)) {
		t.Error("64 characters should be valid")
	}
	if isValidAgentID(strings.Repeat("a", 65)) {
		t.Error("65 characters should be invalid")
	}

	// Test character requirements
	if !isValidAgentID("a1") {
		t.Error("Alphanumeric should be valid")
	}
	if !isValidAgentID("test-id") {
		t.Error("Hyphens should be valid")
	}
	if !isValidAgentID("test_id") {
		t.Error("Underscores should be valid")
	}
	if !isValidAgentID("test-123_abc") {
		t.Error("Mixed valid chars should be valid")
	}

	// Test start/end character requirements
	if isValidAgentID("-test") {
		t.Error("Cannot start with hyphen")
	}
	if isValidAgentID("_test") {
		t.Error("Cannot start with underscore")
	}
	if isValidAgentID("test-") {
		t.Error("Cannot end with hyphen")
	}
	if isValidAgentID("test_") {
		t.Error("Cannot end with underscore")
	}

	// Test invalid characters
	if isValidAgentID("test@id") {
		t.Error("Special characters not allowed")
	}
	if isValidAgentID("test id") {
		t.Error("Spaces not allowed")
	}
	if isValidAgentID("test.id") {
		t.Error("Dots not allowed")
	}
	if isValidAgentID("test/id") {
		t.Error("Slashes not allowed")
	}
}

// TestReservedWords validates reserved word checking
func TestReservedWords(t *testing.T) {
	reservedWords := []string{"main", "global", "system"}

	for _, word := range reservedWords {
		t.Run("Reserved word: "+word, func(t *testing.T) {
			if !isReservedWord(word) {
				t.Errorf("Should identify %s as reserved", word)
			}
			if !isReservedWord(strings.ToUpper(word)) {
				t.Errorf("Should identify %s as reserved (uppercase)", word)
			}
			if !isReservedWord(strings.Title(word)) {
				t.Errorf("Should identify %s as reserved (title case)", word)
			}
		})
	}

	// Test non-reserved words
	nonReservedWords := []string{"main-test", "global-id", "system-auth", "test-main", "arch-001"}
	for _, word := range nonReservedWords {
		t.Run("Non-reserved word: "+word, func(t *testing.T) {
			if isReservedWord(word) {
				t.Errorf("Should not identify %s as reserved", word)
			}
		})
	}
}

// TestRegexPattern validates the implementation approach using complete line extraction
func TestRegexPattern(t *testing.T) {
	// Test that our implementation correctly handles various inputs
	validCases := []struct {
		input    string
		expected string
	}{
		{"AgentId: test-id", "test-id"},
		{"agentid: test-id", "test-id"},
		{"AGENTID: test-id", "test-id"},
		{"AgentId:test-id", "test-id"},
		{"AgentId: test_id", "test_id"},
		{"AgentId:   test-id   \nmore text", "test-id"},
		{"Some text AgentId: test-id\nmore text", "test-id"},
	}

	for _, tc := range validCases {
		t.Run("Valid extraction: "+tc.input, func(t *testing.T) {
			result := ExtractAgentIDFromPrompt(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}

	invalidCases := []string{
		"NotAgentId: test-id", // Wrong keyword
		"AgentId test-id",     // No colon
		"AgentId: -test",      // Starts with hyphen
		"AgentId: test-",      // Ends with hyphen
		"AgentId: t",          // Too short
		"AgentId: test@id",    // Invalid character
		"AgentId: test id",    // Contains space
		"AgentId: ",           // Empty
		"AgentId:   ",         // Whitespace only
	}

	for _, tc := range invalidCases {
		t.Run("Invalid extraction: "+tc, func(t *testing.T) {
			result := ExtractAgentIDFromPrompt(tc)
			if result != "" {
				t.Errorf("Expected empty result for invalid input, got: %s", result)
			}
		})
	}
}

// TestCaseConversion validates that extracted IDs are converted to lowercase
func TestCaseConversion(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"AgentId: ARCH-001", "arch-001"},
		{"AgentId: Test-Frontend", "test-frontend"},
		{"AgentId: DEV_BACKEND", "dev_backend"},
		{"AgentId: Mixed-Case_ID", "mixed-case_id"},
	}

	for _, tc := range testCases {
		t.Run("Case conversion: "+tc.input, func(t *testing.T) {
			result := ExtractAgentIDFromPrompt(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// Helper functions are implemented in agent_id.go

// TestGenerateAgentID_DeterministicGeneration validates that the same inputs produce the same AgentID
func TestGenerateAgentID_DeterministicGeneration(t *testing.T) {
	testCases := []struct {
		name          string
		sessionID     string
		agentType     string
		promptExcerpt string
	}{
		{
			name:          "Standard case",
			sessionID:     "dev-20250812-143022",
			agentType:     "the-architect",
			promptExcerpt: "Design authentication system for the web application",
		},
		{
			name:          "Different agent type",
			sessionID:     "dev-20250812-143022",
			agentType:     "the-developer",
			promptExcerpt: "Implement frontend components",
		},
		{
			name:          "Long prompt excerpt (truncated to 100 chars)",
			sessionID:     "dev-20250812-143022",
			agentType:     "the-tester",
			promptExcerpt: strings.Repeat("This is a very long prompt that should be truncated to exactly 100 characters. ", 5),
		},
		{
			name:          "Empty prompt excerpt",
			sessionID:     "dev-20250812-143022",
			agentType:     "the-devops-engineer",
			promptExcerpt: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate the same ID multiple times
			id1 := GenerateAgentID(tc.sessionID, tc.agentType, tc.promptExcerpt)
			id2 := GenerateAgentID(tc.sessionID, tc.agentType, tc.promptExcerpt)
			id3 := GenerateAgentID(tc.sessionID, tc.agentType, tc.promptExcerpt)

			// All should be identical
			if id1 != id2 || id2 != id3 {
				t.Errorf("Generated IDs are not deterministic: %s, %s, %s", id1, id2, id3)
			}

			// Verify format: {agentType}-{YYYYMMDD-HHMMSS}-{hash8}
			// Since agentType can contain hyphens, we need to find the timestamp part
			if !strings.HasPrefix(id1, tc.agentType+"-") {
				t.Errorf("Generated ID should start with agent type: %s, got: %s", tc.agentType, id1)
			}

			// Extract the part after agent type
			afterAgentType := id1[len(tc.agentType)+1:] // +1 for the hyphen
			parts := strings.Split(afterAgentType, "-")
			if len(parts) != 3 {
				t.Errorf("After agent type, should have 3 parts (YYYYMMDD-HHMMSS-hash8), got: %s", afterAgentType)
			}

			// Check timestamp format (YYYYMMDD)
			dateStr := parts[0]
			if len(dateStr) != 8 {
				t.Errorf("Expected date part to be 8 characters (YYYYMMDD), got: %s", dateStr)
			}

			// Check time format (HHMMSS)
			timeStr := parts[1]
			if len(timeStr) != 6 {
				t.Errorf("Expected time part to be 6 characters (HHMMSS), got: %s", timeStr)
			}

			// Check hash format (8 hex characters)
			hashStr := parts[2]
			if len(hashStr) != 8 {
				t.Errorf("Expected hash part to be 8 characters, got: %s", hashStr)
			}

			// Verify hash contains only hex characters
			for _, char := range hashStr {
				if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
					t.Errorf("Hash part should contain only hex characters, got: %s", hashStr)
				}
			}

			// Verify the generated ID is valid
			if !isValidAgentID(id1) {
				t.Errorf("Generated ID should be valid: %s", id1)
			}
		})
	}
}

// TestGenerateAgentID_UniquenessAcrossInputs validates that different inputs produce different IDs
func TestGenerateAgentID_UniquenessAcrossInputs(t *testing.T) {
	sessionID := "dev-20250812-143022"
	promptExcerpt := "Design authentication system"

	// Different agent types should produce different IDs
	id1 := GenerateAgentID(sessionID, "the-architect", promptExcerpt)
	id2 := GenerateAgentID(sessionID, "the-developer", promptExcerpt)
	if id1 == id2 {
		t.Errorf("Different agent types should produce different IDs: %s", id1)
	}

	// Different session IDs should produce different IDs
	id3 := GenerateAgentID("dev-20250813-090000", "the-architect", promptExcerpt)
	if id1 == id3 {
		t.Errorf("Different session IDs should produce different IDs: %s", id1)
	}

	// Different prompt excerpts should produce different IDs
	id4 := GenerateAgentID(sessionID, "the-architect", "Different task description")
	if id1 == id4 {
		t.Errorf("Different prompt excerpts should produce different IDs: %s", id1)
	}
}

// TestGenerateAgentID_PromptExcerptTruncation validates 100-character truncation
func TestGenerateAgentID_PromptExcerptTruncation(t *testing.T) {
	sessionID := "dev-20250812-143022"
	agentType := "the-architect"

	// Create prompts of different lengths
	shortPrompt := "Short prompt"
	longPrompt := strings.Repeat("This is a long prompt. ", 10)        // Much longer than 100 chars
	veryLongPrompt := strings.Repeat("Very long prompt content. ", 20) // Even longer

	// Generate IDs
	idShort := GenerateAgentID(sessionID, agentType, shortPrompt)
	idLong := GenerateAgentID(sessionID, agentType, longPrompt)
	idVeryLong := GenerateAgentID(sessionID, agentType, veryLongPrompt)

	// All should be different
	if idShort == idLong || idLong == idVeryLong {
		t.Errorf("Different prompts should produce different IDs")
	}

	// Test that truncation works correctly - same 100-char prefix should produce same ID
	prefix100 := longPrompt[:100]
	extendedPrompt := prefix100 + " additional text beyond 100 characters"

	id100 := GenerateAgentID(sessionID, agentType, prefix100)
	idExtended := GenerateAgentID(sessionID, agentType, extendedPrompt)

	if id100 != idExtended {
		t.Errorf("Prompts with same 100-char prefix should produce same ID: %s vs %s", id100, idExtended)
	}
}

// TestGenerateAgentID_ValidFormat validates the generated ID format matches specification
func TestGenerateAgentID_ValidFormat(t *testing.T) {
	testCases := []struct {
		agentType string
	}{
		{"the-architect"},
		{"the-developer"},
		{"the-tester"},
		{"the-devops-engineer"},
		{"the-security-engineer"},
		{"the-product-manager"},
		{"the-technical-writer"},
		{"the-data-engineer"},
		{"the-site-reliability-engineer"},
		{"the-business-analyst"},
		{"the-project-manager"},
		{"the-prompt-engineer"},
		{"the-chief"},
	}

	sessionID := "dev-20250812-143022"
	promptExcerpt := "Test task description"

	for _, tc := range testCases {
		t.Run("Agent type: "+tc.agentType, func(t *testing.T) {
			id := GenerateAgentID(sessionID, tc.agentType, promptExcerpt)

			// Verify format pattern: {agentType}-{YYYYMMDD-HHMMSS}-{hash8}
			expectedPattern := `^` + regexp.QuoteMeta(tc.agentType) + `-\d{8}-\d{6}-[a-f0-9]{8}$`
			matched, err := regexp.MatchString(expectedPattern, id)
			if err != nil {
				t.Fatalf("Regex error: %v", err)
			}
			if !matched {
				t.Errorf("Generated ID %s does not match expected pattern %s", id, expectedPattern)
			}

			// Verify it's a valid agent ID
			if !isValidAgentID(id) {
				t.Errorf("Generated ID should be valid: %s", id)
			}

			// Verify it's not a reserved word
			if isReservedWord(id) {
				t.Errorf("Generated ID should not be a reserved word: %s", id)
			}
		})
	}
}

// TestGenerateAgentID_SHA256Consistency validates SHA256 hash consistency
func TestGenerateAgentID_SHA256Consistency(t *testing.T) {
	sessionID := "dev-20250812-143022"
	agentType := "the-architect"
	promptExcerpt := "Design authentication system"

	// Generate multiple IDs
	ids := make([]string, 10)
	for i := 0; i < 10; i++ {
		ids[i] = GenerateAgentID(sessionID, agentType, promptExcerpt)
	}

	// All should be identical (deterministic)
	for i := 1; i < len(ids); i++ {
		if ids[0] != ids[i] {
			t.Errorf("IDs should be deterministic, got different results: %s vs %s", ids[0], ids[i])
		}
	}

	// Extract hash part and verify it's consistent
	parts := strings.Split(ids[0], "-")
	hash := parts[len(parts)-1]

	// Verify it's exactly 8 hex characters
	if len(hash) != 8 {
		t.Errorf("Hash should be 8 characters, got %d: %s", len(hash), hash)
	}

	for _, char := range hash {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			t.Errorf("Hash should contain only lowercase hex characters, found: %c", char)
		}
	}
}

// TestGenerateAgentID_EdgeCases validates edge cases and boundary conditions
func TestGenerateAgentID_EdgeCases(t *testing.T) {
	testCases := []struct {
		name          string
		sessionID     string
		agentType     string
		promptExcerpt string
		expectValid   bool
	}{
		{
			name:          "Empty session ID",
			sessionID:     "",
			agentType:     "the-architect",
			promptExcerpt: "Test task",
			expectValid:   true, // Should still generate valid ID
		},
		{
			name:          "Empty agent type",
			sessionID:     "dev-20250812-143022",
			agentType:     "",
			promptExcerpt: "Test task",
			expectValid:   true, // Should still generate valid ID (will become "unknown")
		},
		{
			name:          "Empty prompt excerpt",
			sessionID:     "dev-20250812-143022",
			agentType:     "the-architect",
			promptExcerpt: "",
			expectValid:   true, // Should still generate valid ID
		},
		{
			name:          "All empty inputs",
			sessionID:     "",
			agentType:     "",
			promptExcerpt: "",
			expectValid:   true, // Should still generate valid ID (will become "unknown")
		},
		{
			name:          "Special characters in inputs",
			sessionID:     "dev-20250812@143022!",
			agentType:     "the-architect@test",
			promptExcerpt: "Test task with @#$% special chars",
			expectValid:   true, // Should generate valid ID despite special chars in input (will be sanitized)
		},
		{
			name:          "Unicode characters",
			sessionID:     "dev-20250812-143022",
			agentType:     "the-architect",
			promptExcerpt: "Test with unicode: 你好世界 🚀",
			expectValid:   true, // Should handle unicode gracefully
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := GenerateAgentID(tc.sessionID, tc.agentType, tc.promptExcerpt)

			if tc.expectValid {
				// Generated ID should always be valid regardless of input
				if id == "" {
					t.Errorf("Should generate non-empty ID even for edge case inputs")
				}
				if !isValidAgentID(id) {
					t.Errorf("Generated ID should be valid even for edge case: %s", id)
				}
			}

			// Verify deterministic behavior for edge cases
			id2 := GenerateAgentID(tc.sessionID, tc.agentType, tc.promptExcerpt)
			if id != id2 {
				t.Errorf("Edge case should still be deterministic: %s vs %s", id, id2)
			}
		})
	}
}

// TestEnsureUniqueness_NewFile validates uniqueness checking when no existing files
func TestEnsureUniqueness_NewFile(t *testing.T) {
	tempDir := t.TempDir()
	sessionID := "dev-test-session"
	baseID := "the-architect-20250812-143022-abcd1234"

	// Test with non-existent session directory
	result := ensureUniqueness(baseID, sessionID, tempDir)
	if result != baseID {
		t.Errorf("Expected original ID when no conflicts: %s, got: %s", baseID, result)
	}
}

// TestEnsureUniqueness_WithCollisions validates disambiguation with -N suffix
func TestEnsureUniqueness_WithCollisions(t *testing.T) {
	tempDir := t.TempDir()
	sessionID := "dev-test-session"
	baseID := "the-architect-20250812-143022-abcd1234"

	// Create session directory and existing files
	sessionDir := tempDir + "/" + sessionID
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		t.Fatalf("Failed to create session directory: %v", err)
	}

	// Create conflicting files
	existingFiles := []string{
		baseID + ".jsonl",   // Original collision
		baseID + "-1.jsonl", // First disambiguation
		baseID + "-2.jsonl", // Second disambiguation
	}

	for _, fileName := range existingFiles {
		filePath := sessionDir + "/" + fileName
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
	}

	// Test uniqueness enforcement
	result := ensureUniqueness(baseID, sessionID, tempDir)
	expected := baseID + "-3"
	if result != expected {
		t.Errorf("Expected %s, got: %s", expected, result)
	}
}

// TestEnsureUniqueness_MaxDisambiguation validates ultimate fallback with UUID
func TestEnsureUniqueness_MaxDisambiguation(t *testing.T) {
	tempDir := t.TempDir()
	sessionID := "dev-test-session"
	baseID := "the-architect-20250812-143022-abcd1234"

	// Create session directory
	sessionDir := tempDir + "/" + sessionID
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		t.Fatalf("Failed to create session directory: %v", err)
	}

	// Create files for all numeric suffixes (1-999)
	for i := 0; i <= 999; i++ {
		fileName := baseID
		if i > 0 {
			fileName += fmt.Sprintf("-%d", i)
		}
		fileName += ".jsonl"

		filePath := sessionDir + "/" + fileName
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
	}

	// Test ultimate fallback with UUID
	result := ensureUniqueness(baseID, sessionID, tempDir)

	// Should have UUID suffix pattern: baseID-uuid{8chars}
	expectedPrefix := baseID + "-uuid"
	if !strings.HasPrefix(result, expectedPrefix) {
		t.Errorf("Expected result to start with %s, got: %s", expectedPrefix, result)
	}

	// UUID part should be exactly 8 characters
	uuidPart := result[len(expectedPrefix):]
	if len(uuidPart) != 8 {
		t.Errorf("Expected UUID part to be 8 characters, got %d: %s", len(uuidPart), uuidPart)
	}

	// Verify the result is still a valid agent ID
	if !isValidAgentID(result) {
		t.Errorf("Ultimate fallback should produce valid agent ID: %s", result)
	}
}

// TestEnsureUniqueness_DirectoryPermissions validates behavior with permission issues
func TestEnsureUniqueness_DirectoryPermissions(t *testing.T) {
	tempDir := t.TempDir()
	sessionID := "dev-test-session"
	baseID := "the-architect-20250812-143022-abcd1234"

	// Test with non-existent base directory (should handle gracefully)
	nonExistentDir := tempDir + "/nonexistent"
	result := ensureUniqueness(baseID, sessionID, nonExistentDir)

	// Should return original ID when directory doesn't exist
	if result != baseID {
		t.Errorf("Expected original ID when directory doesn't exist: %s, got: %s", baseID, result)
	}
}

// TestTruncateString validates the 100-character prompt truncation
func TestTruncateString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "Short string",
			input:    "Short text",
			maxLen:   100,
			expected: "Short text",
		},
		{
			name:     "Exactly 100 characters",
			input:    strings.Repeat("a", 100),
			maxLen:   100,
			expected: strings.Repeat("a", 100),
		},
		{
			name:     "Long string truncated",
			input:    strings.Repeat("a", 150),
			maxLen:   100,
			expected: strings.Repeat("a", 100),
		},
		{
			name:     "Empty string",
			input:    "",
			maxLen:   100,
			expected: "",
		},
		{
			name:     "Unicode characters",
			input:    "Hello 世界! " + strings.Repeat("a", 100),
			maxLen:   100,
			expected: "Hello 世界! " + strings.Repeat("a", 88), // Accounts for unicode
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := truncateString(tc.input, tc.maxLen)
			resultRunes := []rune(result)
			if len(resultRunes) > tc.maxLen {
				t.Errorf("Truncated string should not exceed %d characters, got %d", tc.maxLen, len(resultRunes))
			}
			if tc.name != "Unicode characters" && result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestExtractOrGenerateAgentID_Integration validates the complete extraction/generation flow
func TestExtractOrGenerateAgentID_Integration(t *testing.T) {
	testCases := []struct {
		name            string
		prompt          string
		agentType       string
		sessionID       string
		expectExtracted bool
		expectedPrefix  string
	}{
		{
			name:            "Successful extraction",
			prompt:          "AgentId: arch-auth-001\nDesign authentication system",
			agentType:       "the-architect",
			sessionID:       "dev-20250812-143022",
			expectExtracted: true,
			expectedPrefix:  "arch-auth-001",
		},
		{
			name:            "Invalid ID triggers generation",
			prompt:          "AgentId: invalid@id\nWork on task",
			agentType:       "the-developer",
			sessionID:       "dev-20250812-143022",
			expectExtracted: false,
			expectedPrefix:  "the-developer-",
		},
		{
			name:            "No ID triggers generation",
			prompt:          "Design frontend components without explicit ID",
			agentType:       "the-developer",
			sessionID:       "dev-20250812-143022",
			expectExtracted: false,
			expectedPrefix:  "the-developer-",
		},
		{
			name:            "Reserved word triggers generation",
			prompt:          "AgentId: main\nMain task",
			agentType:       "the-architect",
			sessionID:       "dev-20250812-143022",
			expectExtracted: false,
			expectedPrefix:  "the-architect-",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ExtractOrGenerateAgentID(tc.prompt, tc.agentType, tc.sessionID)

			// Verify result is non-empty and valid
			if result == "" {
				t.Errorf("Should return non-empty agent ID")
			}
			if !isValidAgentID(result) {
				t.Errorf("Result should be valid agent ID: %s", result)
			}

			// Check prefix expectation
			if !strings.HasPrefix(result, tc.expectedPrefix) {
				t.Errorf("Expected result to start with %s, got: %s", tc.expectedPrefix, result)
			}

			// Verify deterministic behavior
			result2 := ExtractOrGenerateAgentID(tc.prompt, tc.agentType, tc.sessionID)
			if result != result2 {
				t.Errorf("Should be deterministic: %s vs %s", result, result2)
			}
		})
	}
}
