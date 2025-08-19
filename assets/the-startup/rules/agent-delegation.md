**You MUST ALWAYS follow the below rules for decomposing tasks and delegation to sub-agents, with focus on parallel execution, validation, and response handling.**

ALWAYS decompose and delegate to sub-agents proactively to enhance efficiency and maintain quality.

You MUST execute tasks in parallel when they are independent of each other, require different expertise, and can be validated separately.

**Core Principles:**
- Clear boundaries: Specify what's in and out of scope
- Minimal context: Pass only what's needed for the task
- Validate responses: Check for drift before proceeding
- Track progress: Use TodoWrite throughout delegation
- DELEGATE PROACTIVELY: When in doubt, delegate to specialists

**Task decomposition before delegating to sub-agents:**
- Break work into independent, verifiable units with clear inputs/outputs
- Split by expertise, interfaces, or data vs. code vs. validation
- Identify shared prerequisites once; avoid duplication
- Assign one owner per task; if overlap, define a single source of truth
- If decomposition creates heavy cross-talk, merge or run sequentially

**Parallel Agent Execution:**
- Run sub-agent tasks in parallel if independent, needing different expertise, or separable for validation
- Steps: mark tasks as in_progress → assign unique AgentIDs → launch simultaneously → validate each independently → mark completed
- Handle conflicts between results, then consolidate

**Context management:**
- Always include requirements, constraints, success criteria, dependencies
- Explicitly exclude non-requested features, future phases, or other agents’ responsibilities
- AgentIDs use format {agent}-{shortid}, e.g. the-architect-3xy87q

**Response Handling:**
- Always display agent personality and `<commentary>`:
  ```
  <commentary>
  [Agent's personality/thoughts]
  </commentary>
  
  ---
  
  [Response content]
  ```

- Parallel Responses must remain separate (never merge or summarize):
  ```
  <commentary>
  [Agent 1 personality/thoughts]
  <commentary>
  
  <commentary>
  [Agent 2 personality/thoughts]
  <commentary>
  ```

- When agents return `<tasks>`, extract them, confirm with the user, the add approved items to TodoWrite

**Validation & Drift Detection:**
- For each response check: scope adherence, complexity, and pass/drift
- Auto-accept: security fixes, error handling, docs, validation
- Ask user (minor drift): helpful extras, better patterns, more tests
- Require approval (major drift): new features, DB changes, external integrations, performance optimizations
- For drift, offer options to user. If user rejects, re-run agent(s) with tighter boundaries

**Error Recovery:**
- If agent is blocked, provide options: retry with clarifications, skip, or reassign
- Strategies: retry with stricter context, reassign to another specialist, or mark blocked and continue

**TodoWrite Integration:**
- Before delegation: add task
- Start: mark `in_progress`
- After validation: mark `completed` or leave `in_progress` if blocked
- Always update immediately (don’t batch)
- For parallel runs, update individual results

**Phase Transitions:**
- After each phase, provide a short summary of outcomes and ask to continue:
  ```
  📄 [Phase] complete:
  [main changes / key outcomes]
  ```

**Best Practices:**
- Do: run independent tasks or sub-agents in parallel, validate everything, keep commentary visible, update TodoWrite promptly, be explicit about exclusions
- Don’t: merge responses, skip validation, pass excess context, allow unchecked drift, forget tracking

**Remember:** These rules provides patterns to enhance your natural delegation abilities. Apply them when they add value, using your judgment for the specific situation.
