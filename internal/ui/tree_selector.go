package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
	cancelled bool   // Indicates user cancelled with ESC
	showHelp  bool   // Show help overlay
	styles    Styles // Theme styles
}

func NewTreeSelector(title string, root *TreeNode) *TreeSelector {
	ts := &TreeSelector{
		root:   root,
		title:  title,
		cursor: 0,
		width:  80,
		height: 20,
		styles: GetStyles(),
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

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if !ts.showHelp && ts.cursor > 0 {
				ts.cursor--
			}

		case tea.MouseWheelDown:
			if !ts.showHelp && ts.cursor < len(ts.nodes)-1 {
				ts.cursor++
			}

		case tea.MouseLeft:
			if !ts.showHelp {
				// Calculate which item was clicked based on Y position
				// Account for title (2 lines) and top margin
				itemY := msg.Y - 3
				if itemY >= 0 && itemY < len(ts.nodes) {
					ts.cursor = itemY
					// Toggle the clicked item
					if ts.cursor < len(ts.nodes) {
						node := ts.nodes[ts.cursor]
						if node.IsDir {
							newState := !node.Selected
							ts.setNodeAndChildren(node, newState)
						} else {
							node.Selected = !node.Selected
						}
						ts.updateDirectorySelections()
					}
				}
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if ts.showHelp {
				ts.showHelp = false
			} else {
				ts.done = true
				ts.cancelled = true
				return ts, tea.Quit
			}

		case "esc":
			if ts.showHelp {
				ts.showHelp = false
			} else {
				ts.done = true
				ts.cancelled = true
				return ts, tea.Quit
			}

		case "?":
			ts.showHelp = !ts.showHelp

		case "enter":
			if !ts.showHelp {
				ts.done = true
				return ts, tea.Quit
			}

		case "up", "k":
			if !ts.showHelp && ts.cursor > 0 {
				ts.cursor--
			}

		case "down", "j":
			if !ts.showHelp && ts.cursor < len(ts.nodes)-1 {
				ts.cursor++
			}

		case " ":
			if !ts.showHelp && ts.cursor < len(ts.nodes) {
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
			if !ts.showHelp {
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
	}

	return ts, nil
}

func (ts *TreeSelector) View() string {
	if ts.done {
		return ""
	}

	if ts.showHelp {
		return ts.renderHelp()
	}

	var s strings.Builder

	// Title
	s.WriteString(ts.styles.Title.Render(ts.title))
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
			// Primary color ">" for cursor
			line = ts.styles.Cursor.Render(">") + " "
		} else {
			line = "  "
		}

		// Indentation
		line += indent

		// Build rest of line content
		var content string

		// Selection state (filled/empty circle)
		if node.Selected {
			content += IconSelected
		} else {
			content += IconUnselected
		}

		// Add space and name
		content += " " + node.Name

		// Add update indicator if file exists
		if node.Exists && !node.IsDir {
			content += " " + ts.styles.Warning.Render("("+IconUpdate+" update)")
		}

		// Apply appropriate styling to content
		if isCursor {
			// Bright text for cursor line content
			line += ts.styles.CursorLine.Render(content)
		} else {
			// Normal text for non-cursor lines
			line += ts.styles.Normal.Render(content)
		}

		s.WriteString(line)
		s.WriteString("\n")
	}

	// Help text
	s.WriteString("\n")
	s.WriteString(ts.styles.Help.Render("↑↓/jk/mouse: navigate • space/click: toggle • a: all • enter: confirm • ?: help • esc: back"))

	return s.String()
}

func (ts *TreeSelector) renderHelp() string {
	var s strings.Builder

	// Help overlay
	s.WriteString(ts.styles.Title.Render("🔑 Keyboard Shortcuts"))
	s.WriteString("\n\n")

	helpItems := []struct {
		key  string
		desc string
	}{
		{"↑/k", "Move up"},
		{"↓/j", "Move down"},
		{"space", "Toggle selection"},
		{"a", "Toggle all files"},
		{"enter", "Confirm selection"},
		{"esc", "Cancel/Go back"},
		{"?", "Toggle this help"},
		{"q", "Quit"},
	}

	for _, item := range helpItems {
		key := ts.styles.Info.Render(item.key)
		s.WriteString(fmt.Sprintf("  %-20s %s\n", key, item.desc))
	}

	s.WriteString("\n")
	s.WriteString(ts.styles.Help.Render("Press any key to close this help"))

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
	// Remove indentation for consistent flat display
	return ""
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
	p := tea.NewProgram(NewTreeSelector(title, root), tea.WithMouseCellMotion())
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
