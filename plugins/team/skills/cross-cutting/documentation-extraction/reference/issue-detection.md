# Identifying Documentation Issues

Strategies for detecting outdated, conflicting, and missing documentation.

---

## Outdated Documentation

Signals that documentation may be stale:

- **Version Mismatches**: Docs reference v1.x, code is v2.x
- **Missing Features**: Code has capabilities not in docs
- **Dead Links**: References to moved or deleted resources
- **Deprecated Patterns**: Docs use patterns code has abandoned
- **Date Indicators**: "Last updated 2 years ago" on active project

**Verification Steps**:
```
1. Check doc commit history vs code commit history
2. Compare documented API against actual code signatures
3. Run documented examples - do they work?
4. Search code for terms used in docs - are they present?
```

## Conflicting Documentation

When multiple docs disagree:

1. **Identify the conflict explicitly**: Quote both sources
2. **Check timestamps**: Newer usually wins
3. **Check authority**: Official > community, code > docs
4. **Test behavior**: What does the system actually do?
5. **Document the resolution**: Note which source was correct

**Resolution Priority**:
```
1. Actual system behavior (empirical truth)
2. Most recent official documentation
3. Code comments and inline documentation
4. External/community documentation
5. Older official documentation
```

## Missing Documentation

Recognize documentation gaps:

- **Undocumented Endpoints**: Routes exist in code but not docs
- **Hidden Configuration**: Env vars used but not listed
- **Implicit Requirements**: Dependencies not in requirements file
- **Tribal Knowledge**: Processes that exist only in team memory

**Gap Documentation Template**:
```markdown
## Documentation Gap: [Topic]

**Discovered**: [Date]
**Location**: [Where this should be documented]
**Current State**: [What exists now]
**Required Information**: [What's missing]
**Source of Truth**: [Where to get correct info]
```

## Cross-Referencing Documentation with Code

### Tracing Requirements to Implementation
```
1. Extract requirement ID or description
2. Search codebase for requirement reference
3. If not found, search for key domain terms
4. Locate implementation and verify behavior
5. Document mapping: Requirement -> File:Line
```

### Validating API Documentation
```
1. Find endpoint in documentation
2. Locate route definition in code
3. Compare: method, path, parameters
4. Trace to handler implementation
5. Verify response shape matches docs
```

### Configuration Value Tracing
```
1. Identify configuration key in docs
2. Search for key in codebase
3. Find where value is read/consumed
4. Trace through to actual usage
5. Verify documented behavior matches code
```
