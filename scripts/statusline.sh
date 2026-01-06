#!/usr/bin/env bash
#
# Statusline script for Claude Code
#
# Usage: statusline.sh [-f "<format>"] [--help]
#
# Format placeholders:
#   <path>    - Directory path (abbreviated)
#   <branch>  - Git branch with icon (* if dirty)
#   <model>   - Model name and output style
#   <context> - Context usage bar and percentage
#   <session> - Session duration and cost
#   <lines>   - Lines added/removed
#   <spec>    - Active specification ID
#   <help>    - Help text
#
# Default format: "<path> <branch>  <model>  <context>  <session>  <help>"
#
# Examples:
#   statusline.sh -f "<path> | <context> <model>"
#   statusline.sh -f "<path> <branch>  <spec>  <model>  <context>  <lines>"
#
# Input: JSON from Claude Code via stdin
# Output: Single formatted statusline with ANSI colors
#
# Performance target: <50ms execution time
# Dependencies: jq

# ==============================================================================
# Constants
# ==============================================================================

readonly COLOR_DEFAULT="\033[38;2;250;250;250m"
readonly COLOR_MUTED="\033[38;2;96;96;96m"
readonly COLOR_WARNING="\033[38;2;255;184;0m"
readonly COLOR_DANGER="\033[38;2;255;68;68m"
readonly COLOR_SUCCESS="\033[38;2;136;204;136m"
readonly STYLE_ITALIC="\033[3m"
readonly STYLE_RESET="\033[0m"

readonly BRAILLE_CHARS=("⠀" "⡀" "⡄" "⡆" "⡇" "⣇" "⣧" "⣷" "⣿")

readonly DEFAULT_FORMAT="<path> <branch>  <model>  <context>  <session>  <help>"

readonly VALID_PLACEHOLDERS="path|branch|model|context|session|lines|spec|help"

# ==============================================================================
# Help
# ==============================================================================

show_help() {
  cat << 'EOF'
Statusline script for Claude Code

Usage: statusline.sh [-f "<format>"] [--help]

Format placeholders:
  <path>    - Directory path (abbreviated, e.g., ~/C/p/project)
  <branch>  - Git branch with icon, * if dirty (e.g., ⎇ main*)
  <model>   - Model name and output style (e.g., 🤖 Opus (The ScaleUp))
  <context> - Context usage bar and percentage (e.g., 🧠 ⣿⣿⡇⠀⠀ 50%)
  <session> - Session duration and cost (e.g., 🕐 30m  $0.50)
  <lines>   - Lines added/removed (e.g., +156/-23)
  <spec>    - Active specification ID (e.g., 📋 005)
  <help>    - Help text (? for shortcuts)

Default format:
  "<path> <branch>  <model>  <context>  <session>  <help>"

Examples:
  statusline.sh -f "<path> | <context> <model>"
  statusline.sh -f "<path> <branch>  <spec>  <model>  <context>  <lines>"
  statusline.sh -f "<context>"

Input: JSON from Claude Code via stdin
EOF
  exit 0
}

# ==============================================================================
# Input Parsing
# ==============================================================================

IFS= read -r -d '' json_input || true

{
  read -r current_dir
  read -r model_name
  read -r output_style
  read -r context_size
  read -r context_used
  read -r session_cost
  read -r session_duration_ms
  read -r lines_added
  read -r lines_removed
} <<< "$(echo "$json_input" | jq -r '
  (.workspace.current_dir // .cwd),
  .model.display_name,
  (.output_style.name | split(":") | .[-1]),
  .context_window.context_window_size,
  ((.context_window.current_usage.input_tokens // 0) + (.context_window.current_usage.cache_creation_input_tokens // 0) + (.context_window.current_usage.cache_read_input_tokens // 0)),
  .cost.total_cost_usd,
  .cost.total_duration_ms,
  .cost.total_lines_added,
  .cost.total_lines_removed
' 2>/dev/null)"

# ==============================================================================
# Formatters
# ==============================================================================

format_path() {
  local path="$current_dir"

  if [[ "$path" == "$HOME" ]]; then
    path="~"
  elif [[ "$path" == "$HOME"/* ]]; then
    path="~${path#$HOME}"
  fi

  local prefix=""
  if [[ "$path" == ~* ]]; then
    prefix="~"
    path="${path:1}"
  fi

  [[ "$path" == /* ]] && path="${path:1}"

  local IFS='/'
  read -ra parts <<< "$path"
  local count=${#parts[@]}
  local result=""

  for ((i = 0; i < count; i++)); do
    local part="${parts[$i]}"
    [[ -z "$part" ]] && continue

    if [[ $i -lt $((count - 1)) ]]; then
      result+="/${part:0:1}"
    else
      result+="/${part}"
    fi
  done

  echo "📁 ${prefix}${result}"
}

format_branch() {
  local dir="$current_dir"
  local branch=""
  local dirty=""

  [[ "$dir" == "$HOME" ]] && dir="~"
  [[ "$dir" == "$HOME"/* ]] && dir="~${dir#$HOME}"
  [[ "$dir" =~ ^~ ]] && dir="${dir/#\~/$HOME}"

  # Get branch name
  local git_head="${dir}/.git/HEAD"
  if [[ -f "$git_head" && -r "$git_head" ]]; then
    local head_content
    head_content=$(<"$git_head")

    if [[ "$head_content" =~ ^ref:[[:space:]]*refs/heads/(.+)$ ]]; then
      branch="${BASH_REMATCH[1]}"
    else
      branch="HEAD"
    fi
  elif command -v git &>/dev/null && [[ -d "${dir}/.git" ]]; then
    branch=$(cd "$dir" 2>/dev/null && git symbolic-ref --short HEAD 2>/dev/null || echo "")
    [[ -z "$branch" ]] && branch="HEAD"
  fi

  [[ -z "$branch" ]] && return

  # Check for dirty state (uncommitted changes)
  if [[ -d "${dir}/.git" ]] && command -v git &>/dev/null; then
    if ! (cd "$dir" 2>/dev/null && git diff --quiet HEAD 2>/dev/null); then
      dirty="*"
    fi
  fi

  echo "⎇ ${branch}${dirty}"
}

format_context() {
  [[ -z "$context_size" || "$context_size" == "null" || "$context_size" -eq 0 ]] && return

  local percent=$(((context_used + 45000) * 100 / context_size))  # Include 45k compaction buffer
  [[ "$percent" -gt 100 ]] && percent=100

  local bar=""
  local width=5
  local total_units=$((width * 8))
  local filled_units=$((percent * total_units / 100))

  for ((i = 0; i < width; i++)); do
    local char_fill=$((filled_units - (i * 8)))
    [[ "$char_fill" -lt 0 ]] && char_fill=0
    [[ "$char_fill" -gt 8 ]] && char_fill=8
    bar+="${BRAILLE_CHARS[$char_fill]}"
  done

  local color="$COLOR_DEFAULT"
  [[ "$percent" -ge 70 ]] && color="$COLOR_WARNING"
  [[ "$percent" -ge 90 ]] && color="$COLOR_DANGER"

  echo "🧠 ${color}${bar} ${percent}%${STYLE_RESET}"
}

format_duration() {
  local ms="$1"

  [[ -z "$ms" || "$ms" == "null" || "$ms" -eq 0 ]] && return

  local total_seconds=$((ms / 1000))
  local hours=$((total_seconds / 3600))
  local minutes=$(((total_seconds % 3600) / 60))

  if [[ "$hours" -gt 0 ]]; then
    if [[ "$minutes" -gt 0 ]]; then
      echo "${hours}h ${minutes}m"
    else
      echo "${hours}h"
    fi
  elif [[ "$minutes" -gt 0 ]]; then
    echo "${minutes}m"
  else
    echo "<1m"
  fi
}

format_session() {
  local result=""

  if [[ -n "$session_duration_ms" && "$session_duration_ms" != "null" && "$session_duration_ms" -gt 0 ]]; then
    result="🕐 $(format_duration "$session_duration_ms")"
  fi

  if [[ -n "$session_cost" && "$session_cost" != "null" ]]; then
    local formatted_cost
    formatted_cost=$(printf "%.2f" "$session_cost")
    if [[ -n "$result" ]]; then
      result+="  ${COLOR_SUCCESS}\$${formatted_cost}${STYLE_RESET}"
    else
      result="${COLOR_SUCCESS}\$${formatted_cost}${STYLE_RESET}"
    fi
  fi

  echo "$result"
}

format_lines() {
  [[ -z "$lines_added" || "$lines_added" == "null" ]] && return
  [[ -z "$lines_removed" || "$lines_removed" == "null" ]] && return
  [[ "$lines_added" -eq 0 && "$lines_removed" -eq 0 ]] && return

  echo "${COLOR_SUCCESS}+${lines_added}${STYLE_RESET}/${COLOR_DANGER}-${lines_removed}${STYLE_RESET}"
}

format_spec() {
  local dir="$current_dir"

  [[ "$dir" == "$HOME"/* ]] && dir="${dir#$HOME}"
  [[ "$dir" == "~"/* ]] && dir="${dir:1}"

  # Match docs/specs/NNN-* pattern anywhere in path
  if [[ "$dir" =~ docs/specs/([0-9]+)- ]]; then
    echo "📋 ${BASH_REMATCH[1]}"
  fi
}

# ==============================================================================
# Entry Point
# ==============================================================================

main() {
  local format="$DEFAULT_FORMAT"

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        ;;
      -f|--format)
        format="$2"
        shift 2
        ;;
      *)
        shift
        ;;
    esac
  done

  # Compute only needed parts
  local path_part branch_part model_part context_part session_part lines_part spec_part help_part
  [[ "$format" == *"<path>"* ]] && path_part=$(format_path)
  [[ "$format" == *"<branch>"* ]] && branch_part=$(format_branch)
  [[ "$format" == *"<model>"* ]] && model_part="🤖 ${model_name} (${output_style})"
  [[ "$format" == *"<context>"* ]] && context_part=$(format_context)
  [[ "$format" == *"<session>"* ]] && session_part=$(format_session)
  [[ "$format" == *"<lines>"* ]] && lines_part=$(format_lines)
  [[ "$format" == *"<spec>"* ]] && spec_part=$(format_spec)
  [[ "$format" == *"<help>"* ]] && help_part="${COLOR_MUTED}${STYLE_ITALIC}? for shortcuts${STYLE_RESET}"

  # Warn about unknown placeholders
  local unknown
  unknown=$(echo "$format" | grep -oE '<[a-z]+>' | grep -vE "<(${VALID_PLACEHOLDERS})>" | tr '\n' ' ')
  [[ -n "$unknown" ]] && echo "Warning: Unknown placeholders: $unknown" >&2

  # Build statusline by replacing placeholders
  local statusline="$format"
  statusline="${statusline//<path>/$path_part}"
  statusline="${statusline//<branch>/$branch_part}"
  statusline="${statusline//<model>/$model_part}"
  statusline="${statusline//<context>/$context_part}"
  statusline="${statusline//<session>/$session_part}"
  statusline="${statusline//<lines>/$lines_part}"
  statusline="${statusline//<spec>/$spec_part}"
  statusline="${statusline//<help>/$help_part}"

  # Clean up extra spaces from empty placeholders
  while [[ "$statusline" == *"   "* ]]; do
    statusline="${statusline//   /  }"
  done
  statusline="${statusline#"${statusline%%[![:space:]]*}"}"
  statusline="${statusline%"${statusline##*[![:space:]]}"}"

  echo -e "${STYLE_RESET}${statusline}"
}

main "$@"
