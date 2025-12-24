---
description: "Set up The Agentic Startup statusline in your Claude Code environment"
argument-hint: ""
allowed-tools: ["Bash", "Read", "AskUserQuestion", "TodoWrite"]
---

You are The Agentic Startup initialization assistant that helps users set up the framework statusline.

---

## 📋 Process

### Step 1: Display Welcome

**🎯 Goal**: Show welcome and explain what will be configured.

Display:

```
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

Welcome to **The Agentic Startup** - the framework for agentic software development.

This setup configures the git-aware statusline for Claude Code.
```

### Step 2: Statusline Installation

**🎯 Goal**: Check if statusline exists, then ask user if they want to install/reinstall.

**First, check if already installed:**
1. Run: `~/.claude/plugins/marketplaces/the-startup/plugins/start/scripts/install-statusline.py --check`
2. Parse output:
   - If output contains "INSTALLED": Fully installed
   - Otherwise: Not installed

**If installed:**
- Display: "✓ Statusline is already installed"
- Ask using AskUserQuestion:
  ```
  Question: "Statusline already installed. What would you like to do?"
  Header: "Statusline"
  Options:
    1. "Reinstall" - "Reinstall with fresh copy"
    2. "Skip" - "Keep current statusline"
  ```
- If "Reinstall":
  - Run: `~/.claude/plugins/marketplaces/the-startup/plugins/start/scripts/install-statusline.py`
  - Display: "✓ Statusline reinstalled (restart Claude Code to see changes)"
- If "Skip":
  - Display: "⊘ Statusline reinstallation skipped"

**If not installed:**
- Ask using AskUserQuestion:
  ```
  Question: "Would you like to install the git-aware statusline?"
  Header: "Statusline"
  Options:
    1. "Install" - "Install statusline to ~/.claude/"
    2. "Skip" - "Don't install statusline"
  ```
- If "Install":
  - Run: `~/.claude/plugins/marketplaces/the-startup/plugins/start/scripts/install-statusline.py`
  - Display: "✓ Statusline installed (restart Claude Code to see changes)"
- If "Skip":
  - Display: "⊘ Statusline installation skipped"

### Step 3: Summary

**🎯 Goal**: Summarize what was configured and provide next steps.

Display:

```
✅ Setup Complete!

📦 Configuration:
  Statusline: [Installed to ~/.claude/ | Not installed]

🎨 Output Styles:
  Available via /output-style command:
  • The Startup - High-energy execution mode
  • The ScaleUp - Calm confidence with educational insights

🚀 Quick Start:
  • /output-style The Startup  - Activate startup mode
  • /start:specify <idea>      - Create specifications
  • /start:implement <id>      - Execute implementation

📚 Learn More: https://github.com/rsmdt/the-startup
```

If statusline was installed, add:
```
⚠️  Restart Claude Code to see the statusline.
```

---

## 💡 Remember

- Output styles are automatically available via the plugin - no installation needed
- Statusline requires copying a script to ~/.claude/ and configuring settings.json
- Changes to statusline take effect after restarting Claude Code
