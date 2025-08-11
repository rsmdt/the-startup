# Implementation Plan

*[INSTRUCTION: Create a detailed implementation checklist based on the feature requirements, architecture, and discovered patterns. Organize tasks into logical phases that can be executed sequentially, with clear dependencies and validation points. This entire instruction block should not appear in the final PRD.]*

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
