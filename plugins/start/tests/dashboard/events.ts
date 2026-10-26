let clock = Date.parse('2026-10-26T10:00:00.000Z')

export function events(...list: Array<Record<string, unknown>>): string {
  return list.map((event) => JSON.stringify({ ts: new Date((clock += 60000)).toISOString(), ...event })).join('\n') + '\n'
}

export function task(id: string, type: string, status: string, parent: string | null = null, extra: Record<string, unknown> = {}) {
  return { kind: 'task', id, name: `${type} ${id}`, type, status, parent, override: false, ...extra }
}

export const APPROVED_GATES = ['requirements', 'solution', 'plan'].map((part) => ({ kind: 'gate', part, to: 'approved' }))
