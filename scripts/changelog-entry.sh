#!/usr/bin/env bash
# Generate a CHANGELOG.md entry for the current HEAD and print it to stdout.
#
# Usage:
#     scripts/changelog-entry.sh <new-version> [previous-tag]
#
# Reads commits in [previous-tag..HEAD] and pipes them to `claude -p` with
# instructions to synthesize a Keep a Changelog entry — highlighting breaking
# changes, grouping related commits, and skipping trivial noise. The previous
# tag is auto-detected via `git describe` if not supplied.
#
# Output is written to stdout. `scripts/release.sh` captures it and prepends
# to CHANGELOG.md; run this script standalone if you want to preview the entry.

set -euo pipefail

NEW_VERSION="${1:-}"
if [ -z "$NEW_VERSION" ]; then
    echo "Usage: $0 <new-version> [previous-tag]" >&2
    exit 2
fi

PREV_TAG="${2:-$(git describe --tags --abbrev=0 HEAD 2>/dev/null || echo '')}"
if [ -z "$PREV_TAG" ]; then
    echo "No previous tag found. Specify one explicitly as the second argument." >&2
    exit 1
fi

DATE=$(date +%Y-%m-%d)

# Collect commits with bodies so 'BREAKING CHANGE:' footers are captured.
# Separator is unlikely to appear in a real commit.
COMMITS=$(git log --format='%h %s%n%b%n---%%%---' "${PREV_TAG}..HEAD")

if [ -z "$COMMITS" ]; then
    echo "No commits between ${PREV_TAG} and HEAD." >&2
    exit 1
fi

if ! command -v claude >/dev/null 2>&1; then
    echo "claude CLI not found. Install with: npm install -g @anthropic-ai/claude-code" >&2
    exit 1
fi

PROMPT=$(cat <<EOF
You are writing a CHANGELOG.md entry for version ${NEW_VERSION} of
The Agentic Startup, dated ${DATE}.

The Agentic Startup is a Claude Code marketplace distributing plugins
(start, team, consultant) — workflow commands, specialized agents,
skills, and output styles. Users consume it via \`/plugin install\`.

Below is the git commit history since the previous release (${PREV_TAG}).
Commits are separated by '---%---'. Conventional commit prefixes indicate
the type: feat, fix, refactor, chore, docs, test, style, ci, perf.

Write a CHANGELOG entry in Keep a Changelog format. Follow these rules:

1. Header: '## [${NEW_VERSION}] - ${DATE}'
2. Group into these sections (omit empty ones): Added, Changed, Fixed,
   Removed, Deprecated, Security.
3. **Summarize, don't enumerate.** Merge related commits into single
   coherent bullets. Target 2-6 bullets per section. The git log already
   has the per-commit detail — this is the human-readable summary.
4. **Highlight breaking changes** with a '**BREAKING**:' prefix on the
   bullet. Detect them from (a) '!' after the commit type, (b)
   'BREAKING CHANGE:' footers in the body, (c) contextual cues such as
   renamed/removed commands, skills, agents, or changed plugin manifests
   that affect installed users.
5. **Skip trivial commits**: style, chore (except version bumps — still
   skip those), test, ci, and routine docs changes. Include them only
   if they represent user-visible changes (e.g. a README rewrite
   documenting a new workflow, a new agent or skill users should know
   about, or a new output style).
6. **Write for users, not contributors.** Describe what changed from
   the user's perspective ('Added /start:analyze skill for mechanism-level
   codebase analysis'), not how the code was restructured ('Created new
   file plugins/start/skills/analyze/SKILL.md'). Prefer prose that
   explains the outcome over technical implementation detail.
7. Output the markdown entry only. No preamble, no code fences, no
   explanation before or after. Start with the '## [${NEW_VERSION}]'
   header and end with the last bullet of the last section.

Commit history:

${COMMITS}
EOF
)

claude -p "$PROMPT"
