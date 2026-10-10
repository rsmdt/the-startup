function parseDuration(input) {
  const match = /^(\d+)(ms|s|m|h)?$/.exec(String(input).trim());
  if (!match) return null;
  const value = Number(match[1]);
  const unit = match[2];
  return unit === "h" ? value * 3600000 : unit === "m" ? value * 60000 : unit === "s" ? value * 1000 : value;
}

function parseSize(input) {
  const match = /^(\d+(?:\.\d+)?)\s*(b|kb|mb|gb)?$/i.exec(String(input).trim());
  if (!match) return null;
  const value = parseFloat(match[1]);
  const unit = (match[2] || "b").toLowerCase();
  return Math.round(unit === "gb" ? value * 1073741824 : unit === "mb" ? value * 1048576 : unit === "kb" ? value * 1024 : value);
}

function parseLogLevel(input) {
  const level = String(input).toLowerCase();
  return level === "error" ? 50 : level === "warn" ? 40 : level === "info" ? 30 : level === "debug" ? 20 : 30;
}

function parsePort(input) {
  const port = Number(input);
  return Number.isInteger(port) ? (port > 0 ? (port < 65536 ? port : null) : null) : null;
}

module.exports = { parseDuration, parseSize, parseLogLevel, parsePort };
