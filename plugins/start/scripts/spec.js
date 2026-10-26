#!/usr/bin/env node

const fs = require("fs");
const path = require("path");

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
  console.log(`Spec ID: ${specId}`);
  console.log("Specification directory created successfully");
}

function parseCli(argv) {
  const args = { feature: null, read: false };
  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (arg === "--read") args.read = true;
    else if (arg === "-h" || arg === "--help") {
      console.log("usage: spec.js <feature-name> | spec.js <spec-id> --read");
      process.exit(0);
    } else if (args.feature === null) args.feature = arg;
    else fail(`unexpected argument: ${arg}`);
  }
  if (args.feature === null) fail("feature name or spec ID is required");
  return args;
}

const cli = parseCli(process.argv.slice(2));
if (cli.read) readSpec(cli.feature);
else createSpec(cli.feature);
