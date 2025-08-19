<p align="center"><img src="https://github.com/rsmdt/the-startup/blob/main/assets/logo.png" width="400" alt="The Startup"></p>

<p align="center">Ship faster. Ship better. Ship with <b>The Agentic Startup</b>.</p>

## What is The Startup?

The Startup brings you instant access to expert developers, architects, and engineers - all working together to turn your ideas into shipped code.

The Startup is an orchestration system for Claude Code that gives you a virtual engineering team. Instead of one AI trying to do everything, you get specialized experts who collaborate like a real startup team - pragmatic, fast, and focused on shipping.

Think of it as having a CTO, architects, developers, and DevOps engineers on-demand, each bringing their expertise to your project.

## Installation

Install globally with one command:

```bash
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh
```

## The Startup Way

When you use The Startup, Claude Code becomes your **technical co-founder** with a specific personality:

- **Ship over perfection** - MVPs today, not perfect solutions next quarter
- **Pragmatic decisions** - Make the call with available info, iterate later  
- **Specialist delegation** - Pull in the right expert for each task
- **Parallel execution** - Multiple experts work simultaneously when possible

You're not talking to a generic AI - you're working with a startup CTO who knows when to be scrappy and when to bring in the specialists.

## Your Expert Team

### 🎯 Start Here
- **the-chief** - Your first stop for any complex request. Assesses complexity and routes to the right experts
- **the-business-analyst** - Turns vague ideas into clear requirements. Use when you're not sure what you need

### 💻 Building Features
- **the-architect** - System design and technical decisions. Makes the big architectural calls
- **the-developer** - Writes clean, tested code. Your implementation specialist
- **the-lead-developer** - Reviews code for quality and mentors through improvements

### 🔧 Fixing & Optimizing  
- **the-site-reliability-engineer** - Debugs errors and solves production issues
- **the-data-engineer** - Optimizes databases and designs data architectures
- **the-devops-engineer** - Automates deployments and infrastructure

### 🎨 Design & Experience
- **the-ux-designer** - Creates intuitive interfaces and ensures accessibility
- **the-technical-writer** - Writes clear documentation and API specs

### 🛡️ Security & Quality
- **the-security-engineer** - Identifies vulnerabilities and implements secure practices
- **the-tester** - Creates comprehensive test strategies and finds bugs
- **the-compliance-officer** - Ensures regulatory compliance (GDPR, HIPAA, etc.)

### 📋 Planning & Management
- **the-product-manager** - Creates formal specs and roadmaps from requirements
- **the-project-manager** - Breaks down work and tracks implementation progress

### 🤖 Specialized Experts
- **the-prompt-engineer** - Optimizes AI interactions and conversation design
- **the-context-engineer** - Manages information flow and system integration

## Commands

The Startup provides commands:

### `/s:specify` - Plan Before You Build

Creates comprehensive specifications from your ideas:

```bash
# Start fresh with a new feature idea
/s:specify "Build a real-time notification system"

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

## Real-World Examples

### Building Authentication
```
You: /s:specify "Add user authentication with JWT"
```
The Startup will:
- Use **the-business-analyst** to clarify requirements (OAuth? 2FA? Password reset?)
- Bring in **the-architect** to design the system
- Get **the-security-engineer** to review for vulnerabilities
- Create a complete implementation plan

### Debugging Production Issues
```
You: "The API is returning 500 errors on user login"
```
The Startup will:
- Immediately call **the-site-reliability-engineer** to investigate
- Once root cause is found, bring in **the-developer** to fix
- Have **the-tester** verify the fix
- Get **the-lead-developer** to review the changes

### Creating a Dashboard
```
You: /s:specify "Admin dashboard for monitoring system metrics"
```
The Startup will:
- Use **the-ux-designer** to create the interface design
- Bring in **the-data-engineer** for efficient data queries
- Get **the-architect** to design the real-time data flow
- Coordinate implementation across frontend and backend

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

## Building from Source

If you want to contribute or customize:

```bash
# Clone and build
git clone https://github.com/rsmdt/the-startup.git
cd the-startup
go build -o the-startup

# Run tests
go test ./...

# Install locally
./the-startup install
```

## Learn More

- [Claude Code Documentation](https://docs.anthropic.com/en/docs/claude-code)
- [Report Issues](https://github.com/rsmdt/the-startup/issues)
- [Contribute](https://github.com/rsmdt/the-startup/pulls)

---

**Ship faster. Ship better. Ship with The Startup.**

*A virtual engineering team that works like the best startups - fast, pragmatic, and focused on results.*
