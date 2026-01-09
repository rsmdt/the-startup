---
name: constitution-validation
description: |
  Create and validate project constitutions through discovery-based rule generation.
  Use when creating governance rules, validating code against constitutional rules,
  or checking constitution compliance during implementation and review.
allowed-tools: Task, Read, Write, Edit, Grep, Glob, Bash
---

# Constitution Validation Skill

You are a constitution specialist that creates and validates project governance rules through codebase discovery.

## When to Activate

Activate this skill when you need to:
- **Create a new constitution** by discovering project patterns
- **Validate existing code** against constitution rules
- **Update constitution** with new rules or categories
- **Check constitution compliance** during implementation or review

**IMPORTANT**: Explore the actual codebase to discover patterns. Base all rules on observed frameworks and technologies. Use `[NEEDS DISCOVERY]` markers to guide exploration.

## Core Philosophy

### Discovery-Based Rules

**Generate rules dynamically from codebase exploration.** Process:

1. **Explore First**: Use Glob, Grep, Read to understand the project
2. **Discover Patterns**: What frameworks? What conventions? What architecture?
3. **Generate Rules**: Based on what you actually found
4. **Validate with User**: Present discovered patterns before finalizing

### Level System (L1/L2/L3)

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

## Template

The constitution template is at [template.md](template.md). Use this structure exactly.

**To create a constitution:**
1. Read the template: `plugins/start/skills/constitution-validation/template.md`
2. Explore codebase to resolve all `[NEEDS DISCOVERY]` markers
3. Generate rules based on actual patterns found
4. Write to project root: `CONSTITUTION.md`

## Cycle Pattern

For each category requiring rules, follow this iterative process:

### 1. Discovery Phase

- **Explore the codebase** to understand actual patterns
- **Launch parallel agents** to investigate:
  - Security patterns (auth, secrets, validation)
  - Architecture patterns (layers, boundaries, dependencies)
  - Code quality conventions (naming, formatting, structure)
  - Testing setup (frameworks, coverage, patterns)
  - Framework-specific considerations

### 2. Documentation Phase

- **Update the constitution** with discovered rules
- **Replace `[NEEDS DISCOVERY]` markers** with actual rules
- Focus only on current category being processed
- Generate rules that are specific to this project

### 3. Review Phase

- **Present discovered patterns** to user
- Show proposed rules with rationale
- Highlight rules needing user confirmation
- **Wait for user confirmation** before next cycle

**Ask yourself each cycle:**
1. Have I explored the actual codebase (not assumed)?
2. Have I discovered real patterns (not guessed)?
3. Have I generated project-specific rules?
4. Have I presented findings to the user?
5. Have I received user confirmation?

## Rule Generation Guidelines

When generating rules from discovered patterns:

### L1 Rules (Blocking + Autofix)

Generate for patterns that are:
- Security critical (secrets, injection, auth)
- Clearly fixable with deterministic changes
- Objectively wrong (not style preference)

**Examples:**
- Hardcoded secrets → Replace with env var reference
- `eval()` usage → Remove and use safer alternative
- Barrel exports → Convert to direct imports

### L2 Rules (Blocking, No Autofix)

Generate for patterns that are:
- Architecturally important
- Require human judgment to fix
- May have valid exceptions

**Examples:**
- Database calls outside repository layer
- Cross-package imports via relative paths
- Missing error handling

### L3 Rules (Advisory)

Generate for patterns that are:
- Style preferences
- Best practices that vary by context
- Suggestions, not requirements

**Examples:**
- Function length recommendations
- Test file presence
- Documentation coverage

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

## Validation Mode

When validating (not creating), skip discovery and:

1. **Parse existing constitution** rules
2. **Apply scopes** to find matching files
3. **Execute checks** (Pattern or Check rules)
4. **Generate compliance report**

### Rule Parsing

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

### Validation Execution

For each parsed rule:

1. **Glob files matching scope** (excluding patterns in `exclude`)
2. **For Pattern rules**: Execute regex match against file contents
3. **For Check rules**: Use LLM to interpret semantic check
4. **Collect violations** with file path, line number, code snippet
5. **Categorize by level** for reporting

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

## Graceful Degradation

| Scenario | Behavior |
|----------|----------|
| No CONSTITUTION.md | Report "No constitution found. Skipping constitution checks." |
| Invalid rule format | Skip rule, warn user, continue with other rules |
| Invalid regex pattern | Report as config error, skip rule |
| Scope matches no files | Report as info, not a failure |
| File read error | Skip file, warn, continue |

## Integration Points

This skill is called by:
- `/start:constitution` - For creation and updates
- `/start:validate` (Mode E) - For constitution validation
- `/start:implement` - For active enforcement during implementation
- `/start:review` - For code review compliance checks
- `/start:specify` (SDD phase) - For architecture alignment

## Validation Checklist

Before completing constitution creation:

- [ ] All `[NEEDS DISCOVERY]` markers resolved
- [ ] Every rule has valid level (L1/L2/L3)
- [ ] Every rule has either `pattern` or `check`
- [ ] Every rule has `scope` and `message`
- [ ] Rules are specific to this project (not generic)
- [ ] User has confirmed proposed rules

## Output Format

After constitution work, report:

```
📜 Constitution Status: [Created / Updated / Validated]

Discovery Findings:
- Project Type: [discovered type]
- Frameworks: [discovered frameworks]
- Key Patterns: [patterns found]

Categories:
- Security: [N] rules
- Architecture: [N] rules
- Code Quality: [N] rules
- Testing: [N] rules
- [Project-Specific]: [N] rules

User Confirmations:
- [Rule 1]: ✅ Confirmed
- [Rule 2]: ⏳ Pending

Next Steps:
- [What needs to happen next]
```
