package stats_test

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rsmdt/the-startup/internal/stats"
)

// ExampleJSONLParser demonstrates how to use the JSONL parser
func ExampleJSONLParser() {
	// Sample JSONL data
	jsonlData := `{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:00Z","content":{"text":"How do I implement a REST API?"}}
{"type":"assistant","session_id":"sess_001","timestamp":"2025-01-01T10:00:05Z","content":{"text":"I'll help you create a REST API.","tool_uses":[{"id":"tool_001","name":"Write","parameters":{"file_path":"main.go","content":"package main\n\nfunc main() {\n\t// API code here\n}\n"}}]}}
{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:10Z","content":{"tool_use_id":"tool_001","tool_name":"Write","result":{"message":"File created successfully"}}}
{"type":"summary","session_id":"sess_001","timestamp":"2025-01-01T10:00:20Z","content":{"duration":20000000000,"tools_used":1,"tokens_used":{"input":150,"output":200,"total":350}}}`

	// Create a new parser
	parser := stats.NewJSONLParser()

	// Set parsing options
	parser.SetOptions(stats.ParseOptions{
		BufferSize: 64 * 1024,
		ParseTools: true,
		StrictMode: false,
	})

	// Parse the stream
	reader := strings.NewReader(jsonlData)
	entries, errors := parser.ParseStream(reader)

	// Handle errors in a separate goroutine
	go func() {
		for err := range errors {
			log.Printf("Parse error: %v", err)
		}
	}()

	// Process entries
	for entry := range entries {
		switch entry.Type {
		case "user":
			if entry.User.IsToolResult {
				fmt.Printf("Tool result: %s\n", entry.User.ToolName)
			} else {
				fmt.Printf("User: %s\n", entry.User.Text)
			}
		case "assistant":
			fmt.Printf("Assistant: %s\n", entry.Assistant.Text)
			for _, tool := range entry.Assistant.ToolUses {
				fmt.Printf("  Tool used: %s\n", tool.Name)
			}
		case "summary":
			fmt.Printf("Session summary: %d tokens used\n", entry.Summary.TokensUsed.Total)
		}
	}

	// Output:
	// User: How do I implement a REST API?
	// Assistant: I'll help you create a REST API.
	//   Tool used: Write
	// Tool result: Write
	// Session summary: 350 tokens used
}

// ExampleJSONLParser_filtering demonstrates filtering capabilities
func ExampleJSONLParser_filtering() {
	// Sample data with multiple sessions
	jsonlData := `{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:00Z","content":{"text":"Session 1 message"}}
{"type":"user","session_id":"sess_002","timestamp":"2025-01-01T11:00:00Z","content":{"text":"Session 2 message"}}
{"type":"assistant","session_id":"sess_001","timestamp":"2025-01-01T10:00:05Z","content":{"text":"Response","tool_uses":[{"id":"t1","name":"Bash","parameters":{"command":"ls"}}]}}
{"type":"assistant","session_id":"sess_002","timestamp":"2025-01-01T11:00:05Z","content":{"text":"Response","tool_uses":[{"id":"t2","name":"Read","parameters":{"file_path":"test.txt"}}]}}`

	parser := stats.NewJSONLParser()

	// Filter for specific session and tools
	parser.SetFilter(stats.FilterOptions{
		SessionIDs:   []string{"sess_001"},
		IncludeTools: []string{"Bash"},
	})

	reader := strings.NewReader(jsonlData)
	entries, _ := parser.ParseStream(reader)

	count := 0
	for entry := range entries {
		count++
		fmt.Printf("Entry %d: type=%s, session=%s\n", count, entry.Type, entry.SessionID)
	}

	// Output:
	// Entry 1: type=user, session=sess_001
	// Entry 2: type=assistant, session=sess_001
}

// ExampleJSONLParser_largeFile demonstrates streaming large files
func ExampleJSONLParser_largeFile() {
	// Simulate a large file processing scenario
	parser := stats.NewJSONLParser()

	// Configure for optimal large file processing
	parser.SetOptions(stats.ParseOptions{
		BufferSize:     256 * 1024, // 256KB buffer
		ParseTools:     false,       // Skip tool parsing for speed
		SkipSystemLogs: true,        // Skip system messages
	})

	// Set up time-based filtering
	startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	parser.SetFilter(stats.FilterOptions{
		StartTime:  &startTime,
		MaxEntries: 1000, // Process only first 1000 matching entries
	})

	// In a real scenario, you would open a file:
	// file, _ := os.Open("large-log.jsonl")
	// defer file.Close()
	// entries, errors := parser.ParseStream(file)

	fmt.Println("Parser configured for large file processing")
	// Output:
	// Parser configured for large file processing
}