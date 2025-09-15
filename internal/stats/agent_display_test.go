package stats

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatAgentLeaderboard(t *testing.T) {
	formatter := NewDisplayFormatter()

	// Test with nil/empty map - following FormatToolLeaderboard pattern
	result := formatter.FormatAgentLeaderboard(nil, 0)
	if result != "" {
		t.Error("Expected empty string for nil agent stats")
	}

	result = formatter.FormatAgentLeaderboard(map[string]*GlobalAgentStats{}, 0)
	if result != "" {
		t.Error("Expected empty string for empty agent stats")
	}

	// Create test agent stats by simulating UpdateStats calls properly
	agentStats := make(map[string]*GlobalAgentStats)

	// Create "specify" agent
	specify := &GlobalAgentStats{}
	for i := 0; i < 100; i++ {
		specify.UpdateStats(250) // 250ms average
		if i < 95 {
			specify.SuccessCount++
		} else {
			specify.FailureCount++
		}
	}
	agentStats["specify"] = specify

	// Create "build" agent
	build := &GlobalAgentStats{}
	for i := 0; i < 50; i++ {
		build.UpdateStats(150) // 150ms average
		if i < 49 {
			build.SuccessCount++
		} else {
			build.FailureCount++
		}
	}
	agentStats["build"] = build

	// Create "test" agent
	testAgent := &GlobalAgentStats{}
	for i := 0; i < 30; i++ {
		testAgent.UpdateStats(300) // 300ms average
		if i < 27 {
			testAgent.SuccessCount++
		} else {
			testAgent.FailureCount++
		}
	}
	agentStats["test"] = testAgent

	// Test table format with no limit - following FormatToolLeaderboard pattern
	formatter.SetOutputFormat("table")
	result = formatter.FormatAgentLeaderboard(agentStats, 0)
	if !strings.Contains(result, "AGENT USAGE LEADERBOARD") {
		t.Error("Table format missing header")
	}
	if !strings.Contains(result, "specify") {
		t.Error("Table format missing specify agent")
	}
	if !strings.Contains(result, "Sparkline") {
		t.Error("Table format missing sparkline column")
	}

	// Check order (specify should be first with 100 invocations)
	lines := strings.Split(result, "\n")
	specifyFound := false
	buildFound := false
	for _, line := range lines {
		// tabwriter formats with spaces, not literal tabs
		if strings.Contains(line, "1") && strings.Contains(line, "specify") && strings.Contains(line, "100") {
			specifyFound = true
		}
		if strings.Contains(line, "2") && strings.Contains(line, "build") && strings.Contains(line, "50") {
			buildFound = true
		}
	}
	if !specifyFound {
		t.Error("specify should be ranked #1 with 100 invocations")
	}
	if !buildFound {
		t.Error("build should be ranked #2 with 50 invocations")
	}

	// Test with limit
	result = formatter.FormatAgentLeaderboard(agentStats, 2)
	if strings.Contains(result, "test") {
		t.Error("test should not appear when limit is 2")
	}

	// Test JSON format
	formatter.SetOutputFormat("json")
	result = formatter.FormatAgentLeaderboard(agentStats, 0)
	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Errorf("Invalid JSON output: %v", err)
	}
	if jsonResult["total_agents"].(float64) != 3 {
		t.Errorf("JSON format incorrect total_agents: %v", jsonResult["total_agents"])
	}

	// Verify leaderboard structure
	leaderboard, ok := jsonResult["leaderboard"].([]interface{})
	if !ok {
		t.Error("JSON missing leaderboard array")
	}
	if len(leaderboard) != 3 {
		t.Errorf("Expected 3 agents in leaderboard, got %d", len(leaderboard))
	}

	// Check first agent entry structure
	firstAgent := leaderboard[0].(map[string]interface{})
	if firstAgent["rank"].(float64) != 1 {
		t.Errorf("First agent should have rank 1, got %v", firstAgent["rank"])
	}
	if firstAgent["name"].(string) != "specify" {
		t.Errorf("First agent should be specify, got %s", firstAgent["name"])
	}

	// Test CSV format
	formatter.SetOutputFormat("csv")
	result = formatter.FormatAgentLeaderboard(agentStats, 0)
	if !strings.Contains(result, "Rank,Agent,Total Invocations") {
		t.Error("CSV format missing headers")
	}
	if !strings.Contains(result, "1,specify,100") {
		t.Error("CSV format missing specify entry")
	}
	if !strings.Contains(result, "Success Rate") {
		t.Error("CSV format missing success rate column")
	}
}

func TestFormatAgentLeaderboardSuccessRates(t *testing.T) {
	formatter := NewDisplayFormatter()

	// Create test data with different success rates
	agentStats := make(map[string]*GlobalAgentStats)

	// Perfect agent
	perfect := &GlobalAgentStats{}
	for i := 0; i < 10; i++ {
		perfect.UpdateStats(100)
		perfect.SuccessCount++
	}
	agentStats["perfect"] = perfect

	// Failing agent
	failing := &GlobalAgentStats{}
	for i := 0; i < 20; i++ {
		failing.UpdateStats(100)
		if i < 10 {
			failing.SuccessCount++
		} else {
			failing.FailureCount++
		}
	}
	agentStats["failing"] = failing

	// Test table format shows success rates
	formatter.SetOutputFormat("table")
	result := formatter.FormatAgentLeaderboard(agentStats, 0)
	if !strings.Contains(result, "100.0%") {
		t.Error("Table format should show 100% success rate")
	}
	if !strings.Contains(result, "50.0%") {
		t.Error("Table format should show 50% success rate")
	}

	// Test JSON format includes success rates
	formatter.SetOutputFormat("json")
	result = formatter.FormatAgentLeaderboard(agentStats, 0)
	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Errorf("Invalid JSON output: %v", err)
	}

	leaderboard := jsonResult["leaderboard"].([]interface{})
	for _, item := range leaderboard {
		agent := item.(map[string]interface{})
		stats := agent["stats"].(map[string]interface{})
		successRate, ok := stats["success_rate"].(float64)
		if !ok {
			t.Error("JSON format missing success_rate field")
		}
		if agent["name"].(string) == "perfect" && successRate != 1.0 {
			t.Errorf("Perfect agent should have 1.0 success rate, got %v", successRate)
		}
		if agent["name"].(string) == "failing" && successRate != 0.5 {
			t.Errorf("Failing agent should have 0.5 success rate, got %v", successRate)
		}
	}

	// Test CSV format includes success rates
	formatter.SetOutputFormat("csv")
	result = formatter.FormatAgentLeaderboard(agentStats, 0)
	if !strings.Contains(result, ",1.00") { // 100% as 1.00
		t.Error("CSV format should include perfect success rate")
	}
	if !strings.Contains(result, ",0.50") { // 50% as 0.50
		t.Error("CSV format should include failing success rate")
	}
}

func TestFormatAgentLeaderboardSparklines(t *testing.T) {
	formatter := NewDisplayFormatter()
	formatter.SetOutputFormat("table")

	// Create test data with different invocation counts for sparkline testing
	agentStats := make(map[string]*GlobalAgentStats)

	// High usage agent
	high := &GlobalAgentStats{}
	for i := 0; i < 100; i++ {
		high.UpdateStats(100)
		if i < 95 {
			high.SuccessCount++
		} else {
			high.FailureCount++
		}
	}
	agentStats["high"] = high

	// Medium usage agent
	medium := &GlobalAgentStats{}
	for i := 0; i < 50; i++ {
		medium.UpdateStats(100)
		if i < 45 {
			medium.SuccessCount++
		} else {
			medium.FailureCount++
		}
	}
	agentStats["medium"] = medium

	// Low usage agent
	low := &GlobalAgentStats{}
	for i := 0; i < 10; i++ {
		low.UpdateStats(100)
		if i < 9 {
			low.SuccessCount++
		} else {
			low.FailureCount++
		}
	}
	agentStats["low"] = low

	result := formatter.FormatAgentLeaderboard(agentStats, 0)

	// Check sparkline column exists
	if !strings.Contains(result, "Sparkline") {
		t.Error("Table should include Sparkline column")
	}

	// Verify sparklines are generated (should contain sparkline characters)
	lines := strings.Split(result, "\n")
	sparklineFound := false
	for _, line := range lines {
		if strings.Contains(line, "high") {
			// Should have filled sparkline characters for highest count
			if strings.Contains(line, "█") || strings.Contains(line, "▇") || strings.Contains(line, "▆") {
				sparklineFound = true
			}
		}
	}
	if !sparklineFound {
		t.Error("Sparklines should be generated for agent usage visualization")
	}
}

func TestFormatAgentLeaderboardDurationMetrics(t *testing.T) {
	formatter := NewDisplayFormatter()

	// Create test data with different durations
	agentStats := make(map[string]*GlobalAgentStats)

	// Fast agent
	fast := &GlobalAgentStats{}
	for i := 0; i < 10; i++ {
		fast.UpdateStats(100) // 100ms
		fast.SuccessCount++
	}
	agentStats["fast"] = fast

	// Slow agent
	slow := &GlobalAgentStats{}
	for i := 0; i < 10; i++ {
		slow.UpdateStats(500) // 500ms
		slow.SuccessCount++
	}
	agentStats["slow"] = slow

	// Test table format shows duration metrics
	formatter.SetOutputFormat("table")
	result := formatter.FormatAgentLeaderboard(agentStats, 0)
	if !strings.Contains(result, "Avg Duration") {
		t.Error("Table format should include average duration column")
	}

	// Check duration values are present and reasonable
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if strings.Contains(line, "fast") {
			if !strings.Contains(line, "ms") {
				t.Error("Fast agent line should include duration in ms")
			}
		}
	}

	// Test CSV format includes duration metrics
	formatter.SetOutputFormat("csv")
	result = formatter.FormatAgentLeaderboard(agentStats, 0)
	if !strings.Contains(result, "Avg Duration (ms)") {
		t.Error("CSV format should include average duration header")
	}
}

func TestFormatAgentLeaderboardEmptyData(t *testing.T) {
	formatter := NewDisplayFormatter()

	// Test all formats with empty data
	formats := []string{"table", "json", "csv"}

	for _, format := range formats {
		formatter.SetOutputFormat(format)

		// Test with nil
		result := formatter.FormatAgentLeaderboard(nil, 0)
		if result != "" {
			t.Errorf("Format %s should return empty string for nil data", format)
		}

		// Test with empty map
		result = formatter.FormatAgentLeaderboard(map[string]*GlobalAgentStats{}, 0)
		if result != "" {
			t.Errorf("Format %s should return empty string for empty data", format)
		}

		// Test with nil stats in map
		nilStatsMap := map[string]*GlobalAgentStats{
			"agent": nil,
		}
		result = formatter.FormatAgentLeaderboard(nilStatsMap, 0)
		// Should handle gracefully - either empty or show zero values
		if strings.Contains(result, "panic") {
			t.Errorf("Format %s should handle nil stats gracefully", format)
		}
	}
}

func TestFormatAgentLeaderboardLimits(t *testing.T) {
	formatter := NewDisplayFormatter()

	// Create test data with 5 agents
	agentStats := make(map[string]*GlobalAgentStats)
	names := []string{"first", "second", "third", "fourth", "fifth"}
	counts := []int{100, 90, 80, 70, 60} // Decreasing counts for ranking

	for i, name := range names {
		agent := &GlobalAgentStats{}
		for j := 0; j < counts[i]; j++ {
			agent.UpdateStats(int64(100 + i*50)) // Varying durations
			if j < counts[i]-10 {
				agent.SuccessCount++
			} else {
				agent.FailureCount++
			}
		}
		agentStats[name] = agent
	}

	// Test various limits
	testCases := []struct {
		limit    int
		expected int
	}{
		{0, 5},  // No limit
		{3, 3},  // Limit to 3
		{10, 5}, // Limit higher than available
		{1, 1},  // Limit to 1
	}

	for _, tc := range testCases {
		// Test JSON format for easy counting
		formatter.SetOutputFormat("json")
		result := formatter.FormatAgentLeaderboard(agentStats, tc.limit)

		var jsonResult map[string]interface{}
		if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
			t.Errorf("Invalid JSON output for limit %d: %v", tc.limit, err)
			continue
		}

		leaderboard := jsonResult["leaderboard"].([]interface{})
		if len(leaderboard) != tc.expected {
			t.Errorf("Limit %d: expected %d agents, got %d", tc.limit, tc.expected, len(leaderboard))
		}

		// Verify total_agents is always the original count
		if int(jsonResult["total_agents"].(float64)) != 5 {
			t.Errorf("total_agents should always be 5, got %v", jsonResult["total_agents"])
		}

		// Verify showing field reflects actual displayed count
		if int(jsonResult["showing"].(float64)) != tc.expected {
			t.Errorf("showing should be %d, got %v", tc.expected, jsonResult["showing"])
		}
	}
}

func TestFormatAgentLeaderboardRanking(t *testing.T) {
	formatter := NewDisplayFormatter()

	// Create test data with specific counts for ranking verification
	agentStats := make(map[string]*GlobalAgentStats)

	// Top agent (200 invocations)
	top := &GlobalAgentStats{}
	for i := 0; i < 200; i++ {
		top.UpdateStats(100)
		if i < 190 {
			top.SuccessCount++
		} else {
			top.FailureCount++
		}
	}
	agentStats["top"] = top

	// Middle agent (100 invocations)
	middle := &GlobalAgentStats{}
	for i := 0; i < 100; i++ {
		middle.UpdateStats(100)
		if i < 95 {
			middle.SuccessCount++
		} else {
			middle.FailureCount++
		}
	}
	agentStats["middle"] = middle

	// Bottom agent (50 invocations)
	bottom := &GlobalAgentStats{}
	for i := 0; i < 50; i++ {
		bottom.UpdateStats(100)
		if i < 45 {
			bottom.SuccessCount++
		} else {
			bottom.FailureCount++
		}
	}
	agentStats["bottom"] = bottom

	// Test ranking in all formats
	formats := []string{"table", "json", "csv"}

	for _, format := range formats {
		formatter.SetOutputFormat(format)
		result := formatter.FormatAgentLeaderboard(agentStats, 0)

		switch format {
		case "table":
			// Check ranks appear in correct positions
			lines := strings.Split(result, "\n")
			topFound := false
			middleFound := false
			bottomFound := false

			for _, line := range lines {
				if strings.Contains(line, "1") && strings.Contains(line, "top") {
					topFound = true
				}
				if strings.Contains(line, "2") && strings.Contains(line, "middle") {
					middleFound = true
				}
				if strings.Contains(line, "3") && strings.Contains(line, "bottom") {
					bottomFound = true
				}
			}

			if !topFound || !middleFound || !bottomFound {
				t.Errorf("Table format ranking incorrect: top=%v, middle=%v, bottom=%v",
					topFound, middleFound, bottomFound)
			}

		case "json":
			var jsonResult map[string]interface{}
			if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
				t.Errorf("Invalid JSON output: %v", err)
				continue
			}

			leaderboard := jsonResult["leaderboard"].([]interface{})
			if len(leaderboard) != 3 {
				t.Errorf("Expected 3 agents in JSON leaderboard")
				continue
			}

			// Check ranking order
			ranks := []int{1, 2, 3}
			names := []string{"top", "middle", "bottom"}

			for i, item := range leaderboard {
				agent := item.(map[string]interface{})
				if int(agent["rank"].(float64)) != ranks[i] {
					t.Errorf("JSON agent %d should have rank %d, got %v", i, ranks[i], agent["rank"])
				}
				if agent["name"].(string) != names[i] {
					t.Errorf("JSON agent %d should be %s, got %s", i, names[i], agent["name"])
				}
			}

		case "csv":
			lines := strings.Split(result, "\n")
			if len(lines) < 4 { // Header + 3 data rows
				t.Error("CSV should have header + 3 data rows")
				continue
			}

			// Check first data row is top agent
			if !strings.Contains(lines[1], "1,top,200") {
				t.Errorf("CSV first row incorrect: %s", lines[1])
			}
		}
	}
}

func TestFormatAgentLeaderboardTableAlignment(t *testing.T) {
	formatter := NewDisplayFormatter()
	formatter.SetOutputFormat("table")

	// Create test data
	agentStats := make(map[string]*GlobalAgentStats)

	// Short name agent
	short := &GlobalAgentStats{}
	for i := 0; i < 100; i++ {
		short.UpdateStats(100)
		if i < 95 {
			short.SuccessCount++
		} else {
			short.FailureCount++
		}
	}
	agentStats["short"] = short

	// Long name agent
	longName := &GlobalAgentStats{}
	for i := 0; i < 50; i++ {
		longName.UpdateStats(100)
		if i < 45 {
			longName.SuccessCount++
		} else {
			longName.FailureCount++
		}
	}
	agentStats["very-long-agent-name"] = longName

	result := formatter.FormatAgentLeaderboard(agentStats, 0)

	// Verify tabwriter formatting features
	lines := strings.Split(result, "\n")

	// Check header formatting
	headerFound := false
	for _, line := range lines {
		if strings.Contains(line, "AGENT USAGE LEADERBOARD") {
			headerFound = true
			// Should have box drawing characters
			if !strings.Contains(result, "╔") || !strings.Contains(result, "╗") {
				t.Error("Table should use box drawing characters for header")
			}
			break
		}
	}
	if !headerFound {
		t.Error("Table header not found")
	}

	// Check column headers are present
	columnHeadersFound := false
	for _, line := range lines {
		if strings.Contains(line, "Rank") && strings.Contains(line, "Agent") &&
		   strings.Contains(line, "Total Invocations") && strings.Contains(line, "Success Rate") {
			columnHeadersFound = true
			break
		}
	}
	if !columnHeadersFound {
		t.Error("Column headers not found or incomplete")
	}

	// Verify data rows are properly aligned (tabwriter should handle this)
	dataRowCount := 0
	for _, line := range lines {
		if strings.Contains(line, "short") || strings.Contains(line, "very-long-agent-name") {
			dataRowCount++
			// Should contain numeric data
			if !strings.Contains(line, "100") && !strings.Contains(line, "50") {
				t.Errorf("Data row missing numeric data: %s", line)
			}
		}
	}
	if dataRowCount != 2 {
		t.Errorf("Expected 2 data rows, found %d", dataRowCount)
	}
}

func TestFormatAgentLeaderboardJSONStructure(t *testing.T) {
	formatter := NewDisplayFormatter()
	formatter.SetOutputFormat("json")

	// Create comprehensive test data
	agentStats := make(map[string]*GlobalAgentStats)

	// Test agent with varied durations for realistic stats
	testAgent := &GlobalAgentStats{}
	durations := []int64{100, 200, 300, 500, 1000} // Different durations
	for i := 0; i < 100; i++ {
		durationIdx := i % len(durations)
		testAgent.UpdateStats(durations[durationIdx])
		if i < 95 {
			testAgent.SuccessCount++
		} else {
			testAgent.FailureCount++
		}
	}
	agentStats["test-agent"] = testAgent

	result := formatter.FormatAgentLeaderboard(agentStats, 0)

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}

	// Verify top-level structure
	expectedFields := []string{"leaderboard", "total_agents", "showing"}
	for _, field := range expectedFields {
		if _, ok := jsonResult[field]; !ok {
			t.Errorf("JSON missing required field: %s", field)
		}
	}

	// Verify leaderboard structure
	leaderboard, ok := jsonResult["leaderboard"].([]interface{})
	if !ok {
		t.Fatal("leaderboard should be an array")
	}
	if len(leaderboard) != 1 {
		t.Errorf("Expected 1 agent in leaderboard, got %d", len(leaderboard))
	}

	// Verify agent entry structure
	agentEntry := leaderboard[0].(map[string]interface{})
	expectedAgentFields := []string{"rank", "name", "stats"}
	for _, field := range expectedAgentFields {
		if _, ok := agentEntry[field]; !ok {
			t.Errorf("Agent entry missing field: %s", field)
		}
	}

	// Verify stats structure
	statsEntry := agentEntry["stats"].(map[string]interface{})
	expectedStatsFields := []string{"count", "success_count", "failure_count",
		"success_rate", "avg_duration_ms", "total_duration_ms"}
	for _, field := range expectedStatsFields {
		if _, ok := statsEntry[field]; !ok {
			t.Errorf("Stats entry missing field: %s", field)
		}
	}

	// Verify data correctness
	if agentEntry["rank"].(float64) != 1 {
		t.Errorf("Expected rank 1, got %v", agentEntry["rank"])
	}
	if agentEntry["name"].(string) != "test-agent" {
		t.Errorf("Expected name 'test-agent', got %s", agentEntry["name"])
	}
	if statsEntry["count"].(float64) != 100 {
		t.Errorf("Expected count 100, got %v", statsEntry["count"])
	}
	if statsEntry["success_rate"].(float64) != 0.95 {
		t.Errorf("Expected success_rate 0.95, got %v", statsEntry["success_rate"])
	}
}

func TestFormatAgentLeaderboardCSVFormat(t *testing.T) {
	formatter := NewDisplayFormatter()
	formatter.SetOutputFormat("csv")

	// Create test data
	agentStats := make(map[string]*GlobalAgentStats)

	// Agent 1
	agent1 := &GlobalAgentStats{}
	for i := 0; i < 100; i++ {
		agent1.UpdateStats(150)
		if i < 90 {
			agent1.SuccessCount++
		} else {
			agent1.FailureCount++
		}
	}
	agentStats["agent1"] = agent1

	// Agent 2
	agent2 := &GlobalAgentStats{}
	for i := 0; i < 50; i++ {
		agent2.UpdateStats(150)
		if i < 45 {
			agent2.SuccessCount++
		} else {
			agent2.FailureCount++
		}
	}
	agentStats["agent2"] = agent2

	result := formatter.FormatAgentLeaderboard(agentStats, 0)

	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) < 3 { // Header + 2 data rows
		t.Errorf("Expected at least 3 CSV lines, got %d", len(lines))
	}

	// Verify header
	expectedHeaders := []string{"Rank", "Agent", "Total Invocations", "Success Rate", "Avg Duration (ms)"}
	headerFields := strings.Split(lines[0], ",")
	for i, expected := range expectedHeaders {
		if i >= len(headerFields) {
			t.Errorf("Missing CSV header field: %s", expected)
			continue
		}
		if headerFields[i] != expected {
			t.Errorf("CSV header field %d: expected %s, got %s", i, expected, headerFields[i])
		}
	}

	// Verify data rows format
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		fields := strings.Split(lines[i], ",")
		if len(fields) < 5 {
			t.Errorf("CSV data row %d has %d fields, expected at least 5", i, len(fields))
		}

		// Verify rank is numeric
		rank := fields[0]
		if rank != "1" && rank != "2" {
			t.Errorf("CSV rank should be 1 or 2, got %s", rank)
		}

		// Verify agent names are present
		agent := fields[1]
		if agent != "agent1" && agent != "agent2" {
			t.Errorf("CSV agent should be agent1 or agent2, got %s", agent)
		}

		// Verify numeric fields contain numbers
		for j := 2; j < len(fields); j++ {
			if fields[j] == "" {
				t.Errorf("CSV field %d should not be empty in row %d", j, i)
			}
		}
	}
}