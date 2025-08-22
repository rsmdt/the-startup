package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/uninstaller"
)

// RunUninstaller demonstrates how to use the UninstallModel for the uninstall process
// This function shows the integration pattern without actually implementing the full command
func RunUninstaller(dryRun, forceRemove, createBackup, verbose bool) error {
	// Create the main uninstall model with the provided options
	model := NewUninstallModel(dryRun, forceRemove, createBackup, verbose)

	// Create and run the TUI program
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program and wait for completion
	finalModel, err := program.Run()
	if err != nil {
		return fmt.Errorf("uninstall UI failed: %w", err)
	}

	// Check if the user completed the process successfully
	if uninstallModel, ok := finalModel.(*UninstallModel); ok {
		if uninstallModel.Ready() {
			if dryRun {
				fmt.Println("Dry run completed - no files were actually removed")
			} else {
				fmt.Println("Uninstall completed successfully")
			}
		}
	}

	return nil
}

// ExampleUninstallFlow demonstrates the typical usage pattern
func ExampleUninstallFlow() {
	// Example 1: Dry run mode (safe preview)
	fmt.Println("Example 1: Dry run uninstall")
	err := RunUninstaller(true, false, true, true)
	if err != nil {
		fmt.Printf("Dry run failed: %v\n", err)
	}

	// Example 2: Real uninstall with backup
	fmt.Println("\nExample 2: Real uninstall with backup")
	err = RunUninstaller(false, false, true, false)
	if err != nil {
		fmt.Printf("Uninstall failed: %v\n", err)
	}

	// Example 3: Force removal without backup (dangerous)
	fmt.Println("\nExample 3: Force removal without backup")
	err = RunUninstaller(false, true, false, true)
	if err != nil {
		fmt.Printf("Force removal failed: %v\n", err)
	}
}

// CreateMockPreview creates a mock removal preview for testing and demonstration
func CreateMockPreview() *uninstaller.RemovalPreview {
	return &uninstaller.RemovalPreview{
		InstallPath:     "/Users/username/.config/the-startup",
		ClaudePath:      "/Users/username/.claude",
		DiscoverySource: uninstaller.DiscoverySourceLockfile,
		Files: []uninstaller.FileInfo{
			{
				Path:            "/Users/username/.claude/agents/the-architect.md",
				RelativePath:    "agents/the-architect.md",
				Size:            2048,
				Category:        uninstaller.CategoryAgent,
				IsTrackedInLock: true,
			},
			{
				Path:            "/Users/username/.claude/agents/the-frontend-engineer.md",
				RelativePath:    "agents/the-frontend-engineer.md",
				Size:            1536,
				Category:        uninstaller.CategoryAgent,
				IsTrackedInLock: true,
			},
			{
				Path:            "/Users/username/.claude/commands/s/specify.md",
				RelativePath:    "commands/s/specify.md",
				Size:            1024,
				Category:        uninstaller.CategoryCommand,
				IsTrackedInLock: true,
			},
			{
				Path:            "/Users/username/.config/the-startup/bin/the-startup",
				RelativePath:    "bin/the-startup",
				Size:            8388608, // 8MB
				Category:        uninstaller.CategoryBinary,
				IsTrackedInLock: true,
			},
		},
		CategorySummary: []uninstaller.CategorySummary{
			{
				Category:       uninstaller.CategoryAgent,
				Count:          2,
				TotalSize:      3584,
				TrackedFiles:   2,
				UntrackedFiles: 0,
				ModifiedFiles:  0,
			},
			{
				Category:       uninstaller.CategoryCommand,
				Count:          1,
				TotalSize:      1024,
				TrackedFiles:   1,
				UntrackedFiles: 0,
				ModifiedFiles:  0,
			},
			{
				Category:       uninstaller.CategoryBinary,
				Count:          1,
				TotalSize:      8388608,
				TrackedFiles:   1,
				UntrackedFiles: 0,
				ModifiedFiles:  0,
			},
		},
		TotalFiles:       4,
		TotalSize:        8393216, // Sum of all file sizes
		UntrackedFiles:   []uninstaller.FileInfo{},
		OrphanedFiles:    []uninstaller.FileInfo{},
		SettingsFiles:    []uninstaller.FileInfo{},
		SecurityIssues:   []uninstaller.SecurityIssue{},
		ValidationErrors: []uninstaller.ValidationError{},
	}
}

// DemoRemovalPreview shows how to create and use a RemovalPreviewModel
func DemoRemovalPreview() {
	// Create a mock preview
	preview := CreateMockPreview()

	// Create the preview model
	previewModel := NewRemovalPreviewModel(preview, false)

	// Set a reasonable size
	previewModel.SetSize(120, 30)

	// The model is now ready to be used in a TUI application
	fmt.Printf("Created RemovalPreviewModel with %d files to remove\n", preview.TotalFiles)
	fmt.Printf("Total size: %s\n", formatBytes(preview.TotalSize))

	// The model can be integrated into a larger TUI application or used standalone
	// For example, it could be embedded in the main UninstallModel during the preview state
}

// Integration example showing how UninstallModel and RemovalPreviewModel work together
func IntegrationExample() {
	fmt.Println("=== Uninstall UI Models Integration Example ===")
	fmt.Println()

	fmt.Println("1. Creating UninstallModel with various options...")
	uninstallModel := NewUninstallModel(true, false, true, true) // dry-run with backup and verbose
	fmt.Printf("   - Initial state: %s\n", uninstallModel.state.String())
	fmt.Printf("   - Options: dry-run=%v, force=%v, backup=%v, verbose=%v\n",
		uninstallModel.dryRun, uninstallModel.forceRemove, uninstallModel.createBackup, uninstallModel.verbose)

	fmt.Println()
	fmt.Println("2. Creating RemovalPreviewModel with mock data...")
	preview := CreateMockPreview()
	previewModel := NewRemovalPreviewModel(preview, true)
	previewModel.SetSize(80, 24) // Set a default size
	fmt.Printf("   - Preview for: %s\n", preview.InstallPath)
	fmt.Printf("   - Files to remove: %d\n", preview.TotalFiles)
	fmt.Printf("   - Categories: %d\n", len(preview.CategorySummary))
	fmt.Printf("   - Preview model dimensions: %dx%d\n", previewModel.width, previewModel.height)

	fmt.Println()
	fmt.Println("3. Integration flow:")
	fmt.Println("   - UninstallModel starts in Init state")
	fmt.Println("   - Transitions to PathDiscovery → PreviewGeneration → PreviewDisplay")
	fmt.Println("   - During PreviewDisplay, RemovalPreviewModel is embedded")
	fmt.Println("   - User can navigate preview, toggle details, see security issues")
	fmt.Println("   - After confirmation, execution begins with progress tracking")
	fmt.Println("   - Finally shows completion state with results")

	fmt.Println()
	fmt.Println("4. Key Features:")
	fmt.Println("   ✓ Progressive state machine with clear transitions")
	fmt.Println("   ✓ Comprehensive removal preview with categorization")
	fmt.Println("   ✓ Safety features (dry-run, confirmation, backup)")
	fmt.Println("   ✓ Real-time progress tracking during execution")
	fmt.Println("   ✓ Error handling and recovery guidance")
	fmt.Println("   ✓ Responsive layout adapting to terminal size")
	fmt.Println("   ✓ Keyboard navigation and intuitive controls")

	fmt.Println()
	fmt.Println("The models are ready for integration with the uninstall command!")
}