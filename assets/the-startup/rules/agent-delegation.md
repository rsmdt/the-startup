Rules for task decomposition and parallel execution.

1. Task Decomposition Principles:

    Decompose complex work by ACTIVITIES (what needs doing), not roles.
    
    ✅ DO: "Build API endpoints", "Create UI components", "Add security layer"  
    ❌ DON'T: "Backend engineer do X", "Frontend do Y"

    When to Decompose:
    - Multiple distinct activities needed
    - Independent components that can be validated separately  
    - Natural boundaries between system layers
    - Different stakeholder perspectives required
     
    Example:
    ```
    Task: "Add user authentication"
    → Analyze security requirements
    → Design database schema  
    → Create API endpoints
    → Build login/register UI
    ```

    The system automatically matches activities to specialized agents.

2. Parallel Execution Patterns:

    DEFAULT: Always execute in parallel unless tasks depend on each other.

    Parallel Checklist:
    ✅ Independent tasks (no shared state)
    ✅ Different expertise required
    ✅ Separate validation possible
    ✅ Won't block each other

    The FOCUS/EXCLUDE Pattern (Required):
    ```
    FOCUS: [What to do - 2-3 sentences max]
    EXCLUDE: [What NOT to do - prevents scope creep]
    CONTEXT: [Only relevant files/requirements]
    SUCCESS: [Measurable completion criteria]
    ```

    Example:
    ```python
    # Launch simultaneously
    Task(subagent_type="the-backend-engineer", prompt="FOCUS: Build user API...")
    Task(subagent_type="the-software-engineer-database-design", prompt="FOCUS: Design schema...")
    Task(subagent_type="the-security-engineer-authentication-systems", prompt="FOCUS: Review auth...")
    ```

    Result Aggregation:
    - Display each agent response verbatim
    - Synthesize findings after all responses
    - Identify conflicts between results
    - Create unified next steps

3. Validation & Scope Control:

    🟢 Auto-Accept (ship it):
    - Security vulnerability fixes
    - Error handling improvements
    - Input validation additions
    - Performance optimizations under 10 lines
    - Documentation updates
    
    🟡 Requires Review (user confirms):
    - New external dependencies
    - Database schema modifications
    - Public API changes
    - Architectural pattern changes
    - Configuration updates
    
    🔴 Auto-Reject (scope creep):
    - Features not in requirements
    - Breaking changes without migration
    - Untested code modifications
    - Scope expansions beyond FOCUS directive
    
    When agents drift:
    ```
    ⚠️ Scope Alert: Agent added [unexpected feature]
    Options:
    1. Accept and update requirements
    2. Reject and retry with stricter FOCUS
    3. Cherry-pick useful parts
    ```

4. Failure Recovery:

    Fallback Chain:
    ```
    1. Parallel specialists (default)
       ↓ if coordination issues
    2. Sequential with context passing
       ↓ if specialist fails
    3. Broader domain agent
       ↓ if still stuck
    4. DIY (do it yourself)
       ↓ if blocked
    5. Escalate to user with options
    ```

    Recovery Tactics:
    - Refine FOCUS/EXCLUDE and retry
    - Break into smaller tasks
    - Try alternative agents
    - Ship partial results with gaps noted

    Key Principle: One agent fails? Others continue. Collect successes before retrying failures.
