# Implementation Plan

## Overview

This plan implements the enhanced `/s:specify` command with intelligent complexity assessment, user control gates, clarification protocols, and session management capabilities. Since this is prompt engineering (not software development), all changes will be made to the markdown prompt file.

## Phase 1: Preparation & Analysis
**Execution: Sequential**

- [x] Read current `/s:specify` command: @assets/commands/s/specify.md
- [x] Review session protocol template: @assets/templates/SESSION_PROTOCOL.md
- [x] Analyze existing sub-agent prompts for compatibility
- [x] Create backup of current specify.md
- [x] **Validation**: Confirm backup created at specify.md.backup

## Phase 2: Core Prompt Restructuring
**Execution: Sequential**

- [x] Remove "delegate EVERYTHING" mandate (line 13)
- [x] Add Session Initialization section
  - [x] SessionID generation logic
  - [x] Resume detection from $ARGUMENTS
  - [x] State file creation instructions
- [x] Implement Complexity Assessment Protocol
  - [x] Level 1-3 classification criteria
  - [x] Display format with tree structure
  - [x] User override options
- [ ] **Validation**: Test with simple requirement to verify L1 classification

## Phase 3: User Control Implementation
**Execution: Sequential**

- [x] Add User Confirmation Gates section
  - [x] Pre-delegation gate with reasoning display
  - [x] Post-response review gate
  - [x] Document transition approval
  - [x] Format: emoji indicators, lettered choices
- [x] Implement Clarification-First Protocol
  - [x] Ambiguity detection rules
  - [x] Question formatting (numbered, specific)
  - [x] STOP-ASK-WAIT-CONFIRM flow
- [x] Add Anti-Drift Enforcement
  - [x] EXCLUDE section requirements
  - [x] Drift detection and alerting
  - [x] Success criteria validation
- [ ] **Validation**: Test with ambiguous input to trigger clarification

## Phase 4: Session Management Integration (Simplified)
**Execution: Sequential**

- [x] Add Session Management section
  - [x] SessionID format: specify-{timestamp}
  - [x] AgentID format: {type}-{context}-{seq}
  - [x] Agent registry initialization
- [x] Update Sub-Agent Invocation format
  - [x] Include Session Identity block
  - [x] Add Context Retrieval instructions
  - [x] Provide the-startup log command syntax
- [x] Implement Session Continuity
  - [x] Same SessionID for entire Claude Code instance
  - [x] Agent registry for ID reuse
  - [x] State file as audit trail
  - [x] Note: Resume capability deferred
- [x] **Validation**: Verify agents can retrieve previous context

## Phase 5: Bounded Context Protocol
**Execution: Sequential**

- [x] Define Bounded Context Format
  - [x] Session Identity (SessionId, AgentId)
  - [x] Context Retrieval instructions
  - [x] Task structure (TASK, CONTEXT, CONSTRAINTS, etc.)
  - [x] EXCLUDE section mandatory
- [x] Add context examples for each complexity level
  - [x] L1: No delegation needed
  - [x] L2: Brief consultation format
  - [x] L3: Full delegation format
- [x] **Validation**: Verify sub-agent receives proper context format

## Phase 6: State Persistence Implementation
**Execution: Sequential**

- [x] Define State File Structure
  - [x] Location: .the-startup/{session-id}/state.md
  - [x] Sections: Status, Registry, History, Checkpoints
  - [x] Markdown format for @ notation readability
- [x] Add checkpoint saving instructions
  - [x] After complexity assessment
  - [x] At each user gate
  - [x] After sub-agent responses
- [x] Implement state updates
  - [x] Agent registry management
  - [x] Decision history tracking
  - [x] Next steps documentation
- [x] **Validation**: Verify state file created and updated correctly

## Phase 7: Testing & Refinement
**Execution: Parallel where possible**

- [x] Test Level 1 direct execution
  - [x] Simple, single-domain task
  - [x] Verify no delegation occurs
  - [x] Confirm no latency penalty
- [x] Test Level 2 consultation
  - [x] Moderate complexity task
  - [x] Verify brief delegation
  - [x] Check bounded context
- [x] Test Level 3 full delegation
  - [x] Complex, multi-domain task
  - [x] Verify proper routing
  - [x] Validate EXCLUDE enforcement
- [x] Test session continuity
  - [x] Verify same SessionID throughout
  - [x] Check AgentID reuse for same context
  - [x] Confirm agents can retrieve previous work
- [x] **Validation**: All test scenarios pass

## Phase 8: Documentation & Cleanup
**Execution: Sequential**

- [x] Update command help text in YAML frontmatter
  - [x] Change argument-hint from "describe your feature OR provide spec ID to resume"
  - [x] To: "describe your feature or requirement to specify"
- [x] Add inline examples for each feature
- [x] Document measurement points for unknowns
- [x] Remove debug/test code if any
- [x] Format markdown for readability
- [x] **Validation**: Command help displays correctly

## Phase 9: Automated Testing Strategy
**Execution: Sequential**

- [x] Create test scenarios in `tests/` directory
  - [x] L1_simple_task.txt: "Add a submit button to the form"
  - [x] L2_moderate_task.txt: "Add user authentication with email"
  - [x] L3_complex_task.txt: "Design real-time collaboration system"
  - [x] Ambiguous_input.txt: "Make it better"
  - [x] Parallel_conflict.txt: Task requiring multiple agents with conflicts
- [x] Create test harness script using Claude CLI
  ```bash
  #!/bin/bash
  # test_specify.sh
  
  # Test L1 - should not delegate
  echo "Testing Level 1 classification..."
  claude -p --verbose --max-turns 3 "/s:specify $(cat tests/L1_simple_task.txt)"
  
  # Test L2 - should consult briefly
  echo "Testing Level 2 classification..."
  claude -p --verbose --max-turns 5 "/s:specify $(cat tests/L2_moderate_task.txt)"
  
  # Test L3 - should delegate fully
  echo "Testing Level 3 classification..."
  claude -p --verbose --max-turns 7 "/s:specify $(cat tests/L3_complex_task.txt)"
  
  # Test ambiguity detection
  echo "Testing ambiguity detection..."
  claude -p --verbose --max-turns 2 "/s:specify $(cat tests/Ambiguous_input.txt)"
  ```
- [x] Verify test outputs
  - [x] L1 completes without Task tool invocation
  - [x] L2/L3 show proper complexity assessment
  - [x] Ambiguous input triggers clarification
  - [x] Session IDs generated correctly
  - [x] Agent IDs follow format
- [x] Document test execution process
- [x] **Validation**: All test scenarios pass

## Validation Checklist

### Functional Requirements
- [x] Complexity assessment works for L1, L2, L3
- [x] User gates stop at all required points
- [x] Clarification triggered for vague input
- [x] Session management creates proper IDs
- [x] Resume capability reconstructs state
- [x] Sub-agents can retrieve context

### Non-Functional Requirements
- [x] L1 tasks execute without delegation latency
- [x] State files readable via @ notation
- [x] Agent logs accessible via the-startup command
- [x] All user decisions logged to state
- [x] Bounded context stays under limits

### Edge Cases
- [x] Handles non-existent session for resume
- [x] Gracefully manages corrupted state files
- [x] Prevents AgentID collisions
- [x] Detects and alerts on drift
- [x] Handles interruption at any point

## Anti-Patterns to Avoid

### Prompt Engineering Anti-Patterns
- ❌ Making sub-agents dependent on main agent state
- ❌ Passing full context instead of bounded context
- ❌ Auto-progressing without user confirmation
- ❌ Making assumptions instead of asking questions
- ❌ Ignoring EXCLUDE boundaries

### Session Management Anti-Patterns
- ❌ Using @ notation for JSONL logs (use the-startup command)
- ❌ Embedding context data in prompts (provide retrieval instructions)
- ❌ Creating new SessionID mid-workflow (maintain continuity)
- ❌ Overwriting agent IDs (maintain registry)
- ❌ Complex resume logic (deferred to future)

### Implementation Anti-Patterns
- ❌ Modifying Task tool or hook system (prompt-only changes)
- ❌ Creating new file types (use existing .md and .jsonl)
- ❌ Changing core Claude Code behavior (work within constraints)
- ❌ Adding code dependencies (this is prompt engineering)
- ❌ Breaking existing workflows (maintain compatibility)

## Success Criteria

The implementation is complete when:
1. Simple tasks execute directly without delegation
2. Complex tasks properly delegate with bounded context
3. Users maintain control at all decision points
4. Sessions can be interrupted and resumed
5. Sub-agents can retrieve and continue previous work
6. All features documented in PRD are functional
7. Measurement points identified for unknown metrics