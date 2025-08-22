package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rsmdt/the-startup/internal/uninstaller"
)

// UninstallOptions holds the configuration options for the uninstall process
type UninstallOptions struct {
	DryRun       bool
	Force        bool
	KeepLogs     bool
	KeepSettings bool
}

// UninstallState represents the current state of the uninstaller UI
type UninstallState int

const (
	UninstallStateInit UninstallState = iota
	UninstallStatePathDiscovery
	UninstallStatePreviewGeneration
	UninstallStatePreviewDisplay
	UninstallStateConfirmation
	UninstallStateExecuting
	UninstallStateComplete
	UninstallStateError
)

// String returns the string representation of the uninstall state
func (s UninstallState) String() string {
	switch s {
	case UninstallStateInit:
		return "Initializing"
	case UninstallStatePathDiscovery:
		return "Discovering Paths"
	case UninstallStatePreviewGeneration:
		return "Generating Preview"
	case UninstallStatePreviewDisplay:
		return "Preview Display"
	case UninstallStateConfirmation:
		return "Awaiting Confirmation"
	case UninstallStateExecuting:
		return "Executing Removal"
	case UninstallStateComplete:
		return "Complete"
	case UninstallStateError:
		return "Error"
	default:
		return "Unknown"
	}
}

// UninstallMessages define the messages used in the uninstall UI
type PathDiscoveryCompleteMsg struct {
	InstallPath string
	ClaudePath  string
	Source      uninstaller.DiscoverySource
	Error       error
}

type PreviewGenerationCompleteMsg struct {
	Preview *uninstaller.RemovalPreview
	Error   error
}

type UninstallExecutionCompleteMsg struct {
	Plan   *uninstaller.UninstallPlan
	Errors []uninstaller.RemovalError
	Error  error
}

type ProgressUpdateMsg struct {
	Current int
	Total   int
	Message string
}

// UninstallModel orchestrates the complete uninstall process
type UninstallModel struct {
	// UI state
	styles   Styles
	renderer *ProgressiveDisclosureRenderer
	state    UninstallState
	width    int
	height   int

	// Uninstaller and configuration
	uninstaller  *uninstaller.Uninstaller
	dryRun       bool
	forceRemove  bool
	createBackup bool
	verbose      bool

	// State data
	installPath string
	claudePath  string
	source      uninstaller.DiscoverySource
	preview     *uninstaller.RemovalPreview
	plan        *uninstaller.UninstallPlan
	confirmed   bool
	ready       bool

	// Error handling
	err error

	// Progress tracking
	spinner  spinner.Model
	progress progress.Model

	// Sub-models
	previewModel *RemovalPreviewModel

	// Execution tracking
	executionStarted  time.Time
	filesRemoved      int
	totalFiles        int
	currentOperation  string
}

// NewUninstallModel creates a new uninstall model
func NewUninstallModel(dryRun, forceRemove, createBackup, verbose bool) *UninstallModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))

	p := progress.New(progress.WithDefaultGradient())
	p.Width = 60

	uninstallerInstance := uninstaller.NewWithDefaults()
	uninstallerInstance.SetOptions(dryRun, forceRemove, createBackup, verbose)

	return &UninstallModel{
		styles:       GetStyles(),
		renderer:     NewProgressiveDisclosureRenderer(),
		state:        UninstallStateInit,
		width:        80,
		height:       24,
		uninstaller:  uninstallerInstance,
		dryRun:       dryRun,
		forceRemove:  forceRemove,
		createBackup: createBackup,
		verbose:      verbose,
		confirmed:    false,
		ready:        false,
		spinner:      s,
		progress:     p,
		previewModel: nil,
	}
}

// Init initializes the uninstall model and starts path discovery
func (m *UninstallModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.discoverPaths(),
	)
}

// Update handles messages and state transitions
func (m *UninstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = min(60, msg.Width-10)
		if m.previewModel != nil {
			m.previewModel.SetSize(msg.Width, msg.Height-10) // Reserve space for header
		}
		return m, nil

	case spinner.TickMsg:
		newSpinner, cmd := m.spinner.Update(msg)
		m.spinner = newSpinner
		cmds = append(cmds, cmd)

	case PathDiscoveryCompleteMsg:
		return m.handlePathDiscoveryComplete(msg)

	case PreviewGenerationCompleteMsg:
		return m.handlePreviewGenerationComplete(msg)

	case UninstallExecutionCompleteMsg:
		return m.handleExecutionComplete(msg)

	case ProgressUpdateMsg:
		m.filesRemoved = msg.Current
		m.totalFiles = msg.Total
		m.currentOperation = msg.Message
		if m.totalFiles > 0 {
			progressCmd := m.progress.SetPercent(float64(m.filesRemoved) / float64(m.totalFiles))
			cmds = append(cmds, progressCmd)
		}
	}

	// Update sub-models
	if m.previewModel != nil && m.state == UninstallStatePreviewDisplay {
		newPreviewModel, cmd := m.previewModel.Update(msg)
		if previewModel, ok := newPreviewModel.(*RemovalPreviewModel); ok {
			m.previewModel = previewModel
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// handleKeyPress processes key press events based on current state
func (m *UninstallModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if m.state == UninstallStateExecuting {
			// Don't allow quitting during execution for safety
			return m, nil
		}
		m.ready = true
		return m, tea.Quit

	case "esc":
		if m.state == UninstallStateExecuting {
			// Don't allow going back during execution
			return m, nil
		}
		switch m.state {
		case UninstallStatePreviewDisplay, UninstallStateConfirmation:
			// Go back to preview generation
			m.state = UninstallStatePreviewGeneration
			return m, m.generatePreview()
		default:
			m.ready = true
			return m, tea.Quit
		}

	case "enter":
		return m.handleEnterKey()

	case "y", "Y":
		if m.state == UninstallStateConfirmation {
			m.confirmed = true
			m.state = UninstallStateExecuting
			m.executionStarted = time.Now()
			return m, m.executeUninstall()
		}

	case "n", "N":
		if m.state == UninstallStateConfirmation {
			m.ready = true
			return m, tea.Quit
		}

	case "d", "D":
		// Toggle dry-run mode if not already executing
		if m.state != UninstallStateExecuting && m.state != UninstallStateComplete {
			m.dryRun = !m.dryRun
			m.uninstaller.SetOptions(m.dryRun, m.forceRemove, m.createBackup, m.verbose)
		}
		return m, nil
	}

	return m, nil
}

// handleEnterKey processes the Enter key based on current state
func (m *UninstallModel) handleEnterKey() (tea.Model, tea.Cmd) {
	switch m.state {
	case UninstallStatePreviewDisplay:
		// Move to confirmation
		m.state = UninstallStateConfirmation
		return m, nil

	case UninstallStateComplete:
		m.ready = true
		return m, tea.Quit

	case UninstallStateError:
		// Restart from beginning
		m.state = UninstallStateInit
		m.err = nil
		return m, tea.Batch(
			m.spinner.Tick,
			m.discoverPaths(),
		)
	}

	return m, nil
}

// discoverPaths starts path discovery in the background
func (m *UninstallModel) discoverPaths() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		// We need to extract the paths from the plan
		plan, err := m.uninstaller.CreateUninstallPlan()
		if err != nil {
			return PathDiscoveryCompleteMsg{Error: err}
		}

		return PathDiscoveryCompleteMsg{
			InstallPath: plan.InstallPath,
			ClaudePath:  plan.ClaudePath,
			Source:      plan.DiscoverySource,
			Error:       nil,
		}
	})
}

// generatePreview starts preview generation in the background
func (m *UninstallModel) generatePreview() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		preview, err := m.uninstaller.GenerateDetailedPreview()
		return PreviewGenerationCompleteMsg{
			Preview: preview,
			Error:   err,
		}
	})
}

// executeUninstall starts the uninstall execution in the background
func (m *UninstallModel) executeUninstall() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		plan, err := m.uninstaller.CreateUninstallPlan()
		if err != nil {
			return UninstallExecutionCompleteMsg{Error: err}
		}

		// Configure the uninstaller with discovered paths and lockfile
		m.uninstaller.SetPathsAndLockFile(plan.InstallPath, plan.ClaudePath, plan.LockFile)

		execErr := m.uninstaller.ExecuteUninstallPlan(plan)
		return UninstallExecutionCompleteMsg{
			Plan:   plan,
			Errors: plan.RemovalErrors,
			Error:  execErr,
		}
	})
}

// handlePathDiscoveryComplete processes path discovery completion
func (m *UninstallModel) handlePathDiscoveryComplete(msg PathDiscoveryCompleteMsg) (tea.Model, tea.Cmd) {
	if msg.Error != nil {
		m.err = msg.Error
		m.state = UninstallStateError
		return m, nil
	}

	m.installPath = msg.InstallPath
	m.claudePath = msg.ClaudePath
	m.source = msg.Source
	m.state = UninstallStatePreviewGeneration

	return m, m.generatePreview()
}

// handlePreviewGenerationComplete processes preview generation completion
func (m *UninstallModel) handlePreviewGenerationComplete(msg PreviewGenerationCompleteMsg) (tea.Model, tea.Cmd) {
	if msg.Error != nil {
		m.err = msg.Error
		m.state = UninstallStateError
		return m, nil
	}

	m.preview = msg.Preview
	m.previewModel = NewRemovalPreviewModel(msg.Preview, m.dryRun)
	m.previewModel.SetSize(m.width, m.height-10) // Reserve space for header
	m.state = UninstallStatePreviewDisplay

	return m, nil
}

// handleExecutionComplete processes execution completion
func (m *UninstallModel) handleExecutionComplete(msg UninstallExecutionCompleteMsg) (tea.Model, tea.Cmd) {
	m.plan = msg.Plan
	
	if msg.Error != nil {
		m.err = msg.Error
		m.state = UninstallStateError
		return m, nil
	}

	m.state = UninstallStateComplete
	return m, nil
}

// View renders the appropriate view based on current state
func (m *UninstallModel) View() string {
	switch m.state {
	case UninstallStateInit:
		return m.renderInit()
	case UninstallStatePathDiscovery:
		return m.renderPathDiscovery()
	case UninstallStatePreviewGeneration:
		return m.renderPreviewGeneration()
	case UninstallStatePreviewDisplay:
		return m.renderPreviewDisplay()
	case UninstallStateConfirmation:
		return m.renderConfirmation()
	case UninstallStateExecuting:
		return m.renderExecution()
	case UninstallStateComplete:
		return m.renderComplete()
	case UninstallStateError:
		return m.renderError()
	default:
		return "Unknown state"
	}
}

// renderInit renders the initialization view
func (m *UninstallModel) renderInit() string {
	var s strings.Builder
	s.WriteString(m.styles.Title.Render("The (Agentic) Startup - Uninstaller"))
	s.WriteString("\n\n")
	s.WriteString(m.renderer.RenderTitle("Initializing uninstaller..."))
	s.WriteString(m.spinner.View() + " Preparing to analyze installation")
	s.WriteString("\n\n")
	s.WriteString(m.renderOptions())
	return s.String()
}

// renderPathDiscovery renders the path discovery view
func (m *UninstallModel) renderPathDiscovery() string {
	var s strings.Builder
	s.WriteString(m.styles.Title.Render("The (Agentic) Startup - Uninstaller"))
	s.WriteString("\n\n")
	s.WriteString(m.renderer.RenderTitle("Discovering installation paths..."))
	s.WriteString(m.spinner.View() + " Scanning for installation directories")
	s.WriteString("\n\n")
	s.WriteString(m.renderOptions())
	return s.String()
}

// renderPreviewGeneration renders the preview generation view
func (m *UninstallModel) renderPreviewGeneration() string {
	var s strings.Builder
	s.WriteString(m.styles.Title.Render("The (Agentic) Startup - Uninstaller"))
	s.WriteString("\n\n")

	// Show discovered paths
	s.WriteString(m.renderer.RenderTitle("Analyzing installation..."))
	s.WriteString(fmt.Sprintf("Install Path: %s\n", m.installPath))
	s.WriteString(fmt.Sprintf("Claude Path: %s\n", m.claudePath))
	s.WriteString(fmt.Sprintf("Discovery Source: %s\n", m.source.String()))
	s.WriteString("\n")
	
	s.WriteString(m.spinner.View() + " Generating detailed removal preview")
	s.WriteString("\n\n")
	s.WriteString(m.renderOptions())
	return s.String()
}

// renderPreviewDisplay renders the preview display view
func (m *UninstallModel) renderPreviewDisplay() string {
	var s strings.Builder
	s.WriteString(m.styles.Title.Render("The (Agentic) Startup - Uninstaller"))
	s.WriteString("\n\n")
	
	if m.previewModel != nil {
		s.WriteString(m.previewModel.View())
	} else {
		s.WriteString("Error: Preview not available")
	}
	
	s.WriteString("\n\n")
	s.WriteString(m.renderer.RenderHelp("Enter: continue • d: toggle dry-run • Escape: go back • Ctrl+C: quit"))
	return s.String()
}

// renderConfirmation renders the confirmation view
func (m *UninstallModel) renderConfirmation() string {
	var s strings.Builder
	s.WriteString(m.styles.Title.Render("The (Agentic) Startup - Uninstaller"))
	s.WriteString("\n\n")
	
	if m.dryRun {
		s.WriteString(m.styles.Info.Render("🔍 DRY RUN MODE - No files will actually be removed"))
		s.WriteString("\n\n")
		s.WriteString(m.renderer.RenderTitle("Confirm dry run execution"))
		s.WriteString("This will simulate the uninstall process without making any changes.\n")
	} else {
		s.WriteString(m.styles.Warning.Render("⚠️  WARNING: This will permanently remove files!"))
		s.WriteString("\n\n")
		s.WriteString(m.renderer.RenderTitle("Confirm uninstall execution"))
		s.WriteString("This will permanently remove the-startup from your system.\n")
	}

	// Show summary
	if m.preview != nil {
		s.WriteString(fmt.Sprintf("Files to remove: %d\n", m.preview.TotalFiles))
		s.WriteString(fmt.Sprintf("Total size: %s\n", formatBytes(m.preview.TotalSize)))
		
		if m.createBackup && !m.dryRun {
			s.WriteString(m.styles.Info.Render("📦 Backup will be created before removal"))
			s.WriteString("\n")
		}
	}

	s.WriteString("\n")
	s.WriteString(m.styles.Help.Render("Are you sure you want to continue?"))
	s.WriteString("\n\n")
	s.WriteString(m.renderer.RenderHelp("Y: yes, proceed • N: no, cancel • Escape: go back"))
	return s.String()
}

// renderExecution renders the execution view
func (m *UninstallModel) renderExecution() string {
	var s strings.Builder
	s.WriteString(m.styles.Title.Render("The (Agentic) Startup - Uninstaller"))
	s.WriteString("\n\n")
	
	if m.dryRun {
		s.WriteString(m.renderer.RenderTitle("Dry run execution in progress..."))
		s.WriteString(m.styles.Info.Render("🔍 Simulating removal process"))
	} else {
		s.WriteString(m.renderer.RenderTitle("Uninstall execution in progress..."))
		s.WriteString(m.styles.Warning.Render("🗑️  Removing files"))
	}
	
	s.WriteString("\n\n")
	
	// Progress bar
	s.WriteString(m.progress.View())
	s.WriteString("\n\n")
	
	// Current operation
	if m.currentOperation != "" {
		s.WriteString(m.currentOperation)
		s.WriteString("\n")
	}
	
	// Progress stats
	if m.totalFiles > 0 {
		s.WriteString(fmt.Sprintf("Progress: %d/%d files", m.filesRemoved, m.totalFiles))
		s.WriteString("\n")
	}
	
	// Elapsed time
	elapsed := time.Since(m.executionStarted)
	s.WriteString(fmt.Sprintf("Elapsed: %s", elapsed.Round(time.Second)))
	s.WriteString("\n\n")
	
	s.WriteString(m.styles.Help.Render("Please wait... (execution cannot be interrupted)"))
	return s.String()
}

// renderComplete renders the completion view
func (m *UninstallModel) renderComplete() string {
	var s strings.Builder
	s.WriteString(m.styles.Title.Render("The (Agentic) Startup - Uninstaller"))
	s.WriteString("\n\n")
	
	if m.dryRun {
		s.WriteString(m.styles.Success.Render("✓ Dry run completed successfully"))
		s.WriteString("\n\n")
		s.WriteString(m.renderer.RenderTitle("Dry run results"))
		s.WriteString("No files were actually removed. This was a simulation.\n")
	} else {
		s.WriteString(m.styles.Success.Render("✓ Uninstall completed successfully"))
		s.WriteString("\n\n")
		s.WriteString(m.renderer.RenderTitle("Uninstall results"))
		s.WriteString("The (Agentic) Startup has been removed from your system.\n")
	}
	
	if m.plan != nil {
		successCount := len(m.plan.FilesToRemove) - len(m.plan.RemovalErrors)
		s.WriteString(fmt.Sprintf("Files processed: %d\n", len(m.plan.FilesToRemove)))
		s.WriteString(fmt.Sprintf("Successful removals: %d\n", successCount))
		
		if len(m.plan.RemovalErrors) > 0 {
			s.WriteString(m.styles.Warning.Render(fmt.Sprintf("Errors encountered: %d\n", len(m.plan.RemovalErrors))))
		}
		
		if m.plan.BackupCreated {
			s.WriteString(fmt.Sprintf("Backup created: %s\n", m.plan.BackupPath))
		}
	}
	
	s.WriteString("\n")
	s.WriteString(m.renderer.RenderHelp("Enter: exit • Ctrl+C: quit"))
	return s.String()
}

// renderError renders the error view
func (m *UninstallModel) renderError() string {
	var s strings.Builder
	s.WriteString(m.styles.Title.Render("The (Agentic) Startup - Uninstaller"))
	s.WriteString("\n\n")
	
	s.WriteString(m.styles.Error.Render("❌ Uninstall failed"))
	s.WriteString("\n\n")
	s.WriteString(m.renderer.RenderError(m.err, "during uninstall"))
	s.WriteString("\n")
	
	s.WriteString("The uninstall process encountered an error and could not complete.\n")
	s.WriteString("Your system has not been modified.\n\n")
	
	s.WriteString(m.renderer.RenderHelp("Enter: retry • Escape: quit"))
	return s.String()
}

// renderOptions renders current options
func (m *UninstallModel) renderOptions() string {
	var options []string
	
	if m.dryRun {
		options = append(options, m.styles.Info.Render("DRY RUN"))
	}
	
	if m.forceRemove {
		options = append(options, m.styles.Warning.Render("FORCE"))
	}
	
	if m.createBackup {
		options = append(options, m.styles.Success.Render("BACKUP"))
	}
	
	if m.verbose {
		options = append(options, m.styles.Normal.Render("VERBOSE"))
	}
	
	if len(options) > 0 {
		return m.styles.Help.Render("Options: " + strings.Join(options, " • "))
	}
	
	return ""
}

// Ready returns whether the model is ready to quit
func (m *UninstallModel) Ready() bool {
	return m.ready
}

// formatBytes formats byte count as human-readable string
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// NOTE: This old uninstall model is kept for reference but the new simplified
// uninstall flow using the MainModel is in model.go with RunMainUninstaller