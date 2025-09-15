package stats

import (
	"fmt"
	"math"
	"testing"
	"time"
)

// Types are defined in agent_types.go

// TestAgentInvocation_Initialization tests that AgentInvocation initializes correctly
func TestAgentInvocation_Initialization(t *testing.T) {
	now := time.Now()
	inv := AgentInvocation{
		AgentType:       "test-agent",
		SessionID:       "session-123",
		Timestamp:       now,
		DurationMs:      1500,
		Success:         true,
		Confidence:      0.95,
		DetectionMethod: "pattern-match",
	}

	if inv.AgentType != "test-agent" {
		t.Errorf("Expected AgentType 'test-agent', got %s", inv.AgentType)
	}

	if inv.SessionID != "session-123" {
		t.Errorf("Expected SessionID 'session-123', got %s", inv.SessionID)
	}

	if !inv.Timestamp.Equal(now) {
		t.Errorf("Expected Timestamp %v, got %v", now, inv.Timestamp)
	}

	if inv.DurationMs != 1500 {
		t.Errorf("Expected DurationMs 1500, got %d", inv.DurationMs)
	}

	if !inv.Success {
		t.Error("Expected Success to be true")
	}

	if inv.Confidence != 0.95 {
		t.Errorf("Expected Confidence 0.95, got %f", inv.Confidence)
	}

	if inv.DetectionMethod != "pattern-match" {
		t.Errorf("Expected DetectionMethod 'pattern-match', got %s", inv.DetectionMethod)
	}
}

// TestAgentInvocation_ZeroValues tests zero value initialization
func TestAgentInvocation_ZeroValues(t *testing.T) {
	var inv AgentInvocation

	if inv.AgentType != "" {
		t.Errorf("Expected empty AgentType, got %s", inv.AgentType)
	}

	if inv.SessionID != "" {
		t.Errorf("Expected empty SessionID, got %s", inv.SessionID)
	}

	if !inv.Timestamp.IsZero() {
		t.Errorf("Expected zero Timestamp, got %v", inv.Timestamp)
	}

	if inv.DurationMs != 0 {
		t.Errorf("Expected DurationMs 0, got %d", inv.DurationMs)
	}

	if inv.Success {
		t.Error("Expected Success to be false")
	}

	if inv.Confidence != 0.0 {
		t.Errorf("Expected Confidence 0.0, got %f", inv.Confidence)
	}

	if inv.DetectionMethod != "" {
		t.Errorf("Expected empty DetectionMethod, got %s", inv.DetectionMethod)
	}
}

// TestGlobalAgentStats_InitialState tests initial state of GlobalAgentStats
func TestGlobalAgentStats_InitialState(t *testing.T) {
	stats := &GlobalAgentStats{}

	if stats.Count != 0 {
		t.Errorf("Expected Count 0, got %d", stats.Count)
	}

	if stats.SuccessCount != 0 {
		t.Errorf("Expected SuccessCount 0, got %d", stats.SuccessCount)
	}

	if stats.FailureCount != 0 {
		t.Errorf("Expected FailureCount 0, got %d", stats.FailureCount)
	}

	if stats.TotalDurationMs != 0 {
		t.Errorf("Expected TotalDurationMs 0, got %d", stats.TotalDurationMs)
	}

	if stats.Mean() != 0.0 {
		t.Errorf("Expected mean 0.0, got %f", stats.Mean())
	}

	if stats.Variance() != 0.0 {
		t.Errorf("Expected Variance 0.0, got %f", stats.Variance())
	}

	if stats.StdDev() != 0.0 {
		t.Errorf("Expected StdDev 0.0, got %f", stats.StdDev())
	}
}

// TestGlobalAgentStats_WelfordAlgorithm_SingleValue tests Welford's algorithm with one value
func TestGlobalAgentStats_WelfordAlgorithm_SingleValue(t *testing.T) {
	stats := &GlobalAgentStats{}

	stats.UpdateStats(100)

	if stats.Count != 1 {
		t.Errorf("Expected Count 1, got %d", stats.Count)
	}

	if stats.TotalDurationMs != 100 {
		t.Errorf("Expected TotalDurationMs 100, got %d", stats.TotalDurationMs)
	}

	if stats.Mean() != 100.0 {
		t.Errorf("Expected Mean 100.0, got %f", stats.Mean())
	}

	// Variance is undefined for single value
	if stats.Variance() != 0.0 {
		t.Errorf("Expected Variance 0.0 for single value, got %f", stats.Variance())
	}

	if stats.StdDev() != 0.0 {
		t.Errorf("Expected StdDev 0.0 for single value, got %f", stats.StdDev())
	}
}

// TestGlobalAgentStats_WelfordAlgorithm_MultipleValues tests Welford's algorithm accuracy
func TestGlobalAgentStats_WelfordAlgorithm_MultipleValues(t *testing.T) {
	stats := &GlobalAgentStats{}

	// Test data with known mean and variance
	values := []int64{10, 20, 30, 40, 50}
	expectedMean := 30.0
	expectedVariance := 250.0 // Sample variance

	for _, v := range values {
		stats.UpdateStats(v)
	}

	if stats.Count != 5 {
		t.Errorf("Expected Count 5, got %d", stats.Count)
	}

	if stats.TotalDurationMs != 150 {
		t.Errorf("Expected TotalDurationMs 150, got %d", stats.TotalDurationMs)
	}

	if math.Abs(stats.Mean()-expectedMean) > 0.001 {
		t.Errorf("Expected Mean %f, got %f", expectedMean, stats.Mean())
	}

	if math.Abs(stats.Variance()-expectedVariance) > 0.001 {
		t.Errorf("Expected Variance %f, got %f", expectedVariance, stats.Variance())
	}

	expectedStdDev := math.Sqrt(expectedVariance)
	if math.Abs(stats.StdDev()-expectedStdDev) > 0.001 {
		t.Errorf("Expected StdDev %f, got %f", expectedStdDev, stats.StdDev())
	}
}

// TestGlobalAgentStats_WelfordAlgorithm_LargeDataset tests Welford's algorithm with many values
func TestGlobalAgentStats_WelfordAlgorithm_LargeDataset(t *testing.T) {
	stats := &GlobalAgentStats{}

	// Generate 1000 values with known distribution
	n := 1000
	sum := int64(0)
	sumSquares := float64(0)

	for i := 0; i < n; i++ {
		value := int64(i + 1)
		stats.UpdateStats(value)
		sum += value
		sumSquares += float64(value * value)
	}

	expectedMean := float64(sum) / float64(n)
	expectedVariance := (sumSquares/float64(n) - expectedMean*expectedMean) * float64(n) / float64(n-1)

	if stats.Count != n {
		t.Errorf("Expected Count %d, got %d", n, stats.Count)
	}

	if stats.TotalDurationMs != sum {
		t.Errorf("Expected TotalDurationMs %d, got %d", sum, stats.TotalDurationMs)
	}

	if math.Abs(stats.Mean()-expectedMean) > 0.1 {
		t.Errorf("Expected Mean %f, got %f", expectedMean, stats.Mean())
	}

	if math.Abs(stats.Variance()-expectedVariance) > 1.0 {
		t.Errorf("Expected Variance %f, got %f", expectedVariance, stats.Variance())
	}
}

// TestGlobalAgentStats_WelfordAlgorithm_NegativeValues tests with negative durations (edge case)
func TestGlobalAgentStats_WelfordAlgorithm_NegativeValues(t *testing.T) {
	stats := &GlobalAgentStats{}

	// This shouldn't happen in practice, but test edge case
	stats.UpdateStats(-100)
	stats.UpdateStats(200)
	stats.UpdateStats(-50)

	if stats.Count != 3 {
		t.Errorf("Expected Count 3, got %d", stats.Count)
	}

	if stats.TotalDurationMs != 50 {
		t.Errorf("Expected TotalDurationMs 50, got %d", stats.TotalDurationMs)
	}

	expectedMean := 50.0 / 3.0
	if math.Abs(stats.Mean()-expectedMean) > 0.001 {
		t.Errorf("Expected Mean %f, got %f", expectedMean, stats.Mean())
	}
}

// TestConfidenceScore_ValidRange tests confidence score stays within 0.0-1.0
func TestConfidenceScore_ValidRange(t *testing.T) {
	testCases := []struct {
		name       string
		confidence float64
		valid      bool
	}{
		{"Zero confidence", 0.0, true},
		{"Low confidence", 0.25, true},
		{"Medium confidence", 0.5, true},
		{"High confidence", 0.75, true},
		{"Full confidence", 1.0, true},
		{"Negative confidence", -0.1, false},
		{"Over confidence", 1.1, false},
		{"Large negative", -100.0, false},
		{"Large positive", 100.0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid := isValidConfidence(tc.confidence)
			if valid != tc.valid {
				t.Errorf("Expected confidence %f to be valid=%v, got %v",
					tc.confidence, tc.valid, valid)
			}
		})
	}
}

// TestConfidenceScore_EdgeCases tests edge cases for confidence scores
func TestConfidenceScore_EdgeCases(t *testing.T) {
	testCases := []struct {
		name       string
		confidence float64
		valid      bool
	}{
		{"Smallest positive", 0.000001, true},
		{"Almost one", 0.999999, true},
		{"Smallest negative", -0.000001, false},
		{"Just over one", 1.000001, false},
		{"NaN", math.NaN(), false},
		{"Positive infinity", math.Inf(1), false},
		{"Negative infinity", math.Inf(-1), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid := isValidConfidence(tc.confidence)
			if valid != tc.valid {
				t.Errorf("Expected confidence %f to be valid=%v, got %v",
					tc.confidence, tc.valid, valid)
			}
		})
	}
}

// TestConfidenceScore_InAgentInvocation tests confidence validation in context
func TestConfidenceScore_InAgentInvocation(t *testing.T) {
	testCases := []struct {
		name       string
		confidence float64
		shouldFail bool
	}{
		{"Direct Task detection", 1.0, false},
		{"Pattern match high", 0.95, false},
		{"Pattern match medium", 0.90, false},
		{"Pattern match low", 0.75, false},
		{"Command context", 0.70, false},
		{"Uncertain detection", 0.60, false},
		{"No detection", 0.0, false},
		{"Invalid negative", -0.5, true},
		{"Invalid over one", 1.5, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inv := AgentInvocation{
				AgentType:  "test-agent",
				SessionID:  "session-1",
				Timestamp:  time.Now(),
				DurationMs: 100,
				Success:    true,
				Confidence: tc.confidence,
			}

			err := validateAgentInvocation(inv)
			if tc.shouldFail && err == nil {
				t.Errorf("Expected validation to fail for confidence %f", tc.confidence)
			} else if !tc.shouldFail && err != nil {
				t.Errorf("Expected validation to pass for confidence %f, got error: %v",
					tc.confidence, err)
			}
		})
	}
}

// AgentAnalyzer interface is defined in agent_aggregator.go

// MockAgentAnalyzer is just an alias for the actual implementation in tests
type MockAgentAnalyzer = agentAggregator

// NewMockAgentAnalyzer creates a new analyzer (using actual implementation)
func NewMockAgentAnalyzer() AgentAnalyzer {
	return NewAgentAnalyzer()
}

// TestAgentAnalyzer_InterfaceCompliance tests that MockAgentAnalyzer implements the interface
func TestAgentAnalyzer_InterfaceCompliance(t *testing.T) {
	analyzer := NewMockAgentAnalyzer()

	// Verify it implements the interface
	var _ AgentAnalyzer = analyzer

	// Verify it also implements base Aggregator
	var _ Aggregator = analyzer
}

// TestAgentAnalyzer_ProcessAgentInvocation tests processing invocations
func TestAgentAnalyzer_ProcessAgentInvocation(t *testing.T) {
	analyzer := NewMockAgentAnalyzer()

	inv := AgentInvocation{
		AgentType:       "test-agent",
		SessionID:       "session-1",
		Timestamp:       time.Now(),
		DurationMs:      150,
		Success:         true,
		Confidence:      0.95,
		DetectionMethod: "pattern",
	}

	err := analyzer.ProcessAgentInvocation(inv)
	if err != nil {
		t.Fatalf("ProcessAgentInvocation failed: %v", err)
	}

	stats, err := analyzer.GetAgentStats("test-agent")
	if err != nil {
		t.Fatalf("GetAgentStats failed: %v", err)
	}

	if stats.Count != 1 {
		t.Errorf("Expected Count 1, got %d", stats.Count)
	}

	if stats.SuccessCount != 1 {
		t.Errorf("Expected SuccessCount 1, got %d", stats.SuccessCount)
	}

	if stats.TotalDurationMs != 150 {
		t.Errorf("Expected TotalDurationMs 150, got %d", stats.TotalDurationMs)
	}

	if stats.Mean() != 150.0 {
		t.Errorf("Expected Mean 150.0, got %f", stats.Mean())
	}
}

// TestAgentAnalyzer_ProcessAgentInvocation_InvalidConfidence tests invalid confidence handling
func TestAgentAnalyzer_ProcessAgentInvocation_InvalidConfidence(t *testing.T) {
	analyzer := NewMockAgentAnalyzer()

	inv := AgentInvocation{
		AgentType:  "test-agent",
		SessionID:  "session-1",
		Timestamp:  time.Now(),
		DurationMs: 100,
		Success:    true,
		Confidence: 1.5, // Invalid
	}

	err := analyzer.ProcessAgentInvocation(inv)
	if err == nil {
		t.Error("Expected error for invalid confidence, got nil")
	}
}

// TestAgentAnalyzer_GetAgentStats_NotFound tests getting stats for non-existent agent
func TestAgentAnalyzer_GetAgentStats_NotFound(t *testing.T) {
	analyzer := NewMockAgentAnalyzer()

	_, err := analyzer.GetAgentStats("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent agent, got nil")
	}
}

// TestAgentAnalyzer_GetAllAgentStats tests getting all agent stats
func TestAgentAnalyzer_GetAllAgentStats(t *testing.T) {
	analyzer := NewMockAgentAnalyzer()

	// Add multiple agents
	agents := []string{"agent-1", "agent-2", "agent-3"}
	for i, agent := range agents {
		inv := AgentInvocation{
			AgentType:  agent,
			SessionID:  "session-1",
			Timestamp:  time.Now(),
			DurationMs: int64((i + 1) * 100),
			Success:    true,
			Confidence: 0.95,
		}
		err := analyzer.ProcessAgentInvocation(inv)
		if err != nil {
			t.Fatalf("Failed to process invocation for %s: %v", agent, err)
		}
	}

	allStats, err := analyzer.GetAllAgentStats()
	if err != nil {
		t.Fatalf("GetAllAgentStats failed: %v", err)
	}

	if len(allStats) != 3 {
		t.Errorf("Expected 3 agent stats, got %d", len(allStats))
	}

	for _, agent := range agents {
		if _, exists := allStats[agent]; !exists {
			t.Errorf("Missing stats for agent %s", agent)
		}
	}
}

// TestAgentAnalyzer_MultipleInvocations tests multiple invocations for same agent
func TestAgentAnalyzer_MultipleInvocations(t *testing.T) {
	analyzer := NewMockAgentAnalyzer()

	// Process multiple invocations for same agent
	for i := 0; i < 5; i++ {
		inv := AgentInvocation{
			AgentType:  "test-agent",
			SessionID:  fmt.Sprintf("session-%d", i),
			Timestamp:  time.Now(),
			DurationMs: int64((i + 1) * 100),
			Success:    i%2 == 0, // Alternate success/failure
			Confidence: 0.85,
		}
		err := analyzer.ProcessAgentInvocation(inv)
		if err != nil {
			t.Fatalf("Failed to process invocation %d: %v", i, err)
		}
	}

	stats, err := analyzer.GetAgentStats("test-agent")
	if err != nil {
		t.Fatalf("GetAgentStats failed: %v", err)
	}

	if stats.Count != 5 {
		t.Errorf("Expected Count 5, got %d", stats.Count)
	}

	if stats.SuccessCount != 3 {
		t.Errorf("Expected SuccessCount 3, got %d", stats.SuccessCount)
	}

	if stats.FailureCount != 2 {
		t.Errorf("Expected FailureCount 2, got %d", stats.FailureCount)
	}

	expectedTotal := int64(100 + 200 + 300 + 400 + 500)
	if stats.TotalDurationMs != expectedTotal {
		t.Errorf("Expected TotalDurationMs %d, got %d", expectedTotal, stats.TotalDurationMs)
	}

	expectedMean := float64(expectedTotal) / 5.0
	if math.Abs(stats.Mean()-expectedMean) > 0.001 {
		t.Errorf("Expected Mean %f, got %f", expectedMean, stats.Mean())
	}
}

// Helper functions are defined in agent_types.go and agent_aggregator.go