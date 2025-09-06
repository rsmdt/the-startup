# Stats Package - Claude Code Log Analysis

## Overview

The `stats` package provides a high-performance streaming JSONL parser for analyzing Claude Code's native log format. It processes log files of 100MB+ efficiently without loading them entirely into memory.

## Features

### Core Capabilities
- **Streaming Processing**: Handles 100MB+ files with minimal memory footprint
- **Performance**: Processes at ~91 MB/s for large files
- **Error Resilience**: Gracefully handles malformed JSON lines
- **Comprehensive Parsing**: Supports all Claude Code event types:
  - User messages and tool results
  - Assistant messages with tool invocations
  - System events and hooks
  - Session summaries with token usage

### Advanced Features
- **Filtering**: Time ranges, sessions, tools, status filtering
- **Sampling**: Statistical sampling for large datasets
- **Tool Parameter Parsing**: Extracts parameters for common tools (Bash, Read, Edit, Write, Search)
- **Command Extraction**: Identifies command tags in assistant messages
- **Concurrent Processing**: Thread-safe for parallel file processing

## Usage

### Basic Example

```go
import "github.com/rsmdt/the-startup/internal/stats"

// Create parser
parser := stats.NewJSONLParser()

// Parse a file
entries, errors := parser.ParseFile("/path/to/log.jsonl")

// Handle errors
go func() {
    for err := range errors {
        log.Printf("Parse error: %v", err)
    }
}()

// Process entries
for entry := range entries {
    switch entry.Type {
    case "user":
        fmt.Printf("User: %s\n", entry.User.Text)
    case "assistant":
        fmt.Printf("Assistant: %s\n", entry.Assistant.Text)
    }
}
```

### Filtering

```go
// Filter by time and session
startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
parser.SetFilter(stats.FilterOptions{
    StartTime:    &startTime,
    SessionIDs:   []string{"sess_001"},
    IncludeTools: []string{"Bash", "Read"},
})
```

### Large File Optimization

```go
// Configure for optimal large file processing
parser.SetOptions(stats.ParseOptions{
    BufferSize:     256 * 1024,  // 256KB buffer
    ParseTools:     false,        // Skip tool parsing for speed
    SkipSystemLogs: true,         // Skip system messages
})
```

## Performance

Benchmark results on Apple M2:
- **Small files (1MB)**: 31.24 MB/s with full parsing
- **Large files (10MB+)**: 91.17 MB/s with optimized settings
- **Entry processing**: ~127,000 entries/second

## Testing

The package includes comprehensive tests with **88.4% code coverage**:
- Unit tests for all parser functions
- Malformed JSON handling tests
- Concurrent processing tests
- Performance benchmarks
- Example usage tests

Run tests:
```bash
go test ./internal/stats/...
go test ./internal/stats/... -bench=. -benchmem
```

## Architecture

### Parser Interface
```go
type Parser interface {
    ParseStream(reader io.Reader) (<-chan ClaudeLogEntry, <-chan error)
    ParseFile(filePath string) (<-chan ClaudeLogEntry, <-chan error)
    SetOptions(opts ParseOptions)
    SetFilter(filter FilterOptions)
}
```

### Key Types
- `ClaudeLogEntry`: Root structure for each JSONL line
- `UserMessage`: User input and tool results
- `AssistantMessage`: Claude responses and tool invocations
- `SystemMessage`: System events and hooks
- `SummaryMessage`: Session statistics and token usage
- `ToolUse`: Tool invocation details with parsed parameters

## Implementation Details

### Memory Management
- Uses buffered reader with configurable buffer size
- Streaming channels prevent memory accumulation
- Concurrent goroutines for error handling

### Error Handling
- **Strict Mode**: Stops on first error (for validation)
- **Lenient Mode**: Skips malformed lines and continues (default)
- **Error Collection**: Optional error reporting while continuing

### Filtering Strategy
- Early filtering at parse time for efficiency
- Multiple filter criteria combined with AND logic
- Sampling uses modulo for even distribution