Rules for task decomposition and parallel execution.

1. Task Decomposition Principles:

    Before delegating to specialist agents, decompose complex work by ACTIVITIES.

    When to Decompose:
    - Multiple distinct activities needed
    - Independent components that can be validated separately  
    - Natural boundaries between system layers
    - Different stakeholder perspectives required
     
    How to Decompose:
    1. Identify activities: What needs to be DONE (not who does it)
    2. Ensure independence: Each task should have clear inputs/outputs
    3. Avoid duplication: Identify shared prerequisites once
    4. Express as capabilities: Use activity descriptions, not agent names
    5. Check coupling: If heavy cross-talk needed, merge or run sequentially

    Decomposition Example:
    ```
    Task: "Add user authentication"
    Decomposed into:
    - Analyze security requirements for auth system
    - Design database schema for users and sessions
    - Design API endpoints for authentication
    - Create login/register UI components
    ```

    The system automatically matches activities to specialized agents.

2. Parallel task or agent execution patterns:

    ALWAYS execute in parallel when possible - this is startup speed.

    Parallel Execution Criteria:
    - Are all tasks independent (no shared state modifications)?
    - Do different activities require different expertise?
    - Is separate validation possible?
    - Will failure of one block others?

    Context Isolation Pattern (FOCUS/EXCLUDE):
    ```
    FOCUS: [Specific activity and constraints]
    EXCLUDE: [What NOT to do - prevents scope creep]
    CONTEXT: [Only relevant requirements and dependencies]
    SUCCESS: [Clear criteria for completion]
    ```

    Parallel Execution Example:
    ```python
    # Launch multiple specialized agents simultaneously
    Task(subagent_type="specialized-agent-1", prompt="FOCUS: API design...")
    Task(subagent_type="specialized-agent-2", prompt="FOCUS: Database schema...")
    Task(subagent_type="specialized-agent-3", prompt="FOCUS: Security review...")
    ```

    Result Aggregation:
    - Display each agent response verbatim in delimiters
    - Synthesize findings AFTER showing all responses
    - Identify conflicts or dependencies between results
    - Create unified next steps from multiple perspectives

3. Validation & drift detection

    Auto-Accept (continue without review):
    - Security vulnerability fixes
    - Error handling improvements
    - Input validation additions
    - Performance optimizations under 10 lines
    - Documentation updates
    
    Requires Review (need user confirmation):
    - New external dependencies
    - Database schema modifications
    - Public API changes
    - Architectural pattern changes
    - Configuration updates
    
    Auto-Reject (scope creep - block immediately):
    - Features not in requirements
    - Breaking changes without migration path
    - Untested code modifications
    - Scope expansions beyond FOCUS directive
    
    Example - Handling Drift when specialist agent exceeds scope:
    ```
    ⚠️ Scope Alert: {agent} included {unexpected feature}
    
    Options:
    1. Accept and expand scope (update requirements)
    2. Reject and re-run with stricter boundaries
    3. Cherry-pick useful parts, discard rest
    ```

4. Failure Handling & Graceful Degradation

    When Parallel Execution Fails:
    - If one agent fails, others continue independently
    - Collect all successful results before retrying failures
    - Fall back to sequential execution if coordination issues arise
    - Escalate to user if critical path blocked

    Fallback Patterns:
    ```
    Try: Parallel execution of 4 specialized agents
    If coordination issues: Sequential execution with context passing
    If specialist unavailable: Route to broader domain agent
    If all else fails: Escalate to user with clear options
    ```

    Recovery Strategies:
    - Retry with refined FOCUS/EXCLUDE directives
    - Break down task into smaller activities
    - Use alternative specialized agents
    - Provide partial results with clear gaps identified

Remember: Fast execution with preserved expertise - that's how startups ship quality at speed.
