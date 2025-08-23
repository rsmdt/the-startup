# Assumption Prevention Pattern

## Context
Requirements-gathering and validation agents (business-analyst, product-manager, compliance-officer, lead-engineer, qa-lead, chief) tend to make inferences and assumptions instead of asking clarifying questions, leading to misaligned implementations.

## Problem
- Agents infer user intent instead of explicitly asking for clarification
- Assumptions made early compound into larger misunderstandings later
- Scope drift occurs when agents expand beyond original requirements without validation
- Critical details get assumed rather than confirmed with stakeholders

## Solution
Systematic validation protocol embedded in agent workflows:

### 1. Entry Validation Protocol
Every requirements/validation agent must ask before proceeding:
- **Context**: "Show me how you currently handle [specific scenario]"
- **Constraints**: "What are your hard limits on [time/budget/scope/technology]?"
- **Success**: "How will you know this worked? What metrics matter?"
- **Scope**: "What are we explicitly NOT doing in this phase?"
- **Stakeholders**: "Who else needs to approve/validate this decision?"

### 2. "I'm Assuming" Pattern
```markdown
## VALIDATION CHECKPOINT

I'm making these assumptions - please confirm or correct:

1. **[SCOPE]**: I'm assuming [specific assumption] based on [evidence]
   - **Correct?** [Yes/No] 
   - **If No**: [What should I assume instead?]

2. **[PRIORITY]**: I'm assuming [specific assumption] based on [evidence]  
   - **Correct?** [Yes/No]
   - **If No**: [What should I assume instead?]

## SCOPE CONFIRMATION
I understand we're NOT doing:
- [Specific excluded item 1]
- [Specific excluded item 2]

Is this correct? What am I missing?
```

### 3. Context Drift Detection
- Quote original user request before expanding scope
- "You originally asked for X. I'm now considering Y because Z. Should I continue?"
- Flag divergence: "This goes beyond your original scope of X"

## Examples

**Before** (the-business-analyst making assumptions):
```markdown
Based on your e-commerce request, I'll assume you need:
- User registration and login
- Product catalog with search
- Shopping cart functionality
- Payment processing integration
```

**After** (the-business-analyst preventing assumptions):
```markdown
I need to understand your e-commerce requirements better before proceeding:

1. **User Management**: Do you need user accounts, or can customers checkout as guests?
2. **Product Scope**: How many products? Categories? Variants (size, color)?
3. **Payment Method**: What payment methods must you support?
4. **Geographic Scope**: Which countries/regions do you serve?
5. **Integration Constraints**: What existing systems must this work with?

I'm assuming this is a new system (not integrating with existing e-commerce). Correct?
```

## When to Use
- Agent role involves gathering requirements or making decisions
- Agent tends to expand scope beyond user's original request
- Agent makes technical or business assumptions that should be validated
- Agent reviews work that may have drifted from original requirements

## Benefits
- **Clarity**: Explicit validation prevents misaligned implementations
- **Scope Control**: Prevents feature creep and scope expansion without approval
- **Stakeholder Alignment**: Ensures decisions are validated with right people
- **Quality**: Reduces rework from incorrect assumptions

## Implementation Notes
- Embed validation protocols directly in agent approach sections
- Use structured "I'm assuming X, please confirm" format
- Require 3-5 clarifying questions before making recommendations  
- Cross-reference original requirements when detecting scope expansion
- Make assumptions explicit and require user confirmation before proceeding