// Package resume reads the free-text discharge summary and reports mismatches with the coded claim.
// Model output is never trusted: see Verified.
package resume

import (
	"context"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

// Result carries the findings and whether they came from the saved fallback instead of a live model call.
type Result struct {
	Findings []claim.Finding
	Degraded bool
}

type Reader interface {
	Read(ctx context.Context, c claim.Case) (Result, error)
}
