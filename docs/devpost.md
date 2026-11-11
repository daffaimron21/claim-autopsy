# Devpost form text (draft)

Fields marked TODO need the owner. Do not submit claims about Kiro that did not happen (see `docs/kiro-usage.md`).

## Project overview

**Project name** (max 60): `Claim Autopsy`

**Elevator pitch** (max 200, 176 used):
`Returned BPJS claims come back with one vague line. Claim Autopsy dissects them: rules find the cause, an LLM reads the discharge summary, and every finding shows its evidence.`

**Category:** Frankenstein

## Project details (draft)

**Inspiration.** When a hospital's BPJS claim is returned as "pending" or "dispute", the reason is one short line. A coder has to guess which part of the claim is wrong. We wanted a post-mortem for dead claims.

**What it does.** Pick a returned claim and run the autopsy. The report shows the primary cause, other causes, the evidence for each, a suggested fix, and whether the claim can be resubmitted. Every finding opens a drawer showing the rule behind it.

**The Frankenstein part.** Three things that do not usually live together are stitched into one body:
a deterministic rule engine over coded data, a language model that reads free-text discharge summaries, and a forensic report format.
The model is kept on a short leash: it can only report mismatches from a closed list, every finding must quote the summary word for word, and a quote that is not in the text is dropped. The model never ranks causes or decides the verdict.

**How we built it.** Go (stdlib HTTP), Vue 3 with vue-i18n (English and Indonesian), Anthropic API for the resume reader, saved results as a fallback so the demo never depends on a provider being up.
TODO (owner): the Kiro section. Describe only what you really did in Kiro: specs, steering, hooks. See `docs/kiro-usage.md`.

**Challenges.** Keeping model output trustworthy (quote verification), and being honest that the rules are a simplified subset for demo, not official BPJS rules.

**What's next.** Real rule coverage with a coder's review, more claim types, and a worklist view for a whole batch of returned claims.

**Built with:** Go, Vue 3, TypeScript, Vite, vue-i18n, Anthropic API, Docker, Kiro (TODO: confirm).

## Links (TODO)

- Repo: https://github.com/daffaimron21/claim-autopsy
- Live app: TODO after deploy
- Demo video (YouTube, public, under 3 minutes): TODO

## Honest limits to state

- All data is synthetic. Rules are simplified for demo, not official.
- Nothing is submitted to BPJS or SATUSEHAT.
