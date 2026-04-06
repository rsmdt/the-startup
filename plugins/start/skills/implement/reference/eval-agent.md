# Evaluation Agent Prompt Template

The orchestrator constructs this prompt for the Agent tool's `prompt` parameter when spawning an evaluation agent. Variable placeholders are wrapped in `{braces}`.

---

## Prompt Template

```
You are evaluating whether an implementation satisfies the following scenarios.

{all_scenario_content}

The service is running at localhost:{port}.
Test each scenario through the external interface only.

## Evaluation Method (use first available, in priority order)

1. **E2E test automation**: Write and run automated end-to-end tests based on
   the described scenarios. Use the project's existing test framework.
2. **Browser automation**: Use agent-browser or Playwright to perform a manual
   walkthrough of each scenario through the UI or API.
3. **CLI fallback**: Use Bash with curl/httpie to exercise API endpoints directly.

## Restrictions

DO NOT read source code files, unit spec files, or implementation details.
DO NOT access any files under plugins/, src/, or .start/specs/*/units/.

## Reporting

For each scenario, run 3 times. Report:

Satisfaction: {passed}/{total} scenarios ({percentage}%)
Threshold: {threshold}%

Passed:
- {scenario name}: {pass_count}/3

Failed:
- {scenario name}: {one-line description of failure} ({fail_count}/3 failures)
```

## Variable Reference

| Variable | Source | Description |
|----------|--------|-------------|
| `{all_scenario_content}` | Read from `{specDirectory}/scenarios/{unit.id}/*.md` | Concatenated full text of all scenario files for this unit, each separated by a horizontal rule |
| `{port}` | `state.servicePort` | Discovered from AGENTS.md or package.json |
| `{threshold}` | `manifest.threshold * 100` | Displayed as percentage (e.g., 90%) |

## Information Barrier

**Included (evaluation agent sees):**
- Full scenario text (all scenarios for this unit)
- Service URL (localhost with port)
- Evaluation method priority
- Reporting format instructions

**Excluded (evaluation agent never sees):**
- Unit spec content from `units/{id}.md`
- AGENTS.md or codebase orientation
- Code agent output or implementation details
- Other units' scenarios or results

**Behavioral reinforcement:**
- Explicit "DO NOT read source code" instruction
- Explicit "DO NOT access files under plugins/, src/, or .start/specs/*/units/" instruction
- The agent starts with a fresh context — no inherited conversation history
