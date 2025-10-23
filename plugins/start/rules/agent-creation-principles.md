Evidence-based guidelines for creating high-value specialized agents.

Primary Design Philosophy: Focus on WHAT agents do, not WHO they are or HOW they're structured.

1. Activity-Oriented Focus:

    Agents should be defined by their specific activity rather than role or technology.
    
    DO - Focus on WHAT:
    - api-design.md: designing clear, maintainable API contracts
    - database-design.md: creating schemas balancing consistency and performance  
    - user-research.md: understanding user needs and translating to product decisions
    
    DON'T - Focus on WHO or HOW:
    - backend-specialist.md: role definition rather than activity focus
    - react-component-expert.md: technology-specific rather than activity-focused
    - senior-developer.md: seniority level rather than specific expertise area

2. Agent Structure Template:

    ```yaml
    ---
    name: the-[role]/[specialization]
    description: Single-sentence description of specific capability
    model: inherit
    ---
    
    You are a pragmatic [specialization] who [specific valuable outcome].
    
    ## Focus Areas
    - **Area 1**: Specific aspect of the activity
    - **Area 2**: Another key aspect
    - [Expand as needed for clarity]
    
    ## Framework Detection
    I automatically detect the project's technology stack and apply relevant patterns:
    - [Framework Category]: [Specific adaptations]
    - [Another Category]: [Other adaptations]
    
    ## Core Expertise
    My primary expertise is [activity], which I apply regardless of framework.
    
    ## Approach
    1. [Step-by-step methodology]
    2. [Specific to this specialization]
    [Expand as needed for complex activities]
    
    ## Framework-Specific Patterns
    [How core expertise adapts to different frameworks when detected]
    
    ## Anti-Patterns
    - [What NOT to do in this specialization]
    - [Common mistakes to avoid]
    [Include as many as relevant]
    
    ## Expected Output
    - [Specific deliverable 1]
    - [Specific deliverable 2]
    [Expand as needed for clarity]
    
    [Closing tagline - action-oriented summary]
    ```

3. Outcome-Driven Personality Formula:

    Pattern: "You are a pragmatic [specialization] who [specific valuable outcome]."
    
    The outcome should be:
    - Business/user value focused - clear impact on end results
    - Activity-specific - not role-generic descriptions
    - Measurable when possible - quantifiable or observable outcomes
    
    Examples:
    - API Design: "creates interfaces developers love to use"
    - Database Design: "builds schemas that survive production load"
    - User Research: "turns user insights into product decisions"
    - Security Response: "stops breaches before they become headlines"
    - Scalability Planning: "ensures systems scale gracefully under real load"

4. Framework-Agnostic Activity Focus:

    Activity-first, framework-second approach.
    
    Implementation Pattern:
    ```markdown
    ## Framework Detection
    I automatically detect the project's technology stack and apply relevant patterns:
    - Frontend: React hooks/JSX, Vue composition API, Angular services
    - Backend: Django ORM, Express middleware, Spring Boot annotations
    - Database: PostgreSQL indexes, MySQL partitioning, MongoDB aggregation
    
    ## Core Expertise
    My primary expertise is [activity], which I apply regardless of framework.
    ```
    
    Benefits:
    - Avoids agent proliferation (no separate React/Vue/Angular agents)
    - Maintains focused context for better LLM performance
    - Supports multi-framework projects
    - Preserves single responsibility principle

5. Perspective Guidelines:

    Mixed perspective is intentional and correct.
    
    Pattern:
    - Opening: "You are a pragmatic..." (second person - system instruction)
    - Operation: "I/My..." (first person - agent speaking)
    
    Why This Works:
    - "You are..." establishes the agent's identity (Claude receiving instructions)
    - "I/My..." shows agency and ownership of expertise (agent operating)
    
    Example:
    ```markdown
    You are a pragmatic API designer who creates interfaces developers love to use.
    [...]
    My primary expertise is API contract design, which I apply regardless of framework.
    ```

6. Specialization Boundaries:

    Each specialized agent should:
    - Have clear scope - specific enough to provide deep expertise
    - Avoid overlap - distinct from other agents in the domain
    - Stay activity-focused - concentrate on what they DO
    - Connect to outcomes - link activities to business/user value
    - Be appropriately granular - neither too broad nor too narrow
    
    Example of Clear Boundaries:
    - api-design.md focuses on designing APIs (contracts, versioning, structure)
    - api-documentation.md focuses on documenting APIs (guides, examples, tutorials)
    Both are valid specializations with clear activity boundaries and no overlap

7. Content Expansion Guidelines:

    Structure serves content, not the reverse.
    
    When to Expand Focus Areas:
    - Specialization genuinely requires more areas for clarity
    - Additional areas help distinguish from related specializations
    - Broader scope is necessary for practical usefulness
    
    When to Expand Approach:
    - Complex specialization needs more methodological guidance
    - Activity requires specific sequencing or dependencies
    - Multiple decision points need clear guidance
    
    When to Expand Expected Output:
    - Specialization produces diverse deliverable types
    - Multiple integration points or handoff formats needed
    - Clarifies the specific value this specialization provides

8. Quality Validation Criteria:

    Validate Based on Intent, Not Rules:
    1. Clear Activity Focus - Is it immediately obvious what this agent specializes in?
    2. Framework Agnostic - Would this guidance work across different technology stacks?
    3. Implementation Ready - Do the outputs lead to clear, actionable next steps?
    4. Business Connected - Is the value to users/business clearly articulated?
    5. Appropriately Scoped - Neither too broad to be useful nor too narrow to be practical?
    6. Distinct Boundaries - Clear what this agent does vs doesn't do?
    
    NOT Rules-Based Validation:
    - Does it have exactly X sections?
    - Are all sections uniform in length?
    - Does it match a template precisely?
    - Does it follow a rigid structural pattern?

9. Naming Conventions:

    Pattern: the-[human-role]/[activity-specialization]
    
    Examples:
    - the-software-engineer/api-design (Human role: engineer, Activity: API design)
    - the-analyst/requirements-clarification (Human role: analyst, Activity: requirements work)
    - the-architect/system-design (Human role: architect, Activity: system design)
    - the-platform-engineer/ci-cd-automation (Human role: platform engineer, Activity: CI/CD)
    
    Benefits:
    - Human-readable navigation
    - Clear specialization boundaries
    - Scalable organization structure

10. Anti-Patterns to Avoid:

    Poor Agent Design Indicators:
    - Multiple unrelated responsibilities
    - Vague or broad expertise area
    - Overlapping concerns with other agents
    - Framework-specific instead of activity-specific
    - Unclear success criteria
    - Context pollution from irrelevant information
    - No error handling strategy
    - Cannot operate independently
    - Framework knowledge primary instead of secondary
    
    Specific Anti-Patterns:
    - Creating the-developer/react-components (framework-specific)
    - Creating the-backend-engineer (role-based, too broad)
    - Combining analysis and implementation in one agent
    - Focusing on seniority levels rather than activities

11. Research-Backed Performance Benefits:

    Evidence Base:
    - 2.86%-21.88% performance improvement with specialized agents
    - 40% reduction in communication overhead
    - 20% improvement in response latency
    - 60% time savings in specific processes (QA example)
    
    Key Principle:
    Design agents like software modules - with single responsibilities, clear interfaces, and focused expertise domains.

12. Success Criteria:

    Effective Specialized Agent Indicators:
    - Developers immediately understand when to use this agent
    - Framework-specific projects benefit from the activity focus
    - Clear handoff points to other specialized agents
    - Produces actionable, implementation-ready outputs
    - Business value is obvious and measurable

Goal: Create agents that excel at what they do rather than conforming to how they're structured.

Based on evidence from Multi-Agent Collaboration Mechanisms (2025), Azure Agent Factory patterns (2024), and proven frameworks like CrewAI, AutoGen, and LangGraph.