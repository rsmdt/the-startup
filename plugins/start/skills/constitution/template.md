# Project Constitution

> Version: 1.0.0 | Last Updated: [DATE]
> Project Type: [NEEDS DISCOVERY: Explore codebase to determine - single-app | monorepo | library | cli]

## Security

[NEEDS DISCOVERY: Analyze codebase for security patterns]

Explore:
- Authentication mechanisms (JWT, sessions, OAuth)
- Secret handling (env vars, config files)
- Input validation patterns
- API security (CORS, rate limiting)

Generate rules for security patterns actually used in this project.

### No Hardcoded Secrets

```yaml
level: L1
pattern: "(api_key|apikey|secret|password|token|credential)\\s*[:=]\\s*['\"][^'\"]{8,}['\"]"
scope: "**/*.{ts,js,json,yaml,yml}"
exclude: "**/*.test.*, **/*.spec.*, **/*.example.*, .env.example"
message: Hardcoded secret detected. Use environment variables.
```

Secrets must never be committed to source control. Use environment variables or secret management.

### No Eval Usage

```yaml
level: L1
pattern: "\\beval\\s*\\("
scope: "src/**/*.{ts,js}"
message: eval() is prohibited for security reasons.
```

eval() enables code injection attacks. Use safer alternatives like JSON.parse().

## Architecture

[NEEDS DISCOVERY: Analyze codebase for architectural patterns]

Explore:
- Layer structure (presentation, business, data)
- Module boundaries (packages, apps, libs)
- Dependency patterns (injection, imports)
- API patterns (REST, GraphQL, RPC)

Generate rules that enforce the architecture patterns discovered.

## Code Quality

[NEEDS DISCOVERY: Analyze codebase for quality conventions]

Explore:
- Naming conventions (files, variables, functions)
- Import patterns (relative, absolute, aliases)
- Error handling patterns
- Logging conventions

Generate rules that enforce conventions already established.

### No Console Statements in Production

```yaml
level: L2
pattern: "console\\.(log|debug|info)"
scope: "src/**/*.{ts,js}"
exclude: "**/*.test.*, **/*.spec.*, src/utils/logger.ts"
message: Remove console statements. Use logger utility.
```

Console statements should be replaced with the structured logger for production observability.

## Testing

[NEEDS DISCOVERY: Analyze test setup and conventions]

Explore:
- Test framework(s) used
- Test file naming and location
- Coverage expectations
- Testing patterns (unit, integration, e2e)

Generate rules that align with existing test infrastructure.

### Test File Recommended

```yaml
level: L3
check: Every file in src/ should have corresponding .test.ts or .spec.ts
scope: "src/**/*.ts"
exclude: "src/**/*.d.ts, src/**/index.ts, src/**/*.test.ts"
message: Missing test file for this module.
```

Test coverage is encouraged but not enforced at the constitution level.

## [Project-Specific Category]

[NEEDS DISCOVERY: Based on codebase analysis, determine if additional categories are needed]

Consider:
- Domain-specific rules
- Framework-specific patterns (React hooks, FastAPI dependencies)
- Build/deployment conventions
- Documentation requirements

## Custom Rules

This section is for user additions that don't fit standard categories.

<!-- Users can add custom rules here -->
