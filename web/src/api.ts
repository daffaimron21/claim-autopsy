export type Source = 'rule' | 'resume'
export type Verdict = 'resubmit_as_is' | 'resubmit_after_fix' | 'not_resubmittable'

export interface Evidence { field: string; quote: string }
export interface Finding {
  rule_id: string
  source: Source
  severity: 0 | 1 | 2
  message_key: string
  params?: Record<string, string>
  evidence?: Evidence[]
  fix_key: string
  fixable: boolean
}
export interface Report {
  case_id: string
  primary: Finding | null
  secondary: Finding[]
  verdict: Verdict
  degraded: boolean
}
export interface CaseItem { id: string; title: string }
export interface CaseDetail extends CaseItem {
  cbg_code: string
  resume_text: string
  return_reason: string
  diagnoses: { code: string; primary: boolean }[]
  procedures: { code: string }[]
  admit_date: string
  discharge_date: string
  care_class: number
}

/** Carries a message KEY (error.*) that the UI translates. */
export class ApiError extends Error {}

async function call<T>(path: string, init?: RequestInit): Promise<T> {
  let res: Response
  try {
    res = await fetch(path, init)
  } catch {
    throw new ApiError('error.network')
  }
  const body = await res.json().catch(() => null)
  if (!res.ok) throw new ApiError(body?.error ?? 'error.network')
  return body as T
}

export const listCases = () => call<CaseItem[]>('/api/cases')
export const getCase = (id: string) => call<CaseDetail>(`/api/cases/${encodeURIComponent(id)}`)
export const runAutopsy = (id: string) =>
  call<Report>(`/api/cases/${encodeURIComponent(id)}/autopsy`, { method: 'POST' })
