# Constructive Feedback Patterns

Techniques for providing review feedback that improves both code quality and developer skills.

## The Feedback Formula

[Observation] + [Why it matters] + [Suggestion] + [Example if helpful]

## Good Feedback Examples

Instead of: "This is wrong"

Say: "This query runs inside a loop (line 45), which could cause N+1 performance issues as the dataset grows. Consider using a batch query before the loop:

```python
users = User.query.filter(User.id.in_(user_ids)).all()
user_map = {u.id: u for u in users}
```
"

Instead of: "Use better names"

Say: "The variable `d` on line 23 would be clearer as `daysSinceLastLogin` - it helps readers understand the business logic without tracing back to the assignment."

## Feedback Tone Guide

| Avoid | Prefer |
|-------|--------|
| "You should..." | "Consider..." or "What about..." |
| "This is wrong" | "This might cause issues because..." |
| "Why didn't you..." | "Have you considered..." |
| "Obviously..." | "One approach is..." |
| "Always/Never do X" | "In this context, X would help because..." |

## Positive Observations

Always include what's done well:

- "Nice use of the Strategy pattern here - it makes adding new payment methods straightforward."
- "Good error handling - the retry logic with exponential backoff is exactly what we need for this flaky API."
- "Clean separation of concerns between the validation and persistence logic."
