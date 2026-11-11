package rules

import "github.com/daffaimron21/claim-autopsy/internal/claim"

type procSupportRule struct{}

func (procSupportRule) ID() string { return "R-PROC-SUPPORT" }

func (procSupportRule) Check(c claim.Case) []claim.Finding {
	var out []claim.Finding
	for _, p := range c.Procedures {
		entry, known := data.Proc.Procedures[p.Code]
		if !known {
			continue
		}
		supported := false
		for _, d := range c.Diagnoses {
			if hasPrefix(d.Code, entry.DxPrefixes) {
				supported = true
				break
			}
		}
		if supported {
			continue
		}
		out = append(out, claim.Finding{
			RuleID: "R-PROC-SUPPORT", Source: claim.SourceRule, Severity: claim.SeverityMajor,
			MessageKey: "finding.procedure.unsupported", FixKey: "fix.align_procedure_and_dx", Fixable: true,
			Params:   map[string]string{"procedure": p.Code, "label": entry.Label},
			Evidence: []claim.Evidence{{Field: "procedure", Quote: p.Code}},
		})
	}
	return out
}
