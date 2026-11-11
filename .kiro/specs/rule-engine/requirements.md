# Requirements: rule engine

1. WHEN a case is submitted THEN the engine SHALL evaluate every rule and return findings with RuleID, severity, message key, evidence and fix key.
2. WHEN discharge date is before admit date THEN the engine SHALL return a blocking, fixable `R-DATES` finding.
3. WHEN admit date equals discharge date THEN the engine SHALL treat the stay as 1 day and return no `R-DATES` finding.
4. WHEN length of stay is outside the CBG range THEN the engine SHALL return a major `R-LOS` finding with `los` and `max` params.
5. WHEN a case is clean THEN the engine SHALL return no findings.
6. The rule data SHALL carry a `source` and `simplified: true` per entry.
