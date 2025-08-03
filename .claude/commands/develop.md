---
allowed-tools: ["Task", "TodoWrite"]
description: "Iterative development process with parallel specialist execution and user checkpoints"
argument-hint: "describe your project or feature requirements"
---

# 🚀 Iterative Development Pipeline

You will orchestrate an iterative development process for: **$ARGUMENTS**

## 🚫 CRITICAL: Delegation Only Mode

**You are NOT allowed to implement anything directly. You MUST delegate ALL work to specialist agents.**

Your role is strictly to:
- Coordinate specialist agents
- Display their commentary verbatim
- Track tasks in master todo list
- Present findings for user decisions

## Overview

This process uses:
- **Iterative discovery loops** (max 3 iterations per phase)
- **Parallel specialist execution** for independent tasks
- **User checkpoints** after each phase for feedback
- **Dynamic task generation** based on findings

---

## Phase 1: Initial Strategic Assessment

You will start with a high-level strategic assessment:

<Task description="Initial strategic assessment" 
      prompt="Provide initial strategic analysis for: '$ARGUMENTS'. Identify complexity, key areas needing investigation, and which specialists should be engaged for discovery. Focus on what we need to learn, not solutions."
      subagent_type="the-chief" />

### ⚠️ IMPORTANT: Response Handling Protocol

When the Chief responds, you MUST:
1. **Display the `<commentary>` block EXACTLY as provided** (no changes)
2. **Extract and list all tasks from the `<tasks>` block**
3. **Create/update master todo list with these tasks**
4. **STOP and present a user checkpoint** before proceeding

### Chief's Assessment:
[The Chief's commentary will be displayed here verbatim]

### Extracted Tasks:
[Tasks from the Chief's response will be listed here]

### 📋 Strategic Assessment Checkpoint

Based on the Chief's analysis, you need user input:

**Options:**
1. **Proceed with discovery** - Run the recommended discovery tasks
2. **Adjust focus** - Modify the discovery approach
3. **Skip to planning** - If requirements are already clear
4. **Request clarification** - If the project needs refinement

**Please choose an option or provide specific guidance.**

---

## 🔄 Discovery Loop (Iteration 1 of max 3)

Based on the Chief's assessment, you will now run parallel discovery tasks:

### Parallel Discovery Tasks:

#### Business Analysis Track:
<Task description="Requirements discovery - User perspective" 
      prompt="Explore user needs and workflows for: '$ARGUMENTS'. Focus on who will use this, their goals, and success criteria."
      subagent_type="the-business-analyst" />

<Task description="Requirements discovery - Business context" 
      prompt="Investigate business drivers and constraints for: '$ARGUMENTS'. What problem does this solve? What are the business impacts?"
      subagent_type="the-business-analyst" />

<Task description="Requirements discovery - Integration points" 
      prompt="Identify system dependencies and integration needs for: '$ARGUMENTS'. What existing systems must this work with?"
      subagent_type="the-business-analyst" />

#### Technical Discovery Track:
<Task description="Architecture exploration - Scalability" 
      prompt="Assess scalability requirements and patterns for: '$ARGUMENTS'. Consider current and future load, growth patterns."
      subagent_type="the-architect" />

<Task description="Architecture exploration - Technology options" 
      prompt="Evaluate technology choices for: '$ARGUMENTS'. Compare options, consider team skills, existing stack."
      subagent_type="the-architect" />

<Task description="Security landscape assessment" 
      prompt="Identify security considerations for: '$ARGUMENTS'. What are the threats, compliance needs, data sensitivity?"
      subagent_type="the-security-engineer" />

### Discovery Results:

**⚠️ CRITICAL: For EACH agent response above, you MUST:**
1. **Display their `<commentary>` blocks EXACTLY as provided**
2. **Show their key findings clearly**
3. **Extract any new tasks they identified**

[All parallel discovery results will be displayed here with commentary preserved]

---

## Re-evaluation by the Chief

Now you will have the Chief evaluate these findings:

<Task description="Evaluate discovery findings" 
      prompt="Review the discovery findings above for: '$ARGUMENTS'. Are requirements clear enough? Do we need more discovery? What are the key decisions or unknowns? Should we iterate or proceed to planning?"
      subagent_type="the-chief" />

### Chief's Re-evaluation:

**⚠️ You MUST display the Chief's `<commentary>` block here verbatim**

[The Chief's assessment will appear here]

---

## 📋 User Checkpoint

### Summary of Findings:

**Requirements Understanding:**
- [Key requirements discovered]
- [User needs identified]
- [Business constraints]

**Technical Considerations:**
- [Architecture options]
- [Technology recommendations]
- [Security requirements]

**Next Steps Options:**
1. **Continue Discovery** - If requirements still unclear (Iteration 2 of 3)
2. **Proceed to Planning** - If we have sufficient clarity
3. **Pivot Direction** - If discoveries reveal different needs

**What would you like to do?**
- Type "continue" to run another discovery iteration
- Type "proceed" to move to planning phase
- Type "pivot" to adjust the approach
- Or provide specific guidance on areas to explore

---

!if [ "$USER_CHOICE" = "continue" ] && [ $ITERATION -lt 3 ]; then

## 🔄 Discovery Loop (Iteration 2 of max 3)

Based on gaps identified, running targeted discovery:

[Additional parallel discovery tasks based on Chief's evaluation]

!elif [ "$USER_CHOICE" = "proceed" ]; then

## Phase 2: Planning & Design

Now you will move to formal planning with parallel tracks:

### Planning Track:
<Task description="Create product roadmap" 
      prompt="Based on discoveries, create a phased product plan for: '$ARGUMENTS'. Define MVP, future phases, priorities."
      subagent_type="the-product-manager" />

<Task description="Technical design document" 
      prompt="Create detailed technical design for: '$ARGUMENTS' based on discoveries. Include architecture, data flow, APIs."
      subagent_type="the-architect" />

<Task description="Project breakdown" 
      prompt="Create work breakdown structure for: '$ARGUMENTS'. Identify parallel work streams, dependencies, milestones."
      subagent_type="the-project-manager" />

### Planning Results:
[All planning documents will appear here]

---

## 📋 Planning Checkpoint

### Proposed Approach:

**Product Plan:**
[Key features and phases]

**Technical Design:**
[Architecture and technology decisions]

**Project Structure:**
[Work streams and timeline]

**Ready to proceed with implementation?**
- Type "implement" to start development
- Type "refine" to adjust plans
- Or provide specific feedback

!fi

---

## Phase 3: Implementation (When Approved)

### Parallel Implementation Streams:

Based on the project plan, you will coordinate parallel development:

#### Development Teams:
[Multiple parallel developer agents for different components]
[Data engineer agents for data layer]
[DevOps agents for infrastructure]

Each stream will:
1. Implement assigned components
2. Write comprehensive tests
3. Document their work
4. Report blockers

---

## Implementation Protocol

### Key Principles:
1. **No Monolithic Agents** - Each specialist handles focused aspects
2. **Iterative Discovery** - Max 3 loops with Chief re-evaluation
3. **User Checkpoints** - Stop after each phase for feedback
4. **Parallel Execution** - Independent tasks run simultaneously
5. **Dynamic Adaptation** - Workflow adjusts based on findings

### Phase Flow:
1. **Discovery Loop** (iterative with checkpoints)
2. **Planning Phase** (parallel tracks with checkpoint)
3. **Implementation** (parallel streams with checkpoints)
4. **Quality & Launch** (parallel validation with final checkpoint)

---

*This command provides an adaptive, iterative development process with continuous user involvement.*