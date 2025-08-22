package uninstaller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rsmdt/the-startup/internal/config"
)

// SafeFileRemover handles the safe removal of files with verification and cleanup
// This implements the FileRemover interface defined in uninstaller.go
type SafeFileRemover struct {
	dryRun     bool                       // If true, simulate removal without actually deleting
	installPath string                    // Path to .the-startup installation
	claudePath  string                   // Path to .claude directory
	lockFile    *config.LockFile         // Current installation lock file
}

// RemovalResult represents the result of attempting to remove a file
type RemovalResult struct {
	Path         string // Absolute path of the file
	RelativePath string // Relative path as stored in lock file
	Removed      bool   // Whether the file was successfully removed
	Skipped      bool   // Whether the file was skipped (e.g., doesn't exist)
	Modified     bool   // Whether the file was modified since installation
	Error        string // Error message if removal failed
	Warning      string // Warning message (e.g., for modified files)
}

// DirectoryCleanupResult represents the result of cleaning up directories
type DirectoryCleanupResult struct {
	Path    string // Absolute path of the directory
	Removed bool   // Whether the directory was removed
	Error   string // Error message if cleanup failed
	Warning string // Warning message (e.g., directory not empty)
}

// RemovalSummary provides a comprehensive summary of the removal operation
type RemovalSummary struct {
	FilesProcessed      int                      // Total number of files processed
	FilesRemoved        int                      // Number of files successfully removed
	FilesSkipped        int                      // Number of files skipped (not found)
	FilesModified       int                      // Number of files modified since installation
	FilesFailed         int                      // Number of files that failed to remove
	DirectoriesRemoved  int                      // Number of empty directories removed
	DirectoriesFailed   int                      // Number of directories that failed to remove
	FileResults         []RemovalResult          // Detailed results for each file
	DirectoryResults    []DirectoryCleanupResult // Detailed results for directory cleanup
	Errors              []string                 // Collection of all errors encountered
	Warnings            []string                 // Collection of all warnings
}

// NewSafeFileRemover creates a new SafeFileRemover instance
func NewSafeFileRemover(installPath, claudePath string, lockFile *config.LockFile) *SafeFileRemover {
	return &SafeFileRemover{
		dryRun:      false,
		installPath: installPath,
		claudePath:  claudePath,
		lockFile:    lockFile,
	}
}

// SetDryRun enables or disables dry-run mode
func (fr *SafeFileRemover) SetDryRun(enabled bool) {
	fr.dryRun = enabled
}

// RemoveFiles implements the FileRemover interface
func (fr *SafeFileRemover) RemoveFiles(files []string) []RemovalError {
	var errors []RemovalError
	
	for _, filePath := range files {
		if err := os.Remove(filePath); err != nil {
			errors = append(errors, RemovalError{
				FilePath:    filePath,
				Error:       err.Error(),
				Recoverable: !os.IsPermission(err),
			})
		}
	}
	
	return errors
}

// RemoveDirectory implements the FileRemover interface  
func (fr *SafeFileRemover) RemoveDirectory(dirPath string) error {
	// Check if directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Already removed, consider success
		}
		return fmt.Errorf("failed to check directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", dirPath)
	}

	// Check if directory is empty
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if len(entries) > 0 {
		return fmt.Errorf("directory not empty: %s", dirPath)
	}

	// Remove empty directory
	if fr.dryRun {
		return nil // Simulate success in dry-run mode
	}
	
	return os.Remove(dirPath)
}

// RemoveAllFiles safely removes all files tracked in the lock file
func (fr *SafeFileRemover) RemoveAllFiles() (*RemovalSummary, error) {
	if fr.lockFile == nil {
		return nil, fmt.Errorf("no lock file provided")
	}

	summary := &RemovalSummary{
		FileResults:      make([]RemovalResult, 0, len(fr.lockFile.Files)),
		DirectoryResults: make([]DirectoryCleanupResult, 0),
		Errors:          make([]string, 0),
		Warnings:        make([]string, 0),
	}

	// Process each file in the lock file
	for relPath, fileInfo := range fr.lockFile.Files {
		result := fr.removeFile(relPath, fileInfo)
		summary.FileResults = append(summary.FileResults, result)
		summary.FilesProcessed++

		if result.Error != "" {
			summary.FilesFailed++
			summary.Errors = append(summary.Errors, fmt.Sprintf("%s: %s", result.RelativePath, result.Error))
		} else if result.Skipped {
			summary.FilesSkipped++
		} else if result.Removed {
			summary.FilesRemoved++
		}

		if result.Modified {
			summary.FilesModified++
		}

		if result.Warning != "" {
			summary.Warnings = append(summary.Warnings, fmt.Sprintf("%s: %s", result.RelativePath, result.Warning))
		}
	}

	// Clean up empty directories
	if summary.FilesRemoved > 0 {
		fr.cleanupDirectories(summary)
	}

	return summary, nil
}

// removeFile safely removes a single file with verification
func (fr *SafeFileRemover) removeFile(relPath string, fileInfo config.FileInfo) RemovalResult {
	result := RemovalResult{
		RelativePath: relPath,
	}

	// Determine the absolute path
	var absPath string
	if strings.HasPrefix(relPath, "startup/") {
		// Startup assets are stored with "startup/" prefix
		cleanRelPath := strings.TrimPrefix(relPath, "startup/")
		absPath = filepath.Join(fr.installPath, cleanRelPath)
	} else if strings.HasPrefix(relPath, "bin/") {
		// Binary files go to install path directly
		absPath = filepath.Join(fr.installPath, relPath)
	} else {
		// Claude assets go to claude path
		absPath = filepath.Join(fr.claudePath, relPath)
	}
	
	result.Path = absPath

	// Check if file exists
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			result.Skipped = true
			result.Warning = "file not found (already removed or never existed)"
			return result
		}
		result.Error = fmt.Sprintf("failed to check file: %v", err)
		return result
	}

	// Verify file integrity if checksum is available
	if fileInfo.Checksum != "" {
		currentChecksum, err := fr.calculateChecksum(absPath)
		if err != nil {
			result.Error = fmt.Sprintf("failed to calculate checksum: %v", err)
			return result
		}

		if currentChecksum != fileInfo.Checksum {
			result.Modified = true
			result.Warning = "file has been modified since installation"
		}
	}

	// Additional verification: check size and modification time
	if info.Size() != fileInfo.Size {
		result.Modified = true
		if result.Warning != "" {
			result.Warning += "; size differs from installation"
		} else {
			result.Warning = "size differs from installation"
		}
	}

	// In dry-run mode, simulate removal
	if fr.dryRun {
		result.Removed = true
		result.Warning = result.Warning + " (dry-run mode - not actually removed)"
		return result
	}

	// Attempt to remove the file
	if err := os.Remove(absPath); err != nil {
		result.Error = fmt.Sprintf("failed to remove file: %v", err)
		return result
	}

	result.Removed = true
	return result
}

// calculateChecksum computes the SHA256 checksum of a file
func (fr *SafeFileRemover) calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// cleanupDirectories removes empty directories left after file removal
func (fr *SafeFileRemover) cleanupDirectories(summary *RemovalSummary) {
	// Collect unique directories from removed files
	dirsToCheck := make(map[string]bool)
	
	for _, result := range summary.FileResults {
		if result.Removed {
			dir := filepath.Dir(result.Path)
			dirsToCheck[dir] = true
		}
	}

	// Sort directories by depth (deepest first) for proper cleanup order
	dirs := make([]string, 0, len(dirsToCheck))
	for dir := range dirsToCheck {
		dirs = append(dirs, dir)
	}

	// Sort by path length (longer paths = deeper directories)
	for i := 0; i < len(dirs); i++ {
		for j := i + 1; j < len(dirs); j++ {
			if len(dirs[i]) < len(dirs[j]) {
				dirs[i], dirs[j] = dirs[j], dirs[i]
			}
		}
	}

	// Try to remove each directory
	for _, dir := range dirs {
		result := fr.cleanupDirectory(dir)
		summary.DirectoryResults = append(summary.DirectoryResults, result)
		
		if result.Error != "" {
			summary.DirectoriesFailed++
			summary.Errors = append(summary.Errors, fmt.Sprintf("directory cleanup %s: %s", dir, result.Error))
		} else if result.Removed {
			summary.DirectoriesRemoved++
		}

		if result.Warning != "" {
			summary.Warnings = append(summary.Warnings, fmt.Sprintf("directory cleanup %s: %s", dir, result.Warning))
		}
	}
}

// cleanupDirectory attempts to remove a directory if it's empty
func (fr *SafeFileRemover) cleanupDirectory(dirPath string) DirectoryCleanupResult {
	result := DirectoryCleanupResult{
		Path: dirPath,
	}

	// Don't try to remove the root installation or claude directories
	if dirPath == fr.installPath || dirPath == fr.claudePath {
		result.Warning = "skipping root directory"
		return result
	}

	// Check if directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			result.Warning = "directory already removed"
			return result
		}
		result.Error = fmt.Sprintf("failed to check directory: %v", err)
		return result
	}

	if !info.IsDir() {
		result.Warning = "path is not a directory"
		return result
	}

	// Check if directory is empty
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read directory: %v", err)
		return result
	}

	if len(entries) > 0 {
		result.Warning = "directory not empty"
		return result
	}

	// In dry-run mode, simulate removal
	if fr.dryRun {
		result.Removed = true
		result.Warning = "dry-run mode - not actually removed"
		return result
	}

	// Attempt to remove empty directory
	if err := os.Remove(dirPath); err != nil {
		result.Error = fmt.Sprintf("failed to remove directory: %v", err)
		return result
	}

	result.Removed = true
	return result
}

// VerifyFiles checks the integrity of installed files without removing them
func (fr *SafeFileRemover) VerifyFiles() (*RemovalSummary, error) {
	if fr.lockFile == nil {
		return nil, fmt.Errorf("no lock file provided")
	}

	// Temporarily enable dry-run mode to prevent actual removal
	originalDryRun := fr.dryRun
	fr.dryRun = true
	defer func() {
		fr.dryRun = originalDryRun
	}()

	// Use the same logic as RemoveAllFiles but with dry-run enabled
	return fr.RemoveAllFiles()
}

// GetRemovalGuidance provides human-readable guidance based on removal results
func (fr *SafeFileRemover) GetRemovalGuidance(summary *RemovalSummary) []string {
	var guidance []string

	if summary.FilesFailed > 0 {
		guidance = append(guidance, fmt.Sprintf("❌ %d files failed to remove. Check file permissions and ensure files are not in use.", summary.FilesFailed))
	}

	if summary.FilesModified > 0 {
		guidance = append(guidance, fmt.Sprintf("⚠️  %d files were modified since installation. Review changes before removal.", summary.FilesModified))
	}

	if summary.DirectoriesFailed > 0 {
		guidance = append(guidance, fmt.Sprintf("⚠️  %d directories could not be removed (may contain other files).", summary.DirectoriesFailed))
	}

	if len(summary.Errors) > 0 {
		guidance = append(guidance, "")
		guidance = append(guidance, "Detailed errors:")
		for _, err := range summary.Errors {
			guidance = append(guidance, fmt.Sprintf("  • %s", err))
		}
	}

	if len(summary.Warnings) > 0 {
		guidance = append(guidance, "")
		guidance = append(guidance, "Warnings:")
		for _, warning := range summary.Warnings {
			guidance = append(guidance, fmt.Sprintf("  • %s", warning))
		}
	}

	if len(guidance) == 0 {
		guidance = append(guidance, "✅ All files removed successfully!")
	}

	return guidance
}

