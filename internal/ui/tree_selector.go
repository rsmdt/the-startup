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
	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575"))
	
	cursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF06B7")).
		Bold(true)
	
	dirStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3C7EFF")).
		Bold(true)
	
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666"))
	
	checkMark = "✓"
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
		
		// Build the line
		line := ""
		
		// Cursor indicator
		if i == ts.cursor {
			line += "> "
		} else {
			line += "  "
		}
		
		// Indentation
		line += indent
		
		// Selection indicator - dot or check
		if node.Selected {
			line += fmt.Sprintf("%s ", checkMark)
		} else {
			line += "• "
		}
		
		// Name
		name := node.Name
		if node.IsDir {
			// Format directory names without trailing slash, in directory style
			line += dirStyle.Render(name)
		} else {
			line += name
		}
		
		// Apply cursor highlighting to the whole line if this is the current item
		if i == ts.cursor {
			line = cursorStyle.Render(line)
		} else if node.Selected && !node.IsDir {
			// Apply selected style to selected files (but not the cursor indicator)
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				line = parts[0] + " " + selectedStyle.Render(parts[1])
			}
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