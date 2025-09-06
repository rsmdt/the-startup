package stats_test

import (
	"fmt"
	"time"

	"github.com/rsmdt/the-startup/internal/stats"
)

// Example demonstrates how to use the LogDiscovery to find Claude Code log files
func ExampleLogDiscovery_FindLogFiles() {
	// Create a new discovery instance
	discovery := stats.NewLogDiscovery()

	// Find all log files for the current project
	files, err := discovery.FindLogFiles("", stats.FilterOptions{})
	if err != nil {
		fmt.Printf("Error finding log files: %v\n", err)
		return
	}

	fmt.Printf("Found %d log files\n", len(files))
}

// Example demonstrates finding logs with time-based filtering (--since flag support)
func ExampleLogDiscovery_FindLogFiles_withTimeFilter() {
	discovery := stats.NewLogDiscovery()

	// Find logs from the last 7 days (supports --since flag)
	since := time.Now().AddDate(0, 0, -7)
	files, err := discovery.FindLogFiles("", stats.FilterOptions{
		StartTime: &since,
	})
	if err != nil {
		fmt.Printf("Error finding log files: %v\n", err)
		return
	}

	fmt.Printf("Found %d log files from the last 7 days\n", len(files))
}

// Example demonstrates checking if a project has Claude Code logs
func ExampleLogDiscovery_ValidateProjectPath() {
	discovery := stats.NewLogDiscovery()

	// Check if the current project has Claude logs
	currentProject := discovery.GetCurrentProject()
	if currentProject == "" {
		fmt.Println("Not in a Claude project directory")
		return
	}

	if discovery.ValidateProjectPath(currentProject) {
		fmt.Printf("Project '%s' has Claude Code logs\n", currentProject)
	} else {
		fmt.Printf("Project '%s' has no Claude Code logs\n", currentProject)
	}
}

// Example demonstrates working with a specific project path
func ExampleLogDiscovery_FindLogFiles_specificProject() {
	discovery := stats.NewLogDiscovery()

	// Find logs for a specific project
	// The project path is automatically sanitized to match Claude's format
	projectPath := "/Users/irudi/Code/personal/the-startup"
	
	// This will look in ~/.claude/projects/Users-irudi-Code-personal-the-startup/
	files, err := discovery.FindLogFiles(projectPath, stats.FilterOptions{})
	if err != nil {
		fmt.Printf("Error finding log files: %v\n", err)
		return
	}

	fmt.Printf("Found %d log files for project %s\n", len(files), projectPath)
}