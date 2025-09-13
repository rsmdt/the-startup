# Agent Instruction Passing Reference Guide

## Overview

This document provides the EXACT patterns and templates for passing complete instructions to sub-agents, ensuring they inherit all necessary context, constraints, and discovery requirements.

## The Complete Instruction Template

### Full Pattern with All Components

```python
Task(subagent_type="[agent-type]", prompt="""
# SECTION 1: MANDATORY DISCOVERY (Always First)
DISCOVERY_FIRST: Before taking ANY action, you MUST:
[List specific discovery commands]
- find . -name "pattern" -type f | head -20
- grep -r "keyword" --include="*.ext"
- ls -la relevant/directories/

PATTERN_ANALYSIS: From discovery results, identify:
- [What patterns to look for]
- [What conventions to identify]
- [What structures to understand]

# SECTION 2: FOCUS/EXCLUDE BOUNDARIES
FOCUS: [Complete task description with all details]
- [Specific action 1]
- [Specific action 2]
- [Specific deliverable]

EXCLUDE: [All boundaries and restrictions]
- [What NOT to do 1]
- [What NOT to do 2]
- [Out of scope items]

# SECTION 3: INHERITED CONTEXT
CONTEXT: [Full background including prior work, current state, constraints]
        [Include relevant CLAUDE.md rules that apply to this task]
        
INHERITED_CONSTRAINTS:
- Parent command allows: [specific permissions]
- Parent command forbids: [specific restrictions]
- File creation rules: [where files can be created]
- Naming conventions: [established patterns to follow]

# SECTION 4: OUTPUT SPECIFICATIONS
OUTPUT: [EXACT format/structure expected with examples]
- File locations: [exact paths like docs/patterns/name.md]
- Format: [markdown/code/json/etc]
- Structure: [specific sections required]

# SECTION 5: SUCCESS CRITERIA
SUCCESS: [All completion criteria that must be met]
- [Measurable outcome 1]
- [Measurable outcome 2]
- [Quality standard to meet]

# SECTION 6: TERMINATION CONDITIONS
TERMINATION: [Explicit conditions for stopping]
- [When to stop condition 1]
- [Maximum iterations/attempts]
- [Error conditions that halt execution]
""")
```

## Component-by-Component Reference

### 1. Discovery Instructions (MANDATORY)

```python
# Pattern for Test Discovery
DISCOVERY_FIRST: Before writing ANY tests, you MUST:
- find . -name "*test*" -type f | head -20
- find . -name "*.spec.*" -o -name "*.test.*" | head -10
- Identify test framework: look for jest.config, mocha.opts, pytest.ini
- Locate test directories: tests/, __tests__, spec/, or co-located

# Pattern for Code Discovery
DISCOVERY_FIRST: Before writing ANY code, you MUST:
- find src/ -type f -name "*.ts" -o -name "*.js" | head -20
- analyze import statements in 3 files: grep "^import" [files]
- identify naming patterns: ls -la src/components/ src/services/
- check for configuration: cat tsconfig.json package.json

# Pattern for Documentation Discovery
DISCOVERY_FIRST: Before creating ANY documentation, you MUST:
- find docs/ -name "*.md" | grep -i "[topic-keywords]"
- check for existing related docs: grep -r "[topic]" docs/
- identify structure: head -20 docs/patterns/*.md
- verify naming convention: ls -la docs/patterns/ docs/domain/
```

### 2. FOCUS/EXCLUDE Pattern

```python
# Basic Pattern (Required)
FOCUS: [What to do - be comprehensive and specific]
EXCLUDE: [What NOT to do - be equally comprehensive]
CONTEXT: [All relevant background, constraints, and dependencies]
SUCCESS: [Measurable completion criteria]

# Enhanced Pattern (Recommended)
FOCUS: [Complete task description]
       - Primary objective: [main goal]
       - Specific actions: [list of actions]
       - Deliverables: [expected outputs]
       
EXCLUDE: [All boundaries and restrictions]
        - Do not modify: [protected areas]
        - Ignore: [out of scope items]
        - Skip: [non-relevant aspects]
        
CONTEXT: [Full background]
        Project type: [web app/API/library]
        Tech stack: [languages/frameworks]
        Current state: [what exists]
        Constraints: [limitations]
        Dependencies: [what this connects to]
        
OUTPUT: [Exact expectations]
       Location: [exact file paths]
       Format: [file type and structure]
       Naming: [exact naming pattern]
       
SUCCESS: [Measurable criteria]
        - Criteria 1: [specific metric]
        - Criteria 2: [quality standard]
        - Validation: [how to verify]
        
TERMINATION: [When to stop]
            - On success: [completion indicator]
            - On failure: [max attempts]
            - On timeout: [time limit]
```

### 3. Context Inheritance Pattern

```python
# JSON-Structured Context Inheritance
CONTEXT_INHERITANCE:
{json.dumps({
    "instructions": {
        "parent_constraints": "What the parent command requires/forbids",
        "file_creation_rules": "docs/** only, no root files",
        "naming_conventions": "kebab-case for files, PascalCase for components",
        "quality_requirements": "80% test coverage, no console.logs"
    },
    "memory": {
        "similar_tasks": "Previous authentication implementation",
        "recent_patterns": "Repository pattern used in services",
        "failed_attempts": "Avoided direct database access from controllers",
        "success_patterns": "Dependency injection worked well"
    },
    "state": {
        "current_progress": "Database schema complete, API design approved",
        "discovered_structure": "src/features/[feature]/ organization",
        "existing_files": ["src/features/auth/types.ts", "src/features/auth/service.ts"],
        "integration_points": ["POST /api/auth/login", "GET /api/auth/profile"]
    }
}, indent=2)}

# YAML-Structured Context Assembly
CONTEXT_ASSEMBLY_TEMPLATE:
  instructions: "{{inherited_instructions}}"
  memory: "{{compressed_session_state}}"
  world_state: "{{current_environment_context}}"
  delegation_rules: "{{parent_agent_constraints}}"
  discovery_results: "{{existing_patterns_found}}"
```

### 4. CLAUDE.md Context Passing

```python
# Extract and Include Relevant CLAUDE.md Rules
CONTEXT: [Task background] + From CLAUDE.md:
        
        # For Testing Tasks
        - TDD Required: Red-Green-Refactor cycle
        - One behavior per test, test through public interfaces
        - Mock external dependencies only, never internal code
        - Test edge cases: null, empty, boundaries
        - Co-locate tests with source files
        
        # For Code Writing Tasks
        - Functions under 20 lines, single responsibility
        - Intention-revealing names, no abbreviations
        - Validate inputs at boundaries
        - Return early for edge cases
        - Handle errors explicitly
        
        # For Architecture Tasks
        - SOLID principles apply
        - Prefer composition over inheritance
        - Program to interfaces, not implementations
        - Dependency injection for testability
        - Domain-driven design principles
        
        # For Security Tasks
        - Never hardcode secrets
        - Validate all inputs
        - Use parameterized queries
        - Principle of least privilege
        - Sanitize outputs for context
```

### 5. File Creation Specification

```python
# Explicit File Path Instructions
OUTPUT: Create the following files:
       1. Component: src/features/auth/components/LoginForm.tsx
       2. Service: src/features/auth/services/AuthService.ts
       3. Tests: src/features/auth/__tests__/LoginForm.test.tsx
       4. Types: src/features/auth/types/auth.types.ts
       
       NEVER create files in:
       - Project root
       - Outside of src/ directory
       - In wrong feature folder

# Documentation File Paths
OUTPUT: Documentation should be created/updated:
       - Pattern documentation: docs/patterns/[pattern-name].md
       - Business rules: docs/domain/[domain-concept].md
       - API contracts: docs/interfaces/[service-name].md
       
       If similar docs exist, UPDATE them instead of creating new
```

## Complete Real-World Examples

### Example 1: Test Writing with Full Context

```python
Task(subagent_type="test-writer", prompt="""
DISCOVERY_FIRST: Before writing ANY tests, you MUST execute:
- find . -name "*test*" -o -name "*spec*" -type f | head -20
- grep -l "describe\\|test\\|it(" $(find . -name "*.js" -o -name "*.ts") | head -10
- cat package.json | grep -A5 -B5 "test"
- ls -la tests/ test/ __tests__/ src/**/__tests__/ 2>/dev/null

PATTERN_ANALYSIS: From discovery, document:
- Test file locations found
- Test framework identified (Jest/Mocha/Vitest/etc)
- Naming convention (.test.ts vs .spec.ts vs _test.py)
- Test structure patterns (describe/it vs test suites)

FOCUS: Write comprehensive unit tests for the AuthenticationService class
       - Test all public methods
       - Cover success and failure scenarios
       - Include edge cases and error conditions
       - Ensure 90% code coverage

EXCLUDE: Do not write integration tests
        Do not test private methods directly
        Do not modify the implementation code
        Do not create new test infrastructure

CONTEXT: Testing a critical authentication service that handles user login,
        token validation, and session management. The service integrates
        with a PostgreSQL database and Redis cache.
        
        From CLAUDE.md:
        - TDD Required: Write failing tests first
        - One behavior per test case
        - Mock external dependencies (database, cache)
        - Test edge cases: null, empty, invalid inputs
        - Use descriptive test names that explain the scenario

CONTEXT_INHERITANCE:
{
  "instructions": {
    "parent_constraints": "Must use existing test framework only",
    "file_creation_rules": "Tests must be co-located with source files",
    "naming_conventions": "Use *.test.ts pattern found in discovery",
    "quality_requirements": "90% coverage, all tests must pass"
  },
  "memory": {
    "similar_tests": "See UserService.test.ts for pattern examples",
    "test_utilities": "Use test/helpers/mockDatabase.ts for DB mocking"
  },
  "state": {
    "service_location": "src/services/AuthenticationService.ts",
    "current_coverage": "45% - needs improvement",
    "discovered_test_location": "[WILL BE FILLED FROM DISCOVERY]"
  }
}

OUTPUT: Test file at: [DISCOVERED_LOCATION]/AuthenticationService.test.ts
       Using framework: [DISCOVERED_FRAMEWORK]
       Following pattern: [DISCOVERED_PATTERN]

SUCCESS: - All public methods have test coverage
        - Edge cases are thoroughly tested
        - Tests are isolated and fast
        - 90% code coverage achieved
        - All tests pass on first run

TERMINATION: Stop when all public methods are tested
            OR if discovery shows no test framework exists
            OR after 3 failed attempts to make tests pass
""")
```

### Example 2: Code Implementation with Full Context

```python
Task(subagent_type="implementation-specialist", prompt="""
DISCOVERY_FIRST: Before implementing, you MUST:
- find src/ -type f -name "*.ts" -o -name "*.tsx" | grep -E "(component|service|hook)" | head -20
- cat tsconfig.json | grep -E "(paths|baseUrl)"
- ls -la src/features/ src/components/ src/services/ 2>/dev/null
- grep "^import" $(find src/ -name "*.ts" | head -5) | head -20

PATTERN_ANALYSIS: Document the discovered:
- Directory structure (feature-based vs layer-based)
- Import patterns (relative vs absolute)
- File naming conventions
- Component patterns (functional vs class)

FOCUS: Implement a complete user profile management feature including:
       - ProfileView component for displaying user data
       - ProfileEdit component for updating information
       - ProfileService for API interactions
       - Custom useProfile hook for state management

EXCLUDE: Do not implement authentication (already exists)
        Do not create database migrations
        Do not modify existing user model
        Do not add new dependencies to package.json

CONTEXT: Building profile management for an existing React/TypeScript application.
        The app uses Redux for global state, React Query for server state,
        and Material-UI for components. Authentication is already implemented.
        
        From CLAUDE.md:
        - Single responsibility per component
        - Functions under 20 lines
        - TypeScript strict mode is enabled
        - Use existing design system components
        - Handle errors explicitly with user feedback

INHERITED_CONSTRAINTS:
- Architecture: Feature-folder organization required
- API: Must use existing ApiClient service
- State: Integrate with existing Redux store
- Styling: Use only Material-UI components
- Testing: Each component needs a test file

DISCOVERED_STRUCTURE: [TO BE FILLED FROM DISCOVERY]
{
  "component_location": "",
  "service_location": "",
  "hook_location": "",
  "import_pattern": "",
  "naming_convention": ""
}

OUTPUT: Create files in discovered locations:
       1. Components: [DISCOVERED]/ProfileView.tsx, ProfileEdit.tsx
       2. Service: [DISCOVERED]/ProfileService.ts
       3. Hook: [DISCOVERED]/useProfile.ts
       4. Tests: [DISCOVERED]/__tests__/*.test.tsx
       5. Types: [DISCOVERED]/types/profile.types.ts

SUCCESS: - Components render without errors
        - Service methods handle all CRUD operations
        - Hook provides clean state interface
        - TypeScript compilation passes
        - Follows all discovered patterns

TERMINATION: Stop when all components are implemented
            OR if discovery reveals incompatible architecture
            OR after addressing 3 rounds of type errors
""")
```

### Example 3: Documentation with Full Context

```python
Task(subagent_type="documentation-specialist", prompt="""
DISCOVERY_FIRST: Before documenting, you MUST:
- find docs/ -name "*.md" | xargs grep -l "cache\\|caching\\|performance"
- ls -la docs/patterns/ docs/domain/ docs/interfaces/
- head -50 docs/patterns/*.md | grep "^#"
- find . -name "README.md" -o -name "CONTRIBUTING.md"

DEDUPLICATION_CHECK: Search for existing related documentation:
- grep -r "caching strategy" docs/
- grep -r "performance optimization" docs/
- find docs/ -name "*cache*.md" -o -name "*performance*.md"

FOCUS: Document the caching strategy and patterns used in the application
       - Explain the multi-layer caching architecture
       - Document cache invalidation strategies
       - Provide implementation examples
       - Include performance metrics and trade-offs

EXCLUDE: Do not document database-specific caching
        Do not include infrastructure details
        Do not create API reference documentation
        Do not duplicate existing Redis documentation

CONTEXT: The application uses a sophisticated multi-layer caching strategy
        with Redis for distributed cache, in-memory cache for hot data,
        and CDN caching for static assets. Recent performance issues
        require clear documentation of when and how to use each layer.
        
        From project standards:
        - Use docs/patterns/ for reusable technical patterns
        - Use docs/domain/ for business logic documentation
        - Include code examples in markdown
        - Add diagrams where helpful

EXISTING_DOCUMENTATION: [TO BE FILLED FROM DISCOVERY]
{
  "related_files": [],
  "overlapping_content": [],
  "naming_pattern": "",
  "document_structure": ""
}

DEDUPLICATION_RULES:
- If caching patterns exist: UPDATE with new information
- If performance docs exist: ADD caching section
- If neither exists: CREATE new pattern document
- Always: Cross-reference related documentation

OUTPUT: Based on discovery:
       - If updating: Enhance [DISCOVERED_FILE]
       - If creating: docs/patterns/caching-strategy.md
       
       Include sections:
       1. Overview
       2. Architecture Diagram
       3. Implementation Patterns
       4. Code Examples
       5. Performance Metrics
       6. Trade-offs and Considerations
       7. References

SUCCESS: - No duplicate content created
        - Clear, actionable documentation
        - Code examples compile and work
        - Cross-references are accurate
        - Follows discovered document structure

TERMINATION: Stop when documentation is complete
            OR if discovery shows comprehensive docs exist
            OR after 2 rounds of revision
""")
```

## Critical Rules for Instruction Passing

### 1. Order Matters
```
1. DISCOVERY_FIRST (always first - sets the context)
2. FOCUS/EXCLUDE (defines boundaries)
3. CONTEXT (provides background)
4. OUTPUT (specifies deliverables)
5. SUCCESS (defines completion)
6. TERMINATION (sets limits)
```

### 2. Discovery is Mandatory
```python
# NEVER skip discovery
# BAD
prompt = "Write tests for the auth service"

# GOOD
prompt = """
DISCOVERY_FIRST: Before writing ANY tests, you MUST:
[specific discovery commands]
...
Write tests for the auth service
"""
```

### 3. Be Explicit About File Paths
```python
# BAD
OUTPUT: Create appropriate test files

# GOOD
OUTPUT: Create test file at: src/services/__tests__/AuthService.test.ts
       Following pattern: *.test.ts (discovered)
       In directory: src/services/__tests__/ (discovered)
```

### 4. Include Constraints from Parent
```python
# Always inherit parent constraints
INHERITED_CONSTRAINTS:
- From parent command: {parent.allowed_tools}
- File creation limited to: {parent.allowed_paths}
- Must follow conventions: {parent.naming_rules}
- Quality standards: {parent.quality_requirements}
```

### 5. Make Success Measurable
```python
# BAD
SUCCESS: Tests should be good

# GOOD
SUCCESS: - All public methods have tests (measurable)
        - 90% code coverage achieved (measurable)
        - All tests pass (measurable)
        - No console warnings (measurable)
```

## Validation Checklist

Before sending instructions to an agent, verify:

- [ ] Discovery commands are specific and will return useful results
- [ ] FOCUS clearly states what to do
- [ ] EXCLUDE clearly states what NOT to do
- [ ] CONTEXT includes all necessary background
- [ ] CLAUDE.md rules relevant to the task are included
- [ ] OUTPUT specifies exact file paths and formats
- [ ] SUCCESS criteria are measurable
- [ ] TERMINATION conditions are clear
- [ ] Inherited constraints from parent are passed through
- [ ] Deduplication rules are included (for documentation tasks)

## Common Mistakes to Avoid

### 1. Vague Discovery
```python
# BAD
DISCOVERY_FIRST: Look around the codebase

# GOOD
DISCOVERY_FIRST: Before writing tests, you MUST:
- find . -name "*.test.js" -o -name "*.spec.js" | head -20
- cat package.json | grep -A5 "jest"
```

### 2. Missing Context Inheritance
```python
# BAD
Task(subagent_type="test-writer", prompt="Write tests")

# GOOD
Task(subagent_type="test-writer", prompt=f"""
CONTEXT_INHERITANCE:
{json.dumps(parent_context)}
Write tests
""")
```

### 3. Unclear File Locations
```python
# BAD
OUTPUT: Put the tests in the right place

# GOOD
OUTPUT: Create tests at: {discovered_test_location}/AuthService.test.ts
```

## Conclusion

This reference guide provides the EXACT patterns for passing complete instructions to agents. The key is combining:
1. Mandatory discovery protocols
2. Clear FOCUS/EXCLUDE boundaries
3. Inherited context from parent
4. Explicit output specifications
5. Measurable success criteria

Following these patterns ensures sub-agents have all the context they need to succeed while respecting existing patterns and constraints.