package rules

import (
	"strconv"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

type losRule struct{}

func (losRule) ID() string { return "R-LOS" }

func (losRule) Check(c claim.Case) []claim.Finding {
	g, known := data.CBG.Groups[c.CBGCode]
	days, ok := losDays(c.AdmitDate, c.DischargeDate)
	if !known || !ok || (days >= g.LosMin && days <= g.LosMax) {
		return nil
	}
	return []claim.Finding{{
		RuleID: "R-LOS", Source: claim.SourceRule, Severity: claim.SeverityMajor,
		MessageKey: "finding.los.out_of_range", FixKey: "fix.review_los_or_regroup", Fixable: true,
		Params: map[string]string{
			"los": strconv.Itoa(days), "min": strconv.Itoa(g.LosMin), "max": strconv.Itoa(g.LosMax), "group": g.Label,
		},
		Evidence: []claim.Evidence{{Field: "admit_date", Quote: c.AdmitDate}, {Field: "discharge_date", Quote: c.DischargeDate}},
	}}
}
