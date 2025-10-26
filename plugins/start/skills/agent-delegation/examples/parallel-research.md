# Example: Parallel Research Delegation

This example shows how to decompose a research task into parallel specialist activities.

## Scenario

**User Request:** "Research competitive landscape for our B2B SaaS pricing strategy"

## Task Decomposition

```
Original Task: Research competitive landscape for B2B SaaS pricing

Activities Identified:
1. Analyze Competitor A (Salesforce) pricing
   - Expertise: Market research
   - Output: Pricing analysis document
   - Dependencies: None

2. Analyze Competitor B (HubSpot) pricing
   - Expertise: Market research
   - Output: Pricing analysis document
   - Dependencies: None

3. Analyze Competitor C (Zendesk) pricing
   - Expertise: Market research
   - Output: Pricing analysis document
   - Dependencies: None

4. Synthesize competitive findings
   - Expertise: Strategic analysis
   - Output: Unified competitive strategy
   - Dependencies: Activities 1, 2, 3

Execution Strategy: Mixed
- Parallel: Activities 1-3 (independent research)
- Sequential: Activity 4 (requires all results)

Reasoning: Each competitor analysis is independent. Synthesis must wait for all analyses to complete.
```

## Agent Prompts Generated

### Agent 1: Salesforce Analysis

```
FOCUS: Research Salesforce's B2B SaaS pricing strategy
    - Identify all pricing tiers (Essentials, Professional, Enterprise, Unlimited)
    - Map features to each tier
    - Note promotional strategies and discounts
    - Calculate price-per-user for each tier
    - Identify what features drive tier upgrades

EXCLUDE:
    - Don't analyze their technology stack or implementation
    - Don't make pricing recommendations yet (synthesis phase)
    - Don't compare to other competitors (separate analyses)
    - Don't analyze Salesforce Platform (focus on CRM product)

CONTEXT: We're researching competitive landscape for our B2B SaaS pricing strategy.
    - Target market: Mid-market companies (100-500 employees)
    - Looking for pricing patterns and positioning insights
    - Focus on transparent, public pricing only

OUTPUT: Structured analysis document including:
    - Pricing Tiers Table (tier name, price, billing cycle)
    - Feature Matrix by Tier (features mapped to each tier)
    - Key Insights (3-5 bullet points about their strategy)
    - Notable Patterns (what we can learn from their approach)
    - Competitive Positioning (how they position each tier)

SUCCESS: Complete analysis with actionable data
    - All publicly available tiers documented
    - Features accurately mapped to tiers
    - Insights are specific and evidence-based
    - Patterns identified are relevant to our market

TERMINATION:
    - Analysis complete with all sections filled
    - Information not publicly available (document what's missing)
    - Maximum 2 hours of research time reached
```

### Agent 2: HubSpot Analysis

```
FOCUS: Research HubSpot's B2B SaaS pricing strategy
    - Identify all pricing tiers (Starter, Professional, Enterprise)
    - Map features to each tier
    - Note freemium offering and upgrade triggers
    - Calculate price-per-contact for marketing hub
    - Identify bundling strategy across products

EXCLUDE:
    - Don't analyze technology implementation
    - Don't make pricing recommendations yet
    - Don't compare to other competitors
    - Don't analyze HubSpot Academy or services pricing

CONTEXT: Competitive landscape research for B2B SaaS pricing.
    - Target market: Mid-market companies (100-500 employees)
    - Interest in freemium-to-paid conversion strategies
    - Focus on Marketing Hub and Sales Hub pricing

OUTPUT: Structured analysis document including:
    - Pricing Tiers Table
    - Feature Matrix by Tier
    - Freemium Analysis (what's free, what triggers upgrade)
    - Key Insights (3-5 bullet points)
    - Bundling Strategy (how products bundle together)
    - Notable Patterns

SUCCESS: Complete analysis with freemium insights
    - All tiers documented
    - Freemium model understood
    - Bundle pricing analyzed
    - Insights are actionable

TERMINATION:
    - Analysis complete
    - Information unavailable
    - 2 hours research time reached
```

### Agent 3: Zendesk Analysis

```
FOCUS: Research Zendesk's B2B SaaS pricing strategy
    - Identify all pricing tiers (Suite Team, Growth, Professional, Enterprise)
    - Map features to each tier
    - Note per-agent pricing model
    - Analyze add-on pricing strategy
    - Identify what drives customer to higher tiers

EXCLUDE:
    - Don't analyze implementation details
    - Don't make recommendations yet
    - Don't compare to competitors
    - Don't analyze Zendesk Sell (focus on Support product)

CONTEXT: Competitive research for B2B SaaS pricing strategy.
    - Target market: Mid-market companies
    - Interest in per-seat vs per-usage pricing models
    - Focus on customer support software pricing

OUTPUT: Structured analysis including:
    - Pricing Tiers Table
    - Feature Matrix by Tier
    - Per-Agent Economics (cost per agent at scale)
    - Add-On Pricing (what extras cost money)
    - Key Insights (3-5 bullet points)
    - Upgrade Drivers (what makes customers move up tiers)

SUCCESS: Complete analysis with per-seat insights
    - All tiers documented
    - Per-agent model understood
    - Add-on strategy mapped
    - Insights are specific

TERMINATION:
    - Analysis complete
    - Information unavailable
    - 2 hours research reached
```

### Agent 4: Synthesis (Sequential, after 1-3)

```
FOCUS: Synthesize competitive pricing findings into unified strategy insights

EXCLUDE:
    - Don't create our pricing (separate decision)
    - Don't add new competitor research
    - Don't analyze product features (stick to pricing patterns)

CONTEXT: We have three completed competitor analyses.
    - Salesforce analysis: [Results from Agent 1]
    - HubSpot analysis: [Results from Agent 2]
    - Zendesk analysis: [Results from Agent 3]
    - Goal: Extract patterns and strategic insights

OUTPUT: Competitive Pricing Strategy Document including:
    - Cross-Competitor Patterns (themes across all three)
    - Pricing Model Comparison (per-user vs per-contact vs per-usage)
    - Tier Strategy Analysis (how competitors structure tiers)
    - Freemium Insights (who uses it, how it works)
    - Feature Differentiation (what features separate tiers)
    - Positioning Opportunities (gaps or opportunities for us)
    - Strategic Recommendations (3-5 actionable insights)

SUCCESS: Unified synthesis with strategic value
    - Patterns identified across competitors
    - Comparisons are fair and evidence-based
    - Recommendations are actionable
    - No contradiction between analyses

TERMINATION:
    - Synthesis complete
    - Conflicting data needs clarification (escalate)
    - All sections filled
```

## Execution Flow

### Phase 1: Launch Parallel Research (Agents 1-3)

```
🚀 Launching 3 parallel research agents

Agent 1: Salesforce Analysis → RUNNING
Agent 2: HubSpot Analysis → RUNNING
Agent 3: Zendesk Analysis → RUNNING

File Coordination Check:
✅ Agent 1 OUTPUT: research/salesforce-pricing.md
✅ Agent 2 OUTPUT: research/hubspot-pricing.md
✅ Agent 3 OUTPUT: zendesk-pricing.md
✅ All paths unique - No collision risk

Parallel Safety Check:
✅ Independent tasks - No dependencies
✅ Read-only research - No shared state
✅ Separate validation - Each can be reviewed independently
✅ Won't block - Different data sources

Status: ✅ SAFE FOR PARALLEL EXECUTION
```

### Phase 2: Await Completion

```
⏳ Waiting for parallel agents...

Agent 1: COMPLETE ✅ (45 minutes)
  - Salesforce analysis delivered
  - All sections complete
  - 4 tiers documented

Agent 2: COMPLETE ✅ (52 minutes)
  - HubSpot analysis delivered
  - Freemium model analyzed
  - 3 tiers + free tier documented

Agent 3: COMPLETE ✅ (38 minutes)
  - Zendesk analysis delivered
  - Per-agent model mapped
  - 4 tiers + add-ons documented

All parallel agents complete. Proceeding to synthesis.
```

### Phase 3: Launch Synthesis (Agent 4)

```
🚀 Launching synthesis agent

Agent 4: Competitive Synthesis → RUNNING

Context provided:
- Salesforce results (from Agent 1)
- HubSpot results (from Agent 2)
- Zendesk results (from Agent 3)

Agent 4: COMPLETE ✅ (25 minutes)
  - Cross-competitor patterns identified
  - 5 strategic recommendations generated
  - Positioning opportunities highlighted
```

## Results

### Total Time: 52 minutes (parallel) + 25 minutes (synthesis) = 77 minutes

**Compare to sequential:** 45 + 52 + 38 + 25 = 160 minutes
**Time saved:** 83 minutes (52% faster)

### Deliverables

```
📁 research/
├── salesforce-pricing.md (Agent 1)
├── hubspot-pricing.md (Agent 2)
├── zendesk-pricing.md (Agent 3)
└── competitive-strategy.md (Agent 4 synthesis)
```

### Key Insights Generated

From the synthesis agent:

1. **Tiering Pattern:** All three use 3-4 tier structure with similar progression (basic → professional → enterprise)

2. **Pricing Models:** Mixed approaches
   - Salesforce: Per-user, all-inclusive features
   - HubSpot: Per-contact, freemium base
   - Zendesk: Per-agent, add-on marketplace

3. **Feature Gating:** Core features in all tiers, advanced analytics/automation in top tiers

4. **Freemium:** Only HubSpot uses freemium successfully (strong upgrade triggers identified)

5. **Opportunity:** Gap in mid-market transparent pricing - competitors hide "contact sales" behind top tier

## Lessons Learned

### What Worked Well

✅ **Parallel execution:** Saved 52% time
✅ **Independent research:** No coordination overhead
✅ **Synthesis phase:** Unified findings effectively
✅ **Unique file paths:** No collisions
✅ **Explicit FOCUS/EXCLUDE:** Agents stayed on task

### Improvements for Next Time

- Add time limits to prevent research rabbit holes
- Specify exact format (all agents used slightly different table formats)
- Request specific pricing data points (some agents missed cost-per-user calculations)
- Consider adding validation agent before synthesis (check data accuracy)

## Reusable Template

This pattern works for any parallel research:

```
1. Decompose research into independent topics
2. Create identical FOCUS/EXCLUDE templates
3. Customize context and output paths only
4. Launch all in parallel
5. Synthesis agent consolidates findings
```

**Use when:**
- Researching multiple competitors
- Analyzing multiple technologies
- Gathering multiple data sources
- Interviewing multiple stakeholders
