package stats

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// BenchmarkAggregator tests aggregator performance with various dataset sizes
func BenchmarkAggregator(b *testing.B) {
	sizes := []int{100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				agg := NewAggregator()
				now := time.Now()
				
				for j := 0; j < size; j++ {
					duration := int64(100 + j%1000)
					_ = agg.ProcessInvocation(ToolInvocation{
						ID:          fmt.Sprintf("inv%d", j),
						Name:        fmt.Sprintf("Tool%d", j%10),
						InvokedAt:   now.Add(time.Duration(j) * time.Second),
						CompletedAt: ptr(now.Add(time.Duration(j)*time.Second + time.Duration(duration)*time.Millisecond)),
						DurationMs:  &duration,
						Success:     j%10 != 0,
						SessionID:   fmt.Sprintf("session%d", j%100),
					})
				}
				
				// Verify we can get stats
				_, _ = agg.GetOverallStats()
			}
		})
	}
}

// BenchmarkWelfordAlgorithm tests Welford's algorithm performance
func BenchmarkWelfordAlgorithm(b *testing.B) {
	values := make([]float64, 10000)
	for i := range values {
		values[i] = float64(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stats := newWelfordStats()
		for _, v := range values {
			stats.update(v)
		}
		_ = stats.mean
		_ = stats.variance()
	}
}

// BenchmarkTDigest tests t-digest performance
func BenchmarkTDigest(b *testing.B) {
	values := make([]float64, 10000)
	for i := range values {
		values[i] = float64(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		td := newTDigest(1000)
		for _, v := range values {
			td.add(v)
		}
		_ = td.quantile(0.5)
		_ = td.quantile(0.95)
		_ = td.quantile(0.99)
	}
}

// BenchmarkCommandExtraction tests regex performance for command extraction
func BenchmarkCommandExtraction(b *testing.B) {
	agg := NewAggregator()
	text := `I'll run several commands:
		First, <command-name>npm install</command-name> to install dependencies.
		Then <command-name>npm test</command-name> to run tests.
		Finally, <command-name>npm run build</command-name> to build the project.`
	
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "bench",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: text,
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = agg.ProcessEntry(entry)
	}
}

// BenchmarkGetToolStats tests performance of retrieving tool statistics
func BenchmarkGetToolStats(b *testing.B) {
	agg := NewAggregator()
	now := time.Now()
	
	// Setup: Add 1000 invocations for 10 different tools
	for i := 0; i < 1000; i++ {
		duration := int64(100 + i%1000)
		_ = agg.ProcessInvocation(ToolInvocation{
			ID:          fmt.Sprintf("inv%d", i),
			Name:        fmt.Sprintf("Tool%d", i%10),
			InvokedAt:   now,
			CompletedAt: ptr(now.Add(time.Duration(duration) * time.Millisecond)),
			DurationMs:  &duration,
			Success:     true,
			SessionID:   fmt.Sprintf("session%d", i%10),
		})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			_, _ = agg.GetToolStats(fmt.Sprintf("Tool%d", j))
		}
	}
}

// BenchmarkMerge tests performance of merging aggregators
func BenchmarkMerge(b *testing.B) {
	createAggregator := func(offset int) Aggregator {
		agg := NewAggregator()
		now := time.Now()
		
		for i := 0; i < 100; i++ {
			duration := int64(100 + i%100)
			_ = agg.ProcessInvocation(ToolInvocation{
				ID:          fmt.Sprintf("inv%d_%d", offset, i),
				Name:        fmt.Sprintf("Tool%d", i%5),
				InvokedAt:   now,
				CompletedAt: ptr(now.Add(time.Duration(duration) * time.Millisecond)),
				DurationMs:  &duration,
				Success:     true,
				SessionID:   fmt.Sprintf("session%d_%d", offset, i%10),
			})
		}
		return agg
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		agg1 := createAggregator(0)
		agg2 := createAggregator(1)
		_ = agg1.Merge(agg2)
	}
}

// TestMemoryScaling verifies memory usage scales linearly with data
func TestMemoryScaling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory scaling test in short mode")
	}
	
	sizes := []int{1000, 5000, 10000}
	memUsages := make([]uint64, len(sizes))
	
	for idx, size := range sizes {
		// Force GC before measurement
		runtime.GC()
		runtime.GC()
		
		var m1 runtime.MemStats
		runtime.ReadMemStats(&m1)
		
		agg := NewAggregator()
		now := time.Now()
		
		for i := 0; i < size; i++ {
			duration := int64(100 + i%1000)
			_ = agg.ProcessInvocation(ToolInvocation{
				ID:          fmt.Sprintf("inv%d", i),
				Name:        fmt.Sprintf("Tool%d", i%20), // 20 different tools
				InvokedAt:   now,
				CompletedAt: ptr(now.Add(time.Duration(duration) * time.Millisecond)),
				DurationMs:  &duration,
				Success:     true,
				SessionID:   fmt.Sprintf("session%d", i%100), // 100 sessions
			})
		}
		
		// Force GC after processing
		runtime.GC()
		
		var m2 runtime.MemStats
		runtime.ReadMemStats(&m2)
		
		memUsages[idx] = m2.Alloc - m1.Alloc
		
		// Verify stats are correct
		stats, err := agg.GetOverallStats()
		if err != nil {
			t.Fatalf("GetOverallStats failed: %v", err)
		}
		
		if stats.TotalToolCalls != size {
			t.Errorf("Size %d: expected %d tool calls, got %d", size, size, stats.TotalToolCalls)
		}
		
		t.Logf("Size %d: Memory used = %d KB (%.2f bytes/entry)", 
			size, memUsages[idx]/1024, float64(memUsages[idx])/float64(size))
	}
	
	// Check that memory usage scales roughly linearly
	// The ratio between 10k and 1k entries should be less than 15x (allowing for some overhead)
	ratio := float64(memUsages[2]) / float64(memUsages[0])
	if ratio > 15 {
		t.Errorf("Memory usage not scaling linearly: 10k/1k ratio = %.2f", ratio)
	}
}

// TestConcurrentAccess verifies thread safety
func TestConcurrentAccess(t *testing.T) {
	agg := NewAggregator()
	now := time.Now()
	
	// Run concurrent operations
	done := make(chan bool, 3)
	
	// Writer 1: Process invocations
	go func() {
		for i := 0; i < 100; i++ {
			duration := int64(100)
			_ = agg.ProcessInvocation(ToolInvocation{
				ID:          fmt.Sprintf("w1_inv%d", i),
				Name:        "Writer1Tool",
				InvokedAt:   now,
				CompletedAt: ptr(now.Add(100 * time.Millisecond)),
				DurationMs:  &duration,
				Success:     true,
				SessionID:   "writer1",
			})
		}
		done <- true
	}()
	
	// Writer 2: Process entries
	go func() {
		for i := 0; i < 100; i++ {
			_ = agg.ProcessEntry(ClaudeLogEntry{
				Type:      "user",
				SessionID: "writer2",
				Timestamp: now,
				User:      &UserMessage{Text: fmt.Sprintf("Message %d", i)},
			})
		}
		done <- true
	}()
	
	// Reader: Get stats
	go func() {
		for i := 0; i < 50; i++ {
			_, _ = agg.GetOverallStats()
			_, _ = agg.GetSessionStats("writer1")
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()
	
	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
	
	// Verify final state
	stats, err := agg.GetOverallStats()
	if err != nil {
		t.Fatalf("GetOverallStats failed: %v", err)
	}
	
	if stats.TotalToolCalls != 100 {
		t.Errorf("Expected 100 tool calls, got %d", stats.TotalToolCalls)
	}
}