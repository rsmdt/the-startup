# Agent Refactoring System - Implementation Specification

## Technical Stack & Tooling

### Core Technologies
```yaml
Language: Go 1.21+
  Rationale: 
    - Native concurrency for parallel processing
    - Strong typing for transformation rules
    - Excellent performance for batch operations
    - Existing codebase compatibility

Parser Libraries:
  - gopkg.in/yaml.v3: YAML frontmatter parsing
  - github.com/yuin/goldmark: Markdown AST parsing
  - github.com/PuerkitoBio/goquery: HTML/DOM manipulation

NLP & Analysis:
  - github.com/jdkato/prose/v2: NLP tokenization and POS tagging
  - github.com/texttheater/golang-levenshtein: Semantic similarity
  - Custom verb classifier for HOW/WHAT detection

Validation:
  - github.com/go-playground/validator/v10: Struct validation
  - github.com/stretchr/testify: Testing assertions
  - github.com/google/go-cmp: Deep comparison

Infrastructure:
  - github.com/spf13/cobra: CLI interface
  - github.com/charmbracelet/bubbletea: Interactive TUI
  - github.com/olekukonko/tablewriter: Report generation
```

## Data Structures & Models

### Agent File Model
```go
type AgentFile struct {
    // Metadata from frontmatter
    Metadata struct {
        Name        string   `yaml:"name" validate:"required"`
        Description string   `yaml:"description" validate:"required,max=200"`
        Model       string   `yaml:"model" validate:"oneof=inherit custom"`
        Version     string   `yaml:"version"`
        Category    string   `yaml:"category"`
    }
    
    // Parsed content structure
    Content struct {
        RoleStatement    string
        FocusAreas      []FocusArea
        CoreObjectives  []Objective
        Approach        []ApproachStep
        SuccessCriteria []Criterion
        Boundaries      BoundaryDefinition
        AntiPatterns    []string
        ExpectedOutput  []OutputSpec
    }
    
    // Analysis metrics
    Metrics struct {
        LineCount      int
        WordCount      int
        HowRatio       float64
        WhatRatio      float64
        ComplexityScore float64
        QualityScore   float64
    }
    
    // Transformation tracking
    Transformation struct {
        Original       *AgentFile
        TransformedAt  time.Time
        Rules          []TransformationRule
        ValidationResults []ValidationResult
    }
}

type FocusArea struct {
    Title       string
    Description string
    IsHow       bool // Classified as HOW-focused
    IsWhat      bool // Classified as WHAT-focused
}

type Objective struct {
    Statement   string
    Measurable  bool
    Outcome     string
    Constraints []string
}

type BoundaryDefinition struct {
    InScope     []ScopeItem
    OutOfScope  []ScopeItem
    Delegations map[string]string // agent -> responsibility
}
```

### Transformation Rules Engine
```go
type TransformationRule interface {
    Name() string
    Applicable(section Section) bool
    Transform(content string) (string, error)
    Priority() int
}

type RuleEngine struct {
    Rules []TransformationRule
    Config RuleConfig
}

// Example transformation rules
type ConvertStepsToObjectivesRule struct{}
type ExtractMeasurableOutcomesRule struct{}
type SimplifyVerboseDescriptionsRule struct{}
type ConsolidateDuplicateConceptsRule struct{}
type ExtractSuccessCriteriaRule struct{}

type RuleConfig struct {
    MaxObjectives      int     // Default: 5
    MaxLineLength      int     // Default: 100
    TargetWhatRatio    float64 // Default: 0.70
    MinQualityScore    float64 // Default: 0.85
    PreserveKeywords   []string // Domain-specific terms to preserve
}
```

### Validation Framework
```go
type Validator interface {
    Name() string
    Validate(agent *AgentFile) ValidationResult
    Severity() ValidationSeverity
}

type ValidationResult struct {
    Passed      bool
    Score       float64
    Issues      []ValidationIssue
    Suggestions []string
}

type ValidationIssue struct {
    Severity    ValidationSeverity
    Location    string // Section or line reference
    Message     string
    FixHint     string
}

type ValidationSeverity int
const (
    SeverityError ValidationSeverity = iota
    SeverityWarning
    SeverityInfo
)

// Validator implementations
type WhatRatioValidator struct {
    MinRatio float64
}

type LineLengthValidator struct {
    MaxLines int
    Tolerance float64
}

type SemanticPreservationValidator struct {
    SimilarityThreshold float64
    EmbeddingsModel     string
}

type CrossReferenceValidator struct {
    AgentRegistry map[string]*AgentFile
}
```

## Processing Pipeline Implementation

### Analysis Pipeline
```go
type AnalysisPipeline struct {
    Parser    *AgentParser
    Analyzer  *ContentAnalyzer
    Collector *MetricsCollector
}

func (p *AnalysisPipeline) Process(filepath string) (*AnalysisResult, error) {
    // Stage 1: Parse file
    agent, err := p.Parser.ParseFile(filepath)
    if err != nil {
        return nil, fmt.Errorf("parsing failed: %w", err)
    }
    
    // Stage 2: Analyze content
    analysis := p.Analyzer.Analyze(agent)
    
    // Stage 3: Collect metrics
    metrics := p.Collector.Collect(agent, analysis)
    
    return &AnalysisResult{
        Agent:    agent,
        Analysis: analysis,
        Metrics:  metrics,
    }, nil
}

type ContentAnalyzer struct {
    VerbClassifier *VerbClassifier
    PatternMatcher *PatternMatcher
}

func (a *ContentAnalyzer) Analyze(agent *AgentFile) *Analysis {
    howCount := 0
    whatCount := 0
    
    // Analyze each section
    for _, section := range agent.GetSections() {
        if a.isHowFocused(section) {
            howCount++
        } else if a.isWhatFocused(section) {
            whatCount++
        }
    }
    
    return &Analysis{
        HowRatio:  float64(howCount) / float64(howCount+whatCount),
        WhatRatio: float64(whatCount) / float64(howCount+whatCount),
        Sections:  analyzedSections,
    }
}

func (a *ContentAnalyzer) isHowFocused(text string) bool {
    // Detect imperative verbs and procedural patterns
    imperatives := []string{"create", "implement", "build", "check", "scan", "validate"}
    procedural := []string{"first", "then", "next", "finally", "step"}
    
    score := 0
    for _, pattern := range imperatives {
        if strings.Contains(strings.ToLower(text), pattern) {
            score++
        }
    }
    
    return score > threshold
}
```

### Transformation Pipeline
```go
type TransformationPipeline struct {
    Transformer  *ContentTransformer
    Reorganizer  *StructureReorganizer
    Optimizer    *ContentOptimizer
    RuleEngine   *RuleEngine
}

func (p *TransformationPipeline) Transform(agent *AgentFile) (*AgentFile, error) {
    // Create working copy
    transformed := agent.Clone()
    
    // Stage 1: Apply transformation rules
    for _, section := range transformed.GetSections() {
        rules := p.RuleEngine.GetApplicableRules(section)
        for _, rule := range rules {
            newContent, err := rule.Transform(section.Content)
            if err != nil {
                continue // Log and skip failed transformations
            }
            section.Content = newContent
        }
    }
    
    // Stage 2: Reorganize structure
    transformed = p.Reorganizer.Reorganize(transformed)
    
    // Stage 3: Optimize content
    transformed = p.Optimizer.Optimize(transformed)
    
    return transformed, nil
}

type ContentTransformer struct {
    Rules []TransformationRule
}

func (t *ContentTransformer) TransformApproachToObjectives(approach []string) []Objective {
    objectives := []Objective{}
    
    for _, step := range approach {
        // Extract outcome from procedural step
        outcome := t.extractOutcome(step)
        if outcome != "" {
            objectives = append(objectives, Objective{
                Statement:  outcome,
                Measurable: t.isMeasurable(outcome),
                Outcome:    t.extractExpectedResult(step),
            })
        }
    }
    
    return t.consolidateObjectives(objectives)
}
```

### Version Control Implementation
```go
type VersionManager struct {
    BaseDir string
    Git     *git.Repository
}

func (v *VersionManager) CreateSnapshot(agent *AgentFile) (*Version, error) {
    version := &Version{
        ID:        generateVersionID(),
        Agent:     agent,
        Timestamp: time.Now(),
        Metrics:   agent.Metrics,
    }
    
    // Save to filesystem
    versionPath := filepath.Join(v.BaseDir, agent.Metadata.Name, version.ID)
    if err := os.MkdirAll(versionPath, 0755); err != nil {
        return nil, err
    }
    
    // Write agent file
    if err := v.writeAgentFile(versionPath, agent); err != nil {
        return nil, err
    }
    
    // Write metadata
    if err := v.writeMetadata(versionPath, version); err != nil {
        return nil, err
    }
    
    // Git commit for additional safety
    if err := v.commitVersion(version); err != nil {
        return nil, err
    }
    
    return version, nil
}

func (v *VersionManager) Rollback(agentName string, versionID string) error {
    // Load specified version
    version, err := v.LoadVersion(agentName, versionID)
    if err != nil {
        return err
    }
    
    // Validate version integrity
    if err := v.validateVersion(version); err != nil {
        return err
    }
    
    // Restore agent file
    targetPath := filepath.Join(v.BaseDir, "..", "agents", agentName+".md")
    if err := v.restoreFile(version, targetPath); err != nil {
        return err
    }
    
    // Log rollback
    v.logRollback(agentName, versionID, "manual rollback")
    
    return nil
}
```

## Batch Processing System

```go
type BatchProcessor struct {
    Pipeline    *TransformationPipeline
    Validator   *ValidationSuite
    Version     *VersionManager
    WorkerPool  *WorkerPool
}

type BatchConfig struct {
    Size            int           // Agents per batch
    Parallelism     int           // Concurrent workers
    ValidationMode  ValidationMode
    RollbackPolicy  RollbackPolicy
    RetryAttempts   int
    RetryBackoff    time.Duration
}

func (b *BatchProcessor) ProcessBatch(agents []*AgentFile, config BatchConfig) (*BatchResult, error) {
    results := &BatchResult{
        StartTime: time.Now(),
        Agents:    make(map[string]*ProcessingResult),
    }
    
    // Create worker pool
    pool := NewWorkerPool(config.Parallelism)
    defer pool.Close()
    
    // Process agents in parallel
    var wg sync.WaitGroup
    resultChan := make(chan *ProcessingResult, len(agents))
    
    for _, agent := range agents {
        wg.Add(1)
        pool.Submit(func() {
            defer wg.Done()
            result := b.processAgent(agent, config)
            resultChan <- result
        })
    }
    
    // Collect results
    go func() {
        wg.Wait()
        close(resultChan)
    }()
    
    for result := range resultChan {
        results.Agents[result.AgentName] = result
    }
    
    // Validate batch as a whole
    if err := b.validateBatch(results); err != nil {
        if config.RollbackPolicy == RollbackOnBatchFailure {
            b.rollbackBatch(results)
        }
        return results, err
    }
    
    return results, nil
}

func (b *BatchProcessor) processAgent(agent *AgentFile, config BatchConfig) *ProcessingResult {
    result := &ProcessingResult{
        AgentName: agent.Metadata.Name,
        StartTime: time.Now(),
    }
    
    // Create version snapshot
    version, err := b.Version.CreateSnapshot(agent)
    if err != nil {
        result.Error = err
        return result
    }
    result.OriginalVersion = version.ID
    
    // Transform agent
    transformed, err := b.Pipeline.Transform(agent)
    if err != nil {
        result.Error = err
        return result
    }
    
    // Validate transformation
    validation := b.Validator.Validate(transformed)
    result.ValidationScore = validation.Score
    
    if validation.Score < config.MinQualityScore {
        result.Error = fmt.Errorf("quality score %.2f below threshold %.2f", 
            validation.Score, config.MinQualityScore)
        return result
    }
    
    // Save transformed version
    newVersion, err := b.Version.CreateSnapshot(transformed)
    if err != nil {
        result.Error = err
        return result
    }
    result.TransformedVersion = newVersion.ID
    
    result.Success = true
    result.EndTime = time.Now()
    return result
}
```

## Migration Orchestration

```go
type MigrationOrchestrator struct {
    Phases      []MigrationPhase
    Processor   *BatchProcessor
    Monitor     *MigrationMonitor
    Reporter    *ReportGenerator
}

type MigrationPhase struct {
    Name        string
    Agents      []string
    RiskLevel   RiskLevel
    BatchSize   int
    Validation  ValidationConfig
    Criteria    SuccessCriteria
}

func (m *MigrationOrchestrator) ExecutePhase(phase MigrationPhase) (*PhaseResult, error) {
    result := &PhaseResult{
        Phase:     phase.Name,
        StartTime: time.Now(),
    }
    
    // Load agents for this phase
    agents, err := m.loadPhaseAgents(phase.Agents)
    if err != nil {
        return nil, err
    }
    
    // Process in batches
    for i := 0; i < len(agents); i += phase.BatchSize {
        end := min(i+phase.BatchSize, len(agents))
        batch := agents[i:end]
        
        // Process batch with monitoring
        batchResult, err := m.processBatchWithMonitoring(batch, phase)
        if err != nil {
            // Check if we should continue or abort
            if phase.RiskLevel == RiskHigh {
                return result, fmt.Errorf("high-risk phase failed: %w", err)
            }
            // Log and continue for lower risk phases
            result.PartialFailures = append(result.PartialFailures, err)
        }
        
        result.BatchResults = append(result.BatchResults, batchResult)
        
        // Check phase gate criteria
        if !m.checkPhaseGate(result, phase.Criteria) {
            return result, fmt.Errorf("phase gate criteria not met")
        }
    }
    
    result.EndTime = time.Now()
    result.Success = true
    return result, nil
}

func (m *MigrationOrchestrator) checkPhaseGate(result *PhaseResult, criteria SuccessCriteria) bool {
    stats := m.calculatePhaseStats(result)
    
    return stats.SuccessRate >= criteria.MinSuccessRate &&
           stats.QualityScore >= criteria.MinQualityScore &&
           stats.RollbackRate <= criteria.MaxRollbackRate
}

type MigrationMonitor struct {
    Metrics     *MetricsCollector
    Alerting    *AlertManager
    Dashboard   *DashboardServer
}

func (m *MigrationMonitor) StartMonitoring(ctx context.Context) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            metrics := m.Metrics.CollectCurrent()
            m.Dashboard.Update(metrics)
            
            if alert := m.checkAlertConditions(metrics); alert != nil {
                m.Alerting.Send(alert)
            }
            
        case <-ctx.Done():
            return
        }
    }
}
```

## Testing Framework

```go
type RefactoringTestSuite struct {
    Fixtures    *TestFixtures
    Validators  []Validator
    Comparator  *SemanticComparator
}

func (s *RefactoringTestSuite) TestTransformationPreservesMeaning(t *testing.T) {
    // Load test fixture
    original := s.Fixtures.LoadAgent("test-agent-how-focused")
    
    // Transform
    pipeline := NewTransformationPipeline()
    transformed, err := pipeline.Transform(original)
    require.NoError(t, err)
    
    // Compare semantic meaning
    similarity := s.Comparator.Compare(original, transformed)
    assert.GreaterOrEqual(t, similarity, 0.85, "Semantic similarity too low")
    
    // Verify key concepts preserved
    originalConcepts := s.extractConcepts(original)
    transformedConcepts := s.extractConcepts(transformed)
    
    missingConcepts := difference(originalConcepts, transformedConcepts)
    assert.Empty(t, missingConcepts, "Critical concepts lost in transformation")
}

func (s *RefactoringTestSuite) TestBatchProcessingRollback(t *testing.T) {
    agents := s.Fixtures.LoadBatch("rollback-test-batch")
    processor := NewBatchProcessor()
    
    // Inject failure on third agent
    processor.Pipeline.InjectFailure(2)
    
    // Process batch
    result, err := processor.ProcessBatch(agents, BatchConfig{
        RollbackPolicy: RollbackOnBatchFailure,
    })
    
    // Verify rollback occurred
    assert.Error(t, err)
    assert.False(t, result.Success)
    
    // Verify all agents rolled back to original
    for _, agent := range agents {
        current := s.loadCurrentVersion(agent.Metadata.Name)
        assert.Equal(t, agent, current, "Agent not rolled back properly")
    }
}

func (s *RefactoringTestSuite) TestCrossReferenceIntegrity(t *testing.T) {
    // Load all agents in a category
    agents := s.Fixtures.LoadCategory("the-software-engineer")
    
    // Transform all agents
    transformed := make(map[string]*AgentFile)
    for _, agent := range agents {
        result, _ := s.transformAgent(agent)
        transformed[agent.Metadata.Name] = result
    }
    
    // Validate all cross-references
    validator := NewCrossReferenceValidator(transformed)
    results := validator.ValidateAll()
    
    for _, result := range results {
        assert.True(t, result.Passed, "Cross-reference validation failed: %v", result.Issues)
    }
}
```

## CLI Interface

```go
type RefactorCommand struct {
    Pipeline   *TransformationPipeline
    Validator  *ValidationSuite
    Orchestrator *MigrationOrchestrator
}

func (c *RefactorCommand) Execute(cmd *cobra.Command, args []string) error {
    // Parse flags
    phase, _ := cmd.Flags().GetString("phase")
    dryRun, _ := cmd.Flags().GetBool("dry-run")
    interactive, _ := cmd.Flags().GetBool("interactive")
    
    if interactive {
        return c.runInteractive()
    }
    
    // Load phase configuration
    phaseConfig, err := c.loadPhaseConfig(phase)
    if err != nil {
        return err
    }
    
    if dryRun {
        return c.runDryRun(phaseConfig)
    }
    
    // Execute migration
    result, err := c.Orchestrator.ExecutePhase(phaseConfig)
    if err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }
    
    // Generate report
    report := c.generateReport(result)
    fmt.Println(report)
    
    return nil
}

func (c *RefactorCommand) runInteractive() error {
    app := NewRefactoringTUI()
    return app.Run()
}

// TUI for interactive refactoring
type RefactoringTUI struct {
    agents     []*AgentFile
    selected   int
    preview    *TransformationPreview
    validator  *ValidationSuite
}

func (t *RefactoringTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            // Transform selected agent
            return t, t.transformAgent()
        case "p":
            // Preview transformation
            return t, t.previewTransformation()
        case "v":
            // Validate current state
            return t, t.validateAgent()
        case "r":
            // Rollback to previous version
            return t, t.rollbackAgent()
        }
    }
    return t, nil
}
```

## Performance Optimization

```go
type PerformanceOptimizer struct {
    Cache       *TransformationCache
    Parallelism int
    BatchSize   int
}

type TransformationCache struct {
    transformations sync.Map // map[hash]*TransformationResult
    analyses        sync.Map // map[hash]*AnalysisResult
}

func (c *TransformationCache) GetOrCompute(
    key string, 
    compute func() (*TransformationResult, error),
) (*TransformationResult, error) {
    if cached, ok := c.transformations.Load(key); ok {
        return cached.(*TransformationResult), nil
    }
    
    result, err := compute()
    if err != nil {
        return nil, err
    }
    
    c.transformations.Store(key, result)
    return result, nil
}

// Parallel processing with work stealing
type WorkStealingPool struct {
    workers   []*Worker
    globalQueue chan Task
    stealing  chan *Worker
}

func (p *WorkStealingPool) Submit(task Task) {
    select {
    case p.globalQueue <- task:
    default:
        // Try to steal from busy worker
        if worker := p.findIdleWorker(); worker != nil {
            worker.Submit(task)
        } else {
            p.globalQueue <- task // Block if necessary
        }
    }
}
```

## Error Recovery & Monitoring

```go
type ErrorRecovery struct {
    RetryPolicy    RetryPolicy
    CircuitBreaker *CircuitBreaker
    Fallback       FallbackStrategy
}

type CircuitBreaker struct {
    maxFailures  int
    timeout      time.Duration
    state        atomic.Value // closed, open, half-open
    failures     atomic.Int32
    lastFailTime atomic.Value
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if cb.isOpen() {
        return ErrCircuitOpen
    }
    
    err := fn()
    if err != nil {
        cb.recordFailure()
        if cb.failures.Load() >= int32(cb.maxFailures) {
            cb.open()
        }
        return err
    }
    
    cb.recordSuccess()
    return nil
}

type MonitoringCollector struct {
    prometheus.Collector
    
    // Metrics
    transformationDuration  *prometheus.HistogramVec
    validationScore        *prometheus.GaugeVec
    howWhatRatio          *prometheus.GaugeVec
    processingErrors      *prometheus.CounterVec
    rollbackCount         *prometheus.CounterVec
}

func (m *MonitoringCollector) RecordTransformation(agent string, duration time.Duration, score float64) {
    m.transformationDuration.WithLabelValues(agent).Observe(duration.Seconds())
    m.validationScore.WithLabelValues(agent).Set(score)
}
```

## Configuration Management

```yaml
# refactoring-config.yaml
transformation:
  rules:
    - name: convert_steps_to_objectives
      enabled: true
      priority: 1
      config:
        max_objectives: 5
        consolidate_similar: true
        
    - name: extract_success_criteria
      enabled: true
      priority: 2
      config:
        require_measurable: true
        min_criteria: 2
        
validation:
  quality_threshold: 0.85
  semantic_similarity_threshold: 0.80
  max_line_count: 45
  line_count_tolerance: 0.1
  
  validators:
    - what_ratio_validator:
        min_ratio: 0.70
        severity: error
        
    - semantic_preservation_validator:
        threshold: 0.85
        model: "text-embedding-ada-002"
        
    - cross_reference_validator:
        check_circular: true
        check_orphaned: true
        
migration:
  phases:
    - name: "pilot"
      risk_level: "low"
      agents: ["data-protection", "compliance-audit"]
      batch_size: 2
      validation_mode: "strict"
      success_criteria:
        min_success_rate: 1.0
        min_quality_score: 0.90
        
    - name: "main"
      risk_level: "medium"
      agents: ["category:the-analyst", "category:the-designer"]
      batch_size: 5
      validation_mode: "standard"
      success_criteria:
        min_success_rate: 0.95
        min_quality_score: 0.85
        
monitoring:
  metrics_port: 9090
  dashboard_port: 8080
  alert_webhook: "https://hooks.slack.com/services/..."
  
  alerts:
    - name: "high_failure_rate"
      condition: "failure_rate > 0.20"
      severity: "critical"
      
    - name: "quality_degradation"
      condition: "avg_quality_score < 0.80"
      severity: "warning"
```

This implementation specification provides the complete technical foundation for the agent refactoring system with production-ready code structures, comprehensive testing, and robust error handling.