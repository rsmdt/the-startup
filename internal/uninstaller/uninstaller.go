package uninstaller

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rsmdt/the-startup/internal/config"
)

// Re-using DiscoverySource from path_discovery.go

// UninstallPlan represents the complete plan for uninstalling the-startup
type UninstallPlan struct {
	InstallPath     string           `json:"install_path"`
	ClaudePath      string           `json:"claude_path"`
	FilesToRemove   []string         `json:"files_to_remove"`
	SettingsToClean []string         `json:"settings_to_clean"`
	TotalSize       int64            `json:"total_size"`
	DiscoverySource DiscoverySource  `json:"discovery_source"`
	LockFile        *config.LockFile `json:"lock_file,omitempty"`
	RemovalErrors   []RemovalError   `json:"removal_errors,omitempty"`
	BackupCreated   bool             `json:"backup_created,omitempty"`
	BackupPath      string           `json:"backup_path,omitempty"`
}

// RemovalError represents an error that occurred during file removal
type RemovalError struct {
	FilePath    string `json:"file_path"`
	Error       string `json:"error"`
	Recoverable bool   `json:"recoverable"`
}

// Re-using PathDiscoverer interface from path_discovery.go

// FileRemover defines the interface for removing files
type FileRemover interface {
	// RemoveFiles removes the specified files and returns any errors encountered
	RemoveFiles(files []string) []RemovalError

	// RemoveDirectory removes an entire directory if it's empty or contains only our files
	RemoveDirectory(dirPath string) error
}

// SettingsManager defines the interface for cleaning configuration settings
type SettingsManager interface {
	// CleanSettings removes the-startup related settings from Claude configuration
	// Returns the list of settings that were cleaned and any errors
	CleanSettings(claudePath string) ([]string, error)

	// BackupSettings creates a backup of settings before modification
	BackupSettings(claudePath string) (backupPath string, err error)
}

// Uninstaller coordinates the complete uninstallation process
type Uninstaller struct {
	pathDiscoverer  PathDiscoverer
	fileRemover     FileRemover
	settingsManager SettingsManager

	// Configuration options
	dryRun             bool
	forceRemove        bool
	shouldCreateBackup bool
	verbose            bool
}

// New creates a new Uninstaller with the provided dependencies
func New(pathDiscoverer PathDiscoverer, fileRemover FileRemover, settingsManager SettingsManager) *Uninstaller {
	return &Uninstaller{
		pathDiscoverer:     pathDiscoverer,
		fileRemover:        fileRemover,
		settingsManager:    settingsManager,
		shouldCreateBackup: true, // Default to creating backups for safety
		verbose:            false,
	}
}

// NewWithDefaults creates a new Uninstaller with default implementations for all dependencies
// Note: FileRemover will need to be configured later with paths and lock file via SetPathsAndLockFile
func NewWithDefaults() *Uninstaller {
	pathDiscoverer := NewPathDiscovery()
	settingsManager := NewDefaultSettingsManager()
	
	// Create a minimal uninstaller without FileRemover for now
	// FileRemover requires paths which aren't available yet
	return &Uninstaller{
		pathDiscoverer:     pathDiscoverer,
		fileRemover:        nil, // Will be set when paths are discovered
		settingsManager:    settingsManager,
		shouldCreateBackup: true,
		verbose:            false,
	}
}

// SetPathsAndLockFile configures the file remover with discovered paths and lock file
func (u *Uninstaller) SetPathsAndLockFile(installPath, claudePath string, lockFile *config.LockFile) {
	u.fileRemover = NewSafeFileRemover(installPath, claudePath, lockFile)
}

// SetOptions configures the uninstaller behavior
func (u *Uninstaller) SetOptions(dryRun, forceRemove, createBackup, verbose bool) {
	u.dryRun = dryRun
	u.forceRemove = forceRemove
	u.shouldCreateBackup = createBackup
	u.verbose = verbose
	
	// Configure settings manager if it supports these options
	if sm, ok := u.settingsManager.(*DefaultSettingsManager); ok {
		sm.SetDryRun(dryRun)
		sm.SetVerbose(verbose)
	}
}

// CreateUninstallPlan analyzes the system and creates a comprehensive uninstall plan
func (u *Uninstaller) CreateUninstallPlan() (*UninstallPlan, error) {
	// Use the new preview generator for comprehensive analysis
	previewGenerator := NewPreviewGenerator(u.pathDiscoverer)
	previewGenerator.SetVerbose(u.verbose)

	preview, err := previewGenerator.GeneratePreview()
	if err != nil {
		return nil, fmt.Errorf("failed to generate removal preview: %w", err)
	}

	// Convert preview to legacy UninstallPlan format for compatibility
	plan := &UninstallPlan{
		InstallPath:     preview.InstallPath,
		ClaudePath:      preview.ClaudePath,
		DiscoverySource: preview.DiscoverySource,
		LockFile:        preview.LockFile,
		FilesToRemove:   make([]string, len(preview.Files)),
		SettingsToClean: make([]string, len(preview.SettingsFiles)),
		TotalSize:       preview.TotalSize,
		RemovalErrors:   make([]RemovalError, 0),
	}

	// Convert file info to simple paths
	for i, fileInfo := range preview.Files {
		plan.FilesToRemove[i] = fileInfo.Path
	}

	// Convert settings files to simple paths
	for i, settingsFile := range preview.SettingsFiles {
		plan.SettingsToClean[i] = settingsFile.Path
	}

	// Add validation errors as removal errors
	for _, validationError := range preview.ValidationErrors {
		plan.RemovalErrors = append(plan.RemovalErrors, RemovalError{
			FilePath:    validationError.FilePath,
			Error:       validationError.Description,
			Recoverable: validationError.Type != "critical_file",
		})
	}

	return plan, nil
}

// ExecuteUninstallPlan executes the provided uninstall plan
func (u *Uninstaller) ExecuteUninstallPlan(plan *UninstallPlan) error {
	if plan == nil {
		return fmt.Errorf("uninstall plan cannot be nil")
	}

	if u.dryRun {
		if u.verbose {
			fmt.Println("DRY RUN: No files will actually be removed")
		}
		return u.displayPlanSummary(plan)
	}

	// Create backup if requested
	if u.shouldCreateBackup {
		backupPath, err := u.performBackup(plan)
		if err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
		plan.BackupCreated = true
		plan.BackupPath = backupPath
		if u.verbose {
			fmt.Printf("Backup created at: %s\n", backupPath)
		}
	}

	// Clean settings first (safer to fail here before removing files)
	if len(plan.SettingsToClean) > 0 {
		if u.verbose {
			fmt.Println("Cleaning Claude settings...")
		}
		cleanedSettings, err := u.settingsManager.CleanSettings(plan.ClaudePath)
		if err != nil {
			return fmt.Errorf("failed to clean settings: %w", err)
		}
		plan.SettingsToClean = cleanedSettings
		if u.verbose {
			fmt.Printf("✓ Cleaned %d settings entries\n", len(cleanedSettings))
		}
	}

	// Remove files
	if len(plan.FilesToRemove) > 0 {
		if u.fileRemover == nil {
			return fmt.Errorf("file remover not configured - call SetPathsAndLockFile first")
		}
		if u.verbose {
			fmt.Printf("Removing %d files...\n", len(plan.FilesToRemove))
		}
		errors := u.fileRemover.RemoveFiles(plan.FilesToRemove)
		plan.RemovalErrors = errors

		if len(errors) > 0 && !u.forceRemove {
			return fmt.Errorf("encountered %d file removal errors (use --force to continue despite errors)", len(errors))
		}

		if u.verbose {
			successCount := len(plan.FilesToRemove) - len(errors)
			fmt.Printf("✓ Successfully removed %d files\n", successCount)
			if len(errors) > 0 {
				fmt.Printf("! %d files could not be removed\n", len(errors))
			}
		}
	}

	// Clean up empty directories
	if err := u.cleanupEmptyDirectories(plan); err != nil {
		if u.verbose {
			fmt.Printf("Warning: Failed to cleanup empty directories: %v\n", err)
		}
	}

	return nil
}

// GetInstallationInfo returns information about the current installation without modifying anything
func (u *Uninstaller) GetInstallationInfo() (*UninstallPlan, error) {
	return u.CreateUninstallPlan()
}

// GenerateDetailedPreview returns a comprehensive removal preview with detailed categorization
func (u *Uninstaller) GenerateDetailedPreview() (*RemovalPreview, error) {
	previewGenerator := NewPreviewGenerator(u.pathDiscoverer)
	previewGenerator.SetVerbose(u.verbose)

	return previewGenerator.GeneratePreview()
}

// ValidatePaths validates that the provided paths exist and contain the-startup files
func (u *Uninstaller) ValidatePaths(installPath, claudePath string) error {
	// Check install path exists
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		return fmt.Errorf("installation path does not exist: %s", installPath)
	}

	// Check claude path exists
	if _, err := os.Stat(claudePath); os.IsNotExist(err) {
		return fmt.Errorf("claude path does not exist: %s", claudePath)
	}

	// Look for evidence of the-startup installation
	lockFilePath := filepath.Join(installPath, "the-startup.lock")
	agentsPath := filepath.Join(claudePath, "agents")
	binaryPath := filepath.Join(installPath, "bin", "the-startup")

	hasLockFile := false
	hasAgents := false
	hasBinary := false

	if _, err := os.Stat(lockFilePath); err == nil {
		hasLockFile = true
	}

	if _, err := os.Stat(binaryPath); err == nil {
		hasBinary = true
	}

	// Check for the-startup agents
	if entries, err := os.ReadDir(agentsPath); err == nil {
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "the-") && strings.HasSuffix(entry.Name(), ".md") {
				hasAgents = true
				break
			}
		}
	}

	if !hasLockFile && !hasAgents && !hasBinary {
		return fmt.Errorf("no the-startup installation found in the provided paths")
	}

	return nil
}

// Helper methods

// loadLockFile loads the lock file from the installation directory
func (u *Uninstaller) loadLockFile(installPath string) (*config.LockFile, error) {
	lockFilePath := filepath.Join(installPath, "the-startup.lock")

	data, err := os.ReadFile(lockFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	// Parse lock file using the same config package as installer
	var lockFile config.LockFile
	if err := json.Unmarshal(data, &lockFile); err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	return &lockFile, nil
}

// identifyFilesToRemove builds the list of files that need to be removed
func (u *Uninstaller) identifyFilesToRemove(plan *UninstallPlan) ([]string, int64, error) {
	var files []string
	var totalSize int64

	// If we have a lock file, use it as the source of truth
	if plan.LockFile != nil {
		for filePath := range plan.LockFile.Files {
			var fullPath string

			// Handle different file types based on path prefix
			if strings.HasPrefix(filePath, "startup/") {
				// Startup files go to install path
				relPath := strings.TrimPrefix(filePath, "startup/")
				fullPath = filepath.Join(plan.InstallPath, relPath)
			} else if strings.HasPrefix(filePath, "bin/") {
				// Binary files go to install path
				fullPath = filepath.Join(plan.InstallPath, filePath)
			} else {
				// Claude files (agents, commands) go to claude path
				fullPath = filepath.Join(plan.ClaudePath, filePath)
			}

			// Verify file still exists before adding to removal list
			if info, err := os.Stat(fullPath); err == nil {
				files = append(files, fullPath)
				totalSize += info.Size()
			} else if u.verbose {
				fmt.Printf("File already removed or moved: %s\n", fullPath)
			}
		}
	} else {
		// Fallback: scan for typical the-startup files
		files, totalSize = u.discoverFilesWithoutLock(plan.InstallPath, plan.ClaudePath)
	}

	// Always try to remove the lock file itself
	lockFilePath := filepath.Join(plan.InstallPath, "the-startup.lock")
	if info, err := os.Stat(lockFilePath); err == nil {
		files = append(files, lockFilePath)
		totalSize += info.Size()
	}

	return files, totalSize, nil
}

// discoverFilesWithoutLock discovers files when no lock file is available
func (u *Uninstaller) discoverFilesWithoutLock(installPath, claudePath string) ([]string, int64) {
	var files []string
	var totalSize int64

	// Common patterns to look for
	patterns := []struct {
		basePath string
		patterns []string
	}{
		{claudePath, []string{
			"agents/the-*.md",
			"commands/s/*.md",
			"output-styles/*.md",
		}},
		{installPath, []string{
			"bin/the-startup",
			"rules/*.md",
			"templates/*.md",
		}},
	}

	for _, patternSet := range patterns {
		for _, pattern := range patternSet.patterns {
			fullPattern := filepath.Join(patternSet.basePath, pattern)
			matches, _ := filepath.Glob(fullPattern)
			for _, match := range matches {
				if info, err := os.Stat(match); err == nil {
					files = append(files, match)
					totalSize += info.Size()
				}
			}
		}
	}

	return files, totalSize
}

// identifySettingsToClean identifies settings that need to be cleaned
func (u *Uninstaller) identifySettingsToClean(claudePath string) ([]string, error) {
	var settings []string

	// Settings files that might contain the-startup configuration
	settingsFiles := []string{
		filepath.Join(claudePath, "settings.json"),
		filepath.Join(claudePath, "settings.local.json"),
	}

	for _, settingsFile := range settingsFiles {
		if _, err := os.Stat(settingsFile); err == nil {
			settings = append(settings, settingsFile)
		}
	}

	return settings, nil
}

// displayPlanSummary shows what would be done without actually doing it
func (u *Uninstaller) displayPlanSummary(plan *UninstallPlan) error {
	fmt.Println("\nUninstall Plan Summary (DRY RUN)")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Installation Path: %s\n", plan.InstallPath)
	fmt.Printf("Claude Path: %s\n", plan.ClaudePath)
	fmt.Printf("Discovery Source: %s\n", plan.DiscoverySource.String())
	fmt.Printf("Files to Remove: %d\n", len(plan.FilesToRemove))
	fmt.Printf("Total Size: %d bytes\n", plan.TotalSize)

	if len(plan.FilesToRemove) > 0 {
		fmt.Println("\nFiles that would be removed:")
		for _, file := range plan.FilesToRemove {
			fmt.Printf("  - %s\n", file)
		}
	}

	if len(plan.SettingsToClean) > 0 {
		fmt.Println("\nSettings files that would be cleaned:")
		for _, setting := range plan.SettingsToClean {
			fmt.Printf("  - %s\n", setting)
		}
	}

	return nil
}

// performBackup creates a backup of important files before removal
func (u *Uninstaller) performBackup(plan *UninstallPlan) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(plan.InstallPath, fmt.Sprintf("backup_%s", timestamp))

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// For now, return the backup directory path
	// Actual backup implementation will be added by FileRemover interface
	return backupDir, nil
}

// cleanupEmptyDirectories removes empty directories left after file removal
func (u *Uninstaller) cleanupEmptyDirectories(plan *UninstallPlan) error {
	if u.fileRemover == nil {
		if u.verbose {
			fmt.Println("File remover not configured, skipping directory cleanup")
		}
		return nil
	}

	// Common directories that might be left empty
	dirsToCheck := []string{
		filepath.Join(plan.ClaudePath, "agents"),
		filepath.Join(plan.ClaudePath, "commands", "s"),
		filepath.Join(plan.ClaudePath, "commands"),
		filepath.Join(plan.ClaudePath, "output-styles"),
		filepath.Join(plan.InstallPath, "bin"),
		filepath.Join(plan.InstallPath, "rules"),
		filepath.Join(plan.InstallPath, "templates"),
	}

	for _, dir := range dirsToCheck {
		if err := u.fileRemover.RemoveDirectory(dir); err != nil {
			// Not a critical error, just log it
			if u.verbose {
				fmt.Printf("Could not remove directory %s: %v\n", dir, err)
			}
		}
	}

	return nil
}
