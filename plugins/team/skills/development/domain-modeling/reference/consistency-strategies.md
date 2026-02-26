# Consistency Strategies

## Transactional Consistency (ACID)

Use for invariants within an aggregate. Rule: one aggregate per transaction.

```typescript
// Good: Single aggregate updated
async function addItemToOrder(orderId: OrderId, item: OrderItem) {
  const order = await orderRepo.findById(orderId);
  order.addItem(item);  // Business rules enforced
  await orderRepo.save(order);
}

// Bad: Multiple aggregates in one transaction
async function createOrderWithInventory() {
  await db.transaction(async (tx) => {
    await orderRepo.save(order, tx);
    await inventoryRepo.decrement(productId, quantity, tx);  // Don't do this
  });
}
```

## Eventual Consistency

Use for consistency across aggregates. Pattern: domain events + handlers.

```typescript
// Order aggregate publishes event
class Order {
  submit(): void {
    this.status = OrderStatus.Placed;
    this.addEvent(new OrderPlaced(this.id, this.customerId, this.items));
  }
}

// Separate handler updates inventory (eventually)
class InventoryHandler {
  async handle(event: OrderPlaced): Promise<void> {
    for (const item of event.items) {
      await this.inventoryService.reserve(item.productId, item.quantity);
    }
  }
}
```

## Saga Pattern

Coordinate multiple aggregates with compensation. On failure at any step, execute compensation in reverse order.

```
Saga: Order Fulfillment

Create Order --> Reserve Inventory --> Charge Payment --> Ship Order

Compensations (on failure, reverse order):
  Ship Order fails     => Refund Payment
  Charge Payment fails => Release Inventory
  Reserve fails        => Cancel Order
```

## Choosing Consistency

| Scenario | Strategy |
|----------|----------|
| Within single aggregate | Transactional (ACID) |
| Across aggregates, same service | Eventual (domain events) |
| Across services | Saga with compensation |
| Read model updates | Eventual (projection) |
