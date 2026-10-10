const test = require("node:test");
const assert = require("node:assert");
const { parseDuration, parseSize, parseLogLevel, parsePort } = require("../src/utils/parser");

test("parseDuration converts units to milliseconds", () => {
  assert.strictEqual(parseDuration("250"), 250);
  assert.strictEqual(parseDuration("250ms"), 250);
  assert.strictEqual(parseDuration("3s"), 3000);
  assert.strictEqual(parseDuration("2m"), 120000);
  assert.strictEqual(parseDuration("1h"), 3600000);
  assert.strictEqual(parseDuration("soon"), null);
});

test("parseSize converts units to bytes", () => {
  assert.strictEqual(parseSize("512"), 512);
  assert.strictEqual(parseSize("1.5 KB"), 1536);
  assert.strictEqual(parseSize("2mb"), 2097152);
  assert.strictEqual(parseSize("1GB"), 1073741824);
  assert.strictEqual(parseSize("big"), null);
});

test("parseLogLevel maps names to numeric levels", () => {
  assert.strictEqual(parseLogLevel("ERROR"), 50);
  assert.strictEqual(parseLogLevel("warn"), 40);
  assert.strictEqual(parseLogLevel("info"), 30);
  assert.strictEqual(parseLogLevel("debug"), 20);
  assert.strictEqual(parseLogLevel("verbose"), 30);
});

test("parsePort accepts only valid TCP ports", () => {
  assert.strictEqual(parsePort("8080"), 8080);
  assert.strictEqual(parsePort(0), null);
  assert.strictEqual(parsePort(65536), null);
  assert.strictEqual(parsePort("http"), null);
});
