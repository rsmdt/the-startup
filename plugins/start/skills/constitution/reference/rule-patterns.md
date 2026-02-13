# Constitution Validation Reference

Advanced patterns, rule types, and examples for constitution authoring.

## Rule Types

### Pattern Rules

Pattern rules use regex to match violations in source code. These are deterministic and fast.

**When to use:**
- Text patterns that can be matched literally
- Syntax violations (forbidden imports, banned functions)
- Secret detection (API keys, passwords)

**Regex Tips:**
- Escape special characters: `\.` for literal dot
- Use `\s*` for flexible whitespace
- Use `\b` for word boundaries
- Escape backslashes in YAML: `\\b` for `\b`

**Example - Detecting Barrel Exports:**

```yaml
level: L1
pattern: "export \\* from"
scope: "src/**/*.ts"
exclude: "src/index.ts"
message: Barrel exports prohibited. Import from specific files.
```

**Example - Detecting Direct DOM Manipulation:**

```yaml
level: L1
pattern: "\\.(innerHTML|outerHTML)\\s*="
scope: "src/**/*.{ts,tsx,js,jsx}"
exclude: "**/*.test.*"
message: Direct innerHTML assignment risks XSS. Use framework's DOM methods.
```

### Check Rules

Check rules use semantic descriptions that the LLM interprets. These are flexible but non-deterministic.

**When to use:**
- Architectural patterns that require understanding context
- Rules that span multiple lines or files
- Semantic concepts (like "database calls only in repositories")

**Writing Good Check Descriptions:**
- Be specific about what constitutes a violation
- Include examples of valid and invalid patterns
- Mention file locations or naming conventions to check

**Example - Repository Pattern:**

```yaml
level: L1
check: Database queries (Prisma, TypeORM, Knex, raw SQL) only in files matching *Repository.ts or *Repository.js
scope: "src/**/*.{ts,js}"
exclude: "**/repositories/**"
message: Direct database call outside repository layer.
```

**Example - Hook Rules:**

```yaml
level: L2
check: React hooks (useState, useEffect, useContext, custom use* hooks) only called in functional components or custom hooks
scope: "src/**/*.{tsx,jsx}"
message: Hooks must be called at the top level of functional components.
```

## Scope Patterns

The `scope` field uses glob patterns to determine which files to check.

### Common Patterns

| Pattern | Matches |
|---------|---------|
| `**/*.ts` | All TypeScript files |
| `src/**/*.ts` | TypeScript files in src/ |
| `**/*.{ts,js}` | TypeScript and JavaScript |
| `packages/*/src/**` | All packages' src folders |
| `apps/web/**` | Only the web app |
| `!**/*.test.ts` | Exclude test files (use in `exclude`) |

### Monorepo Scoping

For monorepos, scope to specific packages:

```yaml
# Web package only
scope: "packages/web/src/**/*.{ts,tsx}"

# All packages
scope: "packages/*/src/**/*.ts"

# Shared libraries
scope: "libs/*/src/**/*.ts"

# Apps only
scope: "apps/*/src/**/*.ts"
```

## Exclude Patterns

The `exclude` field is comma-separated globs for files to skip.

**Common Excludes:**

```yaml
exclude: "**/*.test.*, **/*.spec.*"              # Test files
exclude: "**/__tests__/**, **/__mocks__/**"      # Test directories
exclude: "**/node_modules/**"                    # Dependencies
exclude: "**/*.d.ts"                             # Type definitions
exclude: "**/*.example.*, .env.example"          # Example files
exclude: "src/generated/**"                      # Generated code
```

## Category ID Prefixes

When parsing rules, IDs are auto-generated from category:

| Category | Prefix | Example |
|----------|--------|---------|
| Security | SEC | SEC-001 |
| Architecture | ARCH | ARCH-001 |
| Code Quality | QUAL | QUAL-001 |
| Testing | TEST | TEST-001 |
| Custom | CUST | CUST-001 |
| [Custom Name] | First 4 letters uppercase | PERF-001 |

## Common Rule Patterns

### Security

```yaml
# SQL Injection Detection
level: L1
pattern: "\\.(query|execute|raw)\\s*\\([^)]*\\$\\{|\\+\\s*['\"]"
scope: "**/*.{ts,js}"
message: Potential SQL injection. Use parameterized queries.

# No Sensitive Data in Logs
level: L2
pattern: "console\\.(log|info|warn|error)\\([^)]*password|secret|token|key"
scope: "src/**/*.{ts,js}"
message: Sensitive data may be logged. Remove or redact.
```

### Architecture

```yaml
# No Cross-Package Relative Imports
level: L2
check: Imports between packages must use package name, not relative path
scope: "packages/*/src/**"
message: Cross-package import must use package name, not relative path.

# Service Layer Boundaries
level: L1
check: HTTP calls (fetch, axios, got) only in files under services/ or api/
scope: "src/components/**"
message: API calls must go through service layer.
```

### Code Quality

```yaml
# No TODO in Main Branch
level: L3
pattern: "TODO|FIXME|XXX|HACK"
scope: "src/**/*.{ts,js}"
message: Unresolved TODO marker. Complete or track in issue.

# Consistent Naming
level: L2
check: React component files must use PascalCase naming
scope: "src/components/**/*.tsx"
message: Component file should use PascalCase naming.
```

### Testing

```yaml
# No console.log in Tests (should use test assertions)
level: L2
pattern: "console\\.log"
scope: "**/*.test.ts, **/*.spec.ts"
message: Use assertions instead of console.log in tests.

# No .only in Committed Tests
level: L1
pattern: "\\.(only|skip)\\s*\\("
scope: "**/*.test.*, **/*.spec.*"
message: Remove .only/.skip before committing tests.
```

### React Specific

```yaml
# No Direct State Mutation
level: L1
pattern: "this\\.state\\s*\\."
scope: "src/**/*.{tsx,jsx}"
exclude: "**/*.test.*"
message: Use setState instead of direct state mutation.

# Key Prop in Lists
level: L2
check: Array.map rendering JSX must include key prop
scope: "src/**/*.{tsx,jsx}"
message: Missing key prop in list rendering.
```

### Node.js Specific

```yaml
# No Sync Methods
level: L2
pattern: "\\.(readFileSync|writeFileSync|existsSync|mkdirSync)"
scope: "src/**/*.ts"
exclude: "scripts/**, cli/**"
message: Use async file operations in application code.

# Require Error Handling for Promises
level: L2
check: Async functions should have try-catch or .catch for error handling
scope: "src/**/*.ts"
message: Unhandled promise rejection possible. Add error handling.
```

## Inheritance and Overrides

### Monorepo Pattern

In monorepos, place CONSTITUTION.md at root. Rules apply to all packages via scope:

```markdown
## Architecture

### Package Boundaries

```yaml
level: L2
check: Imports between packages must use package name
scope: "packages/*/src/**"
message: Cross-package import violation.
```

### Web-Specific Rules

```yaml
level: L1
pattern: "document\\."
scope: "packages/server/src/**"
message: DOM access forbidden in server package.
```
```

### Per-Package Exceptions

Use `exclude` to exempt specific packages:

```yaml
level: L2
pattern: "console\\.log"
scope: "packages/*/src/**"
exclude: "packages/cli/src/**"
message: No console.log except in CLI package.
```

## Validation Modes

### Full Validation

Check all rules against entire codebase:

```
/start:validate constitution
```

### Targeted Validation

Check against specific files or directories:

```
/start:validate constitution src/services/
```

### Implementation Validation

Called automatically during `/start:implement` phases.

## Performance Considerations

1. **Pattern rules are faster** than Check rules (no LLM needed)
2. **Narrow scopes** reduce file scanning
3. **Specific excludes** prevent unnecessary checks
4. **Batch similar rules** under same scope for efficiency

## Troubleshooting

### Rule Not Matching

1. Check regex escaping in YAML (double backslashes)
2. Verify scope matches target files: `ls [scope-pattern]`
3. Check if exclude is too broad
4. For Check rules: verify description is specific enough

### Too Many False Positives

1. Add specific excludes for legitimate uses
2. Narrow scope to problem areas
3. Consider L3 (advisory) instead of L1/L2
4. Refine pattern or check description

### Invalid Regex Error

Common issues:
- Missing escape for special chars: `(`, `)`, `.`, `*`
- Unbalanced groups
- Invalid quantifiers

Test regex at regex101.com before adding to constitution.
