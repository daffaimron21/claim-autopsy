package resume

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

const (
	apiVersion  = "2023-06-01"
	callTimeout = 15 * time.Second
)

type ruleMeta struct {
	severity   claim.Severity
	messageKey string
	fixKey     string
}

// The model may only pick from this closed list. Severity and keys come from here, not from the model.
var closed = map[string]ruleMeta{
	"M-UNCODED-DX":    {claim.SeverityMajor, "finding.resume.uncoded_dx", "fix.add_missing_dx"},
	"M-DATE-MISMATCH": {claim.SeverityBlocking, "finding.resume.date_mismatch", "fix.correct_dates"},
	"M-PROC-MISSING":  {claim.SeverityMajor, "finding.resume.proc_missing", "fix.code_missing_procedure"},
}

const systemPrompt = "You compare an Indonesian hospital discharge summary with the coded data of a BPJS claim. " +
	"Report ONLY mismatches of these kinds: M-UNCODED-DX (a diagnosis in the summary that is not in the coded diagnoses), " +
	"M-DATE-MISMATCH (a stay length or discharge date in the summary that disagrees with the claim dates), " +
	"M-PROC-MISSING (a procedure in the summary that is not coded). " +
	"Every finding MUST include \"quote\": an exact fragment copied from the summary. Never invent text. " +
	"Reply with JSON only: {\"findings\":[{\"rule_id\":\"...\",\"quote\":\"...\",\"params\":{}}]}. If nothing mismatches reply {\"findings\":[]}."

type anthropicReader struct {
	baseURL, apiKey, model string
	client                 *http.Client
}

func NewAnthropic(apiKey, model string, client *http.Client) Reader {
	return newAnthropicAt("https://api.anthropic.com", apiKey, model, client)
}

func newAnthropicAt(baseURL, apiKey, model string, client *http.Client) Reader {
	if client == nil {
		client = http.DefaultClient
	}
	return anthropicReader{baseURL, apiKey, model, client}
}

func (a anthropicReader) Read(ctx context.Context, c claim.Case) (Result, error) {
	if a.apiKey == "" {
		return Result{}, errors.New("resume: missing API key")
	}
	ctx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	structured, _ := json.Marshal(map[string]any{
		"cbg_code": c.CBGCode, "diagnoses": c.Diagnoses, "procedures": c.Procedures,
		"admit_date": c.AdmitDate, "discharge_date": c.DischargeDate,
	})
	body, _ := json.Marshal(map[string]any{
		"model": a.model, "max_tokens": 1024, "system": systemPrompt,
		"messages": []map[string]string{{"role": "user", "content": "Coded claim data:\n" + string(structured) + "\n\nDischarge summary:\n" + c.ResumeText}},
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", apiVersion)
	req.Header.Set("content-type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Result{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("resume: model API status %d", resp.StatusCode)
	}
	var env struct {
		Content []struct{ Type, Text string } `json:"content"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return Result{}, err
	}
	var text strings.Builder
	for _, p := range env.Content {
		if p.Type == "text" {
			text.WriteString(p.Text)
		}
	}
	return parseFindings(text.String())
}

func parseFindings(text string) (Result, error) {
	i, j := strings.Index(text, "{"), strings.LastIndex(text, "}")
	if i < 0 || j < i {
		return Result{}, errors.New("resume: model returned no JSON")
	}
	var out struct {
		Findings []struct {
			RuleID string         `json:"rule_id"`
			Quote  string         `json:"quote"`
			Params map[string]any `json:"params"`
		} `json:"findings"`
	}
	if err := json.Unmarshal([]byte(text[i:j+1]), &out); err != nil {
		return Result{}, fmt.Errorf("resume: bad model JSON: %w", err)
	}
	var res Result
	for _, f := range out.Findings {
		meta, ok := closed[f.RuleID]
		if !ok {
			continue
		}
		params := map[string]string{}
		for k, v := range f.Params {
			params[k] = fmt.Sprint(v)
		}
		res.Findings = append(res.Findings, claim.Finding{
			RuleID: f.RuleID, Source: claim.SourceResume, Severity: meta.severity,
			MessageKey: meta.messageKey, FixKey: meta.fixKey, Fixable: true, Params: params,
			Evidence: []claim.Evidence{{Field: "resume_text", Quote: f.Quote}},
		})
	}
	return res, nil
}
