#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const { EVENTS_FILE, appendEvent } = require("./events");
const { loadTasks, taskEvent } = require("./task");

const SKILLS_DIR = path.resolve(__dirname, "..", "skills");
const SPECS_DIR = path.join(".start", "specs");
const LEGACY_SPECS_DIR = path.join("docs", "specs");
const SPEC_ROOTS = [SPECS_DIR, LEGACY_SPECS_DIR];

const SPEC_FILE = "spec.md";
const SPEC_TEMPLATE = path.join(SKILLS_DIR, "specify-meta", "template.md");

const GATE_FILES = [
  ["definition-of-ready.md", "definition_of_ready"],
  ["definition-of-done.md", "definition_of_done"],
  ["task-definition-of-done.md", "task_definition_of_done"],
];

const GATE_PARTS = ["requirements", "solution", "plan"];
const GATE_VALUES = ["pending", "draft", "approved", "skipped"];
const ASSUMPTION_STATUSES = ["proposed", "confirmed", "corrected", "standing"];
const SLICE_RESULTS = ["completed", "incomplete-coverage", "failed", "stalled", "blocked-environment", "harness-abort"];

const EVENT_SCHEMAS = {
  gate: { part: GATE_PARTS, to: GATE_VALUES },
  requirement: { id: "string", acs: "list" },
  decision: { id: "string", "title?": "string", confirmed: "boolean" },
  question: { id: "string", "text?": "string", resolved: "boolean" },
  assumption: { id: "string", "text?": "string", "observable?": "boolean", status: ASSUMPTION_STATUSES },
  result: { id: "string", result: SLICE_RESULTS, "business?": "number", "technical?": "number", "iterations?": "number" },
  run: { state: ["started", "ended"], "threshold?": "number", "maxIterations?": "number" },
  verify: { constitution: ["pass", "fail", "none"] },
  archived: {},
};

function fail(message) {
  process.stderr.write(`Error: ${message}\n`);
  process.exit(1);
}

function isDir(p) {
  return fs.existsSync(p) && fs.statSync(p).isDirectory();
}

function isFile(p) {
  return fs.existsSync(p) && fs.statSync(p).isFile();
}

function sortedEntries(dir) {
  return fs.readdirSync(dir).sort();
}

function filesIn(dir, pattern = () => true) {
  return sortedEntries(dir).filter((name) => isFile(path.join(dir, name)) && pattern(name));
}

function dirsIn(dir) {
  return sortedEntries(dir).filter((name) => isDir(path.join(dir, name)));
}

function nextSpecId() {
  let max = 0;
  for (const root of SPEC_ROOTS) {
    if (!isDir(root)) continue;
    for (const name of dirsIn(root)) {
      const match = /^(\d{3})-/.exec(name);
      if (match) max = Math.max(max, Number(match[1]));
    }
  }
  return String(max + 1).padStart(3, "0");
}

function sanitizeName(name) {
  return name.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");
}

function findSpecDir(specId) {
  for (const root of SPEC_ROOTS) {
    if (!isDir(root)) continue;
    const name = dirsIn(root).find((n) => n.startsWith(`${specId}-`));
    if (name) return path.join(root, name);
  }
  return null;
}

function firstExisting(specDir, ...candidates) {
  const found = candidates.map((name) => path.join(specDir, name)).find(isFile);
  return found ?? null;
}

function tomlList(items) {
  return `[${items.map((i) => `"${i}"`).join(", ")}]`;
}

function readSpec(specId) {
  const specDir = findSpecDir(specId);
  if (!specDir) fail(`Spec ${specId} not found`);

  const out = [];
  out.push(`id = "${specId}"`);
  out.push(`name = "${path.basename(specDir).slice(specId.length + 1)}"`);
  out.push(`dir = "${specDir}"`);
  out.push("", "[spec]");

  const spec = firstExisting(specDir, SPEC_FILE);
  const prd = spec ?? firstExisting(specDir, "requirements.md", "product-requirements.md");
  const sdd = spec ?? firstExisting(specDir, "solution.md", "solution-design.md");
  const plan = spec ?? firstExisting(specDir, "plan.md");
  if (spec) out.push(`spec = "${spec}"`);
  if (prd) out.push(`prd = "${prd}"`);
  if (sdd) out.push(`sdd = "${sdd}"`);
  if (plan) out.push(`plan_doc = "${plan}"`);

  const planDir = path.join(specDir, "plan");
  if (isDir(planDir)) {
    out.push(`plan_dir = "${planDir}"`);
    const readme = path.join(planDir, "README.md");
    if (isFile(readme)) out.push(`plan = "${readme}"`);
    const phases = filesIn(planDir, (n) => /^phase-.*\.md$/.test(n)).map((n) => path.join(planDir, n));
    if (phases.length) out.push(`phases = ${tomlList(phases)}`);
  }

  const manifest = path.join(specDir, "manifest.md");
  if (isFile(manifest)) out.push(`manifest = "${manifest}"`);

  const unitsDir = path.join(specDir, "units");
  if (isDir(unitsDir)) {
    const units = filesIn(unitsDir, (n) => n.endsWith(".md")).map((n) => path.join(unitsDir, n));
    if (units.length) out.push(`units = ${tomlList(units)}`);
  }

  const scenariosDir = path.join(specDir, "scenarios");
  if (isDir(scenariosDir)) {
    for (const group of dirsIn(scenariosDir)) {
      const groupDir = path.join(scenariosDir, group);
      const scenarios = filesIn(groupDir, (n) => n.endsWith(".md")).map((n) => path.join(groupDir, n));
      if (scenarios.length) out.push(`scenarios_${group} = ${tomlList(scenarios)}`);
    }
  }

  const gates = GATE_FILES.filter(([file]) => isFile(path.join(specDir, file)));
  if (gates.length) {
    out.push("", "[gates]");
    for (const [file, key] of gates) out.push(`${key} = "${path.join(specDir, file)}"`);
  }

  const files = filesIn(specDir);
  for (const sub of ["plan", "units"]) {
    const subDir = path.join(specDir, sub);
    if (isDir(subDir)) files.push(...filesIn(subDir).map((n) => `${sub}/${n}`));
  }
  if (isDir(scenariosDir)) {
    for (const group of dirsIn(scenariosDir)) {
      files.push(...filesIn(path.join(scenariosDir, group)).map((n) => `scenarios/${group}/${n}`));
    }
  }
  out.push("", "files = [");
  files.forEach((file, i) => out.push(`  "${file}"${i < files.length - 1 ? "," : ""}`));
  out.push("]");

  process.stdout.write(out.join("\n") + "\n");
}

function createSpec(featureName) {
  if (!isFile(SPEC_TEMPLATE)) fail(`Template not found: ${SPEC_TEMPLATE}`);
  const specId = nextSpecId();
  const specDir = path.join(SPECS_DIR, `${specId}-${sanitizeName(featureName)}`);
  fs.mkdirSync(path.join(specDir, ".tracker"), { recursive: true });
  fs.writeFileSync(path.join(specDir, SPEC_FILE), fs.readFileSync(SPEC_TEMPLATE, "utf8"));
  console.log(`Created spec directory: ${specDir}`);
  console.log(`Created spec document: ${path.join(specDir, SPEC_FILE)}`);
  console.log(`Created tracker directory: ${path.join(specDir, ".tracker")}`);
  appendEvent(specDir, { kind: "created", name: path.basename(specDir) });
  console.log(`Spec ID: ${specId}`);
  console.log("Specification directory created successfully");
}

function fieldProblem(name, value, rule) {
  if (Array.isArray(rule)) return rule.includes(value) ? null : `${name} must be one of ${rule.join("|")}`;
  if (rule === "list") return Array.isArray(value) && value.every((v) => typeof v === "string") ? null : `${name} must be a list of strings`;
  return typeof value === rule ? null : `${name} must be a ${rule}`;
}

function validateEvent(event) {
  const schema = EVENT_SCHEMAS[event.kind];
  if (!schema) fail(`unknown event kind '${event.kind}' (${Object.keys(EVENT_SCHEMAS).join("|")})`);
  for (const [key, rule] of Object.entries(schema)) {
    const isOptional = key.endsWith("?");
    const name = isOptional ? key.slice(0, -1) : key;
    if (event[name] === undefined) {
      if (isOptional) continue;
      fail(`${event.kind} event requires '${name}'`);
    }
    const problem = fieldProblem(name, event[name], rule);
    if (problem) fail(problem);
  }
}

function recordEvent(specId, json) {
  const specDir = findSpecDir(specId);
  if (!specDir) fail(`Spec ${specId} not found`);
  let event;
  try {
    event = JSON.parse(json);
  } catch (e) {
    fail(`event is not valid JSON: ${e.message}`);
  }
  if (!event || typeof event !== "object" || Array.isArray(event)) fail("event must be a JSON object");
  delete event.ts;
  validateEvent(event);
  appendEvent(specDir, event);
  console.log(`Recorded ${event.kind} event in ${path.join(specDir, EVENTS_FILE)}`);
}

function sectionOf(text, heading) {
  const start = text.indexOf(`\n${heading}\n`);
  if (start < 0) return "";
  const rest = text.slice(start + heading.length + 2);
  const level = heading.indexOf(" ");
  const end = rest.search(new RegExp(`\\n#{1,${level}} `));
  return end < 0 ? rest : rest.slice(0, end);
}

function tableRows(section) {
  return section
    .split("\n")
    .filter((line) => line.startsWith("|") && !/^\|[\s:|-]+\|$/.test(line))
    .map((line) => line.slice(1, line.endsWith("|") ? -1 : undefined).split("|").map((cell) => cell.trim()));
}

function gateEvents(text) {
  const frontmatter = /^---\n([\s\S]*?)\n---/.exec(text);
  if (!frontmatter) return [];
  return [...frontmatter[1].matchAll(/^\s+(requirements|solution|plan):\s*(\w+)\s*$/gm)]
    .filter(([, , to]) => GATE_VALUES.includes(to))
    .map(([, part, to]) => ({ kind: "gate", part, to }));
}

function requirementEvents(text) {
  const acsByFeature = new Map();
  let feature = null;
  for (const line of sectionOf(text, "## Requirements").split("\n")) {
    const heading = /^#{2,6}\s+(FR-\d+)\b/.exec(line);
    if (heading) {
      feature = heading[1];
      acsByFeature.set(feature, []);
    } else if (feature) {
      for (const [, ac] of line.matchAll(/\*\*(AC-\d+\.\d+)\*\*/g)) acsByFeature.get(feature).push(ac);
    }
  }
  return [...acsByFeature].map(([id, acs]) => ({ kind: "requirement", id, acs }));
}

function decisionEvents(text) {
  return text
    .split(/\*\*Decision \(/)
    .slice(1)
    .map((block) => /^(ADR-\d+)\):\s*(.+?)\*\*([\s\S]*?)User confirmed/.exec(block))
    .filter(Boolean)
    .map(([, id, title, between]) => ({ kind: "decision", id, title: title.trim(), confirmed: /\[x\]\s*$/i.test(between.trimEnd()) }));
}

function assumptionEvents(text) {
  return tableRows(sectionOf(text, "### Assumptions"))
    .filter((cells) => /^A\d+$/.test(cells[0]) && cells.length >= 5)
    .map((cells) => ({
      id: cells[0],
      text: cells[1],
      observable: cells[cells.length - 2].toLowerCase() === "yes",
      status: cells[cells.length - 1].split(/\W/)[0].toLowerCase(),
    }))
    .filter((row) => ASSUMPTION_STATUSES.includes(row.status))
    .map((row) => ({ kind: "assumption", ...row }));
}

function questionEvents(text) {
  return tableRows(sectionOf(text, "### Open questions"))
    .filter((cells) => cells[0] !== "Question" && !cells[0].startsWith("[") && /^\[[ x]\]$/i.test(cells[cells.length - 1]))
    .map((cells, index) => ({ kind: "question", id: `Q${index + 1}`, text: cells[0], resolved: /x/i.test(cells[cells.length - 1]) }));
}

function snapshotTaskEvent(task) {
  const event = taskEvent(task);
  if (event.type === "slice") {
    const covers = [...new Set([...task.body.matchAll(/\*\*(AC-\d+\.\d+)\*\*/g)].map(([, ac]) => ac))];
    if (covers.length) event.covers = covers;
  }
  if (event.type === "scenario") {
    const priority = /Priority:\s*(P\d)/.exec(task.body);
    const surface = /Surface:\s*(business|technical)/.exec(task.body);
    if (priority) event.priority = priority[1];
    if (surface) event.surface = surface[1];
  }
  return event;
}

function snapshotSpec(specId) {
  const specDir = findSpecDir(specId);
  if (!specDir) fail(`Spec ${specId} not found`);
  const events = [{ kind: "created", name: path.basename(specDir) }];
  const specFile = path.join(specDir, SPEC_FILE);
  if (isFile(specFile)) {
    const text = fs.readFileSync(specFile, "utf8");
    events.push(...gateEvents(text), ...requirementEvents(text), ...decisionEvents(text), ...assumptionEvents(text), ...questionEvents(text));
  }
  events.push(...loadTasks(specDir).map(snapshotTaskEvent));
  for (const event of events) appendEvent(specDir, event);
  console.log(`Recorded ${events.length} events in ${path.join(specDir, EVENTS_FILE)}`);
}

function parseCli(argv) {
  const args = { feature: null, read: false, event: null, snapshot: false };
  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (arg === "--read") args.read = true;
    else if (arg === "--snapshot") args.snapshot = true;
    else if (arg === "--event") {
      args.event = argv[++i];
      if (args.event === undefined) fail("--event requires a JSON object");
    } else if (arg === "-h" || arg === "--help") {
      console.log("usage: spec.js <feature-name> | spec.js <spec-id> --read | spec.js <spec-id> --event '<json>' | spec.js <spec-id> --snapshot");
      process.exit(0);
    } else if (args.feature === null) args.feature = arg;
    else fail(`unexpected argument: ${arg}`);
  }
  if (args.feature === null) fail("feature name or spec ID is required");
  return args;
}

const cli = parseCli(process.argv.slice(2));
if (cli.read) readSpec(cli.feature);
else if (cli.event !== null) recordEvent(cli.feature, cli.event);
else if (cli.snapshot) snapshotSpec(cli.feature);
else createSpec(cli.feature);
