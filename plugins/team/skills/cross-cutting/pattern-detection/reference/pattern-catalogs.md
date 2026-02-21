# Pattern Catalogs

Detailed catalogs of common patterns to detect across naming, architecture, testing, and code organization.

---

## Naming Convention Catalogs

### File Naming Patterns

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

## Architectural Pattern Catalogs

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

## Testing Pattern Catalogs

### Test Organization

| Pattern | Structure | Example |
|---------|-----------|---------|
| Co-located | `src/user.ts`, `src/user.test.ts` | Common in modern JS/TS |
| Mirror tree | `src/user.ts`, `tests/src/user.test.ts` | Traditional, Java-style |
| Feature-based | `src/user/`, `src/user/__tests__/` | React, organized features |

### Test Naming Conventions

- **BDD style**: `it('should return user when found')`
- **Descriptive**: `test('getUser returns user when id exists')`
- **Function-focused**: `test_get_user_returns_user_when_found`

### Test Structure Patterns

Recognize Arrange-Act-Assert or Given-When-Then patterns:

- Setup block conventions (beforeEach, fixtures, factories)
- Assertion style (expect vs assert vs should)
- Mock/stub patterns (jest.mock vs sinon vs manual)

## Code Organization Catalogs

### Import Organization

```
COMMON IMPORT PATTERNS:
1. External packages first, internal modules second
2. Grouped by type (React, libraries, local)
3. Alphabetized within groups
4. Absolute imports vs relative imports preference
```

### Export Patterns

- Default exports vs named exports preference
- Barrel files (index.ts re-exports) presence
- Public API definition patterns

### Comment and Documentation Patterns

- JSDoc/TSDoc presence and style
- Inline comment frequency and style
- README conventions per module/feature

## Error Handling Patterns

- Custom error classes vs generic errors
- Error hierarchy structure
- HTTP status code conventions
- Error message formatting
- Context inclusion in error objects

## Configuration Patterns

- Config split by domain vs monolithic
- Environment variable naming (SCREAMING_SNAKE_CASE)
- Default value strategies
- Type coercion approaches
- Secret management patterns
