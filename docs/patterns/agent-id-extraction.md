# Pattern: Agent ID Extraction and Management

## Context
Agent instances in the-startup system need consistent identification across multiple interactions to enable context persistence and instance tracking. AgentIds must be extractable from prompts and have fallback generation when not explicitly provided.

## Problem
- Agent prompts may or may not contain explicit AgentId declarations
- System needs unique, consistent identifiers for agent instances
- Generated IDs must be deterministic to avoid context fragmentation
- Multiple agents of same type need distinct identities

## Solution

### Primary: Regex Extraction
```go
var agentIDRegex = regexp.MustCompile(`\bAgentId\s*:\s*([^\s,]+)`)

func ExtractAgentID(prompt string) string {
    matches := agentIDRegex.FindStringSubmatch(prompt)
    if len(matches) > 1 {
        return matches[1]
    }
    return ""
}
```

### Fallback: Deterministic Generation
```go
func GenerateAgentID(agentType, prompt string, timestamp time.Time) string {
    // Create deterministic hash from prompt content
    h := sha256.New()
    h.Write([]byte(agentType + prompt + timestamp.Format("2006-01-02-15"))) // Hour precision
    hash := hex.EncodeToString(h.Sum(nil))[:8]
    
    return fmt.Sprintf("%s-%s", agentType, hash)
}
```

### Complete Implementation
```go
func ExtractOrGenerateAgentID(agentType, prompt string) string {
    // Try extraction first
    if id := ExtractAgentID(prompt); id != "" {
        return id
    }
    
    // Generate fallback
    return GenerateAgentID(agentType, prompt, time.Now())
}
```

## Implementation Example

### Agent Template Usage
```markdown
## Session Context
AgentId: arch-001
SessionId: dev-20240112-1503

I am the-architect agent instance "arch-001" working on system design...
```

### Hook Processing
```go
// In ProcessToolCall function
agentID := ExtractOrGenerateAgentID(subagentType, prompt)
hookData := &HookData{
    AgentID:   agentID,
    AgentType: subagentType,
    // ... other fields
}
```

## Benefits
- **Consistency**: Same agent instance maintains identity across interactions
- **Flexibility**: Works with or without explicit AgentIds
- **Determinism**: Generated IDs are reproducible for same inputs
- **Isolation**: Different agents get unique contexts

## Trade-offs
- **Generated ID Readability**: Fallback IDs are not human-friendly
- **Hash Collision Risk**: Very low but theoretically possible
- **Time Dependency**: Generated IDs depend on hour precision

## Validation
```go
func ValidateAgentID(id string) bool {
    // AgentID must be alphanumeric with hyphens, 3-50 chars
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9-]{3,50}$`, id)
    return matched
}
```