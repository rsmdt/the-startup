# Pattern: Session-Based Context Organization

## Context
Agent context must be organized efficiently to support multiple concurrent sessions, agent instances, and enable quick context retrieval while maintaining backward compatibility with existing logging.

## Problem
- Multiple Claude Code sessions run simultaneously
- Multiple agent instances within each session need isolated context
- Context files can grow large over time
- Need efficient lookup and reading of recent context
- Must maintain backward compatibility with existing agent-instructions.jsonl

## Solution

### Directory Structure
```
.the-startup/
├── [session-id]/                   # Per-session organization
│   ├── main.jsonl                  # Orchestrator context
│   ├── [agent-id-1].jsonl         # Individual agent contexts
│   ├── [agent-id-2].jsonl         
│   └── agent-instructions.jsonl   # Backward compatibility
├── [session-id-2]/
│   └── ...
├── all-agent-instructions.jsonl   # Global log (existing)
└── templates/                      # Templates (existing)
```

### File Organization Logic
```go
type SessionContextManager struct {
    BaseDir   string
    SessionID string
}

func (scm *SessionContextManager) GetAgentContextPath(agentID string) string {
    return filepath.Join(scm.BaseDir, scm.SessionID, agentID+".jsonl")
}

func (scm *SessionContextManager) GetMainContextPath() string {
    return filepath.Join(scm.BaseDir, scm.SessionID, "main.jsonl")
}

func (scm *SessionContextManager) GetCompatibilityPath() string {
    return filepath.Join(scm.BaseDir, scm.SessionID, "agent-instructions.jsonl")
}
```

### Context Writing Strategy
```go
func WriteMultipleContexts(sessionID, agentID string, data *HookData) error {
    scm := &SessionContextManager{
        BaseDir:   GetStartupDir(),
        SessionID: sessionID,
    }
    
    // Create session directory if needed
    sessionDir := filepath.Join(scm.BaseDir, sessionID)
    if err := os.MkdirAll(sessionDir, 0755); err != nil {
        return err
    }
    
    // Write to agent-specific file
    if agentID != "" {
        if err := appendJSONL(scm.GetAgentContextPath(agentID), data); err != nil {
            return err
        }
    }
    
    // Write to main context if orchestrator
    if data.AgentType == "the-chief" || data.AgentType == "orchestrator" {
        if err := appendJSONL(scm.GetMainContextPath(), data); err != nil {
            return err
        }
    }
    
    // Write to compatibility file
    if err := appendJSONL(scm.GetCompatibilityPath(), data); err != nil {
        return err
    }
    
    // Write to global log (existing functionality)
    globalPath := filepath.Join(scm.BaseDir, "all-agent-instructions.jsonl")
    return appendJSONL(globalPath, data)
}
```

## Implementation Example

### Context Lookup
```go
func FindAgentContext(agentID string, sessionID string) (string, error) {
    baseDir := GetStartupDir()
    
    // If sessionID not provided, find latest
    if sessionID == "" {
        sessionID = FindLatestSessionWithAgent(baseDir, agentID)
        if sessionID == "" {
            return "", fmt.Errorf("no context found for agent %s", agentID)
        }
    }
    
    contextPath := filepath.Join(baseDir, sessionID, agentID+".jsonl")
    if !fileExists(contextPath) {
        return "", fmt.Errorf("context file not found: %s", contextPath)
    }
    
    return contextPath, nil
}
```

### Efficient Context Reading
```go
func ReadRecentContext(contextPath string, maxLines int) ([]*HookData, error) {
    lines, err := ReadLastNLines(contextPath, maxLines)
    if err != nil {
        return nil, err
    }
    
    var contexts []*HookData
    for _, line := range lines {
        var hookData HookData
        if err := json.Unmarshal([]byte(line), &hookData); err != nil {
            continue // Skip corrupted lines
        }
        contexts = append(contexts, &hookData)
    }
    
    // Reverse to get chronological order
    for i, j := 0, len(contexts)-1; i < j; i, j = i+1, j-1 {
        contexts[i], contexts[j] = contexts[j], contexts[i]
    }
    
    return contexts, nil
}
```

### Session Discovery
```go
func FindLatestSessionWithAgent(baseDir, agentID string) string {
    entries, err := os.ReadDir(baseDir)
    if err != nil {
        return ""
    }
    
    var latestSession string
    var latestModTime int64
    
    for _, entry := range entries {
        if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "dev-") {
            continue
        }
        
        // Check if this session has the agent
        agentFile := filepath.Join(baseDir, entry.Name(), agentID+".jsonl")
        if !fileExists(agentFile) {
            continue
        }
        
        info, err := entry.Info()
        if err != nil {
            continue
        }
        
        if info.ModTime().Unix() > latestModTime {
            latestModTime = info.ModTime().Unix()
            latestSession = entry.Name()
        }
    }
    
    return latestSession
}
```

## Benefits
- **Isolation**: Each session's contexts are contained
- **Performance**: Quick agent-specific lookups without scanning large files
- **Scalability**: Sessions can be archived or cleaned up independently
- **Compatibility**: Maintains existing file structure expectations

## Trade-offs
- **File Count**: More files created per session
- **Disk Usage**: Some duplication for backward compatibility
- **Complexity**: More complex file management logic

## Cleanup Strategy
```go
func ArchiveOldSessions(baseDir string, keepDays int) error {
    cutoff := time.Now().AddDate(0, 0, -keepDays)
    
    entries, err := os.ReadDir(baseDir)
    if err != nil {
        return err
    }
    
    for _, entry := range entries {
        if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "dev-") {
            continue
        }
        
        info, err := entry.Info()
        if err != nil {
            continue
        }
        
        if info.ModTime().Before(cutoff) {
            sessionPath := filepath.Join(baseDir, entry.Name())
            archivePath := filepath.Join(baseDir, "archive", entry.Name()+".tar.gz")
            
            if err := ArchiveDirectory(sessionPath, archivePath); err != nil {
                continue // Skip failed archives
            }
            
            os.RemoveAll(sessionPath)
        }
    }
    
    return nil
}
```