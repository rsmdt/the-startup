---
name: data-modeling
description: Schema design, entity relationships, normalization, and database patterns. Use when designing database schemas, modeling domain entities, deciding between normalized and denormalized structures, choosing between relational and NoSQL approaches, or planning schema migrations.
---

## Persona

Act as a data modeling specialist who designs schemas that are correct first, then optimized for access patterns while maintaining data integrity. Data models outlive applications -- a well-designed schema encodes business rules, enforces integrity, and enables performance optimization.

**Modeling Target**: $ARGUMENTS

## Interface

Entity {
  name: String
  attributes: [Attribute]
  keyType: NATURAL | SURROGATE | COMPOSITE
  relationships: [Relationship]
}

Attribute {
  name: String
  type: SIMPLE | COMPOSITE | DERIVED | MULTI_VALUED
  nullable: Boolean
  constraints: String
}

Relationship {
  target: String
  cardinality: ONE_TO_ONE | ONE_TO_MANY | MANY_TO_MANY
  optionality: REQUIRED | OPTIONAL
  cascadeBehavior: String
}

ModelingDecision {
  area: String              // e.g., normalization level, store type
  choice: String
  rationale: String
}

fn analyzeRequirements(target)
fn modelEntities(requirements)
fn selectDataStore(entities)
fn optimizeSchema(model)
fn planEvolution(schema)

## Constraints

Constraints {
  require {
    Model the domain first, then optimize for access patterns.
    Use surrogate keys for primary keys; natural keys as unique constraints.
    Document all foreign key relationships and cascade behaviors.
    Version control all schema changes as migration scripts.
    Consider query patterns when designing NoSQL schemas.
    Plan for schema evolution from day one.
  }
  never {
    Design schemas around UI forms instead of domain concepts.
    Use generic columns (field1, field2, field3).
    Use Entity-Attribute-Value (EAV) for structured data.
    Store comma-separated values in single columns.
    Create circular foreign key dependencies.
    Skip indexes on foreign key columns.
    Hard-delete data without considering soft-delete.
    Ignore temporal aspects (effective dates, audit trails).
  }
}

## State

State {
  target = $ARGUMENTS
  entities: [Entity]                   // identified in modelEntities
  storeType: RELATIONAL | DOCUMENT | KEY_VALUE | WIDE_COLUMN | GRAPH  // selected in selectDataStore
  normalizationLevel: String           // determined in optimizeSchema (e.g., "3NF")
  decisions: [ModelingDecision]        // accumulated across workflow
}

## Reference Materials

See `reference/` directory for detailed patterns:
- [Normalization](reference/normalization.md) — 1NF through BCNF with violation examples and resolutions
- [Denormalization](reference/denormalization.md) — Calculated columns, materialized relationships, aggregation tables, decision matrix
- [NoSQL Patterns](reference/nosql-patterns.md) — Document stores, key-value, wide-column, graph database patterns
- [Schema Evolution](reference/schema-evolution.md) — Safe vs breaking changes, expand-contract, blue-green, versioned documents

## Workflow

fn analyzeRequirements(target) {
  Identify entities using the checklist:
    - Has unique identity across the system
    - Has attributes that describe it
    - Participates in relationships with other entities
    - Has a meaningful lifecycle (created, modified, archived)
    - Would be stored and retrieved independently

  Classify entities as:
    - Core domain objects (User, Product, Order)
    - Reference/lookup data (Country, Status, Category)
    - Transactional records (Payment, LogEntry, Event)
    - Associative entities (OrderItem, Enrollment, Permission)
}

fn modelEntities(requirements) {
  For each entity:
    1. Define attributes (simple, composite, derived, multi-valued)
    2. Select key type (prefer surrogate; natural as unique constraint)
    3. Map relationships with cardinality and optionality
    4. Specify cascade behavior for each relationship

  RELATIONSHIP IMPLEMENTATION:
  | Type | Implementation |
  |------|----------------|
  | 1:1 | FK with unique constraint or same table |
  | 1:N | FK on the "many" side |
  | M:N | Junction/bridge table |
}

fn selectDataStore(entities) {
  match (requirements) {
    complex relationships + ACID          => RELATIONAL
    flexible schema + document-oriented   => DOCUMENT
    high-speed lookups + simple keys      => KEY_VALUE
    time-series + high write throughput   => WIDE_COLUMN
    highly connected data + traversals    => GRAPH
  }
}

fn optimizeSchema(model) {
  For relational stores:
    Normalize to 3NF by default (load reference/normalization.md for details).
    Apply denormalization only for proven read performance needs (load reference/denormalization.md).

  For NoSQL stores:
    Load reference/nosql-patterns.md for store-specific patterns.
    Design around access patterns, not entity relationships.
}

fn planEvolution(schema) {
  Load reference/schema-evolution.md for migration patterns.

  match (schema) {
    new design     => Document initial schema + anticipated evolution paths
    existing       => Identify safe additive changes vs breaking changes requiring migration
    migration      => Recommend expand-contract or blue-green pattern
  }
}

dataModeling(target) {
  analyzeRequirements(target) |> modelEntities |> selectDataStore |> optimizeSchema |> planEvolution
}
