# Implementation Plan: Context Management System

## Executive Summary

This plan implements the Context Management System based on complete technical specifications defined in SDD.md. All analysis and design work is complete - this plan contains only implementation tasks that directly execute the already-designed solution.

**Key Implementation Areas**:
- **AgentID Management**: Implement regex extraction and fallback generation as specified in SDD section 3.2
- **Context File Operations**: Implement JSONL writing/reading as specified in SDD sections 3.3.4 and 3.5
- **CLI Enhancement**: Implement --read flag as specified in SDD section 3.3
- **Agent Template Updates**: Apply context loading patterns as specified in SDD section 3.7
- **Testing Implementation**: Execute test specifications defined in SDD section 6

**Estimated Effort**: Medium - Direct implementation of already-designed components with clear specifications.

## Phase 1: AgentID Management System Implementation

*Execute agent identification components as specified in SDD sections 3.2.1-3.2.3*

- [x] Implement AgentID extraction function using regex pattern from SDD section 3.2.1 {agent: `the-developer`}
  - [x] Create `ExtractAgentIDFromPrompt()` with pattern `(?i)\bAgentId\s*:\s*([a-zA-Z0-9][a-zA-Z0-9\-_]{0,62}[a-zA-Z0-9])`
  - [x] Implement validation rules: 2-64 chars, alphanumeric with hyphens/underscores
  - [x] Add reserved word checking for ["main", "global", "system"]
  - [x] Add case conversion to lowercase as specified

- [x] Implement fallback generation algorithm as specified in SDD section 3.2.2 {agent: `the-developer`}
  - [x] Create `GenerateAgentID()` using SHA256 hash of "{sessionID}:{agentType}:{promptExcerpt100}"
  - [x] Format output as "{agentType}-{YYYYMMDD-HHMMSS}-{hash8}"
  - [x] Implement uniqueness checking in session directory
  - [x] Add disambiguation with "-N" suffix (N=1-999)
  - [x] Add ultimate fallback with 8-character UUID segment

- [x] Implement complete extraction function as specified in SDD section 3.2.3 {agent: `the-developer`}
  - [x] Create `ExtractOrGenerateAgentID()` combining extraction and fallback
  - [x] Add truncateString() helper for 100-character prompt excerpts
  - [x] Add ensureUniqueness() helper for session-scoped uniqueness
  - [x] Implement error handling as specified

- [x] **Validation**: Test AgentID management using specifications from SDD section 6.1.1
  - [x] **Command**: `go test ./internal/log/... -run TestExtractAgentID -v`
  - [x] **Command**: `go test ./internal/log/... -run TestDeterministicGeneration -v`

## Phase 2: Context File Operations Implementation

*Implement file management as specified in SDD sections 3.5 and 3.6*

- [x] Implement AgentFileManager structure as specified in SDD section 3.6.1 {agent: `the-developer`}
  - [x] Create struct with SessionID, AgentID, BaseDir, mutex fields
  - [x] Implement GetAgentFilePath() method returning ".the-startup/{sessionID}/{agentID}.jsonl"
  - [x] Implement GetSessionDir() method
  - [x] Implement EnsureSessionDir() with 0755 permissions

- [x] Implement WriteAgentContext function as specified in SDD section 3.3.4 {agent: `the-developer`}
  - [x] Create function with signature `WriteAgentContext(sessionID, agentID string, data *HookData) error`
  - [x] Add parameter validation for non-empty sessionID, agentID, and data
  - [x] Implement directory creation with os.MkdirAll(filepath.Dir(agentFile), 0755)
  - [x] Add atomic write with fileMutex.Lock() protection
  - [x] Use existing appendJSONL() function for writing

- [x] Implement ReadAgentContext function as specified in SDD section 3.3.4 {agent: `the-developer`}
  - [x] Create function with signature `ReadAgentContext(sessionID, agentID string, maxLines int) ([]*HookData, error)`
  - [x] Add parameter validation (agentID required, maxLines 1-1000 with default 50)
  - [x] Implement session resolution using findLatestSessionWithAgent() when sessionID empty
  - [x] Add file existence check returning empty slice (not error) for missing files
  - [x] Use readLastNLines() for efficient tail reading

- [x] Implement file reading helpers as specified in SDD section 3.3.4 {agent: `the-developer`}
  - [x] Create findLatestSessionWithAgent() scanning .the-startup/ for latest session with agent file
  - [x] Create readLastNLines() with 1MB threshold for small vs large file handling
  - [x] Create readTailWithBuffer() for efficient reverse reading of large files (8KB buffer)
  - [x] Add JSONL parsing with error recovery (skip corrupted lines)

- [x] **Validation**: Test file operations using specifications from SDD section 6.1.2
  - [x] **Command**: `go test ./internal/log/... -run TestWriteAgentContext -v`
  - [x] **Command**: `go test ./internal/log/... -run TestReadAgentContext -v`

## Phase 3: CLI Enhancement Implementation

*Implement --read flag as specified in SDD section 3.3*

- [x] Implement LogFlags structure as specified in SDD section 3.6.1 {agent: `the-developer`}
  - [x] Add Read, AgentID, Lines, Session, Format, IncludeMetadata fields
  - [x] Maintain existing Assistant and User fields

- [x] Implement command validation as specified in SDD section 3.3.2 {agent: `the-developer`}
  - [x] Create ValidateLogCommand() enforcing exactly one mode (assistant, user, or read)
  - [x] Add --agent-id requirement when --read is set
  - [x] Add AgentID format validation using isValidAgentID()
  - [x] Add --lines range validation (1-1000)

- [x] Extend log command with --read capability as specified in SDD section 3.3.1 {agent: `the-developer`}
  - [x] Add --read, -r flag (bool, conflicts with --assistant/--user)
  - [x] Add --agent-id flag (string, required when --read set)
  - [x] Add --lines flag (int, default 50, max 1000)
  - [x] Add --session flag (string, optional)
  - [x] Add --format flag (string, default "json", allowed ["json", "text"])
  - [x] Add --include-metadata flag (bool, default false)

- [x] Implement context reading execution path as specified in SDD section 3.3.3 {agent: `the-developer`}
  - [x] Create ContextQuery and ContextResponse structures from SDD section 3.6.1
  - [x] Implement JSON output format with metadata and entries as specified
  - [x] Implement text output format as specified
  - [x] Add silent error handling (exit 0) for hook compatibility

- [x] **Validation**: Test CLI enhancement using integration tests from SDD section 6.2.1
  - [x] **Command**: `echo '{"tool_name":"Task","tool_input":{"subagent_type":"the-architect","prompt":"AgentId: arch-001"}}' | go run . log --assistant`
  - [x] **Command**: `go run . log --read --agent-id arch-001 --lines 10`

## Phase 4: Hook Processing Integration

*Integrate AgentID management with existing hook processing as specified in SDD*

- [x] Enhance ProcessToolCall function to use AgentID extraction as specified {agent: `the-developer`}
  - [x] Replace existing AgentID extraction with ExtractOrGenerateAgentID() from Phase 1
  - [x] Pass agentType (subagent_type), prompt, and sessionID to extraction function
  - [x] Update HookData structure population with extracted/generated AgentID

- [x] Implement agent-specific context writing as specified in SDD sections 3.5.2-3.5.3 {agent: `the-developer`}
  - [x] Call WriteAgentContext() for agent-specific files (.the-startup/{sessionId}/{agentId}.jsonl)
  - [x] Maintain existing writeSessionLog() for backward compatibility (agent-instructions.jsonl)
  - [x] Maintain existing global logging (all-agent-instructions.jsonl)
  - [x] Handle both PreToolUse and PostToolUse events

- [x] Implement orchestrator context routing as specified in SDD section 3.5.2 {agent: `the-developer`}
  - [x] Route main agent interactions (no subagent_type) to main.jsonl
  - [x] Route subagent interactions to individual agent files
  - [x] Ensure all interactions continue to legacy files for compatibility

- [x] **Validation**: Test complete hook integration using specifications from SDD section 6.2.1
  - [x] **Command**: Run complete integration test script from SDD section 6.2.1
  - [x] **Command**: `go test ./internal/log/... -v`

## Phase 5: Agent Template Updates

*Apply context loading instructions as specified in SDD section 3.7*

- [x] Update the-architect.md template as specified in SDD section 3.7.3 {agent: `the-developer`}
  - [x] Insert Context Discovery section near beginning
  - [x] Add context loading command: `CONTEXT=$(the-startup log --read --agent-id arch-{component} --lines 20 --format json 2>/dev/null || echo '{"entries":[]}')`
  - [x] Add conditional context integration with HAS_CONTEXT check
  - [x] Add context summary display and integration notes

- [x] Update the-developer.md template as specified in SDD section 3.7.3 {agent: `the-developer`}
  - [x] Insert Context Discovery section with dev-{feature-area} pattern
  - [x] Apply same context loading and integration pattern
  - [x] Ensure template handles both contextual and fresh development tasks

- [x] Update the-tester.md template as specified in SDD section 3.7.3 {agent: `the-developer`}
  - [x] Insert Context Discovery section with test-{test-type} pattern
  - [x] Apply context loading pattern for test continuity

- [x] Update remaining agent templates as specified in SDD sections 3.7.2-3.7.3 {agent: `the-developer`}
  - [x] Update the-devops-engineer.md with devops-{infrastructure-component} pattern
  - [x] Update the-security-engineer.md with sec-{security-domain} pattern
  - [x] Update the-product-manager.md with pm-{product-area} pattern
  - [x] Update the-technical-writer.md with doc-{documentation-type} pattern
  - [x] Update the-data-engineer.md with data-{data-component} pattern
  - [x] Update the-site-reliability-engineer.md with sre-{reliability-area} pattern
  - [x] Update the-business-analyst.md with ba-{business-domain} pattern
  - [x] Update the-project-manager.md with proj-{project-phase} pattern
  - [x] Update the-prompt-engineer.md with prompt-{prompt-domain} pattern
  - [x] Update the-chief.md with chief-{decision-area} pattern

- [x] **Validation**: Test template context integration using performance tests from SDD section 6.3
  - [x] **Command**: Create test context and measure loading performance (<2 seconds)
  - [x] **Command**: Verify context format compatibility with jq parsing

## Phase 6: Testing Implementation

*Execute test specifications defined in SDD section 6*

- [x] Implement unit tests as specified in SDD section 6.1 {agent: `the-developer`}
  - [x] Create TestExtractAgentID_ValidCases with test cases from SDD section 6.1.1
  - [x] Create TestExtractAgentID_InvalidCases validating fallback generation
  - [x] Create TestDeterministicGeneration validating consistent ID generation
  - [x] Create TestWriteAgentContext validating file operations
  - [x] Create TestReadAgentContext with edge cases (corruption, missing files)

- [x] Implement integration tests as specified in SDD section 6.2 {agent: `the-developer`}
  - [x] Create integration_test_hook_processing.sh from SDD section 6.2.1
  - [x] Test complete PreToolUse/PostToolUse workflow
  - [x] Verify agent-specific and backward compatibility file creation
  - [x] Test context reading functionality

- [x] Implement performance tests as specified in SDD section 6.3 {agent: `the-developer`}
  - [x] Create performance_test_large_files.sh from SDD section 6.3.1
  - [x] Test 50MB file handling with <2 second read requirement
  - [x] Test memory usage <100MB requirement
  - [x] Create concurrent_test.sh from SDD section 6.3.2
  - [x] Test parallel agent execution safety

- [x] Execute test coverage validation as specified in SDD section 6.4 {agent: `the-developer`}
  - [x] Achieve 85% minimum unit test coverage
  - [x] Achieve 100% coverage for ExtractOrGenerateAgentID
  - [x] Achieve 95% coverage for WriteAgentContext and ReadAgentContext
  - [x] Verify all edge cases covered (missing context, corruption, permissions)

- [x] **Final Validation**: Execute complete test suite
  - [x] **Command**: `go test -v ./...`
  - [x] **Command**: `go test -cover ./...`
  - [x] **Command**: `bash integration_test_hook_processing.sh`
  - [x] **Command**: `bash performance_test_large_files.sh`
  - [x] **Command**: `bash concurrent_test.sh`

## Success Criteria

### Functional Requirements (from SDD specifications)
- [x] **AgentID Management**: Regex extraction with fallback generation implemented per SDD section 3.2
- [x] **Context Writing**: Individual agent JSONL files created per SDD section 3.5
- [x] **Context Reading**: CLI --read flag implemented per SDD section 3.3
- [x] **Agent Integration**: All 13 templates updated with context loading per SDD section 3.7
- [x] **Backward Compatibility**: Existing agent-instructions.jsonl and global logging preserved

### Performance Requirements (from SDD specifications)
- [x] **Context Loading**: <2 seconds for files <10MB per SDD section 6.3.1
- [x] **Memory Usage**: <100MB during context processing per SDD section 6.3.1
- [x] **Concurrent Access**: Safe multi-agent execution per SDD section 6.3.2

### Quality Gates
- [x] **Unit Tests**: 85% coverage with critical functions at 95-100% per SDD section 6.4
- [x] **Integration Tests**: Complete hook workflow validated per SDD section 6.2
- [x] **Performance Tests**: All benchmarks meet requirements per SDD section 6.3
- [x] **Build Success**: `go build -o the-startup` completes without errors
- [x] **Static Analysis**: `go fmt ./... && go vet ./...` passes cleanly

## Implementation Notes

### SDD Section References
- AgentID specifications: SDD sections 3.2.1-3.2.3
- CLI specifications: SDD sections 3.3.1-3.3.4
- File format specifications: SDD sections 3.5.1-3.5.4
- Data model specifications: SDD sections 3.6.1-3.6.2
- Template specifications: SDD sections 3.7.1-3.7.4
- Test specifications: SDD sections 6.1-6.4

### Implementation Dependencies
- Existing internal/log package architecture
- Existing appendJSONL() and session detection functionality
- Claude Code hook system (PreToolUse/PostToolUse)
- Go standard library (filepath, os, crypto/sha256, regexp)

### Error Handling Strategy
- Silent failures for hook commands (exit code 0)
- Detailed errors for direct CLI usage
- Empty context returns (not errors) for missing files
- Graceful degradation for corrupted JSONL entries

This implementation plan executes the complete technical design specified in SDD.md without requiring additional analysis or design decisions.