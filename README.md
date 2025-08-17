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

### Core System Agents

| Agent | Purpose | Best For |
|-------|---------|----------|
| **the-chief** | Use FIRST for any new request. Evaluates complexity and routes to the right specialist | Complex multi-step tasks, initial assessment |
| **the-architect** | Deep technical design decisions, architecture analysis, pattern evaluation | System design, technical trade-offs, scalability analysis |
| **the-developer** | Implementation with TDD, clean code practices, translating requirements into software | Coding, API endpoints, refactoring, feature implementation |
| **the-lead-developer** | Code review specialist for AI-generated code, mentorship, refactoring decisions | Code quality assessment, architectural improvements, team mentorship |

### Business & Planning Agents

| Agent | Purpose | Best For |
|-------|---------|----------|
| **the-business-analyst** | Use FIRST when requirements are vague/unclear. Transforms unclear requests into comprehensive BRDs | Requirements discovery, stakeholder analysis |
| **the-product-manager** | Creates formal PRDs, user stories, implementation roadmaps AFTER requirements are gathered | Product specs, feature prioritization, roadmaps |
| **the-project-manager** | Task coordination, progress tracking, blocker removal for complex implementations | Breaking down work, managing dependencies, execution planning |

### Quality & Operations Agents

| Agent | Purpose | Best For |
|-------|---------|----------|
| **the-security-engineer** | Security assessments, vulnerability analysis, compliance reviews, incident response | Security reviews, threat analysis, compliance |
| **the-site-reliability-engineer** | Use for ANY error, bug, crash, performance issue, or production incident | Debugging, root cause analysis, performance optimization |
| **the-tester** | Comprehensive testing, quality assurance, test strategy, bug detection | Test planning, quality assurance, bug hunting |

### Infrastructure & Data Agents

| Agent | Purpose | Best For |
|-------|---------|----------|
| **the-data-engineer** | Database optimization, data modeling, ETL pipeline design, data architecture | Query optimization, schema design, data infrastructure |
| **the-devops-engineer** | Deployment automation, CI/CD pipelines, infrastructure setup (NOT debugging) | Infrastructure automation, containerization, deployments |

### Design & Experience Agents

| Agent | Purpose | Best For |
|-------|---------|----------|
| **the-ux-designer** | User interface design, accessibility compliance, design systems, user interaction patterns | UI/UX design, WCAG compliance, user experience optimization |

### Documentation & Communication Agents

| Agent | Purpose | Best For |
|-------|---------|----------|
| **the-technical-writer** | Technical documentation, API specs, user guides, clear explanations | Documentation, API docs, user guides, specifications |

### Compliance & Risk Agents

| Agent | Purpose | Best For |
|-------|---------|----------|
| **the-compliance-officer** | Regulatory compliance, data privacy laws, AI governance, audit trails | GDPR/CCPA compliance, industry regulations, governance frameworks |

### Specialized Engineering Agents

| Agent | Purpose | Best For |
|-------|---------|----------|
| **the-prompt-engineer** | AI prompt optimization, conversation design, model interaction patterns | Prompt engineering, AI system optimization, conversational interfaces |
| **the-context-engineer** | Context management, information architecture, knowledge organization | Context optimization, information flow, system integration |

## Available Commands

- **`/s:specify`** - Orchestrates development through specialist agents. Creates specifications for new features OR investigates/debugs existing issues. Use with feature description or spec ID to resume (e.g., "001")
- **`/s:implement`** - Executes the implementation plan from a specification. Provide spec ID to implement (e.g., "001" or "001-user-auth")

## Agent Capability Matrix

| Capability | Core | Business | Quality | Infra | Design | Docs | Compliance | Specialized |
|------------|------|----------|---------|--------|--------|------|------------|-------------|
| **Initial Assessment** | ✅ the-chief | | | | | | | |
| **Requirements Analysis** | | ✅ the-business-analyst | | | | | | |
| **System Architecture** | ✅ the-architect | | | | | | | |
| **Code Implementation** | ✅ the-developer | | | | | | | |
| **Code Review** | ✅ the-lead-developer | | | | | | | |
| **Product Planning** | | ✅ the-product-manager | | | | | | |
| **Project Coordination** | | ✅ the-project-manager | | | | | | |
| **Security Assessment** | | | ✅ the-security-engineer | | | | | |
| **Bug Investigation** | | | ✅ the-site-reliability-engineer | | | | | |
| **Quality Assurance** | | | ✅ the-tester | | | | | |
| **Data Architecture** | | | | ✅ the-data-engineer | | | | |
| **Infrastructure Setup** | | | | ✅ the-devops-engineer | | | | |
| **UI/UX Design** | | | | | ✅ the-ux-designer | | | |
| **Documentation** | | | | | | ✅ the-technical-writer | | |
| **Regulatory Compliance** | | | | | | | ✅ the-compliance-officer | |
| **AI Optimization** | | | | | | | | ✅ the-prompt-engineer |
| **Context Management** | | | | | | | | ✅ the-context-engineer |

## Agent Discovery Guide

### 🚀 Starting a New Project

**"I have a vague idea and need to get started"**
→ Start with **the-business-analyst** to clarify requirements
→ Then **the-product-manager** to create formal specifications
→ Finally **the-architect** for technical design

**"I have clear requirements and need technical design"**
→ Start with **the-architect** for system design
→ Then **the-project-manager** to break down implementation
→ Then **the-developer** for coding

### 🔧 During Development

**"I need to implement a feature"**
→ Use **the-developer** for implementation
→ Then **the-lead-developer** for code review
→ Then **the-tester** for quality assurance

**"Something is broken or not working"**
→ Use **the-site-reliability-engineer** for debugging
→ Use **the-security-engineer** if security-related
→ Use **the-data-engineer** if database-related

### 🎨 Design & User Experience

**"I need to design user interfaces"**
→ Use **the-ux-designer** for UI/UX design and accessibility
→ Then **the-developer** to implement the design
→ Then **the-tester** to validate user experience

### 📋 Compliance & Documentation

**"I need to handle regulations or privacy"**
→ Use **the-compliance-officer** for regulatory requirements
→ Use **the-security-engineer** for security compliance
→ Use **the-technical-writer** for documentation

### 🤖 AI & Optimization

**"I need to optimize AI interactions"**
→ Use **the-prompt-engineer** for conversation design
→ Use **the-context-engineer** for information flow
→ Use **the-lead-developer** for AI code review

## Common Agent Collaboration Patterns

### Pattern 1: Full Feature Development
```
the-business-analyst → the-product-manager → the-architect → 
the-developer → the-lead-developer → the-tester → the-technical-writer
```

### Pattern 2: Bug Investigation & Fix
```
the-site-reliability-engineer → the-developer → 
the-lead-developer → the-tester
```

### Pattern 3: Security Implementation
```
the-security-engineer → the-compliance-officer → 
the-architect → the-developer → the-tester
```

### Pattern 4: Design-Led Development
```
the-ux-designer → the-architect → the-developer → 
the-tester → the-technical-writer
```

### Pattern 5: Data-Driven Feature
```
the-data-engineer → the-architect → the-developer → 
the-tester → the-technical-writer
```

### Pattern 6: Infrastructure Setup
```
the-devops-engineer → the-security-engineer → 
the-site-reliability-engineer → the-technical-writer
```

## Training Examples for New Agents

### the-lead-developer Examples

**Code Review Scenario:**
```
Context: AI-generated authentication module
User: "The developer agent just created the auth system"
Assistant: "I'll use the-lead-developer agent to review for security and quality."
Result: Comprehensive code review with specific improvements
```

**Refactoring Scenario:**
```
Context: Legacy codebase with performance issues
User: "This codebase has grown messy with duplicate patterns"
Assistant: "Let me use the-lead-developer agent to identify refactoring opportunities."
Result: Architectural improvements and technical debt reduction
```

### the-ux-designer Examples

**Interface Design Scenario:**
```
Context: New dashboard feature needed
User: "We need an analytics dashboard for administrators"
Assistant: "I'll use the-ux-designer agent to create an intuitive, accessible design."
Result: Complete UI/UX design with accessibility compliance
```

**Accessibility Audit Scenario:**
```
Context: Existing app needs compliance
User: "Our application needs to meet WCAG 2.1 AA standards"
Assistant: "Let me use the-ux-designer agent to audit and improve accessibility."
Result: Comprehensive accessibility improvements and compliance plan
```

### the-compliance-officer Examples

**Data Privacy Scenario:**
```
Context: User data collection system
User: "We're collecting user emails and location data"
Assistant: "I'll use the-compliance-officer agent to ensure GDPR compliance."
Result: Complete privacy framework with consent mechanisms
```

**Industry Regulation Scenario:**
```
Context: Healthcare application development
User: "Building a patient records management system"
Assistant: "Let me use the-compliance-officer agent for HIPAA compliance."
Result: Healthcare-specific compliance framework and audit procedures
```

## Troubleshooting Agent Invocation

### Agent Not Being Used Automatically?

**Problem**: Claude does the work directly instead of using agents
**Solution**: Use explicit invocation patterns:
- "Use the-[agent-name] agent to [task]"
- "I need the-[agent-name] specialist for [scenario]"
- "Let me get the-[agent-name] agent to handle [problem]"

### Wrong Agent Selected?

**Problem**: Claude selects inappropriate agent for the task
**Solution**: Be more specific about the domain:
- For bugs: "Use the-site-reliability-engineer to debug this error"
- For design: "Use the-ux-designer to improve this interface"
- For security: "Use the-security-engineer to assess these risks"

### Agent Missing Context?

**Problem**: Agent doesn't have enough information
**Solution**: Provide context explicitly:
- Reference previous conversations: "Building on the authentication system we discussed"
- Include relevant files: "Looking at the user.js model file"
- Specify constraints: "Working within our React/TypeScript stack"

### Complex Task Not Decomposed?

**Problem**: Single agent trying to handle multi-domain task
**Solution**: Use the-chief first:
- "This seems complex - use the-chief to evaluate and route appropriately"
- "I need multiple specialists for this - start with the-chief"

### Multiple Agents Needed?

**Problem**: Task requires several different specializations
**Solution**: Expect task handoffs:
- Agents will create `<tasks>` blocks assigning work to other agents
- Follow the suggested workflow
- Each agent focuses on their expertise area

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
