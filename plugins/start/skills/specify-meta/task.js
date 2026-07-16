#!/usr/bin/env node
/**
 * The Agentic Startup — Tracker CLI
 *
 * Everything is a task. A task with no `parentId` is a root (≈ a phase); a task with a
 * `parentId` is a child; a task with no children is a leaf (the actionable unit). The tree
 * and the dependency graph are VIEWS this CLI assembles by reading every task file's
 * frontmatter — on disk the files are always flat in `.start/specs/[NNN]-name/.tracker/`.
 *
 * Frontmatter (only what the traversal consumes):
 *   id         alpha-numeric, non-ordinal
 *   name       human label
 *   type       "task" (default) | "bug" | …
 *   status     leaf state machine; parents derive status by rollup (see below)
 *   parentId   optional; absent = root
 *   override   optional; when true, a parent uses its stored status instead of the rollup
 *   dependsOn  [] ; edges to other tasks (inherited down the tree)
 * Everything else — description, bug detail — lives in the markdown body.
 *
 * Derived rules (computed here so the LLM only reads a task and acts):
 *   • Parent status = rollup of children: failed if any child failed, passed if all passed,
 *     pending if all pending, else in_progress. `override` makes the stored status win.
 *   • Only leaves are executable. A leaf is `ready` when its own AND all ancestors' deps
 *     are satisfied (dependsOn is inherited down the tree).
 *
 * Commands:
 *   node task.js create --spec 012 --name "…" [--parent <id>] [--type bug] [--depends-on a,b] [--body "…"]
 *   node task.js read   <id> [--spec 012] [--format toml|json]
 *   node task.js update <id> [--spec 012] [--status …] [--parent <id>] [--type …] [--append-note "…"] [--override|--rollup]
 *   node task.js delete <id> [--spec 012] [--cascade]
 *   node task.js list   --spec 012 [--depth N | --level N] [--format json|toml|tree]
 *   node task.js graph  --spec 012 [--format json|toml|dot]
 *
 * Tracker home resolution: pass `--spec NNN` (resolves .start/specs/NNN-*) OR `--home <dir>`
 * for any directory that owns a `.tracker/` (e.g. `.start/analysis/<name>`). task.js is a
 * generic tracker CLI — specs are just one kind of home.
 *
 * Node, zero dependencies. Invoked via the shell by the /start:* skills and workflow agents
 * (NOT by a dynamic-workflow script — those have no shell; only their agents do).
 */

const fs = require("fs");
const path = require("path");
const crypto = require("crypto");
const { parseArgs } = require("util");

const SPECS_ROOTS = process.env.STARTUP_SPECS_ROOT
  ? [process.env.STARTUP_SPECS_ROOT]
  : [".start/specs", "docs/specs"];

const LEAF_STATUSES = ["pending", "building", "built", "evaluating", "passed", "failed", "blocked"];
const SATISFIED = "passed";

function fail(msg) {
  process.stderr.write(`Error: ${msg}\n`);
  process.exit(1);
}

// ── Spec + file resolution ──────────────────────────────────────────────────────
function findSpecDir(spec) {
  if (!spec) return null;
  if (spec.includes(path.sep) && fs.existsSync(spec)) return spec;
  const isIdPrefix = /^\d+$/.test(spec);
  for (const root of SPECS_ROOTS) {
    if (!fs.existsSync(root)) continue;
    for (const entry of fs.readdirSync(root)) {
      const full = path.join(root, entry);
      if (!fs.statSync(full).isDirectory()) continue;
      if (entry === spec || (isIdPrefix && entry.startsWith(`${spec}-`))) return full;
    }
  }
  return null;
}

const trackerDir = (home) => path.join(home, ".tracker");

// Resolve the tracker home: explicit --home wins, else --spec resolves a spec dir.
function homeOf(o) {
  if (o.home) return o.home;
  if (o.spec) {
    const d = findSpecDir(o.spec);
    if (!d) fail(`Spec '${o.spec}' not found`);
    return d;
  }
  return null;
}

function findTaskFile(id, o) {
  const home = homeOf(o);
  if (home) {
    const file = path.join(trackerDir(home), `${id}.md`);
    return fs.existsSync(file) ? { file, home } : null;
  }
  for (const root of SPECS_ROOTS) {
    if (!fs.existsSync(root)) continue;
    for (const entry of fs.readdirSync(root)) {
      const h = path.join(root, entry);
      const file = path.join(trackerDir(h), `${id}.md`);
      if (fs.existsSync(file)) return { file, home: h };
    }
  }
  return null;
}

// ── Frontmatter (de)serialization ────────────────────────────────────────────────
function parseTask(raw) {
  const m = raw.match(/^---\n([\s\S]*?)\n---\n?([\s\S]*)$/);
  if (!m) fail("Malformed task file (no frontmatter)");
  const data = { dependsOn: [] };
  for (const line of m[1].split("\n")) {
    const kv = line.match(/^([a-zA-Z]+):\s*(.*)$/);
    if (!kv) continue;
    const [, key, raw2] = kv;
    const val = raw2.trim();
    if (key === "dependsOn") {
      const inner = val.replace(/^\[/, "").replace(/\]$/, "").trim();
      data.dependsOn = inner ? inner.split(",").map((s) => s.trim()).filter(Boolean) : [];
    } else if (key === "override") {
      data.override = val === "true";
    } else {
      data[key] = val;
    }
  }
  data.body = m[2].trimEnd();
  return data;
}

function serializeTask(t) {
  const lines = ["---", `id: ${t.id}`, `name: ${t.name}`, `type: ${t.type || "task"}`, `status: ${t.status || "pending"}`];
  if (t.parentId) lines.push(`parentId: ${t.parentId}`);
  if (t.override) lines.push("override: true");
  lines.push(`dependsOn: [${(t.dependsOn || []).join(", ")}]`);
  lines.push("---", "");
  return `${lines.join("\n")}${t.body ? `${t.body}\n` : ""}`;
}

const nowIso = () => new Date().toISOString();
const splitList = (v) => (v == null ? undefined : v.split(",").map((s) => s.trim()).filter(Boolean));

function genId(dir) {
  for (;;) {
    const n = crypto.randomBytes(5).readUIntBE(0, 5);
    const id = n.toString(36).padStart(6, "0").slice(-6);
    if (!fs.existsSync(path.join(dir, `${id}.md`))) return id;
  }
}

function loadTasks(specDir) {
  const dir = trackerDir(specDir);
  if (!fs.existsSync(dir)) return [];
  return fs.readdirSync(dir).filter((f) => f.endsWith(".md")).map((f) => parseTask(fs.readFileSync(path.join(dir, f), "utf8")));
}

// ── Derived graph context (rollup, inherited deps, readiness) ─────────────────────
function buildContext(tasks) {
  const byId = new Map(tasks.map((t) => [t.id, t]));
  const childIds = new Map();
  for (const t of tasks) {
    if (t.parentId && byId.has(t.parentId)) {
      if (!childIds.has(t.parentId)) childIds.set(t.parentId, []);
      childIds.get(t.parentId).push(t.id);
    }
  }
  const isRoot = (t) => !t.parentId || !byId.has(t.parentId);
  const isLeaf = (id) => !(childIds.get(id) || []).length;

  const memo = new Map();
  function effStatus(id, seen = new Set()) {
    if (memo.has(id)) return memo.get(id);
    if (seen.has(id)) return "pending"; // cycle guard
    seen.add(id);
    const t = byId.get(id);
    if (!t) return "pending";
    let result;
    if (isLeaf(id)) {
      result = t.status || "pending";
    } else if (t.override) {
      result = t.status || "pending";
    } else {
      const kids = childIds.get(id).map((c) => effStatus(c, seen));
      if (kids.some((s) => s === "failed")) result = "failed";
      else if (kids.every((s) => s === "passed")) result = "passed";
      else if (kids.every((s) => s === "pending")) result = "pending";
      else result = "in_progress";
    }
    memo.set(id, result);
    return result;
  }

  function ancestors(id) {
    const out = [];
    const seen = new Set();
    let cur = byId.get(id);
    while (cur && cur.parentId && byId.has(cur.parentId) && !seen.has(cur.parentId)) {
      seen.add(cur.parentId);
      cur = byId.get(cur.parentId);
      out.push(cur);
    }
    return out;
  }
  const depth = (id) => ancestors(id).length;
  function effDeps(id) {
    const deps = new Set(byId.get(id).dependsOn || []);
    for (const a of ancestors(id)) for (const d of a.dependsOn || []) deps.add(d);
    return [...deps];
  }
  const blockedBy = (id) => effDeps(id).filter((d) => effStatus(d) !== SATISFIED);
  const ready = (id) => isLeaf(id) && (byId.get(id).status || "pending") === "pending" && blockedBy(id).length === 0;

  const node = (t) => ({
    id: t.id, name: t.name, type: t.type || "task", status: effStatus(t.id),
    parentId: t.parentId || null, depth: depth(t.id), leaf: isLeaf(t.id),
    ready: ready(t.id), dependsOn: t.dependsOn || [], blockedBy: blockedBy(t.id),
  });

  return { byId, childIds, isRoot, isLeaf, effStatus, depth, ready, blockedBy, node };
}

function countBy(nodes) {
  const c = nodes.reduce((a, n) => ((a[n.status] = (a[n.status] || 0) + 1), a), {});
  c.total = nodes.length;
  return c;
}

// ── Commands ────────────────────────────────────────────────────────────────────
function cmdCreate(o) {
  const home = homeOf(o);
  if (!home) fail("create requires --home <dir> or --spec <id>");
  if (!o.name) fail("create requires --name <string>");
  const dir = trackerDir(home);
  fs.mkdirSync(dir, { recursive: true });
  if (o.parent && !fs.existsSync(path.join(dir, `${o.parent}.md`))) fail(`Parent task '${o.parent}' not found`);
  const existing = loadTasks(home);
  const known = new Set(existing.map((t) => t.id));
  for (const d of splitList(o["depends-on"]) || []) if (!known.has(d)) fail(`dependsOn target '${d}' not found`);

  const status = o.status || "pending";
  if (!LEAF_STATUSES.includes(status)) fail(`Invalid status '${status}' (${LEAF_STATUSES.join("|")})`);

  const id = genId(dir);
  const task = {
    id, name: o.name, type: o.type || "task", status,
    parentId: o.parent || undefined, dependsOn: splitList(o["depends-on"]) || [],
    body: o["body-file"] ? fs.readFileSync(o["body-file"], "utf8").trimEnd() : o.body || "",
  };
  fs.writeFileSync(path.join(dir, `${id}.md`), serializeTask(task));
  process.stdout.write(`${id}\n`);
}

function cmdRead(id, o) {
  const found = findTaskFile(id, o);
  if (!found) fail(`Task '${id}' not found`);
  const task = parseTask(fs.readFileSync(found.file, "utf8"));
  const ctx = buildContext(loadTasks(found.home));
  const effective = ctx.effStatus(id);
  if (o.format === "json") {
    process.stdout.write(`${JSON.stringify({ ...task, effectiveStatus: effective, leaf: ctx.isLeaf(id), ready: ctx.ready(id) }, null, 2)}\n`);
    return;
  }
  const out = [`id = "${task.id}"`, `name = "${task.name}"`, `type = "${task.type || "task"}"`, `status = "${task.status || "pending"}"`, `effective_status = "${effective}"`];
  if (task.parentId) out.push(`parentId = "${task.parentId}"`);
  if (task.override) out.push("override = true");
  out.push(`dependsOn = [${(task.dependsOn || []).map((v) => `"${v}"`).join(", ")}]`);
  if (task.body) out.push("", `body = """\n${task.body}\n"""`);
  process.stdout.write(`${out.join("\n")}\n`);
}

function cmdUpdate(id, o) {
  const found = findTaskFile(id, o);
  if (!found) fail(`Task '${id}' not found`);
  const task = parseTask(fs.readFileSync(found.file, "utf8"));
  const siblings = loadTasks(found.home);
  const known = new Set(siblings.map((t) => t.id));
  if (o.name != null) task.name = o.name;
  if (o.type != null) task.type = o.type;
  if (o.parent != null) {
    if (o.parent && !known.has(o.parent)) fail(`Parent task '${o.parent}' not found`);
    if (o.parent && wouldCycle(id, o.parent, siblings)) fail(`Re-parenting '${id}' under '${o.parent}' would create a cycle`);
    task.parentId = o.parent || undefined;
  }
  if (o.status != null) {
    if (!LEAF_STATUSES.includes(o.status)) fail(`Invalid status '${o.status}'`);
    task.status = o.status;
  }
  if (o["depends-on"] != null) {
    const deps = splitList(o["depends-on"]);
    for (const d of deps) if (!known.has(d)) fail(`dependsOn target '${d}' not found`);
    task.dependsOn = deps;
  }
  if (o["body-file"] != null) task.body = fs.readFileSync(o["body-file"], "utf8").trimEnd();
  else if (o.body != null) task.body = o.body;
  if (o.override) task.override = true;
  if (o.rollup) delete task.override;
  if (o["append-note"]) {
    const note = `- ${nowIso()} — ${o["append-note"]}`;
    task.body = task.body ? `${task.body}\n${note}` : note;
  }
  fs.writeFileSync(found.file, serializeTask(task));
  process.stdout.write(`${id}\n`);
}

// True if re-parenting `id` under `newParent` would form a cycle (newParent is id or a descendant).
function wouldCycle(id, newParent, tasks) {
  const byId = new Map(tasks.map((t) => [t.id, t]));
  let cur = newParent;
  const seen = new Set();
  while (cur && !seen.has(cur)) {
    if (cur === id) return true;
    seen.add(cur);
    cur = byId.get(cur)?.parentId;
  }
  return false;
}

function cmdDelete(id, o) {
  const found = findTaskFile(id, o);
  if (!found) fail(`Task '${id}' not found`);
  const tasks = loadTasks(found.home);
  const kids = tasks.filter((t) => t.parentId === id).map((t) => t.id);
  if (kids.length && !o.cascade) fail(`Task '${id}' has ${kids.length} child task(s); pass --cascade to remove the subtree`);
  const dir = trackerDir(found.home);
  const toRemove = [id];
  for (let i = 0; i < toRemove.length; i++) {
    for (const t of tasks) if (t.parentId === toRemove[i]) toRemove.push(t.id);
  }
  for (const rid of toRemove) fs.rmSync(path.join(dir, `${rid}.md`), { force: true });
  process.stdout.write(`deleted ${toRemove.length} task(s)\n`);
}

function cmdList(o) {
  const home = homeOf(o);
  if (!home) fail("list requires --home <dir> or --spec <id>");
  const ctx = buildContext(loadTasks(home));
  let all = [...ctx.byId.values()].map((t) => ctx.node(t));
  if (o.status != null) all = all.filter((n) => n.status === o.status);
  if (o.type != null) all = all.filter((n) => n.type === o.type);
  if (o.parent != null) all = all.filter((n) => n.parentId === o.parent);

  let shown;
  if (o.level != null) shown = all.filter((n) => n.depth === Number(o.level));
  else if (o.depth != null) shown = all.filter((n) => n.depth <= Number(o.depth));
  else if (o.status != null || o.type != null || o.parent != null) shown = all; // explicit filter → all depths
  else shown = all.filter((n) => n.depth === 0);
  shown.sort((a, b) => a.depth - b.depth || a.id.localeCompare(b.id));

  if (o.format === "tree") {
    for (const n of shown) process.stdout.write(`${"  ".repeat(n.depth)}- [${n.status}] ${n.name} (${n.id})${n.ready ? " «ready»" : ""}\n`);
    return;
  }
  if (o.format === "toml") {
    const out = [`home = "${path.basename(home)}"`, "", "[counts]"];
    for (const [k, v] of Object.entries(countBy(shown))) out.push(`${k} = ${v}`);
    process.stdout.write(`${out.join("\n")}\n`);
    return;
  }
  process.stdout.write(`${JSON.stringify({ home: path.basename(home), counts: countBy(shown), tasks: shown }, null, 2)}\n`);
}

function cmdGraph(o) {
  const home = homeOf(o);
  if (!home) fail("graph requires --home <dir> or --spec <id>");
  const ctx = buildContext(loadTasks(home));
  const nodes = [...ctx.byId.values()].map((t) => ctx.node(t));

  if (o.format === "dot") {
    const lines = ["digraph tasks {", "  node [shape=box];"];
    for (const n of nodes) lines.push(`  "${n.id}" [label="${n.name}\\n(${n.status})"];`);
    for (const n of nodes) if (n.parentId) lines.push(`  "${n.parentId}" -> "${n.id}" [style=solid, arrowhead=none];`);
    for (const n of nodes) for (const d of n.dependsOn) lines.push(`  "${d}" -> "${n.id}" [style=dashed, color=gray];`);
    lines.push("}");
    process.stdout.write(`${lines.join("\n")}\n`);
    return;
  }
  if (o.format === "toml") {
    const out = [`home = "${path.basename(home)}"`, "", "[counts]"];
    for (const [k, v] of Object.entries(countBy(nodes))) out.push(`${k} = ${v}`);
    out.push("", `ready = [${nodes.filter((n) => n.ready).map((n) => `"${n.id}"`).join(", ")}]`);
    process.stdout.write(`${out.join("\n")}\n`);
    return;
  }
  process.stdout.write(`${JSON.stringify({ home: path.basename(home), counts: countBy(nodes), ready: nodes.filter((n) => n.ready).map((n) => n.id), tasks: nodes }, null, 2)}\n`);
}

// ── Entry ─────────────────────────────────────────────────────────────────────────
function main() {
  const [, , cmd, ...rest] = process.argv;
  let parsed;
  try {
    parsed = parseArgs({
    args: rest,
    allowPositionals: true,
    options: {
      spec: { type: "string" },
      home: { type: "string" },
      name: { type: "string" },
      type: { type: "string" },
      status: { type: "string" },
      parent: { type: "string" },
      "depends-on": { type: "string" },
      body: { type: "string" },
      "body-file": { type: "string" },
      "append-note": { type: "string" },
      format: { type: "string" },
      depth: { type: "string" },
      level: { type: "string" },
      override: { type: "boolean" },
      rollup: { type: "boolean" },
      cascade: { type: "boolean" },
    },
    });
  } catch (e) {
    fail(e.message);
  }
  const { values, positionals } = parsed;
  const id = positionals[0];
  switch (cmd) {
    case "create": return cmdCreate(values);
    case "read": return id ? cmdRead(id, values) : fail("read requires a task id");
    case "update": return id ? cmdUpdate(id, values) : fail("update requires a task id");
    case "delete": return id ? cmdDelete(id, values) : fail("delete requires a task id");
    case "list": return cmdList(values);
    case "graph": return cmdGraph(values);
    default: return fail(`Unknown command '${cmd || ""}' (create|read|update|delete|list|graph)`);
  }
}

main();
