Before applying SOLID principles, ask yourself:
- Does this class have exactly one reason to change? (Single Responsibility)
- Can I extend behavior without modifying existing code? (Open/Closed)
- Could I substitute any implementation without breaking clients? (Liskov)
- Are clients forced to depend on methods they don't use? (Interface Segregation)
- Do high-level modules depend on low-level details? (Dependency Inversion)

When designing interfaces, ask yourself:
- What's the smallest interface that solves the current problem?
- Can I name this interface based on what it does, not how it works?
- Would a new implementation of this interface require zero changes to existing code?
- Does this interface hide implementation details completely?

Before choosing technology, ask yourself:
- Does the team already know this technology well?
- Will this technology still be supported in 3 years?
- Can we hire people who know this technology?
- Does this solve a real performance or capability gap?
- What's the simplest alternative that would work?

When documenting architectural decisions, ask yourself:
- What problem does this solve that existing patterns don't?
- What will happen if someone ignores this pattern?
- Can I provide a working example others can copy?
- Have I explained when NOT to use this pattern?
- Will someone unfamiliar with the codebase understand why this exists?

Before creating abstractions, ask yourself:
- Have I seen this exact problem at least 3 times?
- Does this abstraction make the simple case simpler?
- Can I remove this abstraction without major code changes?
- Does this abstraction leak implementation details?