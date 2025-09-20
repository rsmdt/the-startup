package ui

import (
	"fmt"
	"os"
	"strings"
)

// performUninstall removes the files that were selected for uninstall
func (m *MainModel) performUninstall() error {
	// Get the files to remove from the file selection model
	filesToRemove := m.fileSelectionModel.GetSelectedFiles()
	
	if len(filesToRemove) == 0 {
		return fmt.Errorf("no files to remove")
	}
	
	var errors []string
	
	// Remove each file
	for _, filePath := range filesToRemove {
		if err := os.Remove(filePath); err != nil {
			errors = append(errors, fmt.Sprintf("Failed to remove %s: %v", filePath, err))
		}
	}
	
	// Try to remove the entire startup directory if it's empty
	m.removeStartupDirectoryIfEmpty()
	
	if len(errors) > 0 {
		return fmt.Errorf("uninstall completed with errors:\n%s", strings.Join(errors, "\n"))
	}
	
	return nil
}

// removeStartupDirectoryIfEmpty removes the entire startup directory if it's empty
// NEVER touches anything inside .claude - only removes individual files there
func (m *MainModel) removeStartupDirectoryIfEmpty() {
	startupPath := m.installer.GetInstallPath()
	
	// Only try to remove the startup directory if it exists and is empty
	if entries, err := os.ReadDir(startupPath); err == nil && len(entries) == 0 {
		os.Remove(startupPath) // Ignore errors - this is best effort cleanup
	}
}