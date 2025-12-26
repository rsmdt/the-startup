---
name: data-modeling
description: Schema design, entity relationships, normalization, and database patterns. Use when designing database schemas, modeling domain entities, deciding between normalized and denormalized structures, choosing between relational and NoSQL approaches, or planning schema migrations. Covers ER modeling, normal forms, and data evolution strategies.
---

# Data Modeling

## When to Use

- Designing new database schemas from domain requirements
- Analyzing existing schemas for optimization opportunities
- Deciding between normalized and denormalized structures
- Choosing appropriate data stores (relational vs NoSQL)
- Planning schema evolution and migration strategies
- Modeling complex entity relationships

## Philosophy

Data models outlive applications. A well-designed schema encodes business rules, enforces integrity, and enables performance optimization. The goal is to create models that are correct first, then optimize for access patterns while maintaining data integrity.

## Entity-Relationship Modeling

### Identifying Entities

Entities represent distinct business concepts that have identity and lifecycle.

**Entity Identification Checklist:**
- Has unique identity across the system
- Has attributes that describe it
- Participates in relationships with other entities
- Has a meaningful lifecycle (created, modified, archived)
- Would be stored and retrieved independently

**Common Entity Patterns:**
- Core domain objects (User, Product, Order)
- Reference/lookup data (Country, Status, Category)
- Transactional records (Payment, LogEntry, Event)
- Associative entities (OrderItem, Enrollment, Permission)

### Relationship Types

| Type | Notation | Example | Implementation |
|------|----------|---------|----------------|
| One-to-One | 1:1 | User - Profile | FK with unique constraint or same table |
| One-to-Many | 1:N | Customer - Orders | FK on the "many" side |
| Many-to-Many | M:N | Students - Courses | Junction/bridge table |

**Relationship Considerations:**
- Cardinality: minimum and maximum on each side
- Optionality: required vs optional participation
- Direction: unidirectional vs bidirectional navigation
- Cascade behavior: what happens on delete/update

### Attribute Analysis

**Attribute Types:**
- Simple: single atomic value (name, price)
- Composite: structured value (address = street + city + postal)
- Derived: calculated from other attributes (age from birthdate)
- Multi-valued: repeating values (phone numbers, tags)

**Key Types:**
- Natural key: business-meaningful identifier (SSN, ISBN)
- Surrogate key: system-generated identifier (UUID, auto-increment)
- Composite key: multiple columns forming identity
- Candidate key: any attribute(s) that could serve as primary key

**Best Practice:** Prefer surrogate keys for primary keys; use natural keys as unique constraints.

## Normalization

### Normal Forms Progression

Each normal form builds on the previous. Normalize until requirements dictate otherwise.

#### First Normal Form (1NF)

**Rule:** Eliminate repeating groups; each cell contains atomic values.

**Violation Example:**
```
Order(id, customer, items: "widget,gadget,thing")
```

**Resolution:**
```
Order(id, customer)
OrderItem(order_id, item_name)
```

#### Second Normal Form (2NF)

**Rule:** Remove partial dependencies on composite keys.

**Violation Example:**
```
OrderItem(order_id, product_id, product_name, quantity)
                                 ^-- depends only on product_id
```

**Resolution:**
```
OrderItem(order_id, product_id, quantity)
Product(product_id, product_name)
```

#### Third Normal Form (3NF)

**Rule:** Remove transitive dependencies; non-key columns depend only on the key.

**Violation Example:**
```
Employee(id, department_id, department_name)
                            ^-- depends on department_id, not employee id
```

**Resolution:**
```
Employee(id, department_id)
Department(id, name)
```

#### Boyce-Codd Normal Form (BCNF)

**Rule:** Every determinant is a candidate key.

**Violation Example:**
```
CourseOffering(student, course, instructor)
-- Constraint: each instructor teaches only one course
-- instructor -> course (but instructor is not a candidate key)
```

**Resolution:**
```
InstructorCourse(instructor, course) -- instructor is key
Enrollment(student, instructor) -- references instructor
```

### When to Stop Normalizing

Stop at 3NF for most OLTP systems. Consider BCNF when:
- Update anomalies cause data corruption
- Data integrity is paramount
- Write frequency is high

## Denormalization Strategies

Denormalize intentionally for read performance, not out of laziness.

### Calculated Columns

Store derived values to avoid repeated computation.

```
Order
  - subtotal (calculated once on item changes)
  - tax_amount (calculated once)
  - total (calculated once)
```

**Trade-off:** Faster reads, more complex writes, potential consistency issues.

### Materialized Relationships

Embed frequently-accessed related data.

```
Post
  - author_id
  - author_name (copied from User.name)
  - author_avatar_url (copied from User.avatar_url)
```

**Trade-off:** Eliminates joins, requires synchronization on source changes.

### Aggregation Tables

Pre-compute summaries for reporting.

```
DailySales
  - date
  - product_id
  - units_sold (sum)
  - revenue (sum)
```

**Trade-off:** Fast analytics, storage overhead, stale until refreshed.

### Denormalization Decision Matrix

| Factor | Normalize | Denormalize |
|--------|-----------|-------------|
| Write frequency | High | Low |
| Read frequency | Low | High |
| Data consistency | Critical | Eventual OK |
| Query complexity | Simple | Complex joins |
| Data size | Small | Large |

## NoSQL Data Modeling Patterns

### Document Stores (MongoDB, DynamoDB)

**Embedding Pattern:**
Embed related data that is read together and has 1:few relationship.

```json
{
  "order_id": "123",
  "customer": {
    "id": "456",
    "name": "Jane Doe",
    "email": "jane@example.com"
  },
  "items": [
    {"product_id": "A1", "name": "Widget", "quantity": 2}
  ]
}
```

**Referencing Pattern:**
Reference related data when it changes independently or is shared.

```json
{
  "order_id": "123",
  "customer_id": "456",
  "item_ids": ["A1", "B2"]
}
```

**Hybrid Pattern:**
Embed summary data, reference for full details.

```json
{
  "order_id": "123",
  "customer_summary": {
    "id": "456",
    "name": "Jane Doe"
  },
  "items": [
    {"product_id": "A1", "name": "Widget", "quantity": 2}
  ]
}
```

### Key-Value Stores

**Access Pattern Design:**
Design keys around query patterns.

```
USER:{user_id} -> user data
USER:{user_id}:ORDERS -> list of order ids
ORDER:{order_id} -> order data
```

**Composite Keys:**
Combine entity type with identifiers for namespacing.

### Wide-Column Stores (Cassandra, HBase)

**Partition Key Design:**
Choose partition keys for even distribution and access locality.

```
Primary Key: (user_id, order_date)
             ^-- partition key (distribution)
                       ^-- clustering column (ordering)
```

**Avoid:**
- High-cardinality partition keys causing hot spots
- Large partitions exceeding recommended sizes
- Scatter-gather queries across partitions

### Graph Databases

**Node and Relationship Design:**
- Nodes: entities with properties
- Relationships: named, directed, with properties
- Labels: categorize nodes for efficient traversal

```
(User)-[:PURCHASED {date, amount}]->(Product)
(User)-[:FOLLOWS]->(User)
(Product)-[:BELONGS_TO]->(Category)
```

## Schema Evolution Strategies

### Additive Changes (Safe)

- Add new nullable columns
- Add new tables
- Add new indexes
- Add new optional fields (NoSQL)

### Breaking Changes (Require Migration)

- Remove columns/tables
- Rename columns/tables
- Change data types
- Add non-nullable columns without defaults

### Migration Patterns

**Expand-Contract Pattern:**
1. Add new column alongside old
2. Backfill new column from old
3. Update application to use new column
4. Remove old column

**Blue-Green Schema:**
1. Create new version of schema
2. Dual-write to both versions
3. Migrate reads to new version
4. Drop old version

**Versioned Documents (NoSQL):**
```json
{
  "_schema_version": 2,
  "name": "Jane",
  "email": "jane@example.com"
}
```

Handle multiple versions in application code during transition.

## Best Practices

- Model the domain first, then optimize for access patterns
- Use surrogate keys for primary keys; natural keys as unique constraints
- Normalize to 3NF for OLTP; denormalize deliberately for read-heavy loads
- Document all foreign key relationships and cascade behaviors
- Version control all schema changes as migration scripts
- Test migrations on production-like data volumes
- Consider query patterns when designing NoSQL schemas
- Plan for schema evolution from day one

## Anti-Patterns

- Designing schemas around UI forms instead of domain concepts
- Using generic columns (field1, field2, field3)
- Entity-Attribute-Value (EAV) for structured data
- Storing comma-separated values in single columns
- Circular foreign key dependencies
- Missing indexes on foreign key columns
- Hard-deleting data without soft-delete consideration
- Ignoring temporal aspects (effective dates, audit trails)

## References

- [templates/schema-design-template.md](templates/schema-design-template.md) - Structured template for schema documentation
