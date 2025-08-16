# Pattern: User Confirmation Gates

## Context
The enhanced `/s:specify` command system requires user control and transparency at critical decision points. Users need the ability to confirm, modify, or override system routing decisions while maintaining workflow momentum.

## Problem
- Users lose control when tasks are automatically delegated without confirmation
- System routing decisions may be suboptimal for user's specific context
- No visibility into why certain delegation choices are made
- Workflow interruptions are jarring when users disagree with system decisions

## Solution
Implement standardized confirmation gates at key decision points that provide:
1. **Clear Context**: What decision is being made and why
2. **Actionable Options**: Specific choices user can make
3. **Reasoning Visibility**: Explanation of system recommendation
4. **Override Capability**: Ability to modify or reject system decision
5. **Timeout Handling**: Graceful defaults when user doesn't respond

### Gate Types

#### Delegation Gate
Appears when system recommends routing task to sub-agent.
```
┌─ Delegation Recommendation ─────────────────────────┐
│ Task: "Design authentication system"               │
│ Recommended: Level 3 → the-architect               │
│ Reason: Multi-system integration, security domain  │
│                                                     │
│ Options:                                           │
│ [C]onfirm delegation                               │
│ [M]odify - choose different agent                  │
│ [D]irect - handle without delegation               │
│ [?] More details                                   │
└─────────────────────────────────────────────────────┘
```

#### Response Gate
Appears after sub-agent completes work, before accepting output.
```
┌─ Sub-Agent Response Review ─────────────────────────┐
│ Agent: the-architect                               │
│ Output: System Design Document (2,847 tokens)      │
│ Quality: High confidence, complete specification   │
│                                                     │
│ Options:                                           │
│ [A]ccept response                                  │
│ [R]evise - request changes                         │
│ [E]xtend - add more detail                         │
│ [V]iew - see full response                         │
└─────────────────────────────────────────────────────┘
```

#### Transition Gate
Appears when moving between major workflow phases.
```
┌─ Workflow Transition ───────────────────────────────┐
│ Completed: Requirements Analysis                   │
│ Next Phase: Solution Design                        │
│ Estimated time: 5-8 minutes                       │
│                                                     │
│ Options:                                           │
│ [C]ontinue to next phase                          │
│ [P]ause - save session for later                  │
│ [M]odify - change approach                         │
│ [S]ummary - review current progress               │
└─────────────────────────────────────────────────────┘
```

## Implementation Example

### Gate Interface
```go
type ConfirmationGate interface {
    Present(context GateContext) (GateDecision, error)
    GetTimeout() time.Duration
    GetDefaultAction() GateDecision
}

type GateContext struct {
    GateType     string
    TaskSummary  string
    Recommendation string
    Reasoning    string
    Options      []GateOption
    SessionData  map[string]interface{}
}

type GateDecision struct {
    Action   string
    Modified bool
    Context  map[string]interface{}
}
```

### BubbleTea Implementation
```go
type confirmationModel struct {
    context     GateContext
    options     []GateOption
    cursor      int
    timer       *time.Timer
    timeLeft    time.Duration
    showDetails bool
}

func (m confirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "c", "C":
            return m, tea.Quit // Confirm
        case "m", "M":
            return m.showModifyOptions(), nil
        case "?":
            m.showDetails = !m.showDetails
            return m, nil
        case "q", "ctrl+c":
            return m, tea.Quit
        }
    case timeoutMsg:
        return m, tea.Quit // Use default action
    }
    return m, nil
}
```

## Consequences

### Benefits
- **User Control**: Users maintain agency over their workflow
- **Transparency**: Clear reasoning for all system decisions
- **Flexibility**: Easy to override or modify system recommendations
- **Workflow Continuity**: Gates integrate smoothly without jarring interruptions

### Trade-offs
- **Additional Latency**: Each gate adds user decision time to workflow
- **Cognitive Load**: Users must understand and evaluate system recommendations
- **Implementation Complexity**: Requires consistent UI patterns across all gate types
- **Timeout Handling**: Must gracefully handle non-responsive users

### Mitigation Strategies
- **Smart Defaults**: Learn from user patterns to reduce unnecessary confirmations
- **Progressive Disclosure**: Show details only when requested
- **Batch Operations**: Group related decisions into single confirmation
- **User Preferences**: Allow users to configure auto-confirmation for trusted scenarios