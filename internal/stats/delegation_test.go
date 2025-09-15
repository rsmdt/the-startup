package stats

import (
	"testing"
	"time"
)

// Test delegation chain construction per PRD lines 89
func TestDelegationChain_Construction(t *testing.T) {
	tests := []struct {
		name            string
		sessionID       string
		agentSequence   []string
		expectedChain   []string
		expectError     bool
	}{
		{
			name:          "Simple A->B delegation",
			sessionID:     "session-1",
			agentSequence: []string{"the-implementor", "the-tester"},
			expectedChain: []string{"the-implementor", "the-tester"},
			expectError:   false,
		},
		{
			name:          "Complex multi-agent chain",
			sessionID:     "session-2",
			agentSequence: []string{"the-specifier", "the-implementor", "the-tester", "the-documenter"},
			expectedChain: []string{"the-specifier", "the-implementor", "the-tester", "the-documenter"},
			expectError:   false,
		},
		{
			name:          "Single agent no delegation",
			sessionID:     "session-3",
			agentSequence: []string{"the-implementor"},
			expectedChain: []string{"the-implementor"},
			expectError:   false,
		},
		{
			name:          "Empty session",
			sessionID:     "session-4",
			agentSequence: []string{},
			expectedChain: []string{},
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewDelegationAnalyzer()

			// Setup delegation chain by processing transitions
			baseTime := time.Now()
			if len(tt.agentSequence) > 1 {
				for i := 0; i < len(tt.agentSequence)-1; i++ {
					fromAgent := tt.agentSequence[i]
					toAgent := tt.agentSequence[i+1]
					timestamp := baseTime.Add(time.Duration(i) * time.Minute)
					err := analyzer.ProcessAgentTransition(tt.sessionID, fromAgent, toAgent, timestamp)
					if err != nil {
						t.Errorf("Failed to process agent transition: %v", err)
					}
				}
			} else if len(tt.agentSequence) == 1 {
				// Handle single agent case - manually create a session with one agent (no transitions)
				realAnalyzer := analyzer.(*realDelegationAnalyzer)
				realAnalyzer.mu.Lock()
				realAnalyzer.sessionChains[tt.sessionID] = &DelegationChain{
					Agents:    []string{tt.agentSequence[0]},
					SessionID: tt.sessionID,
					StartTime: baseTime,
					EndTime:   baseTime,
				}
				realAnalyzer.mu.Unlock()
			}

			// Test chain retrieval
			result, err := analyzer.GetDelegationChain(tt.sessionID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				if len(tt.expectedChain) > 0 {
					t.Errorf("Expected chain with agents, got nil")
				}
				return
			}

			if len(result.Agents) != len(tt.expectedChain) {
				t.Errorf("Expected chain length %d, got %d", len(tt.expectedChain), len(result.Agents))
				return
			}

			for i, agent := range tt.expectedChain {
				if result.Agents[i] != agent {
					t.Errorf("At position %d: expected %s, got %s", i, agent, result.Agents[i])
				}
			}

			if result.SessionID != tt.sessionID {
				t.Errorf("Expected sessionID %s, got %s", tt.sessionID, result.SessionID)
			}
		})
	}
}

// Test agent co-occurrence tracking per PRD lines 87
func TestAgentCoOccurrence_Tracking(t *testing.T) {
	tests := []struct {
		name                string
		sessionAgentPairs   map[string][]string // sessionID -> agents used
		expectedOccurrences []AgentCoOccurrence
	}{
		{
			name: "Two agents in same session",
			sessionAgentPairs: map[string][]string{
				"session-1": {"the-implementor", "the-tester"},
			},
			expectedOccurrences: []AgentCoOccurrence{
				{
					Agent1:   "the-implementor",
					Agent2:   "the-tester",
					Count:    1,
					Sessions: []string{"session-1"},
				},
			},
		},
		{
			name: "Multiple sessions with same agent pair",
			sessionAgentPairs: map[string][]string{
				"session-1": {"the-implementor", "the-tester"},
				"session-2": {"the-implementor", "the-tester"},
				"session-3": {"the-implementor", "the-tester"},
			},
			expectedOccurrences: []AgentCoOccurrence{
				{
					Agent1:   "the-implementor",
					Agent2:   "the-tester",
					Count:    3,
					Sessions: []string{"session-1", "session-2", "session-3"},
				},
			},
		},
		{
			name: "Three agents co-occurrence combinations",
			sessionAgentPairs: map[string][]string{
				"session-1": {"the-specifier", "the-implementor", "the-tester"},
			},
			expectedOccurrences: []AgentCoOccurrence{
				{Agent1: "the-implementor", Agent2: "the-specifier", Count: 1, Sessions: []string{"session-1"}},
				{Agent1: "the-implementor", Agent2: "the-tester", Count: 1, Sessions: []string{"session-1"}},
				{Agent1: "the-specifier", Agent2: "the-tester", Count: 1, Sessions: []string{"session-1"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewDelegationAnalyzer()

			// Process agent sessions to generate co-occurrences
			baseTime := time.Now()
			for sessionID, agents := range tt.sessionAgentPairs {
				// Use ProcessSessionAgents to handle co-occurrence analysis
				realAnalyzer := analyzer.(*realDelegationAnalyzer)
				err := realAnalyzer.ProcessSessionAgents(sessionID, agents, baseTime)
				if err != nil {
					t.Errorf("Failed to process session agents: %v", err)
				}
			}

			// Test co-occurrence retrieval
			result, err := analyzer.GetCoOccurrencePatterns()
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expectedOccurrences) {
				t.Errorf("Expected %d co-occurrences, got %d", len(tt.expectedOccurrences), len(result))
				return
			}

			// Verify each expected co-occurrence (order-independent)
			for _, expected := range tt.expectedOccurrences {
				found := false
				for _, actual := range result {
					if actual.Agent1 == expected.Agent1 && actual.Agent2 == expected.Agent2 {
						found = true
						if actual.Count != expected.Count {
							t.Errorf("Co-occurrence %s-%s: expected Count %d, got %d",
								expected.Agent1, expected.Agent2, expected.Count, actual.Count)
						}
						if len(actual.Sessions) != len(expected.Sessions) {
							t.Errorf("Co-occurrence %s-%s: expected %d sessions, got %d",
								expected.Agent1, expected.Agent2, len(expected.Sessions), len(actual.Sessions))
						}
						// Verify sessions match (order-independent)
						sessionMap := make(map[string]bool)
						for _, s := range actual.Sessions {
							sessionMap[s] = true
						}
						for _, expectedSession := range expected.Sessions {
							if !sessionMap[expectedSession] {
								t.Errorf("Co-occurrence %s-%s: missing expected session %s",
									expected.Agent1, expected.Agent2, expectedSession)
							}
						}
						break
					}
				}
				if !found {
					t.Errorf("Expected co-occurrence %s-%s not found", expected.Agent1, expected.Agent2)
				}
			}
		})
	}
}

// Test session boundary detection per PRD lines 90
func TestSessionBoundary_Detection(t *testing.T) {
	tests := []struct {
		name           string
		sessionID      string
		boundaryEvents []struct {
			eventType string
			timestamp time.Time
		}
		expectedBoundaries int
	}{
		{
			name:      "Session start and end",
			sessionID: "session-1",
			boundaryEvents: []struct {
				eventType string
				timestamp time.Time
			}{
				{"start", time.Now()},
				{"end", time.Now().Add(time.Hour)},
			},
			expectedBoundaries: 2,
		},
		{
			name:      "Multiple session boundaries",
			sessionID: "session-2",
			boundaryEvents: []struct {
				eventType string
				timestamp time.Time
			}{
				{"start", time.Now()},
				{"pause", time.Now().Add(30 * time.Minute)},
				{"resume", time.Now().Add(45 * time.Minute)},
				{"end", time.Now().Add(time.Hour)},
			},
			expectedBoundaries: 4,
		},
		{
			name:               "No boundaries",
			sessionID:          "session-3",
			boundaryEvents:     []struct{ eventType string; timestamp time.Time }{},
			expectedBoundaries: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewDelegationAnalyzer()

			// Process boundary events
			for _, event := range tt.boundaryEvents {
				err := analyzer.ProcessSessionBoundary(tt.sessionID, event.eventType, event.timestamp)
				if err != nil {
					t.Errorf("Unexpected error processing boundary: %v", err)
				}
			}

			// Verify boundaries were recorded using the public API
			realAnalyzer := analyzer.(*realDelegationAnalyzer)
			boundaries, err := realAnalyzer.GetSessionBoundaries(tt.sessionID)
			if err != nil {
				t.Errorf("Unexpected error getting boundaries: %v", err)
				return
			}

			if len(boundaries) != tt.expectedBoundaries {
				t.Errorf("Expected %d boundaries, got %d", tt.expectedBoundaries, len(boundaries))
			}

			// Verify boundary timestamps match
			for i, event := range tt.boundaryEvents {
				if i >= len(boundaries) {
					t.Errorf("Missing boundary at index %d", i)
					continue
				}

				if !boundaries[i].Equal(event.timestamp) {
					t.Errorf("Boundary %d: expected timestamp %v, got %v", i, event.timestamp, boundaries[i])
				}
			}
		})
	}
}

// Test circular delegation handling (edge case)
func TestCircularDelegation_Handling(t *testing.T) {
	tests := []struct {
		name                 string
		sessionID            string
		transitions          []struct {
			from string
			to   string
		}
		expectCircular       bool
		expectedPatternCount int
	}{
		{
			name:      "Simple A->B->A circular delegation",
			sessionID: "session-circular-1",
			transitions: []struct{ from, to string }{
				{"the-implementor", "the-tester"},
				{"the-tester", "the-implementor"},
			},
			expectCircular:       true,
			expectedPatternCount: 2,
		},
		{
			name:      "Complex A->B->C->A circular delegation",
			sessionID: "session-circular-2",
			transitions: []struct{ from, to string }{
				{"the-specifier", "the-implementor"},
				{"the-implementor", "the-tester"},
				{"the-tester", "the-specifier"},
			},
			expectCircular:       true,
			expectedPatternCount: 3,
		},
		{
			name:      "Linear delegation (no cycle)",
			sessionID: "session-linear-1",
			transitions: []struct{ from, to string }{
				{"the-specifier", "the-implementor"},
				{"the-implementor", "the-tester"},
				{"the-tester", "the-documenter"},
			},
			expectCircular:       false,
			expectedPatternCount: 3,
		},
		{
			name:      "Self-delegation edge case",
			sessionID: "session-self-1",
			transitions: []struct{ from, to string }{
				{"the-implementor", "the-implementor"},
			},
			expectCircular:       true,
			expectedPatternCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewDelegationAnalyzer()

			// Process transitions
			baseTime := time.Now()
			for i, transition := range tt.transitions {
				timestamp := baseTime.Add(time.Duration(i) * time.Minute)
				err := analyzer.ProcessAgentTransition(tt.sessionID, transition.from, transition.to, timestamp)
				if err != nil {
					t.Errorf("Unexpected error processing transition: %v", err)
				}
			}

			// Get transition patterns
			patterns, err := analyzer.GetTransitionPatterns()
			if err != nil {
				t.Errorf("Unexpected error getting patterns: %v", err)
				return
			}

			if len(patterns) != tt.expectedPatternCount {
				t.Errorf("Expected %d transition patterns, got %d", tt.expectedPatternCount, len(patterns))
			}

			// Use the built-in circular detection method
			realAnalyzer := analyzer.(*realDelegationAnalyzer)
			hasCircular, err := realAnalyzer.HasCircularDelegation(tt.sessionID)
			if err != nil {
				t.Errorf("Unexpected error checking circular delegation: %v", err)
				return
			}

			if hasCircular != tt.expectCircular {
				t.Errorf("Expected circular: %v, got: %v", tt.expectCircular, hasCircular)
			}
		})
	}
}

// Test agent transition patterns (A → B flows)
func TestAgentTransition_Patterns(t *testing.T) {
	tests := []struct {
		name            string
		transitions     []struct {
			sessionID string
			from      string
			to        string
			timestamp time.Time
		}
		expectedPatterns []DelegationPattern
	}{
		{
			name: "Single transition pattern",
			transitions: []struct {
				sessionID string
				from      string
				to        string
				timestamp time.Time
			}{
				{"session-1", "the-implementor", "the-tester", time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)},
			},
			expectedPatterns: []DelegationPattern{
				{
					SourceAgent: "the-implementor",
					TargetAgent: "the-tester",
					Count:       1,
					LastSeen:    time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
					SessionID:   "session-1",
				},
			},
		},
		{
			name: "Multiple same transitions (frequency counting)",
			transitions: []struct {
				sessionID string
				from      string
				to        string
				timestamp time.Time
			}{
				{"session-1", "the-implementor", "the-tester", time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)},
				{"session-2", "the-implementor", "the-tester", time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)},
				{"session-3", "the-implementor", "the-tester", time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
			},
			expectedPatterns: []DelegationPattern{
				{
					SourceAgent: "the-implementor",
					TargetAgent: "the-tester",
					Count:       3,
					LastSeen:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					SessionID:   "session-3", // Last session ID
				},
			},
		},
		{
			name: "Multiple different transition patterns",
			transitions: []struct {
				sessionID string
				from      string
				to        string
				timestamp time.Time
			}{
				{"session-1", "the-specifier", "the-implementor", time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)},
				{"session-1", "the-implementor", "the-tester", time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC)},
				{"session-2", "the-tester", "the-documenter", time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)},
			},
			expectedPatterns: []DelegationPattern{
				{SourceAgent: "the-specifier", TargetAgent: "the-implementor", Count: 1, LastSeen: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC), SessionID: "session-1"},
				{SourceAgent: "the-implementor", TargetAgent: "the-tester", Count: 1, LastSeen: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC), SessionID: "session-1"},
				{SourceAgent: "the-tester", TargetAgent: "the-documenter", Count: 1, LastSeen: time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC), SessionID: "session-2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewDelegationAnalyzer()

			// Process all transitions
			for _, transition := range tt.transitions {
				err := analyzer.ProcessAgentTransition(transition.sessionID, transition.from, transition.to, transition.timestamp)
				if err != nil {
					t.Errorf("Unexpected error processing transition: %v", err)
				}
			}

			// Get transition patterns
			patterns, err := analyzer.GetTransitionPatterns()
			if err != nil {
				t.Errorf("Unexpected error getting patterns: %v", err)
				return
			}

			if len(patterns) != len(tt.expectedPatterns) {
				t.Errorf("Expected %d patterns, got %d", len(tt.expectedPatterns), len(patterns))
				return
			}

			// Verify each pattern (order-independent)
			for _, expected := range tt.expectedPatterns {
				found := false
				for _, actual := range patterns {
					if actual.SourceAgent == expected.SourceAgent && actual.TargetAgent == expected.TargetAgent {
						found = true
						if actual.Count != expected.Count {
							t.Errorf("Pattern %s->%s: expected Count %d, got %d",
								expected.SourceAgent, expected.TargetAgent, expected.Count, actual.Count)
						}
						if !actual.LastSeen.Equal(expected.LastSeen) {
							t.Errorf("Pattern %s->%s: expected LastSeen %v, got %v",
								expected.SourceAgent, expected.TargetAgent, expected.LastSeen, actual.LastSeen)
						}
						if actual.SessionID != expected.SessionID {
							t.Errorf("Pattern %s->%s: expected SessionID %s, got %s",
								expected.SourceAgent, expected.TargetAgent, expected.SessionID, actual.SessionID)
						}
						break
					}
				}
				if !found {
					t.Errorf("Expected pattern %s->%s not found", expected.SourceAgent, expected.TargetAgent)
				}
			}
		})
	}
}

// Test time period filtering integration
func TestTimePeriodFiltering_Integration(t *testing.T) {
	tests := []struct {
		name            string
		transitions     []struct {
			sessionID string
			from      string
			to        string
			timestamp time.Time
		}
		filterStart     time.Time
		filterEnd       time.Time
		expectedCount   int
	}{
		{
			name: "Filter includes all transitions",
			transitions: []struct {
				sessionID string
				from      string
				to        string
				timestamp time.Time
			}{
				{"session-1", "the-implementor", "the-tester", time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)},
				{"session-2", "the-specifier", "the-implementor", time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)},
			},
			filterStart:   time.Date(2023, 1, 1, 9, 0, 0, 0, time.UTC),
			filterEnd:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedCount: 2,
		},
		{
			name: "Filter excludes early transitions",
			transitions: []struct {
				sessionID string
				from      string
				to        string
				timestamp time.Time
			}{
				{"session-1", "the-implementor", "the-tester", time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC)},  // Before filter
				{"session-2", "the-specifier", "the-implementor", time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)}, // Within filter
				{"session-3", "the-tester", "the-documenter", time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)}, // Within filter
			},
			filterStart:   time.Date(2023, 1, 1, 9, 0, 0, 0, time.UTC),
			filterEnd:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedCount: 2,
		},
		{
			name: "Filter excludes late transitions",
			transitions: []struct {
				sessionID string
				from      string
				to        string
				timestamp time.Time
			}{
				{"session-1", "the-implementor", "the-tester", time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)}, // Within filter
				{"session-2", "the-specifier", "the-implementor", time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)}, // Within filter
				{"session-3", "the-tester", "the-documenter", time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC)}, // After filter
			},
			filterStart:   time.Date(2023, 1, 1, 9, 0, 0, 0, time.UTC),
			filterEnd:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedCount: 2,
		},
		{
			name: "Empty time period excludes all",
			transitions: []struct {
				sessionID string
				from      string
				to        string
				timestamp time.Time
			}{
				{"session-1", "the-implementor", "the-tester", time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)},
			},
			filterStart:   time.Date(2023, 1, 2, 9, 0, 0, 0, time.UTC),
			filterEnd:     time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewDelegationAnalyzer()

			// Process all transitions
			for _, transition := range tt.transitions {
				err := analyzer.ProcessAgentTransition(transition.sessionID, transition.from, transition.to, transition.timestamp)
				if err != nil {
					t.Errorf("Unexpected error processing transition: %v", err)
				}
			}

			// Apply time filtering
			filtered, err := analyzer.GetFilteredPatterns(tt.filterStart, tt.filterEnd)
			if err != nil {
				t.Errorf("Unexpected error filtering patterns: %v", err)
				return
			}

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d filtered patterns, got %d", tt.expectedCount, len(filtered))
			}

			// Verify all filtered patterns are within time range
			for i, pattern := range filtered {
				if pattern.LastSeen.Before(tt.filterStart) || pattern.LastSeen.After(tt.filterEnd) {
					t.Errorf("Pattern %d timestamp %v is outside filter range [%v, %v]",
						i, pattern.LastSeen, tt.filterStart, tt.filterEnd)
				}
			}
		})
	}
}