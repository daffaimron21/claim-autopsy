package rules

import "github.com/daffaimron21/claim-autopsy/internal/claim"

type primaryDxRule struct{}

func (primaryDxRule) ID() string { return "R-PRIMARY-DX" }

func (primaryDxRule) Check(c claim.Case) []claim.Finding {
	g, known := data.CBG.Groups[c.CBGCode]
	f := claim.Finding{
		RuleID: "R-PRIMARY-DX", Source: claim.SourceRule, Severity: claim.SeverityBlocking,
		FixKey: "fix.regroup_cbg", Fixable: true,
		Evidence: []claim.Evidence{{Field: "cbg_code", Quote: c.CBGCode}},
	}
	switch {
	case !known:
		f.MessageKey, f.Params = "finding.cbg.unknown", map[string]string{"cbg": c.CBGCode}
	case !hasPrefix(primaryCode(c), g.DxPrefixes):
		f.MessageKey = "finding.primary_dx.not_in_group"
		f.Params = map[string]string{"code": primaryCode(c), "cbg": c.CBGCode, "group": g.Label}
		f.Evidence = append(f.Evidence, claim.Evidence{Field: "diagnosis", Quote: primaryCode(c)})
	default:
		return nil
	}
	return []claim.Finding{f}
}
