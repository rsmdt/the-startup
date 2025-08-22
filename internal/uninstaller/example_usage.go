package uninstaller

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ExamplePreviewUsage demonstrates how to use the removal preview functionality
func ExamplePreviewUsage() {
	// Create a path discoverer to find the installation
	pathDiscoverer := NewPathDiscovery()

	// Create a preview generator
	previewGenerator := NewPreviewGenerator(pathDiscoverer)
	previewGenerator.SetVerbose(true)

	// Generate comprehensive preview
	preview, err := previewGenerator.GeneratePreview()
	if err != nil {
		fmt.Printf("Error generating preview: %v\n", err)
		return
	}

	// Display basic information
	fmt.Printf("Installation Analysis Complete\n")
	fmt.Printf("================================\n")
	fmt.Printf("Install Path: %s\n", preview.InstallPath)
	fmt.Printf("Claude Path: %s\n", preview.ClaudePath)
	fmt.Printf("Discovery Source: %s\n", preview.DiscoverySource.String())
	fmt.Printf("Total Files: %d\n", preview.TotalFiles)
	fmt.Printf("Total Size: %d bytes (%.2f MB)\n", preview.TotalSize, float64(preview.TotalSize)/1024/1024)

	if preview.LockFile != nil {
		fmt.Printf("Lockfile Version: %s\n", preview.LockFile.Version)
		fmt.Printf("Install Date: %s\n", preview.LockFile.InstallDate)
	}
	fmt.Printf("\n")

	// Display category summary
	if len(preview.CategorySummary) > 0 {
		fmt.Printf("File Categories\n")
		fmt.Printf("===============\n")
		for _, summary := range preview.CategorySummary {
			fmt.Printf("%-12s: %d files, %.2f MB\n",
				summary.Category.String(),
				summary.Count,
				float64(summary.TotalSize)/1024/1024,
			)
			if summary.UntrackedFiles > 0 {
				fmt.Printf("  └─ %d untracked files\n", summary.UntrackedFiles)
			}
			if summary.ModifiedFiles > 0 {
				fmt.Printf("  └─ %d modified files\n", summary.ModifiedFiles)
			}
		}
		fmt.Printf("\n")
	}

	// Display untracked files
	if len(preview.UntrackedFiles) > 0 {
		fmt.Printf("Untracked Files (not in lockfile)\n")
		fmt.Printf("==================================\n")
		for _, file := range preview.UntrackedFiles {
			fmt.Printf("• %s (%s)\n", file.RelativePath, file.Category.String())
		}
		fmt.Printf("\n")
	}

	// Display orphaned files
	if len(preview.OrphanedFiles) > 0 {
		fmt.Printf("Orphaned Files (in lockfile but missing)\n")
		fmt.Printf("========================================\n")
		for _, file := range preview.OrphanedFiles {
			fmt.Printf("• %s\n", file.RelativePath)
		}
		fmt.Printf("\n")
	}

	// Display security issues
	if len(preview.SecurityIssues) > 0 {
		fmt.Printf("Security Issues\n")
		fmt.Printf("===============\n")
		for _, issue := range preview.SecurityIssues {
			fmt.Printf("🚨 %s: %s (%s)\n", issue.Type, issue.Description, issue.Severity)
			fmt.Printf("   File: %s\n", issue.FilePath)
		}
		fmt.Printf("\n")
	}

	// Display validation errors
	if len(preview.ValidationErrors) > 0 {
		fmt.Printf("Validation Errors\n")
		fmt.Printf("=================\n")
		for _, err := range preview.ValidationErrors {
			fmt.Printf("❌ %s: %s\n", err.Type, err.Description)
			fmt.Printf("   File: %s\n", err.FilePath)
		}
		fmt.Printf("\n")
	}

	// Display detailed file list if requested
	if len(preview.Files) <= 20 { // Only show if not too many files
		fmt.Printf("Files to Remove\n")
		fmt.Printf("===============\n")
		for _, file := range preview.Files {
			status := "✓"
			if file.PermissionIssue {
				status = "❌"
			} else if file.IsModified {
				status = "⚠️"
			}

			fmt.Printf("%s %s (%s, %d bytes)\n",
				status,
				file.RelativePath,
				file.Category.String(),
				file.Size,
			)

			if file.PermissionIssue {
				fmt.Printf("  └─ Permission issue: %s\n", file.PermissionError)
			}
			if file.IsModified {
				fmt.Printf("  └─ Modified since installation\n")
			}
			if file.IsSymlink {
				fmt.Printf("  └─ Symlink to: %s\n", file.SymlinkTarget)
			}
		}
	}

	// Export detailed preview as JSON (optional)
	if os.Getenv("EXPORT_PREVIEW") == "1" {
		jsonData, err := json.MarshalIndent(preview, "", "  ")
		if err == nil {
			outputPath := filepath.Join(preview.InstallPath, "removal_preview.json")
			if err := os.WriteFile(outputPath, jsonData, 0644); err == nil {
				fmt.Printf("Detailed preview exported to: %s\n", outputPath)
			}
		}
	}
}

// ExampleIntegratedUsage shows how to use the uninstaller with the new preview functionality
func ExampleIntegratedUsage() {
	// Create uninstaller with default dependencies
	uninstaller := NewWithDefaults()
	uninstaller.SetOptions(
		true,  // dryRun - don't actually remove files
		false, // forceRemove
		true,  // createBackup
		true,  // verbose
	)

	// Generate detailed preview
	preview, err := uninstaller.GenerateDetailedPreview()
	if err != nil {
		fmt.Printf("Error generating detailed preview: %v\n", err)
		return
	}

	fmt.Printf("Generated detailed preview with %d files\n", len(preview.Files))

	// Create legacy uninstall plan (for compatibility)
	plan, err := uninstaller.CreateUninstallPlan()
	if err != nil {
		fmt.Printf("Error creating uninstall plan: %v\n", err)
		return
	}

	fmt.Printf("Legacy plan shows %d files to remove\n", len(plan.FilesToRemove))

	// Execute the plan (dry run)
	if err := uninstaller.ExecuteUninstallPlan(plan); err != nil {
		fmt.Printf("Error executing uninstall plan: %v\n", err)
	} else {
		fmt.Printf("Dry run completed successfully\n")
	}
}