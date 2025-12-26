---
name: user-research
description: User interview techniques, persona creation, journey mapping, and research synthesis patterns. Use when planning research studies, conducting interviews, creating personas, or translating research findings into actionable design recommendations.
---

# User Research Methodology

Systematic approaches for understanding user needs, behaviors, and motivations to inform product decisions.

## When to Activate

- Planning user research studies
- Conducting user interviews
- Creating personas and journey maps
- Synthesizing research findings
- Translating insights into design recommendations
- Validating product concepts with users

## Research Methods

### Method Selection Guide

| Method | Best For | Sample Size | Time Investment |
|--------|----------|-------------|-----------------|
| **User Interviews** | Deep understanding, "why" | 5-12 users | 2-3 weeks |
| **Contextual Inquiry** | Understanding environment | 3-6 users | 1-2 weeks |
| **Usability Testing** | Interface validation | 5 users | 1 week |
| **Surveys** | Quantitative validation | 100+ users | 1-2 weeks |
| **Card Sorting** | Information architecture | 15-30 users | 1 week |
| **Diary Studies** | Longitudinal behavior | 10-15 users | 2-4 weeks |

### User Interviews

One-on-one conversations to understand user perspectives.

#### Interview Structure (60 min)

```
1. INTRODUCTION (5 min)
   - Thank them for participating
   - Explain purpose (learning, not testing)
   - Request permission to record
   - Emphasize no right/wrong answers

2. WARM-UP (5 min)
   - Easy, open questions
   - Build rapport
   - "Tell me about your role..."

3. CONTEXT (10 min)
   - Current situation
   - Tools and processes
   - Goals and challenges
   - "Walk me through a typical day..."

4. DEEP DIVE (30 min)
   - Specific experiences
   - Pain points in detail
   - Workarounds and adaptations
   - "Tell me about a time when..."

5. EXPLORATION (5 min)
   - Reactions to concepts (if applicable)
   - Ideal scenarios
   - "If you could wave a magic wand..."

6. WRAP-UP (5 min)
   - Summary of key points
   - Anything else to add
   - Thank you and next steps
```

#### Question Techniques

| Technique | Purpose | Example |
|-----------|---------|---------|
| **Open-ended** | Encourage stories | "Tell me about..." |
| **Follow-up** | Dig deeper | "Can you say more about that?" |
| **Clarification** | Ensure understanding | "When you say X, what do you mean?" |
| **Contrast** | Explore differences | "How does that compare to...?" |
| **Projection** | Uncover desires | "What would ideal look like?" |
| **Silence** | Let them think | [Wait 5-10 seconds after answers] |

#### Questions to Avoid

| Avoid | Problem | Better |
|-------|---------|--------|
| "Do you like...?" | Yes/no answer | "How do you feel about...?" |
| "Would you use...?" | Hypothetical behavior ≠ real | "When did you last...?" |
| "Don't you think...?" | Leading | "What do you think about...?" |
| "What features...?" | Solution-focused | "What problems do you face?" |

### Contextual Inquiry

Observe users in their natural environment.

#### Protocol

```
PREPARATION:
- Define focus areas
- Prepare observation guide
- Get necessary permissions
- Test recording equipment

DURING OBSERVATION:
1. Arrive early, set up quietly
2. Start with brief introduction
3. Observe first, ask questions after
4. Note everything (actions, environment, emotions)
5. Use "teach me" framing

OBSERVATION GUIDE:
- What are they trying to accomplish?
- What tools are they using?
- What workarounds do they employ?
- What frustrates them?
- What's in their physical environment?
- Who do they interact with?

DEBRIEF:
- Review observations with participant
- Ask clarifying questions
- Confirm interpretations
```

#### Observation Notes Template

```
┌─────────────────────────────────────────────────────────────┐
│ Participant: [ID]     Date: [Date]     Location: [Where]    │
├─────────────────────────────────────────────────────────────┤
│ Task: [What they were doing]                                │
│ Time: [How long it took]                                    │
├─────────────────────────────────────────────────────────────┤
│ Actions Observed:                                           │
│ - [Step 1]                                                  │
│ - [Step 2]                                                  │
├─────────────────────────────────────────────────────────────┤
│ Tools Used:                                                 │
│ - [Tool 1]: [How used]                                      │
│ - [Tool 2]: [How used]                                      │
├─────────────────────────────────────────────────────────────┤
│ Pain Points:                                                │
│ - [Frustration observed]                                    │
├─────────────────────────────────────────────────────────────┤
│ Quotes:                                                     │
│ - "[Direct quote]"                                          │
├─────────────────────────────────────────────────────────────┤
│ Opportunities:                                              │
│ - [Potential improvement]                                   │
└─────────────────────────────────────────────────────────────┘
```

### Think-Aloud Protocol

Have users verbalize thoughts while performing tasks.

```
SETUP:
"I'd like you to complete some tasks while telling me what you're
thinking. There are no wrong answers - I'm testing the design,
not you. Please say out loud whatever you're looking at, thinking,
or feeling as you go through."

PROMPTS DURING SESSION:
- "What are you thinking right now?"
- "What do you expect to happen?"
- "What are you looking for?"
- "Why did you click there?"
- "How does this compare to what you expected?"

AVOID:
- Helping them complete tasks
- Confirming if they're right/wrong
- Explaining how things work
- Interrupting their flow too much
```

## Research Synthesis

### Affinity Mapping

Group observations to find patterns.

```
PROCESS:

1. CAPTURE (Individual)
   - Write one observation per sticky note
   - Use participant quotes
   - Include source identifier

2. CLUSTER (Group)
   - Spread all notes on wall/board
   - Group by similarity
   - Don't pre-define categories
   - Move notes until clusters emerge

3. NAME (Group)
   - Label each cluster
   - Labels should describe the theme
   - Not too broad, not too specific

4. PRIORITIZE
   - Which themes appear most frequently?
   - Which have highest impact?
   - Which are most actionable?

EXAMPLE CLUSTERS:
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ Trust       │ │ Efficiency  │ │ Support     │
│ Concerns    │ │ Pain Points │ │ Needs       │
├─────────────┤ ├─────────────┤ ├─────────────┤
│ "I don't    │ │ "Takes too  │ │ "Wish I     │
│ know if     │ │ many clicks"│ │ could ask   │
│ it's safe"  │ │             │ │ someone"    │
│             │ │ "Have to    │ │             │
│ "Where's    │ │ enter same  │ │ "Help docs  │
│ my data?"   │ │ info twice" │ │ are useless"│
└─────────────┘ └─────────────┘ └─────────────┘
```

### Insight Generation

Transform observations into actionable insights.

```
INSIGHT FORMULA:

[User group] needs [need] because [motivation/context],
but currently [pain point], which means [consequence].

EXAMPLE:

First-time users need clear guidance during setup because
they're unfamiliar with the product, but currently the
onboarding is overwhelming with too many options, which
means they abandon before experiencing value.

VALIDATION CHECKLIST:
- [ ] Based on evidence from multiple participants
- [ ] Identifies a real need (not a solution)
- [ ] Explains the underlying motivation
- [ ] Connects to business impact
- [ ] Is actionable
```

## Personas

### Persona Creation

Research-based archetypes representing user segments.

#### Persona Template

```
┌─────────────────────────────────────────────────────────────┐
│ [PHOTO PLACEHOLDER]                                          │
│                                                              │
│ NAME: [Fictional name]                                       │
│ TITLE: [Role/context]                                        │
│ ARCHETYPE: [2-3 word descriptor]                            │
├─────────────────────────────────────────────────────────────┤
│ QUOTE:                                                       │
│ "[Characteristic quote from research]"                       │
├─────────────────────────────────────────────────────────────┤
│ DEMOGRAPHICS:                                                │
│ Age: [Range]     Experience: [Level]                        │
│ Context: [Work/home environment]                            │
├─────────────────────────────────────────────────────────────┤
│ GOALS:                                                       │
│ - Primary: [Main objective]                                  │
│ - Secondary: [Supporting objective]                          │
├─────────────────────────────────────────────────────────────┤
│ PAIN POINTS:                                                 │
│ - [Frustration 1]                                           │
│ - [Frustration 2]                                           │
│ - [Frustration 3]                                           │
├─────────────────────────────────────────────────────────────┤
│ BEHAVIORS:                                                   │
│ - [How they approach problems]                              │
│ - [Tools/resources they use]                                │
│ - [Decision-making patterns]                                │
├─────────────────────────────────────────────────────────────┤
│ SCENARIO:                                                    │
│ [Brief story of them using your product]                    │
└─────────────────────────────────────────────────────────────┘
```

#### Persona Development Process

```
1. IDENTIFY VARIABLES
   - What attributes differentiate users?
   - Goals, behaviors, pain points, context

2. ANALYZE PATTERNS
   - Cluster research participants
   - Find natural groupings
   - Validate with quantitative data if available

3. CREATE PERSONAS
   - 3-5 personas is typical
   - Each represents a distinct segment
   - Include primary, secondary, negative persona

4. VALIDATE
   - Review with stakeholders
   - Check against additional research
   - Refine based on feedback

5. ACTIVATE
   - Share widely
   - Reference in design discussions
   - Update as you learn more
```

#### Persona Types

| Type | Purpose | When to Create |
|------|---------|----------------|
| **Primary** | Main design target | Always |
| **Secondary** | Important but not primary focus | When segments differ significantly |
| **Negative** | Who we're NOT designing for | When edge cases distract |
| **Proto-persona** | Hypothesis before research | Early exploration |

## Journey Mapping

### Journey Map Structure

```
┌─────────────────────────────────────────────────────────────┐
│ JOURNEY MAP: [User Type] - [Scenario]                       │
├─────────────────────────────────────────────────────────────┤
│ STAGE      │ Awareness │ Consider │ Purchase │ Use │ Renew │
├─────────────────────────────────────────────────────────────┤
│ ACTIONS    │           │          │          │     │       │
│ What they  │ • Sees ad │ • Visits │ • Selects│     │       │
│ do         │ • Asks    │   site   │   plan   │     │       │
│            │   friend  │ • Reads  │ • Enters │     │       │
│            │           │   reviews│   payment│     │       │
├─────────────────────────────────────────────────────────────┤
│ THOUGHTS   │           │          │          │     │       │
│ What they  │ "I need   │ "Is this │ "This    │     │       │
│ think      │ to solve  │ the right│ better be│     │       │
│            │ this      │ choice?" │ worth it"│     │       │
│            │ problem"  │          │          │     │       │
├─────────────────────────────────────────────────────────────┤
│ EMOTIONS   │    😊     │    😐    │    😟    │     │       │
│ How they   │ Hopeful   │ Confused │ Anxious  │     │       │
│ feel       │           │          │          │     │       │
├─────────────────────────────────────────────────────────────┤
│ TOUCH-     │ Social    │ Website  │ Checkout │     │       │
│ POINTS     │ media     │ Reviews  │ Email    │     │       │
├─────────────────────────────────────────────────────────────┤
│ PAIN       │           │ Too many │ Payment  │     │       │
│ POINTS     │           │ options  │ issues   │     │       │
├─────────────────────────────────────────────────────────────┤
│ OPPORT-    │           │ Compare  │ Guest    │     │       │
│ UNITIES    │           │ feature  │ checkout │     │       │
└─────────────────────────────────────────────────────────────┘
```

### Journey Mapping Process

```
1. DEFINE SCOPE
   - Which persona?
   - Which scenario?
   - Start and end points?

2. GATHER DATA
   - Interview transcripts
   - Analytics data
   - Support tickets
   - Observation notes

3. MAP THE STAGES
   - What are the major phases?
   - What triggers transitions?

4. FILL IN LAYERS
   - Actions at each stage
   - Thoughts and questions
   - Emotional state
   - Touchpoints

5. IDENTIFY OPPORTUNITIES
   - Where are the pain points?
   - Where can we improve?
   - What's the priority?

6. VALIDATE & SHARE
   - Review with stakeholders
   - Share findings
   - Define action items
```

## Research Planning

### Research Plan Template

```markdown
# Research Plan: [Study Name]

## Objectives
- Primary: [Main question to answer]
- Secondary: [Additional questions]

## Participants
- Target: [User segment]
- Sample size: [Number]
- Recruitment: [How to find them]
- Screener criteria: [Inclusion/exclusion]

## Methodology
- Method: [Interview/observation/testing]
- Duration: [Session length]
- Location: [Remote/in-person]
- Facilitator: [Who]

## Discussion Guide
- [Link to guide]

## Timeline
| Phase | Dates |
|-------|-------|
| Recruitment | [Dates] |
| Sessions | [Dates] |
| Analysis | [Dates] |
| Reporting | [Date] |

## Deliverables
- [ ] Raw notes
- [ ] Synthesis document
- [ ] Presentation
- [ ] Recommendations
```

## Reporting Research

### Research Report Structure

```markdown
# Research Findings: [Study Name]

## Executive Summary
[1-paragraph overview for stakeholders who won't read details]

## Background
- Objectives
- Methodology
- Participants (demographics, no PII)

## Key Findings

### Finding 1: [Headline]
**Evidence**: [3+ supporting data points]
**Impact**: [Why this matters]
**Recommendation**: [What to do]

### Finding 2: [Headline]
...

## Detailed Observations
[Supporting details, quotes, examples]

## Recommendations Summary
| Priority | Finding | Recommendation | Effort |
|----------|---------|----------------|--------|
| 1 | [Finding] | [Action] | [Est.] |

## Appendix
- Screener
- Discussion guide
- Participant list (anonymized)
```

## Anti-Patterns

| Anti-Pattern | Problem | Solution |
|--------------|---------|----------|
| **Confirmation Bias** | Seeking data that confirms beliefs | Ask open questions, look for disconfirming evidence |
| **Leading Questions** | Influencing responses | Review questions for bias |
| **Recency Effect** | Overweighting last interview | Synthesize across all participants |
| **Sample Bias** | Wrong participants | Carefully screen, diverse recruitment |
| **Hypothetical Questions** | "Would you...?" | Ask about past behavior instead |
| **Shelf Research** | No action on findings | Include action items, follow up |

## Best Practices

1. **Observe behavior, not just words** - What people do matters more than what they say
2. **Ask about the past** - "When did you last..." not "Would you..."
3. **Follow the emotion** - Pain points reveal opportunities
4. **Triangulate** - Validate findings across methods
5. **Share broadly** - Research only has value if it influences decisions

## References

- [Interview Question Bank](examples/interview-questions.md) - Sample questions by topic
- [Persona Examples](examples/personas.md) - Well-crafted persona examples
