// Package api serves the autopsy HTTP API. Error bodies carry message KEYS, the web app translates them.
package api

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/time/rate"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
	"github.com/daffaimron21/claim-autopsy/internal/report"
	"github.com/daffaimron21/claim-autopsy/internal/resume"
	"github.com/daffaimron21/claim-autopsy/internal/rules"
)

type caseView struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	CBGCode       string            `json:"cbg_code"`
	ResumeText    string            `json:"resume_text"`
	ReturnReason  string            `json:"return_reason"`
	Diagnoses     []claim.Diagnosis `json:"diagnoses"`
	Procedures    []claim.Procedure `json:"procedures"`
	AdmitDate     string            `json:"admit_date"`
	DischargeDate string            `json:"discharge_date"`
	CareClass     int               `json:"care_class"`
}

type handler struct {
	cases   []claim.Case
	live    resume.Reader // nil when the live model is off
	saved   resume.Reader
	limiter *rate.Limiter // nil means unlimited
}

// NewHandler wires the API. live may be nil (saved results only). When live is set,
// a denied limiter or a failing model falls back to saved and marks the report degraded.
func NewHandler(cases []claim.Case, live, saved resume.Reader, limiter *rate.Limiter) http.Handler {
	h := &handler{cases, live, saved, limiter}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/cases", h.list)
	mux.HandleFunc("GET /api/cases/{id}", h.get)
	mux.HandleFunc("POST /api/cases/{id}/autopsy", h.autopsy)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, key string) {
	writeJSON(w, status, map[string]string{"error": key})
}

func (h *handler) list(w http.ResponseWriter, _ *http.Request) {
	items := make([]map[string]string, 0, len(h.cases))
	for _, c := range h.cases {
		items = append(items, map[string]string{"id": c.ID, "title": c.Title})
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	c, ok := claim.FindCase(h.cases, r.PathValue("id"))
	if !ok {
		writeErr(w, http.StatusNotFound, "error.case_not_found")
		return
	}
	writeJSON(w, http.StatusOK, caseView{c.ID, c.Title, c.CBGCode, c.ResumeText, c.ReturnReason,
		c.Diagnoses, c.Procedures, c.AdmitDate, c.DischargeDate, c.CareClass})
}

func (h *handler) autopsy(w http.ResponseWriter, r *http.Request) {
	c, ok := claim.FindCase(h.cases, r.PathValue("id"))
	if !ok {
		writeErr(w, http.StatusNotFound, "error.case_not_found")
		return
	}
	reader, limited := h.saved, false
	if h.live != nil {
		if h.limiter == nil || h.limiter.Allow() {
			reader = resume.NewFallback(h.live, h.saved)
		} else {
			limited = true
		}
	}
	res, err := reader.Read(r.Context(), c)
	if err != nil {
		log.Printf("autopsy %s: %v", c.ID, err)
		writeErr(w, http.StatusInternalServerError, "error.autopsy_failed")
		return
	}
	writeJSON(w, http.StatusOK, report.Assemble(c, rules.Run(c), res.Findings, res.Degraded || limited))
}
