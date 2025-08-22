package uninstaller

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rsmdt/the-startup/internal/config"
)

// FileCategory represents different types of files in the installation
type FileCategory int

const (
	CategoryAgent FileCategory = iota
	CategoryCommand
	CategoryTemplate
	CategoryRule
	CategoryBinary
	CategoryLog
	CategorySettings
	CategoryOther
	CategoryUntracked
)

// String returns a human-readable representation of the file category
func (fc FileCategory) String() string {
	switch fc {
	case CategoryAgent:
		return "agent"
	case CategoryCommand:
		return "command"
	case CategoryTemplate:
		return "template"
	case CategoryRule:
		return "rule"
	case CategoryBinary:
		return "binary"
	case CategoryLog:
		return "log"
	case CategorySettings:
		return "settings"
	case CategoryOther:
		return "other"
	case CategoryUntracked:
		return "untracked"
	default:
		return "unknown"
	}
}

// FileInfo represents detailed information about a file to be removed
type FileInfo struct {
	Path             string       `json:"path"`
	RelativePath     string       `json:"relative_path"`
	Size             int64        `json:"size"`
	ModTime          time.Time    `json:"mod_time"`
	Category         FileCategory `json:"category"`
	IsTrackedInLock  bool         `json:"is_tracked_in_lock"`
	IsModified       bool         `json:"is_modified"`
	PermissionIssue  bool         `json:"permission_issue"`
	PermissionError  string       `json:"permission_error,omitempty"`
	IsSymlink        bool         `json:"is_symlink"`
	SymlinkTarget    string       `json:"symlink_target,omitempty"`
}

// CategorySummary provides aggregate information for each file category
type CategorySummary struct {
	Category      FileCategory `json:"category"`
	Count         int          `json:"count"`
	TotalSize     int64        `json:"total_size"`
	TrackedFiles  int          `json:"tracked_files"`
	UntrackedFiles int         `json:"untracked_files"`
	ModifiedFiles int          `json:"modified_files"`
}

// RemovalPreview provides comprehensive information about what will be removed
type RemovalPreview struct {
	InstallPath     string                     `json:"install_path"`
	ClaudePath      string                     `json:"claude_path"`
	DiscoverySource DiscoverySource           `json:"discovery_source"`
	
	// File information
	Files           []FileInfo                `json:"files"`
	CategorySummary []CategorySummary         `json:"category_summary"`
	TotalFiles      int                       `json:"total_files"`
	TotalSize       int64                     `json:"total_size"`
	
	// Lockfile information
	LockFile        *config.LockFile          `json:"lock_file,omitempty"`
	UntrackedFiles  []FileInfo                `json:"untracked_files"`
	OrphanedFiles   []FileInfo                `json:"orphaned_files"`
	
	// Settings information
	SettingsFiles   []FileInfo                `json:"settings_files"`
	
	// Security and validation
	SecurityIssues  []SecurityIssue           `json:"security_issues"`
	ValidationErrors []ValidationError        `json:"validation_errors"`
}

// SecurityIssue represents a potential security concern during removal
type SecurityIssue struct {
	Type        string `json:"type"`
	FilePath    string `json:"file_path"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // "low", "medium", "high", "critical"
}

// ValidationError represents an issue that could prevent successful removal
type ValidationError struct {
	Type        string `json:"type"`
	FilePath    string `json:"file_path"`
	Description string `json:"description"`
}

// PreviewGenerator creates removal previews by scanning directories and analyzing files
type PreviewGenerator struct {
	pathDiscoverer PathDiscoverer
	verbose        bool
}

// NewPreviewGenerator creates a new preview generator
func NewPreviewGenerator(pathDiscoverer PathDiscoverer) *PreviewGenerator {
	return &PreviewGenerator{
		pathDiscoverer: pathDiscoverer,
		verbose:        false,
	}
}

// SetVerbose enables or disables verbose output
func (pg *PreviewGenerator) SetVerbose(verbose bool) {
	pg.verbose = verbose
}

// GeneratePreview creates a comprehensive removal preview
func (pg *PreviewGenerator) GeneratePreview() (*RemovalPreview, error) {
	preview := &RemovalPreview{
		Files:           make([]FileInfo, 0),
		UntrackedFiles:  make([]FileInfo, 0),
		OrphanedFiles:   make([]FileInfo, 0),
		SettingsFiles:   make([]FileInfo, 0),
		SecurityIssues:  make([]SecurityIssue, 0),
		ValidationErrors: make([]ValidationError, 0),
	}

	// Discover installation paths
	installPath, claudePath, source, err := pg.pathDiscoverer.DiscoverPaths()
	if err != nil {
		return nil, fmt.Errorf("failed to discover installation paths: %w", err)
	}

	preview.InstallPath = installPath
	preview.ClaudePath = claudePath
	preview.DiscoverySource = source

	if pg.verbose {
		fmt.Printf("Analyzing installation at: %s\n", installPath)
		fmt.Printf("Analyzing Claude config at: %s\n", claudePath)
	}

	// Load lockfile if available
	if source == DiscoverySourceLockfile {
		lockFile, err := pg.loadLockFile(installPath)
		if err != nil {
			if pg.verbose {
				fmt.Printf("Warning: Failed to load lockfile: %v\n", err)
			}
		} else {
			preview.LockFile = lockFile
		}
	}

	// Scan installation directory
	if err := pg.scanDirectory(installPath, "install", preview); err != nil {
		return nil, fmt.Errorf("failed to scan installation directory: %w", err)
	}

	// Scan Claude directory
	if err := pg.scanDirectory(claudePath, "claude", preview); err != nil {
		return nil, fmt.Errorf("failed to scan claude directory: %w", err)
	}

	// Analyze files against lockfile
	pg.analyzeAgainstLockfile(preview)

	// Generate category summaries
	pg.generateCategorySummaries(preview)

	// Perform security analysis
	pg.performSecurityAnalysis(preview)

	// Validate removal feasibility
	pg.validateRemovalFeasibility(preview)

	return preview, nil
}

// scanDirectory recursively scans a directory for the-startup related files
func (pg *PreviewGenerator) scanDirectory(basePath, dirType string, preview *RemovalPreview) error {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if pg.verbose {
			fmt.Printf("Directory does not exist: %s\n", basePath)
		}
		return nil // Not an error if directory doesn't exist
	}

	// Define expected patterns for documentation (not currently used for filtering)
	// but kept for potential future enhancement of file detection
	_ = dirType // Will be used in future enhancements

	// Do a general scan for any files that might be related
	return filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Log but don't fail on individual file errors
			if pg.verbose {
				fmt.Printf("Warning: Error accessing %s: %v\n", path, err)
			}
			return nil
		}

		// Skip the base directory itself
		if path == basePath {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			relPath = path
		}

		// Check if this is a the-startup related file
		if pg.isStartupRelated(path, relPath, dirType, d) {
			fileInfo, err := pg.analyzeFile(path, relPath, d)
			if err != nil {
				if pg.verbose {
					fmt.Printf("Warning: Failed to analyze file %s: %v\n", path, err)
				}
				return nil // Continue with other files
			}

			// Categorize based on directory type and file characteristics
			fileInfo.Category = pg.categorizeFile(path, relPath, dirType)

			preview.Files = append(preview.Files, *fileInfo)
		}

		return nil
	})
}

// isStartupRelated determines if a file is related to the-startup installation
func (pg *PreviewGenerator) isStartupRelated(fullPath, relPath, dirType string, d fs.DirEntry) bool {
	fileName := d.Name()
	
	// Always include lockfile
	if fileName == "the-startup.lock" {
		return true
	}

	// Check directory-specific patterns
	switch dirType {
	case "install":
		// Installation directory patterns
		return strings.Contains(relPath, "the-startup") ||
			strings.HasPrefix(relPath, "bin/the-startup") ||
			strings.HasPrefix(relPath, "rules/") ||
			strings.HasPrefix(relPath, "templates/") ||
			strings.HasPrefix(relPath, "logs/")

	case "claude":
		// Claude directory patterns
		return strings.HasPrefix(fileName, "the-") ||
			strings.HasPrefix(relPath, "agents/the-") ||
			strings.HasPrefix(relPath, "commands/s/") ||
			strings.HasPrefix(relPath, "output-styles/") ||
			fileName == "settings.json" ||
			fileName == "settings.local.json"
	}

	return false
}

// categorizeFile determines the category of a file based on its path and characteristics
func (pg *PreviewGenerator) categorizeFile(fullPath, relPath, dirType string) FileCategory {
	fileName := filepath.Base(fullPath)

	// Special files
	if fileName == "the-startup.lock" {
		return CategoryOther
	}

	if fileName == "settings.json" || fileName == "settings.local.json" {
		return CategorySettings
	}

	// Binary files
	if strings.Contains(relPath, "bin/") || strings.HasSuffix(fileName, ".exe") {
		return CategoryBinary
	}

	// Based on directory and file patterns
	if strings.HasPrefix(relPath, "agents/") && strings.HasSuffix(fileName, ".md") {
		return CategoryAgent
	}

	if strings.HasPrefix(relPath, "commands/") && strings.HasSuffix(fileName, ".md") {
		return CategoryCommand
	}

	if strings.HasPrefix(relPath, "templates/") && strings.HasSuffix(fileName, ".md") {
		return CategoryTemplate
	}

	if strings.HasPrefix(relPath, "rules/") && strings.HasSuffix(fileName, ".md") {
		return CategoryRule
	}

	if strings.HasPrefix(relPath, "logs/") && strings.HasSuffix(fileName, ".jsonl") {
		return CategoryLog
	}

	return CategoryOther
}

// analyzeFile gets detailed information about a file
func (pg *PreviewGenerator) analyzeFile(fullPath, relPath string, d fs.DirEntry) (*FileInfo, error) {
	info, err := d.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	fileInfo := &FileInfo{
		Path:         fullPath,
		RelativePath: relPath,
		Size:         info.Size(),
		ModTime:      info.ModTime(),
	}

	// Check if it's a symlink
	if info.Mode()&os.ModeSymlink != 0 {
		fileInfo.IsSymlink = true
		if target, err := os.Readlink(fullPath); err == nil {
			fileInfo.SymlinkTarget = target
		}
	}

	// Check permissions
	if err := pg.checkFilePermissions(fullPath); err != nil {
		fileInfo.PermissionIssue = true
		fileInfo.PermissionError = err.Error()
	}

	return fileInfo, nil
}

// checkFilePermissions verifies we can access and potentially remove the file
func (pg *PreviewGenerator) checkFilePermissions(path string) error {
	// Check if we can stat the file
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access file: %w", err)
	}

	// Check parent directory write permissions (needed for removal)
	parentDir := filepath.Dir(path)
	parentInfo, err := os.Stat(parentDir)
	if err != nil {
		return fmt.Errorf("cannot access parent directory: %w", err)
	}

	// Check if parent directory is writable
	// Note: This is a simplified check; real permissions are more complex
	if parentInfo.Mode().Perm()&0200 == 0 {
		return fmt.Errorf("parent directory not writable")
	}

	// Check if file is writable (helps with some removal scenarios)
	if !info.IsDir() && info.Mode().Perm()&0200 == 0 {
		// File is read-only, might need special handling
		// This isn't necessarily an error, but good to know
	}

	return nil
}

// loadLockFile loads the lockfile for analysis
func (pg *PreviewGenerator) loadLockFile(installPath string) (*config.LockFile, error) {
	lockFilePath := filepath.Join(installPath, "the-startup.lock")

	data, err := os.ReadFile(lockFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lockfile: %w", err)
	}

	var lockFile config.LockFile
	if err := json.Unmarshal(data, &lockFile); err != nil {
		return nil, fmt.Errorf("failed to parse lockfile: %w", err)
	}

	return &lockFile, nil
}

// analyzeAgainstLockfile compares discovered files with lockfile entries
func (pg *PreviewGenerator) analyzeAgainstLockfile(preview *RemovalPreview) {
	if preview.LockFile == nil {
		// Mark all files as untracked if no lockfile
		for i := range preview.Files {
			preview.Files[i].IsTrackedInLock = false
			preview.Files[i].Category = CategoryUntracked
		}
		preview.UntrackedFiles = append([]FileInfo{}, preview.Files...)
		return
	}

	// Create a map of lockfile entries for quick lookup
	lockFileMap := make(map[string]config.FileInfo)
	for filePath, fileInfo := range preview.LockFile.Files {
		lockFileMap[filePath] = fileInfo
	}

	for i := range preview.Files {
		file := &preview.Files[i]
		
		// Try to find this file in the lockfile
		var lockEntry config.FileInfo
		var found bool
		
		// Try different path mappings to match lockfile format
		possibleKeys := pg.generateLockfilePaths(file, preview)
		
		for _, key := range possibleKeys {
			if entry, exists := lockFileMap[key]; exists {
				lockEntry = entry
				found = true
				break
			}
		}

		file.IsTrackedInLock = found

		if found {
			// Check if file has been modified since installation
			if lockEntry.Size != file.Size {
				file.IsModified = true
			}
			
			// Parse lockfile timestamp and compare
			if lockTime, err := time.Parse(time.RFC3339, lockEntry.LastModified); err == nil {
				// Allow for small time differences due to filesystem precision
				if file.ModTime.After(lockTime.Add(time.Minute)) {
					file.IsModified = true
				}
			}
		} else {
			// File not in lockfile - it's untracked
			// However, preserve the original category for certain special files
			originalCategory := file.Category
			if originalCategory == CategorySettings || 
			   originalCategory == CategoryOther ||
			   strings.HasSuffix(file.Path, "the-startup.lock") {
				// Keep the original category for these special files
			} else {
				// Mark as untracked for other files
				file.Category = CategoryUntracked
			}
			preview.UntrackedFiles = append(preview.UntrackedFiles, *file)
		}
	}

	// Find orphaned files (in lockfile but not on disk)
	for lockPath := range lockFileMap {
		found := false
		for _, file := range preview.Files {
			possibleKeys := pg.generateLockfilePaths(&file, preview)
			for _, key := range possibleKeys {
				if key == lockPath {
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			// This file was in the lockfile but not found on disk
			orphanedFile := FileInfo{
				RelativePath:    lockPath,
				IsTrackedInLock: true,
				Category:        CategoryOther,
			}
			
			// Try to determine full path
			orphanedFile.Path = pg.reconstructFullPath(lockPath, preview)
			
			preview.OrphanedFiles = append(preview.OrphanedFiles, orphanedFile)
		}
	}
}

// generateLockfilePaths generates possible lockfile key formats for a given file
func (pg *PreviewGenerator) generateLockfilePaths(file *FileInfo, preview *RemovalPreview) []string {
	var keys []string
	
	// The lockfile uses different prefixes for different file types
	relPath := file.RelativePath

	// For files in the installation directory
	if strings.HasPrefix(file.Path, preview.InstallPath) {
		// Try with "startup/" prefix (used for startup assets)
		if !strings.HasPrefix(relPath, "bin/") {
			keys = append(keys, "startup/"+relPath)
		} else {
			// Binary files don't use startup/ prefix
			keys = append(keys, relPath)
		}
	}

	// For files in Claude directory - use path as-is
	if strings.HasPrefix(file.Path, preview.ClaudePath) {
		keys = append(keys, relPath)
	}

	// Also try the path as-is
	keys = append(keys, relPath)

	return keys
}

// reconstructFullPath attempts to reconstruct the full path from a lockfile entry
func (pg *PreviewGenerator) reconstructFullPath(lockPath string, preview *RemovalPreview) string {
	if strings.HasPrefix(lockPath, "startup/") {
		relPath := strings.TrimPrefix(lockPath, "startup/")
		return filepath.Join(preview.InstallPath, relPath)
	}

	if strings.HasPrefix(lockPath, "bin/") {
		return filepath.Join(preview.InstallPath, lockPath)
	}

	// Assume it's a Claude file
	return filepath.Join(preview.ClaudePath, lockPath)
}

// generateCategorySummaries creates aggregate statistics for each file category
func (pg *PreviewGenerator) generateCategorySummaries(preview *RemovalPreview) {
	categoryMap := make(map[FileCategory]*CategorySummary)

	// Initialize summaries for all categories
	for category := CategoryAgent; category <= CategoryUntracked; category++ {
		categoryMap[category] = &CategorySummary{
			Category: category,
		}
	}

	// Aggregate file information
	for _, file := range preview.Files {
		summary := categoryMap[file.Category]
		summary.Count++
		summary.TotalSize += file.Size

		if file.IsTrackedInLock {
			summary.TrackedFiles++
		} else {
			summary.UntrackedFiles++
		}

		if file.IsModified {
			summary.ModifiedFiles++
		}
	}

	// Convert to slice and filter out empty categories
	for category, summary := range categoryMap {
		if summary.Count > 0 {
			preview.CategorySummary = append(preview.CategorySummary, *summary)
		}
		
		// Update totals
		preview.TotalFiles += summary.Count
		preview.TotalSize += summary.TotalSize
		
		_ = category // Keep the compiler happy
	}
}

// performSecurityAnalysis identifies potential security issues with the removal
func (pg *PreviewGenerator) performSecurityAnalysis(preview *RemovalPreview) {
	for _, file := range preview.Files {
		// Check for path traversal attempts
		if strings.Contains(file.Path, "..") {
			preview.SecurityIssues = append(preview.SecurityIssues, SecurityIssue{
				Type:        "path_traversal",
				FilePath:    file.Path,
				Description: "File path contains directory traversal sequences",
				Severity:    "high",
			})
		}

		// Check for files outside expected directories
		if !pg.isPathSafe(file.Path, preview.InstallPath, preview.ClaudePath) {
			preview.SecurityIssues = append(preview.SecurityIssues, SecurityIssue{
				Type:        "unexpected_location",
				FilePath:    file.Path,
				Description: "File is outside expected installation directories",
				Severity:    "medium",
			})
		}

		// Check for suspicious symlinks
		if file.IsSymlink {
			if !pg.isSymlinkSafe(file.Path, file.SymlinkTarget, preview) {
				preview.SecurityIssues = append(preview.SecurityIssues, SecurityIssue{
					Type:        "suspicious_symlink",
					FilePath:    file.Path,
					Description: fmt.Sprintf("Symlink points to suspicious location: %s", file.SymlinkTarget),
					Severity:    "medium",
				})
			}
		}

		// Check for very large files (potential DoS)
		if file.Size > 100*1024*1024 { // 100MB
			preview.SecurityIssues = append(preview.SecurityIssues, SecurityIssue{
				Type:        "large_file",
				FilePath:    file.Path,
				Description: fmt.Sprintf("File is unusually large: %d bytes", file.Size),
				Severity:    "low",
			})
		}
	}
}

// isPathSafe checks if a file path is within the expected directories
func (pg *PreviewGenerator) isPathSafe(filePath, installPath, claudePath string) bool {
	// Resolve any symlinks to get the actual path
	resolvedPath, err := filepath.EvalSymlinks(filePath)
	if err != nil {
		// If we can't resolve, use the original path
		resolvedPath = filePath
	}

	// Convert to absolute paths for comparison
	absFilePath, err := filepath.Abs(resolvedPath)
	if err != nil {
		return false
	}

	absInstallPath, err := filepath.Abs(installPath)
	if err != nil {
		return false
	}

	absClaudePath, err := filepath.Abs(claudePath)
	if err != nil {
		return false
	}

	// Check if file is within either expected directory
	return strings.HasPrefix(absFilePath, absInstallPath) || 
		   strings.HasPrefix(absFilePath, absClaudePath)
}

// isSymlinkSafe checks if a symlink target is safe
func (pg *PreviewGenerator) isSymlinkSafe(linkPath, target string, preview *RemovalPreview) bool {
	// Resolve relative targets
	if !filepath.IsAbs(target) {
		linkDir := filepath.Dir(linkPath)
		target = filepath.Join(linkDir, target)
	}

	// Check if target is within safe directories
	return pg.isPathSafe(target, preview.InstallPath, preview.ClaudePath)
}

// validateRemovalFeasibility checks for issues that could prevent successful removal
func (pg *PreviewGenerator) validateRemovalFeasibility(preview *RemovalPreview) {
	for _, file := range preview.Files {
		// Check permission issues
		if file.PermissionIssue {
			preview.ValidationErrors = append(preview.ValidationErrors, ValidationError{
				Type:        "permission_denied",
				FilePath:    file.Path,
				Description: file.PermissionError,
			})
		}

		// Check for running processes (simplified check)
		if file.Category == CategoryBinary {
			if pg.isBinaryInUse(file.Path) {
				preview.ValidationErrors = append(preview.ValidationErrors, ValidationError{
					Type:        "file_in_use",
					FilePath:    file.Path,
					Description: "Binary file may be currently running",
				})
			}
		}

		// Check for critical system files (paranoid check)
		if pg.isCriticalSystemFile(file.Path) {
			preview.ValidationErrors = append(preview.ValidationErrors, ValidationError{
				Type:        "critical_file",
				FilePath:    file.Path,
				Description: "File appears to be critical system file - removal blocked",
			})
		}
	}
}

// isBinaryInUse checks if a binary file might be currently running
// This is a simplified check - a more robust implementation would check process lists
func (pg *PreviewGenerator) isBinaryInUse(binaryPath string) bool {
	// Try to open the file for exclusive write access
	// If it fails, it might be in use (or we lack permissions)
	file, err := os.OpenFile(binaryPath, os.O_WRONLY, 0)
	if err != nil {
		// Could be in use, or permission issue
		return true
	}
	file.Close()
	return false
}

// isCriticalSystemFile checks if a file path looks like a critical system file
func (pg *PreviewGenerator) isCriticalSystemFile(path string) bool {
	// This is a paranoid safety check
	dangerousPaths := []string{
		"/bin/", "/usr/bin/", "/sbin/", "/usr/sbin/",
		"/etc/", "/var/", "/lib/", "/usr/lib/",
		"/boot/", "/sys/", "/proc/", "/dev/",
		"C:\\Windows\\", "C:\\Program Files\\",
	}

	for _, dangerous := range dangerousPaths {
		if strings.HasPrefix(path, dangerous) {
			return true
		}
	}

	return false
}