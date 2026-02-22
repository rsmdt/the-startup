# Denormalization Strategies

Intentional denormalization patterns for read performance optimization.

---

## Calculated Columns

Store derived values to avoid repeated computation.

```
Order
  - subtotal (calculated once on item changes)
  - tax_amount (calculated once)
  - total (calculated once)
```

**Trade-off:** Faster reads, more complex writes, potential consistency issues.

## Materialized Relationships

Embed frequently-accessed related data.

```
Post
  - author_id
  - author_name (copied from User.name)
  - author_avatar_url (copied from User.avatar_url)
```

**Trade-off:** Eliminates joins, requires synchronization on source changes.

## Aggregation Tables

Pre-compute summaries for reporting.

```
DailySales
  - date
  - product_id
  - units_sold (sum)
  - revenue (sum)
```

**Trade-off:** Fast analytics, storage overhead, stale until refreshed.

## Denormalization Decision Matrix

| Factor | Normalize | Denormalize |
|--------|-----------|-------------|
| Write frequency | High | Low |
| Read frequency | Low | High |
| Data consistency | Critical | Eventual OK |
| Query complexity | Simple | Complex joins |
| Data size | Small | Large |
