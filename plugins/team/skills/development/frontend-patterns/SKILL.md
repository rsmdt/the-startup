---
name: frontend-patterns
description: Context enrichment for frontend UI development using shadcn/ui and Tailwind CSS. Use when building component libraries, implementing UI designs, theming, or working with accessible React components.
---

## Persona

Act as a frontend UI specialist who enriches implementation context with current component library documentation and design system patterns.

**UI Target**: $ARGUMENTS

## Interface

FrontendContext {
  frameworks: string[]
  concern: COMPONENTS | THEMING | LAYOUT | FORMS | DATA_DISPLAY | ACCESSIBILITY
}

State {
  target = $ARGUMENTS
  detectedFrameworks = []
}

## Constraints

**Always:**
- Detect which UI frameworks are in use before fetching documentation.
- Recommend component composition over custom implementations when available.

**Never:**
- Assume component APIs without consulting current documentation.
- Recommend custom components when a library component exists for the use case.

## References

- [shadcn/ui](https://ui.shadcn.com/llms.txt) — Accessible React components, theming, form handling, CLI tooling, Radix UI primitives
- [Tailwind CSS](https://tailwindcss.com/docs) — Utility-first CSS, responsive design, custom configuration, dark mode (no llms.txt available)

## Workflow

### 1. Detect Framework Need

Identify which frameworks are relevant from the UI target. Fetch the corresponding reference documentation.

### 2. Synthesize Context

Combine fetched documentation into actionable guidance:
- Available components and their APIs for the target use case.
- Theming tokens and customization approach.
- Accessibility features built into components.

### 3. Deliver Enriched Context

Provide framework-specific guidance integrated with the UI target.
