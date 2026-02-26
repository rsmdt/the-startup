# Tactical Patterns

## Entities

Objects with identity that persists over time. Equality is based on identity, not attributes.

**Characteristics:**
- Has a unique identifier
- Mutable state
- Lifecycle (created, modified, archived)
- Equality by ID

```typescript
class Order {
  private readonly id: OrderId;      // Identity - immutable
  private status: OrderStatus;        // State - mutable
  private items: OrderItem[];         // State - mutable

  constructor(id: OrderId) {
    this.id = id;
    this.status = OrderStatus.Draft;
    this.items = [];
  }

  equals(other: Order): boolean {
    return this.id.equals(other.id);  // Equality by identity
  }
}
```

## Value Objects

Objects without identity. Equality is based on attributes. Always immutable.

**Characteristics:**
- No unique identifier
- Immutable (all properties readonly)
- Equality by attributes
- Self-validating

```typescript
class Money {
  constructor(
    public readonly amount: number,
    public readonly currency: Currency
  ) {
    if (amount < 0) throw new Error('Amount cannot be negative');
  }

  add(other: Money): Money {
    if (!this.currency.equals(other.currency)) {
      throw new Error('Cannot add different currencies');
    }
    return new Money(this.amount + other.amount, this.currency);
  }

  equals(other: Money): boolean {
    return this.amount === other.amount &&
           this.currency.equals(other.currency);
  }
}
```

### When to Use Value Objects vs Entities

| Use Value Object | Use Entity |
|------------------|------------|
| No need to track over time | Need to track lifecycle |
| Interchangeable instances | Unique identity matters |
| Defined by attributes | Defined by continuity |
| Examples: Money, Address, DateRange | Examples: User, Order, Account |

## Aggregates

A cluster of entities and value objects with a defined boundary. One entity is the aggregate root.

### Aggregate Design Rules

1. PROTECT invariants at aggregate boundary
2. REFERENCE other aggregates by identity only
3. UPDATE one aggregate per transaction
4. DESIGN small aggregates (prefer single entity)

```
Aggregate: Order
Root: Order (entity)

  +------------------+
  | Order (Root)     | <-- Aggregate Root
  | - orderId        |
  | - customerId ----+--> Reference by ID only
  | - status         |
  +--------+---------+
           |
  +--------v---------+
  | OrderItem        | <-- Inside aggregate
  | - productId -----+--> Reference by ID only
  | - quantity       |
  | - price (Money)  | <-- Value Object
  +------------------+
```

### Aggregate Sizing

**Start small**: Begin with single-entity aggregates. Expand only when invariants require it.

**Signs of too-large aggregate:**
- Frequent optimistic lock conflicts
- Loading too much data for simple operations
- Multiple users editing simultaneously
- Transactional failures across unrelated data

**Signs of too-small aggregate:**
- Invariants not protected
- Business rules scattered across services
- Eventual consistency where immediate is required

## Domain Events

Represent something that happened in the domain. Immutable facts about the past.

```typescript
class OrderPlaced implements DomainEvent {
  readonly eventId = uuid();
  readonly occurredAt = new Date();

  constructor(
    readonly orderId: OrderId,
    readonly customerId: CustomerId,
    readonly items: OrderItemData[],
    readonly totalAmount: Money
  ) {}
}
```

### Naming Convention
- Past tense (OrderPlaced, not PlaceOrder)
- Domain language (not technical)
- Include all relevant data (event is immutable)

### Event Patterns

| Pattern | Description | Use Case |
|---------|-------------|----------|
| **Event Notification** | Minimal data, query for details | Loose coupling |
| **Event-Carried State** | Full data in event | Performance, offline |
| **Event Sourcing** | Events as source of truth | Audit, temporal queries |

## Repositories

Abstract persistence, providing collection-like access to aggregates.

**Principles:**
- One repository per aggregate
- Returns aggregate roots only
- Hides persistence mechanism
- Supports aggregate reconstitution

```typescript
interface OrderRepository {
  findById(id: OrderId): Promise<Order | null>;
  findByCustomer(customerId: CustomerId): Promise<Order[]>;
  save(order: Order): Promise<void>;
  delete(order: Order): Promise<void>;
}

// Implementation hides persistence details
class PostgresOrderRepository implements OrderRepository {
  async findById(id: OrderId): Promise<Order | null> {
    const row = await this.db.query('SELECT * FROM orders WHERE id = $1', [id]);
    return row ? this.reconstitute(row) : null;
  }

  private reconstitute(row: OrderRow): Order {
    // Rebuild aggregate from persistence
  }
}
```
