# Constitution Validation Reference

Techniques for validating code against project governance rules.

## Level System (L1/L2/L3)

| Level | Name | Blocking | Autofix | Use Case |
|-------|------|----------|---------|----------|
| **L1** | Must | ✅ Yes | ✅ AI auto-corrects | Critical rules - security, correctness, architecture |
| **L2** | Should | ✅ Yes | ❌ No (needs human judgment) | Important rules requiring manual attention |
| **L3** | May | ❌ No | ❌ No | Advisory/optional - style preferences, suggestions |

**Level Behavior:**

| Level | Validation | Implementation | AI Behavior |
|-------|------------|----------------|-------------|
| `L1` | Fails check, blocks | Blocks phase completion | **Automatically fixes** before proceeding |
| `L2` | Fails check, blocks | Blocks phase completion | Reports violation, **requires human action** |
| `L3` | Reports only | Does not block | Optional improvement, can be ignored |

---

## Rule Schema

Each rule in the constitution uses this YAML structure:

```yaml
level: L1 | L2 | L3
pattern: "regex pattern"    # OR
check: "semantic description for LLM interpretation"
scope: "glob pattern for files to check"
exclude: "glob patterns to skip (comma-separated)"
message: "Human-readable violation message"
```

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `level` | Required | `L1` \| `L2` \| `L3` | Determines blocking and autofix behavior |
| `pattern` | One of | Regex | Pattern to match violations in source code |
| `check` | One of | String | Semantic description for LLM interpretation |
| `scope` | Required | Glob | File patterns to check (supports `**`) |
| `exclude` | Optional | Glob | File patterns to skip (comma-separated) |
| `message` | Required | String | Human-readable violation message |

---

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

### Check Rules

Check rules use semantic descriptions that the LLM interprets. These are flexible but non-deterministic.

**When to use:**
- Architectural patterns that require understanding context
- Rules that span multiple lines or files
- Semantic concepts (like "database calls only in repositories")

**Example - Repository Pattern:**

```yaml
level: L1
check: Database queries (Prisma, TypeORM, Knex, raw SQL) only in files matching *Repository.ts or *Repository.js
scope: "src/**/*.{ts,js}"
exclude: "**/repositories/**"
message: Direct database call outside repository layer.
```

---

## Validation Execution

For each parsed rule:

1. **Glob files matching scope** (excluding patterns in `exclude`)
2. **For Pattern rules**: Execute regex match against file contents
3. **For Check rules**: Use LLM to interpret semantic check
4. **Collect violations** with file path, line number, code snippet
5. **Categorize by level** for reporting

---

## Rule Parsing

```pseudocode
FUNCTION: parse_constitution(markdown_content)
  rules = []
  current_category = null

  FOR EACH section in markdown:
    IF section.header.level == 2:
      current_category = section.header.text  # e.g., "Code Quality", "Security"
    ELSE IF section.header.level == 3:
      yaml_block = extract_yaml_code_block(section.content)
      IF yaml_block:
        rule = {
          id: generate_rule_id(current_category, index),  # e.g., "SEC-001"
          name: section.header.text,                       # e.g., "No Hardcoded Secrets"
          category: current_category,
          level: yaml_block.level,
          pattern: yaml_block.pattern,
          check: yaml_block.check,
          scope: yaml_block.scope,
          exclude: yaml_block.exclude,
          message: yaml_block.message,
        }
        IF rule.pattern OR rule.check:
          # Derive behavior from level
          rule.blocking = (rule.level == "L1" OR rule.level == "L2")
          rule.autofix = (rule.level == "L1")
          rules.append(rule)
  RETURN rules
```

---

## Category ID Prefixes

| Category | Prefix | Example |
|----------|--------|---------|
| Security | SEC | SEC-001 |
| Architecture | ARCH | ARCH-001 |
| Code Quality | QUAL | QUAL-001 |
| Testing | TEST | TEST-001 |
| Custom | CUST | CUST-001 |
| [Custom Name] | First 4 letters uppercase | PERF-001 |

---

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
```

### Testing

```yaml
# No .only in Committed Tests
level: L1
pattern: "\\.(only|skip)\\s*\\("
scope: "**/*.test.*, **/*.spec.*"
message: Remove .only/.skip before committing tests.
```

---

## Compliance Report Format

```markdown
## Constitution Compliance Report

**Constitution:** CONSTITUTION.md
**Target:** [spec-id or file path or "entire codebase"]
**Checked:** [ISO timestamp]

### Summary

- ✅ Passed: [N] rules
- ⚠️ L3 Advisories: [N] rules
- ❌ L2 Blocking: [N] rules
- 🛑 L1 Critical: [N] rules

### Critical Violations (L1 - Autofix Required)

#### 🛑 SEC-001: No Hardcoded Secrets
- **Location:** `src/services/PaymentService.ts:42`
- **Finding:** Hardcoded secret detected. Use environment variables.
- **Code:** `const API_KEY = 'sk_live_xxx...'`
- **Autofix:** Replace with `process.env.PAYMENT_API_KEY`

### Blocking Violations (L2 - Human Action Required)

#### ❌ ARCH-001: Repository Pattern
- **Location:** `src/services/UserService.ts:18`
- **Finding:** Direct database call outside repository.
- **Code:** `await prisma.user.findMany(...)`
- **Action Required:** Extract to UserRepository

### Advisories (L3 - Optional)

#### ⚠️ QUAL-001: Function Length
- **Location:** `src/utils/helpers.ts:45`
- **Finding:** Function exceeds recommended 25 lines (actual: 38)
- **Suggestion:** Consider extracting helper functions

### Recommendations

1. [Prioritized action item based on violations]
2. [Next action item]
```

---

## Graceful Degradation

| Scenario | Behavior |
|----------|----------|
| No CONSTITUTION.md | Report "No constitution found. Skipping constitution checks." |
| Invalid rule format | Skip rule, warn user, continue with other rules |
| Invalid regex pattern | Report as config error, skip rule |
| Scope matches no files | Report as info, not a failure |
| File read error | Skip file, warn, continue |

---

## Scope Patterns

### Common Patterns

| Pattern | Matches |
|---------|---------|
| `**/*.ts` | All TypeScript files |
| `src/**/*.ts` | TypeScript files in src/ |
| `**/*.{ts,js}` | TypeScript and JavaScript |
| `packages/*/src/**` | All packages' src folders |
| `apps/web/**` | Only the web app |

### Monorepo Scoping

```yaml
# Web package only
scope: "packages/web/src/**/*.{ts,tsx}"

# All packages
scope: "packages/*/src/**/*.ts"

# Shared libraries
scope: "libs/*/src/**/*.ts"
```

---

## Performance Considerations

1. **Pattern rules are faster** than Check rules (no LLM needed)
2. **Narrow scopes** reduce file scanning
3. **Specific excludes** prevent unnecessary checks
4. **Batch similar rules** under same scope for efficiency
