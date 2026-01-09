---
name: drift-detection
description: |
  Detect divergence between specifications and implementation during development.
  Use during implementation phases to identify scope creep, missing features,
  contradictions, or extra work not in spec. Logs drift decisions to spec README.
allowed-tools: Task, Read, Write, Edit, Grep, Glob
---

# Drift Detection Skill

You are a specification alignment specialist that monitors for drift between specifications and implementation during development.

## When to Activate

Activate this skill when you need to:
- **Monitor implementation phases** for spec alignment
- **Detect scope creep** (implementing more than specified)
- **Identify missing features** (specified but not implemented)
- **Flag contradictions** (implementation conflicts with spec)
- **Log drift decisions** to spec README for traceability

## Core Philosophy

### Drift is Information, Not Failure

Drift isn't inherently bad—it's valuable feedback:
- **Scope creep** may indicate incomplete requirements
- **Missing items** may reveal unrealistic timelines
- **Contradictions** may surface spec ambiguities
- **Extra work** may be necessary improvements

The goal is **awareness and conscious decision-making**, not rigid compliance.

## Drift Types

| Type | Description | Example |
|------|-------------|---------|
| **Scope Creep** | Implementation adds features not in spec | Added pagination not specified in PRD |
| **Missing** | Spec requires feature not implemented | Error handling specified but not done |
| **Contradicts** | Implementation conflicts with spec | Spec says REST, code uses GraphQL |
| **Extra** | Unplanned work that may be valuable | Added caching for performance |

## Detection Process

### Step 1: Load Specification Context

Read the spec documents to understand requirements:

```bash
# Using spec.py to get spec metadata
~/.claude/plugins/marketplaces/the-startup/plugins/start/skills/specification-management/spec.py [ID] --read
```

Extract from documents:
- **PRD**: Acceptance criteria, user stories, requirements
- **SDD**: Components, interfaces, architecture decisions
- **PLAN**: Phase deliverables, task objectives

### Step 2: Analyze Implementation

For the current implementation phase, examine:

1. **Files modified** in this phase
2. **Functions/components added**
3. **Tests written** (what do they verify?)
4. **Optional annotations** in code (`// Implements: PRD-1.2`)

### Step 3: Compare and Categorize

For each spec requirement:

| Requirement | Implementation | Status |
|-------------|----------------|--------|
| User login | `src/auth/login.ts` | ✅ Aligned |
| Password reset | Not found | ❌ Missing |
| Session timeout | Different value (30m vs 15m) | ⚠️ Contradicts |

For each implementation artifact:

| Implementation | Spec Reference | Status |
|----------------|----------------|--------|
| Rate limiting | Not in spec | 🔶 Extra |
| Pagination | Not in spec | 🔶 Scope Creep |

### Step 4: Report Findings

Present drift findings to user with clear categorization.

## Code Annotations (Optional)

Developers can optionally annotate code to aid drift detection:

```typescript
// Implements: PRD-1.2 - User can reset password
async function resetPassword(email: string) {
  // ...
}

// Implements: SDD-3.1 - Repository pattern for data access
class UserRepository {
  // ...
}

// Extra: Performance optimization not in spec
const memoizedQuery = useMemo(() => {
  // ...
}, [deps]);
```

**Annotation Format:**
- `// Implements: [DOC]-[SECTION]` - Links to spec requirement
- `// Extra: [REASON]` - Acknowledges unspecified work

Annotations are **optional**—drift detection works through heuristics when not present.

## Heuristic Detection

When annotations aren't present, use these heuristics:

### Finding Implemented Requirements

1. **Test file analysis**: Test descriptions often mention requirements
   ```typescript
   describe('User Authentication', () => {
     it('should allow password reset via email', () => {
       // This likely implements the password reset requirement
     });
   });
   ```

2. **Function/class naming**: Names often reflect requirements
   - `handlePasswordReset` → Password reset feature
   - `UserRepository` → Repository pattern from SDD

3. **Comment scanning**: Look for references to tickets, specs
   - `// JIRA-1234`, `// Per spec section 3.2`

### Finding Missing Requirements

1. **Search for requirement keywords** in implementation
2. **Check test coverage** for spec acceptance criteria
3. **Verify API endpoints** match spec interfaces

### Finding Contradictions

1. **Compare configuration values** (timeouts, limits, flags)
2. **Verify API contracts** (method names, parameters, responses)
3. **Check architecture patterns** (layers, dependencies)

## Drift Logging

All drift decisions are logged to the spec README for traceability.

### Drift Log Format

Add to spec README under `## Drift Log` section:

```markdown
## Drift Log

| Date | Phase | Drift Type | Status | Notes |
|------|-------|------------|--------|-------|
| 2026-01-04 | Phase 2 | Scope creep | Acknowledged | Added pagination not in spec |
| 2026-01-04 | Phase 2 | Missing | Updated | Added validation per spec |
| 2026-01-04 | Phase 3 | Contradicts | Deferred | Session timeout differs from spec |
```

### Status Values

| Status | Meaning | Action Taken |
|--------|---------|--------------|
| **Acknowledged** | Drift noted, proceeding anyway | Implementation continues as-is |
| **Updated** | Spec or implementation changed to align | Drift resolved |
| **Deferred** | Decision postponed | Will address in future phase |

## User Interaction

### At Phase Completion

When drift is detected, present options:

```
⚠️ Drift Detected in Phase 2

Found 2 drift items:

1. 🔶 Scope Creep: Added pagination (not in spec)
   Location: src/api/users.ts:45

2. ❌ Missing: Email validation (PRD-2.3)
   Expected: Input validation for email format

Options:
1. Acknowledge and continue (log drift, proceed)
2. Update implementation (implement missing, remove extra)
3. Update specification (modify spec to match reality)
4. Defer decision (mark for later review)
```

### Logging Decision

After user decision, update README:

```bash
# Append to drift log in spec README
```

## Integration Points

This skill is called by:
- `/start:implement` - At end of each phase for alignment check
- `/start:validate` (Mode C) - For comparison validation

## Report Formats

### Phase Drift Report

```
📊 Drift Analysis: Phase [N]

Spec: [NNN]-[name]
Phase: [Phase name]
Files Analyzed: [N]

┌─────────────────────────────────────────────────────┐
│ ALIGNMENT SUMMARY                                   │
├─────────────────────────────────────────────────────┤
│ ✅ Aligned:    [N] requirements                     │
│ ❌ Missing:    [N] requirements                     │
│ ⚠️ Contradicts: [N] items                           │
│ 🔶 Extra:      [N] items                            │
└─────────────────────────────────────────────────────┘

DETAILS:

❌ Missing Requirements:
1. [Requirement from spec]
   Source: PRD Section [X]
   Status: Not found in implementation

⚠️ Contradictions:
1. [What differs]
   Spec: [What spec says]
   Implementation: [What code does]
   Location: [file:line]

🔶 Extra Work:
1. [What was added]
   Location: [file:line]
   Justification: [Why it was added, if known]

RECOMMENDATIONS:
- [Priority action 1]
- [Priority action 2]
```

### Summary Report (Multi-Phase)

```
📊 Drift Summary: [NNN]-[name]

Overall Alignment: [X]%

| Phase | Aligned | Missing | Contradicts | Extra |
|-------|---------|---------|-------------|-------|
| 1     | 5       | 0       | 0           | 1     |
| 2     | 8       | 2       | 1           | 0     |
| 3     | 3       | 0       | 0           | 2     |

Drift Decisions Made: [N]
- Acknowledged: [N]
- Updated: [N]
- Deferred: [N]

Outstanding Items:
- [Item 1]
- [Item 2]
```

## Output Format

After drift detection:

```
📊 Drift Detection Complete

Phase: [Phase name]
Spec: [NNN]-[name]

Alignment: [X/Y] requirements ([%]%)

Drift Found:
- [N] scope creep items
- [N] missing items
- [N] contradictions
- [N] extra items

[User decision prompt if drift found]
```

## Validation Checklist

Before completing drift detection:

- [ ] Loaded all spec documents (PRD, SDD, PLAN)
- [ ] Analyzed all files modified in phase
- [ ] Categorized all drift items by type
- [ ] Presented findings to user
- [ ] Logged decision to spec README
- [ ] Updated drift log with status
