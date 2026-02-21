# Output Format Reference

Task result templates, phase summary, and completion summary for the implement skill.

---

## Task Result Templates

### Success

```
✅ Task [N]: [Name]

Files: src/services/auth.ts, src/routes/auth.ts
Summary: Implemented JWT authentication with bcrypt password hashing
Tests: 5 passing
```

### Blocked

```
⚠️ Task [N]: [Name]

Status: Blocked
Reason: Missing User model - need src/models/User.ts
Options: [present via AskUserQuestion]
```

---

## Phase Summary

Present at each phase checkpoint:

```
Phase [N] Complete

Tasks: [X/Y] completed
Files Changed: [list of paths]
Tests: [All passing / X failing]
Blockers: [none / list]
Drift: [none detected / issues found]
```

---

## Completion Summary

```
✅ Implementation Complete

Spec: [NNN]-[name]
Phases Completed: [N/N]
Tasks Executed: [X] total
Tests: [All passing / X failing]
Mode: [Standard / Team]

Files Changed: [N] files (+[additions] -[deletions])
```
