# Agent Delegation Skill Reference

Complete reference for advanced delegation patterns, edge cases, and optimization strategies.

## Advanced Decomposition Patterns

### Multi-Level Decomposition

For very complex tasks, decompose in layers:

**Layer 1: High-Level Activities**
```
Task: Build e-commerce checkout flow

Activities:
1. Frontend checkout interface
2. Backend payment processing
3. Order fulfillment system
4. Email notifications
```

**Layer 2: Sub-Activity Decomposition**

Take Activity 2 and decompose further:
```
Activity: Backend payment processing

Sub-activities:
2.1 Stripe API integration
2.2 Payment validation logic
2.3 Transaction database schema
2.4 Refund handling
```

**Execution Strategy:**
- Layer 1: Mixed (some parallel, some sequential)
- Layer 2: Decompose only when agent starts Activity 2
- Don't decompose all layers upfront (overwhelming)

### Dependency Graph Decomposition

For complex dependency chains:

```
Task: Deploy new microservice

Activity Map:
A: Write service code
B: Write unit tests (depends on A)
C: Create Docker image (depends on A)
D: Write integration tests (depends on A, C)
E: Deploy to staging (depends on B, C, D)
F: Run smoke tests (depends on E)
G: Deploy to production (depends on F)

Execution Groups:
Group 1: A (sequential)
Group 2: B, C (parallel after A)
Group 3: D (sequential after Group 2)
Group 4: E (sequential after Group 3)
Group 5: F (sequential after E)
Group 6: G (sequential after F)
```

**Pattern:** Identify critical path, parallelize where possible.

### Expertise-Based Decomposition

When multiple domains are involved:

```
Task: Add real-time chat feature

Decompose by expertise:
1. UI/UX design (design expertise)
2. Frontend component (React expertise)
3. WebSocket server (Backend expertise)
4. Message persistence (Database expertise)
5. Security review (Security expertise)
6. Performance testing (Performance expertise)

Execution:
- Phase 1: Activity 1 (sequential)
- Phase 2: Activities 2-4 (parallel, informed by Activity 1)
- Phase 3: Activities 5-6 (parallel review after Phase 2)
```

## Advanced Parallel Patterns

### Fan-Out/Fan-In Pattern

Parallel expansion followed by sequential consolidation:

```
        Start
          ↓
      [Activity 1]
          ↓
    ┌─────┼─────┐
    ↓     ↓     ↓
  [A2]  [A3]  [A4]  ← Fan-out (parallel)
    ↓     ↓     ↓
    └─────┼─────┘
          ↓
     [Synthesize]  ← Fan-in (sequential)
          ↓
        Done
```

**Example:** Competitive analysis
- Fan-out: Analyze competitors A, B, C in parallel
- Fan-in: Synthesize findings into unified strategy

### Pipeline Pattern

Sequential groups where each group can be parallel:

```
Stage 1: Research (parallel within stage)
- Market research
- Competitive analysis
- User interviews
    ↓
Stage 2: Design (parallel within stage)
- UI mockups
- API design
- Database schema
    ↓
Stage 3: Implementation (parallel within stage)
- Frontend build
- Backend build
- Database setup
```

**Pattern:** Stages are sequential, activities within each stage are parallel.

### MapReduce Pattern

Parallel processing with aggregation:

```
Map Phase (parallel):
Agent 1: Process dataset chunk 1
Agent 2: Process dataset chunk 2
Agent 3: Process dataset chunk 3
Agent 4: Process dataset chunk 4

Reduce Phase (sequential):
Aggregate all results into final output
```

**Example:** Code analysis across modules
- Map: Each agent analyzes one module
- Reduce: Aggregate findings into project-wide report

## Advanced Template Patterns

### Context Accumulation Pattern

For sequential tasks, accumulate context:

```
Agent 1:
CONTEXT: Fresh start, no prior context
OUTPUT: Result A

Agent 2:
CONTEXT:
  - Prior results: [Result A from Agent 1]
  - Build on: [Specific insights from A]
OUTPUT: Result B (informed by A)

Agent 3:
CONTEXT:
  - Prior results: [Result A, Result B]
  - Conflicts to resolve: [Any conflicts between A and B]
  - Build on: [Insights from both]
OUTPUT: Result C (synthesizes A and B)
```

**Key:** Each agent gets relevant prior outputs, not everything.

### Constraint Propagation Pattern

Cascade constraints through dependent tasks:

```
Agent 1 (Schema Design):
SUCCESS:
  - Uses PostgreSQL (project standard)
  - Follows naming: snake_case tables
  - All tables have created_at, updated_at

Agent 2 (API Implementation, depends on Agent 1):
CONTEXT:
  - Database constraints from Agent 1:
    * PostgreSQL only
    * snake_case table names
    * created_at/updated_at in all tables
  - Must match schema exactly
```

**Pattern:** SUCCESS criteria from earlier tasks become CONTEXT constraints for later ones.

### Specification Reference Pattern

For implementation tasks, reference specs explicitly:

```
FOCUS: Implement user registration endpoint

CONTEXT:
  - PRD Section 3.1.2: User registration requirements
  - SDD Section 4.2: API endpoint specifications
  - SDD Section 5.3: Database schema for users table
  - PLAN Phase 2, Task 3: Implementation checklist

SDD_REQUIREMENTS:
  - Endpoint: POST /api/auth/register
  - Request body: { email, password, name }
  - Response: { user_id, token }
  - Validation: Email format, password strength (8+ chars)
  - Security: Bcrypt hashing (cost 12)

SPECIFICATION_CHECK: Must match SDD Section 4.2 exactly
```

**Pattern:** Explicit spec references prevent context drift.

## File Coordination Advanced Strategies

### Timestamp-Based Uniqueness

When paths might collide, add timestamps:

```
Agent 1 OUTPUT: logs/analysis-${TIMESTAMP}.md
Agent 2 OUTPUT: logs/research-${TIMESTAMP}.md
Agent 3 OUTPUT: logs/synthesis-${TIMESTAMP}.md

where TIMESTAMP = ISO 8601 format
```

**Result:** No collisions even if agents run simultaneously.

### Directory Hierarchy Assignment

Assign each agent a subdirectory:

```
Agent 1 OUTPUT: results/agent-1/findings.md
Agent 2 OUTPUT: results/agent-2/findings.md
Agent 3 OUTPUT: results/agent-3/findings.md
```

**Result:** Each agent owns a directory, filenames can repeat.

### Atomic File Creation Pattern

For critical files, ensure atomic creation:

```
OUTPUT: Create file at exact path: docs/patterns/auth.md
  - If file exists, FAIL and report (don't overwrite)
  - Use atomic write (temp file + rename)
  - Verify write succeeded before marking complete
```

**Pattern:** Prevents race conditions and corruption.

### Merge Strategy Pattern

When multiple agents create similar content:

```
Strategy: Sequential merge

Agent 1: Create base document
Agent 2: Read base, add section 2
Agent 3: Read base + section 2, add section 3

Each agent:
DISCOVERY_FIRST: Read current state of document
FOCUS: Add my section without modifying others
OUTPUT: Updated document with my section added
```

**Pattern:** Sequential additions to shared document.

## Scope Validation Advanced Patterns

### Severity-Based Acceptance

Categorize scope creep by severity:

**Minor (Auto-accept):**
- Variable name improvements
- Comment additions
- Whitespace formatting
- Import organization

**Medium (Review):**
- Small refactors related to task
- Additional error handling
- Logging additions
- Documentation updates

**Major (Reject):**
- New features
- Architecture changes
- Dependency additions
- Breaking changes

### Value-Based Exception Handling

Sometimes scope creep is valuable:

```
Agent delivered:
✅ Required: Authentication endpoint
⚠️ Extra: Rate limiting on endpoint (not requested)

Analysis:
- Extra work: Rate limiting
- In EXCLUDE? No (not explicitly excluded)
- Valuable? Yes (security best practice)
- Risky? No (standard pattern)
- Increases scope? Minimally

Decision: 🟡 ACCEPT with note
  "Agent proactively added rate limiting for security.
   Aligns with best practices, accepting this valuable addition."
```

**Pattern:** Auto-accept valuable, low-risk extras that align with project goals.

### Specification Drift Detection

For implement tasks, detect drift from specs:

```
Validation:
1. Check FOCUS matches PLAN task description
2. Check implementation matches SDD requirements
3. Check business logic matches PRD rules

Drift detected if:
- Implementation differs from SDD design
- Business rules differ from PRD
- Features not in PLAN added

Report:
📊 Specification Alignment: 85%
✅ Aligned: [aspects that match]
⚠️ Deviations: [aspects that differ]
🔴 Critical drift: [major misalignments]
```

## Retry Strategy Advanced Patterns

### Progressive Refinement

Refine template progressively across retries:

**Attempt 1 (Failed - too vague):**
```
FOCUS: Add caching
```

**Attempt 2 (Failed - still ambiguous):**
```
FOCUS: Add Redis caching for API responses
EXCLUDE: Don't cache user-specific data
```

**Attempt 3 (Success - specific enough):**
```
FOCUS: Add Redis caching for public API endpoints
  - Cache GET requests only
  - TTL: 5 minutes
  - Key format: api:endpoint:params:hash
  - Invalidate on POST/PUT/DELETE to same resource

EXCLUDE:
  - Don't cache authenticated user requests
  - Don't cache admin endpoints
  - Don't implement cache warming
  - Don't add Redis cluster setup (single node for now)

CONTEXT:
  - Redis already configured: localhost:6379
  - Use ioredis client
  - Follow caching pattern: docs/patterns/caching-strategy.md
```

**Pattern:** Each retry adds specificity based on previous failure.

### Agent Type Rotation

If specialist fails, try different angle:

```
Attempt 1: Backend specialist
  - Focused on technical implementation
  - Failed: Too technical, missed user experience

Attempt 2: UX specialist
  - Focused on user flows
  - Failed: Too high-level, missed technical constraints

Attempt 3: Product specialist
  - Balanced user needs with technical reality
  - Success: Right blend of perspectives
```

**Pattern:** Rotate expertise angle based on failure mode.

### Scope Reduction Strategy

If task too complex, reduce scope progressively:

```
Attempt 1 (Failed - too much):
FOCUS: Build complete authentication system
  - Registration, login, logout, password reset
  - OAuth integration
  - Two-factor authentication

Attempt 2 (Failed - still complex):
FOCUS: Build basic authentication
  - Registration, login, logout

Attempt 3 (Success - minimal):
FOCUS: Build login endpoint only
  - POST /auth/login
  - Email + password validation
  - Return JWT token
```

**Pattern:** Reduce scope until agent succeeds, then expand incrementally.

## Edge Cases and Solutions

### Edge Case 1: Circular Dependencies

**Problem:** Agent A needs Agent B's output, Agent B needs Agent A's output

**Detection:**
```
Activity A depends on B
Activity B depends on A
→ Circular dependency detected
```

**Solutions:**

1. **Break the cycle:**
   ```
   Original:
   A (needs B) ↔ B (needs A)

   Refactored:
   C (shared foundation) → A (builds on C) → B (builds on A)
   ```

2. **Iterative approach:**
   ```
   Round 1: A (with assumptions about B)
   Round 2: B (using Round 1 A)
   Round 3: A (refined with actual B)
   ```

3. **Merge activities:**
   ```
   Single agent handles both A and B together
   (They're too coupled to separate)
   ```

### Edge Case 2: Dynamic Dependencies

**Problem:** Don't know dependencies until runtime

**Example:**
```
Task: Analyze codebase

Don't know which modules exist until discovery
Can't plan parallel structure upfront
```

**Solution - Two-phase approach:**

**Phase 1: Discovery (sequential)**
```
Agent 1: Discover project structure
OUTPUT: List of modules

Result: [moduleA, moduleB, moduleC, moduleD]
```

**Phase 2: Analysis (parallel, dynamic)**
```
For each module in result:
  Launch analysis agent

Agent A: Analyze moduleA
Agent B: Analyze moduleB
Agent C: Analyze moduleC
Agent D: Analyze moduleD
```

**Pattern:** Sequential discovery, dynamic parallel execution.

### Edge Case 3: Partial Agent Availability

**Problem:** Some specialist agents unavailable

**Example:**
```
Planned:
- Frontend specialist (available)
- Backend specialist (available)
- DevOps specialist (NOT AVAILABLE)
```

**Solution - Fallback delegation:**

```
If specialist unavailable:
1. Try broader domain agent (general-purpose)
2. Try sequential breakdown (smaller tasks)
3. Handle directly if simple enough
4. Escalate to user if critical
```

**Example execution:**
```
DevOps work:
  Attempt 1: DevOps specialist → UNAVAILABLE
  Attempt 2: Backend specialist with DevOps context → SUCCESS
  Reasoning: Backend specialist has some DevOps overlap
```

### Edge Case 4: Agent Response Conflicts

**Problem:** Parallel agents return conflicting recommendations

**Example:**
```
Agent 1 (Security): "Use bcrypt with cost 14 (maximum security)"
Agent 2 (Performance): "Use bcrypt with cost 10 (reasonable security, better performance)"
```

**Solution - Conflict resolution:**

**1. Present to user:**
```
⚠️ Agent Conflict Detected

Topic: Bcrypt cost factor
Agent 1 (Security): Cost 14 (maximize security)
Agent 2 (Performance): Cost 10 (balance security/performance)

Trade-off:
- Cost 14: ~200ms hashing time, highest security
- Cost 10: ~50ms hashing time, strong security

Recommendation needed: Which priority matters more?
```

**2. Specification arbitration:**
```
Check specs:
- SDD Section 5.2: "Use bcrypt cost factor 12"
→ Use specification value (12)
→ Both agents adjusted to match spec
```

**3. Synthesis agent:**
```
Launch Agent 3 (Architect):
FOCUS: Resolve conflict between security and performance recommendations
CONTEXT:
  - Security recommendation: cost 14
  - Performance recommendation: cost 10
  - Trade-offs: [details]
OUTPUT: Final recommendation with reasoning
```

### Edge Case 5: Resource Constraints

**Problem:** Can't launch all parallel agents simultaneously (rate limits, memory, etc.)

**Solution - Batched parallel execution:**

```
Activities: [A1, A2, A3, A4, A5, A6, A7, A8]
Constraint: Maximum 3 parallel agents

Execution:
Batch 1: A1, A2, A3 (parallel)
  → Wait for completion
Batch 2: A4, A5, A6 (parallel)
  → Wait for completion
Batch 3: A7, A8 (parallel)
  → Complete
```

**Pattern:** Maintain parallelism within constraints.

## Performance Optimization

### Minimize Context Size

Don't pass everything to every agent:

❌ **Bad - Full context:**
```
CONTEXT:
  - Entire PRD (50 pages)
  - Entire SDD (40 pages)
  - All prior agent outputs (30 pages)
```

✅ **Good - Relevant context:**
```
CONTEXT:
  - PRD Section 3.2 (User authentication requirements)
  - SDD Section 4.1 (API endpoint design)
  - Prior output: Authentication flow diagram from Agent 1
```

**Pattern:** Extract only relevant portions, reference docs by section.

### Parallel Batching

Group related parallel tasks:

```
Don't:
Launch 20 individual research agents

Do:
Launch 4 research agents, each handles 5 topics
```

**Benefits:**
- Fewer coordination overhead
- Better context utilization
- Faster overall completion

### Early Termination

Build termination conditions into templates:

```
TERMINATION:
  - Completed successfully
  - Blocked by missing dependency X
  - Information not publicly available
  - Maximum 3 attempts reached
  - ERROR: [specific error conditions]
```

**Pattern:** Let agents fail fast instead of hanging.

## Integration with Other Skills

### With Documentation Skill

When agents discover patterns:

```
Agent completes task
  ↓
Discovers reusable pattern
  ↓
Documentation skill activates
  ↓
Pattern documented in docs/patterns/
  ↓
Reported back to orchestrator
```

**Coordination:** Agent-delegation creates prompts, documentation skill handles pattern storage.

### With Specification Review Skill

For implementation tasks:

```
Agent completes implementation
  ↓
Agent-delegation validates scope
  ↓
Specification review skill validates against PRD/SDD
  ↓
Both validations pass → Complete
```

**Coordination:** Agent-delegation handles scope, spec-review handles alignment.

### With Quality Gates Skill

At phase boundaries:

```
Phase completes
  ↓
Agent-delegation confirms all tasks done
  ↓
Quality gates skill runs DOD checks
  ↓
Both pass → Proceed to next phase
```

**Coordination:** Agent-delegation manages execution, quality-gates validates quality.

## Debugging Failed Delegations

### Symptom: Agents ignore EXCLUDE

**Diagnosis:**
- EXCLUDE too vague
- Agent sees value in excluded work
- Conflict between FOCUS and EXCLUDE

**Fix:**
```
Before:
EXCLUDE: Don't add extra features

After:
EXCLUDE: Do not add these specific features:
  - OAuth integration (separate task)
  - Password reset flow (separate task)
  - Two-factor authentication (not in scope)
  Any feature not explicitly in FOCUS is out of scope.
```

### Symptom: Parallel agents conflict

**Diagnosis:**
- Hidden shared state
- File path collision
- Dependency not identified

**Fix:**
```
Review parallel safety checklist:
- Independent tasks? → Check dependencies again
- Unique file paths? → Verify OUTPUT sections
- No shared state? → Identify what's shared

If any fail → Make sequential or coordinate better
```

### Symptom: Sequential too slow

**Diagnosis:**
- False dependencies
- Over-cautious sequencing
- Could be parallel with coordination

**Fix:**
```
Re-analyze dependencies:
- Must Task B use Task A's output? → True dependency
- Could Task B assume Task A's approach? → False dependency

If false dependency:
  → Make parallel with coordinated assumptions
```

### Symptom: Template too complex

**Diagnosis:**
- Too many constraints
- Context overload
- Agent confused by detail

**Fix:**
```
Simplify:
1. Keep FOCUS to essentials
2. Move details to CONTEXT
3. Provide examples instead of rules

Before (overwhelming):
FOCUS: [20 lines of detailed requirements]

After (simplified):
FOCUS: [2 lines of core task]
CONTEXT: [Details and constraints]
```

## Best Practices Summary

1. **Decompose by activities**, not roles
2. **Parallel by default**, sequential only when necessary
3. **Explicit FOCUS/EXCLUDE**, no ambiguity
4. **Unique file paths**, verify before launching
5. **Minimal context**, only relevant information
6. **Auto-accept safe changes**, review architectural ones
7. **Maximum 3 retries**, then escalate
8. **Early termination**, fail fast when blocked
9. **Validate scope**, check FOCUS/EXCLUDE adherence
10. **Document patterns**, activate documentation skill when discovered

## Common Patterns Quick Reference

| Pattern | When to Use | Structure |
|---------|-------------|-----------|
| Fan-Out/Fan-In | Parallel research → synthesis | Parallel → Sequential |
| Pipeline | Stages with parallel within | Sequential stages, parallel tasks |
| MapReduce | Large dataset processing | Parallel map → Sequential reduce |
| Progressive Refinement | Retry with more detail | Retry N adds specificity |
| Batched Parallel | Resource constraints | Groups of parallel tasks |
| Context Accumulation | Sequential with learning | Each task gets prior outputs |
| Constraint Propagation | Dependent implementations | SUCCESS → next CONTEXT |
| Specification Reference | Implementation tasks | Explicit PRD/SDD references |

---

This reference covers advanced scenarios beyond the main skill. Load this when dealing with complex coordination, optimization, or edge cases.
