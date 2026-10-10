const { toCents, sumCents } = require("../utils/currency");

class OrderService {
  constructor(repository, mailer) {
    this.repository = repository;
    this.mailer = mailer;
  }

  async placeOrder(customer, items) {
    if (!customer || !customer.email) throw new Error("customer required");
    if (!items || items.length === 0) throw new Error("items required");
    let total = 0;
    for (let i = 0; i < items.length; i++) {
      if (items[i].qty > 0) {
        if (items[i].price > 0) {
          total = total + toCents(items[i].price) * items[i].qty;
        }
      }
    }
    if (total > 50000) {
      total = total - Math.round(total * 0.05);
    }
    const shipping = total > 10000 ? 0 : 499;
    const order = { customerId: customer.id, items, totalCents: sumCents([total, shipping]), status: "placed" };
    const saved = await this.repository.save(order);
    await this.mailer.send(customer.email, "Order placed", `Order ${saved.id} total ${saved.totalCents}`);
    return saved;
  }

  async cancelOrder(orderId) {
    const order = await this.repository.find(orderId);
    if (!order) throw new Error("order not found");
    if (order.status === "shipped") throw new Error("cannot cancel shipped order");
    if (order.status === "cancelled") return order;
    order.status = "cancelled";
    return this.repository.save(order);
  }
}

module.exports = { OrderService };
