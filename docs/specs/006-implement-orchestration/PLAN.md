# Implementation Plan

## Context Documents

Before beginning implementation, the following specification documents should be reviewed if they exist:
- **BRD.md**: Business requirements and objectives (if exists)
- **PRD.md**: Product requirements and user stories (if exists)  
- **SDD.md**: Technical design and architecture (required - see above)

Key context from specifications:
- Transform /s:implement into an orchestrator that delegates to specialist agents
- Implement dynamic, context-aware review selection
- Ensure quality through automated review-revision cycles
- Support future agents without code modifications

## Phase 1: Template Enhancement
**Execution: Sequential**

- [ ] **Read current PLAN.md template** [`agent: the-developer`] [`review: false`]
  - Read assets/templates/PLAN.md
  - Understand current structure
  - Identify insertion points for new sections

- [ ] **Add context ingestion section** [`agent: the-developer`] [`review: true`] [`review_focus: template structure, clarity`]
  - Add "## Context Documents" section at the beginning
  - Include instructions to read BRD/PRD/SDD if they exist
  - Add placeholder for extracting key context
  - Ensure backwards compatibility

- [ ] **Enhance task metadata structure** [`agent: the-developer`] [`review: true`] [`review_focus: syntax, usability`]
  - Add documentation for task metadata fields
  - Include examples of agent assignment: [`agent: agent-name`]
  - Include examples of review flags: [`review: true/false`]
  - Include examples of review focus: [`review_focus: areas`]
  - Update instruction blocks

- [ ] **Validation**: Manually review updated template for clarity and completeness

## Phase 2: Command Transformation
**Execution: Sequential**

- [ ] **Analyze current implementation** [`agent: the-architect`] [`review: false`]
  - Read assets/commands/s/implement.md
  - Identify sections to modify
  - Plan transformation approach
  - Document integration points

- [ ] **Add orchestration core** [`agent: the-developer`] [`review: true`] [`review_focus: architecture, patterns`]
  - Replace direct execution with delegation pattern
  - Add agent selection logic based on task metadata
  - Implement parallel/sequential execution handling
  - Maintain TodoWrite synchronization

- [ ] **Implement dynamic review selection** [`agent: the-developer`] [`review: true`] [`review_focus: logic, flexibility`]
  - Add review analysis section
  - Implement context-based reviewer selection
  - Create review prompt templates
  - Add natural language reasoning for selection

- [ ] **Add review cycle management** [`agent: the-developer`] [`review: true`] [`review_focus: flow control, error handling`]
  - Implement feedback parsing
  - Add revision delegation logic
  - Create cycle limits (max 3 attempts)
  - Add user intervention points

- [ ] **Validation**: Test orchestration logic with mock scenarios

## Phase 3: Context Loading Implementation
**Execution: Sequential**

- [ ] **Create context loader section** [`agent: the-developer`] [`review: true`] [`review_focus: completeness, error handling`]
  - Add Phase 0: Context Loading
  - Implement BRD/PRD/SDD detection and reading
  - Extract and format key information
  - Pass context to agents appropriately

- [ ] **Enhance progress tracking** [`agent: the-developer`] [`review: false`]
  - Ensure PLAN.md checkbox updates work correctly
  - Add review cycle tracking to progress reports
  - Maintain clear status reporting

- [ ] **Validation**: Run `go test ./...` to ensure no regressions

## Phase 4: Integration Testing
**Execution: Sequential**

- [ ] **Create test specification** [`agent: the-developer`] [`review: false`]
  - Create simple test spec in docs/specs/test-implement/
  - Include tasks requiring different agents
  - Include review scenarios

- [ ] **Test basic orchestration** [`agent: the-tester`] [`review: false`]
  - Test task delegation to agents
  - Verify parallel execution works
  - Confirm sequential execution works
  - Check TodoWrite updates

- [ ] **Test review cycles** [`agent: the-tester`] [`review: false`]
  - Test successful review (approval)
  - Test review with revision needed
  - Test multiple review cycles
  - Test review cycle limits

- [ ] **Test edge cases** [`agent: the-tester`] [`review: false`]
  - Missing specification documents
  - No review required tasks
  - Agent failures and recovery
  - Progress resumption

- [ ] **Validation**: All test scenarios pass

## Phase 5: Documentation and Cleanup
**Execution: Parallel**

- [ ] **Update command documentation** [`agent: the-technical-writer`] [`review: false`]
  - Document new orchestration behavior
  - Add examples of review selection
  - Update usage instructions

- [ ] **Clean up any deprecated code** [`agent: the-developer`] [`review: false`]
  - Remove old direct execution logic if present
  - Clean up comments
  - Ensure code follows conventions

- [ ] **Final validation** [`agent: the-lead-developer`] [`review: false`]
  - Review all changes for quality
  - Ensure backward compatibility
  - Verify no regressions

## Completion Criteria

- [ ] PLAN.md template includes context ingestion
- [ ] /s:implement successfully orchestrates tasks to agents
- [ ] Review selection works dynamically based on context
- [ ] Review cycles complete successfully
- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Backward compatibility maintained

## Anti-Patterns to Avoid

### Implementation Anti-Patterns
- ❌ Hardcoding agent names for review selection
- ❌ Creating complex state management systems
- ❌ Implementing rollback mechanisms (Git handles this)
- ❌ Over-engineering the review selection logic
- ❌ Breaking existing PLAN.md files

### Review Anti-Patterns  
- ❌ Forcing reviews on trivial tasks
- ❌ Infinite review loops without limits
- ❌ Ignoring review feedback
- ❌ Merging feedback from multiple reviewers
- ❌ Static reviewer assignments

### Testing Anti-Patterns
- ❌ Only testing happy paths
- ❌ Skipping edge case testing
- ❌ Not testing review cycles
- ❌ Ignoring backward compatibility tests