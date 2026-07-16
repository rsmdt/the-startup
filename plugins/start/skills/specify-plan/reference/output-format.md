# Output Format

How to report plan status to the user.

## After building a phase (step 3)

Present the phase with all four facets correlated:

```
## Phase {{N}}: {{title}}   (depends on: {{phases or none}})
{{one-line outcome}}

Features (ACs):     {{AC-N.1, AC-N.2 ...}}
Behavior (EARS):    {{count}} statements
Quality Gates:      {{count}}
Validation:         {{p0}} P0, {{p1}} P1, {{p2}} P2 scenarios

Coverage map:
| Requirement AC | Behavior (EARS) | Proven by scenario |
|----------------|-----------------|--------------------|
| AC-N.1         | When X, shall Y | phase-N/{{name}}   |

Gaps: {{ACs/behaviors without a scenario, or "none"}}
```

## Status report (step 6)

```
## Plan Status — {{FEATURE_TITLE}}

Target State:  Complete | Incomplete
Phases:        {{n}} ({{vertical slices, dependency-ordered}})
Behaviors:     {{n}} EARS statements
Validation:    {{n}} scenarios ({{p0}} P0, {{p1}} P1, {{p2}} P2)
Threshold:     {{0.90}}
AC coverage:   {{x}}/{{y}} acceptance criteria mapped to a phase ({{pct}}%)

Validation: {{pass | list of gaps}}
```

Keep it scannable. Surface gaps prominently — an unmapped acceptance criterion or a phase with no P0 scenario is the most important thing for the user to see.
