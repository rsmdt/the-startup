---
name: code-quality-review
description: Systematic code review patterns, quality dimensions, anti-pattern detection, and constructive feedback techniques. Use when reviewing code changes, assessing codebase quality, identifying technical debt, or mentoring through reviews. Covers correctness, design, security, performance, and maintainability.
---

# Code Quality Review Methodology

Systematic patterns for reviewing code and providing constructive, actionable feedback that improves both code quality and developer skills.

## When to Activate

- Reviewing pull requests or merge requests
- Assessing overall codebase quality
- Identifying and prioritizing technical debt
- Mentoring developers through code review
- Establishing code review standards for teams
- Auditing code for security or compliance

## Review Dimensions

Every code review should evaluate these six dimensions:

### 1. Correctness

Does the code work as intended?

| Check | Questions |
|-------|-----------|
| Functionality | Does it solve the stated problem? |
| Edge Cases | Are boundary conditions handled? |
| Error Handling | Are failures gracefully managed? |
| Data Validation | Are inputs validated at boundaries? |
| Null Safety | Are null/undefined cases covered? |

### 2. Design

Is the code well-structured?

| Check | Questions |
|-------|-----------|
| Single Responsibility | Does each function/class do one thing? |
| Abstraction Level | Is complexity hidden appropriately? |
| Coupling | Are dependencies minimized? |
| Cohesion | Do related things stay together? |
| Extensibility | Can it be modified without major changes? |

### 3. Readability

Can others understand this code?

| Check | Questions |
|-------|-----------|
| Naming | Do names reveal intent? |
| Comments | Is the "why" explained, not the "what"? |
| Formatting | Is style consistent? |
| Complexity | Is cyclomatic complexity reasonable (<10)? |
| Flow | Is control flow straightforward? |

### 4. Security

Is the code secure?

| Check | Questions |
|-------|-----------|
| Input Validation | Are all inputs sanitized? |
| Authentication | Are auth checks present where needed? |
| Authorization | Are permissions verified? |
| Data Exposure | Is sensitive data protected? |
| Dependencies | Are there known vulnerabilities? |

### 5. Performance

Is the code efficient?

| Check | Questions |
|-------|-----------|
| Algorithmic | Is time complexity appropriate? |
| Memory | Are allocations reasonable? |
| I/O | Are database/network calls optimized? |
| Caching | Is caching used where beneficial? |
| Concurrency | Are race conditions avoided? |

### 6. Testability

Can this code be tested?

| Check | Questions |
|-------|-----------|
| Test Coverage | Are critical paths tested? |
| Test Quality | Do tests verify behavior, not implementation? |
| Mocking | Are external dependencies mockable? |
| Determinism | Are tests reliable and repeatable? |
| Edge Cases | Are boundary conditions tested? |

## Anti-Pattern Catalog

Common code smells and their remediation:

### Method-Level Anti-Patterns

| Anti-Pattern | Detection Signs | Remediation |
|--------------|-----------------|-------------|
| **Long Method** | >20 lines, multiple responsibilities | Extract Method |
| **Long Parameter List** | >3-4 parameters | Introduce Parameter Object |
| **Duplicate Code** | Copy-paste patterns | Extract Method, Template Method |
| **Complex Conditionals** | Nested if/else, switch statements | Decompose Conditional, Strategy Pattern |
| **Magic Numbers** | Hardcoded values without context | Extract Constant |
| **Dead Code** | Unreachable or unused code | Delete it |

### Class-Level Anti-Patterns

| Anti-Pattern | Detection Signs | Remediation |
|--------------|-----------------|-------------|
| **God Object** | >500 lines, many responsibilities | Extract Class |
| **Data Class** | Only getters/setters, no behavior | Move behavior to class |
| **Feature Envy** | Method uses another class's data extensively | Move Method |
| **Inappropriate Intimacy** | Classes know too much about each other | Move Method, Extract Class |
| **Refused Bequest** | Subclass doesn't use inherited behavior | Replace Inheritance with Delegation |
| **Lazy Class** | Does too little to justify existence | Inline Class |

### Architecture-Level Anti-Patterns

| Anti-Pattern | Detection Signs | Remediation |
|--------------|-----------------|-------------|
| **Circular Dependencies** | A depends on B depends on A | Dependency Inversion |
| **Shotgun Surgery** | One change requires many file edits | Move Method, Extract Class |
| **Leaky Abstraction** | Implementation details exposed | Encapsulate |
| **Premature Optimization** | Complex code for unproven performance | Simplify, measure first |
| **Over-Engineering** | Abstractions for hypothetical requirements | YAGNI - simplify |

## Review Prioritization

Focus review effort where it matters most:

### Priority 1: Critical (Must Fix)

- Security vulnerabilities (injection, auth bypass)
- Data loss or corruption risks
- Breaking changes to public APIs
- Production stability risks

### Priority 2: High (Should Fix)

- Logic errors affecting functionality
- Performance issues in hot paths
- Missing error handling for likely failures
- Violation of architectural principles

### Priority 3: Medium (Consider Fixing)

- Code duplication
- Missing tests for new code
- Naming that reduces clarity
- Overly complex conditionals

### Priority 4: Low (Nice to Have)

- Style inconsistencies
- Minor optimization opportunities
- Documentation improvements
- Refactoring suggestions

## Constructive Feedback Patterns

### The Feedback Formula

```
[Observation] + [Why it matters] + [Suggestion] + [Example if helpful]
```

### Good Feedback Examples

```markdown
# Instead of:
"This is wrong"

# Say:
"This query runs inside a loop (line 45), which could cause N+1
performance issues as the dataset grows. Consider using a batch
query before the loop:

```python
users = User.query.filter(User.id.in_(user_ids)).all()
user_map = {u.id: u for u in users}
```
"
```

```markdown
# Instead of:
"Use better names"

# Say:
"The variable `d` on line 23 would be clearer as `daysSinceLastLogin` -
it helps readers understand the business logic without tracing back
to the assignment."
```

### Feedback Tone Guide

| Avoid | Prefer |
|-------|--------|
| "You should..." | "Consider..." or "What about..." |
| "This is wrong" | "This might cause issues because..." |
| "Why didn't you..." | "Have you considered..." |
| "Obviously..." | "One approach is..." |
| "Always/Never do X" | "In this context, X would help because..." |

### Positive Observations

Include what's done well:

```markdown
"Nice use of the Strategy pattern here - it makes adding new
payment methods straightforward."

"Good error handling - the retry logic with exponential backoff
is exactly what we need for this flaky API."

"Clean separation of concerns between the validation and persistence logic."
```

## Review Checklists

### Quick Review Checklist (< 100 lines)

- [ ] Code compiles and tests pass
- [ ] Logic appears correct for stated purpose
- [ ] No obvious security issues
- [ ] Naming is clear
- [ ] No magic numbers or strings

### Standard Review Checklist (100-500 lines)

All of the above, plus:
- [ ] Design follows project patterns
- [ ] Error handling is appropriate
- [ ] Tests cover new functionality
- [ ] No significant duplication
- [ ] Performance is reasonable

### Deep Review Checklist (> 500 lines or critical)

All of the above, plus:
- [ ] Architecture aligns with system design
- [ ] Security implications considered
- [ ] Backward compatibility maintained
- [ ] Documentation updated
- [ ] Migration/rollback plan if needed

## Review Workflow

### Before Reviewing

1. Understand the context (ticket, discussion, requirements)
2. Check if CI passes (don't review failing code)
3. Estimate review complexity and allocate time

### During Review

1. First pass: Understand the overall change
2. Second pass: Check correctness and design
3. Third pass: Look for edge cases and security
4. Document findings as you go

### After Review

1. Summarize overall impression
2. Clearly indicate approval status
3. Distinguish blocking vs non-blocking feedback
4. Offer to discuss complex suggestions

## Review Metrics

Track review effectiveness:

| Metric | Target | What It Indicates |
|--------|--------|-------------------|
| Review Turnaround | < 24 hours | Team velocity |
| Comments per Review | 3-10 | Engagement level |
| Defects Found | Decreasing trend | Quality improvement |
| Review Time | < 60 min for typical PR | Right-sized changes |
| Approval Rate | 70-90% first submission | Clear standards |

## Anti-Patterns in Reviewing

Avoid these review behaviors:

| Anti-Pattern | Description | Better Approach |
|--------------|-------------|-----------------|
| **Nitpicking** | Focusing on style over substance | Use linters for style |
| **Drive-by Review** | Quick approval without depth | Allocate proper time |
| **Gatekeeping** | Blocking for personal preferences | Focus on objective criteria |
| **Ghost Review** | Approval without comments | Add at least one observation |
| **Review Bombing** | Overwhelming with comments | Prioritize and limit to top issues |
| **Delayed Review** | Letting PRs sit for days | Commit to turnaround time |

## References

- [Review Dimension Details](reference.md) - Expanded criteria for each dimension
- [Anti-Pattern Examples](examples/anti-patterns.md) - Code examples of each anti-pattern
