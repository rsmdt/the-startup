#!/usr/bin/env bash
#
# Ultra-fast statusline script for Claude Code hooks
# Extracts git branch from JSON input in <10ms
#
# Requirements:
# - Bash 3.2+ or Zsh compatible
# - Reads Claude Code JSON from stdin
# - Outputs git branch name or empty string
# - Performance: <10ms execution time
#
# Fast Path: Direct .git/HEAD file read (fastest)
# Fallback: git symbolic-ref command (slower but reliable)
#
# [ref: PRD; lines: 213-223]
# [ref: SDD; lines: 799-820, 235-240, 1153]

# Read JSON from stdin in one go
IFS= read -r -d '' json_input || true

# Extract "cwd" field using regex (fastest method, no external commands)
# Pattern: "cwd":"value" or "cwd": "value"
cwd=""
if [[ "$json_input" =~ \"cwd\"[[:space:]]*:[[:space:]]*\"([^\"]+)\" ]]; then
  cwd="${BASH_REMATCH[1]}"
fi

# Use current directory if cwd is empty
[[ -z "$cwd" ]] && cwd="."

# Expand tilde to home directory if present
[[ "$cwd" =~ ^~ ]] && cwd="${cwd/#\~/$HOME}"

# Fast path: Direct .git/HEAD file read
git_head="${cwd}/.git/HEAD"
if [[ -f "$git_head" && -r "$git_head" ]]; then
  # Read file content
  head_content=$(<"$git_head")

  # Extract branch from "ref: refs/heads/branch-name"
  if [[ "$head_content" =~ ^ref:[[:space:]]*refs/heads/(.+)$ ]]; then
    echo "${BASH_REMATCH[1]}"
    exit 0
  fi
fi

# Fallback: Use git command if available and in git repo
if command -v git &>/dev/null && [[ -d "${cwd}/.git" ]]; then
  branch=$(cd "$cwd" 2>/dev/null && git symbolic-ref --short HEAD 2>/dev/null || echo "")
  [[ -n "$branch" ]] && echo "$branch" && exit 0
fi

# No git repo or error - return empty string
echo ""
exit 0
