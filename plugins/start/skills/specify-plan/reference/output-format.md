# Output Format

How to report plan status to the user.

## After scenario generation (step 3)

Present completion criteria grouped by capability:

```
## Completion Criteria — {{FEATURE_TITLE}}

### {{Capability A}}
- [P0] {{scenario title}} — {{one-line observable behavior}}
- [P1] {{scenario title}} — {{one-line observable behavior}}

Coverage map:
| Requirement | Proven by |
|-------------|-----------|
| {{req}}     | {{capability/scenario}} |

Gaps: {{requirements with no completion criterion, or "none"}}
```

## Status report (step 6)

```
## Plan Status — {{FEATURE_TITLE}}

Target State:        Complete | Incomplete
Quality Gates:       {{n}} defined
Completion Criteria: {{n}} scenarios across {{m}} capabilities ({{p0}} P0, {{p1}} P1, {{p2}} P2)
Threshold:           {{0.90}}
Requirement coverage: {{x}}/{{y}} requirements ({{pct}}%)

Validation: {{pass | list of gaps}}
```

Keep it scannable. Surface gaps prominently — an uncovered requirement is the most important thing for the user to see.
