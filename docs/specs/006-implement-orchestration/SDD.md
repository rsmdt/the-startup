# Solution Design Document

## Overview

This document outlines the technical design for enhancing the `/s:implement` command to use intelligent orchestration with dynamic review selection, transforming it from a direct executor to a sophisticated delegation system that ensures quality through automated, context-aware review cycles.

## Goals

1. **Intelligent Task Delegation**: Transform `/s:implement` to orchestrate specialist agents rather than execute directly
2. **Dynamic Review Selection**: Automatically select appropriate reviewers based on task context and changes made
3. **Quality Assurance**: Implement review-revision cycles until quality standards are met
4. **Future-Proof Design**: Support new agents automatically without code changes
5. **Context Preservation**: Ensure all agents have access to specification documents (BRD/PRD/SDD)

## Non-Goals

- State persistence mechanisms (Git provides natural state through PLAN.md checkboxes)
- Rollback strategies (Git handles version control)
- Static reviewer mappings (all selection is dynamic)
- Complex error recovery (simple retry/skip/abort is sufficient)

## Technical Architecture

### Component Overview

```
/s:implement [spec-id]
    ├── Context Loader
    │   ├── Read BRD.md (if exists)
    │   ├── Read PRD.md (if exists)
    │   ├── Read SDD.md (if exists)
    │   └── Read PLAN.md (required)
    │
    ├── Task Orchestrator
    │   ├── Parse task structure
    │   ├── Identify execution type (parallel/sequential)
    │   ├── Select appropriate agent
    │   └── Delegate with context
    │
    ├── Review Selector
    │   ├── Analyze implementation output
    │   ├── Identify review needs
    │   ├── Select reviewer dynamically
    │   └── Provide review context
    │
    └── Cycle Manager
        ├── Parse review feedback
        ├── Determine if revision needed
        ├── Re-delegate with feedback
        └── Update PLAN.md checkboxes
```

### Dynamic Review Selection Algorithm

```markdown
FUNCTION selectReviewer(task, implementation, availableAgents):
    // Analyze what was implemented
    changes = analyzeChanges(implementation)
    
    // Identify review priorities
    priorities = []
    IF changes.includes(authentication, authorization, encryption):
        priorities.add("security")
    IF changes.includes(database, queries, migrations):
        priorities.add("data integrity")
    IF changes.includes(newPatterns, architecturalChanges):
        priorities.add("architecture")
    IF changes.includes(performance, caching, optimization):
        priorities.add("performance")
    IF changes.includes(API, contracts, interfaces):
        priorities.add("integration")
    
    // Select best reviewer using natural language reasoning
    reviewer = reasonAboutBestReviewer(
        taskContext: task,
        changeMade: changes,
        reviewPriorities: priorities,
        availableAgents: availableAgents
    )
    
    RETURN reviewer
```

### Task Structure Enhancement

Tasks in PLAN.md will support additional metadata:

```markdown
- [ ] **Task Description** [`agent: the-developer`] [`review: true`] [`review_focus: security, patterns`]
  - Subtask details
  - Implementation notes
```

### Review Cycle Flow

```
1. Agent completes implementation
2. System analyzes output for review needs
3. Dynamically selects appropriate reviewer
4. Reviewer provides feedback:
   - APPROVED: Mark task complete, continue
   - NEEDS_REVISION: Specific feedback provided
5. If revision needed:
   - Original agent receives feedback
   - Implements changes
   - Return to step 2
6. Continue until approved or user intervenes
```

## Implementation Patterns

### Context Ingestion Pattern

```markdown
## Phase 0: Context Loading
Before executing any tasks, load specification context:

1. Check for specification documents:
   - IF exists docs/specs/[ID]/BRD.md: Extract business context
   - IF exists docs/specs/[ID]/PRD.md: Extract product requirements  
   - IF exists docs/specs/[ID]/SDD.md: Extract technical design
   
2. Parse PLAN.md for:
   - Task list with metadata
   - Execution types (parallel/sequential)
   - Validation checkpoints
   - Completion status (checkboxes)
```

### Agent Invocation Pattern

```markdown
## Task Delegation

For each task:
1. Extract task metadata (agent, review requirements)
2. Build context package:
   - Specification context (from BRD/PRD/SDD)
   - Task requirements
   - Dependencies completed
   - Success criteria
   
3. Invoke specified agent with bounded context:
   PROMPT: """
   CONTEXT: [specification summary]
   TASK: [specific task]
   SUCCESS: Task complete when [criteria]
   EXCLUDE: [out of scope items]
   """
```

### Review Selection Pattern

```markdown
## Intelligent Review Selection

After task completion:
1. Analyze implementation:
   - What files were changed?
   - What patterns were used?
   - What risks might exist?
   
2. Select reviewer through reasoning:
   "Given that this task involved [changes],
    and considering [risks/concerns],
    the best agent to review this would be [agent]
    because of their expertise in [relevant area]."
    
3. Invoke reviewer with context:
   PROMPT: """
   REVIEW REQUEST
   Original Task: [task description]
   Implemented by: [agent]
   Changes Made: [summary]
   Focus Areas: [identified concerns]
   Please review for: [specific aspects]
   """
```

## Data Flow

### Input Structure
- Spec ID or path to PLAN.md
- Existing specification documents (BRD/PRD/SDD)
- Current PLAN.md with task checkboxes

### Processing Flow
1. Load all context documents
2. Parse uncompleted tasks from PLAN.md
3. Execute tasks according to phases
4. For each task:
   - Delegate to specified agent
   - If review required, select and invoke reviewer
   - Handle review cycle until approved
   - Update PLAN.md checkbox
5. Continue until all tasks complete

### Output Structure
- Updated PLAN.md with completed checkboxes
- Implementation artifacts from agents
- Review feedback trail
- Final validation results

## Error Handling

### Agent Failures
- If agent reports BLOCKED: Present options to user
- If agent errors: Retry with clarified context
- If repeated failures: Allow skip with user confirmation

### Review Cycle Limits
- Maximum 3 review cycles per task
- After 3 cycles, escalate to user for decision
- User can: accept as-is, manually fix, or skip

### Recovery Strategy
- No special state needed - PLAN.md checkboxes show progress
- On restart, continue from first unchecked task
- Context reloaded fresh each session

## Security Considerations

- Never expose credentials in context passed to agents
- Validate all agent outputs before proceeding
- Review selection prioritizes security when auth/data involved
- Maintain audit trail of all delegations and reviews

## Performance Optimizations

### Parallel Execution
- Identify independent tasks within phases
- Launch multiple agents simultaneously
- Synchronize at phase boundaries
- Batch review requests when possible

### Context Efficiency
- Load specification documents once, reuse for all tasks
- Pass minimal necessary context to each agent
- Cache review decisions for similar tasks

## Testing Strategy

### Unit Tests
- Review selector logic with various scenarios
- Context loading with missing documents
- Task parsing from PLAN.md

### Integration Tests
- Full cycle: implement → review → revise → approve
- Parallel task execution
- Error recovery scenarios
- Progress resumption

### End-to-End Tests
- Complete specification implementation
- Multiple review cycles
- Various agent combinations
- Edge cases (no review needed, all reviews fail)

## Migration Path

### Phase 1: Template Update
- Enhance PLAN.md template with review metadata
- Add context ingestion instructions
- Backward compatible with existing plans

### Phase 2: Command Enhancement
- Add orchestration logic to /s:implement
- Implement dynamic review selection
- Maintain fallback to direct execution

### Phase 3: Full Deployment
- Default to orchestration mode
- Remove direct execution code
- Update all documentation

## Success Metrics

- Review catches issues before they reach main branch
- Reduced implementation errors through specialist delegation  
- Faster implementation through parallel execution
- Improved code quality through automated review cycles
- Support for new agents without code changes

## Example Scenarios

### Scenario 1: Security-Critical Implementation

```
Task: Implement JWT authentication
Agent: the-developer
Implementation: Creates auth middleware

Review Selection:
- Identifies: Authentication, token handling
- Selects: the-security-engineer
- Focus: Token validation, session management

Review Feedback: "Add rate limiting"
Revision: the-developer adds rate limiting
Second Review: Approved
```

### Scenario 2: Performance Optimization

```
Task: Optimize database queries
Agent: the-data-engineer  
Implementation: Adds indexing and query optimization

Review Selection:
- Identifies: Database changes, performance impact
- Selects: the-architect
- Focus: System-wide impact, scalability

Review: Approved with suggestions for monitoring
```

### Scenario 3: New Agent Available

```
Task: Implement AI feature
Agent: the-developer
Implementation: Basic AI integration

System discovers: the-ai-engineer (newly added)
Review Selection:
- Identifies: AI/ML patterns
- Selects: the-ai-engineer (without being hardcoded)
- Focus: Model usage, prompt engineering

Review provides specialized AI feedback
```

## Conclusion

This design transforms `/s:implement` into an intelligent orchestrator that:
- Delegates implementation to appropriate specialists
- Dynamically selects reviewers based on context
- Ensures quality through automated review cycles
- Adapts to new agents without modification
- Maintains simplicity through git-native state management

The solution provides robust quality assurance while remaining flexible and future-proof.