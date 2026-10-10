const test = require("node:test");
const assert = require("node:assert");
const { toCents, formatMoney, sumCents } = require("../src/utils/currency");

test("toCents converts whole amounts", () => {
  assert.strictEqual(toCents(12), 1200);
});

test("toCents rounds fractional amounts to the nearest cent", () => {
  assert.strictEqual(toCents(19.99), 1999);
});

test("formatMoney renders two decimals and the currency", () => {
  assert.strictEqual(formatMoney(1999, "EUR"), "19.99 EUR");
});

test("formatMoney groups thousands", () => {
  assert.strictEqual(formatMoney(123456789, "EUR"), "1,234,567.89 EUR");
});

test("sumCents adds values", () => {
  assert.strictEqual(sumCents([100, 250]), 350);
});

test("sumCents of an empty list is zero", () => {
  assert.strictEqual(sumCents([]), 0);
});
