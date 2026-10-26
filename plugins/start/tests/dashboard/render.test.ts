import { describe, expect, test } from 'claude-code/testing'
import { replay } from '../../mod/replay.ts'
import { footerText, needsYou, view, type Live } from '../../mod/render.ts'
import { APPROVED_GATES, events, task } from './events.ts'
import { SPLITWISE_EVENTS } from './splitwise-events.ts'

const QUIET: Live = { isAsking: false, runningAgents: 0, now: Date.parse('2026-10-26T12:00:00.000Z') }

function needs(log: string, live: Live = QUIET): string[] {
  return needsYou(replay(log), live).map((item) => item.text)
}

describe('needs you', () => {
  test('lists a draft gate, an unconfirmed decision, an open question and a proposed observable assumption', () => {
    const log = events(
      { kind: 'gate', part: 'solution', to: 'draft' },
      { kind: 'decision', id: 'ADR-3', title: 'S3 spool', confirmed: false },
      { kind: 'decision', id: 'ADR-1', title: 'Node', confirmed: true },
      { kind: 'question', id: 'Q1', text: 'Retry cap?', resolved: false },
      { kind: 'question', id: 'Q2', text: 'Locale?', resolved: true },
      { kind: 'assumption', id: 'A4', text: 'Export under 30s', observable: true, status: 'proposed' },
      { kind: 'assumption', id: 'A5', text: 'English only', observable: false, status: 'proposed' },
      { kind: 'assumption', id: 'A6', text: 'Newest first', observable: true, status: 'confirmed' },
    )
    expect(needs(log)).toEqual(['◆ solution draft to approve', '? ADR-3 S3 spool — confirm', '? Q1 Retry cap?', '? A4 Export under 30s'])
  })

  test('shows a question Claude is asking right now first', () => {
    expect(needs(events({ kind: 'gate', part: 'plan', to: 'draft' }), { ...QUIET, isAsking: true })[0]).toBe('? Claude is asking you a question')
  })

  test('lists stalled, blocked and failed slices with their iterations', () => {
    const log = events(
      ...APPROVED_GATES,
      { kind: 'run', state: 'ended', maxIterations: 3 },
      task('s1', 'slice', 'failed', null, { override: true }),
      task('s2', 'slice', 'building', null, { override: true }),
      task('s3', 'slice', 'failed', null, { override: true }),
      { kind: 'result', id: 's1', result: 'stalled', iterations: 2 },
      { kind: 'result', id: 's2', result: 'blocked-environment', iterations: 1 },
    )
    expect(needs(log)).toEqual(['⏸ slice s1  stalled  it 2/3', '⛔ slice s2  blocked-environment  it 1/3', '✗ slice s3  failed'])
  })

  test('flags a slice left building when no run is live, but not while one runs', () => {
    const stuck = events(...APPROVED_GATES, task('s', 'slice', 'building', null, { override: true }))
    expect(needs(stuck)).toEqual(['⚠ slice s  building, no run live'])
    expect(needs(stuck + events({ kind: 'run', state: 'started' }))).toEqual([])
  })

  test('flags criteria no slice covers and slices with no P0 scenario', () => {
    const log = events(
      { kind: 'gate', part: 'plan', to: 'approved' },
      { kind: 'requirement', id: 'FR-1', acs: ['AC-1.1', 'AC-1.2'] },
      task('s', 'slice', 'pending', null, { covers: ['AC-1.1'] }),
      task('c', 'scenario', 'pending', 's', { priority: 'P1' }),
    )
    expect(needs(log)).toContain('✗ 1 AC in no slice: AC-1.2')
    expect(needs(log)).toContain('✗ slice s: no P0 scenario')
  })

  test('points to the next command only in the stage it belongs to', () => {
    const reopened = events(...APPROVED_GATES, task('s', 'slice', 'passed'), { kind: 'gate', part: 'plan', to: 'draft' })
    expect(needs(reopened)).toEqual(['◆ plan draft to approve'])
  })

  test('points to the next command once the build is done', () => {
    const built = events(...APPROVED_GATES, task('s', 'slice', 'passed'))
    expect(needs(built)).toEqual(['◆ run /start:validate drift'])
    expect(needs(built + events({ kind: 'verify', constitution: 'pass' }))).toEqual(['◆ run /start:archive'])
    expect(needs(built + events({ kind: 'verify', constitution: 'pass' }, task('d', 'drift', 'pending', null, { name: 'page size 25' })))).toEqual(['⚠ drift: page size 25'])
  })
})

describe('the pane view', () => {
  test('puts needs you first and follows the current stage', () => {
    const shown = view(replay(SPLITWISE_EVENTS), QUIET, 46)
    expect(shown.title).toBe('start · 001-splitwise-lite')
    expect(shown.stepper).toBe('● Define ── ● Build ── ◐ Verify ── ○ Ship')
    expect(shown.sections.map((section) => section.key)).toEqual(['needs', 'verify', 'tracker', 'ledger'])
    expect(shown.sections[0].lines.map((line) => line.text)).toEqual(['◆ run /start:validate drift'])
  })

  test('counts the tracker from the events', () => {
    const tracker = view(replay(SPLITWISE_EVENTS), QUIET, 46).sections.find((section) => section.key === 'tracker')!
    expect(tracker.lines.map((line) => line.text)).toEqual(['phases 2  slices 5  scenarios 85', 'slices     ✓ 5  ◐ 0  ✗ 0  ○ 0'])
  })

  test('shows the coverage while defining', () => {
    const log = SPLITWISE_EVENTS + events({ kind: 'gate', part: 'plan', to: 'draft' })
    const define = view(replay(log), QUIET, 46).sections.find((section) => section.key === 'define')!
    expect(define.lines.map((line) => line.text)).toEqual(['FR 11 · AC 61', 'ADR 7 (7 confirmed · 0 pending)', 'AC in a slice  61/61', 'slices without P0  0/5', 'scenarios  P0 22 · P1 54 · P2 9'])
  })

  test('builds a slice grid for the current phase that fits the pane', () => {
    const log = events(
      ...APPROVED_GATES,
      { kind: 'run', state: 'started', threshold: 0.9, maxIterations: 3 },
      task('p1', 'phase', 'pending', null, { name: 'Phase 1: Ledger' }),
      task('a', 'slice', 'passed', 'p1', { override: true, name: 'Runnable group and member API with a long name' }),
      task('b', 'slice', 'evaluating', 'p1', { override: true, name: 'Expenses' }),
      { kind: 'result', id: 'a', result: 'completed', business: 1, technical: 0.95, iterations: 1 },
    )
    const build = view(replay(log), QUIET, 46).sections.find((section) => section.key === 'build')!
    expect(build.title).toBe('BUILD · P1  threshold 90%')
    expect(build.lines.every((line) => line.text.length <= 46 - 4)).toBe(true)
    expect(build.lines[1].text).toBe('✓ Runnable group and member…  100/ 95  1/3')
    expect(build.lines[2].text).toBe('◐ Expenses' + ' '.repeat(21) + '--/ --  -/3')
  })

  test('the footer is the stepper and the needs-you count', () => {
    expect(footerText(replay(SPLITWISE_EVENTS), QUIET)).toBe('001 ●●◐○ · 1 needs you')
  })
})
