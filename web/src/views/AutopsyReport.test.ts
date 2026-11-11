import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import AutopsyReport from './AutopsyReport.vue'
import { makeI18n } from '../i18n'
import type { Report } from '../api'

const report: Report = {
  case_id: 'case-06',
  verdict: 'resubmit_after_fix',
  degraded: false,
  saved: false,
  primary: {
    rule_id: 'R-LOS', source: 'rule', severity: 1, message_key: 'finding.los.out_of_range',
    params: { los: '14', min: '2', max: '8', group: 'Pneumonia' },
    evidence: [{ field: 'admit_date', quote: '2026-02-03' }],
    fix_key: 'fix.review_los_or_regroup', fixable: true,
  },
  secondary: [{
    rule_id: 'M-DATE-MISMATCH', source: 'resume', severity: 2, message_key: 'finding.resume.date_mismatch',
    evidence: [{ field: 'resume_text', quote: 'pulang 4 Juli 2026' }],
    fix_key: 'fix.correct_dates', fixable: true,
  }],
}

function render(r: Report, locale: 'en' | 'id' = 'en') {
  return mount(AutopsyReport, { props: { report: r }, global: { plugins: [makeI18n(locale)] } })
}

describe('AutopsyReport', () => {
  it('renders primary, secondary, verdict and rule-vs-model badges', () => {
    const w = render(report)
    expect(w.text()).toContain('The stay lasted 14 days')
    expect(w.text()).toContain('Resubmit after fixing')
    expect(w.find('[data-test="primary"] [data-test="source-rule"]').text()).toBe('Rule')
    expect(w.findAll('[data-test="secondary"] [data-test="finding"]')).toHaveLength(1)
    expect(w.find('[data-test="secondary"] [data-test="source-resume"]').text()).toBe('Model')
    expect(w.text()).toContain('pulang 4 Juli 2026')
  })

  it('shows the degraded notice only when degraded', () => {
    expect(render(report).find('[data-test="degraded"]').exists()).toBe(false)
    expect(render({ ...report, degraded: true }).find('[data-test="degraded"]').exists()).toBe(true)
  })

  it('shows the saved notice for saved-only results, but not twice when degraded', () => {
    expect(render({ ...report, saved: true }).find('[data-test="saved"]').text()).toContain('saved results')
    expect(render(report).find('[data-test="saved"]').exists()).toBe(false)
    const w = render({ ...report, saved: true, degraded: true })
    expect(w.find('[data-test="degraded"]').exists()).toBe(true)
    expect(w.find('[data-test="saved"]').exists()).toBe(false)
  })

  it('handles a clean claim with no primary cause', () => {
    const w = render({ case_id: 'x', primary: null, secondary: [], verdict: 'resubmit_as_is', degraded: false, saved: false })
    expect(w.find('[data-test="no-cause"]').text()).toContain('No cause found')
    expect(w.text()).toContain('Can be resubmitted as is')
  })

  it('renders Indonesian when the locale is id', () => {
    expect(render(report, 'id').text()).toContain('Lama rawat 14 hari')
  })

  it('opens the rule drawer with the rule id and the simplified label', async () => {
    const w = render(report)
    await w.find('[data-test="primary"] [data-test="why"]').trigger('click')
    const drawer = w.find('[data-test="drawer"]')
    expect(drawer.text()).toContain('R-LOS')
    expect(drawer.text()).toContain('Simplified for demo')
  })
})
