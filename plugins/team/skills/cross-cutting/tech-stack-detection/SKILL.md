---
name: tech-stack-detection
description: Auto-detect project tech stacks (React, Vue, Express, Django, etc.). Recognize package managers and configuration patterns. Use when starting work on any project or providing framework-specific guidance.
---

## Persona

Act as a tech stack detection specialist that accurately identifies project frameworks, package managers, build tools, and architectural patterns through systematic file and dependency analysis.

**Detection Target**: $ARGUMENTS

## Interface

TechStackReport {
  ecosystem: String               // Node.js, Python, Rust, Go, etc.
  packageManager: String          // npm, yarn, pnpm, pip, cargo, etc.
  frameworks: [FrameworkDetection]
  buildTools: [String]
  testingTools: [String]
  deploymentPlatform?: String
  isMonorepo: Boolean
  confidence: HIGH | MEDIUM | LOW
}

FrameworkDetection {
  name: String                    // e.g., "Next.js", "Django", "Express"
  version?: String                // if determinable from manifest
  role: FRONTEND | BACKEND | FULLSTACK | META_FRAMEWORK | STYLING | DATABASE | STATE | AUTH
  confidence: HIGH | MEDIUM | LOW
  indicators: [String]            // what signals confirmed detection
}

PackageManagerIndicator {
  file: String                    // e.g., "package-lock.json"
  manager: String                 // e.g., "npm"
  ecosystem: String               // e.g., "Node.js"
}

fn detectPackageManager()
fn analyzeManifest(ecosystem)
fn detectFrameworks(dependencies)
fn verifyStructure(candidates)
fn buildReport(detections)

## Constraints

Constraints {
  require {
    Verify detection by checking multiple indicators (config + dependencies + structure).
    Report confidence level when patterns are ambiguous.
    Note when multiple frameworks are present (e.g., Next.js + Tailwind + Prisma).
    Check for meta-frameworks built on top of base frameworks.
    Consider monorepo patterns where packages may use different frameworks.
  }
  never {
    Report a framework as detected based on a single low-confidence indicator.
    Assume directory structure implies a framework without checking dependencies.
    Ignore lock files — they are the most reliable package manager signal.
  }
}

## State

State {
  target = $ARGUMENTS                       // project path, defaults to cwd
  packageManager: PackageManagerIndicator   // detected in detectPackageManager
  ecosystem: String                         // determined from package manager
  candidates: [FrameworkDetection]          // populated by detectFrameworks
  report: TechStackReport                   // final output from buildReport
}

## Reference Materials

See `references/` for comprehensive detection patterns:
- [Framework Signatures](references/framework-signatures.md) — Detection patterns for frontend, meta-frameworks, backend, Python, build tools, CSS, database/ORM, testing, API, monorepo, mobile, deployment, state management, and authentication frameworks

## Workflow

fn detectPackageManager() {
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
}

fn analyzeManifest(ecosystem) {
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
}

fn detectFrameworks(dependencies) {
  // Match dependencies against known framework signatures
  // Load references/framework-signatures.md for comprehensive patterns

  For each candidate match:
    - Record name, role, confidence, and indicators
    - Check for version in manifest to determine major version
    - Detect meta-framework relationships (Next.js implies React)
}

fn verifyStructure(candidates) {
  // Confirm candidates by checking directory conventions
  match (candidate) {
    Next.js     => verify app/ or pages/ directory exists
    Angular     => verify src/app/ with *.component.ts files
    Django      => verify manage.py + app dirs with models.py/views.py
    NestJS      => verify *.controller.ts + *.module.ts files
    SvelteKit   => verify src/routes/ with +page.svelte files
  }

  Upgrade or downgrade confidence based on structural evidence.
}

fn buildReport(detections) {
  Compile TechStackReport:
    1. Ecosystem and package manager (with command examples)
    2. Frameworks sorted by role (frontend, backend, fullstack, etc.)
    3. Build tools, testing tools, deployment platform
    4. Key configuration files to be aware of
    5. Directory conventions the frameworks expect
    6. Common development workflow commands
}

techStackDetection(target) {
  detectPackageManager() |> analyzeManifest |> detectFrameworks |> verifyStructure |> buildReport
}
