// Package claim holds the shared types for a returned BPJS claim and its autopsy report.
// All data is synthetic.
package claim

type Severity int

const (
	SeverityMinor Severity = iota
	SeverityMajor
	SeverityBlocking
)

type Source string

const (
	SourceRule   Source = "rule"
	SourceResume Source = "resume"
)

type Diagnosis struct {
	Code    string `json:"code"`
	Primary bool   `json:"primary"`
}

type Procedure struct {
	Code string `json:"code"`
}

// Case is one returned claim. Dates are YYYY-MM-DD. Expected lists the true root-cause
// RuleIDs and is the oracle for evaluation; it is never sent to the web app.
type Case struct {
	ID            string      `json:"id"`
	Title         string      `json:"title"`
	CBGCode       string      `json:"cbg_code"`
	ResumeText    string      `json:"resume_text"`
	ReturnReason  string      `json:"return_reason"`
	Diagnoses     []Diagnosis `json:"diagnoses"`
	Procedures    []Procedure `json:"procedures"`
	AdmitDate     string      `json:"admit_date"`
	DischargeDate string      `json:"discharge_date"`
	CareClass     int         `json:"care_class"`
	Expected      []string    `json:"expected"`
}

type Evidence struct {
	Field string `json:"field"`
	Quote string `json:"quote"`
}

type Finding struct {
	RuleID     string            `json:"rule_id"`
	Source     Source            `json:"source"`
	Severity   Severity          `json:"severity"`
	MessageKey string            `json:"message_key"`
	Params     map[string]string `json:"params,omitempty"`
	Evidence   []Evidence        `json:"evidence,omitempty"`
	FixKey     string            `json:"fix_key"`
	Fixable    bool              `json:"fixable"`
}

type Verdict string

const (
	VerdictResubmitAsIs     Verdict = "resubmit_as_is"
	VerdictResubmitAfterFix Verdict = "resubmit_after_fix"
	VerdictNotResubmittable Verdict = "not_resubmittable"
)

type Report struct {
	CaseID    string    `json:"case_id"`
	Primary   *Finding  `json:"primary"`
	Secondary []Finding `json:"secondary"`
	Verdict   Verdict   `json:"verdict"`
	Degraded  bool      `json:"degraded"`
	// Saved is true when the discharge-summary findings came from saved results, not a live model call.
	Saved bool `json:"saved"`
}
