#!/usr/bin/env bash
# Release The Agentic Startup end-to-end.
#
# Usage:
#     scripts/release.sh <new-version> [--yes] [--no-push]
#
# Example:
#     scripts/release.sh 3.7.0
#
# Performs, in order:
#   1. Preflight: clean worktree, on main, claude+jq present, tag not taken.
#   2. Bumps version in:
#        - .claude-plugin/marketplace.json (.metadata.version)
#        - plugins/*/.claude-plugin/plugin.json (.version)
#   3. Generates a CHANGELOG entry via scripts/changelog-entry.sh and
#      prepends it to CHANGELOG.md.
#   4. Shows the diff and asks for confirmation (skip with --yes).
#   5. Commits, tags vX.Y.Z, and pushes main + tag (skip push with --no-push).
#
# The tag push triggers .github/workflows/release.yml, which packages
# archives and creates the GitHub release.

set -euo pipefail

# ---------- args ----------
NEW_VERSION=""
AUTO_YES=0
PUSH=1
for arg in "$@"; do
    case "$arg" in
        --yes|-y) AUTO_YES=1 ;;
        --no-push) PUSH=0 ;;
        -*) echo "Unknown flag: $arg" >&2; exit 2 ;;
        *)
            if [ -z "$NEW_VERSION" ]; then
                NEW_VERSION="$arg"
            else
                echo "Unexpected extra argument: $arg" >&2
                exit 2
            fi
            ;;
    esac
done

if [ -z "$NEW_VERSION" ]; then
    echo "Usage: $0 <new-version> [--yes] [--no-push]" >&2
    exit 2
fi

if ! [[ "$NEW_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Version must be semver (MAJOR.MINOR.PATCH), got: $NEW_VERSION" >&2
    exit 2
fi

TAG="v${NEW_VERSION}"

# ---------- paths ----------
REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

MARKETPLACE="$REPO_ROOT/.claude-plugin/marketplace.json"
CHANGELOG="$REPO_ROOT/CHANGELOG.md"
SCRIPT_DIR="$REPO_ROOT/scripts"

# ---------- preflight ----------
echo "→ Preflight checks"

for cmd in jq claude git; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        echo "  ✗ required command not found: $cmd" >&2
        exit 1
    fi
done

if ! git diff --quiet || ! git diff --cached --quiet; then
    echo "  ✗ worktree is dirty. Commit or stash changes first." >&2
    git status --short >&2
    exit 1
fi

CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "  ✗ must release from main (currently on '$CURRENT_BRANCH')." >&2
    exit 1
fi

if git rev-parse -q --verify "refs/tags/${TAG}" >/dev/null; then
    echo "  ✗ tag ${TAG} already exists." >&2
    exit 1
fi

if [ ! -f "$MARKETPLACE" ]; then
    echo "  ✗ marketplace manifest missing: $MARKETPLACE" >&2
    exit 1
fi

if [ ! -f "$CHANGELOG" ]; then
    echo "  ✗ CHANGELOG.md missing at repo root." >&2
    exit 1
fi

PREV_TAG="$(git describe --tags --abbrev=0 HEAD 2>/dev/null || true)"
if [ -z "$PREV_TAG" ]; then
    echo "  ✗ no previous tag found; cannot compute commit range." >&2
    exit 1
fi

echo "  ✓ releasing ${PREV_TAG} → ${TAG}"

# ---------- bump JSON versions ----------
echo "→ Bumping versions to ${NEW_VERSION}"

bump_json() {
    local file="$1"
    local jq_expr="$2"
    local tmp
    tmp="$(mktemp)"
    jq --indent 2 --arg v "$NEW_VERSION" "$jq_expr" "$file" > "$tmp"
    mv "$tmp" "$file"
    echo "  ✓ ${file#$REPO_ROOT/}"
}

bump_json "$MARKETPLACE" '.metadata.version = $v'

shopt -s nullglob
for manifest in "$REPO_ROOT"/plugins/*/.claude-plugin/plugin.json; do
    bump_json "$manifest" '.version = $v'
done
shopt -u nullglob

# ---------- generate + prepend changelog entry ----------
echo "→ Generating CHANGELOG entry (claude -p, this may take a moment)"

ENTRY_FILE="$(mktemp)"
trap 'rm -f "$ENTRY_FILE"' EXIT

if ! bash "$SCRIPT_DIR/changelog-entry.sh" "$NEW_VERSION" "$PREV_TAG" > "$ENTRY_FILE"; then
    echo "  ✗ changelog-entry.sh failed." >&2
    exit 1
fi

if [ ! -s "$ENTRY_FILE" ]; then
    echo "  ✗ changelog-entry.sh produced empty output." >&2
    exit 1
fi

echo "  ✓ entry generated ($(wc -l < "$ENTRY_FILE" | tr -d ' ') lines)"

# Prepend entry before the first '## [' heading in CHANGELOG.md.
# Intro block (everything up to but not including the first release heading)
# is preserved; the new entry is inserted immediately above the previous top.
PREPEND_TMP="$(mktemp)"
awk -v entry_file="$ENTRY_FILE" '
    BEGIN { inserted = 0 }
    /^## \[/ && !inserted {
        while ((getline line < entry_file) > 0) print line
        close(entry_file)
        print ""
        inserted = 1
    }
    { print }
    END {
        if (!inserted) {
            # No prior release heading — append at EOF.
            while ((getline line < entry_file) > 0) print line
            close(entry_file)
        }
    }
' "$CHANGELOG" > "$PREPEND_TMP"
mv "$PREPEND_TMP" "$CHANGELOG"

echo "  ✓ CHANGELOG.md updated"

# ---------- confirm ----------
echo
echo "→ Review changes"
git --no-pager diff --stat
echo
echo "--- New CHANGELOG entry ---"
cat "$ENTRY_FILE"
echo "--- end entry ---"
echo

if [ "$AUTO_YES" -ne 1 ]; then
    printf "Proceed with commit, tag %s%s? [y/N] " "$TAG" "$([ $PUSH -eq 1 ] && echo ', and push' || echo '')"
    read -r reply
    case "$reply" in
        y|Y|yes|YES) ;;
        *) echo "Aborted. Changes are staged in the worktree; discard with: git checkout -- . && git checkout -- .claude-plugin plugins CHANGELOG.md"; exit 1 ;;
    esac
fi

# ---------- commit + tag + push ----------
echo "→ Committing, tagging, pushing"

git add .claude-plugin/marketplace.json plugins/*/.claude-plugin/plugin.json CHANGELOG.md
git commit -m "chore: Releases ${TAG}"
git tag -a "$TAG" -m "Release ${TAG}"

echo "  ✓ committed and tagged ${TAG}"

if [ "$PUSH" -eq 1 ]; then
    git push origin main
    git push origin "$TAG"
    echo "  ✓ pushed main + ${TAG} (release workflow will run)"
else
    echo "  ⓘ --no-push specified. Push manually with:"
    echo "      git push origin main && git push origin ${TAG}"
fi

echo
echo "✓ Release ${TAG} complete."
