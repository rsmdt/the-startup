Use this format for the-chief invocation:

```
Task: the-chief
FOCUS: Assess complexity and recommend approach
EXCLUDE: Implementation details, specific technical solutions
CONTEXT: [relevant analysis results, requirements, or code assessment]
SUCCESS: Clear complexity score and execution strategy recommendation
```

**DISPLAY PROTOCOL**

Display the-chief's response VERBATIM in delimiters:
```
=== Response from the-chief ===
[COMPLETE UNMODIFIED RESPONSE]
=== End of the-chief response ===
```

After displaying, present routing decision based on complexity score.
STOP and wait for user confirmation before proceeding.

**ROUTING BASED ON COMPLEXITY**

For Specifications:
- Low (1-3): Minimal documentation
- Medium (4-6): Standard documentation  
- High (7-10): Comprehensive documentation

For Refactoring:
- Simple (1-3): Direct execution
- Complex (4+): Create plan for /s:implement