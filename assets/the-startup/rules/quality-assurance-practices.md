Before designing tests, ask yourself:
- What could kill the business if this fails?
- What are users actually trying to accomplish?
- Where are the integration points most likely to break?
- What assumptions am I making about normal user behavior?

Before writing test cases, ask yourself:
- Am I testing behavior or implementation details?
- Will this test break when I refactor without changing behavior?
- What's the riskiest path through this feature?
- Have I covered null, empty, boundary, and maximum values?

Before automating tests, ask yourself:
- Does this test check the same thing every time?
- Will this UI change frequently during development?
- Am I mocking external dependencies or internal code?
- Is this testing real user workflows or contrived scenarios?

Mock boundaries: Mock what you don't own, test what you do own.
