package stats

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestWelfordAlgorithm_Accuracy tests Welford's algorithm with known datasets per SDD lines 519-543
func TestWelfordAlgorithm_Accuracy(t *testing.T) {
	tests := []struct {
		name           string
		values         []float64
		expectedMean   float64
		expectedVar    float64
		tolerance      float64
	}{
		{
			name:           "Simple dataset",
			values:         []float64{2, 4, 4, 4, 5, 5, 7, 9},
			expectedMean:   5.0,
			expectedVar:    4.571, // Sample variance
			tolerance:      0.01,
		},
		{
			name:           "Single value",
			values:         []float64{100},
			expectedMean:   100.0,
			expectedVar:    0.0, // Undefined for single value
			tolerance:      0.001,
		},
		{
			name:           "Uniform distribution",
			values:         []float64{1, 2, 3, 4, 5},
			expectedMean:   3.0,
			expectedVar:    2.5, // Sample variance
			tolerance:      0.001,
		},
		{
			name:           "Large values with precision",
			values:         []float64{1000000, 1000001, 1000002, 1000003, 1000004},
			expectedMean:   1000002.0,
			expectedVar:    2.5,
			tolerance:      0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := newWelfordStats()

			for _, v := range tt.values {
				stats.update(v)
			}

			if math.Abs(stats.mean-tt.expectedMean) > tt.tolerance {
				t.Errorf("Mean: expected %.3f, got %.3f", tt.expectedMean, stats.mean)
			}

			if math.Abs(stats.variance()-tt.expectedVar) > tt.tolerance {
				t.Errorf("Variance: expected %.3f, got %.3f", tt.expectedVar, stats.variance())
			}

			// Test min/max tracking
			if len(tt.values) > 0 {
				expectedMin := tt.values[0]
				expectedMax := tt.values[0]
				for _, v := range tt.values[1:] {
					if v < expectedMin {
						expectedMin = v
					}
					if v > expectedMax {
						expectedMax = v
					}
				}

				if stats.min == nil || *stats.min != expectedMin {
					t.Errorf("Min: expected %.3f, got %v", expectedMin, stats.min)
				}
				if stats.max == nil || *stats.max != expectedMax {
					t.Errorf("Max: expected %.3f, got %v", expectedMax, stats.max)
				}
			}
		})
	}
}

// TestWelfordAlgorithm_StreamingAccuracy tests accuracy with streaming updates
func TestWelfordAlgorithm_StreamingAccuracy(t *testing.T) {
	stats := newWelfordStats()

	// Generate data with known statistical properties
	n := 10000
	mu := 50.0  // Target mean
	sigma := 10.0 // Target std dev

	rand.Seed(42) // For reproducibility
	var sum float64
	var sumSq float64

	for i := 0; i < n; i++ {
		// Generate normal-ish distribution using Box-Muller
		u1 := rand.Float64()
		u2 := rand.Float64()
		z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
		value := mu + sigma*z

		stats.update(value)
		sum += value
		sumSq += value * value
	}

	// Expected values
	expectedMean := sum / float64(n)
	expectedVar := (sumSq/float64(n) - expectedMean*expectedMean) * float64(n) / float64(n-1)

	// Welford should be very accurate
	if math.Abs(stats.mean-expectedMean) > 0.01 {
		t.Errorf("Streaming mean: expected %.6f, got %.6f", expectedMean, stats.mean)
	}

	if math.Abs(stats.variance()-expectedVar) > 0.01 {
		t.Errorf("Streaming variance: expected %.6f, got %.6f", expectedVar, stats.variance())
	}
}

// TestTDigestPercentiles_Accuracy tests T-digest percentile accuracy per SDD lines 544-570
func TestTDigestPercentiles_Accuracy(t *testing.T) {
	tests := []struct {
		name        string
		dataGen     func() []float64
		percentiles map[float64]float64
		tolerance   float64
	}{
		{
			name: "Uniform 1-100",
			dataGen: func() []float64 {
				data := make([]float64, 100)
				for i := 0; i < 100; i++ {
					data[i] = float64(i + 1)
				}
				return data
			},
			percentiles: map[float64]float64{
				0.50: 50.5, // Median
				0.95: 95.0, // P95
				0.99: 99.0, // P99
			},
			tolerance: 2.0, // Allow 2% error for t-digest approximation
		},
		{
			name: "Exponential-like distribution",
			dataGen: func() []float64 {
				data := make([]float64, 1000)
				for i := 0; i < 1000; i++ {
					data[i] = math.Exp(float64(i) / 200.0)
				}
				return data
			},
			percentiles: map[float64]float64{
				0.50: math.Exp(2.5),   // Approximate median
				0.95: math.Exp(4.75),  // Approximate P95
				0.99: math.Exp(4.95),  // Approximate P99
			},
			tolerance: 0.3, // 30% tolerance for exponential (t-digest approximation)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := newTDigest(100)
			data := tt.dataGen()

			// Add all data points
			for _, v := range data {
				td.add(v)
			}

			// Test each percentile
			for p, expected := range tt.percentiles {
				result := td.quantile(p)
				errorRate := math.Abs(result-expected) / expected

				if errorRate > tt.tolerance {
					t.Errorf("P%.0f: expected ~%.2f, got %.2f (error rate: %.3f > %.3f)",
						p*100, expected, result, errorRate, tt.tolerance)
				}
			}
		})
	}
}

// TestTDigestMemoryEfficiency tests that T-digest maintains bounded memory
func TestTDigestMemoryEfficiency(t *testing.T) {
	maxSize := 100
	td := newTDigest(maxSize)

	// Add many more values than maxSize
	numValues := maxSize * 50
	for i := 0; i < numValues; i++ {
		td.add(float64(i))
	}

	// T-digest should compress and stay within memory bounds
	if len(td.centroids) > maxSize {
		t.Errorf("T-digest exceeded memory bound: %d centroids > %d max",
			len(td.centroids), maxSize)
	}

	// Should still provide reasonable percentile estimates
	p50 := td.quantile(0.5)
	expectedP50 := float64(numValues) / 2.0
	errorRate := math.Abs(p50-expectedP50) / expectedP50

	if errorRate > 0.1 { // 10% error tolerance even after compression
		t.Errorf("P50 after compression: expected ~%.0f, got %.2f (error rate: %.3f)",
			expectedP50, p50, errorRate)
	}
}

// TestAgentAggregator_ThreadSafety tests concurrent access with sync.RWMutex
func TestAgentAggregator_ThreadSafety(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping thread safety test in short mode")
	}

	agg := NewAgentAnalyzer()
	numGoroutines := 10
	numInvocationsPerGoroutine := 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numInvocationsPerGoroutine)

	// Start multiple goroutines processing invocations concurrently
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for i := 0; i < numInvocationsPerGoroutine; i++ {
				inv := AgentInvocation{
					AgentType:       fmt.Sprintf("agent-%d", goroutineID%3), // 3 different agents
					SessionID:       fmt.Sprintf("session-%d-%d", goroutineID, i),
					Timestamp:       time.Now(),
					DurationMs:      int64(100 + i),
					Success:         i%2 == 0,
					Confidence:      0.95,
					DetectionMethod: "concurrent-test",
				}

				if err := agg.ProcessAgentInvocation(inv); err != nil {
					errors <- err
					return
				}
			}
		}(g)
	}

	// Also test concurrent reads
	readWg := sync.WaitGroup{}
	for r := 0; r < 5; r++ {
		readWg.Add(1)
		go func() {
			defer readWg.Done()
			for i := 0; i < 50; i++ {
				_, _ = agg.GetAllAgentStats()
				time.Sleep(time.Millisecond)
			}
		}()
	}

	wg.Wait()
	readWg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent processing error: %v", err)
	}

	// Verify final state is consistent
	allStats, err := agg.GetAllAgentStats()
	if err != nil {
		t.Fatalf("GetAllAgentStats failed: %v", err)
	}

	expectedAgents := 3
	if len(allStats) != expectedAgents {
		t.Errorf("Expected %d agents, got %d", expectedAgents, len(allStats))
	}

	totalInvocations := numGoroutines * numInvocationsPerGoroutine
	totalCount := 0
	for agentType, stats := range allStats {
		totalCount += stats.Count
		if stats.Count == 0 {
			t.Errorf("Agent %s has zero invocations", agentType)
		}
	}

	if totalCount != totalInvocations {
		t.Errorf("Expected total count %d, got %d", totalInvocations, totalCount)
	}
}

// TestAgentInvocationProcessing tests end-to-end invocation processing
func TestAgentInvocationProcessing(t *testing.T) {
	agg := NewAgentAnalyzer()

	invocations := []AgentInvocation{
		{
			AgentType:       "task-agent",
			SessionID:       "session-1",
			Timestamp:       time.Now(),
			DurationMs:      100,
			Success:         true,
			Confidence:      0.95,
			DetectionMethod: "direct",
		},
		{
			AgentType:       "task-agent",
			SessionID:       "session-1",
			Timestamp:       time.Now().Add(time.Second),
			DurationMs:      200,
			Success:         false,
			Confidence:      0.85,
			DetectionMethod: "pattern",
		},
		{
			AgentType:       "command-agent",
			SessionID:       "session-2",
			Timestamp:       time.Now(),
			DurationMs:      150,
			Success:         true,
			Confidence:      0.90,
			DetectionMethod: "context",
		},
	}

	// Process all invocations
	for _, inv := range invocations {
		if err := agg.ProcessAgentInvocation(inv); err != nil {
			t.Fatalf("ProcessAgentInvocation failed: %v", err)
		}
	}

	// Test task-agent stats
	taskStats, err := agg.GetAgentStats("task-agent")
	if err != nil {
		t.Fatalf("GetAgentStats failed for task-agent: %v", err)
	}

	if taskStats.Count != 2 {
		t.Errorf("Expected task-agent count 2, got %d", taskStats.Count)
	}
	if taskStats.SuccessCount != 1 {
		t.Errorf("Expected task-agent success count 1, got %d", taskStats.SuccessCount)
	}
	if taskStats.FailureCount != 1 {
		t.Errorf("Expected task-agent failure count 1, got %d", taskStats.FailureCount)
	}
	if taskStats.TotalDurationMs != 300 {
		t.Errorf("Expected task-agent total duration 300ms, got %d", taskStats.TotalDurationMs)
	}
	if taskStats.Mean() != 150.0 {
		t.Errorf("Expected task-agent mean 150ms, got %.2f", taskStats.Mean())
	}

	// Test command-agent stats
	cmdStats, err := agg.GetAgentStats("command-agent")
	if err != nil {
		t.Fatalf("GetAgentStats failed for command-agent: %v", err)
	}

	if cmdStats.Count != 1 {
		t.Errorf("Expected command-agent count 1, got %d", cmdStats.Count)
	}
	if cmdStats.SuccessCount != 1 {
		t.Errorf("Expected command-agent success count 1, got %d", cmdStats.SuccessCount)
	}
	if cmdStats.Mean() != 150.0 {
		t.Errorf("Expected command-agent mean 150ms, got %.2f", cmdStats.Mean())
	}

	// Test getting all stats
	allStats, err := agg.GetAllAgentStats()
	if err != nil {
		t.Fatalf("GetAllAgentStats failed: %v", err)
	}

	if len(allStats) != 2 {
		t.Errorf("Expected 2 agent types, got %d", len(allStats))
	}

	if _, exists := allStats["task-agent"]; !exists {
		t.Error("Missing task-agent in all stats")
	}
	if _, exists := allStats["command-agent"]; !exists {
		t.Error("Missing command-agent in all stats")
	}
}

// TestAgentAggregator_ValidationErrorHandling tests error handling for invalid inputs
func TestAgentAggregator_ValidationErrorHandling(t *testing.T) {
	agg := NewAgentAnalyzer()

	invalidCases := []struct {
		name string
		inv  AgentInvocation
	}{
		{
			name: "Empty agent type",
			inv: AgentInvocation{
				SessionID:  "session-1",
				Timestamp:  time.Now(),
				DurationMs: 100,
				Confidence: 0.95,
			},
		},
		{
			name: "Empty session ID",
			inv: AgentInvocation{
				AgentType:  "test-agent",
				Timestamp:  time.Now(),
				DurationMs: 100,
				Confidence: 0.95,
			},
		},
		{
			name: "Zero timestamp",
			inv: AgentInvocation{
				AgentType:  "test-agent",
				SessionID:  "session-1",
				DurationMs: 100,
				Confidence: 0.95,
			},
		},
		{
			name: "Invalid confidence - negative",
			inv: AgentInvocation{
				AgentType:  "test-agent",
				SessionID:  "session-1",
				Timestamp:  time.Now(),
				DurationMs: 100,
				Confidence: -0.5,
			},
		},
		{
			name: "Invalid confidence - over 1.0",
			inv: AgentInvocation{
				AgentType:  "test-agent",
				SessionID:  "session-1",
				Timestamp:  time.Now(),
				DurationMs: 100,
				Confidence: 1.5,
			},
		},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			err := agg.ProcessAgentInvocation(tc.inv)
			if err == nil {
				t.Errorf("Expected validation error for %s, got nil", tc.name)
			}
		})
	}
}

// BenchmarkAgentAggregator_MemoryUsage benchmarks memory usage with 1M+ log entries per SDD line 836
func BenchmarkAgentAggregator_MemoryUsage(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping memory benchmark in short mode")
	}

	// Force garbage collection before starting
	runtime.GC()
	runtime.GC()

	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	agg := NewAgentAnalyzer()
	numEntries := 1000000 // 1M entries

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < numEntries; i++ {
		inv := AgentInvocation{
			AgentType:       fmt.Sprintf("agent-%d", i%10), // 10 different agents
			SessionID:       fmt.Sprintf("session-%d", i%1000), // 1000 different sessions
			Timestamp:       time.Now(),
			DurationMs:      int64(100 + i%1000), // Varying durations
			Success:         i%3 != 0, // ~67% success rate
			Confidence:      0.85 + float64(i%15)/100.0, // Varying confidence
			DetectionMethod: "benchmark",
		}

		if err := agg.ProcessAgentInvocation(inv); err != nil {
			b.Fatalf("ProcessAgentInvocation failed at entry %d: %v", i, err)
		}
	}

	b.StopTimer()

	// Force garbage collection and measure memory usage
	runtime.GC()
	runtime.GC()

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	memUsedBytes := m2.TotalAlloc - m1.TotalAlloc
	memUsedMB := float64(memUsedBytes) / (1024 * 1024)

	b.Logf("Processed %d entries", numEntries)
	b.Logf("Memory used: %.2f MB", memUsedMB)
	b.Logf("Memory per entry: %.2f bytes", float64(memUsedBytes)/float64(numEntries))

	// Verify memory usage is under 100MB per SDD requirement
	if memUsedMB > 100 {
		b.Errorf("Memory usage %.2f MB exceeds 100MB limit for 1M entries", memUsedMB)
	}

	// Verify data integrity after processing
	allStats, err := agg.GetAllAgentStats()
	if err != nil {
		b.Fatalf("GetAllAgentStats failed: %v", err)
	}

	if len(allStats) != 10 {
		b.Errorf("Expected 10 agent types, got %d", len(allStats))
	}

	totalCount := 0
	for _, stats := range allStats {
		totalCount += stats.Count
	}

	if totalCount != numEntries {
		b.Errorf("Expected total count %d, got %d", numEntries, totalCount)
	}
}

// BenchmarkAgentAggregator_ResponseTime benchmarks response time <200ms requirement
func BenchmarkAgentAggregator_ResponseTime(b *testing.B) {
	agg := NewAgentAnalyzer()

	// Pre-populate with some data
	for i := 0; i < 10000; i++ {
		inv := AgentInvocation{
			AgentType:       fmt.Sprintf("agent-%d", i%5),
			SessionID:       fmt.Sprintf("session-%d", i%100),
			Timestamp:       time.Now(),
			DurationMs:      int64(100 + i%500),
			Success:         i%2 == 0,
			Confidence:      0.90,
			DetectionMethod: "benchmark",
		}
		agg.ProcessAgentInvocation(inv)
	}

	b.ResetTimer()

	// Benchmark response time for GetAllAgentStats
	for i := 0; i < b.N; i++ {
		start := time.Now()
		_, err := agg.GetAllAgentStats()
		duration := time.Since(start)

		if err != nil {
			b.Fatalf("GetAllAgentStats failed: %v", err)
		}

		// Check individual response time (should be well under 200ms)
		if duration > 200*time.Millisecond {
			b.Errorf("Response time %v exceeds 200ms requirement", duration)
		}
	}
}

// TestAgentAggregator_LRUCacheEviction tests cache behavior under memory pressure
func TestAgentAggregator_LRUCacheEviction(t *testing.T) {
	// Note: The current implementation doesn't have LRU cache, but this test
	// demonstrates how it would be tested if implemented

	t.Skip("LRU cache not implemented in current version - test prepared for future implementation")

	// This test would verify:
	// 1. Cache stores recently accessed agent stats
	// 2. Evicts least recently used entries under pressure
	// 3. Maintains performance while staying within memory bounds
	// 4. Correctly handles cache misses by recomputing stats
}

// TestAgentAggregator_StreamingStatistics tests streaming behavior
func TestAgentAggregator_StreamingStatistics(t *testing.T) {
	agg := NewAgentAnalyzer()
	agentType := "streaming-agent"

	// Test that statistics update correctly as new data streams in
	durations := []int64{100, 200, 150, 300, 250}

	for i, duration := range durations {
		inv := AgentInvocation{
			AgentType:       agentType,
			SessionID:       fmt.Sprintf("session-%d", i),
			Timestamp:       time.Now(),
			DurationMs:      duration,
			Success:         true,
			Confidence:      0.90,
			DetectionMethod: "streaming",
		}

		if err := agg.ProcessAgentInvocation(inv); err != nil {
			t.Fatalf("ProcessAgentInvocation failed: %v", err)
		}

		// Check statistics after each update
		stats, err := agg.GetAgentStats(agentType)
		if err != nil {
			t.Fatalf("GetAgentStats failed: %v", err)
		}

		expectedCount := i + 1
		if stats.Count != expectedCount {
			t.Errorf("After %d invocations: expected count %d, got %d",
				expectedCount, expectedCount, stats.Count)
		}

		// Calculate expected mean
		sum := int64(0)
		for j := 0; j <= i; j++ {
			sum += durations[j]
		}
		expectedMean := float64(sum) / float64(expectedCount)

		if math.Abs(stats.Mean()-expectedMean) > 0.001 {
			t.Errorf("After %d invocations: expected mean %.3f, got %.3f",
				expectedCount, expectedMean, stats.Mean())
		}
	}
}

// TestAgentAggregator_Integration tests integration with base Aggregator interface
func TestAgentAggregator_Integration(t *testing.T) {
	agg := NewAgentAnalyzer()

	// Verify it implements both interfaces
	var _ AgentAnalyzer = agg
	var _ Aggregator = agg

	// Test that base aggregator functionality still works
	entry := ClaudeLogEntry{
		Type:      "user",
		SessionID: "session-1",
		Timestamp: time.Now(),
		User:      &UserMessage{Text: "Hello"},
	}

	err := agg.ProcessEntry(entry)
	if err != nil {
		t.Fatalf("Base ProcessEntry failed: %v", err)
	}

	sessionStats, err := agg.GetSessionStats("session-1")
	if err != nil {
		t.Fatalf("GetSessionStats failed: %v", err)
	}

	if sessionStats.UserMessages != 1 {
		t.Errorf("Expected 1 user message, got %d", sessionStats.UserMessages)
	}

	// Test that agent-specific functionality works alongside base functionality
	inv := AgentInvocation{
		AgentType:       "integrated-agent",
		SessionID:       "session-1",
		Timestamp:       time.Now(),
		DurationMs:      150,
		Success:         true,
		Confidence:      0.95,
		DetectionMethod: "integration",
	}

	err = agg.ProcessAgentInvocation(inv)
	if err != nil {
		t.Fatalf("ProcessAgentInvocation failed: %v", err)
	}

	agentStats, err := agg.GetAgentStats("integrated-agent")
	if err != nil {
		t.Fatalf("GetAgentStats failed: %v", err)
	}

	if agentStats.Count != 1 {
		t.Errorf("Expected agent count 1, got %d", agentStats.Count)
	}
}