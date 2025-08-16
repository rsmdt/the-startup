# Implementation Plan

## Overview

This plan implements the enhanced `/s:specify` command with intelligent complexity assessment, user control gates, clarification protocols, and session management capabilities. Since this is prompt engineering (not software development), all changes will be made to the markdown prompt file.

## Phase 1: Preparation & Analysis
**Execution: Sequential**

- [ ] Read current `/s:specify` command: @assets/commands/s/specify.md
- [ ] Review session protocol template: @assets/templates/SESSION_PROTOCOL.md
- [ ] Analyze existing sub-agent prompts for compatibility
- [ ] Create backup of current specify.md
- [ ] **Validation**: Confirm backup created at specify.md.backup

## Phase 2: Core Prompt Restructuring
**Execution: Sequential**

- [ ] Remove "delegate EVERYTHING" mandate (line 13)
- [ ] Add Session Initialization section
  - [ ] SessionID generation logic
  - [ ] Resume detection from $ARGUMENTS
  - [ ] State file creation instructions
- [ ] Implement Complexity Assessment Protocol
  - [ ] Level 1-3 classification criteria
  - [ ] Display format with tree structure
  - [ ] User override options
- [ ] **Validation**: Test with simple requirement to verify L1 classification

## Phase 3: User Control Implementation
**Execution: Sequential**

- [ ] Add User Confirmation Gates section
  - [ ] Pre-delegation gate with reasoning display
  - [ ] Post-response review gate
  - [ ] Document transition approval
  - [ ] Format: emoji indicators, lettered choices
- [ ] Implement Clarification-First Protocol
  - [ ] Ambiguity detection rules
  - [ ] Question formatting (numbered, specific)
  - [ ] STOP-ASK-WAIT-CONFIRM flow
- [ ] Add Anti-Drift Enforcement
  - [ ] EXCLUDE section requirements
  - [ ] Drift detection and alerting
  - [ ] Success criteria validation
- [ ] **Validation**: Test with ambiguous input to trigger clarification

## Phase 4: Session Management Integration (Simplified)
**Execution: Sequential**

- [ ] Add Session Management section
  - [ ] SessionID format: specify-{timestamp}
  - [ ] AgentID format: {type}-{context}-{seq}
  - [ ] Agent registry initialization
- [ ] Update Sub-Agent Invocation format
  - [ ] Include Session Identity block
  - [ ] Add Context Retrieval instructions
  - [ ] Provide the-startup log command syntax
- [ ] Implement Session Continuity
  - [ ] Same SessionID for entire Claude Code instance
  - [ ] Agent registry for ID reuse
  - [ ] State file as audit trail
  - [ ] Note: Resume capability deferred
- [ ] **Validation**: Verify agents can retrieve previous context

## Phase 5: Bounded Context Protocol
**Execution: Sequential**

- [ ] Define Bounded Context Format
  - [ ] Session Identity (SessionId, AgentId)
  - [ ] Context Retrieval instructions
  - [ ] Task structure (TASK, CONTEXT, CONSTRAINTS, etc.)
  - [ ] EXCLUDE section mandatory
- [ ] Add context examples for each complexity level
  - [ ] L1: No delegation needed
  - [ ] L2: Brief consultation format
  - [ ] L3: Full delegation format
- [ ] **Validation**: Verify sub-agent receives proper context format

## Phase 6: State Persistence Implementation
**Execution: Sequential**

- [ ] Define State File Structure
  - [ ] Location: .the-startup/{session-id}/state.md
  - [ ] Sections: Status, Registry, History, Checkpoints
  - [ ] Markdown format for @ notation readability
- [ ] Add checkpoint saving instructions
  - [ ] After complexity assessment
  - [ ] At each user gate
  - [ ] After sub-agent responses
- [ ] Implement state updates
  - [ ] Agent registry management
  - [ ] Decision history tracking
  - [ ] Next steps documentation
- [ ] **Validation**: Verify state file created and updated correctly

## Phase 7: Testing & Refinement
**Execution: Parallel where possible**

- [ ] Test Level 1 direct execution
  - [ ] Simple, single-domain task
  - [ ] Verify no delegation occurs
  - [ ] Confirm no latency penalty
- [ ] Test Level 2 consultation
  - [ ] Moderate complexity task
  - [ ] Verify brief delegation
  - [ ] Check bounded context
- [ ] Test Level 3 full delegation
  - [ ] Complex, multi-domain task
  - [ ] Verify proper routing
  - [ ] Validate EXCLUDE enforcement
- [ ] Test session continuity
  - [ ] Verify same SessionID throughout
  - [ ] Check AgentID reuse for same context
  - [ ] Confirm agents can retrieve previous work
- [ ] **Validation**: All test scenarios pass

## Phase 8: Documentation & Cleanup
**Execution: Sequential**

- [ ] Update command help text in YAML frontmatter
  - [ ] Change argument-hint from "describe your feature OR provide spec ID to resume"
  - [ ] To: "describe your feature or requirement to specify"
- [ ] Add inline examples for each feature
- [ ] Document measurement points for unknowns
- [ ] Remove debug/test code if any
- [ ] Format markdown for readability
- [ ] **Validation**: Command help displays correctly

## Phase 9: Automated Testing Strategy
**Execution: Sequential**

- [ ] Create test scenarios in `tests/` directory
  - [ ] L1_simple_task.txt: "Add a submit button to the form"
  - [ ] L2_moderate_task.txt: "Add user authentication with email"
  - [ ] L3_complex_task.txt: "Design real-time collaboration system"
  - [ ] Ambiguous_input.txt: "Make it better"
  - [ ] Parallel_conflict.txt: Task requiring multiple agents with conflicts
- [ ] Create test harness script using Claude CLI
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
- [ ] Verify test outputs
  - [ ] L1 completes without Task tool invocation
  - [ ] L2/L3 show proper complexity assessment
  - [ ] Ambiguous input triggers clarification
  - [ ] Session IDs generated correctly
  - [ ] Agent IDs follow format
- [ ] Document test execution process
- [ ] **Validation**: All test scenarios pass

## Validation Checklist

### Functional Requirements
- [ ] Complexity assessment works for L1, L2, L3
- [ ] User gates stop at all required points
- [ ] Clarification triggered for vague input
- [ ] Session management creates proper IDs
- [ ] Resume capability reconstructs state
- [ ] Sub-agents can retrieve context

### Non-Functional Requirements
- [ ] L1 tasks execute without delegation latency
- [ ] State files readable via @ notation
- [ ] Agent logs accessible via the-startup command
- [ ] All user decisions logged to state
- [ ] Bounded context stays under limits

### Edge Cases
- [ ] Handles non-existent session for resume
- [ ] Gracefully manages corrupted state files
- [ ] Prevents AgentID collisions
- [ ] Detects and alerts on drift
- [ ] Handles interruption at any point

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