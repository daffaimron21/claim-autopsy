package rules

import (
	"time"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

type datesRule struct{}

func (datesRule) ID() string { return "R-DATES" }

func (datesRule) Check(c claim.Case) []claim.Finding {
	a, err1 := time.Parse(dateLayout, c.AdmitDate)
	d, err2 := time.Parse(dateLayout, c.DischargeDate)
	key := ""
	switch {
	case err1 != nil || err2 != nil:
		key = "finding.dates.invalid"
	case d.Before(a):
		key = "finding.dates.inverted"
	default:
		return nil
	}
	return []claim.Finding{{
		RuleID: "R-DATES", Source: claim.SourceRule, Severity: claim.SeverityBlocking,
		MessageKey: key, FixKey: "fix.correct_dates", Fixable: true,
		Params:   map[string]string{"admit": c.AdmitDate, "discharge": c.DischargeDate},
		Evidence: []claim.Evidence{{Field: "discharge_date", Quote: c.DischargeDate}},
	}}
}
