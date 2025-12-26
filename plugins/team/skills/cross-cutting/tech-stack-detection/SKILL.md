---
name: tech-stack-detection
description: Auto-detect project tech stacks (React, Vue, Express, Django, etc.). Recognize package managers and configuration patterns. Use when starting work on any project, analyzing dependencies, or providing framework-specific guidance.
---

# Framework Detection

## When to Use

- Starting work on an unfamiliar project
- Determining appropriate tooling and patterns for recommendations
- Providing framework-specific guidance and best practices
- Identifying package manager for dependency operations
- Understanding project architecture before making changes

## Detection Methodology

### Step 1: Package Manager Detection

Check for package manager indicators in the project root:

| File | Package Manager | Ecosystem |
|------|-----------------|-----------|
| `package-lock.json` | npm | Node.js |
| `yarn.lock` | Yarn | Node.js |
| `pnpm-lock.yaml` | pnpm | Node.js |
| `bun.lockb` | Bun | Node.js |
| `requirements.txt` | pip | Python |
| `Pipfile.lock` | pipenv | Python |
| `poetry.lock` | Poetry | Python |
| `uv.lock` | uv | Python |
| `Cargo.lock` | Cargo | Rust |
| `go.sum` | Go Modules | Go |
| `Gemfile.lock` | Bundler | Ruby |
| `composer.lock` | Composer | PHP |

### Step 2: Configuration File Analysis

Examine root-level configuration files for framework indicators:

1. **Read `package.json`** - Check `dependencies` and `devDependencies` for framework packages
2. **Read `pyproject.toml`** - Check `[project.dependencies]` or `[tool.poetry.dependencies]`
3. **Read framework-specific configs** - `next.config.js`, `vite.config.ts`, `angular.json`, etc.

### Step 3: Directory Structure Patterns

Identify framework conventions:

- `app/` or `src/app/` - Next.js App Router, Angular
- `pages/` - Next.js Pages Router, Nuxt.js
- `components/` - React/Vue component-based architecture
- `routes/` - Remix, SvelteKit
- `views/` - Django, Rails, Laravel
- `controllers/` - MVC frameworks (Rails, Laravel, NestJS)

### Step 4: Framework-Specific Patterns

Apply detection patterns from the framework signatures reference.

## Detection Workflow

```
START
  |
  v
[Check lock files] --> Identify package manager
  |
  v
[Read manifest] --> package.json / pyproject.toml / Cargo.toml
  |
  v
[Check dependencies] --> Match against known frameworks
  |
  v
[Check config files] --> Framework-specific configuration
  |
  v
[Verify directory structure] --> Confirm framework conventions
  |
  v
[Output] --> Framework, version, package manager, key patterns
```

## Output Format

When reporting detected framework, include:

1. **Framework name and version** (if determinable)
2. **Package manager** (with command examples)
3. **Key configuration files** to be aware of
4. **Directory conventions** the framework expects
5. **Common commands** for development workflow

## Best Practices

- Always verify detection by checking multiple indicators (config + dependencies + structure)
- Report confidence level when patterns are ambiguous
- Note when multiple frameworks are present (e.g., Next.js + Tailwind + Prisma)
- Check for meta-frameworks built on top of base frameworks
- Consider monorepo patterns where different packages may use different frameworks

## References

See `references/framework-signatures.md` for comprehensive detection patterns for all major frameworks.
