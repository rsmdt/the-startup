package stats

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// DelegationPattern represents agent transition patterns per PRD lines 89
type DelegationPattern struct {
	SourceAgent string
	TargetAgent string
	Count       int
	LastSeen    time.Time
	SessionID   string
}

// AgentCoOccurrence represents agents used together in sessions per PRD lines 87
type AgentCoOccurrence struct {
	Agent1    string
	Agent2    string
	Count     int
	Sessions  []string // Session IDs where co-occurrence happened
	LastSeen  time.Time
}

// DelegationChain represents a sequence of agent transitions per PRD lines 89
type DelegationChain struct {
	Agents    []string
	SessionID string
	StartTime time.Time
	EndTime   time.Time
}

// DelegationAnalyzer interface for tracking delegation patterns
type DelegationAnalyzer interface {
	// ProcessAgentTransition processes a transition from one agent to another
	ProcessAgentTransition(sessionID string, fromAgent string, toAgent string, timestamp time.Time) error

	// GetDelegationChain builds delegation chain for a session per PRD lines 89
	GetDelegationChain(sessionID string) (*DelegationChain, error)

	// GetCoOccurrencePatterns gets agent co-occurrence within sessions per PRD lines 87
	GetCoOccurrencePatterns() ([]AgentCoOccurrence, error)

	// GetTransitionPatterns gets agent transition patterns (A → B flows)
	GetTransitionPatterns() ([]DelegationPattern, error)

	// ProcessSessionBoundary handles session start/end detection per PRD lines 90
	ProcessSessionBoundary(sessionID string, eventType string, timestamp time.Time) error

	// GetFilteredPatterns applies time period filtering to patterns
	GetFilteredPatterns(startTime, endTime time.Time) ([]DelegationPattern, error)
}

// realDelegationAnalyzer implements the DelegationAnalyzer interface
type realDelegationAnalyzer struct {
	mu sync.RWMutex

	// Core data structures
	transitions      []DelegationPattern
	sessionChains    map[string]*DelegationChain
	sessionBounds    map[string][]time.Time
	coOccurrences    []AgentCoOccurrence
	sessionAgents    map[string]map[string]bool // sessionID -> set of agents

	// Optimization indexes
	transitionIndex  map[string]map[string]int  // sourceAgent -> targetAgent -> index in transitions
	coOccurrenceIndex map[string]map[string]int  // agent1 -> agent2 -> index in coOccurrences
}

// NewDelegationAnalyzer creates a new delegation analyzer
func NewDelegationAnalyzer() DelegationAnalyzer {
	return &realDelegationAnalyzer{
		transitions:       []DelegationPattern{},
		sessionChains:     make(map[string]*DelegationChain),
		sessionBounds:     make(map[string][]time.Time),
		coOccurrences:     []AgentCoOccurrence{},
		sessionAgents:     make(map[string]map[string]bool),
		transitionIndex:   make(map[string]map[string]int),
		coOccurrenceIndex: make(map[string]map[string]int),
	}
}

// ProcessAgentTransition processes a transition from one agent to another
func (r *realDelegationAnalyzer) ProcessAgentTransition(sessionID string, fromAgent string, toAgent string, timestamp time.Time) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}
	if fromAgent == "" {
		return fmt.Errorf("source agent cannot be empty")
	}
	if toAgent == "" {
		return fmt.Errorf("target agent cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Update transition patterns
	r.updateTransitionPattern(fromAgent, toAgent, sessionID, timestamp)

	// Update delegation chain
	r.updateDelegationChain(sessionID, fromAgent, toAgent, timestamp)

	// Update session agents for co-occurrence tracking
	r.trackSessionAgents(sessionID, fromAgent, toAgent)

	// Update co-occurrences
	r.updateCoOccurrences(sessionID, fromAgent, toAgent, timestamp)

	return nil
}

// updateTransitionPattern updates or creates transition patterns
func (r *realDelegationAnalyzer) updateTransitionPattern(fromAgent, toAgent, sessionID string, timestamp time.Time) {
	// Check if transition already exists using index
	if _, exists := r.transitionIndex[fromAgent]; !exists {
		r.transitionIndex[fromAgent] = make(map[string]int)
	}

	if idx, exists := r.transitionIndex[fromAgent][toAgent]; exists {
		// Update existing transition
		r.transitions[idx].Count++
		r.transitions[idx].LastSeen = timestamp
		r.transitions[idx].SessionID = sessionID
	} else {
		// Create new transition
		newTransition := DelegationPattern{
			SourceAgent: fromAgent,
			TargetAgent: toAgent,
			Count:       1,
			LastSeen:    timestamp,
			SessionID:   sessionID,
		}
		r.transitions = append(r.transitions, newTransition)
		r.transitionIndex[fromAgent][toAgent] = len(r.transitions) - 1
	}
}

// updateDelegationChain updates the delegation chain for a session
func (r *realDelegationAnalyzer) updateDelegationChain(sessionID, fromAgent, toAgent string, timestamp time.Time) {
	chain, exists := r.sessionChains[sessionID]
	if !exists {
		chain = &DelegationChain{
			Agents:    []string{fromAgent},
			SessionID: sessionID,
			StartTime: timestamp,
			EndTime:   timestamp,
		}
		r.sessionChains[sessionID] = chain
	}

	// Always add target agent to chain to preserve transitions including self-delegation
	chain.Agents = append(chain.Agents, toAgent)

	// Update end time
	if timestamp.After(chain.EndTime) {
		chain.EndTime = timestamp
	}
	if timestamp.Before(chain.StartTime) {
		chain.StartTime = timestamp
	}
}

// trackSessionAgents tracks which agents are used in each session for co-occurrence analysis
func (r *realDelegationAnalyzer) trackSessionAgents(sessionID, fromAgent, toAgent string) {
	if r.sessionAgents[sessionID] == nil {
		r.sessionAgents[sessionID] = make(map[string]bool)
	}
	r.sessionAgents[sessionID][fromAgent] = true
	r.sessionAgents[sessionID][toAgent] = true
}

// updateCoOccurrences updates agent co-occurrence patterns
func (r *realDelegationAnalyzer) updateCoOccurrences(sessionID, fromAgent, toAgent string, timestamp time.Time) {
	// Ensure consistent ordering (agent1 < agent2)
	agent1, agent2 := fromAgent, toAgent
	if agent1 > agent2 {
		agent1, agent2 = agent2, agent1
	}

	// Skip if same agent (no co-occurrence with itself)
	if agent1 == agent2 {
		return
	}

	// Check if co-occurrence already exists using index
	if _, exists := r.coOccurrenceIndex[agent1]; !exists {
		r.coOccurrenceIndex[agent1] = make(map[string]int)
	}

	if idx, exists := r.coOccurrenceIndex[agent1][agent2]; exists {
		// Update existing co-occurrence
		coOcc := &r.coOccurrences[idx]

		// Check if session is already tracked
		sessionFound := false
		for _, sid := range coOcc.Sessions {
			if sid == sessionID {
				sessionFound = true
				break
			}
		}

		if !sessionFound {
			coOcc.Sessions = append(coOcc.Sessions, sessionID)
			coOcc.Count = len(coOcc.Sessions)
		}

		if timestamp.After(coOcc.LastSeen) {
			coOcc.LastSeen = timestamp
		}
	} else {
		// Create new co-occurrence
		newCoOcc := AgentCoOccurrence{
			Agent1:   agent1,
			Agent2:   agent2,
			Count:    1,
			Sessions: []string{sessionID},
			LastSeen: timestamp,
		}
		r.coOccurrences = append(r.coOccurrences, newCoOcc)
		r.coOccurrenceIndex[agent1][agent2] = len(r.coOccurrences) - 1
	}
}

// GetDelegationChain builds delegation chain for a session per PRD lines 89
func (r *realDelegationAnalyzer) GetDelegationChain(sessionID string) (*DelegationChain, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	chain, exists := r.sessionChains[sessionID]
	if !exists {
		return nil, nil
	}

	// Return a copy to prevent external modification
	return &DelegationChain{
		Agents:    append([]string{}, chain.Agents...),
		SessionID: chain.SessionID,
		StartTime: chain.StartTime,
		EndTime:   chain.EndTime,
	}, nil
}

// GetCoOccurrencePatterns gets agent co-occurrence within sessions per PRD lines 87
func (r *realDelegationAnalyzer) GetCoOccurrencePatterns() ([]AgentCoOccurrence, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make([]AgentCoOccurrence, len(r.coOccurrences))
	for i, coOcc := range r.coOccurrences {
		result[i] = AgentCoOccurrence{
			Agent1:   coOcc.Agent1,
			Agent2:   coOcc.Agent2,
			Count:    coOcc.Count,
			Sessions: append([]string{}, coOcc.Sessions...),
			LastSeen: coOcc.LastSeen,
		}
	}

	return result, nil
}

// GetTransitionPatterns gets agent transition patterns (A → B flows)
func (r *realDelegationAnalyzer) GetTransitionPatterns() ([]DelegationPattern, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make([]DelegationPattern, len(r.transitions))
	copy(result, r.transitions)

	return result, nil
}

// ProcessSessionBoundary handles session start/end detection per PRD lines 90
func (r *realDelegationAnalyzer) ProcessSessionBoundary(sessionID string, eventType string, timestamp time.Time) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}
	if eventType == "" {
		return fmt.Errorf("event type cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.sessionBounds[sessionID] == nil {
		r.sessionBounds[sessionID] = []time.Time{}
	}
	r.sessionBounds[sessionID] = append(r.sessionBounds[sessionID], timestamp)

	return nil
}

// GetFilteredPatterns applies time period filtering to patterns
func (r *realDelegationAnalyzer) GetFilteredPatterns(startTime, endTime time.Time) ([]DelegationPattern, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []DelegationPattern
	for _, pattern := range r.transitions {
		if pattern.LastSeen.After(startTime) && pattern.LastSeen.Before(endTime) {
			filtered = append(filtered, pattern)
		}
	}

	return filtered, nil
}

// ProcessSessionAgents analyzes all agents in a session and updates co-occurrence patterns
func (r *realDelegationAnalyzer) ProcessSessionAgents(sessionID string, agents []string, timestamp time.Time) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Track all agents in this session
	if r.sessionAgents[sessionID] == nil {
		r.sessionAgents[sessionID] = make(map[string]bool)
	}
	for _, agent := range agents {
		r.sessionAgents[sessionID][agent] = true
	}

	// Generate co-occurrence pairs for all combinations
	for i := 0; i < len(agents); i++ {
		for j := i + 1; j < len(agents); j++ {
			agent1, agent2 := agents[i], agents[j]

			// Ensure consistent ordering
			if agent1 > agent2 {
				agent1, agent2 = agent2, agent1
			}

			// Update co-occurrence
			if _, exists := r.coOccurrenceIndex[agent1]; !exists {
				r.coOccurrenceIndex[agent1] = make(map[string]int)
			}

			if idx, exists := r.coOccurrenceIndex[agent1][agent2]; exists {
				coOcc := &r.coOccurrences[idx]

				// Check if session is already tracked
				sessionFound := false
				for _, sid := range coOcc.Sessions {
					if sid == sessionID {
						sessionFound = true
						break
					}
				}

				if !sessionFound {
					coOcc.Sessions = append(coOcc.Sessions, sessionID)
					coOcc.Count = len(coOcc.Sessions)
				}

				if timestamp.After(coOcc.LastSeen) {
					coOcc.LastSeen = timestamp
				}
			} else {
				// Create new co-occurrence
				newCoOcc := AgentCoOccurrence{
					Agent1:   agent1,
					Agent2:   agent2,
					Count:    1,
					Sessions: []string{sessionID},
					LastSeen: timestamp,
				}
				r.coOccurrences = append(r.coOccurrences, newCoOcc)
				r.coOccurrenceIndex[agent1][agent2] = len(r.coOccurrences) - 1
			}
		}
	}

	return nil
}

// GetSessionBoundaries returns session boundary timestamps
func (r *realDelegationAnalyzer) GetSessionBoundaries(sessionID string) ([]time.Time, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	boundaries, exists := r.sessionBounds[sessionID]
	if !exists {
		return []time.Time{}, nil
	}

	// Return a copy to prevent external modification
	result := make([]time.Time, len(boundaries))
	copy(result, boundaries)

	return result, nil
}

// GetTopTransitionPatterns returns the most frequent transition patterns
func (r *realDelegationAnalyzer) GetTopTransitionPatterns(limit int) ([]DelegationPattern, error) {
	patterns, err := r.GetTransitionPatterns()
	if err != nil {
		return nil, err
	}

	// Sort by count (descending)
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Count > patterns[j].Count
	})

	if limit > 0 && limit < len(patterns) {
		patterns = patterns[:limit]
	}

	return patterns, nil
}

// GetTopCoOccurrencePatterns returns the most frequent co-occurrence patterns
func (r *realDelegationAnalyzer) GetTopCoOccurrencePatterns(limit int) ([]AgentCoOccurrence, error) {
	coOccurrences, err := r.GetCoOccurrencePatterns()
	if err != nil {
		return nil, err
	}

	// Sort by count (descending)
	sort.Slice(coOccurrences, func(i, j int) bool {
		return coOccurrences[i].Count > coOccurrences[j].Count
	})

	if limit > 0 && limit < len(coOccurrences) {
		coOccurrences = coOccurrences[:limit]
	}

	return coOccurrences, nil
}

// HasCircularDelegation detects if there are circular delegation patterns in a session
func (r *realDelegationAnalyzer) HasCircularDelegation(sessionID string) (bool, error) {
	chain, err := r.GetDelegationChain(sessionID)
	if err != nil {
		return false, err
	}

	if chain == nil || len(chain.Agents) == 0 {
		return false, nil
	}

	// Handle single agent case (no transitions)
	if len(chain.Agents) == 1 {
		return false, nil
	}

	// Check for self-delegation (same agent appears consecutively)
	for i := 0; i < len(chain.Agents)-1; i++ {
		if chain.Agents[i] == chain.Agents[i+1] {
			return true, nil
		}
	}

	// Check for circular patterns (agent appearing more than once non-consecutively)
	seen := make(map[string]int)
	for _, agent := range chain.Agents {
		seen[agent]++
		if seen[agent] > 1 {
			return true, nil
		}
	}

	return false, nil
}

// Reset clears all delegation data
func (r *realDelegationAnalyzer) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.transitions = []DelegationPattern{}
	r.sessionChains = make(map[string]*DelegationChain)
	r.sessionBounds = make(map[string][]time.Time)
	r.coOccurrences = []AgentCoOccurrence{}
	r.sessionAgents = make(map[string]map[string]bool)
	r.transitionIndex = make(map[string]map[string]int)
	r.coOccurrenceIndex = make(map[string]map[string]int)
}

// GetStats returns delegation analysis statistics
func (r *realDelegationAnalyzer) GetStats() (*DelegationStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return &DelegationStats{
		TotalTransitions:    len(r.transitions),
		TotalCoOccurrences:  len(r.coOccurrences),
		TotalSessions:       len(r.sessionChains),
		SessionsWithChains:  len(r.sessionChains),
		AvgChainLength:      r.calculateAvgChainLength(),
	}, nil
}

// DelegationStats represents overall delegation statistics
type DelegationStats struct {
	TotalTransitions    int     `json:"total_transitions"`
	TotalCoOccurrences  int     `json:"total_co_occurrences"`
	TotalSessions       int     `json:"total_sessions"`
	SessionsWithChains  int     `json:"sessions_with_chains"`
	AvgChainLength      float64 `json:"avg_chain_length"`
}

// calculateAvgChainLength calculates average delegation chain length
func (r *realDelegationAnalyzer) calculateAvgChainLength() float64 {
	if len(r.sessionChains) == 0 {
		return 0.0
	}

	totalLength := 0
	for _, chain := range r.sessionChains {
		totalLength += len(chain.Agents)
	}

	return float64(totalLength) / float64(len(r.sessionChains))
}