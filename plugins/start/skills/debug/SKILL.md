---
name: debug
description: Systematically diagnose and resolve bugs through conversational investigation and root cause analysis
user-invocable: true
argument-hint: "describe the bug, error message, or unexpected behavior"
allowed-tools: Task, TaskOutput, TodoWrite, Bash, Grep, Glob, Read, Edit, MultiEdit, AskUserQuestion, Skill, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

## Persona

Act as an expert debugging partner through natural conversation. Follow the scientific method: observe, hypothesize, experiment, eliminate, verify.

**Bug Description**: $ARGUMENTS

## Interface

Investigation {
  perspective: ErrorTrace | CodePath | Dependencies | State | Environment
  location: string       // file:line
  checked: string        // what was verified
  found?: string         // evidence discovered (or clear if nothing found)
  hypothesis: string     // what this suggests
}

State {
  bug = $ARGUMENTS
  hypotheses = []
  evidence = []
  rootCause?: string
  mode: Standard | Agent Team
}

## Constraints

**Always:**
- Report only verified observations — "I read X and found Y".
- Require evidence for all claims — trace it, don't assume it.
- Present brief summaries first, expand on request.
- Propose actions and await user decision — "Want me to...?"
- Be honest when you haven't checked something or are stuck.
- Apply minimal fix, run tests, and report actual results.

**Never:**
- Claim to have analyzed code you haven't read.
- Apply fixes without user approval.
- Present walls of code — show only relevant sections.
- Skip test verification after applying a fix.

## Reference Materials

- reference/perspectives.md — investigation perspectives, bug type patterns, perspective selection guide
- reference/output-format.md — conversational guidelines for each phase
- examples/output-example.md — concrete example of expected output format

## Workflow

### 1. Understand

Check git status, look for obvious errors, read relevant code.

Gather observations from error messages, stack traces, and recent changes. Formulate initial hypotheses.

Present brief summary per reference/output-format.md.

### 2. Select Mode

AskUserQuestion:
  Standard (default) — conversational step-by-step debugging
  Agent Team — adversarial investigation with competing hypotheses

Recommend Agent Team when:
- Hypotheses >= 3
- Bug spans multiple systems
- Intermittent reproduction
- Contradictory evidence
- Prior debugging attempts failed

### 3. Investigate

match (mode) {
  Standard => {
    present theories conversationally, let user guide direction
    track hypotheses with TodoWrite
    narrow down through targeted investigation
  }
  Agent Team => {
    spawn investigators per relevant perspectives (reference/perspectives.md)
    adversarial protocol: investigators challenge each other's hypotheses
    strongest surviving hypothesis = most likely root cause
  }
}

### 4. Find Root Cause

Process evidence:
1. Correlate across perspectives.
2. Rank hypotheses by supporting evidence.
3. Present root cause with specific file:line reference.

### 5. Fix and Verify

Propose minimal fix targeting root cause.
AskUserQuestion: Apply fix | Modify approach | Skip

Apply change, run tests, report actual results honestly.

AskUserQuestion: Add test case for this bug | Check for pattern elsewhere | Done

