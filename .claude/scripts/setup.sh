#!/usr/bin/env bash
# Claude Code Hooks Interactive Setup
# Uses Gum for beautiful CLI interactions

set -e

# Colors and styling
CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Project paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLAUDE_DIR="$PROJECT_ROOT/.claude"
HOOKS_DIR="$CLAUDE_DIR/hooks"
STARTUP_DIR="$PROJECT_ROOT/.the-startup"

# Source common functions
source "$SCRIPT_DIR/lib/common.sh" 2>/dev/null || true

# Check if running with Gum
HAS_GUM=false
if command -v gum &> /dev/null; then
    HAS_GUM=true
fi

# Fallback functions for when Gum is not available
print_header() {
    if [ "$HAS_GUM" = true ]; then
        gum style \
            --border double \
            --margin "1" \
            --padding "1 2" \
            --border-foreground "#FF06B7" \
            "$1"
    else
        echo -e "\n${CYAN}═══════════════════════════════════════${NC}"
        echo -e "${CYAN}  $1${NC}"
        echo -e "${CYAN}═══════════════════════════════════════${NC}\n"
    fi
}

print_success() {
    if [ "$HAS_GUM" = true ]; then
        gum style --foreground 212 "✓ $1"
    else
        echo -e "${GREEN}✓ $1${NC}"
    fi
}

print_error() {
    if [ "$HAS_GUM" = true ]; then
        gum style --foreground 196 "✗ $1"
    else
        echo -e "${RED}✗ $1${NC}"
    fi
}

print_info() {
    if [ "$HAS_GUM" = true ]; then
        gum style --foreground 214 "ℹ $1"
    else
        echo -e "${YELLOW}ℹ $1${NC}"
    fi
}

confirm() {
    if [ "$HAS_GUM" = true ]; then
        gum confirm "$1"
    else
        read -p "$1 (y/n) " -n 1 -r
        echo
        [[ $REPLY =~ ^[Yy]$ ]]
    fi
}

input() {
    local prompt="$1"
    local default="$2"
    if [ "$HAS_GUM" = true ]; then
        gum input --placeholder "$prompt (default: $default)" --value "$default"
    else
        read -p "$prompt (default: $default): " value
        echo "${value:-$default}"
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

spin() {
    local title="$1"
    shift
    if [ "$HAS_GUM" = true ]; then
        gum spin --spinner dot --title "$title" -- "$@"
    else
        echo -n "$title"
        "$@" &> /dev/null
        echo " Done!"
    fi
}

# Check for prerequisites
check_prerequisites() {
    local missing_deps=()
    
    # Check for uv
    if ! command -v uv &> /dev/null; then
        missing_deps+=("uv")
    fi
    
    # Check for Python
    if ! command -v python3 &> /dev/null; then
        missing_deps+=("python3")
    fi
    
    if [ ${#missing_deps[@]} -gt 0 ]; then
        print_error "Missing required dependencies: ${missing_deps[*]}"
        
        for dep in "${missing_deps[@]}"; do
            case "$dep" in
                uv)
                    print_info "Install uv with: curl -LsSf https://astral.sh/uv/install.sh | sh"
                    if confirm "Would you like to install uv now?"; then
                        curl -LsSf https://astral.sh/uv/install.sh | sh
                        export PATH="$HOME/.local/bin:$PATH"
                    fi
                    ;;
                python3)
                    print_info "Python 3 is required. Please install it using your package manager."
                    ;;
            esac
        done
        
        # Re-check after potential installation
        if ! command -v uv &> /dev/null || ! command -v python3 &> /dev/null; then
            print_error "Please install all required dependencies and run setup again."
            exit 1
        fi
    fi
    
    print_success "All prerequisites installed"
}

# Install Gum if not present
install_gum() {
    if [ "$HAS_GUM" = false ]; then
        print_info "Gum is not installed. It provides a better setup experience."
        if confirm "Would you like to install Gum for a better experience?"; then
            case "$(uname -s)" in
                Darwin)
                    if command -v brew &> /dev/null; then
                        brew install gum
                    else
                        print_info "Install with: brew install gum"
                    fi
                    ;;
                Linux)
                    if command -v apt-get &> /dev/null; then
                        sudo mkdir -p /etc/apt/keyrings
                        curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg
                        echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list
                        sudo apt update && sudo apt install gum
                    elif command -v pacman &> /dev/null; then
                        sudo pacman -S gum
                    else
                        print_info "Visit https://github.com/charmbracelet/gum for installation instructions"
                    fi
                    ;;
                *)
                    print_info "Visit https://github.com/charmbracelet/gum for installation instructions"
                    ;;
            esac
            
            # Check if installation succeeded
            if command -v gum &> /dev/null; then
                HAS_GUM=true
                print_success "Gum installed successfully!"
            fi
        fi
    fi
}

# Configure hooks
configure_hooks() {
    print_header "Hook Configuration"
    
    # Check current hook status
    local hooks_enabled=false
    if [ -f "$CLAUDE_DIR/settings.local.json" ]; then
        if grep -q '"hooks"' "$CLAUDE_DIR/settings.local.json"; then
            hooks_enabled=true
        fi
    fi
    
    if [ "$hooks_enabled" = true ]; then
        print_info "Hooks are currently ENABLED"
    else
        print_info "Hooks are currently DISABLED"
    fi
    
    # Configuration options
    local enable_hooks
    if confirm "Enable agent instruction logging hooks?"; then
        enable_hooks=true
    else
        enable_hooks=false
    fi
    
    local debug_mode=false
    if [ "$enable_hooks" = true ]; then
        if confirm "Enable debug mode for hooks?"; then
            debug_mode=true
        fi
        
        # Agent filter configuration
        local agent_filter
        agent_filter=$(input "Agent prefix to track" "the-")
        
        # Log retention
        local log_retention
        log_retention=$(choose "Log retention policy:" "Keep all logs" "7 days" "30 days" "90 days")
    fi
    
    # Apply configuration
    if [ "$enable_hooks" = true ]; then
        print_info "Configuring hooks..."
        
        # Update hook scripts with configuration
        if [ "$agent_filter" != "the-" ]; then
            # Update the filter in the Python scripts
            sed -i.bak "s/startswith('the-')/startswith('$agent_filter')/" "$HOOKS_DIR/log_agent_start.py"
            sed -i.bak "s/startswith('the-')/startswith('$agent_filter')/" "$HOOKS_DIR/log_agent_complete.py"
            rm "$HOOKS_DIR"/*.bak
        fi
        
        # Create configuration file
        cat > "$CLAUDE_DIR/hooks/config.sh" << EOF
#!/usr/bin/env bash
# Hook configuration
export HOOK_DEBUG_MODE=$debug_mode
export HOOK_AGENT_FILTER="$agent_filter"
export HOOK_LOG_RETENTION="$log_retention"
EOF
        chmod +x "$CLAUDE_DIR/hooks/config.sh"
        
        print_success "Hooks configured successfully"
    else
        print_info "Hooks will remain disabled"
    fi
    
    # Set up environment
    if [ "$debug_mode" = true ]; then
        echo "export DEBUG_HOOKS=1" >> "$HOME/.bashrc" 2>/dev/null || true
        echo "export DEBUG_HOOKS=1" >> "$HOME/.zshrc" 2>/dev/null || true
        print_info "Debug mode enabled in shell configuration"
    fi
}

# Test hook execution
test_hooks() {
    print_header "Testing Hook Execution"
    
    if confirm "Would you like to test the hooks?"; then
        print_info "Creating test scenario..."
        
        # Create a test JSON input
        local test_input='{
            "tool_name": "Task",
            "tool_input": {
                "subagent_type": "the-test-agent",
                "prompt": "Test prompt. SessionId: test-session-001, AgentId: test123",
                "description": "Test description"
            }
        }'
        
        # Test PreToolUse hook
        print_info "Testing PreToolUse hook..."
        echo "$test_input" | uv run "$HOOKS_DIR/log_agent_start.py"
        
        if [ $? -eq 0 ]; then
            print_success "PreToolUse hook executed successfully"
        else
            print_error "PreToolUse hook failed"
        fi
        
        # Test PostToolUse hook
        local test_output='{
            "tool_name": "Task",
            "tool_input": {
                "subagent_type": "the-test-agent",
                "prompt": "Test prompt. SessionId: test-session-001, AgentId: test123",
                "description": "Test description"
            },
            "output": "Test completed successfully"
        }'
        
        print_info "Testing PostToolUse hook..."
        echo "$test_output" | uv run "$HOOKS_DIR/log_agent_complete.py"
        
        if [ $? -eq 0 ]; then
            print_success "PostToolUse hook executed successfully"
        else
            print_error "PostToolUse hook failed"
        fi
        
        # Check if logs were created
        if [ -f "$STARTUP_DIR/all-agent-instructions.jsonl" ]; then
            print_success "Log files created successfully"
            
            if [ "$HAS_GUM" = true ]; then
                if confirm "View test log entries?"; then
                    tail -n 2 "$STARTUP_DIR/all-agent-instructions.jsonl" | jq . | gum pager
                fi
            else
                print_info "Test log entries:"
                tail -n 2 "$STARTUP_DIR/all-agent-instructions.jsonl" | jq .
            fi
        else
            print_error "Log files not created"
        fi
    fi
}

# Main setup flow
main() {
    clear
    
    # Welcome message
    print_header "🚀 Claude Code Hooks Setup"
    
    if [ "$HAS_GUM" = false ]; then
        print_info "Running in basic mode. Install Gum for a better experience!"
    fi
    
    # Check and install prerequisites
    spin "Checking prerequisites..." check_prerequisites
    
    # Offer to install Gum
    install_gum
    
    # Create necessary directories
    spin "Creating directory structure..." mkdir -p "$HOOKS_DIR" "$STARTUP_DIR"
    
    # Configure hooks
    configure_hooks
    
    # Test hooks
    test_hooks
    
    # Success message
    print_header "✨ Setup Complete!"
    
    print_success "Claude Code hooks have been configured successfully!"
    echo
    print_info "Next steps:"
    echo "  1. Restart your Claude Code session for changes to take effect"
    echo "  2. Run './hook-manager.sh' to manage hooks"
    echo "  3. Check logs in: $STARTUP_DIR"
    echo
    print_info "Hook commands:"
    echo "  • View logs: ./hook-manager.sh logs"
    echo "  • Check status: ./hook-manager.sh status"
    echo "  • Clear logs: ./hook-manager.sh clear"
    echo
    
    if [ "$HAS_GUM" = true ]; then
        gum style \
            --border normal \
            --margin "1" \
            --padding "1" \
            --border-foreground 212 \
            "Thank you for using Claude Code Hooks! 🎉"
    else
        echo -e "${GREEN}Thank you for using Claude Code Hooks! 🎉${NC}"
    fi
}

# Run main function
main "$@"