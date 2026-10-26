# Output Format

How to report plan status to the user.

## Phase approval (step 3)

Present each of the phase's slices with its facets correlated:

```
## Slice: {{title}}   (depends on: {{slice titles or none}})
{{one-line outcome}}

Features (ACs):     {{AC-N.1, AC-N.2 ...}}
Behavior (EARS):    {{count}} statements
Quality Gates:      {{count}}
Acceptance:         {{p0}} P0, {{p1}} P1, {{p2}} P2 scenario tasks

Coverage map:
| Requirement AC | Behavior (EARS) | Proven by scenario |
|----------------|-----------------|--------------------|
| AC-N.1         | When X, shall Y | {{scenario name}} |

Gaps: {{ACs/behaviors without a scenario, or "none"}}
```

## Status report (step 5)

```
## Plan Status — {{FEATURE_TITLE}}

Target State:  Complete | Incomplete
Phases:        {{n}}
Slices:        {{n}} ({{vertical, dependency-ordered}})
Behaviors:     {{n}} EARS statements
Acceptance:    {{n}} scenario tasks ({{p0}} P0, {{p1}} P1, {{p2}} P2)
Threshold:     {{0.90}}   Max iterations/slice: {{3}}
AC coverage:   {{x}}/{{y}} acceptance criteria mapped to a slice ({{pct}}%)

Validation: {{pass | list of gaps}}
```

Show the task tree too (`task.js list --spec {{NNN}} --depth 1 --format tree`). Keep it scannable — surface gaps
prominently: an unmapped acceptance criterion or a slice with no P0 scenario is the most important
thing for the user to see.
