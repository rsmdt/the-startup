#!/usr/bin/env bash
# The Agentic Startup - Welcome Hook
# Displays banner on first plugin session

set -euo pipefail

# Flag file to track first run
FLAG_FILE="${HOME}/.the-startup/.plugin-initialized"

# Check if this is the first run
if [ -f "$FLAG_FILE" ]; then
  # Already initialized - silent exit
  echo '{}'
  exit 0
fi

# Create flag file
mkdir -p "$(dirname "$FLAG_FILE")"
touch "$FLAG_FILE"

# Output welcome banner in additionalContext
cat <<'EOF'
{
  "additionalContext": "🚀 **Welcome to The Agentic Startup!**\n\n**Available Capabilities:**\n\n**📦 50 Specialized Agents** across 9 roles:\n- The Analyst (requirements, feature prioritization, project coordination)\n- The Architect (system architecture, technology research, documentation, quality review)\n- The Designer (design foundation, accessibility, interaction architecture, user research)\n- The ML Engineer (feature operations, prompt optimization, ML operations, context management)\n- The Mobile Engineer (mobile data persistence, operations, development)\n- The Platform Engineer (containerization, pipeline engineering, production monitoring, infrastructure, performance tuning, deployment automation, data architecture)\n- The QA Engineer (test execution, exploratory testing, performance testing)\n- The Security Engineer (security assessment, implementation, incident response)\n- The Software Engineer (performance optimization, component development, service resilience, domain modeling, browser compatibility, API development)\n\n**🔧 6 Slash Commands:**\n- `/s:specify` - Create comprehensive specifications\n- `/s:analyze` - Discover business rules and patterns\n- `/s:implement` - Execute implementation plans\n- `/s:refactor` - Refactor code maintaining behavior\n- `/s:init` - Initialize quality gate templates\n\n**📊 Git Status Integration:**\nYour statusline now shows git branch information\n\n**📚 Documentation:**\n- Rules at `rules/`\n- Templates at `templates/`\n- Agents at `agents/`\n\n**Ready to build something incredible!** 🎯"
}
EOF
