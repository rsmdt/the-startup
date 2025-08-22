Rules for task decomposition and parallel execution.

1. Task Decomposition Principles:

    Before delegating to specialist agents, decompose complex work.

    When to Decompose:
    - Multiple distinct expertise areas needed
    - Independent components that can be validated separately  
    - Natural boundaries between system layers
    - Different stakeholder perspectives required
     
    How to Decompose:
    1. Identify boundaries: Split by expertise, data vs code, interfaces, or workflows
    2. Ensure independence: Each task should have clear inputs/outputs
    3. Avoid duplication: Identify shared prerequisites once
    4. Assign ownership: One sub-agent owns each task - no overlap
    5. Check coupling: If heavy cross-talk needed, merge or run sequentially

    Decomposition Example:
    ```
    Task: "Add user authentication"
    Decomposed into:
    - Security analysis {agent: the-security-engineer}
    - Database schema design {agent: the-data-engineer}  
    - API endpoint implementation {agent: the-developer}
    - UI/UX design {agent: the-ux-designer}
    ```

2. Parallel task or agent execution patterns:

    ALWAYS execute in parallel when possible - this is startup speed.

    Parallel Execution Criteria:
    - Are all tasks are independent (no shared state modifications)?
    - Do I need different expertise domains?
    - Is separate validation possible?
    - Will failure of one block others?

    For each specialist agent, provide the relevant context:
    - Specific task and constraints
    - What NOT to do (prevents scope creep)
    - Relevant requirements and dependencies
    - Clear criteria for completion

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

Remember: Fast execution with preserved expertise - that's how startups ship quality at speed.
