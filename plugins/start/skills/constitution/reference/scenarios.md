# Constitution Examples

Example scenarios showing the constitution skill in action.

---

## Create New Constitution

```
User: /start:constitution

Claude: 📜 Constitution Setup

No CONSTITUTION.md found at project root.

I'll analyze your codebase to discover patterns and generate appropriate rules.

[Discovery process...]

📜 Proposed Constitution

Based on codebase analysis:
- Project Type: Next.js with TypeScript
- Framework: React 18
- Testing: Vitest + React Testing Library
- Data: Prisma ORM

[Proposed rules by category...]

Would you like to:
1. Approve these rules (recommended)
2. Modify before saving
3. Cancel
```

## Create with Focus Areas

```
User: /start:constitution "Focus on security and API patterns"

Claude: 📜 Constitution Setup (Focused)

Focus areas: Security, API patterns

[Targeted discovery...]

📜 Proposed Constitution

Security (5 rules):
- L1: No hardcoded secrets
- L1: No eval/exec usage
- L1: Parameterized SQL queries
- L2: Input validation required
- L2: CORS configuration required

API Patterns (3 rules):
- L1: Error responses use standard format
- L2: Rate limiting on public endpoints
- L3: OpenAPI documentation

[Approval prompt...]
```

## Update Existing Constitution

```
User: /start:constitution "Add testing rules"

Claude: 📜 Constitution Update

Found existing CONSTITUTION.md with 8 rules.

Current categories:
- Security (3 rules)
- Architecture (2 rules)
- Code Quality (3 rules)

Focus: Adding testing rules

[Discovery of test patterns...]

Proposed additions to Testing category:
- L1: No .only in committed tests
- L2: Test descriptions must be meaningful
- L3: Integration tests for API endpoints

Would you like to:
1. Add these rules (recommended)
2. Review and modify
3. Cancel
```
