# Implementation Plan

*[INSTRUCTION: Create a detailed implementation checklist based on the feature requirements, architecture, and discovered patterns. Organize tasks into logical phases that can be executed sequentially, with clear dependencies and validation points. This entire instruction block should not appear in the final PRD.]*

## Context Documents

*[INSTRUCTION: The orchestrator should automatically discover and ingest context from available specification documents. Use intelligent extraction to identify the most relevant information for implementation. Priority should be given to: 1) Technical requirements from SDD, 2) Feature specifications from PRD, 3) Business context from BRD. When documents are missing, note their absence but proceed with available context. This instruction block should not appear in the final document.]*

### Automatic Context Discovery

The orchestrator will search for and analyze the following documents in priority order:

1. **Primary Specifications** (same directory as PLAN.md):
   - **SDD.md** (Solution Design Document): Technical architecture, API contracts, data models, implementation patterns
   - **PRD.md** (Product Requirements Document): User stories, acceptance criteria, UI/UX specifications, feature behavior
   - **BRD.md** (Business Requirements Document): Business objectives, stakeholder needs, success metrics, constraints

2. **Secondary Context** (if primary specs reference them):
   - Related feature specifications in `../` or sibling directories
   - Shared technical documentation in project root
   - API documentation referenced by the SDD
   - Design system documentation for UI components

3. **Codebase Patterns** (discovered through analysis):
   - Existing implementations of similar features
   - Established patterns in the same module/package
   - Test structures that reveal expected behavior
   - Configuration files that define conventions

### Intelligent Context Extraction

*[ORCHESTRATOR: Extract and prioritize information based on relevance to implementation. Focus on concrete requirements over abstract concepts.]*

#### High Priority Context (Must Extract):
- **Technical Requirements**: API endpoints, data schemas, database models, service interfaces
- **Acceptance Criteria**: Specific behavior requirements, validation rules, edge cases
- **Integration Points**: Dependencies, external services, shared components
- **Security Requirements**: Authentication needs, authorization rules, data protection requirements
- **Performance Targets**: Response time requirements, scalability needs, resource constraints

#### Medium Priority Context (Extract if Present):
- **Business Rules**: Domain logic, calculation formulas, workflow sequences
- **User Experience**: UI flows, interaction patterns, error handling requirements
- **Non-functional Requirements**: Logging needs, monitoring points, audit requirements
- **Migration Considerations**: Backward compatibility, data migration, deprecation plans

#### Low Priority Context (Note if Relevant):
- **Background Information**: Problem history, previous attempts, lessons learned
- **Future Considerations**: Potential extensions, roadmap items, architectural evolution
- **Stakeholder Context**: Team preferences, organizational standards, compliance needs

### Context Synthesis

*[ORCHESTRATOR: After discovery, synthesize the context into actionable implementation guidance. Handle missing documents gracefully.]*

**Available Specifications:**
- *[✓ or ✗ for each document type found]*
- *[List actual file paths discovered]*

**Key Implementation Requirements:**
```markdown
// Example format for extracted requirements
From SDD.md:
- REST API with /api/v1/feature endpoint
- PostgreSQL schema with users and sessions tables
- JWT authentication with 24-hour expiry
- Rate limiting at 100 requests/minute per user

From PRD.md:
- Users must be able to login with email/password
- Session timeout after 30 minutes of inactivity
- Display friendly error messages for all failure cases
- Support "Remember Me" for 30-day sessions

From discovered patterns:
- Project uses middleware pattern for auth (see: internal/middleware/auth.go)
- Error responses follow ErrorResponse struct (see: internal/api/types.go)
- Database migrations use golang-migrate (see: migrations/)
```

**Missing Context Handling:**
- *[If SDD.md missing: Note that technical design decisions will be made based on existing patterns]*
- *[If PRD.md missing: Note that implementation will follow standard UX patterns unless specified]*
- *[If BRD.md missing: Note that business context is limited, focus on technical correctness]*

**Implementation Priorities Based on Context:**
1. *[Highest priority requirement or constraint]*
2. *[Second priority]*
3. *[Third priority]*

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
- **`agent`**: Specifies which specialist agent should execute the task
  - Common agents: `the-developer`, `the-architect`, `the-tester`, `the-reviewer`
  - Default: Selected based on task content if not specified
  
#### Review Configuration
- **`review`**: Boolean flag or conditional expression for review requirements
  - `true`: Always requires review
  - `false`: Skip review (default for simple tasks)
  - `auto`: Let orchestrator decide based on complexity/risk
  
- **`review_focus`**: Specific areas the reviewer should examine
  - Examples: `"security, authentication"`, `"performance, scalability"`, `"patterns, architecture"`
  - Multiple areas can be comma-separated
  
- **`reviewer`**: Explicitly specify reviewing agent (overrides dynamic selection)
  - Examples: `the-architect` for design reviews, `the-security-expert` for security reviews
  - Leave blank for dynamic selection based on context

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
- **`depends_on`**: Task IDs that must complete before this task
  - Format: `"task-1, task-3"` or `"previous"` for immediate predecessor
  - Orchestrator ensures proper execution order
  
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
  
- **`validation_cmd`**: Specific command to validate task completion
  - Example: `"npm test auth.spec.js"`, `"go test ./internal/auth/..."`

### Dynamic Review Selection Logic

The orchestrator intelligently selects reviewers based on multiple factors:

1. **Context-Based Selection**: Analyzes task content and metadata
2. **Expertise Matching**: Routes to agents with relevant expertise
3. **Risk Assessment**: Higher risk triggers stricter review
4. **Workload Balancing**: Distributes reviews across available agents

### Metadata Examples

#### Simple Task (No Review)
```markdown
- [ ] **Update README documentation** [`agent: the-developer`] [`complexity: low`]
  - Add installation instructions
  - Update API examples
```

#### Standard Development Task (Conditional Review)
```markdown
- [ ] **Implement user profile endpoint** [`agent: the-developer`] [`review: auto`] [`complexity: medium`] [`validation: integration`]
  - Create GET /api/users/:id endpoint
  - Add database query with proper indexing
  - Include error handling for missing users
```

#### Complex Task with Dependencies
```markdown
- [ ] **Refactor authentication system** [`agent: the-architect`] [`review: true`] [`reviewer: the-security-expert`] [`complexity: high`] [`risk: critical`] [`depends_on: task-2, task-3`] [`validation: integration`]
  - Design new JWT token structure
  - Implement refresh token rotation
  - Add rate limiting for login attempts
  - Ensure backward compatibility
```

#### Parallel Execution Example
```markdown
- [ ] **Frontend updates** [`agent: the-developer`] [`parallel: true`] [`complexity: medium`]
  - Update login form validation
  - Add loading states
  
- [ ] **Backend updates** [`agent: the-developer`] [`parallel: true`] [`complexity: medium`]
  - Add input validation middleware
  - Update error response format
```

#### Dynamic Review Selection Example
```markdown
- [ ] **Process payment integration** [`agent: the-developer`] [`review: auto`] [`complexity: critical`] [`risk: critical`] [`review_focus: security, error-handling, compliance`]
  # Orchestrator will automatically:
  # 1. Detect payment/financial context
  # 2. Select security-focused reviewer
  # 3. Enforce mandatory review due to critical risk
  # 4. Include compliance checks in validation
```

#### Conditional Review Based on Complexity
```markdown
- [ ] **Database migration** [`agent: the-developer`] [`review: auto`] [`complexity: high`] [`validation: manual`] [`validation_cmd: "./scripts/verify-migration.sh"`]
  # Review triggered if:
  # - Migration affects > 1000 records
  # - Includes destructive operations
  # - Modifies critical tables
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

1. **Be Explicit for Critical Tasks**: Always specify review requirements for security/payment tasks
2. **Use Complexity Indicators**: Help orchestrator make intelligent decisions
3. **Define Dependencies**: Prevent race conditions and ensure correct order
4. **Specify Validation**: Include concrete validation steps for quality assurance
5. **Focus Reviews**: Use `review_focus` to guide reviewer attention
6. **Let Orchestrator Optimize**: Use `auto` for review to leverage intelligent routing

### Backward Compatibility

Tasks without metadata continue to work normally:
```markdown
- [ ] **Simple task without metadata**
  - Orchestrator uses content analysis
  - Applies default routing rules
  - Reviews triggered by keyword detection (e.g., "security", "payment")
```

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

*[INSTRUCTION: The number and nature of phases should match the feature complexity. Simple features might need 2-3 phases, complex ones might need 5-7. Always include context file reading as an early task. Always end with integration testing and final validation. This note should not appear in the final PRD.]*

## Validation Checklist

*[INSTRUCTION: Define validation criteria to ensure the implementation meets all requirements. Use project-specific validation commands identified during research. This note should not appear in the final PRD.]*

### Validation Areas to Consider:

- **Code Quality**: Linting, formatting, type checking
- **Functionality**: All test scenarios pass, features work as specified
- **Integration**: Component interactions, API contracts, data flow
- **Performance**: Response times, resource usage, scalability
- **Security**: Input validation, authorization, data protection
- **Standards**: Code conventions, architectural patterns, best practices

Structure validation tasks based on available project commands and the feature's specific requirements.

## Anti-Patterns to Avoid

### Architecture Anti-Patterns
- ❌ Creating new architectural patterns when established ones exist
- ❌ Modifying unrelated systems "while you're there"
- ❌ Adding external dependencies without checking internal capabilities
- ❌ Changing core conventions without explicit approval
- ❌ Implementing business logic in presentation layer
- ❌ Tight coupling between independent components

### Integration Anti-Patterns
- ❌ Hardcoding external service URLs or credentials
- ❌ Ignoring rate limits and retry mechanisms for external services
- ❌ Exposing internal data structures to external systems
- ❌ Synchronous calls to external services in critical paths
- ❌ Assuming external services are always available

### Data Anti-Patterns
- ❌ Direct database access from presentation layer
- ❌ Storing business logic in database triggers or procedures
- ❌ Missing data validation at application boundaries
- ❌ Inconsistent data state across related entities
- ❌ Exposing database structure through API responses

### Testing Anti-Patterns
- ❌ Testing implementation details instead of behavior
- ❌ Skipping tests for "simple" functions
- ❌ Not testing error conditions and edge cases
- ❌ Over-mocking dependencies in integration tests
- ❌ Writing tests that depend on specific execution order
- ❌ Ignoring test failures or marking them as "flaky"

### Process Anti-Patterns
- ❌ Skipping validation steps to move faster
- ❌ Implementing without understanding existing patterns
- ❌ Making assumptions about user requirements
- ❌ Continuing implementation when blocked on critical decisions
- ❌ Deploying changes without proper testing
- ❌ Ignoring performance implications until production
