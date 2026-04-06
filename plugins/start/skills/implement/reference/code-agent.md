# Code Agent Prompt Template

The orchestrator constructs this prompt for the Agent tool's `prompt` parameter when spawning a code agent. Variable placeholders are wrapped in `{braces}`.

---

## Prompt Template

```
You are implementing the following specification.

{unit_spec_content}

Codebase context: read AGENTS.md for project orientation.

DO NOT read or access files in scenarios/ directories.
DO NOT access any files under .start/specs/*/scenarios/.

{retry_block}

When complete:
1. All existing tests must pass
2. New tests must cover all requirements in the spec
3. Run the full test suite and report results
```

## Retry Block

Only include when `unit.iteration > 0`. Omit entirely on the first iteration.

```
Previous evaluation found these issues (iteration {iteration} of max {max_iterations}):
- "{failure_summary_1}"
- "{failure_summary_2}"

Investigate the codebase to find and fix the root causes.
```

## Variable Reference

| Variable | Source | Description |
|----------|--------|-------------|
| `{unit_spec_content}` | Read from `{specDirectory}/units/{unit.id}.md` | Full text of the unit spec file, including YAML frontmatter |
| `{iteration}` | `unit.iteration` | Current retry count (1-indexed when displayed) |
| `{max_iterations}` | `manifest.maxIterations` | Maximum retries from manifest frontmatter |
| `{failure_summary_N}` | `unit.failureSummaries[N]` | One-line extracted from evaluation report |

## Information Barrier

**Included (code agent sees):**
- Full unit spec content (goal, requirements, constraints)
- Instruction to read AGENTS.md for project conventions
- One-line failure summaries on retry (observable symptoms only)

**Excluded (code agent never sees):**
- Scenario text from `scenarios/{id}/*.md`
- Evaluation agent output or satisfaction reports
- Other units' specs or results

**Behavioral reinforcement:**
- Explicit "DO NOT read scenarios/" instruction
- The agent starts with a fresh context — no inherited conversation history
