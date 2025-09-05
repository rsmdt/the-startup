package log

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestProcessHook(t *testing.T) {
	// Create temporary directory for test logs
	tmpDir := t.TempDir()
	originalPath := os.Getenv("THE_STARTUP_PATH")
	os.Setenv("THE_STARTUP_PATH", tmpDir)
	defer os.Setenv("THE_STARTUP_PATH", originalPath)
	
	// Create logs directory
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	tests := []struct {
		name      string
		input     string
		wantEvent string
		wantTool  string
		wantError bool
	}{
		{
			name: "PreToolUse event",
			input: `{
				"hook_event_name": "PreToolUse",
				"tool_name": "Edit",
				"session_id": "session123",
				"cwd": "/test",
				"tool_input": {"file": "test.go", "content": "test"}
			}`,
			wantEvent: "PreToolUse",
			wantTool:  "Edit",
		},
		{
			name: "PostToolUse event with success",
			input: `{
				"hook_event_name": "PostToolUse",
				"tool_name": "Read",
				"session_id": "session456",
				"tool_output": {"content": "file contents"},
				"tool_response": {"success": true}
			}`,
			wantEvent: "PostToolUse",
			wantTool:  "Read",
		},
		{
			name: "PostToolUse event with error",
			input: `{
				"hook_event_name": "PostToolUse",
				"tool_name": "Bash",
				"session_id": "session789",
				"error": "command failed",
				"tool_output": {"error": "exit code 1"}
			}`,
			wantEvent: "PostToolUse",
			wantTool:  "Bash",
		},
		{
			name: "Unknown event type",
			input: `{
				"hook_event_name": "UnknownEvent",
				"tool_name": "Test"
			}`,
			wantError: false, // Silent failure
		},
		{
			name:      "Invalid JSON",
			input:     `{invalid json`,
			wantError: false, // Silent failure
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessHook(strings.NewReader(tt.input))
			
			// ProcessHook always returns nil (silent failure)
			if err != nil {
				t.Errorf("ProcessHook() error = %v, want nil", err)
			}
			
			// Check if metrics were written (for valid events)
			if tt.wantEvent != "" {
				// Read today's log file
				today := time.Now().UTC().Format("20060102")
				logFile := filepath.Join(logsDir, today+".jsonl")
				
				// Check if file exists
				if _, err := os.Stat(logFile); os.IsNotExist(err) && tt.wantEvent != "" {
					t.Errorf("Log file not created: %s", logFile)
				}
			}
		})
	}
}

func TestGenerateToolID(t *testing.T) {
	hook1 := HookPayload{
		ToolName:  "Edit",
		SessionID: "session123456789",
		Timestamp: "2025-09-03T12:34:56Z",
	}
	
	hook2 := HookPayload{
		ToolName:  "Edit",
		SessionID: "session123456789",
		Timestamp: "2025-09-03T12:34:56Z",
	}
	
	id1 := generateToolID(hook1)
	id2 := generateToolID(hook2)
	
	// Same input should generate same ID (for correlation)
	if id1 != id2 {
		t.Errorf("generateToolID() not consistent: %v != %v", id1, id2)
	}
	
	// Check format
	if !strings.HasPrefix(id1, "edit_") {
		t.Errorf("Tool ID should start with lowercase tool name: %v", id1)
	}
	
	if !strings.Contains(id1, "20250903T123456Z") {
		t.Errorf("Tool ID should contain timestamp: %v", id1)
	}
	
	if !strings.Contains(id1, "session1") {
		t.Errorf("Tool ID should contain session prefix: %v", id1)
	}
}

func TestExtractSuccess(t *testing.T) {
	tests := []struct {
		name     string
		hook     HookPayload
		want     *bool
	}{
		{
			name: "explicit success true",
			hook: HookPayload{
				ToolResponse: json.RawMessage(`{"success": true}`),
			},
			want: boolPtr(true),
		},
		{
			name: "explicit success false",
			hook: HookPayload{
				ToolResponse: json.RawMessage(`{"success": false}`),
			},
			want: boolPtr(false),
		},
		{
			name: "error in output",
			hook: HookPayload{
				ToolOutput: json.RawMessage(`{"error": "something failed"}`),
			},
			want: boolPtr(false),
		},
		{
			name: "error field in hook",
			hook: HookPayload{
				Error: "command failed",
			},
			want: boolPtr(false),
		},
		{
			name: "status success",
			hook: HookPayload{
				ToolOutput: json.RawMessage(`{"status": "success"}`),
			},
			want: boolPtr(true),
		},
		{
			name: "status ok",
			hook: HookPayload{
				ToolOutput: json.RawMessage(`{"status": "ok"}`),
			},
			want: boolPtr(true),
		},
		{
			name: "output without error",
			hook: HookPayload{
				ToolOutput: json.RawMessage(`{"result": "data"}`),
			},
			want: boolPtr(true),
		},
		{
			name: "no output",
			hook: HookPayload{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSuccess(tt.hook)
			
			if tt.want == nil {
				if got != nil {
					t.Errorf("extractSuccess() = %v, want nil", *got)
				}
			} else if got == nil {
				t.Errorf("extractSuccess() = nil, want %v", *tt.want)
			} else if *got != *tt.want {
				t.Errorf("extractSuccess() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}