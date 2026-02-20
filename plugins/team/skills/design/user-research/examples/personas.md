# Persona Examples

## Context

This document shows what complete personas look like in practice. It includes two contrasting examples — a proto-persona built from assumptions before research, and a research-backed persona built from interview synthesis — followed by a blank template for new personas.

The key distinction between these types is not format but evidence. A proto-persona is a documented hypothesis. A research-backed persona is a finding. Both are useful. Neither should be mistaken for the other.

---

## Proto-Persona (Assumption-Based)

Proto-personas are created before or instead of research, usually to align a team on a shared working hypothesis or to plan a study. They are honest about their origins — they represent what the team believes, not what users have confirmed.

Label all proto-personas clearly. Treat them as assumptions to be tested.

---

```
┌─────────────────────────────────────────────────────────────┐
│ [PHOTO PLACEHOLDER]                                          │
│                                                              │
│ NAME: Marcus Chen                                            │
│ TITLE: Founder, early-stage B2B SaaS startup                │
│ ARCHETYPE: Resourceful Builder                               │
│                                                              │
│ SOURCE: PROTO-PERSONA — based on team assumptions,          │
│ not validated with research                                  │
├─────────────────────────────────────────────────────────────┤
│ QUOTE:                                                       │
│ "I need to move fast. I'll figure out the proper way        │
│ to do it once we have traction."                            │
├─────────────────────────────────────────────────────────────┤
│ DEMOGRAPHICS:                                                │
│ Age: 28-38     Experience: 2-6 years as founder             │
│ Context: Remote, works across product / sales / ops         │
│ Team size: 1-8 people                                       │
├─────────────────────────────────────────────────────────────┤
│ GOALS:                                                       │
│ - Primary: Ship product quickly to validate the market      │
│ - Secondary: Keep costs low while reaching first $10K MRR   │
├─────────────────────────────────────────────────────────────┤
│ PAIN POINTS (ASSUMED):                                       │
│ - Too many tools pulling attention in different directions   │
│ - No dedicated ops or admin support — does everything        │
│ - Hard to prioritize when everything feels urgent            │
├─────────────────────────────────────────────────────────────┤
│ BEHAVIORS (ASSUMED):                                         │
│ - Heavy Twitter / X and Slack user                           │
│ - Prefers products with quick time-to-value, no onboarding  │
│ - Makes purchase decisions alone, often on impulse           │
│ - Will switch tools if something cheaper appears             │
├─────────────────────────────────────────────────────────────┤
│ TECH COMFORT: High. Will use APIs, configure webhooks,      │
│ read docs to solve problems independently.                  │
├─────────────────────────────────────────────────────────────┤
│ OPEN ASSUMPTIONS (MUST TEST):                               │
│ - Do founders actually buy this themselves, or do they       │
│   delegate it once they see traction?                        │
│ - Is speed really the primary value, or is it trust /       │
│   reliability?                                               │
│ - Does this segment actually experience the pain we think?   │
├─────────────────────────────────────────────────────────────┤
│ SCENARIO:                                                    │
│ Marcus discovers the product through a tweet, signs up for   │
│ a free trial without reading the landing page in full, and  │
│ tries to connect his existing stack in the first 10 minutes. │
│ If he hits a wall, he leaves. If it works, he upgrades.     │
└─────────────────────────────────────────────────────────────┘
```

**When to use a proto-persona:**
- Early in a project before research is possible
- To align stakeholders on who you are and are not designing for
- To scope a research study by making assumptions explicit

**Important constraint from SKILL.md:** Proto-personas must never be presented as research findings. They are a starting position, not a conclusion. Label them explicitly and use them to drive research, not replace it.

---

## Research-Backed Persona

Research-backed personas are synthesized from multiple interview participants. Each attribute is grounded in observed behavior or direct quotes from at least two to three participants. The evidence column is the most important part — it is what separates a persona from a guess.

---

```
┌─────────────────────────────────────────────────────────────┐
│ [PHOTO PLACEHOLDER]                                          │
│                                                              │
│ NAME: Priya Nair                                             │
│ TITLE: Operations Lead, mid-size B2B software company        │
│ ARCHETYPE: The Reluctant Coordinator                         │
│                                                              │
│ SOURCE: RESEARCH-BACKED — synthesized from 9 interviews,     │
│ validated across 3 participant segments                      │
├─────────────────────────────────────────────────────────────┤
│ QUOTE:                                                       │
│ "My job is to make sure nothing falls through the cracks.   │
│ The problem is I'm always finding out too late that          │
│ something already did."                                      │
│ — P4, Operations Lead, 6 years experience                   │
├─────────────────────────────────────────────────────────────┤
│ DEMOGRAPHICS:                                                │
│ Age: 30-45     Experience: 4-10 years in ops roles           │
│ Context: Office or hybrid; manages cross-functional          │
│ workflows across teams of 20-80 people                      │
├─────────────────────────────────────────────────────────────┤
│ GOALS:                                                       │
│ - Primary: Maintain visibility across all in-flight work     │
│   without becoming the bottleneck for every update           │
│ - Secondary: Reduce meeting load by making status visible    │
│   asynchronously                                            │
├─────────────────────────────────────────────────────────────┤
│ PAIN POINTS (EVIDENCE-BACKED):                               │
│                                                              │
│ 1. Status is scattered across too many places               │
│    Evidence: 7/9 participants described checking 3+ tools    │
│    (email, Slack, spreadsheets, project tools) just to       │
│    answer "where does this stand?"                           │
│    Quote: "I have four tabs open just to give one update."   │
│    — P7                                                      │
│                                                              │
│ 2. Escalations arrive too late                               │
│    Evidence: 6/9 participants described situations where a   │
│    problem surfaced at standup or in a meeting that should   │
│    have been flagged a day or more earlier                   │
│    Quote: "By the time I hear about it, the deadline is      │
│    already at risk." — P2                                    │
│                                                              │
│ 3. Reminder and follow-up work is manual and time-consuming  │
│    Evidence: 5/9 participants manually tracked who had       │
│    responded to requests via spreadsheet or flagged emails   │
│    Quote: "I paste the same message into Slack three times   │
│    a week just to get a response." — P4                     │
├─────────────────────────────────────────────────────────────┤
│ BEHAVIORS (OBSERVED):                                        │
│ - Builds personal tracking systems (spreadsheets, Notion     │
│   pages) to compensate for tool fragmentation               │
│ - Sends Slack messages rather than creating formal tickets   │
│   because it's faster, even when a ticket would help more   │
│ - Prefers to check status herself rather than interrupt      │
│   teammates — respects others' focus time                   │
│ - Exports data from tools manually to create summary reports │
│   for leadership on a weekly basis                          │
├─────────────────────────────────────────────────────────────┤
│ TECH COMFORT: Medium. Uses existing tools confidently but    │
│ does not configure integrations or write automations. Relies │
│ on IT or engineering for anything technical.                 │
├─────────────────────────────────────────────────────────────┤
│ SCENARIO:                                                    │
│ It is 3:00 PM on a Thursday. Priya has a leadership sync in  │
│ 90 minutes and needs to confirm that three open deliverables │
│ are on track. She opens Slack to message three different     │
│ people, checks the project tool to see if statuses were     │
│ updated (they were not), then falls back to her spreadsheet. │
│ She composes a "quick update" message that takes 20 minutes  │
│ to send because she doesn't want to seem like she's nagging. │
│ She goes into the sync knowing she has incomplete info.      │
├─────────────────────────────────────────────────────────────┤
│ DESIGN IMPLICATIONS:                                         │
│ - Visibility without friction is the core value proposition  │
│ - Automation must reduce her work, not create new admin      │
│ - Notifications must be smart — she is already overwhelmed  │
│ - Language should not feel corporate — she uses informal     │
│   channels by choice                                         │
└─────────────────────────────────────────────────────────────┘
```

**What makes this persona research-backed:**
- Every pain point cites the number of participants who expressed it
- At least one direct quote supports each attribute
- Behaviors describe observed actions, not inferred preferences
- Design implications follow from the evidence, not from assumptions

---

## Negative Persona Example

Negative personas define who the product is not designed for. They prevent design decisions from being pulled toward the wrong user by stakeholders who advocate for edge cases.

```
┌─────────────────────────────────────────────────────────────┐
│ NAME: "The Power Admin"                                      │
│ ARCHETYPE: Deep configurator, technically proficient         │
├─────────────────────────────────────────────────────────────┤
│ WHY NOT OUR USER:                                            │
│ This person needs granular permission controls, custom        │
│ workflow logic, and API access. Designing for their needs    │
│ would add complexity that makes our product harder to use    │
│ for Priya — our primary persona.                            │
│                                                              │
│ They are better served by [enterprise competitor].           │
├─────────────────────────────────────────────────────────────┤
│ REFERENCE:                                                   │
│ When a stakeholder requests a feature that only a power      │
│ admin would use, refer back to this persona to ground the   │
│ conversation in tradeoffs.                                   │
└─────────────────────────────────────────────────────────────┘
```

---

## Blank Persona Template

Copy this template when creating personas from a new research study. Fill in only what your research supports. Leave fields blank rather than filling them with assumptions — missing evidence is better than fabricated evidence.

```
┌─────────────────────────────────────────────────────────────┐
│ [PHOTO PLACEHOLDER]                                          │
│                                                              │
│ NAME: [Fictional name]                                       │
│ TITLE: [Role / context]                                      │
│ ARCHETYPE: [2-3 word descriptor]                            │
│                                                              │
│ SOURCE: [ ] PROTO-PERSONA (assumptions)                      │
│         [ ] RESEARCH-BACKED (synthesized from N interviews)  │
├─────────────────────────────────────────────────────────────┤
│ QUOTE:                                                       │
│ "[Direct quote from a participant — attribute to P#]"        │
├─────────────────────────────────────────────────────────────┤
│ DEMOGRAPHICS:                                                │
│ Age: [Range]     Experience: [Level / years]                 │
│ Context: [Work environment, team size, role scope]          │
├─────────────────────────────────────────────────────────────┤
│ GOALS:                                                       │
│ - Primary: [Main objective — what success looks like]        │
│ - Secondary: [Supporting objective]                          │
├─────────────────────────────────────────────────────────────┤
│ PAIN POINTS:                                                 │
│                                                              │
│ 1. [Pain point headline]                                    │
│    Evidence: [N/N participants; direct quote — P#]           │
│                                                              │
│ 2. [Pain point headline]                                    │
│    Evidence: [N/N participants; direct quote — P#]           │
│                                                              │
│ 3. [Pain point headline]                                    │
│    Evidence: [N/N participants; direct quote — P#]           │
├─────────────────────────────────────────────────────────────┤
│ BEHAVIORS:                                                   │
│ - [How they approach the problem domain]                    │
│ - [Tools, workarounds, or systems they have built]          │
│ - [Decision-making patterns relevant to your product]       │
│ - [Social / collaborative behavior — who do they involve?]  │
├─────────────────────────────────────────────────────────────┤
│ TECH COMFORT: [Low / Medium / High — describe what they      │
│ can and cannot do independently]                            │
├─────────────────────────────────────────────────────────────┤
│ OPEN QUESTIONS: (for proto-personas only)                    │
│ - [Assumption that still needs validation]                   │
│ - [Assumption that still needs validation]                   │
├─────────────────────────────────────────────────────────────┤
│ SCENARIO:                                                    │
│ [A concrete, present-tense narrative of this persona         │
│ experiencing the core problem your product addresses.        │
│ 3-6 sentences. Specific details, not generalities.]         │
├─────────────────────────────────────────────────────────────┤
│ DESIGN IMPLICATIONS:                                         │
│ - [What this persona requires from the product]             │
│ - [What would break trust or cause abandonment]             │
│ - [Constraints on complexity, workflow, or tone]            │
└─────────────────────────────────────────────────────────────┘
```

---

## Persona Quality Checklist

Before sharing a persona with stakeholders, verify:

- [ ] The persona type (proto vs. research-backed) is clearly labeled
- [ ] Every pain point cites the number of participants who confirmed it
- [ ] At least one direct quote supports each claim
- [ ] The quote is attributed to a participant ID, not paraphrased
- [ ] Behaviors describe observed actions, not inferred preferences
- [ ] The scenario is specific enough to use in a design critique
- [ ] Design implications follow directly from evidence
- [ ] The persona has been reviewed against at least 3 participant transcripts

**Fail condition:** If a team member reads the persona and cannot tell which attributes are observed vs. assumed, the persona needs revision.
