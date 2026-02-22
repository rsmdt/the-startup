# Usability Testing Patterns

Structured approaches for evaluating interface effectiveness through direct user observation.

## Test Types

**Moderated Testing** — Facilitator guides participant
- Best for: Complex tasks, early prototypes, need for probing
- Pros: Rich qualitative data, can adapt on the fly
- Cons: Time-intensive, facilitator can bias

**Unmoderated Testing** — Participant works independently
- Best for: Simple tasks, large samples, geographic spread
- Pros: Scalable, no facilitator bias, natural behavior
- Cons: No probing, participants may give up

**Guerrilla Testing** — Quick tests with available people
- Best for: Early validation, simple concepts, tight timelines
- Pros: Fast, cheap, good for iteration
- Cons: May not match target users, limited depth

## Test Protocol Structure

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

## Writing Effective Task Scenarios

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

## Metrics to Capture

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

## Severity Rating Scale

| Rating | Name | Definition | Action |
|--------|------|------------|--------|
| 1 | Cosmetic | Noticed but no impact | Fix if time permits |
| 2 | Minor | Slight difficulty, recovered easily | Fix in next release |
| 3 | Major | Significant difficulty, delayed success | Fix before release |
| 4 | Critical | Prevented task completion | Must fix immediately |
