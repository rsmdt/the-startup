# System Design Document: Bubbletea Installation Architecture

## Overview

This document outlines the redesign of the installation flow from a fragmented mix of `huh` forms and bubbletea components into a cohesive full-screen bubbletea application with proper state management, progressive disclosure, and enhanced user experience.

## Current State Analysis

### Problems Identified
- **Mixed UI paradigms**: Combination of `huh` forms (non-full-screen) with custom TreeSelector (bubbletea)
- **Fragmented state management**: Multiple forms with complex goto loops and scattered state
- **Poor navigation**: ESC handling inconsistencies and navigation breakdowns
- **No progressive disclosure**: Each form replaces previous selections entirely
- **Styling inconsistencies**: Different indentation and styling approaches

### Current Flow
1. Tool selection (huh form)
2. Path selection (huh form)  
3. Component vs file mode selection (huh form)
4. Component selection (huh) OR tree file selection (bubbletea TreeSelector)
5. Update confirmation (huh form)
6. Installation with progress

## Proposed Architecture

### State Machine Design

```go
type InstallerState int

const (
    StateWelcome InstallerState = iota
    StateToolSelection  
    StatePathSelection
    StateComponentMode
    StateComponentSelection  
    StateFileSelection     // Uses TreeSelector
    StateConfirmation
    StateInstalling
    StateComplete
    StateError
)

type InstallerModel struct {
    state InstallerState
    installer *installer.Installer
    
    // User selections (progressive disclosure)
    selectedTool string
    selectedPath string 
    useAdvancedMode bool
    selectedComponents []string
    selectedFiles []string
    
    // UI state
    cursor int
    choices []string
    treeSelector *TreeSelector
    width, height int
    styles Styles
}
```

### Component Architecture

```
┌─────────────────────────────────────────┐
│              InstallerModel             │
│  ┌─────────────────────────────────┐    │
│  │         State Machine           │    │
│  │  ┌─────┐ ┌─────┐ ┌─────┐       │    │
│  │  │Welcome│Tool│Path│...│       │    │
│  │  └─────┘ └─────┘ └─────┘       │    │
│  └─────────────────────────────────┘    │
│  ┌─────────────────────────────────┐    │
│  │    Progressive Disclosure       │    │
│  │   ┌─ Previous Selections ─┐    │    │
│  │   │ Tool: Claude Code      │    │    │
│  │   │ Path: ~/.config/...   │    │    │
│  │   │ Mode: Advanced         │    │    │
│  │   └────────────────────────┘    │    │
│  └─────────────────────────────────┘    │
│  ┌─────────────────────────────────┐    │
│  │        TreeSelector             │    │
│  │     (embedded component)        │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
```

## User Experience Flow

### 1. Welcome State
- ASCII art banner: "The (Agentic) Startup"
- Brief introduction
- Instructions to proceed

### 2. Progressive Disclosure Pattern
Each subsequent screen shows:
```
┌─ The (Agentic) Startup Installation ─┐
│ Tool: Claude Code                     │
│ Path: ~/.config/the-startup          │  
│ Mode: Advanced (individual files)     │
│ Files: 15 selected                    │
└───────────────────────────────────────┘

[Current Selection UI]

Help: ↑↓ navigate • enter: select • esc: back
```

### 3. State Transitions
- **Forward navigation**: Enter confirms selection, advances state
- **Backward navigation**: ESC returns to previous state  
- **Error handling**: Dedicated error state with recovery options
- **Cancellation**: ESC from welcome state exits cleanly

## Technical Implementation

### New Files Required

1. **`internal/ui/ascii.go`**
   ```go
   package ui
   
   const WelcomeBanner = `
   ████████ ██   ██ ███████     
      ██    ██   ██ ██          
      ██    ███████ █████       
      ██    ██   ██ ██          
      ██    ██   ██ ███████     
   
    █████╗  ██████╗ ███████ ███   ██████ 
   ██   ██ ██      ██      ████  ██   ██ 
   ███████ ██  ███ █████   ██ ██ ██   ██  
   ██   ██ ██   ██ ██      ██  ████   ██ 
   ██   ██  ██████ ███████ ██   ███████  
   
   ███████ ████████  █████  ██████  ████████ ██   ██ ██████  
   ██         ██    ██   ██ ██   ██    ██    ██   ██ ██   ██ 
   ███████    ██    ███████ ██████     ██    ██   ██ ██████  
        ██    ██    ██   ██ ██   ██    ██    ██   ██ ██      
   ███████    ██    ██   ██ ██   ██    ██     █████  ██      
   `
   ```

2. **`internal/ui/installer_model.go`**
   - Main bubbletea model implementation
   - State machine logic
   - Message routing and handling
   - Progressive disclosure rendering

3. **`internal/ui/states.go`**  
   - State definitions and transitions
   - State-specific rendering logic
   - Validation helpers

### Modified Files

1. **`internal/ui/theme.go`**
   ```go
   Title: lipgloss.NewStyle().
       Bold(true).
       Foreground(lipgloss.Color("#000000")). // Black text
       Background(lipgloss.Color("#42FF76")). // Bright green background  
       Padding(0, 1).
       MarginTop(1).
       MarginBottom(1)
   ```

2. **`internal/ui/tree_selector.go`**
   - Remove indentation from `getIndent()` method
   - Update to use new title styling
   - Improve integration with main model

3. **`cmd/install.go`**
   - Replace huh form logic with bubbletea model
   - Simplify to single `tea.NewProgram()` call
   - Remove complex goto navigation logic

## Data Flow

### Installation State Synchronization
```go
// User makes selection in UI
model.selectedTool = "claude-code"

// UI configures installer 
model.installer.SetTool(model.selectedTool)

// Continue until ready for installation
if model.state == StateInstalling {
    model.installer.Install() // Uses existing installer logic
}
```

### Progress Integration
Current installer outputs progress via `fmt.Print`. New approach:
- Installer returns progress channel or accepts callback
- UI model receives progress updates via bubbletea commands  
- Real-time progress rendering in StateInstalling

## Integration Points

### TreeSelector Integration
- Embed TreeSelector as sub-component
- Route messages when in StateFileSelection
- Handle completion/cancellation properly
- Maintain progressive disclosure context

### Installer.Installer Integration  
- Preserve existing business logic
- Configure progressively as user makes selections
- Enhance progress reporting for bubbletea compatibility
- Maintain backward compatibility

## Error Handling

### State-Based Error Recovery
- Dedicated StateError for error display
- Context-aware error messages
- Recovery options based on error type
- Graceful degradation when possible

### Validation Strategy
- Validate selections at state transitions
- Show inline validation feedback
- Prevent invalid state transitions
- Clear error recovery paths

## Testing Strategy

### Unit Testing
- State machine transitions
- Progressive disclosure rendering
- TreeSelector integration  
- Error condition handling

### Integration Testing
- Full installation flow testing
- TreeSelector integration
- Installer.Installer compatibility
- Progress reporting accuracy

## Performance Considerations

### Rendering Optimization
- Efficient terminal size handling
- Minimal re-rendering on updates
- Proper viewport management for large content
- TreeSelector performance preservation

### Memory Management
- Clean state transitions
- Proper model cleanup
- TreeSelector resource management
- Progress tracking efficiency

## Security Considerations

### Input Validation
- Path validation and sanitization
- Component selection validation
- File selection boundary checking
- Installation permission verification

### Error Information
- Avoid exposing sensitive system information
- Sanitize error messages for display
- Proper logging without user data exposure

## Migration Strategy

### Backward Compatibility
- Preserve existing installer.Installer interface
- Maintain CLI flag compatibility
- Keep non-interactive mode support
- Preserve installation output format

### Rollout Plan
1. Implement new bubbletea model
2. Update styling and TreeSelector
3. Replace install command logic
4. Test full integration
5. Deploy with fallback option

## Future Enhancements

### Extensibility
- Plugin architecture for additional states
- Themeable UI components
- Configurable ASCII art
- Advanced progress visualization

### User Experience
- Keyboard shortcuts customization
- Mouse support enhancement
- Progress persistence across sessions
- Installation resume capability