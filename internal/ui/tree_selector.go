package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TreeNode struct {
	Name     string
	Path     string
	IsDir    bool
	Selected bool
	Exists   bool // Whether this file already exists
	Children []*TreeNode
	Parent   *TreeNode
}

type TreeSelector struct {
	root      *TreeNode
	cursor    int
	nodes     []*TreeNode // Flattened view of all nodes
	width     int
	height    int
	title     string
	done      bool
	cancelled bool // Indicates user cancelled with ESC
}

var (
	// Match huh library default theme colors
	cursorIndicatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")) // Pink/magenta for ">" cursor indicator only
	
	cursorLineStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")) // Green for text on cursor line
	
	normalItemStyle = lipgloss.NewStyle() // Default/neutral color for non-cursor lines
	
	existsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")) // Orange/amber for update indicator
	
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")) // Darker gray for help text
	
	filledCircle = "●"  // Selected item indicator
	emptyCircle = "○"   // Unselected item indicator
	updateMark = "↻"    // Indicator for files that will be updated
)

func NewTreeSelector(title string, root *TreeNode) *TreeSelector {
	ts := &TreeSelector{
		root:   root,
		title:  title,
		cursor: 0,
		width:  80,
		height: 20,
	}
	ts.flatten()
	ts.updateDirectorySelections()
	return ts
}

func (ts *TreeSelector) Init() tea.Cmd {
	return nil
}

func (ts *TreeSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ts.width = msg.Width
		ts.height = msg.Height - 6 // Leave room for title and help
	
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			ts.done = true
			ts.cancelled = true
			return ts, tea.Quit
		
		case "esc":
			ts.done = true
			ts.cancelled = true
			return ts, tea.Quit
		
		case "enter":
			ts.done = true
			return ts, tea.Quit
		
		case "up", "k":
			if ts.cursor > 0 {
				ts.cursor--
			}
		
		case "down", "j":
			if ts.cursor < len(ts.nodes)-1 {
				ts.cursor++
			}
		
		case " ":
			if ts.cursor < len(ts.nodes) {
				node := ts.nodes[ts.cursor]
				if node.IsDir {
					// Toggle all children when space is pressed on a directory
					newState := !node.Selected
					ts.setNodeAndChildren(node, newState)
				} else {
					// Toggle individual file
					node.Selected = !node.Selected
				}
				// Update directory selections based on children
				ts.updateDirectorySelections()
			}
		
		case "a", "A":
			// Toggle all files on/off
			allSelected := true
			for _, node := range ts.nodes {
				if !node.IsDir && !node.Selected {
					allSelected = false
					break
				}
			}
			
			// Set all files to opposite of current state
			newState := !allSelected
			ts.setAllFiles(ts.root, newState)
			ts.updateDirectorySelections()
		}
	}
	
	return ts, nil
}

func (ts *TreeSelector) View() string {
	if ts.done {
		return ""
	}
	
	var s strings.Builder
	
	// Title
	s.WriteString(lipgloss.NewStyle().Bold(true).Render(ts.title))
	s.WriteString("\n\n")
	
	// Tree view
	start := 0
	end := len(ts.nodes)
	
	// Viewport scrolling
	if ts.height > 0 && len(ts.nodes) > ts.height {
		// Keep cursor in view
		if ts.cursor < start {
			start = ts.cursor
		} else if ts.cursor >= start+ts.height {
			start = ts.cursor - ts.height + 1
		}
		end = start + ts.height
		if end > len(ts.nodes) {
			end = len(ts.nodes)
			start = end - ts.height
		}
	}
	
	for i := start; i < end && i < len(ts.nodes); i++ {
		node := ts.nodes[i]
		indent := ts.getIndent(node)
		
		// Build the line (matching huh style)
		var line string
		
		// Cursor indicator
		isCursor := i == ts.cursor
		
		if isCursor {
			// Pink ">" for cursor, rest of line will be green
			line = cursorIndicatorStyle.Render(">") + " "
		} else {
			line = "  "
		}
		
		// Indentation
		line += indent
		
		// Build rest of line content
		var content string
		
		// Selection state (filled/empty circle)
		if node.Selected {
			content += filledCircle
		} else {
			content += emptyCircle
		}
		
		// Add space and name
		content += " " + node.Name
		
		// Add update indicator if file exists
		if node.Exists && !node.IsDir {
			content += " " + existsStyle.Render("("+updateMark+" update)")
		}
		
		// Apply appropriate styling to content
		if isCursor {
			// Green for cursor line content
			line += cursorLineStyle.Render(content)
		} else {
			// Neutral/default for non-cursor lines
			line += content
		}
		
		s.WriteString(line)
		s.WriteString("\n")
	}
	
	// Help text
	s.WriteString("\n")
	s.WriteString(helpStyle.Render("↑↓/jk: navigate • space: toggle • a: all on/off • enter: confirm • esc: back"))
	
	return s.String()
}

func (ts *TreeSelector) flatten() {
	ts.nodes = []*TreeNode{}
	ts.flattenNode(ts.root, 0)
}

func (ts *TreeSelector) flattenNode(node *TreeNode, level int) {
	// Skip the root node itself, start with its children
	if node == ts.root {
		for _, child := range node.Children {
			ts.flattenNode(child, level)
		}
	} else {
		ts.nodes = append(ts.nodes, node)
		for _, child := range node.Children {
			ts.flattenNode(child, level+1)
		}
	}
}

func (ts *TreeSelector) getIndent(node *TreeNode) string {
	level := 0
	parent := node.Parent
	// Skip root in indentation calculation
	for parent != nil && parent != ts.root {
		level++
		parent = parent.Parent
	}
	return strings.Repeat("  ", level)
}

func (ts *TreeSelector) setAllFiles(node *TreeNode, selected bool) {
	if !node.IsDir {
		node.Selected = selected
	}
	for _, child := range node.Children {
		ts.setAllFiles(child, selected)
	}
}

func (ts *TreeSelector) setNodeAndChildren(node *TreeNode, selected bool) {
	node.Selected = selected
	for _, child := range node.Children {
		ts.setNodeAndChildren(child, selected)
	}
}

func (ts *TreeSelector) updateDirectorySelections() {
	// Update directory selection status based on children
	ts.updateDirNode(ts.root)
}

func (ts *TreeSelector) updateDirNode(node *TreeNode) {
	if !node.IsDir {
		return
	}
	
	// First, update all child directories
	for _, child := range node.Children {
		if child.IsDir {
			ts.updateDirNode(child)
		}
	}
	
	// Then check if all children are selected
	if len(node.Children) > 0 {
		allSelected := true
		for _, child := range node.Children {
			if !child.Selected {
				allSelected = false
				break
			}
		}
		node.Selected = allSelected
	}
}

func (ts *TreeSelector) GetSelectedPaths() []string {
	paths := []string{}
	ts.collectSelectedPaths(ts.root, &paths)
	return paths
}

func (ts *TreeSelector) collectSelectedPaths(node *TreeNode, paths *[]string) {
	if !node.IsDir && node.Selected {
		*paths = append(*paths, node.Path)
	}
	for _, child := range node.Children {
		ts.collectSelectedPaths(child, paths)
	}
}

func (ts *TreeSelector) GetUpdatingFiles() []string {
	paths := []string{}
	ts.collectUpdatingFiles(ts.root, &paths)
	return paths
}

func (ts *TreeSelector) collectUpdatingFiles(node *TreeNode, paths *[]string) {
	if !node.IsDir && node.Selected && node.Exists {
		*paths = append(*paths, node.Path)
	}
	for _, child := range node.Children {
		ts.collectUpdatingFiles(child, paths)
	}
}

func RunTreeSelector(title string, root *TreeNode) ([]string, error) {
	p := tea.NewProgram(NewTreeSelector(title, root))
	model, err := p.Run()
	if err != nil {
		return nil, err
	}
	
	ts := model.(*TreeSelector)
	if ts.cancelled {
		return nil, fmt.Errorf("selection cancelled")
	}
	if ts.done {
		return ts.GetSelectedPaths(), nil
	}
	return nil, fmt.Errorf("selection cancelled")
}
