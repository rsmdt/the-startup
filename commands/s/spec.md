---
description: "Create numbered spec directories with TOML format"
argument-hint: "feature-name [--add <template>] [--read]"
---

# Spec Generation

!bash ${CLAUDE_PLUGIN_ROOT}/scripts/spec.sh $ARGUMENTS

---

This command generates specification directories with auto-incrementing IDs.

**Usage:**

Create new spec:
```
/s:spec user-authentication
```

Create spec with template:
```
/s:spec user-authentication --add product-requirements
```

Read existing spec:
```
/s:spec 001 --read
```

**Available Templates:**
- `product-requirements` - Product Requirements Document
- `solution-design` - Solution Design Document
- `implementation-plan` - Implementation Plan
- `definition-of-ready` - Quality gate template
- `definition-of-done` - Quality gate template
- `task-definition-of-done` - Task-level quality gate

**Output:**
Creates `docs/specs/NNN-feature-name/` with auto-incremented spec number.
