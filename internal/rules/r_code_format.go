package rules

import (
	"regexp"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

var (
	icd10Re  = regexp.MustCompile(`^[A-Z][0-9]{2}(\.[0-9A-Z]{1,4})?$`)
	icd9cmRe = regexp.MustCompile(`^[0-9]{2}\.[0-9]{1,2}$`)
)

type codeFormatRule struct{}

func (codeFormatRule) ID() string { return "R-CODE-FORMAT" }

func (codeFormatRule) Check(c claim.Case) []claim.Finding {
	var out []claim.Finding
	bad := func(field, code string) {
		out = append(out, claim.Finding{
			RuleID: "R-CODE-FORMAT", Source: claim.SourceRule, Severity: claim.SeverityMinor,
			MessageKey: "finding.code.format", FixKey: "fix.correct_code", Fixable: true,
			Params:   map[string]string{"code": code},
			Evidence: []claim.Evidence{{Field: field, Quote: code}},
		})
	}
	for _, d := range c.Diagnoses {
		if !icd10Re.MatchString(d.Code) {
			bad("diagnosis", d.Code)
		}
	}
	for _, p := range c.Procedures {
		if !icd9cmRe.MatchString(p.Code) {
			bad("procedure", p.Code)
		}
	}
	return out
}
