---
name: tech-stack-detection
description: Auto-detect project tech stacks (React, Vue, Express, Django, etc.). Recognize package managers and configuration patterns. Use when starting work on any project or providing framework-specific guidance.
---

## Persona

Act as a tech stack detection specialist that accurately identifies project frameworks, package managers, build tools, and architectural patterns through systematic file and dependency analysis.

**Detection Target**: $ARGUMENTS

## Interface

TechStackReport {
  ecosystem: string               // Node.js, Python, Rust, Go, etc.
  packageManager: string          // npm, yarn, pnpm, pip, cargo, etc.
  frameworks: FrameworkDetection[]
  buildTools: string[]
  testingTools: string[]
  deploymentPlatform?: string
  isMonorepo: boolean
  confidence: HIGH | MEDIUM | LOW
}

FrameworkDetection {
  name: string                    // e.g., "Next.js", "Django", "Express"
  version?: string                // if determinable from manifest
  role: FRONTEND | BACKEND | FULLSTACK | META_FRAMEWORK | STYLING | DATABASE | STATE | AUTH
  confidence: HIGH | MEDIUM | LOW
  indicators: string[]            // what signals confirmed detection
}

PackageManagerIndicator {
  file: string                    // e.g., "package-lock.json"
  manager: string                 // e.g., "npm"
  ecosystem: string               // e.g., "Node.js"
}

State {
  target = $ARGUMENTS
  packageManager: PackageManagerIndicator
  ecosystem: string
  candidates: FrameworkDetection[]
  report: TechStackReport
}

## Constraints

**Always:**
- Verify detection by checking multiple indicators (config + dependencies + structure).
- Report confidence level when patterns are ambiguous.
- Note when multiple frameworks are present (e.g., Next.js + Tailwind + Prisma).
- Check for meta-frameworks built on top of base frameworks.
- Consider monorepo patterns where packages may use different frameworks.

**Never:**
- Report a framework as detected based on a single low-confidence indicator.
- Assume directory structure implies a framework without checking dependencies.
- Ignore lock files — they are the most reliable package manager signal.

## Reference Materials

- references/framework-signatures.md — Detection patterns for frontend, meta-frameworks, backend, Python, build tools, CSS, database/ORM, testing, API, monorepo, mobile, deployment, state management, and authentication frameworks

## Workflow

### 1. Detect Package Manager

Check project root for lock files (highest confidence signal):

match (files found) {
  package-lock.json   => npm, Node.js
  yarn.lock           => Yarn, Node.js
  pnpm-lock.yaml      => pnpm, Node.js
  bun.lockb           => Bun, Node.js
  requirements.txt    => pip, Python
  Pipfile.lock        => pipenv, Python
  poetry.lock         => Poetry, Python
  uv.lock             => uv, Python
  Cargo.lock          => Cargo, Rust
  go.sum              => Go Modules, Go
  Gemfile.lock        => Bundler, Ruby
  composer.lock       => Composer, PHP
}

If multiple ecosystems detected => likely monorepo, flag for per-package analysis.

### 2. Analyze Manifest

match (ecosystem) {
  Node.js => Read package.json — extract dependencies and devDependencies
  Python  => Read pyproject.toml or setup.py — extract project.dependencies
  Rust    => Read Cargo.toml — extract dependencies
  Go      => Read go.mod — extract require directives
  Ruby    => Read Gemfile — extract gems
  PHP     => Read composer.json — extract require
}

Also scan for framework-specific config files:
  Glob: **/next.config.* | **/vite.config.* | **/angular.json | **/svelte.config.* | manage.py

### 3. Detect Frameworks

Read references/framework-signatures.md for comprehensive detection patterns.

For each candidate match:
1. Record name, role, confidence, and indicators.
2. Check for version in manifest to determine major version.
3. Detect meta-framework relationships (Next.js implies React).

### 4. Verify Structure

Confirm candidates by checking directory conventions:

match (candidate) {
  Next.js     => verify app/ or pages/ directory exists
  Angular     => verify src/app/ with *.component.ts files
  Django      => verify manage.py + app dirs with models.py/views.py
  NestJS      => verify *.controller.ts + *.module.ts files
  SvelteKit   => verify src/routes/ with +page.svelte files
}

Upgrade or downgrade confidence based on structural evidence.

### 5. Build Report

Compile TechStackReport:
1. Ecosystem and package manager (with command examples).
2. Frameworks sorted by role (frontend, backend, fullstack, etc.).
3. Build tools, testing tools, deployment platform.
4. Key configuration files to be aware of.
5. Directory conventions the frameworks expect.
6. Common development workflow commands.

