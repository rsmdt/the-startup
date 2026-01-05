# Drift Detection Reference

Advanced techniques for detecting and managing specification drift.

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

## Annotation Formats

### Standard Annotations

```typescript
// Implements: [Document]-[Section/ID]
// Implements: PRD-AC-1 - User login acceptance criteria
// Implements: SDD-3.2 - Repository pattern
// Implements: PLAN-2.1 - Registration endpoint

// Extra: [Brief justification]
// Extra: Added rate limiting for security
// Extra: Caching layer for performance

// Deferred: [What and why]
// Deferred: Admin dashboard - moved to Phase 3

// TODO: [Spec reference if applicable]
// TODO: PRD-4.2 - Email notifications not yet implemented
```

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

## Heuristic Patterns

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

**Extraction:**
```bash
# Extract test descriptions
grep -r "describe\|it\|test(" tests/ --include="*.test.ts"
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

## README Drift Log Examples

### Initial Log Entry

```markdown
## Drift Log

| Date | Phase | Drift Type | Status | Notes |
|------|-------|------------|--------|-------|
| 2026-01-04 | Phase 1 | Extra | Acknowledged | Added health check endpoint |
```

### After Multiple Phases

```markdown
## Drift Log

| Date | Phase | Drift Type | Status | Notes |
|------|-------|------------|--------|-------|
| 2026-01-04 | Phase 1 | Extra | Acknowledged | Added health check endpoint |
| 2026-01-04 | Phase 2 | Scope creep | Acknowledged | Added pagination to user list |
| 2026-01-04 | Phase 2 | Missing | Updated | Added email validation per PRD-2.3 |
| 2026-01-04 | Phase 3 | Contradicts | Deferred | Session timeout 30m vs spec 15m |
```

### With Resolution Notes

```markdown
## Drift Log

| Date | Phase | Drift Type | Status | Notes |
|------|-------|------------|--------|-------|
| 2026-01-04 | Phase 2 | Missing | Updated | Added email validation per PRD-2.3 |
| 2026-01-05 | Phase 3 | Contradicts | Updated | Changed timeout to 15m per spec |
| 2026-01-05 | Phase 3 | Scope creep | Updated | Spec updated to include pagination |
```

## Integration with Other Skills

### With constitution-validation

After drift detection, check if drifted code complies with constitution:

```
1. Detect drift
2. If new code added (scope creep, extra)
3. Validate new code against constitution
4. Report both drift AND constitution findings
```

### With specification-management

Use spec.py to:
- Read spec metadata
- Locate spec documents
- Understand spec structure

### With implementation-verification

Combine with verification for comprehensive check:
- Drift detection: Is it aligned with spec?
- Implementation verification: Is it technically correct?

## Performance Considerations

### Efficient Searching

1. **Narrow search scope** to modified files
2. **Use ripgrep** for faster text search
3. **Cache spec parsing** across phases
4. **Batch file reads** when possible

### Incremental Detection

Track what was already verified:
- Store verification state between phases
- Only re-verify modified areas
- Mark files as "verified" after check

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

### Large Codebases

- Use targeted searching based on PLAN tasks
- Focus on recently modified files
- Sample verification for comprehensive requirements
- Prioritize high-severity requirements
