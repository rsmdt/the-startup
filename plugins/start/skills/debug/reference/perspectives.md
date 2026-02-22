# Investigation Perspectives

Perspective definitions for parallel hypothesis testing during debugging.

---

## Perspectives

| Perspective | Intent | What to Investigate |
|-------------|--------|---------------------|
| 🔴 **Error Trace** | Follow the error path | Stack traces, error messages, exception handling, error propagation |
| 🔀 **Code Path** | Trace execution flow | Conditional branches, data transformations, control flow, early returns |
| 🔗 **Dependencies** | Check external factors | External services, database queries, API calls, network issues |
| 📊 **State** | Inspect runtime values | Variable values, object states, race conditions, timing issues |
| 🌍 **Environment** | Compare contexts | Configuration, versions, deployment differences, env variables |
| 🕐 **Recent Changes** | Identify regression source | Recent commits, git blame at failure site, dependency updates, recently modified config |

## Bug Type Investigation Patterns

| Bug Type | What to Check | How to Report |
|----------|---------------|---------------|
| Logic errors | Data flow, boundary conditions | "The condition on line X doesn't handle case Y" |
| Integration | API contracts, versions | "The API expects X but we're sending Y" |
| Timing/async | Race conditions, await handling | "There's a race between A and B" |
| Intermittent | Variable conditions, state | "This fails when [condition] because [reason]" |

## Perspective Selection

Not all perspectives are needed for every bug. Select based on hypotheses from Phase 1:

- Error message present → 🔴 Error Trace
- Execution produces wrong result → 🔀 Code Path
- External service involved → 🔗 Dependencies
- Intermittent or timing-related → 📊 State
- Works locally but not in CI/prod → 🌍 Environment
- Regression or "worked before" → 🕐 Recent Changes
