package rules

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

func clean() claim.Case {
	return claim.Case{
		ID: "t", CBGCode: "J-4-10-I",
		Diagnoses:     []claim.Diagnosis{{Code: "J18.9", Primary: true}},
		AdmitDate:     "2026-02-03",
		DischargeDate: "2026-02-08",
		CareClass:     2,
	}
}

func find(fs []claim.Finding, id string) *claim.Finding {
	for i := range fs {
		if fs[i].RuleID == id {
			return &fs[i]
		}
	}
	return nil
}

func TestDates_DischargeBeforeAdmit_Blocking(t *testing.T) {
	c := clean()
	c.AdmitDate, c.DischargeDate = "2026-03-05", "2026-03-01"
	f := find(Run(c), "R-DATES")
	if f == nil || f.Severity != claim.SeverityBlocking || !f.Fixable || f.Source != claim.SourceRule {
		t.Fatalf("want blocking fixable R-DATES rule finding, got %+v", f)
	}
}

func TestDates_SameDayStay_NoFinding(t *testing.T) {
	c := clean()
	c.AdmitDate, c.DischargeDate = "2026-03-05", "2026-03-05"
	if f := find(Run(c), "R-DATES"); f != nil {
		t.Fatalf("same-day stay must not raise R-DATES, got %+v", f)
	}
}

func TestDates_Malformed_Blocking(t *testing.T) {
	c := clean()
	c.DischargeDate = "05/03/2026"
	if f := find(Run(c), "R-DATES"); f == nil || f.Severity != claim.SeverityBlocking {
		t.Fatalf("malformed date must be blocking R-DATES, got %+v", f)
	}
}

func TestLOS_OutsideCBGRange_Major(t *testing.T) {
	c := clean()
	c.DischargeDate = "2026-02-17" // 14 days, J group max is 8
	f := find(Run(c), "R-LOS")
	if f == nil || f.Severity != claim.SeverityMajor {
		t.Fatalf("want major R-LOS, got %+v", f)
	}
	if f.Params["los"] != "14" || f.Params["max"] != "8" {
		t.Fatalf("params want los=14 max=8, got %v", f.Params)
	}
}

func TestLOS_SameDayStayCountsAsOneDay(t *testing.T) {
	c := clean()
	c.AdmitDate, c.DischargeDate = "2026-03-05", "2026-03-05"
	c.CBGCode, c.Diagnoses = "A-4-10-I", []claim.Diagnosis{{Code: "A09", Primary: true}} // min 1
	if f := find(Run(c), "R-LOS"); f != nil {
		t.Fatalf("1 day is inside A group 1..5, got %+v", f)
	}
	c2 := clean() // J group min 2, so 1 day is too short
	c2.AdmitDate, c2.DischargeDate = "2026-03-05", "2026-03-05"
	f := find(Run(c2), "R-LOS")
	if f == nil || f.Params["los"] != "1" {
		t.Fatalf("same-day must count as los=1 and violate min 2, got %+v", f)
	}
}

func TestLOS_SkippedWhenDatesInvalid(t *testing.T) {
	c := clean()
	c.AdmitDate, c.DischargeDate = "2026-03-05", "2026-03-01"
	if f := find(Run(c), "R-LOS"); f != nil {
		t.Fatalf("R-LOS must not fire when dates are invalid, got %+v", f)
	}
}

func TestPrimaryDx_NotInCBGGroup_Blocking(t *testing.T) {
	c := clean()
	c.Diagnoses = []claim.Diagnosis{{Code: "I21.0", Primary: true}}
	f := find(Run(c), "R-PRIMARY-DX")
	if f == nil || f.Severity != claim.SeverityBlocking || !f.Fixable {
		t.Fatalf("want blocking fixable R-PRIMARY-DX, got %+v", f)
	}
}

func TestPrimaryDx_UnknownCBG_Blocking(t *testing.T) {
	c := clean()
	c.CBGCode = "Z-9-99-9"
	if f := find(Run(c), "R-PRIMARY-DX"); f == nil || f.MessageKey != "finding.cbg.unknown" {
		t.Fatalf("unknown CBG must raise R-PRIMARY-DX finding.cbg.unknown, got %+v", f)
	}
}

func TestProcSupport_ProcedureWithoutDiagnosis_Major(t *testing.T) {
	c := clean()
	c.Procedures = []claim.Procedure{{Code: "36.06"}}
	f := find(Run(c), "R-PROC-SUPPORT")
	if f == nil || f.Severity != claim.SeverityMajor {
		t.Fatalf("want major R-PROC-SUPPORT, got %+v", f)
	}
}

func TestProcSupport_SupportedProcedure_NoFinding(t *testing.T) {
	c := clean()
	c.CBGCode = "K-4-10-I"
	c.Diagnoses = []claim.Diagnosis{{Code: "K35.8", Primary: true}}
	c.Procedures = []claim.Procedure{{Code: "47.09"}}
	if f := find(Run(c), "R-PROC-SUPPORT"); f != nil {
		t.Fatalf("supported procedure must not raise, got %+v", f)
	}
}

func TestCodeFormat_InvalidICD10_Minor(t *testing.T) {
	c := clean()
	c.Diagnoses = []claim.Diagnosis{{Code: "J189", Primary: true}}
	f := find(Run(c), "R-CODE-FORMAT")
	if f == nil || f.Severity != claim.SeverityMinor || f.Params["code"] != "J189" {
		t.Fatalf("want minor R-CODE-FORMAT for J189, got %+v", f)
	}
}

func TestRun_CleanCase_NoFindings(t *testing.T) {
	if fs := Run(clean()); len(fs) != 0 {
		t.Fatalf("clean case must have no findings, got %+v", fs)
	}
}

func TestRuleData_EntriesCarrySourceAndSimplified(t *testing.T) {
	for id, g := range data.CBG.Groups {
		if g.Source == "" || !g.Simplified {
			t.Errorf("cbg %s must carry source and simplified=true", id)
		}
	}
	for id, p := range data.Proc.Procedures {
		if p.Source == "" || !p.Simplified {
			t.Errorf("procedure %s must carry source and simplified=true", id)
		}
	}
}

func TestRules_MatchExpectedOnSyntheticCases(t *testing.T) {
	cases, err := claim.LoadCases(os.DirFS("../../testdata/cases"))
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cases {
		var want, got []string
		for _, e := range c.Expected {
			if strings.HasPrefix(e, "R-") {
				want = append(want, e)
			}
		}
		for _, f := range Run(c) {
			got = append(got, f.RuleID)
		}
		sort.Strings(want)
		sort.Strings(got)
		if strings.Join(want, ",") != strings.Join(got, ",") {
			t.Errorf("%s: rule findings want %v got %v", c.ID, want, got)
		}
	}
}
