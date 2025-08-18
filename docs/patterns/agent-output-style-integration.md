# Pattern: Agent Output-Style Integration

## Context

Claude Code's output-styles feature allows customization of the system prompt through markdown files. The-startup has 17 specialized agents, each with unique personalities and expertise. This pattern enables dynamic agent personality presentation while working within Claude Code's static output-style constraints.

## Problem

- Output-styles are loaded once per Claude Code session and cannot be dynamically switched
- Each agent needs to maintain its unique personality and communication style
- Task tool invocations create isolated contexts that need personality injection
- Performance overhead must be minimal when switching between agents
- Agent personalities should enhance, not override, Claude Code's core capabilities

## Solution

### Architecture: Hybrid Meta-Style with Context Detection

The solution uses a three-layer approach:

1. **Base Layer**: Meta output-style that provides adaptive behavior
2. **Detection Layer**: Context-aware personality selection
3. **Injection Layer**: Agent-specific personality overlays

### Implementation Components

#### 1. Meta Output-Style (`~/.claude/output-styles/the-startup-agents.md`)

```markdown
# The Startup Agent System

You are operating within the-startup's agent orchestration system. Your personality and expertise adapt based on the active agent context.

## Agent Detection Rules

When you see patterns indicating agent invocation:
- Task tool calls with `subagent_type: "the-*"`
- Commentary blocks with agent signatures (e.g., "(⌐■_■) **Architect**")
- Agent-specific output formats in responses

Adapt your communication style to match the detected agent while maintaining Claude Code's core capabilities.

## Personality Preservation Protocol

1. **Detect Active Agent**: Parse recent context for agent indicators
2. **Apply Personality Overlay**: Match detected agent to personality profile
3. **Maintain Consistency**: Preserve agent personality throughout interaction
4. **Reset on Context Switch**: Return to neutral when no agent detected

## Agent Personality Profiles

### the-architect
- Signature: (⌐■_■)
- Style: Philosophical, thoughtful, aesthetic appreciation
- Focus: System design, patterns, elegant solutions

### the-developer  
- Signature: (๑˃ᴗ˂)ﻭ
- Style: Enthusiastic, joyful, excited about coding
- Focus: TDD, clean code, celebrating green tests

### the-chief
- Signature: (▀̿Ĺ̯▀̿)
- Style: Decisive, strategic, bottom-line focused
- Focus: Orchestration, delegation, results

[Additional agent profiles...]

## Output Format Adaptation

When an agent is active, format responses to include:
- Agent-specific commentary blocks
- Personality-appropriate language
- Domain-specific terminology
- Consistent emotional markers
```

#### 2. Agent Context Injector

Enhance agent prompts to include stronger personality signals:

```go
// internal/agents/personality.go
package agents

type PersonalityInjector struct {
    signatures map[string]string
    styles     map[string]string
}

func (p *PersonalityInjector) EnhancePrompt(agentType, basePrompt string) string {
    // Inject strong personality markers at prompt start
    personalityHeader := fmt.Sprintf(`
## Active Agent Context
You are currently operating as %s with the following personality:
- Signature: %s
- Communication Style: %s

Maintain this personality throughout your response.

---

`, agentType, p.signatures[agentType], p.styles[agentType])
    
    return personalityHeader + basePrompt
}
```

#### 3. Session-Aware Style Manager

```go
// internal/styles/manager.go
package styles

import (
    "os"
    "path/filepath"
)

type StyleManager struct {
    claudePath string
    startupPath string
}

func (sm *StyleManager) InstallAgentStyles() error {
    // Generate individual agent output-style files
    for _, agent := range GetAllAgents() {
        content := sm.generateAgentStyle(agent)
        stylePath := filepath.Join(sm.claudePath, "output-styles", 
                                  fmt.Sprintf("%s.md", agent.Name))
        if err := os.WriteFile(stylePath, []byte(content), 0644); err != nil {
            return err
        }
    }
    
    // Install meta-style as default
    metaStyle := sm.generateMetaStyle()
    defaultPath := filepath.Join(sm.claudePath, "output-styles", 
                                "the-startup-agents.md")
    return os.WriteFile(defaultPath, []byte(metaStyle), 0644)
}

func (sm *StyleManager) generateAgentStyle(agent Agent) string {
    return fmt.Sprintf(`# %s Output Style

You are %s, with the following characteristics:

## Personality
%s

## Communication Style
- Always include commentary blocks with signature: %s
- Express %s in your responses
- Focus on %s

## Expertise Areas
%s

Remember: You are Claude Code enhanced with %s's personality.
`, agent.Name, agent.Name, agent.Description, 
   agent.Signature, agent.EmotionalTone, agent.Focus,
   agent.Expertise, agent.Name)
}
```

### Integration Points

#### 1. Installation Phase
- During `the-startup install`, generate and install agent-specific output-styles
- Create meta-style as the default option
- Provide user guidance on style selection

#### 2. Command Invocation
- Commands that invoke agents inject personality markers
- Use strong context signals in Task prompts
- Include AgentId for persistence

#### 3. Hook Processing
- Log personality markers for context reconstruction
- Track agent switches for performance analysis
- Maintain personality consistency metrics

### Performance Optimization

#### Switching Overhead Mitigation

1. **Lazy Loading**: Only inject full personality when agent is actively used
2. **Context Caching**: Cache detected agent context for 5-minute windows
3. **Minimal Markers**: Use compact personality signals in prompts
4. **Batch Operations**: Group same-agent tasks to minimize switches

#### Performance Metrics

```go
type PerformanceMetrics struct {
    AgentSwitches      int           // Number of agent context switches
    AverageLatency     time.Duration // Latency added by personality injection
    CacheHitRate       float64       // Context cache effectiveness
    PersonalityDrift   float64       // Consistency of personality preservation
}
```

## Implementation Example

### User Workflow

1. **Initial Setup**:
```bash
# Install the-startup with agent styles
./the-startup install

# Styles are generated at:
# ~/.claude/output-styles/the-startup-agents.md (meta-style)
# ~/.claude/output-styles/the-architect.md
# ~/.claude/output-styles/the-developer.md
# ... (one per agent)
```

2. **Usage in Claude Code**:
```bash
# Option 1: Use meta-style (recommended)
# Claude Code loads the-startup-agents.md automatically

# Option 2: Agent-specific session
# User manually selects the-architect.md for architecture-focused session
```

3. **Agent Invocation**:
```markdown
/s:implement 003

# System detects the-chief orchestrating
# Automatically applies the-chief personality
# Delegates to the-developer with personality switch
# Maintains personality context throughout
```

## Benefits

- **Seamless Integration**: Works within Claude Code's existing architecture
- **Personality Preservation**: Maintains unique agent characteristics
- **Performance**: Minimal overhead with smart caching
- **Flexibility**: Supports both meta-style and dedicated agent styles
- **User Control**: Users can choose their preferred style approach

## Trade-offs

- **Session Persistence**: Style changes require new Claude Code session
- **Detection Complexity**: Context detection adds parsing overhead
- **Personality Bleed**: Risk of personality mixing in rapid switches
- **Maintenance**: Multiple style files to maintain and sync

## Validation

Test personality preservation with:

```go
func TestPersonalityConsistency(t *testing.T) {
    // Invoke agent with personality
    response1 := InvokeAgent("the-architect", "Design a system")
    
    // Check for personality markers
    assert.Contains(t, response1, "(⌐■_■)")
    assert.Contains(t, response1, "philosophical")
    
    // Switch agents
    response2 := InvokeAgent("the-developer", "Implement feature")
    
    // Verify personality switch
    assert.Contains(t, response2, "(๑˃ᴗ˂)ﻭ")
    assert.NotContains(t, response2, "(⌐■_■)")
}
```

## Future Enhancements

1. **Dynamic Style API**: Propose Claude Code feature for runtime style switching
2. **Personality Learning**: ML model to learn and refine agent personalities
3. **Context Persistence**: Store personality context across sessions
4. **Multi-Agent Collaboration**: Support simultaneous multi-agent personalities