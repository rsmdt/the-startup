# Agent Refactoring - Validation & Quality Assurance Framework

## Validation Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                  Validation Framework                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │            Static Analysis Layer                      │  │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐       │  │
│  │  │Syntax  │ │Schema  │ │Format  │ │Metrics │       │  │
│  │  └────────┘ └────────┘ └────────┘ └────────┘       │  │
│  └──────────────────────────────────────────────────────┘  │
│                           ↓                                  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │           Semantic Analysis Layer                     │  │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐       │  │
│  │  │Meaning │ │Coverage│ │Intent  │ │Context │       │  │
│  │  └────────┘ └────────┘ └────────┘ └────────┘       │  │
│  └──────────────────────────────────────────────────────┘  │
│                           ↓                                  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │           Integration Testing Layer                   │  │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐       │  │
│  │  │Cross   │ │Delega- │ │Category│ │System  │       │  │
│  │  │Refs    │ │tion    │ │Cohesion│ │Wide    │       │  │
│  │  └────────┘ └────────┘ └────────┘ └────────┘       │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

## Validation Stages

### Stage 1: Pre-Transformation Validation

```yaml
Purpose: Establish baseline and identify transformation candidates
Components:
  - Current State Analysis
  - Transformation Feasibility
  - Risk Assessment
  - Dependency Mapping

Validators:
  FileIntegrityValidator:
    - Valid YAML frontmatter
    - Proper markdown structure
    - No corruption or encoding issues
    
  BaselineMetricsValidator:
    - Current HOW/WHAT ratio
    - Line count and complexity
    - Cross-reference count
    
  TransformabilityValidator:
    - Sufficient content for transformation
    - Clear section boundaries
    - Identifiable patterns
```

### Stage 2: Transformation Validation

```yaml
Purpose: Validate transformations preserve meaning and improve quality
Components:
  - Rule Application Verification
  - Content Preservation Check
  - Structure Compliance
  - Quality Improvement Metrics

Validators:
  TransformationRuleValidator:
    - All applicable rules executed
    - No rule conflicts
    - Proper rule precedence
    
  ContentPreservationValidator:
    - No critical information lost
    - Domain terms maintained
    - Security requirements intact
    
  StructureComplianceValidator:
    - Matches target template
    - Required sections present
    - Proper hierarchy
```

### Stage 3: Post-Transformation Validation

```yaml
Purpose: Ensure transformed agents meet quality standards
Components:
  - Quality Metrics Verification
  - Functionality Preservation
  - Integration Testing
  - Regression Detection

Validators:
  QualityMetricsValidator:
    - WHAT ratio >= 70%
    - Line count <= 45 (±10%)
    - Clarity score improved
    
  FunctionalityValidator:
    - All capabilities preserved
    - No regression in coverage
    - Delegation paths intact
    
  IntegrationValidator:
    - Cross-references valid
    - Category coherence maintained
    - System-wide consistency
```

## Validation Rules & Criteria

### Content Validation Rules

```go
type ContentValidationRules struct {
    // Structural Rules
    RequiredSections []string {
        "Role",
        "Core Objectives", 
        "Success Criteria",
        "Boundaries",
    }
    
    // Content Rules
    RoleStatement struct {
        MaxSentences   int     // 2
        MaxWords       int     // 30
        MustContain    []string {"expertise", "purpose"}
    }
    
    CoreObjectives struct {
        MinCount       int     // 3
        MaxCount       int     // 5
        MustBeOutcome  bool    // true
        MustBeMeasurable bool  // true
    }
    
    SuccessCriteria struct {
        MinCount       int     // 2
        MustBeSpecific bool    // true
        MustBeTestable bool    // true
    }
    
    Boundaries struct {
        MustHaveInScope  bool   // true
        MustHaveOutScope bool   // true
        MustHaveDelegation bool // true
    }
}
```

### Semantic Validation Rules

```go
type SemanticValidationRules struct {
    // Meaning Preservation
    ConceptRetention struct {
        MinSimilarity     float64  // 0.85
        CriticalConcepts  []string // Domain-specific terms
        MustRetain        []string // Security, compliance terms
    }
    
    // Intent Preservation  
    IntentMapping struct {
        OriginalIntent    string
        TransformedIntent string
        Similarity        float64
        MustMatch         bool
    }
    
    // Context Preservation
    ContextRequirements struct {
        DomainKnowledge   []string
        TechnicalContext  []string
        BusinessContext   []string
    }
}
```

### Quality Validation Rules

```go
type QualityValidationRules struct {
    // Quantitative Metrics
    Metrics struct {
        MaxLineCount      int     // 45
        LineTolerance     float64 // 0.1 (10%)
        MinWhatRatio      float64 // 0.70
        MaxComplexity     float64 // 5.0
        MinClarity        float64 // 0.80
    }
    
    // Qualitative Metrics
    Readability struct {
        MaxSentenceLength int     // 25 words
        MinActiveVoice    float64 // 0.80
        MaxJargonDensity  float64 // 0.15
    }
    
    // Consistency Metrics
    Consistency struct {
        TerminologyConsistent bool
        StyleConsistent       bool
        ToneConsistent        bool
    }
}
```

## Validation Test Scenarios

### Scenario 1: HOW to WHAT Transformation

```go
func TestHowToWhatTransformation(t *testing.T) {
    // Given: HOW-focused content
    original := `
    ## Approach
    1. First, scan the codebase for SQL queries
    2. Then, check each query for injection vulnerabilities
    3. Next, validate input sanitization
    4. Finally, generate security report
    `
    
    // When: Transform to WHAT-focused
    transformed := transformer.Transform(original)
    
    // Then: Validate transformation
    expected := `
    ## Core Objectives
    - SQL injection vulnerability detection across codebase
    - Input validation and sanitization verification
    - Comprehensive security assessment reporting
    `
    
    // Assertions
    assert.Contains(t, transformed, "Core Objectives")
    assert.NotContains(t, transformed, "First")
    assert.NotContains(t, transformed, "Then")
    
    // Validate semantic preservation
    similarity := semantic.Compare(original, transformed)
    assert.GreaterOrEqual(t, similarity, 0.85)
    
    // Validate outcome focus
    outcomes := extractor.ExtractOutcomes(transformed)
    assert.Len(t, outcomes, 3)
    for _, outcome := range outcomes {
        assert.True(t, validator.IsOutcomeFocused(outcome))
    }
}
```

### Scenario 2: Cross-Reference Integrity

```go
func TestCrossReferenceIntegrity(t *testing.T) {
    // Given: Agents with delegations
    agents := map[string]*Agent{
        "security-audit": {
            Boundaries: Boundaries{
                OutOfScope: []string{
                    "Implementation → vulnerability-assessment",
                },
            },
        },
        "vulnerability-assessment": {
            Boundaries: Boundaries{
                InScope: []string{
                    "Security vulnerability scanning",
                },
            },
        },
    }
    
    // When: Transform all agents
    transformed := transformBatch(agents)
    
    // Then: Validate references remain valid
    validator := NewCrossReferenceValidator(transformed)
    
    for agentName, agent := range transformed {
        for _, delegation := range agent.GetDelegations() {
            targetAgent := delegation.Target
            
            // Verify target exists
            assert.Contains(t, transformed, targetAgent,
                "Agent %s delegates to non-existent %s", agentName, targetAgent)
            
            // Verify delegation makes sense
            assert.True(t, validator.IsDelegationValid(agentName, targetAgent),
                "Invalid delegation from %s to %s", agentName, targetAgent)
            
            // Verify no circular references
            assert.False(t, validator.HasCircularReference(agentName, targetAgent),
                "Circular reference detected: %s -> %s", agentName, targetAgent)
        }
    }
}
```

### Scenario 3: Quality Metrics Improvement

```go
func TestQualityMetricsImprovement(t *testing.T) {
    // Given: Original agent with poor metrics
    original := loadAgent("bloated-agent.md")
    originalMetrics := analyzer.Analyze(original)
    
    assert.Equal(t, 85, originalMetrics.LineCount)
    assert.Equal(t, 0.75, originalMetrics.HowRatio)
    assert.Equal(t, 3.2, originalMetrics.ClarityScore)
    
    // When: Transform agent
    transformed := pipeline.Transform(original)
    transformedMetrics := analyzer.Analyze(transformed)
    
    // Then: Validate improvements
    assert.LessOrEqual(t, transformedMetrics.LineCount, 50)
    assert.GreaterOrEqual(t, transformedMetrics.WhatRatio, 0.70)
    assert.GreaterOrEqual(t, transformedMetrics.ClarityScore, 4.0)
    
    // Ensure no functionality lost
    coverage := comparator.CompareFunctionality(original, transformed)
    assert.Equal(t, 1.0, coverage, "Functionality coverage reduced")
}
```

### Scenario 4: Batch Rollback Testing

```go
func TestBatchRollbackOnFailure(t *testing.T) {
    // Given: Batch of 5 agents
    batch := loadBatch("test-batch")
    originals := snapshotBatch(batch)
    
    // When: Process with injected failure
    processor := NewBatchProcessor()
    processor.InjectFailureAt(3) // Fail on 3rd agent
    
    result, err := processor.ProcessBatch(batch, BatchConfig{
        RollbackPolicy: RollbackOnAnyFailure,
    })
    
    // Then: Verify all agents rolled back
    assert.Error(t, err)
    assert.False(t, result.Success)
    
    for i, agent := range batch {
        current := loadAgent(agent.Name)
        assert.Equal(t, originals[i], current,
            "Agent %s not rolled back properly", agent.Name)
    }
    
    // Verify rollback logged
    logs := getRollbackLogs()
    assert.Contains(t, logs, "Batch rollback initiated")
    assert.Contains(t, logs, "5 agents rolled back successfully")
}
```

## Validation Metrics & Scoring

### Quality Score Calculation

```go
type QualityScorer struct {
    Weights QualityWeights
}

type QualityWeights struct {
    WhatRatio      float64 // 0.30
    LineCount      float64 // 0.20
    Clarity        float64 // 0.20
    Completeness   float64 // 0.15
    Consistency    float64 // 0.15
}

func (s *QualityScorer) CalculateScore(agent *Agent) float64 {
    scores := map[string]float64{
        "what_ratio":   s.scoreWhatRatio(agent),
        "line_count":   s.scoreLineCount(agent),
        "clarity":      s.scoreClarity(agent),
        "completeness": s.scoreCompleteness(agent),
        "consistency":  s.scoreConsistency(agent),
    }
    
    weightedScore := 0.0
    weightedScore += scores["what_ratio"] * s.Weights.WhatRatio
    weightedScore += scores["line_count"] * s.Weights.LineCount
    weightedScore += scores["clarity"] * s.Weights.Clarity
    weightedScore += scores["completeness"] * s.Weights.Completeness
    weightedScore += scores["consistency"] * s.Weights.Consistency
    
    return weightedScore
}

func (s *QualityScorer) scoreWhatRatio(agent *Agent) float64 {
    ratio := agent.Metrics.WhatRatio
    if ratio >= 0.70 {
        return 1.0
    } else if ratio >= 0.60 {
        return 0.8
    } else if ratio >= 0.50 {
        return 0.6
    }
    return 0.4
}

func (s *QualityScorer) scoreLineCount(agent *Agent) float64 {
    count := agent.Metrics.LineCount
    target := 45
    tolerance := 5
    
    if count <= target {
        return 1.0
    } else if count <= target+tolerance {
        return 0.9
    } else if count <= target+tolerance*2 {
        return 0.7
    }
    return 0.5
}
```

### Validation Report Generation

```go
type ValidationReport struct {
    Summary     ValidationSummary
    Details     []ValidationDetail
    Metrics     ValidationMetrics
    Issues      []ValidationIssue
    Suggestions []Suggestion
}

func GenerateValidationReport(agent *Agent, results []ValidationResult) *ValidationReport {
    report := &ValidationReport{
        Summary: ValidationSummary{
            AgentName:    agent.Name,
            OverallScore: calculateOverallScore(results),
            Status:       determineStatus(results),
            Timestamp:    time.Now(),
        },
    }
    
    // Aggregate validation details
    for _, result := range results {
        detail := ValidationDetail{
            ValidatorName: result.ValidatorName,
            Passed:       result.Passed,
            Score:        result.Score,
            Issues:       result.Issues,
        }
        report.Details = append(report.Details, detail)
    }
    
    // Calculate metrics
    report.Metrics = ValidationMetrics{
        WhatRatio:        agent.Metrics.WhatRatio,
        LineCount:        agent.Metrics.LineCount,
        ClarityScore:     agent.Metrics.ClarityScore,
        SemanticSimilarity: agent.Metrics.SemanticSimilarity,
        QualityScore:     agent.Metrics.QualityScore,
    }
    
    // Generate suggestions
    report.Suggestions = generateSuggestions(agent, results)
    
    return report
}

func (r *ValidationReport) ToMarkdown() string {
    var sb strings.Builder
    
    sb.WriteString(fmt.Sprintf("# Validation Report: %s\n\n", r.Summary.AgentName))
    sb.WriteString(fmt.Sprintf("**Status:** %s\n", r.Summary.Status))
    sb.WriteString(fmt.Sprintf("**Score:** %.2f/100\n", r.Summary.OverallScore))
    sb.WriteString(fmt.Sprintf("**Date:** %s\n\n", r.Summary.Timestamp.Format(time.RFC3339)))
    
    sb.WriteString("## Metrics\n\n")
    sb.WriteString(fmt.Sprintf("- WHAT Ratio: %.2f%%\n", r.Metrics.WhatRatio*100))
    sb.WriteString(fmt.Sprintf("- Line Count: %d/45\n", r.Metrics.LineCount))
    sb.WriteString(fmt.Sprintf("- Clarity Score: %.2f/5.0\n", r.Metrics.ClarityScore))
    sb.WriteString(fmt.Sprintf("- Semantic Similarity: %.2f%%\n", r.Metrics.SemanticSimilarity*100))
    
    if len(r.Issues) > 0 {
        sb.WriteString("\n## Issues\n\n")
        for _, issue := range r.Issues {
            sb.WriteString(fmt.Sprintf("- **%s:** %s\n", issue.Severity, issue.Message))
            if issue.FixHint != "" {
                sb.WriteString(fmt.Sprintf("  - Fix: %s\n", issue.FixHint))
            }
        }
    }
    
    if len(r.Suggestions) > 0 {
        sb.WriteString("\n## Suggestions\n\n")
        for _, suggestion := range r.Suggestions {
            sb.WriteString(fmt.Sprintf("- %s\n", suggestion.Text))
        }
    }
    
    return sb.String()
}
```

## Regression Testing Framework

### Test Suite Organization

```yaml
test_suites:
  unit_tests:
    - parser_tests
    - transformer_tests
    - validator_tests
    - metrics_tests
    
  integration_tests:
    - pipeline_tests
    - batch_processing_tests
    - rollback_tests
    - cross_reference_tests
    
  system_tests:
    - full_migration_tests
    - performance_tests
    - stress_tests
    - recovery_tests
    
  regression_tests:
    - golden_file_tests
    - backward_compatibility_tests
    - edge_case_tests
    - boundary_tests
```

### Golden File Testing

```go
type GoldenFileTest struct {
    InputPath    string
    ExpectedPath string
    Transformer  *Transformer
    Validator    *Validator
}

func (g *GoldenFileTest) Run(t *testing.T) {
    // Load input
    input, err := os.ReadFile(g.InputPath)
    require.NoError(t, err)
    
    // Transform
    result := g.Transformer.Transform(string(input))
    
    // Load expected output
    expected, err := os.ReadFile(g.ExpectedPath)
    require.NoError(t, err)
    
    // Compare
    if *updateGolden {
        // Update golden file if flag set
        err = os.WriteFile(g.ExpectedPath, []byte(result), 0644)
        require.NoError(t, err)
    } else {
        // Validate against golden file
        assert.Equal(t, string(expected), result,
            "Output doesn't match golden file for %s", g.InputPath)
    }
    
    // Additional validation
    validation := g.Validator.Validate(result)
    assert.True(t, validation.Passed,
        "Validation failed: %v", validation.Issues)
}
```

### Performance Benchmarks

```go
func BenchmarkTransformationPipeline(b *testing.B) {
    pipeline := NewTransformationPipeline()
    agent := loadTestAgent("benchmark-agent.md")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := pipeline.Transform(agent)
        if err != nil {
            b.Fatal(err)
        }
    }
    
    b.ReportMetric(float64(b.Elapsed())/float64(b.N), "ns/transformation")
}

func BenchmarkBatchProcessing(b *testing.B) {
    processor := NewBatchProcessor()
    batch := loadTestBatch("benchmark-batch", 10)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := processor.ProcessBatch(batch, DefaultConfig)
        if err != nil {
            b.Fatal(err)
        }
    }
    
    b.ReportMetric(float64(len(batch)), "agents/batch")
    b.ReportMetric(float64(b.Elapsed())/float64(b.N*len(batch)), "ns/agent")
}
```

## Continuous Validation Pipeline

```yaml
# .github/workflows/validation.yml
name: Agent Validation Pipeline

on:
  pull_request:
    paths:
      - 'assets/claude/agents/**/*.md'
      
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
          
      - name: Run Static Validation
        run: |
          go run ./cmd/validator static \
            --path assets/claude/agents \
            --config validation-config.yaml
            
      - name: Run Semantic Validation
        run: |
          go run ./cmd/validator semantic \
            --original ${{ github.base_ref }} \
            --modified ${{ github.head_ref }}
            
      - name: Run Integration Tests
        run: |
          go test ./tests/integration/... -v
          
      - name: Generate Report
        run: |
          go run ./cmd/validator report \
            --format markdown \
            --output validation-report.md
            
      - name: Comment on PR
        uses: actions/github-script@v6
        with:
          script: |
            const report = fs.readFileSync('validation-report.md', 'utf8');
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: report
            });
```

## Validation Dashboard

```go
type ValidationDashboard struct {
    server *http.Server
    data   *DashboardData
    mu     sync.RWMutex
}

type DashboardData struct {
    CurrentPhase    string
    AgentsProcessed int
    AgentsTotal     int
    SuccessRate     float64
    AverageQuality  float64
    Issues          []DashboardIssue
    Metrics         map[string]float64
}

func (d *ValidationDashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    d.mu.RLock()
    defer d.mu.RUnlock()
    
    tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Agent Refactoring Validation Dashboard</title>
        <meta http-equiv="refresh" content="5">
    </head>
    <body>
        <h1>Validation Dashboard</h1>
        
        <div class="progress">
            <h2>Progress</h2>
            <p>Phase: {{.CurrentPhase}}</p>
            <p>Agents: {{.AgentsProcessed}}/{{.AgentsTotal}}</p>
            <progress value="{{.AgentsProcessed}}" max="{{.AgentsTotal}}"></progress>
        </div>
        
        <div class="metrics">
            <h2>Quality Metrics</h2>
            <table>
                <tr><td>Success Rate:</td><td>{{.SuccessRate}}%</td></tr>
                <tr><td>Average Quality:</td><td>{{.AverageQuality}}/100</td></tr>
                <tr><td>WHAT Ratio:</td><td>{{index .Metrics "what_ratio"}}%</td></tr>
                <tr><td>Line Reduction:</td><td>{{index .Metrics "line_reduction"}}%</td></tr>
            </table>
        </div>
        
        <div class="issues">
            <h2>Recent Issues</h2>
            <ul>
            {{range .Issues}}
                <li>[{{.Severity}}] {{.Agent}}: {{.Message}}</li>
            {{end}}
            </ul>
        </div>
    </body>
    </html>
    `
    
    t := template.Must(template.New("dashboard").Parse(tmpl))
    t.Execute(w, d.data)
}
```

This validation framework provides comprehensive quality assurance for the agent refactoring process, ensuring that transformations maintain functionality while improving structure and clarity.