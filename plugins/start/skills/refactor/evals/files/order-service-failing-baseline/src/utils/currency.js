function toCents(amount) {
  return Math.floor(amount * 100);
}

function formatMoney(cents, currency) {
  return `${(cents / 100).toFixed(2)} ${currency}`;
}

function sumCents(values) {
  return values.reduce((total, value) => total + value);
}

module.exports = { toCents, formatMoney, sumCents };
