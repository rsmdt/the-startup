# Output Format Reference

Conversational templates for each debugging phase.

---

## Phase 1: Understand the Problem

```
"I see you're hitting [brief symptom summary]. Let me take a quick look..."

[Investigation results]

"Here's what I found so far: [1-2 sentence summary]

Want me to dig deeper, or can you tell me more about when this started?"
```

## Phase 2: Narrow It Down

```
"I have a couple of theories:
1. [Most likely] - because I saw [evidence]
2. [Alternative] - though this seems less likely

Want me to dig into the first one?"
```

## Phase 2b: Root Cause Found

```
"Found it. In [file:line], [describe what's wrong].

[Show only relevant code, not walls of text]

The problem: [one sentence explanation]

Should I fix this, or do you want to discuss the approach first?"
```

## Phase 3: Fix and Verify

### Propose Fix

```
"Here's what I'd change:

[Show the proposed fix - just the relevant diff]

This fixes it by [brief explanation].

Want me to apply this?"
```

### After Fix Applied

```
"Applied the fix. Tests are passing now. ✓

Can you verify on your end?"
```

## When Stuck

```
"I've looked at [what you checked] but haven't pinpointed it yet.

A few options:
- I could check [alternative area]
- You could tell me more about [specific question]
- We could take a different angle

What sounds most useful?"
```
