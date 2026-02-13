# Agent Delegation Patterns

Preserved from the `task-delegation` skill. These patterns are useful reference for orchestrating parallel and sequential agent work.

## Parallel vs Sequential Decision Matrix

| Scenario | Dependencies | Shared State | Validation | File Paths | Recommendation |
|----------|--------------|--------------|------------|------------|----------------|
| Research tasks | None | Read-only | Independent | N/A | **PARALLEL** |
| Analysis tasks | None | Read-only | Independent | N/A | **PARALLEL** |
| Documentation | None | Unique paths | Independent | Unique | **PARALLEL** |
| Code creation | None | Unique files | Independent | Unique | **PARALLEL** |
| Build pipeline | Sequential | Shared files | Dependent | Same | **SEQUENTIAL** |
| File editing | None | Same file | Collision risk | Same | **SEQUENTIAL** |
| Dependent tasks | B needs A | Any | Dependent | Any | **SEQUENTIAL** |

### Parallel Execution Checklist

Before launching parallel agents, verify:
- [ ] No dependencies between tasks
- [ ] No shared state modifications
- [ ] Independent validation possible
- [ ] Unique file paths if creating files
- [ ] No resource contention

## File Creation Coordination / Collision Prevention

When multiple agents will create files:

1. Are file paths specified explicitly in each agent's OUTPUT?
2. Are all file paths unique (no two agents write same path)?
3. Do paths follow project conventions?
4. Are paths deterministic (not ambiguous)?

**If any check fails:** Adjust OUTPUT sections to prevent collisions.

### Path Assignment Strategies

**Strategy 1: Explicit Unique Paths**
```
Agent 1 OUTPUT: Create pattern at docs/patterns/authentication-flow.md
Agent 2 OUTPUT: Create interface at docs/interfaces/oauth-providers.md
Agent 3 OUTPUT: Create domain rule at docs/domain/user-permissions.md
```

**Strategy 2: Discovery-Based Paths**
```
Agent 1 OUTPUT: Test file at [DISCOVERED_LOCATION]/AuthService.test.ts
Agent 2 OUTPUT: Test file at [DISCOVERED_LOCATION]/UserService.test.ts
```

**Strategy 3: Hierarchical Paths**
```
Agent 1 OUTPUT: docs/patterns/backend/api-versioning.md
Agent 2 OUTPUT: docs/patterns/frontend/state-management.md
Agent 3 OUTPUT: docs/patterns/database/migration-strategy.md
```

## Scope Validation Framework

### Auto-Accept Criteria

Continue without user review when agent delivers:
- Vulnerability fixes, input validation additions
- Code clarity enhancements, documentation updates
- Exactly matches FOCUS requirements, respects EXCLUDE boundaries

### Requires User Review

Present to user when agent delivers:
- New external dependencies, database schema modifications
- Public API changes, design pattern changes
- Features beyond FOCUS (but valuable)

### Auto-Reject Criteria

Reject when agent delivers:
- Features not in requirements or explicitly in EXCLUDE list
- Breaking changes without migration path
- Missing required OUTPUT format or doesn't meet SUCCESS criteria
- "While I'm here" additions

## Failure Recovery & Retry Chain

When an agent fails, follow this escalation:

1. **Retry with refined prompt** - More specific FOCUS, more explicit EXCLUDE, better CONTEXT
2. **Try different specialist agent** - Different expertise angle, simpler task scope
3. **Break into smaller tasks** - Decompose further, sequential smaller steps
4. **Sequential instead of parallel** - Dependency might exist, coordination issue
5. **Handle directly (DIY)** - Task too specialized, agent limitation
6. **Escalate to user** - Present options, request guidance

### Retry Decision Tree

| Symptom | Likely Cause | Solution |
|---------|--------------|----------|
| Scope creep | FOCUS too vague | Refine FOCUS, expand EXCLUDE |
| Wrong approach | Wrong specialist | Try different agent type |
| Incomplete work | Task too complex | Break into smaller tasks |
| Blocked/stuck | Missing dependency | Check if should be sequential |
| Wrong output | OUTPUT unclear | Specify exact format/path |
| Quality issues | CONTEXT insufficient | Add more constraints/examples |

**Maximum retries: 3 attempts.** After 3 failed attempts, present to user and get guidance.
