package resume

import (
	"context"
	"strings"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

type verified struct{ inner Reader }

// Verified drops every finding whose evidence is missing or whose quote is not
// verbatim in the resume text (whitespace-insensitive). This is the anti-hallucination gate.
func Verified(r Reader) Reader { return verified{r} }

func norm(s string) string { return strings.Join(strings.Fields(s), " ") }

func (v verified) Read(ctx context.Context, c claim.Case) (Result, error) {
	res, err := v.inner.Read(ctx, c)
	if err != nil {
		return res, err
	}
	text := norm(c.ResumeText)
	kept := res.Findings[:0:0]
	for _, f := range res.Findings {
		ok := len(f.Evidence) > 0
		for _, e := range f.Evidence {
			q := norm(e.Quote)
			if q == "" || !strings.Contains(text, q) {
				ok = false
			}
		}
		if ok {
			kept = append(kept, f)
		}
	}
	res.Findings = kept
	return res, nil
}
