# Refactoring Perspectives

Analysis perspectives and focus area mapping for the refactor skill.

---

## Standard Perspectives

| Perspective | Intent | What to Analyze |
|-------------|--------|-----------------|
| 🔧 **Code Smells** | Find improvement opportunities | Long methods, duplication, complexity, deep nesting, magic numbers |
| 🔗 **Dependencies** | Map coupling issues | Circular dependencies, tight coupling, abstraction violations |
| 🧪 **Test Coverage** | Assess safety for refactoring | Existing tests, coverage gaps, test quality, missing assertions |
| 🏗️ **Patterns** | Identify applicable techniques | Design patterns, refactoring recipes, architectural improvements |
| ⚠️ **Risk** | Evaluate change impact | Blast radius, breaking changes, complexity, rollback difficulty |

## Simplification Perspectives

Use when $ARGUMENTS focuses on simplification (e.g., "simplify", "clean up", "reduce complexity"):

| Perspective | Intent | What to Find |
|-------------|--------|--------------|
| 🔧 **Complexity** | Reduce cognitive load | Long methods (>20 lines), deep nesting, complex conditionals, convoluted loops, tangled async/promise chains |
| 📝 **Clarity** | Make intent obvious | Unclear names, magic numbers, inconsistent patterns, overly defensive code, unnecessary ceremony, nested ternaries |
| 🏗️ **Structure** | Improve organization | Mixed concerns, tight coupling, bloated interfaces, god objects, too many parameters, hidden dependencies |
| 🧹 **Waste** | Eliminate what shouldn't exist | Duplication, dead code, unused abstractions, speculative generality, copy-paste patterns, unreachable paths |
