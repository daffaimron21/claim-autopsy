# Handover to Kiro

Read this first. Owner speaks Bahasa Indonesia; code, docs and UI strings in this repo are English (UI is bilingual en/id).

## 1. What this is

**Claim Autopsy**: a web app that dissects a returned BPJS claim (synthetic data only) and reports the primary cause,
other causes, evidence, a suggested fix, and whether the claim can be resubmitted.
Entry for **Kiroween 2026**, category **Frankenstein**. Submission deadline: **4 Dec 2026, 14:00 PT**. Target: submit by **3 Dec**.

Kiroween rules that constrain us: new project, public GitHub repo with an open-source license, `.kiro/` at the repo root and
**not gitignored**, demo video under 3 minutes showing the app working, a written explanation of how Kiro was used
(specs, steering, hooks, MCP), a working app URL. Judging: Potential Impact, Technological Implementation, Design, Creativity, Presentation.

## 2. Honest state (do not rewrite history)

- Everything in this repo so far (Go backend, Vue app, tests, Dockerfile, docs, and the files in `.kiro/`) was written by **Claude Code, outside Kiro**.
- The `.kiro/` files (steering, the `rule-engine` spec, one hook) were authored from memory of Kiro's file formats. **They have never been opened in Kiro.** The hook format is the least certain part.
- So far, nothing can truthfully be claimed as "built with Kiro". **Your work in Kiro is what makes that claim true.**
- Keep `docs/kiro-usage.md` as a dated, honest log: what you asked Kiro, what it did, what worked, what did not. Never claim something that did not happen. Judges read it.

## 3. Run and test

```bash
make test            # go test ./... and the web tests (bun)
make dev             # API :8080 + Vite :5173
make check           # .kiro is tracked and not ignored
go run ./cmd/record-fixtures   # re-record saved model results (needs ANTHROPIC_API_KEY in the environment)
```

- `.env` (git-ignored, never commit) holds `ANTHROPIC_API_KEY`. Load it with `set -a; . ./.env; set +a`.
- The app reads `LIVE_MODEL=true`, `ANTHROPIC_API_KEY`, `LIVE_MODEL_NAME` (default `claude-haiku-5-5`), `PORT`, `WEB_DIR`. It ignores `ANTHROPIC_MODEL`.
- Last verified: 51 Go tests (with `-race`), 11 web tests, Docker image builds and serves, gitleaks clean.

## 4. Map

| Path | Job |
|---|---|
| `internal/claim` | Shared types (`Case`, `Finding`, `Report`) and the case loader |
| `internal/rules` | Deterministic rule engine. Data in `ruledata/*.json`, one file per rule |
| `internal/resume` | Reads the free-text discharge summary: Anthropic reader, saved-fixture reader, `Verified` quote check, `Fallback` |
| `internal/report` | `Assemble`: merges, ranks, picks the primary cause, decides the verdict. `eval_test.go` is the oracle over the 6 cases |
| `internal/api` | HTTP handlers, rate limiter, per-case cache of live results |
| `cmd/server`, `cmd/record-fixtures` | Server entry, fixture recorder |
| `testdata/cases`, `testdata/resume_fixtures` | 6 synthetic cases and recorded model findings (embedded into the binary by `data.go`) |
| `web/` | Vue 3 + vue-i18n. `src/i18n/{en,id}.json`, `views/AutopsyReport.vue`, `components/*` |

## 5. Invariants (do not break; they are the point of the project)

1. The model is on a short leash: it may only report `M-UNCODED-DX`, `M-DATE-MISMATCH`, `M-PROC-MISSING`. Severity, message key and fix key come from our table (`internal/resume/anthropic.go`), not from the model.
2. Every model finding must quote the discharge summary verbatim. `Verified` drops any finding whose quote is not in the text. Keep it wrapping both the live and the saved reader.
3. The rule engine is deterministic. No model calls in `internal/rules`. Ranking and the verdict are decided in `internal/report`, never by the model.
4. The backend sends message **keys** (`finding.*`, `fix.*`, `error.*`); the web app translates. `src/i18n/i18n.test.ts` fails if en/id differ or a key the backend can send is missing. Add keys to both languages.
5. Rule data is a **simplified subset designed for this demo**, not an official BPJS rule set. Every data entry carries `source` and `simplified: true`. Never present it as official.
6. Synthetic data only. No real patient data, no credentials, no keys in git. Never gitignore `.kiro/`.
7. Honest labelling: when results come from saved fixtures the report carries `saved: true`, and when the live model failed it also carries `degraded: true`. Do not remove the on-screen notice.
8. Paid-key protection: a successful live result is cached once per case, a limiter guards misses, a failure falls back to saved results.

## 6. Known gaps (deferred from the code review)

- `internal/resume/verify.go`: a one-word quote ("hari") passes the quote check. Add a minimum length.
- `internal/resume/fixture.go`: saved findings are trusted as-is; normalise them through the same closed rule table as live findings.
- The verdict `not_resubmittable` can never occur: every finding is fixable. (See task B below.)
- `web/src/views/AutopsyReport.vue`, `web/src/components/CaseFile.vue`: `v-for` keys can repeat.
- The printed report does not say which claim it is for (print CSS hides the case file). Show case id and title in the report header.
- `internal/rules/r_primary_dx.go`: a claim with no primary diagnosis gives a message with a blank code.
- `web/src/components/RuleDrawer.vue`: no focus management (focus in on open, back on close).

## 7. Work queue for Kiro (in order)

Do the first one before anything else. Paste the quoted prompts into Kiro and adapt them.

**A. Make `.kiro/` real (day 1).**
Open this repo in Kiro. Check: both steering files load, the `rule-engine` spec appears in the Specs panel, the hook fires when a file in `internal/rules/` changes. Fix whatever does not work (the hook JSON is the likely one). Log what you saw in `docs/kiro-usage.md`.

**B. Spec-first: a rule that makes `not_resubmittable` reachable.**
> "Create a spec for a new rule R-FILING-DEADLINE. A claim has a `submitted_date`. If it is submitted more than N days after the discharge date, raise a blocking, NOT fixable finding, so the verdict becomes `not_resubmittable`. N and the rule are designed for this demo and must be labelled 'simplified for demo', not presented as an official BPJS deadline. Add a 7th synthetic case, update the eval oracle, en and id strings, and the rule drawer text."

This touches `claim.Case`, `internal/rules`, `testdata/cases`, `eval_test.go`, both i18n files. Good material for the write-up: requirements, design, tasks, then implementation.

**C. Close the known gaps in section 6** as one spec ("Harden quote check and polish"). Each item gets a failing test first.

**D. Steering and hooks that earn their place.** A steering file that actually changes Kiro's output (for example the invariants above), and a second hook that runs the web tests when an i18n file changes. Verify each one fires, and record it.

**E. Re-record fixtures** (`go run ./cmd/record-fixtures`) after B and C change cases or prompts, review the diff, and re-run `make test`.

**F. Optional, only if it is truly useful:** an MCP server. Do not add one just to fill a checklist.

## 8. Not Kiro's job (owner does these)

- **Deploy** (Docker image; `ANTHROPIC_API_KEY` goes in the host's secret settings, never in a file in git). Set a low monthly spend limit on the Anthropic key.
- **Demo video** (under 3 minutes, real app, no copyrighted music) using `docs/demo-script.md`.
- **Devpost submission** using `docs/devpost.md`. The Kiro section must be written from `docs/kiro-usage.md`.

## 9. Suggested dates

| Date | Goal |
|---|---|
| 12–14 Nov | A, then start B |
| 15–22 Nov | B and C done, E, tests green |
| 23–26 Nov | Deploy, live URL, spend limit |
| 27–30 Nov | Demo video, Devpost text |
| 1–3 Dec | Buffer, final checks, submit (target 3 Dec) |

## 10. Owner preferences

- No `Co-Authored-By` or "Generated with ..." lines in commit messages.
- Tests first for behaviour changes. Run `make test` before committing.
- Never commit `.env`. Run `gitleaks git .` before pushing.
