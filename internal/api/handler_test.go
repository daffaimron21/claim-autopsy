package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"golang.org/x/time/rate"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
	"github.com/daffaimron21/claim-autopsy/internal/resume"
)

type liveFake struct {
	calls *int32
	res   resume.Result
	err   error
}

func (l liveFake) Read(context.Context, claim.Case) (resume.Result, error) {
	atomic.AddInt32(l.calls, 1)
	return l.res, l.err
}

type brokenSaved struct{}

func (brokenSaved) Read(context.Context, claim.Case) (resume.Result, error) {
	return resume.Result{}, errors.New("disk gone")
}

func setup(t *testing.T, live resume.Reader, saved resume.Reader, lim *rate.Limiter) http.Handler {
	t.Helper()
	cases, err := claim.LoadCases(os.DirFS("../../testdata/cases"))
	if err != nil {
		t.Fatal(err)
	}
	if saved == nil {
		saved = resume.Verified(resume.NewFixture(os.DirFS("../../testdata/resume_fixtures")))
	}
	return NewHandler(cases, live, saved, lim)
}

func do(h http.Handler, method, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func decodeReport(t *testing.T, rec *httptest.ResponseRecorder) claim.Report {
	t.Helper()
	var r claim.Report
	if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil {
		t.Fatalf("bad report JSON: %v\n%s", err, rec.Body.String())
	}
	return r
}

func TestListCases_ReturnsAll(t *testing.T) {
	rec := do(setup(t, nil, nil, nil), "GET", "/api/cases")
	var items []map[string]string
	if rec.Code != 200 || json.Unmarshal(rec.Body.Bytes(), &items) != nil || len(items) != 6 || items[0]["id"] == "" || items[0]["title"] == "" {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestGetCase_HidesExpected(t *testing.T) {
	rec := do(setup(t, nil, nil, nil), "GET", "/api/cases/case-05")
	if rec.Code != 200 || strings.Contains(rec.Body.String(), `"expected"`) || !strings.Contains(rec.Body.String(), "resume_text") {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestAutopsy_ReturnsReportWithPrimary(t *testing.T) {
	rec := do(setup(t, nil, nil, nil), "POST", "/api/cases/case-01/autopsy")
	r := decodeReport(t, rec)
	if rec.Code != 200 || r.Primary == nil || r.Primary.RuleID != "R-LOS" || r.Degraded {
		t.Fatalf("code=%d report=%+v", rec.Code, r)
	}
}

func TestAutopsy_UnknownCase_404(t *testing.T) {
	rec := do(setup(t, nil, nil, nil), "POST", "/api/cases/nope/autopsy")
	if rec.Code != 404 || !strings.Contains(rec.Body.String(), "error.case_not_found") {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body.String())
	}
	if rec2 := do(setup(t, nil, nil, nil), "GET", "/api/cases/nope"); rec2.Code != 404 {
		t.Fatalf("GET unknown case want 404, got %d", rec2.Code)
	}
}

func TestAutopsy_ModelDown_StillReturnsReportDegraded(t *testing.T) {
	var calls int32
	h := setup(t, liveFake{calls: &calls, err: errors.New("provider down")}, nil, nil)
	rec := do(h, "POST", "/api/cases/case-06/autopsy")
	r := decodeReport(t, rec)
	if rec.Code != 200 || !r.Degraded || r.Primary == nil || r.Primary.RuleID != "M-DATE-MISMATCH" || calls != 1 {
		t.Fatalf("code=%d calls=%d report=%+v", rec.Code, calls, r)
	}
}

func TestAutopsy_RateLimited_FallsBackToFixture(t *testing.T) {
	var calls int32
	h := setup(t, liveFake{calls: &calls}, nil, rate.NewLimiter(0, 0)) // never allows
	rec := do(h, "POST", "/api/cases/case-06/autopsy")
	r := decodeReport(t, rec)
	if rec.Code != 200 || calls != 0 || !r.Degraded || r.Primary == nil || r.Primary.RuleID != "M-DATE-MISMATCH" {
		t.Fatalf("code=%d calls=%d report=%+v", rec.Code, calls, r)
	}
}

func TestAutopsy_LiveUsedWhenAllowed(t *testing.T) {
	var calls int32
	h := setup(t, liveFake{calls: &calls}, nil, rate.NewLimiter(rate.Inf, 1))
	rec := do(h, "POST", "/api/cases/case-01/autopsy")
	if rec.Code != 200 || calls != 1 || decodeReport(t, rec).Degraded {
		t.Fatalf("code=%d calls=%d", rec.Code, calls)
	}
}

func TestAutopsy_AllReadersFail_500(t *testing.T) {
	var calls int32
	h := setup(t, liveFake{calls: &calls, err: errors.New("x")}, brokenSaved{}, nil)
	rec := do(h, "POST", "/api/cases/case-01/autopsy")
	if rec.Code != 500 || !strings.Contains(rec.Body.String(), "error.autopsy_failed") {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body.String())
	}
}
