# The Agentic Startup - Complete Agent Reference

## Agent Architecture

The Agentic Startup provides 40 specialized agents organized by domain expertise. Each agent focuses on specific activities rather than arbitrary roles, enabling better performance and parallel execution.

## Agent Domains

### 🎯 **the-chief**

Eliminates bottlenecks through smart routing and complexity assessment. The Chief acts as your technical orchestrator, rapidly assessing complexity across multiple dimensions and routing work to appropriate specialists. It enables parallel execution by identifying independent activities and eliminating bottlenecks in your development workflow.

**Key Capabilities:**
- Complexity assessment across security, performance, and architecture dimensions
- Intelligent routing to specialist agents based on task requirements
- Bottleneck identification and elimination
- Parallel execution planning

---

### 📊 **the-analyst**

Transforms vague requirements into actionable specifications. The Analyst team clarifies ambiguous ideas, documents comprehensive specifications, and ensures nothing gets lost in translation.

#### `requirements-clarification`
Uncovers hidden needs and resolves ambiguities through systematic analysis. Identifies edge cases, validates assumptions, and ensures complete understanding before implementation begins.

#### `requirements-documentation`
Creates comprehensive PRDs and requirements specifications. Documents user stories, acceptance criteria, and success metrics with clarity that prevents scope creep.

#### `feature-prioritization`
Data-driven feature prioritization using frameworks like RICE and MoSCoW. Evaluates value vs effort, defines KPIs, and creates roadmaps aligned with strategic objectives.

#### `solution-research`
Researches proven approaches and patterns for solving specific problems. Evaluates alternatives, documents trade-offs, and recommends optimal solutions based on constraints.

#### `project-coordination`
Breaks down complex projects into manageable tasks with clear dependencies. Creates work breakdown structures, maps technical dependencies, and establishes execution timelines.

---

### 🏗️ **the-architect**

Balances elegance with pragmatic business reality. The Architect team ensures your system is built on solid foundations, scalable patterns, and maintainable code.

#### `system-design`
Designs scalable system architectures with comprehensive planning. Creates service designs, selects appropriate technologies, implements scalability patterns, and develops evolutionary roadmaps.

#### `system-documentation`
Creates architectural documentation including design decision records, system diagrams, integration guides, and operational runbooks. Maintains living documentation that evolves with the system.

#### `architecture-review`
Validates design patterns and ensures architectural compliance. Reviews implementations for scalability, maintainability, and alignment with architectural principles.

#### `code-review`
Elevates team capabilities through constructive feedback. Reviews for quality, security, patterns, and best practices while mentoring through actionable suggestions.

#### `scalability-planning`
Ensures systems scale gracefully from MVP to enterprise. Plans for 10x, 100x, and 1000x growth with appropriate architectural evolution strategies.

#### `technology-evaluation`
Makes framework and tool decisions based on comprehensive analysis. Evaluates vendors, creates proof-of-concepts, analyzes trade-offs, and provides recommendations.

#### `technology-standards`
Prevents technology chaos through consistent standards. Establishes coding conventions, architectural patterns, tooling consistency, and cross-team alignment strategies.

---

### 💻 **the-software-engineer**

Ships features that actually work. The Software Engineer team builds robust, maintainable code across the full stack with focus on quality and user experience.

#### `api-design`
REST/GraphQL APIs with clear contracts and excellent developer experience. Designs endpoints, defines schemas, implements versioning strategies, and creates SDKs.

#### `api-documentation`
Comprehensive API documentation that developers love. Creates interactive docs with examples, implements OpenAPI/GraphQL schemas, and maintains postman collections.

#### `database-design`
Balanced schemas for any database type. Designs relational models, NoSQL structures, implements migrations, and optimizes for specific query patterns.

#### `service-integration`
Reliable service communication patterns. Implements REST/GraphQL clients, message queues, event streaming, and webhook handlers with proper error handling.

#### `component-architecture`
Reusable UI components with clear interfaces. Designs component hierarchies, implements design systems, manages props/state, and ensures accessibility compliance.

#### `business-logic`
Domain rules and validation that match real business needs. Implements complex workflows, validation rules, calculations, and business constraints correctly.

#### `reliability-engineering`
Error handling and resilience patterns. Implements circuit breakers, retries with backoff, graceful degradation, and comprehensive error recovery.

#### `performance-optimization`
Bundle size and Core Web Vitals optimization. Profiles applications, implements lazy loading, optimizes critical paths, and improves user-perceived performance.

#### `state-management`
Client and server state patterns that scale. Implements Redux/MobX/Zustand patterns, server state with React Query/SWR, and optimistic updates.

#### `browser-compatibility`
Cross-browser support without compromises. Implements polyfills, handles vendor prefixes, tests across browsers, and provides progressive enhancement.

---

### 🚀 **the-platform-engineer**

Makes systems that don't wake you at 3am. The Platform Engineer team ensures your infrastructure is reliable, scalable, and observable.

#### `system-performance`
Handles 10x load without 10x cost. Profiles applications, optimizes databases, implements caching strategies, and reduces resource consumption.

#### `observability`
Monitoring that catches problems early. Implements metrics, logging, tracing, and alerting with proper SLI/SLO definitions.

#### `containerization`
Consistent deployment everywhere. Creates optimized Docker images, Kubernetes manifests, and container orchestration configurations.

#### `pipeline-engineering`
Reliable data processing at scale. Builds ETL/ELT pipelines, stream processing systems, and batch jobs with monitoring and error handling.

#### `ci-cd-automation`
Safe deployments at scale. Designs CI/CD pipelines, implements quality gates, automates testing, and ensures reliable releases.

#### `deployment-strategies`
Progressive rollouts that minimize risk. Implements blue-green deployments, canary releases, feature flags, and automated rollbacks.

#### `incident-response`
Production fire debugging expertise. Creates runbooks, implements on-call procedures, performs root cause analysis, and improves system resilience.

#### `infrastructure-as-code`
Reproducible infrastructure through code. Writes Terraform/CloudFormation, manages state, implements modules, and automates provisioning.

#### `storage-architecture`
Scalable storage solutions for any need. Designs object storage, implements CDNs, optimizes file systems, and manages data lifecycle.

#### `query-optimization`
Fast database queries through systematic optimization. Analyzes execution plans, creates indexes, rewrites queries, and implements partitioning strategies.

#### `data-modeling`
Balanced data models for performance and maintainability. Designs schemas, implements denormalization strategies, and optimizes for specific access patterns.

---

### 🎨 **the-designer**

Creates products people actually want to use. The Designer team ensures your product is usable, accessible, and delightful.

#### `accessibility-implementation`
WCAG 2.1 AA compliance and beyond. Implements ARIA attributes, keyboard navigation, screen reader support, and ensures universal usability.

#### `user-research`
Real user needs, not assumptions. Conducts interviews, usability testing, creates personas, and translates insights into actionable recommendations.

#### `interaction-design`
Minimal friction user flows. Designs navigation, user journeys, wireframes, and interaction patterns that feel intuitive.

#### `visual-design`
Brand-enhancing UI aesthetics. Creates design systems, color palettes, typography scales, and visual hierarchies that communicate effectively.

#### `design-systems`
Consistent component libraries. Builds design tokens, component libraries, style guides, and documentation for consistent experiences.

#### `information-architecture`
Intuitive content hierarchies. Organizes information, designs taxonomies, creates sitemaps, and improves findability.

---

### 🧪 **the-qa-engineer**

Catches bugs before users do. The QA Engineer team ensures quality through comprehensive testing strategies.

#### `test-strategy`
Risk-based testing approaches. Plans test coverage, designs test cases, identifies critical paths, and creates quality assurance processes.

#### `test-implementation`
Comprehensive test suites that catch real bugs. Writes unit tests, integration tests, E2E tests with proper mocking and assertions.

#### `exploratory-testing`
Creative defect discovery beyond automation. Performs manual testing, finds edge cases, validates user experience, and identifies gaps in automated coverage.

#### `performance-testing`
Load and stress validation. Designs load tests, implements stress testing, validates scalability, and ensures performance under production conditions.

---

### 🔐 **the-security-engineer**

Keeps the bad guys out. The Security Engineer team protects your application and data from threats.

#### `vulnerability-assessment`
OWASP-based security checks. Performs penetration testing, vulnerability scanning, threat modeling, and provides remediation guidance.

#### `authentication-systems`
OAuth, JWT, SSO, MFA implementation. Builds secure authentication flows, manages sessions, implements authorization, and handles identity management.

#### `security-incident-response`
Rapid containment and recovery. Handles breach response, conducts forensics, implements containment, and develops incident procedures.

#### `compliance-audit`
GDPR, SOX, HIPAA compliance verification. Audits systems for regulatory compliance, identifies gaps, and implements necessary controls.

#### `data-protection`
Encryption and privacy controls. Implements encryption at rest/transit, manages keys, ensures data privacy, and handles secure deletion.

---

### 📱 **the-mobile-engineer**

Ships apps users love. The Mobile Engineer team creates native and cross-platform mobile experiences.

#### `mobile-interface-design`
Platform-specific UI patterns. Implements iOS Human Interface Guidelines, Material Design, creates responsive layouts, and handles device variations.

#### `mobile-data-persistence`
Offline-first strategies. Implements local databases, sync mechanisms, conflict resolution, and handles spotty connectivity gracefully.

#### `cross-platform-integration`
Native and hybrid bridges. Integrates React Native with native modules, implements platform-specific features, and manages code sharing.

#### `mobile-deployment`
App store submissions. Handles certificates, provisioning, store metadata, screenshots, and manages the submission process.

#### `mobile-performance`
Battery and memory optimization. Profiles apps, reduces memory usage, optimizes battery consumption, and improves app startup time.

---

### 🤖 **the-ml-engineer**

Makes AI that actually ships. The ML Engineer team brings machine learning from research to production.

#### `model-deployment`
Production-ready inference. Containerizes models, implements serving infrastructure, manages versioning, and monitors performance.

#### `ml-monitoring`
Drift detection systems. Tracks model performance, detects data drift, monitors predictions, and triggers retraining pipelines.

#### `prompt-optimization`
LLM prompt engineering. Designs system prompts, implements few-shot examples, tests variations, and manages prompt versions.

#### `mlops-automation`
Reproducible ML pipelines. Automates training, implements experiment tracking, manages model registry, and ensures reproducibility.

#### `context-management`
AI memory architectures. Builds RAG systems, implements conversation memory, manages context windows, and ensures coherent interactions.

#### `feature-engineering`
Model-ready data pipelines. Creates feature pipelines, implements feature stores, validates data quality, and monitors feature drift.

---

### 🔧 **the-meta-agent**

Creates new specialized agents. The Meta-Agent designs and generates new Claude Code sub-agents, validates specifications, and refactors existing agents following evidence-based design principles.

**Key Capabilities:**
- Agent generation following Claude Code specifications
- Agent validation and compliance checking
- Refactoring existing agents for improved performance
- Applying evidence-based agent design patterns

---

## Agent Delegation Patterns

### Parallel Execution

Most tasks benefit from parallel agent execution. The system automatically identifies independent activities and launches appropriate specialists simultaneously.

```
Example: "Add user authentication"
→ Parallel execution:
   - database-design: User schema
   - api-design: Auth endpoints
   - component-architecture: Login UI
   - authentication-systems: OAuth flow
```

### Sequential Dependencies

When tasks have dependencies, agents execute in proper sequence with context passing.

```
Example: "Migrate database schema"
→ Sequential execution:
   1. database-design: New schema
   2. data-modeling: Migration plan
   3. test-implementation: Migration tests
   4. deployment-strategies: Rollout plan
```

### Specialist Selection

The system selects agents based on activity matching, not role assumptions:

- ✅ "Build API endpoints" → `api-design` + `api-documentation`
- ✅ "Optimize queries" → `query-optimization` + `system-performance`
- ❌ "Backend work" → Too vague, needs activity specification

## Using Agents Effectively

### Best Practices

1. **Be specific about activities**: "Implement caching" is better than "improve performance"
2. **Let the system parallelize**: Don't manually sequence unless dependencies exist
3. **Trust specialist selection**: The system matches activities to appropriate agents
4. **Review agent outputs**: Agents provide recommendations, you make decisions

### Common Patterns

#### Feature Development
```
requirements-clarification → system-design → parallel(api-design, component-architecture, database-design) → test-implementation
```

#### Performance Optimization
```
system-performance → parallel(query-optimization, performance-optimization) → performance-testing → deployment-strategies
```

#### Security Hardening
```
vulnerability-assessment → parallel(authentication-systems, data-protection) → security-incident-response planning
```

#### Refactoring
```
code-review → architecture-review → parallel(targeted refactoring by specialists) → test-implementation
```

## Agent Evolution

Agents continuously improve through:

- **Community feedback**: Real usage shapes agent behavior
- **Pattern recognition**: Common workflows become optimized
- **Capability expansion**: New activities added based on needs
- **Performance tuning**: Response quality and speed improvements

Share your agent experiences and suggestions at [github.com/rsmdt/the-startup](https://github.com/rsmdt/the-startup).