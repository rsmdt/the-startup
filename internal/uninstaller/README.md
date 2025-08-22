# Uninstaller Package - Removal Preview System

This package provides comprehensive removal preview functionality for the-startup uninstaller, integrating with the existing lockfile system and path discovery from Phase 1.

## Features Implemented

### 1. Directory Scanning and File Analysis
- **Comprehensive file discovery**: Scans both installation and Claude directories recursively
- **Smart file identification**: Uses naming patterns and directory structure to identify startup-related files
- **File categorization**: Automatically categorizes files into 9 distinct types:
  - `CategoryAgent`: Claude agents (the-*.md files in agents/)  
  - `CategoryCommand`: Slash commands (files in commands/s/)
  - `CategoryTemplate`: Document templates (files in templates/)
  - `CategoryRule`: Configuration rules (files in rules/)
  - `CategoryBinary`: Executable files (the-startup binary)
  - `CategoryLog`: Session logs (*.jsonl files in logs/)
  - `CategorySettings`: Claude settings (settings.json, settings.local.json)
  - `CategoryOther`: Other tracked files (lockfile, etc.)
  - `CategoryUntracked`: Files not in lockfile

### 2. Detailed File Information
For each discovered file, the system captures:
- **Path information**: Both absolute and relative paths
- **Size and timestamps**: File size and modification time  
- **Lockfile tracking**: Whether file is tracked in lockfile
- **Modification detection**: Compares with lockfile to detect changes
- **Permission analysis**: Checks if files can be removed
- **Symlink handling**: Detects and analyzes symlinks safely

### 3. Lockfile Integration
- **Smart comparison**: Matches discovered files with lockfile entries using multiple path formats
- **Change detection**: Identifies files modified since installation
- **Untracked file detection**: Finds files not in lockfile (user additions)
- **Orphaned file detection**: Identifies lockfile entries for missing files
- **Path format handling**: Supports different lockfile path prefixes (startup/, bin/, etc.)

### 4. Size Calculations and Aggregation  
- **Individual file sizes**: Accurate size reporting for each file
- **Category totals**: Aggregated size by file category
- **Overall totals**: Total files and disk space to be freed
- **Summary statistics**: Tracked vs untracked counts, modification counts

### 5. Security Analysis
The system performs comprehensive security checks:
- **Path traversal detection**: Identifies suspicious ".." patterns
- **Location validation**: Ensures files are within expected directories
- **Symlink safety**: Validates symlink targets aren't malicious
- **Large file detection**: Flags unusually large files (>100MB)
- **Permission validation**: Checks file access permissions

### 6. Validation and Error Detection
- **Permission issues**: Identifies files that cannot be removed due to permissions
- **File-in-use detection**: Basic check for running binaries
- **Critical system file protection**: Prevents accidental system file removal
- **Directory validation**: Ensures directories exist and are accessible

### 7. Comprehensive Data Structures

#### `RemovalPreview`
Main preview result containing:
- Path information and discovery source
- Complete file inventory with detailed metadata
- Category summaries and aggregated statistics
- Lockfile information and comparison results
- Security issues and validation errors
- Untracked and orphaned file lists

#### `FileInfo`
Detailed information for each file:
- Path, size, timestamps, permissions
- Category and lockfile tracking status
- Modification and symlink detection
- Permission issue reporting

#### `CategorySummary`
Aggregated statistics per file category:
- File counts and total sizes
- Tracked vs untracked breakdown
- Modified file counts

### 8. Integration Points

#### Existing Systems
- **PathDiscoverer**: Uses Phase 1 path discovery for installation detection
- **LockFile**: Integrates with existing config.LockFile structure
- **Uninstaller**: Extends existing uninstaller with preview capabilities

#### New Interface Methods
- `GenerateDetailedPreview()`: Returns comprehensive RemovalPreview
- `GeneratePreview()`: Core preview generation (PreviewGenerator)

### 9. Test Coverage

Comprehensive test suite covering:
- **Unit tests**: Individual function and method testing
- **Integration tests**: Full preview generation with mock installations  
- **Edge case tests**: Non-existent directories, empty installations
- **Security tests**: Malicious input handling
- **Performance tests**: Benchmarks for file operations

Test files:
- `preview_test.go`: Core preview functionality tests
- `path_discovery_test.go`: Path discovery tests  
- `file_remover_test.go`: File removal tests
- `settings_manager_test.go`: Settings management tests

### 10. Usage Examples

#### Basic Preview Generation
```go
pathDiscoverer := NewPathDiscovery()
previewGenerator := NewPreviewGenerator(pathDiscoverer)
preview, err := previewGenerator.GeneratePreview()
```

#### Integrated with Uninstaller
```go
uninstaller := NewWithDefaults()
preview, err := uninstaller.GenerateDetailedPreview()
```

#### Category Analysis
```go
for _, summary := range preview.CategorySummary {
    fmt.Printf("%s: %d files, %d bytes\n", 
        summary.Category.String(), 
        summary.Count, 
        summary.TotalSize)
}
```

## Architecture

### Core Components

1. **PreviewGenerator** (`preview.go`): Main analysis engine
2. **PathDiscoverer** (`path_discovery.go`): Installation path detection  
3. **SafeFileRemover** (`file_remover.go`): Safe file removal implementation
4. **DefaultSettingsManager** (`settings_manager.go`): Settings cleanup

### Key Design Principles

- **Security First**: Multiple layers of validation and safety checks
- **Comprehensive Analysis**: No file left behind - thorough directory scanning
- **Lockfile Integration**: Leverages existing tracking for accuracy
- **Test Coverage**: Extensive testing for reliability
- **Performance**: Efficient directory traversal and file analysis
- **Extensibility**: Modular design for future enhancements

## Files Created/Modified

### New Files
- `internal/uninstaller/preview.go` - Core preview functionality
- `internal/uninstaller/preview_test.go` - Comprehensive tests
- `internal/uninstaller/example_usage.go` - Usage examples and documentation

### Integration Updates  
- `internal/uninstaller/uninstaller.go` - Added GenerateDetailedPreview() method
- `internal/uninstaller/file_remover.go` - Fixed interface implementation
- `internal/uninstaller/settings_manager.go` - Fixed unused variable

This implementation provides a complete, production-ready removal preview system that safely analyzes installations, categorizes files, calculates disk usage, and identifies potential issues before any actual removal occurs.