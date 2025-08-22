package uninstaller

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rsmdt/the-startup/internal/config"
)

// DiscoverySource indicates how installation paths were discovered
type DiscoverySource int

const (
	// DiscoverySourceLockfile indicates paths were read from lockfile
	DiscoverySourceLockfile DiscoverySource = iota
	// DiscoverySourceAutoDetect indicates paths were auto-detected
	DiscoverySourceAutoDetect
	// DiscoverySourceUserInput indicates paths were provided by user
	DiscoverySourceUserInput
)

// String returns a human-readable representation of the discovery source
func (ds DiscoverySource) String() string {
	switch ds {
	case DiscoverySourceLockfile:
		return "lockfile"
	case DiscoverySourceAutoDetect:
		return "auto-detected"
	case DiscoverySourceUserInput:
		return "user-provided"
	default:
		return "unknown"
	}
}

// PathDiscoverer interface for discovering installation paths
type PathDiscoverer interface {
	DiscoverPaths() (installPath, claudePath string, source DiscoverySource, err error)
}

// DefaultPathDiscovery implements PathDiscoverer with lockfile reading and fallback logic
type DefaultPathDiscovery struct {
	// Optional: custom lockfile path for testing
	lockfilePath string
}

// NewPathDiscovery creates a new DefaultPathDiscovery instance
func NewPathDiscovery() *DefaultPathDiscovery {
	return &DefaultPathDiscovery{}
}

// NewPathDiscoveryWithLockfile creates a new DefaultPathDiscovery with custom lockfile path
func NewPathDiscoveryWithLockfile(lockfilePath string) *DefaultPathDiscovery {
	return &DefaultPathDiscovery{
		lockfilePath: lockfilePath,
	}
}

// DiscoverPaths discovers installation paths using multiple strategies
func (d *DefaultPathDiscovery) DiscoverPaths() (installPath, claudePath string, source DiscoverySource, err error) {
	// Strategy 1: Try to read from lockfile
	installPath, claudePath, err = d.readFromLockfile()
	if err == nil {
		// Successfully read from lockfile
		return installPath, claudePath, DiscoverySourceLockfile, nil
	}

	// Strategy 2: Auto-detect common installation paths
	installPath, claudePath, err = d.autoDetectPaths()
	if err == nil {
		// Successfully auto-detected
		return installPath, claudePath, DiscoverySourceAutoDetect, nil
	}

	// Strategy 3: Return default paths for user input fallback
	installPath, claudePath = d.getDefaultPaths()
	return installPath, claudePath, DiscoverySourceUserInput, nil
}

// readFromLockfile attempts to read installation paths from the lockfile
func (d *DefaultPathDiscovery) readFromLockfile() (installPath, claudePath string, err error) {
	lockfilePath := d.getLockfilePath()

	// Check if lockfile exists
	if _, err := os.Stat(lockfilePath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("lockfile not found at %s", lockfilePath)
	}

	// Read lockfile content
	data, err := os.ReadFile(lockfilePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read lockfile: %w", err)
	}

	// Parse JSON
	var lockFile config.LockFile
	if err := json.Unmarshal(data, &lockFile); err != nil {
		return "", "", fmt.Errorf("failed to parse lockfile JSON: %w", err)
	}

	// Extract paths
	installPath = lockFile.InstallPath
	claudePath = lockFile.ClaudePath

	// Validate paths are not empty
	if installPath == "" || claudePath == "" {
		return "", "", fmt.Errorf("lockfile contains empty paths: install='%s', claude='%s'",
			installPath, claudePath)
	}

	// Expand tilde paths to full paths
	installPath, err = expandTildePath(installPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to expand install path: %w", err)
	}

	claudePath, err = expandTildePath(claudePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to expand claude path: %w", err)
	}

	// Validate paths exist and are accessible
	if err := d.validatePath(installPath, "installation"); err != nil {
		return "", "", err
	}

	if err := d.validatePath(claudePath, "claude"); err != nil {
		return "", "", err
	}

	return installPath, claudePath, nil
}

// autoDetectPaths attempts to auto-detect installation paths by looking for evidence
func (d *DefaultPathDiscovery) autoDetectPaths() (installPath, claudePath string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Common installation paths to check
	candidateInstallPaths := []string{
		".the-startup",                                   // project-local
		filepath.Join(homeDir, ".the-startup"),           // home directory
		filepath.Join(homeDir, ".config", "the-startup"), // XDG config
	}

	// Common claude paths to check
	candidateClaudePaths := []string{
		filepath.Join(homeDir, ".claude"),           // default claude path
		filepath.Join(homeDir, ".config", "claude"), // XDG config
	}

	// Look for evidence of installation
	var foundInstallPath, foundClaudePath string

	// Check for install paths by looking for evidence
	for _, path := range candidateInstallPaths {
		if d.hasInstallationEvidence(path) {
			foundInstallPath = path
			break
		}
	}

	// Check for claude paths by looking for our agents
	for _, path := range candidateClaudePaths {
		if d.hasClaudeEvidence(path) {
			foundClaudePath = path
			break
		}
	}

	// Both paths must be found for successful auto-detection
	if foundInstallPath == "" || foundClaudePath == "" {
		return "", "", fmt.Errorf("could not auto-detect installation paths")
	}

	return foundInstallPath, foundClaudePath, nil
}

// hasInstallationEvidence checks if a path contains evidence of the-startup installation
func (d *DefaultPathDiscovery) hasInstallationEvidence(path string) bool {
	// Check for lockfile
	lockfilePath := filepath.Join(path, "the-startup.lock")
	if _, err := os.Stat(lockfilePath); err == nil {
		return true
	}

	// Check for binary
	binaryPath := filepath.Join(path, "bin", "the-startup")
	if _, err := os.Stat(binaryPath); err == nil {
		return true
	}

	// Check for templates directory
	templatesPath := filepath.Join(path, "templates")
	if _, err := os.Stat(templatesPath); err == nil {
		return true
	}

	// Check for rules directory
	rulesPath := filepath.Join(path, "rules")
	if _, err := os.Stat(rulesPath); err == nil {
		return true
	}

	return false
}

// hasClaudeEvidence checks if a path contains evidence of our claude installation
func (d *DefaultPathDiscovery) hasClaudeEvidence(path string) bool {
	// Check for agents directory with our agents
	agentsPath := filepath.Join(path, "agents")
	entries, err := os.ReadDir(agentsPath)
	if err != nil {
		return false
	}

	// Look for files starting with "the-" (our agent naming convention)
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "the-") && strings.HasSuffix(entry.Name(), ".md") {
			return true
		}
	}

	return false
}

// getDefaultPaths returns sensible default paths when auto-detection fails
func (d *DefaultPathDiscovery) getDefaultPaths() (installPath, claudePath string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to relative paths if home directory unavailable
		return ".the-startup", ".claude"
	}

	return filepath.Join(homeDir, ".the-startup"), filepath.Join(homeDir, ".claude")
}

// getLockfilePath returns the path to the lockfile
func (d *DefaultPathDiscovery) getLockfilePath() string {
	// Use custom path if provided (for testing)
	if d.lockfilePath != "" {
		return d.lockfilePath
	}

	// Try common lockfile locations
	candidates := []string{
		"the-startup.lock",              // current directory
		".the-startup/the-startup.lock", // project-local
	}

	// Add home directory candidates
	if homeDir, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates,
			filepath.Join(homeDir, ".the-startup", "the-startup.lock"),
			filepath.Join(homeDir, ".config", "the-startup", "the-startup.lock"),
		)
	}

	// Return first candidate that exists
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// Default to project-local if none found
	return ".the-startup/the-startup.lock"
}

// validatePath validates that a path exists and is accessible
func (d *DefaultPathDiscovery) validatePath(path, pathType string) error {
	// Check if path exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s path does not exist: %s", pathType, path)
	}
	if err != nil {
		return fmt.Errorf("cannot access %s path %s: %w", pathType, path, err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("%s path is not a directory: %s", pathType, path)
	}

	// Check if we have read access (needed for scanning files)
	testFile := filepath.Join(path, ".the-startup-access-test")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("insufficient permissions for %s path %s: %w", pathType, path, err)
	}
	file.Close()
	os.Remove(testFile) // Clean up test file

	return nil
}

// expandTildePath expands ~ to the user's home directory
func expandTildePath(path string) (string, error) {
	if !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, path[2:]), nil
}
