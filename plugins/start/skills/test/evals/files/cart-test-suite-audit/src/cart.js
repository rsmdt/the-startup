class Cart {
  constructor(pricing) {
    this.pricing = pricing;
    this.items = [];
  }

  add(sku, quantity) {
    if (quantity <= 0) throw new RangeError("quantity must be positive");
    this.items.push({ sku, quantity });
  }

  total() {
    return this.items.reduce((sum, item) => sum + this.pricing.priceOf(item.sku) * item.quantity, 0);
  }
}

module.exports = { Cart };
