package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/uninstaller"
)

func TestNewUninstallModel(t *testing.T) {
	model := NewUninstallModel(true, false, true, false)
	
	if model == nil {
		t.Fatal("NewUninstallModel returned nil")
	}
	
	if model.state != UninstallStateInit {
		t.Errorf("Expected initial state to be UninstallStateInit, got %v", model.state)
	}
	
	if !model.dryRun {
		t.Error("Expected dryRun to be true")
	}
	
	if model.forceRemove {
		t.Error("Expected forceRemove to be false")
	}
	
	if !model.createBackup {
		t.Error("Expected createBackup to be true")
	}
}

func TestUninstallModelInit(t *testing.T) {
	model := NewUninstallModel(true, false, true, false)
	cmd := model.Init()
	
	if cmd == nil {
		t.Error("Init() should return a command to start the process")
	}
}

func TestUninstallModelStateTransitions(t *testing.T) {
	model := NewUninstallModel(true, false, true, false)
	
	// Test key handling
	quitMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, cmd := model.Update(quitMsg)
	
	uninstallModel, ok := newModel.(*UninstallModel)
	if !ok {
		t.Fatal("Update should return *UninstallModel")
	}
	
	if !uninstallModel.Ready() {
		t.Error("Model should be ready to quit after Ctrl+C")
	}
	
	if cmd == nil {
		t.Error("Update should return tea.Quit command")
	}
}

func TestUninstallModelOptions(t *testing.T) {
	model := NewUninstallModel(false, true, false, true)
	
	if model.dryRun {
		t.Error("Expected dryRun to be false")
	}
	
	if !model.forceRemove {
		t.Error("Expected forceRemove to be true")
	}
	
	if model.createBackup {
		t.Error("Expected createBackup to be false")
	}
	
	if !model.verbose {
		t.Error("Expected verbose to be true")
	}
}

func TestNewRemovalPreviewModel(t *testing.T) {
	// Create a minimal preview for testing
	preview := &uninstaller.RemovalPreview{
		InstallPath:     "/test/install",
		ClaudePath:      "/test/claude",
		DiscoverySource: uninstaller.DiscoverySourceLockfile,
		Files:           []uninstaller.FileInfo{},
		CategorySummary: []uninstaller.CategorySummary{},
		TotalFiles:      0,
		TotalSize:       0,
		UntrackedFiles:  []uninstaller.FileInfo{},
		OrphanedFiles:   []uninstaller.FileInfo{},
		SettingsFiles:   []uninstaller.FileInfo{},
		SecurityIssues:  []uninstaller.SecurityIssue{},
		ValidationErrors: []uninstaller.ValidationError{},
	}
	
	model := NewRemovalPreviewModel(preview, true)
	
	if model == nil {
		t.Fatal("NewRemovalPreviewModel returned nil")
	}
	
	if model.preview != preview {
		t.Error("Preview was not set correctly")
	}
	
	if !model.dryRun {
		t.Error("Expected dryRun to be true")
	}
}

func TestRemovalPreviewModelInit(t *testing.T) {
	preview := &uninstaller.RemovalPreview{
		InstallPath:      "/test/install",
		ClaudePath:       "/test/claude",
		Files:            []uninstaller.FileInfo{},
		CategorySummary:  []uninstaller.CategorySummary{},
		UntrackedFiles:   []uninstaller.FileInfo{},
		OrphanedFiles:    []uninstaller.FileInfo{},
		SettingsFiles:    []uninstaller.FileInfo{},
		SecurityIssues:   []uninstaller.SecurityIssue{},
		ValidationErrors: []uninstaller.ValidationError{},
	}
	
	model := NewRemovalPreviewModel(preview, false)
	cmd := model.Init()
	
	// Init should return nil for RemovalPreviewModel
	if cmd != nil {
		t.Error("RemovalPreviewModel Init() should return nil")
	}
}

func TestRemovalPreviewModelResize(t *testing.T) {
	preview := &uninstaller.RemovalPreview{
		InstallPath:      "/test/install",
		ClaudePath:       "/test/claude",
		Files:            []uninstaller.FileInfo{},
		CategorySummary:  []uninstaller.CategorySummary{},
		UntrackedFiles:   []uninstaller.FileInfo{},
		OrphanedFiles:    []uninstaller.FileInfo{},
		SettingsFiles:    []uninstaller.FileInfo{},
		SecurityIssues:   []uninstaller.SecurityIssue{},
		ValidationErrors: []uninstaller.ValidationError{},
	}
	
	model := NewRemovalPreviewModel(preview, false)
	
	originalWidth := model.width
	originalHeight := model.height
	
	model.SetSize(120, 40)
	
	if model.width == originalWidth || model.height == originalHeight {
		t.Error("SetSize should update model dimensions")
	}
	
	if model.width != 120 || model.height != 40 {
		t.Errorf("Expected dimensions 120x40, got %dx%d", model.width, model.height)
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}
	
	for _, test := range tests {
		result := formatBytes(test.bytes)
		if result != test.expected {
			t.Errorf("formatBytes(%d) = %s, expected %s", test.bytes, result, test.expected)
		}
	}
}