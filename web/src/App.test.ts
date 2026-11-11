import { afterEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import App from './App.vue'
import { makeI18n } from './i18n'

const reply = (body: unknown, ok = true) => ({ ok, json: async () => body }) as Response
const detail = (id: string) => ({
  id, title: `Title ${id}`, cbg_code: 'J-4-10-I', resume_text: 'x', return_reason: 'r',
  diagnoses: [{ code: 'J18.9', primary: true }], procedures: [], admit_date: '2026-01-01', discharge_date: '2026-01-02', care_class: 1,
})
const report = (id: string) => ({ case_id: id, primary: null, secondary: [], verdict: 'resubmit_as_is', degraded: false, saved: true })

function deferred<T>() {
  let resolve!: (v: T) => void
  const promise = new Promise<T>((r) => { resolve = r })
  return { promise, resolve }
}

afterEach(() => vi.unstubAllGlobals())

async function setup(routes: Record<string, () => Promise<Response>>) {
  vi.stubGlobal('fetch', vi.fn((url: string) => {
    if (url === '/api/cases') return Promise.resolve(reply([{ id: 'a', title: 'A' }, { id: 'b', title: 'B' }]))
    const r = routes[url]
    if (!r) throw new Error(`unexpected ${url}`)
    return r()
  }))
  const w = mount(App, { global: { plugins: [makeI18n('en')] } })
  await flushPromises() // loads the list and auto-selects the first case
  return w
}

describe('App', () => {
  it('discards a report that arrives after another case was selected', async () => {
    const slow = deferred<Response>()
    const w = await setup({
      '/api/cases/a': () => Promise.resolve(reply(detail('a'))),
      '/api/cases/b': () => Promise.resolve(reply(detail('b'))),
      '/api/cases/a/autopsy': () => slow.promise,
    })
    await w.find('button.primary').trigger('click') // autopsy for a, still pending
    await w.findAll('.case-list li button')[1].trigger('click') // user moves to b
    await flushPromises()
    slow.resolve(reply(report('a')))
    await flushPromises()
    expect(w.find('.report').exists()).toBe(false)
    expect(w.find('.case-file').text()).toContain('Title b')
  })

  it('does not keep showing the old case file when the new case fails to load', async () => {
    const w = await setup({
      '/api/cases/a': () => Promise.resolve(reply(detail('a'))),
      '/api/cases/b': () => Promise.resolve(reply({ error: 'error.case_not_found' }, false)),
    })
    await w.findAll('.case-list li button')[1].trigger('click')
    await flushPromises()
    expect(w.find('.case-file').exists()).toBe(false)
    expect(w.find('[role="alert"]').text()).toContain('That claim was not found')
  })
})
