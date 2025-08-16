package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewTreeSelector(t *testing.T) {
	root := &TreeNode{
		Name:  "Root",
		IsDir: true,
		Children: []*TreeNode{
			{Name: "dir1", IsDir: true, Children: []*TreeNode{
				{Name: "file1.txt", IsDir: false, Selected: true},
			}},
			{Name: "file2.txt", IsDir: false, Selected: false},
		},
	}

	selector := NewTreeSelector("Test Selector", root)

	if selector == nil {
		t.Fatal("NewTreeSelector returned nil")
	}

	if selector.title != "Test Selector" {
		t.Errorf("Expected title 'Test Selector', got '%s'", selector.title)
	}

	if len(selector.nodes) == 0 {
		t.Error("Expected nodes to be flattened")
	}

	// Should have flattened: dir1, file1.txt, file2.txt (3 nodes)
	expectedNodes := 3
	if len(selector.nodes) != expectedNodes {
		t.Errorf("Expected %d flattened nodes, got %d", expectedNodes, len(selector.nodes))
	}
}

func TestTreeSelectorMouseInput(t *testing.T) {
	root := &TreeNode{
		Name:  "Root",
		IsDir: true,
		Children: []*TreeNode{
			{Name: "file1.txt", IsDir: false, Selected: true},
			{Name: "file2.txt", IsDir: false, Selected: false},
		},
	}

	selector := NewTreeSelector("Test", root)

	// Test mouse wheel up
	wheelUpMsg := tea.MouseMsg{Type: tea.MouseWheelUp}
	selector.Update(wheelUpMsg)
	// Cursor should stay at 0 (can't go negative)
	if selector.cursor != 0 {
		t.Errorf("Expected cursor to stay at 0, got %d", selector.cursor)
	}

	// Test mouse wheel down
	wheelDownMsg := tea.MouseMsg{Type: tea.MouseWheelDown}
	selector.Update(wheelDownMsg)
	if selector.cursor != 1 {
		t.Errorf("Expected cursor to move to 1, got %d", selector.cursor)
	}

	// Test mouse click
	clickMsg := tea.MouseMsg{Type: tea.MouseLeft, X: 5, Y: 4} // Click on item 1 (0-indexed + 3 offset for title)
	selector.Update(clickMsg)

	// Should have toggled selection of second file (from false to true)
	if !selector.nodes[1].Selected {
		t.Error("Expected second file to be selected after click")
	}
}

func TestTreeSelectorKeyboardShortcuts(t *testing.T) {
	root := &TreeNode{
		Name:  "Root",
		IsDir: true,
		Children: []*TreeNode{
			{Name: "file1.txt", IsDir: false, Selected: true},
			{Name: "file2.txt", IsDir: false, Selected: false},
			{Name: "file3.txt", IsDir: false, Selected: false},
		},
	}

	selector := NewTreeSelector("Test", root)

	// Test 'a' key (toggle all)
	aMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	selector.Update(aMsg)

	// Since not all files were selected initially, 'a' should select all
	allSelected := true
	for _, node := range selector.nodes {
		if !node.IsDir && !node.Selected {
			allSelected = false
			break
		}
	}
	if !allSelected {
		t.Error("Expected all files to be selected after 'a' key")
	}

	// Test 'a' again (should deselect all)
	selector.Update(aMsg)
	anySelected := false
	for _, node := range selector.nodes {
		if !node.IsDir && node.Selected {
			anySelected = true
			break
		}
	}
	if anySelected {
		t.Error("Expected all files to be deselected after second 'a' key")
	}

	// Test 'j' and 'k' keys (vim-style navigation)
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	selector.Update(jMsg)
	if selector.cursor != 1 {
		t.Errorf("Expected cursor at 1 after 'j', got %d", selector.cursor)
	}

	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	selector.Update(kMsg)
	if selector.cursor != 0 {
		t.Errorf("Expected cursor at 0 after 'k', got %d", selector.cursor)
	}
}

func TestTreeSelectorDirectoryHandling(t *testing.T) {
	root := &TreeNode{
		Name:  "Root",
		IsDir: true,
		Children: []*TreeNode{
			{
				Name:     "components",
				IsDir:    true,
				Selected: false,
				Children: []*TreeNode{
					{Name: "comp1.txt", IsDir: false, Selected: true},
					{Name: "comp2.txt", IsDir: false, Selected: true},
				},
			},
		},
	}

	// Set parent references
	for _, child := range root.Children {
		child.Parent = root
		for _, grandchild := range child.Children {
			grandchild.Parent = child
		}
	}

	selector := NewTreeSelector("Test", root)

	// Initially, directory should be selected because all children are selected
	dirNode := selector.nodes[0] // First node should be the directory
	if !dirNode.Selected {
		t.Error("Expected directory to be selected when all children are selected")
	}

	// Toggle directory (should affect all children)
	selector.cursor = 0 // Position on directory
	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
	selector.Update(spaceMsg)

	// All children should now be deselected
	for _, node := range selector.nodes {
		if !node.IsDir && node.Selected {
			t.Error("Expected all files to be deselected after toggling directory")
		}
	}
}

func TestTreeSelectorViewRendering(t *testing.T) {
	root := &TreeNode{
		Name:  "Root",
		IsDir: true,
		Children: []*TreeNode{
			{Name: "selected.txt", IsDir: false, Selected: true},
			{Name: "unselected.txt", IsDir: false, Selected: false},
			{Name: "existing.txt", IsDir: false, Selected: true, Exists: true},
		},
	}

	selector := NewTreeSelector("File Selection", root)
	view := selector.View()

	// Should contain title
	if !containsString(view, "File Selection") {
		t.Error("Expected view to contain title")
	}

	// Should contain file names
	if !containsString(view, "selected.txt") {
		t.Error("Expected view to contain selected file")
	}

	if !containsString(view, "unselected.txt") {
		t.Error("Expected view to contain unselected file")
	}

	// Should show update indicator for existing files
	if !containsString(view, "update") {
		t.Error("Expected view to show update indicator for existing files")
	}

	// Should contain help text
	if !containsString(view, "navigate") {
		t.Error("Expected view to contain navigation help")
	}

	// Test cursor indication
	if !containsString(view, ">") {
		t.Error("Expected view to show cursor indicator")
	}
}

func TestTreeSelectorViewportScrolling(t *testing.T) {
	// Create a tree with many files to test scrolling
	children := make([]*TreeNode, 20)
	for i := 0; i < 20; i++ {
		name := "file" + string(rune('0'+i/10)) + string(rune('0'+i%10)) + ".txt"
		children[i] = &TreeNode{
			Name:     name,
			IsDir:    false,
			Selected: false,
		}
	}

	root := &TreeNode{
		Name:     "Root",
		IsDir:    true,
		Children: children,
	}

	selector := NewTreeSelector("Large Tree", root)
	selector.height = 10 // Small viewport

	// Move cursor to bottom
	selector.cursor = 19

	view := selector.View()

	// Should still render without errors
	if view == "" {
		t.Error("Expected non-empty view for large tree")
	}

	// Should contain some files when scrolled
	if len(view) < 100 {
		t.Error("Expected substantial view content for large tree")
	}
}

func TestTreeSelectorCompletion(t *testing.T) {
	root := &TreeNode{
		Name:  "Root",
		IsDir: true,
		Children: []*TreeNode{
			{Name: "file1.txt", Path: "agents/file1.txt", IsDir: false, Selected: true},
			{Name: "file2.txt", Path: "commands/file2.txt", IsDir: false, Selected: false},
		},
	}

	selector := NewTreeSelector("Test", root)

	// Complete selection
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	selector.Update(enterMsg)

	if !selector.done {
		t.Error("Expected selector to be done after Enter")
	}

	if selector.cancelled {
		t.Error("Expected selector not to be cancelled after Enter")
	}

	// Get selected paths
	paths := selector.GetSelectedPaths()
	if len(paths) != 1 {
		t.Errorf("Expected 1 selected path, got %d", len(paths))
	}

	if paths[0] != "agents/file1.txt" {
		t.Errorf("Expected path 'agents/file1.txt', got '%s'", paths[0])
	}
}

func TestTreeSelectorCancellation(t *testing.T) {
	root := &TreeNode{
		Name:     "Root",
		IsDir:    true,
		Children: []*TreeNode{},
	}

	selector := NewTreeSelector("Test", root)

	// Cancel with ESC
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	selector.Update(escMsg)

	if !selector.done {
		t.Error("Expected selector to be done after ESC")
	}

	if !selector.cancelled {
		t.Error("Expected selector to be cancelled after ESC")
	}

	// Cancel with Ctrl+C
	selector = NewTreeSelector("Test", root) // Reset
	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	selector.Update(ctrlCMsg)

	if !selector.cancelled {
		t.Error("Expected selector to be cancelled after Ctrl+C")
	}
}

func TestTreeSelectorUpdatingFiles(t *testing.T) {
	root := &TreeNode{
		Name:  "Root",
		IsDir: true,
		Children: []*TreeNode{
			{Name: "new.txt", Path: "agents/new.txt", IsDir: false, Selected: true, Exists: false},
			{Name: "existing.txt", Path: "agents/existing.txt", IsDir: false, Selected: true, Exists: true},
			{Name: "unselected.txt", Path: "agents/unselected.txt", IsDir: false, Selected: false, Exists: true},
		},
	}

	selector := NewTreeSelector("Test", root)

	// Get updating files (selected + existing)
	updatingFiles := selector.GetUpdatingFiles()

	if len(updatingFiles) != 1 {
		t.Errorf("Expected 1 updating file, got %d", len(updatingFiles))
	}

	if updatingFiles[0] != "agents/existing.txt" {
		t.Errorf("Expected 'agents/existing.txt', got '%s'", updatingFiles[0])
	}
}

func TestTreeSelectorHelpInteraction(t *testing.T) {
	root := &TreeNode{Name: "Root", IsDir: true, Children: []*TreeNode{}}
	selector := NewTreeSelector("Test", root)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	selector.Update(helpMsg)

	if !selector.showHelp {
		t.Error("Expected help to be shown")
	}

	// Try navigation while help is shown (should be ignored)
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	oldCursor := selector.cursor
	selector.Update(downMsg)

	if selector.cursor != oldCursor {
		t.Error("Expected cursor to not move while help is shown")
	}

	// Try selection while help is shown (should be ignored)
	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
	selector.Update(spaceMsg)
	// Should not affect any selections (tested by no crash)

	// Close help with ESC
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	selector.Update(escMsg)

	if selector.showHelp {
		t.Error("Expected help to be hidden after ESC")
	}
}
