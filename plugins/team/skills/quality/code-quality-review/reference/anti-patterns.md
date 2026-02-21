# Anti-Pattern Catalog

Common code smells organized by scope, with detection signs and remediation strategies.

## Method-Level Anti-Patterns

| Anti-Pattern | Detection Signs | Remediation |
|--------------|-----------------|-------------|
| **Long Method** | >20 lines, multiple responsibilities | Extract Method |
| **Long Parameter List** | >3-4 parameters | Introduce Parameter Object |
| **Duplicate Code** | Copy-paste patterns | Extract Method, Template Method |
| **Complex Conditionals** | Nested if/else, switch statements | Decompose Conditional, Strategy Pattern |
| **Magic Numbers** | Hardcoded values without context | Extract Constant |
| **Dead Code** | Unreachable or unused code | Delete it |

## Class-Level Anti-Patterns

| Anti-Pattern | Detection Signs | Remediation |
|--------------|-----------------|-------------|
| **God Object** | >500 lines, many responsibilities | Extract Class |
| **Data Class** | Only getters/setters, no behavior | Move behavior to class |
| **Feature Envy** | Method uses another class's data extensively | Move Method |
| **Inappropriate Intimacy** | Classes know too much about each other | Move Method, Extract Class |
| **Refused Bequest** | Subclass doesn't use inherited behavior | Replace Inheritance with Delegation |
| **Lazy Class** | Does too little to justify existence | Inline Class |

## Architecture-Level Anti-Patterns

| Anti-Pattern | Detection Signs | Remediation |
|--------------|-----------------|-------------|
| **Circular Dependencies** | A depends on B depends on A | Dependency Inversion |
| **Shotgun Surgery** | One change requires many file edits | Move Method, Extract Class |
| **Leaky Abstraction** | Implementation details exposed | Encapsulate |
| **Premature Optimization** | Complex code for unproven performance | Simplify, measure first |
| **Over-Engineering** | Abstractions for hypothetical requirements | YAGNI - simplify |
