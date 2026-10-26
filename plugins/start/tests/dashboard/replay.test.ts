import { describe, expect, test } from 'claude-code/testing'
import { currentStage, effectiveStatus, replay, stageMarks, tasksOfType } from '../../mod/replay.ts'
import { APPROVED_GATES, events, task } from './events.ts'
import { SPLITWISE_EVENTS } from './splitwise-events.ts'

describe('replaying a snapshot of a finished spec', () => {
  const state = replay(SPLITWISE_EVENTS)

  test('restores the gates, decisions, assumptions and requirements', () => {
    expect(state.name).toBe('001-splitwise-lite')
    expect(state.gates).toEqual({ requirements: 'approved', solution: 'approved', plan: 'approved' })
    expect([...state.decisions.values()].filter((decision) => decision.confirmed).length).toBe(7)
    expect(state.assumptions.size).toBe(30)
    expect(state.requirements.size).toBe(11)
    expect([...state.requirements.values()].flat().length).toBe(61)
  })

  test('restores the tracker tree', () => {
    expect(tasksOfType(state, 'phase').length).toBe(2)
    expect(tasksOfType(state, 'slice').length).toBe(5)
    expect(tasksOfType(state, 'scenario').length).toBe(85)
  })

  test('places a built spec with no verification in the Verify stage', () => {
    expect(currentStage(state)).toBe('Verify')
  })
})

describe('the status rollup matches task.js', () => {
  const tree = (...children: string[]) =>
    replay(events(task('p', 'phase', 'pending'), ...children.map((status, index) => task(`s${index}`, 'slice', status, 'p'))))

  test('a parent with a failed child is failed', () => {
    expect(effectiveStatus(tree('passed', 'failed', 'building'), 'p')).toBe('failed')
  })

  test('a parent whose children all passed has passed', () => {
    expect(effectiveStatus(tree('passed', 'passed'), 'p')).toBe('passed')
  })

  test('a parent whose children are all pending is pending', () => {
    expect(effectiveStatus(tree('pending', 'pending'), 'p')).toBe('pending')
  })

  test('a parent with mixed children is in progress', () => {
    expect(effectiveStatus(tree('passed', 'pending'), 'p')).toBe('in_progress')
  })

  test('an override keeps the stored status over the rollup', () => {
    const state = replay(events(task('s', 'slice', 'passed', null, { override: true }), task('c', 'scenario', 'pending', 's')))
    expect(effectiveStatus(state, 's')).toBe('passed')
  })
})

describe('replaying task changes', () => {
  test('an update keeps fields only the create carried', () => {
    const state = replay(events(task('s', 'slice', 'pending', null, { covers: ['AC-1.1'] }), task('s', 'slice', 'building', null, { from: 'pending' })))
    expect(state.tasks.get('s')).toEqual({ id: 's', name: 'slice s', type: 'slice', status: 'building', parent: null, override: false, covers: ['AC-1.1'] })
  })

  test('a removal drops the tasks', () => {
    const state = replay(events(task('b', 'bug', 'pending'), { kind: 'removed', ids: ['b'] }))
    expect(state.tasks.size).toBe(0)
  })

  test('a torn last line is ignored', () => {
    const state = replay(events(task('b', 'bug', 'pending')) + '{"kind":"ta')
    expect(state.tasks.size).toBe(1)
  })
})

describe('the stage follows the spec through its life', () => {
  const define = events({ kind: 'gate', part: 'requirements', to: 'draft' })
  const build = events(...APPROVED_GATES, task('p', 'phase', 'pending'), task('s', 'slice', 'building', 'p'))
  const verify = events(...APPROVED_GATES, task('s', 'slice', 'passed'))
  const ship = verify + events({ kind: 'verify', constitution: 'pass' })
  const archived = ship + events({ kind: 'archived' })

  test('each stage is current in turn', () => {
    expect([define, build, verify, ship].map((log) => currentStage(replay(log)))).toEqual(['Define', 'Build', 'Verify', 'Ship'])
  })

  test('an archived spec has every stage done', () => {
    expect(Object.values(stageMarks(replay(archived)))).toEqual(['done', 'done', 'done', 'done'])
  })

  test('an open drift finding keeps the spec in Verify', () => {
    const state = replay(ship + events(task('d', 'drift', 'pending')))
    expect(currentStage(state)).toBe('Verify')
  })
})
