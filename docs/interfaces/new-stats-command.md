# New Stats Command Interface

## Overview
Complete replacement of the CLI stats command with an interactive fullscreen dashboard. No backward compatibility with previous table/JSON/CSV output formats.

## Command Interface

### Primary Command
```bash
the-startup stats
```
**Behavior**: Launches interactive fullscreen dashboard (no flags needed)

### Removed Commands
The following CLI commands are **completely removed**:
```bash
# These no longer exist
the-startup stats tools
the-startup stats agents
the-startup stats commands
the-startup stats sessions
the-startup stats --format json
the-startup stats --format csv
the-startup stats --since 7d
the-startup stats -g
```

### Exit Methods
- Press `q` to quit
- Press `Ctrl+C` to force exit
- Press `ESC` to quit

## Dashboard Layout

### Fullscreen Interface
```
┌─────────────────────────────────────────────────────────┐
│                    Timeline (Top 15%)                    │
├─────────────────────┬──────────────────┬────────────────┤
│   Agent Stats       │   Delegation     │  Co-occurrence │
│     (Left 33%)      │   (Center 33%)   │   (Right 33%)  │
│     (70% height)    │                  │                │
├─────────────────────┴──────────────────┴────────────────┤
│     Filters & Help (Bottom 15%)                         │
└─────────────────────────────────────────────────────────┘
```

### Panel Contents

#### Timeline Panel (Top)
- Horizontal timeline of user/assistant message exchanges
- Color-coded by message type and success/failure
- Shows agent usage over time
- Click or arrow keys to jump to specific time periods

#### Agent Stats Panel (Bottom Left)
- Leaderboard of agents by usage count
- Success rate indicators
- Average duration metrics
- Trend sparklines (simple ASCII)

#### Delegation Panel (Bottom Center)
- Flow diagram showing agent-to-agent delegation patterns
- Frequency and success rate of delegation chains
- Visual representation of workflow patterns

#### Co-occurrence Panel (Bottom Right)
- Matrix or network view of agents used together
- Session correlation analysis
- Agent relationship insights

## Navigation

### Keyboard Controls
- `Tab` / `Shift+Tab`: Navigate between panels
- `Arrow Keys`: Navigate within panels
- `Enter`: Drill down into selected item
- `ESC`: Go back / Exit detail view
- `q`: Quit dashboard
- `r`: Refresh data
- `h` / `?`: Show help overlay
- `1-6`: Quick filter by time range (1w, 1m, 3m, 6m, 1y, all)

### Filter Controls
- `f`: Open filter menu
- Time range: 1w, 1m, 3m, 6m, 1y, all
- Agent filter: Select specific agents
- Success filter: Show only successful/failed interactions
- Session filter: Filter by specific sessions

## Data Loading

### Default Behavior
- Launches with 1-month lookback filter
- Loads all available data within time range
- Shows loading indicator during data processing

### Filter Changes
- Time filter changes trigger full data reload
- Other filters apply to existing loaded data
- Progress indicator for slow operations

### Error Handling
- Missing log files: Show "No data found" with guidance
- Corrupted data: Show parse error count, continue with valid data
- Empty results: Show "No data in time range" with expand suggestion
- Terminal too small: Show minimum size requirement

## Data Export (Optional Future Enhancement)
- Press `e`: Export current view
- Formats: CSV, JSON (for programmatic access)
- Exports filtered dataset only

## Help System
- Press `h` or `?`: Toggle help overlay
- Shows all keyboard shortcuts
- Context-sensitive help based on focused panel
- Quick start guide for new users

## Migration from CLI
**Breaking Change Notice**:
- All previous CLI flags and output formats are removed
- Scripts using `the-startup stats --format json` will break
- No migration path - users must adapt to interactive interface
- Consider this a major version change in user experience

## Technical Interface

### Data Structures
```go
// Main dashboard state
type DashboardState struct {
    CurrentFilter  FilterState
    LoadedData     *StatsData
    FocusedPanel   PanelType
    ViewMode       ViewMode
}

// Simplified filter state
type FilterState struct {
    TimeRange    TimePreset  // 1w, 1m, 3m, 6m, 1y, all
    Agents       []string    // Empty = all agents
    Success      *bool       // nil = all, true = success only, false = failures only
    Sessions     []string    // Empty = all sessions
}

// Panel types
type PanelType int
const (
    TimelinePanel PanelType = iota
    AgentStatsPanel
    DelegationPanel
    CoOccurrencePanel
)
```

### Integration Points
- Reuses all existing `internal/stats/` functions
- Replaces `cmd/stats.go` completely
- Maintains compatibility with JSONL log format
- No external API changes