# Reading Strategies by Document Type

Detailed strategies for extracting actionable information from each documentation type.

---

## README Files

READMEs are entry points. Extract these elements in order:

1. **Project Purpose**: First paragraph usually states what the project does
2. **Quick Start**: Look for "Getting Started", "Installation", or "Usage" sections
3. **Prerequisites**: Dependencies, environment requirements, version constraints
4. **Architecture Hints**: Links to other docs, directory structure descriptions
5. **Maintenance Status**: Last updated date, badges, contribution activity

**Reading Pattern**:
```
1. Scan headings to build mental map (30 seconds)
2. Read purpose/description section fully
3. Locate quick start commands - test if they work
4. Note any "gotchas" or "known issues" sections
5. Identify links to deeper documentation
```

**Red Flags**:
- No update in 12+ months on active project
- Quick start commands that fail
- References to deprecated dependencies
- Missing license or security sections

## API Documentation

Extract information in this priority:

1. **Authentication**: How to authenticate (API keys, OAuth, tokens)
2. **Base URL / Endpoints**: Entry points and environment variations
3. **Request Format**: Headers, body structure, content types
4. **Response Format**: Success/error shapes, status codes
5. **Rate Limits**: Throttling, quotas, retry policies
6. **Versioning**: How versions are specified, deprecation timeline

**Reading Pattern**:
```
1. Find authentication section first - nothing works without it
2. Locate a simple endpoint (health check, list operation)
3. Trace a complete request/response cycle
4. Note pagination patterns for list endpoints
5. Identify error response structure
6. Check for SDK/client library availability
```

**Cross-Reference Checks**:
- Compare documented endpoints against actual network calls
- Verify response schemas match real responses
- Test documented error codes actually occur

## Technical Specifications

Specifications define expected behavior. Extract:

1. **Requirements List**: Numbered requirements, acceptance criteria
2. **Constraints**: Technical limitations, compatibility requirements
3. **Data Models**: Entity definitions, relationships, constraints
4. **Interfaces**: API contracts, message formats, protocols
5. **Non-Functional Requirements**: Performance, security, scalability targets

**Reading Pattern**:
```
1. Identify document type (PRD, SDD, RFC, ADR)
2. Locate requirements or acceptance criteria section
3. Extract testable assertions (MUST, SHALL, SHOULD language)
4. Map requirements to implementation locations
5. Note any open questions or TBD items
```

**Verification Approach**:
- Create checklist from requirements
- Mark each as: Implemented / Partial / Missing / Contradicted
- Document gaps for follow-up

## Configuration Files

Configuration files control runtime behavior. Approach by file type:

### Package Manifests (package.json, Cargo.toml, pyproject.toml)
```
1. Project metadata: name, version, description
2. Entry points: main, bin, exports
3. Dependencies: runtime vs dev, version constraints
4. Scripts/commands: available automation
5. Engine requirements: Node version, Python version
```

### Environment Configuration (.env, config.yaml, settings.json)
```
1. Required variables (those without defaults)
2. Environment-specific overrides
3. Secret references (never actual values)
4. Feature flags and toggles
5. Service URLs and connection strings
```

### Build/Deploy Configuration (Dockerfile, CI configs, terraform)
```
1. Base images or providers
2. Build stages and dependencies
3. Environment variable injection points
4. Secret management approach
5. Output artifacts and destinations
```

### General Reading Pattern
```
1. Identify configuration format and schema (if available)
2. List all configurable options
3. Determine which have defaults vs require values
4. Trace where configuration values are consumed in code
5. Note any environment-specific overrides
```

## Architecture Decision Records (ADRs)

ADRs capture why decisions were made. Extract:

1. **Context**: What problem prompted the decision
2. **Decision**: What was chosen
3. **Consequences**: Trade-offs accepted
4. **Status**: Accepted, Deprecated, Superseded
5. **Related Decisions**: Links to related ADRs

**Reading Pattern**:
```
1. Read context to understand the problem space
2. Note alternatives that were considered
3. Understand why current approach was chosen
4. Check if decision is still active or superseded
5. Consider if context has changed since decision
```
