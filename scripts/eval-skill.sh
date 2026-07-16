#!/usr/bin/env bash
#
# eval-skill.sh — thin CLI entrypoint for running a skill's eval set.
#
# This does NOT reimplement an eval harness. It standardizes the invocation of
# the agent-driven skill-creator benchmark workflow (which reads evals/evals.json,
# runs with_skill / without_skill executors, grades against the assertions, and
# aggregates a benchmark). See docs/research/running-skill-evals.md.
#
# Interactivity: skills under test call AskUserQuestion / wait for approval. A
# headless run has nobody to answer, so we bake an "autonomous defaults" preamble
# into every executor run — assume approval, pick recommended options, never pause.
#
# Usage:
#   scripts/eval-skill.sh <skill>            # e.g. specify-plan, implement
#   scripts/eval-skill.sh <skill> --runs 3   # variance: N runs per configuration
#   scripts/eval-skill.sh --list             # list skills that have an eval set
#
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SKILLS_DIR="$REPO_ROOT/plugins/start/skills"

AUTONOMOUS_PREAMBLE="Run autonomously and to completion: assume the user approves at \
every gate, choose the recommended option for any AskUserQuestion, and never pause \
to ask the user — this is a non-interactive benchmark run. Record in the trace which \
options you would have presented and which you chose."

list_skills() {
  echo "Skills with an eval set:"
  for d in "$SKILLS_DIR"/*/; do
    if [ -f "${d}evals/evals.json" ]; then
      name="$(basename "$d")"
      n="$(grep -c '"id"' "${d}evals/evals.json" 2>/dev/null || echo '?')"
      printf "  %-24s %s evals\n" "$name" "$n"
    fi
  done
}

if [ "${1:-}" = "--list" ] || [ "${1:-}" = "-l" ]; then
  list_skills
  exit 0
fi

SKILL="${1:-}"
if [ -z "$SKILL" ]; then
  echo "usage: scripts/eval-skill.sh <skill> [--runs N] [--model NAME]" >&2
  echo "       scripts/eval-skill.sh --list" >&2
  exit 2
fi
shift || true

RUNS=1
MODEL=""
while [ $# -gt 0 ]; do
  case "$1" in
    --runs)  RUNS="$2"; shift 2 ;;
    --model) MODEL="$2"; shift 2 ;;
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

if ! command -v claude >/dev/null 2>&1; then
  echo "The 'claude' CLI is required (skill-creator is agent-driven)." >&2
  exit 1
fi

PROMPT="Use the skill-creator skill to run the eval set and benchmark for the skill at \
plugins/start/skills/$SKILL (its evals are in $EVALS). Run $RUNS run(s) per configuration \
(with_skill and without_skill) and write results to the standard \
${SKILL}-workspace/iteration-N/ layout. This is a headless environment: pass --static to \
the eval viewer to write a standalone HTML report instead of opening a browser.

For EVERY executor run (both with_skill and without_skill), prepend this directive to the \
eval prompt so the run does not hang waiting for input:

\"$AUTONOMOUS_PREAMBLE\"

When done, report the per-eval pass rates and the with_skill vs without_skill delta."

echo "▶ Benchmarking skill '$SKILL' ($RUNS run(s)/config) via skill-creator…"
echo "  evals: $EVALS"
[ -n "$MODEL" ] && echo "  model: $MODEL"
echo

cd "$REPO_ROOT"
if [ -n "$MODEL" ]; then
  exec claude -p "$PROMPT" --model "$MODEL"
else
  exec claude -p "$PROMPT"
fi
