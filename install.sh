#!/usr/bin/env bash
# The Startup - Agent System Installer
# Repository: https://github.com/the-startup/the-startup

set -e

# Configuration
REPO_OWNER="the-startup"
REPO_NAME="the-startup"
REPO_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}"
RAW_URL="https://raw.githubusercontent.com/${REPO_OWNER}/${REPO_NAME}/main"
CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/the-startup"
VERSION="1.0.0"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Global variables
INSTALL_TYPE=""
INSTALL_DIR=""
TOOL_TYPE=""
UPDATE_MODE=false
HAS_GUM=false
HAS_HUH=false
EXISTING_LOCK_FILE=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --update|-u)
            UPDATE_MODE=true
            shift
            ;;
        --help|-h)
            echo "The Startup Installer"
            echo ""
            echo "Usage: curl -LsSf ${RAW_URL}/install.sh | sh"
            echo "       curl -LsSf ${RAW_URL}/install.sh | sh -s -- --update"
            echo ""
            echo "Options:"
            echo "  --update, -u    Update existing installation"
            echo "  --help, -h      Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Helper functions
print_header() {
    echo -e "\n${CYAN}╔════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║${NC}  ${BOLD}$1${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════╝${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

confirm() {
    local prompt="$1"
    if [ "$HAS_GUM" = true ]; then
        gum confirm "$prompt"
    else
        read -p "$prompt (y/n) " -n 1 -r
        echo
        [[ $REPLY =~ ^[Yy]$ ]]
    fi
}

choose() {
    local prompt="$1"
    shift
    if [ "$HAS_GUM" = true ]; then
        gum choose "$@"
    else
        echo "$prompt"
        select choice in "$@"; do
            echo "$choice"
            break
        done
    fi
}

input() {
    local prompt="$1"
    local default="$2"
    if [ "$HAS_GUM" = true ]; then
        if [ -n "$default" ]; then
            gum input --placeholder "$prompt" --value "$default"
        else
            gum input --placeholder "$prompt"
        fi
    else
        if [ -n "$default" ]; then
            read -p "$prompt [$default]: " value
            echo "${value:-$default}"
        else
            read -p "$prompt: " value
            echo "$value"
        fi
    fi
}

# Check for prerequisites
check_prerequisites() {
    print_header "🔍 Checking Prerequisites"
    
    # Check for curl
    if ! command -v curl &> /dev/null; then
        print_error "curl is required but not installed"
        exit 1
    fi
    print_success "curl found"
    
    # Check for gum (optional but recommended)
    if command -v gum &> /dev/null; then
        HAS_GUM=true
        print_success "gum found (enhanced UI enabled)"
    else
        print_info "gum not found (using basic UI)"
        print_info "Install gum for better experience: https://github.com/charmbracelet/gum"
    fi
    
    # Check for huh (for multi-select)
    if command -v huh &> /dev/null; then
        HAS_HUH=true
        print_success "huh found (multi-select enabled)"
    fi
    
    # Check for jq (for JSON processing)
    if ! command -v jq &> /dev/null; then
        print_warning "jq not found - some features may be limited"
        print_info "Install jq for full functionality"
    fi
    
    # Check for Python (for hooks)
    if ! command -v python3 &> /dev/null; then
        print_warning "Python 3 not found - hooks will require Python"
    else
        print_success "Python 3 found"
    fi
    
    # Check for uv (for Python package management)
    if ! command -v uv &> /dev/null; then
        print_warning "uv not found - Python hooks will need uv"
        if confirm "Would you like to install uv now?"; then
            curl -LsSf https://astral.sh/uv/install.sh | sh
            export PATH="$HOME/.local/bin:$PATH"
            if command -v uv &> /dev/null; then
                print_success "uv installed successfully"
            fi
        fi
    else
        print_success "uv found"
    fi
}

# Select tool type
select_tool() {
    print_header "🛠️  Select Tool"
    
    TOOL_TYPE=$(choose "Select the code tool you're using:" \
        "Claude Code" \
        "Cancel")
    
    if [ "$TOOL_TYPE" = "Cancel" ]; then
        print_info "Installation cancelled"
        exit 0
    fi
    
    print_success "Selected: $TOOL_TYPE"
}

# Select installation type
select_install_type() {
    print_header "📍 Select Installation Location"
    
    echo "Where would you like to install the agent system?"
    echo ""
    echo "  ${BOLD}Global${NC}: Install to ~/.claude (available for all projects)"
    echo "  ${BOLD}Local${NC}:  Install to ./.claude (current directory only)"
    echo ""
    
    INSTALL_TYPE=$(choose "Select installation type:" \
        "Global (~/.claude)" \
        "Local (./.claude)" \
        "Cancel")
    
    case "$INSTALL_TYPE" in
        "Global (~/.claude)")
            INSTALL_TYPE="global"
            INSTALL_DIR="$HOME/.claude"
            ;;
        "Local (./.claude)")
            INSTALL_TYPE="local"
            INSTALL_DIR="$(pwd)/.claude"
            ;;
        *)
            print_info "Installation cancelled"
            exit 0
            ;;
    esac
    
    print_success "Installation directory: $INSTALL_DIR"
}

# Check for existing installation
check_existing_installation() {
    # Check for lock file in config directory
    if [ -f "$CONFIG_DIR/the-startup.lock" ]; then
        EXISTING_LOCK_FILE="$CONFIG_DIR/the-startup.lock"
        print_info "Found existing installation (lock file: $CONFIG_DIR/the-startup.lock)"
        return 0
    fi
    
    # Check for existing agents/hooks/commands
    local has_existing=false
    if [ -d "$INSTALL_DIR/agents" ] || [ -d "$INSTALL_DIR/hooks" ] || [ -d "$INSTALL_DIR/commands" ]; then
        has_existing=true
    fi
    
    if [ "$has_existing" = true ]; then
        print_warning "Found existing files in $INSTALL_DIR"
        if ! confirm "Continue with installation?"; then
            exit 0
        fi
    fi
    
    return 1
}

# Download file from GitHub
download_file() {
    local remote_path="$1"
    local local_path="$2"
    local file_url="${RAW_URL}/${remote_path}"
    
    # Create directory if needed
    mkdir -p "$(dirname "$local_path")"
    
    # Download file
    if curl -fsSL "$file_url" -o "$local_path" 2>/dev/null; then
        return 0
    else
        return 1
    fi
}

# Get list of available files from GitHub
get_available_files() {
    local component_type="$1"  # agents, hooks, commands
    
    # For now, we'll use a predefined list
    # In production, this would query GitHub API
    case "$component_type" in
        agents)
            echo "the-architect.md the-business-analyst.md the-chief.md the-data-engineer.md the-developer.md the-devops-engineer.md the-product-manager.md the-project-manager.md the-security-engineer.md the-site-reliability-engineer.md the-technical-writer.md the-tester.md"
            ;;
        hooks)
            echo "log_agent_start.py log_agent_complete.py"
            ;;
        commands)
            echo "develop.md start.md"
            ;;
        rules)
            echo "context-management.md"
            ;;
        *)
            echo ""
            ;;
    esac
}

# Select components to install
select_components() {
    print_header "📦 Select Components"
    
    local selected_agents=""
    local selected_hooks=""
    local selected_commands=""
    
    # Select agents
    print_info "Available agents:"
    local available_agents=$(get_available_files "agents")
    
    if [ "$HAS_HUH" = true ] && [ "$UPDATE_MODE" = true ]; then
        # Use huh for multi-select during updates
        selected_agents=$(echo "$available_agents" | tr ' ' '\n' | huh multi-select --header "Select agents to install/update")
    else
        # For initial install, select all by default
        if confirm "Install all agents?"; then
            selected_agents="$available_agents"
        else
            # Manual selection
            for agent in $available_agents; do
                if confirm "Install $agent?"; then
                    selected_agents="$selected_agents $agent"
                fi
            done
        fi
    fi
    
    # Select hooks
    print_info "Available hooks:"
    local available_hooks=$(get_available_files "hooks")
    
    if confirm "Install hooks for agent logging?"; then
        selected_hooks="$available_hooks"
    fi
    
    # Select commands
    print_info "Available commands:"
    local available_commands=$(get_available_files "commands")
    
    if confirm "Install all commands?"; then
        selected_commands="$available_commands"
    fi
    
    # Return selections
    echo "$selected_agents|$selected_hooks|$selected_commands"
}

# Install components
install_components() {
    local components="$1"
    
    IFS='|' read -r agents hooks commands <<< "$components"
    
    print_header "📥 Installing Components"
    
    # Install agents
    if [ -n "$agents" ]; then
        print_info "Installing agents..."
        for agent in $agents; do
            if download_file "agents/$agent" "$INSTALL_DIR/agents/$agent"; then
                print_success "Installed: agents/$agent"
            else
                print_error "Failed to install: agents/$agent"
            fi
        done
    fi
    
    # Install hooks
    if [ -n "$hooks" ]; then
        print_info "Installing hooks..."
        for hook in $hooks; do
            if download_file "hooks/$hook" "$INSTALL_DIR/hooks/$hook"; then
                chmod +x "$INSTALL_DIR/hooks/$hook"
                print_success "Installed: hooks/$hook"
            else
                print_error "Failed to install: hooks/$hook"
            fi
        done
    fi
    
    # Install commands
    if [ -n "$commands" ]; then
        print_info "Installing commands..."
        for command in $commands; do
            if download_file "commands/$command" "$INSTALL_DIR/commands/$command"; then
                print_success "Installed: commands/$command"
            else
                print_error "Failed to install: commands/$command"
            fi
        done
    fi
    
    # Install rules (always install)
    print_info "Installing rules..."
    for rule in $(get_available_files "rules"); do
        if download_file "rules/$rule" "$INSTALL_DIR/rules/$rule"; then
            print_success "Installed: rules/$rule"
        else
            print_error "Failed to install: rules/$rule"
        fi
    done
}

# Configure hooks in settings
configure_hooks() {
    print_header "⚙️  Configuring Hooks"
    
    local settings_file="$INSTALL_DIR/settings.json"
    
    # Check if settings file exists
    if [ ! -f "$settings_file" ]; then
        print_info "Creating new settings.json"
        echo '{}' > "$settings_file"
    fi
    
    # Update settings with Python (more reliable JSON handling)
    python3 - <<EOF
import json
import os

settings_file = "$settings_file"

# Read existing settings
with open(settings_file, 'r') as f:
    settings = json.load(f)

# Ensure hooks section exists
if 'hooks' not in settings:
    settings['hooks'] = {}

# Configure PreToolUse hooks
if 'PreToolUse' not in settings['hooks']:
    settings['hooks']['PreToolUse'] = []

# Find existing the-startup hooks
startup_hooks = [h for h in settings['hooks']['PreToolUse'] 
                 if not any(hook.get('_source') == 'the-startup' 
                           for hook in h.get('hooks', []))]

# Add our hook
startup_hook = {
    "matcher": "Task",
    "hooks": [{
        "type": "command",
        "command": "uv run \$CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_start.py",
        "_source": "the-startup"
    }]
}

settings['hooks']['PreToolUse'] = startup_hooks + [startup_hook]

# Configure PostToolUse hooks
if 'PostToolUse' not in settings['hooks']:
    settings['hooks']['PostToolUse'] = []

# Find existing the-startup hooks
startup_hooks = [h for h in settings['hooks']['PostToolUse'] 
                 if not any(hook.get('_source') == 'the-startup' 
                           for hook in h.get('hooks', []))]

# Add our hook
startup_hook = {
    "matcher": "Task",
    "hooks": [{
        "type": "command",
        "command": "uv run \$CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_complete.py",
        "_source": "the-startup"
    }]
}

settings['hooks']['PostToolUse'] = startup_hooks + [startup_hook]

# Write updated settings
with open(settings_file, 'w') as f:
    json.dump(settings, f, indent=2)

print("Hooks configured successfully")
EOF
    
    if [ $? -eq 0 ]; then
        print_success "Hooks configured in settings.json"
    else
        print_error "Failed to configure hooks"
    fi
}

# Create lock file
create_lock_file() {
    print_header "🔒 Creating Lock File"
    
    mkdir -p "$CONFIG_DIR"
    
    local lock_file="$CONFIG_DIR/the-startup.lock"
    
    # Create lock file with Python
    python3 - <<EOF
import json
import os
from datetime import datetime

lock_data = {
    "version": "$VERSION",
    "install_date": datetime.utcnow().isoformat() + "Z",
    "install_type": "$INSTALL_TYPE",
    "install_path": "$INSTALL_DIR",
    "tool": "$TOOL_TYPE",
    "installed_files": {}
}

# List installed files
for root, dirs, files in os.walk("$INSTALL_DIR"):
    for file in files:
        if file.endswith(('.md', '.py', '.sh')):
            rel_path = os.path.relpath(os.path.join(root, file), "$INSTALL_DIR")
            lock_data["installed_files"][rel_path] = {
                "version": "$VERSION",
                "last_modified": datetime.utcnow().isoformat() + "Z"
            }

# Write lock file
with open("$lock_file", 'w') as f:
    json.dump(lock_data, f, indent=2)

print(f"Lock file created: {lock_file}")
EOF
    
    if [ $? -eq 0 ]; then
        print_success "Lock file created: $lock_file"
    else
        print_error "Failed to create lock file"
    fi
}

# Main installation flow
main() {
    clear
    
    # Welcome message
    echo -e "${BOLD}${MAGENTA}"
    echo "╔══════════════════════════════════════════╗"
    echo "║                                          ║"
    echo "║        🚀 THE STARTUP INSTALLER 🚀       ║"
    echo "║                                          ║"
    echo "║    Agent System for Development Tools    ║"
    echo "║                                          ║"
    echo "╚══════════════════════════════════════════╝"
    echo -e "${NC}"
    
    # Check prerequisites
    check_prerequisites
    
    # Select tool
    select_tool
    
    # Select installation type
    select_install_type
    
    # Check for existing installation
    if check_existing_installation; then
        UPDATE_MODE=true
        print_info "Switching to update mode"
    fi
    
    # Select components
    components=$(select_components)
    
    # Install components
    install_components "$components"
    
    # Configure hooks
    if [[ "$components" == *"log_agent"* ]]; then
        configure_hooks
    fi
    
    # Create/update lock file
    create_lock_file
    
    # Success message
    print_header "✨ Installation Complete!"
    
    print_success "The Startup agent system has been installed successfully!"
    echo ""
    print_info "Installation details:"
    echo "  • Type: $INSTALL_TYPE"
    echo "  • Location: $INSTALL_DIR"
    echo "  • Tool: $TOOL_TYPE"
    echo "  • Config: $CONFIG_DIR"
    echo ""
    print_info "Next steps:"
    echo "  1. Restart your $TOOL_TYPE session"
    echo "  2. The agents and commands are now available"
    echo "  3. Check logs in: ~/.the-startup/ (if hooks installed)"
    echo ""
    print_info "To update in the future, run:"
    echo "  curl -LsSf ${RAW_URL}/install.sh | sh -s -- --update"
    echo ""
    
    if [ "$HAS_GUM" = true ]; then
        gum style \
            --border normal \
            --margin "1" \
            --padding "1" \
            --border-foreground 212 \
            "Thank you for installing The Startup! 🎉"
    else
        echo -e "${GREEN}Thank you for installing The Startup! 🎉${NC}"
    fi
}

# Run main function
main "$@"