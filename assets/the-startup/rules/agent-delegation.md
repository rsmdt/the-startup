Rules for task decomposition and parallel execution.

1. Task Decomposition Principles:

    **Ask yourself before decomposing**:
    - What distinct activities make up this task?
    - Can these activities run independently?
    - Do they require different expertise?
    - Where are the natural boundaries?
    - What specific output format do I need from each agent?
    - Have I provided enough detail for unambiguous execution?

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

    **Ask yourself before launching parallel agents**:
    - Will these tasks block each other?
    - Do they share state or dependencies?
    - Can I validate each independently?
    - Am I providing exhaustive detail in FOCUS/EXCLUDE?
    - Would another orchestrator understand exactly what's needed?

    DEFAULT: Always execute in parallel unless tasks depend on each other.

    Parallel Checklist:
    ✅ Independent tasks (no shared state)
    ✅ Different expertise required
    ✅ Separate validation possible
    ✅ Won't block each other

    **Ask yourself before writing each prompt**:
    - Have I described the complete task in FOCUS?
    - Have I listed ALL things to avoid in EXCLUDE?
    - Is the CONTEXT sufficient for independent execution?
    - Will the OUTPUT format prevent ambiguity?
    - Are SUCCESS criteria measurable and clear?
    - Would another orchestrator get identical results?
    
    The FOCUS/EXCLUDE Pattern (Required):
    ```
    FOCUS: [What to do - be comprehensive and specific]
    EXCLUDE: [What NOT to do - be equally comprehensive]
    CONTEXT: [All relevant background, constraints, and dependencies]
    SUCCESS: [Measurable completion criteria]
    ```
    
    Enhanced version (recommended for all agents):
    ```
    FOCUS: [What to do - provide complete task description with all details]
    EXCLUDE: [What NOT to do - list all boundaries and restrictions]
    CONTEXT: [Full background including prior work, current state, constraints]
    OUTPUT: [EXACT format/structure expected with examples if helpful]
           [If creating files: specify exact paths like docs/patterns/auth-pattern.md]
    SUCCESS: [All completion criteria that must be met]
    TERMINATION: [Explicit conditions for stopping]
    ```
    
    Example:
    ```python
    # Launch simultaneously
    Task(subagent_type="api-specialist", prompt="FOCUS: Build user API...")
    Task(subagent_type="database-specialist", prompt="FOCUS: Design schema...")
    Task(subagent_type="security-specialist", prompt="FOCUS: Review auth...")
    ```
    
    Result Aggregation:
    - Display each agent response verbatim
    - Synthesize findings after all responses
    - Identify conflicts between results
    - Create unified next steps



3. File Creation Coordination:

    **Ask yourself when agents will create files**:
    - Have I specified exact file paths for each agent?
    - Are all file paths unique (no collisions)?
    - Do the paths follow project conventions?
    - Will parallel agents overwrite each other?
    
    When multiple agents create documentation:
    - Specify exact file paths in OUTPUT section
    - Example: "Create pattern at docs/patterns/caching-strategy.md"
    - Never let multiple agents write to same file path
    - Use descriptive names to prevent accidental overlaps

4. Validation & Scope Control:

    **Ask yourself when reviewing agent responses**:
    - Did the agent stay within FOCUS boundaries?
    - Is this a security/quality improvement (auto-accept)?
    - Does this need user review (new dependencies)?
    - Is this scope creep (auto-reject)?

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
    - Missing required OUTPUT format
    - "While I'm here" additions
    - Unrequested improvements
    
    **Ask yourself when agent response seems off**:
    - Did I provide ambiguous instructions?
    - Should I have been more explicit in EXCLUDE?
    - Is this actually valuable despite being out of scope?
    - Will stricter FOCUS help or just waste time?
    
    When agents drift:
    ```
    ⚠️ Scope Alert: Agent added [unexpected feature]
    Options:
    1. Accept and update requirements
    2. Reject and retry with stricter FOCUS
    3. Cherry-pick useful parts
    ```

5. Failure Recovery:

    **Ask yourself when an agent fails**:
    - Was my FOCUS/EXCLUDE clear enough?
    - Should I try a different agent?
    - Can I break this into smaller tasks?
    - Should other parallel agents continue?

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

    **Ask yourself before retrying**:
    - What specifically caused the failure?
    - Would a different agent be better suited?
    - Should this be multiple smaller tasks?
    - Is partial success acceptable here?
    
    Recovery Tactics:
    - Refine FOCUS/EXCLUDE and retry
    - Break into smaller tasks
    - Try alternative agents
    - Ship partial results with gaps noted

    Key Principle: One agent fails? Others continue. Collect successes before retrying failures.
