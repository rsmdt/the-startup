# Delegation Validation Protocol

## Post-Response Validation

After EVERY sub-agent response or batch of parallel responses, perform validation:

### Step 1: Drift Detection

Scan response for the following drift indicators:

- **Scope creep**: Features or functionality not in original request
- **Over-engineering**: Unnecessary complexity or abstraction
- **Assumption escalation**: Small assumptions becoming major features
- **Technology sprawl**: Introducing unspecified dependencies or tools
- **Feature bloat**: "Nice to have" or "while we're at it" additions
- **Premature optimization**: Performance improvements not requested
- **Architecture inflation**: More complex patterns than needed

### Step 2: Validation Check

Perform systematic validation:

```
🔍 Validating Response from [agent-name]
├─ Task adherence: [On-track/Drift detected]
├─ Scope check: [Within bounds/Exceeded]
├─ Complexity: [Appropriate/Over-engineered]
├─ Dependencies: [Specified only/New additions]
└─ Result: [PASS/MINOR DRIFT/MAJOR DRIFT]
```

### Step 3: Handle Validation Results

#### ✅ PASSED (no issues detected)
- Continue with workflow
- No user interruption needed
- Log validation success

#### 📝 MINOR DRIFT (reasonable additions)
Additions that might be beneficial but weren't explicitly requested:

```
📝 Minor additions detected:

The agent suggested these additions:
- [Addition 1]: [Brief description]
- [Addition 2]: [Brief description]

These seem reasonable and may improve the solution.
Include them? [Y/n]: _
```

Examples of minor drift:
- Better error handling
- Improved code structure
- Additional validation
- Helpful documentation
- Test coverage suggestions

#### ⚠️ MAJOR DRIFT (significant deviation)
Substantial changes to scope or approach:

```
⚠️ Significant drift detected:

Original scope: [what was requested]
Agent added: [what was added]
Impact: [consequences of the addition]

This represents a significant change to the requirements.

Options:
a) Accept additions (update scope)
b) Reject and revise (strict scope)
c) Partial accept (you specify what to keep)
d) Cancel and reassign task

Your choice [a/b/c/d]: _
```

Examples of major drift:
- New user-facing features
- Additional database tables
- External service integrations
- Authentication systems not requested
- New API endpoints beyond scope

### Step 4: Corrective Actions

If drift is rejected:

1. **Create revised context** with stricter boundaries:
   ```
   PREVIOUS ATTEMPT DRIFT:
   - Agent added [X] which was not requested
   
   STRICT REQUIREMENT:
   - ONLY implement [exactly what's needed]
   - Do NOT add [specific exclusion]
   
   EXCLUDE:
   - [Previous drift item 1]
   - [Previous drift item 2]
   ```

2. **Re-delegate** with clarified context
3. **Document** the drift for future reference

## Validation Rules

### Always Flag for Review:
- New user-facing features not in requirements
- Additional database schema changes
- New external dependencies or services
- Authentication/authorization modifications
- Performance optimizations (unless requested)
- Caching layers (unless requested)
- Additional API endpoints
- New configuration requirements

### Usually Accept Without Interruption:
- Error handling improvements
- Input validation enhancements
- Security best practices
- Code structure improvements
- Documentation additions
- Type safety improvements
- Test suggestions
- Logging additions

### Gray Areas (Use Judgment):
- Refactoring suggestions
- Alternative approaches
- Technology recommendations
- Architecture patterns

## Batch Validation

When multiple agents run in parallel:

1. **Validate each response individually**
2. **Check for conflicts** between responses
3. **Present consolidated validation** if issues found:
   ```
   🔍 Batch Validation Results:
   
   the-architect: ✅ PASSED
   the-developer: 📝 MINOR DRIFT - Added error handling
   the-data-engineer: ⚠️ MAJOR DRIFT - Added caching layer
   
   Address drift issues? [Y/n]: _
   ```
4. **Handle each drift type** appropriately
5. **Ensure consistency** across all responses

## Context Preservation

When re-delegating after validation failure:

```
=== REVISED DELEGATION ===

CONTEXT: [Original context]

PREVIOUS ATTEMPT ISSUE:
The previous response included [specific drift] which is outside scope.

STRICT REQUIREMENTS:
1. ONLY implement [specific requirement]
2. Do NOT add [drift item 1]
3. Do NOT include [drift item 2]

EXCLUDE explicitly:
- [All items from previous drift]
- [Any related features to avoid]
```

## Validation Metrics

Track validation patterns for improvement:

- Frequency of drift by agent type
- Common drift patterns
- Success rate after re-delegation
- User acceptance rate of minor drift

## Quick Reference Checklist

For each response, ask:

1. ✓ Does it solve the requested task?
2. ✓ Does it stay within boundaries?
3. ✓ Are all additions justified?
4. ✓ Is complexity appropriate?
5. ✓ Are dependencies minimal?
6. ✓ Would the user be surprised by anything?

If any answer is "no", validation action is required.