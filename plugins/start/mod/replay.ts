export type GatePart = 'requirements' | 'solution' | 'plan'
export type GateValue = 'pending' | 'draft' | 'approved' | 'skipped'
export type StageName = 'Define' | 'Build' | 'Verify' | 'Ship'
export type StageMark = 'done' | 'current' | 'todo'

export type Task = {
  id: string
  name: string
  type: string
  status: string
  parent: string | null
  override: boolean
  covers?: string[]
  priority?: string
  surface?: string
}

export type SliceResult = { result: string; business?: number; technical?: number; iterations?: number }
export type Decision = { id: string; title?: string; confirmed: boolean }
export type Question = { id: string; text?: string; resolved: boolean }
export type Assumption = { id: string; text?: string; observable?: boolean; status: string }
export type Run = { state: 'started' | 'ended'; at: string; threshold?: number; maxIterations?: number }
export type Verification = { constitution: 'pass' | 'fail' | 'none'; at: string }

export type SpecState = {
  name: string
  gates: Record<GatePart, GateValue>
  requirements: Map<string, string[]>
  decisions: Map<string, Decision>
  questions: Map<string, Question>
  assumptions: Map<string, Assumption>
  tasks: Map<string, Task>
  results: Map<string, SliceResult>
  run?: Run
  verify?: Verification
  createdAt?: string
  archivedAt?: string
}

export const GATE_PARTS: readonly GatePart[] = ['requirements', 'solution', 'plan']
export const STAGES: readonly StageName[] = ['Define', 'Build', 'Verify', 'Ship']
const SETTLED_GATES: readonly GateValue[] = ['approved', 'skipped']

function emptyState(name: string): SpecState {
  return {
    name,
    gates: { requirements: 'pending', solution: 'pending', plan: 'pending' },
    requirements: new Map(),
    decisions: new Map(),
    questions: new Map(),
    assumptions: new Map(),
    tasks: new Map(),
    results: new Map(),
  }
}

function parseLine(line: string): Record<string, any> | undefined {
  try {
    const event = JSON.parse(line)
    return event && typeof event === 'object' && typeof event.kind === 'string' ? event : undefined
  } catch {
    return undefined
  }
}

function fieldsOf(event: Record<string, any>): Record<string, any> {
  const { kind, ts, from, ...fields } = event
  return fields
}

function merge<T>(map: Map<string, T>, id: string, fields: Record<string, any>): void {
  map.set(id, { ...(map.get(id) ?? {}), ...fields } as T)
}

function apply(state: SpecState, event: Record<string, any>): void {
  const fields = fieldsOf(event)
  switch (event.kind) {
    case 'created':
      state.createdAt ??= event.ts
      if (typeof event.name === 'string') state.name = event.name
      return
    case 'gate':
      if (GATE_PARTS.includes(event.part)) state.gates[event.part as GatePart] = event.to
      return
    case 'requirement':
      state.requirements.set(event.id, Array.isArray(event.acs) ? event.acs : [])
      return
    case 'decision':
      return merge(state.decisions, event.id, fields)
    case 'question':
      return merge(state.questions, event.id, fields)
    case 'assumption':
      return merge(state.assumptions, event.id, fields)
    case 'task':
      return merge(state.tasks, event.id, fields)
    case 'removed':
      for (const id of event.ids ?? []) state.tasks.delete(id)
      return
    case 'result':
      return merge(state.results, event.id, fields)
    case 'run':
      state.run = { ...fields, at: event.ts } as Run
      return
    case 'verify':
      state.verify = { constitution: event.constitution, at: event.ts }
      return
    case 'archived':
      state.archivedAt = event.ts
      return
  }
}

export function replay(text: string, name = ''): SpecState {
  const state = emptyState(name)
  for (const line of text.split('\n')) {
    const event = line.trim() === '' ? undefined : parseLine(line)
    if (event !== undefined) apply(state, event)
  }
  return state
}

export function childrenOf(state: SpecState, id: string): Task[] {
  return [...state.tasks.values()].filter((task) => task.parent === id)
}

export function effectiveStatus(state: SpecState, id: string, seen: Set<string> = new Set()): string {
  const task = state.tasks.get(id)
  if (task === undefined || seen.has(id)) return 'pending'
  seen.add(id)
  const children = childrenOf(state, id)
  if (children.length === 0 || task.override) return task.status
  const statuses = children.map((child) => effectiveStatus(state, child.id, seen))
  if (statuses.some((status) => status === 'failed')) return 'failed'
  if (statuses.every((status) => status === 'passed')) return 'passed'
  if (statuses.every((status) => status === 'pending')) return 'pending'
  return 'in_progress'
}

export function tasksOfType(state: SpecState, type: string): Task[] {
  return [...state.tasks.values()].filter((task) => task.type === type)
}

export function isOpen(task: Task): boolean {
  return task.status !== 'passed'
}

export function isDefined(state: SpecState): boolean {
  return GATE_PARTS.every((part) => SETTLED_GATES.includes(state.gates[part]))
}

export function isBuilt(state: SpecState): boolean {
  const slices = tasksOfType(state, 'slice')
  return slices.length > 0 && slices.every((slice) => effectiveStatus(state, slice.id) === 'passed')
}

export function isVerified(state: SpecState): boolean {
  return state.verify !== undefined && state.verify.constitution !== 'fail' && !tasksOfType(state, 'drift').some(isOpen)
}

export function stageMarks(state: SpecState): Record<StageName, StageMark> {
  const done = [isDefined(state), isBuilt(state), isVerified(state), state.archivedAt !== undefined]
  const current = state.archivedAt !== undefined ? STAGES.length : done.findIndex((isDone) => !isDone)
  const at = current < 0 ? STAGES.length : current
  return Object.fromEntries(STAGES.map((stage, index) => [stage, index < at ? 'done' : index === at ? 'current' : 'todo'])) as Record<StageName, StageMark>
}

export function currentStage(state: SpecState): StageName {
  const marks = stageMarks(state)
  return STAGES.find((stage) => marks[stage] === 'current') ?? 'Ship'
}

export function isRunLive(state: SpecState): boolean {
  return state.run?.state === 'started'
}

export function uncoveredCriteria(state: SpecState): string[] {
  const slices = tasksOfType(state, 'slice')
  if (!slices.some((slice) => slice.covers !== undefined)) return []
  const covered = new Set(slices.flatMap((slice) => slice.covers ?? []))
  return [...state.requirements.values()].flat().filter((ac) => !covered.has(ac))
}

export function slicesWithoutP0(state: SpecState): Task[] {
  const scenarios = tasksOfType(state, 'scenario')
  if (!scenarios.some((scenario) => scenario.priority !== undefined)) return []
  return tasksOfType(state, 'slice').filter((slice) => !scenarios.some((scenario) => scenario.parent === slice.id && scenario.priority === 'P0'))
}
