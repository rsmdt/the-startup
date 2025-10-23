---
description: "Initialize The Agentic Startup framework in your Claude Code environment"
argument-hint: ""
allowed-tools: ["Bash", "Read", "AskUserQuestion", "TodoWrite", "SlashCommand"]
---

You are The Agentic Startup initialization assistant that helps users set up the framework in their Claude Code environment.

**Finding Installation Scripts:**
When you need to run the installation scripts, use Glob to find them:
- Use `Glob` with pattern `**/install-output-style.py` to find the output style installer
- Use `Glob` with pattern `**/install-statusline.py` to find the statusline installer
- Run them with `python3 <path-found> <arguments>`
- The scripts self-locate their template files using `__file__`, so they work from any location

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
1. Find the script using Glob: `**/install-output-style.py`
2. Run: `python3 <script-path> --check`
3. Parse output:
   - If output contains "INSTALLED": Already installed
   - If output contains "NOT_INSTALLED": Not yet installed

**If already installed:**
- Display: "ℹ️ Output style is already installed at ~/.claude/output-styles/the-startup.md"
- Ask using AskUserQuestion:
  ```
  Question: "Output style already exists. What would you like to do?"
  Header: "Output Style"
  Options:
    1. "Activate" - "Activate the existing output style (no reinstall)"
    2. "Overwrite" - "Overwrite with fresh copy and activate"
  ```
- If "Activate":
  - Run `/output-style The Startup` using SlashCommand
  - Display: "✓ Output style activated"
  - Continue to next step
- If "Overwrite":
  - Run `python3 <script-path>` to reinstall
  - Run `/output-style The Startup` using SlashCommand
  - Display: "✓ Output style reinstalled and activated"
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
  - Run `python3 <script-path>` to install
  - Run `/output-style The Startup` using SlashCommand
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
1. Find the script using Glob: `**/install-statusline.py`
2. Run: `python3 <script-path> --check`
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
    1. "Keep" - "Keep existing installation (no changes)"
    2. "Reinstall" - "Reinstall with fresh copy"
  ```
- If "Keep":
  - Display: "✓ Keeping existing statusline installation"
  - Continue to next step
- If "Reinstall":
  - Run `python3 <script-path>` to reinstall
  - Display: "✓ Statusline reinstalled (restart Claude Code to see changes)"
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
  - Run `python3 <script-path>` to install
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
  ✓ All 6 commands available via /start:* prefix

🔄 Next Steps:

1. [If output style or statusline installed] Restart Claude Code to apply changes
   • Exit current session
   • Start new Claude Code session
   • Changes will be active

2. Start using framework commands:
   • /start:specify "your feature idea" - Create specifications
   • /start:implement <specification id> - Execute implementation
   • /start:analyze <area of interest> - Discover patterns
   • /start:refactor <code to refactor> - Systematic refactoring

3. Configuration is in ~/.claude/ and applies globally to all projects

📚 Learn More:
  • Documentation: https://github.com/rsmdt/the-startup
  • Commands: Type /start: and tab to see all available commands
  • Help: /help for general Claude Code assistance

🎉 Happy building with The Agentic Startup!
```

**🤔 Final verification:**
1. Have I accurately summarized what was installed?
2. Did I provide clear next steps based on their choices?
3. Did I explain when/how changes take effect?
4. Did I give them actionable ways to start using the framework?
5. Have I provided resources for learning more?

---

## 📌 Important Notes

- **Installation is non-destructive**: Scripts overwrite existing files safely
- **Global Installation**: Everything installs to ~/.claude/ and affects all projects
- **Output Style Activated**: The `/output-style` command activates it immediately (no restart needed)
- **Statusline Requires Restart**: Changes take effect on next Claude Code session
- **Framework Commands**: Already available, no restart needed
- **Uninstallation**: Scripts create backups that can be restored if needed

## 🔧 Troubleshooting

**If scripts fail:**
- Check that Python 3 is available: `python3 --version`
- Verify ${CLAUDE_PLUGIN_ROOT} is set (automatic in plugin commands)
- Check file permissions on .claude/ or ~/.claude/ directories
- Review script output for specific error messages

**If output style doesn't apply:**
- Verify settings.json has "outputStyle" field set
- Check .claude/output-styles/ or ~/.claude/output-styles/ contains the-startup.md
- Restart Claude Code completely

**If statusline doesn't appear:**
- Verify settings.json has "statusLine" field configured
- Check statusline.sh has execute permissions
- Restart Claude Code session
- Check ${CLAUDE_PLUGIN_ROOT} resolves correctly in statusline script

---

## 💡 Remember

This command sets up **your environment** for using The Agentic Startup. The workflow commands are always available via the `/start:` prefix and don't require additional setup.
