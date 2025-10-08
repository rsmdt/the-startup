# Domain Documentation

This directory contains business rules, workflows, and domain patterns discovered through analysis.

## Purpose

Domain documentation captures the **business logic and rules** that drive the system. This includes:

- **Business Rules**: Validation rules, business constraints, and domain invariants
- **Workflows**: Business processes and operational procedures
- **Domain Patterns**: Recurring business patterns and domain-specific solutions
- **Terminology**: Ubiquitous language and domain vocabulary

## When to Create Domain Documentation

Create domain documentation when using `/s:analyze` to discover:
- Business validation rules in the codebase
- Domain-specific workflows and processes
- Business constraints and requirements
- Domain modeling patterns

## File Naming Convention

Use descriptive, searchable names that reflect the business concept:
- `user-registration-rules.md` - Business rules for user registration
- `order-fulfillment-workflow.md` - Order processing workflow
- `pricing-calculation-rules.md` - Pricing business logic
- `subscription-lifecycle.md` - Subscription state management

## Structure

Each domain document should include:
1. **Overview**: What business concept this addresses
2. **Business Rules**: Explicit rules and constraints
3. **Workflows**: Process flows and state transitions
4. **Examples**: Real-world scenarios and use cases
5. **Related Patterns**: Links to relevant technical patterns in `docs/patterns/`

## Related Directories

- `../patterns/` - Technical implementation patterns
- `../interfaces/` - External system integration contracts
