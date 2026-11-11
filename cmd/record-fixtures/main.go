// Command record-fixtures calls the real model once per synthetic case and saves the
// verified findings to testdata/resume_fixtures. Needs ANTHROPIC_API_KEY. Review the output before committing.
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
	"github.com/daffaimron21/claim-autopsy/internal/resume"
)

func main() {
	model := os.Getenv("LIVE_MODEL_NAME")
	if model == "" {
		model = "claude-haiku-5-5"
	}
	cases, err := claim.LoadCases(os.DirFS("testdata/cases"))
	if err != nil {
		log.Fatal(err)
	}
	r := resume.Verified(resume.NewAnthropic(os.Getenv("ANTHROPIC_API_KEY"), model, http.DefaultClient))
	for _, c := range cases {
		res, err := r.Read(context.Background(), c)
		if err != nil {
			log.Fatalf("%s: %v", c.ID, err)
		}
		if res.Findings == nil {
			res.Findings = []claim.Finding{}
		}
		b, _ := json.MarshalIndent(res.Findings, "", "  ")
		if err := os.WriteFile(filepath.Join("testdata/resume_fixtures", c.ID+".json"), append(b, '\n'), 0o644); err != nil {
			log.Fatal(err)
		}
		log.Printf("%s: %d findings", c.ID, len(res.Findings))
	}
}
