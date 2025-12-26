---
name: pattern-detection
description: Identify existing codebase patterns (naming conventions, architectural patterns, testing patterns) to maintain consistency. Use when generating code, reviewing changes, or understanding established practices. Ensures new code aligns with project conventions.
---

# Pattern Recognition

## When to Use

- Before writing new code to ensure consistency with existing patterns
- During code review to verify alignment with established conventions
- When onboarding to understand project-specific practices
- Before refactoring to preserve intentional design decisions

## Core Methodology

### Pattern Discovery Process

1. **Survey representative files**: Read 3-5 files of the type you will create or modify
2. **Identify recurring structures**: Note repeated patterns in naming, organization, imports
3. **Verify intentionality**: Check if patterns are documented or consistently applied
4. **Apply discovered patterns**: Use the same conventions in new code

### Priority Order for Pattern Sources

1. **Existing code in the same module/feature** - Most authoritative
2. **Project style guides or CONTRIBUTING.md** - Explicit documentation
3. **Test files** - Often reveal expected patterns and naming
4. **Similar files in adjacent modules** - Fallback when no direct examples exist

## Naming Convention Recognition

### File Naming Patterns

Detect and follow the project's file naming style:

| Pattern | Example | Common In |
|---------|---------|-----------|
| kebab-case | `user-profile.ts` | Node.js, Vue, Angular |
| PascalCase | `UserProfile.tsx` | React components |
| snake_case | `user_profile.py` | Python |
| camelCase | `userProfile.js` | Legacy JS, Java |

### Function/Method Naming

Identify the project's verb conventions:

- **get** vs **fetch** vs **retrieve** for data access
- **create** vs **add** vs **new** for creation
- **update** vs **set** vs **modify** for mutations
- **delete** vs **remove** vs **destroy** for deletion
- **is/has/can/should** prefixes for booleans

### Variable Naming

Detect pluralization and specificity patterns:

- Singular vs plural for collections (`user` vs `users` vs `userList`)
- Hungarian notation presence (`strName`, `iCount`)
- Private member indicators (`_private`, `#private`, `mPrivate`)

## Architectural Pattern Recognition

### Layer Identification

Recognize how the codebase separates concerns:

```
COMMON LAYERING PATTERNS:
- MVC: controllers/, models/, views/
- Clean Architecture: domain/, application/, infrastructure/
- Hexagonal: core/, adapters/, ports/
- Feature-based: features/auth/, features/billing/
- Type-based: components/, services/, utils/
```

### Dependency Direction

Identify import patterns that reveal architecture:

- Which modules import from which (dependency flow)
- Shared vs feature-specific code boundaries
- Framework code vs application code separation

### State Management Patterns

Recognize how state flows through the application:

- Global stores (Redux, Vuex, MobX patterns)
- React Context usage patterns
- Service layer patterns for backend state
- Event-driven vs request-response patterns

## Testing Pattern Recognition

### Test Organization

Identify how tests are structured:

| Pattern | Structure | Example |
|---------|-----------|---------|
| Co-located | `src/user.ts`, `src/user.test.ts` | Common in modern JS/TS |
| Mirror tree | `src/user.ts`, `tests/src/user.test.ts` | Traditional, Java-style |
| Feature-based | `src/user/`, `src/user/__tests__/` | React, organized features |

### Test Naming Conventions

Detect the project's test description style:

- **BDD style**: `it('should return user when found')`
- **Descriptive**: `test('getUser returns user when id exists')`
- **Function-focused**: `test_get_user_returns_user_when_found`

### Test Structure Patterns

Recognize Arrange-Act-Assert or Given-When-Then patterns:

- Setup block conventions (beforeEach, fixtures, factories)
- Assertion style (expect vs assert vs should)
- Mock/stub patterns (jest.mock vs sinon vs manual)

## Code Organization Patterns

### Import Organization

Identify import ordering and grouping:

```
COMMON IMPORT PATTERNS:
1. External packages first, internal modules second
2. Grouped by type (React, libraries, local)
3. Alphabetized within groups
4. Absolute imports vs relative imports preference
```

### Export Patterns

Recognize module boundary conventions:

- Default exports vs named exports preference
- Barrel files (index.ts re-exports) presence
- Public API definition patterns

### Comment and Documentation Patterns

Identify documentation conventions:

- JSDoc/TSDoc presence and style
- Inline comment frequency and style
- README conventions per module/feature

## Best Practices

- **Follow existing patterns even if imperfect** - Consistency trumps personal preference
- **Document deviations explicitly** - When breaking patterns intentionally, explain why
- **Pattern changes require migration** - Dont introduce new patterns without updating existing code
- **Check tests for patterns too** - Test code often reveals expected conventions
- **Prefer explicit over implicit** - When patterns are unclear, ask or document assumptions

## Anti-Patterns to Avoid

- Mixing naming conventions in the same codebase
- Introducing new architectural patterns without team consensus
- Assuming patterns from other projects apply here
- Ignoring test patterns when writing implementation
- Creating "special" files that dont follow established structure

## References

- `examples/common-patterns.md` - Concrete examples of pattern recognition in action
