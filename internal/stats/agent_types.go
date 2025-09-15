package stats

import (
	"math"
	"time"
)

// AgentInvocation represents an agent invocation from logs
type AgentInvocation struct {
	AgentType       string
	SessionID       string
	Timestamp       time.Time
	DurationMs      int64
	Success         bool
	Confidence      float64
	DetectionMethod string
}

// GlobalAgentStats tracks agent usage statistics with streaming algorithms
type GlobalAgentStats struct {
	Count           int
	SuccessCount    int
	FailureCount    int
	TotalDurationMs int64

	// Welford's algorithm for streaming mean and variance
	durationStats *welfordStats

	// T-digest for streaming percentile calculation
	percentiles *tDigest
}

// UpdateStats updates statistics using streaming algorithms
func (s *GlobalAgentStats) UpdateStats(durationMs int64) {
	s.Count++
	s.TotalDurationMs += durationMs

	// Initialize streaming stats if needed
	if s.durationStats == nil {
		s.durationStats = newWelfordStats()
	}
	if s.percentiles == nil {
		s.percentiles = newTDigest(100) // 100 centroids for good accuracy
	}

	// Update streaming statistics
	duration := float64(durationMs)
	s.durationStats.update(duration)
	s.percentiles.add(duration)
}

// Variance returns the sample variance
func (s *GlobalAgentStats) Variance() float64 {
	if s.durationStats == nil {
		return 0
	}
	return s.durationStats.variance()
}

// StdDev returns the standard deviation
func (s *GlobalAgentStats) StdDev() float64 {
	return math.Sqrt(s.Variance())
}

// Mean returns the mean duration
func (s *GlobalAgentStats) Mean() float64 {
	if s.durationStats == nil {
		return 0
	}
	return s.durationStats.mean
}

// Min returns the minimum duration
func (s *GlobalAgentStats) Min() *float64 {
	if s.durationStats == nil {
		return nil
	}
	return s.durationStats.min
}

// Max returns the maximum duration
func (s *GlobalAgentStats) Max() *float64 {
	if s.durationStats == nil {
		return nil
	}
	return s.durationStats.max
}

// Quantile returns the specified quantile (percentile) of durations
func (s *GlobalAgentStats) Quantile(q float64) float64 {
	if s.percentiles == nil {
		return 0
	}
	return s.percentiles.quantile(q)
}

// P50 returns the median duration (50th percentile)
func (s *GlobalAgentStats) P50() float64 {
	return s.Quantile(0.5)
}

// P95 returns the 95th percentile duration
func (s *GlobalAgentStats) P95() float64 {
	return s.Quantile(0.95)
}

// P99 returns the 99th percentile duration
func (s *GlobalAgentStats) P99() float64 {
	return s.Quantile(0.99)
}

// isValidConfidence checks if confidence score is in valid range [0.0, 1.0]
func isValidConfidence(confidence float64) bool {
	if math.IsNaN(confidence) || math.IsInf(confidence, 0) {
		return false
	}
	return confidence >= 0.0 && confidence <= 1.0
}