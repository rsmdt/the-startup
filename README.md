# The Startup - Agent System for Development Tools

A specialized AI agent system that brings expert developers, architects, and project managers to your development workflow. Install once, use everywhere across all your projects.

## What This Does

**The Startup** transforms your development workflow by providing:
- **Expert AI Agents** - Each specialized in different aspects of software development
- **Specialized Commands** - Agents work together on complex tasks

Perfect for developers who want structured, expert-level guidance on architecture, implementation, testing, and project management.

## Quick Start

### 1. Installation

```bash
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh
```

This installs globally to `~/.config/the-startup/` and configures Claude Code to use the agents.

### 2. Start Using Agents

```bash
# Create specifications for new features
/s:specify "Create a user authentication system with JWT tokens"

# Resume working on existing specifications
/s:specify 001

# Implement a completed specification
/s:implement 001-user-auth

# Or work with specific agents directly
/the-architect "Design a microservices architecture for e-commerce"
/the-developer "Implement the user login endpoint with validation"
```

## Available Agents

| Agent | Purpose | Best For |
|-------|---------|----------|
| **the-chief** | Use FIRST for any new request. Evaluates complexity and routes to the right specialist | Complex multi-step tasks, initial assessment |
| **the-architect** | Deep technical design decisions, architecture analysis, pattern evaluation | System design, technical trade-offs, scalability analysis |
| **the-developer** | Implementation with TDD, clean code practices, translating requirements into software | Coding, API endpoints, refactoring, feature implementation |
| **the-business-analyst** | Use FIRST when requirements are vague/unclear. Transforms unclear requests into comprehensive BRDs | Requirements discovery, stakeholder analysis |
| **the-product-manager** | Creates formal PRDs, user stories, implementation roadmaps AFTER requirements are gathered | Product specs, feature prioritization, roadmaps |
| **the-project-manager** | Task coordination, progress tracking, blocker removal for complex implementations | Breaking down work, managing dependencies, execution planning |
| **the-security-engineer** | Security assessments, vulnerability analysis, compliance reviews, incident response | Security reviews, threat analysis, compliance |
| **the-site-reliability-engineer** | Use for ANY error, bug, crash, performance issue, or production incident | Debugging, root cause analysis, performance optimization |
| **the-data-engineer** | Database optimization, data modeling, ETL pipeline design, data architecture | Query optimization, schema design, data infrastructure |
| **the-devops-engineer** | Deployment automation, CI/CD pipelines, infrastructure setup (NOT debugging) | Infrastructure automation, containerization, deployments |
| **the-technical-writer** | Technical documentation, API specs, user guides, clear explanations | Documentation, API docs, user guides, specifications |
| **the-tester** | Comprehensive testing, quality assurance, test strategy, bug detection | Test planning, quality assurance, bug hunting |

## Available Commands

- **`/s:specify`** - Orchestrates development through specialist agents. Creates specifications for new features OR investigates/debugs existing issues. Use with feature description or spec ID to resume (e.g., "001")
- **`/s:implement`** - Executes the implementation plan from a specification. Provide spec ID to implement (e.g., "001" or "001-user-auth")

## Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/rsmdt/the-startup.git
cd the-startup

# Build the binary
go build -o the-startup

# Run tests
go test ./...

# Install locally
./the-startup install
```

### Run commands directly

```bash
# run install command
go run . install
```

### Project Structure

This is a Go project using:
- **Cobra** for CLI commands
- **BubbleTea** for interactive TUI during installation
- **Embedded assets** for agents, hooks, and templates
- **JSONL logging** for agent interaction tracking

## Contributing

To contribute new agents, hooks, or commands:
1. Fork the repository
2. Add your component to the appropriate `assets/` directory
3. Test with a local build and installation
4. Submit a pull request with clear description

## Resources

- [Claude Code Documentation](https://docs.anthropic.com/en/docs/claude-code)
- [Agent Development Guide](https://docs.anthropic.com/en/docs/claude-code/agents)
- [Hooks Reference](https://docs.anthropic.com/en/docs/claude-code/hooks)

---

*Built for developers who want expert guidance and systematic approaches to complex software challenges.*
