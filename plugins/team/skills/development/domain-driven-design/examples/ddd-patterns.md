# DDD Pattern Implementation Examples

## Context

This guide provides concrete implementation examples for each DDD tactical pattern. Use these when translating domain concepts into code, reviewing existing domain models, or teaching DDD practices to a team. Examples use TypeScript pseudocode with the e-commerce domain as a consistent thread.

---

## Entities

Entities have identity that persists over time. Two entities are equal if they share the same ID, regardless of attribute differences.

### Correct Implementation

```typescript
class Order {
  private readonly id: OrderId;
  private status: OrderStatus;
  private items: OrderItem[];
  private customerId: CustomerId;  // Reference by ID — never embed Customer object

  constructor(id: OrderId, customerId: CustomerId) {
    this.id = id;
    this.customerId = customerId;
    this.status = OrderStatus.Draft;
    this.items = [];
  }

  addItem(productId: ProductId, quantity: Quantity, price: Money): void {
    this.guardDraftStatus('add items');
    const existing = this.items.find(i => i.productId.equals(productId));
    if (existing) {
      existing.increaseQuantity(quantity);
    } else {
      this.items.push(new OrderItem(productId, quantity, price));
    }
  }

  submit(): void {
    this.guardDraftStatus('submit');
    if (this.items.length === 0) throw new Error('Cannot submit empty order');
    this.status = OrderStatus.Placed;
    this.addEvent(new OrderPlaced(this.id, this.customerId, this.items, this.total));
  }

  get total(): Money {
    return this.items.reduce(
      (sum, item) => sum.add(item.subtotal),
      Money.zero('USD')
    );
  }

  equals(other: Order): boolean {
    return this.id.equals(other.id);  // Identity-based equality
  }

  private guardDraftStatus(operation: string): void {
    if (this.status !== OrderStatus.Draft) {
      throw new Error(`Cannot ${operation} on a ${this.status} order`);
    }
  }
}
```

### Anti-Pattern: Anemic Entity

```typescript
// WRONG: State exposed, no behavior, logic lives elsewhere
class Order {
  id: string;
  status: string;
  items: any[];
  customerId: string;
}

// WRONG: OrderService doing what Order should do
class OrderService {
  submit(order: Order): void {
    if (order.status !== 'draft') throw new Error('...');
    if (order.items.length === 0) throw new Error('...');
    order.status = 'placed';  // Mutating from outside
  }
}
```

**Why it matters**: Business rule "cannot submit an empty order" lives in `OrderService`. When a second service also needs to submit orders, the rule either gets duplicated or the service becomes a dependency. The entity cannot protect its own invariants.

---

## Value Objects

Value objects describe characteristics. Two value objects are equal if all their attributes match. They are always immutable — operations return new instances.

### Correct Implementation: Money

```typescript
class Money {
  private readonly amount: number;
  private readonly currency: Currency;

  constructor(amount: number, currency: Currency) {
    if (amount < 0) throw new Error(`Amount cannot be negative: ${amount}`);
    if (!currency) throw new Error('Currency is required');
    this.amount = Math.round(amount * 100) / 100;  // Store in cents precision
    this.currency = currency;
  }

  static zero(currency: string): Money {
    return new Money(0, Currency.of(currency));
  }

  add(other: Money): Money {
    this.guardSameCurrency(other);
    return new Money(this.amount + other.amount, this.currency);
  }

  subtract(other: Money): Money {
    this.guardSameCurrency(other);
    const result = this.amount - other.amount;
    return new Money(result, this.currency);
  }

  multiply(factor: number): Money {
    if (factor < 0) throw new Error('Cannot multiply by negative factor');
    return new Money(this.amount * factor, this.currency);
  }

  isGreaterThan(other: Money): boolean {
    this.guardSameCurrency(other);
    return this.amount > other.amount;
  }

  format(): string {
    return `${this.currency.symbol}${this.amount.toFixed(2)}`;
  }

  equals(other: Money): boolean {
    return this.amount === other.amount && this.currency.equals(other.currency);
  }

  private guardSameCurrency(other: Money): void {
    if (!this.currency.equals(other.currency)) {
      throw new Error(`Currency mismatch: ${this.currency} vs ${other.currency}`);
    }
  }
}
```

### Correct Implementation: Address

```typescript
class Address {
  constructor(
    public readonly street: string,
    public readonly city: string,
    public readonly postalCode: PostalCode,
    public readonly country: Country
  ) {
    if (!street.trim()) throw new Error('Street is required');
    if (!city.trim()) throw new Error('City is required');
  }

  // Returns new instance — does not mutate
  withCity(city: string): Address {
    return new Address(this.street, city, this.postalCode, this.country);
  }

  equals(other: Address): boolean {
    return (
      this.street === other.street &&
      this.city === other.city &&
      this.postalCode.equals(other.postalCode) &&
      this.country.equals(other.country)
    );
  }
}
```

### Anti-Pattern: Primitive Obsession

```typescript
// WRONG: Primitive types for domain concepts
function placeOrder(
  customerId: string,         // Is this a UUID? An email? A user number?
  shippingStreet: string,
  shippingCity: string,
  shippingZip: string,        // Is zip validation done? Where?
  price: number,              // USD? EUR? Cents? Dollars?
  currency: string
): void { ... }

// CORRECT: Value objects carry meaning and validate themselves
function placeOrder(
  customerId: CustomerId,
  shippingAddress: Address,   // Validated, complete, typed
  price: Money                // Currency-aware, non-negative
): void { ... }
```

**Why it matters**: `string` for a postal code accepts `"hello"`. `PostalCode` validates format at construction. The type system communicates intent and prevents entire classes of bugs.

---

## Aggregates

An aggregate is a cluster of domain objects treated as a single unit. The aggregate root is the only entry point — external code cannot reach inside to modify children directly.

### Correct Implementation: Order Aggregate

```typescript
class Order {  // Aggregate root
  private readonly id: OrderId;
  private items: OrderItem[];  // Part of aggregate — not exposed directly
  private status: OrderStatus;
  private readonly events: DomainEvent[] = [];

  // Only the root exposes behavior — callers never touch OrderItem directly
  addItem(productId: ProductId, quantity: Quantity, unitPrice: Money): void {
    this.guardCanModify();
    const item = this.findItem(productId);
    if (item) {
      item.increaseQuantity(quantity);
    } else {
      this.items.push(OrderItem.create(productId, quantity, unitPrice));
    }
  }

  removeItem(productId: ProductId): void {
    this.guardCanModify();
    const index = this.items.findIndex(i => i.productId.equals(productId));
    if (index === -1) throw new Error(`Item ${productId} not in order`);
    this.items.splice(index, 1);
  }

  // Invariant: cannot submit with zero items
  // Invariant: cannot submit non-draft order
  submit(): void {
    this.guardCanModify();
    if (this.items.length === 0) {
      throw new DomainError('Order must have at least one item to be submitted');
    }
    this.status = OrderStatus.Placed;
    this.events.push(new OrderPlaced(this.id, this.snapshot()));
  }

  // Read model — callers can observe but not mutate children
  getItems(): ReadonlyArray<OrderItemSnapshot> {
    return this.items.map(i => i.toSnapshot());
  }

  pullEvents(): DomainEvent[] {
    const events = [...this.events];
    this.events.length = 0;
    return events;
  }

  private findItem(productId: ProductId): OrderItem | undefined {
    return this.items.find(i => i.productId.equals(productId));
  }

  private guardCanModify(): void {
    if (this.status !== OrderStatus.Draft) {
      throw new DomainError(`Cannot modify a ${this.status} order`);
    }
  }
}
```

### Aggregate: Cross-Reference by ID

```typescript
// WRONG: Holding object reference across aggregate boundary
class Order {
  customer: Customer;  // Entire Customer aggregate embedded — creates coupling
}

// CORRECT: Reference by identity
class Order {
  customerId: CustomerId;  // Lookup when needed, no coupling to Customer aggregate
}

// WRONG: Order directly modifying inventory
class OrderService {
  async submitOrder(orderId: OrderId): Promise<void> {
    const order = await this.orderRepo.findById(orderId);
    const inventory = await this.inventoryRepo.findByProduct(order.productId);
    order.submit();
    inventory.decrement(order.quantity);          // Two aggregates in one transaction
    await this.db.transaction(() => {
      this.orderRepo.save(order);
      this.inventoryRepo.save(inventory);         // Violates one-aggregate-per-transaction
    });
  }
}

// CORRECT: Order raises event, inventory responds eventually
class OrderService {
  async submitOrder(orderId: OrderId): Promise<void> {
    const order = await this.orderRepo.findById(orderId);
    order.submit();                               // Produces OrderPlaced event
    await this.orderRepo.save(order);             // Single aggregate persisted
    // OrderPlaced event → InventoryHandler reduces stock (eventual consistency)
  }
}
```

---

## Repositories

Repositories abstract persistence. They provide collection-like semantics — you put aggregates in, you get aggregates out — without leaking database concerns into the domain.

### Correct Implementation

```typescript
// Interface in domain layer — no persistence technology references
interface OrderRepository {
  findById(id: OrderId): Promise<Order | null>;
  findByCustomer(customerId: CustomerId, options?: QueryOptions): Promise<Order[]>;
  findPendingFulfillment(): Promise<Order[]>;
  save(order: Order): Promise<void>;
  delete(id: OrderId): Promise<void>;
}

// Implementation in infrastructure layer
class PostgresOrderRepository implements OrderRepository {
  constructor(private readonly db: DatabaseClient) {}

  async findById(id: OrderId): Promise<Order | null> {
    const row = await this.db.queryOne(
      'SELECT * FROM orders WHERE id = $1',
      [id.value]
    );
    return row ? this.reconstitute(row) : null;
  }

  async save(order: Order): Promise<void> {
    const data = this.decompose(order);
    await this.db.upsert('orders', data);
    // Dispatch domain events after successful save
    const events = order.pullEvents();
    await this.eventBus.publishAll(events);
  }

  // Reconstitution rebuilds the full aggregate from raw data
  private reconstitute(row: OrderRow): Order {
    return Order.reconstitute({
      id: OrderId.of(row.id),
      customerId: CustomerId.of(row.customer_id),
      status: OrderStatus[row.status],
      items: row.items.map(this.reconstituteItem),
    });
  }

  private decompose(order: Order): OrderRow {
    // Snapshot aggregate to persistence format
    const snapshot = order.toSnapshot();
    return {
      id: snapshot.id,
      customer_id: snapshot.customerId,
      status: snapshot.status,
      items: snapshot.items,
      updated_at: new Date(),
    };
  }
}
```

### Anti-Pattern: Repository as Query Dumping Ground

```typescript
// WRONG: Repository leaking infrastructure concerns into callers
interface OrderRepository {
  findById(id: string): Promise<OrderRow>;         // Returns DB row, not domain object
  executeQuery(sql: string): Promise<any[]>;       // Domain layer writing SQL
  findWithJoin(table: string, on: string): Promise<any>;  // Structural leakage
}

// WRONG: Domain service building queries
class OrderService {
  async getOrdersForDashboard(): Promise<Order[]> {
    return this.db.query(`
      SELECT o.*, c.name, p.title
      FROM orders o
      JOIN customers c ON o.customer_id = c.id
      ...
    `);
  }
}

// CORRECT: Named queries that reveal domain intent
interface OrderRepository {
  findById(id: OrderId): Promise<Order | null>;
  findPendingFulfillment(): Promise<Order[]>;       // Domain language, not SQL
  findByCustomer(customerId: CustomerId): Promise<Order[]>;
}
```

---

## Domain Events

Domain events record facts about the past. They are immutable and named in past tense using domain vocabulary.

### Correct Implementation

```typescript
// Base contract — all events share this structure
interface DomainEvent {
  readonly eventId: string;
  readonly occurredAt: Date;
  readonly aggregateId: string;
  readonly eventType: string;
}

class OrderPlaced implements DomainEvent {
  readonly eventId = crypto.randomUUID();
  readonly occurredAt = new Date();
  readonly eventType = 'OrderPlaced';

  constructor(
    readonly aggregateId: string,    // orderId
    readonly customerId: string,
    readonly items: ReadonlyArray<OrderItemData>,
    readonly totalAmount: MoneyData,
    readonly placedAt: Date          // Domain-meaningful timestamp
  ) {}
}

class PaymentFailed implements DomainEvent {
  readonly eventId = crypto.randomUUID();
  readonly occurredAt = new Date();
  readonly eventType = 'PaymentFailed';

  constructor(
    readonly aggregateId: string,    // orderId
    readonly reason: PaymentFailureReason,
    readonly attemptedAmount: MoneyData,
    readonly failedAt: Date
  ) {}
}

// Aggregate collects events, repository dispatches them after save
class Order {
  private events: DomainEvent[] = [];

  submit(): void {
    this.status = OrderStatus.Placed;
    this.events.push(new OrderPlaced(
      this.id.value,
      this.customerId.value,
      this.items.map(i => i.toData()),
      this.total.toData(),
      new Date()
    ));
  }

  pullEvents(): DomainEvent[] {
    const events = [...this.events];
    this.events.length = 0;
    return events;
  }
}
```

### Event Handlers

```typescript
// Handler in the same service — synchronous or via in-process bus
class InventoryHandler {
  async handle(event: OrderPlaced): Promise<void> {
    for (const item of event.items) {
      const inventory = await this.inventoryRepo.findByProduct(
        ProductId.of(item.productId)
      );
      inventory.reserve(Quantity.of(item.quantity));
      await this.inventoryRepo.save(inventory);
    }
  }
}

// Handler across service boundary — integration event via message broker
class NotificationService {
  async handle(event: OrderPlaced): Promise<void> {
    const customer = await this.customerRepo.findById(
      CustomerId.of(event.customerId)
    );
    await this.emailGateway.sendOrderConfirmation(customer.email, event);
  }
}
```

### Anti-Pattern: Command Masquerading as Event

```typescript
// WRONG: Event telling other contexts what to do
class OrderSubmitted {
  command = 'RESERVE_INVENTORY';   // Not an event — it's an instruction
  productId: string;
  quantity: number;
}

// WRONG: Mutable event (events are immutable facts)
class OrderPlaced {
  orderId: string;
  status: string;      // Mutable — callers could change it
  items: Item[];       // Mutable array
}

// CORRECT: Immutable fact with all data needed for consumers
class OrderPlaced {
  constructor(
    readonly orderId: string,
    readonly items: ReadonlyArray<OrderItemData>,  // Immutable snapshot
    readonly totalAmount: MoneyData,
    readonly placedAt: Date
  ) {}
}
```

---

## Domain Services

Domain services encapsulate business logic that naturally involves multiple aggregates or domain concepts but does not belong to any single one.

### When to Use Domain Services

Use a domain service when:
- The operation spans multiple aggregates
- The operation requires external data to enforce a business rule
- Placing the logic in an entity would create an unnatural dependency

### Correct Implementation

```typescript
// QUESTION: Who calculates shipping cost?
// Not Order — it doesn't know shipping rates.
// Not ShippingRate — it doesn't know order contents.
// A domain service bridges them.

class ShippingCostCalculator {
  constructor(private readonly rateRepository: ShippingRateRepository) {}

  async calculate(order: Order, destination: Address): Promise<Money> {
    const weight = order.totalWeight;
    const zone = ShippingZone.forAddress(destination);
    const rate = await this.rateRepository.findRate(weight, zone);

    if (!rate) {
      throw new DomainError(`No shipping rate for zone ${zone} at weight ${weight}`);
    }

    return rate.applyTo(weight);
  }
}

// Transfer between accounts — neither Account alone can enforce the rule
class FundsTransferService {
  async transfer(
    sourceId: AccountId,
    destinationId: AccountId,
    amount: Money
  ): Promise<void> {
    const source = await this.accountRepo.findById(sourceId);
    const destination = await this.accountRepo.findById(destinationId);

    if (!source.hasSufficientFunds(amount)) {
      throw new InsufficientFundsError(sourceId, amount);
    }

    source.debit(amount);
    destination.credit(amount);

    // Each save publishes its own events — no cross-aggregate transaction
    await this.accountRepo.save(source);
    await this.accountRepo.save(destination);
  }
}
```

### Anti-Pattern: Domain Service as Catch-All

```typescript
// WRONG: Domain service doing what the entity should do
class OrderService {
  calculateTotal(order: Order): Money {   // Order should own this
    return order.items.reduce((sum, item) => sum + item.price, 0);
  }

  isEligibleForDiscount(order: Order): boolean {  // Business rule — belongs in Order
    return order.items.length > 5;
  }

  canBeShipped(order: Order): boolean {   // Order invariant — belongs in Order
    return order.status === 'placed';
  }
}
```

**Why it matters**: When business logic lives in a service instead of the entity, every caller must import the service to perform basic operations. The domain model becomes a passive data structure — the anemic domain model anti-pattern at the service level.

---

## Bounded Contexts

Bounded contexts establish explicit boundaries where a specific domain model and ubiquitous language apply.

### Modeling the Same Concept Differently

```
Scenario: "Product" in an e-commerce platform

┌──────────────────────────┐  ┌──────────────────────────┐  ┌──────────────────────────┐
│   Catalog Context        │  │   Inventory Context      │  │   Pricing Context        │
├──────────────────────────┤  ├──────────────────────────┤  ├──────────────────────────┤
│ Product:                 │  │ StockKeepingUnit:        │  │ PricedItem:              │
│ - name                   │  │ - sku                    │  │ - productId              │
│ - description            │  │ - warehouseLocation      │  │ - basePrice              │
│ - images                 │  │ - quantity               │  │ - discountRules          │
│ - categories             │  │ - reorderThreshold       │  │ - taxCategory            │
│ - specifications         │  │ - supplierId             │  │ - effectiveFrom          │
└──────────────────────────┘  └──────────────────────────┘  └──────────────────────────┘
         │                              │                              │
         └──────────────────────────────┴──────────────────────────────┘
                               Shared identifier: productId
                               Different model, different language, different team
```

### Anti-Corruption Layer

```typescript
// External supplier API speaks a different language
// ACL translates without polluting the domain model

interface SupplierApi {
  getProductData(ean: string): Promise<SupplierProductRecord>;
}

// Raw supplier vocabulary — not your domain's language
interface SupplierProductRecord {
  GTIN: string;
  PROD_NAME: string;
  WAREHOUSE_CODE: string;
  AVAIL_QTY: number;
  REORDER_LVL: number;
}

// ACL translates to your ubiquitous language
class SupplierInventoryAdapter {
  constructor(private readonly supplierApi: SupplierApi) {}

  async getStockLevel(productId: ProductId): Promise<StockLevel> {
    const raw = await this.supplierApi.getProductData(productId.ean);
    return new StockLevel(
      ProductId.of(raw.GTIN),
      Quantity.of(raw.AVAIL_QTY),
      Quantity.of(raw.REORDER_LVL),
      WarehouseCode.of(raw.WAREHOUSE_CODE)
    );
  }
}

// Domain service uses your language — supplier details hidden behind adapter
class InventoryService {
  async checkReorderNeeded(productId: ProductId): Promise<boolean> {
    const stock = await this.supplierAdapter.getStockLevel(productId);
    return stock.isBelowReorderThreshold();
  }
}
```

---

## Common Mistakes Summary

| Mistake | Symptom | Fix |
|---------|---------|-----|
| **Anemic domain model** | Services contain all business logic, entities are data bags | Move business rules into entity methods |
| **Oversized aggregates** | Lock contention, slow loads, unrelated data changes together | Split on invariant boundaries, use eventual consistency |
| **Primitive obsession** | `string` for email, `number` for price, no validation | Create value objects for domain concepts |
| **Cross-aggregate references** | Entity holds reference to another aggregate's object | Reference by ID only |
| **Multiple aggregates per transaction** | Transaction spans two repositories | Use domain events and eventual consistency |
| **Repository as query service** | Repository has 30 methods, returns DTOs or raw rows | Separate read models from aggregate repositories |
| **Commands masquerading as events** | Event names like `SendConfirmationEmail` | Name events in past tense: `OrderPlaced` |
| **Business logic in event handlers** | Handlers make business decisions, not just reactions | Handlers update state only; business rules stay in aggregates |
| **Missing ubiquitous language** | Code uses technical terms, domain experts can't read it | Rename classes and methods to match domain vocabulary |
| **Shared database between contexts** | Two services join tables across context boundaries | Each context owns its data; integrate via events or APIs |
