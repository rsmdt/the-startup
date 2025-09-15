package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rsmdt/the-startup/internal/stats"
)

// PanelType represents different dashboard panels
type PanelType int

const (
	TimelinePanel PanelType = iota
	StatsPanel
	DelegationPanel
	CoOccurrencePanel
)

// TimeRange represents predefined time filter ranges
type TimeRange struct {
	Label  string
	Filter string
	Desc   string
}

// Predefined time ranges matching interface specification
var TimeRanges = []TimeRange{
	{Label: "1w", Filter: "7d", Desc: "Last 7 days"},
	{Label: "1m", Filter: "30d", Desc: "Last 30 days"},
	{Label: "3m", Filter: "90d", Desc: "Last 3 months"},
	{Label: "6m", Filter: "180d", Desc: "Last 6 months"},
	{Label: "1y", Filter: "365d", Desc: "Last year"},
	{Label: "all", Filter: "", Desc: "All time"},
}

// FilterMenuState represents the state of the filter menu
type FilterMenuState struct {
	IsOpen         bool
	SelectedIndex  int
	SelectedFilter TimeRange
}

// SortColumn represents the different columns we can sort by
type SortColumn int

const (
	SortByCount SortColumn = iota
	SortBySuccessRate
	SortByDuration
)

// SortDirection represents ascending or descending sort
type SortDirection int

const (
	SortAscending SortDirection = iota
	SortDescending
)

// DashboardModel represents the agent dashboard UI model
type DashboardModel struct {
	// Data
	agentStats         map[string]*stats.GlobalAgentStats
	delegationPatterns []stats.DelegationPattern
	coOccurrences      []stats.AgentCoOccurrence

	// Data loading state
	isLoading       bool
	lastLoadError   error
	timeFilter      string // Default "30d" for 1-month lookback
	logFiles        []string
	lastLoadTime    time.Time

	// Filter state
	currentTimeRange TimeRange
	filterMenu       FilterMenuState
	loadingProgress  string

	// UI State
	focusedPanel              PanelType
	selectedTimelineIndex     int
	selectedStatsIndex        int
	selectedDelegationIndex   int
	selectedCoOccurrenceIndex int
	selectedFiltersHelpIndex  int
	showHelp                  bool

	// Sort state for stats panel
	sortColumn    SortColumn
	sortDirection SortDirection

	// Leaderboard component
	leaderboard *LeaderboardComponent

	// Timeline component
	timeline *MessageTimelineComponent

	// Enhanced delegation and co-occurrence components
	delegationFilterOptions    DelegationFilterOptions
	delegationFlowComponent    *DelegationFlowComponent
	coOccurrenceMatrixComponent *CoOccurrenceMatrixComponent

	// Display state
	width  int
	height int

	// Theme state
	currentTheme     ColorTheme
	currentThemeName string
	styles           *DashboardStyles
}

// RefreshDataMsg is a message sent when data needs to be refreshed
type RefreshDataMsg struct{}

// DataLoadedMsg is sent when data loading completes
type DataLoadedMsg struct {
	AgentStats         map[string]*stats.GlobalAgentStats
	DelegationPatterns []stats.DelegationPattern
	CoOccurrences      []stats.AgentCoOccurrence
	MessageEvents      []MessageEvent
	Error              error
}

// LoadDataMsg is sent to trigger data loading
type LoadDataMsg struct {
	TimeFilter string
}

// TimeFilterChangedMsg is sent when time filter changes
type TimeFilterChangedMsg struct {
	TimeRange TimeRange
}

// FilterMenuToggleMsg is sent when filter menu is toggled
type FilterMenuToggleMsg struct{}

// LoadProgressMsg is sent during data loading to show progress
type LoadProgressMsg struct {
	Message string
}

// NewDashboardModel creates a new dashboard model with the provided data
func NewDashboardModel(
	agentStats map[string]*stats.GlobalAgentStats,
	delegationPatterns []stats.DelegationPattern,
	coOccurrences []stats.AgentCoOccurrence,
) *DashboardModel {
	return &DashboardModel{
		agentStats:                agentStats,
		delegationPatterns:        delegationPatterns,
		coOccurrences:             coOccurrences,
		timeFilter:                "30d", // Default 1-month lookback
		currentTimeRange:          TimeRanges[1], // Default to 1m (30d)
		filterMenu: FilterMenuState{
			IsOpen:         false,
			SelectedIndex:  1, // Default to 1m
			SelectedFilter: TimeRanges[1],
		},
		focusedPanel:              TimelinePanel,
		selectedTimelineIndex:     0,
		selectedStatsIndex:        0,
		selectedDelegationIndex:   0,
		selectedCoOccurrenceIndex: 0,
		selectedFiltersHelpIndex:  0,
		showHelp:                  false,
		sortColumn:                SortByCount,    // Default sort by count
		sortDirection:             SortDescending, // Default descending
		leaderboard:               NewLeaderboardComponent("Agent Statistics"),
		timeline:                  NewMessageTimelineComponent(),
		delegationFilterOptions: DelegationFilterOptions{
			SortBy:     SortByFrequency,
			MinCount:   0,
			MaxResults: 10,
		},
		delegationFlowComponent:     NewDelegationFlowComponent(),
		coOccurrenceMatrixComponent: NewCoOccurrenceMatrixComponent(),
		width:                       80,
		height:                      24,
		currentThemeName:            "dracula",
		currentTheme:                GetTheme("dracula"),
		styles:                      ApplyTheme(GetTheme("dracula")),
	}
}

// NewDashboardModelWithLoading creates a dashboard model that loads data using existing stats engine
func NewDashboardModelWithLoading() *DashboardModel {
	return &DashboardModel{
		agentStats:                make(map[string]*stats.GlobalAgentStats),
		delegationPatterns:        []stats.DelegationPattern{},
		coOccurrences:             []stats.AgentCoOccurrence{},
		timeFilter:                "30d", // Default 1-month lookback
		currentTimeRange:          TimeRanges[1], // Default to 1m (30d)
		filterMenu: FilterMenuState{
			IsOpen:         false,
			SelectedIndex:  1, // Default to 1m
			SelectedFilter: TimeRanges[1],
		},
		isLoading:                 false,
		focusedPanel:              TimelinePanel,
		selectedTimelineIndex:     0,
		selectedStatsIndex:        0,
		selectedDelegationIndex:   0,
		selectedCoOccurrenceIndex: 0,
		selectedFiltersHelpIndex:  0,
		showHelp:                  false,
		sortColumn:                SortByCount,    // Default sort by count
		sortDirection:             SortDescending, // Default descending
		leaderboard:               NewLeaderboardComponent("Agent Statistics"),
		timeline:                  NewMessageTimelineComponent(),
		delegationFilterOptions: DelegationFilterOptions{
			SortBy:     SortByFrequency,
			MinCount:   0,
			MaxResults: 10,
		},
		delegationFlowComponent:     NewDelegationFlowComponent(),
		coOccurrenceMatrixComponent: NewCoOccurrenceMatrixComponent(),
		width:                       80,
		height:                      24,
		currentThemeName:            "dracula",
		currentTheme:                GetTheme("dracula"),
		styles:                      ApplyTheme(GetTheme("dracula")),
	}
}

// Init initializes the dashboard model
func (m *DashboardModel) Init() tea.Cmd {
	// Load data immediately when dashboard starts
	m.isLoading = true
	m.loadingProgress = fmt.Sprintf("Loading data for %s...", m.currentTimeRange.Desc)
	return m.loadData(m.timeFilter)
}

// Update handles messages and updates the model state
func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case DataLoadedMsg:
		if msg.Error != nil {
			m.lastLoadError = msg.Error
			m.isLoading = false
			m.loadingProgress = ""
		} else {
			m.agentStats = msg.AgentStats
			m.delegationPatterns = msg.DelegationPatterns
			m.coOccurrences = msg.CoOccurrences
			m.timeline.SetEvents(msg.MessageEvents)
			m.lastLoadError = nil
			m.isLoading = false
			m.loadingProgress = ""
			m.lastLoadTime = time.Now()
		}
		return m, nil
	case TimeFilterChangedMsg:
		m.currentTimeRange = msg.TimeRange
		m.timeFilter = msg.TimeRange.Filter
		m.filterMenu.IsOpen = false
		m.filterMenu.SelectedFilter = msg.TimeRange
		m.isLoading = true
		m.loadingProgress = fmt.Sprintf("Loading data for %s...", msg.TimeRange.Desc)
		return m, m.loadData(msg.TimeRange.Filter)
	case FilterMenuToggleMsg:
		m.filterMenu.IsOpen = !m.filterMenu.IsOpen
		return m, nil
	case LoadProgressMsg:
		m.loadingProgress = msg.Message
		return m, nil
	case RefreshDataMsg:
		m.isLoading = true
		m.loadingProgress = fmt.Sprintf("Refreshing data for %s...", m.currentTimeRange.Desc)
		return m, m.loadData(m.timeFilter)
	}
	return m, nil
}

// handleKeypress processes keyboard input
func (m *DashboardModel) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys (work regardless of help state)
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}
		return m, tea.Quit
	}

	// Help toggle
	if msg.Type == tea.KeyRunes {
		switch msg.Runes[0] {
		case '?':
			m.showHelp = !m.showHelp
			return m, nil
		case 'q':
			return m, tea.Quit
		}
	}

	// If help is shown, block other navigation
	if m.showHelp {
		return m, nil
	}

	// If filter menu is open, handle filter menu navigation
	if m.filterMenu.IsOpen {
		return m.handleFilterMenuNavigation(msg)
	}

	// Handle navigation keys
	switch msg.Type {
	case tea.KeyTab:
		m.switchPanel()
		return m, nil
	case tea.KeyDown:
		if m.focusedPanel == TimelinePanel {
			m.timeline.ZoomTimeline(1) // Zoom out
		} else if m.focusedPanel == CoOccurrencePanel && m.coOccurrenceMatrixComponent.DisplayMode == MatrixView {
			m.coOccurrenceMatrixComponent.NavigateMatrix("down")
		} else {
			m.moveDown()
		}
		return m, nil
	case tea.KeyUp:
		if m.focusedPanel == TimelinePanel {
			m.timeline.ZoomTimeline(-1) // Zoom in
		} else if m.focusedPanel == CoOccurrencePanel && m.coOccurrenceMatrixComponent.DisplayMode == MatrixView {
			m.coOccurrenceMatrixComponent.NavigateMatrix("up")
		} else {
			m.moveUp()
		}
		return m, nil
	case tea.KeyLeft:
		if m.focusedPanel == TimelinePanel {
			m.timeline.NavigateTime(-1) // Move backward in time
			return m, nil
		} else if m.focusedPanel == CoOccurrencePanel && m.coOccurrenceMatrixComponent.DisplayMode == MatrixView {
			m.coOccurrenceMatrixComponent.NavigateMatrix("left")
			return m, nil
		}
	case tea.KeyRight:
		if m.focusedPanel == TimelinePanel {
			m.timeline.NavigateTime(1) // Move forward in time
			return m, nil
		} else if m.focusedPanel == CoOccurrencePanel && m.coOccurrenceMatrixComponent.DisplayMode == MatrixView {
			m.coOccurrenceMatrixComponent.NavigateMatrix("right")
			return m, nil
		}
	case tea.KeySpace:
		if m.focusedPanel == TimelinePanel {
			m.timeline.SelectNearestEvent()
			return m, nil
		}
	case tea.KeyEnter:
		// Trigger detail view or action
		return m, m.handleEnter()
	}

	// Handle rune keys (j, k, r, c, s, d, f, 1-6)
	if msg.Type == tea.KeyRunes {
		switch msg.Runes[0] {
		case 'j':
			m.moveDown()
			return m, nil
		case 'k':
			m.moveUp()
			return m, nil
		case 'r':
			// Trigger data refresh or reset timeline view
			if m.focusedPanel == TimelinePanel {
				m.timeline.ResetView()
				return m, nil
			}
			return m, func() tea.Msg {
				return RefreshDataMsg{}
			}
		case 'c':
			// Sort by count (stats panel)
			if m.focusedPanel == StatsPanel {
				m.setSortColumn(SortByCount)
				return m, nil
			}
		case 's':
			// Sort by success rate (stats panel) or cycle sort options (delegation panel)
			if m.focusedPanel == StatsPanel {
				m.setSortColumn(SortBySuccessRate)
				return m, nil
			} else if m.focusedPanel == DelegationPanel {
				m.cycleDelegationSortOptions()
				return m, nil
			}
		case 'd':
			// Sort by duration
			if m.focusedPanel == StatsPanel {
				m.setSortColumn(SortByDuration)
				return m, nil
			}
		case 'f':
			// Toggle filter menu for time range selection
			if m.filterMenu.IsOpen {
				m.filterMenu.IsOpen = false
			} else {
				return m, func() tea.Msg {
					return FilterMenuToggleMsg{}
				}
			}
			return m, nil
		case 'm':
			// Switch matrix display mode (co-occurrence panel only)
			if m.focusedPanel == CoOccurrencePanel {
				m.coOccurrenceMatrixComponent.SwitchDisplayMode()
				return m, nil
			}
		case 'i':
			// Switch to insights view (co-occurrence panel only)
			if m.focusedPanel == CoOccurrencePanel {
				m.coOccurrenceMatrixComponent.DisplayMode = InsightsView
				return m, nil
			}
		// Quick time filter shortcuts (1-6 keys)
		case '1':
			return m, m.handleQuickTimeFilter(0) // 1w
		case '2':
			return m, m.handleQuickTimeFilter(1) // 1m
		case '3':
			return m, m.handleQuickTimeFilter(2) // 3m
		case '4':
			return m, m.handleQuickTimeFilter(3) // 6m
		case '5':
			return m, m.handleQuickTimeFilter(4) // 1y
		case '6':
			return m, m.handleQuickTimeFilter(5) // all
		case 't':
			// Cycle through themes
			m.cycleTheme()
			return m, nil
		}
	}

	return m, nil
}

// switchPanel cycles through the available panels
func (m *DashboardModel) switchPanel() {
	// Simplified navigation for 2-column layout
	switch m.focusedPanel {
	case TimelinePanel:
		m.focusedPanel = CoOccurrencePanel // Go to left panel (matrix)
	case CoOccurrencePanel:
		m.focusedPanel = StatsPanel // Go to right panel (stats)
	case StatsPanel:
		m.focusedPanel = TimelinePanel // Back to timeline
	default:
		m.focusedPanel = TimelinePanel
	}
}

// moveDown moves the selection down in the current panel
func (m *DashboardModel) moveDown() {
	switch m.focusedPanel {
	case StatsPanel:
		maxIndex := len(m.agentStats) - 1
		if m.selectedStatsIndex < maxIndex {
			m.selectedStatsIndex++
		}
	case DelegationPanel:
		if m.selectedDelegationIndex < len(m.delegationPatterns)-1 {
			m.selectedDelegationIndex++
		}
	case CoOccurrencePanel:
		if m.selectedCoOccurrenceIndex < len(m.coOccurrences)-1 {
			m.selectedCoOccurrenceIndex++
		}
	}
}

// moveUp moves the selection up in the current panel
func (m *DashboardModel) moveUp() {
	switch m.focusedPanel {
	case StatsPanel:
		if m.selectedStatsIndex > 0 {
			m.selectedStatsIndex--
		}
	case DelegationPanel:
		if m.selectedDelegationIndex > 0 {
			m.selectedDelegationIndex--
		}
	case CoOccurrencePanel:
		if m.selectedCoOccurrenceIndex > 0 {
			m.selectedCoOccurrenceIndex--
		}
	}
}

// handleEnter processes the enter key press
func (m *DashboardModel) handleEnter() tea.Cmd {
	// Return a command to show detail view or perform action
	return func() tea.Msg {
		return RefreshDataMsg{} // Placeholder for detail view
	}
}

// setSortColumn sets the sort column and toggles direction if it's the same column
func (m *DashboardModel) setSortColumn(column SortColumn) {
	if m.sortColumn == column {
		// Toggle direction if same column
		if m.sortDirection == SortDescending {
			m.sortDirection = SortAscending
		} else {
			m.sortDirection = SortDescending
		}
	} else {
		// New column, default to descending
		m.sortColumn = column
		m.sortDirection = SortDescending
	}

	// Reset selection to top when sorting changes
	m.selectedStatsIndex = 0
}

// handleDelegationFilterToggle toggles delegation filter options
func (m *DashboardModel) handleDelegationFilterToggle() tea.Cmd {
	// Cycle through min count filters: 0 → 2 → 5 → 10 → 0
	switch m.delegationFilterOptions.MinCount {
	case 0:
		m.delegationFilterOptions.MinCount = 2
	case 2:
		m.delegationFilterOptions.MinCount = 5
	case 5:
		m.delegationFilterOptions.MinCount = 10
	default:
		m.delegationFilterOptions.MinCount = 0
	}
	// Reset selection when filters change
	m.selectedDelegationIndex = 0
	return nil
}

// cycleDelegationSortOptions cycles through available sort options for delegation patterns
func (m *DashboardModel) cycleDelegationSortOptions() {
	switch m.delegationFilterOptions.SortBy {
	case SortByFrequency:
		m.delegationFilterOptions.SortBy = SortByRecency
	case SortByRecency:
		m.delegationFilterOptions.SortBy = SortBySource
	case SortBySource:
		m.delegationFilterOptions.SortBy = SortByTarget
	case SortByTarget:
		m.delegationFilterOptions.SortBy = SortByFrequency
	}
	// Reset selection when sorting changes
	m.selectedDelegationIndex = 0
}

// handleQuickTimeFilter handles quick time filter shortcuts (1-6 keys)
func (m *DashboardModel) handleQuickTimeFilter(index int) tea.Cmd {
	if index >= 0 && index < len(TimeRanges) {
		return func() tea.Msg {
			return TimeFilterChangedMsg{
				TimeRange: TimeRanges[index],
			}
		}
	}
	return nil
}

// handleFilterMenuNavigation handles navigation within the filter menu
func (m *DashboardModel) handleFilterMenuNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.filterMenu.IsOpen = false
		return m, nil
	case tea.KeyEnter:
		// Apply selected filter
		selectedRange := TimeRanges[m.filterMenu.SelectedIndex]
		m.filterMenu.SelectedFilter = selectedRange
		m.filterMenu.IsOpen = false
		return m, func() tea.Msg {
			return TimeFilterChangedMsg{
				TimeRange: selectedRange,
			}
		}
	case tea.KeyUp:
		if m.filterMenu.SelectedIndex > 0 {
			m.filterMenu.SelectedIndex--
		}
		return m, nil
	case tea.KeyDown:
		if m.filterMenu.SelectedIndex < len(TimeRanges)-1 {
			m.filterMenu.SelectedIndex++
		}
		return m, nil
	}

	// Handle number keys in filter menu
	if msg.Type == tea.KeyRunes {
		switch msg.Runes[0] {
		case '1', '2', '3', '4', '5', '6':
			index := int(msg.Runes[0] - '1')
			if index >= 0 && index < len(TimeRanges) {
				m.filterMenu.SelectedIndex = index
				selectedRange := TimeRanges[index]
				m.filterMenu.SelectedFilter = selectedRange
				m.filterMenu.IsOpen = false
				return m, func() tea.Msg {
					return TimeFilterChangedMsg{
						TimeRange: selectedRange,
					}
				}
			}
		}
	}

	return m, nil
}

// getSortIndicator returns the sort indicator symbol for display
func (m *DashboardModel) getSortIndicator(column SortColumn) string {
	if m.sortColumn != column {
		return ""
	}

	if m.sortDirection == SortDescending {
		return " ↓"
	}
	return " ↑"
}

// getSortColumnName returns the human-readable name of the sort column
func (m *DashboardModel) getSortColumnName() string {
	switch m.sortColumn {
	case SortByCount:
		return "Count"
	case SortBySuccessRate:
		return "Success Rate"
	case SortByDuration:
		return "Duration"
	default:
		return "Count"
	}
}

// renderPanelWithTitle creates a btop-style bordered panel with title in the top border
func (m *DashboardModel) renderPanelWithTitle(content string, title string, width int, height int, borderColor lipgloss.Color, isFocused bool) string {
	// Create the title with proper styling
	titleColor := m.currentTheme.Primary
	if isFocused {
		titleColor = m.currentTheme.Highlight
	}

	// Build the top border with embedded title (btop style)
	titleLen := len(title)
	borderChar := "─"
	if isFocused {
		borderChar = "═"
	}

	topBorder := "┌" + borderChar + title
	remainingWidth := width - titleLen - 4 // account for corners and title
	if remainingWidth > 0 {
		topBorder += strings.Repeat(borderChar, remainingWidth)
	}
	topBorder += "┐"

	// Style the top border with embedded title
	topBorderStyled := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(topBorder[:1]) + // left corner
		lipgloss.NewStyle().
			Foreground(titleColor).
			Bold(true).
			Render(title) +
		lipgloss.NewStyle().
			Foreground(borderColor).
			Render(topBorder[1+titleLen:]) // rest of border

	// Create the bordered content
	lines := strings.Split(content, "\n")
	var result strings.Builder
	result.WriteString(topBorderStyled + "\n")

	// Add content lines with side borders
	sideChar := "│"
	if isFocused {
		sideChar = "║"
	}

	borderStyle := lipgloss.NewStyle().Foreground(borderColor)
	contentWidth := width - 2 // account for side borders

	for i := 0; i < height-2 && i < len(lines); i++ {
		line := lines[i]
		if len(line) > contentWidth {
			line = line[:contentWidth]
		} else if len(line) < contentWidth {
			line = line + strings.Repeat(" ", contentWidth-len(line))
		}
		result.WriteString(borderStyle.Render(sideChar) + line + borderStyle.Render(sideChar) + "\n")
	}

	// Fill remaining height if needed
	for i := len(lines); i < height-2; i++ {
		result.WriteString(borderStyle.Render(sideChar) + strings.Repeat(" ", contentWidth) + borderStyle.Render(sideChar) + "\n")
	}

	// Add bottom border
	bottomChar := "└"
	if isFocused {
		bottomChar = "╚"
	}
	bottomBorder := bottomChar + strings.Repeat(borderChar, width-2) + "┘"
	if isFocused {
		bottomBorder = bottomChar + strings.Repeat(borderChar, width-2) + "╝"
	}
	result.WriteString(borderStyle.Render(bottomBorder))

	return result.String()
}

// cycleTheme cycles through available themes
func (m *DashboardModel) cycleTheme() {
	themeOrder := []string{"dracula", "nord", "monokai", "github", "solarized-dark", "one-dark"}

	// Find current theme index
	currentIndex := 0
	for i, name := range themeOrder {
		if name == m.currentThemeName {
			currentIndex = i
			break
		}
	}

	// Move to next theme
	nextIndex := (currentIndex + 1) % len(themeOrder)
	m.currentThemeName = themeOrder[nextIndex]
	m.currentTheme = GetTheme(m.currentThemeName)
	m.styles = ApplyTheme(m.currentTheme)

	// Update component themes
	m.timeline.ApplyTheme(m.currentTheme)
	m.coOccurrenceMatrixComponent.ApplyTheme(m.currentTheme)
}

// View renders the dashboard view
func (m *DashboardModel) View() string {
	if m.showHelp {
		return m.renderHelp()
	}

	if m.filterMenu.IsOpen {
		return m.renderFilterMenu()
	}

	var content strings.Builder

	// Enhanced title with status indicators using theme
	titleStyle := m.styles.Title.
		Align(lipgloss.Center).
		Width(m.width)

	// Add loading indicator if data is loading
	var statusIndicator string
	if m.isLoading {
		statusIndicator = m.styles.Warning.
			Render(" ⟳ Loading...")
	} else if m.lastLoadError != nil {
		statusIndicator = m.styles.Error.
			Render(" ⚠ Error")
	} else {
		// Show data freshness
		if !m.lastLoadTime.IsZero() {
			timeSince := time.Since(m.lastLoadTime)
			if timeSince < time.Minute {
				statusIndicator = m.styles.Success.
					Render(" ✓ Fresh")
			} else if timeSince < 5*time.Minute {
				statusIndicator = m.styles.Warning.
					Render(" ○ Recent")
			} else {
				statusIndicator = m.styles.Muted.
					Render(" ● Stale")
			}
		}
	}

	// Add filter indicator
	filterInfo := fmt.Sprintf(" [%s]", m.currentTimeRange.Desc)
	filterIndicator := m.styles.Info.
		Render(filterInfo)

	// Add theme name to title
	themeIndicator := m.styles.Muted.Render(fmt.Sprintf(" [%s]", m.currentTheme.Name))
	title := titleStyle.Render("🔧 Agent Dashboard" + statusIndicator + filterIndicator + themeIndicator)
	content.WriteString(title + "\n\n")

	// Calculate heights for new layout
	// Title takes up ~2 lines
	availableHeight := m.height - 4 // Account for title and padding

	// New layout: Timeline (30%), Main panels (55%), Footer (15%)
	timelineHeight := int(float64(availableHeight) * 0.30)
	mainPanelHeight := int(float64(availableHeight) * 0.55)
	footerHeight := int(float64(availableHeight) * 0.15)

	// Ensure minimum heights
	if timelineHeight < 6 {
		timelineHeight = 6 // Minimum for graph visibility
	}
	if mainPanelHeight < 15 {
		mainPanelHeight = 15
	}
	if footerHeight < 3 {
		footerHeight = 3
	}

	// 1. Render timeline panel at the top with project name
	// Get current project path for display
	projectName := "the-startup" // Default
	if cwd, err := os.Getwd(); err == nil {
		projectName = filepath.Base(cwd)
	}

	timelineContent := m.renderTimelinePanelContent(timelineHeight)
	borderColor := m.currentTheme.Border
	if m.focusedPanel == TimelinePanel {
		borderColor = m.currentTheme.Highlight
	}

	// Create btop-style border with title in the top border
	timelineStyled := m.renderPanelWithTitle(
		timelineContent,
		fmt.Sprintf(" timeline [%s] ", projectName),
		m.width-2,
		timelineHeight,
		borderColor,
		m.focusedPanel == TimelinePanel,
	)
	content.WriteString(timelineStyled + "\n")

	// 2. Two-column layout: Co-occurrence Matrix (left 50%) and Agent Stats (right 50%)
	leftWidth := (m.width - 4) / 2
	rightWidth := m.width - 4 - leftWidth

	// Co-occurrence Matrix Panel (left) - with full agent names
	coOccurrencePanel := m.renderCoOccurrencePanelContent()
	matrixBorderColor := m.currentTheme.Border
	if m.focusedPanel == CoOccurrencePanel {
		matrixBorderColor = m.currentTheme.Highlight
	}

	// Create btop-style border with title
	coOccurrenceStyled := m.renderPanelWithTitle(
		coOccurrencePanel,
		" matrix ",
		leftWidth,
		mainPanelHeight,
		matrixBorderColor,
		m.focusedPanel == CoOccurrencePanel,
	)

	// Agent Stats Panel (right)
	statsPanel := m.renderStatsPanelContent()
	statsBorderColor := m.currentTheme.Border
	if m.focusedPanel == StatsPanel {
		statsBorderColor = m.currentTheme.Highlight
	}

	// Create btop-style border with title
	statsStyled := m.renderPanelWithTitle(
		statsPanel,
		fmt.Sprintf(" stats [%d] ", len(m.agentStats)),
		rightWidth,
		mainPanelHeight,
		statsBorderColor,
		m.focusedPanel == StatsPanel,
	)

	// Join the two panels horizontally
	mainPanels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		coOccurrenceStyled,
		statsStyled,
	)
	content.WriteString(mainPanels + "\n")

	// 3. Render footer/filters panel at bottom (15% height)
	footerStyle := m.styles.BorderNormal.
		Padding(0, 1).
		Width(m.width - 2).
		Height(footerHeight)

	// Context-sensitive footer based on focused panel
	var footerText string
	switch m.focusedPanel {
	case TimelinePanel:
		footerText = "Timeline: ←→ time nav • ↑↓ zoom • Space select • r reset │ Global: t theme • 1-6 time filter • f filter menu • Tab panels • ? help • q quit"
	case StatsPanel:
		footerText = "Stats: c/s/d sort • ↑↓ navigate │ Global: 1-6 time filter • f filter menu • Tab panels • ? help • q quit"
	case DelegationPanel:
		footerText = "Delegation: s sort • ↑↓ navigate │ Global: 1-6 time filter • f filter menu • Tab panels • ? help • q quit"
	case CoOccurrencePanel:
		footerText = "Collaboration: m mode • i insights • ↑↓←→ navigate │ Global: 1-6 time filter • f filter menu • Tab panels • ? help • q quit"
	}

	footer := footerStyle.Render(footerText)
	content.WriteString(footer)

	return content.String()
}

// renderTimelinePanel renders the message timeline panel with enhanced styling (DEPRECATED - use renderTimelinePanelContent)
func (m *DashboardModel) renderTimelinePanel(maxHeight int) string {
	// Set the timeline component height
	m.timeline.MaxHeight = maxHeight
	m.timeline.MaxWidth = m.width - 8

	// Get current project path for display
	projectName := "the-startup" // Default
	if cwd, err := os.Getwd(); err == nil {
		projectName = filepath.Base(cwd)
	}

	// Timeline title with project name
	titleText := fmt.Sprintf("📊 Message Timeline - Project: %s", projectName)
	if m.focusedPanel == TimelinePanel {
		titleText += " (Active)"
	}

	// Render the timeline
	timelineContent := m.timeline.RenderTimeline(
		m.focusedPanel == TimelinePanel,
		titleText,
	)

	return timelineContent
}

// renderTimelinePanelContent renders just the timeline content without title (for btop-style borders)
func (m *DashboardModel) renderTimelinePanelContent(maxHeight int) string {
	// Set the timeline component height
	m.timeline.MaxHeight = maxHeight - 2 // Account for border
	m.timeline.MaxWidth = m.width - 8

	// Render the timeline without title (title is in border now)
	timelineContent := m.timeline.RenderTimelineContent(
		m.focusedPanel == TimelinePanel,
	)

	return timelineContent
}

// renderStatsPanelContent renders just the stats content without title (for btop-style borders)
func (m *DashboardModel) renderStatsPanelContent() string {
	var content strings.Builder

	// Add sort status indicator
	sortStatusStyle := m.styles.Muted.Faint(true)
	sortStatus := fmt.Sprintf("Sort: %s %s",
		m.getSortColumnName(),
		map[SortDirection]string{SortAscending: "↑", SortDescending: "↓"}[m.sortDirection])
	content.WriteString(sortStatusStyle.Render(sortStatus) + "\n\n")

	// Create sort indicators map
	sortIndicators := make(map[SortColumn]string)
	sortIndicators[SortByCount] = m.getSortIndicator(SortByCount)
	sortIndicators[SortBySuccessRate] = m.getSortIndicator(SortBySuccessRate)
	sortIndicators[SortByDuration] = m.getSortIndicator(SortByDuration)

	// Render leaderboard
	if len(m.agentStats) == 0 {
		content.WriteString(m.styles.EmptyState.Render("No agent data available\nRun some commands to see statistics"))
	} else {
		leaderboardContent := m.leaderboard.RenderAgentLeaderboardWithSort(
			m.agentStats,
			m.selectedStatsIndex,
			m.focusedPanel == StatsPanel,
			m.sortColumn,
			m.sortDirection,
			sortIndicators,
		)
		content.WriteString(leaderboardContent)
	}

	return content.String()
}

// renderStatsPanel renders the agent statistics panel with professional styling (DEPRECATED - use renderStatsPanelContent)
func (m *DashboardModel) renderStatsPanel() string {
	var content strings.Builder

	// Enhanced panel title with focus indicator and stats summary
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	if m.focusedPanel == StatsPanel {
		titleStyle = titleStyle.Foreground(lipgloss.Color("46"))
	}

	totalAgents := len(m.agentStats)
	titleText := fmt.Sprintf("📈 Agent Statistics (%d agents)", totalAgents)
	if m.focusedPanel == StatsPanel {
		titleText += " (Active)"
	}

	content.WriteString(titleStyle.Render(titleText) + "\n")

	// Add sort status indicator
	sortStatusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	sortStatus := fmt.Sprintf("Sorted by %s %s",
		m.getSortColumnName(),
		map[SortDirection]string{SortAscending: "(ascending)", SortDescending: "(descending)"}[m.sortDirection])
	content.WriteString(sortStatusStyle.Render(sortStatus) + "\n\n")

	// Create sort indicators map
	sortIndicators := make(map[SortColumn]string)
	sortIndicators[SortByCount] = m.getSortIndicator(SortByCount)
	sortIndicators[SortBySuccessRate] = m.getSortIndicator(SortBySuccessRate)
	sortIndicators[SortByDuration] = m.getSortIndicator(SortByDuration)

	// Render leaderboard with enhanced styling
	if len(m.agentStats) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(2, 4).
			Margin(1, 0)

		content.WriteString(emptyStyle.Render("No agent data available\nRun some commands to see statistics"))
	} else {
		leaderboardContent := m.leaderboard.RenderAgentLeaderboardWithSort(
			m.agentStats,
			m.selectedStatsIndex,
			m.focusedPanel == StatsPanel,
			m.sortColumn,
			m.sortDirection,
			sortIndicators,
		)
		content.WriteString(leaderboardContent)
	}

	return content.String()
}

// renderDelegationPanel renders the delegation patterns panel
func (m *DashboardModel) renderDelegationPanel() string {
	// Enhanced delegation patterns section
	delegationTitle := "🔄 Delegation Flow Patterns"
	if m.focusedPanel == DelegationPanel {
		delegationTitle += " (Active)"
	}

	delegationContent := m.delegationFlowComponent.RenderDelegationFlows(
		m.delegationPatterns,
		m.delegationFilterOptions,
		m.selectedDelegationIndex,
		m.focusedPanel == DelegationPanel,
		delegationTitle,
	)

	return delegationContent
}

// renderCoOccurrencePanelContent renders just the co-occurrence content without title (for btop-style borders)
func (m *DashboardModel) renderCoOccurrencePanelContent() string {
	// Render the matrix without title (title is in border now)
	coOccurrenceContent := m.coOccurrenceMatrixComponent.RenderCoOccurrenceMatrixContent(
		m.coOccurrences,
		m.selectedCoOccurrenceIndex,
		m.focusedPanel == CoOccurrencePanel,
	)

	return coOccurrenceContent
}

// renderCoOccurrencePanel renders the co-occurrence patterns panel (DEPRECATED - use renderCoOccurrencePanelContent)
func (m *DashboardModel) renderCoOccurrencePanel() string {
	// Enhanced co-occurrence patterns section
	coOccurrenceTitle := "🤝 Agent Collaboration Matrix"
	if m.focusedPanel == CoOccurrencePanel {
		coOccurrenceTitle += " (Active)"
	}

	coOccurrenceContent := m.coOccurrenceMatrixComponent.RenderCoOccurrenceMatrix(
		m.coOccurrences,
		m.selectedCoOccurrenceIndex,
		m.focusedPanel == CoOccurrencePanel,
		coOccurrenceTitle,
	)

	return coOccurrenceContent
}

// renderDelegationAndCoOccurrence renders the right panel with professional styling
func (m *DashboardModel) renderDelegationAndCoOccurrence() string {
	var content strings.Builder

	// Calculate available height for each sub-panel
	availableHeight := m.height - 15 // Account for title, timeline, and footer
	delegationHeight := availableHeight / 2
	coOccurrenceHeight := availableHeight - delegationHeight

	// Enhanced delegation patterns section
	delegationTitle := "🔄 Delegation Flow Patterns"
	if m.focusedPanel == DelegationPanel {
		delegationTitle += " (Active)"
	}

	delegationContent := m.delegationFlowComponent.RenderDelegationFlows(
		m.delegationPatterns,
		m.delegationFilterOptions,
		m.selectedDelegationIndex,
		m.focusedPanel == DelegationPanel,
		delegationTitle,
	)

	// Style delegation panel with consistent border
	delegationStyled := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("99")).
		Padding(0, 0, 1, 0).
		Height(delegationHeight).
		Render(delegationContent)

	content.WriteString(delegationStyled + "\n")

	// Enhanced co-occurrence patterns section
	coOccurrenceTitle := "🤝 Agent Collaboration Matrix"
	if m.focusedPanel == CoOccurrencePanel {
		coOccurrenceTitle += " (Active)"
	}

	coOccurrenceContent := m.coOccurrenceMatrixComponent.RenderCoOccurrenceMatrix(
		m.coOccurrences,
		m.selectedCoOccurrenceIndex,
		m.focusedPanel == CoOccurrencePanel,
		coOccurrenceTitle,
	)

	// Style co-occurrence panel
	coOccurrenceStyled := lipgloss.NewStyle().
		Height(coOccurrenceHeight).
		Render(coOccurrenceContent)

	content.WriteString(coOccurrenceStyled)

	return content.String()
}

// renderHelp renders the comprehensive help screen with context-sensitive help
func (m *DashboardModel) renderHelp() string {
	var help strings.Builder

	// Title with enhanced styling
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		Width(80)

	help.WriteString(titleStyle.Render("🔧 Agent Dashboard - Complete Help Guide") + "\n\n")

	// Context-sensitive help section based on focused panel
	contextHelp := m.renderContextSensitiveHelp()
	if contextHelp != "" {
		help.WriteString(contextHelp + "\n\n")
	}

	// Core navigation section
	help.WriteString(m.renderHelpSection("🎯 Core Navigation", []HelpItem{
		{"Tab", "Switch between panels (Timeline → Stats → Delegation → Co-occurrence)"},
		{"↑/↓ or j/k", "Navigate within current panel (move selection up/down)"},
		{"Enter", "View details for selected item or trigger action"},
		{"?", "Toggle this comprehensive help screen"},
		{"q", "Quit dashboard"},
		{"Esc", "Quit dashboard (or close help if open)"},
		{"Ctrl+C", "Force quit application"},
	}))

	// Time filter help
	help.WriteString(m.renderHelpSection("🕰 Time Filter Commands", []HelpItem{
		{"1", "Filter to last 7 days (1w)"},
		{"2", "Filter to last 30 days (1m) - Default"},
		{"3", "Filter to last 3 months (3m)"},
		{"4", "Filter to last 6 months (6m)"},
		{"5", "Filter to last year (1y)"},
		{"6", "Show all data (all)"},
		{"f", "Open interactive filter menu"},
		{"r", "Refresh data with current filter"},
	}))

	// Timeline panel help
	help.WriteString(m.renderHelpSection("📊 Timeline Panel Commands", []HelpItem{
		{"←/→", "Navigate backward/forward in time"},
		{"↑/↓", "Zoom in/out timeline (Hour ↔ Day ↔ Week ↔ Month)"},
		{"Space", "Select nearest event to viewport center"},
		{"r", "Reset timeline view to show recent events"},
	}))

	// Stats panel help
	help.WriteString(m.renderHelpSection("📈 Agent Statistics Panel", []HelpItem{
		{"c", "Sort by Count (toggle ascending/descending)"},
		{"s", "Sort by Success Rate (toggle ascending/descending)"},
		{"d", "Sort by Duration (toggle ascending/descending)"},
		{"↑/↓", "Navigate through agent list"},
	}))

	// Delegation panel help
	help.WriteString(m.renderHelpSection("🔄 Delegation Flow Panel", []HelpItem{
		{"f", "Filter patterns (cycles: none → min 2 → min 5 → min 10)"},
		{"s", "Sort patterns (frequency → recency → source → target)"},
		{"↑/↓", "Navigate through delegation patterns"},
	}))

	// Co-occurrence panel help
	help.WriteString(m.renderHelpSection("🤝 Agent Collaboration Panel", []HelpItem{
		{"m", "Switch display mode (matrix → relationships → insights)"},
		{"i", "Jump directly to insights view"},
		{"←→↑↓", "Navigate matrix cells (matrix view only)"},
		{"↑/↓", "Navigate relationships list (relationship view)"},
	}))

	// Data commands
	help.WriteString(m.renderHelpSection("💾 Data Operations", []HelpItem{
		{"r", "Refresh data (reload from logs)"},
	}))

	// Panel descriptions
	help.WriteString(m.renderPanelDescriptions())

	// Status and visual indicators
	help.WriteString(m.renderVisualIndicatorsHelp())

	// Footer
	footerStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center).
		Width(80)

	help.WriteString("\n" + footerStyle.Render("Press any key to continue or ? to close help") + "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(84).
		Render(help.String())
}

// HelpItem represents a single help entry
type HelpItem struct {
	Key         string
	Description string
}

// renderContextSensitiveHelp provides help specific to the currently focused panel
func (m *DashboardModel) renderContextSensitiveHelp() string {
	if !m.showHelp {
		return ""
	}

	var panelName, panelTips string
	switch m.focusedPanel {
	case TimelinePanel:
		panelName = "📊 Timeline Panel Active"
		panelTips = "Currently viewing message timeline. Use ←→ to navigate time, ↑↓ to zoom, Space to select events."
	case StatsPanel:
		panelName = "📈 Agent Statistics Panel Active"
		panelTips = "Currently viewing agent leaderboard. Use c/s/d to sort by different metrics, ↑↓ to navigate."
	case DelegationPanel:
		panelName = "🔄 Delegation Flow Panel Active"
		panelTips = "Currently viewing agent delegation patterns. Use f to filter, s to sort, ↑↓ to navigate flows."
	case CoOccurrencePanel:
		panelName = "🤝 Collaboration Panel Active"
		panelTips = "Currently viewing agent collaboration matrix. Use m to switch modes, ↑↓←→ to navigate."
	}

	contextStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46")).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("46")).
		Padding(0, 1).
		Width(78)

	return contextStyle.Render(panelName+"\n"+panelTips)
}

// renderHelpSection renders a help section with consistent formatting
func (m *DashboardModel) renderHelpSection(title string, items []HelpItem) string {
	var section strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		MarginTop(1)

	section.WriteString(titleStyle.Render(title) + "\n")

	for _, item := range items {
		keyStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Width(12)

		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

		line := fmt.Sprintf("  %s %s",
			keyStyle.Render(item.Key),
			descStyle.Render(item.Description))
		section.WriteString(line + "\n")
	}

	return section.String()
}

// renderPanelDescriptions provides detailed panel descriptions
func (m *DashboardModel) renderPanelDescriptions() string {
	var desc strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		MarginTop(1)

	desc.WriteString(titleStyle.Render("📋 Panel Descriptions") + "\n")

	panels := []struct {
		Name        string
		Description string
		Features    []string
	}{
		{
			Name:        "📊 Message Timeline",
			Description: "Horizontal timeline showing message exchanges over time",
			Features: []string{
				"Color-coded by role: Blue (User), Pink (Assistant), Yellow (Tool), Gray (System)",
				"Success indicators: ● (success), ○ (failure), ◆ (selected)",
				"Interactive time navigation and zoom controls",
				"Multi-track display for different message types",
			},
		},
		{
			Name:        "📈 Agent Statistics",
			Description: "Leaderboard of agent usage, success rates, and performance metrics",
			Features: []string{
				"Success rates: Green (≥90%), Yellow (≥70%), Red (<70%)",
				"Duration indicators: Green (<1s), Red (>5s)",
				"Sortable by count, success rate, and duration",
				"Real-time performance tracking",
			},
		},
		{
			Name:        "🔄 Delegation Flow Patterns",
			Description: "Enhanced A → B flows with frequency visualization",
			Features: []string{
				"ASCII flow diagrams with frequency bars",
				"Pattern filtering and sorting capabilities",
				"Temporal flow analysis",
				"Agent transition tracking",
			},
		},
		{
			Name:        "🤝 Agent Collaboration Matrix",
			Description: "Enhanced relationship visualization with multiple modes",
			Features: []string{
				"Matrix view: Visual correlation matrix with strength indicators",
				"Relationship view: Detailed correlation analysis",
				"Insights view: Statistical analysis of collaboration patterns",
				"Strength indicators: Strong/Medium/Weak/Minimal",
				"Pattern classification: Sequential/Concurrent/Alternative",
			},
		},
	}

	for _, panel := range panels {
		panelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Italic(true)

		featureStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

		desc.WriteString("  " + panelStyle.Render(panel.Name) + "\n")
		desc.WriteString("    " + descStyle.Render(panel.Description) + "\n")
		for _, feature := range panel.Features {
			desc.WriteString("    • " + featureStyle.Render(feature) + "\n")
		}
		desc.WriteString("\n")
	}

	return desc.String()
}

// renderVisualIndicatorsHelp explains the visual indicators used throughout the dashboard
func (m *DashboardModel) renderVisualIndicatorsHelp() string {
	var indicators strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		MarginTop(1)

	indicators.WriteString(titleStyle.Render("🎨 Visual Indicators") + "\n")

	visualItems := []struct {
		Category string
		Items    []HelpItem
	}{
		{
			Category: "Timeline Symbols",
			Items: []HelpItem{
				{"●", "Successful message/event"},
				{"○", "Failed message/event"},
				{"◆", "Currently selected event"},
				{"─", "Time axis line"},
				{"┬", "Time marker"},
			},
		},
		{
			Category: "Matrix Relationship Strength",
			Items: []HelpItem{
				{"█", "Strong relationship (10+ sessions)"},
				{"▓", "Medium relationship (5-9 sessions)"},
				{"▒", "Weak relationship (2-4 sessions)"},
				{"░", "Minimal relationship (1 session)"},
				{"·", "No relationship"},
			},
		},
		{
			Category: "Status Colors",
			Items: []HelpItem{
				{"Green", "High success rates (≥90%) or fast duration (<1s)"},
				{"Yellow", "Medium success rates (70-89%)"},
				{"Red", "Low success rates (<70%) or slow duration (>5s)"},
				{"Pink", "Selected items or focused panels"},
				{"Blue", "User messages or primary actions"},
			},
		},
		{
			Category: "Navigation Indicators",
			Items: []HelpItem{
				{"▶", "Currently selected item"},
				{"↑↓", "Sort direction indicators"},
				{"→", "Delegation flow direction"},
				{"⟷", "Concurrent collaboration pattern"},
				{"↔", "Alternative collaboration pattern"},
			},
		},
	}

	for _, category := range visualItems {
		categoryStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("220"))

		indicators.WriteString("  " + categoryStyle.Render(category.Category) + "\n")
		for _, item := range category.Items {
			symbolStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				Width(8)

			descStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

			line := fmt.Sprintf("    %s %s",
				symbolStyle.Render(item.Key),
				descStyle.Render(item.Description))
			indicators.WriteString(line + "\n")
		}
		indicators.WriteString("\n")
	}

	return indicators.String()
}

// loadData creates a command to load data using existing stats engine
func (m *DashboardModel) loadData(timeFilter string) tea.Cmd {
	return func() tea.Msg {
		// Parse time filter using the same logic as cmd/stats.go
		now := time.Now()
		var startTime *time.Time

		if timeFilter != "" {
			t, err := parseTimeFilter(timeFilter, now)
			if err != nil {
				return DataLoadedMsg{Error: fmt.Errorf("invalid time filter: %w", err)}
			}
			startTime = &t
		}

		// Discover log files (reusing logic from cmd/stats.go)
		discovery := stats.NewLogDiscovery()
		projectPath := discovery.GetCurrentProject()
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return DataLoadedMsg{Error: fmt.Errorf("failed to get current directory: %w", err)}
			}
			projectPath = cwd
		}

		filter := stats.FilterOptions{StartTime: startTime}
		logFiles, err := discovery.FindLogFiles(projectPath, filter)
		if err != nil {
			return DataLoadedMsg{Error: fmt.Errorf("failed to discover log files: %w", err)}
		}

		if len(logFiles) == 0 {
			return DataLoadedMsg{Error: fmt.Errorf("no log files found in current project")}
		}

		// Parse log files using existing stats engine
		parser := stats.NewJSONLParser()
		var entries []stats.ClaudeLogEntry

		for _, logFile := range logFiles {
			entryChan, errorChan := parser.ParseFile(logFile)

			for entry := range entryChan {
				entries = append(entries, entry)
			}

			// Drain error channel (ignore errors for now)
			for range errorChan {
			}
		}

		if len(entries) == 0 {
			return DataLoadedMsg{Error: fmt.Errorf("no log entries found matching the criteria")}
		}

		// Extract agent analytics (reusing logic from cmd/stats.go)
		agentStats, delegationPatterns, coOccurrences, err := extractAgentAnalytics(entries)
		if err != nil {
			return DataLoadedMsg{Error: fmt.Errorf("failed to extract agent analytics: %w", err)}
		}

		// Extract timeline events
		messageEvents := extractMessageEvents(entries)


		return DataLoadedMsg{
			AgentStats:         agentStats,
			DelegationPatterns: delegationPatterns,
			CoOccurrences:      coOccurrences,
			MessageEvents:      messageEvents,
		}
	}
}

// parseTimeFilter converts a time string to a time.Time (copied from cmd/stats.go)
func parseTimeFilter(since string, now time.Time) (time.Time, error) {
	// Handle special keywords
	switch strings.ToLower(since) {
	case "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		return time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location()), nil
	}

	// Handle duration formats (1h, 24h, 7d, 30d)
	if strings.HasSuffix(since, "h") {
		hours := strings.TrimSuffix(since, "h")
		if duration, err := time.ParseDuration(hours + "h"); err == nil {
			return now.Add(-duration), nil
		}
	}

	if strings.HasSuffix(since, "d") {
		days := strings.TrimSuffix(since, "d")
		var d int
		if _, err := fmt.Sscanf(days, "%d", &d); err == nil {
			return now.AddDate(0, 0, -d), nil
		}
	}

	// Try parsing as date (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", since); err == nil {
		return t, nil
	}

	// Try parsing as datetime (YYYY-MM-DD HH:MM:SS)
	if t, err := time.Parse("2006-01-02 15:04:05", since); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unrecognized time format: %s (use 1h, 24h, 7d, 30d, or YYYY-MM-DD)", since)
}

// extractAgentAnalytics processes log entries to extract agent statistics and patterns (copied from cmd/stats.go)
func extractAgentAnalytics(entries []stats.ClaudeLogEntry) (
	map[string]*stats.GlobalAgentStats,
	[]stats.DelegationPattern,
	[]stats.AgentCoOccurrence,
	error,
) {
	// Create analyzers
	agentAnalyzer := stats.NewAgentAnalyzer()
	delegationAnalyzer := stats.NewDelegationAnalyzer()

	// Process entries to detect agents and build patterns
	sessionAgents := make(map[string][]string)

	for _, entry := range entries {
		// Detect agent in the entry using existing function
		agentType, confidence := stats.DetectAgent(entry)

		if agentType != "" && confidence > 0.5 { // Only process high confidence detections
			// Create agent invocation
			inv := stats.AgentInvocation{
				AgentType:       agentType,
				SessionID:       entry.SessionID,
				Timestamp:       entry.Timestamp,
				DurationMs:      1000, // Default duration, could be extracted from logs
				Success:         true, // Simplified, could be improved with error detection
				Confidence:      confidence,
				DetectionMethod: "pattern_detection",
			}

			// Process with agent analyzer
			agentAnalyzer.ProcessAgentInvocation(inv)

			// Track session agents for delegation analysis
			if sessionAgents[entry.SessionID] == nil {
				sessionAgents[entry.SessionID] = []string{}
			}
			sessionAgents[entry.SessionID] = append(sessionAgents[entry.SessionID], agentType)
		}
	}

	// Build delegation patterns from session agents
	for sessionID, agents := range sessionAgents {
		// Process agent transitions
		for i := 0; i < len(agents)-1; i++ {
			delegationAnalyzer.ProcessAgentTransition(
				sessionID,
				agents[i],
				agents[i+1],
				time.Now(), // Simplified timestamp
			)
		}

		// Process co-occurrence patterns manually by creating transitions for all pairs
		for i := 0; i < len(agents); i++ {
			for j := i + 1; j < len(agents); j++ {
				// Create bidirectional transitions for co-occurrence
				delegationAnalyzer.ProcessAgentTransition(sessionID, agents[i], agents[j], time.Now())
				delegationAnalyzer.ProcessAgentTransition(sessionID, agents[j], agents[i], time.Now())
			}
		}
	}

	// Get results
	agentStats, err := agentAnalyzer.GetAllAgentStats()
	if err != nil {
		return nil, nil, nil, err
	}

	delegationPatterns, err := delegationAnalyzer.GetTransitionPatterns()
	if err != nil {
		return nil, nil, nil, err
	}

	coOccurrences, err := delegationAnalyzer.GetCoOccurrencePatterns()
	if err != nil {
		return nil, nil, nil, err
	}

	return agentStats, delegationPatterns, coOccurrences, nil
}

// extractMessageEvents converts log entries to MessageEvent structures for timeline visualization
func extractMessageEvents(entries []stats.ClaudeLogEntry) []MessageEvent {
	var events []MessageEvent


	for _, entry := range entries {
		// Process all entries with timestamps
		if entry.Timestamp.IsZero() {
			continue
		}

		// Determine role - be more flexible
		var role MessageRole
		switch strings.ToLower(entry.Type) {
		case "user":
			role = UserMessage
		case "assistant":
			role = AssistantMessage
		case "summary": // Include summaries as system messages
			role = SystemMessage
		case "tool", "tool_use":
			role = ToolMessage
		default:
			// Try to infer from content
			content := string(entry.Content)
			if strings.Contains(content, "Human:") || strings.Contains(content, "user:") {
				role = UserMessage
			} else if strings.Contains(content, "Assistant:") || strings.Contains(content, "assistant:") {
				role = AssistantMessage
			} else {
				// Skip if we can't determine the role
				continue
			}
		}

		// Detect agent if any
		agentType, _ := stats.DetectAgent(entry)

		// Extract tools used if any
		toolsUsed := extractToolsUsed(entry)

		// Create message event
		event := MessageEvent{
			ID:          fmt.Sprintf("%s-%d", entry.SessionID, entry.Timestamp.UnixNano()),
			Timestamp:   entry.Timestamp,
			SessionID:   entry.SessionID,
			Role:        role,
			Content:     truncateContent(string(entry.Content), 200),
			TokenCount:  0, // Could be extracted from logs if available
			AgentUsed:   agentType,
			ToolsUsed:   toolsUsed,
			Success:     true, // Default to success, could be improved with error detection
			Duration:    1000 * time.Millisecond, // Default duration
		}

		events = append(events, event)
	}


	return events
}

// determineMessageRole analyzes log entry to determine message role
func determineMessageRole(entry stats.ClaudeLogEntry) MessageRole {
	// Use the Type field from the log entry
	switch strings.ToLower(entry.Type) {
	case "user":
		return UserMessage
	case "assistant":
		return AssistantMessage
	case "system":
		return SystemMessage
	case "tool", "function":
		return ToolMessage
	default:
		// For unknown types, try to infer from content
		content := strings.ToLower(string(entry.Content))
		if strings.Contains(content, "tool") || strings.Contains(content, "function") {
			return ToolMessage
		}
		// Default to assistant for most Claude Code entries
		return AssistantMessage
	}
}

// extractToolsUsed extracts tool usage information from log entry
func extractToolsUsed(entry stats.ClaudeLogEntry) []string {
	var tools []string
	content := strings.ToLower(string(entry.Content))

	// Common tools that might appear in logs
	toolPatterns := []string{
		"bash", "read", "write", "edit", "grep", "glob",
		"webfetch", "websearch", "task", "todo",
	}

	for _, tool := range toolPatterns {
		if strings.Contains(content, tool) {
			tools = append(tools, tool)
		}
	}

	return tools
}

// truncateContent truncates content to specified length
func truncateContent(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	if maxLen <= 3 {
		return content[:maxLen]
	}
	return content[:maxLen-3] + "..."
}

// GetTimeFilter returns the current time filter
func (m *DashboardModel) GetTimeFilter() string {
	return m.timeFilter
}

// GetAgentStats returns the current agent statistics
func (m *DashboardModel) GetAgentStats() map[string]*stats.GlobalAgentStats {
	return m.agentStats
}

// renderFilterMenu renders the interactive filter menu
func (m *DashboardModel) renderFilterMenu() string {
	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		Width(m.width)

	content.WriteString(titleStyle.Render("🔧 Time Range Filter") + "\n\n")

	// Filter options
	filterStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Margin(0, 0, 1, 0).
		Width(40)

	selectedStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Margin(0, 0, 1, 0).
		Width(40).
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230"))

	currentStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Margin(0, 0, 1, 0).
		Width(40).
		Foreground(lipgloss.Color("46"))

	for i, timeRange := range TimeRanges {
		label := fmt.Sprintf("%d. %s - %s", i+1, timeRange.Label, timeRange.Desc)

		if i == m.filterMenu.SelectedIndex {
			// Currently selected in menu
			content.WriteString(selectedStyle.Render("▶ " + label) + "\n")
		} else if timeRange.Label == m.currentTimeRange.Label {
			// Currently active filter
			content.WriteString(currentStyle.Render("✓ " + label) + "\n")
		} else {
			// Regular option
			content.WriteString(filterStyle.Render("  " + label) + "\n")
		}
	}

	// Instructions
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Align(lipgloss.Center).
		Width(m.width).
		Margin(2, 0, 0, 0)

	instructions := "\n" + instructionStyle.Render("Use ↑↓ or 1-6 keys to select, Enter to apply, Esc to cancel")
	content.WriteString(instructions)

	// Wrap in border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(2).
		Align(lipgloss.Center).
		Width(m.width - 4).
		Height(m.height - 4)

	return borderStyle.Render(content.String())
}