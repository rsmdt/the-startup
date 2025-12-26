# Documentation Skill Reference

Complete reference for the documentation skill including advanced patterns, edge cases, and detailed protocols.

## Advanced Categorization Rules

### Gray Areas and Edge Cases

#### When Business and Technical Overlap

**Authentication Example:**
- `docs/domain/user-roles.md` - WHO can access WHAT (business rule)
- `docs/patterns/authentication-flow.md` - HOW authentication works (technical)
- `docs/interfaces/oauth-providers.md` - EXTERNAL services used (integration)

**Guideline:** If it affects WHAT users can do → domain. If it affects HOW we build it → patterns.

#### When Pattern Becomes Interface

**Caching Example:**
- Local in-memory caching → `docs/patterns/caching-strategy.md`
- Redis/Memcached integration → `docs/interfaces/redis-cache.md`

**Guideline:** Self-contained code patterns → patterns. External service dependencies → interfaces.

#### When Multiple Categories Apply

**Payment Processing Example:**
Could span all three:
- `docs/domain/payment-rules.md` - Refund policies, pricing rules
- `docs/patterns/payment-processing.md` - Internal payment handling
- `docs/interfaces/stripe-api.md` - Stripe integration specifics

**Guideline:** Create separate documents for each perspective. Cross-reference heavily.

## Naming Conventions

### Pattern: `[noun]-[noun/verb].md`

**Good Examples:**
- `error-handling.md`
- `database-migrations.md`
- `api-versioning.md`
- `event-sourcing.md`

**Avoid:**
- Single words: `cache.md`, `auth.md`
- Abbreviations: `db-mig.md`, `err-hdl.md`
- Generic terms: `utilities.md`, `helpers.md`

### Interface: `[service-name]-[integration-type].md`

**Good Examples:**
- `stripe-payments.md`
- `sendgrid-webhooks.md`
- `github-api.md`
- `aws-s3-storage.md`

**Avoid:**
- Generic: `payment-gateway.md` (which one?)
- Vague: `email.md` (what about email?)
- Tech-only: `rest-api.md` (which service?)

### Domain: `[entity/concept]-[aspect].md`

**Good Examples:**
- `user-permissions.md`
- `order-workflow.md`
- `inventory-tracking.md`
- `pricing-rules.md`

**Avoid:**
- Implementation details: `user-table.md` (that's technical)
- Generic: `rules.md` (which rules?)
- Too broad: `business-logic.md` (everything?)

## Update vs Create Decision Matrix

| Scenario | Existing Doc | Action |
|----------|--------------|--------|
| New payment provider | `stripe-payments.md` exists | **Create** `paypal-payments.md` (different service) |
| Additional caching layer | `caching-strategy.md` exists | **Update** existing (same pattern, new details) |
| New user role type | `user-permissions.md` exists | **Update** existing (extends same rule set) |
| Different auth method | `jwt-authentication.md` exists | **Create** `oauth-authentication.md` (different approach) |
| API version change | `github-api.md` exists | **Update** existing (same service, evolved) |
| New business constraint | `order-workflow.md` exists | **Update** if related, **Create** if distinct |

**Guiding Principle:** Same topic/service = update. Different topic/service = create new.

## Template Usage Guidelines

### Pattern Template
Use for:
- Architectural decisions (MVC, microservices, event-driven)
- Code organization patterns (repository, factory, singleton)
- Data handling approaches (caching, validation, serialization)
- Testing strategies (unit, integration, e2e)

### Interface Template
Use for:
- Third-party API integrations
- Webhook implementations
- External service authentication
- Data exchange protocols
- Partner system integrations

### Domain Template
Use for:
- Business rules and constraints
- User permission systems
- Workflow state machines
- Validation requirements
- Domain entity behaviors

## Deduplication Techniques

### Technique 1: Keyword Search
```bash
# Search filenames
find docs -type f -name "*.md" | grep -i keyword

# Search content
grep -ri "search term" docs/
```

### Technique 2: Category Listing
```bash
# List all patterns
ls docs/patterns/

# List all interfaces
ls docs/interfaces/

# List all domain docs
ls docs/domain/
```

### Technique 3: Content Scanning
```bash
# Show first 5 lines of each file
find docs/patterns -name "*.md" -exec head -5 {} \; -print

# Search for specific concept
grep -l "authentication" docs/**/*.md
```

### Technique 4: Related Term Mapping

For a new document about "caching":
- Check for: cache, caching, cached, memoization, storage
- Check categories: patterns (implementation), interfaces (Redis/Memcached)
- Read related files before deciding

## Merge vs Separate Guidelines

### Merge When:
- Same category and closely related topic
- Information enhances without confusing
- Single cohesive narrative possible
- Total length stays under 500 lines

**Example:** Merging "JWT tokens" into existing `authentication-flow.md`

### Keep Separate When:
- Different approaches to same problem
- Distinct services/technologies
- Would make document unfocused
- Exceeds reasonable length

**Example:** `jwt-authentication.md` and `oauth-authentication.md` as separate files

## Cross-Reference Patterns

### Within Same Category
```markdown
## Related Patterns
- [Repository Pattern](./repository-pattern.md) - Data access layer
- [Service Layer](./service-layer.md) - Business logic organization
```

### Across Categories
```markdown
## Related Documentation
- **Domain:** [User Permissions](../domain/user-permissions.md) - Authorization rules
- **Patterns:** [Authentication Flow](../patterns/authentication-flow.md) - Technical implementation
- **Interfaces:** [OAuth Providers](../interfaces/oauth-providers.md) - External auth services
```

### To Specifications
```markdown
## Implementations
- [User Authentication](../specs/001-user-auth/SDD.md) - Technical specification
- [OAuth Integration](../specs/015-oauth/PRD.md) - Product requirements
```

## Version Management

### When Patterns Evolve

**Approach 1: Update in Place**
- Add "Version History" section
- Document what changed and when
- Keep current approach primary

**Approach 2: Separate Documents**
- `authentication-v1.md` (legacy)
- `authentication-v2.md` (current)
- Clear migration path documented

**Guideline:** Update in place unless breaking change makes old version still relevant for existing code.

### Deprecation

When a pattern/interface is superseded:

```markdown
# Old Authentication Pattern

> **⚠️ DEPRECATED:** This pattern is no longer recommended.
> See [New Authentication Flow](./authentication-flow.md) for current approach.
>
> This document is maintained for reference by legacy code in modules X, Y, Z.

[Original content preserved...]
```

## Quality Standards

### Completeness Checklist
- [ ] Title clearly states what is documented
- [ ] Context explains when/why this applies
- [ ] Examples show real usage
- [ ] Edge cases are covered
- [ ] Related docs are linked
- [ ] Code snippets use real project conventions

### Clarity Checklist
- [ ] New team member could understand it
- [ ] Technical terms are explained
- [ ] Assumptions are stated explicitly
- [ ] Steps are in logical order
- [ ] Diagrams included for complex flows (if applicable)

### Maintainability Checklist
- [ ] Searchable filename
- [ ] Correct category
- [ ] No duplicate content
- [ ] Cross-references are bidirectional
- [ ] Version history if evolved

## Common Mistakes to Avoid

### ❌ Mistake 1: Creating Without Checking
**Problem:** Duplicate documentation proliferates
**Solution:** Always search first - multiple ways (grep, find, ls)

### ❌ Mistake 2: Wrong Category
**Problem:** Business rules in patterns/, technical details in domain/
**Solution:** Ask "Is this about WHAT (domain) or HOW (patterns)?"

### ❌ Mistake 3: Too Generic Names
**Problem:** Can't find documentation later
**Solution:** Full descriptive names, not abbreviations

### ❌ Mistake 4: No Cross-References
**Problem:** Related knowledge stays siloed
**Solution:** Link liberally between related docs

### ❌ Mistake 5: Template Ignored
**Problem:** Inconsistent structure makes scanning hard
**Solution:** Follow templates for consistency

### ❌ Mistake 6: No Examples
**Problem:** Abstract descriptions don't help
**Solution:** Include real code snippets and scenarios

## Edge Case Handling

### What if Nothing Fits the Categories?

**Option 1:** Expand categories (rare, think hard first)
**Option 2:** Create `docs/architecture/` for cross-cutting concerns
**Option 3:** Add to specification docs if feature-specific

**Example:** ADRs (Architecture Decision Records) might warrant `docs/decisions/`

### What if It's Too Small to Document?

**Guideline:** If it's reusable or non-obvious, document it.

**Too small:**
- "We use camelCase" (coding standard, not pattern)
- "API returns JSON" (obvious, not worth documenting)

**Worth documenting:**
- "We use optimistic locking for inventory" (non-obvious pattern)
- "Rate limiting uses token bucket algorithm" (specific approach)

### What if It's Extremely Specific?

**Guideline:** Very feature-specific logic goes in specs, not shared docs.

**Spec-level:**
- `specs/023-checkout/SDD.md` - Checkout flow specifics

**Shared docs:**
- `docs/patterns/state-machines.md` - Reusable state machine pattern
- `docs/domain/order-workflow.md` - General order rules

## Performance Considerations

### Keep Docs Focused
- Single file shouldn't exceed 1000 lines
- Split large topics into multiple focused docs
- Use cross-references instead of duplicating

### Optimize for Searchability
- Use keywords in filename
- Include synonyms in content
- Add tags/topics section at top

### Progressive Detail
```markdown
# Caching Strategy

Quick overview: We use Redis for session and API response caching.

## Details
[Detailed implementation...]

## Advanced Configuration
[Complex edge cases...]
```

## Integration with Specifications

### During Analysis (`/start:analyze`)
Documentation skill captures discovered patterns:
- Code analysis reveals patterns → Document in `docs/patterns/`
- Business rules discovered → Document in `docs/domain/`
- External APIs found → Document in `docs/interfaces/`

### During Specification (`/start:specify`)
- PRD/SDD references existing documentation
- New patterns discovered → Document them
- Specifications live in `docs/specs/`, reference shared docs

### During Implementation (`/start:implement`)
- Implementation follows documented patterns
- Deviations discovered → Update documentation
- New patterns emerge → Document for reuse

## Automation Support

### Pre-documentation Checks
Automate the search process:
```bash
# Check if topic exists
./scripts/check-doc-exists.sh "authentication"

# List related docs
./scripts/find-related-docs.sh "payment"
```

### Post-documentation Validation
```bash
# Verify no duplicates
./scripts/validate-docs.sh

# Check cross-references
./scripts/check-links.sh
```

## Summary

The documentation skill ensures:
1. **No duplication** - Always check before creating
2. **Correct categorization** - Business vs Technical vs External
3. **Discoverability** - Descriptive names and cross-references
4. **Consistency** - Template-based structure
5. **Maintainability** - Clear, complete, and up-to-date

When in doubt, ask:
- Does related documentation already exist?
- Which category fits best?
- What name would I search for?
- What template applies?
- How does this connect to other knowledge?
