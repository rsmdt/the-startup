const test = require("node:test");
const assert = require("node:assert");
const { Cart } = require("../src/cart");

const pricing = { priceOf: (sku) => ({ apple: 50, pear: 80 })[sku] };

test("total sums price times quantity", () => {
  const cart = new Cart(pricing);
  cart.add("apple", 2);
  cart.add("pear", 1);
  assert.strictEqual(cart.total(), 180);
});

test("add rejects non-positive quantities", () => {
  const cart = new Cart(pricing);
  assert.throws(() => cart.add("apple", 0), RangeError);
});

test("total matches its own computation", () => {
  const cart = new Cart(pricing);
  cart.add("apple", 3);
  const expected = cart.items.reduce((sum, item) => sum + pricing.priceOf(item.sku) * item.quantity, 0);
  assert.strictEqual(cart.total(), expected);
});

test("Array.prototype.push appends items", () => {
  const items = [];
  items.push({ sku: "apple", quantity: 1 });
  assert.strictEqual(items.length, 1);
});

test("total asks pricing for each sku", (t) => {
  const priceOf = t.mock.fn(() => 10);
  const cart = new Cart({ priceOf });
  cart.add("apple", 1);
  cart.add("pear", 1);
  cart.total();
  assert.deepStrictEqual(priceOf.mock.calls.map((call) => call.arguments[0]), ["apple", "pear"]);
});
