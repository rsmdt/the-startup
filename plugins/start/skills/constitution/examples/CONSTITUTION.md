# Project Constitution

> Version: 1.0.0 | Last Updated: 2026-01-04
> Project Type: single-app (Next.js with TypeScript)

## Security

### No Hardcoded Secrets

```yaml
level: L1
pattern: "(api_key|apikey|secret|password|token|credential)\\s*[:=]\\s*['\"][^'\"]{8,}['\"]"
scope: "**/*.{ts,js,json,yaml,yml}"
exclude: "**/*.test.*, **/*.spec.*, **/*.example.*, .env.example"
message: Hardcoded secret detected. Use environment variables.
```

Secrets must never be committed to source control. Use environment variables via `process.env` or a secret management solution.

### No Eval Usage

```yaml
level: L1
pattern: "\\beval\\s*\\("
scope: "src/**/*.{ts,js}"
message: eval() is prohibited for security reasons.
```

eval() enables code injection attacks. Use safer alternatives like JSON.parse() for data parsing.

### No innerHTML Assignment

```yaml
level: L1
pattern: "\\.(innerHTML|outerHTML)\\s*="
scope: "src/**/*.{ts,tsx,js,jsx}"
exclude: "**/*.test.*"
message: Direct innerHTML assignment risks XSS. Use React's dangerouslySetInnerHTML with sanitization or DOM methods.
```

Direct DOM manipulation with user content can lead to cross-site scripting vulnerabilities.

## Architecture

### Repository Pattern Required

```yaml
level: L1
check: Database queries (Prisma, TypeORM, Knex, raw SQL) only in files matching *Repository.ts or *Repository.js
scope: "src/**/*.{ts,js}"
exclude: "**/repositories/**, **/prisma/**"
message: Direct database call outside repository layer.
```

All database operations must go through repository classes to maintain separation of concerns and enable testing.

### No Direct API Calls in Components

```yaml
level: L1
check: HTTP calls (fetch, axios) only in files under services/ or api/ or hooks/
scope: "src/components/**"
message: API calls must go through service layer.
```

Components should be pure UI; data fetching belongs in the service layer or custom hooks.

### No Barrel Exports

```yaml
level: L1
pattern: "export \\* from"
scope: "src/**/*.ts"
exclude: "src/index.ts"
message: Barrel exports prohibited. Import from specific files.
```

Barrel exports cause circular dependency issues and make tree-shaking ineffective. Always import from specific module files.

## Code Quality

### No Console Statements in Production

```yaml
level: L2
pattern: "console\\.(log|debug|info)"
scope: "src/**/*.{ts,js}"
exclude: "**/*.test.*, **/*.spec.*, src/utils/logger.ts"
message: Remove console statements. Use logger utility.
```

Console statements should be replaced with the structured logger for production observability.

### Functions Under 25 Lines

```yaml
level: L3
check: Functions should not exceed 25 lines
scope: "src/**/*.{ts,js}"
message: Function too long. Consider extracting into smaller functions.
```

Smaller functions improve readability and testability, but this is a guideline not a strict requirement.

### No Magic Numbers

```yaml
level: L3
pattern: "[^0-9]\\b[2-9]\\d{2,}\\b(?![0-9])"
scope: "src/**/*.{ts,js}"
exclude: "**/*.test.*, **/*.spec.*, **/constants/**"
message: Magic number detected. Consider extracting to named constant.
```

Numbers other than 0, 1, or small constants should be named for clarity.

## Testing

### No .only in Committed Tests

```yaml
level: L1
pattern: "\\.(only|skip)\\s*\\("
scope: "**/*.test.*, **/*.spec.*"
message: Remove .only/.skip before committing tests.
```

Focused tests prevent the full suite from running and may hide failures.

### No console.log in Tests

```yaml
level: L2
pattern: "console\\.log"
scope: "**/*.test.ts, **/*.spec.ts"
message: Use assertions instead of console.log in tests.
```

Tests should use assertions; console.log often indicates incomplete test development.

### Test File Recommended

```yaml
level: L3
check: Every file in src/ should have corresponding .test.ts or .spec.ts
scope: "src/**/*.ts"
exclude: "src/**/*.d.ts, src/**/index.ts, src/**/*.test.ts, src/**/*.spec.ts"
message: Missing test file for this module.
```

Test coverage is encouraged. Critical business logic should have corresponding tests.

## React Patterns

### No Direct State Mutation

```yaml
level: L1
pattern: "this\\.state\\.[a-zA-Z]+\\s*="
scope: "src/**/*.{tsx,jsx}"
exclude: "**/*.test.*"
message: Use setState instead of direct state mutation.
```

Direct state mutation bypasses React's rendering lifecycle and causes bugs.

### Key Prop Required in Lists

```yaml
level: L2
check: Array.map rendering JSX must include key prop with unique value (not array index)
scope: "src/**/*.{tsx,jsx}"
message: Missing or invalid key prop in list rendering. Use unique identifiers.
```

Missing or index-based keys cause inefficient re-renders and bugs with stateful children.

### No useEffect Dependencies Warning

```yaml
level: L2
check: useEffect hooks should have explicit dependency arrays (not missing the second argument)
scope: "src/**/*.{tsx,jsx}"
message: useEffect without dependency array runs on every render.
```

Missing dependency arrays often indicate unintended behavior.

## Custom Rules

<!-- Project-specific rules can be added here -->

### No Direct Environment Access in Components

```yaml
level: L2
pattern: "process\\.env\\."
scope: "src/components/**/*.{tsx,jsx}"
message: Access environment variables through config module, not directly in components.
```

Environment variables should be centralized for validation and default handling.
