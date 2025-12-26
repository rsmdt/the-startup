---
name: exploratory-testing
description: Test charter structure, SFDPOT and FEW HICCUPPS heuristics, edge case discovery patterns, and session-based test management. Use when conducting exploratory testing, discovering edge cases, validating user journeys, or finding issues automated tests miss.
---

# Exploratory Testing Heuristics

Systematic heuristics and techniques for discovering defects through creative exploration that complements automated testing.

## When to Activate

- Testing new features manually
- Discovering edge cases and boundary conditions
- Validating user journeys end-to-end
- Finding issues automated tests miss
- Evaluating usability and user experience
- Security probing and input validation testing

## Test Charters

A charter provides focus and structure for exploration.

### Charter Template

```
┌─────────────────────────────────────────────────────────────┐
│ CHARTER: [Short descriptive title]                          │
├─────────────────────────────────────────────────────────────┤
│ TARGET:                                                      │
│ [Feature, component, or area under test]                    │
│                                                             │
│ DURATION:                                                    │
│ [Time box: 30, 60, 90 minutes]                              │
│                                                             │
│ MISSION:                                                     │
│ Explore [target]                                            │
│ With [resources, techniques, data]                          │
│ To discover [types of issues, risks]                        │
├─────────────────────────────────────────────────────────────┤
│ SESSION NOTES:                                               │
│ - [Timestamp] [Observation]                                 │
│ - [Timestamp] [Bug found]                                   │
│ - [Timestamp] [Question raised]                             │
├─────────────────────────────────────────────────────────────┤
│ BUGS FOUND:                                                  │
│ [ ] BUG-001: [Title]                                        │
│ [ ] BUG-002: [Title]                                        │
├─────────────────────────────────────────────────────────────┤
│ QUESTIONS/CONCERNS:                                          │
│ - [Question needing clarification]                          │
│                                                             │
│ COVERAGE NOTES:                                              │
│ - Tested: [What was covered]                                │
│ - Not tested: [What remains]                                │
└─────────────────────────────────────────────────────────────┘
```

### Example Charters

```
CHARTER: Checkout Edge Cases
TARGET: Shopping cart and checkout flow
DURATION: 60 minutes
MISSION:
  Explore checkout with unusual cart combinations
  With empty carts, single items, max items, mixed currencies
  To discover payment failures and data integrity issues

---

CHARTER: Mobile Login Stress
TARGET: Authentication on mobile devices
DURATION: 45 minutes
MISSION:
  Explore login/logout under poor network conditions
  With airplane mode, 2G simulation, wifi switching
  To discover session handling and offline behavior issues

---

CHARTER: Search Manipulation
TARGET: Product search functionality
DURATION: 30 minutes
MISSION:
  Explore search with special characters and edge inputs
  With SQL injection, XSS payloads, Unicode, emojis
  To discover security vulnerabilities and input handling issues
```

## SFDPOT Heuristic

A mnemonic for comprehensive coverage across six dimensions.

### The Framework

| Letter | Dimension | Focus |
|--------|-----------|-------|
| **S** | Structure | Data and code organization |
| **F** | Function | What the product does |
| **D** | Data | Information flowing through |
| **P** | Platform | Environment dependencies |
| **O** | Operations | How it's used in practice |
| **T** | Time | Temporal aspects |

### Detailed Questions

#### S - Structure

```
Questions:
- What are the components and how do they connect?
- What does the database schema look like?
- What are the file/folder structures?
- What configuration options exist?
- What are the integration points?

Exploration:
- Map the system architecture
- Trace data through components
- Test component isolation
- Verify configuration handling
```

#### F - Function

```
Questions:
- What are all the features?
- What are the inputs and outputs?
- What transformations occur?
- What are the business rules?
- What calculations are performed?

Exploration:
- Test each feature systematically
- Verify all stated functionality
- Check feature interactions
- Validate business logic
```

#### D - Data

```
Questions:
- What data types are handled?
- What are the valid ranges?
- How is data validated?
- How is data stored and retrieved?
- What happens with invalid data?

Exploration:
- Boundary value testing
- Invalid input testing
- Data type mismatches
- Empty/null/missing data
- Large data volumes
```

#### P - Platform

```
Questions:
- What operating systems?
- What browsers and versions?
- What devices and screen sizes?
- What network conditions?
- What dependencies exist?

Exploration:
- Cross-browser testing
- Mobile device testing
- Different OS versions
- Network throttling
- Dependency version testing
```

#### O - Operations

```
Questions:
- How is it installed/deployed?
- How is it configured?
- How is it maintained?
- How is it monitored?
- How is it backed up?

Exploration:
- Installation procedures
- Upgrade/downgrade paths
- Backup and recovery
- Configuration changes
- Log and monitoring verification
```

#### T - Time

```
Questions:
- What happens over time?
- Are there timeouts?
- What about scheduling?
- How does concurrency work?
- What about time zones?

Exploration:
- Long-running sessions
- Timeout behaviors
- Scheduled tasks
- Race conditions
- Time zone handling
- Daylight saving transitions
```

## FEW HICCUPPS Heuristic

Consistency oracles for judging whether behavior is correct.

### The Framework

Each letter represents a source of expectations:

| Letter | Oracle | Question |
|--------|--------|----------|
| **F** | Familiarity | Is it consistent with similar products? |
| **E** | Explainability | Can we explain this behavior to users? |
| **W** | World | Is it consistent with the real world? |
| **H** | History | Is it consistent with past versions? |
| **I** | Image | Is it consistent with the brand/reputation? |
| **C** | Comparable | Is it consistent with competitors? |
| **C** | Claims | Is it consistent with documentation? |
| **U** | Users | Is it consistent with user expectations? |
| **P** | Purpose | Is it consistent with stated goals? |
| **P** | Product | Is it internally consistent? |
| **S** | Statutes | Is it consistent with regulations? |

### Using the Oracles

```
For each observed behavior, ask:

FAMILIARITY
"I've used similar products. This feels wrong because..."

EXPLAINABILITY
"If a user asked why it works this way, I couldn't explain because..."

WORLD
"In the real world, this doesn't make sense because..."

HISTORY
"The previous version didn't work this way. The change seems..."

IMAGE
"This doesn't match what I'd expect from a quality product because..."

COMPARABLE
"Competitor X handles this by... Our approach seems..."

CLAIMS
"The documentation says X, but the behavior is Y..."

USERS
"Users would expect... but instead they get..."

PURPOSE
"The stated goal is X, but this feature undermines it by..."

PRODUCT
"This part works like X, but that part works like Y. Inconsistent..."

STATUTES
"This might violate [regulation/law] because..."
```

## Edge Case Patterns

### Boundary Value Testing

```
INPUT BOUNDARIES:

For numeric inputs:
┌────────────────────────────────────────────────┐
│ Value  │ Type                                  │
├────────────────────────────────────────────────┤
│ -1     │ Below minimum                         │
│ 0      │ Minimum boundary                      │
│ 1      │ Just above minimum                    │
│ ...    │ Typical values                        │
│ max-1  │ Just below maximum                    │
│ max    │ Maximum boundary                      │
│ max+1  │ Above maximum                         │
└────────────────────────────────────────────────┘

For string inputs:
- Empty string: ""
- Single character: "a"
- Maximum length
- Maximum length + 1
- Only whitespace: "   "
- Leading/trailing whitespace: " test "
```

### Special Input Values

```
EMPTY/NULL/MISSING:
- null
- undefined
- Empty string ""
- Empty array []
- Empty object {}
- Whitespace only
- Missing required fields

NUMERIC EDGE CASES:
- 0, -0
- Negative numbers
- Decimals (0.1, 0.01)
- Scientific notation (1e10)
- MAX_INT, MIN_INT
- MAX_INT + 1 (overflow)
- Infinity, -Infinity
- NaN

STRING EDGE CASES:
- Unicode: "Ṁöñᵗý Ṗẏẗḧöñ"
- Emojis: "😀🎉🔥"
- RTL text: "مرحبا"
- Null byte: "test\x00injection"
- Very long strings (10KB+)
- Format strings: "%s %n %d"
- Path traversal: "../../../etc/passwd"

DATE/TIME EDGE CASES:
- Leap years (Feb 29)
- Daylight saving transitions
- End of month (Jan 31 → Feb 28)
- End of year (Dec 31 → Jan 1)
- Time zone boundaries
- Unix epoch (1970-01-01)
- Y2K-style dates (2038 problem)
- Far future dates
```

### State Transition Testing

```
Test all valid transitions AND invalid ones:

ORDER STATUS:
┌────────┐    ┌───────────┐    ┌────────┐
│ Draft  │───►│ Confirmed │───►│Shipped │
└────────┘    └───────────┘    └────────┘
     │              │              │
     ▼              ▼              ▼
┌────────┐    ┌───────────┐    ┌────────┐
│Cancelled│   │ Cancelled │    │Delivered│
└────────┘    └───────────┘    └────────┘

Test:
✓ Valid: Draft → Confirmed → Shipped → Delivered
✓ Valid: Confirmed → Cancelled
✗ Invalid: Shipped → Draft (should fail)
✗ Invalid: Delivered → Confirmed (should fail)
```

### Error Condition Testing

```
NETWORK ERRORS:
- Connection timeout
- Connection refused
- DNS failure
- SSL certificate errors
- Partial response
- Corrupted response

RESOURCE ERRORS:
- Disk full
- Out of memory
- Too many open files
- Permission denied
- File locked

CONCURRENCY ERRORS:
- Race conditions
- Deadlocks
- Lost updates
- Double submissions
```

## Session-Based Test Management

### Session Metrics

```
SESSION REPORT:
┌─────────────────────────────────────────────────────────────┐
│ Tester: [Name]                                              │
│ Charter: [Title]                                            │
│ Date: [Date]           Duration: [Actual time]              │
├─────────────────────────────────────────────────────────────┤
│ TIME BREAKDOWN:                                              │
│ ┌────────────────────────────────────────────────────────┐ │
│ │ Testing:        ████████████████████ 60%              │ │
│ │ Bug Investigation: ██████████ 25%                      │ │
│ │ Setup/Overhead:    █████ 15%                           │ │
│ └────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ OUTCOMES:                                                    │
│ Bugs filed: [count]                                         │
│ Questions raised: [count]                                   │
│ Test ideas generated: [count]                               │
│ Coverage areas: [list]                                      │
├─────────────────────────────────────────────────────────────┤
│ FOLLOW-UP:                                                   │
│ - [Next charter suggestion]                                 │
│ - [Area needing more coverage]                              │
└─────────────────────────────────────────────────────────────┘
```

### Coverage Visualization

```
FEATURE COVERAGE MAP:

Feature                Testing Coverage
─────────────────────────────────────────
Login                  ████████████████ 100%
Registration           ████████████░░░░  75%
Search                 ████████░░░░░░░░  50%
Checkout               ████████████████ 100%
Profile                ████░░░░░░░░░░░░  25%
Admin                  ░░░░░░░░░░░░░░░░   0%

Legend: █ Tested  ░ Not tested
```

## Bug Reporting

### Bug Report Template

```markdown
# BUG-[ID]: [Descriptive Title]

## Summary
[One sentence description]

## Environment
- Browser/Device: [e.g., Chrome 120 on macOS]
- URL: [where it occurred]
- User: [role/account type]

## Steps to Reproduce
1. [First step]
2. [Second step]
3. [Third step]

## Expected Result
[What should happen]

## Actual Result
[What actually happened]

## Evidence
- Screenshot: [attached]
- Video: [link]
- Console log: [attached]

## Severity
[Critical / High / Medium / Low]

## Notes
[Additional context, related issues, workarounds]
```

### Severity Guidelines

| Level | Definition | Example |
|-------|------------|---------|
| **Critical** | System unusable, data loss | Payment double-charging |
| **High** | Major feature broken | Cannot complete checkout |
| **Medium** | Feature impaired, workaround exists | Search requires exact match |
| **Low** | Minor issue, cosmetic | Alignment off by 2px |

## Best Practices

### Do

- **Time-box sessions** - 60-90 minutes max before break
- **Take notes continuously** - Capture observations as you go
- **Follow your curiosity** - Interesting behavior often hides bugs
- **Vary your approach** - Don't test the same way every time
- **Question everything** - "Is this right?" should be constant
- **Document coverage** - Know what you've tested and what remains

### Avoid

- **Random clicking** - Exploration should be purposeful
- **Only happy paths** - Edge cases hide the interesting bugs
- **Ignoring instincts** - "That's weird" often leads to bugs
- **Testing in isolation** - Consider the full user journey
- **Rushing** - Quality exploration takes focused time

## References

- [Charter Templates](examples/charters.md) - Ready-to-use charter templates
- [Heuristic Quick Reference](reference.md) - Printable heuristic cards
