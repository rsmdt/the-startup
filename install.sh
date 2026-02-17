#!/bin/sh
# The Agentic Startup - Installation Script
# https://github.com/rsmdt/the-startup
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh
#   ./install.sh [--help]

set -e

# -----------------------------------------------------------------------------
# Configuration
# -----------------------------------------------------------------------------

MARKETPLACE="rsmdt/the-startup"
PLUGINS="team@the-startup start@the-startup"
STATUSLINE_URL="https://raw.githubusercontent.com/rsmdt/the-startup/main/scripts/statusline.sh"
CONFIG_URL="https://raw.githubusercontent.com/rsmdt/the-startup/main/scripts/statusline.toml"
STATUSLINE_DIR="$HOME/.config/the-agentic-startup"
STATUSLINE_PATH="$STATUSLINE_DIR/statusline.sh"
CONFIG_PATH="$STATUSLINE_DIR/statusline.toml"
SETTINGS_FILE="$HOME/.claude/settings.json"
OUTPUT_STYLE="start:The Startup"

# -----------------------------------------------------------------------------
# Colors
# -----------------------------------------------------------------------------

GREEN='\033[0;32m'
BRIGHT_GREEN='\033[1;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
DIM='\033[2m'
RESET='\033[0m'

# -----------------------------------------------------------------------------
# Logging
# -----------------------------------------------------------------------------

info()    { printf "${DIM}→${RESET} %s\n" "$*"; }
warn()    { printf "${YELLOW}!${RESET} %s\n" "$*" >&2; }
error()   { printf "${RED}✗${RESET} %s\n" "$*" >&2; }
success() { printf "${GREEN}✓${RESET} %s\n" "$*"; }

# -----------------------------------------------------------------------------
# Banner
# -----------------------------------------------------------------------------

banner() {
  printf "${BRIGHT_GREEN}"
  cat << 'EOF'

████████ ██   ██ ███████
   ██    ██   ██ ██
   ██    ███████ █████
   ██    ██   ██ ██
   ██    ██   ██ ███████

 █████  ██████  ███████ ███   ██ ████████ ██  ██████
██   ██ ██      ██      ████  ██    ██    ██ ██
███████ ██  ███ █████   ██ ██ ██    ██    ██ ██
██   ██ ██   ██ ██      ██  ████    ██    ██ ██
██   ██  ██████ ███████ ██   ███    ██    ██  ██████

███████ ████████  █████  ██████  ████████ ██   ██ ██████
██         ██    ██   ██ ██   ██    ██    ██   ██ ██   ██
███████    ██    ███████ ██████     ██    ██   ██ ██████
     ██    ██    ██   ██ ██   ██    ██    ██   ██ ██
███████    ██    ██   ██ ██   ██    ██     █████  ██

EOF
  printf "${RESET}"
  echo "The framework for agentic software development"
  echo ""
}

# -----------------------------------------------------------------------------
# Functions
# -----------------------------------------------------------------------------

install() {
  # Check Claude CLI is available
  if ! command -v claude >/dev/null 2>&1; then
    error "Claude CLI is not installed"
    echo "  Install: curl -fsSL https://claude.ai/install.sh | sh"
    exit 1
  fi

  # Add or update marketplace
  info "Configuring marketplace..."
  if claude plugin marketplace add "$MARKETPLACE" >/dev/null 2>&1; then
    success "Marketplace added"
  else
    # Already exists - update it
    if claude plugin marketplace update "$MARKETPLACE" >/dev/null 2>&1; then
      success "Marketplace updated"
    else
      success "Marketplace configured"
    fi
  fi

  info "Installing plugins..."
  for plugin in $PLUGINS; do
    if claude plugin install "$plugin" >/dev/null 2>&1; then
      success "$plugin"
    else
      error "Failed to install $plugin"
      exit 2
    fi
  done
}

configure() {
  # Check jq is available
  if ! command -v jq >/dev/null 2>&1; then
    error "jq is not installed"
    echo "  macOS:  brew install jq"
    echo "  Ubuntu: sudo apt install jq"
    exit 1
  fi

  # Ensure settings directory exists
  mkdir -p "$(dirname "$SETTINGS_FILE")"

  # Create empty settings if missing
  if [ ! -f "$SETTINGS_FILE" ]; then
    echo '{}' > "$SETTINGS_FILE"
  fi

  # Validate existing JSON, backup and reset if invalid
  if ! jq empty "$SETTINGS_FILE" 2>/dev/null; then
    warn "settings.json malformed, creating backup"
    cp "$SETTINGS_FILE" "${SETTINGS_FILE}.bak"
    echo '{}' > "$SETTINGS_FILE"
  fi

  # Set outputStyle
  info "Configuring output style..."
  tmp_file=$(mktemp)
  if jq --arg style "$OUTPUT_STYLE" '.outputStyle = $style' "$SETTINGS_FILE" > "$tmp_file"; then
    mv "$tmp_file" "$SETTINGS_FILE"
    success "Output style: $OUTPUT_STYLE"
  else
    rm -f "$tmp_file"
    error "Failed to configure output style"
    exit 4
  fi
}

agent_teams() {
  # Check if agent teams env var is already configured
  if jq -e '.env.CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS == "1"' "$SETTINGS_FILE" >/dev/null 2>&1; then
    success "Agent Teams already enabled, excellent!"
    AGENT_TEAMS_CONFIGURED=true
    return
  fi

  # Ask user if they want to enable agent teams
  printf "\n"
  printf "Agent Teams ${DIM}(experimental)${RESET} enables multi-agent collaboration\n"
  printf "where specialized agents work together on complex tasks.\n"
  printf "\n"
  printf "Enable Agent Teams? [Y/n] "
  read -r answer </dev/tty

  # Default to yes
  case "$answer" in
    [nN]|[nN][oO])
      AGENT_TEAMS_CONFIGURED=false
      info "Skipping Agent Teams"
      ;;
    *)
      info "Enabling Agent Teams..."
      tmp_file=$(mktemp)
      if jq '.env = (.env // {}) + {"CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"}' "$SETTINGS_FILE" > "$tmp_file"; then
        mv "$tmp_file" "$SETTINGS_FILE"
        success "Agent Teams enabled"
        AGENT_TEAMS_CONFIGURED=true
      else
        rm -f "$tmp_file"
        error "Failed to enable Agent Teams"
        AGENT_TEAMS_CONFIGURED=false
      fi
      ;;
  esac
}

statusline() {
  # Check curl is available
  if ! command -v curl >/dev/null 2>&1; then
    error "curl is not installed"
    echo "  macOS:  curl is pre-installed"
    echo "  Ubuntu: sudo apt install curl"
    exit 1
  fi

  # Check jq is available
  if ! command -v jq >/dev/null 2>&1; then
    error "jq is not installed"
    echo "  macOS:  brew install jq"
    echo "  Ubuntu: sudo apt install jq"
    exit 1
  fi

  # Create target directory
  mkdir -p "$STATUSLINE_DIR"

  # Download statusline script
  info "Downloading statusline..."
  if ! curl -fsSL "$STATUSLINE_URL" -o "$STATUSLINE_PATH" 2>/dev/null; then
    error "Failed to download statusline"
    exit 3
  fi
  chmod +x "$STATUSLINE_PATH"

  # Download config if it doesn't exist (preserve user customizations)
  if [ ! -f "$CONFIG_PATH" ]; then
    info "Creating default config..."
    if curl -fsSL "$CONFIG_URL" -o "$CONFIG_PATH" 2>/dev/null; then
      success "Config created at $CONFIG_PATH"
    else
      warn "Could not download config (statusline will use defaults)"
    fi
  else
    info "Config exists, preserving customizations"
  fi

  # Configure statusLine in settings.json
  tmp_file=$(mktemp)
  if jq --arg cmd "$STATUSLINE_PATH" '.statusLine = {"type": "command", "command": $cmd}' "$SETTINGS_FILE" > "$tmp_file"; then
    mv "$tmp_file" "$SETTINGS_FILE"
    success "Statusline installed"
  else
    rm -f "$tmp_file"
    error "Failed to configure statusline"
    exit 4
  fi
}

usage() {
  cat << EOF
Usage: $0 [OPTIONS]

Options:
  -h, --help    Show this help message

Installs The Agentic Startup plugins, configures output style, and sets up the statusline.
EOF
}

# -----------------------------------------------------------------------------
# Argument Parsing
# -----------------------------------------------------------------------------

parse_args() {
  while [ $# -gt 0 ]; do
    case "$1" in
      -h|--help)
        usage
        exit 0
        ;;
      *)
        error "Unknown option: $1"
        usage
        exit 1
        ;;
    esac
  done
}

# -----------------------------------------------------------------------------
# Main
# -----------------------------------------------------------------------------

main() {
  parse_args "$@"

  AGENT_TEAMS_CONFIGURED=false

  banner
  install
  configure
  agent_teams
  statusline

  echo ""
  printf "${BRIGHT_GREEN}Installation complete!${RESET}\n"

  if [ "$AGENT_TEAMS_CONFIGURED" = false ]; then
    printf "\n"
    printf "${YELLOW}💡${RESET} You can enable Agent Teams later by adding to ${DIM}~/.claude/settings.json${RESET}:\n"
    printf "${DIM}   \"env\": { \"CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS\": \"1\" }${RESET}\n"
  fi

  printf "\n"
  printf "${DIM}Learn more: https://github.com/rsmdt/the-startup${RESET}\n"
}

main "$@"
