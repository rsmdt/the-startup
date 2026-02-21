# Anti-Patterns and Implementation Checklist

## Anemic Domain Model

Logic outside domain objects — entities become data bags, services contain all behavior.

```typescript
// Anti-pattern: Logic outside domain objects
class Order {
  id: string;
  items: Item[];
  status: string;
}

class OrderService {
  calculateTotal(order: Order): number { ... }
  validate(order: Order): boolean { ... }
  submit(order: Order): void { ... }
}

// Better: Logic inside domain objects
class Order {
  private items: OrderItem[];
  private status: OrderStatus;

  get total(): Money {
    return this.items.reduce((sum, item) => sum.add(item.subtotal), Money.zero());
  }

  submit(): void {
    this.validate();
    this.status = OrderStatus.Submitted;
  }
}
```

## Large Aggregates

Everything in one aggregate leads to lock conflicts, excessive data loading, and poor performance.

```typescript
// Anti-pattern: Everything in one aggregate
class Customer {
  orders: Order[];           // Could be thousands
  addresses: Address[];
  paymentMethods: PaymentMethod[];
  preferences: Preferences;
  activityLog: Activity[];   // Could be millions
}

// Better: Separate aggregates referenced by ID
class Customer {
  id: CustomerId;
  defaultAddressId: AddressId;
  defaultPaymentMethodId: PaymentMethodId;
}

class Order {
  customerId: CustomerId;    // Reference by ID
}
```

## Primitive Obsession

Using primitive types where value objects express domain concepts more clearly and safely.

```typescript
// Anti-pattern: Primitive types for domain concepts
function createOrder(
  customerId: string,
  productId: string,
  quantity: number,
  price: number,
  currency: string
) { ... }

// Better: Value objects
function createOrder(
  customerId: CustomerId,
  productId: ProductId,
  quantity: Quantity,
  price: Money
) { ... }
```

## Implementation Checklist

### Aggregate Design
- [ ] Single entity can be aggregate root
- [ ] Invariants are protected at boundary
- [ ] Other aggregates referenced by ID only
- [ ] Fits in memory comfortably
- [ ] One transaction per aggregate

### Entity Implementation
- [ ] Has unique identifier
- [ ] Equality based on ID
- [ ] Encapsulates business rules
- [ ] State changes through methods

### Value Object Implementation
- [ ] All properties immutable
- [ ] Equality based on attributes
- [ ] Self-validating
- [ ] Operations return new instances

### Repository Implementation
- [ ] One per aggregate
- [ ] Returns aggregate roots only
- [ ] Hides persistence details
- [ ] Supports queries needed by domain
