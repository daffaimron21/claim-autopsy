---
inclusion: always
---
# Conventions

- Backend sends message KEYS (`finding.*`, `fix.*`), never display text. The web app translates (en default, id).
- Every finding points to a RuleID (`R-*` rule engine, `M-*` resume reader). Model findings must quote the resume text verbatim.
- Rule engine is deterministic: no model calls in `internal/rules`.
- Synthetic data only. Never add real patient data, BPJS or SATUSEHAT credentials, or API keys.
- Never gitignore `.kiro/`.
