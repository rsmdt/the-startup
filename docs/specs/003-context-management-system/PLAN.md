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

- [ ] Implement AgentID extraction function using regex pattern from SDD section 3.2.1 {agent: `the-developer`}
  - [ ] Create `ExtractAgentIDFromPrompt()` with pattern `(?i)\bAgentId\s*:\s*([a-zA-Z0-9][a-zA-Z0-9\-_]{0,62}[a-zA-Z0-9])`
  - [ ] Implement validation rules: 2-64 chars, alphanumeric with hyphens/underscores
  - [ ] Add reserved word checking for ["main", "global", "system"]
  - [ ] Add case conversion to lowercase as specified

- [ ] Implement fallback generation algorithm as specified in SDD section 3.2.2 {agent: `the-developer`}
  - [ ] Create `GenerateAgentID()` using SHA256 hash of "{sessionID}:{agentType}:{promptExcerpt100}"
  - [ ] Format output as "{agentType}-{YYYYMMDD-HHMMSS}-{hash8}"
  - [ ] Implement uniqueness checking in session directory
  - [ ] Add disambiguation with "-N" suffix (N=1-999)
  - [ ] Add ultimate fallback with 8-character UUID segment

- [ ] Implement complete extraction function as specified in SDD section 3.2.3 {agent: `the-developer`}
  - [ ] Create `ExtractOrGenerateAgentID()` combining extraction and fallback
  - [ ] Add truncateString() helper for 100-character prompt excerpts
  - [ ] Add ensureUniqueness() helper for session-scoped uniqueness
  - [ ] Implement error handling as specified

- [ ] **Validation**: Test AgentID management using specifications from SDD section 6.1.1
  - [ ] **Command**: `go test ./internal/log/... -run TestExtractAgentID -v`
  - [ ] **Command**: `go test ./internal/log/... -run TestDeterministicGeneration -v`

## Phase 2: Context File Operations Implementation

*Implement file management as specified in SDD sections 3.5 and 3.6*

- [ ] Implement AgentFileManager structure as specified in SDD section 3.6.1 {agent: `the-developer`}
  - [ ] Create struct with SessionID, AgentID, BaseDir, mutex fields
  - [ ] Implement GetAgentFilePath() method returning ".the-startup/{sessionID}/{agentID}.jsonl"
  - [ ] Implement GetSessionDir() method
  - [ ] Implement EnsureSessionDir() with 0755 permissions

- [ ] Implement WriteAgentContext function as specified in SDD section 3.3.4 {agent: `the-developer`}
  - [ ] Create function with signature `WriteAgentContext(sessionID, agentID string, data *HookData) error`
  - [ ] Add parameter validation for non-empty sessionID, agentID, and data
  - [ ] Implement directory creation with os.MkdirAll(filepath.Dir(agentFile), 0755)
  - [ ] Add atomic write with fileMutex.Lock() protection
  - [ ] Use existing appendJSONL() function for writing

- [ ] Implement ReadAgentContext function as specified in SDD section 3.3.4 {agent: `the-developer`}
  - [ ] Create function with signature `ReadAgentContext(sessionID, agentID string, maxLines int) ([]*HookData, error)`
  - [ ] Add parameter validation (agentID required, maxLines 1-1000 with default 50)
  - [ ] Implement session resolution using findLatestSessionWithAgent() when sessionID empty
  - [ ] Add file existence check returning empty slice (not error) for missing files
  - [ ] Use readLastNLines() for efficient tail reading

- [ ] Implement file reading helpers as specified in SDD section 3.3.4 {agent: `the-developer`}
  - [ ] Create findLatestSessionWithAgent() scanning .the-startup/ for latest session with agent file
  - [ ] Create readLastNLines() with 1MB threshold for small vs large file handling
  - [ ] Create readTailWithBuffer() for efficient reverse reading of large files (8KB buffer)
  - [ ] Add JSONL parsing with error recovery (skip corrupted lines)

- [ ] **Validation**: Test file operations using specifications from SDD section 6.1.2
  - [ ] **Command**: `go test ./internal/log/... -run TestWriteAgentContext -v`
  - [ ] **Command**: `go test ./internal/log/... -run TestReadAgentContext -v`

## Phase 3: CLI Enhancement Implementation

*Implement --read flag as specified in SDD section 3.3*

- [ ] Implement LogFlags structure as specified in SDD section 3.6.1 {agent: `the-developer`}
  - [ ] Add Read, AgentID, Lines, Session, Format, IncludeMetadata fields
  - [ ] Maintain existing Assistant and User fields

- [ ] Implement command validation as specified in SDD section 3.3.2 {agent: `the-developer`}
  - [ ] Create ValidateLogCommand() enforcing exactly one mode (assistant, user, or read)
  - [ ] Add --agent-id requirement when --read is set
  - [ ] Add AgentID format validation using isValidAgentID()
  - [ ] Add --lines range validation (1-1000)

- [ ] Extend log command with --read capability as specified in SDD section 3.3.1 {agent: `the-developer`}
  - [ ] Add --read, -r flag (bool, conflicts with --assistant/--user)
  - [ ] Add --agent-id flag (string, required when --read set)
  - [ ] Add --lines flag (int, default 50, max 1000)
  - [ ] Add --session flag (string, optional)
  - [ ] Add --format flag (string, default "json", allowed ["json", "text"])
  - [ ] Add --include-metadata flag (bool, default false)

- [ ] Implement context reading execution path as specified in SDD section 3.3.3 {agent: `the-developer`}
  - [ ] Create ContextQuery and ContextResponse structures from SDD section 3.6.1
  - [ ] Implement JSON output format with metadata and entries as specified
  - [ ] Implement text output format as specified
  - [ ] Add silent error handling (exit 0) for hook compatibility

- [ ] **Validation**: Test CLI enhancement using integration tests from SDD section 6.2.1
  - [ ] **Command**: `echo '{"tool_name":"Task","tool_input":{"subagent_type":"the-architect","prompt":"AgentId: arch-001"}}' | go run . log --assistant`
  - [ ] **Command**: `go run . log --read --agent-id arch-001 --lines 10`

## Phase 4: Hook Processing Integration

*Integrate AgentID management with existing hook processing as specified in SDD*

- [ ] Enhance ProcessToolCall function to use AgentID extraction as specified {agent: `the-developer`}
  - [ ] Replace existing AgentID extraction with ExtractOrGenerateAgentID() from Phase 1
  - [ ] Pass agentType (subagent_type), prompt, and sessionID to extraction function
  - [ ] Update HookData structure population with extracted/generated AgentID

- [ ] Implement agent-specific context writing as specified in SDD sections 3.5.2-3.5.3 {agent: `the-developer`}
  - [ ] Call WriteAgentContext() for agent-specific files (.the-startup/{sessionId}/{agentId}.jsonl)
  - [ ] Maintain existing writeSessionLog() for backward compatibility (agent-instructions.jsonl)
  - [ ] Maintain existing global logging (all-agent-instructions.jsonl)
  - [ ] Handle both PreToolUse and PostToolUse events

- [ ] Implement orchestrator context routing as specified in SDD section 3.5.2 {agent: `the-developer`}
  - [ ] Route main agent interactions (no subagent_type) to main.jsonl
  - [ ] Route subagent interactions to individual agent files
  - [ ] Ensure all interactions continue to legacy files for compatibility

- [ ] **Validation**: Test complete hook integration using specifications from SDD section 6.2.1
  - [ ] **Command**: Run complete integration test script from SDD section 6.2.1
  - [ ] **Command**: `go test ./internal/log/... -v`

## Phase 5: Agent Template Updates

*Apply context loading instructions as specified in SDD section 3.7*

- [ ] Update the-architect.md template as specified in SDD section 3.7.3 {agent: `the-developer`}
  - [ ] Insert Context Discovery section near beginning
  - [ ] Add context loading command: `CONTEXT=$(the-startup log --read --agent-id arch-{component} --lines 20 --format json 2>/dev/null || echo '{"entries":[]}')`
  - [ ] Add conditional context integration with HAS_CONTEXT check
  - [ ] Add context summary display and integration notes

- [ ] Update the-developer.md template as specified in SDD section 3.7.3 {agent: `the-developer`}
  - [ ] Insert Context Discovery section with dev-{feature-area} pattern
  - [ ] Apply same context loading and integration pattern
  - [ ] Ensure template handles both contextual and fresh development tasks

- [ ] Update the-tester.md template as specified in SDD section 3.7.3 {agent: `the-developer`}
  - [ ] Insert Context Discovery section with test-{test-type} pattern
  - [ ] Apply context loading pattern for test continuity

- [ ] Update remaining agent templates as specified in SDD sections 3.7.2-3.7.3 {agent: `the-developer`}
  - [ ] Update the-devops-engineer.md with devops-{infrastructure-component} pattern
  - [ ] Update the-security-engineer.md with sec-{security-domain} pattern
  - [ ] Update the-product-manager.md with pm-{product-area} pattern
  - [ ] Update the-technical-writer.md with doc-{documentation-type} pattern
  - [ ] Update the-data-engineer.md with data-{data-component} pattern
  - [ ] Update the-site-reliability-engineer.md with sre-{reliability-area} pattern
  - [ ] Update the-business-analyst.md with ba-{business-domain} pattern
  - [ ] Update the-project-manager.md with proj-{project-phase} pattern
  - [ ] Update the-prompt-engineer.md with prompt-{prompt-domain} pattern
  - [ ] Update the-chief.md with chief-{decision-area} pattern

- [ ] **Validation**: Test template context integration using performance tests from SDD section 6.3
  - [ ] **Command**: Create test context and measure loading performance (<2 seconds)
  - [ ] **Command**: Verify context format compatibility with jq parsing

## Phase 6: Testing Implementation

*Execute test specifications defined in SDD section 6*

- [ ] Implement unit tests as specified in SDD section 6.1 {agent: `the-developer`}
  - [ ] Create TestExtractAgentID_ValidCases with test cases from SDD section 6.1.1
  - [ ] Create TestExtractAgentID_InvalidCases validating fallback generation
  - [ ] Create TestDeterministicGeneration validating consistent ID generation
  - [ ] Create TestWriteAgentContext validating file operations
  - [ ] Create TestReadAgentContext with edge cases (corruption, missing files)

- [ ] Implement integration tests as specified in SDD section 6.2 {agent: `the-developer`}
  - [ ] Create integration_test_hook_processing.sh from SDD section 6.2.1
  - [ ] Test complete PreToolUse/PostToolUse workflow
  - [ ] Verify agent-specific and backward compatibility file creation
  - [ ] Test context reading functionality

- [ ] Implement performance tests as specified in SDD section 6.3 {agent: `the-developer`}
  - [ ] Create performance_test_large_files.sh from SDD section 6.3.1
  - [ ] Test 50MB file handling with <2 second read requirement
  - [ ] Test memory usage <100MB requirement
  - [ ] Create concurrent_test.sh from SDD section 6.3.2
  - [ ] Test parallel agent execution safety

- [ ] Execute test coverage validation as specified in SDD section 6.4 {agent: `the-developer`}
  - [ ] Achieve 85% minimum unit test coverage
  - [ ] Achieve 100% coverage for ExtractOrGenerateAgentID
  - [ ] Achieve 95% coverage for WriteAgentContext and ReadAgentContext
  - [ ] Verify all edge cases covered (missing context, corruption, permissions)

- [ ] **Final Validation**: Execute complete test suite
  - [ ] **Command**: `go test -v ./...`
  - [ ] **Command**: `go test -cover ./...`
  - [ ] **Command**: `bash integration_test_hook_processing.sh`
  - [ ] **Command**: `bash performance_test_large_files.sh`
  - [ ] **Command**: `bash concurrent_test.sh`

## Success Criteria

### Functional Requirements (from SDD specifications)
- [ ] **AgentID Management**: Regex extraction with fallback generation implemented per SDD section 3.2
- [ ] **Context Writing**: Individual agent JSONL files created per SDD section 3.5
- [ ] **Context Reading**: CLI --read flag implemented per SDD section 3.3
- [ ] **Agent Integration**: All 13 templates updated with context loading per SDD section 3.7
- [ ] **Backward Compatibility**: Existing agent-instructions.jsonl and global logging preserved

### Performance Requirements (from SDD specifications)
- [ ] **Context Loading**: <2 seconds for files <10MB per SDD section 6.3.1
- [ ] **Memory Usage**: <100MB during context processing per SDD section 6.3.1
- [ ] **Concurrent Access**: Safe multi-agent execution per SDD section 6.3.2

### Quality Gates
- [ ] **Unit Tests**: 85% coverage with critical functions at 95-100% per SDD section 6.4
- [ ] **Integration Tests**: Complete hook workflow validated per SDD section 6.2
- [ ] **Performance Tests**: All benchmarks meet requirements per SDD section 6.3
- [ ] **Build Success**: `go build -o the-startup` completes without errors
- [ ] **Static Analysis**: `go fmt ./... && go vet ./...` passes cleanly

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