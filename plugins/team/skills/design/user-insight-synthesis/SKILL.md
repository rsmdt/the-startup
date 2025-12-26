---
name: user-insight-synthesis
description: Interview techniques, persona creation, journey mapping, and usability testing patterns. Use when planning research, conducting user interviews, creating personas, mapping user journeys, or designing usability tests. Essential for user-research, requirements-analysis, and interaction-architecture agents.
---

# User Research Methods

A domain-specific skill providing structured methodologies for understanding users through research. This skill covers the complete research lifecycle from planning through synthesis and reporting.

## When to Use

- Planning user research studies for new or existing products
- Conducting user interviews to understand needs and behaviors
- Creating behavioral personas based on research data
- Mapping user journeys to identify pain points and opportunities
- Designing and running usability tests
- Synthesizing research findings into actionable insights

## Research Method Selection

Choose research methods based on what you need to learn and where you are in the product lifecycle.

### Method Selection Matrix

| Question Type | Early Discovery | Design Validation | Post-Launch |
|--------------|-----------------|-------------------|-------------|
| What do users need? | Contextual inquiry, Diary studies | Concept testing | Support ticket analysis |
| How do users behave? | Field observation, Shadowing | Prototype testing | Analytics, Session recordings |
| What do users think? | Depth interviews | Preference testing | Surveys, NPS |
| Can users complete tasks? | Card sorting | Usability testing | A/B testing |

### Generative vs Evaluative Research

**Generative Research** - Use when exploring problem spaces:
- Contextual inquiry (observe users in their environment)
- Diary studies (longitudinal behavior patterns)
- Participatory design workshops (co-create with users)
- Jobs-to-be-done interviews (understand underlying motivations)

**Evaluative Research** - Use when validating solutions:
- Usability testing (can users complete tasks?)
- A/B testing (which variant performs better?)
- Preference testing (which option do users prefer?)
- Accessibility audits (does it work for everyone?)

### Sample Size Guidelines

| Method | Minimum Sample | Recommended | Diminishing Returns |
|--------|---------------|-------------|---------------------|
| Depth interviews | 5 | 8-12 | 15+ |
| Usability testing | 5 | 5-8 | 10+ |
| Card sorting | 15 | 30 | 50+ |
| Surveys | 100 | 300-500 | Depends on segments |
| A/B tests | Statistical power calculation required | - | - |

## User Interview Techniques

### Interview Structure

**1. Opening (5 minutes)**
- Build rapport with neutral topics
- Explain the purpose without biasing
- Confirm recording consent
- Set expectations for the session

**2. Context Gathering (10 minutes)**
- Understand their role and background
- Map their typical day or workflow
- Identify tools and touchpoints they use

**3. Core Exploration (25-35 minutes)**
- Use open-ended questions
- Follow the participant's lead
- Probe deeper on interesting topics
- Ask for specific examples and stories

**4. Closing (5 minutes)**
- Ask if anything was missed
- Request permission for follow-up
- Thank them for their time

### Question Types

**Opening Questions** - Establish context:
- "Walk me through a typical day when you..."
- "Tell me about the last time you..."
- "How did you first start...?"

**Probing Questions** - Go deeper:
- "What do you mean by...?"
- "Can you give me a specific example?"
- "What happened next?"
- "How did that make you feel?"

**Clarifying Questions** - Ensure understanding:
- "So if I understand correctly..."
- "You mentioned X, can you tell me more?"
- "When you say X, do you mean...?"

**Projective Questions** - Surface hidden needs:
- "If you had a magic wand, what would you change?"
- "What would your ideal experience look like?"
- "What would have to be true for you to switch?"

### Common Interview Pitfalls

| Pitfall | Problem | Solution |
|---------|---------|----------|
| Leading questions | Biases responses | Ask neutral, open questions |
| Asking about future behavior | People predict poorly | Ask about past behavior instead |
| Accepting vague answers | Loses detail | Probe for specific examples |
| Talking too much | Reduces user input | Embrace silence, let them think |
| Solving during interview | Shifts to validation mode | Save solutions for later |

### Observing Behavior vs Statements

Users often say one thing but do another. Watch for discrepancies:

- **Workarounds**: Creative solutions reveal unmet needs
- **Hesitation**: Confusion or friction points
- **Skipped steps**: What they consider unnecessary
- **Emotional reactions**: Frustration, delight, confusion
- **Tool switching**: Integration pain points

## Persona Creation Framework

### Behavioral Personas vs Demographic Personas

**Avoid demographic personas** that describe who users are (age, income, job title). These often become stereotypes that don't predict behavior.

**Create behavioral personas** that describe what users do, why they do it, and what barriers they face. These drive design decisions.

### Persona Components

```
[Persona Name]

BEHAVIORAL ARCHETYPE
One sentence describing their core behavior pattern

GOALS
- Primary goal (what they're ultimately trying to achieve)
- Secondary goals (supporting objectives)
- Emotional goals (how they want to feel)

BEHAVIORS
- Key behavior 1 (observed pattern with frequency)
- Key behavior 2 (observed pattern with context)
- Key behavior 3 (observed pattern with trigger)

PAIN POINTS
- Frustration 1 (specific problem with impact)
- Frustration 2 (specific problem with workaround)
- Frustration 3 (specific problem with frequency)

DECISION FACTORS
- What influences their choices
- What trade-offs they make
- What they prioritize

CONTEXT
- When they engage with the product
- Where they use it
- What else competes for attention

QUOTE
"Verbatim quote from research that captures their perspective"
```

### Creating Personas from Research

**Step 1: Identify Behavior Patterns**
- Review all interview notes and observations
- Tag recurring behaviors, goals, and pain points
- Look for clusters of similar behavior

**Step 2: Define Behavioral Variables**
- List the key dimensions that differentiate users
- Place participants along each dimension
- Identify clusters that represent archetypes

**Step 3: Build Persona Profiles**
- Write narrative based on research data only
- Include direct quotes from participants
- Validate that persona represents multiple participants

**Step 4: Prioritize Personas**
- Identify primary persona (design for first)
- Secondary personas (accommodate)
- Edge cases (consider but don't optimize for)

### Persona Anti-Patterns

| Anti-Pattern | Why It Fails | Better Approach |
|--------------|--------------|-----------------|
| Fictional details | Creates false confidence | Use only observed data |
| Photos of real people | Triggers stereotypes | Use illustrations or initials |
| Too many personas | Dilutes focus | Limit to 3-5 maximum |
| Demographic focus | Doesn't predict behavior | Focus on goals and behaviors |
| No pain points | Misses design opportunities | Ground in observed frustrations |

## Journey Mapping Methodology

### Journey Map Components

```
JOURNEY MAP: [User Goal]

PERSONA: [Which persona this represents]
SCENARIO: [Specific context and trigger]

| Stage | [Stage 1] | [Stage 2] | [Stage 3] | [Stage 4] |
|-------|-----------|-----------|-----------|-----------|
| Actions | What user does | ... | ... | ... |
| Touchpoints | Channels/interfaces | ... | ... | ... |
| Thoughts | What they're thinking | ... | ... | ... |
| Emotions | How they feel (scale) | ... | ... | ... |
| Pain Points | Frustrations | ... | ... | ... |
| Opportunities | Design possibilities | ... | ... | ... |
```

### Journey Mapping Process

**1. Define Scope**
- Which persona is this for?
- What goal or scenario are we mapping?
- Where does the journey start and end?
- What level of detail do we need?

**2. Map the Current State**
- List all stages from trigger to completion
- Document what users do at each stage
- Identify all touchpoints (channels, systems, people)
- Note what users think and feel
- Mark pain points and moments of delight

**3. Validate with Data**
- Cross-reference with analytics data
- Validate with additional user interviews
- Check assumptions against support data
- Ensure map represents typical experience

**4. Identify Opportunities**
- Where are the highest-impact pain points?
- Where do users drop off or get stuck?
- What moments could be transformed?
- Where can we exceed expectations?

### Journey Map Types

**Current State Maps** - Document how things work today
- Based on research observations
- Reveals improvement opportunities
- Aligns stakeholders on reality

**Future State Maps** - Envision improved experience
- Based on current state insights
- Shows target experience
- Guides design decisions

**Service Blueprints** - Include backend processes
- Shows frontstage and backstage
- Reveals operational dependencies
- Identifies system requirements

### Emotional Curve Mapping

Track emotional state across the journey using a simple scale:

```
Very Positive  +2  ----*----
Positive       +1  ----*---------*----
Neutral         0  *---------*----
Negative       -1  ----*----
Very Negative  -2  ----*----
                   |Stage1|Stage2|Stage3|Stage4|
```

Key moments to identify:
- **Peak moments**: Highest positive emotion (protect and amplify)
- **Valley moments**: Lowest emotion (priority for improvement)
- **Transition points**: Where emotion shifts (critical touchpoints)
- **Ending moments**: Final impression (strong impact on memory)

## Usability Testing Patterns

### Test Types

**Moderated Testing** - Facilitator guides participant
- Best for: Complex tasks, early prototypes, need for probing
- Pros: Rich qualitative data, can adapt on the fly
- Cons: Time-intensive, facilitator can bias

**Unmoderated Testing** - Participant works independently
- Best for: Simple tasks, large samples, geographic spread
- Pros: Scalable, no facilitator bias, natural behavior
- Cons: No probing, participants may give up

**Guerrilla Testing** - Quick tests with available people
- Best for: Early validation, simple concepts, tight timelines
- Pros: Fast, cheap, good for iteration
- Cons: May not match target users, limited depth

### Test Protocol Structure

**1. Pre-Test Setup**
- Confirm participant matches screener
- Prepare test environment (prototype, recording)
- Review tasks and questions
- Test the test (pilot run)

**2. Introduction (5 minutes)**
- Explain the purpose (testing the design, not them)
- Describe think-aloud protocol
- Confirm recording consent
- Encourage honest feedback

**3. Background Questions (5 minutes)**
- Relevant experience with similar products
- Current tools and workflows
- Expectations for this type of product

**4. Task Scenarios (30-40 minutes)**
- Present tasks one at a time
- Use realistic scenarios, not instructions
- Observe without helping
- Probe after task completion

**5. Post-Test Questions (10 minutes)**
- Overall impressions
- Comparison to expectations
- Suggestions for improvement
- Follow-up on observed issues

### Writing Effective Task Scenarios

**Bad task**: "Click on Settings and change your notification preferences"
- Reveals the solution
- Uses UI terminology
- No realistic context

**Good task**: "You're getting too many email notifications. How would you reduce them?"
- Goal-oriented
- User's language
- Realistic motivation

### Task Scenario Template

```
SCENARIO: [Context that makes the task realistic]
GOAL: [What the user is trying to accomplish]
SUCCESS CRITERIA: [How you'll know they succeeded]
```

### Metrics to Capture

**Effectiveness Metrics**
- Task success rate (completed / attempted)
- Error rate (errors / task)
- Recovery rate (recovered from errors / total errors)

**Efficiency Metrics**
- Time on task (seconds to completion)
- Number of steps (compared to optimal path)
- Help requests (times asked for assistance)

**Satisfaction Metrics**
- Post-task rating (1-7 scale per task)
- System Usability Scale (SUS) score
- Net Promoter Score (NPS)
- Qualitative feedback themes

### Severity Rating Scale

| Rating | Name | Definition | Action |
|--------|------|------------|--------|
| 1 | Cosmetic | Noticed but no impact | Fix if time permits |
| 2 | Minor | Slight difficulty, recovered easily | Fix in next release |
| 3 | Major | Significant difficulty, delayed success | Fix before release |
| 4 | Critical | Prevented task completion | Must fix immediately |

## Synthesis and Reporting

### Affinity Mapping Process

**1. Gather Raw Data**
- Write each observation on a separate note
- Include quotes, behaviors, and pain points
- One insight per note

**2. Cluster Related Notes**
- Group similar observations together
- Don't force categories, let them emerge
- Move notes as patterns become clear

**3. Name the Clusters**
- Create descriptive labels for each group
- Labels should capture the theme
- Higher-level groups may emerge

**4. Identify Patterns**
- What themes appear across multiple participants?
- What surprises challenge assumptions?
- What opportunities become clear?

### Research Report Structure

```
RESEARCH REPORT: [Study Name]

EXECUTIVE SUMMARY
- Research objective
- Methods used
- Key findings (3-5 bullets)
- Recommended actions

METHODOLOGY
- Research questions
- Participant criteria and recruitment
- Methods and sample size
- Limitations

KEY FINDINGS
Finding 1: [Title]
- Evidence: What we observed
- Impact: Why it matters
- Quote: "Supporting verbatim"

Finding 2: [Title]
...

RECOMMENDATIONS
- Priority 1: [Action] - Addresses [finding]
- Priority 2: [Action] - Addresses [finding]
- Priority 3: [Action] - Addresses [finding]

NEXT STEPS
- Immediate actions
- Further research needed
- Stakeholder follow-up

APPENDIX
- Participant details (anonymized)
- Full data sets
- Methodology details
```

### Presenting Findings

**Lead with insights, not methodology**
- Start with what you learned, not how you learned it
- Save methodology details for appendix

**Use participant voices**
- Include direct quotes that bring findings to life
- Video clips are more powerful than text

**Connect to business outcomes**
- Tie findings to metrics stakeholders care about
- Quantify impact where possible

**Provide clear recommendations**
- Don't just report problems, suggest solutions
- Prioritize by impact and feasibility

## Best Practices

- Recruit participants who actually face the problem you're solving
- Focus on understanding behavior, not validating your ideas
- Look for patterns across participants, not individual opinions
- Always observe what users do, not just what they say
- Use multiple methods to triangulate findings
- Create lightweight deliverables that teams will actually use
- Involve stakeholders in research to build empathy
- Measure research impact through design decisions influenced

## References

- `templates/research-plan-template.md` - Planning template for research studies
