#!/usr/bin/env bash
#
# eval-skill.sh — thin CLI entrypoint for running a skill's eval set.
#
# This does NOT reimplement an eval harness. It standardizes the invocation of
# the agent-driven skill-creator benchmark workflow (which reads evals/evals.json,
# runs with_skill executors — plus without_skill ones with --baseline — grades
# against the assertions, and aggregates a benchmark). See
# docs/research/running-skill-evals.md.
#
# Skills that launch a dynamic workflow are excluded: skill-creator's executors are
# subagents without the Workflow tool, so those evals cannot pass here.
#
# Interactivity: skills under test call AskUserQuestion / wait for approval. A
# headless run has nobody to answer, so we bake an "autonomous defaults" preamble
# into every executor run — assume approval, pick recommended options, never pause.
#
# Usage:
#   scripts/eval-skill.sh <skill>            # e.g. specify-plan, implement
#   scripts/eval-skill.sh <skill> --runs 3   # variance: N runs per configuration
#   scripts/eval-skill.sh <skill> --baseline # also run without_skill and report the delta
#   scripts/eval-skill.sh --list             # list skills that have an eval set
#
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SKILLS_DIR="$REPO_ROOT/plugins/start/skills"

WORKFLOW_SKILLS="analyze implement review reflect"

AUTONOMOUS_PREAMBLE="Run to completion without pausing for input: treat every approval \
gate as approved and take the recommended option wherever you would ask a question. \
In your final response, list each question or choice you would have put to the user \
and the option you took."

is_workflow_skill() {
  case " $WORKFLOW_SKILLS " in *" $1 "*) return 0 ;; esac
  return 1
}

list_skills() {
  echo "Skills with an eval set:"
  for d in "$SKILLS_DIR"/*/; do
    if [ -f "${d}evals/evals.json" ]; then
      name="$(basename "$d")"
      n="$(grep -c '"id"' "${d}evals/evals.json" 2>/dev/null || echo '?')"
      if is_workflow_skill "$name"; then
        printf "  %-24s %s evals  (excluded: needs the Workflow tool)\n" "$name" "$n"
      else
        printf "  %-24s %s evals\n" "$name" "$n"
      fi
    fi
  done
}

if [ "${1:-}" = "--list" ] || [ "${1:-}" = "-l" ]; then
  list_skills
  exit 0
fi

SKILL="${1:-}"
if [ -z "$SKILL" ]; then
  echo "usage: scripts/eval-skill.sh <skill> [--runs N] [--model NAME] [--baseline]" >&2
  echo "       scripts/eval-skill.sh --list" >&2
  exit 2
fi
shift || true

RUNS=1
MODEL=""
BASELINE=0
while [ $# -gt 0 ]; do
  case "$1" in
    --runs)     RUNS="$2"; shift 2 ;;
    --model)    MODEL="$2"; shift 2 ;;
    --baseline) BASELINE=1; shift ;;
    *) echo "unknown flag: $1" >&2; exit 2 ;;
  esac
done

SKILL_DIR="$SKILLS_DIR/$SKILL"
EVALS="$SKILL_DIR/evals/evals.json"
if [ ! -f "$EVALS" ]; then
  echo "No eval set found at: $EVALS" >&2
  echo "Run 'scripts/eval-skill.sh --list' to see skills that have evals." >&2
  exit 1
fi

if is_workflow_skill "$SKILL"; then
  echo "Skipping '$SKILL': it launches a dynamic workflow, and skill-creator's executors have no Workflow tool." >&2
  echo "Evaluate it in a top-level session instead (one 'claude -p' per eval case)." >&2
  exit 3
fi

if ! command -v claude >/dev/null 2>&1; then
  echo "The 'claude' CLI is required (skill-creator is agent-driven)." >&2
  exit 1
fi

if [ "$BASELINE" = 1 ]; then
  CONFIGS="with_skill and without_skill"
  REPORT="the per-eval pass rates and the with_skill vs without_skill delta"
else
  CONFIGS="with_skill only — do not run a without_skill baseline"
  REPORT="the per-eval pass rates"
fi

PROMPT="Use the skill-creator skill to run the eval set for the skill at \
plugins/start/skills/$SKILL (its evals are in $EVALS). Run $RUNS run(s) per configuration \
($CONFIGS) and write results to the standard ${SKILL}-workspace/iteration-N/ layout. \
Where an eval lists \`files\` (paths under evals/files/<case>/), copy each into that run's working \
directory at its path below evals/files/<case>/ before the run starts. \
Pass --static to the eval viewer to write a standalone HTML report.

This session is headless and ends the moment you end your turn:
- Run every executor and grader in the foreground and wait for all of them to finish \
before you end your turn. Never end your turn while any run is pending.
- Never ask the user anything. When a decision is needed (retry, skip, re-run), take the \
option that completes the benchmark and note it in the final report.

For every executor run, prepend this directive to the eval prompt:

\"$AUTONOMOUS_PREAMBLE\"

When done, report $REPORT."

echo "▶ Evaluating skill '$SKILL' ($RUNS run(s)/config$([ "$BASELINE" = 1 ] && echo ', with baseline')) via skill-creator…"
echo "  evals: $EVALS"
[ -n "$MODEL" ] && echo "  model: $MODEL"
echo

cd "$REPO_ROOT"
if [ -n "$MODEL" ]; then
  exec claude -p "$PROMPT" --model "$MODEL"
else
  exec claude -p "$PROMPT"
fi
