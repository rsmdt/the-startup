#!/usr/bin/env bash
# Claude Code Hooks Manager
# Interactive tool for managing and monitoring hooks

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
SETTINGS_FILE="$CLAUDE_DIR/settings.local.json"

# Check if running with Gum
HAS_GUM=false
if command -v gum &> /dev/null; then
    HAS_GUM=true
fi

# Source common functions
source "$SCRIPT_DIR/lib/common.sh" 2>/dev/null || true

# Styling functions
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

# Check hook status
check_status() {
    print_header "Hook Status"
    
    # Check if hooks are configured
    if [ -f "$SETTINGS_FILE" ] && grep -q '"hooks"' "$SETTINGS_FILE"; then
        print_success "Hooks are ENABLED"
        
        # Check specific hooks
        if grep -q "PreToolUse" "$SETTINGS_FILE"; then
            print_success "PreToolUse hook: Active"
        else
            print_error "PreToolUse hook: Not configured"
        fi
        
        if grep -q "PostToolUse" "$SETTINGS_FILE"; then
            print_success "PostToolUse hook: Active"
        else
            print_error "PostToolUse hook: Not configured"
        fi
    else
        print_error "Hooks are DISABLED"
    fi
    
    # Check log files
    echo
    print_info "Log Files:"
    if [ -f "$STARTUP_DIR/all-agent-instructions.jsonl" ]; then
        local log_size=$(du -h "$STARTUP_DIR/all-agent-instructions.jsonl" | cut -f1)
        local log_lines=$(wc -l < "$STARTUP_DIR/all-agent-instructions.jsonl")
        print_success "Global log: $log_size ($log_lines entries)"
    else
        print_info "Global log: Not created yet"
    fi
    
    # Check session logs
    local session_count=$(find "$STARTUP_DIR" -name "agent-instructions.jsonl" 2>/dev/null | wc -l)
    if [ "$session_count" -gt 0 ]; then
        print_success "Session logs: $session_count sessions"
    else
        print_info "Session logs: No sessions yet"
    fi
    
    # Check debug mode
    echo
    if [ "$DEBUG_HOOKS" = "1" ]; then
        print_info "Debug mode: ENABLED"
    else
        print_info "Debug mode: DISABLED"
    fi
}

# View logs
view_logs() {
    print_header "View Logs"
    
    if [ ! -f "$STARTUP_DIR/all-agent-instructions.jsonl" ]; then
        print_error "No logs found yet"
        return
    fi
    
    if [ "$HAS_GUM" = true ]; then
        local view_option=$(gum choose \
            "View all logs" \
            "View recent logs (last 20)" \
            "View by session" \
            "View by agent type" \
            "Search logs" \
            "Back")
    else
        echo "1) View all logs"
        echo "2) View recent logs (last 20)"
        echo "3) View by session"
        echo "4) View by agent type"
        echo "5) Search logs"
        echo "6) Back"
        read -p "Select option: " choice
        case $choice in
            1) view_option="View all logs" ;;
            2) view_option="View recent logs (last 20)" ;;
            3) view_option="View by session" ;;
            4) view_option="View by agent type" ;;
            5) view_option="Search logs" ;;
            6) view_option="Back" ;;
            *) view_option="Back" ;;
        esac
    fi
    
    case "$view_option" in
        "View all logs")
            if [ "$HAS_GUM" = true ]; then
                cat "$STARTUP_DIR/all-agent-instructions.jsonl" | jq . | gum pager
            else
                cat "$STARTUP_DIR/all-agent-instructions.jsonl" | jq . | less
            fi
            ;;
        "View recent logs (last 20)")
            if [ "$HAS_GUM" = true ]; then
                tail -n 20 "$STARTUP_DIR/all-agent-instructions.jsonl" | jq . | gum pager
            else
                tail -n 20 "$STARTUP_DIR/all-agent-instructions.jsonl" | jq .
            fi
            ;;
        "View by session")
            local sessions=$(find "$STARTUP_DIR" -type d -name "dev-*" 2>/dev/null | sort -r)
            if [ -z "$sessions" ]; then
                print_error "No sessions found"
                return
            fi
            
            if [ "$HAS_GUM" = true ]; then
                local selected_session=$(echo "$sessions" | xargs -n1 basename | gum choose)
                if [ -n "$selected_session" ]; then
                    local session_log="$STARTUP_DIR/$selected_session/agent-instructions.jsonl"
                    if [ -f "$session_log" ]; then
                        cat "$session_log" | jq . | gum pager
                    else
                        print_error "No logs for this session"
                    fi
                fi
            else
                echo "Available sessions:"
                echo "$sessions" | xargs -n1 basename | nl
                read -p "Select session number: " session_num
                local selected_session=$(echo "$sessions" | sed -n "${session_num}p" | xargs basename)
                if [ -n "$selected_session" ]; then
                    local session_log="$STARTUP_DIR/$selected_session/agent-instructions.jsonl"
                    if [ -f "$session_log" ]; then
                        cat "$session_log" | jq . | less
                    else
                        print_error "No logs for this session"
                    fi
                fi
            fi
            ;;
        "View by agent type")
            local agent_types=$(cat "$STARTUP_DIR/all-agent-instructions.jsonl" | jq -r '.agent_type' | sort -u)
            if [ "$HAS_GUM" = true ]; then
                local selected_agent=$(echo "$agent_types" | gum choose)
                if [ -n "$selected_agent" ]; then
                    cat "$STARTUP_DIR/all-agent-instructions.jsonl" | jq "select(.agent_type == \"$selected_agent\")" | gum pager
                fi
            else
                echo "Available agent types:"
                echo "$agent_types" | nl
                read -p "Select agent type number: " agent_num
                local selected_agent=$(echo "$agent_types" | sed -n "${agent_num}p")
                if [ -n "$selected_agent" ]; then
                    cat "$STARTUP_DIR/all-agent-instructions.jsonl" | jq "select(.agent_type == \"$selected_agent\")" | less
                fi
            fi
            ;;
        "Search logs")
            if [ "$HAS_GUM" = true ]; then
                local search_term=$(gum input --placeholder "Enter search term")
                if [ -n "$search_term" ]; then
                    grep -i "$search_term" "$STARTUP_DIR/all-agent-instructions.jsonl" | jq . | gum pager
                fi
            else
                read -p "Enter search term: " search_term
                if [ -n "$search_term" ]; then
                    grep -i "$search_term" "$STARTUP_DIR/all-agent-instructions.jsonl" | jq . | less
                fi
            fi
            ;;
    esac
}

# Toggle hooks
toggle_hooks() {
    print_header "Toggle Hooks"
    
    if [ -f "$SETTINGS_FILE" ] && grep -q '"hooks"' "$SETTINGS_FILE"; then
        print_info "Hooks are currently ENABLED"
        if [ "$HAS_GUM" = true ]; then
            if gum confirm "Disable hooks?"; then
                # Remove hooks section from settings
                python3 -c "
import json
with open('$SETTINGS_FILE', 'r') as f:
    settings = json.load(f)
if 'hooks' in settings:
    del settings['hooks']
with open('$SETTINGS_FILE', 'w') as f:
    json.dump(settings, f, indent=2)
"
                print_success "Hooks disabled"
            fi
        else
            read -p "Disable hooks? (y/n) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                python3 -c "
import json
with open('$SETTINGS_FILE', 'r') as f:
    settings = json.load(f)
if 'hooks' in settings:
    del settings['hooks']
with open('$SETTINGS_FILE', 'w') as f:
    json.dump(settings, f, indent=2)
"
                print_success "Hooks disabled"
            fi
        fi
    else
        print_info "Hooks are currently DISABLED"
        if [ "$HAS_GUM" = true ]; then
            if gum confirm "Enable hooks?"; then
                # Add hooks section to settings
                python3 -c "
import json
with open('$SETTINGS_FILE', 'r') as f:
    settings = json.load(f)
settings['hooks'] = {
    'PreToolUse': [{
        'matcher': 'Task',
        'hooks': [{
            'type': 'command',
            'command': 'uv run \$CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_start.py'
        }]
    }],
    'PostToolUse': [{
        'matcher': 'Task',
        'hooks': [{
            'type': 'command',
            'command': 'uv run \$CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_complete.py'
        }]
    }]
}
with open('$SETTINGS_FILE', 'w') as f:
    json.dump(settings, f, indent=2)
"
                print_success "Hooks enabled"
            fi
        else
            read -p "Enable hooks? (y/n) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                python3 -c "
import json
with open('$SETTINGS_FILE', 'r') as f:
    settings = json.load(f)
settings['hooks'] = {
    'PreToolUse': [{
        'matcher': 'Task',
        'hooks': [{
            'type': 'command',
            'command': 'uv run \$CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_start.py'
        }]
    }],
    'PostToolUse': [{
        'matcher': 'Task',
        'hooks': [{
            'type': 'command',
            'command': 'uv run \$CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_complete.py'
        }]
    }]
}
with open('$SETTINGS_FILE', 'w') as f:
    json.dump(settings, f, indent=2)
"
                print_success "Hooks enabled"
            fi
        fi
    fi
}

# Clear logs
clear_logs() {
    print_header "Clear Logs"
    
    if [ "$HAS_GUM" = true ]; then
        local clear_option=$(gum choose \
            "Clear all logs" \
            "Clear global log only" \
            "Clear session logs only" \
            "Clear old sessions (>7 days)" \
            "Back")
    else
        echo "1) Clear all logs"
        echo "2) Clear global log only"
        echo "3) Clear session logs only"
        echo "4) Clear old sessions (>7 days)"
        echo "5) Back"
        read -p "Select option: " choice
        case $choice in
            1) clear_option="Clear all logs" ;;
            2) clear_option="Clear global log only" ;;
            3) clear_option="Clear session logs only" ;;
            4) clear_option="Clear old sessions (>7 days)" ;;
            5) clear_option="Back" ;;
            *) clear_option="Back" ;;
        esac
    fi
    
    case "$clear_option" in
        "Clear all logs")
            if [ "$HAS_GUM" = true ]; then
                if gum confirm "Are you sure you want to clear ALL logs?"; then
                    rm -f "$STARTUP_DIR/all-agent-instructions.jsonl"
                    rm -rf "$STARTUP_DIR"/dev-*/
                    print_success "All logs cleared"
                fi
            else
                read -p "Are you sure you want to clear ALL logs? (y/n) " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    rm -f "$STARTUP_DIR/all-agent-instructions.jsonl"
                    rm -rf "$STARTUP_DIR"/dev-*/
                    print_success "All logs cleared"
                fi
            fi
            ;;
        "Clear global log only")
            if [ -f "$STARTUP_DIR/all-agent-instructions.jsonl" ]; then
                rm -f "$STARTUP_DIR/all-agent-instructions.jsonl"
                print_success "Global log cleared"
            else
                print_info "No global log to clear"
            fi
            ;;
        "Clear session logs only")
            if ls "$STARTUP_DIR"/dev-*/ &>/dev/null; then
                rm -rf "$STARTUP_DIR"/dev-*/
                print_success "Session logs cleared"
            else
                print_info "No session logs to clear"
            fi
            ;;
        "Clear old sessions (>7 days)")
            local old_sessions=$(find "$STARTUP_DIR" -type d -name "dev-*" -mtime +7 2>/dev/null)
            if [ -n "$old_sessions" ]; then
                echo "$old_sessions" | xargs rm -rf
                local count=$(echo "$old_sessions" | wc -l)
                print_success "Cleared $count old sessions"
            else
                print_info "No old sessions to clear"
            fi
            ;;
    esac
}

# Test hooks
test_hooks() {
    print_header "Test Hooks"
    
    print_info "Creating test scenario..."
    
    # Create test input
    local test_input='{
        "tool_name": "Task",
        "tool_input": {
            "subagent_type": "the-test-agent",
            "prompt": "Test prompt for hook testing. SessionId: test-'$(date +%Y%m%d-%H%M%S)'-test, AgentId: test'$(openssl rand -hex 3)'",
            "description": "Hook test"
        }
    }'
    
    print_info "Testing PreToolUse hook..."
    echo "$test_input" | uv run "$HOOKS_DIR/log_agent_start.py"
    
    if [ $? -eq 0 ]; then
        print_success "PreToolUse hook executed successfully"
    else
        print_error "PreToolUse hook failed"
    fi
    
    # Add output for PostToolUse test
    local test_output=$(echo "$test_input" | jq '. + {"output": "Test completed successfully"}')
    
    print_info "Testing PostToolUse hook..."
    echo "$test_output" | uv run "$HOOKS_DIR/log_agent_complete.py"
    
    if [ $? -eq 0 ]; then
        print_success "PostToolUse hook executed successfully"
    else
        print_error "PostToolUse hook failed"
    fi
    
    # Show recent log entries
    if [ -f "$STARTUP_DIR/all-agent-instructions.jsonl" ]; then
        print_info "Recent test entries:"
        tail -n 2 "$STARTUP_DIR/all-agent-instructions.jsonl" | jq .
    fi
}

# Main menu
main_menu() {
    while true; do
        clear
        print_header "🔧 Claude Code Hook Manager"
        
        if [ "$HAS_GUM" = true ]; then
            local action=$(gum choose \
                "View Status" \
                "View Logs" \
                "Toggle Hooks" \
                "Clear Logs" \
                "Test Hooks" \
                "Exit")
        else
            echo "1) View Status"
            echo "2) View Logs"
            echo "3) Toggle Hooks"
            echo "4) Clear Logs"
            echo "5) Test Hooks"
            echo "6) Exit"
            echo
            read -p "Select option: " choice
            case $choice in
                1) action="View Status" ;;
                2) action="View Logs" ;;
                3) action="Toggle Hooks" ;;
                4) action="Clear Logs" ;;
                5) action="Test Hooks" ;;
                6) action="Exit" ;;
                *) continue ;;
            esac
        fi
        
        case "$action" in
            "View Status")
                check_status
                ;;
            "View Logs")
                view_logs
                ;;
            "Toggle Hooks")
                toggle_hooks
                ;;
            "Clear Logs")
                clear_logs
                ;;
            "Test Hooks")
                test_hooks
                ;;
            "Exit")
                print_info "Goodbye!"
                exit 0
                ;;
        esac
        
        echo
        if [ "$HAS_GUM" = false ]; then
            read -p "Press Enter to continue..."
        else
            gum input --placeholder "Press Enter to continue..." > /dev/null
        fi
    done
}

# Handle command-line arguments
case "${1:-}" in
    status)
        check_status
        ;;
    logs)
        view_logs
        ;;
    toggle)
        toggle_hooks
        ;;
    clear)
        clear_logs
        ;;
    test)
        test_hooks
        ;;
    *)
        main_menu
        ;;
esac