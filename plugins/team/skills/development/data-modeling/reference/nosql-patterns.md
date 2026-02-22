# NoSQL Data Modeling Patterns

Patterns for document stores, key-value stores, wide-column stores, and graph databases.

---

## Document Stores (MongoDB, DynamoDB)

### Embedding Pattern

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

### Referencing Pattern

Reference related data when it changes independently or is shared.

```json
{
  "order_id": "123",
  "customer_id": "456",
  "item_ids": ["A1", "B2"]
}
```

### Hybrid Pattern

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

## Key-Value Stores

### Access Pattern Design

Design keys around query patterns.

```
USER:{user_id} -> user data
USER:{user_id}:ORDERS -> list of order ids
ORDER:{order_id} -> order data
```

### Composite Keys

Combine entity type with identifiers for namespacing.

## Wide-Column Stores (Cassandra, HBase)

### Partition Key Design

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

## Graph Databases

### Node and Relationship Design

- Nodes: entities with properties
- Relationships: named, directed, with properties
- Labels: categorize nodes for efficient traversal

```
(User)-[:PURCHASED {date, amount}]->(Product)
(User)-[:FOLLOWS]->(User)
(Product)-[:BELONGS_TO]->(Category)
```
