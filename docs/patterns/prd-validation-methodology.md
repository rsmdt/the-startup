# PRD Validation Methodology

## Overview

This document defines the three-phase validation methodology for Product Requirement Documents (PRDs) that ensures consistency, clarity, and completeness through structured frameworks.

## Validation Framework Components

### 1. SCQA-Guided Dialogue (Addresses Methodology Gap)

**Origin**: Barbara Minto's Pyramid Principle (McKinsey, 1960s)

**Purpose**: Provides structured dialogue framework to systematically validate PRD logic and uncover gaps, ambiguities, and inconsistencies.

#### SCQA Components

1. **Situation (S)**: Context and baseline - What is the current state?
2. **Complication (C)**: Problem and urgency - Why does this matter now?
3. **Question (Q)**: Core inquiry - What needs to be decided/solved?
4. **Answer (A)**: Proposed solution - How do we address this?

#### Application to PRD Validation

SCQA validates PRD documents by ensuring logical flow from context → problem → solution:

**Situation Validation Questions**:
- Is the current state clearly established?
- Are user contexts and personas well-defined?
- Is the baseline documented with evidence?
- Are all stakeholders aligned on the starting point?

**Complication Validation Questions**:
- Is the problem statement specific and measurable?
- Is the urgency or impact clearly articulated?
- Are we addressing root causes or symptoms?
- Is this problem validated by user research?
- Does the complication justify the proposed solution scope?

**Question Validation Questions**:
- Is the core question clearly framed?
- Do the proposed features align with the question?
- Is the scope appropriate to answer the question?
- Are we solving the right problem?

**Answer Validation Questions**:
- Do the features collectively address the question?
- Are acceptance criteria testable and complete?
- Are success metrics aligned with the solution?
- Is the solution feasible within constraints?
- Are edge cases and business rules comprehensive?

#### Benefits for PRD Quality

1. **Logical Consistency**: Validates vertical logic (S→C→Q→A flow) and horizontal logic (within each component)
2. **Gap Detection**: Identifies missing context, weak problems, misaligned requirements, incomplete solutions
3. **Structured Review**: Provides systematic framework for PRD reviewers
4. **Stakeholder Alignment**: Creates shared understanding of logic chain
5. **Quality Confidence**: Ensures documents are logically sound before implementation

---

### 2. Divergent-Convergent Cycles (Manages Exploration)

**Origin**: J.P. Guilford's cognitive research (1950s), refined by Design Thinking methodologies

**Purpose**: Manages thorough document exploration (divergent) while ensuring synthesis into clear validation results (convergent), preventing both premature conclusion and endless exploration.

#### Cycle Phases (Double Diamond Framework)

##### Phase 1: DISCOVER (Divergent)
**Goal**: Explore PRD comprehensively

**Activities**:
- Read document section by section
- Identify all claims, assumptions, and requirements
- Note questions, ambiguities, and inconsistencies
- Gather relevant context and background
- List all potential issues without filtering

**Duration**: 40% of validation time

**Exit Criteria**:
- Information saturation (no new findings)
- All sections explored
- Pattern recognition emerging
- Comprehensive issue list created

##### Phase 2: DEFINE (Convergent)
**Goal**: Synthesize findings and frame validation questions

**Activities**:
- Categorize discovered issues
- Prioritize findings by severity
- Frame specific validation questions
- Identify patterns in issues
- Define scope for deeper investigation

**Duration**: 20% of validation time

**Exit Criteria**:
- Clear validation questions defined
- Issues prioritized and categorized
- Investigation scope determined
- Handoff to next phase ready

##### Phase 3: DEVELOP (Divergent)
**Goal**: Investigate validation questions and gather evidence

**Activities**:
- Research each validation question thoroughly
- Gather supporting evidence
- Explore alternative solutions
- Test assumptions
- Document findings with evidence

**Duration**: 25% of validation time

**Exit Criteria**:
- All validation questions investigated
- Evidence gathered and documented
- Alternative approaches explored
- Sufficient data for conclusions

##### Phase 4: DELIVER (Convergent)
**Goal**: Reach conclusions and provide recommendations

**Activities**:
- Synthesize investigation findings
- Make validation determinations
- Create actionable recommendations
- Document rationale for decisions
- Prepare validation report

**Duration**: 15% of validation time

**Exit Criteria**:
- Clear validation conclusions
- Actionable recommendations provided
- Rationale documented
- Report ready for stakeholders

#### Benefits for Document Exploration

1. **Prevents Premature Conclusion**: Mandated exploration ensures thoroughness before evaluation
2. **Prevents Endless Exploration**: Bounded phases with time limits ensure timely conclusions
3. **Maintains Focus**: Explicit mode declarations keep validation on track
4. **Manages Cognitive Load**: Single-mode operation reduces mental switching costs
5. **Improves Quality**: Systematic process ensures thoroughness, objectivity, and traceability
6. **Natural Rhythm**: Alternating expansion/contraction mirrors effective problem-solving

#### Transition Triggers

**Switch to Convergent when**:
- Information saturation reached
- Time boundary hit
- Pattern recognition emerges
- Questions exhausted

**Switch to Divergent when**:
- Insufficient information discovered
- Invalid assumptions found
- New questions emerge from analysis
- Current solutions seem suboptimal

---

### 3. Automated MECE Validation (Builds Confidence)

**Origin**: McKinsey MECE Principle (Barbara Minto, 1960s)

**Purpose**: Ensures PRD completeness (Collectively Exhaustive) and non-redundancy (Mutually Exclusive) through automated validation, building confidence in document quality.

#### MECE Principles

**Mutually Exclusive (ME)**: No overlap between categories
- Each feature belongs to only one category
- No duplicate requirements
- No redundant acceptance criteria
- No overlapping personas

**Collectively Exhaustive (CE)**: Complete coverage
- All user needs addressed
- All success metrics tracked
- All personas have journeys
- All features have acceptance criteria

#### Automated Validation Approaches

##### 1. Structural Analysis (Rule-Based)

**Purpose**: Detect structural gaps and validate hierarchical organization

**Techniques**:
- Schema validation (all required sections present)
- Hierarchical consistency (parent-child relationships valid)
- Cross-reference validation (all references resolve)
- Enumeration completeness (categorical breakdowns complete)

**Example Detection**:
```
✗ Missing Section: "Success Metrics" (CE violation)
✗ Duplicate Section ID: "feature-001" (ME violation)
✗ Orphaned Reference: Feature "X" references non-existent persona "Y"
```

##### 2. Semantic Similarity Analysis (NLP-Based)

**Purpose**: Detect content overlaps and semantic redundancy

**Technology**: Sentence Transformers (e.g., Sentence-BERT)

**Technique**: Calculate cosine similarity between text embeddings

**Thresholds**:
- **> 0.95**: Critical overlap (near-duplicate content)
- **0.85-0.95**: High overlap (significant redundancy)
- **0.75-0.85**: Moderate overlap (review recommended)
- **< 0.75**: Acceptable similarity (related but distinct)

**Example Detection**:
```
⚠ Feature Overlap Detected:
  Feature A: "Allow users to export data to CSV"
  Feature B: "Enable CSV download functionality"
  Similarity: 0.92 (High)
  Recommendation: Merge or clarify distinction
```

##### 3. Coverage Analysis (Gap Detection)

**Purpose**: Identify missing requirements and incomplete specifications

**Techniques**:
- Traceability matrix analysis (persona → journey → feature → criteria)
- Cross-reference validation (metric → tracking event)
- Completeness checks (all required fields populated)

**Example Detection**:
```
✗ Coverage Gap:
  Persona: "Power User"
  Associated Journeys: [] (empty)
  Impact: User needs not captured
  Fix: Document at least one user journey
```

#### MECE Confidence Scoring

**Multi-Dimensional Confidence Score**:

```
Confidence = 0.25×Structural + 0.25×Semantic + 0.30×Coverage + 0.20×Traceability

Where:
- Structural Score: % of required sections present
- Semantic Score: 1 - (overlap violations / total comparisons)
- Coverage Score: 1 - (gap violations / total entities)
- Traceability Score: % of required links established
```

**Confidence Levels**:
- **95-100%**: Excellent (high confidence)
- **85-94%**: Good (acceptable confidence)
- **70-84%**: Fair (review recommended)
- **< 70%**: Poor (significant issues)

#### Common PRD MECE Violations

**Mutually Exclusive Violations**:
- Duplicate features (same functionality in multiple sections)
- Overlapping user stories (similar goals)
- Redundant acceptance criteria (testing same behavior)
- Persona overlap (similar characteristics/goals)
- Metric redundancy (measuring same behavior)

**Collectively Exhaustive Violations**:
- Missing user journeys (critical workflows not documented)
- Incomplete feature coverage (user needs without features)
- Edge case gaps (business rules don't cover all scenarios)
- Missing personas (user segments unrepresented)
- Tracking gaps (metrics without corresponding events)

---

## Integrated Validation Workflow

### Phase 1: SCQA Mapping
1. Map PRD sections to SCQA components
2. Validate logical flow (S→C→Q→A)
3. Identify structural issues using SCQA framework

### Phase 2: Divergent Exploration (Discover)
1. Read PRD comprehensively
2. Apply SCQA validation questions systematically
3. Document all findings without filtering
4. Use MECE structural analysis to detect gaps

### Phase 3: Convergent Synthesis (Define)
1. Categorize issues by SCQA component
2. Prioritize by severity (Critical/High/Medium/Low)
3. Frame specific validation questions
4. Run semantic similarity analysis for overlaps

### Phase 4: Divergent Investigation (Develop)
1. Research each validation question
2. Gather evidence for each issue
3. Run MECE coverage analysis
4. Explore alternative solutions

### Phase 5: Convergent Resolution (Deliver)
1. Calculate MECE confidence scores
2. Make validation determinations
3. Create actionable recommendations
4. Generate validation report with confidence metrics

---

## Implementation Recommendations

### For Manual Review (Human Validators)

1. **Use SCQA Checklist**: Apply SCQA validation questions section by section
2. **Follow Double Diamond**: Explicitly declare which phase (Discover/Define/Develop/Deliver)
3. **Time-Box Phases**: Allocate 40%/20%/25%/15% of validation time
4. **Document Findings**: Track issues in structured format for MECE analysis

### For Automated Validation (Tool Implementation)

1. **Structural Validation**: Implement schema checks and hierarchy validation
2. **Semantic Analysis**: Integrate Sentence Transformers for overlap detection
3. **Coverage Analysis**: Build traceability matrix and gap detection
4. **Confidence Scoring**: Calculate multi-dimensional confidence scores
5. **Reporting**: Generate actionable reports with prioritized recommendations

### Integration with PRD Template

Add validation instructions to PRD template:

```markdown
## Validation Instructions

This PRD should be validated using the three-phase methodology:

1. **SCQA Validation**: Ensure logical flow from Situation → Complication → Question → Answer
2. **Divergent-Convergent Cycles**: Follow Discover → Define → Develop → Deliver phases
3. **MECE Validation**: Check for completeness (CE) and non-redundancy (ME)

Target Confidence Score: 85%+ before implementation
```

---

## Success Metrics

**Validation Effectiveness**:
- **Adoption**: % of PRDs validated before implementation
- **Quality**: Reduction in implementation-phase requirement clarifications
- **Confidence**: Average confidence score of validated PRDs
- **Speed**: Time from PRD draft to approval with validation
- **Defect Prevention**: % reduction in PRD-related implementation issues

**Expected Impact**:
- 40-60% reduction in requirement clarification issues during implementation
- 80%+ confidence in PRD completeness before development begins
- 20-30% faster PRD approval process with systematic validation

---

## References

1. **SCQA Framework**: Minto, B. (1996). "The Pyramid Principle"
2. **Divergent-Convergent Cycles**: British Design Council (2004). "Double Diamond Framework"
3. **MECE Principle**: McKinsey & Company consulting methodology
4. **Semantic Similarity**: Reimers & Gurevych (2019). "Sentence-BERT"
5. **Requirements Engineering**: IEEE/ACM standards for requirements validation

---

## Related Documentation

- `docs/domain/prd-structure.md`: PRD domain model and business rules
- `docs/patterns/document-validation-patterns.md`: General document validation patterns
- `docs/interfaces/validation-tool-api.md`: API specification for validation tools