# Schema Design Template

## Purpose

Structured documentation for database schema designs. Use this template when designing new schemas, documenting existing ones, or planning migrations. Ensures consistent coverage of entities, relationships, constraints, and evolution strategy.

## Template

```markdown
# Schema: [Schema Name]

## Overview

**Domain:** [Business domain this schema serves]
**Data Store:** [PostgreSQL | MySQL | MongoDB | DynamoDB | etc.]
**Primary Access Patterns:** [List main query patterns]
**Consistency Requirements:** [Strong | Eventual | Mixed]

## Entities

### [Entity Name]

**Purpose:** [What this entity represents in the domain]

**Table/Collection:** `[table_name]`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK | Unique identifier |
| [column] | [type] | [constraints] | [description] |

**Indexes:**
- `idx_[name]` on ([columns]) - [purpose]

**Constraints:**
- [constraint description]

---

[Repeat for each entity]

## Relationships

### [Relationship Name]

**Type:** [1:1 | 1:N | M:N]
**From:** [Source Entity]
**To:** [Target Entity]
**Implementation:** [FK location or junction table]
**Cascade Behavior:** [ON DELETE/UPDATE actions]
**Business Rule:** [Why this relationship exists]

---

[Repeat for each relationship]

## Normalization Analysis

### Current Normal Form: [1NF | 2NF | 3NF | BCNF]

**Deliberate Denormalizations:**

| Location | Pattern | Justification |
|----------|---------|---------------|
| [table.column] | [calculated | embedded | aggregated] | [performance reason] |

## Access Patterns

### [Pattern Name]

**Query:** [Description of what is being retrieved]
**Frequency:** [Requests per second / day]
**Latency Requirement:** [Target response time]
**Supporting Indexes:** [Index names used]

Example Query:
```sql
[Sample SQL or query]
```

---

[Repeat for each major access pattern]

## Data Lifecycle

### Retention Policy

| Entity | Retention Period | Archive Strategy | Deletion Strategy |
|--------|-----------------|------------------|-------------------|
| [entity] | [period] | [approach] | [soft/hard delete] |

### Audit Requirements

- [What changes are tracked]
- [Where audit logs are stored]
- [Retention of audit data]

## Migration Plan

### Version: [version number]

**Changes from Previous:**
- [Change 1]
- [Change 2]

**Migration Steps:**
1. [Step with rollback procedure]
2. [Step with rollback procedure]

**Rollback Procedure:**
1. [Step to undo migration]

**Data Backfill Requirements:**
- [What data needs to be populated]
- [Estimated duration]

## Performance Considerations

### Expected Data Volumes

| Entity | Initial | 1 Year | 3 Years |
|--------|---------|--------|---------|
| [entity] | [count] | [count] | [count] |

### Partitioning Strategy

[Describe partitioning approach if applicable]

### Caching Strategy

[Describe what data is cached and where]

## Security

### Sensitive Data

| Column | Classification | Protection |
|--------|---------------|------------|
| [table.column] | [PII | Financial | PHI] | [encryption | masking | access control] |

### Access Control

[Row-level security, column-level permissions, etc.]
```

## Usage Instructions

1. Copy the template above
2. Fill in each section based on your domain requirements
3. Remove sections that do not apply to your data store type
4. Add custom sections for domain-specific concerns
5. Keep this document in version control alongside migration scripts
6. Update when schema evolves

## Section Guidelines

### Entities Section
- List all entities, even lookup/reference tables
- Include all columns with accurate types
- Document all constraints including CHECK constraints
- Note default values and auto-generation rules

### Relationships Section
- Capture business meaning, not just technical linkage
- Document cascade behaviors explicitly
- Note any orphan-prevention strategies

### Access Patterns Section
- Focus on the 80% use cases
- Include realistic query examples
- Note if patterns require denormalization

### Migration Plan Section
- Always include rollback procedures
- Estimate data backfill times
- Note any downtime requirements

## Examples

### Simple Entity Definition

```markdown
### User

**Purpose:** Registered users of the application

**Table:** `users`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| email | VARCHAR(255) | NOT NULL, UNIQUE | Login email address |
| password_hash | VARCHAR(255) | NOT NULL | Bcrypt hashed password |
| display_name | VARCHAR(100) | NOT NULL | Public display name |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | Account creation time |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | Last modification time |
| deleted_at | TIMESTAMP | NULL | Soft delete timestamp |

**Indexes:**
- `idx_users_email` on (email) - Login lookup
- `idx_users_deleted_at` on (deleted_at) WHERE deleted_at IS NULL - Active user queries

**Constraints:**
- Email must be valid format (CHECK constraint or application-level)
- Display name between 2-100 characters
```

### Relationship Definition

```markdown
### User Orders

**Type:** 1:N
**From:** User
**To:** Order
**Implementation:** `orders.user_id` references `users.id`
**Cascade Behavior:** ON DELETE RESTRICT (prevent user deletion with orders)
**Business Rule:** Every order must be placed by a registered user. Users may have zero or more orders.
```

### NoSQL Variant

```markdown
### Order (Document)

**Purpose:** Customer purchase transactions

**Collection:** `orders`

```json
{
  "_id": "ObjectId",
  "order_number": "string (unique)",
  "customer": {
    "id": "string (reference)",
    "name": "string (embedded snapshot)",
    "email": "string (embedded snapshot)"
  },
  "items": [
    {
      "product_id": "string (reference)",
      "name": "string (embedded)",
      "quantity": "number",
      "unit_price": "decimal"
    }
  ],
  "totals": {
    "subtotal": "decimal",
    "tax": "decimal",
    "total": "decimal"
  },
  "status": "string (enum)",
  "created_at": "date",
  "updated_at": "date"
}
```

**Indexes:**
- `order_number_1` (unique) - Order lookup by number
- `customer.id_1` - Customer order history
- `status_1_created_at_-1` - Order status queries

**Embedding Rationale:**
- Customer name/email embedded as snapshot at order time (may differ from current)
- Items embedded because they belong exclusively to this order
- Product details embedded because price/name at purchase may differ from current
```
