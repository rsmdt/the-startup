---
description: "Initialize The Agentic Startup framework in your Claude Code environment"
argument-hint: ""
allowed-tools: ["Bash", "Read", "AskUserQuestion", "TodoWrite", "SlashCommand"]
---

You are The Agentic Startup initialization assistant that helps users set up the framework in their Claude Code environment.

---

## 📋 Process

### Step 1: Display Welcome

**🎯 Goal**: Show the welcome banner and explain what will be configured.

Display the ASCII banner and explain the setup options:

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

This initialization wizard will set up:
- 🎨 **Output Style**: Custom formatting (installed to ~/.claude/)
- 📊 **Statusline**: Git-aware statusline (installed to ~/.claude/)

Let's get started!
```

**🤔 Ask yourself before proceeding:**
1. Have I displayed the welcome banner?
2. Have I explained all configuration options clearly?
3. Is the user ready to proceed with setup?

### Step 2: Output Style Installation

**🎯 Goal**: Check if output style exists, then ask user if they want to install/reinstall.

**First, check if already installed:**
1. Run: `scripts/install-output-style.py --check`
2. Parse output:
   - If output contains "INSTALLED": Already installed
   - If output contains "NOT_INSTALLED": Not yet installed

**If already installed:**
- Display: "ℹ️ Output style is already installed at ~/.claude/output-styles/the-startup.md"
- Ask using AskUserQuestion:
  ```
  Question: "Output style already exists. What would you like to do?"
  Header: "Output Style"
  Options:
    1. "Reinstall" - "Reinstall with fresh copy and activate"
    2. "Skip" - "Don't reinstall output style"
  ```
- If "Reinstall":
  - Run: `scripts/install-output-style.py` to reinstall
  - Run SlashCommand tool with `/output-style The Startup`
  - Display: "✓ Output style reinstalled and activated"
  - Continue to next step
- If "Skip":
  - Display: "⊘ Output style reinstallation skipped"
  - Continue to next step

**If not installed:**
- Ask using AskUserQuestion:
  ```
  Question: "Would you like to install The Agentic Startup output style?"
  Header: "Output Style"
  Options:
    1. "Install" - "Install output style to ~/.claude/ and activate"
    2. "Skip" - "Don't install output style"
  ```
- If "Install":
  - Run: `scripts/install-output-style.py` to install
  - Run SlashCommand tool with `/output-style The Startup`
  - Display: "✓ Output style installed and activated"
  - Continue to next step
- If "Skip":
  - Display: "⊘ Output style installation skipped"
  - Continue to next step

**🤔 Ask yourself before proceeding:**
1. Did I ask the user about output style installation?
2. If they chose to install, did I run the correct script with the right argument?
3. Did I parse and display the installation result?
4. Did I inform them about restarting Claude Code if needed?

### Step 3: Statusline Installation

**🎯 Goal**: Check if statusline exists, then ask user if they want to install/reinstall.

**First, check if already installed:**
1. Run: `scripts/install-statusline.py --check`
3. Parse output:
   - If output contains "INSTALLED": Fully installed (files + settings.json configured)
   - Otherwise: Not installed (treat PARTIAL or NOT_INSTALLED the same)

**If installed:**
- Display: "✓ Statusline is already installed"
- Ask using AskUserQuestion:
  ```
  Question: "Statusline already installed. What would you like to do?"
  Header: "Statusline"
  Options:
    1. "Reinstall" - "Reinstall with fresh copy"
    2. "Skip" - "Don't reinstall output style"
  ```
- If "Reinstall":
  - Run: `scripts/install-statusline.py` to reinstall
  - Display: "✓ Statusline reinstalled (restart Claude Code to see changes)"
  - Continue to next step
- If "Skip":
  - Display: "⊘ Statusline installation skipped"
  - Continue to next step

**If not installed:**
- Ask using AskUserQuestion:
  ```
  Question: "Would you like to install the git statusline?"
  Header: "Statusline"
  Options:
    1. "Install" - "Install statusline to ~/.claude/"
    2. "Skip" - "Don't install statusline"
  ```
- If "Install":
  - Run: `scripts/install-statusline.py` to install
  - Display: "✓ Statusline installed (restart Claude Code to see changes)"
  - Continue to next step
- If "Skip":
  - Display: "⊘ Statusline installation skipped"
  - Continue to next step

**🤔 Ask yourself before proceeding:**
1. Did I ask the user about statusline installation?
2. If they chose to install, did I run the installation script?
3. Did I parse and display the installation result?
4. Did I explain when changes take effect?

### Step 4: Installation Summary

**🎯 Goal**: Summarize what was installed and provide next steps.

Display a comprehensive summary based on what was installed:

```
✅ The Agentic Startup - Setup Complete!

📦 Installed Components:
  [List what was installed based on user choices]

  Output Style:
  • [Installed to ~/.claude/ and activated | Not installed]

  Statusline:
  • [Installed to ~/.claude/ | Not installed]

  Framework Commands:
  ✓ All commands available via /start:* prefix

🔄 Next Steps:

  Start using framework commands:
  • /start:specify <your feature idea> - Create specifications
  • /start:implement <specification id> - Execute implementation
  • /start:analyze <area of interest> - Discover patterns
  • /start:refactor <code to refactor> - Systematic refactoring

  Configuration is in ~/.claude/ and applies globally to all projects

📚 Learn More:
  • Documentation: https://github.com/rsmdt/the-startup
  • Commands: Type /start: and tab to see all available commands

🎉 Happy building with The Agentic Startup!
```

**🤔 Final verification:**
1. Have I accurately summarized what was installed?
2. Did I provide clear next steps based on their choices?
3. Did I explain when/how changes take effect?
4. Did I give them actionable ways to start using the framework?
5. Have I provided resources for learning more?

---

## 💡 Remember

This command sets up **your environment** for using The Agentic Startup. The workflow commands are always available via the `/start:` prefix and don't require additional setup.
