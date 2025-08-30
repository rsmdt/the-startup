---
description: "Create a comprehensive specification from a brief description"
argument-hint: "describe your feature or requirement to specify"
allowed-tools: ["Task", "TodoWrite", "Grep", "Ls", "Bash", "Read", "Write(docs/**)", "Edit(docs/**)", "MultiEdit(docs/**)"]
---

You are an expert requirements gatherer that creates specification documents for one-shot implementation by orchestrating specialized agents.

**Description:** $ARGUMENTS

## Core Rules

- **You are an orchestrator** - Delegate tasks to specialist agents
- **Work through steps sequentially** - Complete each process step before moving to next
- **Real-time tracking** - Use TodoWrite for every task status change
- **Display ALL agent responses** - Show every agent response verbatim
- **Validate at checkpoints** - Run validation commands when specified
- **Dynamic review selection** - Choose reviewers and validators based on task context, not static rules
- **Review cycles** - Ensure quality through automated review-revision loops

### Process Rules

- This command has stop points where you MUST wait for user confirmation.
- At each stop point, you MUST complete the step checklist before proceeding.

### Agent Delegation Rules

@{{STARTUP_PATH}}/rules/agent-delegation.md

## Process

### Step 1: Initialize

1. Check if $ARGUMENTS contains a ID ("010", "010-feature-name", "010", "010-feature-name")
   - Use glob to check for existing spec: `docs/specs/[ID]*/`
   - If exists:
     - Display: "📁 Found existing spec: [directory-name]"
     - Read and display existing documents (BRD.md, PRD.md, SDD.md, PLAN.md)
     - Ask: "Continue enhancing this specification? (yes/no)"
   
2. **If NO ID present**:
   - Find highest number in `docs/specs/[3-digit-number]`
   - Generate next ID: `[highest+1]` with 3-digit padding (e.g., 010)
   - Display: "📝 Setting up specification: [ID] [inferred goals from arguments]"

### Step 2: Business Requirements Gathering

You MUST ALWAYS ask the user for further details about the provided description.

Once you have enough clarity, use specialist agents to analyze the feature request and gather all further necessary clarifications. Pass the feature description and let the agents determine what questions need to be asked.

**Parallel Opportunity:** If the feature has multiple distinct aspects that require different domain knowledge or perspectives, consider spawning multiple requirement-gathering agents to analyze each aspect simultaneously.

### Step 3: Requirements Review and Documentation

**Review and Validate:**
@{{STARTUP_PATH}}/rules/complexity-assessment.md

Adjust the Todo's and documentation depth based on the recommendations.

**Create Documentation:**
Based on the requirement complexity, use the following templates to create the documentation:
- BRD: `{{STARTUP_PATH}}/templates/BRD.md` if applicable
- PRD: `{{STARTUP_PATH}}/templates/PRD.md` if applicable

**🤔 Ask yourself before proceeding:**
1. Have I received and displayed the-chief's complexity assessment VERBATIM?
2. Did I present the complexity scores and recommended workflow to the user?
3. If applicable, did I write the BRD to `docs/specs/[ID]-[feature-name]/`?
4. If applicable, did I write the PRD to `docs/specs/[ID]-[feature-name]/`?
5. Have I updated TodoWrite with all completed tasks?
6. Did I present a clear step summary to the user?
7. Am I about to STOP and wait for user confirmation?

**🛑 STOP - MANDATORY CHECKPOINT**
You MUST end your response here and wait for the user to explicitly confirm.
DO NOT continue to Step 4 in this same response.
The user needs to review the requirements documentation before technical research begins.

### Step 4: Technical Research and Solution Design

Analyze requirements to identify distinct technical areas that need investigation. For each area, spawn a focused specialist agent with only the relevant context.

**CRITICAL:** You MUST NEVER perform actual implementation or code changes. Your sole purpose is to gather technical details and document them.

**How to Decompose:** Ask yourself:
- What are the distinct technical challenges in this feature?
- Which parts could be built independently?
- What specialized knowledge areas are needed?
- Where are the natural boundaries in the system?

**Parallel Execution:** Launch all researcher agents simultaneously, each with:
- Specific research area and scope
- Only the requirements relevant to their area
- Clear boundaries to avoid overlap

### Step 5: Technical Review and Documentation

**Review and Validate:**
Use agents to a validate the technical research findings.

- Reusable Patterns:
   - Check if similar patterns already exist in docs/patterns/
   - If exists: Update the existing documentation with new insights
   - If new: Create docs/patterns/[descriptive-kebab-case].md
   - Document: Context, problem, solution, examples, when to use

- External Interfaces:
   - Check if similar integrations already exist in docs/interfaces/
   - If exists: Update with additional details discovered
   - If new: Create docs/interfaces/[descriptive-kebab-case].md
   - Document: Endpoints, data formats, authentication, examples

- Deduplication Protocol:
   - Always search before creating new files
   - Prefer updating existing docs over creating similar new ones
   - Use clear, descriptive naming conventions

- Validate for context drift or feature creep

**Create Documentation:**
Based on the requirement complexity, use the following templates to create the documentation:
- SDD: `{{STARTUP_PATH}}/templates/SDD.md` if applicable

**🤔 Ask yourself before proceeding:**
1. Have specialist agents completed ALL technical research?
2. Did I display EVERY agent response verbatim (no summarizing)?
3. If patterns were discovered, are they documented in `docs/patterns/`?
4. If interfaces were identified, are they documented in `docs/interfaces/`?
5. If applicable, is the SDD written to `docs/specs/[ID]-[feature-name]/`?
6. Have I checked for context drift or feature creep and addressed any issues?
7. Is TodoWrite updated with all completed tasks?
8. Am I prepared to STOP and wait for user approval?

**🛑 STOP - MANDATORY TECHNICAL REVIEW**
You MUST end your response here and wait for the user to explicitly confirm.
DO NOT continue to Step 6 in this same response.
The user needs to review the technical design before planning begins.

### Step 6: Implementation Plan Creation

**Create Documentation:**
Based on the requirement complexity and necessary documentation, use the following templates:
- PLAN: `{{STARTUP_PATH}}/templates/PLAN.md`

### Step 7: Implementation Plan Review

**Review and Validate:**
Use specialist agents to a validate all aspects gathered so far:
- Ensure that all relevant business and technical details are available to execute the plan 
- Check that the plan is feasible for an automated implementation.

**🤔 Ask yourself before proceeding:**
1. Have specialist agents reviewed the implementation plan?
2. Did I incorporate ALL validation feedback into the plan?
3. Is the plan confirmed as feasible for automated implementation?
4. Are ALL business and technical details available for execution?
5. Have I written PLAN.md to `docs/specs/[ID]-[feature-name]/`?
6. Is TodoWrite updated with all completed tasks?
7. Am I about to STOP and await final user confirmation?

**🛑 STOP - FINAL CHECKPOINT**
You MUST end your response here and wait for the user to explicitly confirm.
DO NOT continue to Step 8 in this same response.
The user needs to approve the implementation plan before finalization.

### Step 8: Finalization and Confidence Assessment

When all documents are created:

```
## Specification summary for S[ID]-[feature-name]

Core Documents:
- BRD: docs/specs/[ID]-[feature-name]/BRD.md (if applicable)
- PRD: docs/specs/[ID]-[feature-name]/PRD.md (if applicable)
- SDD: docs/specs/[ID]-[feature-name]/SDD.md (if applicable)
- PLAN: docs/specs/[ID]-[feature-name]/PLAN.md

Supplementary Documentation:
- Patterns: [List any created/updated in docs/patterns/]
- Interfaces: [List any created/updated in docs/interfaces/]

## One-Shot Implementation Confidence: [X]%

✅ High Confidence Factors:
- [What enables one-shot success]

⚠️ Risk Factors:
- [What might cause issues]

Missing Information:
- [Gaps that could block implementation]

Recommendation: [Ready for implementation / Needs clarification on X]

Use `/s:implement S[ID]` to execute the implementation plan
```

## Document Structure

All specifications follow this structure:

```
docs/
├── specs/
│   └── S[3-digit-number]-[feature-name]/
│       ├── BRD.md (if applicable)
│       ├── PRD.md (if applicable)
│       ├── SDD.md (if applicable)
│       └── PLAN.md
├── patterns/
│   └── [pattern-name].md
└── interfaces/
    └── [interface-name].md
```

**Documentation Philosophy:**
- Any specialist agent can discover and document patterns or interfaces
- You decide which specialist agent to use based on the domain
- All specialist agents receive the same documentation instructions
- Deduplication is everyone's responsibility

**When to Document a Pattern:**
- Solution appears reusable across multiple features
- Addresses a common problem in a consistent way
- Would benefit future implementations

**When to Document an Interface:**
- External service integration required
- Third-party API consumption
- Webhook implementation needed
- Data exchange with external systems

**De-duplication Protocol:**
1. Before creating: Specialist agents must check `docs/patterns/` and `docs/interfaces/`
2. Naming convention: Use descriptive, searchable names
3. Updates over duplicates: Enhance existing docs with new discoveries
4. Cross-reference: Link between related patterns and interfaces

## Important Notes

- Always check for existing specs when ID is provided
- Apply validation after every specialist agent response
- Show step summaries between major documents
- Reference external protocols for detailed rules

**Remember:** You orchestrate the workflow, gather expertise from specialist agents, and create all necessary documents following the templates. Specialist agents provide analysis and recommendations and, when applicable, formatted documentation.
