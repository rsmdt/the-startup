# ADR-001: Use PostgreSQL as the Primary Datastore

## Status

Accepted

## Context

The platform needs a primary datastore for user accounts, orders, products, and inventory. The data is relational: orders reference users and products, inventory adjustments are tied to orders, and reporting requires joins across multiple entities.

The team evaluated options against the following constraints:

- ACID transactions are required. An order creation must atomically decrement inventory and record the order — partial failures are not acceptable.
- The schema is not yet fully stabilized. Early product development will produce migrations frequently.
- The team has strong SQL skills but limited experience with document or wide-column stores.
- The system is expected to serve up to 10,000 requests per minute within the first 18 months, based on growth projections.
- Hosting budget is approximately $500/month for infrastructure at launch.

Three options were evaluated: PostgreSQL, MongoDB, and DynamoDB.

## Decision

Use PostgreSQL as the primary datastore, hosted on AWS RDS with a single read replica.

The decision is based on:

1. **Transactional correctness.** PostgreSQL provides full ACID transactions across tables. The order creation workflow (insert order, decrement inventory, record payment intent) requires this. MongoDB offers multi-document transactions since 4.0 but they carry a performance penalty and require deliberate session management that the team is not familiar with. DynamoDB transactions are limited to 25 items and do not generalize well to the relational access patterns of this domain.

2. **Team fluency.** All five developers on the current team have production PostgreSQL experience. MongoDB and DynamoDB would require a ramp-up period and introduce risk during the critical early launch phase.

3. **Relational access patterns.** Reporting requirements include queries like "orders by user with product details and shipment status." These are natural SQL joins. Document stores require either embedding (which causes data duplication) or application-level joins (which push complexity into service code).

4. **Operational familiarity.** RDS for PostgreSQL has predictable operational characteristics, automated backups, and point-in-time recovery. The team knows how to monitor, tune, and migrate it. DynamoDB's capacity planning model (provisioned vs. on-demand) and partition key selection require specialized knowledge to get right.

5. **Cost.** An RDS `db.t3.medium` instance with a read replica comes to approximately $180/month. DynamoDB costs are harder to predict at unknown access volumes. MongoDB Atlas at comparable durability settings costs approximately $250/month.

## Consequences

### Positive

- Developers can write and review schema changes using standard SQL migrations (using golang-migrate).
- ACID transactions eliminate a class of consistency bugs that would require compensating transactions or saga patterns with other datastores.
- Rich query capabilities mean reporting can be done directly against the database without an ETL pipeline or separate analytics store for the first year.
- Read replica allows reporting queries to be routed off the primary, protecting transactional write performance.
- pgvector extension is available if the product roadmap adds semantic search or recommendation features.

### Negative

- PostgreSQL is not horizontally scalable for writes. If write throughput exceeds what a single primary can handle, the migration path is complex: sharding, moving to CockroachDB, or extracting high-write domains into separate services with their own datastores.
- The schema is a shared artifact. All services that evolve from this monolith will need to coordinate schema migrations or extract their own databases — this is a known constraint of the modular monolith starting point.
- Full-text search capabilities are limited compared to a dedicated search engine. If product search becomes a core user experience feature, Elasticsearch or OpenSearch will need to be introduced.
- Column-oriented queries (large aggregations over historical data) will be slow at scale. A columnar store or materialized views will be needed for analytics beyond the first year.

### Neutral

- PostgreSQL's connection-per-process model requires connection pooling (PgBouncer or RDS Proxy) before the application scales to more than ~100 concurrent application instances.
- JSON/JSONB columns are available for semi-structured data, but using them to avoid schema definition should be treated as a deliberate trade-off, not a default, since it sacrifices query performance and type safety.

## Alternatives Considered

### MongoDB Atlas

- Pros: Flexible schema useful during early product discovery; horizontal write scaling via sharding; native document model fits product catalog with variable attribute sets.
- Cons: Multi-document transactions add complexity to the order creation workflow; team has no production MongoDB experience; costs more at equivalent durability; joins require application-level code.
- Why rejected: The consistency requirements of the financial transaction workflow outweigh the schema flexibility benefits. The team would spend more time learning the operational model than building product.

### DynamoDB

- Pros: Fully managed with no connection limits; seamlessly handles unpredictable traffic spikes; pay-per-use pricing works well for uncertain early workloads.
- Cons: Requires upfront access pattern design that is difficult to change later; no ad hoc queries; 25-item transaction limit does not fit the domain; requires DynamoDB-specific expertise for capacity planning and key design.
- Why rejected: The rigid access pattern requirement makes it unsuitable for a domain in early discovery. The cost of a wrong key design is a full table rebuild. PostgreSQL's flexibility is worth paying for at this stage.
