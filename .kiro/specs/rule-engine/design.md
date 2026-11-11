# Design: rule engine

- `type Rule interface { ID() string; Check(c claim.Case) []claim.Finding }`
- `rules.All()` lists rules; `rules.Run(c)` concatenates all findings.
- Rules: `R-LOS`, `R-DATES`, `R-PRIMARY-DX`, `R-PROC-SUPPORT`, `R-CODE-FORMAT`.
- Data: `internal/rules/ruledata/cbg_map.json` (simplified subset, sourced).
