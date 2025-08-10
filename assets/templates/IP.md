# Implementation Plan (IP)
**Project:** [Project Name]  
**Version:** 1.0  
**Date:** [Date]  
**Author:** the-project-manager  
**Complexity:** [Simple/Medium/Complex]

## Execution Rules for Main Agent

### How to Read This Document
This implementation plan uses a checklist format optimized for LLM execution while remaining human-readable.

### Phase Execution Types
- **parallel**: All tasks in the phase can be executed simultaneously
  - Invoke multiple agents at once using batch Task tool calls
  - Monitor all parallel executions
  - Wait for all to complete before proceeding
  
- **sequential**: Tasks must be completed in the order listed
  - Complete each task fully before starting the next
  - Pass outputs from previous tasks as context to subsequent tasks
  - Stop on any failures

### Task Annotations
Each task line contains:
1. **Checkbox**: `- [ ]` tracks completion status
2. **Task description**: What needs to be done
3. **Agent assignment**: `{agent: specialist-name}` 
4. **Source reference**: `[→ doc#section]` links to requirements

Example:
```
- [ ] Implement user authentication {agent: developer} [→ PRD#auth-requirements]
```

### Subtasks
Indented items are subtasks that should be included in the agent's prompt:
```
- [ ] Main task {agent: developer}
  - [ ] Subtask 1 (include in prompt)
  - [ ] Subtask 2 (include in prompt)
```

### Validation Checkpoints
**Validation** tasks use available project commands to verify progress:
```
- [ ] **Validation**: npm test, npm run lint
```

### Status Tracking
As you execute:
- Mark `- [ ]` as `- [x]` when complete
- Update the status in memory/context
- Report progress to user periodically

### Available Commands
```bash
# Setup
npm install              # Install dependencies
npm run dev              # Start development

# Validation
npm run lint             # Code quality
npm run typecheck        # Type safety
npm test                 # Run tests
npm run build            # Build project

# Database
npm run db:migrate       # Run migrations
npm run db:seed          # Seed test data
```

## Phase 1: Foundation & Analysis
**Execution**: parallel  
**Dependencies**: Project approved

- [ ] Read and analyze existing codebase {agent: architect} [→ SDD#context]
  - [ ] Identify architectural patterns
  - [ ] Map component structure
  - [ ] Document integration points
- [ ] Clarify business requirements {agent: business-analyst} [→ BRD#requirements]
  - [ ] Validate assumptions
  - [ ] Identify missing requirements
  - [ ] Get stakeholder confirmation
- [ ] Setup development environment {agent: devops} [→ SDD#deployment]
  - [ ] Initialize repository
  - [ ] Configure CI/CD pipeline
  - [ ] Setup environments (dev/staging/prod)
- [ ] **Validation**: All setup commands pass, requirements documented

## Phase 2: Core Infrastructure
**Execution**: sequential  
**Dependencies**: Phase 1 complete

- [ ] Design system architecture {agent: architect} [→ SDD#architecture]
  - [ ] Create component diagrams
  - [ ] Define API contracts
  - [ ] Plan data models
- [ ] Implement data layer {agent: data-engineer} [→ SDD#data-design]
  - [ ] Create database schema
  - [ ] Setup migrations
  - [ ] Add seed data
- [ ] Setup authentication {agent: developer} [→ PRD#auth-requirements]
  - [ ] Implement auth flow
  - [ ] Add authorization checks
  - [ ] Create user management
- [ ] Security hardening {agent: security-engineer} [→ SDD#security]
  - [ ] Audit auth implementation
  - [ ] Add input validation
  - [ ] Configure security headers
- [ ] **Validation**: npm test, security scan passes

## Phase 3: Feature Implementation
**Execution**: parallel  
**Dependencies**: Phase 2 complete, auth working

- [ ] Backend API development {agent: developer} [→ SDD#api-design]
  - [ ] Implement REST endpoints
    - [ ] User endpoints [→ PRD#user-stories-1]
    - [ ] Resource endpoints [→ PRD#user-stories-2]
  - [ ] Add business logic [→ SDD#business-rules]
  - [ ] Implement error handling
  - [ ] **Validation**: API tests pass
  
- [ ] Frontend components {agent: developer} [→ PRD#ui-requirements]
  - [ ] Create base components
    - [ ] Navigation [→ PRD#navigation]
    - [ ] Forms [→ PRD#forms]
    - [ ] Data displays [→ PRD#displays]
  - [ ] Implement state management
  - [ ] Add loading/error states
  - [ ] **Validation**: Component tests pass

- [ ] Integration layer {agent: developer} [→ SDD#integration]
  - [ ] Connect frontend to API
  - [ ] Add caching strategy
  - [ ] Handle offline scenarios
  - [ ] **Validation**: Integration tests pass

## Phase 4: Quality Assurance
**Execution**: sequential  
**Dependencies**: Phase 3 complete

- [ ] Comprehensive testing {agent: tester} [→ PRD#test-scenarios]
  - [ ] Unit test coverage >80%
  - [ ] Integration test critical paths
  - [ ] E2E test user journeys
    - [ ] Happy path flows [→ PRD#scenario-1]
    - [ ] Error scenarios [→ PRD#scenario-2]
    - [ ] Edge cases [→ PRD#scenario-3]
  - [ ] Accessibility testing
  - [ ] **Validation**: All test suites pass

- [ ] Performance optimization {agent: site-reliability-engineer} [→ SDD#performance]
  - [ ] Run load tests
  - [ ] Optimize slow queries
  - [ ] Implement caching
  - [ ] Monitor resource usage
  - [ ] **Validation**: Meets performance targets

- [ ] Documentation {agent: technical-writer} [→ PRD#documentation]
  - [ ] API documentation
  - [ ] User guides
  - [ ] Deployment guide
  - [ ] **Validation**: Docs review complete

## Phase 5: Deployment
**Execution**: sequential  
**Dependencies**: All tests passing, documentation complete

- [ ] Production setup {agent: devops} [→ SDD#infrastructure]
  - [ ] Provision infrastructure
  - [ ] Configure monitoring
  - [ ] Setup backups
  - [ ] Deploy to staging
  - [ ] **Validation**: Staging deployment successful

- [ ] Security review {agent: security-engineer} [→ BRD#security-requirements]
  - [ ] Penetration testing
  - [ ] Compliance check
  - [ ] Security scan
  - [ ] **Validation**: Security approval received

- [ ] Production deployment {agent: devops} [→ SDD#deployment-strategy]
  - [ ] Deploy with rollback plan
  - [ ] Monitor metrics
  - [ ] Verify functionality
  - [ ] **Validation**: Production healthy, metrics normal

## Completion Criteria
- [ ] All tasks marked complete
- [ ] Test coverage meets targets
- [ ] Performance benchmarks achieved
- [ ] Security review passed
- [ ] Documentation approved
- [ ] Production deployment successful
- [ ] Stakeholder sign-off received

## Dynamic Task Addition

When adding new tasks during execution:
```markdown
- [ ] [New task description] {agent: specialist} [→ source#ref]
  - **Reason**: [Why this was added]
  - **Dependencies**: [What must complete first]
```

## Rollback Plan

If deployment fails:
1. {agent: devops} Execute rollback procedure
2. {agent: site-reliability-engineer} Diagnose failure  
3. {agent: project-manager} Update plan with fixes
4. Resume from appropriate phase