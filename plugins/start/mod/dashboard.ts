import { replay, type SpecState } from './replay.ts'
import { footerText, TONE_STYLE, view, type Live } from './render.ts'

export const DASHBOARD_PANE = 'start-dashboard'
export const DASHBOARD_COMMAND = 'start-dashboard'
const PANE_COLUMNS = 48
const EVENTS_FILE = 'events.jsonl'
const ARCHIVED_EVENTS_SUFFIX = '.events.jsonl'
const SKILL_PREFIX = 'start:'
const LIFECYCLE_SKILLS = new Set(['specify', 'specify-meta', 'specify-requirements', 'specify-solution', 'specify-plan', 'implement', 'validate', 'archive'])
const SPEC_PATH_MARKERS = ['.start/', 'task.js', 'spec.js']
const NO_SPEC = 'No spec has recorded events yet. A spec records them in .start/specs/<id>/events.jsonl; run spec.js <id> --snapshot to start one for an existing spec.'

type EventLog = { path: string; name: string; modifiedAt: number }

let spec: SpecState | undefined
let pinnedSpec: string | undefined
let isDismissed = false
let isAsking = false
let runningAgents = 0
const collapsed = new Set<string>(['ledger'])

export function register(on) {
  on('session.start', async ($, e, next) => {
    await $.command.register({ name: DASHBOARD_COMMAND, description: 'Open or close the start dashboard; pass a spec id to show that spec', argumentHint: '[spec-id]' }).catch(() => undefined)
    await refresh($).catch(() => undefined)
    return next(e)
  })

  on('command.run', async ($, e, next) => {
    if (e.command !== DASHBOARD_COMMAND) return next(e)
    return { text: await runCommand($, e.args.trim()) }
  })

  on('skill.prompt', async ($, e, next) => {
    const prompted = await next(e)
    if (isLifecycleSkill(e.skill)) await refresh($).then(() => openPane($)).catch(() => undefined)
    return prompted
  })

  on('tool.call', async ($, e, next) => {
    const isQuestion = e.tool === 'AskUserQuestion'
    if (isQuestion) {
      isAsking = true
      $.ui.invalidate('ui.render')
    }
    try {
      return await next(e)
    } finally {
      if (isQuestion) isAsking = false
      if (isQuestion || touchesSpecs(e)) await refresh($).catch(() => undefined)
    }
  })

  on('agent.spawn', async ($, e, next) => {
    const spawned = await next(e)
    await refresh($).catch(() => undefined)
    return spawned
  })

  on('turn.complete', async ($, e, next) => {
    const completed = await next(e)
    await refresh($).catch(() => undefined)
    return completed
  })

  on('ui.close', ($, e, next) => {
    if (e.id === DASHBOARD_PANE && e.origin?.kind === 'person') isDismissed = true
    return next(e)
  })

  on('ui.render', { component: 'Pane' }, ($, e, next) => {
    if (e.requestId !== DASHBOARD_PANE) return next(e)
    const { Box, Text, Button } = $.ui.resolve(e)
    if (spec === undefined) return h(Text, { dimColor: true }, NO_SPEC)
    const shown = view(spec, live(), e.props.bodyColumns)
    const sections = shown.sections.map((section) => {
      const isCollapsed = collapsed.has(section.key)
      const toggle = () => {
        if (!collapsed.delete(section.key)) collapsed.add(section.key)
        $.ui.invalidate('ui.render')
      }
      const summary = isCollapsed && section.lines.length === 1 ? `  ${section.lines[0].text}` : ''
      return h(
        Box,
        { key: `section:${section.key}`, flexDirection: 'column', marginTop: 1 },
        h(
          Box,
          { flexDirection: 'row', gap: 1 },
          h(Button, { key: `toggle:${section.key}`, label: isCollapsed ? '▸' : '▾', plain: true, onPress: toggle }),
          h(Text, { bold: true, wrap: 'truncate-end' }, `${section.title}${summary}`),
        ),
        ...(isCollapsed
          ? []
          : section.lines.map((row, index) => h(Box, { key: `${section.key}:${index}`, paddingLeft: 2 }, h(Text, { ...TONE_STYLE[row.tone], wrap: 'truncate-end' }, row.text)))),
      )
    })
    return h(
      Box,
      { flexDirection: 'column' },
      h(Box, { key: 'title' }, h(Text, { bold: true, wrap: 'truncate-end' }, shown.title)),
      h(Box, { key: 'stepper' }, h(Text, { wrap: 'truncate-end' }, shown.stepper)),
      h(Box, { key: 'stage', paddingLeft: 2 }, h(Text, { ...TONE_STYLE[shown.stageLine.tone], wrap: 'truncate-end' }, shown.stageLine.text)),
      ...sections,
    )
  })

  on('ui.render', { component: 'SessionMode' }, async ($, e, next) => {
    const drawn = await next(e)
    if (spec === undefined) return drawn
    const { Box, Text } = $.ui.resolve(e)
    return h(Box, { flexDirection: 'row', gap: 2 }, drawn, h(Box, { key: 'start-progress' }, h(Text, { dimColor: true }, footerText(spec, live()))))
  })
}

function live(): Live {
  return { isAsking, runningAgents, now: Date.now() }
}

function isLifecycleSkill(skill: string): boolean {
  const isOwn = skill.startsWith(SKILL_PREFIX)
  return (isOwn || !skill.includes(':')) && LIFECYCLE_SKILLS.has(isOwn ? skill.slice(SKILL_PREFIX.length) : skill)
}

function touchesSpecs(e): boolean {
  const call = JSON.stringify(e)
  return SPEC_PATH_MARKERS.some((marker) => call.includes(marker))
}

async function runCommand($, specId: string): Promise<string> {
  if (specId !== '') {
    pinnedSpec = specId
    isDismissed = false
    await refresh($)
    if (spec === undefined) return `No event log found for spec ${specId}`
    await openPane($)
    return `Dashboard showing ${spec.name}`
  }
  if ((await $.ui.panes()).some((pane) => pane.id === DASHBOARD_PANE)) {
    isDismissed = true
    await $.ui.close({ id: DASHBOARD_PANE })
    return 'Dashboard closed'
  }
  isDismissed = false
  await refresh($)
  await openPane($)
  return spec === undefined ? NO_SPEC : `Dashboard showing ${spec.name}`
}

async function openPane($): Promise<void> {
  if (isDismissed) return
  await $.ui.open({ id: DASHBOARD_PANE, title: 'start', columns: PANE_COLUMNS })
}

async function refresh($): Promise<void> {
  const logs = await eventLogs($)
  const chosen = pinnedSpec === undefined ? newest(logs) : logs.find((log) => log.name.startsWith(pinnedSpec!))
  spec = chosen === undefined ? undefined : replay(await $.fs.read(chosen.path), chosen.name)
  runningAgents = (await $.agent.list().catch(() => [])).filter((agent) => agent.status === 'running').length
  $.ui.invalidate('ui.render')
}

function newest(logs: EventLog[]): EventLog | undefined {
  return logs.reduce<EventLog | undefined>((latest, log) => (latest === undefined || log.modifiedAt > latest.modifiedAt ? log : latest), undefined)
}

async function eventLogs($): Promise<EventLog[]> {
  const root = `${await $.session.cwd()}/.start`
  const specs = await $.fs.list(`${root}/specs`).catch(() => [])
  const archived = await $.fs.list(`${root}/archive`).catch(() => [])
  const candidates = [
    ...specs.filter((entry) => entry.kind === 'dir').map((entry) => ({ path: `${root}/specs/${entry.name}/${EVENTS_FILE}`, name: entry.name })),
    ...archived
      .filter((entry) => entry.kind === 'file' && entry.name.endsWith(ARCHIVED_EVENTS_SUFFIX))
      .map((entry) => ({ path: `${root}/archive/${entry.name}`, name: entry.name.slice(0, -ARCHIVED_EVENTS_SUFFIX.length) })),
  ]
  const found = await Promise.all(
    candidates.map(async (candidate) => {
      const stat = await $.fs.stat(candidate.path).catch(() => undefined)
      return stat === undefined ? undefined : { ...candidate, modifiedAt: stat.mtimeMs }
    }),
  )
  return found.filter((log): log is EventLog => log !== undefined)
}
