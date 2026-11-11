// Package report merges rule findings and resume findings into one autopsy report.
package report

import (
	"sort"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

// Assemble orders findings (severity desc, rule before resume, then RuleID), picks the
// primary cause and decides the resubmit verdict. The model never decides the ranking.
func Assemble(c claim.Case, ruleFindings, resumeFindings []claim.Finding, degraded bool) claim.Report {
	all := make([]claim.Finding, 0, len(ruleFindings)+len(resumeFindings))
	all = append(all, ruleFindings...)
	all = append(all, resumeFindings...)
	sort.SliceStable(all, func(i, j int) bool {
		a, b := all[i], all[j]
		if a.Severity != b.Severity {
			return a.Severity > b.Severity
		}
		if a.Source != b.Source {
			return a.Source == claim.SourceRule
		}
		return a.RuleID < b.RuleID
	})

	r := claim.Report{CaseID: c.ID, Secondary: []claim.Finding{}, Verdict: claim.VerdictResubmitAsIs, Degraded: degraded}
	if len(all) == 0 {
		return r
	}
	r.Primary = &all[0]
	r.Secondary = append(r.Secondary, all[1:]...)
	r.Verdict = claim.VerdictResubmitAfterFix
	for _, f := range all {
		if f.Severity == claim.SeverityBlocking && !f.Fixable {
			r.Verdict = claim.VerdictNotResubmittable
		}
	}
	return r
}
