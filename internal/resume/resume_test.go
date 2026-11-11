package resume

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"testing/fstest"
	"time"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

type fake struct {
	res   Result
	err   error
	block bool
}

func (f fake) Read(ctx context.Context, _ claim.Case) (Result, error) {
	if f.block {
		<-ctx.Done()
		return Result{}, ctx.Err()
	}
	return f.res, f.err
}

func finding(quote string) claim.Finding {
	return claim.Finding{RuleID: "M-UNCODED-DX", Source: claim.SourceResume,
		Evidence: []claim.Evidence{{Field: "resume_text", Quote: quote}}}
}

var theCase = claim.Case{ID: "c1", ResumeText: "Pasien  dirawat 3 hari,\npulang 4 Juli 2026 dalam kondisi baik."}

func TestVerified_DropsFindingWhoseQuoteNotInResume(t *testing.T) {
	r := Verified(fake{res: Result{Findings: []claim.Finding{finding("pulang 4 Juli 2026"), finding("kalimat karangan model")}}})
	res, err := r.Read(context.Background(), theCase)
	if err != nil || len(res.Findings) != 1 || res.Findings[0].Evidence[0].Quote != "pulang 4 Juli 2026" {
		t.Fatalf("want only the real quote kept, got %+v err=%v", res.Findings, err)
	}
}

func TestVerified_AcceptsQuoteWithDifferentWhitespace(t *testing.T) {
	r := Verified(fake{res: Result{Findings: []claim.Finding{finding("Pasien dirawat 3 hari, pulang 4 Juli 2026")}}})
	res, _ := r.Read(context.Background(), theCase)
	if len(res.Findings) != 1 {
		t.Fatalf("whitespace differences must not drop a real quote, got %+v", res.Findings)
	}
}

func TestVerified_DropsFindingWithEmptyQuoteOrNoEvidence(t *testing.T) {
	noEv := claim.Finding{RuleID: "M-UNCODED-DX", Source: claim.SourceResume}
	r := Verified(fake{res: Result{Findings: []claim.Finding{finding(""), finding("   "), noEv}}})
	res, _ := r.Read(context.Background(), theCase)
	if len(res.Findings) != 0 {
		t.Fatalf("want none kept, got %+v", res.Findings)
	}
}

func fixtureFS() fstest.MapFS {
	return fstest.MapFS{"c1.json": {Data: []byte(`[{"rule_id":"M-DATE-MISMATCH","source":"resume","severity":2,"message_key":"finding.resume.date_mismatch","evidence":[{"field":"resume_text","quote":"pulang 4 Juli 2026"}],"fix_key":"fix.correct_dates","fixable":true}]`)}}
}

func TestFixture_ReturnsSavedFindingsForCaseID(t *testing.T) {
	res, err := NewFixture(fixtureFS()).Read(context.Background(), theCase)
	if err != nil || len(res.Findings) != 1 || res.Findings[0].RuleID != "M-DATE-MISMATCH" || res.Degraded {
		t.Fatalf("got %+v err=%v", res, err)
	}
}

func TestFixture_MissingCase_EmptyNoError(t *testing.T) {
	res, err := NewFixture(fixtureFS()).Read(context.Background(), claim.Case{ID: "other"})
	if err != nil || len(res.Findings) != 0 {
		t.Fatalf("got %+v err=%v", res, err)
	}
}

func TestFallback_UsesFixtureWhenPrimaryErrors(t *testing.T) {
	r := NewFallback(fake{err: errors.New("boom")}, NewFixture(fixtureFS()))
	res, err := r.Read(context.Background(), theCase)
	if err != nil || !res.Degraded || len(res.Findings) != 1 {
		t.Fatalf("want fixture findings and Degraded, got %+v err=%v", res, err)
	}
}

func TestFallback_UsesFixtureOnTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	res, err := NewFallback(fake{block: true}, NewFixture(fixtureFS())).Read(ctx, theCase)
	if err != nil || !res.Degraded || len(res.Findings) != 1 {
		t.Fatalf("want fixture on timeout, got %+v err=%v", res, err)
	}
}

func TestFallback_NotDegradedWhenPrimarySucceeds(t *testing.T) {
	r := NewFallback(fake{res: Result{Findings: []claim.Finding{finding("x")}}}, NewFixture(fixtureFS()))
	res, err := r.Read(context.Background(), theCase)
	if err != nil || res.Degraded || len(res.Findings) != 1 {
		t.Fatalf("got %+v err=%v", res, err)
	}
}

func modelServer(t *testing.T, text string, status int, hits *int32) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(hits, 1)
		if r.Header.Get("x-api-key") != "k" || r.URL.Path != "/v1/messages" {
			http.Error(w, "bad request", 400)
			return
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{"content":[{"type":"text","text":` + quoteJSON(text) + `}]}`))
	}))
}

func quoteJSON(s string) string {
	b := []byte{'"'}
	for _, r := range s {
		switch r {
		case '"', '\\':
			b = append(b, '\\', byte(r))
		case '\n':
			b = append(b, '\\', 'n')
		default:
			b = append(b, string(r)...)
		}
	}
	return string(append(b, '"'))
}

func TestAnthropic_EmptyKey_ReturnsErrorWithoutNetwork(t *testing.T) {
	var hits int32
	srv := modelServer(t, `{"findings":[]}`, 200, &hits)
	defer srv.Close()
	if _, err := newAnthropicAt(srv.URL, "", "m", srv.Client()).Read(context.Background(), theCase); err == nil {
		t.Fatal("empty key must be an error")
	}
	if hits != 0 {
		t.Fatalf("no network call expected, got %d", hits)
	}
}

func TestAnthropic_ParsesModelJSON(t *testing.T) {
	var hits int32
	body := "```json\n" + `{"findings":[{"rule_id":"M-DATE-MISMATCH","quote":"pulang 4 Juli 2026","params":{"resume_date":"2026-07-04"}}]}` + "\n```"
	srv := modelServer(t, body, 200, &hits)
	defer srv.Close()
	res, err := newAnthropicAt(srv.URL, "k", "m", srv.Client()).Read(context.Background(), theCase)
	if err != nil || len(res.Findings) != 1 {
		t.Fatalf("got %+v err=%v", res, err)
	}
	f := res.Findings[0]
	if f.RuleID != "M-DATE-MISMATCH" || f.Source != claim.SourceResume || f.Severity != claim.SeverityBlocking ||
		f.Evidence[0].Quote != "pulang 4 Juli 2026" || f.MessageKey != "finding.resume.date_mismatch" || f.FixKey != "fix.correct_dates" {
		t.Fatalf("unexpected finding %+v", f)
	}
}

func TestAnthropic_DropsUnknownRuleID(t *testing.T) {
	var hits int32
	srv := modelServer(t, `{"findings":[{"rule_id":"M-MADE-UP","quote":"x"}]}`, 200, &hits)
	defer srv.Close()
	res, err := newAnthropicAt(srv.URL, "k", "m", srv.Client()).Read(context.Background(), theCase)
	if err != nil || len(res.Findings) != 0 {
		t.Fatalf("unknown rule id must be dropped, got %+v err=%v", res, err)
	}
}

func TestAnthropic_HTTPError_ReturnsError(t *testing.T) {
	var hits int32
	srv := modelServer(t, `{}`, 500, &hits)
	defer srv.Close()
	if _, err := newAnthropicAt(srv.URL, "k", "m", srv.Client()).Read(context.Background(), theCase); err == nil {
		t.Fatal("HTTP 500 must be an error")
	}
}

func TestAnthropic_GarbageOutput_ReturnsError(t *testing.T) {
	var hits int32
	srv := modelServer(t, `sorry I cannot help`, 200, &hits)
	defer srv.Close()
	if _, err := newAnthropicAt(srv.URL, "k", "m", srv.Client()).Read(context.Background(), theCase); err == nil {
		t.Fatal("non-JSON model output must be an error so the fallback engages")
	}
}
