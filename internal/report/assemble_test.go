package report

import (
	"testing"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

func f(id string, src claim.Source, sev claim.Severity, fixable bool) claim.Finding {
	return claim.Finding{RuleID: id, Source: src, Severity: sev, Fixable: fixable}
}

func ids(r claim.Report) []string {
	var out []string
	if r.Primary != nil {
		out = append(out, r.Primary.RuleID)
	}
	for _, s := range r.Secondary {
		out = append(out, s.RuleID)
	}
	return out
}

func TestAssemble_OrdersBySeverityThenSourceThenRuleID(t *testing.T) {
	rf := []claim.Finding{
		f("R-CODE-FORMAT", claim.SourceRule, claim.SeverityMinor, true),
		f("R-PROC-SUPPORT", claim.SourceRule, claim.SeverityMajor, true),
		f("R-LOS", claim.SourceRule, claim.SeverityMajor, true),
	}
	mf := []claim.Finding{
		f("M-UNCODED-DX", claim.SourceResume, claim.SeverityMajor, true),
		f("M-DATE-MISMATCH", claim.SourceResume, claim.SeverityBlocking, true),
	}
	got := ids(Assemble(claim.Case{ID: "x"}, rf, mf, false))
	want := []string{"M-DATE-MISMATCH", "R-LOS", "R-PROC-SUPPORT", "M-UNCODED-DX", "R-CODE-FORMAT"}
	if len(got) != len(want) {
		t.Fatalf("want %v got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("want %v got %v", want, got)
		}
	}
}

func TestAssemble_NoFindings_PrimaryNilAndResubmitAsIs(t *testing.T) {
	r := Assemble(claim.Case{ID: "x"}, nil, nil, false)
	if r.Primary != nil || r.Verdict != claim.VerdictResubmitAsIs {
		t.Fatalf("got %+v", r)
	}
	if r.Secondary == nil || len(r.Secondary) != 0 {
		t.Fatalf("Secondary must be an empty non-nil slice for clean JSON, got %#v", r.Secondary)
	}
}

func TestAssemble_UnfixableBlocking_NotResubmittable(t *testing.T) {
	r := Assemble(claim.Case{}, []claim.Finding{f("R-X", claim.SourceRule, claim.SeverityBlocking, false)}, nil, false)
	if r.Verdict != claim.VerdictNotResubmittable {
		t.Fatalf("got %v", r.Verdict)
	}
}

func TestAssemble_FixableBlocking_ResubmitAfterFix(t *testing.T) {
	r := Assemble(claim.Case{}, []claim.Finding{f("R-X", claim.SourceRule, claim.SeverityBlocking, true)}, nil, false)
	if r.Verdict != claim.VerdictResubmitAfterFix {
		t.Fatalf("got %v", r.Verdict)
	}
}

func TestAssemble_PassesDegradedAndCaseID(t *testing.T) {
	r := Assemble(claim.Case{ID: "case-9"}, nil, nil, true)
	if !r.Degraded || r.CaseID != "case-9" {
		t.Fatalf("got %+v", r)
	}
}

func TestAssemble_MergesResumeFindingsOfTheSameRule(t *testing.T) {
	ev := func(q string) []claim.Evidence { return []claim.Evidence{{Field: "resume_text", Quote: q}} }
	m1 := f("M-DATE-MISMATCH", claim.SourceResume, claim.SeverityBlocking, true)
	m1.Evidence = ev("Dirawat 3 hari")
	m2 := f("M-DATE-MISMATCH", claim.SourceResume, claim.SeverityBlocking, true)
	m2.Evidence = ev("pulang 4 Juli 2026")
	r1, r2 := f("R-CODE-FORMAT", claim.SourceRule, claim.SeverityMinor, true), f("R-CODE-FORMAT", claim.SourceRule, claim.SeverityMinor, true)

	r := Assemble(claim.Case{}, []claim.Finding{r1, r2}, []claim.Finding{m1, m2}, false)

	got := ids(r)
	if len(got) != 3 || got[0] != "M-DATE-MISMATCH" {
		t.Fatalf("want 1 merged model finding + 2 separate rule findings, got %v", got)
	}
	if len(r.Primary.Evidence) != 2 || r.Primary.Evidence[1].Quote != "pulang 4 Juli 2026" {
		t.Fatalf("merged finding must carry both quotes, got %+v", r.Primary.Evidence)
	}
	if len(m1.Evidence) != 1 {
		t.Fatalf("input finding must not be mutated (it may be cached), got %+v", m1.Evidence)
	}
}
