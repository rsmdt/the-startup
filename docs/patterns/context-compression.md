# Pattern: Bounded Context Compression

## Context
Sub-agent communication in the enhanced `/s:specify` system must achieve 70% token reduction while maintaining sufficient context for high-quality output. Current full-context delegation creates 10-15x token overhead.

## Problem
- Full context sharing consumes excessive tokens (1,500+ for 100-token tasks)
- Sub-agents receive irrelevant information that dilutes focus
- Communication overhead dominates actual processing time
- Token costs scale linearly with delegation frequency

## Solution
Implement bounded context compression that extracts only essential information for sub-agent execution:

1. **Core Requirement Distillation**: Maximum 50 words capturing essential task
2. **Constraint Extraction**: Only relevant limitations and requirements  
3. **Context Minimization**: Compressed background information
4. **Format Specification**: Clear output expectations
5. **Success Criteria**: Concrete validation measures

### Context Compression Algorithm

```
INPUT: full_session_context, task_description, target_agent
OUTPUT: bounded_context (JSON)

1. EXTRACT_CORE_REQUIREMENT:
   core = summarize(task_description, max_words=50)
   validate_completeness(core, task_description)

2. IDENTIFY_CONSTRAINTS:
   constraints = extract_limitations(session_context)
   filter_by_relevance(constraints, target_agent.domain)

3. COMPRESS_BACKGROUND:
   relevant_context = filter_context_by_agent_needs(full_context, target_agent)
   compressed = summarize_context(relevant_context, max_tokens=200)

4. SPECIFY_OUTPUT_FORMAT:
   format = determine_expected_output(task_type, agent_capabilities)
   
5. DEFINE_SUCCESS_CRITERIA:
   criteria = extract_validation_requirements(task_description)
   make_measurable(criteria)

6. VALIDATE_TOKEN_BUDGET:
   total_tokens = calculate_tokens(bounded_context)
   assert total_tokens <= (original_tokens * 0.3)  # 70% reduction
```

## Implementation Example

### Bounded Context Schema
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "BoundedContext",
  "type": "object",
  "properties": {
    "core_requirement": {
      "type": "string",
      "maxLength": 250,
      "description": "Essential task description in 50 words or less"
    },
    "constraints": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "maxItems": 5,
      "description": "Critical limitations and requirements"
    },
    "assumptions": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "maxItems": 3,
      "description": "Explicit assumptions to validate or challenge"
    },
    "expected_format": {
      "type": "string",
      "enum": ["document", "code", "diagram", "analysis", "checklist"],
      "description": "Required output format"
    },
    "success_criteria": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "minItems": 1,
      "maxItems": 5,
      "description": "Measurable success indicators"
    },
    "minimal_context": {
      "type": "object",
      "properties": {
        "project_type": {"type": "string"},
        "technology_stack": {"type": "array", "items": {"type": "string"}},
        "team_constraints": {"type": "array", "items": {"type": "string"}},
        "timeline": {"type": "string"}
      },
      "additionalProperties": false
    },
    "token_budget": {
      "type": "integer",
      "minimum": 100,
      "maximum": 500,
      "description": "Maximum tokens for agent response"
    }
  },
  "required": ["core_requirement", "expected_format", "success_criteria"],
  "additionalProperties": false
}
```

### Go Implementation
```go
type BoundedContext struct {
    CoreRequirement  string            `json:"core_requirement"`
    Constraints      []string          `json:"constraints,omitempty"`
    Assumptions      []string          `json:"assumptions,omitempty"`
    ExpectedFormat   string            `json:"expected_format"`
    SuccessCriteria  []string          `json:"success_criteria"`
    MinimalContext   map[string]interface{} `json:"minimal_context,omitempty"`
    TokenBudget      int               `json:"token_budget"`
}

type ContextCompressor struct {
    tokenCounter TokenCounter
    summarizer   TextSummarizer
    extractor    ConstraintExtractor
}

func (c *ContextCompressor) Compress(
    fullContext SessionContext, 
    task string, 
    targetAgent AgentSpec,
) (*BoundedContext, error) {
    
    // Core requirement distillation
    core, err := c.summarizer.Distill(task, 50)
    if err != nil {
        return nil, fmt.Errorf("failed to distill core requirement: %w", err)
    }
    
    // Constraint extraction
    constraints := c.extractor.ExtractRelevant(fullContext, targetAgent.Domain)
    
    // Context compression
    minimalContext := c.compressContext(fullContext, targetAgent, 200)
    
    // Validate token budget
    bounded := &BoundedContext{
        CoreRequirement: core,
        Constraints:     constraints,
        ExpectedFormat:  determineFormat(task, targetAgent),
        SuccessCriteria: extractCriteria(task),
        MinimalContext:  minimalContext,
        TokenBudget:     calculateBudget(task, targetAgent),
    }
    
    if c.tokenCounter.Count(bounded) > fullContext.TokenCount*0.3 {
        return nil, fmt.Errorf("compression failed: %d tokens exceeds 30%% budget", 
            c.tokenCounter.Count(bounded))
    }
    
    return bounded, nil
}
```

### Example Compression

**Before (Full Context - 1,247 tokens):**
```
User wants to create a user authentication system for their e-commerce platform. 
The platform is built using React frontend with TypeScript, Node.js backend with 
Express, PostgreSQL database, and is deployed on AWS. The team consists of 3 
developers with varying experience levels. They need to support social login 
(Google, Facebook), email/password auth, password reset, account verification, 
and role-based access control. The system should integrate with their existing 
user management system and support their mobile app in the future. Security is 
a top priority due to PCI compliance requirements. They have a tight deadline 
of 6 weeks and limited budget for external services. The current user table 
has 10,000+ existing users that need to be migrated. Performance requirements 
include sub-200ms login response times and support for 1000 concurrent users...
```

**After (Bounded Context - 312 tokens):**
```json
{
  "core_requirement": "Design secure authentication system supporting social login, email/password, and RBAC for e-commerce platform with PCI compliance requirements.",
  "constraints": [
    "Must integrate with existing user management system",
    "PCI compliance required for payment processing",
    "6-week timeline with limited budget for external services",
    "Support 1000 concurrent users with sub-200ms response"
  ],
  "assumptions": [
    "Current user table migration is feasible",
    "Mobile app integration can be phase 2"
  ],
  "expected_format": "document",
  "success_criteria": [
    "Secure authentication flows documented",
    "Integration points clearly defined",
    "PCI compliance approach specified",
    "Migration strategy included"
  ],
  "minimal_context": {
    "technology_stack": ["React", "TypeScript", "Node.js", "Express", "PostgreSQL"],
    "deployment": "AWS",
    "team_size": 3,
    "existing_users": "10000+"
  },
  "token_budget": 450
}
```

## Consequences

### Benefits
- **Dramatic Token Reduction**: 70%+ reduction in communication overhead
- **Focused Agent Attention**: Agents receive only relevant information
- **Faster Processing**: Reduced context leads to quicker agent responses  
- **Cost Efficiency**: Linear reduction in token-based costs

### Trade-offs
- **Information Loss**: Some nuance may be lost in compression
- **Context Reconstruction**: Agents may need to request additional details
- **Compression Overhead**: CPU time spent analyzing and compressing context
- **Validation Complexity**: Ensuring compression maintains task completeness

### Mitigation Strategies
- **Iterative Refinement**: Allow agents to request specific additional context
- **Domain-Specific Compression**: Tailor compression rules to agent specializations
- **Quality Validation**: Automated checks for essential information preservation
- **Fallback Mechanism**: Expand context if agent requests more information