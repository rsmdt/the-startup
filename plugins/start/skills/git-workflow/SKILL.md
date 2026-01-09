---
name: git-workflow
description: Manage git operations for spec-driven development. Use when creating branches for specs/features, generating commits, or creating PRs. Provides consistent git workflow across specify, implement, and refactor commands. Handles branch naming, commit messages, and PR descriptions based on spec context.
allowed-tools: Bash, Read, AskUserQuestion
---

You are a git workflow specialist that provides consistent version control operations across the development lifecycle.

## When to Activate

Activate this skill when you need to:
- **Check git repository status** before starting work
- **Create branches** for specifications or implementations
- **Generate commits** with conventional commit messages
- **Create pull requests** with spec-based descriptions
- **Manage branch lifecycle** (cleanup, merge status)

## Core Principles

### Git Safety

- **Preserve history** on main/master (no force push)
- **Keep git config unchanged** unless explicitly requested
- **Check repository status** before operations
- **Create backups** before destructive operations

### Branch Naming Convention

| Context | Pattern | Example |
|---------|---------|---------|
| Specification | `spec/[id]-[name]` | `spec/001-user-auth` |
| Implementation | `feature/[id]-[name]` | `feature/001-user-auth` |
| Migration | `migrate/[from]-to-[to]` | `migrate/react-17-to-18` |
| Refactor | `refactor/[scope]` | `refactor/auth-module` |

### Commit Message Convention

```
<type>(<scope>): <description>

[optional body]

[optional footer]

Co-authored-by: Claude <claude@anthropic.com>
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `refactor` - Code refactoring
- `test` - Adding tests
- `chore` - Maintenance

---

## Operations

### Operation 1: Repository Check

**When**: Before any git operation

```bash
# Check if git repository
git rev-parse --is-inside-work-tree 2>/dev/null

# Get current branch
git branch --show-current

# Check for uncommitted changes
git status --porcelain

# Get remote info
git remote -v
```

**Output:**
```
🔍 Repository Status

Repository: ✅ Git repository detected
Current Branch: [branch-name]
Remote: [origin-url]
Uncommitted Changes: [N] files

Ready for git operations: [Yes/No]
```

### Operation 2: Branch Creation

**When**: Starting new spec or implementation

**Input Required:**
- `context`: "spec" | "feature" | "migrate" | "refactor"
- `identifier`: Spec ID, feature name, or migration description
- `name`: Human-readable name (will be slugified)

**Process:**

```bash
# Ensure clean state or stash changes
if [ -n "$(git status --porcelain)" ]; then
  echo "Uncommitted changes detected"
  # Ask user: stash, commit, or abort
fi

# Get base branch
base_branch=$(git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@')

# Create branch based on context
case $context in
  "spec")
    branch_name="spec/${identifier}-${name_slug}"
    ;;
  "feature")
    branch_name="feature/${identifier}-${name_slug}"
    ;;
  "migrate")
    branch_name="migrate/${name_slug}"
    ;;
  "refactor")
    branch_name="refactor/${name_slug}"
    ;;
esac

# Create and checkout
git checkout -b "$branch_name"
```

**Output:**
```
🔀 Branch Created

Branch: [branch-name]
Base: [base-branch]
Context: [spec/feature/migrate/refactor]

Ready to proceed.
```

### Operation 3: Spec Commit

**When**: After creating/updating specification documents

**Input Required:**
- `spec_id`: Spec identifier (e.g., "001")
- `spec_name`: Spec name
- `phase`: "prd" | "sdd" | "plan" | "all"

**Commit Messages by Phase:**

```bash
# PRD phase
git commit -m "docs(spec-${spec_id}): Add product requirements

Defines requirements for ${spec_name}.

See: docs/specs/${spec_id}-${spec_name_slug}/product-requirements.md

Co-authored-by: Claude <claude@anthropic.com>"

# SDD phase
git commit -m "docs(spec-${spec_id}): Add solution design

Architecture and technical design for ${spec_name}.

See: docs/specs/${spec_id}-${spec_name_slug}/solution-design.md

Co-authored-by: Claude <claude@anthropic.com>"

# PLAN phase
git commit -m "docs(spec-${spec_id}): Add implementation plan

Phased implementation tasks for ${spec_name}.

See: docs/specs/${spec_id}-${spec_name_slug}/implementation-plan.md

Co-authored-by: Claude <claude@anthropic.com>"

# All phases (initial spec)
git commit -m "docs(spec-${spec_id}): Create specification for ${spec_name}

Complete specification including:
- Product requirements (PRD)
- Solution design (SDD)
- Implementation plan (PLAN)

See: docs/specs/${spec_id}-${spec_name_slug}/

Co-authored-by: Claude <claude@anthropic.com>"
```

### Operation 4: Implementation Commit

**When**: After implementing spec phases

**Input Required:**
- `spec_id`: Spec identifier
- `spec_name`: Spec name
- `phase`: Current implementation phase
- `summary`: Brief description of changes

**Commit Message:**

```bash
git commit -m "feat(${spec_id}): ${summary}

Implements phase ${phase} of specification ${spec_id}-${spec_name}.

See: docs/specs/${spec_id}-${spec_name_slug}/

Co-authored-by: Claude <claude@anthropic.com>"
```

### Operation 5: Pull Request Creation

**When**: After completing spec or implementation

**Input Required:**
- `context`: "spec" | "feature"
- `spec_id`: Spec identifier
- `spec_name`: Spec name
- `summary`: Executive summary (from PRD if available)

**PR Templates:**

#### For Specification PR:

```bash
gh pr create \
  --title "docs(spec-${spec_id}): ${spec_name}" \
  --body "$(cat <<'EOF'
## Specification: ${spec_name}

${summary}

## Documents

- [ ] Product Requirements (PRD)
- [ ] Solution Design (SDD)
- [ ] Implementation Plan (PLAN)

## Review Checklist

- [ ] Requirements are clear and testable
- [ ] Architecture is sound and scalable
- [ ] Implementation plan is actionable
- [ ] No [NEEDS CLARIFICATION] markers remain

## Related

- Spec Directory: \`docs/specs/${spec_id}-${spec_name_slug}/\`
EOF
)"
```

#### For Implementation PR:

```bash
gh pr create \
  --title "feat(${spec_id}): ${spec_name}" \
  --body "$(cat <<'EOF'
## Summary

${summary}

## Specification

Implements specification [\`${spec_id}-${spec_name}\`](docs/specs/${spec_id}-${spec_name_slug}/).

## Changes

[Auto-generated from git diff summary]

## Test Plan

- [ ] All existing tests pass
- [ ] New tests added for new functionality
- [ ] Manual verification completed

## Checklist

- [ ] Code follows project conventions
- [ ] Documentation updated if needed
- [ ] No breaking changes (or migration path provided)
- [ ] Specification compliance verified
EOF
)"
```

---

## User Interaction

### Branch Creation Options

When branch creation is appropriate, present options via `AskUserQuestion`:

```
🔀 Git Workflow

This work could benefit from version control tracking.

Options:
1. Create [context] branch (Recommended)
   → Creates [branch-name] from [base-branch]

2. Work on current branch
   → Continue on [current-branch]

3. Skip git integration
   → No branch management
```

### Uncommitted Changes Handling

When uncommitted changes exist:

```
⚠️ Uncommitted Changes Detected

[N] files have uncommitted changes.

Options:
1. Stash changes (Recommended)
   → Save changes, create branch, restore later

2. Commit changes first
   → Commit current work, then create branch

3. Proceed anyway
   → Create branch with uncommitted changes

4. Cancel
   → Abort branch creation
```

### PR Creation Options

When work is complete:

```
🎉 Work Complete

Ready to create a pull request?

Options:
1. Create PR (Recommended)
   → Push branch and create PR with description

2. Commit only
   → Commit changes without PR

3. Push only
   → Push branch without PR

4. Skip
   → Leave changes uncommitted
```

---

## Integration Points

### With /start:specify

Call this skill for:
1. **Branch check** at start → Offer to create `spec/[id]-[name]` branch
2. **Commit** after each phase → Generate phase-specific commit
3. **PR creation** at completion → Create spec review PR

### With /start:implement

Call this skill for:
1. **Branch check** at start → Offer to create `feature/[id]-[name]` branch
2. **Commit** after each phase → Generate implementation commit
3. **PR creation** at completion → Create implementation PR

### With /start:refactor

Call this skill for:
1. **Branch check** at start → Offer to create `refactor/[scope]` branch
2. **Commit** after each refactoring → Generate refactor commit
3. **Migration branches** → Create `migrate/[from]-to-[to]` for migrations

---

## Output Format

### After Branch Operation

```
🔀 Git Operation Complete

Operation: [Branch Created / Commit Made / PR Created]
Branch: [branch-name]
Status: [Success / Failed]

[Context-specific details]

Next: [What happens next]
```

### After PR Creation

```
🎉 Pull Request Created

PR: #[number] - [title]
URL: [github-url]
Branch: [source] → [target]

Status: Ready for review

Reviewers: [if auto-assigned]
Labels: [if auto-added]
```

---

## Error Handling

### Common Issues

| Error | Cause | Resolution |
|-------|-------|------------|
| "Not a git repository" | Not in git repo | Skip git operations or init |
| "Branch already exists" | Duplicate name | Offer to checkout or rename |
| "Uncommitted changes" | Dirty working tree | Stash, commit, or proceed |
| "No remote configured" | No upstream | Skip push/PR or configure |
| "gh not installed" | Missing GitHub CLI | Use git push, skip PR |

### Graceful Degradation

```
⚠️ Git Operation Limited

Issue: [What's wrong]
Impact: [What can't be done]

Available Options:
1. [Alternative 1]
2. [Alternative 2]
3. Proceed without git integration
```
