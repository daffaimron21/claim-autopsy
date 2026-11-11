import { describe, expect, it } from 'vitest'
import { readFileSync, readdirSync } from 'node:fs'
import { join } from 'node:path'
import en from './en.json'
import id from './id.json'

function flatten(o: Record<string, unknown>, prefix = ''): string[] {
  return Object.entries(o).flatMap(([k, v]) =>
    typeof v === 'object' && v !== null ? flatten(v as Record<string, unknown>, `${prefix}${k}.`) : [`${prefix}${k}`])
}

const enKeys = new Set(flatten(en))
const idKeys = new Set(flatten(id))

function backendKeys(): string[] {
  const root = join(__dirname, '../../..')
  const files = [
    ...['internal/rules', 'internal/resume', 'internal/api'].flatMap((d) =>
      readdirSync(join(root, d)).filter((f) => f.endsWith('.go') && !f.endsWith('_test.go')).map((f) => join(root, d, f))),
    ...readdirSync(join(root, 'testdata/resume_fixtures')).map((f) => join(root, 'testdata/resume_fixtures', f)),
  ]
  const keys = new Set<string>()
  for (const f of files) {
    for (const m of readFileSync(f, 'utf8').matchAll(/"((?:finding|fix|error)\.[a-z_.]+)"/g)) keys.add(m[1])
  }
  return [...keys]
}

describe('i18n', () => {
  it('en and id have identical keys', () => {
    expect([...enKeys].sort()).toEqual([...idKeys].sort())
  })

  it('every message key the backend can send exists in en', () => {
    const keys = backendKeys()
    expect(keys.length).toBeGreaterThan(10)
    expect(keys.filter((k) => !enKeys.has(k))).toEqual([])
  })

  it('has labels for severities, sources, verdicts and every rule id', () => {
    const need = [
      'severity.minor', 'severity.major', 'severity.blocking', 'source.rule', 'source.resume',
      'verdict.resubmit_as_is', 'verdict.resubmit_after_fix', 'verdict.not_resubmittable',
      'rule.R_LOS.title', 'rule.R_DATES.title', 'rule.R_PRIMARY_DX.title', 'rule.R_PROC_SUPPORT.title',
      'rule.R_CODE_FORMAT.title', 'rule.M_UNCODED_DX.title', 'rule.M_DATE_MISMATCH.title', 'rule.M_PROC_MISSING.title',
    ]
    expect(need.filter((k) => !enKeys.has(k))).toEqual([])
  })
})
