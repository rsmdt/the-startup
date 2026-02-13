# Ambiguity Detection Reference

Techniques for detecting and scoring ambiguous language in specifications.

## Vague Language Patterns

| Pattern | Example | Recommendation |
|---------|---------|----------------|
| Hedge words | "should", "might", "could" | Use "must" or "will" |
| Vague quantifiers | "fast", "many", "various" | Specify metrics |
| Open-ended lists | "etc.", "and so on" | Enumerate all items |
| Undefined terms | "the system", "appropriate" | Define specifically |
| Passive voice | "errors are handled" | Specify who/what |
| Weak verbs | "support", "allow" | Use concrete actions |

## Ambiguity Score

```
ambiguity_score = vague_patterns / total_statements * 100

  0-5%:   ✅ Excellent clarity
  5-15%:  🟡 Acceptable
  15-25%: 🟠 Recommend clarification
  25%+:   🔴 High ambiguity
```

## Ambiguity Red Flags

- "should", "might", "could", "may"
- "fast", "slow", "many", "few"
- "etc.", "and so on", "..."
- "appropriate", "reasonable"
- "some", "several", "a few"
- "as needed", "when necessary"
- "properly", "correctly"

## Automated Detection

### Ambiguity Scan

```bash
grep -inE "(should|might|could|may|various|etc\.|and so on|appropriate|reasonable|fast|slow|many|few)" [file]
```

### Counting Script

```bash
# Count vague patterns
vague_count=$(grep -icE "(should|might|could|may|various|etc\.|appropriate|reasonable)" [file])

# Count total lines (rough statement count)
total=$(wc -l < [file])

# Calculate percentage
echo "Ambiguity: $((vague_count * 100 / total))%"
```

## Category-Specific Patterns

### Requirements Ambiguity

| Vague | Specific |
|-------|----------|
| "The system should be fast" | "Response time < 200ms p95" |
| "Handle many users" | "Support 10,000 concurrent users" |
| "User-friendly interface" | "WCAG 2.1 AA compliant, 5 clicks max to any feature" |
| "Secure authentication" | "OAuth 2.0 with JWT, 15-min token expiry" |

### Architecture Ambiguity

| Vague | Specific |
|-------|----------|
| "Scalable design" | "Horizontal scaling via K8s, stateless services" |
| "Proper error handling" | "Errors caught at service boundary, logged with correlation ID" |
| "Standard patterns" | "Repository pattern for data access, Service layer for business logic" |

### Implementation Ambiguity

| Vague | Specific |
|-------|----------|
| "Validate input" | "Check email format (RFC 5322), length 5-254 chars, sanitize HTML" |
| "Handle edge cases" | "Null user → 404, Empty list → empty array, Invalid ID → 400" |
| "Add appropriate logging" | "Log INFO for requests, WARN for retries, ERROR with stack trace" |

## Remediation Strategies

### For Requirements

1. **Add metrics**: Replace qualitative with quantitative
2. **Define boundaries**: Specify min/max/exact values
3. **List explicitly**: Replace "etc." with complete list
4. **Name actors**: Replace "the system" with specific component

### For Design

1. **Reference standards**: Link to design patterns, RFCs, specs
2. **Show examples**: Include code snippets or diagrams
3. **Define interfaces**: Specify method signatures, not just descriptions
4. **Enumerate options**: List all valid states/values

### For Implementation

1. **Write tests first**: Tests define unambiguous behavior
2. **Use types**: Let type system enforce constraints
3. **Add assertions**: Make implicit assumptions explicit
4. **Document edge cases**: Comment unusual handling

## Report Format

```
⚠️ Ambiguity Analysis

File: [path]
Score: [X]% ([level])

High-Priority (should → must):
- Line 23: "should validate" → "must validate"
- Line 45: "may include" → "includes" or "does not include"

Medium-Priority (vague quantifiers):
- Line 67: "fast response" → "< 200ms"
- Line 89: "many records" → "up to 10,000 records"

Low-Priority (style):
- Line 12: "etc." → list all items
- Line 34: "appropriate" → define criteria

Recommendations:
1. Address high-priority items before implementation
2. Clarify quantifiers with stakeholders
3. Replace open-ended lists with explicit enumerations
```
