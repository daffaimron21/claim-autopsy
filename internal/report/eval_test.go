package report

import (
	"context"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
	"github.com/daffaimron21/claim-autopsy/internal/resume"
	"github.com/daffaimron21/claim-autopsy/internal/rules"
)

// The oracle: every synthetic case's true root causes must be found, nothing extra,
// and the right cause must rank first.
func TestEval_AllSyntheticCasesRecoverExpectedRootCauses(t *testing.T) {
	cases, err := claim.LoadCases(os.DirFS("../../testdata/cases"))
	if err != nil {
		t.Fatal(err)
	}
	reader := resume.Verified(resume.NewFixture(os.DirFS("../../testdata/resume_fixtures")))
	wantPrimary := map[string]string{
		"case-01": "R-LOS", "case-02": "R-PROC-SUPPORT", "case-03": "R-PRIMARY-DX",
		"case-04": "R-DATES", "case-05": "M-UNCODED-DX", "case-06": "M-DATE-MISMATCH",
	}
	for _, c := range cases {
		res, err := reader.Read(context.Background(), c)
		if err != nil {
			t.Fatalf("%s: %v", c.ID, err)
		}
		r := Assemble(c, rules.Run(c), res.Findings, res.Degraded)
		got := ids(r)
		sorted := append([]string(nil), got...)
		want := append([]string(nil), c.Expected...)
		sort.Strings(sorted)
		sort.Strings(want)
		t.Logf("%s | expected %v | found %v | verdict %s", c.ID, c.Expected, got, r.Verdict)
		if strings.Join(sorted, ",") != strings.Join(want, ",") {
			t.Errorf("%s: want causes %v got %v", c.ID, want, sorted)
		}
		if r.Primary == nil || r.Primary.RuleID != wantPrimary[c.ID] {
			t.Errorf("%s: want primary %s got %+v", c.ID, wantPrimary[c.ID], r.Primary)
		}
	}
}
