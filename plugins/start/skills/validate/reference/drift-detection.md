# Drift Detection Reference

Techniques for detecting and managing divergence between specifications and implementation.

## Core Philosophy

**Drift is Information, Not Failure.**

Drift isn't inherently bad—it's valuable feedback:
- **Scope creep** may indicate incomplete requirements
- **Missing items** may reveal unrealistic timelines
- **Contradictions** may surface spec ambiguities
- **Extra work** may be necessary improvements

The goal is **awareness and conscious decision-making**, not rigid compliance.

---

## Drift Types

| Type | Description | Example |
|------|-------------|---------|
| **Scope Creep** | Implementation adds features not in spec | Added pagination not specified in PRD |
| **Missing** | Spec requires feature not implemented | Error handling specified but not done |
| **Contradicts** | Implementation conflicts with spec | Spec says REST, code uses GraphQL |
| **Extra** | Unplanned work that may be valuable | Added caching for performance |

---

## Detection Strategies

### Strategy 1: Acceptance Criteria Mapping

Map PRD acceptance criteria to implementation evidence.

**Process:**
1. Extract acceptance criteria from PRD
2. Search implementation for matching behavior
3. Verify through test assertions

**Example PRD Criteria:**
```markdown
### AC-1: User Login
Given a registered user
When they enter valid credentials
Then they receive an access token
```

**Search Patterns:**
```bash
# Search for login implementation
grep -r "login\|authenticate\|access.token" src/

# Search for test coverage
grep -r "should.*login\|given.*registered.*user" tests/
```

### Strategy 2: Interface Contract Validation

Compare SDD interfaces with actual implementation.

**SDD Interface:**
```typescript
// From SDD
interface UserService {
  login(email: string, password: string): Promise<AuthToken>;
  logout(token: string): Promise<void>;
  resetPassword(email: string): Promise<void>;
}
```

**Validation:**
```bash
# Find actual interface
grep -A 20 "interface UserService\|class UserService" src/

# Compare method signatures
```

### Strategy 3: Architecture Pattern Verification

Verify SDD architectural decisions in code.

**SDD Decision:**
> ADR-3: Use repository pattern for data access

**Verification:**
1. Check for `*Repository` classes
2. Verify no direct database calls outside repositories
3. Confirm dependency injection of repositories

```bash
# Find repositories
find src -name "*Repository*"

# Check for direct DB calls outside repositories
grep -r "prisma\.\|db\.\|query(" src/ --include="*.ts" | grep -v Repository
```

### Strategy 4: PLAN Task Completion

Verify PLAN tasks against implementation.

**PLAN Task:**
```markdown
- [ ] Implement user registration endpoint
  - Route: POST /api/users
  - Validation: Email format, password strength
  - Response: Created user object
```

**Verification:**
```bash
# Find route
grep -r "POST.*\/api\/users\|router.post.*users" src/

# Find validation
grep -r "email.*valid\|password.*strength" src/
```

---

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

### Annotation Scanning

```bash
# Find all implementation references
grep -r "// Implements:" src/

# Find all extra work
grep -r "// Extra:" src/

# Find deferred items
grep -r "// Deferred:" src/

# Find spec-related TODOs
grep -r "// TODO:.*PRD\|SDD\|PLAN" src/
```

---

## Heuristic Detection

When annotations aren't present, use these heuristics:

### Naming Convention Analysis

| Naming Pattern | Likely Spec Source |
|----------------|-------------------|
| `handle[Action]` | PRD user action |
| `validate[Entity]` | PRD validation rule |
| `[Entity]Repository` | SDD repository pattern |
| `[Entity]Service` | SDD service layer |
| `use[Feature]` | SDD/PRD feature hook |

### Test Description Analysis

```typescript
// Test descriptions often reflect requirements
describe('Authentication', () => {
  it('should return 401 for invalid credentials');  // Security requirement
  it('should rate limit failed attempts');          // Security requirement
  it('should issue JWT token on success');          // Token specification
});
```

### Import Analysis

Imports reveal architectural patterns:

```typescript
// Repository usage indicates data layer separation
import { UserRepository } from '@/repositories/user';

// Service usage indicates business layer
import { AuthService } from '@/services/auth';

// Direct ORM usage may indicate pattern violation
import { prisma } from '@/lib/prisma';  // Should only be in repositories
```

---

## Contradiction Detection

### Configuration Mismatches

**Common areas:**
- Timeout values
- Rate limits
- Pagination sizes
- Cache durations
- Retry counts

**Detection:**
```bash
# Find configuration values
grep -r "timeout\|limit\|size\|duration\|retry" src/config/

# Compare against spec values
```

### API Contract Mismatches

**Check:**
- HTTP methods (GET vs POST)
- Route paths (/users vs /user)
- Request/response shapes
- Status codes
- Error formats

### Type Mismatches

**Compare:**
- Spec data types vs implementation types
- Optional vs required fields
- Enum values
- Validation constraints

---

## Drift Severity Assessment

### High Severity (Address Immediately)

- Security requirement missing
- Core functionality not implemented
- Breaking API contract change
- Data integrity issue

### Medium Severity (Address Before Release)

- Non-critical feature missing
- Performance requirement unmet
- Documentation mismatch
- Test coverage gap

### Low Severity (Track for Future)

- Style/preference differences
- Nice-to-have features
- Optimization opportunities
- Documentation improvements

---

## Drift Logging

All drift decisions should be logged to the spec README for traceability.

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

---

## User Decision Workflow

When drift is detected, present options:

```
⚠️ Drift Detected

Found [N] drift items:

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

---

## Report Format

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

---

## Troubleshooting

### False Positives (Detecting Drift That Isn't)

- **Naming mismatch**: Spec says "user" but code says "account"
  - Solution: Build synonym mapping
- **Abstraction level**: Spec is high-level, code is detailed
  - Solution: Consider implementation details as aligned
- **Reorganization**: Same feature in different location
  - Solution: Search more broadly before flagging missing

### False Negatives (Missing Real Drift)

- **Subtle differences**: Close but not exact
  - Solution: Fuzzy matching on requirements
- **Hidden in complexity**: Feature buried in large functions
  - Solution: Deeper code analysis
- **Different terminology**: Spec and code use different terms
  - Solution: Keyword expansion
