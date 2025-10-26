#!/usr/bin/env bash
#
# Complete statusline script for Claude Code - Shell implementation
# Replicates the functionality of the Go implementation from rsmdt/the-startup
#
# Features:
# - Shows current directory (with ~ for home)
# - Shows git branch (if in a git repo)
# - Shows model name and output style
# - Shows help text with styling
# - Terminal width aware
# - ANSI color support
#
# Input: JSON from Claude Code via stdin
# Output: Single formatted statusline with ANSI colors
#
# Performance target: <50ms execution time
#

# ANSI color codes
# Main text color: #FAFAFA (very light gray/white)
MAIN_COLOR="\033[38;2;250;250;250m"
# Help text color: #606060 (gray, muted)
HELP_COLOR="\033[38;2;96;96;96m"
# Italic style
ITALIC="\033[3m"
# Reset all styles
RESET="\033[0m"

# Read JSON from stdin in one go
IFS= read -r -d '' json_input || true

# Extract fields from JSON using regex (no jq dependency for speed)
# Pattern: "field": "value" or "field":"value"

# Extract current_dir from workspace.current_dir
current_dir=""
if [[ "$json_input" =~ \"workspace\"[^}]*\"current_dir\"[[:space:]]*:[[:space:]]*\"([^\"]+)\" ]]; then
  current_dir="${BASH_REMATCH[1]}"
fi

# Fallback to cwd if current_dir not found
if [[ -z "$current_dir" && "$json_input" =~ \"cwd\"[[:space:]]*:[[:space:]]*\"([^\"]+)\" ]]; then
  current_dir="${BASH_REMATCH[1]}"
fi

# Use current directory if still empty
[[ -z "$current_dir" ]] && current_dir="$PWD"

# Extract model display_name
model_name=""
if [[ "$json_input" =~ \"model\"[^}]*\"display_name\"[[:space:]]*:[[:space:]]*\"([^\"]+)\" ]]; then
  model_name="${BASH_REMATCH[1]}"
fi
[[ -z "$model_name" ]] && model_name="Claude"

# Extract output_style name
output_style=""
if [[ "$json_input" =~ \"output_style\"[^}]*\"name\"[[:space:]]*:[[:space:]]*\"([^\"]+)\" ]]; then
  output_style="${BASH_REMATCH[1]}"
fi
[[ -z "$output_style" ]] && output_style="default"

# Home directory substitution
# Replace /Users/username or /home/username with ~
home_dir="$HOME"
if [[ -n "$home_dir" && "$current_dir" == "$home_dir" ]]; then
  # Exact match: /Users/username -> ~
  current_dir="~"
elif [[ -n "$home_dir" && "$current_dir" == "$home_dir"/* ]]; then
  # Prefix match: /Users/username/Documents -> ~/Documents
  current_dir="~${current_dir#$home_dir}"
fi

# Get git branch information
get_git_branch() {
  local dir="$1"

  # Expand tilde to home directory if present
  [[ "$dir" =~ ^~ ]] && dir="${dir/#\~/$HOME}"

  # Fast path: Direct .git/HEAD file read
  local git_head="${dir}/.git/HEAD"
  if [[ -f "$git_head" && -r "$git_head" ]]; then
    # Read file content
    local head_content
    head_content=$(<"$git_head")

    # Extract branch from "ref: refs/heads/branch-name"
    if [[ "$head_content" =~ ^ref:[[:space:]]*refs/heads/(.+)$ ]]; then
      echo "${BASH_REMATCH[1]}"
      return 0
    fi

    # If HEAD is detached, return HEAD
    echo "HEAD"
    return 0
  fi

  # Fallback: Use git command if available and in git repo
  if command -v git &>/dev/null && [[ -d "${dir}/.git" ]]; then
    local branch
    branch=$(cd "$dir" 2>/dev/null && git symbolic-ref --short HEAD 2>/dev/null || echo "")
    if [[ -n "$branch" ]]; then
      echo "$branch"
      return 0
    fi
    # Check if in detached HEAD state
    if (cd "$dir" 2>/dev/null && git rev-parse --git-dir &>/dev/null); then
      echo "HEAD"
      return 0
    fi
  fi

  # No git repo
  echo ""
}

# Get git info with branch symbol
git_branch=$(get_git_branch "$current_dir")
git_info=""
if [[ -n "$git_branch" ]]; then
  git_info="⎇ $git_branch"
fi

# Get terminal width
get_term_width() {
  local width

  # Method 1: COLUMNS environment variable (most reliable in hooks/scripts)
  if [[ -n "$COLUMNS" && "$COLUMNS" =~ ^[0-9]+$ && "$COLUMNS" -gt 0 ]]; then
    echo "$COLUMNS"
    return 0
  fi

  # Method 2: tput cols command (if available)
  if command -v tput &>/dev/null; then
    width=$(tput cols 2>/dev/null)
    if [[ -n "$width" && "$width" =~ ^[0-9]+$ && "$width" -gt 0 ]]; then
      echo "$width"
      return 0
    fi
  fi

  # Method 3: stty size command (if available)
  if command -v stty &>/dev/null; then
    local size
    size=$(stty size 2>/dev/null)
    if [[ -n "$size" ]]; then
      width=$(echo "$size" | cut -d' ' -f2)
      if [[ -n "$width" && "$width" =~ ^[0-9]+$ && "$width" -gt 0 ]]; then
        echo "$width"
        return 0
      fi
    fi
  fi

  # Default fallback
  echo "120"
}

term_width=$(get_term_width)

# Build the statusline parts
# Format: 📁 <dir> <git>  🤖 <model> (<style>)  ? for shortcuts

# Part 1: Directory with git info
if [[ -n "$git_info" ]]; then
  dir_part="📁 $current_dir $git_info"
else
  dir_part="📁 $current_dir"
fi

# Part 2: Model and output style
model_part="🤖 $model_name ($output_style)"

# Part 3: Help text (styled differently)
help_text="? for shortcuts"

# Calculate spacing (2 spaces between each part)
# We need to account for emoji character widths and ANSI codes
# For simplicity, we'll use fixed spacing like the Go implementation

# Build the complete statusline
# Main parts with MAIN_COLOR, help text with HELP_COLOR and italic
statusline="${MAIN_COLOR}${dir_part}  ${model_part}  ${HELP_COLOR}${ITALIC}${help_text}${RESET}"

# Output the statusline
echo -e "$statusline"

exit 0
