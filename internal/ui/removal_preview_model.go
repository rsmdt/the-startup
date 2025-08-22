package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rsmdt/the-startup/internal/uninstaller"
)

// RemovalPreviewModel displays detailed removal preview with file categorization
type RemovalPreviewModel struct {
	// UI components
	styles   Styles
	renderer *ProgressiveDisclosureRenderer
	viewport viewport.Model

	// Data
	preview *uninstaller.RemovalPreview
	dryRun  bool

	// Display options
	showDetails          bool
	selectedCategory     uninstaller.FileCategory
	showSecurityIssues   bool
	showValidationErrors bool
	showUntrackedFiles   bool
	showOrphanedFiles    bool

	// Dimensions
	width  int
	height int

	// Ready state
	ready bool
}

// NewRemovalPreviewModel creates a new removal preview model
func NewRemovalPreviewModel(preview *uninstaller.RemovalPreview, dryRun bool) *RemovalPreviewModel {
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#606060")).
		Padding(1)

	model := &RemovalPreviewModel{
		styles:               GetStyles(),
		renderer:             NewProgressiveDisclosureRenderer(),
		viewport:             vp,
		preview:              preview,
		dryRun:               dryRun,
		showDetails:          false,
		selectedCategory:     uninstaller.CategoryAgent, // Default to first category
		showSecurityIssues:   len(preview.SecurityIssues) > 0,
		showValidationErrors: len(preview.ValidationErrors) > 0,
		showUntrackedFiles:   len(preview.UntrackedFiles) > 0,
		showOrphanedFiles:    len(preview.OrphanedFiles) > 0,
		width:                80,
		height:               24,
	}

	model.updateContent()
	return model
}

// Init initializes the removal preview model
func (m *RemovalPreviewModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the removal preview model
func (m *RemovalPreviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		case "pgup":
			m.viewport.HalfViewUp()
		case "pgdown":
			m.viewport.HalfViewDown()
		case "home":
			m.viewport.GotoTop()
		case "end":
			m.viewport.GotoBottom()
		case "tab":
			// Cycle through categories
			m.cycleCategory()
			m.updateContent()
		case "d", "D":
			// Toggle details view
			m.showDetails = !m.showDetails
			m.updateContent()
		case "s", "S":
			// Toggle security issues display
			m.showSecurityIssues = !m.showSecurityIssues
			m.updateContent()
		case "v", "V":
			// Toggle validation errors display
			m.showValidationErrors = !m.showValidationErrors
			m.updateContent()
		case "u", "U":
			// Toggle untracked files display
			m.showUntrackedFiles = !m.showUntrackedFiles
			m.updateContent()
		case "o", "O":
			// Toggle orphaned files display
			m.showOrphanedFiles = !m.showOrphanedFiles
			m.updateContent()
		}

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// SetSize updates the model dimensions
func (m *RemovalPreviewModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	
	// Reserve space for header and help
	viewportHeight := height - 6
	if viewportHeight < 10 {
		viewportHeight = 10
	}
	
	viewportWidth := width - 4
	if viewportWidth < 40 {
		viewportWidth = 40
	}
	
	m.viewport.Width = viewportWidth
	m.viewport.Height = viewportHeight
	
	m.updateContent()
}

// cycleCategory moves to the next category with files
func (m *RemovalPreviewModel) cycleCategory() {
	categories := m.getCategoriesWithFiles()
	if len(categories) == 0 {
		return
	}
	
	// Find current category index
	currentIndex := -1
	for i, cat := range categories {
		if cat == m.selectedCategory {
			currentIndex = i
			break
		}
	}
	
	// Move to next category
	nextIndex := (currentIndex + 1) % len(categories)
	m.selectedCategory = categories[nextIndex]
}

// getCategoriesWithFiles returns categories that have files
func (m *RemovalPreviewModel) getCategoriesWithFiles() []uninstaller.FileCategory {
	categoryMap := make(map[uninstaller.FileCategory]bool)
	
	for _, file := range m.preview.Files {
		categoryMap[file.Category] = true
	}
	
	var categories []uninstaller.FileCategory
	for category := range categoryMap {
		categories = append(categories, category)
	}
	
	// Sort for consistent order
	sort.Slice(categories, func(i, j int) bool {
		return categories[i] < categories[j]
	})
	
	return categories
}

// updateContent refreshes the viewport content
func (m *RemovalPreviewModel) updateContent() {
	content := m.renderContent()
	m.viewport.SetContent(content)
}

// renderContent generates the full content for the viewport
func (m *RemovalPreviewModel) renderContent() string {
	var s strings.Builder

	// Header
	s.WriteString(m.renderHeader())
	s.WriteString("\n")

	// Summary section
	s.WriteString(m.renderSummary())
	s.WriteString("\n")

	// Category breakdown
	s.WriteString(m.renderCategoryBreakdown())
	s.WriteString("\n")

	// Detailed view if requested
	if m.showDetails {
		s.WriteString(m.renderDetailedView())
		s.WriteString("\n")
	}

	// Security issues
	if m.showSecurityIssues && len(m.preview.SecurityIssues) > 0 {
		s.WriteString(m.renderSecurityIssues())
		s.WriteString("\n")
	}

	// Validation errors
	if m.showValidationErrors && len(m.preview.ValidationErrors) > 0 {
		s.WriteString(m.renderValidationErrors())
		s.WriteString("\n")
	}

	// Untracked files
	if m.showUntrackedFiles && len(m.preview.UntrackedFiles) > 0 {
		s.WriteString(m.renderUntrackedFiles())
		s.WriteString("\n")
	}

	// Orphaned files
	if m.showOrphanedFiles && len(m.preview.OrphanedFiles) > 0 {
		s.WriteString(m.renderOrphanedFiles())
		s.WriteString("\n")
	}

	return s.String()
}

// renderHeader renders the preview header
func (m *RemovalPreviewModel) renderHeader() string {
	var s strings.Builder
	
	if m.dryRun {
		s.WriteString(m.styles.Info.Render("🔍 REMOVAL PREVIEW (DRY RUN)"))
	} else {
		s.WriteString(m.styles.Warning.Render("🗑️  REMOVAL PREVIEW"))
	}
	
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("Installation Path: %s\n", m.preview.InstallPath))
	s.WriteString(fmt.Sprintf("Claude Path: %s\n", m.preview.ClaudePath))
	s.WriteString(fmt.Sprintf("Discovery Source: %s\n", m.preview.DiscoverySource.String()))
	
	return s.String()
}

// renderSummary renders the removal summary
func (m *RemovalPreviewModel) renderSummary() string {
	var s strings.Builder
	
	s.WriteString(m.styles.Title.Render("📋 Removal Summary"))
	s.WriteString("\n")
	
	s.WriteString(fmt.Sprintf("Total Files: %d\n", m.preview.TotalFiles))
	s.WriteString(fmt.Sprintf("Total Size: %s\n", formatBytes(m.preview.TotalSize)))
	
	if m.preview.LockFile != nil {
		s.WriteString(fmt.Sprintf("Tracked Files: %d\n", len(m.preview.LockFile.Files)))
	}
	
	if len(m.preview.UntrackedFiles) > 0 {
		s.WriteString(m.styles.Warning.Render(fmt.Sprintf("Untracked Files: %d\n", len(m.preview.UntrackedFiles))))
	}
	
	if len(m.preview.OrphanedFiles) > 0 {
		s.WriteString(m.styles.Info.Render(fmt.Sprintf("Orphaned Files: %d\n", len(m.preview.OrphanedFiles))))
	}
	
	if len(m.preview.SecurityIssues) > 0 {
		s.WriteString(m.styles.Error.Render(fmt.Sprintf("Security Issues: %d\n", len(m.preview.SecurityIssues))))
	}
	
	if len(m.preview.ValidationErrors) > 0 {
		s.WriteString(m.styles.Error.Render(fmt.Sprintf("Validation Errors: %d\n", len(m.preview.ValidationErrors))))
	}
	
	return s.String()
}

// renderCategoryBreakdown renders the file category breakdown
func (m *RemovalPreviewModel) renderCategoryBreakdown() string {
	var s strings.Builder
	
	s.WriteString(m.styles.Title.Render("📂 Files by Category"))
	s.WriteString("\n")
	
	if len(m.preview.CategorySummary) == 0 {
		s.WriteString(m.styles.Help.Render("No files found to categorize"))
		return s.String()
	}
	
	// Sort categories by name for consistent display
	categories := make([]uninstaller.CategorySummary, len(m.preview.CategorySummary))
	copy(categories, m.preview.CategorySummary)
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Category < categories[j].Category
	})
	
	for _, summary := range categories {
		isSelected := summary.Category == m.selectedCategory
		prefix := "  "
		if isSelected {
			prefix = m.styles.Cursor.Render("> ")
		}
		
		categoryName := strings.Title(summary.Category.String())
		line := fmt.Sprintf("%s%s: %d files (%s)",
			prefix,
			categoryName,
			summary.Count,
			formatBytes(summary.TotalSize))
		
		if isSelected {
			line = m.styles.CursorLine.Render(line)
		} else {
			line = m.styles.Normal.Render(line)
		}
		
		s.WriteString(line)
		s.WriteString("\n")
		
		// Show additional details for selected category
		if isSelected && (summary.UntrackedFiles > 0 || summary.ModifiedFiles > 0) {
			details := []string{}
			if summary.TrackedFiles > 0 {
				details = append(details, fmt.Sprintf("%d tracked", summary.TrackedFiles))
			}
			if summary.UntrackedFiles > 0 {
				details = append(details, m.styles.Warning.Render(fmt.Sprintf("%d untracked", summary.UntrackedFiles)))
			}
			if summary.ModifiedFiles > 0 {
				details = append(details, m.styles.Info.Render(fmt.Sprintf("%d modified", summary.ModifiedFiles)))
			}
			
			if len(details) > 0 {
				s.WriteString(m.styles.Help.Render("    (" + strings.Join(details, ", ") + ")"))
				s.WriteString("\n")
			}
		}
	}
	
	return s.String()
}

// renderDetailedView renders detailed file listings for the selected category
func (m *RemovalPreviewModel) renderDetailedView() string {
	var s strings.Builder
	
	categoryName := strings.Title(m.selectedCategory.String())
	s.WriteString(m.styles.Title.Render(fmt.Sprintf("📄 %s Files (Detailed)", categoryName)))
	s.WriteString("\n")
	
	// Filter files by selected category
	categoryFiles := []uninstaller.FileInfo{}
	for _, file := range m.preview.Files {
		if file.Category == m.selectedCategory {
			categoryFiles = append(categoryFiles, file)
		}
	}
	
	if len(categoryFiles) == 0 {
		s.WriteString(m.styles.Help.Render("No files in this category"))
		return s.String()
	}
	
	// Sort files by path for consistent display
	sort.Slice(categoryFiles, func(i, j int) bool {
		return categoryFiles[i].RelativePath < categoryFiles[j].RelativePath
	})
	
	// Limit display to avoid overwhelming output
	maxFiles := 20
	for i, file := range categoryFiles {
		if i >= maxFiles {
			remaining := len(categoryFiles) - maxFiles
			s.WriteString(m.styles.Help.Render(fmt.Sprintf("... and %d more files", remaining)))
			s.WriteString("\n")
			break
		}
		
		s.WriteString(m.renderFileDetails(file))
		s.WriteString("\n")
	}
	
	return s.String()
}

// renderFileDetails renders detailed information about a single file
func (m *RemovalPreviewModel) renderFileDetails(file uninstaller.FileInfo) string {
	var s strings.Builder
	
	// File path and size
	s.WriteString(m.styles.Normal.Render(file.RelativePath))
	s.WriteString(m.styles.Help.Render(fmt.Sprintf(" (%s)", formatBytes(file.Size))))
	
	// Status indicators
	var indicators []string
	
	if !file.IsTrackedInLock {
		indicators = append(indicators, m.styles.Warning.Render("untracked"))
	}
	
	if file.IsModified {
		indicators = append(indicators, m.styles.Info.Render("modified"))
	}
	
	if file.PermissionIssue {
		indicators = append(indicators, m.styles.Error.Render("permission issue"))
	}
	
	if file.IsSymlink {
		indicators = append(indicators, m.styles.Info.Render("symlink"))
	}
	
	if len(indicators) > 0 {
		s.WriteString(" [" + strings.Join(indicators, ", ") + "]")
	}
	
	// Show modification time if available
	if !file.ModTime.IsZero() {
		modTime := file.ModTime.Format("2006-01-02 15:04:05")
		s.WriteString(m.styles.Help.Render(fmt.Sprintf("\n    Modified: %s", modTime)))
	}
	
	// Show symlink target if applicable
	if file.IsSymlink && file.SymlinkTarget != "" {
		s.WriteString(m.styles.Help.Render(fmt.Sprintf("\n    Target: %s", file.SymlinkTarget)))
	}
	
	// Show permission error if applicable
	if file.PermissionIssue && file.PermissionError != "" {
		s.WriteString(m.styles.Error.Render(fmt.Sprintf("\n    Error: %s", file.PermissionError)))
	}
	
	return s.String()
}

// renderSecurityIssues renders security issue warnings
func (m *RemovalPreviewModel) renderSecurityIssues() string {
	var s strings.Builder
	
	s.WriteString(m.styles.Error.Render("🛡️  Security Issues"))
	s.WriteString("\n")
	
	if len(m.preview.SecurityIssues) == 0 {
		s.WriteString(m.styles.Success.Render("✓ No security issues detected"))
		return s.String()
	}
	
	for _, issue := range m.preview.SecurityIssues {
		severity := m.renderSeverity(issue.Severity)
		s.WriteString(fmt.Sprintf("%s %s: %s\n", severity, issue.Type, issue.Description))
		if issue.FilePath != "" {
			s.WriteString(m.styles.Help.Render(fmt.Sprintf("  File: %s", issue.FilePath)))
			s.WriteString("\n")
		}
	}
	
	return s.String()
}

// renderValidationErrors renders validation error messages
func (m *RemovalPreviewModel) renderValidationErrors() string {
	var s strings.Builder
	
	s.WriteString(m.styles.Error.Render("⚠️  Validation Errors"))
	s.WriteString("\n")
	
	if len(m.preview.ValidationErrors) == 0 {
		s.WriteString(m.styles.Success.Render("✓ No validation errors detected"))
		return s.String()
	}
	
	for _, validationError := range m.preview.ValidationErrors {
		s.WriteString(m.styles.Error.Render(fmt.Sprintf("❌ %s: %s\n", validationError.Type, validationError.Description)))
		if validationError.FilePath != "" {
			s.WriteString(m.styles.Help.Render(fmt.Sprintf("  File: %s", validationError.FilePath)))
			s.WriteString("\n")
		}
	}
	
	return s.String()
}

// renderUntrackedFiles renders untracked files section
func (m *RemovalPreviewModel) renderUntrackedFiles() string {
	var s strings.Builder
	
	s.WriteString(m.styles.Warning.Render("📋 Untracked Files"))
	s.WriteString("\n")
	s.WriteString(m.styles.Help.Render("Files not present in the original installation lock file"))
	s.WriteString("\n\n")
	
	if len(m.preview.UntrackedFiles) == 0 {
		s.WriteString(m.styles.Success.Render("✓ No untracked files"))
		return s.String()
	}
	
	maxFiles := 10
	for i, file := range m.preview.UntrackedFiles {
		if i >= maxFiles {
			remaining := len(m.preview.UntrackedFiles) - maxFiles
			s.WriteString(m.styles.Help.Render(fmt.Sprintf("... and %d more untracked files", remaining)))
			s.WriteString("\n")
			break
		}
		
		s.WriteString(fmt.Sprintf("  %s (%s)\n", file.RelativePath, formatBytes(file.Size)))
	}
	
	return s.String()
}

// renderOrphanedFiles renders orphaned files section
func (m *RemovalPreviewModel) renderOrphanedFiles() string {
	var s strings.Builder
	
	s.WriteString(m.styles.Info.Render("👻 Orphaned Files"))
	s.WriteString("\n")
	s.WriteString(m.styles.Help.Render("Files listed in lock file but not found on disk"))
	s.WriteString("\n\n")
	
	if len(m.preview.OrphanedFiles) == 0 {
		s.WriteString(m.styles.Success.Render("✓ No orphaned files"))
		return s.String()
	}
	
	for _, file := range m.preview.OrphanedFiles {
		s.WriteString(fmt.Sprintf("  %s\n", file.RelativePath))
	}
	
	return s.String()
}

// renderSeverity renders security issue severity with appropriate styling
func (m *RemovalPreviewModel) renderSeverity(severity string) string {
	switch strings.ToLower(severity) {
	case "critical":
		return m.styles.Error.Render("🔴 CRITICAL")
	case "high":
		return m.styles.Error.Render("🟠 HIGH")
	case "medium":
		return m.styles.Warning.Render("🟡 MEDIUM")
	case "low":
		return m.styles.Info.Render("🟢 LOW")
	default:
		return m.styles.Help.Render("⚪ UNKNOWN")
	}
}

// View renders the removal preview
func (m *RemovalPreviewModel) View() string {
	var s strings.Builder
	
	// Viewport with scrollable content
	s.WriteString(m.viewport.View())
	s.WriteString("\n")
	
	// Help text
	s.WriteString(m.renderHelp())
	
	return s.String()
}

// renderHelp renders help text for navigation and controls
func (m *RemovalPreviewModel) renderHelp() string {
	var helpItems []string
	
	helpItems = append(helpItems, "↑↓ scroll")
	helpItems = append(helpItems, "Tab: cycle categories")
	
	if m.showDetails {
		helpItems = append(helpItems, "d: hide details")
	} else {
		helpItems = append(helpItems, "d: show details")
	}
	
	if len(m.preview.SecurityIssues) > 0 {
		if m.showSecurityIssues {
			helpItems = append(helpItems, "s: hide security")
		} else {
			helpItems = append(helpItems, "s: show security")
		}
	}
	
	if len(m.preview.ValidationErrors) > 0 {
		if m.showValidationErrors {
			helpItems = append(helpItems, "v: hide validation")
		} else {
			helpItems = append(helpItems, "v: show validation")
		}
	}
	
	if len(m.preview.UntrackedFiles) > 0 {
		if m.showUntrackedFiles {
			helpItems = append(helpItems, "u: hide untracked")
		} else {
			helpItems = append(helpItems, "u: show untracked")
		}
	}
	
	if len(m.preview.OrphanedFiles) > 0 {
		if m.showOrphanedFiles {
			helpItems = append(helpItems, "o: hide orphaned")
		} else {
			helpItems = append(helpItems, "o: show orphaned")
		}
	}
	
	return m.styles.Help.Render(strings.Join(helpItems, " • "))
}

// Ready returns whether the model is ready to proceed
func (m *RemovalPreviewModel) Ready() bool {
	return m.ready
}