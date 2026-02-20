# Interview Question Bank

## Context

This document provides a structured bank of questions for stakeholder interviews across different phases of requirements elicitation. Questions are organized by purpose and stakeholder type. Select questions based on your session goals — not all questions apply to every interview.

Use these in conjunction with the interview structure from SKILL.md (Context → Current State → Desired State → Constraints → Wrap-Up).

---

## Discovery Questions

These questions establish context and reveal the real problem. Ask them early to avoid solving the wrong thing.

### Understanding the Initiative

- What's the business problem you're trying to solve?
- What happens if you don't solve this? What's the cost of inaction?
- What prompted this initiative now? Why not six months ago?
- Who requested this work, and why do they care?
- How does this fit into the broader product or company strategy?
- Is there a deadline, and what's driving it?

### Uncovering the Real Need (The 5 Whys Starters)

- Walk me through why this matters to you personally.
- What would success look like in six months?
- If you had a magic wand and this was solved perfectly, what would be different?
- What's the worst thing that could happen if we build the wrong thing here?
- Is this a new problem or one that's been around for a while?

### Stakeholder Landscape

- Who else is affected by this problem?
- Who do you think I should talk to next?
- Are there people who benefit from the current situation staying the same?
- Who signs off on this decision?
- Are there stakeholders outside your team — customers, partners, regulators?

---

## Current State Questions

These questions document what exists today so you can identify pain points and understand workflows before proposing solutions.

### Process and Workflow

- Walk me through how you handle this today, step by step.
- How often do you do this? Daily, weekly, monthly?
- How long does this take from start to finish?
- Who is involved at each step?
- What triggers this process to start? What marks it as done?

### Tools and Systems

- What tools, systems, or software do you currently use for this?
- How do these systems talk to each other? Do they, or do you manually move data?
- How many places do you have to go to complete this task?
- What data do you need, and where does it come from?

### Pain Points

- What's the most frustrating part of the current process?
- Where do things break down or slow down?
- Where do errors happen most often? What causes them?
- What workarounds have you built to cope with the current system?
- What do you spend time on that you shouldn't have to?
- If you could change one thing about how this works today, what would it be?

### Volume and Scale

- How many people do this task?
- How much data is involved — records, transactions, files?
- What happens during peak periods? How does the system cope?
- Has this grown significantly in the past year? Where do you expect it to be in two years?

---

## Desired State Questions

These questions clarify the goal and prevent gold plating by anchoring requirements to outcomes rather than features.

### Defining Success

- Describe the ideal state — what does "done" look like?
- What specific outcome tells you this worked?
- What metrics would you measure to know it's working?
- What's the minimum you need to declare this a success?
- If we had to cut scope, what absolutely must be there for this to be useful?

### Workflows and Interactions

- Walk me through a typical scenario using the new solution.
- What should happen when things go right? When things go wrong?
- Who starts the process? Who hands it off? Who finishes it?
- What decisions does a person have to make, and which ones should the system make?

### Priority and Trade-offs

- If you had to rank these capabilities — fast, cheap, correct — which matters most?
- What can we defer to a later phase without losing the core value?
- Are there features you think would be nice but aren't essential?
- Would you rather have something basic by [date] or something complete later?

---

## Validation Questions

These questions confirm your understanding and surface gaps before you commit to a direction.

### Checking Comprehension

- Let me summarize what I've heard — please correct me where I'm wrong.
- Is there anything important I haven't asked?
- What assumptions are you making that I might not share?
- What would you add if you had more time?

### Probing Completeness

- Are there edge cases we haven't discussed?
- What happens if a user does something unexpected?
- What about users with accessibility needs or limited technical ability?
- Are there regulatory, legal, or compliance requirements that affect this?

### Surfacing Conflict

- Have you talked to [other stakeholder]? Do you know if they agree with this?
- Are there people who see this differently than you do?
- What's the most controversial aspect of what you've described?

---

## Stakeholder-Specific Questions

Different roles have different perspectives. Tailor your questions to surface insights specific to each stakeholder's vantage point.

### End Users

These questions focus on lived experience and day-to-day workflow. End users reveal what actually happens, not what the process says should happen.

- Tell me about the last time you did this task. Walk me through it.
- What would you do first if this tool disappeared tomorrow?
- What do you wish the system just knew, so you didn't have to tell it?
- When does the system get in your way?
- Are there parts of your job this doesn't touch that it probably should?
- What do new people struggle to learn about this process?

### Product Owners / Business Stakeholders

These questions focus on value, priority, and business outcomes.

- What is the business impact of getting this right? Of getting it wrong?
- What customer segments does this affect?
- How does this compare in priority to other things on the roadmap?
- What are customers telling you about this problem?
- What does the competition do here? Is that relevant?
- What would make you confident enough to demo this to a customer?

### Engineering and Technical Teams

These questions uncover implementation constraints and integration complexity.

- Are there existing systems this needs to integrate with?
- What parts of the current architecture make this harder than it should be?
- Are there data quality or availability issues we need to account for?
- What are the constraints we can't change — legacy systems, contracts, infrastructure?
- What's the riskiest part of building this from a technical standpoint?
- Are there security, privacy, or compliance requirements that aren't obvious?

### Executives and Sponsors

These questions focus on strategic alignment, risk tolerance, and go/no-go criteria.

- What does success look like in terms of business outcomes, not features?
- What's the risk you're most concerned about?
- What would make you pull the plug on this initiative?
- Who are the key decision-makers, and what do they care about?
- Is there a point at which this becomes "too expensive" relative to the value?
- What external pressures — market, regulatory, competitive — are driving this timeline?

---

## Domain-Specific Question Sets

### For Internal Tools and Operations

- Who owns the process this tool supports?
- How is performance measured today, and how would you measure it with the new tool?
- What reporting or audit requirements exist?
- What happens during off-hours or when people are out of the office?

### For Customer-Facing Products

- Who is the target user — are they tech-savvy, occasional users, power users?
- What devices and environments will they use this on?
- What does the user already know when they arrive at this feature?
- What's the user's goal, and what might distract them from it?
- What's the cost to the user if they make a mistake?

### For Data and Reporting Features

- What decisions will be made using this data?
- Who will look at this, and how often?
- What's the acceptable latency — real-time, hourly, daily?
- What's the source of truth, and how reliable is it?
- What does "wrong" look like — what data errors have caused problems before?

### For Integration and API Features

- What systems will call this, and what do they expect back?
- What are the upstream and downstream dependencies?
- What's the acceptable error rate and latency?
- How will failures be handled — retry, fallback, alert?
- Who is responsible for the other side of this integration?

---

## Wrap-Up Questions

Always close interviews with these to surface what you missed and identify next steps.

- What haven't I asked that I should have?
- Who else should I talk to about this?
- Is there any documentation, existing specs, or prior work I should read?
- Can I follow up with you if I have clarifying questions?
- What's the best way to share my summary of this conversation with you for review?
