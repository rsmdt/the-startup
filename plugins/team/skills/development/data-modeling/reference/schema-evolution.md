# Schema Evolution Strategies

Safe and breaking changes, plus migration patterns for evolving schemas.

---

## Additive Changes (Safe)

- Add new nullable columns
- Add new tables
- Add new indexes
- Add new optional fields (NoSQL)

## Breaking Changes (Require Migration)

- Remove columns/tables
- Rename columns/tables
- Change data types
- Add non-nullable columns without defaults

## Migration Patterns

### Expand-Contract Pattern

1. Add new column alongside old
2. Backfill new column from old
3. Update application to use new column
4. Remove old column

### Blue-Green Schema

1. Create new version of schema
2. Dual-write to both versions
3. Migrate reads to new version
4. Drop old version

### Versioned Documents (NoSQL)

```json
{
  "_schema_version": 2,
  "name": "Jane",
  "email": "jane@example.com"
}
```

Handle multiple versions in application code during transition.
