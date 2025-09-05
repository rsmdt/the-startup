# Solution Design Document: Tool Usage Statistics Collection

## Executive Summary

Complete replacement of the existing log package with a simplified metrics collection system that tracks ALL Claude Code tool invocations via hooks, storing them in daily JSONL files for analysis.

## System Architecture

### Component Overview

```
┌─────────────────────────────────────────────────┐
│                 Claude Code                      │
├─────────────────────────────────────────────────┤
│  PreToolUse Hook  │  PostToolUse Hook           │
└────────┬──────────┴──────────┬──────────────────┘
         │                     │
         ↓                     ↓
┌─────────────────────────────────────────────────┐
│           the-startup log command                │
├─────────────────────────────────────────────────┤
│  --pre handler    │  --post handler             │
│  (create entry)   │  (update entry)             │
└────────┬──────────┴──────────┬──────────────────┘
         │                     │
         ↓                     ↓
┌─────────────────────────────────────────────────┐
│          {{STARTUP_PATH}}/logs/                  │
│         20250903.jsonl (daily files)             │
└─────────────────────────────────────────────────┘
         │
         ↓
┌─────────────────────────────────────────────────┐
│        the-startup log [analysis]                │
│   (summary, tools, timeline, errors)             │
└─────────────────────────────────────────────────┘
```

### Package Structure

```
internal/
├── log/               # Completely replaced implementation
│   ├── collector.go   # Hook data processing
│   ├── writer.go      # JSONL file operations  
│   ├── reader.go      # Metrics reading/streaming
│   ├── aggregator.go  # Statistics calculation
│   ├── types.go       # Data structures
│   └── display.go     # Terminal output formatting
│
cmd/
├── log.go            # Updated with new subcommands
└── install.go        # Updated to configure hooks
```

## Data Model

### MetricsEntry Structure

Based on actual Claude Code hook payloads, we store a simplified structure:

```go
type MetricsEntry struct {
    // Core fields from hook payload
    ToolName       string          `json:"tool_name"`
    SessionID      string          `json:"session_id"`
    TranscriptPath string          `json:"transcript_path,omitempty"`
    CWD            string          `json:"cwd,omitempty"`
    HookEvent      string          `json:"hook_event"` // PreToolUse or PostToolUse
    
    // Timing (added by us)
    Timestamp      time.Time       `json:"timestamp"`
    DurationMs     *int64          `json:"duration_ms,omitempty"` // Calculated on PostToolUse
    
    // Tool data (stored as-is from hooks)
    ToolInput      json.RawMessage `json:"tool_input"`
    ToolOutput     json.RawMessage `json:"tool_output,omitempty"` // From PostToolUse
    
    // Correlation
    ToolID         string          `json:"tool_id"` // Generated to match pre/post events
    
    // Status (derived from output)
    Success        *bool           `json:"success,omitempty"` // Extracted from tool_output if available
}
```

Note: We keep the structure minimal and close to what Claude Code actually sends, avoiding over-engineering.

### Aggregated Statistics Structure

```go
type ToolStatistics struct {
    Name            string
    TotalCalls      int
    SuccessCount    int
    FailureCount    int
    TotalDurationMs int64
    AvgDurationMs   float64
    MinDurationMs   int64
    MaxDurationMs   int64
    ErrorTypes      map[string]int
    LastUsed        time.Time
}

type MetricsSummary struct {
    Period          TimePeriod
    TotalCalls      int
    UniqueSessions  int
    SuccessRate     float64
    ToolStats       map[string]*ToolStatistics
    TopErrors       []ErrorPattern
    HourlyActivity  []HourlyStats
}
```

## Implementation Details

### Hook Processing Flow

#### Unified Hook Handler
```go
func HandleHook(input io.Reader) error {
    // 1. Parse hook data from stdin
    var hookData HookPayload
    if err := json.NewDecoder(input).Decode(&hookData); err != nil {
        return nil // Silent failure
    }
    
    // 2. Check which event type
    switch hookData.HookEventName {
    case "PreToolUse":
        return handlePreToolUse(hookData)
    case "PostToolUse":
        return handlePostToolUse(hookData)
    default:
        return nil // Ignore other events
    }
}

func handlePreToolUse(hookData HookPayload) error {
    // Generate unique tool ID for correlation
    toolID := fmt.Sprintf("%s_%s_%s", 
        strings.ToLower(hookData.ToolName),
        time.Now().Format("20060102T150405Z"),
        hookData.SessionID[:8])
    
    entry := MetricsEntry{
        ToolID:      toolID,
        ToolName:    hookData.ToolName,
        HookEvent:   "PreToolUse",
        Timestamp:   time.Now().UTC(),
        SessionID:   hookData.SessionID,
        ToolInput:   hookData.ToolInput,
        CWD:         hookData.CWD,
    }
    
    return appendMetrics(entry)
}

func handlePostToolUse(hookData HookPayload) error {
    // For PostToolUse, we create a separate entry
    // (simpler than trying to update the PreToolUse entry)
    entry := MetricsEntry{
        ToolID:      generateToolID(hookData), // Same ID as PreToolUse
        ToolName:    hookData.ToolName,
        HookEvent:   "PostToolUse",
        Timestamp:   time.Now().UTC(),
        SessionID:   hookData.SessionID,
        ToolInput:   hookData.ToolInput,
        ToolOutput:  hookData.ToolResponse,
        Success:     extractSuccess(hookData.ToolResponse),
    }
    
    return appendMetrics(entry)
}
```

### File Operations

#### Daily File Writing
```go
func appendMetrics(entry MetricsEntry) error {
    // Generate filename based on current date: 20250903.jsonl
    today := time.Now().Format("20060102")
    metricsPath := filepath.Join(getStartupPath(), "logs", today+".jsonl")
    
    // Ensure directory exists
    os.MkdirAll(filepath.Dir(metricsPath), 0755)
    
    // Open file for append
    file, err := os.OpenFile(metricsPath, 
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil // Silent failure
    }
    defer file.Close()
    
    // Use file lock for concurrent writes
    syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
    defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
    
    // Write JSON line
    encoder := json.NewEncoder(file)
    encoder.Encode(entry)
    
    return nil
}
```

#### Streaming Reader for Analysis
```go
func StreamMetrics(filter MetricsFilter) (<-chan MetricsEntry, error) {
    logsDir := filepath.Join(getStartupPath(), "logs")
    
    // Get list of log files to process based on date range
    files, err := getLogFiles(logsDir, filter.StartDate, filter.EndDate)
    if err != nil {
        return nil, err
    }
    
    entries := make(chan MetricsEntry, 100)
    
    go func() {
        defer close(entries)
        
        // Process each daily file
        for _, filename := range files {
            file, err := os.Open(filename)
            if err != nil {
                continue // Skip unreadable files
            }
            
            scanner := bufio.NewScanner(file)
            for scanner.Scan() {
                var entry MetricsEntry
                if json.Unmarshal(scanner.Bytes(), &entry) != nil {
                    continue // Skip malformed lines
                }
                
                if filter.Matches(entry) {
                    entries <- entry
                }
            }
            file.Close()
        }
    }()
    
    return entries, nil
}

func getLogFiles(dir string, start, end time.Time) ([]string, error) {
    var files []string
    
    // Generate list of dates to check
    for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
        filename := filepath.Join(dir, d.Format("20060102")+".jsonl")
        if _, err := os.Stat(filename); err == nil {
            files = append(files, filename)
        }
    }
    
    return files, nil
}
```

### Metrics Aggregation

```go
func AggregateMetrics(entries <-chan MetricsEntry) *MetricsSummary {
    summary := &MetricsSummary{
        ToolStats: make(map[string]*ToolStatistics),
    }
    
    for entry := range entries {
        // Update tool statistics
        stats := summary.ToolStats[entry.ToolName]
        if stats == nil {
            stats = &ToolStatistics{Name: entry.ToolName}
            summary.ToolStats[entry.ToolName] = stats
        }
        
        stats.TotalCalls++
        if entry.Status == "success" {
            stats.SuccessCount++
        } else if entry.Status == "failure" {
            stats.FailureCount++
            stats.ErrorTypes[entry.Error]++
        }
        
        if entry.DurationMs != nil {
            stats.TotalDurationMs += *entry.DurationMs
            // Update min/max
        }
        
        summary.TotalCalls++
    }
    
    // Calculate derived metrics
    for _, stats := range summary.ToolStats {
        if stats.TotalCalls > 0 {
            stats.AvgDurationMs = float64(stats.TotalDurationMs) / 
                                float64(stats.TotalCalls)
        }
    }
    
    return summary
}
```

### Display Formatting

```go
func DisplaySummary(summary *MetricsSummary) {
    // Header
    fmt.Println(strings.Repeat("═", 75))
    fmt.Printf("                    Metrics Summary - %s\n", 
               summary.Period.String())
    fmt.Println(strings.Repeat("═", 75))
    
    // Overview
    fmt.Printf("\nOverview:\n")
    fmt.Printf("  Total Invocations:     %d\n", summary.TotalCalls)
    fmt.Printf("  Success Rate:          %.1f%%\n", 
               summary.SuccessRate*100)
    
    // Top tools table
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprintln(w, "\nTop Tools by Usage:")
    fmt.Fprintln(w, "Tool\tCount\tSuccess\tAvg Time\tErrors")
    fmt.Fprintln(w, "────\t─────\t───────\t────────\t──────")
    
    for _, stats := range getTopTools(summary.ToolStats, 10) {
        successRate := float64(stats.SuccessCount) / 
                      float64(stats.TotalCalls) * 100
        fmt.Fprintf(w, "%s\t%d\t%.1f%%\t%.1fs\t%d\n",
            stats.Name,
            stats.TotalCalls,
            successRate,
            stats.AvgDurationMs/1000,
            stats.FailureCount)
    }
    w.Flush()
}
```

## Hook Configuration

### Settings Template Update

```json
{
  "permissions": {
    "additionalDirectories": ["{{STARTUP_PATH}}"]
  },
  "statusLine": {
    "type": "command",
    "command": "{{STARTUP_PATH}}/bin/the-startup statusline"
  },
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "{{STARTUP_PATH}}/bin/the-startup log"
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "{{STARTUP_PATH}}/bin/the-startup log"
          }
        ]
      }
    ]
  }
}
```

## Migration Plan

### Phase 1: Remove Old Log Package
1. Delete `internal/log/` directory completely
2. Remove `cmd/log.go` and related tests
3. Remove log-related configuration from installation

### Phase 2: Implement New Metrics Package
1. Create `internal/metrics/` with new implementation
2. Add `cmd/metrics.go` with hook handlers and analysis
3. Update installation to create `logs/` directory

### Phase 3: Update Installation
1. Modify `assets/claude/settings.json` template
2. Update installer to configure hooks
3. Remove old hook configurations

## Command Interface

### Hook Handler
```bash
# Called by Claude Code hooks (auto-detects PreToolUse vs PostToolUse)
the-startup log < hook_data.json
```

### Analysis Commands
```bash
# Summary view (default)
the-startup log

# Time-based filtering  
the-startup log --since 24h
the-startup log --today
the-startup log --yesterday

# Tool-specific analysis
the-startup log tools --top 20
the-startup log tools --tool Edit

# Error analysis
the-startup log errors --since 7d

# Timeline view
the-startup log timeline --group-by hour

# JSON output for scripts
the-startup log --format json --since 1h

# CSV export
the-startup log --format csv --output report.csv
```

## Error Handling

All errors are handled silently to avoid disrupting Claude Code:

```go
func main() {
    // Global panic recovery
    defer func() {
        if r := recover(); r != nil {
            debugLog("Panic: %v", r)
        }
        os.Exit(0) // Always exit 0
    }()
    
    if err := runMetricsCommand(); err != nil {
        debugLog("Error: %v", err)
        os.Exit(0) // Silent failure
    }
    
    os.Exit(0)
}

func debugLog(format string, args ...interface{}) {
    if os.Getenv("DEBUG_HOOKS") != "" {
        fmt.Fprintf(os.Stderr, format+"\n", args...)
    }
}
```

## Testing Strategy

1. **Unit Tests**: Test each component in isolation
2. **Integration Tests**: Test hook processing end-to-end
3. **Manual Testing**: Use test JSON payloads
4. **Performance Testing**: Ensure fast hook processing

## Security Considerations

1. **Input Validation**: Validate all JSON input
2. **File Permissions**: Restrict metrics file to user only
3. **Path Traversal**: Validate all file paths
4. **Command Injection**: Never execute user input
5. **Silent Failures**: Don't leak errors to Claude Code

## Performance Requirements

- Hook processing: < 100ms per invocation
- File writes: Use OS buffering, no forced sync
- Analysis queries: Stream processing for large files
- Memory usage: Constant memory for streaming

## Success Criteria

1. ✅ All existing log package code removed and replaced
2. ✅ New system captures ALL Claude Code tool invocations
3. ✅ Daily JSONL files created in {{STARTUP_PATH}}/logs/ (e.g., 20250903.jsonl)
4. ✅ Hooks configured automatically during install
5. ✅ Analysis commands provide useful insights
6. ✅ Silent failure mode maintained (always exit 0)
7. ✅ No performance impact on Claude Code
8. ✅ Simple correlation of PreToolUse and PostToolUse events