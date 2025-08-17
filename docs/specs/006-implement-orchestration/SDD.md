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

- Complex state management beyond review cycles (Git provides natural state through PLAN.md checkboxes)
- Rollback strategies (Git handles version control)
- Static reviewer mappings (all selection is dynamic)
- Complex error recovery beyond timeouts (simple retry/skip/abort is sufficient)

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

### State Management

#### Review Cycle Persistence

The system maintains review cycle state in `.the-startup/review-cycles.json` to survive process restarts:

```json
{
  "sessions": {
    "session-id-123": {
      "tasks": {
        "task-hash-abc": {
          "description": "Implement JWT authentication",
          "current_cycle": 2,
          "max_cycles": 3,
          "implementer": "the-developer",
          "reviewers": [
            {
              "agent": "the-security-engineer",
              "cycle": 1,
              "status": "NEEDS_REVISION",
              "feedback": "Add rate limiting",
              "timestamp": "2024-01-15T10:30:00Z"
            },
            {
              "agent": "the-security-engineer",
              "cycle": 2,
              "status": "pending",
              "timestamp": "2024-01-15T10:45:00Z"
            }
          ],
          "pattern_tracking": {
            "same_issue_count": 0,
            "recurring_patterns": []
          }
        }
      },
      "global_patterns": {
        "reviewer_failure_rates": {
          "the-security-engineer": {
            "total_reviews": 15,
            "revisions_required": 12,
            "failure_rate": 0.8
          }
        },
        "common_issues": [
          {
            "pattern": "missing_rate_limiting",
            "count": 5,
            "suggested_reviewer": "the-architect"
          }
        ]
      }
    }
  },
  "updated_at": "2024-01-15T10:45:00Z"
}
```

#### TodoWrite Integration

When using TodoWrite, include review cycle metadata:

```javascript
todos: [
  {
    "id": "task-1",
    "content": "Review implementation for security concerns",
    "status": "in_progress",
    "metadata": {
      "review_cycle": "2/3",
      "reviewer": "the-security-engineer",
      "task_hash": "abc123"
    }
  }
]
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
3. Check review history for patterns:
   - If high failure rate with current reviewer type
   - Suggest alternative reviewer
4. Dynamically selects appropriate reviewer
5. Start review with timeout (default: 5 minutes)
6. Reviewer provides feedback:
   - APPROVED: Mark task complete, continue
   - NEEDS_REVISION: Specific feedback provided
   - TIMEOUT: Apply default action (skip_with_warning)
7. If revision needed:
   - Update pattern tracking (same_issue_count)
   - Original agent receives feedback
   - Implements changes
   - Return to step 2
8. Continue until approved, max cycles reached, or timeout
9. Persist state to `.the-startup/review-cycles.json`
```

### User Escalation Timeout Configuration

```json
{
  "escalation": {
    "timeout_minutes": 5,
    "default_action": "skip_with_warning",
    "notification_method": "console",
    "actions": [
      "skip_with_warning",
      "accept_as_is",
      "abort_process"
    ]
  }
}
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
   
2. Check pattern history:
   - Load `.the-startup/review-cycles.json`
   - Check reviewer failure rates
   - Identify recurring issues (same_issue_count)
   - If failure_rate > 0.7 for intended reviewer:
     * Suggest alternative reviewer
     * Log pattern for future reference
   
3. Select reviewer through reasoning:
   "Given that this task involved [changes],
    and considering [risks/concerns],
    and noting [historical patterns if any],
    the best agent to review this would be [agent]
    because of their expertise in [relevant area]."
    
4. Invoke reviewer with context:
   PROMPT: """
   REVIEW REQUEST [Cycle: {current}/{max}]
   Original Task: [task description]
   Implemented by: [agent]
   Changes Made: [summary]
   Focus Areas: [identified concerns]
   Previous Issues: [if same_issue_count > 0]
   Please review for: [specific aspects]
   """
```

### Status Detection Patterns

#### Fuzzy Approval Detection

The system uses pattern matching to detect approval signals:

```javascript
const APPROVAL_PATTERNS = [
  /^APPROVED$/i,
  /^LOOKS\s+GOOD$/i,
  /^LGTM$/i,
  /^SHIP\s+IT$/i,
  /^\+1$/,
  /^READY\s+TO\s+(MERGE|SHIP)$/i,
  /^ALL\s+GOOD$/i,
  /^PASSED\s+REVIEW$/i,
  /^✅/,
  /^👍/
];

const REVISION_PATTERNS = [
  /^NEEDS[\s_]REVISION$/i,
  /^REQUIRES?\s+CHANGES?$/i,
  /^NEEDS?\s+WORK$/i,
  /^NOT\s+READY$/i,
  /^-1$/,
  /^BLOCKED$/i,
  /^FIX\s+REQUIRED$/i,
  /^❌/,
  /^👎/
];

function detectReviewStatus(feedback) {
  // Check first line or overall sentiment
  const firstLine = feedback.split('\n')[0].trim();
  
  for (const pattern of APPROVAL_PATTERNS) {
    if (pattern.test(firstLine)) {
      return 'APPROVED';
    }
  }
  
  for (const pattern of REVISION_PATTERNS) {
    if (pattern.test(firstLine)) {
      return 'NEEDS_REVISION';
    }
  }
  
  // Analyze content for implicit signals
  const lowerFeedback = feedback.toLowerCase();
  const hasIssues = /\b(issue|problem|error|bug|wrong|incorrect|missing)\b/.test(lowerFeedback);
  const hasApproval = /\b(good|great|excellent|perfect|works|correct)\b/.test(lowerFeedback);
  
  if (hasIssues && !hasApproval) return 'NEEDS_REVISION';
  if (hasApproval && !hasIssues) return 'APPROVED';
  
  return 'UNCLEAR'; // Requires user clarification
}
```

#### Pattern Tracking Logic

```javascript
function trackRecurringPatterns(taskHash, feedback, cycleData) {
  const patterns = extractPatterns(feedback);
  const tracking = cycleData.pattern_tracking;
  
  // Check if same issue appearing again
  for (const pattern of patterns) {
    if (tracking.recurring_patterns.includes(pattern)) {
      tracking.same_issue_count++;
      
      // Suggest alternative reviewer after threshold
      if (tracking.same_issue_count >= 2) {
        return suggestAlternativeReviewer(pattern, cycleData);
      }
    } else {
      tracking.recurring_patterns.push(pattern);
    }
  }
  
  return null; // No alternative needed yet
}

function suggestAlternativeReviewer(pattern, cycleData) {
  // Map common patterns to specialized reviewers
  const specializations = {
    'rate_limiting': 'the-architect',
    'authentication': 'the-security-engineer',
    'performance': 'the-site-reliability-engineer',
    'data_integrity': 'the-data-engineer',
    'api_design': 'the-architect'
  };
  
  // Find best alternative based on pattern
  for (const [key, agent] of Object.entries(specializations)) {
    if (pattern.includes(key) && agent !== cycleData.current_reviewer) {
      return {
        agent,
        reason: `Recurring ${key} issues detected, switching to specialist`
      };
    }
  }
  
  // Default to architect for persistent issues
  if (cycleData.current_reviewer !== 'the-architect') {
    return {
      agent: 'the-architect',
      reason: 'Escalating to architect after recurring issues'
    };
  }
  
  return null;
}
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
- Maximum 3 review cycles per task (configurable)
- After 3 cycles OR timeout (5 minutes default):
  - Escalate to user for decision
  - Apply default action if no response
- User options:
  - accept as-is: Continue with current state
  - manually fix: Pause for user intervention
  - skip: Move to next task with warning
  - abort: Stop entire process

### Timeout Handling
- Each review has a 5-minute timeout (configurable)
- On timeout:
  1. Check configured default_action
  2. Log timeout event with context
  3. Apply action (skip_with_warning by default)
  4. Notify user via configured method
  5. Continue process unless aborted

### Pattern-Based Recovery
- Track same_issue_count across cycles
- After 2 occurrences of same issue:
  1. Suggest alternative reviewer
  2. Log pattern for future optimization
  3. Update global failure rates
- High failure rate (>70%) triggers:
  1. Automatic reviewer switching
  2. Pattern analysis for root cause
  3. Recommendation for process improvement

### Recovery Strategy
- State persisted in `.the-startup/review-cycles.json`
- On restart:
  1. Load review cycle state from JSON
  2. Check PLAN.md for checkbox status
  3. Resume from last known state
  4. Recover in-progress reviews
- Context includes:
  - Previous review feedback
  - Cycle count (current/max)
  - Pattern history
  - Timeout status

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