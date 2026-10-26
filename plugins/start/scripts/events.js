const fs = require("fs");
const path = require("path");

const EVENTS_FILE = "events.jsonl";

function appendEvent(home, event) {
  fs.mkdirSync(home, { recursive: true });
  fs.appendFileSync(path.join(home, EVENTS_FILE), `${JSON.stringify({ ts: new Date().toISOString(), ...event })}\n`);
}

module.exports = { EVENTS_FILE, appendEvent };
