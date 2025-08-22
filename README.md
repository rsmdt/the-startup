<p align="center"><img src="https://github.com/rsmdt/the-startup/blob/main/assets/logo.png" width="400" alt="The Startup"></p>

<p align="center">Ship faster. Ship better. Ship with <b>The Agentic Startup</b>.</p>

## What is The Startup?

The Startup brings you instant access to expert developers, architects, and engineers - all working together to turn your ideas into shipped code.

The Startup is an orchestration system for Claude Code that gives you a virtual engineering team. Instead of one AI trying to do everything, you get specialized experts who collaborate like a real startup team - pragmatic, fast, and focused on shipping.

Think of it as having a CTO, architects, developers, and DevOps engineers on-demand, each bringing their expertise to your project.

## Quick Start

Install and start using The Startup:

```bash
# Install (interactive)
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh

# Plan a feature
/s:specify "Add user authentication"

# Build it 
/s:implement 001-user-auth
```

**More installation options**: See the [Installation & Uninstall](#installation--uninstall) section below.

## The Startup Way

When you use The Startup, Claude Code becomes your **technical co-founder** with a specific personality:

- **Ship over perfection** - MVPs today, not perfect solutions next quarter
- **Pragmatic decisions** - Make the call with available info, iterate later  
- **Specialist delegation** - Pull in the right expert for each task
- **Parallel execution** - Multiple experts work simultaneously when possible

You're not talking to a generic AI - you're working with a startup CTO who knows when to be scrappy and when to bring in the specialists.

## Your Expert Team

### 🚀 Leadership & Orchestration
- [**the-chief**](assets/claude/agents/the-chief.md) ¯\\_(ツ)_/¯ - Routes requests based on complexity assessment
- [**the-project-manager**](assets/claude/agents/the-project-manager.md) (⌐■_■) - Breaks down work, removes blockers, coordinates tasks
- [**the-product-manager**](assets/claude/agents/the-product-manager.md) (＾-＾)ノ - Prioritizes features, creates user stories, defines metrics

### 🏗️ Architecture & Design
- [**the-software-architect**](assets/claude/agents/the-software-architect.md) (⌐■_■) - System design, service boundaries, technical trade-offs
- [**the-staff-engineer**](assets/claude/agents/the-staff-engineer.md) (⚡◡⚡) - Sets technical standards, defines patterns, mentors teams
- [**the-business-analyst**](assets/claude/agents/the-business-analyst.md) (◔_◔) - Clarifies vague requirements through targeted questioning

### 💻 Engineering Team
- [**the-lead-engineer**](assets/claude/agents/the-lead-engineer.md) (▰˘◡˘▰) - Reviews code, provides mentorship, ensures quality
- [**the-frontend-engineer**](assets/claude/agents/the-frontend-engineer.md) (◕‿◕) - React/Vue/Angular, components, performance optimization
- [**the-backend-engineer**](assets/claude/agents/the-backend-engineer.md) (⚙◡⚙) - APIs, services, business logic, database design
- [**the-mobile-engineer**](assets/claude/agents/the-mobile-engineer.md) (📱◡📱) - iOS/Android, React Native, app store deployment
- [**the-ml-engineer**](assets/claude/agents/the-ml-engineer.md) (🤖◡🤖) - Model integration, MLOps, inference optimization
- [**the-developer**](assets/claude/agents/the-developer.md) 🚫 - *[DEPRECATED - Use specialized engineers above]*

### 🚦 Infrastructure & Operations
- [**the-devops-engineer**](assets/claude/agents/the-devops-engineer.md) (◉_◉) - CI/CD pipelines, containerization, infrastructure as code
- [**the-site-reliability-engineer**](assets/claude/agents/the-site-reliability-engineer.md) (╯°□°)╯ - Incident response, debugging, root cause analysis
- [**the-data-engineer**](assets/claude/agents/the-data-engineer.md) (⊙_⊙) - Database optimization, ETL pipelines, query performance
- [**the-performance-engineer**](assets/claude/agents/the-performance-engineer.md) (⚡◡⚡) - Core Web Vitals, bundle optimization, load times

### 🎨 Design & Documentation
- [**the-ux-designer**](assets/claude/agents/the-ux-designer.md) (◍•ᴗ•◍) - User interfaces, accessibility, interaction patterns
- [**the-principal-designer**](assets/claude/agents/the-principal-designer.md) (◉◡◉) - Design systems, design review, strategic vision
- [**the-technical-writer**](assets/claude/agents/the-technical-writer.md) (◕‿◕) - API docs, user guides, system documentation

### 🛡️ Quality & Security
- [**the-qa-lead**](assets/claude/agents/the-qa-lead.md) (✓◡✓) - Test strategy, risk prioritization, release decisions
- [**the-qa-engineer**](assets/claude/agents/the-qa-engineer.md) (¬_¬) - Test implementation, automation, bug hunting
- [**the-security-engineer**](assets/claude/agents/the-security-engineer.md) (ಠ_ಠ) - Vulnerability assessment, secure practices, incident response
- [**the-compliance-officer**](assets/claude/agents/the-compliance-officer.md) (⚖◡⚖) - GDPR/HIPAA compliance, data privacy, audit trails

### 🤖 AI & Specialized
- [**the-prompt-engineer**](assets/claude/agents/the-prompt-engineer.md) (◐‿◑) - Claude prompt optimization, agent instructions
- [**the-context-engineer**](assets/claude/agents/the-context-engineer.md) (ʘ_ʘ) - AI memory systems, context windows, inter-agent communication

## Commands

The Startup provides commands:

### `/s:specify` - Plan Before You Build

Creates comprehensive specifications from your ideas:

```bash
# Start fresh with a new feature idea
/s:specify Build a real-time notification system

# Resume working on a specification
/s:specify 001
```

This command will:
1. Gather requirements through targeted questions
2. Create business and technical specifications  
3. Design the system architecture
4. Generate an implementation plan

### `/s:implement` - Execute the Plan

Takes a specification and builds it with the right experts:

```bash
# Implement a completed specification
/s:implement 001-notifications
```

This command will:
1. Load the implementation plan
2. Assign tasks to appropriate specialists
3. Execute phase-by-phase with validation
4. Track progress through completion

### `/s:refactor` - Improve Code Quality

Analyzes and refactors existing code for better maintainability:

```bash
# Refactor specific code or modules
/s:refactor improve the authentication module for better testability

# Refactor for specific goals
/s:refactor reduce complexity in the payment processing logic
```

This command will:
1. Analyze existing code structure and patterns
2. Identify refactoring opportunities
3. Preserve behavior while improving quality
4. Apply industry best practices

## Real-World Examples

### Building Authentication
```
/s:specify Add user authentication with JWT
```
The Startup will:
- Use **the-business-analyst** to clarify requirements (OAuth? 2FA? Password reset?)
- Bring in **the-software-architect** to design the system
- Get **the-security-engineer** to review for vulnerabilities
- Create a complete implementation plan

### Debugging Production Issues
```
The API is returning 500 errors on user login
```
The Startup will:
- Immediately call **the-site-reliability-engineer** to investigate
- Once root cause is found, bring in **the-backend-engineer** to fix
- Have **the-qa-engineer** verify the fix
- Get **the-lead-engineer** to review the changes

### Creating a Dashboard
```
/s:specify Admin dashboard for monitoring system metrics
```
The Startup will:
- Use **the-ux-designer** to create the interface design
- Bring in **the-data-engineer** for efficient data queries
- Get **the-software-architect** to design the real-time data flow
- Deploy **the-frontend-engineer** and **the-backend-engineer** for implementation

## How It Works

1. **You make a request** - Either directly or through commands
2. **The Startup assesses** - Determines complexity and required expertise
3. **Specialists are called** - The right experts for your specific need
4. **Parallel execution** - Multiple experts work simultaneously when possible
5. **Results are synthesized** - Expert input becomes actionable next steps
6. **You ship faster** - With the confidence of a full team behind you

## Quick Start Examples

```bash
# Fix a bug
"Error: Cannot read property 'user' of undefined in auth.js"

# Build a feature  
/s:specify "Add CSV export functionality to reports"

# Optimize performance
"The dashboard takes 10 seconds to load"

# Review code
"Review my authentication implementation for security issues"

# Get unstuck
"I don't know how to structure this microservices architecture"
```

## Installation

The Startup provides easy installation via script and clean uninstall capabilities.

### Installation Options

Install The Startup agents, commands, and configuration using the install script:

```bash
# Interactive installation (shows all options)
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh

# Quick global installation (recommended paths, no prompts)
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh -s -- -y

# Local installation (project-specific paths, with file selection)
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh -s -- -l

# Quick local installation (project-specific, no prompts)
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh -s -- -ly
```

## Building from Source

If you want to contribute or customize:

```bash
# Clone and build
git clone https://github.com/rsmdt/the-startup.git
cd the-startup
go build -o the-startup

# Run tests
go test ./...

# Install from local binary (for development/offline use)
./the-startup install              # Interactive
./the-startup install -y           # Quick global
./the-startup install -ly          # Quick local
```

## Learn More

- [Claude Code Documentation](https://docs.anthropic.com/en/docs/claude-code)
- [Report Issues](https://github.com/rsmdt/the-startup/issues)
- [Contribute](https://github.com/rsmdt/the-startup/pulls)

---

**Ship faster. Ship better. Ship with The Startup.**

*A virtual engineering team that works like the best startups - fast, pragmatic, and focused on results.*
