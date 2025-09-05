package cmd

import (
	"testing"
	// Integration tests for the old agent logging functionality have been removed
	// as the log command now implements metrics collection instead.
	// The new metrics functionality is tested in internal/log package.
)

// TestLogCommandIntegration has been deprecated
// The log command now processes metrics, not agent instructions
func TestLogCommandIntegration(t *testing.T) {
	t.Skip("Test deprecated: log command now implements metrics collection")
}

// TestLogCommandUserIntegration has been deprecated  
// The log command now processes metrics, not agent instructions
func TestLogCommandUserIntegration(t *testing.T) {
	t.Skip("Test deprecated: log command now implements metrics collection")
}

// TestLogCommandWithoutFlags has been deprecated
// The log command now processes hooks without requiring flags
func TestLogCommandWithoutFlags(t *testing.T) {
	t.Skip("Test deprecated: log command now processes hooks automatically")
}

// TestLogCommandFiltersNonTaskTools has been deprecated
// Filtering is now handled differently in the metrics implementation
func TestLogCommandFiltersNonTaskTools(t *testing.T) {
	t.Skip("Test deprecated: filtering logic moved to metrics implementation")
}

// TestLogCommandFailsSilentlyOnError has been deprecated
// Silent failure is still implemented but for metrics processing
func TestLogCommandFailsSilentlyOnError(t *testing.T) {
	t.Skip("Test deprecated: silent failure now handled in metrics processing")
}

// TestLogCommandInvalidFlags has been deprecated
// Flag validation is no longer needed as hooks are auto-detected
func TestLogCommandInvalidFlags(t *testing.T) {
	t.Skip("Test deprecated: hook type auto-detection replaces flag validation")
}

// TestLogCommandReadMode has been deprecated
// Read functionality has been replaced with metrics analysis subcommands
func TestLogCommandReadMode(t *testing.T) {
	t.Skip("Test deprecated: read mode replaced with metrics subcommands")
}