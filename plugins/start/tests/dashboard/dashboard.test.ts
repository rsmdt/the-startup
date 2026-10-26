import { expect, test } from 'claude-code/testing'
import { APPROVED_GATES, events, task } from './events.ts'

const PANE = 'start-dashboard'
const COMMAND = 'start-dashboard'
const CWD = '/repo'
const SPEC_LOG = `${CWD}/.start/specs/001-ledger/events.jsonl`
const PANE_PROPS = { title: 'start', isFocused: false, bodyColumns: 46, placement: 'docked', scroll: { offset: 0, bodyRows: 40 }, view: {} }

type Stubs = { panes: Array<{ id: string }>; commands: string[] }

function stubMods(on, files: Record<string, string>, mtimes: Record<string, number> = {}): Stubs {
  const panes: Array<{ id: string }> = []
  const commands: string[] = []
  const directories = () => new Set(Object.keys(files).flatMap((path) => path.split('/').slice(1, -1).map((_, index, parts) => `/${parts.slice(0, index + 1).join('/')}`)))
  const list = (dir: string) => {
    const prefix = `${dir}/`
    const names = new Map<string, 'file' | 'dir'>()
    for (const path of Object.keys(files)) {
      if (!path.startsWith(prefix)) continue
      const [name, ...rest] = path.slice(prefix.length).split('/')
      names.set(name, rest.length === 0 ? 'file' : 'dir')
    }
    if (names.size === 0 && !directories().has(dir)) throw new Error(`ENOENT: ${dir}`)
    return [...names].map(([name, kind]) => ({ name, kind, size: 0, mtimeMs: 1, isLink: false }))
  }
  const answer = (value: () => unknown) => {
    try {
      return { value: value() }
    } catch (error) {
      return { deny: (error as Error).message }
    }
  }
  on('session.cwd', () => answer(() => CWD))
  on('session.start', () => ({ cwd: CWD }))
  on('fs.list', ($, e) => answer(() => list(e.path)))
  on('fs.stat', ($, e) =>
    answer(() => {
      if (!(e.path in files)) throw new Error(`ENOENT: ${e.path}`)
      return { kind: 'file', size: files[e.path].length, mtimeMs: mtimes[e.path] ?? 1, isLink: false }
    }),
  )
  on('fs.read', ($, e) =>
    answer(() => {
      if (!(e.path in files)) throw new Error(`ENOENT: ${e.path}`)
      return files[e.path]
    }),
  )
  on('agent.list', () => answer(() => []))
  on('command.register', ($, e) => answer(() => (commands.push(e.name), { command: e.name })))
  on('ui.open', ($, e) => answer(() => void (panes.some((pane) => pane.id === e.id) || panes.push({ id: e.id }))))
  on('ui.close', ($, e) => answer(() => void panes.splice(0, panes.length, ...panes.filter((pane) => pane.id !== e.id))))
  on('ui.panes', () => answer(() => [...panes]))
  on('tool.call', () => ({ result: 'ran' }))
  on('skill.prompt', ($, e) => ({ text: e.text }))
  on('ui.render', ($, e) => h($.ui.resolve(e).Box, { key: 'engine-band' }))
  return { panes, commands }
}

function pane($) {
  return $.ui.mount({ plugin: 'start', surface: 'terminal', component: 'Pane', requestId: PANE, props: PANE_PROPS })
}

const BUILDING = events(...APPROVED_GATES, task('p', 'phase', 'pending', null, { name: 'Phase 1: Ledger' }), task('s', 'slice', 'building', 'p', { override: true, name: 'Expenses' }))

test('the command opens the pane on the newest spec and closes it again', async ($, on) => {
  const stubs = stubMods(on, { [SPEC_LOG]: BUILDING })
  await $.session.start({ cwd: CWD })
  expect(stubs.commands).toContain(COMMAND)
  await $.command.run({ command: COMMAND, args: '' })
  expect(stubs.panes.map((open) => open.id)).toEqual([PANE])
  const ui = await pane($)
  expect((await ui.find({ key: 'title' }))?.text).toBe('start · 001-ledger')
  expect((await ui.find({ key: 'stepper' }))?.text).toBe('● Define ── ◐ Build ── ○ Verify ── ○ Ship')
  await $.command.run({ command: COMMAND, args: '' })
  expect(stubs.panes).toEqual([])
})

test('a lifecycle skill opens the pane and other skills do not', async ($, on) => {
  const stubs = stubMods(on, { [SPEC_LOG]: BUILDING })
  await $.session.start({ cwd: CWD })
  await $.skill.prompt({ skill: 'start:review', text: 'review' })
  expect(stubs.panes).toEqual([])
  await $.skill.prompt({ skill: 'start:implement', text: 'implement 001' })
  expect(stubs.panes.map((open) => open.id)).toEqual([PANE])
})

test('closing the pane with the command keeps later skills from reopening it', async ($, on) => {
  const stubs = stubMods(on, { [SPEC_LOG]: BUILDING })
  await $.session.start({ cwd: CWD })
  await $.skill.prompt({ skill: 'start:implement', text: 'implement 001' })
  await $.command.run({ command: COMMAND, args: '' })
  await $.skill.prompt({ skill: 'start:validate', text: 'validate 001' })
  expect(stubs.panes).toEqual([])
})

test('a tracker change made by a tool call shows in the pane', async ($, on) => {
  const files = { [SPEC_LOG]: BUILDING }
  stubMods(on, files)
  await $.session.start({ cwd: CWD })
  await $.command.run({ command: COMMAND, args: '' })
  let ui = await pane($)
  expect((await ui.find({ key: 'needs:0' }))?.text).toBe('⚠ Expenses  building, no run live')
  await ui.unmount()
  files[SPEC_LOG] += events(task('s', 'slice', 'passed', 'p', { override: true, name: 'Expenses', from: 'building' }))
  await $.tool.call({ tool: 'Bash', command: 'node task.js update s --spec 001 --status passed --override', tool_use_id: 'tu-1' })
  ui = await pane($)
  expect((await ui.find({ key: 'stepper' }))?.text).toBe('● Define ── ● Build ── ◐ Verify ── ○ Ship')
})

test('the newest event log wins, and a spec id pins another one', async ($, on) => {
  const older = `${CWD}/.start/specs/002-export/events.jsonl`
  stubMods(on, { [SPEC_LOG]: BUILDING, [older]: events({ kind: 'created', name: '002-export' }) }, { [SPEC_LOG]: 5, [older]: 2 })
  await $.session.start({ cwd: CWD })
  await $.command.run({ command: COMMAND, args: '' })
  const ui = await pane($)
  expect((await ui.find({ key: 'title' }))?.text).toBe('start · 001-ledger')
  await ui.unmount()
  await $.command.run({ command: COMMAND, args: '002' })
  expect((await (await pane($)).find({ key: 'title' }))?.text).toBe('start · 002-export')
})

test('the footer shows the stepper and what needs you', async ($, on) => {
  stubMods(on, { [SPEC_LOG]: BUILDING })
  await $.session.start({ cwd: CWD })
  const ui = await $.ui.mount({ plugin: 'start', surface: 'terminal', component: 'SessionMode', props: { modes: [] } })
  expect((await ui.find({ key: 'start-progress' }))?.text).toBe('001 ●◐○○ · 1 needs you')
})

test('without any event log the pane says how to start one', async ($, on) => {
  stubMods(on, {})
  await $.session.start({ cwd: CWD })
  expect(await $.command.run({ command: COMMAND, args: '' })).toEqual({ text: expect.stringContaining('No spec has recorded events yet') })
})
