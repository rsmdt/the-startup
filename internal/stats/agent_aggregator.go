package stats

import (
	"container/list"
	"fmt"
	"sync"
)

// AgentAnalyzer interface - extends existing Aggregator
type AgentAnalyzer interface {
	// Embed existing aggregator
	Aggregator

	// Agent-specific methods
	ProcessAgentInvocation(inv AgentInvocation) error
	GetAgentStats(agentType string) (*GlobalAgentStats, error)
	GetAllAgentStats() (map[string]*GlobalAgentStats, error)
}

// LRU cache entry for agent statistics
type agentCacheEntry struct {
	agentType string
	stats     *GlobalAgentStats
}

// LRU cache for agent statistics (max 1000 entries)
type agentLRUCache struct {
	maxSize  int
	cache    map[string]*list.Element
	orderList *list.List
}

// newAgentLRUCache creates a new LRU cache for agent stats
func newAgentLRUCache(maxSize int) *agentLRUCache {
	return &agentLRUCache{
		maxSize:   maxSize,
		cache:     make(map[string]*list.Element),
		orderList: list.New(),
	}
}

// get retrieves agent stats from cache, returns nil if not found
func (lru *agentLRUCache) get(agentType string) *GlobalAgentStats {
	if elem, found := lru.cache[agentType]; found {
		// Move to front (most recently used)
		lru.orderList.MoveToFront(elem)
		return elem.Value.(*agentCacheEntry).stats
	}
	return nil
}

// put adds or updates agent stats in cache
func (lru *agentLRUCache) put(agentType string, stats *GlobalAgentStats) {
	if elem, found := lru.cache[agentType]; found {
		// Update existing entry and move to front
		elem.Value.(*agentCacheEntry).stats = stats
		lru.orderList.MoveToFront(elem)
		return
	}

	// Add new entry
	entry := &agentCacheEntry{agentType: agentType, stats: stats}
	elem := lru.orderList.PushFront(entry)
	lru.cache[agentType] = elem

	// Evict if necessary
	if lru.orderList.Len() > lru.maxSize {
		// Remove least recently used (from back)
		oldest := lru.orderList.Back()
		if oldest != nil {
			oldestEntry := oldest.Value.(*agentCacheEntry)
			delete(lru.cache, oldestEntry.agentType)
			lru.orderList.Remove(oldest)
		}
	}
}

// getAllStats returns all cached agent stats as a map
func (lru *agentLRUCache) getAllStats() map[string]*GlobalAgentStats {
	result := make(map[string]*GlobalAgentStats)
	for agentType, elem := range lru.cache {
		result[agentType] = elem.Value.(*agentCacheEntry).stats
	}
	return result
}

// agentAggregator implements AgentAnalyzer interface with LRU cache
type agentAggregator struct {
	*aggregator
	mu         sync.RWMutex
	agentCache *agentLRUCache
}

// NewAgentAnalyzer creates a new agent analyzer with LRU cache
func NewAgentAnalyzer() AgentAnalyzer {
	return &agentAggregator{
		aggregator:  NewAggregator().(*aggregator),
		agentCache: newAgentLRUCache(1000), // Max 1000 agents in cache
	}
}

// ProcessAgentInvocation processes an agent invocation with thread-safe LRU cache
func (a *agentAggregator) ProcessAgentInvocation(inv AgentInvocation) error {
	if err := validateAgentInvocation(inv); err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	// Get or create agent stats from cache
	stats := a.agentCache.get(inv.AgentType)
	if stats == nil {
		stats = &GlobalAgentStats{}
	}

	// Update statistics
	stats.UpdateStats(inv.DurationMs)
	if inv.Success {
		stats.SuccessCount++
	} else {
		stats.FailureCount++
	}

	// Update cache with modified stats
	a.agentCache.put(inv.AgentType, stats)

	return nil
}

// GetAgentStats returns stats for a specific agent from LRU cache
func (a *agentAggregator) GetAgentStats(agentType string) (*GlobalAgentStats, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	stats := a.agentCache.get(agentType)
	if stats == nil {
		return nil, fmt.Errorf("no stats for agent type: %s", agentType)
	}
	return stats, nil
}

// GetAllAgentStats returns all agent stats from LRU cache
func (a *agentAggregator) GetAllAgentStats() (map[string]*GlobalAgentStats, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Get all stats from cache (creates a copy)
	return a.agentCache.getAllStats(), nil
}

// validateAgentInvocation validates an agent invocation
func validateAgentInvocation(inv AgentInvocation) error {
	if inv.AgentType == "" {
		return fmt.Errorf("agent type cannot be empty")
	}

	if inv.SessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	if inv.Timestamp.IsZero() {
		return fmt.Errorf("timestamp cannot be zero")
	}

	if !isValidConfidence(inv.Confidence) {
		return fmt.Errorf("confidence score must be between 0.0 and 1.0, got %f", inv.Confidence)
	}

	return nil
}