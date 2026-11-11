# How Kiro was used (honest log)

This file feeds the Devpost write-up. Record only what actually happened. Judges read this.

## 2026-11-11 (state at the end of the first build session)

- The design, the implementation plan, and ALL application code in this repo (Go backend, Vue web app, tests, Dockerfile) were written with **Claude Code**, not with Kiro.
- The files in `.kiro/` (two steering files, the `rule-engine` spec, one hook) were authored by Claude Code from its memory of Kiro's file formats.
  They have **not** been opened or run in Kiro. The hook format in particular is unverified.
- Kiroween asks for projects built with Kiro and for a write-up of how Kiro was used. As things stand, that claim cannot honestly be made.

## What the owner needs to do (TODO)

1. Open this repo in Kiro. Check that the steering files load, the spec shows in the Specs panel, and the hook fires. Fix whatever does not, and write down what you saw.
2. Do real work in Kiro on remaining tasks, for example: a new rule written spec-first (requirements, design, tasks), a steering file that changes the output, a hook that blocks a bad change.
3. Replace this file with a dated log of what you actually did and what Kiro did, including what did not work.
