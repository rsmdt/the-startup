# Implementation Plan

*[INSTRUCTION: Create a detailed implementation checklist based on the feature requirements, architecture, and discovered patterns. Organize tasks into logical phases that can be executed sequentially, with clear dependencies and validation points. This entire instruction block should not appear in the final PRD.]*

## Context Documents

*[INSTRUCTION: Automatically discover and ingest context from available specification documents. Use intelligent extraction to identify the most relevant information for implementation. Priority should be given to: 1) Technical requirements from SDD, 2) Feature specifications from PRD, 3) Business context from BRD. When documents are missing, note their absence but proceed with available context. This instruction block should not appear in the final document.]*

### Context Validation Checkpoints

Before proceeding with implementation, verify:
- [ ] All referenced specifications have been located and read
- [ ] Key requirements have been extracted and understood
- [ ] Conflicts between specifications have been identified and resolved
- [ ] Missing context has been noted with mitigation strategies
- [ ] Implementation approach aligns with discovered patterns

## Task Metadata Guidelines

*[INSTRUCTION: Each task can include metadata to guide intelligent orchestration and dynamic review selection. The orchestrator uses this metadata to route tasks to appropriate agents, determine review requirements, and manage dependencies. All metadata fields are optional but enhance execution quality. This instruction should not appear in the final document.]*

### Core Metadata Fields

Tasks support the following metadata fields for orchestration:

#### Execution Control
- **`agent`** (optional): Specifies which specialist sub-agent should execute the task
  - Only use when specific expertise is required
  - If omitted, orchestrator selects best agent based on task requirements
  - Examples: `the-architect` for design tasks, `the-security-engineer` for auth tasks
  
#### Review Configuration
- **`review`**: Areas to focus on during review (presence indicates review required)
  - Examples: `"security, authentication"`, `"performance, scalability"`, `"patterns, architecture"`
  - Multiple areas can be comma-separated
  - Omit entirely if no review needed
  - Can be empty to indicate automatic review area selection
  - Reviewer specialist sub-agent is selected automatically

#### Complexity and Risk Indicators
- **`complexity`**: Task complexity level driving review decisions
  - `low`: Simple, straightforward tasks (review optional)
  - `medium`: Standard implementation tasks (review recommended)
  - `high`: Complex logic or critical functionality (review required)
  - `critical`: Security/payment/data-critical tasks (mandatory review)
  
- **`risk`**: Business/technical risk assessment
  - `minimal`: Low impact on system/users
  - `moderate`: Standard business logic
  - `high`: Core functionality, user-facing features
  - `critical`: Security, payments, data integrity

#### Task Dependencies and Ordering
- **`blocks`**: Task IDs that cannot start until this completes
  - Useful for identifying critical path items
  
- **`parallel`**: Can execute simultaneously with other parallel tasks
  - Boolean: `true` allows concurrent execution with other parallel tasks

#### Validation Strategies
- **`validation`**: Type of validation required after task completion
  - `unit`: Run unit tests for affected code
  - `integration`: Run integration test suite
  - `manual`: Requires human verification
  - `automated`: Run automated validation scripts
  - `performance`: Execute performance benchmarks

### Dynamic **Review** Selection Logic

Intelligently select reviewer sub-agent based on multiple factors:
- **Context-Based Selection**: Analyzes task content and metadata
- **Expertise Matching**: Routes to agents with relevant expertise
- **Risk Assessment**: Higher risk triggers stricter review
- **Workload Balancing**: Distributes reviews across available agents


### Dynamic **Validate** Logic

Intelligently consider:
- **Code Quality**: Linting, formatting, type checking
- **Functionality**: All test scenarios pass, features work as specified
- **Integration**: Component interactions, API contracts, data flow
- **Performance**: Response times, resource usage, scalability
- **Security**: Input validation, authorization, data protection
- **Standards**: Code conventions, architectural patterns, best practices

### Metadata Examples

Simple Task (No Review):
```markdown
- [ ] **Phase 1:** Update README documentation [`complexity: low`]
  - [ ] Add installation instructions
  - [ ] Update API examples
  # Orchestrator will select appropriate agent
```

Standard Development Task (With Review):
```markdown
- [ ] **Phase 2:** Implement user profile endpoint [`complexity: medium`]
  - [ ] Create GET /api/users/:id endpoint
  - [ ] Add database query with proper indexing
  - [ ] Include error handling for missing users
  - [ ] **Validate** [specific validation command or check]
  - [ ] **Review** [specific review agent and area]
  # Review required with focus on database and error handling
```

Complex Task with Dependencies:
```markdown
- [ ] **Phase 3:** Refactor authentication system [`complexity: high`] [`risk: critical`]
  - [ ] Design new JWT token structure
  - [ ] Implement refresh token rotation
  - [ ] Add rate limiting for login attempts
  - [ ] Ensure backward compatibility
  - [ ] **Validate** [specific validation command or check]
  - [ ] **Review** [specific review agent and area]
  # Critical task with security review required
```

Complex Parallel Execution Example:
```markdown
- [ ] **Phase 1:** Login Validation
  - [ ] Frontend updates [`agent: the-developer`] [`parallel: true`] [`complexity: medium`]
    - [ ] Update login form validation
    - [ ] Add loading states
    - [ ] **Validate** [specific validation command or check]
    - [ ] **Review** [specific review agent and area]
    
  - [ ] Backend updates [`agent: the-developer`] [`parallel: true`] [`complexity: medium`]
    - [ ] Add input validation middleware
    - [ ] Update error response format
    - [ ] **Validate** [specific validation command or check]
    - [ ] **Review** [specific review agent and area]
```

Payment Task (Critical Review)
```markdown
- [ ] **Phase 1:** Process payment integration** [`review: security, compliance, error-handling`] [`complexity: critical`] [`risk: critical`]
  # Critical risk + payment context = mandatory security review
  # Orchestrator selects security-focused reviewer
```

### Orchestration Behavior

Based on metadata, the orchestrator will:
1. **Route Tasks**: Send to specified agent or select based on content
2. **Manage Dependencies**: Ensure proper execution order
3. **Trigger Reviews**: Apply review logic based on metadata and context
4. **Select Reviewers**: Choose appropriate reviewer dynamically
5. **Run Validation**: Execute specified validation after completion
6. **Handle Failures**: Retry with different agents or escalate to human

### Best Practices

- Only specify agent when specific expertise is mandatory
- Use Complexity Indicators to help make intelligent decisions
- Always Specify Validation and Review for quality assurance

## Checklist Structure Guidelines

Organize implementation tasks into phases that make sense for this specific feature:
- Group related tasks that can be worked on together
- Identify dependencies between tasks
- Include validation points after each significant milestone
- Consider the feature's architecture when determining phases

Each task should be:
- [ ] Specific and actionable
- [ ] Independently verifiable as complete
- [ ] Sized appropriately (not too large, not too granular)

Include validation commands from the Project Commands section at appropriate checkpoints.

## Example Phase Structure (adapt based on feature):

**Phase X: [Descriptive Phase Name]**
- [ ] [Specific task with clear completion criteria]
- [ ] [Another related task]
- [ ] **Validation**: [Specific validation command or check]
- [ ] **Review**: [specific review agent and area]

*[INSTRUCTION: The number and nature of phases should match the feature complexity. Always include context file reading as Phase 1. Always end with integration testing and final validation. This note should not appear in the final PRD.]*

Structure validation tasks based on available project commands and the feature's specific requirements.

## Anti-Patterns to Avoid

- Skipping validation steps to move faster
- Implementing without understanding existing patterns
- Making assumptions about user requirements
- Continuing implementation when blocked on critical decisions
