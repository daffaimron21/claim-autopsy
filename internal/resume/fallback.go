package resume

import (
	"context"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

type fallback struct{ primary, saved Reader }

// NewFallback tries primary; on any error (including timeout) it serves saved
// and marks the result Degraded so the report can say so.
func NewFallback(primary, saved Reader) Reader { return fallback{primary, saved} }

func (f fallback) Read(ctx context.Context, c claim.Case) (Result, error) {
	res, err := f.primary.Read(ctx, c)
	if err == nil {
		return res, nil
	}
	res, err = f.saved.Read(context.WithoutCancel(ctx), c)
	res.Degraded = true
	return res, err
}
