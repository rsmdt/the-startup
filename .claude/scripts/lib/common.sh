#!/usr/bin/env bash
# Common functions for Claude Code Hooks scripts

# Export common variables
export CLAUDE_HOOKS_VERSION="1.0.0"

# Color codes
export COLOR_RESET='\033[0m'
export COLOR_BLACK='\033[0;30m'
export COLOR_RED='\033[0;31m'
export COLOR_GREEN='\033[0;32m'
export COLOR_YELLOW='\033[0;33m'
export COLOR_BLUE='\033[0;34m'
export COLOR_MAGENTA='\033[0;35m'
export COLOR_CYAN='\033[0;36m'
export COLOR_WHITE='\033[0;37m'
export COLOR_BOLD='\033[1m'

# Log levels
export LOG_LEVEL_DEBUG=0
export LOG_LEVEL_INFO=1
export LOG_LEVEL_WARN=2
export LOG_LEVEL_ERROR=3

# Current log level (can be overridden)
export CURRENT_LOG_LEVEL=${CURRENT_LOG_LEVEL:-$LOG_LEVEL_INFO}

# Logging functions
log_debug() {
    if [ "$CURRENT_LOG_LEVEL" -le "$LOG_LEVEL_DEBUG" ]; then
        echo -e "${COLOR_CYAN}[DEBUG]${COLOR_RESET} $*" >&2
    fi
}

log_info() {
    if [ "$CURRENT_LOG_LEVEL" -le "$LOG_LEVEL_INFO" ]; then
        echo -e "${COLOR_GREEN}[INFO]${COLOR_RESET} $*" >&2
    fi
}

log_warn() {
    if [ "$CURRENT_LOG_LEVEL" -le "$LOG_LEVEL_WARN" ]; then
        echo -e "${COLOR_YELLOW}[WARN]${COLOR_RESET} $*" >&2
    fi
}

log_error() {
    if [ "$CURRENT_LOG_LEVEL" -le "$LOG_LEVEL_ERROR" ]; then
        echo -e "${COLOR_RED}[ERROR]${COLOR_RESET} $*" >&2
    fi
}

# Check if a command exists
command_exists() {
    command -v "$1" &> /dev/null
}

# Check if running on macOS
is_macos() {
    [[ "$OSTYPE" == "darwin"* ]]
}

# Check if running on Linux
is_linux() {
    [[ "$OSTYPE" == "linux-gnu"* ]]
}

# Check if running on Windows (Git Bash, WSL, etc.)
is_windows() {
    [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]
}

# Get the absolute path of a file or directory
get_absolute_path() {
    local path="$1"
    if [ -d "$path" ]; then
        (cd "$path" && pwd)
    elif [ -f "$path" ]; then
        echo "$(cd "$(dirname "$path")" && pwd)/$(basename "$path")"
    else
        echo "$path"
    fi
}

# Create a backup of a file
backup_file() {
    local file="$1"
    local backup_suffix="${2:-.bak}"
    
    if [ -f "$file" ]; then
        local timestamp=$(date +%Y%m%d_%H%M%S)
        local backup_file="${file}.${timestamp}${backup_suffix}"
        cp "$file" "$backup_file"
        log_debug "Created backup: $backup_file"
        echo "$backup_file"
    else
        log_error "File not found: $file"
        return 1
    fi
}

# Check if a file contains a string
file_contains() {
    local file="$1"
    local search_string="$2"
    
    if [ -f "$file" ]; then
        grep -q "$search_string" "$file"
    else
        return 1
    fi
}

# Get the size of a file in human-readable format
get_file_size() {
    local file="$1"
    
    if [ -f "$file" ]; then
        if is_macos; then
            stat -f%z "$file" | awk '{ 
                if ($1 < 1024) print $1 " B"
                else if ($1 < 1048576) printf "%.1f KB\n", $1/1024
                else if ($1 < 1073741824) printf "%.1f MB\n", $1/1048576
                else printf "%.1f GB\n", $1/1073741824
            }'
        else
            stat -c%s "$file" | awk '{ 
                if ($1 < 1024) print $1 " B"
                else if ($1 < 1048576) printf "%.1f KB\n", $1/1024
                else if ($1 < 1073741824) printf "%.1f MB\n", $1/1048576
                else printf "%.1f GB\n", $1/1073741824
            }'
        fi
    else
        echo "0 B"
    fi
}

# Count lines in a file
count_lines() {
    local file="$1"
    
    if [ -f "$file" ]; then
        wc -l < "$file"
    else
        echo "0"
    fi
}

# Check if JSON is valid
is_valid_json() {
    local json="$1"
    
    if command_exists jq; then
        echo "$json" | jq empty 2>/dev/null
        return $?
    else
        # Basic check if jq is not available
        echo "$json" | python3 -m json.tool > /dev/null 2>&1
        return $?
    fi
}

# Parse JSON value
get_json_value() {
    local json="$1"
    local key="$2"
    
    if command_exists jq; then
        echo "$json" | jq -r "$key"
    else
        echo "$json" | python3 -c "import sys, json; print(json.loads(sys.stdin.read())$key)"
    fi
}

# Check if a port is in use
is_port_in_use() {
    local port="$1"
    
    if is_macos; then
        lsof -i ":$port" &> /dev/null
    else
        netstat -tuln 2>/dev/null | grep -q ":$port "
    fi
}

# Generate a random string
generate_random_string() {
    local length="${1:-8}"
    
    if command_exists openssl; then
        openssl rand -hex "$((length/2))" | cut -c1-"$length"
    elif [ -f /dev/urandom ]; then
        tr -dc 'a-zA-Z0-9' < /dev/urandom | head -c "$length"
    else
        # Fallback to less random method
        echo "$RANDOM$RANDOM$RANDOM" | md5sum | cut -c1-"$length"
    fi
}

# Create a temporary file
create_temp_file() {
    local prefix="${1:-claude-hooks}"
    local suffix="${2:-.tmp}"
    
    if command_exists mktemp; then
        mktemp "/tmp/${prefix}.XXXXXX${suffix}"
    else
        local temp_file="/tmp/${prefix}.$$${suffix}"
        touch "$temp_file"
        echo "$temp_file"
    fi
}

# Clean up old files
cleanup_old_files() {
    local directory="$1"
    local days="${2:-7}"
    local pattern="${3:-*}"
    
    if [ -d "$directory" ]; then
        find "$directory" -name "$pattern" -type f -mtime +"$days" -delete
        log_debug "Cleaned up files older than $days days in $directory"
    fi
}

# Check system requirements
check_requirements() {
    local requirements=("$@")
    local missing=()
    
    for req in "${requirements[@]}"; do
        if ! command_exists "$req"; then
            missing+=("$req")
        fi
    done
    
    if [ ${#missing[@]} -gt 0 ]; then
        log_error "Missing required commands: ${missing[*]}"
        return 1
    fi
    
    return 0
}

# Format timestamp
format_timestamp() {
    local timestamp="${1:-$(date +%s)}"
    local format="${2:-%Y-%m-%d %H:%M:%S}"
    
    if is_macos; then
        date -r "$timestamp" +"$format"
    else
        date -d "@$timestamp" +"$format"
    fi
}

# Get relative time (e.g., "2 hours ago")
get_relative_time() {
    local timestamp="$1"
    local now=$(date +%s)
    local diff=$((now - timestamp))
    
    if [ $diff -lt 60 ]; then
        echo "$diff seconds ago"
    elif [ $diff -lt 3600 ]; then
        echo "$((diff / 60)) minutes ago"
    elif [ $diff -lt 86400 ]; then
        echo "$((diff / 3600)) hours ago"
    else
        echo "$((diff / 86400)) days ago"
    fi
}

# Export functions for use in other scripts
export -f log_debug log_info log_warn log_error
export -f command_exists is_macos is_linux is_windows
export -f get_absolute_path backup_file file_contains
export -f get_file_size count_lines
export -f is_valid_json get_json_value
export -f is_port_in_use generate_random_string
export -f create_temp_file cleanup_old_files
export -f check_requirements format_timestamp get_relative_time